/*
 * Created Date: March 17th 2022, 11:22:57 am
 * Author: zhoupengcheng
 * -----
 * Last Modified: March 17th 2022, 1:57:16 pm
 */
import { MutationTree } from 'vuex';
import { NoteLocationInfo } from '~/src/api/setting';
import { ToolBarType } from '@idea/pdf-annotate-core';
import { DocumentsState } from './type';
import { merge } from 'lodash-es'
import { RightSideBarType } from '~/src/components/Right/TabPanel/type';

interface RightTabBarParams {
  key: RightSideBarType;
  title: string;
  shown: boolean;
  disabled?: boolean;
}

export enum DocumentsMutationTypes {
  SET_TOOLBAR_TYPE = 'setToolbarType',
  SET_SETTING = 'setSetting',
  SET_FULL_PAGE = 'setFullPage',
  DISABLE_RIGHT_TAB = 'disableRightTab',
  SET_RIGHT_TAB_BARS = 'setRightTabBars'
}

export type Mutations<S = DocumentsState> = {
  [DocumentsMutationTypes.SET_TOOLBAR_TYPE](
    state: S,
    payload: ToolBarType
  ): void;
  [DocumentsMutationTypes.SET_SETTING](
    state: S,
    payload: NoteLocationInfo
  ): void;
  [DocumentsMutationTypes.SET_FULL_PAGE](state: S, payload: boolean): void;
  [DocumentsMutationTypes.DISABLE_RIGHT_TAB](state: S, payload: RightSideBarType | RightSideBarType[]): void;
};

export const mutations: MutationTree<DocumentsState> & Mutations = {
  [DocumentsMutationTypes.SET_TOOLBAR_TYPE](state, payload) {
    state.toolBarType = payload;
  },
  [DocumentsMutationTypes.SET_SETTING](state, payload) {
    state.userSettingInfo = merge(true, state.userSettingInfo, payload);
  },
  [DocumentsMutationTypes.SET_FULL_PAGE](state, payload) {
    state.isFullPage = payload;
  },
  [DocumentsMutationTypes.DISABLE_RIGHT_TAB](state, payload) {
    if (Array.isArray(payload)) {
      state.disabledRightTabs = state.disabledRightTabs.concat(payload)
    } else {
      state.disabledRightTabs.push(payload)
    }
  },
  [DocumentsMutationTypes.SET_RIGHT_TAB_BARS](state, payload) {
    state.userSettingInfo.rightTabBars = state.userSettingInfo.rightTabBars.filter((tab) => {
      const findIndex = payload.findIndex((item: RightTabBarParams) => tab.key === item.key)
      if(findIndex !== -1) return true

      return false
    })
  },
};
