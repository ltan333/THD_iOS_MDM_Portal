import type { AuthProvider } from "@refinedev/core";
import { cookies } from "next/headers";
import type { UserRole } from "@/types";

/**
 * Server-side auth provider for Next.js App Router
 * 
 * Backend đặt HTTP-only cookies sau khi VerifyOTP:
 * - rcycles_token: Access token JWT
 * - rcycles_refresh: Refresh token
 * 
 * Vì cookies là HTTP-only, frontend không thể đọc value
 * nhưng có thể kiểm tra xem cookie có tồn tại không
 */

// Extended check result with role
export type AuthCheckResult = {
  authenticated: boolean;
  role?: UserRole;
  redirectTo?: string;
  logout?: boolean;
};

// Type for the auth provider that includes our extended check method
type ExtendedAuthProvider = {
  check: () => Promise<AuthCheckResult>;
};

export const authProviderServer: ExtendedAuthProvider = {
  check: async (): Promise<AuthCheckResult> => {
    const cookieStore = await cookies();

    // Kiểm tra xem cookie rcycles_token có tồn tại không
    // Cookie name được config trong backend: COOKIE_NAME (default: "rcycles_token")
    const authToken = cookieStore.get("rcycles_token");
    const refreshToken = cookieStore.get("rcycles_refresh");

    // Nếu có ít nhất 1 trong 2 cookie thì coi như authenticated
    // Backend middleware sẽ tự động refresh nếu access token hết hạn
    if (authToken || refreshToken) {
      // Fetch user role from API
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
        const response = await fetch(`${apiUrl}/users/me`, {
          headers: {
            Cookie: cookieStore.toString(),
          },
          cache: "no-store",
        });
        console.log(response);
        if (response.ok) {
          const data = await response.json();
          return {
            authenticated: true,
            role: data.data?.role as UserRole,
          };
        }
      } catch {
        // If API call fails, still return authenticated but without role
        // This allows the user to access public-authenticated routes
      }

      return {
        authenticated: true,
      };
    }

    return {
      authenticated: false,
      logout: true,
      redirectTo: "/login",
    };
  },
};
