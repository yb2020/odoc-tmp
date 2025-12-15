import ClipboardJS from 'clipboard';
import { message } from 'ant-design-vue';

export const initClipBoradButton = (cls: string | Element, info?: string) => {
  const clipboard = new ClipboardJS(cls);
  clipboard.on('success', function (e: any) {
    message.success(info || 'Copied!');
    e.clearSelection();
  });

  clipboard.on('error', function () {
    message.error('复制失败');
  });
};
