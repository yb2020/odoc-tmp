import { createSharedComposable } from '@vueuse/core';
import { computed, watch } from 'vue';
import { useRequest } from 'ahooks-vue';
import { getActivityCardsStatus } from '../api/vip';
import { useAIBeans } from './useAIBeans';
import { ActivityTypeEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipActivitiesInfo';

function useActivityCardsRaw() {
  const { confirmable } = useAIBeans();

  const { data, ...rest } = useRequest(getActivityCardsStatus, {});
  const cards = computed(() => {
    return (data.value?.activitiesPrivilegeStatusList ?? []).reduce(
      (acc, item) => {
        acc[item.activityTypeEnum] = item.hasActivityPrivilege;

        return acc;
      },
      {} as Record<ActivityTypeEnum, boolean>
    );
  });

  watch(
    cards,
    () => {
      confirmable.value = !cards.value[ActivityTypeEnum.UN_LIMIT_POLISH_CARD];
    },
    {
      immediate: true,
    }
  );

  return {
    data,
    cards,
    ...rest,
  };
}

export const useActivityCards = createSharedComposable(useActivityCardsRaw);
