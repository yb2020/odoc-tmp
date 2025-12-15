import scrollIntoView from 'scroll-into-view-if-needed';
import { connect } from '@/dom/arrow';
import { PDF_ANNOTATE_ID, attrSelector } from '@idea/pdf-annotate-core';
import { ViewerController } from '@idea/pdf-annotate-viewer';

export function getQueryString(name: string) {
  return new URLSearchParams(window.location.search).get(name);
}

export function connectComment(pdfAnnotateId: string) {
  setTimeout(() => {
    const comment = document.querySelector(
      `.note${attrSelector(PDF_ANNOTATE_ID, pdfAnnotateId)}`
    );

    if (!comment) return;

    scrollIntoView(comment, {
      scrollMode: 'if-needed',
      block: 'nearest',
      inline: 'nearest',
    });

    connect(pdfAnnotateId);
  }, 300);
}

export function scrollMark(
  pageNumber: number,
  pdfAnnotateId: string,
  pdfViewer: ViewerController
) {
  const scroller = pdfViewer.getScrollController();

  scroller.goToPage(pageNumber);

  const annotateEle = document.querySelector(
    attrSelector(PDF_ANNOTATE_ID, pdfAnnotateId)
  );

  annotateEle &&
    scrollIntoView(annotateEle, {
      scrollMode: 'if-needed',
      block: 'center',
      inline: 'center',
    });
}
