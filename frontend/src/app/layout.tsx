import { Metadata } from "next";
import { cookies } from "next/headers";
import React from "react";

import { AppProvider } from "@providers/app-provider";
import "@styles/globals.css";

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
    <html lang="vi" className={defaultMode}>
      <body>
        <AppProvider defaultColorMode={defaultMode}>
          {children}
        </AppProvider>
      </body>
    </html>
  );
}
