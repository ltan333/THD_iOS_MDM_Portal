export function compileParamToUrl(
  endpoint: string,
  params?: Record<string, string | number>
): string {
  if (!params) return endpoint;
  return endpoint.replace(/:([a-zA-Z0-9_]+)/g, (_, key) => {
    if (params[key] === undefined) {
      throw new Error(`Missing value for URL param: ${key}`);
    }
    return encodeURIComponent(String(params[key]));
  });
}

export function filterQueryParams(
  params?: Record<string, string | number | boolean | undefined>
): Record<string, string | number | boolean> {
  if (!params) return {};

  return Object.fromEntries(
    Object.entries(params).filter(([, value]) => {
      return value !== null && value !== undefined && value !== "";
    })
  ) as Record<string, string | number | boolean>;
}
