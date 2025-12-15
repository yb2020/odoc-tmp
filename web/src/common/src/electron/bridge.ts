export const ELECTRON_CHANNEL_NAME = 'electron-client';
export const ELECTRON_CHANNEL_EVENT_OPEN_URL = `${ELECTRON_CHANNEL_NAME}-event-open-url`;

export const invoke = (k: string, params: unknown) => {
  window.electron?.readpaperBridge?.invoke?.(k, params);
};
