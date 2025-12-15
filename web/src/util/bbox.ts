import { PdfBBox } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import { Nullable } from '../typings/global';
import scrollIntoViewIfNeeded from 'scroll-into-view-if-needed';

export const goBbox = (
  pageNum: number,
  bbox: PdfBBox,
  pdfViewer: ViewerController
) => {
  if (!pdfViewer) {
    throw new Error('Invalid pdfViewer');
  }

  const scrollController = pdfViewer.getScrollController();

  scrollController.updatePdfViewerLastPosition();

  scrollController.goToPage(pageNum, true);

  const pageEl = pdfViewer
    .getDocumentViewer()
    .container.querySelector(
      `.page[data-page-number="${pageNum}"]`
    ) as Nullable<HTMLDivElement>;
  if (pageEl) {
    const div = document.createElement('div');
    const ratio = pageEl.offsetWidth / bbox.originWidth;
    div.style.width = (bbox.x1 - bbox.x0) * ratio + 'px';
    div.style.height = (bbox.y1 - bbox.y0) * ratio + 'px';
    div.style.backgroundColor = 'rgb(255, 249, 86, 50%)';
    div.style.position = 'absolute';
    div.style.left = bbox.x0 * ratio - 8 + 'px';
    div.style.top = bbox.y0 * ratio - 8 + 'px';
    div.style.padding = '8px';
    div.classList.add('bbox-animation');
    div.style.boxSizing = 'content-box';
    pageEl.appendChild(div);
    window.setTimeout(() => {
      pageEl.removeChild(div);
    }, 3000);

    const position =
      div.offsetHeight >= window.innerHeight ? 'start' : 'center';

    scrollIntoViewIfNeeded(div, {
      block: position,
      inline: position,
    });
  }
};
