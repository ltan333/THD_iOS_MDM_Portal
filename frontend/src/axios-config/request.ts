import type { MethodType } from "./constants";
import type { FormDataParams } from "./utils/form-data-compiler";
import type { AxiosRequestConfig, AxiosResponse } from "axios";

import axiosClient from "./apiClient";
import { METHODS } from "./constants";
import { ObjectToFormData } from "./utils/form-data-compiler";
import { compileParamToUrl, filterQueryParams } from "./utils/url-parser";

interface RequestOptions {
  headers?: AxiosRequestConfig["headers"];
  useFormData?: boolean;
  pathParams?: Record<string, string | number>;
  queryParams?: Record<string, string | number | boolean | undefined>;
  timeout?: number;
}

export async function request<T, D = unknown>(config: {
  method: MethodType;
  endpoint: string;
  data?: D;
  options?: RequestOptions;
}): Promise<T> {
  try {
    const { method, endpoint, data, options } = config;
    const {
      headers = {},
      useFormData = false,
      pathParams,
      queryParams,
      timeout,
    } = options || {};

    const url = compileParamToUrl(endpoint, pathParams);

    const requestConfig: AxiosRequestConfig = {
      method,
      url,
      params: filterQueryParams(queryParams),
      headers,
      data: useFormData ? ObjectToFormData(data as FormDataParams) : data,
      timeout,
    };

    const response: AxiosResponse<T> = await axiosClient.request(requestConfig);

    return response.data;
  } catch (error) {
    throw error;
  }
}

export const get = <T>(endpoint: string, options?: RequestOptions) =>
  request<T>({ method: METHODS.GET, endpoint, options });

export const post = <T, D>(
  endpoint: string,
  data?: D,
  options?: RequestOptions
) => request<T, D>({ method: METHODS.POST, endpoint, data, options });

export const put = <T, D>(
  endpoint: string,
  data?: D,
  options?: RequestOptions
) => request<T, D>({ method: METHODS.PUT, endpoint, data, options });

export const patch = <T, D>(
  endpoint: string,
  data?: D,
  options?: RequestOptions
) => request<T, D>({ method: METHODS.PATCH, endpoint, data, options });

export const del = <T>(endpoint: string, options?: RequestOptions) =>
  request<T>({ method: METHODS.DELETE, endpoint, options });
