/*
 * Created Date: May 12th 2021, 3:17:53 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: March 22nd 2022, 5:46:59 pm
 */
export enum ToolBarEndEvent {
  EventEnd = 'eventEnd',
}

export enum ToolBarType {
  select = 0,
  rect = 1,
  ios = 2,
  hot = 3,

  Cursor = 'cursor',
  Selector = 'selector',
  Draw = 'drawing',
  Text = 'textbox',
  // Rect = 'rect',
  Line = 'line',
  Arrow = 'arrow',
  Highlight = 'highlight',
  AIHighlight = 'aihighlight',
  Underline = 'underline',
  Vocabulary = 'vocabulary',
  // Comment = 'select',
  Delete = 'delete',
  None = 'none',
  Circle = 'circle',
}

export const READPAPER_PAGE_CONTAINER_ = 'readpaper-page-container-';
export const DATA_PAGE_NUMBER = 'data-page-number';
export const PDF_ANNOTATE_PAGENUMBER = 'data-pdf-annotate-pageNumber';
export const PDF_ANNOTATE_CONTAINER = 'data-pdf-annotate-container';
export const PDF_ANNOTATE_VIEWPORT = 'data-pdf-annotate-viewport';
export const PDF_ANNOTATE_DOCUMENT = 'data-pdf-annotate-document';
export const PDF_ANNOTATE_ID = 'data-pdf-annotate-id';
export const PDF_ANNOTATE_TYPE = 'data-pdf-annotate-type';
export const PDF_ANNOTATE_HANDWRITE_CANVAS =
  'data-pdf-annotate-handwrite-canvas';

export const PDF_ANNOTATE_SVG_READONLY = 'pdf-annotate-svg-readonly';
export const PDF_ANNOTATE_SHAPE_READONLY = 'pdf-annotate-shape-readonly';
export const PDF_ANNOTATE_SHAPE_DIV = 'data-pdf-annotate-shape-div';
export const PDF_ANNOTATE_TEXTSHAPE_SCALE = 'pdf-annotate-textshape-scale';
export const PDF_ANNOTATE_TEXTSHAPE_BOX = 'data-pdf-annotate-textshape-box';
export const PDF_ANNOTATE_TEXTSHAPE_SPAN = 'data-pdf-annotate-textshape-span';
export const PDF_ANNOTATE_TEXTSHAPE_STYLE = 'pdf-annotate-textshape-style';

export const PDF_ANNOTATE_USERID = 'data-pdf-annotate-userId';
export const PDF_ANNOTATE_USER_SELECT = 'data-pdf-annotate-user-select';
export const PDF_ANNOTATIONLAYER = 'annotationLayer';
export const PDF_ANNOTATIONLAYER_GROUP_AI = 'g-ai';
export const PDF_ANNOTATIONLAYER_GROUP_NOTES = 'g-notes';
export const PDF_ANNOTATIONLAYER_GROUP_VOCABULARY = 'g-vocabulary';
export const PDF_CANVASWRAPPER = 'canvasWrapper';

export const attrSelector = (...args: [key: string, value?: string]) => {
  const [key, value] = args;
  const { length } = args;

  const valueSelector = length === 1 ? '' : `="${value}"`;
  return `[${key}${valueSelector}]`;
};
