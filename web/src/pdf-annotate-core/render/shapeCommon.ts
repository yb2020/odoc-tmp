import { Canvas, Ellipse, Line, Rect } from 'fabric';
import { AnnotationColor } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import {
  PDF_ANNOTATE_SHAPE_DIV,
  PDF_ANNOTATE_TEXTSHAPE_SCALE,
} from '../constants';
import { getDivByPageNumberFrom1 } from '../UI/utils';

export const PageFcanvasMap = new Map<number, Canvas>();

export const stopPropagation = (event: MouseEvent) => {
  event.stopPropagation();
};

export interface ShapeToolbarProps {
  originColor: AnnotationColor;
  onColorEnter(color: AnnotationColor): void;
  onColorLeave(): void;
  selectColor(color: AnnotationColor): void;
  originSize?: number;
  onSizeEnter?(size: number): void;
  onSizeLeave?(): void;
  selectSize?(fontSize: number): void;
}

const SHAPE_TOOLBAR = 'shape-toolbar';
export const SHAPE_DELETE = 'shape-delete';

export const getShapeDiv = (
  pageNumberFrom1: number,
  container?: HTMLElement
) => {
  const shapeDiv = getDivByPageNumberFrom1(
    pageNumberFrom1,
    container
  )?.querySelector('.' + PDF_ANNOTATE_SHAPE_DIV) as HTMLDivElement;

  return shapeDiv;
};

export const getShapeScale = (
  pageNumberFrom1: number,
  container?: HTMLElement
) => {
  return getShapeDiv(pageNumberFrom1, container).querySelector(
    '.' + PDF_ANNOTATE_TEXTSHAPE_SCALE
  ) as HTMLDivElement;
};

export function createToolbarDiv(scale = 1) {
  const toolbar = document.createElement('div');
  toolbar.classList.add(SHAPE_TOOLBAR);

  Object.assign(toolbar.style, {
    position: 'absolute',
    width: '128px',
    height: '0',
    transform: `scale(${1 / scale})`,
    transformOrigin: 'top left',
    zIndex: '300',
  } as CSSStyleDeclaration);

  toolbar.addEventListener('pointerdown', stopPropagation);
  toolbar.addEventListener('click', stopPropagation);

  return toolbar;
}

export const createDeleteButton = (scale = 1) => {
  const deleteButton = document.createElement('div');
  deleteButton.classList.add(SHAPE_DELETE);

  deleteButton.innerHTML = `
  <span
    role="img"
    aria-label="close"
    class="anticon anticon-close"
  >
    <svg
      focusable="false"
      class=""
      data-icon="close"
      width="1em"
      height="1em"
      fill="currentColor"
      aria-hidden="true"
      viewBox="64 64 896 896"
    >
      <path
        d="M563.8 512l262.5-312.9c4.4-5.2.7-13.1-6.1-13.1h-79.8c-4.7 0-9.2 2.1-12.3 5.7L511.6 449.8 295.1 191.7c-3-3.6-7.5-5.7-12.3-5.7H203c-6.8 0-10.5 7.9-6.1 13.1L459.4 512 196.9 824.9A7.95 7.95 0 00203 838h79.8c4.7 0 9.2-2.1 12.3-5.7l216.5-258.1 216.5 258.1c3 3.6 7.5 5.7 12.3 5.7h79.8c6.8 0 10.5-7.9 6.1-13.1L563.8 512z"
      ></path>
    </svg>
  </span> 
  `;

  deleteButton.addEventListener('pointerdown', stopPropagation);
  deleteButton.addEventListener('click', stopPropagation);
  Object.assign(deleteButton.style, {
    position: 'absolute',
    top: '0',
    height: '20px',
    width: '20px',
    cursor: 'pointer',
    background: '#dadfe6',
    color: 'black',
    fontSize: '12px',
    boxShadow: '0px 2px 4px 0px rgba(0, 0, 0, 0.2)',
    borderRadius: '2px',
    border: '1px solid #d3d6d8',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    transform: `scale(${1 / scale})`,
    transformOrigin: 'top left',
    userSelect: 'none',
    zIndex: '300',
  } as CSSStyleDeclaration);

  return deleteButton;
};

export const clearShapeToolbar = (container: HTMLElement) => {
  const toolbars = container.querySelectorAll(`.${SHAPE_TOOLBAR}`);

  toolbars.forEach((toolbar) => {
    toolbar.remove();
  });

  const deletes = container.querySelectorAll(`.${SHAPE_DELETE}`);
  deletes.forEach((del) => {
    del.remove();
  });
};

export const optimizeRectEllipsePosition = (
  fobj: Rect | Ellipse,
  pageNumber: number
) => {
  const fcanvas = PageFcanvasMap.get(pageNumber) as Canvas;
  const fwidth = fcanvas.getWidth();
  const fheight = fcanvas.getHeight();

  const halfWidth =
    fobj instanceof Ellipse ? fobj.getRx() : fobj.getScaledWidth() / 2;
  const halfHeight =
    fobj instanceof Ellipse ? fobj.getRy() : fobj.getScaledHeight() / 2;
  const middleX = fobj.left + halfWidth;
  const middleY = fobj.top + halfHeight;

  const optimizePosition: {
    left?: number;
    top?: number;
  } = {};

  if (middleX > fwidth) {
    optimizePosition.left = fwidth - halfWidth;
  } else if (middleX < 0) {
    optimizePosition.left = -halfWidth;
  }

  if (middleY > fheight) {
    optimizePosition.top = fheight - halfHeight;
  } else if (middleY < 0) {
    optimizePosition.top = -halfHeight;
  }

  if (
    optimizePosition.left !== undefined ||
    optimizePosition.top !== undefined
  ) {
    fobj.set(optimizePosition);
  }
};

export const optimizeTextArrowPosition = (
  fobj: Rect | Line,
  pageNumber: number
) => {
  const fcanvas = PageFcanvasMap.get(pageNumber) as Canvas;
  const fwidth = fcanvas.getWidth();
  const fheight = fcanvas.getHeight();

  const lwidth = fobj.getScaledWidth();
  const lheight = fobj.getScaledHeight();

  const optimizePosition: {
    left?: number;
    top?: number;
  } = {};

  if (fobj.left < 0) {
    optimizePosition.left = 0;
  } else if (fobj.left + lwidth > fwidth) {
    optimizePosition.left = fwidth - lwidth;
  }

  if (fobj.top < 0) {
    optimizePosition.top = 0;
  } else if (fobj.top + lheight > fheight) {
    optimizePosition.top = fheight - lheight;
  }

  if (
    optimizePosition.left !== undefined ||
    optimizePosition.top !== undefined
  ) {
    fobj.set(optimizePosition);
  }
};
