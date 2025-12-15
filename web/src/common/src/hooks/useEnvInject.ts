import { isObject } from 'lodash';
import { useUrlSearchParams } from '@vueuse/core';
import { AxiosInstance, AxiosRequestConfig } from 'axios';
import axios from '@common/api/axios';

export function useEnvInject(axiosInstances?: AxiosInstance[]) {
  const params = useUrlSearchParams<{
    env: string;
  }>();
  const { env } = params;

  if (typeof env === 'string' && ['dev', 'gray'].includes(env)) {
    // $axios.defaults.params = {
    //   ...$axios.defaults.params,
    //   env,
    // }
    const addEnv = (config: AxiosRequestConfig) => {
      if (config.url?.startsWith('/')) {
        if (config.params instanceof URLSearchParams) {
          config.params.append('env', env);
        } else if (isObject(config.params)) {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          config.params.env = env;
        } else {
          config.params = {
            ...config.params,
            env,
          };
        }
      }

      return config;
    };

    // @TODO 待browser/package.json依赖的axios升级到1.x
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    axios.interceptors.request.use(addEnv);
    axiosInstances?.forEach((instance) => {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      instance.interceptors.request.use(addEnv);
    });
  }
}

export default useEnvInject;
