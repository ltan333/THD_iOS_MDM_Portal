"use client";

import { Spin } from "antd";
import { motion } from "framer-motion";

/**
 * Full page loading spinner
 * Used for Suspense fallback and page transitions
 */
export function FullPageLoader() {
    return (
        <div className="fixed inset-0 flex items-center justify-center bg-white dark:bg-gray-900 z-50">
            <motion.div
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.3 }}
                className="flex flex-col items-center gap-4"
            >
                <Spin size="large" />
                <span className="text-gray-500 dark:text-gray-400 text-sm">
                    Đang tải...
                </span>
            </motion.div>
        </div>
    );
}

/**
 * Inline loader for smaller sections
 */
export function InlineLoader({ text = "Đang tải..." }: { text?: string }) {
    return (
        <div className="flex items-center justify-center p-8">
            <Spin />
            <span className="ml-2 text-gray-500 dark:text-gray-400">{text}</span>
        </div>
    );
}

/**
 * Card skeleton loader
 */
export function CardSkeleton() {
    return (
        <div className="animate-pulse">
            <div className="h-48 bg-gray-200 dark:bg-gray-700 rounded-lg" />
        </div>
    );
}

/**
 * Table row skeleton loader
 */
export function TableRowSkeleton({ columns = 5 }: { columns?: number }) {
    return (
        <div className="flex gap-4 py-4 animate-pulse">
            {Array.from({ length: columns }).map((_, i) => (
                <div
                    key={i}
                    className="h-4 bg-gray-200 dark:bg-gray-700 rounded flex-1"
                />
            ))}
        </div>
    );
}
