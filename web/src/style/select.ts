import { ColorKey, ColorStyle } from '@common/components/Notes/types';

export * from '@common/components/Notes/types';

export type NoteFilter = ColorKey | 'ref';

export const rectStyleMap: Record<ColorKey, ColorStyle> = {
  [ColorKey.blue]: {
    fill: '#338AFF',
    fillOpacity: 0.1,
    color: '#338AFF',
    text: '蓝色',
    i18n: 'viewer.blue',
  },
  [ColorKey.green]: {
    fill: '#80E639',
    fillOpacity: 0.1,
    color: '#80E639',
    text: '绿色',
    i18n: 'viewer.green',
  },
  [ColorKey.yellow]: {
    fill: '#FFFF00',
    fillOpacity: 0.1,
    color: '#FFFF00',
    text: '黄色',
    i18n: 'viewer.yellow',
  },
  [ColorKey.orange]: {
    fill: '#FF8C19',
    fillOpacity: 0.1,
    color: '#FF8C19',
    text: '橙色',
    i18n: 'viewer.orange',
  },
  [ColorKey.red]: {
    fill: '#F24030',
    fillOpacity: 0.1,
    color: '#F24030',
    text: '红色',
    i18n: 'viewer.red',
  },
};
