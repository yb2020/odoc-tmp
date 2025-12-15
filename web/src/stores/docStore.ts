import { defineStore } from 'pinia';
import { ref } from 'vue';
import { getUserDoc } from '~/src/api/material';
import { DocDetailInfo } from 'go-sea-proto/gen/ts/doc/ClientDoc';
import { UserDocParsedStatusEnum } from 'go-sea-proto/gen/ts/doc/UserDocParsedStatus';

export const useDocStore = defineStore('docStore', () => {
  const docInfo = ref<DocDetailInfo | null>(null);
  const isLoading = ref(false);
  const error = ref<Error | null>(null);
  const isEmbeddingReady = ref(false);

  const fetchDocInfo = async (noteId: bigint) => {
    if (!noteId) {
      docInfo.value = null;
      return;
    }

    isLoading.value = true;
    error.value = null;
    try {
      const response = await getUserDoc({ id: noteId, noteId });
      docInfo.value = response;
      // @ts-ignore
      isEmbeddingReady.value = (response as any)?.embeddingStatus === UserDocParsedStatusEnum.EMBEDDED;
    } catch (e) {
      error.value = e as Error;
      docInfo.value = null;
    } finally {
      isLoading.value = false;
    }
  };

  /**
   * 更新文档状态（用于轮询时同步状态）
   * @param parsedStatus - 解析状态
   * @param embeddingStatus - embedding 状态
   */
  const updateDocStatus = (parsedStatus: number, embeddingStatus: number) => {
    if (docInfo.value) {
      // @ts-ignore - 更新 parsedStatus
      docInfo.value.parsedStatus = parsedStatus;
      // @ts-ignore - 更新 embeddingStatus
      docInfo.value.embeddingStatus = embeddingStatus;
      // 同时更新 isEmbeddingReady
      isEmbeddingReady.value = embeddingStatus === UserDocParsedStatusEnum.EMBEDDED;
    }
  };

  return {
    docInfo,
    isLoading,
    error,
    isEmbeddingReady,
    fetchDocInfo,
    updateDocStatus,
  };
});
