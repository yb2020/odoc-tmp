import {
  ViewerController,
  ViewerEvent,
  ScrollEvent,
  PageMouseEventPayload,
} from '@idea/pdf-annotate-viewer';
import { computed, watch } from 'vue';
import {
  useFullTextTranslateStore,
  PDFWebviewScrollMode,
  PDFWebviewPreviewMode,
} from '../stores/fullTextTranslateStore';
import { debounce } from 'lodash-es';
import { SVG, Svg } from '@svgdotjs/svg.js';
import { PDFPageView } from '@idea/pdfjs-dist/web/pdf_viewer';
import { switchThemeColorToDefault } from '../theme';

const usePDFWebviewParallelHighlight = ({
  pdfWebviewPointers,
}: {
  pdfWebviewPointers: {
    mainPDFWebview: null | ViewerController;
    tiedPDFWebview: null | ViewerController;
  };
}) => {
  const fullTextTranslateStore = useFullTextTranslateStore();

  const clearSvg = (pdfWebview: ViewerController) => {
    const pdfViewer = pdfWebview.getDocumentViewer().getPdfViewer();
    const pages = pdfViewer.pagesCount;
    for (let i = 0; i < pages; i += 1) {
      const page = pdfViewer.getPageView(i);
      if (page?.extraSvgInstance) {
        page.extraSvgInstance.clear();
      }
    }
  };

  const parallelHightlight = (
    payload: PageMouseEventPayload,
    type: 'main' | 'tied'
  ) => {
    if (!fullTextTranslateStore.alignment) {
      return;
    }
    const mainPDFWebview = pdfWebviewPointers.mainPDFWebview;
    const tiedPDFWebview = pdfWebviewPointers.tiedPDFWebview;

    mainPDFWebview && clearSvg(mainPDFWebview);
    tiedPDFWebview && clearSvg(tiedPDFWebview);

    const { pageIndex } = payload.point;

    if (!mainPDFWebview || !tiedPDFWebview) {
      return;
    }
    const mainPage = mainPDFWebview
      .getDocumentViewer()
      .getPdfViewer()
      .getPageView(pageIndex);
    const tiedPage = tiedPDFWebview
      .getDocumentViewer()
      .getPdfViewer()
      .getPageView(pageIndex);

    if (
      fullTextTranslateStore.previewMode !==
      PDFWebviewPreviewMode.withOriginalPDF
    ) {
      return;
    }

    const alignments =
      type === 'main'
        ? fullTextTranslateStore.alignment.src
        : fullTextTranslateStore.alignment.translated;
    const currentAlignment = alignments[pageIndex];
    if (!currentAlignment || !currentAlignment.bboxes?.length) {
      return;
    }

    const radio = payload.point.viewport.width / currentAlignment.width;
    let i = 0;
    for (; i < currentAlignment.bboxes.length; i += 1) {
      const bbox = currentAlignment.bboxes[i];
      const [x0, y0, x1, y1] = bbox[0];
      if (
        payload.point.left > x0 * radio &&
        payload.point.left < x1 * radio &&
        payload.point.top > y0 * radio &&
        payload.point.top < y1 * radio
      ) {
        break;
      }
    }

    const currentBbox = currentAlignment.bboxes[i];
    if (!currentBbox) {
      return;
    }

    if (mainPage.extraSvgInstance && tiedPage.extraSvgInstance) {
      const [x0, y0, , y1] = currentBbox[0];
      (mainPage.extraSvgInstance as Svg)
        .rect(4, (y1 - y0) * radio)
        .fill('#52C41A8C')
        .move((x0 - 6) * radio, y0 * radio);
      currentBbox[1].forEach((bbox) => {
        const [x0, y0, x1, y1] = bbox;
        (tiedPage.extraSvgInstance as Svg)
          .rect((x1 - x0) * radio, (y1 - y0) * radio)
          .fill('#FADB144C')
          .move(x0 * radio, y0 * radio);
      });
    }
  };

  const parallelMainHightlight = (payload: PageMouseEventPayload) => {
    parallelHightlight(payload, 'main');
  };

  const parallelTiedHightlight = (payload: PageMouseEventPayload) => {
    parallelHightlight(payload, 'tied');
  };

  const createSvgInstance = (payload: { source: PDFPageView }) => {
    if ((payload.source as any).extraSvgInstance) {
      ((payload.source as any).extraSvgInstance as Svg).node.remove();
    }
    const draw = SVG().size('100%', '100%');
    payload.source.div.prepend(draw.node);
    draw.addClass('js-fulltranslate-svg');
    draw.node.style.position = 'absolute';
    draw.node.style.top = '0';
    (payload.source as any).extraSvgInstance = draw;
  };

  const enablePDFParallelHighlight = (
    pdfWebview: ViewerController,
    type: 'main' | 'tied'
  ) => {
    pdfWebview.addEventListener(
      ViewerEvent.PAGE_MOUSE_EVENT,
      type === 'main' ? parallelMainHightlight : parallelTiedHightlight
    );
    pdfWebview
      .getDocumentViewer()
      .container.querySelectorAll('div.page[data-loaded="true"]')
      ?.forEach((div) => {
        const pageNumber =
          parseInt(div.getAttribute('data-page-number') || '', 10) || 0;
        const pageView = pdfWebview
          .getDocumentViewer()
          .getPdfViewer()
          .getPageView(pageNumber - 1);
        if (pageView) {
          createSvgInstance({ source: pageView });
        }
      });
    pdfWebview.addEventListener(ViewerEvent.PAGE_RENDERED, createSvgInstance);
  };

  const enableMainPDFParallelHighlight = () => {
    if (pdfWebviewPointers.mainPDFWebview) {
      enablePDFParallelHighlight(pdfWebviewPointers.mainPDFWebview, 'main');
    }
  };

  const clearMainPDFParallelHighlight = () => {
    const mainPDFWebview = pdfWebviewPointers.mainPDFWebview;
    if (mainPDFWebview) {
      mainPDFWebview.removeEventListener(
        ViewerEvent.PAGE_MOUSE_EVENT,
        parallelMainHightlight
      );
      mainPDFWebview.removeEventListener(
        ViewerEvent.PAGE_RENDERED,
        createSvgInstance
      );
      clearSvg(mainPDFWebview);
    }
  };

  const enableTiedPDFParallelHighlight = () => {
    if (pdfWebviewPointers.tiedPDFWebview) {
      enablePDFParallelHighlight(pdfWebviewPointers.tiedPDFWebview, 'tied');
    }
  };

  const clearTiedPDFParallelHighlight = () => {
    pdfWebviewPointers.tiedPDFWebview?.removeEventListener(
      ViewerEvent.PAGE_MOUSE_EVENT,
      parallelTiedHightlight
    );
    pdfWebviewPointers.tiedPDFWebview?.removeEventListener(
      ViewerEvent.PAGE_RENDERED,
      createSvgInstance
    );
  };

  return {
    enableMainPDFParallelHighlight,
    clearMainPDFParallelHighlight,
    enableTiedPDFParallelHighlight,
    clearTiedPDFParallelHighlight,
  };
};

export const usePDFWEbviewPreviewMode = () => {
  const fullTextTranslateStore = useFullTextTranslateStore();
  const showOriginalPDF = computed(() => {
    return (
      !fullTextTranslateStore.pdfId ||
      fullTextTranslateStore.previewMode !==
        PDFWebviewPreviewMode.onlyTranslatePDF
    );
  });
  const enableFullTextTranslate = computed(
    () => !!fullTextTranslateStore.pdfId
  );
  return {
    showOriginalPDF,
    enableFullTextTranslate,
  };
};

export const usePDFWebviewScrollMode = () => {
  const fullTextTranslateStore = useFullTextTranslateStore();
  const pdfWebviewPointers: {
    mainPDFWebview: null | ViewerController;
    tiedPDFWebview: null | ViewerController;
  } = {
    mainPDFWebview: null,
    tiedPDFWebview: null,
  };

  let whereIsMouse = '';

  const {
    enableMainPDFParallelHighlight,
    clearMainPDFParallelHighlight,
    enableTiedPDFParallelHighlight,
    clearTiedPDFParallelHighlight,
  } = usePDFWebviewParallelHighlight({ pdfWebviewPointers });

  const mainScroll = (e: any) => {
    if (
      fullTextTranslateStore.scrollMode !== PDFWebviewScrollMode.lock ||
      fullTextTranslateStore.previewMode ===
        PDFWebviewPreviewMode.onlyTranslatePDF
    ) {
      return;
    }
    if (whereIsMouse === 'tied') {
      return;
    }
    const tiedContainer =
      pdfWebviewPointers.tiedPDFWebview?.getDocumentViewer().container;
    if (tiedContainer) {
      tiedContainer.scrollTop = e.target.scrollTop;
      tiedContainer.scrollLeft = e.target.scrollLeft;
    }
  };

  const mainMousemove = () => {
    whereIsMouse = 'main';
  };

  const tiedScroll = (e: any) => {
    if (
      fullTextTranslateStore.scrollMode !== PDFWebviewScrollMode.lock ||
      fullTextTranslateStore.previewMode ===
        PDFWebviewPreviewMode.onlyTranslatePDF
    ) {
      return;
    }
    if (whereIsMouse === 'main') {
      return;
    }
    const mainContainer =
      pdfWebviewPointers.mainPDFWebview?.getDocumentViewer().container;
    if (mainContainer) {
      mainContainer.scrollLeft = e.target.scrollLeft;
      mainContainer.scrollTop = e.target.scrollTop;
    }
  };

  const tiedMousemove = () => {
    whereIsMouse = 'tied';
  };

  const mainScaleChange = debounce(
    (scale: string | number) => {
      const pdfViewer = pdfWebviewPointers.tiedPDFWebview
        ?.getDocumentViewer()
        .getPdfViewer();
      if (pdfViewer) {
        if (scale === 'page-width') {
          scale = pdfWebviewPointers
            .mainPDFWebview!.getDocumentViewer()
            .getPdfViewer().currentScale;
        }
        pdfViewer.currentScaleValue = `${scale}`;
      }
    },
    100,
    { leading: true }
  );

  const tiedScaleChange = debounce(
    (scale: string | number) => {
      if (
        fullTextTranslateStore.previewMode ===
        PDFWebviewPreviewMode.onlyTranslatePDF
      ) {
        return;
      }
      const pdfViewer = pdfWebviewPointers.mainPDFWebview
        ?.getDocumentViewer()
        .getPdfViewer();
      if (pdfViewer) {
        pdfViewer.currentScaleValue = `${scale}`;
      }
    },
    100,
    { leading: true }
  );

  const initTiedScrollPos = () => {
    if (
      fullTextTranslateStore.previewMode ===
      PDFWebviewPreviewMode.onlyTranslatePDF
    ) {
      return;
    }
    const mainContainer =
      pdfWebviewPointers.mainPDFWebview?.getDocumentViewer().container;
    const tiedContainer =
      pdfWebviewPointers.tiedPDFWebview?.getDocumentViewer().container;

    if (!mainContainer || !tiedContainer) {
      return;
    }
    setTimeout(() => {
      tiedContainer.scrollTop = mainContainer.scrollTop;
      tiedContainer.scrollLeft = mainContainer.scrollLeft;
    }, 100);
  };

  const initMainScrollPos = () => {
    const mainContainer =
      pdfWebviewPointers.mainPDFWebview?.getDocumentViewer().container;
    const tiedContainer =
      pdfWebviewPointers.tiedPDFWebview?.getDocumentViewer().container;

    if (!mainContainer || !tiedContainer) {
      return;
    }
    setTimeout(() => {
      mainContainer.scrollTop = tiedContainer.scrollTop;
      mainContainer.scrollLeft = tiedContainer.scrollLeft;
    }, 100);
  };

  const setPDFWebviewInstance = (
    type: 'main' | 'tied',
    pdfWebview: ViewerController
  ) => {
    if (type === 'main') {
      pdfWebviewPointers.mainPDFWebview = pdfWebview;
      if (pdfWebviewPointers.tiedPDFWebview) {
        const mainContainer =
          pdfWebviewPointers.mainPDFWebview.getDocumentViewer().container;
        mainContainer.addEventListener('scroll', mainScroll, {
          passive: true,
        });

        pdfWebviewPointers.mainPDFWebview.addEventListener(
          ViewerEvent.TRIGGER_SCALE_CHANGE,
          mainScaleChange
        );
        pdfWebviewPointers.mainPDFWebview
          .getScrollController()
          .addEventListener(
            ScrollEvent.TRIGGER_SCROLL_CONTROLLER,
            initTiedScrollPos
          );
      }
      if (fullTextTranslateStore.pdfId) {
        enableMainPDFParallelHighlight();
      }
    } else {
      pdfWebviewPointers.tiedPDFWebview = pdfWebview;
      const tiedContainer = pdfWebview.getDocumentViewer().container;
      tiedContainer.addEventListener('scroll', tiedScroll, {
        passive: true,
      });
      tiedContainer.addEventListener('mousemove', tiedMousemove, {
        passive: true,
      });
      pdfWebview.addEventListener(
        ViewerEvent.TRIGGER_SCALE_CHANGE,
        tiedScaleChange
      );
      if (fullTextTranslateStore.scrollMode === PDFWebviewScrollMode.lock) {
        pdfWebview.addEventListener(ViewerEvent.PAGES_INIT, initTiedScrollPos);
      }
      if (pdfWebviewPointers.mainPDFWebview) {
        pdfWebview.addEventListener(ViewerEvent.PAGES_INIT, () => {
          mainScaleChange(
            pdfWebviewPointers
              .mainPDFWebview!.getDocumentViewer()
              .getPdfViewer().currentScale
          );
        });
        pdfWebviewPointers.mainPDFWebview
          .getDocumentViewer()
          .container.addEventListener('mousemove', mainMousemove, {
            passive: true,
          });
      }
      enableTiedPDFParallelHighlight();
      switchThemeColorToDefault();
    }
  };

  const clearTiedPDFWebview = () => {
    pdfWebviewPointers.tiedPDFWebview = null;
    const mainContainer =
      pdfWebviewPointers.mainPDFWebview?.getDocumentViewer().container;
    mainContainer?.removeEventListener('scroll', mainScroll);
    document.removeEventListener('mousemove', mainMousemove);

    pdfWebviewPointers.mainPDFWebview?.removeEventListener(
      ViewerEvent.SCALE_CHANGING,
      mainScaleChange
    );

    pdfWebviewPointers.mainPDFWebview
      ?.getScrollController()
      .removeEventListener(
        ScrollEvent.TRIGGER_SCROLL_CONTROLLER,
        initTiedScrollPos
      );

    clearMainPDFParallelHighlight();
    clearTiedPDFParallelHighlight();
  };

  watch(
    () => fullTextTranslateStore.pdfId,
    (newVal, oldVal) => {
      if (!newVal && oldVal) {
        clearTiedPDFWebview();
      }
      if (newVal && pdfWebviewPointers.mainPDFWebview) {
        const mainContainer =
          pdfWebviewPointers.mainPDFWebview.getDocumentViewer().container;
        mainContainer.addEventListener('scroll', mainScroll, {
          passive: true,
        });
        pdfWebviewPointers.mainPDFWebview.addEventListener(
          ViewerEvent.TRIGGER_SCALE_CHANGE,
          mainScaleChange
        );
        pdfWebviewPointers.mainPDFWebview
          .getScrollController()
          .addEventListener(
            ScrollEvent.TRIGGER_SCROLL_CONTROLLER,
            initTiedScrollPos
          );
        enableMainPDFParallelHighlight();
      }
    }
  );

  watch(
    () => fullTextTranslateStore.previewMode,
    (newVal) => {
      if (
        newVal === PDFWebviewPreviewMode.withOriginalPDF &&
        fullTextTranslateStore.pdfId &&
        fullTextTranslateStore.scrollMode === PDFWebviewScrollMode.lock
      ) {
        initMainScrollPos();
      }
    }
  );

  return {
    setPDFWebviewInstance,
    clearTiedPDFWebview,
  };
};
