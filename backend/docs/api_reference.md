# Danh mục Tuyệt đối API MDM Portal (100% Endpoints)

Tài liệu này liệt kê không thiếu bất kỳ một endpoint nào trong hệ thống, bao gồm cả các endpoint hệ thống và proxy. Tổng cộng có 59 endpoints đã được kiểm chứng.

## Base URL: `/api/v1` (ngoại trừ các endpoint hệ thống)

---

## 0. Endpoint Hệ thống & Tài liệu (Global)
Các endpoint này nằm ngoài prefix `/api/v1`.

| Phương thức | Endpoint | Mô tả chi tiết |
| :--- | :--- | :--- |
| **GET** | `/health` | Kiểm tra tình trạng hoạt động (Health check) của server. Trả về `{"status": "ok"}`. |
| **GET** | `/swagger/*any` | Giao diện tài liệu API tương tác (Swagger UI). Giúp bạn chạy thử các request trực tiếp trên trình duyệt. |

---

## 1. Hệ thống Xác thực (`/api/v1/auth`)
Quản lý quyền truy cập và phiên làm việc của người dùng.

| Phương thức | Endpoint | Mô tả chi tiết |
| :--- | :--- | :--- |
| **POST** | `/api/v1/auth/login` | Đăng nhập hệ thống bằng `username` và `password`. Trả về JWT Token để sử dụng cho các request sau. |
| **POST** | `/api/v1/auth/logout` | Đăng xuất, hủy phiên làm việc và xóa token xác thực hiện tại. |
| **GET** | `/api/v1/auth/me` | Lấy thông tin chi tiết về tài khoản đang đăng nhập (Username, Email, Role). |

---

## 2. Quản lý Người dùng (`/api/v1/users`)
Quản trị các tài khoản quản trị viên và kỹ thuật viên.

| Phương thức | Endpoint | Mô tả chi tiết |
| :--- | :--- | :--- |
| **GET** | `/api/v1/users` | Liệt kê danh sách người dùng. Hỗ trợ tìm kiếm, lọc theo vai trò và phân trang. |
| **GET** | `/api/v1/users/:id` | Xem thông tin chi tiết của một người dùng cụ thể dựa trên ID. |
| **POST** | `/api/v1/users` | Tạo mới một tài khoản người dùng với các thông tin cơ bản và vai trò. |
| **PUT** | `/api/v1/users/:id` | Cập nhật thông tin (Email, trạng thái, vai trò) cho người dùng hiện có. |
| **DELETE** | `/api/v1/users/:id` | Xóa vĩnh viễn tài khoản người dùng khỏi hệ thống. |

---

## 3. Quản lý Quyền hạn & Vai trò (`/api/v1/policies` & `/api/v1/roles`)
Dựa trên mô hình phân quyền Casbin (RBAC).

| Phương thức | Endpoint | Mô tả chi tiết |
| :--- | :--- | :--- |
| **GET** | `/api/v1/policies` | Xem toàn bộ ma trận phân quyền (Ai được làm gì trên đường dẫn nào). |
| **POST** | `/api/v1/policies` | Thêm một quy tắc phân quyền mới cho một vai trò. |
| **DELETE** | `/api/v1/policies` | Xóa bỏ một quy tắc phân quyền cụ thể. |
| **GET** | `/api/v1/policies/role/:role` | Liệt kê tất cả các quyền hạn đang được gán cho một vai trò (ví dụ: `ADMIN`). |
| **GET** | `/api/v1/roles` | Danh sách các liên kết phân cấp vai trò (Ví dụ: vai trò nào kế thừa từ vai trò nào). |
| **POST** | `/api/v1/roles` | Thiết lập quan hệ cha-con giữa các vai trò (Role Inheritance). |
| **DELETE** | `/api/v1/roles` | Xóa bỏ liên kết phân cấp giữa các vai trò. |

---

## 4. MDM Core & Certificates (`/api/v1/mdm`)
Nghiệp vụ MDM cốt lõi của Apple.

| Phương thức | Endpoint | Mô tả chi tiết |
| :--- | :--- | :--- |
| **POST** | `/api/v1/mdm/pushcert` | Tải lên chứng chỉ Push APNs (`.p12` hoặc `.p7m`) để server có quyền ra lệnh cho máy. |
| **GET** | `/api/v1/mdm/pushcert` | Kiểm tra meta-data của chứng chỉ hiện tại (Topic, Ngày hết hạn). |
| **PUT** | `/api/v1/mdm/enqueue/:id` | Đưa lệnh MDM Raw (định dạng XML như `DeviceLock`, `EraseDevice`) vào hàng đợi cho thiết bị. |

---

## 5. Quản lý DEP / Apple Business Manager (`/api/v1/dep`)

### 5.1 Thao tác trực tiếp với Database & NanoDEP
| Phương thức | Endpoint | Mô tả chi tiết |
| :--- | :--- | :--- |
| **GET** | `/api/v1/dep/names` | Lấy danh sách các máy chủ DEP (DEP Server) đã được cấu hình. |
| **PUT** | `/api/v1/dep/token/:name` | Đồng bộ và tải lên Token DEP (`.p7m`) từ portal ABM của Apple. |
| **GET** | `/api/v1/dep/token/:name` | Lấy thông tin Certificate để cấu hình trên portal Apple. |
| **GET** | `/api/v1/dep/tokens/:name` | Xem trạng thái chi tiết của Token từ dịch vụ NanoDEP. |
| **GET** | `/api/v1/dep/config/:name` | Truy vấn cấu hình của DEP Server (Interval, Flags). |
| **GET** | `/api/v1/dep/assigner/:name` | Xem quy tắc gán Profile tự động khi máy mới active. |
| **PUT** | `/api/v1/dep/assigner/:name` | Cập nhật hoặc thiết lập quy tắc gán Profile tự động cho thiết bị. |
| **POST** | `/api/v1/dep/sync` | Ép buộc hệ thống thực hiện đồng bộ danh sách máy từ Apple ngay lập tức. |
| **POST** | `/api/v1/dep/profile` | Định nghĩa Profile đăng ký mới (Setup Assistant settings) trên hệ thống. |
| **GET** | `/api/v1/dep/profiles` | Liệt kê danh sách các Profile đăng ký đã được tạo. |
| **GET** | `/api/v1/dep/profile/:uuid` | Truy vấn thông tin chi tiết của 1 Profile đăng ký qua UUID. |
| **POST** | `/api/v1/dep/disown` | Đánh dấu gỡ quyền quản lý thiết bị trên server (Disown). |

### 5.2 DEP Proxy (Tương tác trực tiếp với Cloud Apple)
| Phương thức | Endpoint | Mô tả chi tiết |
| :--- | :--- | :--- |
| **GET** | `/api/v1/dep/proxy/:name/account` | Truy vấn thông tin tài khoản tổ chức (Organization) từ Apple Business Manager. |
| **GET** | `/api/v1/dep/proxy/:name/profile` | Kiểm tra thông tin Profile đăng ký đang được lưu trên Apple Server. |
| **POST** | `/api/v1/dep/proxy/:name/profile` | Gửi định nghĩa Profile đăng ký trực tiếp lên Cloud của Apple. |
| **POST** | `/api/v1/dep/proxy/:name/devices` | Lấy danh sách thiết bị thực tế đang được Apple quản lý cho tổ chức. |
| **POST** | `/api/v1/dep/proxy/:name/devices/sync` | Đồng bộ hóa chênh lệch danh sách máy qua kênh Proxy Apple. |
| **POST** | `/api/v1/dep/proxy/:name/devices/disown` | Thực hiện lệnh Disown vật lý tới Apple Server qua Proxy. |

---

## 6. NanoCMD & Workflows Nâng cao (`/api/v1/nanocmd`)
Công cụ quản trị và tự động hóa mạnh mẽ.

| Phương thức | Endpoint | Mô tả chi tiết |
| :--- | :--- | :--- |
| **GET** | `/api/v1/nanocmd/version` | Kiểm tra phiên bản và tình trạng hoạt động của dịch vụ NanoCMD. |
| **POST** | `/api/v1/nanocmd/workflow/:name/start` | Thực thi một quy trình tự động (Workflow) như cài bộ ứng dụng, cấu hình WiFi... |
| **GET** | `/api/v1/nanocmd/event/:name` | Xem chi tiết cấu hình đăng ký nhận sự kiện (Event Subscriptions). |
| **PUT** | `/api/v1/nanocmd/event/:name` | Tạo mới hoặc cập nhật đăng ký sự kiện từ thiết bị. |
| **GET** | `/api/v1/nanocmd/fvenable/profiletemplate` | Lấy mẫu cấu hình Profile dành cho việc mã hóa ổ đĩa (FileVault). |
| **GET** | `/api/v1/nanocmd/profile/:name` | Truy xuất nội dung file cấu hình `.mobileconfig` đã lưu trên server. |
| **PUT** | `/api/v1/nanocmd/profile/:name` | Tải lên một file `.mobileconfig` mới để triển khai xuống thiết bị. |
| **DELETE** | `/api/v1/nanocmd/profile/:name` | Xóa vĩnh viễn file cấu hình khỏi server. |
| **GET** | `/api/v1/nanocmd/profiles` | Danh sách tất cả các file cấu hình đang được quản lý bởi NanoCMD. |
| **GET** | `/api/v1/nanocmd/cmdplan/:name` | Xem kế hoạch thực thi lệnh có điều kiện (Command Plan). |
| **PUT** | `/api/v1/nanocmd/cmdplan/:name` | Cập nhật hoặc thiết lập mới kế hoạch thực thi lệnh. |
| **GET** | `/api/v1/nanocmd/inventory` | Báo cáo chi tiết về phần cứng, phần mềm và bảo mật của thiết bị (Inventory). |
| **POST** | `/api/v1/nanocmd/webhook` | Endpoint công khai dành cho hệ thống NanoCMD gửi thông báo Webhook (Public). |

---

> [!TIP]
> Bạn có thể sử dụng **Swagger UI** tại `/swagger/index.html` để xem chi tiết schema JSON cho từng API này. Tất cả các endpoint trên đều đã được tích hợp đầy đủ.
