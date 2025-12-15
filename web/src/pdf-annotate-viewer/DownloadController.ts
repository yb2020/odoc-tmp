import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer';
import * as pdfjsLib from '@idea/pdfjs-dist';
import { PDFDocumentProxy } from "@idea/pdfjs-dist/types/src/display/api";

export default class DownloadController {
  private pdfDocument?: PDFDocumentProxy;

  private downloadManager: pdfjsViewer.DownloadManager;
  private pdfScriptingManager: pdfjsViewer.PDFScriptingManager;

  private downloadComplete = false

  private pdfUrl: string;

  // private saveInProgress = false;

  constructor({ url, pdfScriptingManager }: { url: string, pdfScriptingManager: pdfjsViewer.PDFScriptingManager }) {
    this.downloadManager = new pdfjsViewer.DownloadManager()
    this.pdfUrl = url
    this.pdfScriptingManager = pdfScriptingManager
  }

  setPDFDocument(pdfDocument: PDFDocumentProxy) {
    this.pdfDocument = pdfDocument;
  }

  private async ensureDownloadComplete() {
    if (!this.pdfDocument) {
      throw Error('Call setPDFDocument first.')
    }
    if (this.downloadComplete) {
      return true
    }
    await this.pdfDocument.getDownloadInfo()
    this.downloadComplete = true
    return true
  }

  async download(filename?: string) {
    filename = filename || pdfjsLib.getFilenameFromUrl(this.pdfUrl) || 'export.pdf';
    try {
      await this.ensureDownloadComplete()
      const data = await this.pdfDocument!.getData()
      const blob = new Blob([data], { type: "application/pdf" });
      await this.downloadManager.download(blob, '', filename);
    } catch (error) {
      // cdn增加响应头 https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Disposition
      const url = new URL(this.pdfUrl)
      url.searchParams.set('attname', '')
      await this.downloadManager.downloadUrl(url.href, filename)
    }
  }

  // private async save(filename?: string) {
  //   const pdfDoc = this.pdfDocument!;
  //   pdfDoc.annotationStorage.setValue('1046R', { value: 'Hello World', })
  //   if (this.saveInProgress) {
  //     return
  //   }
  //   filename = filename || getFilenameFromUrl(this.pdfUrl) || 'exportWithAnnotations.pdf';
  //   await this.ensureDownloadComplete()
  //   this.saveInProgress = true;
  //   await this.pdfScriptingManager.dispatchWillSave(null);
  //   try {
  //     const data = await this.pdfDocument!.saveDocument();
  //     const blob = new Blob([data], { type: "application/pdf", });

  //     await this.downloadManager.download(blob, '', filename);
  //   } catch (error) {
  //     await this.download(filename);
  //   } finally {
  //     await this.pdfScriptingManager.dispatchDidSave(null);
  //     this.saveInProgress = false;
  //   }
   
  // }
}