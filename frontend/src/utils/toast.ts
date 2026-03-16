import { toast as sonnerToast } from 'sonner';

/**
 * Toast Service - Quản lý thông báo toast cho toàn bộ dự án
 * Sử dụng Sonner toast library
 */

// ==================== TOAST FUNCTIONS ====================

interface ToastOptions {
  duration?: number; // milliseconds
  id?: string;
  onDismiss?: () => void;
  onAutoClose?: () => void;
  description?: string;
}

const DEFAULT_DURATION = 3000; // milliseconds

/**
 * Hiển thị toast thành công
 */
export const toastSuccess = (message: string, options?: ToastOptions) => {
  return sonnerToast.success(message, {
    duration: options?.duration ?? DEFAULT_DURATION,
    id: options?.id,
    onDismiss: options?.onDismiss,
    onAutoClose: options?.onAutoClose,
    description: options?.description,
  });
};

/**
 * Hiển thị toast lỗi
 */
export const toastError = (message: string, options?: ToastOptions) => {
  return sonnerToast.error(message, {
    duration: options?.duration ?? 5000, // Error hiển thị lâu hơn
    id: options?.id,
    onDismiss: options?.onDismiss,
    onAutoClose: options?.onAutoClose,
    description: options?.description,
  });
};

/**
 * Hiển thị toast cảnh báo
 */
export const toastWarning = (message: string, options?: ToastOptions) => {
  return sonnerToast.warning(message, {
    duration: options?.duration ?? 4000,
    id: options?.id,
    onDismiss: options?.onDismiss,
    onAutoClose: options?.onAutoClose,
    description: options?.description,
  });
};

/**
 * Hiển thị toast thông tin
 */
export const toastInfo = (message: string, options?: ToastOptions) => {
  return sonnerToast.info(message, {
    duration: options?.duration ?? DEFAULT_DURATION,
    id: options?.id,
    onDismiss: options?.onDismiss,
    onAutoClose: options?.onAutoClose,
    description: options?.description,
  });
};

/**
 * Hiển thị toast loading
 */
export const toastLoading = (message: string, id?: string) => {
  return sonnerToast.loading(message, {
    id: id ?? 'loading',
  });
};

/**
 * Đóng toast theo id
 */
export const toastDismiss = (id?: string | number) => {
  if (id) {
    sonnerToast.dismiss(id);
  } else {
    sonnerToast.dismiss();
  }
};

// ==================== HELPER FUNCTIONS ====================

/**
 * Hiển thị toast lỗi từ API response
 */
export const toastApiError = (error: any, fallbackMessage = 'Đã có lỗi xảy ra') => {
  const errorMessage =
    error?.response?.data?.message ||
    error?.message ||
    fallbackMessage;

  return toastError(errorMessage);
};

/**
 * Hiển thị toast lỗi validation form
 */
export const toastValidationError = (message = 'Vui lòng điền đầy đủ thông tin bắt buộc!') => {
  return toastError(message);
};

/**
 * Hiển thị toast khi thao tác CRUD thành công
 */
export const toastCrudSuccess = (action: 'create' | 'update' | 'delete', entity: string) => {
  const messages = {
    create: `Tạo ${entity} thành công!`,
    update: `Cập nhật ${entity} thành công!`,
    delete: `Xóa ${entity} thành công!`,
  };
  return toastSuccess(messages[action]);
};

/**
 * Hiển thị toast khi thao tác CRUD thất bại
 */
export const toastCrudError = (action: 'create' | 'update' | 'delete', entity: string, error?: any) => {
  const messages = {
    create: `Tạo ${entity} thất bại`,
    update: `Cập nhật ${entity} thất bại`,
    delete: `Xóa ${entity} thất bại`,
  };

  const errorDetail = error?.response?.data?.message || error?.message;
  const fullMessage = errorDetail
    ? `${messages[action]}: ${errorDetail}`
    : messages[action];

  return toastError(fullMessage);
};

/**
 * Promise toast - hiển thị loading, success, error tự động
 */
export const toastPromise = <T>(
  promise: Promise<T>,
  messages: {
    loading: string;
    success: string | ((data: T) => string);
    error: string | ((error: any) => string);
  }
) => {
  return sonnerToast.promise(promise, messages);
};

// ==================== TOAST OBJECT (Alternative API) ====================

/**
 * Toast object với các methods tiện dụng
 * Usage: toast.success('Message'), toast.error('Message')
 */
export const toast = {
  success: toastSuccess,
  error: toastError,
  warning: toastWarning,
  info: toastInfo,
  loading: toastLoading,
  dismiss: toastDismiss,
  promise: toastPromise,

  // API helpers
  apiError: toastApiError,
  validationError: toastValidationError,
  crudSuccess: toastCrudSuccess,
  crudError: toastCrudError,
};

export default toast;
