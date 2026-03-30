"use client";

import {
    LayoutDashboard,
    FileText,
    ClipboardList,
    FolderKanban,
    Building2,
    Users,
    Shield,
    Bell,
    User,
    History,
} from "lucide-react";
import type { NavItem } from "@/types";

/**
 * Unified navigation items for sidebar.
 * Items are grouped by sections and filtered by user role.
 */
export const navItems: NavItem[] = [
    // ==================== ADMIN ITEMS ====================
    // Dashboard
    {
        key: "dashboard",
        label: "Dashboard",
        icon: <LayoutDashboard size={24} strokeWidth={1.5} />,
        href: "/dashboard",
        roles: ["ADMIN", "SYSTEM_ADMIN"],
    },

    // Reports
    {
        key: "reports",
        label: "Báo cáo",
        icon: <FileText size={24} strokeWidth={1.5} />,
        href: "/reports",
        roles: ["ADMIN", "SYSTEM_ADMIN"],
    },

    // Templates
    {
        key: "templates",
        label: "Mẫu báo cáo",
        icon: <ClipboardList size={24} strokeWidth={1.5} />,
        href: "/templates",
        roles: ["ADMIN", "SYSTEM_ADMIN"],
        children: [
            {
                key: "templates-list",
                label: "Danh sách mẫu",
                href: "/templates",
            },
            {
                key: "templates-active",
                label: "Mẫu đang áp dụng",
                href: "/templates/active",
            },
        ],
    },

    // Projects
    {
        key: "projects",
        label: "Dự án",
        icon: <FolderKanban size={24} strokeWidth={1.5} />,
        href: "/projects",
        roles: ["ADMIN", "SYSTEM_ADMIN"],
    },

    // Departments
    {
        key: "departments",
        label: "Phòng ban",
        icon: <Building2 size={24} strokeWidth={1.5} />,
        href: "/departments",
        roles: ["ADMIN", "SYSTEM_ADMIN"],
    },

    // Users
    {
        key: "users",
        label: "Người dùng",
        icon: <Users size={24} strokeWidth={1.5} />,
        href: "/users",
        roles: ["ADMIN", "SYSTEM_ADMIN"],
        children: [
            {
                key: "users-list",
                label: "Danh sách",
                href: "/users",
            },
            {
                key: "users-import",
                label: "Import người dùng",
                href: "/users/import",
            },
        ],
    },

    // Permissions
    {
        key: "permissions",
        label: "Phân quyền",
        icon: <Shield size={24} strokeWidth={1.5} />,
        href: "/permissions",
        roles: ["ADMIN", "SYSTEM_ADMIN"],
    },

    // Notifications
    {
        key: "notifications",
        label: "Thông báo",
        icon: <Bell size={24} strokeWidth={1.5} />,
        href: "/notifications",
        roles: ["ADMIN", "SYSTEM_ADMIN"],
        children: [
            {
                key: "notifications-schedules",
                label: "Lịch gửi thông báo",
                href: "/notifications/schedules",
            },
            {
                key: "notifications-send",
                label: "Gửi ngay",
                href: "/notifications/send",
            },
        ],
    },

    // ==================== PERSONAL ITEMS (No section) ====================
    // Personal Information - shows for ALL roles
    {
        key: "user-profile",
        label: "Thông tin cá nhân",
        icon: <User size={24} strokeWidth={1.5} />,
        href: "/profile",
    },

    // Work Report
    {
        key: "user-report",
        label: "Báo cáo công việc",
        icon: <FileText size={24} strokeWidth={1.5} />,
        href: "/my-report",
        roles: ["USER"],
    },

    // Report History
    {
        key: "user-history",
        label: "Lịch sử báo cáo",
        icon: <History size={24} strokeWidth={1.5} />,
        href: "/my-history",
        roles: ["USER"],
    },
];

/**
 * Filter navigation items by user role
 */
export function filterNavItemsByRole(
    items: NavItem[],
    userRole?: "USER" | "ADMIN" | "SYSTEM_ADMIN"
): NavItem[] {
    return items.filter((item) => {
        // If no roles specified, everyone can see
        if (!item.roles || item.roles.length === 0) {
            return true;
        }
        // Check if user role is in allowed roles
        return userRole && item.roles.includes(userRole);
    });
}

// Legacy exports for backward compatibility - includes admin items + items without roles
export const defaultNavItems = navItems.filter(
    (item) => !item.roles || item.roles.length === 0 || item.roles.some((r) => ["ADMIN", "SYSTEM_ADMIN"].includes(r))
);
