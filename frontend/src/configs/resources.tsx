import { ResourceProps } from "@refinedev/core";
import {
    FormOutlined,
    BellOutlined,
    FileTextOutlined,
    SettingOutlined,
    TeamOutlined,
    ProjectOutlined,
    LockOutlined,
    UserOutlined,
} from "@ant-design/icons";

export const resources: ResourceProps[] = [
    // Main Resources
    {
        name: "reports",
        list: "/reports",
        show: "/reports/:id",
        meta: {
            label: "Tổng hợp báo cáo",
            icon: <FileTextOutlined />,
        },
    },

    // Settings Parent
    {
        name: "settings",
        meta: {
            label: "Cài đặt",
            icon: <SettingOutlined />,
        },
    },

    // Settings Children
    {
        name: "templates",
        identifier: "report-templates",
        list: "/settings/report-templates",
        create: "/settings/report-templates/create",
        edit: "/settings/report-templates/:id/edit",
        show: "/settings/report-templates/:id",
        meta: {
            label: "Mẫu câu hỏi báo cáo",
            icon: <FormOutlined />,
            parent: "settings",
            canDelete: true,
        },
    },
    {
        name: "notifications/schedules",
        list: "/notifications/schedules",
        create: "/notifications/schedules/create",
        edit: "/notifications/schedules/:id/edit",
        show: "/notifications/schedules/:id",
        meta: {
            label: "Lịch gửi thông báo",
            icon: <BellOutlined />,
            parent: "settings",
            canDelete: true,
        },
    },

    // Master Data
    {
        name: "departments",
        list: "/departments",
        create: "/departments/create",
        edit: "/departments/:id/edit",
        meta: {
            label: "Phòng ban",
            icon: <TeamOutlined />,
        },
    },
    {
        name: "projects",
        list: "/projects",
        create: "/projects/create",
        edit: "/projects/:id/edit",
        meta: {
            label: "Dự án",
            icon: <ProjectOutlined />,
        },
    },

    // Permissions (SysAdmin only)
    {
        name: "permissions",
        list: "/permissions",
        create: "/permissions/create",
        edit: "/permissions/:id/edit",
        meta: {
            label: "Phân quyền truy cập",
            icon: <LockOutlined />,
            canDelete: true,
        },
    },

    // Users management resource
    {
        name: "users",
        list: "/users",
        create: "/users/create",
        edit: "/users/:id/edit",
        show: "/users/:id",
        meta: {
            label: "Người dùng",
            icon: <UserOutlined />,
            canDelete: true,
        },
    },
];

