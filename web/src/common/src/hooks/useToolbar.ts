import { ref } from 'vue';
import {
  ScaleType,
  usePDFViewerController,
} from '@common/hooks/useCreatePDFViewer';
import { ViewerEvent } from '@idea/pdf-annotate-viewer';

export const useToolbar = (
  getPdfViewerController: ReturnType<
    typeof usePDFViewerController
  >['getPdfViewerController']
) => {
  const total = ref(0);
  const pageNum = ref(1);
  const scale = ref<ScaleType>(1.0);

  const addEvent = async () => {
    const pdfWebviewController = await getPdfViewerController();
    total.value = pdfWebviewController.getPdfDocument()?.numPages || 0;
    scale.value = pdfWebviewController.getDocumentViewer().getPdfViewer()
      .currentScaleValue as ScaleType;
    pdfWebviewController.addEventListener(ViewerEvent.PAGES_INIT, (event) => {
      total.value = event.source.pagesCount;
    });
    pdfWebviewController.addEventListener(
      ViewerEvent.PAGE_CHANGING,
      (payload) => {
        pageNum.value = payload.pageNumber;
      }
    );
    pdfWebviewController.addEventListener(
      ViewerEvent.SCALE_CHANGING,
      (payload) => {
        scale.value = (payload.presetValue as ScaleType) || payload.scale;
      }
    );
  };

  const onPageChange = async (page: number) => {
    console.log('onPageChange', page);
    const pdfWebviewController = await getPdfViewerController();
    pdfWebviewController?.getScrollController().goToPage(page);
  };

  const onScaleChange = async (scale: ScaleType | 'increase' | 'decrease') => {
    const pdfWebviewController = await getPdfViewerController();
    if (scale === 'increase') {
      pdfWebviewController?.getDocumentViewer().zoomIn(1);
    } else if (scale === 'decrease') {
      pdfWebviewController?.getDocumentViewer().zoomOut(1);
    } else {
      pdfWebviewController?.getDocumentViewer().changeScale(scale);
    }
  };

  addEvent();

  return {
    onPageChange,
    onScaleChange,
    pageNum,
    scale,
    total,
  };
};
