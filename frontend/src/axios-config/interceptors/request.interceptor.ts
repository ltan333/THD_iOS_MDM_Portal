import type {
  AxiosInterceptorManager,
  InternalAxiosRequestConfig,
} from "axios";

import { HEADER_KEYS } from "../constants";
import { tokenManager } from "../utils/token-manager";
import { Enum } from "@/configs";
import { DecryptBasic } from "@/utils/hashAes";

export const requestInterceptor = (
  request: AxiosInterceptorManager<InternalAxiosRequestConfig>
) => {
  request.use(
    (config: InternalAxiosRequestConfig) => {
      const accessToken = tokenManager.getAccessToken();
      if (accessToken) {
        config.headers.set(HEADER_KEYS.AUTHORIZATION, `Bearer ${accessToken}`);
      }

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
