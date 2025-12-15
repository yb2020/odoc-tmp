import { defineStore } from 'pinia';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import { PDFJSAnnotate } from '@idea/pdf-annotate-core';
import { toRaw } from 'vue';

export interface PdfState {
  pdfViewerInstanceMap: Map<string, ViewerController>;
  pdfAnnotateInstanceMap: Map<string, PDFJSAnnotate>;
  curPdfViewerInstance: ViewerController | null;
  curPdfAnnotateInstance: PDFJSAnnotate | null;
}

export const usePdfStore = defineStore('pdf', {
  state: (): PdfState => ({
    pdfViewerInstanceMap: new Map(),
    pdfAnnotateInstanceMap: new Map(),
    curPdfAnnotateInstance: null,
    curPdfViewerInstance: null,
  }),
  getters: {
    getViewer(state) {
      return (pdfId: string) => {
        return toRaw(state.pdfViewerInstanceMap.get(pdfId));
      };
    },
    getAnnotater(state) {
      return (noteId: string) => {
        return toRaw(state.pdfAnnotateInstanceMap.get(noteId));
      };
    },
  },
  actions: {
    setViewer(pdfId: string, viewer: ViewerController) {
      this.pdfViewerInstanceMap.set(pdfId, viewer);
      this.curPdfViewerInstance = viewer;
    },
    setAnnotater(noteId: string, annotater: PDFJSAnnotate) {
      this.pdfAnnotateInstanceMap.set(noteId, annotater);
      this.curPdfAnnotateInstance = annotater;
    },
  },
});
