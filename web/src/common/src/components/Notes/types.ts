// import {
//   // NoteManageDocInfo,
//   // WordColor as VocabularyColorKey,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import { NoteManageDocInfo } from 'go-sea-proto/gen/ts/note/NoteManage'
import { WordColor as VocabularyColorKey, } from 'go-sea-proto/gen/ts/note/NoteWord'
import { GetMarkTagListByFolderIdResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/doc/response/_GetMarkTagListByFolderIdResponse';
import FileText from '~common/assets/images/notes/file-text.png';
import FileTextSelected from '~common/assets/images/notes/file-text-selected.png';
import HintTrans from '~common/assets/images/notes/hint-translation.png';
import HintTransSelected from '~common/assets/images/notes/hint-translation-selected.png';
import FeedEdit from '~common/assets/images/notes/feedback-edit.png';
import FeedEditSelected from '~common/assets/images/notes/feedback-edit-selected.png';

export enum NoteSubTypes {
  Summary,
  Vocabulary,
  Annotation,
}

export const NoteSubType2ModuleType = {
  [NoteSubTypes.Summary]: 'note_summary',
  [NoteSubTypes.Vocabulary]: 'note_word',
  [NoteSubTypes.Annotation]: 'note_excerpt',
};

export const NoteSubType2Icons = {
  [NoteSubTypes.Summary]: [FileText, FileTextSelected],
  [NoteSubTypes.Vocabulary]: [HintTrans, HintTransSelected],
  [NoteSubTypes.Annotation]: [FeedEdit, FeedEditSelected],
};

export const NoteSubType2I18nKey = {
  [NoteSubTypes.Summary]: 'summary',
  [NoteSubTypes.Vocabulary]: 'vocabulary',
  [NoteSubTypes.Annotation]: 'annotation',
};

export enum ColorKey {
  blue = 1,
  green = 2,
  yellow = 3,
  orange = 4,
  red = 5,
}

export interface ColorStyle {
  fill: string;
  fillOpacity: number;
  color: string;
  text: string;
  i18n: string;
}

export const styleMap: Record<ColorKey, ColorStyle> = {
  [ColorKey.blue]: {
    fill: '#338AFF',
    fillOpacity: 0.3,
    color: '#338AFF',
    text: '蓝色',
    i18n: 'common.notes.blue',
  },
  [ColorKey.green]: {
    fill: '#80E639',
    fillOpacity: 0.3,
    color: '#80E639',
    text: '绿色',
    i18n: 'common.notes.green',
  },
  [ColorKey.yellow]: {
    fill: '#FFFF00',
    fillOpacity: 0.3,
    color: '#FFFF00',
    text: '黄色',
    i18n: 'common.notes.yellow',
  },
  [ColorKey.orange]: {
    fill: '#FF8C19',
    fillOpacity: 0.3,
    color: '#FF8C19',
    text: '橙色',
    i18n: 'common.notes.orange',
  },
  [ColorKey.red]: {
    fill: '#F24030',
    fillOpacity: 0.3,
    color: '#F24030',
    text: '红色',
    i18n: 'common.notes.red',
  },
};

export { VocabularyColorKey };

export const vocabularyStyleMap: Record<string, Omit<ColorStyle, 'color'>> = {
  [VocabularyColorKey.BLUE]: styleMap[ColorKey.blue],
  // {
  //   fill: '#3A57E6',
  //   fillOpacity: 0.3,
  //   text: '极客蓝',
  //   i18n: 'viewer.geekblue',
  // },
  [VocabularyColorKey.CYAN]: styleMap[ColorKey.green],
  // {
  //   fill: '#0DC1C1',
  //   fillOpacity: 0.3,
  //   text: '青色',
  //   i18n: 'viewer.cyan',
  // },
  [VocabularyColorKey.PURPLE]: styleMap[ColorKey.orange],
  // {
  //   fill: '#813DE1',
  //   fillOpacity: 0.3,
  //   text: '紫色',
  //   i18n: 'viewer.purple',
  // },
  [VocabularyColorKey.AURATUS]: styleMap[ColorKey.yellow],
  // {
  //   fill: '#FAAD14',
  //   fillOpacity: 0.3,
  //   text: '金黄色',
  //   i18n: 'viewer.orange',
  // },
  [VocabularyColorKey.ROSE_RED]: styleMap[ColorKey.red],
  // {
  //   fill: '#E4379C',
  //   fillOpacity: 0.3,
  //   text: '玫瑰红色',
  //   i18n: 'viewer.red',
  // },
};

export interface NoteStyleDefine extends ColorStyle {
  type: ColorKey;
}

export interface NoteBreadcrumb {
  key: string;
  title: string;
}

export interface NoteFolder extends NoteBreadcrumb {
  count?: number;
  paperId?: string;
  pdfId?: string;
  docId?: string;
  children: NoteFolder[];
  isDoc?: boolean;
  docInfos?: NoteManageDocInfo[];
  noteWordCount?: number;
  noteAnnotateCount?: number;
}

export interface NoteFolderExtra {
  path: number[];
  width: number;
  title: string;
  isDoc: boolean;
  docId?: string;
  docInfos?: NoteManageDocInfo[];
  noteWordCount?: number;
  noteAnnotateCount?: number;
}

export type NoteTag = GetMarkTagListByFolderIdResponse['markTagList'][0];
