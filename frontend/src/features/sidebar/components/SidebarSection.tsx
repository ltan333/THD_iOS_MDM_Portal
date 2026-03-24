'use client';

import { useState } from 'react';
import Link from 'next/link';
import { ChevronDown, ChevronRight } from 'lucide-react';
import { cn } from '@/utils/cn';
import { useSidebarStore } from '../stores/useSidebarStore';

interface SidebarSectionProps {
  icon: React.ElementType;
  label: string;
  items: { label: string; path: string }[];
  isActive?: boolean;
}

export function SidebarSection({ icon: Icon, label, items, isActive }: SidebarSectionProps) {
  const { isCollapsed } = useSidebarStore();
  const [isOpen, setIsOpen] = useState(isActive || false);

  const toggleSection = () => {
    if (!isCollapsed) setIsOpen(!isOpen);
  };

  return (
    <div className="relative group px-3 py-1">
      <button
        onClick={toggleSection}
        className={cn(
          'w-full flex items-center rounded-lg cursor-pointer transition-colors duration-200 text-left',
          'hover:bg-indigo-50 dark:hover:bg-slate-800',
          'focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500',
          isCollapsed ? 'justify-center p-2' : 'px-3 py-2',
          isActive
            ? 'bg-indigo-50 text-indigo-600 dark:bg-slate-800 dark:text-indigo-400 font-medium'
            : 'text-slate-600 dark:text-slate-400'
        )}
        aria-expanded={isOpen}
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

        {!isCollapsed && (
          <div className="ml-auto flex-shrink-0">
            {isOpen ? (
              <ChevronDown className="w-4 h-4 text-slate-400" />
            ) : (
              <ChevronRight className="w-4 h-4 text-slate-400" />
            )}
          </div>
        )}
      </button>

      {/* Expanded Submenu */}
      {!isCollapsed && (
        <div
          className={cn(
            'overflow-hidden transition-all duration-300 ease-in-out',
            isOpen ? 'max-h-48 opacity-100 mt-1' : 'max-h-0 opacity-0'
          )}
        >
          <ul className="pl-10 space-y-1 py-1 border-l-2 border-slate-100 dark:border-slate-800 ml-4">
            {items.map((item) => (
              <li key={item.path}>
                <Link
                  href={item.path}
                  className={cn(
                    'block px-3 py-1.5 rounded-md text-sm transition-colors duration-200',
                    'hover:bg-indigo-50 dark:hover:bg-slate-800',
                    'text-slate-500 dark:text-slate-400 hover:text-indigo-600 dark:hover:text-indigo-400'
                  )}
                >
                  {item.label}
                </Link>
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Tooltip for collapsed state */}
      {isCollapsed && (
        <div
          className={cn(
            'absolute left-full top-0 ml-2 py-2 w-48',
            'bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-lg shadow-xl',
            'opacity-0 invisible group-hover:opacity-100 group-hover:visible',
            'transition-all duration-200 z-50 pointer-events-none group-hover:pointer-events-auto',
            'after:content-[""] after:absolute after:-left-2 after:top-0 after:h-full after:w-2'
          )}
        >
          <div className="px-3 pb-2 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
            {label}
          </div>
          <ul className="space-y-1 pointer-events-auto">
            {items.map((item) => (
              <li key={item.path}>
                <Link
                  href={item.path}
                  className="block px-3 py-1.5 text-sm text-slate-600 dark:text-slate-300 hover:bg-indigo-50 dark:hover:bg-slate-800 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
                >
                  {item.label}
                </Link>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
