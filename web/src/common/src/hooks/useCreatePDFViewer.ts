import { onMounted, onUnmounted, ref } from 'vue';
import { ViewerController, createPDFWebview } from '@idea/pdf-annotate-viewer';
import pdfjsWorkerSrc from '@idea/pdfjs-dist/build/pdf.worker.js?url';
import debounce from 'lodash-es/debounce';

export type ScaleType = number | 'page-width';
export const getScalePresetValue = (value: ScaleType) => {
  if (typeof value === 'number') {
    // 百分比整数，不要小数点
    return `${(value * 100).toFixed(0)}%`;
  }
  return '适应大小';
};

const resizeScale = (
  container: HTMLDivElement,
  pdfWebviewController: ViewerController
) => {
  if (!container) {
    console.warn('container is not ready');
    return null;
  }
  try {
    /**
     * 如果是自适应的设置，需要动态调整大小
     */
    const resizeObserver = new ResizeObserver(
      debounce(
        () => {
          const documentViewer = pdfWebviewController.getDocumentViewer();
          if (
            documentViewer.getPdfViewer().currentScaleValue === 'page-width'
          ) {
            documentViewer.changeScale('page-width');
          }
        },
        100,
        {
          leading: false,
          trailing: true,
        }
      )
    );

    resizeObserver.observe(container);

    return () => {
      resizeObserver.disconnect();
    };
  } catch (error) {
    console.error(error);
    return null;
  }
};

export interface PDFViewerOptions {
  scale?: ScaleType;
  onInitSuccess?: (pdfViewerController: ViewerController) => void;
  onCreateSuccess?: (pdfViewerController: ViewerController) => void;
  onCreateError?: (error: Error) => void;
}

export const useCreatePDFViewer = (
  pdfUrl: string,
  options: PDFViewerOptions
) => {
  const documentWrapper = ref<HTMLDivElement>();

  let pdfWebviewController: ViewerController | null = null;
  let cancelResizeObserver: (() => void) | null = null;

  const errorInfo = ref<Error | null>(null);

  const createPDFViewer = async (pdfUrl: string) => {
    if (!documentWrapper.value) {
      throw new Error('documentWrapper is not ready');
    }
    const pdfWebview = createPDFWebview(
      {
        pdfDocumentParams: {
          url: pdfUrl,
          cMapUrl:
            'https://nuxt.cdn.readpaper.com/pdfjs-dist%402.13.216/cmaps/',
        },
        containers: {
          documentWrapper: documentWrapper.value,
        },
      },
      pdfjsWorkerSrc
    );

    pdfWebviewController = pdfWebview;

    options.onInitSuccess?.(pdfWebview);

    try {
      await pdfWebview.build(options.scale || '1.0', {
        annotationMode: 1,
      });

      options.onCreateSuccess?.(pdfWebview);

      // 自适应大小
      cancelResizeObserver = resizeScale(documentWrapper.value, pdfWebview);
    } catch (error) {
      console.error(error);
      errorInfo.value = error as Error;
      options.onCreateError?.(error as Error);
      destroy();
    }
  };

  onMounted(async () => {
    await createPDFViewer(pdfUrl);
  });

  onUnmounted(() => {
    destroy();
  });

  const destroy = () => {
    cancelResizeObserver?.();
    pdfWebviewController?.destroy();
    pdfWebviewController = null;
    if (documentWrapper.value) {
      documentWrapper.value.innerHTML = '';
    }
  };

  const open = async (pdfUrl: string) => {
    destroy();
    await createPDFViewer(pdfUrl);
  };

  const find = async (
    query: string,
    opts?: object,
    highlightOpts?: { twinkling?: boolean; color?: string }
  ) => {
    const finderViewer = pdfWebviewController?.getFinderViewer();
    const finderController = finderViewer?.getPdfFinderController();
    const { twinkling, color } = highlightOpts || {};
    if (typeof twinkling === 'boolean') {
      finderController?.toggleHighlightTwinkling(twinkling);
    }
    if (typeof color === 'string' && color) {
      finderController?.setHighlightColor(color);
    }
    finderController?._eventBus.dispatch('find', {
      query,
      type: '',
      phraseSearch: true,
      highlightAll: true,
      jumpToMatch: true,
      ...opts,
    });
  };

  return {
    documentWrapper,
    open,
    find,
    error: errorInfo,
  };
};

// TypeError: Cannot read from private field
// proxy + #privateField 有问题 不能存入ref中
export const usePDFViewerController = (timeout = 180000) => {
  // 实现一个Deferred函数
  const Deferred = () => {
    let _resolve: (value: ViewerController) => void;
    let _reject: (reason?: any) => void;
    const promise = new Promise<ViewerController>((resolve, reject) => {
      _resolve = (value: ViewerController) => {
        window.clearTimeout(timeoutTimer);
        resolve(value);
      };
      _reject = (reason?: any) => {
        window.clearTimeout(timeoutTimer);
        reject(reason);
      };
    });
    return {
      resolve: _resolve!,
      reject: _reject!,
      promise,
    };
  };

  const deferred = Deferred();

  const getPdfViewerController = () => {
    return deferred.promise;
  };

  const timeoutTimer = window.setTimeout(() => {
    deferred.reject(new Error('getPdfViewerController timeout'));
  }, timeout);

  return {
    deferred,
    getPdfViewerController,
  };
};
