const scriptsProperties = [
  'type',
  'src',
  'htmlFor',
  'event',
  'charset',
  'async',
  'defer',
  'crossOrigin',
  'text',
  'onerror',
] as const;

type Params = {
  [K in (typeof scriptsProperties)[number]]: string;
};

export const loadScript = function (
  url: string,
  params?: Params & { lazyLoad?: boolean }
) {
  return new Promise((resolve) => {
    const script = document.createElement('script');
    script.type = 'text/javascript';

    script.onload = function () {
      resolve(true);
    };
    script.onerror = function () {
      resolve(false);
    };

    if (typeof params === 'object') {
      for (const key in params) {
        if (
          Object.prototype.hasOwnProperty.call(params, key) &&
          scriptsProperties.includes(key as keyof Params)
        ) {
          (script as any)[key] = params[key as keyof Params];
        }
      }
    }

    //当设置script的text属性时，需要在加载script DOM节点前设置src属性，否则会执行text文本内容
    script.src = url;
    document
      .getElementsByTagName(params?.['lazyLoad'] ? 'body' : 'head')[0]
      .appendChild(script);
  });
};

export const loadCss = (cssId: string, url: string) => {
  if (!document.getElementById(cssId)) {
    const head = document.getElementsByTagName('head')[0];
    const link = document.createElement('link');
    link.id = cssId;
    link.rel = 'stylesheet';
    link.type = 'text/css';
    link.href = url;
    link.media = 'all';
    head.appendChild(link);
  }
};

interface Options {
  paperId: string;
  pdfId: string;
  pageType: string;
}

interface Handlers {
  addVisibleChangeHandler?: (visible: boolean) => void;
  addUpdateSuccessHandler?: () => void;
}

declare global {
  interface ShadomDiv extends HTMLDivElement {
    getDiv: () => HTMLDivElement;
    onReady: (fn: (ready: boolean) => void) => void;
  }

  interface Window {
    readpaperUpdateMeta: {
      init: (
        options: Options,
        rootElement: HTMLDivElement | string,
        handlers?: Handlers
      ) => {
        toggle: (visible: boolean, local: string) => void;
        unmount: () => void;
      };
      createShadow: (id?: string, rootElement?: HTMLDivElement) => ShadomDiv;
    };
  }
}

export const loadlib = (
  container: HTMLDivElement,
  props: Options,
  handlers: Handlers
) => {
  let instance: ReturnType<(typeof window)['readpaperUpdateMeta']['init']>;
  let shadowDom: ShadomDiv;

  return async () => {
    if (instance) {
      return {
        instance,
        destroy() {
          if (instance) {
            instance.unmount();
          }
          if (shadowDom) {
            shadowDom.remove();
          }
        },
      };
    }
    const libId = 'readpaperUpdateMeta' as const;
    if (!window[libId]) {
      const success = await loadScript(
        'https://nuxt.cdn.readpaper.com/readpaper-ai/lib/meta/readpaper-update-meta.umd.v2.js?t=20240321',
        undefined
      );
      if (!success) {
        return null;
      }
    }
    shadowDom = (
      window[libId] as (typeof window)['readpaperUpdateMeta']
    ).createShadow('update-meta', container);
    await new Promise((resolve) => {
      shadowDom.onReady((ready) => {
        resolve(ready);
      });
    });
    instance = (window[libId] as (typeof window)['readpaperUpdateMeta']).init(
      props,
      shadowDom.getDiv(),
      handlers
    );
    return {
      instance,
      destroy() {
        if (instance) {
          instance.unmount();
        }
        if (shadowDom) {
          shadowDom.remove();
        }
      },
    };
  };
};
