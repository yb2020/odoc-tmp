import { ref } from 'vue';
import { RightSideBarType } from '~/src/components/Right/TabPanel/type';
import { TAB_PARAM } from '~/src/constants';
import { getQueryString } from '~/src/util/scroll';
import { UserSettingInfo } from '~/src/api/setting';
import { createSharedComposable } from '@vueuse/core';

const TAB_SETTINGS_FACTORY: {
  [k: string]: () => Partial<UserSettingInfo>;
} = {
  [RightSideBarType.Copilot.toLocaleLowerCase()]: () => ({
    rightWidth: window.innerWidth / 3,
    rightShow: true,
    rightTab: RightSideBarType.Copilot,
    scalePresetValue: 'page-width',
  }),
};

function usePresetTabSettings() {
  const presetTab = ref(getQueryString(TAB_PARAM));

  const presetTabConfig = ref(
    (() => {
      const factory =
        presetTab.value && presetTab.value in TAB_SETTINGS_FACTORY
          ? TAB_SETTINGS_FACTORY[presetTab.value]
          : null;

      return factory?.();
    })()
  );

  const clearPresetTab = () => {
    presetTab.value = null;
  };

  const setPresetTabConfig = (cfg: object) => {
    presetTabConfig.value = presetTabConfig.value
      ? {
          ...presetTabConfig.value,
          ...cfg,
        }
      : undefined;
  };

  return {
    presetTab,
    presetTabConfig,
    clearPresetTab,
    setPresetTabConfig,
  };
}

export const useSharedPresetTabSettings =
  createSharedComposable(usePresetTabSettings);
