import { defineStore } from 'pinia';

export interface LimitDialogState {
  visible: boolean;
  code: number;
  message: string;
}

export const useLimitDialogStore = defineStore('limitDialog', {
  state: (): LimitDialogState => ({
    visible: false,
    code: 0,
    message: '',
  }),
  
  actions: {
    /**
     * 显示限制提示弹窗
     * @param code 状态码
     * @param message 错误信息
     */
    show(code: number, message: string) {
      this.code = code;
      this.message = message;
      this.visible = true;
    },
    
    /**
     * 关闭限制提示弹窗
     */
    close() {
      this.visible = false;
    },
  },
});
