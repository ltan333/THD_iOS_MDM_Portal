'use client';

import { usePathname } from 'next/navigation';
import { Bell, Search, Menu } from 'lucide-react';
import { Button } from '@heroui/button';
import { useSidebarStore } from '@/features/sidebar/stores/useSidebarStore';
import { useLanguageStore } from '@/stores/languageStore';
import { useState, useEffect } from 'react';

export function Header() {
  const pathname = usePathname();
  const { toggleSidebar } = useSidebarStore();
  const { language, toggleLanguage } = useLanguageStore();
  const [mounted, setMounted] = useState(false);

  // Prevent hydration mismatch
  useEffect(() => {
    setMounted(true);
  }, []);

  const t = {
    vi: {
      dashboard: 'Bảng điều khiển',
      profile: 'Cấu hình',
      devices: 'Thiết bị',
      users: 'Người dùng',
      applications: 'Ứng dụng',
      policies: 'Chính sách',
      reports: 'Báo cáo',
      alerts: 'Cảnh báo',
      settings: 'Cài đặt',
      search: 'Tìm kiếm...'
    },
    en: {
      dashboard: 'Dashboard',
      profile: 'Profile',
      devices: 'Devices',
      users: 'Users',
      applications: 'Applications',
      policies: 'Policies',
      reports: 'Reports',
      alerts: 'Alerts',
      settings: 'Settings',
      search: 'Search...'
    }
  };

  const currentLang = mounted ? language : 'vi';

  // Create a readable title from pathname
  const getTitle = () => {
    if (pathname === '/') return t[currentLang].dashboard;
    
    const segment = pathname.split('/')[1];
    if (!segment) return t[currentLang].dashboard;

    // Try to translate the segment if it exists in our dictionary
    const key = segment as keyof typeof t.en;
    if (t[currentLang][key]) {
      return t[currentLang][key];
    }

    // Fallback: capitalize first letter
    return segment.charAt(0).toUpperCase() + segment.slice(1);
  };

  return (
    <header className="sticky top-0 z-30 flex items-center justify-between h-16 px-4 sm:px-6 lg:px-8 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 transition-colors">
      <div className="flex items-center gap-4">
        {/* Mobile menu toggle */}
        <Button
          isIconOnly
          variant="light"
          className="lg:hidden text-slate-500 dark:text-slate-400"
          onPress={toggleSidebar}
          aria-label="Open sidebar"
        >
          <Menu className="w-5 h-5" />
        </Button>

        {/* Page Title / Breadcrumbs */}
        <div className="flex items-center">
          <h1 className="text-lg font-semibold text-slate-900 dark:text-white truncate max-w-[200px] sm:max-w-md transition-opacity duration-300" style={{ opacity: mounted ? 1 : 0 }}>
            {getTitle()}
          </h1>
        </div>
      </div>

      {/* Right side actions */}
      <div className="flex items-center gap-2 sm:gap-4">
        {/* Search */}
        <div className="hidden md:flex relative group">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search className="h-4 w-4 text-slate-400 group-focus-within:text-indigo-500 transition-colors" />
          </div>
          <input
            type="text"
            placeholder={mounted ? t[language].search : t.vi.search}
            className="block w-full sm:w-64 pl-10 pr-3 py-2 border border-slate-200 dark:border-slate-700 rounded-lg text-sm bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
          />
        </div>

        <Button
          isIconOnly
          variant="light"
          className="md:hidden text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200"
          aria-label="Search"
        >
          <Search className="w-5 h-5" />
        </Button>

        {/* Language Toggle */}
        <div className="flex items-center gap-2 bg-slate-100 dark:bg-slate-800 px-3 py-1.5 rounded-full border border-slate-200 dark:border-slate-700 opacity-0 transition-opacity duration-300" style={{ opacity: mounted ? 1 : 0 }}>
          <span className={`text-[10px] font-bold tracking-wider transition-colors ${mounted && language === 'en' ? 'text-slate-800 dark:text-slate-200' : 'text-slate-400'}`}>
            EN
          </span>
          <button 
            onClick={toggleLanguage}
            className="relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500/50 bg-slate-300 dark:bg-slate-600 cursor-pointer"
            role="switch"
            aria-checked={mounted && language === 'vi'}
          >
            <span className="sr-only">Toggle Language</span>
            <span
              className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white transition-transform duration-300 ease-spring shadow-sm ${
                mounted && language === 'vi' ? 'translate-x-4' : 'translate-x-1'
              }`}
            />
          </button>
          <span className={`text-[10px] font-bold tracking-wider transition-colors ${mounted && language === 'vi' ? 'text-slate-800 dark:text-slate-200' : 'text-slate-400'}`}>
            VI
          </span>
        </div>

        {/* Notifications */}
        <div className="relative">
          <Button
            isIconOnly
            variant="light"
            className="text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200"
            aria-label="Notifications"
          >
            <Bell className="w-5 h-5" />
          </Button>
          <span className="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full border-2 border-white dark:border-slate-900"></span>
        </div>
      </div>
    </header>
  );
}
