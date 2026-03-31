import { DeviceResponse } from "./device.type";

export interface DeviceGroupResponse {
    id: number;
    name?: string;
    description?: string;
    device_count?: number;
    devices?: DeviceResponse[];
    created_at?: string;
    updated_at?: string;
}

export interface CreateDeviceGroupRequest {
    name: string;
    description?: string;
}

export interface UpdateDeviceGroupRequest {
    name?: string;
    description?: string;
}

export interface ManageGroupDevicesRequest {
    device_ids: string[];
}
