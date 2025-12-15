import { BaseUseRequestOptions, useRequest } from 'ahooks-vue';
import {
  createSharedComposable,
  UseOffsetPaginationOptions,
  useOffsetPagination,
} from '@vueuse/core';
import { ComputedRef, Ref, computed, ref, triggerRef } from 'vue';
import {
  getWordNotes,
  addWordNote,
  updateWordNote,
  delWordNote,
  setNoteWordsConfig,
} from '~/src/api/note';
import {
  displayMode as DisplayMode,
  // ChangeWordConfigReq,
  // DeleteWordReq,
  // SaveWordReq,
  // SaveWordResponse,
  // UpdateWordReq,
  // WordInfo,
  // GetWordsResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import { GetNoteWordsByNoteIdResponse, WordInfo, ChangeWordConfigRequest, DeleteNoteWordRequest, UpdateNoteWordRequest, SaveNoteWordRequest, SaveNoteWordResponse } from 'go-sea-proto/gen/ts/note/NoteWord'
import { fetchTranslate, fetchUniTranslate, TranslateTabKey } from '~/src/api/translate';
import { VocabularyColorKey, vocabularyStyleMap } from './types';
import { useUserStore } from '@common/stores/user';



export { DisplayMode };

export const useWordNotes = createSharedComposable(
  (
    pdfId: Ref<string>,
    noteId: Ref<string>,
    paginationOpts?: Partial<UseOffsetPaginationOptions>,
    addOpts?: Partial<
      BaseUseRequestOptions<SaveNoteWordResponse, [Omit<SaveNoteWordRequest, 'noteId'>]>
    >,
    prefOpts?: Partial<
      BaseUseRequestOptions<object, [Omit<ChangeWordConfigRequest, 'noteId'>]>
    >
  ) => {
    
    const ready = computed(() => !!noteId.value);
    const data = ref<WordInfo[]>([]);
    const gmode = ref(DisplayMode.GLOBAL);
    const gcolor = ref(VocabularyColorKey.PURPLE);
    const gcolorHex = computed(
      () => vocabularyStyleMap[gcolor.value ?? VocabularyColorKey.PURPLE].fill
    );
    const total = ref(1);
    const lastNoteId = computed(() => data.value?.[data.value.length - 1]?.id);
    const { run, ...rest } = useRequest(
      async () => {
        // useOffsetPagination有BUG，待升级
        if (currentPage.value === 0) {
          return;
        }

        const res: GetNoteWordsByNoteIdResponse = await getWordNotes({
          noteId: noteId.value,
          pageSize: currentPageSize.value,
          currentPage: currentPage.value,
          minLoadedId: lastNoteId.value,
        });

        total.value = res.total;
        gmode.value = res.displayMode ?? DisplayMode.GLOBAL;
        gcolor.value = res.color ?? VocabularyColorKey.PURPLE;

        if (lastNoteId.value) {
          data.value = [...data.value, ...res.words];
        } else {
          data.value.splice(
            (currentPage.value - 1) * currentPageSize.value,
            currentPageSize.value,
            ...res.words
          );
          triggerRef(data);
        }

        return res;
      },
      {
        ready,
        refreshDeps: [noteId],
        manual: true,
      }
    );
    const { currentPage, currentPageSize, isLastPage, next } =
      useOffsetPagination({
        pageSize: 10,
        onPageChange: run,
        // onPageSizeChange: run,
        ...paginationOpts,
        total,
      });

    const { run: config, loading: isConfiging } = useWordNotesPrefOp(noteId, {
      ...prefOpts,
      onSuccess(_, params) {
        gmode.value = params[0]?.displayMode ?? gmode.value;
        gcolor.value = params[0]?.color ?? gcolor.value;
        prefOpts?.onSuccess?.(_, params);
      },
    });

    const {
      data: added,
      run: add,
      loading: isAdding,
    } = useWordNoteAddOp(pdfId, noteId, {
      ...addOpts,
      onSuccess(res, params) {
        if (res?.id) {
          total.value += 1;
          if (data.value && params[0]?.wordInfo) {
            data.value = [
              {
                id: res.id,
                word: params[0].wordInfo.word!,
                rectangle: params[0].wordInfo.rectangle!,
                translateInfo: params[0].wordInfo.translateInfo,
              },
              ...data.value,
            ];
          }
        }
        addOpts?.onSuccess?.(res, params);
      },
    });

    const mutate = (id: string, v: string) => {
      const item = data.value?.find((x) => x.id === id);
      if (item?.translateInfo) {
        item.translateInfo = {
          ...item.translateInfo,
          targetContent: [v],
          targetResp: [],
        };
        triggerRef(data);
      }
    };

    const remove = (id: string) => {
      data.value = data.value.filter((x) => x.id !== id);
      total.value -= 1;
    };

    return {
      ...rest,
      data,
      added,
      gmode,
      gcolor,
      gcolorHex,
      lastNoteId,
      total,
      isLastPage,
      isConfiging,
      isAdding,
      run,
      config,
      add,
      mutate,
      remove,
      next,
    };
  }
);

export function useWordNotesPrefOp(
  noteId: Ref<string>,
  opts?: Partial<
    BaseUseRequestOptions<object, [Omit<ChangeWordConfigRequest, 'noteId'>]>
  >
) {
  return useRequest(
    async (params: Omit<ChangeWordConfigRequest, 'noteId'>) => {
      await setNoteWordsConfig({
        noteId: noteId.value,
        ...params,
      });
    },
    {
      ...opts,
      manual: true,
    }
  );
}

export function useWordNoteAddOp(
  pdfId: Ref<string>,
  noteId: Ref<string>,
  opts?: Partial<
    BaseUseRequestOptions<SaveNoteWordResponse, [Omit<SaveNoteWordRequest, 'noteId'>]>
  >
) {
  return useRequest(
    async (params) => {
      const { wordInfo } = params;

      if (!wordInfo || !wordInfo?.word) {
        throw new Error('invalid word');
      }

      if (!wordInfo.translateInfo) {
        // const {
        //   targetContent,
        //   targetResp = [],
        //   ...rest
        // } = await fetchUniTranslate({
        //   channel: TranslateTabKey.idea,
        //   content: wordInfo?.word,
        //   pdfId: pdfId.value,
        // });
        const {
          targetContent,
          targetResp = [],
          ...rest
        } = await fetchTranslate({
          type: TranslateTabKey.google,
          content: wordInfo?.word,
          pdfId: pdfId.value,
          param: {},
          useGlossary: false,
        });

        wordInfo.translateInfo = {
          ...rest,
          targetResp,
          targetContent: Array.isArray(targetContent)
            ? targetContent
            : [targetContent],
        };

        
        //刷新积分
        const userStore = useUserStore();
        userStore.refreshUserCredits()
      }

      const res = await addWordNote({
        ...params,
        noteId: noteId.value,
      });

      return res;
    },
    {
      manual: true,
      ...opts,
    }
  );
}

export function useWordNoteOp(wid: Ref<string> | ComputedRef<string>) {
  const { run: updateWord, loading: isUpdating } = useRequest(
    async (params: Omit<UpdateNoteWordRequest, 'wordId'>) => {
      await updateWordNote({
        wordId: wid.value,
        ...params,
      });
    },
    {
      manual: true,
    }
  );

  const { run: delWord, loading: isDeleting } = useRequest(
    async (params?: Omit<DeleteNoteWordRequest, 'wordId'>) => {
      await delWordNote({
        wordId: wid.value,
        ...params,
      });
    },
    {
      manual: true,
    }
  );

  return {
    updateWord,
    isUpdating,
    delWord,
    isDeleting,
  };
}
