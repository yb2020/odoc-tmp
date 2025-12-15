import { defineStore } from 'pinia';
import { ElementClick, PageType, reportElementClick } from '@/api/report';

export enum PDFWebviewScrollMode {
  unlock,
  lock,
}

export enum PDFWebviewPreviewMode {
  onlyTranslatePDF,
  withOriginalPDF,
}

export type Bbox = [
  [number, number, number, number],
  [number, number, number, number][]
];

export interface FullTextTranslateState {
  pdfId: string;
  fullTextTranslatePDFUrl: string;
  scrollMode: PDFWebviewScrollMode;
  previewMode: PDFWebviewPreviewMode;
  alignment?: {
    src: Record<string, { bboxes: Bbox[]; width: number; height: number }>;
    translated: Record<
      string,
      { bboxes: Bbox[]; width: number; height: number }
    >;
  };
}

export const useFullTextTranslateStore = defineStore('fullTextTranslate', {
  state: (): FullTextTranslateState => ({
    pdfId: '',
    fullTextTranslatePDFUrl: '',
    scrollMode: PDFWebviewScrollMode.lock,
    previewMode: PDFWebviewPreviewMode.withOriginalPDF,
  }),
  actions: {
    setFullTextTranslatePDFUrl(url: string) {
      this.fullTextTranslatePDFUrl = url;
    },
    setFullTextTranslateAlignment(str: string) {
      try {
        const json = JSON.parse(str);
        this.alignment = json;
      } catch (error) {}
    },
    setScrollMode(mode: PDFWebviewScrollMode) {
      this.scrollMode = mode;
    },
    setPreviewMode(mode: PDFWebviewPreviewMode) {
      this.previewMode = mode;
    },
    toggleFullTextTranslate(pdfId?: string) {
      if (pdfId && !this.pdfId) {
        this.pdfId = pdfId;
      } else {
        this.pdfId = '';
      }
      // 全文翻译按钮点击上报
      reportElementClick({
        page_type: PageType.note,
        type_parameter: pdfId || '',
        element_name: ElementClick.full_translationg_reading,
        status: this.pdfId ? 'on' : 'off',
      });
    },
  },
});
