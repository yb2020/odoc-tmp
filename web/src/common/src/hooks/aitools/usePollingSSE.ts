import ssePostFetcher from './sse';
import { message } from 'ant-design-vue';
import { getCurrentInstance, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { BeanScenes, useAIBeans, useAIBeansBuy } from '../useAIBeans';
import { PageType } from '@common/utils/report';
import { ResponseError } from '../../api/type';
import {
  ERROR_CODE_BEANS_NOT_ENOUGH,
  ERROR_CODE_NEED_VIP,
} from '../../api/const';
import { NeedAiBeanException } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';

export class ResetError extends Error {
  constructor() {
    super('reset error');
  }
}

export class OnlyMessageError extends Error {
  constructor(msg: string) {
    super(msg);
  }
}

export class AbortError extends Error {
  constructor() {
    super('abort error');
  }
}

export class IgnoreWithMsgError extends Error {
  constructor(msg: string) {
    super(msg);
  }
}

const usePollingSSE = <P, R>(
  url: string,
  scene: BeanScenes,
  callback: (res: Error | R) => void,
  reportParams?: object
) => {
  const loading = ref(false);
  const error = ref<Error | null>(null);
  const requestId = ref('');

  // let abortSignal: AbortController | null = null

  const handleText = (res: Error | R) => {
    if (res instanceof Error && !(res instanceof ResetError)) {
      refundBeans(BeanScenes.POLISH);
      const flag = checkBeansErr(res, {
        page_type: PageType.POLISH,
        ...reportParams,
      });
      if (flag) {
        return;
      }
    }

    callback(res);
  };

  const modifyText = (str: string) => {
    try {
      const res = JSON.parse(str) as {
        data: unknown;
        status: number;
        message: string;
      };
      if (res.status === 0) {
        if (res.data) {
          defered?.reject(new OnlyMessageError(res.data as string));
        } else {
          defered?.reject(new Error(res.message));
        }
        return;
      }
      if (
        [ERROR_CODE_BEANS_NOT_ENOUGH, ERROR_CODE_NEED_VIP].includes(res.status)
      ) {
        defered?.reject(
          new ResponseError({
            code: res.status,
            message: res.message,
            extra: res.data as any as NeedAiBeanException,
          })
        );
        return;
      }
      handleText(res.data as R);
    } catch (error) {
      handleText(error as Error);
    }
  };

  const { t } = useI18n();
  const { consumeBeans, refundBeans } = useAIBeans();
  const { checkBeansErr } = useAIBeansBuy(getCurrentInstance()?.appContext);

  let defered: {
    promise: Promise<void>;
    resolve: () => void;
    reject: (err: Error) => void;
  } | null = null;

  const startPolling = (chat: P, isRetry = false) => {
    const createDefered = () => {
      let _resolve: () => void;
      let _reject: (err: Error) => void;
      const abortSignal = new AbortController();
      const promise = new Promise<void>((resolve, reject) => {
        _resolve = () => {
          resolve();
          defered = null;
          loading.value = false;
        };
        _reject = (err: Error) => {
          if (err instanceof IgnoreWithMsgError) {
            message.warning(err.message);
            reject(err);
            return;
          }
          if (err instanceof AbortError) {
            abortSignal.abort();
          }
          defered = null;
          loading.value = false;
          handleText(err);
          error.value = err as Error;
          reject(err);
        };
        if (loading.value) {
          _reject(new IgnoreWithMsgError(t('common.tips.processing')));
          return;
        }
        handleText(new ResetError());
        error.value = null;
        loading.value = true;
        (isRetry ? Promise.resolve(true) : consumeBeans(scene))
          .then(async (isContinue) => {
            if (isContinue) {
              return ssePostFetcher<P, string>(
                url,
                chat,
                modifyText,
                abortSignal
              );
            }

            return false;
          })
          .then(() => {
            console.log('stream ends');
            _resolve();
            return;
          })
          .catch((err) => {
            console.log(err);
            _reject(err);
          });
      });
      return {
        promise,
        resolve: () => {
          _resolve();
        },
        reject: (err: Error) => {
          _reject(err);
        },
      };
    };
    defered = createDefered();

    return defered;
  };

  // const startPollingAnswer = async (chat: P) => {
  //   if (loading.value) {
  //     message.warning(t('common.tips.processing'))
  //     return;
  //   }
  //   handleText(new ResetError())
  //   error.value = null
  //   loading.value = true
  //   abortSignal = new AbortController()
  //   try {
  //     // requestId.value = new Date().getTime().toString()
  //     await ssePostFetcher<P, string>(url, chat, modifyText, abortSignal)
  //     console.log('stream end')
  //   } catch (err) {
  //     console.log(err)
  //     handleText(err as Error)
  //     error.value = err as Error
  //   }
  //   loading.value = false
  // }

  return {
    loading,
    startPolling,
    // startPollingAnswer,
    abortRequest: () => {
      if (defered) {
        defered.reject(new AbortError());
      }
      // if (abortSignal) {
      //   abortSignal.abort()
      // }
      // loading.value = false
      // handleText(new ResetError())
    },
    error,
    requestId,
  };
};

export default usePollingSSE;
