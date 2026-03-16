"use client";

import React, { createContext, useContext, useState, useCallback, useEffect } from "react";
import { usePathname, useSearchParams } from "next/navigation";

interface LoadingContextType {
    isLoading: boolean;
    startLoading: () => void;
    stopLoading: () => void;
}

const LoadingContext = createContext<LoadingContextType>({
    isLoading: false,
    startLoading: () => { },
    stopLoading: () => { },
});

export const useLoading = () => useContext(LoadingContext);

export function LoadingProvider({ children }: { children: React.ReactNode }) {
    const [isLoading, setIsLoading] = useState(false);
    const pathname = usePathname();
    const searchParams = useSearchParams();

    const startLoading = useCallback(() => setIsLoading(true), []);
    const stopLoading = useCallback(() => setIsLoading(false), []);

    // Stop loading when route changes complete
    useEffect(() => {
        setIsLoading(false);
    }, [pathname, searchParams]);

    return (
        <LoadingContext.Provider value={{ isLoading, startLoading, stopLoading }}>
            {children}

            {/* Page Transition Loader */}
            {isLoading && (
                <div className="fixed inset-0 z-9999 pointer-events-none">
                    {/* Top progress bar */}
                    <div className="absolute top-0 left-0 right-0 h-1 bg-primary/20 overflow-hidden">
                        <div className="h-full bg-primary animate-loading-bar" />
                    </div>

                    {/* Subtle overlay */}
                    <div className="absolute inset-0 bg-background/30 backdrop-blur-[1px] animate-fade-in" />
                </div>
            )}
        </LoadingContext.Provider>
    );
}
