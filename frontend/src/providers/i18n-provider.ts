"use client";

import type { I18nProvider } from "@refinedev/core";

/**
 * Vietnamese translations for Refine notifications and common messages
 * This overrides the default English messages from Refine
 */
const translations: Record<string, string> = {
    // Common actions
    "actions.create": "Tạo mới",
    "actions.edit": "Chỉnh sửa",
    "actions.delete": "Xóa",
    "actions.show": "Xem",
    "actions.list": "Danh sách",
    "actions.save": "Lưu",
    "actions.cancel": "Hủy",
    "actions.refresh": "Làm mới",
    "actions.close": "Đóng",
    "actions.confirm": "Xác nhận",

    // Notifications - Success
    "notifications.success": "Thành công",
    "notifications.createSuccess": "Tạo mới thành công!",
    "notifications.editSuccess": "Cập nhật thành công!",
    "notifications.deleteSuccess": "Xóa thành công!",

    // Notifications - Error  
    "notifications.error": "Có lỗi xảy ra",
    "notifications.createError": "Tạo mới thất bại",
    "notifications.editError": "Cập nhật thất bại",
    "notifications.deleteError": "Xóa thất bại",

    // Page titles
    "pages.create": "Tạo mới",
    "pages.edit": "Chỉnh sửa",
    "pages.show": "Chi tiết",
    "pages.list": "Danh sách",
    "pages.error.404": "Không tìm thấy trang",
    "pages.error.resource404": "Không tìm thấy",
    "pages.error.info": "Bạn có thể đã nhập sai địa chỉ, hoặc trang đã được di chuyển",
    "pages.error.backHome": "Về trang chủ",

    // Buttons
    "buttons.create": "Tạo mới",
    "buttons.delete": "Xóa",
    "buttons.edit": "Chỉnh sửa",
    "buttons.save": "Lưu",
    "buttons.cancel": "Hủy",
    "buttons.confirm": "Xác nhận",
    "buttons.filter": "Lọc",
    "buttons.clear": "Xóa bộ lọc",
    "buttons.refresh": "Làm mới",
    "buttons.show": "Xem",
    "buttons.undo": "Hoàn tác",
    "buttons.import": "Nhập",
    "buttons.export": "Xuất",
    "buttons.clone": "Nhân bản",

    // Table
    "table.actions": "Thao tác",

    // Form
    "warnWhenUnsavedChanges": "Bạn có thay đổi chưa lưu. Bạn có chắc muốn rời khỏi trang?",

    // Auth
    "pages.login.title": "Đăng nhập",
    "pages.login.signin": "Đăng nhập",
    "pages.login.signup": "Đăng ký",
    "pages.login.forgotPassword": "Quên mật khẩu?",
    "pages.login.rememberMe": "Ghi nhớ đăng nhập",

    // Undoable
    "notifications.undoable": "Bạn có",
    "notifications.seconds": "giây để hoàn tác",
};

/**
 * Simple i18n provider for Vietnamese translations
 * Used to override Refine's default English notification messages
 */
export const i18nProvider: I18nProvider = {
    translate: (key: string, params?: Record<string, unknown>) => {
        // Try to get translation from our translations object
        let translation = translations[key];

        if (!translation) {
            // Return the key itself if no translation found
            // Refine will use this as fallback
            return key;
        }

        // Replace params like {{name}} with actual values
        if (params) {
            Object.entries(params).forEach(([paramKey, value]) => {
                translation = translation?.replace(
                    new RegExp(`{{${paramKey}}}`, "g"),
                    String(value)
                );
            });
        }

        return translation;
    },
    changeLocale: async () => Promise.resolve(),
    getLocale: () => "vi",
};

export default i18nProvider;
