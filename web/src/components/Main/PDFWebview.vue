<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import {
  createPDFWebview,
  ViewerEvent,
  PDFViewerScale,
  ViewerController,
} from '@idea/pdf-annotate-viewer';
import pdfjsWorkerSrc from '@idea/pdfjs-dist/build/pdf.worker.js?url';

import PDFError from './PDFError.vue';
import { debounce } from 'lodash-es';
import PDFViewerThemeAddonInstance from '~/src/theme';
import { PDFViewerColorToneAddon } from '@idea/pdf-annotate-viewer';
import reporter from '@idea/aiknowledge-report';
import { EventCode } from '~/src/api/report';
import i18n from '~/src/locals/i18n';
import { useMouse } from './mouseCore';
import { useI18n } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { standardToLanguageEnum } from '~/src/shared/language/service';
import { isInTauri } from '@/util/env';

const props = defineProps<{
  pdfUrl: string;
  scale?: PDFViewerScale;
  pdfId: string;
  showAnnotations?: boolean;
}>();

const { t } = useI18n();

const documentWrapper = ref();

const emit = defineEmits<{
  (event: 'created', payload: { pdfWebview: ViewerController }): void;
  (event: 'finished', payload: { pdfWebview: ViewerController }): void;
}>();

const errorInfo = ref({
  title: '',
  log: '',
});

const errorVisible = ref(false);

let pdfWebviewController: ViewerController;

let redrawConfig = {
  visible: true,
  needRedraw: false,
};

const redraw = () => {
  if (!pdfWebviewController) {
    return;
  }
  if (!redrawConfig.visible) {
    redrawConfig.needRedraw = true;
    return;
  }
  pdfWebviewController.redraw();
  redrawConfig.needRedraw = false;
};

const getPDFRetryUrl = (pdfUrl: string) => {
  const url = new URL(pdfUrl);
  if (url.host === 'static.cdn.readpaper.com') {
    url.host = 'staticcdn.readpaper.com';
  } else {
    url.host = 'downcdn.readpaper.com';
  }
  return url.href;
};

/**
 * 判断是否为本地文件路径
 */
const isLocalFilePath = (path: string) => {
  return path.startsWith('/') || path.startsWith('file://') || /^[A-Za-z]:[\\/]/.test(path);
};

/**
 * Tauri 环境下读取本地文件并转换为 Blob URL
 */
const getLocalPdfUrl = async (localPath: string): Promise<string> => {
  const { readFile } = await import('@tauri-apps/plugin-fs');
  const fileData = await readFile(localPath);
  const blob = new Blob([fileData], { type: 'application/pdf' });
  return URL.createObjectURL(blob);
};

onMounted(async () => {
  try {
    /**
     * pdfWebview不能放入store里面，因为这里加上proxy后，会导致pdfjs-dist调用#这种private的方法报错
     * TypeError: Cannot read from private field
     * https://esdiscuss.org/topic/why-does-a-javascript-class-getter-for-a-private-field-fail-using-a-proxy
     */
    
    // Tauri 环境下，如果是本地路径，先转换为 Blob URL
    let pdfUrl = props.pdfUrl;
    if (isInTauri() && isLocalFilePath(props.pdfUrl)) {
      pdfUrl = await getLocalPdfUrl(props.pdfUrl);
    }

    const pdfWebview = createPDFWebview(
      {
        pdfDocumentParams: {
          url: pdfUrl,
          // url: 'https://pdf.cdn.readpaper.com/aiKnowledge/pdf/112021-12-22/e82a7d52c316479e92b20f43f51fb56e-/userUpload/45207566b9a5ab22f51c801821da1138.pdf',
          cMapUrl:
            'https://nuxt.cdn.readpaper.com/pdfjs-dist%402.13.216/cmaps/',
          retryUrl: isLocalFilePath(props.pdfUrl) ? undefined : getPDFRetryUrl(props.pdfUrl),
          language: (() => {
            const langEnum = standardToLanguageEnum(i18n.global.locale.value);
            return langEnum !== null ? langEnum : Language.EN_US;
          })(),
        },
        containers: {
          documentWrapper: documentWrapper.value,
        },
      },
      pdfjsWorkerSrc
    );

    pdfWebview.addEventListener(ViewerEvent.PDFURL_RETRY, (payload) => {
      reporter.report(
        {
          event_code: EventCode.readpaperNoteCrashRetrySuccess,
        },
        {
          pdf_id: props.pdfId,
          pdf_url: payload.pdfUrl,
          note_id: new URL(window.location.href).searchParams.get('noteId'),
        }
      );
    });

    pdfWebviewController = pdfWebview;
    emit('created', { pdfWebview });

    const { documentViewer } = await pdfWebview.build(props.scale || '1.0', {
      annotationMode: props.showAnnotations === false ? 4 : 1,
    });


    useMouse(documentViewer);

    emit('finished', {
      pdfWebview,
    });

    try {
      /**
       * 如果是自适应的设置，需要动态调整大小
       */
      const resizeObserver = new ResizeObserver(
        debounce(
          (_: ResizeObserverEntry[]) => {
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

      resizeObserver.observe(documentWrapper.value);
    } catch (error) {
      console.error(error);
    }

    PDFViewerThemeAddonInstance.addEventListener(
      PDFViewerColorToneAddon.ForceRedrawEventType,
      redraw
    );

    // 首次加载的时候是不是就是不可见，如果不可见需要在可见后redraw
    let needFirstRedraw: boolean;

    const intersectionObserver = new IntersectionObserver(
      debounce((args: IntersectionObserverEntry[]) => {
        if (args[0]?.isIntersecting === false) {
          // 不可见，那么就不能redraw，需要等到可见的时候再redraw
          if (needFirstRedraw === undefined) {
            redrawConfig.needRedraw = true;
          }
          needFirstRedraw = true;
          redrawConfig.visible = false;
        } else {
          needFirstRedraw = false;
          redrawConfig.visible = true;
          if (redrawConfig.needRedraw) {
            redraw();
          }
        }
      }),
      {
        root: null,
        threshold: [0],
      }
    );

    intersectionObserver.observe(documentWrapper.value);
  } catch (error) {
    console.error('PDF initialization error:', error);
    
    // 临时禁用错误弹窗，用于调试
    console.warn('PDF error popup temporarily disabled for debugging:', {
      message: (error as Error).message,
      pdfId: props.pdfId,
      stack: (error as Error).stack
    });
    
    // 保留错误信息设置，但不显示弹窗
    if ((error as Error).message === 'Invalid PDF structure.') {
      errorInfo.value = {
        log: `error: ${(error as Error).message} pdfId: ${props.pdfId}`,
        title: t('viewer.pdfError.formatError'),
      };
    } else {
      errorInfo.value = {
        log: `error: ${(error as Error).message} pdfId: ${props.pdfId}`,
        title: t('viewer.pdfError.downloadError'),
      };
    }
    
    reporter.report(
      {
        event_code: EventCode.readpaperNoteCrashPopupImpression,
      },
      {
        pdf_id: props.pdfId,
        note_id: new URL(window.location.href).searchParams.get('noteId'),
        error_message: (error as Error).message,
      }
    );
    
    errorVisible.value = true;
  }
});

onUnmounted(() => {
  PDFViewerThemeAddonInstance.removeListener(
    PDFViewerColorToneAddon.ForceRedrawEventType,
    redraw
  );
});
</script>

<template>
  <div class="document-wrap">
    <div
      ref="documentWrapper"
      class="document"
    />
    <PDFError
      v-model:visible="errorVisible"
      :title="errorInfo.title"
      :log="errorInfo.log"
    />
  </div>
</template>

<style scoped lang="less">
.document-wrap {
  height: 100%;
  flex: 1;
  margin: 0 2px;
  display: flex;
  position: relative;
  z-index: 1;
  
  .document {
    height: 100%;
    flex: 1;
    background-color: var(--site-theme-bg-dark);
  }
}
</style>
