//import { AnnotateTag } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { AnnotateTag } from 'go-sea-proto/gen/ts/common/AnnotateTag';
import setAttributes from '../utils/setAttributes';
import { ToolBarType } from '../constants';

/**
 * Create SVGRectElements from an annotation definition.
 * This is used for anntations of type `area` and `highlight`.
 *
 * @param {Object} a The annotation definition
 * @return {SVGGElement|SVGRectElement} A group of all rects to be rendered
 */
export interface Rectangle {
  height: number;
  width: number;
  x: number;
  y: number;
  pageNumber: number;
}
export interface SelectAnnotation {
  type:
    | ToolBarType.Underline
    | ToolBarType.select
    | ToolBarType.Highlight
    | ToolBarType.AIHighlight
    | ToolBarType.hot
    | ToolBarType.Vocabulary;
  rectangles: Array<Rectangle>;
  rectRaw?: boolean;
  rectStr: string;
  idea: string;
  color: string;
  pageNumber: number;
  tags: AnnotateTag[];
}

export const renderUnderline = (a: SelectAnnotation) => {
  const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');

  setAttributes(group, {
    cursor: 'pointer',
    pointerEvents: 'all',
    fill: a.color,
    ...a,
  });

  a.rectangles.forEach((r) => {
    group.appendChild(createTopline(r));
    group.appendChild(createUnderline(r));
  });

  return group;
};

export const renderHighlight = (a: SelectAnnotation) => {
  const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');

  setAttributes(group, {
    fillOpacity: 0.3,
    cursor: [ToolBarType.Vocabulary].includes(a.type) ? 'text' : 'pointer',
    pointerEvents: 'all',
    fill: a.color,
    ...a,
  });

  a.rectangles.forEach((r) => {
    group.appendChild(createRect(r));
  });

  return group;
};

export const renderComment = (a: SelectAnnotation) => {
  const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');

  if (!a.rectangles || a.rectangles.length === 0) {
    return group;
  }

  setAttributes(group, {
    fillOpacity: 0.3,
    cursor: 'pointer',
    pointerEvents: 'all',
    fill: a.color,
    ...a,
  });

  a.rectangles
    .filter((item) => item.pageNumber === a.pageNumber)
    .forEach((r) => {
      group.appendChild(createRect(r));
    });

  return group;
};

export default function renderSelect(a: SelectAnnotation) {
  switch (a.type) {
    case ToolBarType.hot:
    case ToolBarType.Underline:
      return renderUnderline(a);

    case ToolBarType.Highlight:
    case ToolBarType.AIHighlight:
    case ToolBarType.Vocabulary:
      return renderHighlight(a);

    case ToolBarType.select:
      return renderComment(a);
  }
}

function createRect(r: Rectangle) {
  const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');

  setAttributes(rect, {
    x: r.x,
    y: r.y,
    width: r.width,
    height: r.height,
  });

  return rect;
}

function createUnderline(r: Rectangle) {
  const rect = document.createElementNS('http://www.w3.org/2000/svg', 'line');

  setAttributes(rect, {
    x1: r.x,
    y1: r.y + r.height,
    x2: r.x + r.width,
    y2: r.y + r.height,
    stroke: '#A8AFBA',
    strokeDasharray: '1 1',
  });

  return rect;
}

function createTopline(r: Rectangle) {
  const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');

  setAttributes(rect, {
    x: r.x,
    y: r.y,
    width: r.width,
    height: 1,
    fill: 'none',
  });

  return rect;
}
