import { isInElectron } from '@common/utils/env';
import { ELECTRON_CHANNEL_EVENT_OPEN_URL } from '@common/electron/bridge';
import { getAuthorizationCode, getAuthorizationUrl } from '@common/api/auth';
import type { Router } from 'vue-router';
import semver from 'semver';

export const getVersion = (ua: string) => {
  // eslint-disable-next-line no-useless-escape
  const m = ua.match(/readpaper\/([0-9\.]+)/gi);
  if (m) {
    const version = m[0].split('/')[1];
    return version;
  }
  return '';
};

export const openProjectPage = async (
  router: Router,
  params: { id: string },
  method: 'resolve' | 'push' | 'replace'
) => {
  // const href = router.resolve({
  //   path: '/',
  //   query: params,
  // }).href

  const href = `/?id=${params.id}`;

  const route = /^polish/.test(window.location.host)
    ? `https://${window.location.host}${href}`
    : `https://polish.${window.location.host}${href}`;

  const ver = getVersion(navigator.userAgent);
  if (isInElectron() && ver && semver.gte(ver, '1.25.2')) {
    const code = await getAuthorizationCode();
    const url = getAuthorizationUrl({
      authorizationCode: code,
      redirectUrl: route,
    });
    window.electron.readpaperBridge.invoke(ELECTRON_CHANNEL_EVENT_OPEN_URL, {
      url,
    });
  } else if (method === 'resolve') {
    window.open(route, '_blank');
  } else {
    router[method](
      router.resolve({
        path: '/',
        query: params,
      })
    );
  }
};

export function openBeginnerGuide(isWebEn?: boolean) {
  if (isWebEn) {
    window.open('https://readpaper.com/new');
    return;
  }
  window.open('https://docs.qq.com/doc/DRWhaQURJS3lHTXps');
}
