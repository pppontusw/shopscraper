function timeAgo(date) {
  const now = new Date();
  const diffInSeconds = (date - now) / 1000; // If negative, date is in the past
  const absDiffInSeconds = Math.abs(diffInSeconds);

  const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' });

  const ranges = {
    year: 3600 * 24 * 365,
    month: 3600 * 24 * 30,
    week: 3600 * 24 * 7,
    day: 3600 * 24,
    hour: 3600,
    minute: 60,
    second: 1,
  };

  for (const key in ranges) {
    const amount = Math.floor(absDiffInSeconds / ranges[key]);
    if (amount >= 1) {
      if (key === 'day' && amount === 1) {
        return diffInSeconds < 0 ? 'yesterday' : 'tomorrow';
      }
      if (key === 'week' && amount === 1) {
        return diffInSeconds < 0 ? 'last week' : 'next week';
      }
      if (key === 'month' && amount === 1) {
        return diffInSeconds < 0 ? 'last month' : 'next month';
      }
      if (key === 'year' && amount === 1) {
        return diffInSeconds < 0 ? 'last year' : 'next year';
      }
      return rtf.format(diffInSeconds < 0 ? -amount : amount, key);
    }
  }

  return 'just now';
}

export default timeAgo