import { NoteLocationInfo, UserSettingInfo } from '../../api/setting';
import { useLocalStorage } from '@vueuse/core';
import merge from 'lodash-es/merge';
import pick from 'lodash-es/pick';
import { useStore } from '../../store';
import { DocumentsActionTypes, DocumentsMutationTypes } from '../../store/documents';
import { defaultCommonSettings } from './const';
import { computed } from 'vue';
import { PDF_READER } from '@/common/src/constants/storage-keys';

export enum SaveMode {
  local = 'local',
  store = 'store',
  remote = 'remote'
}

export default function useRemoteSettings() {
  // 通用属性配置存本地一份
  const localSettings = useLocalStorage<Required<UserSettingInfo>>(
    PDF_READER.SETTINGS,
    defaultCommonSettings,
    {
      listenToStorageChanges: true,
    }
  );
  const store = useStore();

  const saveRemoteUserSettings = (values: NoteLocationInfo) => {
    store.dispatch(`documents/${DocumentsActionTypes.SAVE_SETTING}`, values);
  };

  const updateStoreUserSetting = (values: UserSettingInfo) => {
    localSettings.value = merge(
      localSettings.value,
      pick(values, Object.keys(defaultCommonSettings))
    ) as Required<UserSettingInfo>;

    store.commit(`documents/${DocumentsMutationTypes.SET_SETTING}`, values);
  };

  const userSettings = computed(() => store.state.documents.userSettingInfo);

  const saveByMode = (settings: NoteLocationInfo, mode: SaveMode) => {
    if (mode === SaveMode.remote) {
      saveRemoteUserSettings(settings);
    } else if (mode === SaveMode.store) {
      updateStoreUserSetting(settings);
    }
  };

  return {
    localSettings,
    userSettings,
    saveRemoteUserSettings,
    updateStoreUserSetting,
    saveByMode,
  };
}
