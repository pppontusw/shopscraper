import timeAgo from './timeAgo';
import { advanceTo, clear } from 'jest-date-mock';

describe('timeAgo Function', () => {
  beforeEach(() => {
    // Set a specific point in time (e.g., Jan 1, 2020 00:00:00 GMT+0000)
    advanceTo(new Date(2020, 0, 1, 0, 0, 0));
  });

  afterEach(() => {
    // Clear the mocked date after each test to avoid leakage between tests
    clear();
  });

  test('returns "just now" for times less than a minute ago', () => {
    const now = new Date();
    expect(timeAgo(now)).toBe('just now');
  });

  test('returns correct time ago for various units', () => {
    const now = new Date();
    const oneMinuteAgo = new Date(now.getTime() - 60 * 1000);
    const oneHourAgo = new Date(now.getTime() - 60 * 60 * 1000);
    const oneDayAgo = new Date(now.getTime() - 24 * 60 * 60 * 1000);
    const oneWeekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
    const oneMonthAgo = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
    const oneYearAgo = new Date(now.getTime() - 365 * 24 * 60 * 60 * 1000);

    expect(timeAgo(oneMinuteAgo)).toBe('1 minute ago');
    expect(timeAgo(oneHourAgo)).toBe('1 hour ago');
    expect(timeAgo(oneDayAgo)).toBe('yesterday');
    expect(timeAgo(oneWeekAgo)).toBe('last week');
    expect(timeAgo(oneMonthAgo)).toBe('last month');
    expect(timeAgo(oneYearAgo)).toBe('last year');
  });

  // Test future dates
  test('handles future dates correctly', () => {
    const now = new Date();
    const oneMinuteInFuture = new Date(now.getTime() + 60 * 1000);
    expect(timeAgo(oneMinuteInFuture)).toBe('in 1 minute');
  });

  // Testing edge cases
  test('returns correct time ago for edge cases', () => {
    const now = new Date();
    const justUnderAMinute = new Date(now.getTime() - 59 * 1000);
    const justUnderAnHour = new Date(now.getTime() - 3599 * 1000);
    const justUnderADay = new Date(now.getTime() - 86399 * 1000);
    const justUnderAWeek = new Date(now.getTime() - 6 * 24 * 3600 * 1000 - 86399 * 1000);

    expect(timeAgo(justUnderAMinute)).toBe('59 seconds ago');
    expect(timeAgo(justUnderAnHour)).toBe('59 minutes ago');
    expect(timeAgo(justUnderADay)).toBe('23 hours ago');
    expect(timeAgo(justUnderAWeek)).toBe('6 days ago');
  });
});

