import isMobile from 'is-mobile';
import semver from 'semver';

export const isInElectron = () => {
  if (typeof window === 'undefined') {
    return false;
  }
  return !!window.navigator.userAgent.match(/readpaper/i);
};


export const isInTauri = () => {
  if (typeof window === 'undefined') {
    return false;
  }
  return typeof window !== 'undefined' && '__TAURI__' in window;
};

/**
 * @description TauriMode指的是是否在 Tauri 环境
 */
export const isTauriMode = () => {
  return isInTauri() && document.referrer === '';
};

/**
 * @description ElectronMode指的是是否在客户端阅读模式，与客户端搜索、小组打开笔记页区分
 */
export const isElectronMode = () => {
  return isInElectron() && document.referrer === '';
};

export const getDomainOrigin = () => {
  return window.location.origin;
};

export const getPdfAnnotateNoteUrl = (params: {
  pdfId?: string;
  noteId?: string;
  groupId?: string;
}) => {
  const url = new URL(getDomainOrigin());
  url.pathname = '/note';
  url.search = String(new URLSearchParams(params));
  return url;
};

export const getHostname = () => {
  return window.location.hostname;
};

export const IS_MOBILE = isMobile();
export const IS_ELECTRON = isInElectron();
export const IS_ELECTRON_MODE = isElectronMode();

export const isInOverseaseElectron = () => {
  return navigator.userAgent.indexOf('ReadPaperI18n') >= 0;
};

export const gteElectronVersion = (version: string) => {
  const electronVersion = navigator.userAgent.match(
    /ReadPaper\/(\d+\.\d+\.\d+)/
  );
  if (!electronVersion) {
    return false;
  }
  return semver.gte(electronVersion[1], version);
};
