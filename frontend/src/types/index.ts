/**
 * Types Module - Global DTOs and API contracts
 * 
 * This folder contains global type definitions shared across features.
 * Each subfolder/file represents one module (users, auth, etc.)
 * 
 * Includes:
 * - Request types
 * - Response types
 * - Enum types from backend
 * - Shared interfaces
 */

// Common API response types
export interface ResponseAPI<T> {
    data: T;
    error?: string;
    is_success: boolean;
    message: string;
}

export interface ListResponseAPI<T> {
    data: {
        items: T[];
        pagination: {
            page: number;
            limit: number;
            total: number;
            total_pages: number;
        };
    };
    error?: string;
    is_success: boolean;
    message: string;
}

// Re-export all types
export * from "./auth.type";
export * from "./user.type";
export * from "./modal.type";
export * from "./navigation.type";
export * from "./template.type";
export * from "./report.type";

