# Tài liệu API: Dashboard (Thống kê tổng quan)
Base URL: `/api/v1/dashboard`

Các endpoints sau được thiết kế tối giản để phục vụ render Dữ liệu trên màn Home Dashboard của Frontend Admin mà không cần gọi từng CRUD API và tự tính toán đếm số liệu lặp vòng phía máy khách.

---

### 1. `GET /stats`
- **Ý nghĩa**: Cung cấp bức tranh toàn cảnh bằng các con số Card Widgets. (Tổng User, Tổng số thiết bị đăng ký, tổng Application...).
- **Response Payload**:
  - `total_devices`: Lệnh `Count` trên bảng `devices`.
  - `active_devices`: Các thiết bị có trạng thái đang ping.
  - `enrolled_devices`: Số máy đã Enrollment thành công và lấy token đầy đủ.
  - `total_users`: Gộp Count toàn bộ user/admin trên hệ thống.
  - `total_applications`: Đếm Entity applications.
  - `total_alerts`: Đếm Alert log rủi ro gửi về.
- **Workflow / cURL**:
```json
// Lệnh gọi: curl -X GET "http://localhost:8080/api/v1/dashboard/stats"
{
  "total_devices": 150,
  "active_devices": 112,
  "enrolled_devices": 140,
  "total_users": 5,
  "total_profiles": 10,
  "total_applications": 20,
  "total_alerts": 34
}
```

### 2. `GET /device-stats`
- **Ý nghĩa**: Thường dùng để xây các biểu đồ tròn (Pie Charts). Cung cấp sự phân mảnh theo các chiều:
  - Nền tảng (Platform: iOS, Android, macOS).
  - Trạng thái đăng ký.
  - Tỉ lệ tuân thủ chính sách bảo mật (Compliance Rules).
- **Data Source**: Cơ chế Backend sẽ gọi `GroupBy` query vào thư viện Ent: `client.Device.Query().GroupBy(device.FieldPlatform).Aggregate(ent.Count()).Scan(ctx, &v)`.
- **Workflow / cURL**:
```json
// Lệnh gọi: curl -X GET "http://localhost:8080/api/v1/dashboard/device-stats"
{
  "by_platform": {
    "ios": 90,
    "android": 60
  },
  "by_status": {
    "active": 112,
    "inactive": 38
  },
  "compliance": {
    "compliant": 100,
    "non_compliant": 50
  }
}
```

### 3. `GET /alerts-summary`
- **Ý nghĩa**: Báo cáo tổng thể tình trạng hệ thống bị báo động theo cấp bậc độ nguy hiểm.
- **Response Payload**:
  - Mức độ (`critical`, `high`, `medium`, `low`).
  - Phân loại (Type): `security` (bị Root máy IOS jailbreak), `connectivity` (mất mạng kéo dài)...
- **Data Source**: `GroupBy` trên bảng `alerts` dựa theo field `Severity` và `Type`.

### 4. `GET /charts/:type`
- **Ý nghĩa**: Xuất Vector Arrays dữ liệu dành cho Line Chart phát triển theo chu kỳ ngày/tháng/năm để nhìn biểu đồ đường lượng truy cập. Trả về format `labels` (x-axis) và `data` (y-axis).
