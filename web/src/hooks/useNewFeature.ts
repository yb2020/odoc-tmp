import { useLocalStorage } from '@vueuse/core';
import { computed } from 'vue';

export const useNewFeature = (
  name: string,
  year: number,
  normalMonth: number,
  date: number
) => {
  const deadline = new Date();
  deadline.setFullYear(year);
  deadline.setMonth(normalMonth - 1);
  deadline.setDate(date);
  deadline.setHours(0);
  deadline.setMinutes(0);
  deadline.setSeconds(0);
  deadline.setMilliseconds(0);

  const today = new Date();
  const key = `${deadline.toLocaleDateString()}-${name}`;

  // 产品逻辑：一次只能有一个功能展示提示！
  const iGotIt = useLocalStorage('pdf-annotate/2.0/new-feature-i-got-it', '');
  const tipsVisible = computed(() => today < deadline && iGotIt.value !== key);

  const hideTips = () => {
    iGotIt.value = key;
  };

  return {
    tipsVisible,
    hideTips,
  };
};
