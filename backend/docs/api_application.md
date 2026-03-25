# Tài liệu API: Application Management
Base URL: `/api/v1/applications`

Quản lý ứng dụng (Store App, Enterprise App), quản lý version và push deployment (triển khai) cài ứng dụng từ xa.

---

## 1. Application Core

### 1.1 `GET /applications`
- **Ý nghĩa**: Danh sách app trên hệ thống. Dữ liệu từ Entity `Application`.
- **Query Params**: `page`, `limit`, `search`, `platform` (ios, android), `type` (app_store, enterprise, web_clip).

### 1.2 `POST /applications`
- **Ý nghĩa**: Thêm một app mới vào danh mục MDM.
- **Payload**:
  - `name`: Tên app
  - `bundle_id`: Định danh bắt buộc mã nhúng (VD: `com.google.chrome`)
  - `platform`: ios/android/windows
  - `type`: `app_store` (app public) hoặc `enterprise`
- **Workflow / cURL**:
```bash
curl -X POST "http://localhost:8080/api/v1/applications" \
     -H "Content-Type: application/json" -H "Authorization: Bearer <TOKEN>" \
     -d '{
           "name": "Chrome",
           "bundle_id": "com.google.chrome",
           "platform": "ios",
           "type": "app_store"
         }'
```

---

## 2. Application Versions (Phiên bản App)

Một Ứng dụng (Application) có thể có nhiều phiên bản (AppVersion - Dành cho Enterprise App).

### 2.1 `GET /applications/:id/versions`
- **Ý nghĩa**: Liệt kê các version con. Data từ entity `AppVersion`.

### 2.2 `POST /applications/:id/versions`
- **Ý nghĩa**: Đẩy lên version mới.
- **Payload**:
  - `version` (string)
  - `build_number` (string)
  - `file_url` (đường dẫn S3/GCP down file IPA/APK)
  - `size` (int)

---

## 3. Deployments (Triển khai App)

### 3.1 `POST /applications/deployments`
- **Ý nghĩa**: Push yêu cầu cài đặt App xuống thiết bị hoặc Nhóm thiết bị.
- **Request Payload**:
  - `app_version_id`: ID của App Version cụ thể cần push.
  - `target_type`: `device` / `group`.
  - `target_id`: ID của thiết bị hoặc nhóm.
- **Data Source**: Tạo record trong bảng `app_deployments` với status `pending`. Process nền sẽ dịch AppDeployment thành MDM Command gửi qua Apple.
- **Workflow / cURL**:
```bash
curl -X POST "http://localhost:8080/api/v1/applications/deployments" \
     -H "Content-Type: application/json" -H "Authorization: Bearer <TOKEN>" \
     -d '{
           "app_version_id": 5,
           "target_type": "device",
           "target_id": "12"
         }'
```

### 3.2 `GET /applications/:id/versions/:versionId/deployments`
- **Ý nghĩa**: Kiểm tra tiến độ cài đặt (Installing, Failed, Success... của những thiết bị đã nhận lệnh).
