import {
  SVG_CONTAINER_SELECTOR,
  PDF_ANNOTATE_ID,
  attrSelector,
} from '@idea/pdf-annotate-core';
import { getAnnotationRect } from '@idea/pdf-annotate-core/UI/utils';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import EditOverlay from '~/src/components/ToolTip/EditOverlay.vue';
import { createApp } from 'vue';
import { clearWhenScroll } from './arrow';
import { PDF_ANNOTATE_EDIT_OVERLAY } from '../constants';
import i18n from '../locals/i18n';
import { useAnnotationStore } from '../stores/annotationStore';

const OVERLAY_BORDER_SIZE = 1;

export function destroyEditOverlay(pdfViewer?: ViewerController) {
  const container = pdfViewer?.getDocumentViewer().container;
  if (!container) {
    return;
  }

  const list = container.getElementsByClassName(PDF_ANNOTATE_EDIT_OVERLAY);
  Array.from(list).forEach((el) => {
    el.parentNode?.removeChild(el);
  });
}

export function createEditOverlay(
  annotateId: string,
  pageNumber: number,
  style: Partial<CSSStyleDeclaration>,
  pdfViewer?: ViewerController
) {
  destroyEditOverlay(pdfViewer);

  const container = pdfViewer?.getDocumentViewer().container as HTMLDivElement;

  useAnnotationStore().activeOverlayPageNumber = pageNumber;

  const groupList = container.querySelectorAll<SVGGElement>(
    `g${attrSelector('uuid', annotateId)}`
  );
  const rectList = container.querySelectorAll<SVGRectElement>(
    `rect${attrSelector('uuid', annotateId)}`
  );
  const allList = [...Array.from(groupList), ...Array.from(rectList)];

  const annotatePageNumber = Number(allList[0]?.getAttribute('page-number'));

  const list = allList.map((group) => {
    return createOne(annotatePageNumber, group, style);
  });

  return list;
}

export function createOne(
  annotatePageNumber: number,
  target: SVGGElement,
  style: Partial<CSSStyleDeclaration>
) {
  const parentNode = target.closest(SVG_CONTAINER_SELECTOR)
    ?.parentElement as HTMLDivElement;
  const id = target.getAttribute(PDF_ANNOTATE_ID);
  const rectPageNumber = Number(target.getAttribute('page-number') as string);

  const rect = getAnnotationRect(target);

  const styleLeft = rect.left - OVERLAY_BORDER_SIZE;
  const styleTop = rect.top - OVERLAY_BORDER_SIZE;

  const app = createApp(EditOverlay, {
    annotatePageNumber,
    rectPageNumber,
    top: styleTop,
    left: styleLeft,
    width: rect.width,
    height: rect.height,
    annotateId: id,
    style,
  });

  app.use(i18n);

  const instance = app.mount(document.createElement('div'));

  parentNode.appendChild(instance.$el);

  return {
    rectPageNumber,
    div: instance.$el as HTMLDivElement,
  };
}

export const appendPopperElement = (
  popperElement: HTMLElement,
  style: { x: number; y: number; width: number; height: number }
) => {
  if (!popperElement) {
    return;
  }

  const { x, y, width, height } = style;

  popperElement.style.top = y + height + 'px';
  popperElement.style.left = x + width + 'px';
  popperElement.style.transform = `translate(0, -5px)`;

  popperElement.style.display = 'block';

  setTimeout(() => {
    clearWhenScroll(() => {
      popperElement.style.display = 'none';
    });
  }, 300);
};

export const removePopperElement = (popperElement: HTMLElement) => {
  if (!popperElement) {
    return;
  }

  popperElement.style.display = 'none';
};

export const hideTooltip = () => {
  const toolTips = document.getElementsByClassName(
    TOOLTIP_CLASSNAME
  ) as HTMLCollectionOf<HTMLDivElement>;

  Array.from(toolTips).forEach((tooltip) => {
    tooltip.style.display = 'none';
  });
};

export const TOOLTIP_CLASSNAME = 'pdf-content-tooltip';
