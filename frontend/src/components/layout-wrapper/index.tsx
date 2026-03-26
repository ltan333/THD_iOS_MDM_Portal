'use client';

import { usePathname, useRouter } from 'next/navigation';
import { Sidebar } from '@/features/sidebar';
import { Header } from '@/features/header';
import { useSidebarStore } from '@/features/sidebar/stores/useSidebarStore';
import { useEffect, useState } from 'react';
import { cn } from '@/utils/cn';

export function LayoutWrapper({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const { isCollapsed } = useSidebarStore();
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  
  const isLoginPage = pathname === '/login';

  useEffect(() => {
    // Check authentication status
    const token = localStorage.getItem('auth_token');
    
    if (!token && !isLoginPage) {
      router.push('/login');
    } else if (token && isLoginPage) {
      router.push('/dashboard');
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
    <div className="flex h-screen overflow-hidden bg-slate-50 dark:bg-slate-900">
      <Sidebar />
      <div 
        className={cn(
          "flex flex-col flex-1 min-w-0 overflow-hidden transition-all duration-300 ease-in-out bg-slate-50 dark:bg-slate-900",
          isCollapsed ? "lg:ml-[72px]" : "lg:ml-[260px]"
        )}
      >
        <Header />
        <main className="flex-1 overflow-y-auto">
          <div className="w-full h-full max-w-[1600px] mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}