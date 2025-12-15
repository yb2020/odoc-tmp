import {
  Module,
  Store as VuexStore,
  DispatchOptions,
  CommitOptions,
} from 'vuex';
import { ToolBarType } from '@idea/pdf-annotate-core';
import { RootState } from '../types';
import { DocumentsState } from './type';
import { Mutations, mutations } from './mutations';
import { Actions, actions } from './actions';
import { defaultCommonSettings } from '~/src/hooks/UserSettings/const';

export { DocumentsMutationTypes } from './mutations';
export { DocumentsActionTypes } from './actions';

export type DocumentsStore<S = DocumentsState> = Omit<
  VuexStore<S>,
  'dispatch' | 'commit'
> & {
  commit<K extends keyof Mutations, P extends Parameters<Mutations[K]>[1]>(
    key: K,
    payload: P,
    options?: CommitOptions
  ): ReturnType<Mutations[K]>;
} & {
  dispatch<K extends keyof Actions>(
    key: K,
    payload: Parameters<Actions[K]>[1],
    options?: DispatchOptions
  ): ReturnType<Actions[K]>;
};

export const DocumentsModule: Module<DocumentsState, RootState> = {
  namespaced: true,
  mutations,
  actions,
  state: () => {
    return {
      isFullPage: false,
      toolBarType: ToolBarType.None,
      userSettingInfo: defaultCommonSettings,
      rightSideTabSettings: null,
      disabledRightTabs: [],
    };
  },
};
