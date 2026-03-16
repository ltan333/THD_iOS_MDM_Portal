/**
 * User Types - User related types
 */

import type { UserRole, Permission } from "./auth.type";

export interface UserResponse {
    id: number;
    username: string;
    email: string;
    full_name?: string;
    phone?: string;
    role: UserRole;
    department_id?: number;
    department?: {
        id: number;
        name: string;
    };
    project_id?: number;
    project?: {
        id: number;
        name: string;
    };
    permissions?: Permission[];
    avatar_url?: string;
    is_active: boolean;
    created_at: string;
    updated_at?: string;
}

export interface UserCreateRequest {
    username: string;
    email: string;
    password: string;
    full_name?: string;
    phone?: string;
    role: UserRole;
    department_id?: number;
    project_id?: number;
    permission_ids?: number[];
}

export interface UserUpdateRequest {
    email?: string;
    full_name?: string;
    phone?: string;
    role?: UserRole;
    department_id?: number;
    project_id?: number;
    permission_ids?: number[];
    is_active?: boolean;
}

export interface UserProfileUpdateRequest {
    full_name?: string;
    phone?: string;
    avatar_url?: string;
}
