/**
 * Modal Types - Modal and Sheet state types
 */

import type { ReactNode } from "react";

// Modal sizes
export type ModalSize = "sm" | "md" | "lg" | "xl" | "full";

// Sheet positions
export type SheetPosition = "left" | "right" | "top" | "bottom";

// Base modal config
export interface ModalConfig {
    title?: string;
    content?: ReactNode;
    size?: ModalSize;
    closable?: boolean;
    maskClosable?: boolean;
    onClose?: () => void;
    onConfirm?: () => void | Promise<void>;
    confirmText?: string;
    cancelText?: string;
    showFooter?: boolean;
}

// Sheet config
export interface SheetConfig {
    title?: string;
    content?: ReactNode;
    position?: SheetPosition;
    size?: number | string;
    closable?: boolean;
    maskClosable?: boolean;
    onClose?: () => void;
}

// Confirm dialog config
export interface ConfirmConfig {
    title: string;
    description?: string;
    confirmText?: string;
    cancelText?: string;
    type?: "info" | "warning" | "danger";
    onConfirm?: () => void | Promise<void>;
    onCancel?: () => void;
}

// ============ Confirm Delete Modal ============
export interface ConfirmDeleteModalState {
    isOpen: boolean;
    title: string;
    description: string;
    isLoading: boolean;
    onConfirm: () => void | Promise<void>;
    onCancel?: () => void;
}

export interface OpenConfirmDeletePayload {
    title: string;
    description: string;
    onConfirm: () => void | Promise<void>;
    onCancel?: () => void;
}

// ============ Generic Modal ============
export interface GenericModalState {
    isOpen: boolean;
    title?: string;
    content: ReactNode | null;
    footer?: ReactNode;
    width: number;
    maskClosable: boolean;
    onClose?: () => void;
}

export interface OpenGenericModalPayload {
    title?: string;
    content: ReactNode;
    footer?: ReactNode;
    width?: number;
    maskClosable?: boolean;
    onClose?: () => void;
}

// ============ Sheet ============
export interface SheetState {
    isOpen: boolean;
    title?: string;
    content: ReactNode | null;
    position: SheetPosition;
    size?: number | string;
    maskClosable: boolean;
    onClose?: () => void;
}

export interface OpenSheetPayload {
    title?: string;
    content: ReactNode;
    position?: SheetPosition;
    size?: number | string;
    maskClosable?: boolean;
    onClose?: () => void;
}

// ============ Resource Sheet (for CRUD operations) ============
export type ResourceSheetType = "create" | "edit" | "view";

export interface ResourceSheetState {
    isOpen: boolean;
    type: ResourceSheetType;
    resource: string;
    recordId?: number | string;
    defaultValues?: Record<string, unknown>;
    title?: string;
    width: number;
    onSuccess?: () => void;
    onClose?: () => void;
}

export interface OpenResourceSheetPayload {
    type: ResourceSheetType;
    resource: string;
    recordId?: number | string;
    defaultValues?: Record<string, unknown>;
    title?: string;
    width?: number;
    onSuccess?: () => void;
    onClose?: () => void;
}

