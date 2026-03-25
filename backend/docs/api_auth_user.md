# Tài liệu API: System Administration (Quản trị Hệ thống)

Cụm API phục vụ việc đăng nhập, quản lý Admin, và phân quyền Casbin Roles/Policies bọc bên ngoài lớp Business API.

---

## 1. Authentication (`/api/v1/auth`)

### 1.1 `POST /auth/login`
- **Ý nghĩa**: Đăng nhập hệ thống (bằng Username/Password).
- **Request Body**: `{"username": "admin", "password": "123"}`
- **Response**: Trả về Access Token định dạng JWT.
- **Workflow**: Nhận được chuỗi JWT, Frontend cần lưu vào LocalStorage và attach vào `Authorization: Bearer <TOKEN>` cho các API phía sau.

### 1.2 `GET /auth/me`
- **Ý nghĩa**: Trả về Profile chi tiết của User đang được auth bởi token (Role hiện tại là gì).

### 1.3 `POST /auth/logout`
- **Ý nghĩa**: Xóa Token hiện tại (Session invalidate) và lưu log.

---

## 2. Users Management (`/api/v1/users`)

- API CRUD Quản trị cấu trúc bảng `Users` trên DB.
- Bao gồm:
  - `GET /users` (Danh sách thành viên quản trị trị Portal)
  - `POST /users` (Thêm Admin mới)
  - `PUT /users/:id` (Cập nhật Password/Thông tin liên hệ)
  - `DELETE /users/:id` (Khóa tài khoản admin)

---

## 3. RBAC Policies & Roles (`/api/v1/policies` & `/api/v1/roles`)

Phân quyền chặt chẽ sử dụng engine **Casbin**.

### 3.1 `GET /roles` & `POST /roles`
- **Ý nghĩa**: Quản trị cấu trúc tên role. Vd: `SuperAdmin`, `Readonly`, `Helpdesk`.

### 3.2 `GET /policies/role/:role`
- **Ý nghĩa**: Liệt kê tất cả các phân quyền (vd: Cho phép `GET /api/v1/devices`, nhưng cấm lệnh `POST /api/v1/devices/*/wipe`).
- **Data Source**: Phân quyền này không lưu vào RDBMS thuần mà ghi thẳng vào policy definitions của Enforcer.
