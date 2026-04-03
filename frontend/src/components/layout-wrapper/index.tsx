'use client';

import { usePathname, useRouter } from 'next/navigation';
import { Sidebar } from '@/features/sidebar';
import { Header } from '@/features/header';
import { useSidebarStore } from '@/features/sidebar/stores/useSidebarStore';
import { useEffect, useState } from 'react';
import { cn } from '@/utils/cn';
import { tokenManager } from '@/axios-config/utils/token-manager';

export function LayoutWrapper({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const { isCollapsed } = useSidebarStore();
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  
  const isLoginPage = pathname === '/login';

  useEffect(() => {
    // Check authentication status using standard token manager
    const token = tokenManager.getAccessToken();
    
    // Nếu chưa có token và không ở trang login -> Đẩy về login
    if (!token && !isLoginPage) {
      router.replace('/login');
    } 
    // CHỈ đẩy về dashboard nếu CÓ TOKEN và đang cố vào trang login
    else if (token && isLoginPage) {
      router.replace('/dashboard');
    } else {
      setIsAuthenticated(!!token);
    }
  }, [pathname, router, isLoginPage]);

  // Don't render anything while checking auth to prevent hydration mismatch
  // and flickering of unauthorized content
  if (isAuthenticated === null && !isLoginPage) {
    return null;
  }

  if (isLoginPage) {
    return <>{children}</>;
  }

  return (
    <div className="app-liquid-ui flex h-screen overflow-hidden bg-gradient-to-br from-slate-100 via-slate-50 to-blue-50 dark:from-gray-950 dark:via-slate-950 dark:to-slate-900">
      <Sidebar />
      <div 
        className={cn(
          "flex flex-col flex-1 min-w-0 overflow-hidden transition-all duration-300 ease-in-out",
          isCollapsed ? "lg:ml-[72px]" : "lg:ml-[260px]"
        )}
      >
        <Header />
        <main className="flex-1 overflow-y-auto">
          <div className="w-full h-full max-w-[1600px] mx-auto p-5 liquid-layout-main liquid-glass">
            <div key={pathname} className="motion-safe-page h-full">
              {children}
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
