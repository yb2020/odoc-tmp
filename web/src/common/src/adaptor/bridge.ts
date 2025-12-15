// import { useUserStore } from '@common/stores/user'

import { getDomainOrigin, isInElectron } from '../utils/env';

export const ELECTRON_CHANNEL_NAME = 'electron-client';
export const ELECTRON_CHANNEL_NAME_LOGIN = `${ELECTRON_CHANNEL_NAME}-login`;

class BridgeAdaptor {
  private isElectron = isInElectron();
  login(redirectUrl?: string) {
    if (this.isElectron) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (window as any).electron?.readpaperBridge?.invoke(
        ELECTRON_CHANNEL_NAME_LOGIN,
        {
          method: 'openLogin',
        }
      );
    } else if (redirectUrl) {
      if (window.location.host === 'polish.readpaper.com') {
        window.location.replace(
          `https://readpaper.com/login?redirect_url=${encodeURIComponent(
            redirectUrl
          )}`
        );
        return;
      }
      window.location.replace(
        `${getDomainOrigin()}/login?redirect_url=${encodeURIComponent(
          redirectUrl
        )}`
      );
    } else {
      // const userStore = useUserStore();
      // userStore.openLogin();
    }
  }
}

export const bridgeAdaptor = new BridgeAdaptor();
