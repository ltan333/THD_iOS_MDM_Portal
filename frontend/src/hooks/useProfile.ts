import { useCustomMutation } from "@refinedev/core";
import type { HttpError } from "@refinedev/core";
import type { UserResponse } from "@/types";

// Types for profile update
interface UpdateProfileRequest {
    name: string;
    email: string;
}

// Types for password change
interface ChangePasswordRequest {
    current_password: string;
    new_password: string;
    confirm_password: string;
}

/**
 * Hook to update current user's profile
 * Uses Refine's useCustomMutation with PUT /users/me
 */
export function useUpdateProfile() {
    const { mutate, mutation } = useCustomMutation<UserResponse, HttpError, UpdateProfileRequest>();

    const updateProfile = (
        data: UpdateProfileRequest,
        callbacks?: {
            onSuccess?: (data: UserResponse) => void;
            onError?: (error: HttpError) => void;
        }
    ) => {
        mutate(
            {
                url: "/users/me",
                method: "put",
                values: data,
                errorNotification: false,
                successNotification: false,
            },
            {
                onSuccess: (response) => {
                    callbacks?.onSuccess?.(response.data);
                },
                onError: (error) => {
                    callbacks?.onError?.(error);
                },
            }
        );
    };

    return {
        updateProfile,
        isPending: mutation.isPending,
        isError: mutation.isError,
        error: mutation.error,
    };
}

/**
 * Hook to change current user's password
 * Uses Refine's useCustomMutation with PUT /auth/password
 */
export function useChangePassword() {
    const { mutate, mutation } = useCustomMutation<Record<string, never>, HttpError, ChangePasswordRequest>();

    const changePassword = (
        data: ChangePasswordRequest,
        callbacks?: {
            onSuccess?: () => void;
            onError?: (error: HttpError) => void;
        }
    ) => {
        mutate(
            {
                url: "/auth/password",
                method: "put",
                values: data,
                errorNotification: false,
                successNotification: false,
            },
            {
                onSuccess: () => {
                    callbacks?.onSuccess?.();
                },
                onError: (error) => {
                    callbacks?.onError?.(error);
                },
            }
        );
    };

    return {
        changePassword,
        isPending: mutation.isPending,
        isError: mutation.isError,
        error: mutation.error,
    };
}
