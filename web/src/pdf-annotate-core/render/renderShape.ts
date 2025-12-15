import * as fabric from 'fabric';

import { PDFPageView } from '@idea/pdfjs-dist/web/pdf_viewer';
// import {
//   AnnotationColor,
//   ShapeAnnotation,
//   ShapeType,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { AnnotationColor, ShapeAnnotation, ShapeType } from 'go-sea-proto/gen/ts/common/ShapeAnnotation';
import { TRANSPARENT, colorMap } from '../constants/color';
import { PointOnPdf } from '@idea/pdf-annotate-viewer';
import {
  PageFcanvasMap,
  createDeleteButton,
  createToolbarDiv,
  getShapeDiv,
  getShapeScale,
  optimizeRectEllipsePosition,
  optimizeTextArrowPosition,
} from './shapeCommon';
import { numberToPx } from '../utils/setAttributes';

const STROKE_WIDTH = Number(localStorage.getItem('strokeWidth')) || 1.5;

export const shapeCommonProps = {
  fill: TRANSPARENT,
  strokeWidth: STROKE_WIDTH,
  strokeUniform: true,
  perPixelTargetFind: true,
  transparentCorners: false,
  borderColor: 'transparent',
  cornerColor: 'white',
  cornerStyle: 'circle' as const,
  cornerSize: Number(localStorage.getItem('cornerSize')) || 10,
};

export const newShape = (item: ShapeAnnotation, scale: number) => {
  const commonProps = {
    ...shapeCommonProps,
    stroke: colorMap[item.strokeColor],
    cornerStrokeColor: colorMap[item.strokeColor],
  };

  commonProps.strokeWidth *= scale;

  if (item.type === ShapeType.rectangle) {
    return new fabric.Rect({
      left: item.x * scale,
      top: item.y * scale,
      height: item.height * scale,
      width: item.width * scale,
      ...commonProps,
    });
  } else if (item.type === ShapeType.ellipse) {
    return new fabric.Ellipse({
      top: (item.y - item.radiusY) * scale,
      left: (item.x - item.radiusX) * scale,
      rx: item.radiusX * scale,
      ry: item.radiusY * scale,
      ...commonProps,
      borderColor: commonProps.stroke,
    });
  } else if (item.type === ShapeType.arrow) {
    return new fabric.Line(
      [item.x * scale, item.y * scale, item.endX * scale, item.endY * scale],
      {
        ...commonProps,
        borderColor: commonProps.stroke,
      }
    );
  } else {
    throw '不支持的图形类型';
  }
};

export const getArrowheadTransform = (arrowXY: ArrrowXY) => {
  return `translate(-50%, -50%) rotate(${Math.atan2(
    arrowXY.endY - arrowXY.y,
    arrowXY.endX - arrowXY.x
  )}rad)`;
};

export const getArrowheadStyle = (line: fabric.Line, scale: number) => {
  const arrowXY = convertLineToArrow(line);
  return {
    left: arrowXY.endX / scale + 'px',
    top: arrowXY.endY / scale + 'px',
    transform: getArrowheadTransform(arrowXY),
  };
};

export const ARROWHEAD_CLASSNAME = 'shape-arrowhead';

export const createArrowhead = (
  id: string,
  line: fabric.Line,
  scale: number
) => {
  const arrowhead = document.createElement('div');
  arrowhead.setAttribute(`data-${ARROWHEAD_CLASSNAME}`, id);
  arrowhead.classList.add(ARROWHEAD_CLASSNAME);
  arrowhead.innerHTML = '<div></div>';

  const arrowheadStyle = getArrowheadStyle(line, scale);
  Object.assign(arrowhead.style, arrowheadStyle);
  (arrowhead.children[0] as HTMLDivElement).style.backgroundColor =
    line.stroke as string;

  const sync = () => {
    syncArrow(arrowhead, line, scale);
  };

  line.on('moving', sync);
  line.on('scaling', sync);

  return arrowhead;
};

export const getArrowHead = (id: string) => {
  // @TODO 替换document
  return document.querySelector(
    `.${ARROWHEAD_CLASSNAME}[data-${ARROWHEAD_CLASSNAME}="${id}"]`
  ) as HTMLDivElement;
};

export interface ShapeCallback {
  onUpdate(sa: ShapeAnnotation): void;
  onDelete(id: string): void;
  createToolbar(
    div: HTMLDivElement,
    sa: ShapeAnnotation,
    onChange: (color: AnnotationColor) => void,
  ): void;
}

export const initShape = (
  item: ShapeAnnotation,
  fobj: fabric.Rect | fabric.Ellipse | fabric.Line,
  pageView: PDFPageView,
  callback: ShapeCallback
) => {
  const { pageNumber } = item;
  const fcanvas = PageFcanvasMap.get(pageNumber) as fabric.Canvas;

  const container = pageView.div.closest<HTMLDivElement>(
    '.rp-pdf-viewer-viewer'
  );
  const renderFrame = () => {
    const shapeBound = shapeDiv.getBoundingClientRect();
    const fbound = getMrect();

    const y = shapeBound.top - container!.getBoundingClientRect().top;
    const x = shapeBound.left - container!.getBoundingClientRect().left;

    {
      const top = y + fbound.top + fbound.height + 20;
      const left = x + fbound.left;
      Object.assign(toolbarDiv.style, {
        top: numberToPx(top),
        left: numberToPx(left),
      });
    }

    {
      const top = y + fbound.top - 10;
      const left = x + fbound.left + fbound.width - 10;
      Object.assign(deleteButton.style, {
        top: numberToPx(top),
        left: numberToPx(left),
      });
    }

    if (!toolbarDiv.parentElement) {
      container?.appendChild(toolbarDiv);
    }

    if (!deleteButton.parentElement) {
      container?.appendChild(deleteButton);
    }

    function getMrect() {
      const { left, top } = fobj;
      const width =
        fobj instanceof fabric.Ellipse
          ? (fobj as fabric.Ellipse).getRx() * 2
          : fobj.getScaledWidth();
      const height =
        fobj instanceof fabric.Ellipse
          ? (fobj as fabric.Ellipse).getRy() * 2
          : fobj.getScaledHeight();

      return {
        left,
        top,
        width,
        height,
      };
    }
  };

  const shapeDiv = getShapeDiv(pageNumber, container!);
  const toolbarDiv = createToolbarDiv();
  const deleteButton = createDeleteButton();
  callback.createToolbar(toolbarDiv, item, (color) => {
    fobj.set({
      stroke: colorMap[color],
      cornerStrokeColor: colorMap[color],
    });
    fcanvas.requestRenderAll();

    if (fobj instanceof fabric.Line) {
      const arrowhead = getArrowHead(
        item.shapeId || item.uuid
      ) as HTMLDivElement;
      (arrowhead.children[0] as HTMLDivElement).style.backgroundColor =
        colorMap[color];
    }
  });

  deleteButton.addEventListener('click', () => {
    callback.onDelete(item.shapeId || item.uuid);
    fcanvas.remove(fobj);
    toolbarDiv.remove();
    deleteButton.remove();

    if (fobj instanceof fabric.Line) {
      const arrowhead = getArrowHead(item.shapeId || item.uuid);
      (arrowhead as HTMLDivElement).remove();
    }
  });

  fobj.on('modified', () => {
    renderFrame();

    const scale = pageView.viewport.scale;

    if (item.type === ShapeType.rectangle) {
      item.x = fobj.left;
      item.y = fobj.top;
      item.width = fobj.getScaledWidth();
      item.height = fobj.getScaledHeight();

      item.width /= scale;
      item.height /= scale;
    } else if (fobj instanceof fabric.Ellipse) {
      item.x = fobj.left + (fobj as fabric.Ellipse).getRx();
      item.y = fobj.top + (fobj as fabric.Ellipse).getRy();
      item.radiusX = (fobj as fabric.Ellipse).getRx();
      item.radiusY = (fobj as fabric.Ellipse).getRy();

      item.radiusX /= scale;
      item.radiusY /= scale;
    } else if (fobj instanceof fabric.Line) {
      const arrowXY = convertLineToArrow(fobj);
      Object.assign(item, arrowXY);

      item.endX /= scale;
      item.endY /= scale;
    }

    item.x /= scale;
    item.y /= scale;
    callback.onUpdate(item);
  });
  fobj.on('selected', () => {
    fobj.perPixelTargetFind = false;
    renderFrame();
    Array.from(PageFcanvasMap.values())
      .filter((fc) => fc !== fcanvas)
      .forEach((fc) => {
        fc.discardActiveObject();
      });
  });
  fobj.on('deselected', () => {
    fobj.perPixelTargetFind = true;
    toolbarDiv.remove();
    deleteButton.remove();
  });

  // let moving = false;

  const onMoving = () => {
    // fobj.set({
    //   height: fobj.getScaledHeight(),
    //   width: fobj.getScaledWidth(),
    //   scaleX: 1,
    //   scaleY: 1,
    // });

    // if (!moving) {
    //   moving = true;
    //   clearModelFrame();
    // }

    if (fobj instanceof fabric.Line) {
      optimizeTextArrowPosition(fobj, pageNumber);
    } else {
      optimizeRectEllipsePosition(fobj, pageNumber);
    }

    if (toolbarDiv?.parentElement) {
      toolbarDiv.remove();
    }

    if (deleteButton?.parentElement) {
      deleteButton.remove();
    }
  };

  fobj.on('moving', onMoving);
  fobj.on('scaling', onMoving);
  fobj.on('mousedown', (event) => {
    event.e.stopPropagation();
  });
  fobj.setControlVisible('mtr', false);
};

export const syncArrow = (
  arrowhead: HTMLDivElement,
  line: fabric.Line,
  scale: number
) => {
  const arrowheadStyle = getArrowheadStyle(line, scale);
  Object.assign(arrowhead.style, arrowheadStyle);
};

export type ArrrowXY = Pick<ShapeAnnotation, 'x' | 'y' | 'endX' | 'endY'>;

export const convertLineToArrow = (line: fabric.Line): ArrrowXY => {
  const endX = line.left + line.getScaledWidth();
  const endY = line.top + line.getScaledHeight();

  const arrow: ArrrowXY = {
    x: 0,
    y: 0,
    endX: 0,
    endY: 0,
  };

  [arrow.x, arrow.endX] =
    line.x1 < line.x2
      ? line.flipX
        ? [endX, line.left]
        : [line.left, endX]
      : line.flipX
        ? [line.left, endX]
        : [endX, line.left];

  [arrow.y, arrow.endY] =
    line.y1 < line.y2
      ? line.flipY
        ? [endY, line.top]
        : [line.top, endY]
      : line.flipY
        ? [line.top, endY]
        : [endY, line.top];

  return arrow;
};

export const createShapeAtPoint = (
  point: PointOnPdf,
  type: ShapeType | 'text',
  color: AnnotationColor,
  scale: number,
  container: HTMLElement
) => {
  const fcanvas = PageFcanvasMap.get(point.cur + 1) as fabric.Canvas;

  const strokeCommonProps = {
    stroke: colorMap[color],
    strokeWidth: shapeCommonProps.strokeWidth * scale,
    cornerStrokeColor: colorMap[color],
  };

  let shape: fabric.Rect | fabric.Ellipse | fabric.Line;

  if (type === ShapeType.ellipse) {
    shape = new fabric.Ellipse({
      left: point.left,
      top: point.top,
      rx: 0,
      ry: 0,
      ...shapeCommonProps,
      ...strokeCommonProps,
      borderColor: strokeCommonProps.stroke,
    });
  } else if (type === ShapeType.arrow) {
    shape = new fabric.Line([point.left, point.top, point.left, point.top], {
      ...shapeCommonProps,
      ...strokeCommonProps,
      borderColor: strokeCommonProps.stroke,
    });

    const arrowhead = createArrowhead(NEW_ARROW_ID, shape, scale);
    const shapeScale = getShapeScale(point.cur + 1, container);
    shapeScale.appendChild(arrowhead);
  } else {
    const rectCommonProps = {
      left: point.left,
      top: point.top,
      width: 0,
      height: 0,
    };

    if (type === ShapeType.rectangle) {
      shape = new fabric.Rect({
        ...rectCommonProps,
        ...shapeCommonProps,
        ...strokeCommonProps,
      });
    } else {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { strokeUniform, strokeWidth, ...rest } = shapeCommonProps;

      shape = new fabric.Rect({
        ...rest,
        ...rectCommonProps,
        strokeUniform: false,
        strokeWidth: 1,
        stroke: colorMap[color],
        borderColor: colorMap[color],
        cornerStrokeColor: colorMap[color],
      });
    }
  }

  fcanvas.add(shape);

  return shape;
};

export const NEW_ARROW_ID = 'new-arrow-id';

export const updateShapeWithPoint = (
  startPoint: PointOnPdf,
  currentPoint: PointOnPdf,
  shape: fabric.Rect | fabric.Ellipse | fabric.Line,
  scale: number
) => {
  const fcanvas = PageFcanvasMap.get(startPoint.cur + 1) as fabric.Canvas;
  const fwidth = fcanvas.getWidth();
  const fheight = fcanvas.getHeight();

  const optimizedCurrentPoint = { ...currentPoint };

  if (optimizedCurrentPoint.left < 0) {
    optimizedCurrentPoint.left = 0;
  } else if (optimizedCurrentPoint.left > fwidth) {
    optimizedCurrentPoint.left = fwidth;
  }

  if (optimizedCurrentPoint.top < 0) {
    optimizedCurrentPoint.top = 0;
  } else if (optimizedCurrentPoint.top > fheight) {
    optimizedCurrentPoint.top = fheight;
  }

  const { left, top, width, height } = getShapeBound(
    startPoint,
    optimizedCurrentPoint
  );

  if (shape instanceof fabric.Rect) {
    shape.set({
      left,
      top,
      width,
      height,
    });
  } else if (shape instanceof fabric.Ellipse) {
    const rx = width / 2;
    const ry = height / 2;

    shape.set({
      left,
      top,
      rx,
      ry,
    });
  } else if (shape instanceof fabric.Line) {
    shape.set({
      x2: optimizedCurrentPoint.left,
      y2: optimizedCurrentPoint.top,
    });

    const arrowhead = getArrowHead(NEW_ARROW_ID);
    syncArrow(arrowhead, shape, scale);
  }

  fcanvas.requestRenderAll();
};

export const getShapeBound = (
  startPoint: PointOnPdf,
  currentPoint: PointOnPdf
) => {
  const height = Math.abs(currentPoint.top - startPoint.top);
  const width = Math.abs(currentPoint.left - startPoint.left);
  const left = Math.min(currentPoint.left, startPoint.left);
  const top = Math.min(currentPoint.top, startPoint.top);

  return {
    left,
    top,
    width,
    height,
  };
};
