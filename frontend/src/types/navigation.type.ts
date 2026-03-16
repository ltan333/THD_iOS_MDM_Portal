/**
 * Navigation Types - Sidebar and navigation types
 */

import type { ReactNode } from "react";
import type { UserRole } from "./auth.type";

export interface NavItem {
    key: string;
    label: string;
    icon?: ReactNode;
    href?: string;
    path?: string;
    children?: NavItem[];
    roles?: UserRole[];
    permissions?: string[];
    badge?: number | string;
    isNew?: boolean;
}

export interface NavGroup {
    key: string;
    label: string;
    items: NavItem[];
    roles?: UserRole[];
}

export interface BreadcrumbItem {
    label: string;
    path?: string;
    icon?: ReactNode;
}
