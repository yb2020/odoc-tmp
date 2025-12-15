import { message } from 'ant-design-vue';
import copyTextToClipboard from 'copy-text-to-clipboard';
import i18n from '../locals/i18n';

export const copyToPaste = (text: string) => {
  const success = copyTextToClipboard(text);
  if (success) {
    message.success(i18n.global.t('message.copyToPasteboard.success'));
  } else {
    message.warn(i18n.global.t('message.copyToPasteboard.failed'));
  }
};
