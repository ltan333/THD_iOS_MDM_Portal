import React, { Suspense } from "react";
import { Metadata } from "next";
import { cookies } from "next/headers";
import { Fira_Sans, Fira_Code } from "next/font/google";

import { AppProvider } from "@providers/app-provider";
import "@styles/globals.css";
import { LayoutWrapper } from "../components/layout-wrapper";

const firaSans = Fira_Sans({
  subsets: ["latin", "vietnamese"],
  weight: ["300", "400", "500", "600", "700"],
  display: "swap",
  variable: "--font-fira-sans",
});

const firaCode = Fira_Code({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  display: "swap",
  variable: "--font-fira-code",
});

export const metadata: Metadata = {
  title: "THD MDM Portal",
  description: "THD MDM Portal System",
  icons: {
    icon: "/favicon.ico",
  },
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const cookieStore = await cookies();
  const theme = cookieStore.get("theme");
  const defaultMode = theme?.value === "dark" ? "dark" : "light";

  return (
    <html lang="vi" className={`${defaultMode} ${firaSans.variable} ${firaCode.variable}`}>
      <body className={firaSans.className}>
        <AppProvider defaultColorMode={defaultMode}>
          <LayoutWrapper>
            {children}
          </LayoutWrapper>
        </AppProvider>
      </body>
    </html>
  );
}
