import { invert } from 'lodash-es';
import { SelectTextEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';

export enum ParagraphPluginCmd {
  level0 = 'none',
  level1 = 'improve',
  level2 = 'reduceSimilarity',
  level3 = 'shorten',
  level4 = 'expand',
  level5 = 'simple',
}

export enum ParagraphMode {
  standard = 'standard',
  shorten = 'shorten',
  expand = 'expand',
  simple = 'simple',
  // 优化
  improve = 'improve',
  // 降重
  reduceSimilar = 'reduceSimilar',
}

export interface ParagraphOpts {
  mode: ParagraphMode;
  synonyms: number;
}

export const ParagraphParameters = {
  [ParagraphPluginCmd.level0]: {
    mode: ParagraphMode.standard,
    synonyms: 1,
    cmd: ParagraphPluginCmd.level0,
  },
  [ParagraphPluginCmd.level1]: {
    mode: ParagraphMode.standard,
    synonyms: 1,
    cmd: ParagraphPluginCmd.level1,
  },
  [ParagraphPluginCmd.level2]: {
    mode: ParagraphMode.standard,
    synonyms: 3,
    cmd: ParagraphPluginCmd.level2,
  },
  [ParagraphPluginCmd.level3]: {
    mode: ParagraphMode.shorten,
    synonyms: 1,
    cmd: ParagraphPluginCmd.level3,
  },
  [ParagraphPluginCmd.level4]: {
    mode: ParagraphMode.expand,
    synonyms: 1,
    cmd: ParagraphPluginCmd.level4,
  },
  [ParagraphPluginCmd.level5]: {
    mode: ParagraphMode.simple,
    synonyms: 1,
    cmd: ParagraphPluginCmd.level5,
  },
};

export const ParagraphWordsLimit = 500;

export const ParagraphMode2SelectTextEnum = {
  [ParagraphMode.shorten]: SelectTextEnum.SHORTEN,
  [ParagraphMode.expand]: SelectTextEnum.EXPAND,
  [ParagraphMode.simple]: SelectTextEnum.SIMPLE,
  [ParagraphMode.standard]: SelectTextEnum.STANDARD,
};

export const SelectTextEnum2ParagraphMode = invert(
  ParagraphMode2SelectTextEnum
) as {
  [k in SelectTextEnum]: ParagraphMode;
};
