import { getDomainOrigin } from './env';

// 检测是否在 Tauri 环境
const isTauri = () => typeof window !== 'undefined' && '__TAURI__' in window;

export const goPathPage = (path: string) => {
  const host = /https?:\/\//.test(path) ? '' : getDomainOrigin();

  clickUrl(`${host}${path}`);
};

const clickUrl = async (url: string, currentWindow = false) => {
  // Tauri 环境：在当前窗口导航（后续可改为 Tab 管理）
  if (isTauri()) {
    try {
      const urlObj = new URL(url);
      window.location.href = urlObj.pathname + urlObj.search + urlObj.hash;
    } catch {
      window.location.href = url;
    }
    return;
  }

  // Web 环境：使用 a 标签
  const a = document.createElement('a');

  a.setAttribute('href', url);

  if (!currentWindow) {
    a.setAttribute('target', '_blank');
  }

  a.setAttribute('id', 'open-new-page');

  document.body.appendChild(a);

  a.click();

  a.remove();
};

export const goPdfPage = (params: Record<string, string>) => {
  const searchParams = new URLSearchParams(params);
  const origin = getDomainOrigin();
  // const path = '/pdf-annotate/note';
  const path = '/note';
  const url = `${origin}${path}?${searchParams}`;
  clickUrl(url);
};

export function downloadUrl(blobUrl: string, filename: string) {
  const url = new URL(blobUrl);
  url.searchParams.set('attname', '');
  const a = document.createElement('a');
  if (!a.click) {
    throw new Error('DownloadManager: "a.click()" is not supported.');
  }
  a.href = url.href;
  a.target = '_blank';
  // Use a.download if available. This increases the likelihood that
  // the file is downloaded instead of opened by another PDF plugin.
  if ('download' in a) {
    a.download = filename;
  }
  // <a> must be in the document for recent Firefox versions,
  // otherwise .click() is ignored.
  (document.body || document.documentElement).append(a);
  a.click();
  a.remove();
}

export const goPersonPage = (userId: string) => {
  goPathPage(`${getDomainOrigin()}/user/${userId}`);
};

export const goSummaryPage = (params: Record<string, string>) => {
  const searchParams = new URLSearchParams(params);
  const origin = getDomainOrigin();
  // const path = '/pdf-annotate/summary.html';
  const path = '/summary.html';
  const url = `${origin}${path}?${searchParams}`;
  clickUrl(url);
};
