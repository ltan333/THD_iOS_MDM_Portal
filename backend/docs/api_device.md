# Tài liệu API: Quản lý Thiết Bị (Devices & Device Groups)
Base URL: `/api/v1`

Đây là nhóm API điều khiển luồng dữ liệu trung tâm của MDM Portal, nơi frontend giao tiếp với Server để liệt kê thiết bị, phân nhóm và truyền lệnh khóa (Lock/Wipe) xuống thiết bị vật lý.

---

## 1. Devices (Thiết bị)

### 1.1 `GET /devices`
- **Ý nghĩa**: Lấy danh sách tất cả các thiết bị đã và đang Enrollment vào hệ thống. Hỗ trợ phân trang, lọc và tìm kiếm text.
- **Request Parameters**:
  - `page` (int, default=1), `limit` (int, default=20): Phân trang.
  - `search` (string): Tìm kiếm Text (Map từ DB field: `name`, `serial_number`, `model`).
  - `status`, `platform`, `is_raw` (string/bool): Bộ lọc tùy chọn.
- **Data Source**: Truy vấn trực tiếp từ bảng `devices` (Ent schema `Device`).
- **Response**: Trả về struct `ListResponse[dto.DeviceResponse]`.
  - `platform` (OS: ios/android/windows), `mac_address`, `ip_address`, `battery_level`... do thiết bị đẩy lên Apple sau đó Server listen qua NanoCMD Webhook và cập nhật lại DB.
- **Workflow / cURL**:
```bash
curl -X GET "http://localhost:8080/api/v1/devices?page=1&limit=50&search=iPhone" \
     -H "Authorization: Bearer <TOKEN_HERE>"
```

### 1.2 `GET /devices/:id`
- **Ý nghĩa**: Lấy toàn bộ thông tin chi tiết một thiết bị theo Khóa chính (Primary Key `id`).
- **Data Source**: Câu lệnh `client.Device.Query().WithGroups().WithProfiles().Where(device.IDEQ(id)).Only(ctx)`. Kết quả bao gồm cả các DeviceGroup thiết bị này thuộc về, và các Profiles đang được gán.
- **Workflow / cURL**:
```bash
curl -X GET "http://localhost:8080/api/v1/devices/1" -H "Authorization: Bearer <TOKEN_HERE>"
```

### 1.3 `POST /devices/:id/lock`
- **Ý nghĩa**: Kích hoạt khóa thiết bị từ xa (DeviceLock). Action này không lưu trực tiếp thông tin vào bảng Device, mà tương tác với `NanoCMDService`. Nó sẽ sinh ra một lệnh XML (Command XML) đẩy vào Event Queue để Apple APNs Server bốc gửi xuống điện thoại.
- **Response**: Trả về `{ "message": "Command initiated" }` nếu lệnh vào Hàng đợi (Queue) thành công.
- **Workflow / cURL**:
```bash
curl -X POST "http://localhost:8080/api/v1/devices/1/lock" -H "Authorization: Bearer <TOKEN>"
```

### 1.4 `POST /devices/:id/wipe`
- **Ý nghĩa**: Kích hoạt quá trình tự hủy xóa ổ cứng diện rộng (EraseDevice).
- **Data Source**: Phía sau cũng gọi qua hệ thống Push Command tới APNs giống `lock`. Thiết bị nhận được sẽ bị Factory Reset vĩnh viễn và không thể khôi phục (nếu không có backup).

---

## 2. Device Groups (Nhóm thiết bị)

### 2.1 `GET /device-groups`
- **Ý nghĩa**: Liệt kê các tổ chức hoặc nhóm gom chung thiết bị (ví dụ: "Nhóm Kế toán", "Thiết bị Công ty cấp"). Đôi khi có thể lồng nhóm cha-con (`parent_id`).
- **Query Params**: `page`, `limit`, `search` (tìm theo `name`).
- **Workflow / cURL**:
```bash
curl -X GET "http://localhost:8080/api/v1/device-groups" -H "Authorization: Bearer <TOKEN>"
```

### 2.2 `POST /device-groups`
- **Ý nghĩa**: Tạo mới một Device Group rỗng.
- **Request Body JSON**:
  - `name`: Tên group (bắt buộc).
  - `description`: Mô tả.
  - `parent_id`: Group cha tương ứng.
- **Data Source**: Ghi vào bảng `device_groups` (Ent schema).
- **Workflow / cURL**:
```bash
curl -X POST "http://localhost:8080/api/v1/device-groups" \
     -H "Content-Type: application/json" -H "Authorization: Bearer <TOKEN>" \
     -d '{"name": "Developer iOS", "description": "Nhóm cho Dev Test"}'
```

### 2.3 `POST /device-groups/:id/devices`
- **Ý nghĩa**: Đẩy hàng loạt Thiết bị vào trong một Group (để từ Group đó có thể Deploy Ứng dụng & Cấu hình tập thể).
- **Request Body JSON**: Định dạng `{"device_ids": [1, 2, 5]}`.
- **Data Source**: Cập nhật cầu nối n-n (Edge quan hệ) giữa `device_groups` và `devices` trong Ent. Hành động này sẽ gọi `AddDevices` lên đối tượng Graph.
- **Workflow / cURL**:
```bash
curl -X POST "http://localhost:8080/api/v1/device-groups/5/devices" \
     -H "Content-Type: application/json" -H "Authorization: Bearer <TOKEN>" \
     -d '{"device_ids": [12, 14, 15]}'
```

### 2.4 `DELETE /device-groups/:id/devices/:deviceId`
- **Ý nghĩa**: Bốc một Thiết bị nào đó ra khỏi cụm Group.
- **Workflow / cURL**:
```bash
curl -X DELETE "http://localhost:8080/api/v1/device-groups/5/devices/12" -H "Authorization: Bearer <TOKEN>"
```
