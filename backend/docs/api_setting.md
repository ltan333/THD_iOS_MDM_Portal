# Tài liệu API: System Settings (Cấu hình hệ thống)
Base URL: `/api/v1/settings`

Bao gồm các thao tác quản lý dạng Key-Value để tuỳ biến Backend, cấu hình SMTP, Webhook kết nối hệ sinh thái mà không phải sửa Source Code.

---

### 1. `GET /settings`
- **Ý nghĩa**: Fetch toàn bộ settings hiện tại của ứng dụng.
- **Data Source**: Bảng `settings`.
- **Workflow / cURL**:
```bash
curl -X GET "http://localhost:8080/api/v1/settings" -H "Authorization: Bearer <TOKEN>"
```

### 2. `POST /settings`
- **Ý nghĩa**: Đăng ký biến môi trường/key cấu hình mới.
- **Request Parameters**:
  - `key` (string, required): Tên định danh (vd: `smtp_host`).
  - `value` (string, required): Nội dung lưu (có thể là chuỗi hoặc JSON string).
  - `description` (optional): Diễn giải.
- **Data Source**: Lưu vào bảng `settings`.
- **Workflow / cURL**:
```bash
curl -X POST "http://localhost:8080/api/v1/settings" \
     -H "Content-Type: application/json" -H "Authorization: Bearer <TOKEN>" \
     -d '{
           "key": "webhook_slack_url",
           "value": "https://hooks.slack.com/services/...",
           "description": "Slack Alert Webhook"
         }'
```

### 3. `GET /settings/:key`
- **Ý nghĩa**: Truy xuất duy nhất 1 biến môi trường theo `key` định danh thay vì ID.
- **Data Source**: `client.Setting.Query().Where(setting.KeyEQ(key)).Only(ctx)`

### 4. `PUT /settings/:key`
- **Ý nghĩa**: Cập nhật giá trị của biến.
- **Request Parameters**:
  - `value` (string)
  - `description` (string)
- **Workflow / cURL**:
```bash
curl -X PUT "http://localhost:8080/api/v1/settings/webhook_slack_url" \
     -H "Content-Type: application/json" -H "Authorization: Bearer <TOKEN>" \
     -d '{"value": "https://new_url"}'
```

### 5. `DELETE /settings/:key`
- **Ý nghĩa**: Xóa hẳn biến khỏi DB.
