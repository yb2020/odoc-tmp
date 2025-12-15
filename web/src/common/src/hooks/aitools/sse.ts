import {
  fetchEventSource,
  EventStreamContentType,
} from '@microsoft/fetch-event-source';
import { getLanguageCookie } from '../../../../shared/language/service';

class RetriableError extends Error {
  constructor(message?: string) {
    super(message || 'Connection Failed');
    this.name = 'RetriableError';
  }
}
class FatalError extends Error {}

const ssePostFetcher = <T, R>(
  url: string,
  data: T,
  onMessage: (r: R) => void,
  abortSignal: AbortController
) => {
  // 获取当前语言设置
  const currentLang = getLanguageCookie();  
  // 构建请求头
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  };
  
  // 添加 Accept-Language 头
  if (currentLang) {
    headers['Accept-Language'] = currentLang;
  } else {
    console.warn(`[ssePostFetcher] 未读取到有效的语言Cookie，将不设置Accept-Language头`);
  }
  
  return fetchEventSource(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(data),
    async onopen(response) {
      console.log('sse onopen message', response);
      if (
        response.ok &&
        response.headers.get('content-type') === EventStreamContentType
      ) {
        return; // everything's good
      } else if (
        response.status >= 400 &&
        response.status < 500 &&
        response.status !== 429
      ) {
        // client-side errors are usually non-retriable:
        throw new FatalError(`HTTP ${response.status} ${response.statusText}`);
      } else {
        throw new RetriableError(
          `ContentType Error: HTTP ${response.status} ${response.statusText} ${
            response.headers.get('content-type') || ''
          }`
        );
      }
    },
    onmessage(msg) {
      //console.log('sse onmessage', Date.now(), msg);
      // if the server emits an error message, throw an exception
      // so it gets handled by the onerror callback below:
      if (msg.event === 'FatalError') {
        throw new FatalError(msg.data);
      }
      onMessage(msg.data as any);
    },
    onclose() {
      console.log('sse onclose message');
      // if the server closes the connection unexpectedly, retry:
      abortSignal.abort();
    },
    onerror(err) {
      console.log('sse onerror message', err);
      throw err;
    },
    signal: abortSignal.signal,
    openWhenHidden: false,
  });
};

export default ssePostFetcher;
