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
    <div className="relative group">
      <button
        onClick={toggleSection}
        className={cn(
          'w-full flex items-center rounded-lg cursor-pointer transition-all duration-200 text-left',
          'hover:bg-primary-50 dark:hover:bg-primary-950/30',
          'focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2',
          isCollapsed ? 'justify-center p-3 mx-1' : 'px-3 py-2.5 gap-3',
          isActive
            ? 'bg-primary-100 text-primary-700 dark:bg-primary-950/50 dark:text-primary-400 font-semibold shadow-sm border border-primary-200 dark:border-primary-900'
            : 'text-gray-700 dark:text-gray-300 hover:text-primary-700 dark:hover:text-primary-300'
        )}
        aria-expanded={isOpen}
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

        {!isCollapsed && (
          <div className="ml-auto flex-shrink-0">
            {isOpen ? (
              <ChevronDown className="w-4 h-4 text-gray-500 dark:text-gray-400 transition-transform duration-200" />
            ) : (
              <ChevronRight className="w-4 h-4 text-gray-500 dark:text-gray-400 transition-transform duration-200" />
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
          <ul className="pl-11 space-y-0.5 py-1 border-l-2 border-primary-200 dark:border-primary-900 ml-4">
            {items.map((item) => (
              <li key={item.path}>
                <Link
                  href={item.path}
                  className={cn(
                    'block px-3 py-2 rounded-md text-sm transition-all duration-200',
                    'hover:bg-primary-50 dark:hover:bg-primary-950/30',
                    'text-gray-600 dark:text-gray-400 hover:text-primary-700 dark:hover:text-primary-300',
                    'hover:translate-x-0.5'
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
            'bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-lg shadow-xl',
            'opacity-0 invisible group-hover:opacity-100 group-hover:visible',
            'transition-all duration-200 z-50 pointer-events-none group-hover:pointer-events-auto',
            'after:content-[""] after:absolute after:-left-2 after:top-0 after:h-full after:w-2'
          )}
        >
          <div className="px-3 pb-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
            {label}
          </div>
          <ul className="space-y-0.5 pointer-events-auto">
            {items.map((item) => (
              <li key={item.path}>
                <Link
                  href={item.path}
                  className="block px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-primary-50 dark:hover:bg-primary-950/30 hover:text-primary-700 dark:hover:text-primary-300 transition-all rounded-md mx-1"
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
