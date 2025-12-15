import { WebDrawItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { ToolType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { createSimplePath, createPathList, setLineCapLineJoin } from './renderHandwrite';
import { PDF_ANNOTATE_SHAPE_DIV, attrSelector } from '../constants';
import { createDeleteButton } from './shapeCommon';

export function getModelByPoint (
  canvas: HTMLCanvasElement,
  x: number,
  y: number,
  modelList: WebDrawItem[],
  ipadUnit = 1
): WebDrawItem | null {
  const context = canvas.getContext('2d') as CanvasRenderingContext2D;
  setLineCapLineJoin(context);

  for (let i = 0; i < modelList.length; i += 1) {
    const model = modelList[i];

    if (!model.points.length) {
      continue;
    }

    if (model.tooltype === ToolType.IDEADrawToolTypeEraser) {
      continue;
    }

    if (!isPointInModel(model, false)) {
      continue;
    }

    if (model.erase) {
      const path = createSimplePath(model.erase.points, 1);
      context.lineWidth = model.erase.lineWidth * ipadUnit;
      if (context.isPointInStroke(path, x, y)) {
        continue;
      }
    }

    let cleared = false;

    for (let j = i + 1; j < modelList.length; j += 1) {
      const afterModel = modelList[j];

      if (afterModel.tooltype !== ToolType.IDEADrawToolTypeEraser) {
        continue;
      }

      if (isPointInModel(afterModel, true)) {
        cleared = true;
        break;
      }
    }

    if (!cleared) {
      return model;
    }
  }

  return null;

  function isPointInModel (model: WebDrawItem, isErase: boolean): boolean {
    const list = createPathList(model, 1);
    return list.some(([path, lineWidth]) => {
      const width = lineWidth * ipadUnit;
      context.lineWidth = isErase
        ? width
        : Math.max(width, 1) * 2;
      return context.isPointInStroke(path, x, y);
    });
  }
}

export interface ModelRect {
  x: number;
  y: number;
  height: number;
  width: number;
}

export const DATA_OLD_HANDWRITE_FRAME = 'data-old-handwrite-frame';

export function clearModelFrame () {
  Array
    .from(document
      .getElementsByClassName(DATA_OLD_HANDWRITE_FRAME))
    .forEach(element => {
      element.remove();
    });
}

// 新图形：fb边框，div删除
// 旧图形：div边框，div删除
// 文本：div展示，fb边框，textarea编辑，div删除

export function renderOldFrame (
  pageNumber: number,
  scale: number,
  model: WebDrawItem,
  onDelete: () => void
) {
  let left = Infinity;
  let right = -Infinity;
  let top = Infinity;
  let bottom = -Infinity;

  model.points.forEach((point) => {
    if (point.x < left) {
      left = point.x;
    }

    if (point.x > right) {
      right = point.x;
    }

    if (point.y < top) {
      top = point.y;
    }

    if (point.y > bottom) {
      bottom = point.y;
    }
  });

  left -= model.lineWidth / 2;
  right += model.lineWidth / 2;
  top -= model.lineWidth / 2;
  bottom += model.lineWidth / 2;

  const modelRect = {
    x: left,
    y: top,
    width: right - left,
    height: bottom - top,
  };

  const frame = renderFrameCommon(pageNumber, scale, model.id, modelRect, modelRect.height, onDelete)
  frame.style.border = '1px dashed red';
  return frame;
}

function renderFrameCommon (
  pageNumber: number,
  scale: number,
  modelId: WebDrawItem['id'],
  modelRect: ModelRect,
  height: number,
  onDelete: () => void
) {
  const frame = document.createElement('div');
  frame.classList.add(DATA_OLD_HANDWRITE_FRAME);
  frame.setAttribute(DATA_OLD_HANDWRITE_FRAME, modelId);
  frame.style.height = `${height}px`;

  const deleteButton = createDeleteButton(scale);
  deleteButton.addEventListener('click', onDelete);
  Object.assign(deleteButton.style, {
    top: '0',
    left: '100%',
    transformOrigin: 'center',
    transform: 'translate(-50%, -50%) ' + deleteButton.style.transform,
  } as CSSStyleDeclaration);

  frame.appendChild(deleteButton);

  Object.assign(frame.style, {
    top: `${modelRect.y}px`,
    left: `${modelRect.x}px`,
    width: `${modelRect.width}px`,
  } as CSSStyleDeclaration);

  frame.addEventListener('click', (event) => {
    console.log('click frame todo', event);
    event.stopPropagation();
  });

  frame.addEventListener('mouesdown', (event) => {
    event.stopPropagation();
  });

  const handtextDiv = document
    .body
    .querySelector(
      attrSelector(
        PDF_ANNOTATE_SHAPE_DIV,
        String(pageNumber)
      )
    ) as HTMLDivElement;

  handtextDiv?.children[0].appendChild(frame);

  return frame;
}
