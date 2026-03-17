export const TOKEN_KEYS = {
  // IMPORTANT: These must match backend cookie names in config.go
  // COOKIE_NAME=rcycles_token, COOKIE_REFRESH_NAME=rcycles_refresh
  ACCESS_TOKEN: "rcycles_token",
  REFRESH_TOKEN: "rcycles_refresh",
  ORG_TOKEN: "org_token",
} as const;

export const API_VERSION = "/api";

export const HEADER_KEYS = {
  CONTENT_TYPE: "Content-Type",
  AUTHORIZATION: "Authorization",
  X_ORGANIZATION_KEY: "X-Organization-Key",
  NGROK_SKIP_WARNING: "ngrok-skip-browser-warning",
} as const;

export const METHODS = {
  GET: "GET",
  POST: "POST",
  PUT: "PUT",
  DELETE: "DELETE",
  PATCH: "PATCH",
} as const;

export type MethodType = (typeof METHODS)[keyof typeof METHODS];

export const RETRY_CONFIG = {
  MAX_RETRIES: 2, // Only for network timeouts
  BASE_DELAY_MS: 1000,
} as const;

export const AUTH_CONFIG = {
  // Core authentication
  LOGIN_ENDPOINT: `/auth/login`,
  LOGOUT_ENDPOINT: `/auth/logout`,
  REFRESH_TOKEN_ENDPOINT: `/auth/refresh-token`,
  ME_ENDPOINT: `/users/me`,

  // 2FA (Two-Factor Authentication)
  SETUP_2FA_ENDPOINT: `/auth/setup-2fa`,
  VERIFY_2FA_SETUP_ENDPOINT: `/auth/verify-2fa-setup`,
  VERIFY_OTP_ENDPOINT: `/auth/verify-otp`,

  // Password management
  CHANGE_PASSWORD_ENDPOINT: `/auth/password`,
  REQUEST_PASSWORD_RESET_ENDPOINT: `/auth/password/reset/request`,
  VERIFY_PASSWORD_RESET_OTP_ENDPOINT: `/auth/password/reset/verify`,
  COMPLETE_PASSWORD_RESET_ENDPOINT: `/auth/password/reset/complete`,
} as const;

export const ANALYTICS_ENDPOINTS = {
  OVERVIEW: `/analytics/overview`,
  EMPLOYEES: `/analytics/employees`,
  ATTENDANCE: `/analytics/attendance`,
  SCHEDULE: `/analytics/schedule`,
  SALARY: `/analytics/salary`,
  TRENDS: `/analytics/trends`,
  COMPARISON: `/analytics/comparison`,
  PERFORMANCE_RANKING: `/analytics/performance-ranking`,
  EXPORT: `/analytics/export`,
} as const;

export const REPORT_TEMPLATE_ENDPOINTS = {
  BASE: `/templates`,
  ACTIVE: `/templates/active`,
  BY_ID: (id: number) => `/templates/${id}`,
  ACTIVATE: (id: number) => `/templates/${id}/status`,
} as const;

export const QUESTION_ENDPOINTS = {
  BY_TEMPLATE: (templateId: number) => `/templates/${templateId}/questions`,
  BY_ID: (templateId: number, questionId: number) =>
    `/templates/${templateId}/questions/${questionId}`,
} as const;

export const REPORT_ENDPOINTS = {
  BASE: `/reports`,
  BY_ID: (id: number) => `/reports/${id}`,
} as const;