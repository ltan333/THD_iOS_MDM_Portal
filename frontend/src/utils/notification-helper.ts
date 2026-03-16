/**
 * Success Notification Helper
 * 
 * Helper function để sử dụng với Refine mutation hooks.
 * Lấy message từ API response và hiển thị trong notification.
 * 
 * Usage:
 * ```tsx
 * const { mutate } = useCreate({
 *   successNotification: successNotificationFromResponse,
 * });
 * ```
 * 
 * Hoặc với custom fallback message:
 * ```tsx
 * const { mutate } = useCreate({
 *   successNotification: successNotificationFromResponse("Tạo thành công!"),
 * });
 * ```
 */

import type { OpenNotificationParams } from "@refinedev/core";

type SuccessNotificationFn = (
    data?: unknown,
    values?: unknown,
    resource?: string
) => OpenNotificationParams | false;

/**
 * Tạo success notification function lấy message từ API response
 * 
 * @param fallbackMessage - Message mặc định nếu không có message trong response
 * @returns Function để pass vào successNotification của mutation hooks
 */
export const successNotificationFromResponse = (
    fallbackMessage: string = "Thao tác thành công!"
): SuccessNotificationFn => {
    return (data, _values, _resource) => {
        // Try to get message from response meta
        const responseData = data as { meta?: { message?: string } } | undefined;
        const message = responseData?.meta?.message || fallbackMessage;

        return {
            message: "Thành công",
            description: message,
            type: "success",
        };
    };
};

/**
 * Tạo error notification function lấy message từ API error response  
 * 
 * @param fallbackMessage - Message mặc định nếu không có message trong error
 * @returns Function để pass vào errorNotification của mutation hooks
 */
export const errorNotificationFromResponse = (
    fallbackMessage: string = "Có lỗi xảy ra"
): SuccessNotificationFn => {
    return (error, _values, _resource) => {
        const errorData = error as { message?: string } | undefined;
        const message = errorData?.message || fallbackMessage;

        return {
            message: "Lỗi",
            description: message,
            type: "error",
        };
    };
};

/**
 * Default notification configs để reuse trong các forms
 */
export const defaultNotifications = {
    create: {
        successNotification: successNotificationFromResponse("Tạo mới thành công!"),
        errorNotification: errorNotificationFromResponse("Lỗi khi tạo mới"),
    },
    update: {
        successNotification: successNotificationFromResponse("Cập nhật thành công!"),
        errorNotification: errorNotificationFromResponse("Lỗi khi cập nhật"),
    },
    delete: {
        successNotification: successNotificationFromResponse("Xóa thành công!"),
        errorNotification: errorNotificationFromResponse("Lỗi khi xóa"),
    },
};

export default successNotificationFromResponse;
