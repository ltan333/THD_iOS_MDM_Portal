import { create } from "zustand";
import type {
    ConfirmDeleteModalState,
    OpenConfirmDeletePayload,
    GenericModalState,
    OpenGenericModalPayload,
} from "@/types/modal.type";

// ==================== Store State Interface ====================

interface ModalStore {
    // Confirm Delete Modal
    confirmDelete: ConfirmDeleteModalState;
    openConfirmDelete: (payload: OpenConfirmDeletePayload) => void;
    closeConfirmDelete: () => void;
    setConfirmDeleteLoading: (isLoading: boolean) => void;

    // Generic Modal
    genericModal: GenericModalState;
    openGenericModal: (payload: OpenGenericModalPayload) => void;
    closeGenericModal: () => void;

    // Reset all
    resetAll: () => void;
}

// ==================== Initial States ====================

const initialConfirmDeleteState: ConfirmDeleteModalState = {
    isOpen: false,
    title: "",
    description: "",
    isLoading: false,
    onConfirm: () => { },
    onCancel: undefined,
};

const initialGenericModalState: GenericModalState = {
    isOpen: false,
    title: undefined,
    content: null,
    footer: undefined,
    width: 520,
    maskClosable: true,
    onClose: undefined,
};

// ==================== Store Implementation ====================

export const useModalStore = create<ModalStore>((set, get) => ({
    // Confirm Delete Modal State & Actions
    confirmDelete: initialConfirmDeleteState,

    openConfirmDelete: (payload) =>
        set({
            confirmDelete: {
                isOpen: true,
                title: payload.title,
                description: payload.description,
                isLoading: false,
                onConfirm: payload.onConfirm,
                onCancel: payload.onCancel,
            },
        }),

    closeConfirmDelete: () =>
        set({
            confirmDelete: {
                ...get().confirmDelete,
                isOpen: false,
            },
        }),

    setConfirmDeleteLoading: (isLoading) =>
        set({
            confirmDelete: {
                ...get().confirmDelete,
                isLoading,
            },
        }),

    // Generic Modal State & Actions
    genericModal: initialGenericModalState,

    openGenericModal: (payload) =>
        set({
            genericModal: {
                isOpen: true,
                title: payload.title,
                content: payload.content,
                footer: payload.footer,
                width: payload.width ?? 520,
                maskClosable: payload.maskClosable ?? true,
                onClose: payload.onClose,
            },
        }),

    closeGenericModal: () => {
        const currentState = get().genericModal;
        currentState.onClose?.();
        set({
            genericModal: {
                ...initialGenericModalState,
                isOpen: false,
            },
        });
    },

    // Reset All
    resetAll: () =>
        set({
            confirmDelete: initialConfirmDeleteState,
            genericModal: initialGenericModalState,
        }),
}));
