import { get, post, put, del } from "@/axios-config/request";
import { DeviceGroupResponse, CreateDeviceGroupRequest, UpdateDeviceGroupRequest, ManageGroupDevicesRequest } from "@/types/device-group.type";
import { ResponseAPI, ListResponse } from "@/types/api.type";

const BASE_URL = "/device-groups";

export const deviceGroupService = {
  getGroups: (params?: { page?: number; limit?: number; search?: string }) => {
    return get<ResponseAPI<ListResponse<DeviceGroupResponse>>>(BASE_URL, { params });
  },

  getGroupById: (id: number) => {
    return get<ResponseAPI<DeviceGroupResponse>>(`${BASE_URL}/${id}`);
  },

  createGroup: (payload: CreateDeviceGroupRequest) => {
    return post<ResponseAPI<DeviceGroupResponse>, CreateDeviceGroupRequest>(BASE_URL, payload);
  },

  updateGroup: (id: number, payload: UpdateDeviceGroupRequest) => {
    return put<ResponseAPI<DeviceGroupResponse>, UpdateDeviceGroupRequest>(`${BASE_URL}/${id}`, payload);
  },

  deleteGroup: (id: number) => {
    return del<ResponseAPI<any>>(`${BASE_URL}/${id}`);
  },

  addDevicesToGroup: (id: number, payload: ManageGroupDevicesRequest) => {
    return post<ResponseAPI<any>, ManageGroupDevicesRequest>(`${BASE_URL}/${id}/devices`, payload);
  },

  removeDeviceFromGroup: (groupId: number, deviceId: string) => {
    return del<ResponseAPI<any>>(`${BASE_URL}/${groupId}/devices/${deviceId}`);
  }
};
