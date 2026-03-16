import {
  type AxiosError,
  type AxiosInstance,
  type AxiosInterceptorManager,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from "axios";
import { toast } from "sonner";

import { AUTH_CONFIG } from "../constants";
import { createErrorResponse } from "../utils/api-error-response";
import { handleRetryLogic, shouldRetry } from "../utils/retry-handler";
import { handleTokenRefresh } from "../utils/refresh-token-handler";

interface RetryableRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
  _retryCount?: number;
}

// Flag to prevent multiple 401 redirects causing infinite loop
let isRedirecting = false;

export const responseInterceptor = (
  response: AxiosInterceptorManager<AxiosResponse>,
  axiosInstance: AxiosInstance
) => {
  response.use(
    (data: AxiosResponse) => {
      return data;
    },
    async (error: AxiosError<BaseResponseAPI<any>>) => {
      const originalRequest = error.config as RetryableRequestConfig;

      if (!originalRequest) {
        return Promise.reject(createErrorResponse(error));
      }

      const isLoginRequest = originalRequest.url?.includes(
        AUTH_CONFIG.LOGIN_ENDPOINT
      );
      const isLogoutRequest = originalRequest.url?.includes(
        AUTH_CONFIG.LOGOUT_ENDPOINT
      );
      const isRefreshRequest = originalRequest.url?.includes(
        AUTH_CONFIG.REFRESH_TOKEN_ENDPOINT
      );

      // Handle 401 - Unauthorized
      if (error.response?.status === 401) {
        // Don't handle 401 for login/refresh requests
        if (isLoginRequest || isRefreshRequest) {
          return Promise.reject(createErrorResponse(error));
        }

        // Try to refresh token silently first
        if (!originalRequest._retry && !isRedirecting) {
          originalRequest._retry = true;

          try {
            // Thử refresh token trong âm thầm
            console.log("🔄 [Auth] Attempting silent token refresh...");
            const result = await handleTokenRefresh(originalRequest, axiosInstance);
            console.log("✅ [Auth] Silent token refresh successful");
            return result;
          } catch (refreshError) {
            // Refresh token thất bại - redirect về login
            console.log("❌ [Auth] Silent token refresh failed, redirecting to login");

            if (typeof window !== "undefined") {
              const currentPath = window.location.pathname;
              const publicPaths = ["/login", "/forgot-password", "/reset-password", "/setup-2fa", "/verify-otp", "/login/otp"];
              const isPublicPage = publicPaths.some(p => currentPath.startsWith(p));

              if (!isPublicPage) {
                // Set flag to prevent multiple redirects
                isRedirecting = true;
                toast.error("Phiên đăng nhập hết hạn. Vui lòng đăng nhập lại!");
                // Use replace instead of href to prevent history pollution
                window.location.replace("/login");
                return Promise.reject(createErrorResponse(error));
              }
            }
          }
        }

        return Promise.reject(createErrorResponse(error));
      }

      // Handle retry logic for 5xx errors
      if (shouldRetry(error)) {
        if (isLogoutRequest) {
          return Promise.reject(createErrorResponse(error));
        }
        return handleRetryLogic(error, originalRequest, axiosInstance);
      }

      // Handle other status codes - log for debugging
      // Note: Refine notification provider will handle showing toasts
      const status = error.response?.status;
      const errorData = error.response?.data;

      if (status === 403) {
        console.warn(
          "Access forbidden - insufficient permissions",
          errorData
        );
      } else if (status === 404) {
        console.warn("Resource not found", errorData);
      } else if (status === 500) {
        console.error("Server error", errorData);
      }

      return Promise.reject(createErrorResponse(error));
    }
  );
};
