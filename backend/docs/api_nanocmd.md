# Tài liệu API: NanoCMD (Engine Lệnh & Sự kiện)

Nền tảng giao tiếp Micro-services tích hợp dự án mã nguồn mở NanoMDM, chuyển dịch các Logic Business về phía Backend giao thức thuần (Device Channel Protocol).

---

## 1. Webhook Lắng Nghe Khách (The Event Engine)

### `POST /api/v1/nanocmd/webhook` 
*(hoặc cấu hình public tùy môi trường)*
- **Ý nghĩa**: API cực kỳ sống còn. Apple thiết bị khi nhận Push Wakeup sẽ tự phát HTTP `PUT/POST` về đường link này chứa trạng thái `Acknowledge`, `NotNow` (Từ chối cài do máy khóa cất túi quá lâu), hoặc `Error`.
- **Luồng xử lý ngầm**: NanoCMD giải quyết luồng bytes này, phát tín hiệu (Event) sang cho các Module như Application, Alert (để tự update trạng thái Deployment).

## 2. The Internal Commands (`/api/v1/nanocmd`)

### 2.1 Workflow & Event
- **`POST /workflow/:name/start`**: Đẩy đồng loạt chuỗi logic. Vd: 1. Đổi tên máy -> 2. Cài Application X -> 3. Lock máy -> 4. Báo cáo. (Chains of commands).
- **`GET /event/:name`**: Lấy ra trạng thái hiện hành của sự kiện giao tiếp trên kênh Nano.

### 2.2 Inventory (Tài sản)
- **`GET /inventory`**: Trích xuất snapshot RAW (dữ liệu thô) thiết bị đã ping, thường dùng để gỡ lỗi nội bộ sâu hoặc làm trigger quét cấu hình OS Version.

> Lưu ý: Nhóm API này thuần túy hướng system-to-system integration. Web App Dashboard không gọi trực tiếp các API này mà gọi thông qua Business API (vd: `POST /devices/{id}/wipe`) — Business API sẽ tự bọc lệnh lại và gọi sang NanoCMD.
