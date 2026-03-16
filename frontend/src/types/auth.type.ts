/**
 * Auth Types - Authentication related types
 */

// User roles
export type UserRole = "ADMIN" | "SYSTEM_ADMIN" | "USER" | "MANAGER";

// Permission type - matches existing usage
export interface Permission {
    id: number;
    name: string;
    code: string;
    description?: string;
    action?: string;
    collection?: string;
}

// User identity for auth context - extended to match existing usage
export interface UserIdentity {
    id: number;
    username: string;
    email: string;
    full_name?: string;
    role: UserRole | { name?: string };
    permissions?: Permission[];
    avatar_url?: string;
    is_admin?: boolean;
    can_access_app?: boolean;
    employee_id?: number;
    employee?: {
        id: number;
        name: string;
    };
}

// Login request/response
export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse {
    requires_2fa?: boolean;
    requires_otp?: boolean;
    user?: {
        id: number;
        username: string;
        role: UserRole;
    };
}

// OTP verification
export interface VerifyOTPRequest {
    otp_code: string;
    username: string;
}

export interface VerifyOTPResponse {
    user?: UserIdentity;
    token?: string;
}

// 2FA Setup
export interface Setup2FARequest {
    username: string;
}

export interface Setup2FAResponse {
    secret: string;
    qr_code_url: string;
}

export interface Verify2FASetupRequest {
    otp_code: string;
    secret?: string;
    username?: string;
}

export interface Verify2FASetupResponse {
    success: boolean;
    message?: string;
}

// Password management
export interface ChangePasswordRequest {
    current_password: string;
    new_password: string;
    confirm_password: string;
}

export interface VerifyPasswordResetOTPRequest {
    email: string;
    otp_code: string;
}

export interface VerifyPasswordResetOTPResponse {
    reset_token: string;
}
