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

interface CircleOptions {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface CircleAnnotation {
  color: string;
  width: number;
  opacity: number;
  rectangles: CircleOptions[];
  type: ToolBarType.Circle;
}

export default function renderCircle(options: CircleAnnotation) {
  const a = options.rectangles[0];

  const circle = createCircle(a);

  setAttributes(circle, {
    stroke: normalizeColor(options.color || '#f00'),
    fill: 'none',
    strokeWidth: options.width || 1,
  });

  return circle;
}

function createCircle(r: CircleOptions) {
  const circle = document.createElementNS(
    'http://www.w3.org/2000/svg',
    'circle'
  );

  setAttributes(circle, {
    r: r.width,
    cx: r.x,
    cy: r.y,
  });

  return circle;
}

export function repaintCircle(
  circle: SVGCircleElement,
  attributes: CircleOptions
) {
  setAttributes(circle, {
    r: attributes.width,
  });
}
