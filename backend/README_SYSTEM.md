# THD iOS MDM Portal - Tài liệu Hệ thốn
## 1. Kiến trúc Hệ thống

Hệ thống bao gồm các thành phần chính sau:

*   **Portal Backend (Go)**: Xử lý logic nghiệp vụ, quản lý người dùng và API Gateway.
*   **Database (PostgreSQL)**: Lưu trữ dữ liệu người dùng, chính sách (policy), DEP token và cấu hình APNs.
*   **NanoMDM**: Server MDM tối giản cho các thiết bị Apple (điều phối các lệnh MDM cấp thấp).
*   **NanoCMD**: Engine điều khiển dựa trên API để tương tác với NanoMDM cho các workflow và profile.
*   **Apple Services**: APNs (Thông báo đẩy) và DEP (Chương trình đăng ký thiết bị).

### Luồng Giao tiếp:
`Người dùng` -> `Portal Backend` -> `NanoCMD` -> `NanoMDM` -> `Thiết bị Apple`

---

## 2. Cài đặt & Cấu hình

### Yêu cầu hệ thống:
*   Go 1.23+
*   PostgreSQL 15+
*   `swag` (để tạo tài liệu API) và `atlas` (để quản lý migration)

### Bắt đầu nhanh:
1.  **Môi trường**: Sao chép file `.env.example` thành `.env` và cấu hình thông tin database của bạn.
2.  **Dependencies**: Chạy `make deps` để cài đặt các công cụ cần thiết.
3.  **Database**: Chạy `make migrate/apply` để khởi tạo cấu trúc database.
4.  **Khởi chạy**: Chạy `make run` để bắt đầu server tại `http://localhost:8000`.
5.  **Tài liệu API**: Truy cập tài liệu API tương tác tại `http://localhost:8000/swagger/index.html`.

---

## 3. Các Luồng Chính & Ví dụ

### A. Xác thực (Authentication)
Tất cả các API được bảo vệ đều yêu cầu Bearer JWT token trong header `Authorization`.

**1. Đăng nhập:**
```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "yourpassword"}'
```
*   **Kết quả**: `{"access_token": "...", "refresh_token": "..."}`

---

### B. Cấu hình MDM (APNs)
Trước khi quản lý thiết bị, bạn cần tải lên chứng chỉ APNs.

**1. Tải lên Push Certificate:**
```bash
curl -X POST http://localhost:8000/api/v1/mdm/pushcert \
  -H "Authorization: Bearer <TOKEN>" \
  -F "cert=@/path/to/push_cert.pem"
```

---

### C. Đăng ký Thiết bị (DEP)
DEP giúp tự động hóa quá trình thiết lập thiết bị.

**1. Tải lên DEP Token:**
```bash
curl -X PUT http://localhost:8000/api/v1/dep/tokenpki/my-org \
  -H "Authorization: Bearer <TOKEN>" \
  -F "token=@/path/to/dep_token.p7m"
```

**2. Đồng bộ Thiết bị:**
```bash
curl -X POST http://localhost:8000/api/v1/dep/sync \
  -H "Authorization: Bearer <TOKEN>"
```

---

### D. Quản lý Lệnh (NanoCMD)
Sử dụng NanoCMD để thực hiện các tác vụ trên thiết bị đã đăng ký.

**1. Kiểm tra Phiên bản:**
```bash
curl -X GET http://localhost:8000/api/v1/nanocmd/version \
  -H "Authorization: Bearer <TOKEN>"
```

**2. Bắt đầu một Workflow (ví dụ: Setup):**
```bash
curl -X POST "http://localhost:8000/api/v1/nanocmd/workflow/setup/start?id=DEVICE-SERIAL-123" \
  -H "Authorization: Bearer <TOKEN>"
```

---

## 4. Danh sách API (Full API Reference)

Hệ thống cung cấp các endpoint sau (tất cả đều có tiền tố `/api` trừ khi có ghi chú khác):

### Xác thực (`/auth`)
| Phương thức | Đường dẫn | Mô tả | Yêu cầu Token |
| :--- | :--- | :--- | :--- |
| POST | `/login` | Đăng nhập người dùng (nhận JWT) | Không |
| POST | `/logout` | Đăng xuất người dùng | Có |
| GET | `/me` | Lấy thông tin người dùng hiện tại | Có |

### Quản lý Người dùng (`/users`)
| Phương thức | Đường dẫn | Mô tả | Yêu cầu Token |
| :--- | :--- | :--- | :--- |
| GET | `/` | Liệt kê tất cả người dùng | Có |
| GET | `/:id` | Xem chi tiết người dùng theo ID | Có |
| POST | `/` | Tạo người dùng mới | Có |
| PUT | `/:id` | Cập nhật thông tin người dùng | Có |
| DELETE | `/:id` | Xóa (tạm thời) người dùng | Có |

### MDM Services (`/v1/mdm`)
| Phương thức | Đường dẫn | Mô tả | Yêu cầu Token |
| :--- | :--- | :--- | :--- |
| POST | `/pushcert` | Tải lên/Cập nhật chứng chỉ APNs | Có |
| GET | `/pushcert` | Lấy cấu hình APNs hiện tại | Có |

### DEP Services (`/v1/dep`)
| Phương thức | Đường dẫn | Mô tả | Yêu cầu Token |
| :--- | :--- | :--- | :--- |
| PUT | `/tokenpki/:name` | Tải lên DEP token (.p7m) | Có |
| GET | `/tokens/:name` | Xem chi tiết DEP token | Có |
| POST | `/sync` | Đồng bộ thiết bị từ Apple DEP | Có |
| POST | `/profiles` | Định nghĩa một DEP profile | Có |
| GET | `/profiles/:uuid` | Xem chi tiết DEP profile | Có |
| POST | `/devices/disown` | Gỡ bỏ thiết bị khỏi quản lý DEP | Có |

### NanoCMD Integration (`/v1/nanocmd`)
| Phương thức | Đường dẫn | Mô tả | Yêu cầu Token |
| :--- | :--- | :--- | :--- |
| GET | `/version` | Lấy phiên bản server NanoCMD | Có |
| POST | `/workflow/:name/start` | Bắt đầu một workflow cụ thể | Có |
| GET | `/inventory` | Lấy thông tin phần cứng/phần mềm thiết bị | Có |
| GET | `/profile/:name` | Lấy nội dung cấu hình (XML) | Có |
| PUT | `/profile/:name` | Tải lên cấu hình mới (XML) | Có |
| DELETE | `/profile/:name` | Xóa cấu hình | Có |
| GET | `/cmdplan/:name` | Xem kế hoạch lệnh cho thiết bị | Có |
| POST | `/nanocmd/webhook`* | Webhook công khai nhận sự kiện từ NanoCMD | Không |

---

## 5. Các Lệnh Phát triển (Makefile)

Sử dụng `make` để thực hiện các tác vụ nhanh:

| Lệnh | Mô tả |
| :--- | :--- |
| `make run` | Khởi chạy server |
| `make dev` | Chạy server với chế độ hot-reload (tự tải lại khi code thay đổi) |
| `make swagger` | Cập nhật tài liệu API tương tác |
| `make lint` | Kiểm tra chất lượng mã nguồn |
| `make deps` | Cài đặt các thư viện/công cụ cần thiết |
| `make migrate/apply` | Cập nhật cấu trúc database |
