import { ActionTree, ActionContext } from 'vuex';
import { RootState } from '../types';
import { DocumentsState } from './type';
import {
  getSetting,
  setSetting,
  getLocation,
  setLocation,
  NoteLocationInfo,
} from '~/src/api/setting';
import { Mutations, DocumentsMutationTypes } from './mutations';
import { merge, uniqWith } from 'lodash-es';
import { RightSideBarType } from '~/src/components/Right/TabPanel/type';
import { defaultCommonSettings } from '~/src/hooks/UserSettings/const';
import { useCopilotStore } from '~/src/stores/copilotStore';
import { useLocalStorage } from '@vueuse/core';
import { LangType } from '~/src/stores/copilotType';
import { isOwner } from '~/src/store';

import { SELF_NOTEINFO_GROUPID } from '../base/type';
import { useEnvStore } from '~/src/stores/envStore';
import { useVipStore } from '@common/stores/vip';
import { UI } from '~/src/common/src/constants/storage-keys';

const copilotTabPositioned = useLocalStorage<boolean>(
  UI.COPILOT_TAB_POSITION,
  false
);

export enum DocumentsActionTypes {
  SAVE_SETTING = 'saveSetting',
  INIT_PDFVIEWER_WITH_SETTING = 'initPdfViewerWithSetting',
}

type AugmentedActionContext = {
  commit<K extends keyof Mutations>(
    key: K,
    payload: Parameters<Mutations[K]>[1]
  ): ReturnType<Mutations[K]>;
} & Omit<ActionContext<DocumentsState, RootState>, 'commit'>;

export interface Actions {
  [DocumentsActionTypes.SAVE_SETTING](
    context: AugmentedActionContext,
    payload: NoteLocationInfo
  ): void;
  [DocumentsActionTypes.INIT_PDFVIEWER_WITH_SETTING](
    context: AugmentedActionContext,
    payload: string
  ): Promise<NoteLocationInfo>;
}

export const actions: ActionTree<DocumentsState, RootState> & Actions = {
  [DocumentsActionTypes.SAVE_SETTING]({ commit, state, rootState }, payload) {
    if (!state.userSettingInfo) {
      return;
    }
    const setting = merge({}, state.userSettingInfo, payload);

    commit(DocumentsMutationTypes.SET_SETTING, { ...setting });

    const { currentPage, currentGroupId, ..._setting } = setting;

    const selfNoteInfo = rootState.base.noteInfoMap[SELF_NOTEINFO_GROUPID];
    // eslint-disable-next-line promise/catch-or-return
    Promise.all([
      setLocation({
        noteId: selfNoteInfo!.noteId,
        location: { currentPage, currentGroupId },
      }),
      setSetting(_setting),
    ]);
  },

  async [DocumentsActionTypes.INIT_PDFVIEWER_WITH_SETTING]({ commit }, noteId) {
    const [location, setting] = await Promise.all([
      getLocation({ noteId }),
      getSetting(),
    ]);

    const userSettingInfo = { ...setting, ...location };

    const rightTabBars = (userSettingInfo.rightTabBars || []).filter(
      (tab) =>
        !!defaultCommonSettings.rightTabBars.find(
          (item) => item.key === tab.key
        )
    );

    const tabBars = uniqWith(rightTabBars, (a, b) => a.key === b.key);

    if (tabBars.length < defaultCommonSettings.rightTabBars.length) {
      defaultCommonSettings.rightTabBars.forEach((item) => {
        const find = tabBars.find((tab) => tab.key === item.key);
        if (!find) {
          tabBars.push(item);
        }
      });
    }

    const queryGroupId = new URL(window.location.href).searchParams.get(
      'groupId'
    );

    let activeTabKey = userSettingInfo.rightTab || RightSideBarType.Matirial;

    const shownTab = tabBars.find(
      (item) => item.shown && item.key !== RightSideBarType.Group
    );

    const vipStore = useVipStore();
    vipStore.fetchVipConfig();
    if (isOwner.value) {
      vipStore.fetchVipProfile();
      vipStore.fetchOcrCount();
    }

    const copilotStore = useCopilotStore();
    if (!isOwner.value) {
      commit(
        DocumentsMutationTypes.DISABLE_RIGHT_TAB,
        RightSideBarType.Copilot
      );
    } else {
      userSettingInfo.copilotLanguage =
        userSettingInfo.copilotLanguage || LangType.CHINESE;
      if (!(userSettingInfo.copilotLanguage in LangType)) {
        userSettingInfo.copilotLanguage = LangType.CHINESE;
      }
    }

    const envStore = useEnvStore();
    commit(
      DocumentsMutationTypes.DISABLE_RIGHT_TAB,
      (envStore.viewerConfig.hiddenRightTabBars ?? []) as RightSideBarType[]
    );

    if (
      copilotStore.accessAiCopilot &&
      !copilotTabPositioned.value &&
      isOwner.value
    ) {
      activeTabKey = RightSideBarType.Copilot;

      copilotTabPositioned.value = true;
    } else if (queryGroupId) {
      activeTabKey = RightSideBarType.Group;
    } else if (activeTabKey === RightSideBarType.Group) {
      activeTabKey = shownTab?.key || RightSideBarType.Matirial;
    }

    const activeIdx = tabBars.findIndex((item) => item.key === activeTabKey);

    const activeTab = tabBars[activeIdx >= 0 ? activeIdx : 0];

    activeTab.shown = true;

    if (activeTab.key === RightSideBarType.Group) {
      userSettingInfo.rightShow = true;
    }

    userSettingInfo.rightTabBars = tabBars;

    userSettingInfo.rightTab = activeTab.key;

    if (queryGroupId) {
      userSettingInfo.currentGroupId = queryGroupId;
    }

    commit(DocumentsMutationTypes.SET_SETTING, userSettingInfo);

    return userSettingInfo;
  },
};
