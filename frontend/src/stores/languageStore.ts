import { create } from 'zustand';
import { persist } from 'zustand/middleware';

type Language = 'vi' | 'en';

interface LanguageState {
  language: Language;
  setLanguage: (lang: Language) => void;
  toggleLanguage: () => void;
}

export const useLanguageStore = create<LanguageState>()(
  persist(
    (set) => ({
      language: 'vi',
      setLanguage: (lang) => set({ language: lang }),
      toggleLanguage: () => set((state) => ({ language: state.language === 'vi' ? 'en' : 'vi' })),
    }),
    {
      name: 'language-storage',
    }
  )
);
