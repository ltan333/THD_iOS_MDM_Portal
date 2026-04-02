import { get, post, put, del } from "@/axios-config/request";
import {
  ProfileResponse,
  CreateProfileRequest,
  UpdateProfileRequest,
  AssignProfileRequest,
  ProfileAssignmentResponse,
  UpdateProfileStatusRequest,
  ProfileDeploymentStatusResponse,
} from "@/types/profile.type";
import { ListResponseAPI, ResponseAPI } from "@/types";

const BASE_URL = "/profiles";

export const profileService = {
  getProfiles: (queryParams?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    platform?: string;
    name?: string;
    scope?: string;
  }) => {
    return get<ListResponseAPI<ProfileResponse>>(BASE_URL, { queryParams });
  },

  getProfileById: (id: string | number) => {
    return get<ResponseAPI<ProfileResponse>>(`${BASE_URL}/${id}`);
  },

  createProfile: (payload: CreateProfileRequest) => {
    return post<ResponseAPI<ProfileResponse>, CreateProfileRequest>(BASE_URL, payload);
  },

  updateProfile: (id: string | number, payload: UpdateProfileRequest) => {
    return put<ResponseAPI<ProfileResponse>, UpdateProfileRequest>(`${BASE_URL}/${id}`, payload);
  },

  updateProfileStatus: (id: string | number, status: "active" | "draft" | "archived") => {
    return put<ResponseAPI<any>, UpdateProfileStatusRequest>(`${BASE_URL}/${id}/status`, { status });
  },

  assignProfile: (id: string | number, payload: AssignProfileRequest) => {
    return post<ResponseAPI<any>, AssignProfileRequest>(`${BASE_URL}/${id}/assignments`, payload);
  },

  getProfileAssignments: (id: string | number) => {
    return get<ResponseAPI<ProfileAssignmentResponse[]>>(`${BASE_URL}/${id}/assignments`);
  },

  deleteAssignment: (profileId: string | number, assignmentId: string | number) => {
    return del<ResponseAPI<any>>(`${BASE_URL}/${profileId}/assignments/${assignmentId}`);
  },

  getDeploymentStatus: (id: string | number) => {
    return get<ResponseAPI<ProfileDeploymentStatusResponse[]>>(`${BASE_URL}/${id}/deployment-status`);
  },

  repushProfile: (id: string | number) => {
    return post<ResponseAPI<any>, {}>(`${BASE_URL}/${id}/repush`, {});
  },

  getProfileVersions: (id: string | number) => {
    return get<ResponseAPI<any[]>>(`${BASE_URL}/${id}/versions`);
  },

  duplicateProfile: (id: string | number) => {
    return post<ResponseAPI<ProfileResponse>, {}>(`${BASE_URL}/${id}/duplicate`, {});
  },

  deleteProfile: (id: string | number) => {
    return del<ResponseAPI<any>>(`${BASE_URL}/${id}`);
  },
};
