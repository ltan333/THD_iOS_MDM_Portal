# Hướng dẫn Luồng vận hành MDM qua API

Tài liệu này hướng dẫn cách kết nối các API để hoàn thiện một chu kỳ quản lý thiết bị iOS, từ lúc nhập máy đến khi điều khiển từ xa.

---

## 🚀 1. Luồng vận hành Tổng quát
Quy trình chuyên nghiệp (Zero-Touch) bao gồm:
1. **Sync DEP**: Đồng bộ danh sách máy từ Apple.
2. **Assign Group**: Gán thiết bị (Serial Number) vào một Nhóm (Group).
3. **Assign Profile**: Gán các chính sách (Wifi, Passcode...) vào Nhóm đó.
4. **Enrollment**: Người dùng bật máy, máy tự cấu hình dựa trên Nhóm đã gán.
5. **Management**: Khóa/Xóa máy khi cần.

---

## 🛠 2. Chi tiết các bước và API

### Bước 1: Đồng bộ thiết bị từ Apple (DEP Sync)
Hệ thống sẽ lấy danh sách Serial Number về và lưu ở trạng thái `Pending`.

*   **Endpoint:** `POST /api/v1/dep/proxy/:name/devices/sync`
*   **Body:** `{}` (Trống)
*   **Kết quả:** Các máy mới sẽ xuất hiện trong `GET /api/v1/devices` với ID dạng `dep-SERIALNUMBER`.

### Bước 2: Tạo Nhóm và Gán máy vào Nhóm
Việc gán nhóm giúp tự động hóa cấu hình sau này.

*   **Tạo Nhóm:** `POST /api/v1/device-groups`
    ```json
    {
      "name": "Phòng Kỹ Thuật",
      "description": "Nhóm thiết bị cho nhân viên IT"
    }
    ```
*   **Gán thiết bị:** `POST /api/v1/device-groups/:id/devices`
    ```json
    {
      "device_ids": ["dep-C7GGH123456"]
    }
    ```

### Bước 3: Tạo và Gán Profile cho Nhóm
Profile chứa các cài đặt bảo mật.

*   **Tạo Profile:** `POST /api/v1/profiles`
    ```json
    {
      "name": "Chính sách Bảo mật IT",
      "platform": "ios",
      "scope": "system",
      "security_settings": {
        "passcode_required": true,
        "min_length": 6
      },
      "restrictions": {
        "allow_camera": false,
        "allow_screenshot": true
      }
    }
    ```
*   **Gán cho Nhóm:** `POST /api/v1/profiles/:id/assignments`
    ```json
    {
      "profile_id": 10,
      "target_type": "group",
      "group_id": 1,
      "schedule_type": "immediate"
    }
    ```

### Bước 4: Đăng ký thiết bị (Enrollment Webhook)
Đây là bước **tự động** khi NanoCMD gọi về Backend. Bạn có thể giả lập để test.

*   **Endpoint:** `POST /api/v1/nanocmd/webhook`
*   **Body (Giả lập TokenUpdate):**
    ```json
    {
      "topic": "mdm.TokenUpdate",
      "checkin_event": {
        "udid": "REAL-UDID-789-XYZ",
        "serial_number": "C7GGH123456",
        "model": "iPhone 15 Pro",
        "os_version": "17.4"
      }
    }
    ```
*   **Hành động tự động:** 
    1. Hệ thống tìm thấy máy `dep-C7GGH123456` đang ở nhóm "Phòng Kỹ Thuật".
    2. Chuyển đổi ID sang `REAL-UDID-789-XYZ`.
    3. Tự động đẩy Profile "Chính sách Bảo mật IT" xuống máy.

### Bước 5: Điều khiển từ xa (Remote Control)
Khi máy đã ở trạng thái `Active`, bạn có thể ra lệnh.

*   **Khóa máy:** `POST /api/v1/devices/:udid/lock`
*   **Xóa máy:** `POST /api/v1/devices/:udid/wipe`

---

## 📋 3. Một số API hữu ích khác

| Chức năng | Method | Endpoint |
|:---|:---|:---|
| Đẩy lại Profile (Repush) | POST | `/api/v1/profiles/:id/repush` |
| Kiểm tra trạng thái Deploy | GET | `/api/v1/profiles/:id/deployment-status` |
| Xem XML thực tế | GET | `/api/v1/mobile-configs/:id/xml` |

---

> [!TIP]
> Luôn kiểm tra **Log** của Backend khi thực hiện Bước 4 để xác nhận luồng "Identity Migration" và "Auto Deploy" hoạt động chính xác.
