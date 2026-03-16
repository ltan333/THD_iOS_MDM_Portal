import { useState, useEffect, useRef, useCallback } from "react";

/**
 * Debounce a value
 */
export function useDebounce<T>(value: T, delay: number = 300): T {
    const [debouncedValue, setDebouncedValue] = useState<T>(value);

    useEffect(() => {
        const timer = setTimeout(() => setDebouncedValue(value), delay);
        return () => clearTimeout(timer);
    }, [value, delay]);

    return debouncedValue;
}

/**
 * Modal/Drawer disclosure control
 */
export function useDisclosure(initialState = false) {
    const [isOpen, setIsOpen] = useState(initialState);

    const open = useCallback(() => setIsOpen(true), []);
    const close = useCallback(() => setIsOpen(false), []);
    const toggle = useCallback(() => setIsOpen((prev) => !prev), []);

    return { isOpen, open, close, toggle };
}

/**
 * Track if component is mounted
 */
export function useMounted() {
    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        setMounted(true);
    }, []);

    return mounted;
}

/**
 * Previous value
 */
export function usePrevious<T>(value: T): T | undefined {
    const ref = useRef<T | undefined>(undefined);

    useEffect(() => {
        ref.current = value;
    }, [value]);

    return ref.current;
}

// Auth sync hooks
export { useAuthSync, useCurrentUser } from "./useAuthSync";

// Report, Template & Question hooks
export {
    // Template hooks
    useActiveTemplate,
    useTemplates,
    useTemplate,
    // Question hooks
    useTemplateQuestions,
    // Report hooks
    useReports,
    useReport,
    useCreateReport,
    useUpdateReport,
    useDeleteReport,
    // Permission hooks
    useMyPermissions,
} from "./useReport";

// Profile hooks
export { useUpdateProfile, useChangePassword } from "./useProfile";
