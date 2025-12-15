import { computed } from 'vue';
import { isOwner } from '~/src/store';
import useRemoteSettings from './useRemoteUserSettings';

export interface GroupSettings {
  currentGroupId: string;
}

export default function useGroupSettings() {
  const { userSettings, saveRemoteUserSettings, updateStoreUserSetting } =
    useRemoteSettings();

  const groupSettings = computed(() => {
    return {
      currentGroupId: userSettings.value.currentGroupId,
    };
  });

  const setCurrentGroupId = (id: string) => {
    const save = isOwner.value
      ? saveRemoteUserSettings
      : updateStoreUserSetting;

    save({ currentGroupId: id });
  };

  return {
    groupSettings,
    setCurrentGroupId,
  };
}
