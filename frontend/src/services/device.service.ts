import { get, post, del } from "@/axios-config/request";
import { DeviceResponse, DeviceActionResponse, DeviceStatsResponse } from "@/types/device.type";
import { ResponseAPI, ListResponseAPI } from "@/types";

const BASE_URL = "/devices";

export const deviceService = {
  // Fetch devices with pagination and filters
  getDevices: (queryParams?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    platform?: string;
    model?: string;
    serial_number?: string;
  }) => {
    return get<ListResponseAPI<DeviceResponse>>(BASE_URL, { queryParams });
  },

  // Get single device by ID
  getDeviceById: (id: string) => {
    return get<ResponseAPI<DeviceResponse>>(`${BASE_URL}/${id}`);
  },

  // Device Actions
  lockDevice: (id: string, payload: { message?: string; phone_number?: string; pin?: string }) => {
    return post<ResponseAPI<DeviceActionResponse>, typeof payload>(`${BASE_URL}/${id}/lock`, payload);
  },

  wipeDevice: (id: string, payload: { pin?: string; preserve_data_plan?: boolean; disallow_proximity_setup?: boolean; obliteration_behavior?: string }) => {
    return post<ResponseAPI<DeviceActionResponse>, typeof payload>(`${BASE_URL}/${id}/wipe`, payload);
  },

  restartDevice: (id: string, notify_user: boolean = false) => {
    return post<ResponseAPI<DeviceActionResponse>, { notify_user: boolean }>(`${BASE_URL}/${id}/restart`, { notify_user });
  },

  shutdownDevice: (id: string) => {
    return post<ResponseAPI<DeviceActionResponse>, {}>(`${BASE_URL}/${id}/shutdown`, {});
  },

  requestInfo: (id: string, queries: string[]) => {
    return post<ResponseAPI<DeviceActionResponse>, { queries: string[] }>(`${BASE_URL}/${id}/request-info`, { queries });
  },

  installProfile: (id: string, profile_id: number) => {
    return post<ResponseAPI<DeviceActionResponse>, { profile_id: number }>(`${BASE_URL}/${id}/install-profile`, { profile_id });
  },

  removeProfile: (id: string, profile_identifier: string) => {
    return post<ResponseAPI<DeviceActionResponse>, { profile_identifier: string }>(`${BASE_URL}/${id}/remove-profile`, { profile_identifier });
  },

  // Get device statistics
  getDeviceStats: () => {
    return get<ResponseAPI<DeviceStatsResponse>>(`${BASE_URL}/stats`);
  }
};
