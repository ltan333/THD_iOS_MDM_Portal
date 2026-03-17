import React from 'react';
import Link from 'next/link';
import { Mail, Lock, ArrowRight, ShieldCheck, Smartphone, Globe } from 'lucide-react';

export default function LoginPage() {
  return (
    <div className="min-h-screen grid grid-cols-1 lg:grid-cols-2">
      {/* Left Side - Hero Section (Enterprise Gateway Pattern) */}
      <div className="hidden lg:flex flex-col justify-center items-center relative overflow-hidden bg-slate-900 text-white p-12">
        {/* Background Elements */}
        <div className="absolute top-0 left-0 w-full h-full overflow-hidden z-0">
          <div className="absolute -top-24 -left-24 w-96 h-96 bg-primary-500/20 rounded-full blur-3xl animate-pulse-slow"></div>
          <div className="absolute top-1/2 -right-24 w-80 h-80 bg-blue-500/20 rounded-full blur-3xl animate-float-delayed"></div>
          <div className="absolute bottom-0 left-1/4 w-64 h-64 bg-purple-500/20 rounded-full blur-3xl animate-float"></div>
        </div>

        {/* Content */}
        <div className="relative z-10 max-w-lg text-center space-y-8">
          <div className="flex justify-center mb-8">
            <div className="bg-white/10 backdrop-blur-md p-6 rounded-2xl border border-white/20 shadow-2xl animate-bounce-subtle">
              <ShieldCheck className="w-16 h-16 text-primary-500" />
            </div>
          </div>
          
          <h1 className="text-4xl font-bold tracking-tight animate-fade-in-up">
            THD iOS MDM Portal
          </h1>
          <div className="grid grid-cols-3 gap-4 mt-12 animate-fade-in-up" style={{ animationDelay: '0.2s' }}>
            <div className="flex flex-col items-center gap-2 p-4 rounded-xl bg-white/5 backdrop-blur-sm border border-white/10 hover:bg-white/10 transition-colors">
              <Smartphone className="w-6 h-6 text-blue-400" />
              <span className="text-sm font-medium">Kiểm soát thiết bị</span>
            </div>
            <div className="flex flex-col items-center gap-2 p-4 rounded-xl bg-white/5 backdrop-blur-sm border border-white/10 hover:bg-white/10 transition-colors">
              <Lock className="w-6 h-6 text-green-400" />
              <span className="text-sm font-medium">Bảo mật</span>
            </div>
            <div className="flex flex-col items-center gap-2 p-4 rounded-xl bg-white/5 backdrop-blur-sm border border-white/10 hover:bg-white/10 transition-colors">
              <Globe className="w-6 h-6 text-purple-400" />
              <span className="text-sm font-medium">Truy cập toàn cầu</span>
            </div>
          </div>
        </div>
        
        {/* Footer */}
        <div className="absolute bottom-8 text-sm text-slate-500">
          © 2024 THD iOS MDM Portal. All rights reserved.
        </div>
      </div>

      {/* Right Side - Login Form */}
      <div className="flex flex-col justify-center items-center p-8 bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-50">
        <div className="w-full max-w-md space-y-8">
          <div className="text-center lg:text-left animate-fade-in-up">
            <h2 className="text-3xl font-bold tracking-tight">Đăng Nhập</h2>
        
          </div>

          <form className="mt-8 space-y-6 animate-fade-in-up" style={{ animationDelay: '0.1s' }}>
            <div className="space-y-5">
              <div>
                <label htmlFor="email" className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
                  Email hoặc tên đăng nhập
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Mail className="h-5 w-5 text-slate-400" />
                  </div>
                  <input
                    id="email"
                    name="email"
                    type="email"
                    autoComplete="email"
                    required
                    className="block w-full pl-10 pr-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg bg-white dark:bg-slate-900 text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all h-10 sm:text-sm"
                    placeholder="admin@thd.com"
                  />
                </div>
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
                  Mật Khẩu
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Lock className="h-5 w-5 text-slate-400" />
                  </div>
                  <input
                    id="password"
                    name="password"
                    type="password"
                    autoComplete="current-password"
                    required
                    className="block w-full pl-10 pr-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg bg-white dark:bg-slate-900 text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all h-10 sm:text-sm"
                    placeholder="••••••••"
                  />
                </div>
              </div>
            </div>

            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <input
                  id="remember-me"
                  name="remember-me"
                  type="checkbox"
                  className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-slate-300 rounded cursor-pointer"
                />
                <label htmlFor="remember-me" className="ml-2 block text-sm text-slate-700 dark:text-slate-300 cursor-pointer select-none">
                  Ghi nhớ đăng nhập
                </label>
              </div>

              <div className="text-sm">
                <Link href="#" className="font-medium text-primary-600 hover:text-primary-500 transition-colors">
                  Quên mật khẩu?
                </Link>
              </div>
            </div>

            <div>
              <button
                type="submit"
                className="group relative w-full flex justify-center py-2.5 px-4 border border-transparent text-sm font-semibold rounded-lg text-white bg-primary-500 hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all duration-200 shadow-lg hover:shadow-primary-500/30 cursor-pointer"
              >
                <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                  <ArrowRight className="h-5 w-5 text-primary-300 group-hover:text-primary-100 transition-colors" />
                </span>
                <span className="text-white">Đăng Nhập</span>
              </button>
            </div>
          </form>

          <div className="mt-6">
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-slate-200 dark:border-slate-700"></div>
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-slate-50 dark:bg-slate-950 text-slate-500">
                  Hoặc tiếp tục với
                </span>
              </div>
            </div>

            <div className="mt-6 grid grid-cols-2 gap-3">
              <button className="w-full inline-flex justify-center py-2 px-4 border border-slate-300 dark:border-slate-700 rounded-lg shadow-sm bg-white dark:bg-slate-800 text-sm font-medium text-slate-700 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-700 transition-colors cursor-pointer">
                <span className="sr-only">Sign in with SSO</span>
                SSO
              </button>
              <button className="w-full inline-flex justify-center py-2 px-4 border border-slate-300 dark:border-slate-700 rounded-lg shadow-sm bg-white dark:bg-slate-800 text-sm font-medium text-slate-700 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-700 transition-colors cursor-pointer">
                <span className="sr-only">Sign in with Microsoft</span>
                Microsoft
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
