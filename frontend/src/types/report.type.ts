/**
 * Report Types - Work report related types
 */

import type { UserResponse } from "./user.type";
import type { QuestionResponse } from "./template.type";

export type ReportSession = "R1" | "R2";

export interface ReportAnswerResponse {
    id: number;
    report_id: number;
    question_id: number;
    question?: QuestionResponse;
    answer_text: string;
    created_at: string;
}

export interface ReportResponse {
    id: number;
    user_id: number;
    user?: UserResponse;
    template_id: number;
    project_id?: number;
    project?: {
        id: number;
        name: string;
    };
    department_id?: number;
    department?: {
        id: number;
        name: string;
    };
    session: ReportSession;
    report_date: string;
    completion_rate: number;
    notes?: string;
    answers?: ReportAnswerResponse[];
    created_at: string;
    updated_at?: string;
}

export interface ReportCreateRequest {
    template_id: number;
    project_id?: number;
    session: ReportSession;
    report_date: string;
    completion_rate: number;
    notes?: string;
    answers: Array<{
        question_id: number;
        answer_text: string;
    }>;
}

export interface ReportUpdateRequest {
    completion_rate?: number;
    notes?: string;
    answers?: Array<{
        question_id: number;
        answer_text: string;
    }>;
}

export interface ReportFilterParams {
    project_id?: number;
    department_id?: number;
    session?: ReportSession;
    user_id?: number;
    start_date?: string;
    end_date?: string;
}
