import { ViewerController } from '@idea/pdf-annotate-viewer';

export const getPDFCurPages = (ctrl?: ViewerController) => {
  const pageEls =
    ctrl
      ?.getDocumentViewer()
      .container.querySelectorAll<HTMLDivElement>(
        '.page[data-loaded="true"]'
      ) ?? [];
  return [...pageEls]
    .map((el) => Number(el.getAttribute('data-page-number')))
    .filter(Boolean);
};

export const goToPDFPage = (
  page: number,
  noLastPosition = false,
  pdfWebviewController?: ViewerController
) => {
  pdfWebviewController?.getScrollController().goToPage(page, noLastPosition);
};

export const changePDFScale = (
  scale: 'page-width' | number,
  pdfWebviewController?: ViewerController
) => {
  if (!pdfWebviewController) {
    return;
  }
  if (scale === 'page-width') {
    pdfWebviewController.getDocumentViewer().changeScale(scale);
    return;
  }
  let newScale = Math.max(scale, MIN_SCALE);
  newScale = Math.min(newScale, MAX_SCALE);
  pdfWebviewController.getDocumentViewer().changeScale(newScale);
};

// const DEFAULT_SCALE_DELTA = 1.1;
const MIN_SCALE = 0.5;
const MAX_SCALE = 4.0;
export const changePDFScaleBySteps = (
  type: 'increase' | 'decrease',
  pdfWebviewController: ViewerController
) => {
  const pdfDocumentViewer = pdfWebviewController.getDocumentViewer();
  if (!pdfDocumentViewer) {
    return;
  }
  if (type === 'increase') {
    pdfDocumentViewer.zoomIn(1);
  } else {
    pdfDocumentViewer.zoomOut(1);
  }
};

export const scrollToPDFLastPosition = (
  pdfWebviewController?: ViewerController
) => {
  pdfWebviewController?.getScrollController().scrollToPdfViewerLastPosition();
};

export const togglePDFSearchViewer = (
  pdfWebviewController?: ViewerController
) => {
  return pdfWebviewController?.getFinderViewer()?.toggleSearchPanel();
};
