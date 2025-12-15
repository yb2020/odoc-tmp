import { computed } from 'vue';
import { createGlobalState } from '@vueuse/core';
import { isOwner } from '~/src/store';
import useRemoteSettings from './useRemoteUserSettings';
import _ from 'lodash';
import { UserSettingInfo } from '~/src/api/setting';

export type ToolBarSettings = Pick<
  UserSettingInfo,
  'toolBarVisible' | 'toolBarHeadVisible' | 'toolBarNoteVisible'
>;

export function useToolBarSettings() {
  const { userSettings, saveRemoteUserSettings, updateStoreUserSetting } =
    useRemoteSettings();

  const toolbarSettings = computed(() => {
    return _.pick(userSettings.value, [
      'toolBarVisible',
      'toolBarHeadVisible',
      'toolBarNoteVisible',
    ]);
  });

  const setToolbarSettings = (x: ToolBarSettings) => {
    const save = isOwner.value
      ? saveRemoteUserSettings
      : updateStoreUserSetting;

    save(x);
  };
  
  return {
    toolbarSettings,
    setToolbarSettings,
  };
}

export default createGlobalState(useToolBarSettings);
