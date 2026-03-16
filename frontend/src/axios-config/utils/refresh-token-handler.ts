import axios, {
  type AxiosInstance,
  type InternalAxiosRequestConfig,
} from "axios";
import { tokenManager } from "./token-manager";
import { AUTH_CONFIG } from "../constants";
import { post } from "../request";

type RefreshTokenResponse = ResponseAPI<{
  token: string;
  refresh_token: string;
}>;

interface RetryableRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
  _retryCount?: number;
}

// ============================================
// Mutex lock pattern for refresh token (cross-tab safe)
// ============================================
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value: unknown) => void;
  reject: (error: unknown) => void;
}> = [];

// Key for cross-tab synchronization
const REFRESH_LOCK_KEY = "hrms_refresh_lock";
const REFRESH_LOCK_TIMEOUT = 15000; // 15 seconds max lock time

const processQueue = (error: Error | null, token: string | null = null) => {
  console.log(`📤 [Queue] Processing ${failedQueue.length} queued requests, error: ${!!error}`);
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

// Check if another tab is refreshing
const isAnotherTabRefreshing = (): boolean => {
  if (typeof window === "undefined") return false;
  const lockData = localStorage.getItem(REFRESH_LOCK_KEY);
  if (!lockData) return false;

  try {
    const { timestamp } = JSON.parse(lockData);
    // Lock expired?
    if (Date.now() - timestamp > REFRESH_LOCK_TIMEOUT) {
      localStorage.removeItem(REFRESH_LOCK_KEY);
      return false;
    }
    return true;
  } catch {
    return false;
  }
};

// Set refresh lock for cross-tab
const setRefreshLock = () => {
  if (typeof window === "undefined") return;
  localStorage.setItem(REFRESH_LOCK_KEY, JSON.stringify({ timestamp: Date.now() }));
};

// Clear refresh lock
const clearRefreshLock = () => {
  if (typeof window === "undefined") return;
  localStorage.removeItem(REFRESH_LOCK_KEY);
};

// Base URL for refresh - bypass interceptors
const getBaseUrl = () => process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000/api/v1";

async function doRefresh(): Promise<string> {
  const refreshToken = tokenManager.getRefreshToken();
  if (!refreshToken) {
    throw new Error("No refresh token available");
  }

  console.log("🔄 [Refresh Token] Starting refresh process...");

  try {
    // IMPORTANT: Use raw axios instead of axiosInstance to avoid interceptor loop
    const response = await post<RefreshTokenResponse, { refresh_token: string }>(
      `${getBaseUrl()}${AUTH_CONFIG.REFRESH_TOKEN_ENDPOINT}`,
      { refresh_token: refreshToken },
      {
        headers: {
          "Content-Type": "application/json",
        },
        timeout: 10000, // 10s timeout for refresh
      }
    );

    const responseData = response?.data
    const token = responseData?.token;
    const refresh_token = responseData?.refresh_token;

    if (!token || !refresh_token) {
      console.error("❌ [Refresh Token] Invalid response structure:", response.data);
      throw new Error("Invalid refresh response");
    }

    console.log("✅ [Refresh Token] Successfully refreshed");
    tokenManager.setTokens(token, refresh_token);

    return token;
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message;
    console.error("❌ [Refresh Token] Failed:", errorMsg);

    // If invalid credentials, clear tokens immediately
    if (errorMsg?.includes("Invalid") || error?.response?.status === 401) {
      console.log("🗑️ [Refresh Token] Clearing invalid tokens");
      tokenManager.clearTokens();
    }

    throw error;
  }
}

async function logout() {
  console.log("🚪 [Auth] Logging out...");
  tokenManager.clearTokens();
  clearRefreshLock();

  if (typeof window !== "undefined") {
    window.location.href = "/login";
  }
}

export async function handleTokenRefresh(
  originalRequest: RetryableRequestConfig,
  axiosInstance: AxiosInstance
): Promise<unknown> {
  // Already retried this request - don't try again
  if (originalRequest._retry) {
    console.log(`⏭️ [Auth] Request already retried, skipping: ${originalRequest.url}`);
    return Promise.reject(new Error("Token refresh failed"));
  }

  // Check if we have tokens at all
  const currentToken = tokenManager.getAccessToken();
  const currentRefreshToken = tokenManager.getRefreshToken();

  if (!currentRefreshToken) {
    console.log(`🚫 [Auth] No refresh token, redirecting to login`);
    await logout();
    return Promise.reject(new Error("No refresh token"));
  }

  console.log(`🔐 [Auth] Token refresh needed for ${originalRequest.url}`);
  originalRequest._retry = true;

  // If already refreshing in this tab, add to queue
  if (isRefreshing) {
    console.log(`⏳ [Auth] Refresh in progress, queuing request: ${originalRequest.url}`);
    return new Promise((resolve, reject) => {
      failedQueue.push({ resolve, reject });
    }).then((token) => {
      // Got new token from queue, retry request
      originalRequest.headers.set("Authorization", `Bearer ${token}`);
      console.log(`🔄 [Auth] Retrying queued request: ${originalRequest.url}`);
      return axiosInstance(originalRequest);
    }).catch((err) => {
      return Promise.reject(err);
    });
  }

  // Check if another tab is refreshing - wait a bit and use new token
  if (isAnotherTabRefreshing()) {
    console.log(`🔄 [Auth] Another tab is refreshing, waiting...`);
    // Wait for other tab to finish (max 5 seconds)
    await new Promise(resolve => setTimeout(resolve, 2000));

    // Check if we got new token from other tab
    const newToken = tokenManager.getAccessToken();
    if (newToken && newToken !== currentToken) {
      console.log(`✅ [Auth] Got new token from another tab`);
      originalRequest.headers.set("Authorization", `Bearer ${newToken}`);
      return axiosInstance(originalRequest);
    }
  }

  // Start refresh
  isRefreshing = true;
  setRefreshLock();

  try {
    const newToken = await doRefresh();

    // Process all queued requests with new token
    processQueue(null, newToken);

    // Retry original request with new token
    originalRequest.headers.set("Authorization", `Bearer ${newToken}`);
    console.log(`🔄 [Auth] Retrying original request: ${originalRequest.url}`);
    return axiosInstance(originalRequest);
  } catch (refreshError: any) {
    console.error(`❌ [Auth] Refresh failed, logging out`);

    // Reject all queued requests
    processQueue(refreshError, null);

    // Logout
    await logout();

    return Promise.reject(refreshError);
  } finally {
    isRefreshing = false;
    clearRefreshLock();
  }
}
