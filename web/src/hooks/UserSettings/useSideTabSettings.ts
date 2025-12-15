import { computed, ref } from 'vue';
import { UserStatusEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';

import useRemoteSettings, { SaveMode } from './useRemoteUserSettings';
import { store, isOwner, pdfStatusInfo, currentGroupId } from '~/src/store';
import { defaultCommonSettings } from './const';
import { RightSideBarType } from '~/src/components/Right/TabPanel/type';
import { NoteLocationInfo, UserSettingInfo } from '~/src/api/setting';
import { LeftSideBarType } from '~/src/components/Left/type';
import { merge, pickBy, isNil } from 'lodash-es';
import { remove } from '@vue/shared';
import { DocumentsMutationTypes } from '~/src/store/documents';
import { useSharedPresetTabSettings } from './usePresetTab';
import { NoteSubTypes } from '~/src/store/note/types';
import { useTranslateLock } from '../useTranslation';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';

export type SideTabCommonSettings = {
  shown: boolean;
  tab: string;
  subTab?: number | string;
  width: number;
};

const pickSideTabUserSettings = (
  dir: 'left' | 'right',
  localSettings: Required<UserSettingInfo>
) => {
  if (dir === 'left') {
    return {
      shown: localSettings.sideBarShow,
      tab: localSettings.sideBarTab,
      width: Math.max(
        localSettings.sideBarWidth,
        defaultCommonSettings.sideBarWidth
      ),
    };
  }

  return {
    shown: localSettings.rightShow,
    tab: localSettings.rightTab,
    width: Math.max(
      localSettings.rightWidth,
      defaultCommonSettings.sideBarWidth
    ),
  };
};

const allTabPanes = {
  [RightSideBarType.Matirial]: 'viewer.info',
  [RightSideBarType.Note]: 'viewer.notes',
  [RightSideBarType.Question]: 'viewer.discussions',
  [RightSideBarType.Group]: 'viewer.teams',
  [RightSideBarType.Learn]: 'viewer.tasks',
  [RightSideBarType.Copilot]: 'viewer.aiCopilot',
  [RightSideBarType.Translate]: 'translate.translate',
};

export interface RightSideTabItem {
  key: RightSideBarType;
  title: string;
  shown: boolean;
  disabled?: boolean;
  sortable: boolean;
}

let beforeChangeTab: null | (() => void) = null;
export const executeAndSetBeforeChangeTab = (
  callback: typeof beforeChangeTab
) => {
  if (beforeChangeTab) {
    beforeChangeTab();
  }

  beforeChangeTab = callback;
};

export const removeBeforeChangeTab = (callback?: () => void) => {
  if (callback === beforeChangeTab) {
    beforeChangeTab = null;
  }
};

export function useRightSideTabSettings() {
  const { presetTabConfig, setPresetTabConfig } = useSharedPresetTabSettings();
  const { userSettings, saveByMode } = useRemoteSettings();

  const sideTabSettings = computed(() => {
    return pickSideTabUserSettings(
      'right',
      (presetTabConfig.value as Required<UserSettingInfo>) ?? userSettings.value
    );
  });

  const { showTranslateTabInRight } = useTranslateLock();

  const tabs = computed<RightSideTabItem[]>(() => {
    const tabBars = userSettings.value.rightTabBars;

    const tabPanes: {
      key: RightSideBarType;
      title: string;
      shown: boolean;
      disabled?: boolean;
      sortable: boolean;
    }[] = [];

    const disabledTabs = store.state.documents.disabledRightTabs;

    tabBars.forEach((tab) => {
      const title = allTabPanes[tab.key];

      if (disabledTabs.includes(tab.key)) {
        return;
      }

      // Hide Learn tab
      if (tab.key === RightSideBarType.Learn) {
        return;
      }

      if (
        !isOwner.value &&
        (tab.key === RightSideBarType.Group ||
          tab.key === RightSideBarType.Learn)
      ) {
        return;
      }

      tabPanes.push({
        key: tab.key,
        title,
        shown: !!tab.shown,
        disabled: false,
        sortable: true,
      });
    });

    store.commit(
      `documents/${DocumentsMutationTypes.SET_RIGHT_TAB_BARS}`,
      tabPanes
    );

    if (
      showTranslateTabInRight.value &&
      !tabPanes.find((tab) => tab.key === RightSideBarType.Translate)
    ) {
      tabPanes.unshift({
        key: RightSideBarType.Translate,
        title: allTabPanes[RightSideBarType.Translate],
        shown: true,
        disabled: false,
        sortable: false,
      });
    }

    return tabPanes;
  });

  const saveByUserType = (settings: NoteLocationInfo) => {
    saveByMode(
      settings,
      pdfStatusInfo.value.noteUserStatus === UserStatusEnum.TOURIST
        ? SaveMode.store
        : SaveMode.remote
    );
  };

  const toggleTab = (key: RightSideBarType) => {
    executeAndSetBeforeChangeTab(null);
    const tabBars = tabs.value;
    const selected = tabBars.filter((item) => item.shown && !item.disabled);
    if (selected.length === 1 && key === selected[0].key) {
      return;
    }
    const tabIdx = tabBars.findIndex((item) => item.key === key);

    const tab = tabBars[tabIdx];

    tab.shown = !tab.shown;

    let activeTab = userSettings.value.rightTab;

    const shownLength = tabBars.filter((item) => item.shown).length;

    if (!tab.shown && activeTab === key) {
      activeTab = tabBars[(tabIdx + shownLength - 1) % shownLength].key;
    }

    saveByUserType({
      rightTabBars: tabBars,
      rightTab: activeTab,
    });
  };

  const setSideTabSetting = (values: Partial<SideTabCommonSettings>) => {
    executeAndSetBeforeChangeTab(null);

    const settings = {
      rightTab: values.tab as RightSideBarType,
      rightSubTab: values.subTab,
      rightWidth: values.width,
      rightShow: values.shown,
    };
    if (presetTabConfig.value) {
      setPresetTabConfig(pickBy(settings, (x) => !isNil(x)));
    } else {
      saveByUserType(settings);
    }
  };

  const switchAndShowTab = (
    tab: RightSideBarType,
    values: Omit<Partial<SideTabCommonSettings>, 'tab' | 'shown'>
  ) => {
    setSideTabSetting({
      tab,
      shown: true,
      ...values,
    });

    const tabShow = tabs.value.find((t) => t.key === tab);

    if (!tabShow?.shown) {
      toggleTab(tab);
    }
  };

  const activeTab = computed(() => {
    const one = tabs.value.find(
      (tab) =>
        tab.key.toLowerCase() === presetTabConfig.value?.rightTab?.toLowerCase()
    );

    if (one) {
      return one.key;
    }

    const tab = tabs.value.find(
      (tab) =>
        tab.key === userSettings.value.rightTab && !tab.disabled && tab.shown
    );

    if (tab) {
      return tab.key;
    }

    const shown = tabs.value.find((tab) => tab.shown && !tab.disabled);

    if (shown) {
      return shown.key;
    }

    const tabBars = userSettings.value.rightTabBars;
    const materialTab = tabBars.find(
      (tab) => tab.key === RightSideBarType.Matirial
    )!;

    materialTab.shown = true;

    const { rightTabBars } = userSettings.value;
    saveByUserType({ rightTabBars });

    return RightSideBarType.Matirial;
  });

  const activeSubTab = computed(
    () =>
      presetTabConfig.value?.rightSubTab ??
      userSettings.value.rightSubTab ??
      NoteSubTypes.Summary
  );

  const toggleDisabledTab = (tab: RightSideBarType, disabled: boolean) => {
    executeAndSetBeforeChangeTab(null);

    const disabledTabs = store.state.documents.disabledRightTabs;
    if (!disabled) {
      remove(disabledTabs, tab);
    } else {
      disabledTabs.push(tab);
    }
  };

  const updateRightTabBars = (values: RightSideTabItem[]) => {
    executeAndSetBeforeChangeTab(null);

    const bars = values.map((item) => {
      return {
        key: item.key,
        shown: item.shown,
      };
    });
    saveByUserType({
      rightTabBars: bars,
    });
  };

  const checkNoteTabVisible = () => {
    return (
      sideTabSettings.value.shown &&
      (activeTab.value === RightSideBarType.Group ||
        (activeTab.value === RightSideBarType.Note &&
          activeSubTab.value === NoteSubTypes.Annotation))
    );
  };

  const checkSwitchGroupTab = () => {
    if (
      currentGroupId.value !== SELF_NOTEINFO_GROUPID &&
      activeTab.value !== RightSideBarType.Group
    ) {
      switchAndShowTab(RightSideBarType.Group, {});
    }
  };

  return {
    tabs,
    toggleTab,
    sideTabSettings,
    setSideTabSetting,
    switchAndShowTab,
    activeTab,
    activeSubTab,
    toggleDisabledTab,
    updateRightTabBars,
    checkNoteTabVisible,
    checkSwitchGroupTab,
  };
}

export function useLeftSideTabSettings(mode: SaveMode) {
  const { userSettings, saveByMode } = useRemoteSettings();

  const sideTabSettings = ref<SideTabCommonSettings>(
    pickSideTabUserSettings('left', userSettings.value)
  );

  const setSideTabSetting = (values: Partial<SideTabCommonSettings>) => {
    sideTabSettings.value = merge(sideTabSettings.value, values);
    const settings = {
      sideBarTab: values.tab as LeftSideBarType,
      sideBarWidth: values.width,
      sideBarShow: values.shown,
    };
    saveByMode(settings, mode);
  };

  return {
    sideTabSettings,
    setSideTabSetting,
  };
}
