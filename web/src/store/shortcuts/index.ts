import {
  FullscreenOutlined,
  SearchOutlined,
  VerticalAlignTopOutlined,
} from '@ant-design/icons-vue';
import { Module } from 'vuex';
import { h } from 'vue';
import { PAGE_ROUTE_NAME } from '../../routes/type';
import { RootState } from '../types';
import { PlatformKey, ShortcutsState } from './type';

let platformKey: PlatformKey;
export const getPlatformKey = () => {
  const platform =
    (navigator as any)?.userAgentData?.platform || navigator?.platform || '';

  if (!platformKey) {
    platformKey = /Mac/i.test(platform) ? 'darwin' : 'win32';
  }

  return platformKey;
};

export const shortcutTxtMap: { [k: string]: string } = {
  Command: 'Cmd',
  Option: 'Opt',
};
export const getShortcutTxt = (s?: string) => {
  if (!s) {
    return null;
  }

  return h(
    'span',
    {
      style: 'margin-left:8px;',
    },
    s
      .split('+')
      .map((x) => (x in shortcutTxtMap ? shortcutTxtMap[x] : x))
      .join(' + ')
  );
};

export const ShortcutsModule: Module<ShortcutsState, RootState> = {
  namespaced: true,
  state: (): ShortcutsState => ({
    [PAGE_ROUTE_NAME.NOTE]: {
      shortcuts: {
        FULLSCREEN: {
          order: 30,
          icon: FullscreenOutlined,
          name: '全屏',
          value: {
            win32: 'F11',
            darwin: 'Command+Ctrl+F',
          },
          i18n: 'viewer.fullscreen',
        },
        TOGGLE_NOTE: {
          order: 60,
          icon: VerticalAlignTopOutlined,
          iconAttrs: {
            style: {
              transform: 'rotate(90deg)',
            },
          },
          name: '收起/展开笔记栏（右侧）',
          value: {
            win32: 'Alt+W',
            darwin: 'Option+W',
          },
          i18n: 'viewer.toggleRightbar',
        },
        TOGGLE_CATALOG: {
          order: 50,
          icon: VerticalAlignTopOutlined,
          iconAttrs: {
            style: {
              transform: 'rotate(-90deg)',
            },
          },
          name: '收起/展开目录栏（左侧）',
          value: {
            win32: 'Alt+Q',
            darwin: 'Option+Q',
          },
          i18n: 'viewer.toggleLeftbar',
        },
        PDF_SIZE_SELF_ADAPTION: {
          order: 40,
          icon: 'i',
          iconAttrs: {
            class: 'aiknowledge-icon icon-resize-autofit',
            'aria-hidden': 'true',
          },
          name: 'PDF大小自适应',
          value: {
            win32: 'Ctrl+0',
            darwin: 'Command+0',
          },
          i18n: 'viewer.zoomToPageWidth',
        },
        SEARCH: {
          order: 20,
          icon: SearchOutlined,
          name: '全文搜索',
          value: {
            win32: 'Ctrl+F',
            darwin: 'Command+F',
          },
          i18n: 'viewer.findInDocument',
        },
        NOTE_SCREENSHOT: {
          order: 10,
          icon: 'i',
          iconAttrs: {
            class: 'aiknowledge-icon icon-crop',
            'aria-hidden': 'true',
          },
          name: '截图',
          value: {
            win32: 'Alt+R',
            darwin: 'Option+R',
          },
          i18n: 'viewer.screenshot',
        },
        MULTI_SEGMENT: {
          order: 70,
          icon: 'i',
          iconAttrs: {
            class: 'aiknowledge-icon icon-sort-descending',
            'aria-hidden': 'true',
          },
          name: '多段选取',
          value: {
            win32: 'viewer.pressing+Ctrl',
            darwin: 'viewer.pressing+Command',
          },
          i18n: 'viewer.selectMultiSegment',
        },
      },
    },
  }),
};
