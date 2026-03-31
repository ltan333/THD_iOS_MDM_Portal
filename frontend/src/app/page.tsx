"use client";

import { Suspense, useEffect } from "react";
import { Authenticated } from "@refinedev/core";
import { useRouter } from "next/navigation";

function RedirectToDashboard() {
  const router = useRouter();
  useEffect(() => {
    router.replace("/dashboard");
  }, [router]);
  return null;
}

export default function IndexPage() {
  return (
    <Suspense>
      <Authenticated key="home-page">
        <RedirectToDashboard />
      </Authenticated>
    </Suspense>
  );
}
