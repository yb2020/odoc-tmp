import { PDF_ANNOTATE_ID } from '@idea/pdf-annotate-core';
import { useAnnotationStore } from '../stores/annotationStore';
import { usePdfStore } from '../stores/pdfStore';
import { currentNoteInfo } from '../store';
import { computed } from 'vue';
import { ViewerController } from '@idea/pdf-annotate-viewer';

export interface Arrow {
  annotationId: null | string;
  line: null | {
    color: string;
    remove(): void;
    position(): void;
  };
}

export const arrow: Arrow = {
  annotationId: null,
  line: null,
};

export const clearWhenScroll = (callback: () => void) => {
  const pdfStore = usePdfStore();
  const pdfViewer = pdfStore.getViewer(currentNoteInfo.value?.pdfId);
  const scrollContainer = pdfViewer?.getDocumentViewer().container;

  setTimeout(() => {
    const cb = () => {
      callback();
      scrollContainer?.removeEventListener('scroll', cb);
    };

    scrollContainer?.addEventListener('scroll', cb);
  }, 200);
};

const scrollCb = () => {
  clearArrow();
};

const clearArrowCb = () => {
  const pdfStore = usePdfStore();
  const pdfViewer = pdfStore.getViewer(currentNoteInfo.value?.pdfId);
  const scrollContainer = pdfViewer?.getDocumentViewer().container;

  clearArrow();

  document.removeEventListener('click', clearArrowCb);
  scrollContainer?.removeEventListener('scroll', scrollCb);

  const annotationStore = useAnnotationStore();
  annotationStore.currentAnnotationId = '';
};

const drawArrow = (
  pdfViewer: ViewerController,
  annotationId: string,
  from: Element,
  to: Element,
  color?: string
) => {
  const scrollContainer = pdfViewer.getDocumentViewer().container;

  clearArrow();

  arrow.annotationId = annotationId;
  arrow.line = new LeaderLine(from, to, {
    path: 'grid',
    size: 2,
    color: color || 'rgba(177, 74, 56, 1)',
    dash: true,
  });

  setTimeout(() => {
    document.addEventListener('click', clearArrowCb);
    scrollContainer.addEventListener('scroll', scrollCb);
  }, 200);
};

export const clearArrow = () => {
  if (!arrow.line) return;

  arrow.line.remove();
  arrow.line = null;
  arrow.annotationId = null;
};

export const connect = (pdfAnnotateId: string, isEdit = false) => {
  // 颜色编辑时没有连线的情况下不处理
  if (isEdit && (!arrow.line || arrow.annotationId !== pdfAnnotateId)) {
    return;
  }

  clearArrow();

  const pdfStore = usePdfStore();
  const pdfViewer = pdfStore.getViewer(currentNoteInfo.value?.pdfId);
  const container = pdfViewer?.getDocumentViewer().container;

  const commentRef = document.querySelector(
    `.note[${PDF_ANNOTATE_ID}="${pdfAnnotateId}"]`
  );

  const selectGroup =
    container?.querySelector<SVGGElement>(
      `g[${PDF_ANNOTATE_ID}="${pdfAnnotateId}"]`
    ) ||
    container?.querySelector<SVGRectElement>(
      `rect[${PDF_ANNOTATE_ID}="${pdfAnnotateId}"]`
    );

  if (!commentRef || !selectGroup) {
    return;
  }

  // 天培的代码；svg元素的display: none真的有效吗？
  if (selectGroup.style.display === 'none') {
    return;
  }

  const { width, height } = selectGroup?.getBoundingClientRect();

  if (width === 0 && height === 0) {
    return;
  }

  if (selectGroup.closest('svg')?.style.display === 'none') {
    return;
  }

  if (pdfViewer) {
    drawArrow(
      pdfViewer,
      pdfAnnotateId,
      commentRef,
      selectGroup,
      selectGroup.getAttribute('color')!
    );
  }
};
