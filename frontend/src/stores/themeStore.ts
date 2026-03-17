import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

type ThemeMode = "light" | "dark";

interface ThemeStore {
    mode: ThemeMode;
    setMode: (mode: ThemeMode) => void;
    toggle: () => void;
}

export const useThemeStore = create<ThemeStore>()(
    persist(
        (set, get) => ({
            mode: "light",

            setMode: (mode) => set({ mode }),

            toggle: () =>
                set({ mode: get().mode === "light" ? "dark" : "light" }),
        }),
        {
            name: "theme-storage",
            storage: createJSONStorage(() => localStorage),
        }
    )
);
