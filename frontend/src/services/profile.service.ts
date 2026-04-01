import { get, post, put, del } from "@/axios-config/request";
import { ResponseAPI } from "@/types";
import {
    ProfileResponse,
    ProfileListData,
    CreateProfileRequest,
    UpdateProfileRequest,
    UpdateProfileStatusRequest,
    AssignProfileRequest,
    ProfileDeploymentStatusResponse,
} from "@/types/profile.type";

const BASE_URL = "/profiles";

export const profileService = {
    getProfiles: (params?: {
        page?: number;
        limit?: number;
        search?: string;
        status?: string;
        platform?: string;
        scope?: string;
    }) => {
        return get<ResponseAPI<ProfileListData>>(BASE_URL, { params });
    },

    getProfileById: (id: number) => {
        return get<ResponseAPI<ProfileResponse>>(`${BASE_URL}/${id}`);
    },

    createProfile: (payload: CreateProfileRequest) => {
        return post<ResponseAPI<ProfileResponse>, CreateProfileRequest>(BASE_URL, payload);
    },

    updateProfile: (id: number, payload: UpdateProfileRequest) => {
        return put<ResponseAPI<ProfileResponse>, UpdateProfileRequest>(`${BASE_URL}/${id}`, payload);
    },

    deleteProfile: (id: number) => {
        return del<ResponseAPI<unknown>>(`${BASE_URL}/${id}`);
    },

    updateStatus: (id: number, payload: UpdateProfileStatusRequest) => {
        return put<ResponseAPI<ProfileResponse>, UpdateProfileStatusRequest>(
            `${BASE_URL}/${id}/status`,
            payload
        );
    },

    repush: (id: number) => {
        return post<ResponseAPI<unknown>, Record<string, never>>(`${BASE_URL}/${id}/repush`, {});
    },

    duplicate: (id: number) => {
        return post<ResponseAPI<ProfileResponse>, Record<string, never>>(
            `${BASE_URL}/${id}/duplicate`,
            {}
        );
    },

    assign: (id: number, payload: AssignProfileRequest) => {
        return post<ResponseAPI<unknown>, AssignProfileRequest>(
            `${BASE_URL}/${id}/assignments`,
            payload
        );
    },

    getDeploymentStatus: (id: number) => {
        return post<ResponseAPI<ProfileDeploymentStatusResponse[]>, Record<string, never>>(
            `${BASE_URL}/${id}/deployment-status`,
            {}
        );
    },
};
