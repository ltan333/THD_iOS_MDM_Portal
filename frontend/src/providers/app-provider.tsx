"use client";

import React, { Suspense } from "react";
import { Refine } from "@refinedev/core";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";
import routerProvider from "@refinedev/nextjs-router";
import { AntdRegistry } from "@ant-design/nextjs-registry";
import { ConfigProvider, theme as antdTheme } from "antd";
import viVN from "antd/locale/vi_VN";

import { ColorModeContextProvider } from "@contexts/color-mode";
import { LoadingProvider } from "@/contexts/loading";
import { DevtoolsProvider } from "@providers/devtools";
import { authProviderClient } from "@providers/auth-provider/auth-provider.client";
import { dataProvider } from "@providers/data-provider";
import { sonnerNotificationProvider } from "@providers/sonner-notification-provider";
import { i18nProvider } from "@providers/i18n-provider";
import { resources } from "@configs/resources";
import { FullPageLoader } from "@/components/loading";
import { HeroUIProvider } from "./hero-ui-provider";

interface AppProviderProps {
    children: React.ReactNode;
    defaultColorMode?: "light" | "dark";
}

export function AppProvider({
    children,
    defaultColorMode = "light",
}: AppProviderProps) {
    return (
        <RefineKbarProvider>
            <AntdRegistry>
                <ColorModeContextProvider defaultMode={defaultColorMode}>
                    <HeroUIProvider>
                        <Suspense fallback={<FullPageLoader />}>
                            <LoadingProvider>
                            <ConfigProvider
                                locale={viVN}
                                theme={{
                                    algorithm:
                                        defaultColorMode === "dark"
                                            ? antdTheme.darkAlgorithm
                                            : antdTheme.defaultAlgorithm,
                                    token: {
                                        // THD Brand Colors
                                        colorPrimary: "#1890ff",
                                        colorSuccess: "#52c41a",
                                        colorWarning: "#faad14",
                                        colorError: "#ff4d4f",
                                        colorInfo: "#1890ff",
                                        borderRadius: 6,
                                        fontFamily: "Inter, system-ui, sans-serif",
                                    },
                                    components: {
                                        Button: {
                                            borderRadius: 6,
                                        },
                                        Card: {
                                            borderRadius: 8,
                                        },
                                        Table: {
                                            borderRadius: 8,
                                        },
                                        Modal: {
                                            borderRadius: 8,
                                        },
                                    },
                                }}
                            >
                                <DevtoolsProvider>
                                    <Refine
                                        routerProvider={routerProvider}
                                        dataProvider={dataProvider}
                                        notificationProvider={sonnerNotificationProvider}
                                        i18nProvider={i18nProvider}
                                        authProvider={authProviderClient}
                                        resources={resources}
                                        options={{
                                            syncWithLocation: true,
                                            warnWhenUnsavedChanges: true,
                                            projectId: "45Om88-hpfu3D-qlH1I2",
                                            reactQuery: {
                                                clientConfig: {
                                                    defaultOptions: {
                                                        queries: {
                                                            retry: false, // Disable auto-retry, let user retry via UI button
                                                            refetchOnWindowFocus: false,
                                                        },
                                                    },
                                                },
                                            },
                                        }}
                                    >
                                        {children}
                                        <RefineKbar />
                                    </Refine>
                                </DevtoolsProvider>
                            </ConfigProvider>
                        </LoadingProvider>
                    </Suspense>
                </HeroUIProvider>
            </ColorModeContextProvider>
            </AntdRegistry>
        </RefineKbarProvider>
    );
}
