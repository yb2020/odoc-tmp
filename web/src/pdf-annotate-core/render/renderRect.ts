//import { AnnotateTag } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { AnnotateTag } from 'go-sea-proto/gen/ts/common/AnnotateTag';
import { ToolBarType } from '../constants';
import setAttributes from '../utils/setAttributes';
import normalizeColor from '../utils/normalizeColor';

/**
 * Create SVGRectElements from an annotation definition.
 * This is used for anntations of type rect.
 *
 * @param {Object} a The annotation definition
 * @return {SVGGElement|SVGRectElement} A group of all rects to be rendered
 */

interface RectOptions {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface RectAnnotation {
  color: string;
  strokeWidth: number;
  opacity: number;
  rectangles: RectOptions[];
  type: ToolBarType.rect;
  fill: string;
  idea: string;
  tags: AnnotateTag[];
}

export default function renderRect(options: RectAnnotation) {
  const a = options.rectangles[0];

  const rect = createRect(a);

  setAttributes(rect, {
    ...options,
    stroke: normalizeColor(options.color || '#f00'),
    strokeWidth: options.strokeWidth || 1,
    fill: options.fill || 'none',
  });

  return rect;
}

function createRect(r: RectOptions) {
  const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');

  setAttributes(rect, {
    x: r.x,
    y: r.y,
    width: r.width,
    height: r.height,
  });

  return rect;
}

export function repaintRect(rect: SVGRectElement, attributes: RectOptions) {
  setAttributes(rect, {
    width: attributes.width,
    height: attributes.height,
  });
}
