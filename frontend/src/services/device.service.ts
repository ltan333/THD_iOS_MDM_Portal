import { get, post, del } from "@/axios-config/request";
import {
  DeviceResponse,
  DeviceActionResponse,
  DeviceStatsResponse,
  DeviceLockRequest,
  DeviceWipeRequest,
  DeviceRestartRequest,
  DeviceCommandResult,
} from "@/types/device.type";
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
  lockDevice: (id: string, payload?: DeviceLockRequest) => {
    return post<ResponseAPI<DeviceCommandResult>, DeviceLockRequest | undefined>(`${BASE_URL}/${id}/lock`, payload);
  },

  unlockDevice: (id: string) => {
    return post<ResponseAPI<DeviceCommandResult>, undefined>(`${BASE_URL}/${id}/unlock`);
  },

  wipeDevice: (id: string, payload?: DeviceWipeRequest) => {
    return post<ResponseAPI<DeviceCommandResult>, DeviceWipeRequest | undefined>(`${BASE_URL}/${id}/wipe`, payload);
  },

  restartDevice: (id: string, payload?: DeviceRestartRequest) => {
    return post<ResponseAPI<DeviceCommandResult>, DeviceRestartRequest | undefined>(`${BASE_URL}/${id}/restart`, payload);
  },

  shutdownDevice: (id: string) => {
    return post<ResponseAPI<DeviceCommandResult>, undefined>(`${BASE_URL}/${id}/shutdown`);
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
