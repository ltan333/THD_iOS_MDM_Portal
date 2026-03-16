import Cookies from "js-cookie";

export const CookieUtil = {
  set: (name: string, value: string, days?: number) => {
    Cookies.set(name, value, {
      expires: days || 365, // Hết hạn sau 1 năm
      path: "/",            // Quan trọng: Để cookie dùng được cho toàn bộ web
      
      // ⚠️ ÉP CỨNG FALSE ĐỂ CHẠY ĐƯỢC TRÊN HTTP (IP AZURE)
      // Khi nào có domain HTTPS thật thì mới sửa thành true sau
      secure: false, 
      
      sameSite: "Lax",
    });
  },

  get: (name: string): string | undefined => {
    return Cookies.get(name);
  },

  remove: (name: string) => {
    Cookies.remove(name, { path: "/" });
  },

  exists: (name: string): boolean => {
    return !!Cookies.get(name);
  },
};