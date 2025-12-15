import { defineStore } from 'pinia';
import { useUserStore } from '@common/stores/user';

export interface CopilotState {
  gptGrayTip: string;
  currentComponent: string;
  components: string[];
}

export const useCopilotStore = defineStore('copilot', {
  state: (): CopilotState => ({
    gptGrayTip: '',
    currentComponent: 'summary',
    components: ['summary','chat'],
  }),
  getters: {
    accessAiCopilot: () => {
      const userStore = useUserStore();
      const credits = userStore.getTotalCredits()
      const isEnable = userStore.getCurrentMembershipPermission()?.ai?.copilot?.isEnable

      if(!credits || !isEnable) { // 用于判断是否有值
        return false;
      }

      return isEnable && credits > 0;
    },
    isSummary: (state) => state.currentComponent === 'summary',
    isChat: (state) => state.currentComponent === 'chat',
  },
  actions: {
    setGptInfo(gptGrayTip: string) {
      this.gptGrayTip = gptGrayTip;
    },
    gotoSummary() {
      this.gotoComponent('summary');
    },
    gotoChat() {
      this.gotoComponent('chat');
    },
    gotoComponent(component: string) {
      if (!this.components.includes(component)) {
        console.error('Invalid component:', component);
        return;
      }
      this.currentComponent = component;
    },
  },
});
