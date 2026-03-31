/** @type {import('next').NextConfig} */
const nextConfig = {
  transpilePackages: ["@refinedev/antd"],
  output: "standalone",
  images: {
    qualities: [25, 50, 75, 100],
  },
  // Rewrite proxy để tránh lỗi CORS khi gọi API
  async rewrites() {
    // Lấy API URL từ biến môi trường hoặc dùng mặc định
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "https://mdm-9554.dichvu-it.vn/api/v1";
    // Đảm bảo không có trailing slash
    const cleanApiUrl = apiUrl.endsWith('/') ? apiUrl.slice(0, -1) : apiUrl;
    
    return [
      {
        source: "/api/v1/:path*",
        destination: `${cleanApiUrl}/:path*`, 
      },
    ];
  },
};

export default nextConfig;
