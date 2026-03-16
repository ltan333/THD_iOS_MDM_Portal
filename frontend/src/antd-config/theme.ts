/**
 * Ant Design Theme Configuration
 * 
 * Centralized theme tokens matching TailwindCSS theme
 */

import type { ThemeConfig } from "antd";

// THD Brand Colors matching globals.css
export const THD_COLORS = {
    primary: {
        50: "#fff1f0",
        100: "#ffccc7",
        200: "#ffa39e",
        300: "#ff7875",
        400: "#ff4d4f",
        500: "#c00f0c",
        600: "#a00a08",
        700: "#800806",
        800: "#600604",
        900: "#400402",
    },
    secondary: {
        50: "#f0f5ff",
        100: "#d6e4ff",
        200: "#adc6ff",
        300: "#85a5ff",
        400: "#597ef7",
        500: "#2f54eb",
        600: "#1d39c4",
        700: "#10239e",
        800: "#061178",
        900: "#030852",
    },
    success: {
        500: "#52c41a",
        600: "#389e0d",
    },
    warning: {
        500: "#faad14",
        600: "#d48806",
    },
    error: {
        500: "#ff4d4f",
        600: "#cf1322",
    },
    gray: {
        50: "#fafafa",
        100: "#f5f5f5",
        200: "#e8e8e8",
        300: "#d9d9d9",
        400: "#bfbfbf",
        500: "#8c8c8c",
        600: "#595959",
        700: "#434343",
        800: "#262626",
        900: "#1f1f1f",
    },
};

// Light theme configuration
export const lightTheme: ThemeConfig = {
    token: {
        colorPrimary: THD_COLORS.primary[500],
        colorSuccess: THD_COLORS.success[500],
        colorWarning: THD_COLORS.warning[500],
        colorError: THD_COLORS.error[500],
        colorInfo: "#1890ff",
        borderRadius: 6,
        fontFamily: "Montserrat, system-ui, sans-serif",
        colorBgContainer: "#ffffff",
        colorBgElevated: "#ffffff",
        colorBgLayout: "#f5f5f5",
        colorText: "#262626",
        colorTextSecondary: "#595959",
        colorBorder: "#d9d9d9",
    },
    components: {
        Button: {
            borderRadius: 6,
        },
        Card: {
            borderRadius: 8,
        },
        Table: {
            borderRadius: 8,
        },
        Modal: {
            borderRadius: 8,
        },
        Input: {
            borderRadius: 6,
        },
        Select: {
            borderRadius: 6,
        },
    },
};

// Dark theme configuration
export const darkTheme: ThemeConfig = {
    token: {
        colorPrimary: THD_COLORS.primary[500],
        colorSuccess: THD_COLORS.success[500],
        colorWarning: THD_COLORS.warning[500],
        colorError: THD_COLORS.error[500],
        colorInfo: "#1890ff",
        borderRadius: 6,
        fontFamily: "Montserrat, system-ui, sans-serif",
        colorBgContainer: "#141414",
        colorBgElevated: "#1a1a1a",
        colorBgLayout: "#0a0a0a",
        colorText: "#f0f0f0",
        colorTextSecondary: "#8c8c8c",
        colorBorder: "#333333",
    },
    components: {
        Button: {
            borderRadius: 6,
        },
        Card: {
            borderRadius: 8,
        },
        Table: {
            borderRadius: 8,
        },
        Modal: {
            borderRadius: 8,
        },
        Input: {
            borderRadius: 6,
        },
        Select: {
            borderRadius: 6,
        },
    },
};

// Get theme config based on mode
export function getThemeConfig(mode: "light" | "dark"): ThemeConfig {
    return mode === "dark" ? darkTheme : lightTheme;
}
