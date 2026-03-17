/**
 * Template Types - Report template related types
 */

export interface QuestionResponse {
    id: number;
    template_id: number;
    question_text: string;
    question_order: number;
    is_required: boolean;
    question_type: "text" | "number" | "textarea" | "select" | "checkbox";
    options?: string[];
    created_at: string;
    updated_at?: string;
}

export interface TemplateResponse {
    id: number;
    name: string;
    description?: string;
    is_active: boolean;
    questions?: QuestionResponse[];
    created_at: string;
    updated_at?: string;
}

export interface TemplateCreateRequest {
    name: string;
    description?: string;
    is_active?: boolean;
}

export interface TemplateUpdateRequest {
    name?: string;
    description?: string;
    is_active?: boolean;
}

export interface QuestionCreateRequest {
    template_id: number;
    question_text: string;
    question_order?: number;
    is_required?: boolean;
    question_type?: "text" | "number" | "textarea" | "select" | "checkbox";
    options?: string[];
}

export interface QuestionUpdateRequest {
    question_text?: string;
    question_order?: number;
    is_required?: boolean;
    question_type?: "text" | "number" | "textarea" | "select" | "checkbox";
    options?: string[];
}
