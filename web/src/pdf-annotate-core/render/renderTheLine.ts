import { ToolBarType } from '../constants';
import setAttributes from '../utils/setAttributes';
import normalizeColor from '../utils/normalizeColor';
import uuid from '../utils/uuid';

/**
 * Create SVGLineElements from an annotation definition.
 * This is used for anntations of type `strikeout`.
 *
 * @param {Object} a The annotation definition
 * @return {SVGGElement} A group of all lines to be rendered
 */
export interface LineAnnotation {
  color: string;
  width: number;
  opacity: number;
  rectangles: { x1: number; y1: number; x2: number; y2: number }[];
  type: ToolBarType.Arrow | ToolBarType.Line;
}
export default function renderLine(options: LineAnnotation) {
  const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
  setAttributes(group, {
    stroke: normalizeColor(options.color || '#f00'),
    strokeWidth: options.width || 1,
    strokeOpacity: options.opacity,
  });

  const r = options.rectangles[0];

  const line = document.createElementNS('http://www.w3.org/2000/svg', 'line');

  setAttributes(line, {
    x1: r.x1,
    y1: r.y1,
    x2: r.x2,
    y2: r.y2,
  });

  if (options.type === ToolBarType.Arrow) {
    const { def, id } = renderArrow(options);

    group.appendChild(def);

    setAttributes(line, {
      markerEnd: `url(#${id})`,
    });
  }

  group.appendChild(line);

  return group;
}

export function repaintLine(
  line: SVGLineElement,
  attributes: { x2: number; y2: number }
) {
  const { x2, y2 } = attributes;

  setAttributes(line, {
    x2,
    y2,
  });
}

export function renderArrow(options: LineAnnotation) {
  const def = document.createElementNS('http://www.w3.org/2000/svg', 'defs');
  const marker = document.createElementNS(
    'http://www.w3.org/2000/svg',
    'marker'
  );

  const id = uuid();

  setAttributes(marker, {
    markerWidth: 5,
    markerHeight: 5,
    refX: 0.5,
    refY: 2,
    orient: 'auto',
    markerUnits: 'strokeWidth',
    stroke: normalizeColor(options.color || '#f00'),
    strokeWidth: options.width || 1,
    strokeOpacity: options.opacity,
    fill: 'transparent',
    id,
  });

  const polyline = document.createElementNS(
    'http://www.w3.org/2000/svg',
    'polyline'
  );
  setAttributes(polyline, {
    points: '0.25,0.5 3.25,2 0.25,3.5 3.25,2 -0.25,2',
    stroke: normalizeColor(options.color || '#f00'),
    strokeWidth: options.width || 1,
    strokeOpacity: options.opacity,
    style:
      'pointer-events: none; fill: none; stroke-width: 1; stroke-dasharray: 100;',
  });

  marker.appendChild(polyline);
  def.appendChild(marker);

  return { def, id };
}
