import { ref } from 'vue';
import { NoteLocationInfo } from '../../api/setting';
import merge from 'lodash-es/merge';
import useRemoteSettings, { SaveMode } from './useRemoteUserSettings';
import { useSharedPresetTabSettings } from './usePresetTab';

export interface PDFWebviewSettings {
  currentPage: number;
  scale: number;
  scalePresetValue: string;
}

export default function usePDFWebviewSettings(isSelfNote?: boolean) {
  const { presetTabConfig, setPresetTabConfig } = useSharedPresetTabSettings();
  const { userSettings, saveByMode } = useRemoteSettings();

  /**
   * 这个不能是全局的，因为不同的pdf需要不同的setting
   */
  const pdfWebviewSettings = ref<PDFWebviewSettings>({
    currentPage: isSelfNote ? userSettings.value.currentPage : 1,
    scale: userSettings.value.scale,
    scalePresetValue:
      presetTabConfig.value?.scalePresetValue ??
      userSettings.value.scalePresetValue,
  });

  const save = (settings: NoteLocationInfo, mode: SaveMode) => {
    pdfWebviewSettings.value = merge(pdfWebviewSettings.value, settings);
    saveByMode(settings, mode);
  };

  const setScaleSetting = (
    values: { scale: number; presetValue?: string },
    mode: SaveMode
  ) => {
    const settings = {
      scale: values.scale,
      scalePresetValue: values.presetValue,
    };
    if (presetTabConfig.value) {
      setPresetTabConfig(settings);
    } else {
      save(settings, mode);
    }
  };

  const setCurrentPageSetting = (page: number, mode: SaveMode) => {
    save({ currentPage: page }, mode);
  };

  return {
    pdfWebviewSettings,
    setScaleSetting,
    setCurrentPageSetting,
  };
}
