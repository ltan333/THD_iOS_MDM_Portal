/**
 * Employee Code Generator
 * Tự động tạo mã nhân viên theo quy tắc: NV + số thứ tự (ví dụ: NV001, NV002, ...)
 */

/**
 * Tạo mã nhân viên mới dựa trên ngày tháng năm hiện tại
 * Format: NV + YYYYMMDD + số thứ tự (ví dụ: NV20241205001)
 * 
 * @returns Mã nhân viên được tạo tự động
 */
export const generateEmployeeCode = (): string => {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, '0');
  const day = String(now.getDate()).padStart(2, '0');
  const timestamp = Math.floor(Math.random() * 10000); // Random 0-9999
  const sequence = String(timestamp).padStart(4, '0');
  
  return `NV${year}${month}${day}${sequence}`;
};

/**
 * Alternative: Tạo mã nhân viên đơn giản hơn
 * Format: NV + Timestamp + Random (ví dụ: NV1701787200123)
 * 
 * @returns Mã nhân viên được tạo tự động
 */
export const generateSimpleEmployeeCode = (): string => {
  const timestamp = Date.now();
  const random = Math.floor(Math.random() * 1000);
  return `NV${timestamp}${random}`;
};

/**
 * Tạo mã nhân viên dạng tuần tự
 * Format: NV + số thứ tự (ví dụ: NV0001, NV0002, ...)
 * Lưu ý: Cách này cần lấy số thứ tự từ BE
 * 
 * @param sequence Số thứ tự từ BE
 * @returns Mã nhân viên với định dạng NV + số thứ tự
 */
export const generateSequentialEmployeeCode = (sequence: number): string => {
  const paddedSequence = String(sequence).padStart(4, '0');
  return `NV${paddedSequence}`;
};

export default generateEmployeeCode;
