'use client';

import Link from 'next/link';
import { cn } from '@/utils/cn';
import { useSidebarStore } from '../stores/useSidebarStore';

interface SidebarItemProps {
  icon: React.ElementType;
  label: string;
  path: string;
  isActive?: boolean;
}

export function SidebarItem({ icon: Icon, label, path, isActive }: SidebarItemProps) {
  const { isCollapsed } = useSidebarStore();

  return (
    <div className="relative group">
      <Link
        href={path}
        className={cn(
          'flex items-center rounded-lg cursor-pointer transition-all duration-200',
          'hover:bg-primary-50 dark:hover:bg-primary-950/30',
          'focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2',
          isCollapsed ? 'justify-center p-3 mx-1' : 'px-3 py-2.5 gap-3',
          isActive
            ? 'bg-primary-100 text-primary-700 dark:bg-primary-950/50 dark:text-primary-400 font-semibold shadow-sm border border-primary-200 dark:border-primary-900'
            : 'text-gray-700 dark:text-gray-300 hover:text-primary-700 dark:hover:text-primary-300'
        )}
        aria-label={label}
        aria-current={isActive ? 'page' : undefined}
      >
        <Icon
          className={cn(
            'flex-shrink-0 transition-all duration-200',
            isCollapsed ? 'w-6 h-6' : 'w-5 h-5',
            isActive ? 'text-primary-600 dark:text-primary-400' : 'text-gray-500 dark:text-gray-400 group-hover:text-primary-600 dark:group-hover:text-primary-400'
          )}
        />
        
        {!isCollapsed && (
          <span className="truncate flex-1 text-sm">{label}</span>
        )}
        
        {!isCollapsed && isActive && (
          <div className="w-1.5 h-1.5 rounded-full bg-primary-600 dark:bg-primary-400 shrink-0" />
        )}
      </Link>

      {/* Tooltip for collapsed state */}
      {isCollapsed && (
        <div
          className={cn(
            'absolute left-full top-1/2 -translate-y-1/2 ml-2',
            'px-3 py-2 bg-gray-900 dark:bg-gray-800 text-white text-sm rounded-lg shadow-lg whitespace-nowrap',
            'opacity-0 invisible group-hover:opacity-100 group-hover:visible',
            'transition-all duration-200 z-50 pointer-events-none',
            'border border-gray-700 dark:border-gray-600'
          )}
        >
          {label}
          <div className="absolute right-full top-1/2 -translate-y-1/2 border-4 border-transparent border-r-gray-900 dark:border-r-gray-800" />
        </div>
      )}
    </div>
  );
}
