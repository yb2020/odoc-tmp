import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import Cookie from 'js-cookie';
import {
  REQUEST_APPID,
  REQUEST_ORGID,
  ERROR_CODE_ILLEGAL_TOKEN,
  ERROR_CODE_REPEAT_LOGIN,
  ERROR_CODE_TOKEN_EXPIRED,
  ERROR_CODE_UNLOGIN,
  UNKNOWN_ERROR_CODE,
  UNKNOWN_ERROR_MESSAGE,
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_USER,
  COOKIE_REFRESH_TOKEN,
  ERROR_CODE_ILLEGAL_IDENTITY_REQUEST_REFRESH,
  ERROR_CODE_EXPIRED_IDENTITY_REQUEST_REFRESH,
} from './const';
import {
  Response,
  isSuccessResponse,
  ResponseError,
  FailResponse,
  RefreshTokenReq,
  RefreshTokenResp,
  SuccessResponse,
} from './type';
import { message } from 'ant-design-vue';
import { setAuthorization } from '../utils/login';
import { v4 as uuidv4 } from 'uuid';
import { bridgeAdaptor } from '../adaptor/bridge';
import { getHostname, isInElectron } from '../utils/env';
import { isObject } from '@vueuse/core';
import { isUnCert } from '@common/stores/cert';

export const getAPIRoot = () => {
  if (
    !isInElectron() ||
    ['http:', 'https:'].includes(window.location.protocol)
  ) {
    return '/api';
  }

  return `https://${getHostname()}/api`;
};

interface RequestParam {
  orgId?: string;
  appId?: string;
}

const options: AxiosRequestConfig<RequestParam> = {
  baseURL: getAPIRoot(),
  headers: {
    Authorization: setAuthorization('aiKnowledge', 'readpaper_com'),
  },
};

export const api = axios.create(options);

api.interceptors.request.use(async (config) => {
  config.timeout = config.timeout || 10000;
  config.data = config.data || {};
  if (isObject(config.data as RequestParam)) {
    config.data.orgId = config.data.orgId || REQUEST_ORGID;
    config.data.appId = config.data.appId || REQUEST_APPID;
  }

  config.headers = config.headers || {};
  config.headers['x-traceId-header'] = uuidv4();

  // 携带跨域cookie
  config.withCredentials = true;

  return config;
});

function throwLoginError(response: AxiosResponse) {
  const data: FailResponse = response.data;

  if (data.code === ERROR_CODE_TOKEN_EXPIRED) {
    throw new ResponseError({
      code: ERROR_CODE_TOKEN_EXPIRED,
      message: data.message || 'token expired',
    });
  } else if (
    data.code === ERROR_CODE_ILLEGAL_TOKEN ||
    data.code === ERROR_CODE_ILLEGAL_IDENTITY_REQUEST_REFRESH ||
    data.code === ERROR_CODE_EXPIRED_IDENTITY_REQUEST_REFRESH ||
    data.code === ERROR_CODE_REPEAT_LOGIN
  ) {
    const beforeClear = Cookie.get(COOKIE_REFRESH_TOKEN);
    console.log({ beforeClear });
    Cookie.remove(COOKIE_REFRESH_TOKEN);
    throwUnLogin();
  }

  function throwUnLogin() {
    if (!response.config.headers?.[HEADER_CANCLE_AUTO_ERROR]) {
      bridgeAdaptor.login();
    }

    throw new ResponseError({
      code: ERROR_CODE_UNLOGIN,
      message: data.message,
    });
  }
}

api.interceptors.response.use(
  async (response) => {
    const data: Response<any> = response.data;

    if (
      response.config.responseType === 'blob' ||
      response.config.responseType === 'arraybuffer' ||
      isSuccessResponse(data)
    ) {
      return response;
    }

    const p = isUnCert(response, api);
    if (p !== false) {
      return p;
    }

    throwLoginError(response);

    if (!response.config.headers?.[HEADER_CANCLE_AUTO_ERROR]) {
      message.error(data.message || UNKNOWN_ERROR_MESSAGE);
    }

    let code = data.code;
    if (!code && /无权限/.test(data.message || '')) {
      code = 403;
    }

    throw new ResponseError({
      code: code || data.status || UNKNOWN_ERROR_CODE,
      message: data.message || UNKNOWN_ERROR_MESSAGE,
      extra: data.data,
    });
  },
  async (error) => {
    // onResponse中reject的error，这里会捕获到
    if (error instanceof ResponseError) {
      throw error;
    }

    const response = error.response;

    const code = response?.status || UNKNOWN_ERROR_CODE;

    if (code === 401) {
      throwLoginError(response);
    }

    // TODO 将错误上报给sentry
    if (!error.config?.headers[HEADER_CANCLE_AUTO_ERROR]) {
      message.error(error.message || UNKNOWN_ERROR_MESSAGE);
    }

    throw new ResponseError({
      code,
      message: error.message || UNKNOWN_ERROR_MESSAGE,
    });
  }
);

export default api;

let tokenPromise: null | Promise<void> = null;
function updateToken(refreshToken: string) {
  if (!tokenPromise) {
    tokenPromise = fetchToken();
  }

  return tokenPromise;

  async function fetchToken() {
    const platform =
      (navigator as any).userAgentData?.platform ?? navigator.platform ?? '';

    const scope = !isInElectron()
      ? 'WEB'
      : /Win32/i.test(platform) || /Windows/i.test(navigator.userAgent)
        ? 'WINPC'
        : 'MACOS';

    const params: RefreshTokenReq = {
      refreshToken,
      scope,
    };

    const search = new URLSearchParams(
      params as unknown as Record<string, string>
    );

    let data: RefreshTokenResp;

    try {
      const response = await api.post<SuccessResponse<RefreshTokenResp>>(
        `${REQUEST_SERVICE_NAME_USER}/oauth/refreshToken`,
        String(search),
        { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }
      );

      ({ data } = response.data);
    } catch (error) {
      tokenPromise = null;
      throw error;
    }

    const expires = new Date();
    expires.setDate(expires.getDate() + 1);
    expires.setMonth(expires.getMonth() + 1);
    Cookie.set(COOKIE_REFRESH_TOKEN, data.refreshToken, {
      expires,
      path: '/',
      secure: location.protocol !== 'http:',
    });

    tokenPromise = null;
  }
}

export async function useToken<T>(
  fetchFunction: () => Promise<T>
): ReturnType<typeof fetchFunction> {
  const originRefreshToken = Cookie.get(COOKIE_REFRESH_TOKEN);
  let latestRefreshToken: string | undefined;
  try {
    return await fetchFunction();
  } catch (error) {
    if (
      !(
        error instanceof ResponseError &&
        error.code === ERROR_CODE_TOKEN_EXPIRED
      )
    ) {
      throw error;
    }

    latestRefreshToken = Cookie.get(COOKIE_REFRESH_TOKEN);

    if (!latestRefreshToken) {
      bridgeAdaptor.login();
      throw error;
    }
  }

  if (originRefreshToken === latestRefreshToken && !isInElectron()) {
    try {
      await updateToken(latestRefreshToken);
    } catch (error) {
      bridgeAdaptor.login();
      throw error;
    }
  }

  return fetchFunction();
}

{
  const { get, post, put, delete: del } = api;

  api.get = <T>(...args: Parameters<typeof get>) => {
    return useToken<T>(() => {
      return get.apply(api, args) as Promise<T>;
    });
  };

  api.post = <T>(...args: Parameters<typeof post>) => {
    return useToken<T>(() => {
      return post.apply(api, args) as Promise<T>;
    });
  };

  api.put = <T>(...args: Parameters<typeof put>) => {
    return useToken<T>(() => {
      return put.apply(api, args) as Promise<T>;
    });
  };

  api.delete = <T>(...args: Parameters<typeof del>) => {
    return useToken<T>(() => {
      return del.apply(api, args) as Promise<T>;
    });
  };
}
