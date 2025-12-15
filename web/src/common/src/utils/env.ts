import isMobile from 'is-mobile';

export const isInElectron = () => {
  if (typeof window === 'undefined') {
    return false;
  }
  return !!window.navigator.userAgent.match(/readpaper/i);
};

export const isElectronMode = () => {
  return isInElectron();
};

export const getDomainOrigin = () => {
  if (!isInElectron()) {
    return window.location.origin;
  }

  return `https://${getWebHost()}`;
};

export const getWebHost = () => {
  if (import.meta.env.VITE_API_ENV === 'dev') {
    return 'paper.dev.aiteam.cc';
  }

  if (import.meta.env.VITE_API_ENV === 'uat') {
    return 'paper.uat.aiteam.cc';
  }

  return 'readpaper.com';
};

export const getHostname = getWebHost;

export const IS_MOBILE = isMobile();
export const IS_ELECTRON = isInElectron();
export const IS_ELECTRON_MODE = isElectronMode();

export const isInOverseaseElectron = () => {
  return navigator.userAgent.indexOf('ReadPaperI18n') >= 0;
};

export const isDev = process.env.NODE_ENV === 'development';
