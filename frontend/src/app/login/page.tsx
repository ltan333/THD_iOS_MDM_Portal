"use client";

import React, { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { Mail, Lock, ArrowRight, Smartphone, Globe, Languages } from 'lucide-react';

export default function LoginPage() {
  const [language, setLanguage] = useState<'vi' | 'en'>('vi');

  const toggleLanguage = () => {
    setLanguage(prev => prev === 'vi' ? 'en' : 'vi');
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
      globalAccess: 'Truy cập toàn cầu',
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
    <div className="h-screen w-screen relative flex items-center justify-center overflow-hidden bg-white">
      
      {/* Background Image */}
      <div className="absolute inset-0 z-0">
        <Image
          src="/assets/images/abstract-light-pink-wallpaper-background-image.jpg"
          alt="Background"
          fill
          priority
          className="object-cover object-center opacity-100 brightness-110"
          quality={100}
        />
        {/* Subtle light overlay to maintain airy feel while allowing text contrast */}
        <div className="absolute inset-0 bg-white/20 backdrop-blur-[2px]"></div>
      </div>

      {/* Main Layout Container */}
      <div className="relative z-10 w-full max-w-[1200px] h-[90vh] max-h-[850px] min-h-[600px] p-4 lg:p-8 flex items-center justify-center">
        
        {/* Unified Glass Container - Lighter Theme */}
        <div className="w-full h-full flex flex-col lg:flex-row bg-white/40 dark:bg-black/20 backdrop-blur-[40px] border border-white/60 shadow-[0_8px_32px_0_rgba(222,42,21,0.1)] rounded-[32px] overflow-hidden relative">
          
          {/* Inner highlight ring */}
          <div className="absolute inset-0 rounded-[32px] border border-white/80 pointer-events-none mix-blend-overlay"></div>

          {/* Left Side - Hero Section */}
          <div className="w-full lg:w-[45%] p-8 lg:p-12 flex flex-col justify-between relative border-b lg:border-b-0 lg:border-r border-white/30">
            
            <div className="flex flex-col items-center lg:items-start text-center lg:text-left space-y-6">
              <div className="bg-white/60 backdrop-blur-md p-3 rounded-3xl border border-white/80 shadow-[0_4px_12px_rgba(222,42,21,0.15)] flex items-center justify-center h-20 w-20 relative group transition-transform duration-500 hover:scale-105">
                <div className="absolute inset-0 bg-gradient-to-tr from-[#de2a15]/20 to-transparent rounded-3xl blur-md opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <Image 
                  src="/assets/images/favicon.png" 
                  alt="THD Logo" 
                  width={60} 
                  height={60} 
                  className="object-contain relative z-10"
                  style={{
                    filter: 'drop-shadow(0 2px 4px rgba(222,42,21,0.3))'
                  }}
                />
              </div>
              
              <div className="space-y-2">
                <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-slate-800 drop-shadow-sm leading-tight">
                  {t[language].heroTitle}
                </h1>
                <p className="text-base text-slate-600 font-medium max-w-sm">
                  {t[language].heroDesc}
                </p>
              </div>
            </div>

            <div className="hidden lg:grid grid-cols-1 gap-3 mt-8">
              <div className="flex items-center gap-4 p-3 rounded-2xl bg-white/40 backdrop-blur-md border border-white/50 shadow-sm hover:bg-white/60 transition-colors duration-300">
                <div className="p-2.5 rounded-full bg-rose-100 text-rose-600 shadow-inner">
                  <Smartphone className="w-4 h-4" />
                </div>
                <span className="text-sm font-semibold text-slate-700">{t[language].deviceControl}</span>
              </div>
              <div className="flex items-center gap-4 p-3 rounded-2xl bg-white/40 backdrop-blur-md border border-white/50 shadow-sm hover:bg-white/60 transition-colors duration-300">
                <div className="p-2.5 rounded-full bg-pink-100 text-pink-600 shadow-inner">
                  <Lock className="w-4 h-4" />
                </div>
                <span className="text-sm font-semibold text-slate-700">{t[language].securityPolicy}</span>
              </div>
              <div className="flex items-center gap-4 p-3 rounded-2xl bg-white/40 backdrop-blur-md border border-white/50 shadow-sm hover:bg-white/60 transition-colors duration-300">
                <div className="p-2.5 rounded-full bg-red-100 text-[#de2a15] shadow-inner">
                  <Globe className="w-4 h-4" />
                </div>
                <span className="text-sm font-semibold text-slate-700">{t[language].globalAccess}</span>
              </div>
            </div>
            
            <div className="mt-8 text-xs text-slate-500 font-medium hidden lg:block">
              {t[language].footer}
            </div>
          </div>

          {/* Right Side - Login Form */}
          <div className="w-full lg:w-[55%] p-8 lg:p-12 flex flex-col justify-center relative bg-white/20 h-full overflow-y-auto custom-scrollbar">
            
            {/* Language Toggle Button */}
            <div className="absolute top-6 right-6 flex items-center gap-2 z-20 bg-white/60 backdrop-blur-md px-3 py-1.5 rounded-full border border-white/80 shadow-sm">
              <span className={`text-[10px] font-bold tracking-wider transition-colors ${language === 'en' ? 'text-slate-800' : 'text-slate-400'}`}>
                EN
              </span>
              <button 
                onClick={toggleLanguage}
                className="relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-[#de2a15]/50 bg-slate-200 shadow-inner cursor-pointer"
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
                <h2 className="text-2xl lg:text-3xl font-bold tracking-tight text-slate-800">{t[language].title}</h2>
                <p className="mt-2 text-sm text-slate-600 font-medium">
                  {t[language].subtitle}
                </p>
              </div>

              <form className="space-y-5">
                <div className="space-y-4">
                  <div className="space-y-1.5">
                    <label htmlFor="email" className="block text-sm font-semibold text-slate-700 ml-1">
                      {t[language].email}
                    </label>
                    <div className="relative group">
                      <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors group-focus-within:text-[#de2a15] text-slate-400">
                        <Mail className="h-4 w-4" />
                      </div>
                      <input
                        id="email"
                        name="email"
                        type="email"
                        autoComplete="email"
                        required
                        className="block w-full pl-10 pr-4 py-3 bg-white/60 border border-white/80 rounded-xl text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-[#de2a15]/50 focus:bg-white/80 transition-all backdrop-blur-md shadow-[inset_0_2px_4px_rgba(0,0,0,0.02)]"
                        placeholder="admin@thd.com"
                      />
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <label htmlFor="password" className="block text-sm font-semibold text-slate-700 ml-1">
                      {t[language].password}
                    </label>
                    <div className="relative group">
                      <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors group-focus-within:text-[#de2a15] text-slate-400">
                        <Lock className="h-4 w-4" />
                      </div>
                      <input
                        id="password"
                        name="password"
                        type="password"
                        autoComplete="current-password"
                        required
                        className="block w-full pl-10 pr-4 py-3 bg-white/60 border border-white/80 rounded-xl text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-[#de2a15]/50 focus:bg-white/80 transition-all backdrop-blur-md shadow-[inset_0_2px_4px_rgba(0,0,0,0.02)]"
                        placeholder="••••••••"
                      />
                    </div>
                  </div>
                </div>

                <div className="flex items-center justify-between px-1">
                  <div className="flex items-center">
                    <input
                      id="remember-me"
                      name="remember-me"
                      type="checkbox"
                      className="h-3.5 w-3.5 rounded border-slate-300 bg-white text-[#de2a15] focus:ring-[#de2a15]/50 cursor-pointer"
                    />
                    <label htmlFor="remember-me" className="ml-2 block text-xs font-medium text-slate-600 cursor-pointer select-none">
                      {t[language].remember}
                    </label>
                  </div>

                  <div className="text-xs">
                    <Link href="#" className="font-semibold text-[#de2a15] hover:text-red-700 transition-colors">
                      {t[language].forgot}
                    </Link>
                  </div>
                </div>

                <div className="pt-2">
                  <button
                    type="submit"
                    className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-bold rounded-xl text-white bg-gradient-to-r from-[#de2a15] to-[#f03b26] hover:from-[#c22412] hover:to-[#de2a15] focus:outline-none focus:ring-2 focus:ring-[#de2a15]/50 transition-all duration-300 shadow-[0_4px_12px_rgba(222,42,21,0.3)] hover:shadow-[0_6px_16px_rgba(222,42,21,0.4)] overflow-hidden cursor-pointer"
                  >
                    <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-[100%] group-hover:translate-x-[100%] transition-transform duration-1000 ease-out"></div>
                    <span className="absolute left-0 inset-y-0 flex items-center pl-3.5">
                      <ArrowRight className="h-4 w-4 text-white/80 group-hover:text-white group-hover:translate-x-1 transition-all" />
                    </span>
                    {t[language].submit}
                  </button>
                </div>
              </form>
            </div>
            
            {/* Mobile Footer */}
            <div className="mt-6 text-center text-xs text-slate-500 font-medium lg:hidden pb-4">
              {t[language].footer}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
