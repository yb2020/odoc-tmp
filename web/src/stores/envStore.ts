import isMobile from 'is-mobile';
import { defineStore } from 'pinia';
import { isInElectron, isElectronMode, isInOverseaseElectron } from '../util/env';
import viewerConfig from '../../viewer.properties.json'

export interface EnvState {
  isMobile: boolean;
  isElectron: boolean;
  isElectronMode: boolean;
  viewerConfig: Partial<typeof viewerConfig.overseas>
}

export const useEnvStore = defineStore('env', {
  state: (): EnvState => ({
    isMobile: isMobile(),
    isElectron: isInElectron(),
    isElectronMode: isElectronMode(),
    viewerConfig: {},
  }),
  actions: {
    initViewerConfig() {
      if (isInOverseaseElectron()) {
        this.viewerConfig = viewerConfig.overseas;
      } else {
        this.viewerConfig = viewerConfig.domestic
      }
    },
  },
});
