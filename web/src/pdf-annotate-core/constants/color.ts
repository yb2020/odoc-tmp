import { AnnotationColor } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';

export const colorMap: Record<AnnotationColor, string> = {
  [AnnotationColor.blue]: '#1F71E0',
  [AnnotationColor.green]: '#73D13D',
  [AnnotationColor.yellow]: '#FFC53D',
  [AnnotationColor.red]: '#E66045',
  [AnnotationColor.black]: '#262625',
  [AnnotationColor.other]: '',
  [AnnotationColor.UNRECOGNIZED]: '',
};

export const colorKeyMap: Record<string, AnnotationColor> = {
  [colorMap[AnnotationColor.blue]]: AnnotationColor.blue,
  [colorMap[AnnotationColor.green]]: AnnotationColor.green,
  [colorMap[AnnotationColor.yellow]]: AnnotationColor.yellow,
  [colorMap[AnnotationColor.black]]: AnnotationColor.black,
  [colorMap[AnnotationColor.red]]: AnnotationColor.red,
};

export const colorList = [
  AnnotationColor.blue,
  AnnotationColor.yellow,
  AnnotationColor.red,
  AnnotationColor.green,
  AnnotationColor.black,
];

export const colorNameMap: Record<AnnotationColor, string> = {
  [AnnotationColor.blue]: '蓝色',
  [AnnotationColor.green]: '绿色',
  [AnnotationColor.yellow]: '黄色',
  [AnnotationColor.black]: '黑色',
  [AnnotationColor.red]: '红色',
  [AnnotationColor.other]: '',
  [AnnotationColor.UNRECOGNIZED]: '',
};

export const colorI18nMap: Record<AnnotationColor, string> = {
  [AnnotationColor.blue]: 'viewer.blue',
  [AnnotationColor.green]: 'viewer.green',
  [AnnotationColor.yellow]: 'viewer.yellow',
  [AnnotationColor.black]: 'viewer.black',
  [AnnotationColor.red]: 'viewer.red',
  [AnnotationColor.other]: '',
  [AnnotationColor.UNRECOGNIZED]: '',
};

export const ANNOTATION_COLOR_OPACITY = 0.3;
export const SCREENSHOT_COLOR_OPACITY = 0.1;

export const TRANSPARENT: CSSStyleDeclaration['color'] = 'transparent';
