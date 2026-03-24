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
    <div className="relative group px-3 py-1">
      <Link
        href={path}
        className={cn(
          'flex items-center rounded-lg cursor-pointer transition-colors duration-200',
          'hover:bg-indigo-50 dark:hover:bg-slate-800',
          'focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500',
          isCollapsed ? 'justify-center p-2' : 'px-3 py-2',
          isActive
            ? 'bg-indigo-50 text-indigo-600 dark:bg-slate-800 dark:text-indigo-400 font-medium'
            : 'text-slate-600 dark:text-slate-400'
        )}
        aria-label={label}
        aria-current={isActive ? 'page' : undefined}
      >
        <Icon
          className={cn(
            'flex-shrink-0 transition-colors duration-200',
            isCollapsed ? 'w-6 h-6' : 'w-5 h-5 mr-3',
            isActive ? 'text-indigo-600 dark:text-indigo-400' : 'text-slate-500 dark:text-slate-400'
          )}
        />
        
        {!isCollapsed && (
          <span className="truncate flex-1">{label}</span>
        )}
      </Link>

      {/* Tooltip for collapsed state */}
      {isCollapsed && (
        <div
          className={cn(
            'absolute left-full top-1/2 -translate-y-1/2 ml-2',
            'px-2 py-1 bg-slate-800 text-white text-xs rounded shadow-lg whitespace-nowrap',
            'opacity-0 invisible group-hover:opacity-100 group-hover:visible',
            'transition-all duration-200 z-50 pointer-events-none'
          )}
        >
          {label}
        </div>
      )}
    </div>
  );
}
