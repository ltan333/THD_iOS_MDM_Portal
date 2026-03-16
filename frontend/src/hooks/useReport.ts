"use client";

import { useState, useEffect, useCallback } from "react";
import { useList, useOne, useCreate, useUpdate, useDelete, CrudFilters } from "@refinedev/core";
import { get } from "@/axios-config";
import { REPORT_TEMPLATE_ENDPOINTS, QUESTION_ENDPOINTS } from "@/axios-config/constants";
import type { TemplateResponse, QuestionResponse } from "@/types/template.type";
import type { ReportResponse } from "@/types/report.type";

// ============ Template Hooks ============

interface UseQueryResult<T> {
    data: T | undefined;
    isLoading: boolean;
    error: Error | null;
    refetch: () => void;
}

/**
 * Get the currently active template with questions
 * Uses custom endpoint: GET /templates/active
 */
export function useActiveTemplate(): UseQueryResult<TemplateResponse> {
    const [data, setData] = useState<TemplateResponse | undefined>();
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchData = useCallback(async () => {
        try {
            setIsLoading(true);
            setError(null);

            const response = await get<ResponseAPI<TemplateResponse>>(REPORT_TEMPLATE_ENDPOINTS.ACTIVE);
            setData(response.data);
        } catch (err) {
            setError(err instanceof Error ? err : new Error("Failed to fetch active template"));
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchData();
    }, [fetchData]);

    return { data, isLoading, error, refetch: fetchData };
}

/**
 * Get list of templates with pagination
 * Uses: GET /templates
 */
export function useTemplates(options?: {
    pagination?: { current?: number; pageSize?: number };
}) {
    return useList<TemplateResponse>({
        resource: "templates",
        pagination: options?.pagination,
    });
}

/**
 * Get a template by ID
 * Uses: GET /templates/:id
 */
export function useTemplate(id: number | undefined) {
    return useOne<TemplateResponse>({
        resource: "templates",
        id: id as number,
        queryOptions: {
            enabled: !!id,
        },
    });
}

// ============ Question Hooks ============

/**
 * Get questions for a template
 * Uses custom endpoint: GET /templates/:id/questions
 */
export function useTemplateQuestions(templateId: number | undefined): UseQueryResult<QuestionResponse[]> {
    const [data, setData] = useState<QuestionResponse[] | undefined>();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);

    const fetchData = useCallback(async () => {
        if (!templateId) return;

        try {
            setIsLoading(true);
            setError(null);

            const response = await get<ResponseAPI<QuestionResponse[]>>(
                QUESTION_ENDPOINTS.BY_TEMPLATE(templateId)
            );
            setData(response.data);
        } catch (err) {
            setError(err instanceof Error ? err : new Error("Failed to fetch questions"));
        } finally {
            setIsLoading(false);
        }
    }, [templateId]);

    useEffect(() => {
        if (templateId) {
            fetchData();
        }
    }, [templateId, fetchData]);

    return { data, isLoading, error, refetch: fetchData };
}

// ============ Report Hooks ============

/**
 * Get list of reports with filters
 * Uses: GET /reports
 * 
 * Available filters (from swagger):
 * - project_id: Filter by project ID
 * - department_id: Filter by department ID
 * - session: Filter by session (R1, R2)
 * - user_id: Filter by user ID
 * - start_date: Filter by start date
 * - end_date: Filter by end date
 */
export function useReports(options?: {
    pagination?: { current?: number; pageSize?: number };
    filters?: CrudFilters;
    projectId?: number;
    departmentId?: number;
    session?: "R1" | "R2";
    userId?: number;
    startDate?: string;
    endDate?: string;
}) {
    const filters: CrudFilters = options?.filters || [];

    // Add custom filters
    if (options?.projectId) {
        filters.push({ field: "project_id", operator: "eq", value: options.projectId });
    }
    if (options?.departmentId) {
        filters.push({ field: "department_id", operator: "eq", value: options.departmentId });
    }
    if (options?.session) {
        filters.push({ field: "session", operator: "eq", value: options.session });
    }
    if (options?.userId) {
        filters.push({ field: "user_id", operator: "eq", value: options.userId });
    }
    if (options?.startDate) {
        filters.push({ field: "start_date", operator: "eq", value: options.startDate });
    }
    if (options?.endDate) {
        filters.push({ field: "end_date", operator: "eq", value: options.endDate });
    }

    return useList<ReportResponse>({
        resource: "reports",
        pagination: options?.pagination,
        filters,
    });
}

/**
 * Get a report by ID
 * Uses: GET /reports/:id
 */
export function useReport(id: number | undefined) {
    return useOne<ReportResponse>({
        resource: "reports",
        id: id as number,
        queryOptions: {
            enabled: !!id,
        },
    });
}

/**
 * Create a new report
 * Uses: POST /reports
 */
export function useCreateReport() {
    return useCreate<ReportResponse>();
}

/**
 * Update an existing report
 * Uses: PUT /reports/:id
 */
export function useUpdateReport() {
    return useUpdate<ReportResponse>();
}

/**
 * Delete a report
 * Uses: DELETE /reports/:id
 */
export function useDeleteReport() {
    return useDelete<ReportResponse>();
}

// ============ Permission Hooks ============

export interface UserPermissionsResponse {
    user_id: number;
    projects: Array<{ id: number; name: string }>;
    departments: Array<{ id: number; name: string }>;
}

/**
 * Get current user's permissions (accessible projects & departments)
 * Uses: GET /permissions/me
 */
export function useMyPermissions(): UseQueryResult<UserPermissionsResponse> {
    const [data, setData] = useState<UserPermissionsResponse | undefined>();
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchData = useCallback(async () => {
        try {
            setIsLoading(true);
            setError(null);

            const response = await get<ResponseAPI<UserPermissionsResponse>>("/permissions/me");
            setData(response.data);
        } catch (err) {
            setError(err instanceof Error ? err : new Error("Failed to fetch permissions"));
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchData();
    }, [fetchData]);

    return { data, isLoading, error, refetch: fetchData };
}
