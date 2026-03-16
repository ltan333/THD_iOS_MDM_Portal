"use client";

import { toast } from 'sonner';
import type { NotificationProvider } from '@refinedev/core';

/**
 * Vietnamese translations for Refine's default notification messages
 */
const VIETNAMESE_MESSAGES: Record<string, string> = {
    // Success messages 
    "Successful": "Thao tác thành công!",
    "Successfully created": "Tạo mới thành công!",
    "Successfully updated": "Cập nhật thành công!",
    "Successfully deleted": "Xóa thành công!",
    "Created successfully": "Tạo mới thành công!",
    "Updated successfully": "Cập nhật thành công!",
    "Deleted successfully": "Xóa thành công!",

    // Error messages
    "Error": "Có lỗi xảy ra",
    "Error when fetching": "Lỗi khi tải dữ liệu",
    "Something went wrong": "Đã có lỗi xảy ra",
    "There was an error creating": "Lỗi khi tạo mới",
    "There was an error updating": "Lỗi khi cập nhật",
    "There was an error deleting": "Lỗi khi xóa",
};

/**
 * Translate message to Vietnamese if available
 */
const translateMessage = (message: string | undefined): string | undefined => {
    if (!message) return message;

    // Check for exact match first
    if (VIETNAMESE_MESSAGES[message]) {
        return VIETNAMESE_MESSAGES[message];
    }

    // Check for partial matches (e.g., "Successfully created user")
    for (const [key, translation] of Object.entries(VIETNAMESE_MESSAGES)) {
        if (message.toLowerCase().includes(key.toLowerCase())) {
            return translation;
        }
    }

    return message;
};

/**
 * Custom notification provider using Sonner
 * Replaces the default Ant Design notification provider from @refinedev/antd
 * Automatically translates Refine's English messages to Vietnamese
 */
export const sonnerNotificationProvider: NotificationProvider = {
    open: ({ message, description, type, key }) => {
        // Translate both message and description to Vietnamese
        const translatedMessage = translateMessage(message);
        const translatedDescription = translateMessage(description);

        // For errors, prefer description (usually more detailed)
        // For success, prefer description if available, otherwise message
        const content = type === 'error'
            ? (translatedDescription || translatedMessage)
            : (translatedDescription || translatedMessage);

        switch (type) {
            case 'success':
                toast.success(content, { id: key });
                break;
            case 'error':
                toast.error(content, { id: key });
                break;
            case 'progress':
                toast.loading(content, { id: key });
                break;
            default:
                toast.info(content, { id: key });
        }
    },
    close: (key) => {
        toast.dismiss(key);
    },
};

export default sonnerNotificationProvider;

