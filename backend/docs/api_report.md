# Tài liệu API: Reports (Xuất báo cáo)
Base URL: `/api/v1/reports`

Phase này phục vụ nhu cầu tạo file Export, trích xuất dữ liệu gốc sang những định dạng văn bản Data (như *.csv) để Admin gửi báo cáo hoặc Pivot Table trong MS Excel.

Tất cả Endpoint báo cáo đều cấu hình response HEADER là `Content-Disposition: attachment; filename="..."` kết hợp với MIME Type `text/csv`.

---

### 1. `GET /devices/export`
- **Ý nghĩa**: Tải xuống file CSV chứa toàn bộ thiết bị.
- **Query Params**: Hỗ trợ Filter, dùng `search` để truyền keyword nếu muốn xuất báo cáo lọc theo một cụm từ. 
- **Data Source**: Gọi trực tiếp service Device để quét thông tin từ bảng `devices`. Dữ liệu duyệt qua vòng lặp, mỗi máy tương đương một hàng định dạng chuỗi: `ID, Name, Serial Number, Platform, Model, OS Version, Status, Compliance Status, Is Enrolled, MAC Address, IP Address, Created At`.
- **Workflow / cURL**:
```bash
curl -X GET "http://localhost:8080/api/v1/reports/devices/export?search=Admin" \
     -H "Authorization: Bearer <TOKEN>" \
     -o "/path/to/downloaded_devices.csv"
```
*(Option `-o` trong cURL chỉ định nới lưu file dữ liệu byte-stream tải về)*

### 2. `GET /alerts/export`
- **Ý nghĩa**: Báo cáo thống kê những biến cố, cảnh báo bảo mật nguy hiểm lưu trên DB. Các cột gồm Severity, Type, ID, và thời điểm Resolved lúc nào.
- **Data Source**: Query database trên table `alerts` và ghi đè string thông qua package lõi `encoding/csv`.
- **Workflow / cURL**:
```bash
curl -X GET "http://localhost:8080/api/v1/reports/alerts/export" \
     -H "Authorization: Bearer <TOKEN>" \
     -o alerts_log.csv
```

### 3. `GET /applications/export`
- **Ý nghĩa**: Báo cáo danh sách mọi Ứng dụng đã được ghi danh trên Portal. (Không phân biệt Store app và Enterprise App).
- **Data Source**: Query theo entity `applications`.
- **Workflow / cURL**:
```bash
curl -X GET "http://localhost:8080/api/v1/reports/applications/export" \
     -H "Authorization: Bearer <TOKEN>" \
     -o applications.csv
```
