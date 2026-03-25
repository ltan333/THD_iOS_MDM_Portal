# Tài liệu API: Apple DEP (Device Enrollment Program)
Base URL: `/api/v1/dep`

Tích hợp cổng Apple Business Manager (ABM). Cho phép thiết bị do Công ty mua nguyên seal, bóc hộp khởi động lên có Wifi là tự động dính ngay bộ giám sát MDM mà không cần Setup bằng tay.

---

## 1. Authentication & Token
- **`PUT /dep/tokenpki/:name`**: Upload Public Key Infrastructure đăng ký với Apple ABM.
- **`PUT /dep/tokens/:name`**: Cập nhật Server Token tải trực tiếp từ Portal của Apple vào Backend để tạo link Authorize 2 phía.

## 2. DEP Configuration (Cấu hình máy nguyên Seal)
- **`GET /dep/config/:name`** & **`PUT /dep/config/:name`**: 
  Cấu hình trải nghiệm "Mở máy lần đầu" (Setup Assistant). Quản trị viên truyền payload yêu cầu: Gắn thông tin Công ty, Bỏ qua màn hình setup Apple ID, Force cấu hình Location Service. Tự động lock màn hình bắt user đăng nhập SSO.

## 3. Quản lý Thiết Bị ABM (Proxy lệnh)
Nhóm API này hoạt động như gạch nối (Proxy) giữa Server của ta và API gốc của Apple.
- **`POST /dep/proxy/:name/devices/sync`**: API gọi lên Apple, yêu cầu trả về danh sách toàn bộ các SN thiết bị mà Apple ghi nhận thuộc quyền sở hữu của Công ty (thường gọi khi có lô hàng mua mới).
- **`POST /dep/proxy/:name/profile`**: Trả Profile thiết lập Setup Assistant mặc định (trầm phía trên) xuống thiết bị.
- **`POST /dep/proxy/:name/devices/disown`**: Trả tự do cho máy, xóa vĩnh viễn khỏi ABM (khi thanh lý / bán đồ cũ công ty).
