# Tài liệu API: Core MDM Protocol & MobileConfigs
Base URL: `/api/v1/mdm` và `/api/v1/mobile-configs`

Nhóm mạng lưới nguyên thủy đảm nhiệm việc thiết lập giao tiếp nền tảng với Apple APNs. 
*Lưu ý: Thường các API này tích hợp chạy tự động, hiếm khi Web Frontend cần tương tác.*

---

## 1. Lớp MDM Core (`/mdm`)

### 1.1 Quản lý Push Certificate (Chứng chỉ APNs)
- **`PUT /mdm/pushcert`**: Upload file .pem chứa chứng thư bảo mật để Server backend có quyền ping tới Server Apple (đây là bắt buộc để gửi lệnh Wakeup xuống thiết bị iOS).
- **`GET /mdm/pushcert`**: Trả về meta thông tin ngày hết hạn của Cert để Admin trên Portal biết đường renew.

### 1.2 Command Queue (Trạm trung chuyển lệnh)
- **`PUT /mdm/enqueue/:id`**: Nhét một lệnh thuần thiết bị (Device Command) vào Hàng chờ Queue. 
- **`GET /mdm/push/:id`**: Gửi một notification mỏng (Push Wakeup) xuống Device `id` thông qua APNs. Cú Ping này nhằm thông báo tới thiết bị rằng: "Ê, tao có lệnh đang chờ ở Backend đó, gọi Check-in lên tao lấy về mà xử lý đi".

### 1.3 `POST /mdm/escrowkeyunlock`
- **Ý nghĩa**: Mở khóa FileVault 2 (đối với MacOS) dựa trên Recovery Escrow Key lưu an toàn trên DB.

---

## 2. Lớp Mobile Config XML (`/mobile-configs`)

Một Profile (Phase 2) khi muốn áp xuống máy Apple thì nó bắt buộc phải được dịch ra định dạng *.mobileconfig XML (Property List XML).

- **`GET /mobile-configs`**: Tương tự CRUD, nơi chứa toàn bộ payload thô.
- **`GET /mobile-configs/:id/xml`**: Quan trọng! API này xuất file văn bản XML trả về cho máy khách tải & tự động cài đặt Configuration. Mọi thông tin (WiFi, Restrictions, CA Cert) do Business Layer sinh ra đều chui qua endpoint này để giao cho Apple.
