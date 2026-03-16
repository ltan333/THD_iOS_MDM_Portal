"use client";

import { HeroUIProvider as HeroUIProviderBase } from "@heroui/system";
import { useRouter } from "next/navigation";

export interface HeroUIProviderProps {
  children: React.ReactNode;
}

export function HeroUIProvider({ children }: HeroUIProviderProps) {
  const router = useRouter();

  return <HeroUIProviderBase navigate={router.push}>{children}</HeroUIProviderBase>;
}
