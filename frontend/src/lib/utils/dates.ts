import { getDateFormat } from '$lib/stores/dateFormat.svelte';

/**
 * Format a date for display (e.g., "Monday, January 20, 2025")
 */
export function formatDate(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

/**
 * Format a date short (e.g., "Jan 20")
 */
export function formatDateShort(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric'
  })
}

/**
 * Format a time (e.g., "2:00 PM" or "14:00" based on user preference)
 */
export function formatTime(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date
  const dateFormat = getDateFormat()
  return d.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: dateFormat.timeFormat === '12h'
  })
}

/**
 * Format a time range (e.g., "2:00 PM - 3:00 PM")
 */
export function formatTimeRange(start: Date | string, end: Date | string): string {
  return `${formatTime(start)} - ${formatTime(end)}`
}

/**
 * Format a date range (e.g., "Jan 20 - Jan 22")
 */
export function formatDateRange(start: Date | string, end: Date | string): string {
  return `${formatDateShort(start)} - ${formatDateShort(end)}`
}

/**
 * Get relative time (e.g., "2 hours ago", "in 3 days")
 */
export function formatRelative(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date
  const now = new Date()
  const diffMs = d.getTime() - now.getTime()
  const diffSec = Math.round(diffMs / 1000)
  const diffMin = Math.round(diffSec / 60)
  const diffHour = Math.round(diffMin / 60)
  const diffDay = Math.round(diffHour / 24)

  if (Math.abs(diffMin) < 1) return 'just now'
  if (Math.abs(diffMin) < 60) {
    return diffMin > 0 ? `in ${diffMin} min` : `${Math.abs(diffMin)} min ago`
  }
  if (Math.abs(diffHour) < 24) {
    return diffHour > 0 ? `in ${diffHour} hours` : `${Math.abs(diffHour)} hours ago`
  }
  return diffDay > 0 ? `in ${diffDay} days` : `${Math.abs(diffDay)} days ago`
}

/**
 * Check if a date is today
 */
export function isToday(date: Date | string): boolean {
  const d = typeof date === 'string' ? new Date(date) : date
  const today = new Date()
  return d.toDateString() === today.toDateString()
}

/**
 * Get ISO date string (YYYY-MM-DD) for a date
 */
export function toISODateString(date: Date): string {
  return date.toISOString().split('T')[0]
}

/**
 * Get day names ordered by week start preference from the dateFormat store
 */
export function getWeekDays(format: 'short' | 'narrow' | 'long' = 'short'): string[] {
  const dateFormat = getDateFormat()
  return dateFormat.getWeekDays(format)
}

/**
 * Get day index adjusted for week start preference (0 = first day of week)
 */
export function getDayIndex(date: Date): number {
  const dateFormat = getDateFormat()
  return dateFormat.getDayIndex(date)
}
