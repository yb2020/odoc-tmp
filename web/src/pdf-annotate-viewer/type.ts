import * as pdfjsLib from "@idea/pdfjs-dist";
import { TextItem } from '@idea/pdfjs-dist/types/src/display/api';

export * from './ViewerController';
import type _ThumbnailViewer from './ThumbnailViewer'
import { Language } from 'go-sea-proto/gen/ts/lang/Language';

export type ThumbnailViewer = _ThumbnailViewer

export interface BaseMousePoint {
  x: number;
  y: number;
}
export interface MousePoint extends BaseMousePoint {
  index: number;
}

export interface OffsetCoordinate {
  left: number;
  top: number;
  right: number;
  bottom: number;
  size: [number, number];
  m?: number[];
  angle: number;
}

export interface TextContentChar {
  char: string;
  dim: number;
  transfrom: number[];
  unicode: string;
  str: string;
  spaceWidth: number;
}

export interface TextContentBound {
  geom: TextItem & { charsArray: TextContentChar[] };
  str: string;
  offset: {
    angle: number;
    bottom: number;
    left: number;
    m?: number[];
    right: number;
    size: [number, number];
    top: number;
  };
  style: {
    ascent: number;
    descent: number;
    fontFamily: string;
    vertical: boolean;
    fontSize: number;
  };
  shouldScaleText: boolean;
}

export interface PointOnPdf {
  left: number;
  top: number;
  cur: number;
  pv: pdfjsLib.PageViewport;
}

export enum MousePointDir {
  forward,
  backward,
}

export interface TextRectCoordinate {
  x: number;
  y: number;
  width: number;
  height: number;
  offset?: {
    angle: number;
    dx: number;
    dy: number;
  };
  text: string;
  shouldScaleText: boolean;
}

export interface PageSelectText {
  multiSegment: boolean;
  crossing: boolean;
  pageNum: number;
  text: string;
  rects: TextRectCoordinate[];
  viewport: pdfjsLib.PageViewport;
}

export enum ViewerEvent {
  PROGRESS_CHANGE = 'x-progress-change',
  TEXT_SELECT = 'x-text-select',
  EMPTY_CLICK = 'emptyclick',
  SCALE_CHANGING = 'scalechanging',
  PAGE_RENDERED = 'pagerendered',
  TEXT_LAYER_RENDERED = 'textlayerrendered',
  PAGES_INIT = 'pagesinit',
  PAGE_CHANGING = 'pagechanging',
  TRIGGER_SCALE_CHANGE = 'x-trigger-scale-change',
  PAGE_MOUSE_EVENT = 'x-mouse-event',
  PDFURL_RETRY = 'x-pdfurl-retry'
}

export type PageMouseEventPayload = {
  event: MouseEvent;
  point: {
      left: number;
      top: number;
      pageIndex: number;
      viewport: pdfjsLib.PageViewport;
  };
}


export type PDFViewerScale = string | number;
