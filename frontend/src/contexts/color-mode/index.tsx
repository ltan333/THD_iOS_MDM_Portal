"use client";

import { RefineThemes } from "@refinedev/antd";
import { App as AntdApp, ConfigProvider, theme } from "antd";
import Cookies from "js-cookie";
import React, {
  type PropsWithChildren,
  createContext,
  useCallback,
  useEffect,
  useState,
} from "react";

type ColorModeContextType = {
  mode: "light" | "dark";
  setMode: (mode: "light" | "dark") => void;
  toggleMode: () => void;
};

export const ColorModeContext = createContext<ColorModeContextType>(
  {} as ColorModeContextType
);

type ColorModeContextProviderProps = {
  defaultMode?: "light" | "dark";
};

export const ColorModeContextProvider: React.FC<
  PropsWithChildren<ColorModeContextProviderProps>
> = ({ children, defaultMode = "light" }) => {
  const [isMounted, setIsMounted] = useState(false);
  const [mode, setModeState] = useState<"light" | "dark">(defaultMode);

  // Apply dark class to HTML element
  const applyThemeClass = useCallback((newMode: "light" | "dark") => {
    if (typeof document !== "undefined") {
      const root = document.documentElement;
      if (newMode === "dark") {
        root.classList.add("dark");
      } else {
        root.classList.remove("dark");
      }
    }
  }, []);

  useEffect(() => {
    setIsMounted(true);
    // Apply initial theme
    applyThemeClass(defaultMode);
  }, [defaultMode, applyThemeClass]);

  useEffect(() => {
    if (isMounted) {
      const savedTheme = Cookies.get("theme") as "light" | "dark" | undefined;
      if (savedTheme) {
        setModeState(savedTheme);
        applyThemeClass(savedTheme);
      }
    }
  }, [isMounted, applyThemeClass]);

  const setMode = useCallback(
    (newMode: "light" | "dark") => {
      setModeState(newMode);
      Cookies.set("theme", newMode, { expires: 365 });
      applyThemeClass(newMode);
    },
    [applyThemeClass]
  );

  const toggleMode = useCallback(() => {
    const newMode = mode === "light" ? "dark" : "light";
    setMode(newMode);
  }, [mode, setMode]);

  const { darkAlgorithm, defaultAlgorithm } = theme;

  // Define component tokens based on current mode
  const componentTokens = {
    Card: {
      colorBgContainer: mode === "dark" ? "#141414" : "#ffffff",
      colorBorderSecondary: mode === "dark" ? "#333333" : "#f0f0f0",
      colorText: mode === "dark" ? "#f0f0f0" : "#262626",
    },
    Upload: {
      colorFillAlter: mode === "dark" ? "#141414" : "#fafafa",
      colorBorder: mode === "dark" ? "#333333" : "#d9d9d9",
      colorText: mode === "dark" ? "#f0f0f0" : "#262626",
      colorTextDescription: mode === "dark" ? "#8c8c8c" : "#8c8c8c",
    },
    Alert: {
      colorWarningBg: mode === "dark" ? "rgba(250, 173, 20, 0.15)" : "#fffbe6",
      colorWarningBorder: mode === "dark" ? "rgba(250, 173, 20, 0.3)" : "#ffe58f",
      colorText: mode === "dark" ? "#f0f0f0" : "#262626",
    },
    Progress: {
      colorText: mode === "dark" ? "#f0f0f0" : "#262626",
      remainingColor: mode === "dark" ? "#333333" : "#e5e5e5",
    },
    Typography: {
      colorText: mode === "dark" ? "#f0f0f0" : "#262626",
      colorTextSecondary: mode === "dark" ? "#8c8c8c" : "#595959",
      colorTextHeading: mode === "dark" ? "#f0f0f0" : "#262626",
    },
    Button: {
      // Default button
      defaultBorderColor: mode === "dark" ? "#424242" : "#d9d9d9",
      defaultColor: mode === "dark" ? "#f0f0f0" : "#262626",
      defaultBg: mode === "dark" ? "transparent" : "#ffffff",
      // Primary button - remove white border
      primaryShadow: "none",
      // Disabled state
      colorBgContainerDisabled: mode === "dark" ? "#262626" : "#f5f5f5",
      colorTextDisabled: mode === "dark" ? "#595959" : "#bfbfbf",
      borderColorDisabled: mode === "dark" ? "#333333" : "#d9d9d9",
    },
    // Note: Don't add Tag tokens here as they override Ant Design preset colors (blue, red, success, etc.)
  };

  // Global token config
  const tokenConfig = {
    colorText: mode === "dark" ? "#f0f0f0" : "#262626",
    colorTextSecondary: mode === "dark" ? "#8c8c8c" : "#595959",
    colorTextTertiary: mode === "dark" ? "#595959" : "#8c8c8c",
    colorBgContainer: mode === "dark" ? "#141414" : "#ffffff",
    colorBgElevated: mode === "dark" ? "#1a1a1a" : "#ffffff",
    colorBgLayout: mode === "dark" ? "#0a0a0a" : "#f5f5f5",
    colorBorder: mode === "dark" ? "#333333" : "#d9d9d9",
    colorBorderSecondary: mode === "dark" ? "#262626" : "#f0f0f0",
  };

  return (
    <ColorModeContext.Provider
      value={{
        mode,
        setMode,
        toggleMode,
      }}
    >
      <ConfigProvider
        theme={{
          ...RefineThemes.Blue,
          algorithm: mode === "light" ? defaultAlgorithm : darkAlgorithm,
          token: tokenConfig,
          components: componentTokens,
        }}
      >
        <AntdApp>{children}</AntdApp>
      </ConfigProvider>
    </ColorModeContext.Provider>
  );
};

// Custom hook for easy access
export const useColorMode = () => {
  const context = React.useContext(ColorModeContext);
  if (!context) {
    throw new Error(
      "useColorMode must be used within ColorModeContextProvider"
    );
  }
  return context;
};
