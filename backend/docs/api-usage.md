---
description: Cách gọi API MDM Portal một cách hệ thống và bảo mật (Exhaustive Workflow)
---

# Luồng Vận Hành API MDM Portal Toàn diện

Quy trình này hướng dẫn bạn cách sử dụng **tất cả 59 endpoints** của hệ thống để quản lý vòng đời thiết bị Apple.

### Giai đoạn 1: Thiết lập & Hạ tầng (Infrastructure)
Bạn không thể quản lý thiết bị nếu chưa hoàn thành các bước "móng" này.

1. **Xác thực**:
   - Gọi `POST /api/v1/auth/login`. Nhận `token`.
   - Mọi request tiếp theo phải có Header: `Authorization: Bearer <token>`.
   - Kiểm tra quyền hạn của mình qua `GET /api/v1/auth/me`.

2. **Kích hoạt Kênh điều khiển MDM**:
   - Gọi `POST /api/v1/mdm/pushcert`. Tải lên chứng chỉ APNs.
   - Xác nhận chứng chỉ đã hoạt động qua `GET /api/v1/mdm/pushcert`.
   - **Lưu ý**: Nếu bước này hỏng, mọi lệnh MDM ở Giai đoạn 3 sẽ bị treo.

3. **Kết nối Apple Business Manager (ABM)**:
   - Gọi `PUT /api/v1/dep/token/default` để nạp Token ABM.
   - Kiểm tra thông tin tổ chức trả về từ Apple qua `GET /api/v1/dep/proxy/default/account`.

### Giai đoạn 2: Tự động hóa & Khởi động (Automation & Onboarding)
Chuẩn bị sẵn sàng để máy "tự đăng ký" vào hệ thống.

4. **Kéo danh sách máy về Database**:
   - Gọi `POST /api/v1/dep/sync` (Đồng bộ nhanh) hoặc `POST /api/v1/dep/proxy/default/devices/sync` (Đồng bộ sâu qua Apple).
   - Xem danh sách máy đang nằm trên server Apple qua `POST /api/v1/dep/proxy/default/devices`.

5. **Thiết lập quy tắc bóc hộp (Enrollment)**:
   - Tạo Profile đăng ký tại `POST /api/v1/dep/profile`.
   - Thiết lập luật tự động gán Profile này cho máy mới qua `PUT /api/v1/dep/assigner/default`.
   - Kiểm tra lại các profile đã tạo tại `GET /api/v1/dep/profiles`.

6. **Chuẩn bị cấu hình chi tiết**:
   - Tải lên các file `.mobileconfig` (WiFi, Email, Bảo mật) tại `PUT /api/v1/nanocmd/profile/:name`.
   - Xem danh sách các file cấu hình đang sẵn sàng tại `GET /api/v1/nanocmd/profiles`.

### Giai đoạn 3: Quản trị Thực tế (Operation)
Thực hiện điều khiển thiết bị khi nhân viên đã bật máy và kết nối WiFi.

7. **Kích hoạt sự kiện theo dõi**:
   - Thiết lập để máy thông báo về server khi có thay đổi qua `PUT /api/v1/nanocmd/event/:name`.

8. **Gửi lệnh bảo mật**:
   - Khi cần Khóa máy, Xóa máy từ xa: Gọi `PUT /api/v1/mdm/enqueue/:id` với XML payload tương ứng.
   - Sử dụng các workflow tự động (như "Gia nhập miền", "Kiểm tra bảo mật") qua `POST /api/v1/nanocmd/workflow/:name/start`.

9. **Báo cáo & Kiểm kê (Inventory)**:
   - Lấy báo cáo chi tiết Serial, Pin, App đã cài... qua `GET /api/v1/nanocmd/inventory?id=<udid>`.
   - Theo dõi các kế hoạch thực thi dài hạn qua `GET /api/v1/nanocmd/cmdplan/:name`.

---

> [!CAUTION]
> Luôn sử dụng `/health` để kiểm tra tình hình server trước khi thực hiện các đợt đồng bộ danh sách thiết bị lớn (Sync).
