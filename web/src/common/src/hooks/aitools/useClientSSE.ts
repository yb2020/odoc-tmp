import { TranslateSentence } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { Ref, computed, ref } from 'vue';
import usePollingSSE, {
  AbortError,
  OnlyMessageError,
  ResetError,
} from '@common/hooks/aitools/usePollingSSE';
import { message } from 'ant-design-vue';
import {
  ElementName,
  EventCode,
  ModuleType,
  reportEvent,
} from '@common/utils/report';
import { BeanScenes } from '../useAIBeans';

const PolishType2ElementName: Record<string, string> = {
  [ModuleType.POLISH_REWRITE_ZH]: ElementName.upper_polish_rewrite_zh,
  [ModuleType.POLISH_REWRITE]: ElementName.upper_polish_rewrite,
  [ModuleType.ZH_TO_EN]: ElementName.upper_zh_to_en,
};

declare module '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText' {
  interface TranslateSentence {
    createTime: string;
    isRefundAiBean?: boolean;
  }
}

export const useDataSSE = <P, R, D>(
  url: string,
  scene: BeanScenes,
  dataRef: Ref<R>,
  updateDataRef: (d?: D) => void,
  reportParams?: object
) => {
  const error = ref<Error | null>(null);

  // const startTime = 0
  // let firstTime = 0
  // const lastTime = 0

  const handleText = (res: Error | D) => {
    if (res instanceof ResetError || res instanceof AbortError) {
      updateDataRef();
      error.value = null;
      return;
    }
    if (res instanceof OnlyMessageError) {
      updateDataRef();
      error.value = res;
      message.error(res.message);
      return;
    }
    if (res instanceof Error) {
      updateDataRef();
      error.value = res;
      message.error(res.message);
      return;
    }
    updateDataRef(res);
  };

  const {
    startPolling,
    loading,
    // error,
    abortRequest,
  } = usePollingSSE<P, D>(url, scene, handleText, reportParams);

  // const startPollingData = (chat: P) => {
  //   startTime = Date.now()
  //   firstTime = 0
  //   lastTime = 0
  //   const p = startPolling(chat)
  //   // eslint-disable-next-line promise/catch-or-return
  //   p.promise
  //     .then(() => {
  //       lastTime = Date.now()
  //       if (reportParams) {
  //         const { event, ...rest } = reportParams
  //         reportEvent(event, {
  //           ...rest,
  //           first_word_response_time: firstTime,
  //           last_word_response_time: lastTime,
  //           requestion_statr_time: startTime,
  //           polish_task_id: requestId.value,
  //         })
  //       }
  //       return
  //     })
  //     .finally(() => {
  //       startTime = 0
  //       firstTime = 0
  //       lastTime = 0
  //     })

  //   return p
  // }

  return {
    dataRef,
    startPollingData: startPolling,
    loading,
    error,
    abortRequest,
  };
};

export const useTextSSE = <P>(
  url: string,
  reportParams?: {
    event: EventCode;
    page_type: string;
    [k: string]: string;
  }
) => {
  const text = ref('');
  const requestId = ref('');
  const result = ref<TranslateSentence>();
  const { startPollingData, ...rest } = useDataSSE<
    P,
    string,
    TranslateSentence
  >(
    url,
    BeanScenes.POLISH,
    text,
    (d) => {
      if (d === undefined) {
        text.value = '';
        return;
      }
      if (!firstTime) {
        firstTime = Date.now();
      }
      text.value += d.target;
      requestId.value = d.taskId;
      result.value = d;
    },
    {
      page_type: reportParams?.page_type,
      element_name: PolishType2ElementName[reportParams?.polish_type as string],
    }
  );

  const isResolvedError = computed(() => {
    if (rest.error?.value instanceof OnlyMessageError) {
      return null;
    }
    return rest.error?.value;
  });

  let startTime = 0;
  let firstTime = 0;
  let lastTime = 0;

  const startPollingText = (chat: P, isRetry = false) => {
    startTime = Date.now();
    firstTime = 0;
    lastTime = 0;
    const p = startPollingData(chat, isRetry);
    // eslint-disable-next-line promise/catch-or-return
    p.promise
      .then(() => {
        lastTime = Date.now();
        if (reportParams) {
          const { event, ...rest } = reportParams;
          reportEvent(event, {
            ...rest,
            first_word_response_time: firstTime,
            last_word_response_time: lastTime,
            requestion_statr_time: startTime,
            polish_task_id: requestId.value,
          });
        }
        return;
      })
      .finally(() => {
        startTime = 0;
        firstTime = 0;
        lastTime = 0;
      });

    return p;
  };

  return {
    text,
    startPollingText,
    isResolvedError,
    requestId,
    result,
    ...rest,
  };
};

interface SentenceResult {
  sentences: TranslateSentence[];
  createTime: string;
  isRefundAiBean?: boolean;
}

export const useSentenceSSE = <P>(
  url: string,
  reportParams?: {
    event: EventCode;
    page_type: string;
    [k: string]: string;
  }
) => {
  const sentences = ref<TranslateSentence[]>([]);
  const requestId = ref('');
  const result = ref<SentenceResult>();
  const { startPollingData, ...rest } = useDataSSE<
    P,
    TranslateSentence[],
    SentenceResult
  >(
    url,
    BeanScenes.POLISH,
    sentences,
    (d) => {
      if (d === undefined) {
        sentences.value = [];
        return;
      }
      if (!firstTime) {
        firstTime = Date.now();
      }
      sentences.value = d.sentences;
      requestId.value = d.sentences?.[0]?.taskId || '';
      result.value = d;
    },
    {
      page_type: reportParams?.page_type,
      element_name: PolishType2ElementName[reportParams?.polish_type as string],
    }
  );

  let startTime = 0;
  let firstTime = 0;
  let lastTime = 0;
  const startPollingSentences = (chat: P) => {
    startTime = Date.now();
    firstTime = 0;
    lastTime = 0;
    const p = startPollingData(chat);
    // eslint-disable-next-line promise/catch-or-return
    p.promise
      .then(() => {
        lastTime = Date.now();
        if (reportParams) {
          const { event, ...rest } = reportParams;
          reportEvent(event, {
            ...rest,
            first_word_response_time: firstTime,
            last_word_response_time: lastTime,
            requestion_statr_time: startTime,
            polish_task_id: requestId.value,
          });
        }
        return;
      })
      .finally(() => {
        startTime = 0;
        firstTime = 0;
        lastTime = 0;
      });

    return p;
  };

  return {
    sentences,
    startPollingSentences,
    requestId,
    result,
    ...rest,
  };
};

// export const useTextSSE = <P>(url: string, reportParams?: {
//   event: EventCode
//   page_type: string
//   [k: string]: string
// }) => {
//   const text = ref('')
//   const requestId = ref('')
//   const error = ref<Error | null>(null)

//   let startTime = 0
//   let firstTime = 0
//   let lastTime = 0

//   const handleText = (res: Error | string) => {
//     if (res instanceof ResetError || res instanceof AbortError) {
//       text.value = ''
//       error.value = null
//       return
//     }
//     if (res instanceof OnlyMessageError) {
//       text.value = ''
//       error.value = null
//       message.error(res.message)
//       return
//     }
//     if (res instanceof Error) {
//       text.value = ''
//       error.value = res
//       message.error(res.message)
//       return
//     }
//     if (!firstTime) {
//       firstTime = Date.now()
//     }
//     text.value = text.value + res
//   }

//   const {
//     startPolling,
//     loading,
//     // error,
//     abortRequest,
//   } = usePollingSSE<P, string>(url, handleText)

//   const startPollingText = (chat: P) => {
//     startTime = Date.now()
//     firstTime = 0
//     lastTime = 0
//     const p = startPolling(chat)
//     // eslint-disable-next-line promise/catch-or-return
//     p.promise.then(() => {
//       lastTime = Date.now()
//       if (reportParams) {
//         const { event, ...rest } = reportParams
//         reportEvent(event, {
//           ...rest,
//           first_word_response_time: firstTime,
//           last_word_response_time: lastTime,
//           requestion_statr_time: startTime,
//           polish_task_id: requestId.value,
//         })
//       }
//       return
//     }).finally(() => {
//       startTime = 0
//       firstTime = 0
//       lastTime = 0
//     })

//     return p;
//   }

//   return {
//     text,
//     startPollingText,
//     loading,
//     error,
//     abortRequest,
//     requestId,
//   }
// }

// export const useSentenceSSE = <P>(url: string) => {
//   const sentences = ref<TranslateSentence[]>([])
//   const requestId = ref('')
//   const error = ref<Error | null>(null)

//   const handleSentences = (res: Error | { sentences: TranslateSentence[] }) => {
//     if (res instanceof ResetError || res instanceof AbortError) {
//       sentences.value = []
//       error.value = null
//       return
//     }
//     if (res instanceof Error) {
//       sentences.value = []
//       error.value = res
//       message.error(res.message)
//       return
//     }
//     console.log(res)
//     sentences.value = res.sentences
//   }

//   const {
//     startPolling: startPollingSentences,
//     loading,
//     // error,
//     abortRequest,
//   } = usePollingSSE<P, { sentences: TranslateSentence[] }>(url, handleSentences)

//   return {
//     sentences,
//     startPollingSentences,
//     loading,
//     error,
//     abortRequest,
//     requestId,
//   }
// }
