import hotkeys, { KeyHandler } from 'hotkeys-js';
import { onUnmounted, watchEffect, ComputedRef } from 'vue';

declare module 'hotkeys-js' {
  interface Hotkeys {
    getPressedKeyString(): string[];
  }
}

export type HotKeysOptions = Parameters<typeof hotkeys>[1];

export const hk = hotkeys.noConflict();
export const hkfilter = hk.filter;
export const hkfilters: Array<typeof hk.filter> = [];

hk.filter = (e: KeyboardEvent) => {
  if (hkfilters.length && hkfilters.some((f) => f?.(e) === false)) {
    return false;
  }

  return hkfilter(e);
};

export default function useShortcuts(
  shortcuts: ComputedRef<string>,
  handler: KeyHandler,
  opts: ComputedRef<HotKeysOptions>
) {
  const enableShortcut = () => hk(shortcuts.value, opts.value, handler);
  const disableShortcut = () => hk.unbind(shortcuts.value, handler);

  watchEffect((onCleanup) => {
    enableShortcut();
    onCleanup(() => {
      disableShortcut();
    });
  });

  onUnmounted(() => {
    hk.unbind(shortcuts.value, opts.value.scope || 'all');
  });

  return {
    trigger: hk.trigger,
  };
}
