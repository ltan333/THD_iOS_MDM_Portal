export interface ProfileResponse {
    id: number;
    name: string;
    platform: string; // ios, android, windows, macos, all
    scope: string;    // device, user, group
    status: string;   // active, draft, archived
    security_settings?: Record<string, unknown>;
    network_config?: Record<string, unknown>;
    restrictions?: Record<string, unknown>;
    content_filter?: Record<string, unknown>;
    compliance_rules?: Record<string, unknown>;
    payloads?: Record<string, unknown>;
    version: number;
    created_at: string;
    updated_at: string;
}

export interface ProfileListData {
    items: ProfileResponse[];
    total: number;
    page: number;
    limit: number;
    total_pages: number;
}

export interface CreateProfileRequest {
    name: string;
    platform?: string;
    scope?: string;
    security_settings?: Record<string, unknown>;
    network_config?: Record<string, unknown>;
    restrictions?: Record<string, unknown>;
    content_filter?: Record<string, unknown>;
    compliance_rules?: Record<string, unknown>;
}

export interface UpdateProfileRequest {
    name?: string;
    platform?: string;
    scope?: string;
    security_settings?: Record<string, unknown>;
    network_config?: Record<string, unknown>;
    restrictions?: Record<string, unknown>;
    content_filter?: Record<string, unknown>;
}

export interface UpdateProfileStatusRequest {
    status: "active" | "draft" | "archived";
}

export interface AssignProfileRequest {
    target_type: "device" | "group";
    device_id?: string;
    group_id?: number;
    schedule_type?: string;
}

export interface ProfileDeploymentStatusResponse {
    id: number;
    profile_id: number;
    device_id: string;
    status: string;   // pending, success, failed
    error_message?: string;
    applied_at?: string;
    created_at: string;
    updated_at: string;
}
