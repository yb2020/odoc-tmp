export type PlatformKey = 'darwin' | 'win32';

export interface Shortcuts {
  darwin: string;
  win32: string;
}

process.platform;

export interface ShortcutInfo {
  [command: string]: {
    order?: number;
    icon: any;
    iconAttrs?: object;
    name: string;
    value: Shortcuts;
    i18n: string;
  };
}

export interface ShortcutsState {
  [k: string]: {
    scope?: string;
    // routes: Set<string>;
    shortcuts: ShortcutInfo;
  };
}
