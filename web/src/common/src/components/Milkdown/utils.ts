import { MilkdownPlugin } from '@milkdown/core';
import { AtomList } from '@milkdown/utils';
import {
  ConfigItem,
  defaultConfig as menuDefaultConfig,
} from '@milkdown/plugin-menu';
import { math } from '@milkdown/plugin-math';
import { block } from '@milkdown/plugin-block';
import { history } from '@milkdown/plugin-history';
import { clipboard } from './plugins/clipboard';
import { placeholder } from 'milkdown-plugin-placeholder';
import { highlight } from './plugins/highlight';

import 'milkdown-plugin-placeholder/styles/index.css';

export const getMenuConfig = (lang = '', ignoreKeys: string[] = []) => {
  const menuConfig = menuDefaultConfig.map((arr) => {
    return arr.reduce((slice, item) => {
      if (
        item.type === 'button' &&
        (['Undo', 'Redo'].includes(item.key as string) ||
          ignoreKeys.includes(item.key as string))
      ) {
        return slice;
      }

      if (
        /^zh/.test(lang) &&
        item.type === 'select' &&
        item.text === 'Heading'
      ) {
        item = {
          ...item,
          text: '标题',
          options: [
            { id: '1', text: '一级标题' },
            { id: '2', text: '二级标题' },
            { id: '3', text: '三级标题' },
            { id: '0', text: '正文' },
          ],
        };
      }

      slice.push(item);

      return slice;
    }, [] as ConfigItem[]);
  });

  return menuConfig;
};

export const DefaultPlugins: Record<string, MilkdownPlugin | AtomList> = {
  math,
  block,
  history,
  clipboard,
  highlight,
  placeholder,
};

export const getDefaultPlugins = (plugins = DefaultPlugins) => {
  return Object.values(plugins);
};
