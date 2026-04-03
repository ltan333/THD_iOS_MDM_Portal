"use client";

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { Mail, Lock, ArrowRight, Smartphone, Globe, Eye, EyeOff } from 'lucide-react';
import { create } from 'zustand';
import { authProviderClient } from '@/providers/auth-provider/auth-provider.client';

type Language = 'en' | 'vi';

interface LanguageState {
  language: Language;
  toggleLanguage: () => void;
}

const useLanguageStore = create<LanguageState>((set) => ({
  language: 'en',
  toggleLanguage: () =>
    set((state) => ({ language: state.language === 'en' ? 'vi' : 'en' })),
}));

export default function LoginPage() {
  const { language, toggleLanguage } = useLanguageStore();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [mounted, setMounted] = useState(false);
  const router = useRouter();

  // Prevent hydration mismatch
  useEffect(() => {
    setMounted(true);
  }, []);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    const isMockEnabled = process.env.NEXT_PUBLIC_MOCK_LOGIN_ENABLED === 'true';
    const mockEmail = process.env.NEXT_PUBLIC_MOCK_USERNAME || 'admin@thd.com';
    const mockPassword = process.env.NEXT_PUBLIC_MOCK_PASSWORD || 'password123';

    if (isMockEnabled) {
      if (email === mockEmail && password === mockPassword) {
        window.location.href = '/dashboard';
        return;
      }

      setError(language === 'vi' ? 'Email hoặc mật khẩu không chính xác.' : 'Invalid email or password.');
      return;
    }

    try {
      const result = await authProviderClient.login({ username: email, password });

      if (result.success) {
        // Use window.location.href for a full page reload to clear any stale state
        window.location.href = result.redirectTo || '/dashboard';
        return;
      }

      setError(result.error?.message || (language === 'vi' ? 'Đăng nhập thất bại.' : 'Login failed.'));
    } catch {
      setError(language === 'vi' ? 'Đăng nhập thất bại.' : 'Login failed.');
    }
  };

  const t = {
    vi: {
      title: 'Đăng Nhập',
      subtitle: 'Vui lòng nhập thông tin đăng nhập để truy cập trang quản trị.',
      email: 'Địa chỉ Email hoặc Tài khoản',
      password: 'Mật khẩu',
      remember: 'Ghi nhớ đăng nhập',
      forgot: 'Quên mật khẩu?',
      submit: 'Đăng nhập',
      or: 'Hoặc tiếp tục với',
      heroTitle: 'THD iOS MDM Portal',
      heroDesc: 'MDM Portal for iOS devices',
      deviceControl: 'Kiểm soát thiết bị',
      securityPolicy: 'Chính sách bảo mật',
      globalAccess: 'Quản lý đa nền tảng',
      footer: '© 2024 Cổng thông tin THD iOS MDM. Đã đăng ký Bản quyền.'
    },
    en: {
      title: 'Welcome back',
      subtitle: 'Please enter your credentials to access the admin dashboard.',
      email: 'Email address or username',
      password: 'Password',
      remember: 'Remember me',
      forgot: 'Forgot your password?',
      submit: 'Sign in',
      or: 'Or continue with',
      heroTitle: 'THD iOS MDM Portal',
      heroDesc: 'MDM Portal for iOS devices',
      deviceControl: 'Device Control',
      securityPolicy: 'Security Policy',
      globalAccess: 'Global Access',
      footer: '© 2024 THD iOS MDM Portal. All rights reserved.'
    }
  };

  return (
    <div className="h-screen w-screen relative flex items-center justify-center overflow-hidden bg-slate-950">
      
      {/* Custom background image */}
      <div className="absolute inset-0 z-0">
        <Image
          src="/assets/images/456c795b-d275-44d2-b958-fa6ba4f5e5ec.jpg"
          alt="Login Background"
          fill
          priority
          className="object-cover object-center"
          quality={100}
        />
        <div className="absolute inset-0 bg-slate-900/35"></div>
      </div>

      {/* Main Layout Container */}
      <div className="relative z-10 w-full max-w-[1200px] h-[90vh] max-h-[850px] min-h-[600px] p-4 lg:p-8 flex items-center justify-center opacity-0 transition-opacity duration-300" style={{ opacity: mounted ? 1 : 0 }}>
        
        {/* Unified Glass Container - Lighter Theme */}
        <div className="w-full h-full flex flex-col lg:flex-row bg-white/12 dark:bg-white/5 backdrop-blur-[32px] border border-white/25 shadow-[0_20px_60px_rgba(2,6,23,0.45)] rounded-[32px] overflow-hidden relative">
          
          {/* Inner highlight ring */}
          <div className="absolute inset-0 rounded-[32px] border border-white/80 pointer-events-none mix-blend-overlay"></div>

          {/* Left Side - Hero Section */}
          <div className="w-full lg:w-[45%] p-8 lg:p-12 flex flex-col justify-between relative border-b lg:border-b-0 lg:border-r border-white/30">
            
            <div className="flex flex-col items-center lg:items-start text-center lg:text-left space-y-6">
              <div className="bg-white/15 backdrop-blur-md p-3 rounded-3xl border border-white/30 shadow-[0_8px_18px_rgba(15,23,42,0.35)] flex items-center justify-center h-20 w-20 relative group transition-transform duration-500 hover:scale-105">
                <div className="absolute inset-0 bg-gradient-to-tr from-blue-400/30 to-transparent rounded-3xl blur-md opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <Image 
                  src="/assets/images/favicon.png" 
                  alt="THD Logo" 
                  width={60} 
                  height={60} 
                  className="object-contain relative z-10"
                  style={{
                    filter: 'drop-shadow(0 4px 8px rgba(59,130,246,0.35))'
                  }}
                />
              </div>
              
              <div className="space-y-2">
                <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-white drop-shadow-sm leading-tight">
                  {mounted ? t[language].heroTitle : t.vi.heroTitle}
                </h1>
                <p className="text-base text-slate-200/90 font-medium max-w-sm">
                  {mounted ? t[language].heroDesc : t.vi.heroDesc}
                </p>
              </div>
            </div>

            <div className="hidden lg:grid grid-cols-1 gap-3 mt-8">
              <div className="flex items-center gap-4 p-3 rounded-2xl bg-white/10 backdrop-blur-md border border-white/20 shadow-sm hover:bg-white/15 transition-colors duration-300">
                <div className="p-2.5 rounded-full bg-blue-500/20 text-blue-200 shadow-inner">
                  <Smartphone className="w-4 h-4" />
                </div>
                <span className="text-sm font-semibold text-slate-100">{mounted ? t[language].deviceControl : t.vi.deviceControl}</span>
              </div>
              <div className="flex items-center gap-4 p-3 rounded-2xl bg-white/10 backdrop-blur-md border border-white/20 shadow-sm hover:bg-white/15 transition-colors duration-300">
                <div className="p-2.5 rounded-full bg-indigo-500/20 text-indigo-200 shadow-inner">
                  <Lock className="w-4 h-4" />
                </div>
                <span className="text-sm font-semibold text-slate-100">{mounted ? t[language].securityPolicy : t.vi.securityPolicy}</span>
              </div>
              <div className="flex items-center gap-4 p-3 rounded-2xl bg-white/10 backdrop-blur-md border border-white/20 shadow-sm hover:bg-white/15 transition-colors duration-300">
                <div className="p-2.5 rounded-full bg-cyan-500/20 text-cyan-200 shadow-inner">
                  <Globe className="w-4 h-4" />
                </div>
                <span className="text-sm font-semibold text-slate-100">{mounted ? t[language].globalAccess : t.vi.globalAccess}</span>
              </div>
            </div>
            
            <div className="mt-8 text-xs text-slate-300/80 font-medium hidden lg:block">
              {mounted ? t[language].footer : t.vi.footer}
            </div>
          </div>

          {/* Right Side - Login Form */}
          <div className="w-full lg:w-[55%] p-8 lg:p-12 flex flex-col justify-center relative bg-white/10 h-full overflow-y-auto custom-scrollbar">
            
            {/* Language Toggle Button */}
            <div className="absolute top-6 right-6 flex items-center gap-2 z-20 bg-white/15 backdrop-blur-md px-3 py-1.5 rounded-full border border-white/30 shadow-sm">
              <span className={`text-[10px] font-bold tracking-wider transition-colors ${language === 'en' ? 'text-slate-800' : 'text-slate-400'}`}>
                EN
              </span>
              <button 
                onClick={toggleLanguage}
                className="relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-400/50 bg-slate-300/80 shadow-inner cursor-pointer"
                role="switch"
                aria-checked={language === 'vi'}
              >
                <span className="sr-only">Toggle Language</span>
                <span
                  className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white transition-transform duration-300 ease-spring ${
                    language === 'vi' ? 'translate-x-4' : 'translate-x-1'
                  } shadow-md`}
                />
              </button>
              <span className={`text-[10px] font-bold tracking-wider transition-colors ${language === 'vi' ? 'text-slate-800' : 'text-slate-400'}`}>
                VI
              </span>
            </div>

            <div className="w-full max-w-sm mx-auto space-y-8 my-auto">
              <div className="text-center lg:text-left">
                <h2 className="text-2xl lg:text-3xl font-bold tracking-tight text-white">{mounted ? t[language].title : t.vi.title}</h2>
                <p className="mt-2 text-sm text-slate-200/90 font-medium">
                  {mounted ? t[language].subtitle : t.vi.subtitle}
                </p>
              </div>

              <form className="space-y-5" onSubmit={handleLogin}>
                <div className="space-y-4">
                  {error && (
                    <div className="p-3 bg-red-50 border border-red-200 text-red-600 text-sm rounded-xl">
                      {error}
                    </div>
                  )}
                  <div className="space-y-1.5">
                    <label htmlFor="email" className="block text-sm font-semibold text-slate-100 ml-1">
                      {mounted ? t[language].email : t.vi.email}
                    </label>
                    <div className="relative group">
                      <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors group-focus-within:text-blue-500 text-slate-400">
                        <Mail className="h-4 w-4" />
                      </div>
                      <input
                        id="email"
                        name="email"
                        type="text"
                        autoComplete="username"
                        required
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className="block w-full pl-10 pr-4 py-3 bg-white/85 border border-white/70 rounded-xl text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/35 focus:bg-white transition-all backdrop-blur-md shadow-[inset_0_2px_4px_rgba(0,0,0,0.02)]"
                        placeholder="admin@thd.com hoặc admin"
                      />
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <label htmlFor="password" className="block text-sm font-semibold text-slate-100 ml-1">
                      {mounted ? t[language].password : t.vi.password}
                    </label>
                    <div className="relative group">
                      <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors group-focus-within:text-blue-500 text-slate-400">
                        <Lock className="h-4 w-4" />
                      </div>
                      <input
                        id="password"
                        name="password"
                        type={showPassword ? "text" : "password"}
                        autoComplete="current-password"
                        required
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className="block w-full pl-10 pr-10 py-3 bg-white/85 border border-white/70 rounded-xl text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/35 focus:bg-white transition-all backdrop-blur-md shadow-[inset_0_2px_4px_rgba(0,0,0,0.02)]"
                        placeholder="••••••••"
                      />
                      <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-600 transition-colors focus:outline-none"
                      >
                        {showPassword ? (
                          <EyeOff className="h-4 w-4" />
                        ) : (
                          <Eye className="h-4 w-4" />
                        )}
                      </button>
                    </div>
                  </div>
                </div>

                <div className="flex items-center justify-between px-1">
                  <div className="flex items-center">
                    <input
                      id="remember-me"
                      name="remember-me"
                      type="checkbox"
                      className="h-3.5 w-3.5 rounded border-slate-300 bg-white text-blue-600 focus:ring-blue-500/50 cursor-pointer"
                    />
                    <label htmlFor="remember-me" className="ml-2 block text-xs font-medium text-slate-600 cursor-pointer select-none">
                      {mounted ? t[language].remember : t.vi.remember}
                    </label>
                  </div>

                  <div className="text-xs">
                    <Link href="#" className="font-semibold text-blue-300 hover:text-blue-200 transition-colors">
                      {mounted ? t[language].forgot : t.vi.forgot}
                    </Link>
                  </div>
                </div>

                <div className="pt-2">
                  <button
                    type="submit"
                    className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-bold rounded-xl text-white bg-gradient-to-r from-blue-600 to-cyan-500 hover:from-blue-700 hover:to-cyan-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all duration-300 shadow-[0_6px_16px_rgba(37,99,235,0.35)] hover:shadow-[0_8px_20px_rgba(37,99,235,0.45)] overflow-hidden cursor-pointer"
                  >
                    <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-[100%] group-hover:translate-x-[100%] transition-transform duration-1000 ease-out"></div>
                    <span className="absolute left-0 inset-y-0 flex items-center pl-3.5">
                      <ArrowRight className="h-4 w-4 text-white/80 group-hover:text-white group-hover:translate-x-1 transition-all" />
                    </span>
                    {mounted ? t[language].submit : t.vi.submit}
                  </button>
                </div>
              </form>
            </div>
            
            {/* Mobile Footer */}
            <div className="mt-6 text-center text-xs text-slate-300/80 font-medium lg:hidden pb-4">
              {mounted ? t[language].footer : t.vi.footer}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
