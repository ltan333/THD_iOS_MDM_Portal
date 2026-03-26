# Toàn bộ Danh sách API Endpoint (Exhaustive & Prioritized)

Tài liệu này bao gồm **tất cả** 131 endpoint hiện có trong hệ thống, được trích xuất đầy đủ từ Router và sắp xếp theo thứ tự ưu tiên cho việc phát triển **Portal (Frontend)**.

---

## 🛑 PHẦN 1: CÁC API PHỤC VỤ PORTAL (Ưu tiên Cao)
Nhóm này bao gồm các chức năng cốt lõi Front-end cần thực hiện để hoàn thiện giao diện chính.

### 1.1 Xác thực & Tài khoản (Auth & Me)
| Chức năng | Method | Endpoint |
|:---|:---|:---|
| Login | POST | `/api/v1/auth/login` |
| Refresh Token | POST | `/api/v1/auth/refresh` |
| Logout | POST | `/api/v1/auth/logout` |
| Get My Info | GET | `/api/v1/auth/me` |
| Health Check | GET | `/health` |
| Swagger Docs | GET | `/swagger/*any` |

### 1.2 Quản trị Profile (Business Logic)
| Chức năng | Method | Endpoint |
|:---|:---|:---|
| Danh sách Profiles | GET | `/api/v1/profiles` |
| Tạo Profile | POST | `/api/v1/profiles` |
| Chi tiết Profile | GET | `/api/v1/profiles/:id` |
| Cập nhật Profile | PUT | `/api/v1/profiles/:id` |
| Xóa Profile | DELETE | `/api/v1/profiles/:id` |
| Cập nhật Trạng thái | PUT | `/api/v1/profiles/:id/status` |
| **Sửa Security Settings** | PUT | `/api/v1/profiles/:id/settings/security` |
| **Sửa Network Config** | PUT | `/api/v1/profiles/:id/settings/network` |
| **Sửa Restrictions** | PUT | `/api/v1/profiles/:id/settings/restrictions` |
| **Sửa Content Filter** | PUT | `/api/v1/profiles/:id/settings/content-filter` |
| **Sửa Compliance Rules** | PUT | `/api/v1/profiles/:id/settings/compliance` |
| Danh sách Assignment | GET | `/api/v1/profiles/:id/assignments` |
| Gán Profile | POST | `/api/v1/profiles/:id/assignments` |
| Gỡ gán Profile | DELETE | `/api/v1/profiles/:id/assignments/:assignmentId` |
| Danh sách Phiên bản | GET | `/api/v1/profiles/:id/versions` |
| Rollback Phiên bản | POST | `/api/v1/profiles/:id/versions/:versionId/rollback` |
| Deployment Status | GET | `/api/v1/profiles/:id/deployment-status` |
| Gửi lại Profile | POST | `/api/v1/profiles/:id/repush` |
| Nhân bản Profile | POST | `/api/v1/profiles/:id/duplicate` |

### 1.3 Thiết bị & Nhóm (Inventory & Groups)
| Chức năng | Method | Endpoint |
|:---|:---|:---|
| Danh sách Thiết bị | GET | `/api/v1/devices` |
| Xuất danh sách Excel | GET | `/api/v1/devices/export` |
| Chi tiết Thiết bị | GET | `/api/v1/devices/:id` |
| Khóa Thiết bị | POST | `/api/v1/devices/:id/lock` |
| Wipe Thiết bị | POST | `/api/v1/devices/:id/wipe` |
| Danh sách Nhóm | GET | `/api/v1/device-groups` |
| Tạo Nhóm | POST | `/api/v1/device-groups` |
| Chi tiết Nhóm | GET | `/api/v1/device-groups/:id` |
| Cập nhật Nhóm | PUT | `/api/v1/device-groups/:id` |
| Xóa Nhóm | DELETE | `/api/v1/device-groups/:id` |
| Thêm Thiết bị vào Nhóm | POST | `/api/v1/device-groups/:id/devices` |
| Xóa Thiết bị khỏi Nhóm | DELETE | `/api/v1/device-groups/:id/devices/:deviceId` |

### 1.4 Dashboard, Alerts & Rules
| Chức năng | Method | Endpoint |
|:---|:---|:---|
| Dashboard Stats | GET | `/api/v1/dashboard/stats` |
| Device Sync Stats | GET | `/api/v1/dashboard/device-stats` |
| Alerts Summary | GET | `/api/v1/dashboard/alerts-summary` |
| Chart Data | GET | `/api/v1/dashboard/charts/:type` |
| Danh sách Cảnh báo | GET | `/api/v1/alerts` |
| Tạo Cảnh báo thủ công | POST | `/api/v1/alerts` |
| Thống kê Cảnh báo | GET | `/api/v1/alerts/stats` |
| Chi tiết Cảnh báo | GET | `/api/v1/alerts/:id` |
| Acknowledge Alert | PUT | `/api/v1/alerts/:id/acknowledge` |
| Resolve Alert | PUT | `/api/v1/alerts/:id/resolve` |
| Bulk Resolve | POST | `/api/v1/alerts/bulk-resolve` |
| Action: Lock Device | POST | `/api/v1/alerts/:id/actions/lock` |
| Action: Wipe Device | POST | `/api/v1/alerts/:id/actions/wipe` |
| Action: Push Policy | POST | `/api/v1/alerts/:id/actions/push-policy` |
| Action: Send Message | POST | `/api/v1/alerts/:id/actions/message` |
| Danh sách Alert Rules | GET | `/api/v1/alerts/rules` |
| Tạo Rule | POST | `/api/v1/alerts/rules` |
| Chi tiết Rule | GET | `/api/v1/alerts/rules/:id` |
| Cập nhật Rule | PUT | `/api/v1/alerts/rules/:id` |
| Xóa Rule | DELETE | `/api/v1/alerts/rules/:id` |
| Bật/Tắt Rule | PUT | `/api/v1/alerts/rules/:id/toggle` |

### 1.5 Ứng dụng & Báo cáo
| Chức năng | Method | Endpoint |
|:---|:---|:---|
| CRUD Applications | GET, POST, PUT, DELETE | `/api/v1/applications` (và /:id) |
| Quản lý Application Versions | GET, POST, DELETE | `/api/v1/applications/:id/versions/...` |
| Xem trạng thái Deploy App | GET | `/api/v1/applications/:id/versions/:vId/deployments` |
| Triển khai App (Deploy) | POST | `/api/v1/applications/deployments` |
| Xuất Báo cáo Thiết bị | GET | `/api/v1/reports/devices/export` |
| Xuất Báo cáo Cảnh báo | GET | `/api/v1/reports/alerts/export` |
| Xuất Báo cáo App | GET | `/api/v1/reports/applications/export` |

### 1.6 Quản trị User, RBAC & Settings
| Chức năng | Method | Endpoint |
|:---|:---|:---|
| CRUD Users | GET, POST, PUT, DELETE | `/api/v1/users` (và /:id) |
| Danh sách Policy & Role | GET, POST, DELETE | `/api/v1/policies` / `/api/v1/roles` |
| Policy theo Role | GET | `/api/v1/policies/role/:role` |
| CRUD Mobile-Configs (Raw) | GET, POST, PUT, DELETE | `/api/v1/mobile-configs` (và /:id) |
| Lấy RAW XML Mobile-Config | GET | `/api/v1/mobile-configs/:id/xml` |
| CRUD System Settings | GET, POST, PUT, DELETE | `/api/v1/settings` (và /:key) |

---

## ⚙️ PHẦN 2: CÁC API TÍCH HỢP HẠ TẦNG (Ưu tiên thấp hơn)
Nhóm này phục vụ đồng bộ với Apple DEP và các server MDM/CMD nền tảng.

### 2.1 Apple DEP (Device Enrollment Program)
| Chức năng | Method | Endpoint |
|:---|:---|:---|
| List DEP Names | GET | `/api/v1/dep/dep_names` |
| Token PKI Ops | GET, PUT | `/api/v1/dep/tokenpki/:name` |
| Tokens Ops | GET, PUT | `/api/v1/dep/tokens/:name` |
| Config Ops | GET, PUT | `/api/v1/dep/config/:name` |
| Assigner Ops | GET, PUT | `/api/v1/dep/assigner/:name` |
| Get MAID JWT | GET | `/api/v1/dep/maidjwt/:name` |
| Get Bypass Code | GET | `/api/v1/dep/bypasscode` |
| Get DEP Version | GET | `/api/v1/dep/version` |
| **Proxy: Account Info** | GET | `/api/v1/dep/proxy/:name/account` |
| **Proxy: DEP Profile** | GET, POST | `/api/v1/dep/proxy/:name/profile` |
| **Proxy: Sync/Disown/List**| POST | `/api/v1/dep/proxy/:name/devices/...` |

### 2.2 NanoCMD Workflows
| Chức năng | Method | Endpoint |
|:---|:---|:---|
| NanoCMD Version | GET | `/api/v1/nanocmd/version` |
| Start Workflow | POST | `/api/v1/nanocmd/workflow/:name/start` |
| Event Subscriptions | GET, PUT | `/api/v1/nanocmd/event/:name` |
| FileVault Template | GET | `/api/v1/nanocmd/fvenable/profiletemplate` |
| Raw Profile Template | GET, PUT, DELETE | `/api/v1/nanocmd/profile/:name` |
| List Templates | GET | `/api/v1/nanocmd/profiles` |
| CMD Plan Ops | GET, PUT | `/api/v1/nanocmd/cmdplan/:name` |
| System Inventory | GET | `/api/v1/nanocmd/inventory` |
| **NanoCMD Webhook** | POST | `/api/v1/nanocmd/webhook` |

### 2.3 MDM Core Operations
| Chức năng | Method | Endpoint |
|:---|:---|:---|
| Push Certificate Ops | PUT, GET | `/api/v1/mdm/pushcert` |
| Trigger MDM Push | GET | `/api/v1/mdm/push/:id` |
| Enqueue Command | PUT | `/api/v1/mdm/enqueue/:id` |
| Escrow Key Unlock | POST | `/api/v1/mdm/escrowkeyunlock` |
| NanoMDM Version | GET | `/api/v1/mdm/version` |
