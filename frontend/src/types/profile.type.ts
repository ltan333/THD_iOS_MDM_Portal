export interface ProfileResponse {
    id: number;
    name: string;
    platform: string;
    scope: string;
    status: string;
    version: number;
    compliance_rules?: any;
    content_filter?: any;
    network_config?: any;
    payloads?: any;
    restrictions?: any;
    security_settings?: any;
    created_at: string;
    updated_at: string;
}

export interface CreateProfileRequest {
    name: string;
    platform: string;
    scope: string;
    compliance_rules?: any;
    content_filter?: any;
    network_config?: any;
    payloads?: any;
    restrictions?: any;
    security_settings?: any;
}

export interface UpdateProfileRequest {
    name: string;
    platform: string;
    scope: string;
    compliance_rules?: any;
    content_filter?: any;
    network_config?: any;
    payloads?: any;
    restrictions?: any;
    security_settings?: any;
}

export interface AssignProfileRequest {
    target_type: "device" | "group";
    device_id?: string;
    group_id?: number;
    schedule_type?: "immediate" | "scheduled";
    scheduled_at?: string;
}

export interface ProfileAssignmentResponse {
    id: number;
    profile_id: number;
    target_type: "device" | "group" | string;
    device_id?: string;
    group_id?: number;
    schedule_type?: string;
    scheduled_at?: string;
    created_at?: string;
}
