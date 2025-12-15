/*
 * Created Date: March 22nd 2022, 3:14:01 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: March 22nd 2022, 5:46:59 pm
 */
import { createGlobalState, useLocalStorage } from '@vueuse/core';
import { ColorKey } from '../style/select';
import { AnnotationColor } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { PDF_READER } from '@/common/src/constants/storage-keys';

const COMMENT_KEY = PDF_READER.COMMENT;

export const DEFAULT_FONT_SIZE = 12;
export const fontSizeList = [DEFAULT_FONT_SIZE, 16, 22];

export const useCommentGlobalState = createGlobalState(() =>
  useLocalStorage(COMMENT_KEY, {
    styleId: ColorKey.blue,
    shapeStyleId: AnnotationColor.blue,
    shapeFontSize: DEFAULT_FONT_SIZE,
  })
);
