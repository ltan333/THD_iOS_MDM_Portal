"use client";

import type { DataProvider, CrudFilter, LogicalFilter } from "@refinedev/core";
import { get, post, put, del, patch } from "@/axios-config";

/**
 * Custom Data Provider for R-Cycles Backend API
 * Uses axios-config methods for consistent API handling
 * 
 * API Response format (from swagger):
 * - Single: { data: T, error?: string, is_success: boolean, message: string }
 * - List: { data: { items: T[], pagination: { page, limit, total, total_pages } }, ... }
 * 
 * Note: We return the API message in the response so Refine can use it for notifications.
 * Components should configure successNotification to use this message.
 */

// Type guard for LogicalFilter
function isLogicalFilter(filter: CrudFilter): filter is LogicalFilter {
    return "field" in filter;
}

// Map Refine operators to backend operators
const OPERATOR_MAP: Record<string, string> = {
    eq: "eq",
    ne: "ne",
    lt: "lt",
    lte: "lte",
    gt: "gt",
    gte: "gte",
    contains: "like",
    containss: "like",
    ncontains: "ne",
    in: "in",
    nin: "ne",
};

// Helper to build query params for list requests
// Backend format: filter[field:operator]=value, sort=-field1,field2
function buildListParams(params: {
    pagination?: { current?: number; pageSize?: number };
    sorters?: Array<{ field: string; order: "asc" | "desc" }>;
    filters?: CrudFilter[];
}): Record<string, string | number | boolean | undefined> {
    const queryParams: Record<string, string | number | boolean | undefined> = {};

    // Pagination
    if (params.pagination) {
        queryParams.page = params.pagination.current || 1;
        queryParams.limit = params.pagination.pageSize || 20;
    }

    // Sorting - format as "-field1,field2" (- prefix means DESC)
    // Default to "id" if no sorters provided
    if (params.sorters && params.sorters.length > 0) {
        const sortFields = params.sorters.map((sorter) => {
            return sorter.order === "desc" ? `-${sorter.field}` : sorter.field;
        });
        queryParams.sort = sortFields.join(",");
    } else {
        // Default sort by id ascending
        queryParams.sort = "id";
    }

    // Filters - convert to filter[field:operator]=value format
    if (params.filters && params.filters.length > 0) {
        params.filters.forEach((filter) => {
            if (!isLogicalFilter(filter)) return;

            const { field, operator, value } = filter;
            if (value === undefined || value === null || value === "") return;

            const backendOperator = OPERATOR_MAP[operator] || "eq";

            // Build filter key: filter[field] for eq, filter[field:operator] for others
            const filterKey =
                backendOperator === "eq"
                    ? `filter[${field}]`
                    : `filter[${field}:${backendOperator}]`;

            // Handle value formatting
            if (operator === "in" && Array.isArray(value)) {
                queryParams[filterKey] = value.join(",");
            } else {
                queryParams[filterKey] = String(value);
            }
        });
    }

    return queryParams;
}

// Custom data provider implementation
const customDataProvider = {
    getList: async ({ resource, pagination, sorters, filters, meta }: {
        resource: string;
        pagination?: { current?: number; pageSize?: number };
        sorters?: Array<{ field: string; order: "asc" | "desc" }>;
        filters?: CrudFilter[];
        meta?: Record<string, unknown>;
    }) => {
        const queryParams = buildListParams({ pagination, sorters, filters });

        // Merge additional query params from meta
        if (meta?.queryParams) {
            Object.assign(queryParams, meta.queryParams);
        }

        const response = await get<ListResponseAPI<Record<string, unknown>>>(`/${resource}`, { queryParams });

        return {
            data: response.data.items,
            total: response.data.pagination?.total || response.data.items.length,
        };
    },

    getOne: async ({ resource, id, meta }: {
        resource: string;
        id: string | number;
        meta?: Record<string, unknown>;
    }) => {
        const queryParams = (meta?.queryParams || {}) as Record<string, string | number | boolean | undefined>;
        const response = await get<ResponseAPI<Record<string, unknown>>>(`/${resource}/${id}`, { queryParams });

        return { data: response.data };
    },

    create: async ({ resource, variables, meta }: {
        resource: string;
        variables: Record<string, unknown>;
        meta?: Record<string, unknown>;
    }) => {
        const options = { useFormData: (meta?.useFormData as boolean) || false };
        const response = await post<ResponseAPI<Record<string, unknown>>, typeof variables>(`/${resource}`, variables, options);

        return {
            data: response.data,
            // Include API message for notification
            meta: { message: response.message },
        };
    },

    update: async ({ resource, id, variables, meta }: {
        resource: string;
        id: string | number;
        variables: Record<string, unknown>;
        meta?: Record<string, unknown>;
    }) => {
        const options = { useFormData: (meta?.useFormData as boolean) || false };
        const method = (meta?.method as string) || "put";

        const response = method === "patch"
            ? await patch<ResponseAPI<Record<string, unknown>>, typeof variables>(`/${resource}/${id}`, variables, options)
            : await put<ResponseAPI<Record<string, unknown>>, typeof variables>(`/${resource}/${id}`, variables, options);

        return {
            data: response.data,
            // Include API message for notification
            meta: { message: response.message },
        };
    },

    deleteOne: async ({ resource, id }: {
        resource: string;
        id: string | number;
    }) => {
        const response = await del<ResponseAPI<Record<string, unknown>>>(`/${resource}/${id}`);

        return {
            data: response.data,
            // Include API message for notification
            meta: { message: response.message },
        };
    },

    getMany: async ({ resource, ids, meta }: {
        resource: string;
        ids: (string | number)[];
        meta?: Record<string, unknown>;
    }) => {
        // Fetch items one by one since API doesn't support bulk get
        const promises = ids.map((id) =>
            get<ResponseAPI<Record<string, unknown>>>(`/${resource}/${id}`)
        );

        const results = await Promise.all(promises);

        return { data: results.map((r) => r.data) };
    },

    deleteMany: async ({ resource, ids }: {
        resource: string;
        ids: (string | number)[];
    }) => {
        const promises = ids.map((id) =>
            del<ResponseAPI<Record<string, unknown>>>(`/${resource}/${id}`)
        );

        const results = await Promise.all(promises);

        return { data: results.map((r) => r.data) };
    },

    updateMany: async ({ resource, ids, variables }: {
        resource: string;
        ids: (string | number)[];
        variables: Record<string, unknown>;
    }) => {
        const promises = ids.map((id) =>
            put<ResponseAPI<Record<string, unknown>>, typeof variables>(`/${resource}/${id}`, variables)
        );

        const results = await Promise.all(promises);

        return { data: results.map((r) => r.data) };
    },

    createMany: async ({ resource, variables }: {
        resource: string;
        variables: Record<string, unknown>[];
    }) => {
        const promises = variables.map((item) =>
            post<ResponseAPI<Record<string, unknown>>, typeof item>(`/${resource}`, item)
        );

        const results = await Promise.all(promises);

        return { data: results.map((r) => r.data) };
    },

    custom: async ({ url, method, payload, query, headers }: {
        url: string;
        method: string;
        payload?: unknown;
        query?: Record<string, unknown>;
        headers?: Record<string, string>;
    }) => {
        const options = {
            queryParams: query as Record<string, string | number | boolean | undefined>,
            headers: headers as Record<string, string>,
        };

        let response: ResponseAPI<unknown>;

        switch (method) {
            case "get":
                response = await get<ResponseAPI<unknown>>(url, options);
                break;
            case "post":
                response = await post<ResponseAPI<unknown>, typeof payload>(url, payload, options);
                break;
            case "put":
                response = await put<ResponseAPI<unknown>, typeof payload>(url, payload, options);
                break;
            case "patch":
                response = await patch<ResponseAPI<unknown>, typeof payload>(url, payload, options);
                break;
            case "delete":
                response = await del<ResponseAPI<unknown>>(url, options);
                break;
            default:
                response = await get<ResponseAPI<unknown>>(url, options);
        }

        return { data: response.data };
    },

    getApiUrl: () => process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000/api/v1",
};

// Cast to DataProvider to satisfy Refine's type requirements
export const dataProvider = customDataProvider as unknown as DataProvider;
