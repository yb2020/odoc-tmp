/* eslint-disable no-underscore-dangle */
import { defineStore } from 'pinia';
import type {
  AxiosStatic,
  AxiosRequestConfig,
  AxiosResponse,
  AxiosInstance,
} from 'axios';
import axios from 'axios';
import { ResponseError } from '@common/api/type';

declare module 'axios' {
  interface AxiosRequestConfig {
    isCertRetry?: boolean;
    noCertRetry?: boolean;
  }
}

const DEFAULT_SCENE_ID = 'workflowSceneNormal';

export interface CertProps {
  sceneId?: string;
}
export interface CertResult {
  status: 0 | 1;
}

export type CertHandler = (x?: CertProps) => Promise<CertResult>;

export type CertCallback = (x: CertResult) => void;

export type CertPayload = CertProps & {
  callback: CertCallback;
};

export type CertState = Required<CertProps> & {
  visible: boolean;
  result?: CertResult;
  callbacks: CertCallback[];
};

export const ERROR_CODE_UNCERT = 4200;

export const state = () => ({
  sceneId: '',
  visible: false,
  callbacks: [] as CertCallback[],
});

export const useCertStore = defineStore('cert', {
  state,
  actions: {
    showCertDialog({ sceneId, callback }: CertPayload) {
      this.sceneId = sceneId || DEFAULT_SCENE_ID;
      this.visible = true;
      this.callbacks.push(callback);
    },
    hideCertDialog(payload: CertResult) {
      this.visible = false;
      this.callbacks.forEach((cb) => cb?.(payload));
      this.callbacks = [];
    },
  },
});

const weakMap = new WeakMap<object, CertHandler>();

export function register($axios: AxiosInstance, fn: CertHandler) {
  weakMap.set($axios, fn);
}

function fixConfig(rootAxios: AxiosStatic, config: AxiosRequestConfig) {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  if (rootAxios.defaults.agent === config.agent) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    delete config.agent;
  }
  if (rootAxios.defaults.httpAgent === config.httpAgent) {
    delete config.httpAgent;
  }
  if (rootAxios.defaults.httpsAgent === config.httpsAgent) {
    delete config.httpsAgent;
  }
}

export function isUnCert(response: AxiosResponse, $axios: AxiosInstance) {
  const { data, config } = response;
  if (data.code === ERROR_CODE_UNCERT) {
    const onNeedCert = weakMap.get($axios);
    const sceneId = data.data;
    const p =
      onNeedCert?.({
        sceneId,
      }) ??
      Promise.resolve({
        status: 0,
      });

    return p.then(({ status }) => {
      const isNeedRetry =
        'noCertRetry' in config
          ? !config.noCertRetry
          : !sceneId || sceneId === DEFAULT_SCENE_ID;
      if (status !== 0 || config.isCertRetry || !isNeedRetry) {
        throw new ResponseError({
          code: ERROR_CODE_UNCERT,
          message: data.message,
          extra: {
            sceneId,
          },
        });
      }

      // Copied from axios-retry
      // Axios fails merging this configuration to the default configuration because it has an issue
      // with circular structures: https://github.com/mzabriskie/axios/issues/370
      // eslint-disable-next-line @typescript-eslint/no-use-before-define
      fixConfig(axios, config);

      // 重试
      return $axios({
        ...config,
        isCertRetry: true,
      });
    });
  }

  return false;
}
