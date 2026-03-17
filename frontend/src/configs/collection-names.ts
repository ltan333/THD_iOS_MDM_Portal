/**
 * Mapping tên collection tiếng Anh sang tiếng Việt
 * Sử dụng cho các trang quản lý phân quyền
 */
export const collectionNameMap: Record<string, string> = {
  // Nhân sự
  employees: "Nhân viên",
  users: "Người dùng",
  roles: "Vai trò",
  policies: "Chính sách",
  permissions: "Quyền hạn",
  positions: "Chức vụ",
  departments: "Phòng ban",

  // Lương & Hợp đồng
  contracts: "Hợp đồng",
  salary_schemes: "Bậc lương",
  monthly_payrolls: "Bảng lương tháng",
  salary_requests: "Yêu cầu lương",
  deductions: "Khấu trừ",

  // Chấm công
  attendance_logs: "Nhật ký chấm công",
  attendance_shifts: "Ca chấm công",
  attendance_adjustments: "Điều chỉnh chấm công",
  monthly_employee_stats: "Thống kê nhân viên tháng",

  // Lịch làm việc
  shifts: "Ca làm việc",
  shift_types: "Loại ca làm",
  shift_position_requirements: "Yêu cầu vị trí ca làm",
  weekly_schedules: "Lịch tuần",
  schedule_assignments: "Phân ca",
  schedule_change_requests: "Yêu cầu đổi ca",
  employee_availability: "Đăng ký lịch rảnh",
  employee_availability_positions: "Vị trí ưu tiên nhân viên",

  // Thiết bị & RFID
  devices: "Thiết bị",
  device_events: "Sự kiện thiết bị",
  rfid_cards: "Thẻ RFID",

  // Thông báo
  notifications: "Thông báo",

  // Báo cáo & Phân tích
  reports: "Báo cáo",
  analytics: "Phân tích",
  analysis: "Phân tích dữ liệu",

  // Files
  files: "Tập tin",
  directus_files: "Tập tin hệ thống",

  // Directus system collections
  directus_users: "Người dùng hệ thống",
  directus_roles: "Vai trò hệ thống",
  directus_policies: "Chính sách hệ thống",
  directus_permissions: "Quyền hệ thống",
  directus_collections: "Danh sách bảng",
  directus_fields: "Trường dữ liệu",
  directus_relations: "Quan hệ dữ liệu",
  directus_activity: "Nhật ký hoạt động",
  directus_revisions: "Lịch sử chỉnh sửa",
  directus_webhooks: "Webhooks",
  directus_flows: "Quy trình tự động",
  directus_operations: "Thao tác tự động",
  directus_notifications: "Thông báo hệ thống",
  directus_dashboards: "Bảng điều khiển",
  directus_panels: "Bảng hiển thị",
  directus_shares: "Chia sẻ",
  directus_presets: "Cấu hình mặc định",
  directus_settings: "Cài đặt hệ thống",
  directus_folders: "Thư mục",
  directus_translations: "Bản dịch",
  directus_versions: "Phiên bản",
  directus_extensions: "Tiện ích mở rộng",
  directus_access: "Quyền truy cập",

  // Collections
  collections: "Danh sách bảng",
};

/**
 * Lấy tên hiển thị tiếng Việt của collection kèm tên gốc
 * Format: "Tên tiếng Việt (tên_gốc)"
 * @param collectionName - Tên collection tiếng Anh
 * @returns Tên tiếng Việt kèm tên gốc trong ngoặc
 */
export const getCollectionDisplayName = (collectionName: string): string => {
  const vietnameseName = collectionNameMap[collectionName];
  if (vietnameseName) {
    return `${vietnameseName} (${collectionName})`;
  }
  // Nếu không có mapping, trả về tên gốc
  return collectionName;
};
