import { browser } from '$app/environment';

export type WeekStartDay = 'sunday' | 'monday';
export type TimeFormat = '12h' | '24h';

export interface DateFormatSettings {
  weekStartDay: WeekStartDay;
  timeFormat: TimeFormat;
}

const STORAGE_KEY = 'dateFormat';

// Default settings - will be overridden by browser detection
let weekStartDay = $state<WeekStartDay>('sunday');
let timeFormat = $state<TimeFormat>('12h');
let initialized = $state(false);

function detectBrowserPreferences(): DateFormatSettings {
  if (!browser) {
    return { weekStartDay: 'sunday', timeFormat: '12h' };
  }

  // Detect time format preference from locale
  const locale = navigator.language || 'en-US';
  const testDate = new Date(2000, 0, 1, 13, 0);
  const formatted = testDate.toLocaleTimeString(locale, { hour: 'numeric' });
  const is24h = !formatted.includes('PM') && !formatted.includes('AM');

  // Detect week start day preference
  // Most locales use Monday, US/Canada/Japan/etc use Sunday
  const sundayLocales = ['en-US', 'en-CA', 'ja-JP', 'ko-KR', 'zh-TW', 'he-IL'];
  const isSundayStart = sundayLocales.some(l => locale.startsWith(l.split('-')[0]) && locale.includes(l.split('-')[1]))
    || locale === 'en-US' || locale.startsWith('en-US');

  return {
    weekStartDay: isSundayStart ? 'sunday' : 'monday',
    timeFormat: is24h ? '24h' : '12h'
  };
}

function loadFromStorage(): DateFormatSettings | null {
  if (!browser) return null;
  const stored = localStorage.getItem(STORAGE_KEY);
  if (!stored) return null;
  try {
    return JSON.parse(stored) as DateFormatSettings;
  } catch {
    return null;
  }
}

function saveToStorage(settings: DateFormatSettings) {
  if (!browser) return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
}

export function getDateFormat() {
  return {
    get weekStartDay() { return weekStartDay; },
    get timeFormat() { return timeFormat; },
    get initialized() { return initialized; },

    init() {
      if (initialized) return;
      const stored = loadFromStorage();
      if (stored) {
        weekStartDay = stored.weekStartDay;
        timeFormat = stored.timeFormat;
      } else {
        const detected = detectBrowserPreferences();
        weekStartDay = detected.weekStartDay;
        timeFormat = detected.timeFormat;
      }
      initialized = true;
    },

    setWeekStartDay(day: WeekStartDay) {
      weekStartDay = day;
      saveToStorage({ weekStartDay, timeFormat });
    },

    setTimeFormat(format: TimeFormat) {
      timeFormat = format;
      saveToStorage({ weekStartDay, timeFormat });
    },

    reset() {
      const detected = detectBrowserPreferences();
      weekStartDay = detected.weekStartDay;
      timeFormat = detected.timeFormat;
      if (browser) {
        localStorage.removeItem(STORAGE_KEY);
      }
    },

    // Get day names ordered by week start preference
    getWeekDays(format: 'short' | 'narrow' | 'long' = 'short'): string[] {
      const days = [];
      const baseDate = new Date(2024, 0, 7); // Known Sunday
      const start = weekStartDay === 'monday' ? 1 : 0;

      for (let i = 0; i < 7; i++) {
        const d = new Date(baseDate);
        d.setDate(d.getDate() + ((start + i) % 7));
        days.push(d.toLocaleDateString('en-US', { weekday: format }));
      }
      return days;
    },

    // Get day index adjusted for week start (0 = first day of week)
    getDayIndex(date: Date): number {
      const jsDay = date.getDay(); // 0=Sun, 6=Sat
      if (weekStartDay === 'monday') {
        return jsDay === 0 ? 6 : jsDay - 1;
      }
      return jsDay;
    },

    // Format time according to preference
    formatTime(date: Date | string): string {
      const d = typeof date === 'string' ? new Date(date) : date;
      return d.toLocaleTimeString('en-US', {
        hour: 'numeric',
        minute: '2-digit',
        hour12: timeFormat === '12h'
      });
    },

    // Format date (weekday, month day)
    formatDate(date: Date | string): string {
      const d = typeof date === 'string' ? new Date(date) : date;
      return d.toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
    },

    // Format date short (Jan 20)
    formatDateShort(date: Date | string): string {
      const d = typeof date === 'string' ? new Date(date) : date;
      return d.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric'
      });
    },

    // Format time range
    formatTimeRange(start: Date | string, end: Date | string): string {
      return `${this.formatTime(start)} - ${this.formatTime(end)}`;
    },

    // Format date and time together
    formatDateTime(date: Date | string): string {
      const d = typeof date === 'string' ? new Date(date) : date;
      const dateStr = d.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric'
      });
      return `${dateStr}, ${this.formatTime(d)}`;
    },
  };
}
