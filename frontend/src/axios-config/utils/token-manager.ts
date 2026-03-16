import { TOKEN_KEYS } from "../constants";
import { CookieUtil } from "@/utils/cookieUtils";

class TokenManager {
  private getFromCookie(key: string): string | null {
    try {
      return CookieUtil.get(key) || null;
    } catch (error) {
      console.warn(`Failed to get ${key} from cookies:`, error);
      return null;
    }
  }

  private setToCookie(key: string, value: string, hours: number = 6.5): void {
    try {
      CookieUtil.set(key, value, hours);
    } catch (error) {
      console.warn(`Failed to set ${key} to cookies:`, error);
    }
  }

  private removeFromCookie(key: string): void {
    try {
      CookieUtil.remove(key);
    } catch (error) {
      console.warn(`Failed to remove ${key} from cookies:`, error);
    }
  }

  getAccessToken(): string | null {
    return this.getFromCookie(TOKEN_KEYS.ACCESS_TOKEN);
  }

  getRefreshToken(): string | null {
    return this.getFromCookie(TOKEN_KEYS.REFRESH_TOKEN);
  }

  getOrgToken(): string | null {
    return this.getFromCookie(TOKEN_KEYS.ORG_TOKEN);
  }

  setAccessToken(token: string): void {
    this.setToCookie(TOKEN_KEYS.ACCESS_TOKEN, token, 6.5);
  }

  setRefreshToken(token: string): void {
    this.setToCookie(TOKEN_KEYS.REFRESH_TOKEN, token, 6.5);
  }

  setTokens(accessToken: string, refreshToken: string): void {
    this.setAccessToken(accessToken);
    this.setRefreshToken(refreshToken);
  }

  clearTokens(): void {
    this.removeFromCookie(TOKEN_KEYS.ACCESS_TOKEN);
    this.removeFromCookie(TOKEN_KEYS.REFRESH_TOKEN);
    this.removeFromCookie(TOKEN_KEYS.ORG_TOKEN);
    this.removeFromCookie("us_r");

    // Also clear old cookie names for backward compatibility
    this.removeFromCookie("ac_token");
    this.removeFromCookie("rf_token");

    // Clear localStorage
    if (typeof window !== "undefined") {
      localStorage.removeItem("user_ifo");
    }
  }

  isAuthenticated(): boolean {
    return !!this.getAccessToken();
  }
}

export const tokenManager = new TokenManager();
