import { create } from "zustand";
import type {
    ResourceSheetState,
    OpenResourceSheetPayload,
} from "@/types/modal.type";

// ==================== Store State Interface ====================

interface SheetStore {
    // Resource Sheet (Drawer)
    resourceSheet: ResourceSheetState;
    openSheet: (payload: OpenResourceSheetPayload) => void;
    closeSheet: () => void;

    // Helper to update sheet state
    updateSheetState: (updates: Partial<ResourceSheetState>) => void;
}

// ==================== Initial State ====================

const initialResourceSheetState: ResourceSheetState = {
    isOpen: false,
    type: "create",
    resource: "",
    recordId: undefined,
    defaultValues: undefined,
    title: undefined,
    width: 600,
    onSuccess: undefined,
    onClose: undefined,
};

// ==================== Store Implementation ====================

export const useSheetStore = create<SheetStore>((set, get) => ({
    resourceSheet: initialResourceSheetState,

    openSheet: (payload) =>
        set({
            resourceSheet: {
                isOpen: true,
                type: payload.type,
                resource: payload.resource,
                recordId: payload.recordId,
                defaultValues: payload.defaultValues,
                title: payload.title,
                width: payload.width ?? 600,
                onSuccess: payload.onSuccess,
                onClose: payload.onClose,
            },
        }),

    closeSheet: () => {
        const currentState = get().resourceSheet;
        currentState.onClose?.();
        set({
            resourceSheet: {
                ...initialResourceSheetState,
                isOpen: false,
            },
        });
    },

    updateSheetState: (updates) =>
        set({
            resourceSheet: {
                ...get().resourceSheet,
                ...updates,
            },
        }),
}));
