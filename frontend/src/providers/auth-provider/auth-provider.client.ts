"use client";

import type { AuthProvider } from "@refinedev/core";
import { post, get, put } from "@/axios-config";
import { tokenManager } from "@/axios-config";
import { AUTH_CONFIG } from "@/axios-config/constants";
import type {
  LoginRequest,
  LoginResponse,
  VerifyOTPRequest,
  VerifyOTPResponse,
  Setup2FARequest,
  Setup2FAResponse,
  Verify2FASetupRequest,
  Verify2FASetupResponse,
  ChangePasswordRequest,
  VerifyPasswordResetOTPRequest,
  VerifyPasswordResetOTPResponse,
} from "@/types";
import type { UserResponse, UserRole } from "@/types";
import { getRedirectPath } from "@/utils/role-guard";

// ==================== HELPER FUNCTIONS ====================

/**
 * Get pending username from sessionStorage
 */
const getPendingUsername = (): string | null => {
  if (typeof window === "undefined") return null;
  return sessionStorage.getItem("pending_2fa_username");
};

/**
 * Get pending 2FA secret from sessionStorage
 */
const getPending2FASecret = (): string | null => {
  if (typeof window === "undefined") return null;
  return sessionStorage.getItem("pending_2fa_secret");
};

/**
 * Clear pending auth data from sessionStorage
 */
const clearPendingAuthData = () => {
  if (typeof window === "undefined") return;
  sessionStorage.removeItem("pending_2fa_username");
  sessionStorage.removeItem("pending_2fa_secret");
};

/**
 * Extract error message from API error response
 */
const getErrorMessage = (error: unknown, defaultMessage: string): string => {
  const err = error as { response?: { data?: ResponseAPI<unknown> } };
  return err?.response?.data?.message || defaultMessage;
};

/**
 * Initiate 2FA setup for first-time users
 * Calls setup-2fa endpoint and stores secret
 */
const initiate2FASetup = async (username: string): Promise<{
  success: boolean;
  redirectTo?: string;
  error?: { name: string; message: string };
}> => {
  try {
    const response = await post<ResponseAPI<Setup2FAResponse>, Setup2FARequest>(
      AUTH_CONFIG.SETUP_2FA_ENDPOINT,
      { username }
    );

    const data = response.data;

    // Store secret for verification later
    if (data?.secret && typeof window !== "undefined") {
      sessionStorage.setItem("pending_2fa_secret", data.secret);
    }

    return {
      success: true,
      redirectTo: `/setup-2fa?qr=${encodeURIComponent(data?.qr_code_url || "")}`,
    };
  } catch (error: unknown) {
    return {
      success: false,
      error: {
        name: "Lỗi thiết lập 2FA",
        message: getErrorMessage(error, "Không thể khởi tạo xác thực 2 bước. Vui lòng thử lại."),
      },
    };
  }
};

// ==================== AUTH PROVIDER ====================

export const authProviderClient = {
  /**
   * Login with username and password
   * API: POST /auth/login
   * 
   * Response cases:
   * - requires_2fa: true → User needs to setup 2FA (first time)
   * - requires_otp: true → User has 2FA enabled, needs to verify OTP
   * - Neither → Login successful, redirect to dashboard
   */
  login: async ({ username, password }: LoginRequest) => {
    try {
      const response = await post<ResponseAPI<LoginResponse>, LoginRequest>(
        AUTH_CONFIG.LOGIN_ENDPOINT,
        { username, password }
      );

      const data = response.data;

      // Store username for 2FA/OTP flow
      if (typeof window !== "undefined") {
        sessionStorage.setItem("pending_2fa_username", username);
      }

      // TESTING - Force requires_2fa = true to test 2FA setup flow
      // const forceTest2FA = true;
      // if (forceTest2FA) {
      //   return initiate2FASetup(username);
      // }

      // Case 1: First time 2FA setup required
      if (data.requires_2fa) {
        return initiate2FASetup(username);
      }

      // Case 2: 2FA enabled, verify OTP
      if (data.requires_otp) {
        return {
          success: true,
          redirectTo: "/login/otp",
        };
      }

      // Case 3: Login successful - redirect based on role
      clearPendingAuthData();
      // Get role from login response if available
      const role = data.user?.role as UserRole | undefined;
      const redirectPath = getRedirectPath(role);
      return {
        success: true,
        redirectTo: redirectPath,
      };
    } catch (error: unknown) {
      return {
        success: false,
        error: {
          name: "Đăng nhập thất bại",
          message: getErrorMessage(error, "Tên đăng nhập hoặc mật khẩu không đúng"),
        },
      };
    }
  },

  /**
   * Verify OTP for 2FA login
   * API: POST /auth/verify-otp
   */
  verifyOtp: async ({ otp_code }: { otp_code: string }) => {
    const username = getPendingUsername();

    if (!username) {
      return {
        success: false,
        error: {
          name: "Lỗi xác thực",
          message: "Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.",
        },
      };
    }

    try {
      const response = await post<ResponseAPI<VerifyOTPResponse>, VerifyOTPRequest>(
        AUTH_CONFIG.VERIFY_OTP_ENDPOINT,
        { otp_code, username }
      );

      clearPendingAuthData();
      
      // Get role from response or default to profile
      const role = response.data?.user?.role as UserRole | undefined;
      const redirectPath = role === "ADMIN" || role === "SYSTEM_ADMIN" ? "/dashboard" : "/profile";
      
      return {
        success: true,
        redirectTo: redirectPath,
      };
    } catch (error: unknown) {
      return {
        success: false,
        error: {
          name: "Lỗi xác thực",
          message: getErrorMessage(error, "Mã OTP không hợp lệ"),
        },
      };
    }
  },

  /**
   * Verify 2FA setup with OTP
   * API: POST /auth/verify-2fa-setup
   */
  verify2FASetup: async ({ otp_code, secret, username }: Verify2FASetupRequest) => {
    const storedUsername = username || getPendingUsername() || "";
    const storedSecret = secret || getPending2FASecret() || "";

    try {
      const response = await post<ResponseAPI<Verify2FASetupResponse & { token?: string; refresh_token?: string }>, Verify2FASetupRequest>(
        AUTH_CONFIG.VERIFY_2FA_SETUP_ENDPOINT,
        { otp_code, secret: storedSecret, username: storedUsername }
      );


      // Save tokens if returned by backend
      if (response.data?.token) {
        tokenManager.setAccessToken(response.data.token);
      }
      if (response.data?.refresh_token) {
        tokenManager.setRefreshToken(response.data.refresh_token);
      }

      clearPendingAuthData();
      
      // After verify 2FA setup, user needs to login again
      // The Setup2FAPage handles this by redirecting to /login
      return {
        success: true,
        data: response.data,
        redirectTo: "/login",
      };
    } catch (error: unknown) {
      return {
        success: false,
        error: {
          name: "Lỗi xác thực 2FA",
          message: getErrorMessage(error, "Mã OTP không hợp lệ"),
        },
      };
    }
  },

  /**
   * Change password (authenticated users)
   * API: PUT /auth/password
   */
  changePassword: async ({ current_password, new_password, confirm_password }: ChangePasswordRequest) => {
    try {
      await put<ResponseAPI<unknown>, ChangePasswordRequest>(
        AUTH_CONFIG.CHANGE_PASSWORD_ENDPOINT,
        { current_password, new_password, confirm_password }
      );
      return {
        success: true,
        successNotification: {
          message: "Đổi mật khẩu thành công",
          description: "Mật khẩu của bạn đã được cập nhật",
        },
      };
    } catch (error: unknown) {
      return {
        success: false,
        error: {
          name: "Lỗi đổi mật khẩu",
          message: getErrorMessage(error, "Không thể đổi mật khẩu"),
        },
      };
    }
  },

  /**
   * Verify password reset OTP
   * API: POST /auth/password/reset/verify
   */
  verifyPasswordResetOTP: async ({ email, otp_code }: VerifyPasswordResetOTPRequest) => {
    try {
      const response = await post<ResponseAPI<VerifyPasswordResetOTPResponse>, VerifyPasswordResetOTPRequest>(
        AUTH_CONFIG.VERIFY_PASSWORD_RESET_OTP_ENDPOINT,
        { email, otp_code }
      );
      return {
        success: true,
        data: response.data,
      };
    } catch (error: unknown) {
      return {
        success: false,
        error: {
          name: "Lỗi xác thực OTP",
          message: getErrorMessage(error, "Mã OTP không hợp lệ"),
        },
      };
    }
  },

  /**
   * Logout
   * API: POST /auth/logout
   */
  logout: async () => {
    try {
      await post<ResponseAPI<unknown>, object>(AUTH_CONFIG.LOGOUT_ENDPOINT, {});
    } catch {
      // Ignore logout errors
    } finally {
      tokenManager.clearTokens();
      clearPendingAuthData();
    }

    return {
      success: true,
      redirectTo: "/login",
    };
  },

  /**
   * Check if user is authenticated
   * API: GET /users/me
   * Backend sử dụng session cookies - không cần check token locally
   */
  check: async () => {
    // Skip auth check on public pages
    if (typeof window !== "undefined") {
      const path = window.location.pathname;
      const publicPaths = ["/login", "/forgot-password", "/reset-password", "/setup-2fa", "/verify-otp", "/login/otp"];
      const isPublicPage = publicPaths.some(p => path.startsWith(p));

      if (isPublicPage) {
        return { authenticated: false };
      }
    }

    // Gọi API kiểm tra session - backend sẽ check session cookies tự động
    try {
      const response = await get<ResponseAPI<UserResponse>>(AUTH_CONFIG.ME_ENDPOINT);

      if (response.is_success && response.data) {
        return { authenticated: true };
      }

      return {
        authenticated: false,
        redirectTo: "/login",
      };
    } catch {
      return {
        authenticated: false,
        redirectTo: "/login",
      };
    }
  },

  /**
   * Get user permissions/roles
   * API: GET /users/me
   */
  getPermissions: async () => {
    try {
      const response = await get<ResponseAPI<UserResponse>>(AUTH_CONFIG.ME_ENDPOINT);
      return response.data?.role || null;
    } catch {
      return null;
    }
  },

  /**
   * Get current user identity
   * API: GET /users/me
   */
  getIdentity: async () => {
    try {
      const response = await get<ResponseAPI<UserResponse>>(AUTH_CONFIG.ME_ENDPOINT);
      return response.data;
    } catch {
      return null;
    }
  },

  /**
   * Handle authentication errors
   */
  onError: async (error: unknown) => {
    const status = (error as { status?: number })?.status;

    if (status === 401) {
      return {
        logout: true,
        redirectTo: "/login",
      };
    }

    if (status === 403) {
      return {
        error: {
          name: "Không có quyền truy cập",
          message: "Bạn không có quyền truy cập tài nguyên này",
        },
      };
    }

    return { error };
  },

  /**
   * Request password reset
   * API: POST /auth/password/reset/request
   */
  forgotPassword: async ({ email }: { email: string }) => {
    try {
      await post<ResponseAPI<unknown>, { email: string }>(
        AUTH_CONFIG.REQUEST_PASSWORD_RESET_ENDPOINT,
        { email }
      );
      return {
        success: true,
      };
    } catch (error: unknown) {
      return {
        success: false,
        error: {
          name: "Lỗi đặt lại mật khẩu",
          message: getErrorMessage(error, "Không thể gửi email đặt lại mật khẩu"),
        },
      };
    }
  },

  /**
   * Complete password reset
   * API: POST /auth/password/reset/complete
   */
  updatePassword: async ({ password, confirmPassword, token }: {
    password: string;
    confirmPassword: string;
    token?: string;
  }) => {
    try {
      await post<ResponseAPI<unknown>, { new_password: string; confirm_password: string; reset_token: string }>(
        AUTH_CONFIG.COMPLETE_PASSWORD_RESET_ENDPOINT,
        {
          new_password: password,
          confirm_password: confirmPassword,
          reset_token: token || ""
        }
      );
      return {
        success: true,
        redirectTo: "/login",
      };
    } catch (error: unknown) {
      return {
        success: false,
        error: {
          name: "Lỗi cập nhật mật khẩu",
          message: getErrorMessage(error, "Không thể cập nhật mật khẩu"),
        },
      };
    }
  },
} as AuthProvider & {
  verifyOtp: (params: { otp_code: string }) => Promise<{
    success: boolean;
    redirectTo?: string;
    error?: { name: string; message: string };
  }>;

  verify2FASetup: (params: Verify2FASetupRequest) => Promise<{
    success: boolean;
    data?: Verify2FASetupResponse;
    redirectTo?: string;
    error?: { name: string; message: string };
  }>;

  changePassword: (params: ChangePasswordRequest) => Promise<{
    success: boolean;
    successNotification?: { message: string; description: string };
    error?: { name: string; message: string };
  }>;

  verifyPasswordResetOTP: (params: VerifyPasswordResetOTPRequest) => Promise<{
    success: boolean;
    data?: VerifyPasswordResetOTPResponse;
    error?: { name: string; message: string };
  }>;
};
