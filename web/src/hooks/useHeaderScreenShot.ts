import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  ToolBarEndEvent,
  PDF_CANVASWRAPPER,
  PDF_ANNOTATE_ID,
  PDFJSAnnotate,
} from '@idea/pdf-annotate-core';
import { JS_IGNORE_MOUSE_OUTSIDE } from '@idea/pdf-annotate-viewer';
import { UserInfo } from '../api/user';
import { currentGroupId, currentNoteInfo } from '@/store';
import { ImageStorageType, uploadImage } from '../api/upload';
import { rectStyleMap } from '@/style/select';
import { AnnotationRect } from '../../src/stores/annotationStore/BaseAnnotationController';
import setAttributes from '../../src/pdf-annotate-core/utils/setAttributes';
import { useCommentGlobalState } from './useNoteState';
import {
  CLASS_NAME_ANNOTATION_SCREENSHOT_CAPTURING,
  useAnnotationStore,
} from '../stores/annotationStore';
import { Defer, ImageMimeType, compressImageBase64 } from '../../src/util/image';
import { usePdfStore } from '../stores/pdfStore';
import { SELF_NOTEINFO_GROUPID } from '../store/base/type';
import { uploadBase64File } from '@/utils/pdf-upload';
import { OSSBucketEnum, OSSKeyPolicyEnum } from 'go-sea-proto/gen/ts/oss/OSS';

export const ANNOTATION_SCREENSHOT = 'annotation:screenShot';

export interface ScreenShotPayload {
  rect: DOMRect;
  pageNum: number;
  visibleBtns: Partial<{
    ai: boolean;
    highlight: boolean;
    translate: boolean;
    comment: boolean;
  }>;
}

export interface ClipPayload {
  newRectAnnotation: AnnotationRect | null;
  newRectElement: SVGRectElement | null;
  newCanvasElement: HTMLCanvasElement | null;
}

export const SUBMIT_CLIP_PAYLOAD_TO_COPILTO =
  'submit_clip_payload_to_copilot' as const;

export const useClip = () => {
  const pdfStore = usePdfStore();
  const pdfViewerRef = computed(() => {
    return pdfStore.getViewer(currentNoteInfo.value?.pdfId);
  });
  const pdfAnnotater = computed(() => {
    return pdfStore.getAnnotater(currentNoteInfo.value?.noteId);
  });

  const commentState = useCommentGlobalState();

  const clipSelecting = ref(false);
  let newRectAnnotation: ClipPayload['newRectAnnotation'] = null;
  let newRectElement: ClipPayload['newRectElement'] = null;
  let newCanvasElement: ClipPayload['newCanvasElement'] = null;
  const getClipPayload = () => ({
    newRectElement,
    newCanvasElement,
    newRectAnnotation,
  });
  let defer: Defer<ClipPayload | null> | null = null;
  let noDashed = false;

  return {
    clipSelecting,
    clipAction: {
      init,
      cancelCut,
      clearClip,
      getClipPayload,
      onMouseDownClearClip,
    },
  };

  function clearClip() {
    newCanvasElement = null;
    newRectAnnotation = null;
    newRectElement = null;
  }

  async function screenShotEnd(
    annotation: AnnotationRect,
    svg: SVGSVGElement,
    rectElement: SVGRectElement
  ) {
    console.log('screenShotEnd');
    if (!pdfViewerRef.value || !pdfAnnotater.value) {
      clearDefer();
      return;
    }

    if (noDashed) {
      setAttributes(rectElement, {
        lineColor: rectStyleMap[commentState.value.styleId].color,
        strokeDasharray: 'none',
      });
    }

    const rect = rectElement.getBoundingClientRect();

    if (!rect || rect.width < 2 || rect.height < 2) {
      clearDefer();
      return;
    }

    newCanvasElement = getCanvas(svg, rect);
    newRectAnnotation = annotation;
    newRectElement = rectElement;

    defer?.resolve({
      newRectAnnotation,
      newRectElement,
      newCanvasElement,
    });

    cancelCut();
  }

  function getCanvas(svg: SVGSVGElement, rect: DOMRect) {
    const canvas = (svg.closest('.page') as HTMLDivElement).querySelector(
      `.${PDF_CANVASWRAPPER} canvas`
    ) as HTMLCanvasElement;

    const scale = canvas.width / parseInt(canvas.style.width);

    const svgRect = svg.getBoundingClientRect();

    const areaX = (rect.x - svgRect.x) * scale;
    const areaY = (rect.y - svgRect.y) * scale;
    const areaW = rect.width * scale;
    const areaH = rect.height * scale;
    const drawCanvasCtx = canvas.getContext('2d') as CanvasRenderingContext2D;

    const helper = document.createElement('canvas');

    const context = helper.getContext('2d') as CanvasRenderingContext2D;

    const data = drawCanvasCtx.getImageData(areaX, areaY, areaW, areaH);

    helper.width = areaW;

    helper.height = areaH;

    context.putImageData(data, 0, 0);

    return helper;
  }

  function init(noDashedBorder = false): Promise<ClipPayload | null> | void {
    noDashed = noDashedBorder;

    if (!pdfAnnotater.value) {
      clearDefer();
      return;
    }

    pdfAnnotater.value.onOffAnnotation(false);

    if (!pdfViewerRef.value) {
      clearDefer();
      return;
    }

    const UI = pdfAnnotater.value.UI;

    UI.on(ToolBarEndEvent.EventEnd, screenShotEnd);

    pdfViewerRef.value.getDocumentViewer().enableSelection(false);

    UI.rectController.enable();

    UI.rectController.setOptions({
      ...rectStyleMap[commentState.value.styleId],
      fill: rectStyleMap[commentState.value.styleId].color,
      fillOpacity: '0.1',
      strokeDasharray: '5,5',
      styleId: commentState.value.styleId,
    });

    const container = pdfViewerRef.value.getDocumentViewer().container;

    container.classList.add(CLASS_NAME_ANNOTATION_SCREENSHOT_CAPTURING);

    clipSelecting.value = true;
    defer = new Defer();
    return defer.promise;
  }

  function cancelCut() {
    if (!pdfAnnotater.value) {
      clearDefer();
      return;
    }

    pdfAnnotater.value.onOffAnnotation(true);

    if (!pdfViewerRef.value) {
      clearDefer();
      return;
    }

    const UI = pdfAnnotater.value.UI;

    pdfViewerRef.value.getDocumentViewer().enableSelection(true);

    UI.rectController.disable();

    const container = pdfViewerRef.value.getDocumentViewer().container;

    container.classList.remove(CLASS_NAME_ANNOTATION_SCREENSHOT_CAPTURING);

    UI.off(ToolBarEndEvent.EventEnd, screenShotEnd);
    clearDefer();
  }

  function clearDefer() {
    if (clipSelecting.value) {
      clipSelecting.value = false;
    }

    if (defer) {
      defer.resolve(null);
      defer = null;
    }
  }

  function onMouseDownClearClip() {
    const clear = (event: MouseEvent) => {
      if (
        !newRectElement ||
        newRectElement.getAttribute(PDF_ANNOTATE_ID) ||
        (event.target instanceof HTMLElement &&
          event.target.closest('.' + JS_IGNORE_MOUSE_OUTSIDE))
      ) {
        return;
      }

      newRectElement.remove();
      clearClip();
    };

    onMounted(() => {
      document.body.addEventListener('pointerdown', clear, { passive: true });
    });

    onUnmounted(() => {
      document.body.removeEventListener('pointerdown', clear);
    });
  }
};

const LOADING_ANNOTATION_ID = 'loading_annotation_id';

export async function saveRectAnnotation(
  pdfAnnotater: PDFJSAnnotate,

  userInfo: UserInfo | null,
  annotationStore: ReturnType<typeof useAnnotationStore>,

  newRectAnnotation: AnnotationRect,
  newRectElement: SVGRectElement,
  newCanvasElement: HTMLCanvasElement,

  idea?: string
) {
  let uuid: string | null = null;

  if (typeof idea === 'string') {
    newRectAnnotation.idea = idea;
  }

  const setRectElementId = (id: string) => {
    pdfAnnotater.UI.rectController.updateUuid(newRectElement, id);
  };

  setRectElementId(LOADING_ANNOTATION_ID);

  try {
    uuid = await annotationStore.controller.onlineSaveAnnotation(
      newRectAnnotation.documentId,
      newRectAnnotation.pageNumber,
      newRectAnnotation
    );
  } catch (error) {
    console.error(error);
  }

  if (!uuid) {
    setRectElementId('');
    return false;
  }

  setRectElementId(newRectAnnotation.uuid);

  annotationStore.controller.localAddAnnotation({
    ...newRectAnnotation,
    picUrl: newCanvasElement?.toDataURL(ImageMimeType.PNG, 1),
    commentatorInfoView: {
      nickName: userInfo?.nickName ?? '',
      userId: userInfo?.id ?? '',
      avatarCdnUrl: userInfo?.avatarUrl ?? '',
    },
  } as AnnotationRect);

  if (currentGroupId.value === SELF_NOTEINFO_GROUPID) {
    annotationStore.showHoverNote({
      uuid: newRectAnnotation.uuid,
      page: newRectAnnotation.pageNumber,
    });
  }

  await nextTick();

  annotationStore.currentAnnotationId = newRectAnnotation.uuid;

  return true;
}

export async function uploadImageAndUpdateAnnotation(
  annotationStore: ReturnType<typeof useAnnotationStore>,

  newRectAnnotation: AnnotationRect,
  newCanvasElement: HTMLCanvasElement,

  compressRatio = 1
) {
  const image = compressImageBase64(
    newCanvasElement as HTMLCanvasElement,
    compressRatio
  );

  try {
    const updateResult = await uploadBase64File(image, 'png', OSSBucketEnum.PUBLIC, OSSKeyPolicyEnum.UPLOAD_PUBLIC_SHORT);
    console.log('updateResult', updateResult);
    const picUrl = updateResult.publicUrl;

    annotationStore.controller.patchAnnotation(newRectAnnotation.uuid, {
      picUrl,
      strokeDasharray: 'none',
    } as Partial<AnnotationRect>);
  } catch (e) {
    annotationStore.controller.localDeleteAnnotation(
      newRectAnnotation.uuid,
      newRectAnnotation.pageNumber
    );

    message.error('截图上传异常，请稍后重试!');
  }
}
