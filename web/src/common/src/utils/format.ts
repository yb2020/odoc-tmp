import { useI18n } from 'vue-i18n';
import dayjs from 'dayjs';

export function formatRecentDate(ts?: number | string) {
  if (!ts) {
    return '';
  }
  const time = dayjs(ts);
  const now = dayjs();
  const { t } = useI18n();
  if (time.year() < now.year()) {
    return time.format('YYYY-MM-DD');
  }

  const offsetDay = now.date() - time.date();
  if (!time.isSame(now, 'M') || offsetDay > 1 || offsetDay < 0) {
    return time.format('MM-DD HH:mm');
  }
  const prefix =
    offsetDay === 1 ? t('common.time.yesterday') : t('common.time.today');
  return prefix + ' ' + time.format(`HH:mm`);
}
