import type { AxiosInstance } from "axios";
import axios from "axios";

import { HEADER_KEYS } from "./constants";
import { requestInterceptor } from "./interceptors/request.interceptor";
import { responseInterceptor } from "./interceptors/response.interceptor";

// Trong development, sử dụng proxy qua Next.js rewrites để cookies hoạt động same-origin
// Nếu NEXT_PUBLIC_USE_PROXY=true, gọi qua /api/v1 (proxy)
// Nếu không, gọi trực tiếp backend URL (production hoặc khi có HTTPS)
const useProxy = process.env.NEXT_PUBLIC_USE_PROXY === "true";
const directUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost-fallback:8000/api/v1";
const baseURL = useProxy ? "/api/v1" : directUrl;

const axiosClient: AxiosInstance = axios.create({
  baseURL,
  timeout: 30000,
  withCredentials: true, // Gửi cookies với mọi request
  headers: {
    [HEADER_KEYS.CONTENT_TYPE]: "application/json",
  },
});

requestInterceptor(axiosClient.interceptors.request);
responseInterceptor(axiosClient.interceptors.response, axiosClient);

export default axiosClient;
