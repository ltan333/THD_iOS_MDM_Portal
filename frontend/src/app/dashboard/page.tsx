'use client';

import { useLanguageStore } from '@/stores/languageStore';
import { useEffect, useState } from 'react';

export default function DashboardPage() {
  const { language } = useLanguageStore();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const t = {
    vi: {
      totalDevices: 'Tổng số thiết bị',
      activeUsers: 'Người dùng hoạt động',
      complianceAlerts: 'Cảnh báo tuân thủ',
      appsDeployed: 'Ứng dụng đã triển khai',
      deviceStatus: 'Tổng quan trạng thái thiết bị',
      chartPlaceholder: 'Thành phần biểu đồ'
    },
    en: {
      totalDevices: 'Total Devices',
      activeUsers: 'Active Users',
      complianceAlerts: 'Compliance Alerts',
      appsDeployed: 'Apps Deployed',
      deviceStatus: 'Device Status Overview',
      chartPlaceholder: 'Chart Component Placeholder'
    }
  };

  if (!mounted) return null; // Avoid hydration mismatch

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Quick Stats Cards */}
        <div className="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-slate-200 dark:border-slate-700">
          <h3 className="text-sm font-medium text-slate-500 dark:text-slate-400">{t[language].totalDevices}</h3>
          <p className="mt-2 text-3xl font-semibold text-slate-900 dark:text-white">1,248</p>
        </div>
        <div className="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-slate-200 dark:border-slate-700">
          <h3 className="text-sm font-medium text-slate-500 dark:text-slate-400">{t[language].activeUsers}</h3>
          <p className="mt-2 text-3xl font-semibold text-slate-900 dark:text-white">856</p>
        </div>
        <div className="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-slate-200 dark:border-slate-700">
          <h3 className="text-sm font-medium text-slate-500 dark:text-slate-400">{t[language].complianceAlerts}</h3>
          <p className="mt-2 text-3xl font-semibold text-rose-600 dark:text-rose-400">12</p>
        </div>
        <div className="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-slate-200 dark:border-slate-700">
          <h3 className="text-sm font-medium text-slate-500 dark:text-slate-400">{t[language].appsDeployed}</h3>
          <p className="mt-2 text-3xl font-semibold text-slate-900 dark:text-white">45</p>
        </div>
      </div>

      <div className="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-slate-200 dark:border-slate-700 min-h-[400px]">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-white mb-4">{t[language].deviceStatus}</h2>
        <div className="flex items-center justify-center h-64 text-slate-400 border-2 border-dashed border-slate-200 dark:border-slate-700 rounded-lg">
          {t[language].chartPlaceholder}
        </div>
      </div>
    </div>
  );
}