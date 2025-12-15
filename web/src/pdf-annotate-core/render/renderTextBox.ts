import { Canvas, Rect } from 'fabric';

import {
  AnnotateTextBox,
  Font,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Persistence';
// import { WebNoteAnnotationModel } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { WebNoteAnnotationModel } from 'go-sea-proto/gen/ts/note/Web'
import {
  PDF_ANNOTATE_TEXTSHAPE_BOX,
  PDF_ANNOTATE_TEXTSHAPE_SPAN,
  attrSelector,
} from '../constants';
import { numberToPx, pxToNumber } from '../utils/setAttributes';
import {
  PageFcanvasMap,
  ShapeToolbarProps,
  createDeleteButton,
  createToolbarDiv,
  getShapeDiv,
  getShapeScale,
  optimizeTextArrowPosition,
  stopPropagation,
} from './shapeCommon';
import { clearModelFrame } from './editHandwrite';
import { colorKeyMap, colorMap } from '../constants/color';
import { AnnotationColor } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { shapeCommonProps } from './renderShape';

export const PageTextRectMap = new Map<number, Map<string, Rect>>();

export type TextAnnotation = Pick<
  WebNoteAnnotationModel,
  'pageNumber' | 'textBox' | 'type'
>;

export interface TextCallback {
  edited(textItem: TextAnnotation): void;
  changeColor(color: AnnotationColor): void;
  deleted(id: string): void;
  confirmDelete(textItem: TextAnnotation): Promise<boolean>;
  createToolbarApp(
    toobarProps: ShapeToolbarProps,
    toolbarDiv: HTMLDivElement,
  ): void;
}

const lineHeight = '1.2';
export const HAND_BOX_EDIT_STYLE = 'data-pdf-annotate-textshape-box-edit';
export const HAND_BOX_SELECT_STYLE = 'data-pdf-annotate-textshape-box-select';

export const getTextBody = (annotation: TextAnnotation) => {
  return annotation.textBox?.fonts?.[0];
};

const getTextarea = (textDiv: HTMLDivElement) => {
  return textDiv.querySelector('textarea') as HTMLTextAreaElement;
};

export const getTextRectByAnnotation = (textItem: TextAnnotation) => {
  return PageTextRectMap.get(textItem.pageNumber)?.get(
    (textItem.textBox as AnnotateTextBox).id
  ) as Rect;
};

export const setTextRectByAnnotation = (
  textItem: TextAnnotation,
  frect: Rect
) => {
  if (!PageTextRectMap.has(textItem.pageNumber)) {
    PageTextRectMap.set(textItem.pageNumber, new Map());
  }

  PageTextRectMap.get(textItem.pageNumber)?.set(
    (textItem.textBox as AnnotateTextBox).id,
    frect
  );
};

export const deleteTextRectByAnnotation = (textItem: TextAnnotation) => {
  const frect = getTextRectByAnnotation(textItem);

  PageTextRectMap.get(textItem.pageNumber)?.delete(
    (textItem.textBox as AnnotateTextBox).id
  );

  PageFcanvasMap.get(textItem.pageNumber)?.remove(frect);
};

export const editingTextarea = {
  id: '',
  value: '',
};

export const editTextarea = (textDiv: HTMLDivElement) => {
  textDiv.classList.add(HAND_BOX_EDIT_STYLE);
  textDiv.dataset.dirty = String(1);
  const textarea = getTextarea(textDiv);

  setTimeout(() => {
    textarea.focus();
  }, 0);
};

export const setEditingAndEditTextarea = (
  id: string,
  textDiv: HTMLDivElement
) => {
  editingTextarea.id = id;
  editingTextarea.value = getTextarea(textDiv).value;
  editTextarea(textDiv);
};

export const conitnueEditTextarea = (textDiv: HTMLDivElement) => {
  const textarea = getTextarea(textDiv);
  textarea.value = editingTextarea.value;
  editTextarea(textDiv);
};

export const createTextDiv = (
  textItem: TextAnnotation,
  container: HTMLElement
) => {
  const textBox = textItem.textBox as AnnotateTextBox;
  const textBody = getTextBody(textItem) as Font;
  const textDiv = document.createElement('div');
  textDiv.classList.add(PDF_ANNOTATE_TEXTSHAPE_BOX);
  textDiv.setAttribute(PDF_ANNOTATE_TEXTSHAPE_BOX, textBox.id);
  Object.assign(textDiv.style, {
    top: numberToPx(textBox.y),
    left: numberToPx(textBox.x),
    fontSize: textBody.fontSize + 'em',
    color: textBody.color,
  } as CSSStyleDeclaration);

  const viewDiv = document.createElement('div');

  Object.assign(viewDiv.style, {
    width: numberToPx(textBox.width),
    height: numberToPx(textBox.height),
  } as CSSStyleDeclaration);

  const shapeScale = getShapeScale(textItem.pageNumber, container);
  const scaleFontSize = pxToNumber(shapeScale.style.fontSize);
  const span = document.createElement('span');
  const hws = Number(localStorage.getItem('hws')) || 1.2;
  span.classList.add(PDF_ANNOTATE_TEXTSHAPE_SPAN);
  span.style.margin = numberToPx((scaleFontSize * 6.6666) / hws);
  span.innerText = textBody.content;

  Object.assign(span.style, {
    lineHeight,
  } as CSSStyleDeclaration);

  span.addEventListener('pointerdown', stopPropagation);

  viewDiv.appendChild(span);

  const textarea = document.createElement('textarea');
  textarea.textContent = textBody.content;
  textarea.placeholder = '输入文本';
  Object.assign(textarea.style, {
    width: numberToPx(textBox.width),
    lineHeight,
    border: '0',
    backgroundColor: 'transparent',
    padding: span.style.margin,
  } as CSSStyleDeclaration);
  textarea.addEventListener('pointerdown', stopPropagation);

  textDiv.appendChild(viewDiv);
  textDiv.appendChild(textarea);

  return textDiv;
};

export const syncTextareaHeightToRect = (
  scale: number,
  textItem: TextAnnotation,
  textDiv: HTMLDivElement
) => {
  const textarea = getTextarea(textDiv);
  const fcanvas = PageFcanvasMap.get(textItem.pageNumber) as Canvas;
  const frect = getTextRectByAnnotation(textItem);

  const rectHeight = frect.getScaledHeight();
  const textareaHeight = textarea.getBoundingClientRect().height;

  if (textareaHeight > rectHeight) {
    frect.set({ height: textareaHeight / frect.scaleY });
    fcanvas.renderAll();
  }
};

export const onMouseDownStartEditText = (
  textDiv: HTMLDivElement,
  startEditText: () => void
) => {
  const span = textDiv.querySelector('span') as HTMLSpanElement;
  span.addEventListener('pointerdown', startEditText);
};

export const useStartEditText = (
  scale: number,
  textItem: TextAnnotation,
  textDiv: HTMLDivElement,
  callback: TextCallback,
  prepareShapes?: () => void
) => {
  return () => {
    prepareShapes?.();
    clearModelFrame();
    createTextRect(scale, textItem);
    initTextRect(scale, textItem, textDiv, callback);

    const fcanvas = PageFcanvasMap.get(textItem.pageNumber) as Canvas;
    const frect = getTextRectByAnnotation(textItem);
    fcanvas.add(frect);
    fcanvas.setActiveObject(frect);
    fcanvas.requestRenderAll();
  };
};

const autoSetHeight = (element: HTMLDivElement | HTMLTextAreaElement) => {
  element.style.height = '5px';
  element.style.height = numberToPx(element.scrollHeight);
};

export const initTextRect = (
  scale: number,
  textItem: TextAnnotation,
  textDiv: HTMLDivElement,
  callback: TextCallback
) => {
  Array.from(PageFcanvasMap.values()).forEach((fcanvas) => {
    fcanvas.discardActiveObject();
  });

  textDiv.classList.add(HAND_BOX_SELECT_STYLE);

  const textBox = textItem.textBox as AnnotateTextBox;
  const container = textDiv.closest<HTMLDivElement>('.rp-pdf-viewer-viewer');

  const fcanvas = PageFcanvasMap.get(textItem.pageNumber) as Canvas;
  const wrapper = fcanvas.getElement().parentElement as HTMLElement;
  const frect = getTextRectByAnnotation(textItem);

  frect.setControlVisible('mtr', false);

  const stopPropagation = (e: Event) => {
    if (fcanvas._hoveredTarget) {
      e.stopPropagation();
    }
  };
  wrapper.addEventListener('pointerdown', stopPropagation);
  frect.on('mousedown:before', (event) => stopPropagation(event.e));
  frect.once('removed', () =>
    wrapper.removeEventListener('pointerdown', stopPropagation)
  );

  let moved = false;
  const onMoving = () => {
    moved = true;

    optimizeTextArrowPosition(frect, textItem.pageNumber);

    if (toolbarDiv?.parentElement) {
      toolbarDiv.remove();
    }

    if (deleteButton?.parentElement) {
      deleteButton.remove();
    }

    textDiv.style.top = frect.top / scale + 'px';
    textDiv.style.left = frect.left / scale + 'px';
  };

  frect.on('moving', onMoving);

  frect.on('scaling', (event) => {
    console.log({ event, textDiv, scale });

    onMoving();

    Array.from(textDiv.children).forEach((child) => {
      const element = child as HTMLDivElement | HTMLTextAreaElement;
      const width = frect.getScaledWidth() / scale;
      element.style.width = numberToPx(width);
      autoSetHeight(element);
    });
  });

  const syncTextareaRectToolbarDelete = async () => {
    syncTextareaHeightToRect(scale, textItem, textDiv);

    const shapeBound = getShapeDiv(
      textItem.pageNumber,
      container!
    ).getBoundingClientRect();
    const y = shapeBound.top - container!.getBoundingClientRect().top;
    const x = shapeBound.left - container!.getBoundingClientRect().left;

    await new Promise((resolve) => setTimeout(resolve, 1));

    {
      const top = y + frect.top + frect.getScaledHeight() + 20;
      const left = x + frect.left;
      Object.assign(toolbarDiv.style, {
        top: numberToPx(top),
        left: numberToPx(left),
      } as CSSStyleDeclaration);
    }

    {
      const top = y + frect.top - 10;
      const left = x + frect.left + frect.getScaledWidth() - 10;
      Object.assign(deleteButton.style, {
        top: numberToPx(top),
        left: numberToPx(left),
      } as CSSStyleDeclaration);
    }
  };

  frect.on('modified', async () => {
    await syncTextareaRectToolbarDelete();

    setTextBoxByTextArea(textItem, textDiv);

    callback.edited(textItem);

    if (!deleteButton.parentElement) {
      container?.appendChild(deleteButton);
    }

    if (!toolbarDiv.parentElement) {
      container?.appendChild(toolbarDiv);
    }
  });

  frect.on('deselected', () => {
    deleteTextRectByAnnotation(textItem);
    cancelEditText(textItem, callback, textDiv);

    toolbarDiv.remove();
    deleteButton.remove();
  });

  frect.on('mouseup', () => {
    if (moved) {
      moved = false;
      return;
    }

    setEditingAndEditTextarea(textBox.id, textDiv);
  });

  const textarea = getTextarea(textDiv);
  autoSetHeight(textarea);
  textarea.addEventListener('input', () => {
    editingTextarea.id = textBox.id;
    editingTextarea.value = textarea.value;

    const originHeight = textarea.getBoundingClientRect().height;

    autoSetHeight(textarea);

    const newHeight = textarea.getBoundingClientRect().height;

    if (newHeight > originHeight) {
      syncTextareaRectToolbarDelete();
    }
  });

  const shapeDiv = getShapeDiv(textItem.pageNumber, container!);
  const toolbarDiv = createToolbarDiv();
  const deleteButton = createDeleteButton();

  {
    const shapeBound = shapeDiv.getBoundingClientRect();
    const y = shapeBound.top - container!.getBoundingClientRect().top;
    const x = shapeBound.left - container!.getBoundingClientRect().left;

    {
      const top = y + (textBox.y + textBox.height) * scale + 10;
      const left = x + textBox.x * scale;

      Object.assign(toolbarDiv.style, {
        top: numberToPx(top),
        left: numberToPx(left),
        width: '216px',
      } as CSSStyleDeclaration);
    }

    {
      const top = y + textBox.y * scale - 10;
      const left = x + (textBox.x + textBox.width) * scale - 10;
      Object.assign(deleteButton.style, {
        top: numberToPx(top),
        left: numberToPx(left),
      } as CSSStyleDeclaration);
    }
  }

  setTimeout(() => {
    container?.appendChild(toolbarDiv);
    container?.appendChild(deleteButton);
  }, 2);

  deleteButton.addEventListener('click', async () => {
    const deleted = await callback.confirmDelete(textItem);

    if (!deleted) {
      return;
    }

    deleteTextRectByAnnotation(textItem);
    textDiv.remove();
    toolbarDiv?.remove();
    deleteButton?.remove();
  });

  {
    const textBody = getTextBody(textItem) as Font;
    const fcanvas = PageFcanvasMap.get(textItem.pageNumber) as Canvas;

    const setColorToTextDivAndRect = (color: string) => {
      textDiv.style.color = color;
      frect.set({
        borderColor: color,
        cornerStrokeColor: color,
      });
      fcanvas.requestRenderAll();
    };

    callback.createToolbarApp(
      {
        originColor: colorKeyMap[textBody.color],
        onColorEnter(color: AnnotationColor) {
          setColorToTextDivAndRect(colorMap[color]);
        },
        onColorLeave() {
          setColorToTextDivAndRect(textBody.color);
        },
        selectColor(color: AnnotationColor) {
          textBody.color = colorMap[color];
          if (!Number(textDiv.dataset.dirty)) {
            textDiv.dataset.dirty = String(1);
          }
          callback.changeColor(color);
          fcanvas.discardActiveObject();
          setColorToTextDivAndRect(textBody.color);
        },
        originSize: textBody.fontSize,
        onSizeEnter() {
          // textDiv.style.fontSize = size + 'em';
        },
        onSizeLeave() {
          // textDiv.style.fontSize = textBody.fontSize + 'em';
        },
        selectSize(size: number) {
          textBody.fontSize = size;
          textDiv.style.fontSize = size + 'em';
          autoSetHeight(textarea);
          syncTextareaRectToolbarDelete();
          callback.edited(textItem);
        },
      },
      toolbarDiv
    );
  }
};

export const getTextDivByAnnotation = (
  container: HTMLElement,
  textItem: TextAnnotation
) => {
  return container.querySelector(
    `.${PDF_ANNOTATE_TEXTSHAPE_BOX}${attrSelector(
      PDF_ANNOTATE_TEXTSHAPE_BOX,
      (textItem.textBox as AnnotateTextBox).id
    )}`
  ) as HTMLDivElement;
};

export function cancelEditText(
  textItem: TextAnnotation,
  callback: Omit<TextCallback, 'confirmDelete' | 'createToolbarApp'>,
  textDiv: HTMLDivElement
) {
  editingTextarea.id = '';
  editingTextarea.value = '';

  const textarea = getTextarea(textDiv);
  const textBox = textItem.textBox as AnnotateTextBox;

  if (!textarea.value) {
    textDiv.remove();
    callback.deleted(textBox.id);
    return;
  }

  const viewDiv = textDiv.querySelector('div') as HTMLDivElement;

  viewDiv.style.width = textarea.style.width;
  viewDiv.style.height = textarea.style.height;
  (viewDiv.querySelector('span') as HTMLSpanElement).innerText = textarea.value;
  (getTextBody(textItem) as Font).content = textarea.value;

  setTextBoxByTextArea(textItem, textDiv);

  if (textDiv.classList.contains(HAND_BOX_EDIT_STYLE)) {
    textDiv.classList.remove(HAND_BOX_EDIT_STYLE);
  }
  textDiv.classList.remove(HAND_BOX_SELECT_STYLE);

  if (Number(textDiv.dataset.dirty)) {
    textDiv.dataset.dirty = String(0);
    callback.edited(textItem);
  }
}

const setTextBoxByTextArea = (
  textItem: TextAnnotation,
  textDiv: HTMLDivElement
) => {
  const textBox = textItem.textBox as AnnotateTextBox;
  const textarea = getTextarea(textDiv);
  textBox.height = pxToNumber(textarea.style.height);
  textBox.width = pxToNumber(textarea.style.width);
  textBox.x = pxToNumber(textDiv.style.left);
  textBox.y = pxToNumber(textDiv.style.top);
};

const createTextRect = (scale: number, textItem: TextAnnotation) => {
  const textBox = textItem.textBox as AnnotateTextBox;
  const color = (getTextBody(textItem) as Font).color;

  const frect = new Rect({
    top: textBox.y * scale,
    left: textBox.x * scale,
    width: textBox.width * scale,
    height: textBox.height * scale,
    fill: 'transparent',
    stroke: 'transparent',
    perPixelTargetFind: false,
    transparentCorners: false,
    hasRotatingPoint: false,
    borderColor: color,
    cornerColor: 'white',
    cornerStrokeColor: color,
    cornerStyle: 'circle',
    cornerSize: shapeCommonProps.cornerSize,
  });

  setTextRectByAnnotation(textItem, frect);
};
