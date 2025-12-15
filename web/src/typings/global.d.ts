/* eslint-disable @typescript-eslint/triple-slash-reference */
/// <reference path="../../../readpaper-ai/vite-global.d.ts" />

export type PLATFORM_ELECTRON = 'electron';
export type PLATFORM_WEB = 'web';

export type API_DEV = 'dev';
export type API_UAT = 'uat';
export type API_PROD = 'prod';
declare global {
  interface Window {
    store: any;
  }

  interface ImportMetaEnv {
    readonly VITE_DEBUG_PROD: string;
    readonly VITE_PLATFORM: PLATFORM_ELECTRON | PLATFORM_WEB;
    readonly VITE_API_ENV: API_DEV | API_UAT | API_PROD;
    // 更多环境变量...
  }

  interface ImportMeta {
    readonly env: ImportMetaEnv;
  }
}

declare type Nullable<T> = T | null;

declare interface Fn<T = any, R = T> {
  (...arg: T[]): R;
}
