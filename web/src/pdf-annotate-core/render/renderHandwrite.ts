// import { WebDrawItem, WebDrawV2 } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { WebDrawItem, WebDrawV2 } from 'go-sea-proto/gen/ts/note/web';
// import { 
//   ToolType,
//   DrawShapeType,
//   BasePoint,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { ToolType, DrawShapeType } from 'go-sea-proto/gen/ts/common/DrawShapeType';
import { BasePoint } from 'go-sea-proto/gen/ts/common/Point';
import { PDFPageView } from '@idea/pdfjs-dist/web/pdf_viewer';
import { PDFJSAnnotate } from '..';

/** lineWidth、fontSize、height、width单位和xy的单位不一致 */
const IPAD_PDF_100 = 950;
export const getIpadUnit = (instance: PDFJSAnnotate, pageNumberFrom0 = 0) => {
  const pdfViewer = instance.pdfWebview.getDocumentViewer().getPdfViewer();
  const pageView: PDFPageView = pdfViewer.getPageView(pageNumberFrom0);
  const { width, scale } = pageView.viewport;
  return width / scale / IPAD_PDF_100;
};

export function createPathList(model: WebDrawItem, scale: number) {
  const pathAndLineWidthList: [Path2D, number][] = [];

  switch (model.drawShapeType) {
    case DrawShapeType.IDEADrawShapeTypeCircle:
    case DrawShapeType.IDEADrawShapeTypeEllipse:
      const circlePath = renderCircle();
      pushSinglePath(circlePath);
      break;

    case DrawShapeType.IDEADrawShapeTypeArrow:
      const arrowPath = renderArrow();
      pushSinglePath(arrowPath);
      break;

    case DrawShapeType.IDEADrawShapeTypePentagon:
      const pentagonPath = renderPentagon();
      pushSinglePath(pentagonPath);
      break;

    default:
      const old = model
        .points
        .some(({ size }) => !size?.height && !size?.width)

        || (
          model.lineWidth !== 0
          && !sessionStorage.getItem('debugNewHandwrite')
        );

      if (old) {
        const oldNormalPath = createSimplePath(model.points, scale);
        pushSinglePath(oldNormalPath);
      } else {
        renderNewNormal();
      }
  }

  return pathAndLineWidthList;

  function createSinglePath() {
    const path = new Path2D();
    const [point] = model.points;
    path.moveTo(point.x * scale, point.y * scale);
    return path;
  }

  function pushSinglePath(path: Path2D | null) {
    if (path) {
      pathAndLineWidthList.push([path, model.lineWidth * scale]);
    }
  }

  function renderCircle() {
    const path = createSinglePath();

    for (let index = 1; index < model.points.length; ) {
      const cp1 = model.points[index++];
      const cp2 = model.points[index++];
      const end = model.points[index++];

      if (!cp1 || !cp2 || !end) {
        console.error('points illegal', model.points);
        return null;
      }

      path.bezierCurveTo(
        cp1.x * scale,
        cp1.y * scale,
        cp2.x * scale,
        cp2.y * scale,
        end.x * scale,
        end.y * scale
      );
    }

    return path;
  }

  function renderArrow() {
    const path = createSinglePath();

    const [, cp, end, ...rest] = model.points;
    if (!cp || !end) {
      console.error('points illegal', model.points);
      return null;
    }

    const cpx = cp.x * scale;
    const cpy = cp.y * scale;

    path.bezierCurveTo(cpx, cpy, cpx, cpy, end.x * scale, end.y * scale);

    setPathLine(path, rest, scale);

    return path;
  }

  function renderPentagon() {
    const path = createSinglePath();

    const [point, ...pointList] = model.points;
    pointList.push(point);

    setPathLine(path, pointList, scale);

    return path;
  }

  function renderNewNormal() {
    let [start, end] = model.points;
    let index = 1;
    while(end) {
      const path = new Path2D();
      path.moveTo(start.x * scale, start.y * scale);
      path.lineTo(end.x * scale, end.y * scale);

      let lineWidth = (
        heightWidthToLineWidth(
          start.size?.height as number,
          start.size?.width as number
        ) +
        heightWidthToLineWidth(
          end.size?.height as number,
          end.size?.width as number
        )
      ) / 2 * scale;

      // if (typeof start.force === 'number' && typeof end.force === 'number') {
      //   lineWidth *= (start.force + end.force) / 2;
      // }

      // if (typeof start.azimuth === 'number' && typeof end.azimuth === 'number') {
      //   lineWidth *= (start.azimuth + end.azimuth) / 2;
      // }

      if (typeof start.altitude === 'number' && typeof end.altitude === 'number') {
        // lineWidth *= (start.altitude + end.altitude) / Math.PI;
        // lineWidth *= (Math.PI - (start.altitude + end.altitude)) / Math.PI;
        lineWidth /= 2;
      }

      pathAndLineWidthList.push([path, lineWidth]);

      index += 1;
      start = end;
      end = model.points[index];
    }

    function heightWidthToLineWidth(height: number, width: number) {
      return Math.sqrt(height * height + width * width);
    }
  }
}

export function createSimplePath(pointList: BasePoint[], scale: number) {
  const path = new Path2D();
  const [point, ...rest] = pointList;

  path.moveTo(point.x * scale, point.y * scale);
  setPathLine(path, rest, scale);

  return path;
}

function setPathLine(path: Path2D, pointList: BasePoint[], scale: number) {
  for (let index = 0; index < pointList.length; index += 1) {
    const { x, y } = pointList[index];
    path.lineTo(x * scale, y * scale);
  }
}

export function setLineCapLineJoin(context: CanvasRenderingContext2D) {
  const debugLineCap = 'debugLineCap';
  const debugLineJoin = 'debugLineJoin';
  if (debugLineCap in sessionStorage) {
    if (sessionStorage.getItem(debugLineCap)) {
      context.lineCap = sessionStorage.getItem(debugLineCap) as CanvasLineCap;
    }
  } else {
    context.lineCap = 'round';
  }
 
  if (debugLineJoin in sessionStorage) {
    if (sessionStorage.getItem(debugLineJoin)) {
      context.lineJoin = sessionStorage.getItem(debugLineJoin) as CanvasLineJoin;
    }
  } else {
    context.lineJoin = 'round';
  }
}

function renderHandwrite(
  canvas: HTMLCanvasElement,
  scale: number,
  ipadUnit: number,
  model: WebDrawItem
) {
  if (!model.points?.length) {
    console.warn('no points');
    return;
  }

  const context = canvas.getContext('2d') as CanvasRenderingContext2D;

  if (!model.erase) {
    render();
    return;
  }

  const { points, lineWidth } = model.erase;

  const path = createSimplePath(points, scale);
  path.lineTo(points[0].x * scale, points[0].y * scale);
  setLineCapLineJoin(context);
  context.lineWidth = lineWidth * scale * ipadUnit;
  context.save();
  context.clip(path);
  render();
  context.restore();

  function render() {
    switch (model.tooltype) {
      case ToolType.IDEADrawToolTypePen:
      case ToolType.IDEADrawToolTypeMarker:
        withDebug(() => {
          context.strokeStyle = getColor();
          renderPath();
        });
        break;

      case ToolType.IDEADrawToolTypeEraser:
        withDebug(() => {
          context.strokeStyle = '#000000';
          context.globalCompositeOperation = 'destination-out';
          renderPath();
          context.globalCompositeOperation = 'source-over';
        });
        break;
  
      default:
    }
  }

  function renderPath() {
    setLineCapLineJoin(context);

    const list = createPathList(model, scale);
    list.forEach(([path, lineWidth]) => {
      context.lineWidth = lineWidth * ipadUnit;
      context.stroke(path);
    });
  }

  function getColor() {
    if (model.lineHexColor.length !== 7) {
      console.error('手写数据颜色长度异常');
    }

    const hexAlpha =
      model.lineAlpha === 1
        ? ''
        : Math.round(256 * model.lineAlpha).toString(16);

    return model.lineHexColor + hexAlpha;
  }

  function debugPrintPoints() {
    context.fillStyle = getColor();
    for (const point of model.points) {
      context.beginPath();
      context.arc(point.x * scale, point.y * scale, scale, 0, Math.PI * 2);
      context.fill();
      context.closePath();
    }
  }

  function withDebug(renderFunc: () => void) {
    const key = `debugHand${model.tooltype}`
    const debugType = sessionStorage.getItem(key) as HandwriteDebugType;
    if (debugType === HandwriteDebugType.single) {
      debugPrintPoints();
    } else if (debugType === HandwriteDebugType.both) {
      renderFunc();
      debugPrintPoints();
    } else {
      renderFunc();
    }
  }
}

enum HandwriteDebugType {
  single = '1',
  both = '2',
}

const getScale = (instance: PDFJSAnnotate, pageNumberPlus1: number) => {
  const pdfViewer = instance.pdfWebview.getDocumentViewer().getPdfViewer();
  const pageView: PDFPageView = pdfViewer.getPageView(pageNumberPlus1 - 1);
  return pageView.viewport.scale; 
}

export const renderSinglePageHandwrite = (
  pageNumberPlus1: number, 
  list: WebDrawV2[],
  instance: PDFJSAnnotate,
  clear = false,
  scale = getScale(instance, pageNumberPlus1),
  canvas = instance.canvasElements.get(pageNumberPlus1 - 1) as HTMLCanvasElement
) => {
  if (clear) {
    const context = canvas.getContext('2d') as CanvasRenderingContext2D;
    context.clearRect(0, 0, canvas.width, canvas.height);
  }

  const ipadUnit = getIpadUnit(instance, pageNumberPlus1 - 1);

  list.forEach((webDrawV2) => {
    webDrawV2.webDrawItems.forEach((model) => {
      renderHandwrite(canvas, scale, ipadUnit, model);
    });
  });
};
