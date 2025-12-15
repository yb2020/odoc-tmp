import { useEventListener, useUrlSearchParams } from '@vueuse/core';
import { UseRequestOptions, useRequest } from 'ahooks-vue';
import { Ref, computed, ref, watch } from 'vue';
// import { updateSummaryNote } from '@common/api/note';
import { getSummaryNote, updateSummaryNote } from '~/src/api/note';
import { GetNoteSummaryByNoteIdResponse } from 'go-sea-proto/gen/ts/note/NoteSummary'

export interface SummaryEvent {
  type: 'unlock';
  noteId: string;
}

export const KEY_DETACHED = 'detached';
export const EVENT_DETACHED = `note:summary`;

/**
 * @description
 * 状态由客户端控制
 */
export function useSummaryDetachedStatus(noteId: string, defaultValue = false) {
  const query = useUrlSearchParams();
  const detached = ref(
    KEY_DETACHED in query ? query.detached === '1' : defaultValue
  );

  useEventListener(window, EVENT_DETACHED, (e: CustomEvent<SummaryEvent>) => {
    if (e.detail?.type === 'unlock' && e.detail?.noteId === noteId) {
      detached.value = false;
    }
  });

  watch(detached, () => {
    const url = new URL(window.location.href);
    if (detached.value) {
      url.searchParams.set(KEY_DETACHED, '1');
    } else {
      url.searchParams.delete(KEY_DETACHED);
    }
    window.history.replaceState({}, '', url);
  });

  return {
    detached,
  };
}

export function useSummary(
  noteId: Ref<string>,
  opts?: Partial<UseRequestOptions<GetNoteSummaryByNoteIdResponse>> & {
    onConflict?: () => Promise<boolean>;
  }
) {
  const ready = computed(() => !!noteId.value);

  const summary = ref('');
  const prevData = ref<GetNoteSummaryByNoteIdResponse>();
  const prevTs = ref(Number.POSITIVE_INFINITY);
  const {
    data,
    run: refresh,
    ...rest
  } = useRequest(
    async (isCover = true) => {
      const res = await getSummaryNote({
        noteId: noteId.value, // 使用字符串ID，API层处理转换
      });

      if (isCover) {
        prevData.value = res;
        prevTs.value = Date.now();
        summary.value = res?.content ?? '';
      }

      return res;
    },
    {
      ready,
      refreshDeps: [noteId],
      ...opts,
    }
  );

  const { run: save, loading: isSaving } = useRequest(
    async (isNeedCheck = true) => {
      if (rest.loading.value) {
        return;
      }

      if (isNeedCheck) {
        await refresh(false);

        const preCt = prevData.value?.content; // 本地更改前内容
        const curCt = data.value?.content; // 远程内容
        const curTs = Number(data.value?.modifyDate);
        const nexCt = summary.value; // 本地更改后内容
        if (
          curCt &&
          preCt &&
          curCt !== preCt &&
          curCt !== nexCt &&
          curTs > prevTs.value
        ) {
          const flag = await opts?.onConflict?.();

          if (flag !== true) {
            summary.value = data.value?.content ?? '';
            throw new Error('Conflicted');
          }
        }
      }

      prevData.value = {
        ...data.value!,
        content: summary.value,
      };

      return updateSummaryNote({
        noteId: noteId.value,
        content: summary.value,
      });
    },
    {
      manual: true,
    }
  );

  return {
    ...rest,
    data,
    prevData,
    summary,
    refresh,
    save,
    isSaving,
  };
}
