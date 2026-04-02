export interface DeviceResponse {
    id: string;
    udid?: string;
    name?: string;
    serial_number?: string;
    model?: string;
    os_version?: string;
    platform?: string;
    device_type?: string;
    battery_level?: number;
    storage_capacity?: number;
    storage_used?: number;
    status?: string;
    compliance_status?: string;
    enrollment_type?: string;
    is_enrolled?: boolean;
    is_jailbroken?: boolean;
    mac_address?: string;
    ip_address?: string;
    owner_id?: number;
    enrolled_at?: string;
    last_seen?: string;
    created_at?: string;
    updated_at?: string;
}

export interface DeviceActionResponse {
    command_uuid?: string;
    status?: string;
    request_type?: string;
    message?: string;
}

export type DeviceWipeObliterationBehavior =
    | "Default"
    | "DoNotObliterate"
    | "ObliterateWithWarning"
    | "Always";

export interface DeviceLockRequest {
    footnote?: string;
    message?: string;
    phone_number?: string;
    pin?: string;
}

export interface DeviceWipeRequest {
    disallow_proximity_setup?: boolean;
    obliteration_behavior?: DeviceWipeObliterationBehavior;
    pin?: string;
    preserve_data_plan?: boolean;
}

export interface DeviceRestartRequest {
    notify_user?: boolean;
}

export interface DeviceCommandResult {
    command_uuid?: string;
    request_type?: string;
    no_push?: boolean;
    push_error?: string;
    command_error?: string;
}

export interface DeviceStatsResponse {
    total?: number;
    active?: number;
    inactive?: number;
    enrolled?: number;
    compliant?: number;
    non_compliant?: number;
    by_platform?: Record<string, number>;
    by_status?: Record<string, number>;
}
