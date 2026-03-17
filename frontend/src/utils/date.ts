import dayjs from "dayjs";
import "dayjs/locale/vi";
import relativeTime from "dayjs/plugin/relativeTime";

dayjs.extend(relativeTime);
dayjs.locale("vi");

/**
 * Format date to display string
 */
export function formatDate(date: string | Date, format = "DD/MM/YYYY"): string {
    return dayjs(date).format(format);
}

/**
 * Format date with time
 */
export function formatDateTime(
    date: string | Date,
    format = "DD/MM/YYYY HH:mm"
): string {
    return dayjs(date).format(format);
}

/**
 * Get relative time (e.g., "2 hours ago")
 */
export function formatRelativeTime(date: string | Date): string {
    return dayjs(date).fromNow();
}

/**
 * Format time only
 */
export function formatTime(date: string | Date, format = "HH:mm"): string {
    return dayjs(date).format(format);
}

export { dayjs };
