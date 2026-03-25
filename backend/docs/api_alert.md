# Tài liệu API: Alerts & Rules (Cảnh báo và Định tuyến Sự cố)
Base URL: `/api/v1/alerts`

Bao gồm danh sách sự cố vi phạm chính sách của thiết bị (Alert Logs), bộ hành động nhanh (Lock, Wipe, Notify) và quản trị tự động bằng Rules (Automation).

---

## 1. Màn hình Cảnh Báo (Alert Logs)

### 1.1 `GET /alerts`
- **Ý nghĩa**: Xem log sự kiện cần giải quyết.
- **Query Params**: `severity` (mức độ: critical/high/medium), `status` (open, resolved, acknowledged).
- **Data Source**: Table `alerts`.

### 1.2 `GET /alerts/stats`
- **Ý nghĩa**: Thống kê số lượng theo nhóm (Alert Analytics). Thường dùng cho Widget Dashboard.

### 1.3 `PUT /alerts/:id/acknowledge` & `PUT /alerts/:id/resolve`
- **Ý nghĩa**: Đánh dấu Alert là đang duyệt (Acknowledge) hoặc đã giải quyết xong (Resolve).
- **Workflow / cURL**:
```bash
curl -X PUT "http://localhost:8080/api/v1/alerts/1/resolve" -H "Authorization: Bearer <TOKEN>"
```

---

## 2. Quick Actions (Tác nghiệp tức thì theo sự cố)

Là nhóm API tích hợp sâu vào nền tảng thiết bị, kích hoạt bảo vệ dữ liệu ngay từ Portal khi Admin đọc một Alert.

### 2.1 `POST /alerts/:id/actions/lock`
- **Ý nghĩa**: Khóa máy tính/điện thoại gây ra cảnh báo `id` này.
- **Data Source**: Gọi trực tiếp Module MDM Push Command (`DeviceLock`), tự động gắn log này vào Alert.
- **Workflow / cURL**:
```bash
curl -X POST "http://localhost:8080/api/v1/alerts/12/actions/lock" -H "Authorization: Bearer <TOKEN>"
```

### 2.2 `POST /alerts/:id/actions/wipe`
- **Ý nghĩa**: Reset Factory máy thiết bị vi phạm cực độ.

### 2.3 `POST /alerts/:id/actions/push-policy`
- **Ý nghĩa**: Áp đặt một Policy/Profile khẩn cấp (vd chặn mạng). Cần `policy_id` trong payload.

---

## 3. Automation Alert Rules (Quy tắc tự động)

API Route: `/api/v1/alerts/rules`
- **Ý nghĩa**: CRUD logic tạo rule phát cảnh báo. Ví dụ: "Nếu thiết bị báo cáo OSVersion < 15.0 -> Gửi Alert Critical tạo Action Lock Device"
- Dữ liệu `condition` và `actions` được lưu trữ dạng stringified JSON vào table `alert_rules`.

### 3.1 `POST /alerts/rules`
- **Workflow / cURL**:
```json
// POST "http://localhost:8080/api/v1/alerts/rules"
{
  "name": "Jailbreak Detected Security Flow",
  "condition": { "field": "is_jailbroken", "operator": "eq", "value": true },
  "actions": { "severity": "critical", "auto_action": "lock" }
}
```
