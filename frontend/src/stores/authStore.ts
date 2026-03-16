import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import type { UserResponse } from "@/types";

/**
 * Auth State - lưu trữ user data sau khi xác thực
 * 
 * Flow:
 * 1. Sau OTP verify → redirect đến dashboard
 * 2. Dashboard gọi useGetIdentity → lấy user data từ API
 * 3. Sync user data vào Zustand store
 * 4. Khi refresh page → Zustand có thể mất → useGetIdentity fetch lại
 */

interface AuthState {
    user: UserResponse | null;
    isAuthenticated: boolean;
    isLoading: boolean;
}

interface AuthStore extends AuthState {
    // Actions
    setUser: (user: UserResponse | null) => void;
    setAuthenticated: (isAuthenticated: boolean) => void;
    setLoading: (isLoading: boolean) => void;
    login: (user: UserResponse) => void;
    logout: () => void;
    reset: () => void;
}

const initialState: AuthState = {
    user: null,
    isAuthenticated: false,
    isLoading: true,
};

export const useAuthStore = create<AuthStore>()(
    persist(
        (set) => ({
            ...initialState,

            setUser: (user) => set({ user, isAuthenticated: !!user }),

            setAuthenticated: (isAuthenticated) => set({ isAuthenticated }),

            setLoading: (isLoading) => set({ isLoading }),

            login: (user) =>
                set({
                    user,
                    isAuthenticated: true,
                    isLoading: false,
                }),

            logout: () =>
                set({
                    ...initialState,
                    isLoading: false,
                }),

            reset: () => set(initialState),
        }),
        {
            name: "auth-storage",
            storage: createJSONStorage(() => localStorage),
            partialize: (state) => ({
                user: state.user,
                isAuthenticated: state.isAuthenticated,
            }),
        }
    )
);
