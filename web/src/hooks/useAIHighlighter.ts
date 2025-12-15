import { useRequest } from 'ahooks-vue';
import {
  getHighlightResult,
  getRightInfo,
  startHighlight,
} from '../api/aiHighlighter';
import { createSharedComposable } from '@vueuse/core';
import { Ref, computed, ref } from 'vue';
import { SCIMItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/scim/AiSCIMInfo';

export interface HighlightsMap {
  [k: string]: {
    type: string;
    color: string;
    checked: boolean;
    items: SCIMItem[];
  };
}

function useAIHighligterRaw() {
  return useRequest(getRightInfo, {});
}

export const useAIHighligter = createSharedComposable(useAIHighligterRaw);

function useAIHighlightsRaw(pdfId: Ref<string>, noteId: Ref<string>) {
  const { data: rights } = useAIHighligter();

  const start = async () => {
    if (rights.value?.canUse && rights.value.remainingTimes > 0) {
      rights.value.remainingTimes -= 1;
      await startHighlight({
        pdfId: pdfId.value,
        noteId: noteId.value,
      });
      polling.value = true;
    }
  };

  const polling = ref(true);
  const errtimes = ref(0);
  const ready = computed(
    () =>
      !!rights.value?.canUse && !!pdfId.value && !!noteId.value && polling.value
  );
  const { data, ...rest } = useRequest(
    async () => {
      const res = await getHighlightResult({
        pdfId: pdfId.value,
        noteId: noteId.value,
      });

      polling.value = res.progress > 0 && res.progress !== 100;

      return {
        ...res,
        traceid: res.items?.[0]?.traceId,
        map: (res.items ?? []).reduce((acc, cur) => {
          if (!acc[cur.scimType]) {
            acc[cur.scimType] = {
              type: cur.scimType,
              color: cur.color,
              checked: cur.checked,
              items: [],
            };
          }
          acc[cur.scimType].items.push(cur);
          return acc;
        }, {} as HighlightsMap),
      };
    },
    {
      ready,
      refreshDeps: [pdfId, noteId],
      pollingInterval: 3000,
      pollingSinceLastFinished: true,
      pollingWhenHidden: false,
      onSuccess() {
        errtimes.value = 0;
      },
      onError() {
        errtimes.value += 1;
        polling.value = errtimes.value <= 3;
      },
    }
  );

  return {
    start,
    data,
    ...rest,
  };
}

export const useAIHighlights = createSharedComposable(useAIHighlightsRaw);
