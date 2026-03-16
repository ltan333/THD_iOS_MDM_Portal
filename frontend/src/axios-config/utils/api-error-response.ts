import { AxiosError } from "axios";

/**
 * Creates a standardized error response from an Axios error
 * Handles various API error formats:
 * - { is_success: false, error: { code, message } }  - Standard error
 * - { is_success: false, message: "..." }            - Simple error
 * - Network errors without response
 */
export const createErrorResponse = (error: AxiosError): BaseResponseAPI<null> => {
  if (!error.response) {
    return {
      message: "Lỗi kết nối. Vui lòng kiểm tra đường truyền.",
      data: null,
      is_success: false,
    };
  }

  const responseData = error.response.data as {
    is_success?: boolean;
    message?: string;
    error?: string | { code?: string; message?: string };
    data?: unknown;
  } | undefined;

  // Extract message from nested error object if available
  // API format: { is_success: false, error: { code, message } }
  let errorMessage = "Đã có lỗi xảy ra";

  if (responseData) {
    if (typeof responseData.error === "object" && responseData.error?.message) {
      // Nested error: { error: { message: "..." } }
      errorMessage = responseData.error.message;
    } else if (typeof responseData.error === "string") {
      // Simple error string: { error: "..." }
      errorMessage = responseData.error;
    } else if (responseData.message) {
      // Top-level message: { message: "..." }
      errorMessage = responseData.message;
    } else if (error.message) {
      // Fallback to Axios error message
      errorMessage = error.message;
    }
  }

  return {
    message: errorMessage,
    data: null,
    is_success: false,
  };
};
