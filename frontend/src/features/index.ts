/**
 * Features Module
 * 
 * This folder contains all feature-specific modules.
 * Each feature is a self-contained module with its own:
 * - components/  (feature-only UI components)
 * - hooks/       (feature-only hooks, including refine hooks)
 * - stores/      (feature-only Zustand stores)
 * - services/    (API / data logic)
 * - types/       (feature-specific types)
 * - pages/       (page-level components)
 * - index.ts     (public exports)
 * 
 * RULES:
 * - Feature A must NOT import internal files of Feature B
 * - Cross-feature access must go through index.ts
 * - All refine hooks (useList, useCreate, etc.) MUST live inside feature hooks
 */

export {};
