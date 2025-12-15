import PDFPageView from '@idea/pdf-annotate-viewer/typing/PDFPageView';
import {
  PDFJSAnnotate,
  PageHandwriteMap,
  PageShapeMap,
  PageShapeTextMap,
  ShapeCallback,
  TextCallback,
} from '..';

/*
 * Created Date: March 8th 2022, 11:08:19 am
 * Author: zhoupengcheng
 * -----
 * Last Modified: March 9th 2022, 10:33:32 am
 */
export interface RenderAnnotateProps {
  documentId: string;
  pageNumber: number;
  source: PDFPageView;
  instance: PDFJSAnnotate;
  handwriteAndShapeVisible: boolean;
  handwriteBuffer?(): Promise<PageHandwriteMap> | null;
  shapeBuffer?(): Promise<PageShapeMap> | null;
  shapeTextBuffer?(): Promise<PageShapeTextMap | null>;
  shapeEditable?: boolean;
  shapeCallback?: ShapeCallback;
  textCallback?: TextCallback;
}

export type RenderCanvasProps = RenderAnnotateProps & {
  pageView: PDFPageView;
};

export type PrepareCanvasProps = Pick<
  RenderCanvasProps,
  'instance' | 'pageNumber' | 'pageView' | 'shapeEditable'
>;

export enum ToolBarType {
  Cursor = 'cursor',
  Selector = 'selector',
  Draw = 'drawing',
  Text = 'textbox',
  Rect = 'rect',
  Line = 'line',
  Arrow = 'arrow',
  Highlight = 'highlight',
  Underline = 'underline',
  Comment = 'comment',
  Delete = 'delete',
  None = 'none',
  Circle = 'circle',
}

// 通用的标记类型
export interface CommonAnnotation {
  uuid: string;
  documentId: string;
}

// // 圆形标记
// interface CircleOptions {
//   x: number;
//   y: number;
//   width: number;
//   height: number;
// }
// export type CircleAnnotation = CommonAnnotation & {
//   color: string;
//   width: number;
//   opacity: number;
//   rectangles: CircleOptions[];
//   type: ToolBarType.Circle;
// };

// // ✏️ 画画
// export type PathAnnotation = CommonAnnotation & {
//   lines: Array<Array<number>>;
//   color: string;
//   width?: number;
//   strokeOpacity?: number;
//   type: ToolBarType.Draw;
// };

// 矩形
interface RectOptions {
  x: number;
  y: number;
  width: number;
  height: number;
}
export type RectAnnotation = CommonAnnotation & {
  color: string;
  strokeWidth: number;
  opacity: number;
  rectangles: RectOptions[];
  type: ToolBarType.Rect;
  fill: string;
  idea: string;
  picUrl: string;
};

// 选区 => 高亮，评论，下划线
interface Rectangle {
  height: number;
  width: number;
  x: number;
  y: number;
}
export type SelectAnnotation = CommonAnnotation & {
  type: ToolBarType.Underline | ToolBarType.Comment | ToolBarType.Highlight;
  rectangles: Array<Rectangle>;
  rectStr: string;
  idea: string;
};

// // 文字标记
// export type TextAnnotation = CommonAnnotation & {
//   type: ToolBarType.Text;
//   content: string;
//   size: string;
//   x: number;
//   y: number;
//   color: string;
// };

// // 直线或者箭头标记
// export type LineAnnotation = CommonAnnotation & {
//   color: string;
//   width: number;
//   opacity: number;
//   rectangles: { x1: number; y1: number; x2: number; y2: number }[];
//   type: ToolBarType.Arrow | ToolBarType.Line;
// };

export type CommonProps<A, B> = Pick<
  A,
  {
    [K in keyof A & keyof B]: A[K] extends B[K]
      ? B[K] extends A[K]
        ? K
        : never
      : never;
  }[keyof A & keyof B]
>;
