import { UserRole } from "@/types";
import { redirect } from "next/navigation";

/**
 * Result type from auth provider check
 */
export type AuthCheckResult = {
    authenticated: boolean;
    role?: UserRole;
    redirectTo?: string;
};

/**
 * Roles that have admin access (can access all routes)
 */
const ADMIN_ROLES: UserRole[] = ["ADMIN", "SYSTEM_ADMIN"];

/**
 * Get default redirect path based on user role
 */
export function getRedirectPath(role?: UserRole): string {
    if (!role) return "/login";
    
    // Admin users go to dashboard
    if (ADMIN_ROLES.includes(role)) {
        return "/dashboard";
    }
    
    // Regular users go to their profile
    return "/profile";
}

/**
 * Check if user has admin access
 * Redirects to /profile if user is authenticated but not admin
 * Redirects to /login if not authenticated
 */
export function checkAdminAccess(auth: AuthCheckResult): void {
    if (!auth.authenticated) {
        redirect(auth.redirectTo || "/login");
    }

    // If role is USER (not ADMIN or SYSTEM_ADMIN), redirect to user dashboard
    if (auth.role && !ADMIN_ROLES.includes(auth.role)) {
        redirect("/profile");
    }
}

/**
 * Check if user has basic user access
 * All authenticated users can access /user/* routes (including ADMIN/SYSTEM_ADMIN)
 * Redirects to /login if not authenticated
 */
export function checkUserAccess(auth: AuthCheckResult): void {
    if (!auth.authenticated) {
        redirect(auth.redirectTo || "/login");
    }
    // All authenticated users can access /user/* routes
}

/**
 * Check if a role is an admin role
 */
export function isAdminRole(role?: UserRole): boolean {
    return role ? ADMIN_ROLES.includes(role) : false;
}
