import { get, post, del } from "@/axios-config/request";
import { ProfileResponse, CreateProfileRequest } from "@/types/profile.type";
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

  deleteProfile: (id: string | number) => {
    return del<ResponseAPI<any>>(`${BASE_URL}/${id}`);
  }
};
