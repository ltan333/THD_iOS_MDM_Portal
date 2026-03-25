# Tài liệu API: Profiles & Policies (Cấu hình)
Base URL: `/api/v1/profiles`

Quản lý Payload & Configuration Profile dùng để rải xuống thiết bị (vd chặn Camera, bắt buộc cài Passcode).

---

## 1. Profile Core CRUD

### 1.1 `GET /profiles`
- **Ý nghĩa**: Liệt kê các Profile đang có trên hệ thống, kèm theo các filter `type`, `status`.
- **Response**: Trả về `ListResponse[ProfileResponse]`, map từ Entity `Profile`.

### 1.2 `POST /profiles`
- **Ý nghĩa**: Tạo "vỏ" Profile.
- **Request Body JSON**:
  - `name`, `identifier`, `description`, `type` (device/user), `scope` (system, user).
- **Data Source**: Lưu vào table `profiles`.

---

## 2. Profile Settings (Tùy biến rules)

Dữ liệu sẽ được lưu dưới dạng JSON blob vào các trường Tương Ứng trên Entity `Profile`.

### 2.1 `PUT /profiles/:id/settings/security`
- **Ý nghĩa**: Cập nhật chuẩn an ninh mạng (Passcode, Encryption).
- **Request Payload**: 
  - `passcode` (object: `require_passcode`, `min_length`, `max_failed_attempts`...)
  - `encryption` (object: `require_storage_encryption`...)

### 2.2 `PUT /profiles/:id/settings/restrictions`
- **Ý nghĩa**: Cấm các tính năng của người dùng cuối. 
- **Request Payload**:
  - `allow_camera` (bool)
  - `allow_app_installation` (bool)
  - `allow_screen_capture` (bool)

### 2.3 `PUT /profiles/:id/settings/compliance`
- **Ý nghĩa**: Cấu hình quy chuẩn rủi ro (Compliance Rules).
- **Request Payload**: Dạng Map key-value như `{"min_os_version": "16.0", "block_jailbroken": true}`.
- **Workflow / cURL**:
```bash
curl -X PUT "http://localhost:8080/api/v1/profiles/1/settings/compliance" \
     -H "Content-Type: application/json" -H "Authorization: Bearer <TOKEN>" \
     -d '{"min_os_version": "16.0", "block_jailbroken": true}'
```

---

## 3. Assignment & Versioning (Triển khai & Cập nhật)

### 3.1 `POST /profiles/:id/assignments`
- **Ý nghĩa**: Gán Profile cho Device/DeviceGroup.
- **Request Payload**: `{"target_type": "device" | "group", "target_id": "1", "priority": 1}`
- **Data Source**: Lưu vào Entity `ProfileAssignment`. Hành động này cũng kéo theo Event cập nhật xuống NanoCMD Webhook (MDM Command).

### 3.2 `POST /profiles/:id/repush`
- **Ý nghĩa**: Cưỡng ép cài đặt lại Profile lên tất cả thiết bị đang được gán (assignment). 
- **Workflow / cURL**:
```bash
curl -X POST "http://localhost:8080/api/v1/profiles/1/repush" -H "Authorization: Bearer <TOKEN>"
```

### 3.3 `POST /profiles/:id/versions/:versionId/rollback`
- **Ý nghĩa**: Quay lùi version của cấu hình profile.
- **Data Source**: Phục hồi dữ liệu từ bảng `profile_versions` ngược về `profiles`.
