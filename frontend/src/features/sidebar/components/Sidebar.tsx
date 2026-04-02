'use client';

import { usePathname } from 'next/navigation';
import { cn } from '@/utils/cn';
import { useSidebarStore } from '../stores/useSidebarStore';
import { useLanguageStore } from '@/stores/languageStore';
import { tokenManager } from '@/axios-config/utils/token-manager';
import { authProviderClient } from '@/providers/auth-provider/auth-provider.client';
import { useEffect, useState } from 'react';
import {
  LayoutDashboard,
  Smartphone,
  Users,
  AppWindow,
  Shield,
  BarChart3,
  Bell,
  Settings,
  ChevronLeft,
  ChevronRight,
  LogOut,
  User
} from 'lucide-react';
import { SidebarItem } from './SidebarItem';
import { SidebarSection } from './SidebarSection';
import { Button } from '@heroui/button';

export function Sidebar() {
  const { isCollapsed, toggleSidebar } = useSidebarStore();
  const { language } = useLanguageStore();
  const [mounted, setMounted] = useState(false);
  const pathname = usePathname();

  useEffect(() => {
    setMounted(true);
  }, []);

  const isRouteActive = (path: string) => pathname === path || (pathname || '').startsWith(path + '/');

  const t = {
    vi: {
      dashboard: 'Bảng điều khiển',
      profile: 'Cấu hình',
      devices: 'Thiết bị',
      allDevices: 'Tất cả thiết bị',
      deviceGroups: 'Nhóm thiết bị',
      users: 'Người dùng',
      applications: 'Ứng dụng',
      policies: 'Chính sách',
      reports: 'Báo cáo',
      alerts: 'Cảnh báo',
      settings: 'Cài đặt',
      collapseMenu: 'Thu gọn Menu'
    },
    en: {
      dashboard: 'Dashboard',
      profile: 'Profile',
      devices: 'Devices',
      allDevices: 'All Devices',
      deviceGroups: 'Device Groups',
      users: 'Users',
      applications: 'Applications',
      policies: 'Policies',
      reports: 'Reports',
      alerts: 'Alerts',
      settings: 'Settings',
      collapseMenu: 'Collapse Menu'
    }
  };

  const currentLang = mounted ? language : 'vi';

  return (
    <>
      {/* Mobile Backdrop */}
      {!isCollapsed && (
        <div
          className="fixed inset-0 bg-black/50 z-30 lg:hidden transition-opacity"
          onClick={() => toggleSidebar()}
          aria-hidden="true"
        />
      )}

      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-40 flex flex-col',
          'bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800',
          'shadow-sidebar',
          'transition-all duration-300 ease-in-out',
          isCollapsed ? '-translate-x-full lg:translate-x-0 w-[260px] lg:w-[72px]' : 'translate-x-0 w-[260px]'
        )}
        aria-label="Sidebar Navigation"
      >
      {/* Header / Logo */}
      <div className="flex items-center justify-between h-16 px-4 border-b border-gray-200 dark:border-gray-800 shrink-0 bg-gradient-to-r from-primary-600 to-primary-700 dark:from-primary-700 dark:to-primary-800">
        <div className={cn('flex items-center', isCollapsed && 'justify-center w-full')}>
          <div className="w-9 h-9 rounded-lg bg-white/10 backdrop-blur-sm flex items-center justify-center shrink-0 border border-white/20">
            <Shield className="w-5 h-5 text-white" />
          </div>
          {!isCollapsed && (
            <span className="ml-3 font-semibold text-white truncate text-base">
              MDM Portal
            </span>
          )}
        </div>
      </div>

      {/* Main Navigation */}
      <nav className="flex-1 overflow-y-auto py-4 px-3 space-y-1 [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none]">
        {/* 1. Dashboard */}
        <SidebarItem
          icon={LayoutDashboard}
          label={t[currentLang].dashboard}
          path="/dashboard"
          isActive={isRouteActive('/dashboard')}
        />
        
        {/* 2. Devices */}
        <SidebarSection
          icon={Smartphone}
          label={t[currentLang].devices}
          isActive={isRouteActive('/devices')}
          items={[
            { label: t[currentLang].allDevices, path: '/devices/all' },
            { label: t[currentLang].deviceGroups, path: '/devices/groups' },
          ]}
        />

        {/* 3. Profile */}
        <SidebarItem
          icon={Settings}
          label={t[currentLang].profile}
          path="/profiles"
          isActive={isRouteActive('/profiles')}
        />

        {/* 4. Policies */}
        <SidebarItem
          icon={Shield}
          label={t[currentLang].policies}
          path="/policies"
          isActive={isRouteActive('/policies')}
        />

        {/* 5. Applications */}
        <SidebarItem
          icon={AppWindow}
          label={t[currentLang].applications}
          path="/applications"
          isActive={isRouteActive('/applications')}
        />

        {/* 6. Users */}
        <SidebarItem
          icon={Users}
          label={t[currentLang].users}
          path="/users"
          isActive={isRouteActive('/users')}
        />
        
        <div className="pt-2 pb-1">
          <div className={cn(
            'h-px bg-gray-200 dark:bg-gray-800 mx-2',
            isCollapsed && 'mx-0'
          )} />
        </div>

        {/* Rest remains the same */}
        <SidebarItem
          icon={BarChart3}
          label={t[currentLang].reports}
          path="/reports"
          isActive={isRouteActive('/reports')}
        />

        <SidebarItem
          icon={Bell}
          label={t[currentLang].alerts}
          path="/alerts"
          isActive={isRouteActive('/alerts')}
        />

        <SidebarItem
          icon={Settings}
          label={t[currentLang].settings}
          path="/settings"
          isActive={isRouteActive('/settings')}
        />
      </nav>

      {/* Footer / Profile & Toggle */}
      <div className="p-4 border-t border-gray-200 dark:border-gray-800 shrink-0 flex flex-col gap-3 bg-gray-50 dark:bg-gray-900/50">
        {/* User Profile */}
        <div className={cn('flex items-center gap-3', isCollapsed && 'justify-center')}>
          <div className="flex items-center min-w-0 flex-1">
            <div className="w-9 h-9 rounded-full bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center shrink-0 border-2 border-primary-200 dark:border-primary-900">
              <User className="w-4 h-4 text-white" />
            </div>
            {!isCollapsed && (
              <div className="ml-3 flex flex-col min-w-0 flex-1">
                <span className="text-sm font-semibold text-gray-900 dark:text-white truncate">
                  Admin User
                </span>
                <span className="text-xs text-gray-500 dark:text-gray-400 truncate">
                  {process.env.NEXT_PUBLIC_MOCK_USERNAME || 'admin@thd.com'}
                </span>
              </div>
            )}
          </div>
          
          {!isCollapsed && (
            <Button
              isIconOnly
              variant="light"
              className="text-gray-500 hover:text-error-600 dark:text-gray-400 dark:hover:text-error-400 hover:bg-error-50 dark:hover:bg-error-950/30 transition-all duration-200"
              aria-label="Logout"
              onPress={async () => {
                await authProviderClient.logout({});
                window.location.href = '/login';
              }}
            >
              <LogOut className="w-4 h-4" />
            </Button>
          )}
        </div>

        {/* Collapse Toggle Button */}
        <Button
          onPress={toggleSidebar}
          variant="flat"
          className={cn(
            'w-full flex items-center justify-center bg-gray-100 hover:bg-gray-200 dark:bg-gray-800 dark:hover:bg-gray-700',
            'text-gray-700 dark:text-gray-300 transition-all duration-200',
            'font-medium',
            isCollapsed ? 'h-10 px-0' : 'h-10'
          )}
          aria-label={isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        >
          {isCollapsed ? (
            <ChevronRight className="w-5 h-5" />
          ) : (
            <>
              <ChevronLeft className="w-5 h-5 mr-2" />
              <span className="text-sm">{t[currentLang].collapseMenu}</span>
            </>
          )}
        </Button>
      </div>
    </aside>
    </>
  );
}
