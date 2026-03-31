import type { AxiosInstance } from "axios";
import axios from "axios";

import { HEADER_KEYS } from "./constants";
import { requestInterceptor } from "./interceptors/request.interceptor";
import { responseInterceptor } from "./interceptors/response.interceptor";

// Trong development, sử dụng proxy qua Next.js rewrites để tránh lỗi CORS
// Tùy chỉnh thông qua biến NEXT_PUBLIC_USE_PROXY=true (mặc định nên dùng proxy khi dev local)
const useProxy = process.env.NEXT_PUBLIC_USE_PROXY !== "false"; // Mặc định là true nếu không set false
const directUrl = process.env.NEXT_PUBLIC_API_URL || "https://mdm-9554.dichvu-it.vn/api/v1";
const baseURL = useProxy ? "/api/v1" : directUrl;

const axiosClient: AxiosInstance = axios.create({
  baseURL,
  timeout: 30000,
  withCredentials: true, // Gửi cookies nếu có
  headers: {
    [HEADER_KEYS.CONTENT_TYPE]: "application/json",
  },
});

requestInterceptor(axiosClient.interceptors.request);
responseInterceptor(axiosClient.interceptors.response, axiosClient);

export default axiosClient;
