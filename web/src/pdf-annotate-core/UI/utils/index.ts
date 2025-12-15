import createStyleSheet from 'create-stylesheet';
import isElVisible from 'element-visible';
import {
  attrSelector,
  DATA_PAGE_NUMBER,
  PDF_ANNOTATE_CONTAINER,
  PDF_ANNOTATE_DOCUMENT,
  PDF_ANNOTATE_PAGENUMBER,
  PDF_ANNOTATE_TYPE,
  PDF_ANNOTATE_USER_SELECT,
  PDF_ANNOTATE_VIEWPORT,
  PDF_ANNOTATIONLAYER,
  PDF_ANNOTATIONLAYER_GROUP_NOTES,
  // READPAPER_PAGE_CONTAINER_,
} from '../../constants';
import { Rectangle, SelectAnnotation } from '../../render/renderSelect';
import { ViewerController } from '@idea/pdf-annotate-viewer';

export const BORDER_COLOR = '#00BFFF';

const userSelectStyleSheet = createStyleSheet({
  body: {
    '-webkit-user-select': 'none',
    '-moz-user-select': 'none',
    '-ms-user-select': 'none',
    'user-select': 'none',
  } as unknown as CSSStyleDeclaration,
});
userSelectStyleSheet.setAttribute(PDF_ANNOTATE_USER_SELECT, 'true');

export const SVG_CONTAINER_SELECTOR = `svg${attrSelector(PDF_ANNOTATE_CONTAINER, String(true))}`;
/**
 * Find an SVGElement container at a given point
 */
export function findSVGsAtPoint(x: number, y: number) {
  // @TODO 替换document
  const elements = document.querySelectorAll<SVGSVGElement>(
    // 这里需要兼容单词笔记，获取所有svg
    SVG_CONTAINER_SELECTOR
  );
  const res = [];

  for (let i = 0, l = elements.length; i < l; i++) {
    const el = elements[i];
    const rect = el.getBoundingClientRect();

    if (pointIntersectsRect(x, y, rect)) {
      res.push(el);
    }
  }

  return res;
}

export function findSVGAtPoint(x: number, y: number) {
  const svgs = findSVGsAtPoint(x, y);
  const svg = svgs.find((x) => x.classList.contains(PDF_ANNOTATIONLAYER));

  return svg ?? null;
}

/**
 * Find Elements that represents an annotation at a given point
 */
export function findAnnotationListAtPoint(x: number, y: number) {
  const annotationList: (SVGGElement | SVGRectElement)[] = [];

  const svgs = findSVGsAtPoint(x, y);

  if (!svgs.length) {
    return annotationList;
  }

  const elementList = svgs.reduce(
    (arr, svg) => {
      return [
        ...arr,
        ...Array.from(
          svg.querySelectorAll<SVGGElement | SVGRectElement>(
            attrSelector(PDF_ANNOTATE_TYPE)
          )
        ),
      ];
    },
    [] as Array<SVGGElement | SVGRectElement>
  );

  for (let i = 0, l = elementList.length; i < l; i += 1) {
    const element = elementList[i];

    const elementOffset = getOffsetAnnotationRect(element);
    if (!pointIntersectsRect(x, y, elementOffset)) {
      continue;
    }

    if (element instanceof SVGRectElement && isElVisible(element)) {
      annotationList.push(element);
      continue;
    }

    for (let j = 0; j < element.children.length; j += 1) {
      const childOffset = getOffsetAnnotationRect(
        element.children[j] as SVGRectElement,
        true
      );
      if (pointIntersectsRect(x, y, childOffset)) {
        annotationList.push(element);
      }
    }
  }

  return annotationList;
}

/**
 * Determine if a point intersects a rect
 *
 * @param {Number} x The x coordinate of the point
 * @param {Number} y The y coordinate of the point
 * @param {Object} rect The points of a rect (likely from getBoundingClientRect)
 * @return {Boolean} True if a collision occurs, otherwise false
 */
export function pointIntersectsRect(
  x: number,
  y: number,
  rect: Pick<DOMRect, 'top' | 'bottom' | 'right' | 'left'>
) {
  if (!rect) {
    return false;
  }

  return y >= rect.top && y <= rect.bottom && x >= rect.left && x <= rect.right;
}

/**
 * Get the rect of an annotation element accounting for offset.
 */
export function getOffsetAnnotationRect(el: SVGGElement, rectAsG = false) {
  const rect = getAnnotationRect(el, rectAsG);

  if (!rect) {
    return {
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
    };
  }

  const { offsetLeft, offsetTop } = getOffset(el);
  return {
    top: rect.top + offsetTop,
    left: rect.left + offsetLeft,
    right: rect.right + offsetLeft,
    bottom: rect.bottom + offsetTop,
  };
}

/**
 * Get the rect of an annotation element.
 */
export function getAnnotationRect(el: Element, rectAsG = false) {
  let h = 0,
    w = 0,
    x = 0,
    y = 0;
  const rect = el.getBoundingClientRect();

  // TODO this should be calculated somehow
  const LINE_OFFSET = 16;

  const nodeName = el.nodeName.toLowerCase();

  function calcGroupOrRect() {
    const { offsetLeft, offsetTop } = getOffset(
      el as SVGGElement | SVGRectElement
    );
    h = rect.height;
    w = rect.width;
    x = rect.left - offsetLeft;
    y = rect.top - offsetTop;

    // 引文是一整段文字，鼠标滑过单词左右、上下间隙时，需要防止批注tips闪烁
    if (nodeName === 'rect' && rectAsG) {
      h += 6;
      w += 6;
      x -= 3;
      y -= 3;
    } else if (el.getAttribute(PDF_ANNOTATE_TYPE) === 'strikeout') {
      h += LINE_OFFSET;
      y -= LINE_OFFSET / 2;
    }
  }

  switch (nodeName) {
    case 'path':
      let minX: number, maxX: number, minY: number, maxY: number;

      el.getAttribute('d')!
        .replace(/Z/, '')
        .split('M')
        .splice(1)
        .forEach((p) => {
          const s = p.split(' ').map((i) => parseInt(i, 10));

          if (typeof minX === 'undefined' || s[0] < minX) {
            minX = s[0];
          }
          if (typeof maxX === 'undefined' || s[2] > maxX) {
            maxX = s[2];
          }
          if (typeof minY === 'undefined' || s[1] < minY) {
            minY = s[1];
          }
          if (typeof maxY === 'undefined' || s[3] > maxY) {
            maxY = s[3];
          }
        });

      h = maxY! - minY!;
      w = maxX! - minX!;
      x = minX!;
      y = minY!;
      break;

    case 'line':
      h =
        parseInt(el.getAttribute('y2')!, 10) -
        parseInt(el.getAttribute('y1')!, 10);
      w =
        parseInt(el.getAttribute('x2')!, 10) -
        parseInt(el.getAttribute('x1')!, 10);
      x = parseInt(el.getAttribute('x1')!, 10);
      y = parseInt(el.getAttribute('y1')!, 10);

      if (h === 0) {
        h += LINE_OFFSET;
        y -= LINE_OFFSET / 2;
      }
      break;

    case 'text':
      h = rect.height;
      w = rect.width;
      x = parseInt(el.getAttribute('x')!, 10);
      y = parseInt(el.getAttribute('y')!, 10) - h;
      break;

    case 'g':
      calcGroupOrRect();
      break;

    case 'rect':
      if (rectAsG) {
        calcGroupOrRect();
        break;
      }

      h = parseInt(el.getAttribute('height') as string, 10);
      w = parseInt(el.getAttribute('width') as string, 10);
      x = parseInt(el.getAttribute('x') as string, 10);
      y = parseInt(el.getAttribute('y') as string, 10);

      h += 3;
      w += 3;
      x += 1.5;
      y += 1.5;

      break;
    case 'circle':
      h = parseInt(
        ((el.getAttribute('r') as unknown as number) * 2) as unknown as string,
        10
      );
      w = parseInt(
        ((el.getAttribute('r') as unknown as number) * 2) as unknown as string,
        10
      );
      x = parseInt(
        ((el.getAttribute('cx') as unknown as number) -
          h / 2) as unknown as string,
        10
      );
      y = parseInt(
        ((el.getAttribute('cy') as unknown as number) -
          h / 2) as unknown as string,
        10
      );

      break;
    case 'svg':
      h = parseInt(el.getAttribute('height')!, 10);
      w = parseInt(el.getAttribute('width')!, 10);
      x = parseInt(el.getAttribute('x')!, 10);
      y = parseInt(el.getAttribute('y')!, 10);
      break;
  }

  // Result provides same properties as getBoundingClientRect
  let result = {
    top: y,
    left: x,
    width: w,
    height: h,
    right: x + w,
    bottom: y + h,
  };

  // For the case of nested SVG (point annotations) and grouped
  // lines or rects no adjustment needs to be made for scale.
  // I assume that the scale is already being handled
  // natively by virtue of the `transform` attribute.
  if (
    nodeName !== 'svg' &&
    nodeName !== 'g' &&
    (nodeName !== 'rect' || !rectAsG)
  ) {
    const svg = findSVGAtPoint(rect.left, rect.top);
    if (svg) {
      result = scaleUp(svg, result) || result;
    }
  }

  return result;
}

interface Rect {
  [key: string]: number;
}

/**
 * Adjust scale from normalized scale (100%) to rendered scale.
 */
export function scaleUp(svg: SVGElement, rect: Rect) {
  if (!svg) {
    return;
  }
  const { viewport } = getMetadata(svg);

  return scaleUpRaw(viewport.scale, rect);
}

export function scaleUpRaw(scale: number, rect: Rect) {
  const result: any = {};

  Object.keys(rect).forEach((key) => {
    result[key] = rect[key] * scale;
  });

  return result;
}

/**
 * Adjust scale from rendered scale to a normalized scale (100%).
 *
 * @param {SVGElement} svg The SVG to gather metadata from
 * @param {Object} rect A map of numeric values to scale
 * @return {Object} A copy of `rect` with values scaled down
 */
export function scaleDown(svg: SVGElement, rect: Rect) {
  const { viewport } = getMetadata(svg);

  return scaleDownRaw(viewport.scale, rect);
}

export function scaleDownRaw(scale: number, rect: Rect) {
  const result: Rect = {} as Rect;

  Object.keys(rect).forEach((key) => {
    result[key as keyof Rect] = rect[key as keyof Rect] / scale;
  });

  return result;
}

export function scaleRect(rect: Rect) {
  const result: Rect = {} as Rect;

  const { viewport } = getMetadata(
    document.querySelector(`svg.${PDF_ANNOTATIONLAYER}`)!
  );

  Object.keys(rect).forEach((key) => {
    result[key as keyof Rect] = rect[key as keyof Rect] * viewport.scale;
  });

  return result;
}

/**
 * Get the scroll position of an element, accounting for parent elements
 *
 * @param {Element} el The element to get the scroll position for
 * @return {Object} The scrollTop and scrollLeft position
 */
export function getScroll(el: Element) {
  let scrollTop = 0;
  let scrollLeft = 0;
  let parentNode = el;

  while (
    (parentNode = parentNode.parentNode as HTMLElement) &&
    parentNode !== (document as any)
  ) {
    scrollTop += parentNode.scrollTop;
    scrollLeft += parentNode.scrollLeft;
  }

  return { scrollTop, scrollLeft };
}

/**
 * Get the offset position of an element, accounting for parent elements
 *
 * @param {Element} el The element to get the offset position for
 * @return {Object} The offsetTop and offsetLeft position
 */
function getOffset(el: SVGGElement) {
  const svg = el.closest('svg') as SVGSVGElement;
  const { left, top } = svg.getBoundingClientRect();
  return {
    offsetLeft: left,
    offsetTop: top,
  };
}

/**
 * Disable user ability to select text on page
 */
export function disableUserSelect() {
  if (!userSelectStyleSheet.parentNode) {
    // 直接 appendChild 好像不起作用，先 setTimeout
    // document.head.appendChild(userSelectStyleSheet);

    setTimeout(() => {
      document.head.appendChild(userSelectStyleSheet);
    }, 0);
  }
}

/**
 * Enable user ability to select text on page
 */
export function enableUserSelect() {
  if (userSelectStyleSheet.parentNode) {
    userSelectStyleSheet.parentNode.removeChild(userSelectStyleSheet);
  }
}

/**
 * Get the metadata for a SVG container
 *
 * @param {SVGElement} svg The SVG container to get metadata for
 */
export function getMetadata(svg: SVGElement) {
  return {
    documentId: svg.getAttribute(PDF_ANNOTATE_DOCUMENT)!,
    pageNumber: parseInt(svg.getAttribute(PDF_ANNOTATE_PAGENUMBER)!, 10),
    viewport: JSON.parse(svg.getAttribute(PDF_ANNOTATE_VIEWPORT)!),
  };
}

export const setDisplay = (
  element: SVGElement | HTMLCanvasElement | HTMLDivElement,
  visible: boolean
) => {
  const display = visible ? 'block' : 'none';
  if (element.style.display !== display) {
    element.style.display = display;
  }
};

export interface SvgBound {
  svg: SVGSVGElement;
  g: SVGGElement;
  bound: DOMRect;
}

export const getDivByPageNumberFrom1 = (
  pageNumberFrom1: number,
  container = document.body
) => {
  return container.querySelector<HTMLDivElement>(
    // '#' + READPAPER_PAGE_CONTAINER_ + pageNumberFrom1
    `.page${attrSelector(DATA_PAGE_NUMBER, String(pageNumberFrom1))}`
  );
};

export class CacheSvg {
  public constructor(public pdfWebView: ViewerController) {}

  private cache: Record<number, SvgBound> = {};

  public findByRect(rect: Rectangle) {
    return this.findByPage(rect.pageNumber);
  }

  public findByPage(pageNumber: number) {
    if (!this.cache[pageNumber]) {
      const svg = getDivByPageNumberFrom1(
        pageNumber,
        this.pdfWebView.getDocumentViewer().container
      )?.querySelector<SVGSVGElement>(`svg.${PDF_ANNOTATIONLAYER}`);
      const g = svg?.querySelector<SVGGElement>(
        `.${PDF_ANNOTATIONLAYER_GROUP_NOTES}`
      );

      if (!svg || !g) {
        return null;
      }

      this.setCache(pageNumber, svg, g);
    }

    return this.cache[pageNumber];
  }

  private setCache(pageNumber: number, svg: SVGSVGElement, g: SVGGElement) {
    const bound = svg.getBoundingClientRect();
    this.cache[pageNumber] = { svg, g, bound };
  }
}

export const splitAnnotation = (annotation: SelectAnnotation) => {
  const list: SelectAnnotation[] = [];

  annotation.rectangles.forEach((rect) => {
    let index = list.findIndex((anno) => anno.pageNumber === rect.pageNumber);
    if (index === -1) {
      list.push({
        ...annotation,
        pageNumber: rect.pageNumber,
        rectangles: [],
      });
      index += list.length;
    }

    list[index].rectangles.push(rect);
  });

  return list;
};
