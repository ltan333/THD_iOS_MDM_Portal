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
          'bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800',
          'transition-all duration-300 ease-in-out',
          isCollapsed ? '-translate-x-full lg:translate-x-0 w-[260px] lg:w-[72px]' : 'translate-x-0 w-[260px]'
        )}
        aria-label="Sidebar Navigation"
      >
      {/* Header / Logo */}
      <div className="flex items-center justify-between h-16 px-4 border-b border-slate-200 dark:border-slate-800 shrink-0">
        <div className={cn('flex items-center', isCollapsed && 'justify-center w-full')}>
          <div className="w-8 h-8 rounded-lg bg-indigo-600 flex items-center justify-center shrink-0">
            <Shield className="w-5 h-5 text-white" />
          </div>
          {!isCollapsed && (
            <span className="ml-3 font-semibold text-slate-900 dark:text-white truncate">
              MDM Portal
            </span>
          )}
        </div>
      </div>

      {/* Main Navigation */}
      <nav className="flex-1 overflow-y-auto py-4 px-2 space-y-1 [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none]">
        <SidebarItem
          icon={LayoutDashboard}
          label={t[currentLang].dashboard}
          path="/dashboard"
          isActive={isRouteActive('/dashboard')}
        />
        
        <SidebarItem
          icon={Settings}
          label={t[currentLang].profile}
          path="/profiles"
          isActive={isRouteActive('/profiles')}
        />
        
        <SidebarSection
          icon={Smartphone}
          label={t[currentLang].devices}
          isActive={isRouteActive('/devices')}
          items={[
            { label: t[currentLang].allDevices, path: '/devices/all' },
            { label: t[currentLang].deviceGroups, path: '/devices/groups' },
          ]}
        />

        <SidebarItem
          icon={Users}
          label={t[currentLang].users}
          path="/users"
          isActive={isRouteActive('/users')}
        />

        <SidebarItem
          icon={AppWindow}
          label={t[currentLang].applications}
          path="/applications"
          isActive={isRouteActive('/applications')}
        />

        <SidebarItem
          icon={Shield}
          label={t[currentLang].policies}
          path="/policies"
          isActive={isRouteActive('/policies')}
        />

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
      <div className="p-4 border-t border-slate-200 dark:border-slate-800 shrink-0 flex flex-col gap-4">
        {/* User Profile */}
        <div className={cn('flex items-center', isCollapsed ? 'justify-center' : 'justify-between')}>
          <div className="flex items-center">
            <div className="w-8 h-8 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center shrink-0 border border-slate-200 dark:border-slate-700">
              <User className="w-4 h-4 text-slate-600 dark:text-slate-400" />
            </div>
            {!isCollapsed && (
              <div className="ml-3 flex flex-col">
                <span className="text-sm font-medium text-slate-900 dark:text-white truncate max-w-[120px]">
                  Admin User
                </span>
                <span className="text-xs text-slate-500 dark:text-slate-400 truncate max-w-[120px]">
                  {process.env.NEXT_PUBLIC_MOCK_USERNAME || 'admin@thd.com'}
                </span>
              </div>
            )}
          </div>
          
          {!isCollapsed && (
            <Button
              isIconOnly
              variant="light"
              className="text-slate-500 hover:text-red-600 dark:text-slate-400 dark:hover:text-red-400"
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
            'w-full flex items-center justify-center bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-600 dark:text-slate-300 transition-colors',
            isCollapsed ? 'h-10 px-0' : 'h-10'
          )}
          aria-label={isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        >
          {isCollapsed ? (
            <ChevronRight className="w-5 h-5" />
          ) : (
            <>
              <ChevronLeft className="w-5 h-5 mr-2" />
              <span>{t[currentLang].collapseMenu}</span>
            </>
          )}
        </Button>
      </div>
    </aside>
    </>
  );
}
