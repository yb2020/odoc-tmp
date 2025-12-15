import axios from 'axios';
import { api } from './axios';
import { useLimitDialogStore } from '@/common/src/stores/limitDialog';

// 需要拦截并显示限制提示弹窗的状态码范围（4000-4999）
const isLimitStatusCode = (status: number) => status >= 4000 && status <= 4999;

/**
 * 处理响应数据，检查是否需要显示限制提示弹窗
 */
const handleResponse = (response) => {
  // 处理业务状态码在响应体中的情况
  const data = response.data;
  
  // 添加日志输出，调试拦截器
  //console.log('[Interceptor] Response data:', data);
  
  if (data && typeof data === 'object' && 'status' in data) {
    const status = data.status;
    //console.log('[Interceptor] Status code:', status);
    
    // 检查是否是需要拦截的状态码范围
    if (isLimitStatusCode(status)) {
      console.log('[Interceptor] Limit status code detected:', status);
      try {
        const limitDialogStore = useLimitDialogStore();
        console.log('[Interceptor] Store accessed:', !!limitDialogStore);
        limitDialogStore.show(status, data.message || '操作受限，请升级会员或购买积分包');
        console.log('[Interceptor] Dialog shown');
        
        // 在响应数据中添加limitHandled标记，表示已被拦截器处理
        if (data) {
          data.limitHandled = true;
        }
      } catch (e) {
        console.error('[Interceptor] Error showing dialog:', e);
      }
    }
  }
  
  return response;
};

/**
 * 处理错误响应，检查是否需要显示限制提示弹窗
 */
const handleError = (error) => {
  // 处理 HTTP 错误
  if (error.response && error.response.data) {
    const { data } = error.response;
    
    // 如果响应体中包含业务状态码
    if (data && typeof data === 'object' && 'status' in data) {
      const businessStatus = data.status;
      
      // 检查是否是需要拦截的状态码范围
      if (isLimitStatusCode(businessStatus)) {
        const limitDialogStore = useLimitDialogStore();
        limitDialogStore.show(businessStatus, data.message || '操作受限，请升级会员或购买积分包');
        
        // 在响应数据中添加limitHandled标记，表示已被拦截器处理
        if (data) {
          data.limitHandled = true;
        }
      }
    }
  }
  
  return Promise.reject(error);
};

/**
 * 配置 Axios 拦截器
 */
export const setupInterceptors = () => {
  // 全局 axios 拦截器
  axios.interceptors.response.use(handleResponse, handleError);
  
  // 本地 api 实例拦截器
  api.interceptors.response.use(handleResponse, handleError);
};
