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
