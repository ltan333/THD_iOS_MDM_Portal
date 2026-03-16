"use client";

import { useEffect, useRef } from "react";
import { useGetIdentity } from "@refinedev/core";
import { useAuthStore } from "@/stores/authStore";
import type { UserResponse } from "@/types";

/**
 * Hook để sync dữ liệu user từ useGetIdentity vào Zustand store
 * 
 * Sử dụng trong các layout/page chính sau khi authenticated:
 * - Dashboard
 * - Main layouts
 * 
 * Flow:
 * 1. Component mount → gọi useGetIdentity (fetch /users/me)
 * 2. Khi có data → sync vào Zustand store
 * 3. Khi refresh page → Zustand persist từ localStorage
 * 4. Nếu localStorage rỗng → fetch lại từ API
 */
export function useAuthSync() {
    const { data: identity, isLoading, isError, error, refetch } = useGetIdentity<UserResponse>();
    const { user, setUser, setLoading, logout } = useAuthStore();
    
    // Track if logout was already called to prevent multiple calls
    const hasLoggedOut = useRef(false);

    // Check if error is a server error (500) vs auth error (401)
    // Support multiple error formats
    const err = error as any;
    const errorStatus =
        err?.statusCode ||
        err?.status ||
        err?.response?.status;
    const isServerError = isError && errorStatus && errorStatus >= 500;
    const isAuthError = isError && (!errorStatus || errorStatus === 401 || errorStatus === 403);

    useEffect(() => {
        setLoading(isLoading);
    }, [isLoading, setLoading]);

    useEffect(() => {
        // Khi có identity data từ API → sync vào store
        if (identity && !isLoading) {
            setUser(identity);
            // Reset logout flag when successfully authenticated
            hasLoggedOut.current = false;
        }
    }, [identity, isLoading, setUser]);

    useEffect(() => {
        // Chỉ logout khi lỗi xác thực (401, etc.), không logout khi lỗi server (500)
        // Và chỉ logout 1 lần để tránh vòng lặp
        if (isAuthError && !isLoading && !hasLoggedOut.current) {
            hasLoggedOut.current = true;
            logout();
        }
    }, [isAuthError, isLoading, logout]);

    return {
        user: user || identity,
        isLoading,
        isError,
        isServerError, // Để hiển thị error UI cho lỗi 500
        error,
        refetch,
        refetchUser: refetch, // Alias for semantic usage in profile updates
        isAuthenticated: !!user || !!identity,
    };
}

/**
 * Hook đơn giản để lấy user từ store (không fetch API)
 * Dùng khi chắc chắn user đã được sync rồi
 */
export function useCurrentUser() {
    const user = useAuthStore((state) => state.user);
    const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
    return { user, isAuthenticated };
}
