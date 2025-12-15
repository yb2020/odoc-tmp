import { REQUEST_SERVICE_NAME_AI_POLISH } from '@common/api/const';
import { Ref, computed, ref, unref } from 'vue';
import {
  PolishTextResponse,
  SelectTextEnum,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { ParagraphMode } from '@common/components/AIParagraph/type';
import usePollingSSE, { AbortError } from './usePollingSSE';
import { ElementName } from '@common/utils/report';
import { BeanScenes } from '../useAIBeans';
// import { message } from 'ant-design-vue'
// import { useI18n } from 'vue-i18n'

export enum PollingStatus {
  SUCCESS = 'success',
  FAIL = 'fail',
  PENDING = 'pending',
}

interface PollingParams {
  text: string;
  wordCount: number;
  mode: ParagraphMode;
  synonyms: number;
  pageId?: string;
}

export const convertMode = (mode: ParagraphMode) => {
  switch (mode) {
    case ParagraphMode.shorten:
      return SelectTextEnum.SHORTEN;
    case ParagraphMode.expand:
      return SelectTextEnum.EXPAND;
    case ParagraphMode.simple:
      return SelectTextEnum.SIMPLE;
    case ParagraphMode.improve:
      return SelectTextEnum.OPTIMIZATION;
    case ParagraphMode.reduceSimilar:
      return SelectTextEnum.REDUCE_REPETITION;
    default:
      return SelectTextEnum.STANDARD;
  }
};

export const usePollingParagraphSSE = (curVersionId?: string | Ref<string>) => {
  const text = ref('');

  const handleText = (res: PolishTextResponse | Error) => {
    if (res instanceof Error) {
      text.value = '';
      requestId.value = '';
      return;
    }
    text.value = (text.value || '') + res.content;
    requestId.value = res.sessionId;
  };

  const { loading, startPolling, abortRequest, error, requestId } =
    usePollingSSE<
      {
        pageId: string;
        text: string;
        requestId: string;
        wordCount: number;
        selectTextEnum: SelectTextEnum;
        synonymsEnum: number;
      },
      PolishTextResponse
    >(
      `/api${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/text`,
      BeanScenes.POLISH,
      handleText,
      {
        element_name: ElementName.upper_polish_rewrite,
      }
    );

  // const { t } = useI18n()

  const start = (chat: PollingParams) => {
    // if (loading.value) {
    //   message.warning(t('common.tips.processing'))
    //   return
    // }
    const p = {
      pageId: chat.pageId || unref(curVersionId!),
      text: chat.text,
      requestId: requestId.value,
      wordCount: chat.wordCount,
      selectTextEnum: convertMode(chat.mode),
      synonymsEnum: chat.synonyms + 1,
    };
    return startPolling(p).promise;
  };

  return {
    loading,
    startPollingAnswer: start,
    abortRequest: () => {
      text.value = '';
      return abortRequest();
    },
    text,
    error: computed(() => {
      if (error.value instanceof AbortError) {
        return null;
      }
      return error.value;
    }),
    requestId,
  };
};
