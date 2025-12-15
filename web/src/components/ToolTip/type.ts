import {  TextRectCoordinate } from '@idea/pdf-annotate-viewer';

export type PageTextRects = (TextRectCoordinate & {pageNumber: number})[]
