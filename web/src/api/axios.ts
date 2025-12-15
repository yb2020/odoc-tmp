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
// import { setAuthorization } from '../util/login';
import { v4 as uuidv4 } from 'uuid';
import { bridgeAdaptor } from '../adaptor/bridge';
import { isInElectron, getHostname, isTauriMode } from '../util/env';
import { isUnCert } from '../store/cert';
import { isObject } from '@vueuse/core';
import { getI18nCookie } from '../util/cookie';
import { getLanguageCookie } from '../shared/language/service';

/**
 * 获取API根路径
 * 根据环境自动选择：
 * - Tauri 本地版本 → localhost:8081
 * - Electron 客户端 → 云端 API
 * - Web 版本 → 相对路径（由 nginx/vite 代理）
 */
export const getAPIRoot = () => {
  if (isTauriMode()) {
    return 'http://localhost:8081/api';
  }
  
  if (
    !isInElectron() ||
    ['http:', 'https:'].includes(window.location.protocol)
  ) {
    return '/api';
  }

  return `https://${getHostname()}/api`;
};

const options: AxiosRequestConfig = {
  baseURL: getAPIRoot(),
  headers: {
    // Authorization: setAuthorization('aiKnowledge', 'readpaper_com'),
  },
};

export const api = axios.create(options);

api.interceptors.request.use(async (config) => {
  config.timeout = config.timeout || 10000;
  config.data = config.data || {};
  
  config.headers = config.headers || {};
  config.headers['x-traceId-header'] = uuidv4();

  // 设置 Accept-Language 头，使用标准格式（RFC 5646）
  const currentLang = getLanguageCookie();
  
  if (currentLang) {
    config.headers['Accept-Language'] = currentLang; // 直接使用标准格式 'en-US' 或 'zh-CN'
  } else {
    console.warn(`[axios拦截器] 未读取到有效的语言Cookie，将不设置Accept-Language头`);
  }

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

    // 检查是否是限制相关的状态码范围（4000-4999）
    const responseStatus = data.status;
    const isLimitStatusCode = responseStatus >= 4000 && responseStatus <= 4999;
    
    // 只有非限制状态码才显示错误消息
    if (!isLimitStatusCode && !response.config.headers?.[HEADER_CANCLE_AUTO_ERROR]) {
      message.error(data.message || UNKNOWN_ERROR_MESSAGE);
    }

    let code = data.code;
    if (!code && /无权限/.test(data.message || '')) {
      code = 403;
    }
    
    // 拦截 4000-4999 之间的状态码
    if (isLimitStatusCode) {
      try {
        // 动态导入 Pinia store，避免循环依赖
        const { useLimitDialogStore } = await import('@/common/src/stores/limitDialog');
        const limitDialogStore = useLimitDialogStore();
        limitDialogStore.show(responseStatus, data.message || '操作受限，请升级会员或购买积分包');
        console.log('[axios.ts] 显示限制提示弹窗，状态码:', responseStatus);
        
        // 返回一个特殊的响应对象，不再抛出异常
        // 这样上传组件将不会显示错误信息
        return {
          data: {
            status: responseStatus,
            message: data.message || '操作受限，请升级会员或购买积分包',
            data: {
              limitHandled: true,  // 标记该响应已经被限制处理逻辑处理过
              ...data.data,
            },
          },
          status: 200,  // 返回200状态码，避免触发错误处理
          statusText: 'OK',
          headers: response.headers,
          config: response.config,
        };
      } catch (e) {
        console.error('[axios.ts] 显示限制提示弹窗失败:', e);
      }
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
  const throwNotTokenExpiredError = (error: unknown) => {
    if (
      !(
        error instanceof ResponseError &&
        error.code === ERROR_CODE_TOKEN_EXPIRED
      )
    ) {
      throw error;
    }
  };

  try {
    return await fetchFunction();
  } catch (error) {
    throwNotTokenExpiredError(error);

    latestRefreshToken = Cookie.get(COOKIE_REFRESH_TOKEN);

    if (!latestRefreshToken) {
      bridgeAdaptor.login();
      throw '';
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

  try {
    return await fetchFunction();
  } catch (error) {
    throwNotTokenExpiredError(error);

    bridgeAdaptor.login();
    throw '';
  }
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
