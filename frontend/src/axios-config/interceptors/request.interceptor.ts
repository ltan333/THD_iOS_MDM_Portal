import type {
  AxiosInterceptorManager,
  InternalAxiosRequestConfig,
} from "axios";

import { AUTH_CONFIG, HEADER_KEYS } from "../constants";
import { tokenManager } from "../utils/token-manager";
import { Enum } from "@/configs";
import { DecryptBasic } from "@/utils/hashAes";

export const requestInterceptor = (
  request: AxiosInterceptorManager<InternalAxiosRequestConfig>
) => {
  request.use(
    (config: InternalAxiosRequestConfig) => {
      // Backend sử dụng session cookies - browser tự động gửi qua withCredentials: true
      // Không cần thêm Authorization header vì backend check session cookies

      // Thêm organization key nếu có
      const orgToken = tokenManager.getOrgToken();
      if (orgToken) {
        const orgKey = DecryptBasic(orgToken, Enum.secretKey);
        if (orgKey) {
          config.headers.set(HEADER_KEYS.X_ORGANIZATION_KEY, orgKey);
        }
      }

      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );
};
