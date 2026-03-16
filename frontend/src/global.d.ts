/**
 * Global Type Definitions
 * 
 * This file contains global types that are available throughout the application
 * without needing to import them. These types are primarily for API communication.
 * 
 * For domain-specific types (UserRole, Status, Session, etc.), 
 * see src/types/common.type.ts
 */

declare global {
  /**
   * Base API Response wrapper
   * Matches: rcycles-backend_pkg_response.APIResponse
   */
  type BaseResponseAPI<T> = {
    data: T;
    error?: string;
    is_success: boolean;
    message: string;
  };

  /**
   * Single item response
   * Use this for most API responses
   */
  type ResponseAPI<T> = BaseResponseAPI<T>;

  /**
   * Pagination object
   * Matches: rcycles-backend_internal_dto.Pagination
   */
  type Pagination = {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };

  /**
   * List response with pagination
   * All list endpoints return { items: T[], pagination: Pagination }
   */
  type ListResponseAPI<T> = BaseResponseAPI<{
    items: T[];
    pagination: Pagination;
  }>;

  /**
   * Common list request params for GET endpoints
   */
  type ListPaginationRequest = {
    page?: number;
    limit?: number;
    search?: string;
    sort?: string;
  };
}

export { };
