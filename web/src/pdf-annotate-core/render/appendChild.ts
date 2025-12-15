import {
  PDF_ANNOTATE_ID,
  PDF_ANNOTATE_TYPE,
  PDF_ANNOTATE_USERID,
  PDF_ANNOTATE_VIEWPORT,
  ToolBarType,
} from '../constants';
import renderPath, { PathAnnotation } from './renderPath';
// import renderPoint from './renderPoint';
import renderSelect, { SelectAnnotation } from './renderSelect';
import renderText, { TextAnnotation } from './renderText';
import renderTheLine, { LineAnnotation } from './renderTheLine';
import renderRect, { RectAnnotation } from './renderRect';
import renderCircle, { CircleAnnotation } from './renderCircle';
import { setDisplay } from '../UI/utils';

const isFirefox = /firefox/i.test(navigator.userAgent);

export interface Viewport {
  rotation: number;
  width: number;
  height: number;
  scale: number;
}

/**
 * Get the x/y translation to be used for transforming the annotations
 * based on the rotation of the viewport.
 *
 * @param {Object} viewport The viewport data from the page
 * @return {Object}
 */
function getTranslation(viewport: Viewport) {
  let x;
  let y;

  // Modulus 360 on the rotation so that we only
  // have to worry about four possible values.
  switch (viewport.rotation % 360) {
    case 0:
      x = y = 0;
      break;
    case 90:
      x = 0;
      y = (viewport.width / viewport.scale) * -1;
      break;
    case 180:
      x = (viewport.width / viewport.scale) * -1;
      y = (viewport.height / viewport.scale) * -1;
      break;
    case 270:
      x = (viewport.height / viewport.scale) * -1;
      y = 0;
      break;
  }

  return { x, y };
}

/**
 * Transform the rotation and scale of a node using SVG's native transform attribute.
 *
 * @param {Node} node The node to be transformed
 * @param {Object} viewport The page's viewport data
 * @return {Node}
 */
function transform(node: SVGElement, viewport: Viewport) {
  const trans = getTranslation(viewport);

  // Let SVG natively transform the element
  node.setAttribute(
    'transform',
    `scale(${viewport.scale}) rotate(${viewport.rotation}) translate(${trans.x}, ${trans.y})`
  );

  // Manually adjust x/y for nested SVG nodes
  if (!isFirefox && node instanceof SVGSVGElement) {
    node.setAttribute(
      'x',
      parseInt(node.getAttribute('x')!, 10) * viewport.scale + ''
    );
    node.setAttribute(
      'y',
      parseInt(node.getAttribute('y')!, 10) * viewport.scale + ''
    );

    const x = parseInt(node.getAttribute('x')!, 10);
    const y = parseInt(node.getAttribute('y')!, 10);
    const width = parseInt(node.getAttribute('width')!, 10);
    const height = parseInt(node.getAttribute('height')!, 10);
    const path = node.querySelector('path')!;

    const svg: SVGAElement = path.parentNode as SVGAElement;

    // Scale width/height
    [node, svg, path, node.querySelector('rect')].forEach((n) => {
      if (!n) {
        return;
      }

      n.setAttribute(
        'width',
        parseInt(n.getAttribute('width')!, 10) * viewport.scale + ''
      );
      n.setAttribute(
        'height',
        parseInt(n.getAttribute('height')!, 10) * viewport.scale + ''
      );
    });

    // Transform path but keep scale at 100% since it will be handled natively
    transform(path, {
      ...viewport,
      scale: 1,
    });

    switch (viewport.rotation % 360) {
      case 90:
        node.setAttribute('x', viewport.width - y - width + '');
        node.setAttribute('y', x + '');
        svg.setAttribute('x', '1');
        svg.setAttribute('y', '0');
        break;
      case 180:
        node.setAttribute('x', viewport.width - x - width + '');
        node.setAttribute('y', viewport.height - y - height + '');
        svg.setAttribute('y', '2');
        break;
      case 270:
        node.setAttribute('x', y + '');
        node.setAttribute('y', viewport.height - x - height + '');
        svg.setAttribute('x', '-1');
        svg.setAttribute('y', '0');
        break;
    }
  }
}

/**
 * Append an annotation as a child of an SVG.
 *
 * @param {SVGElement} svg The SVG element to append the annotation to
 * @param {Object} annotation The annotation definition to render and append
 * @param {Object} viewport The page's viewport data
 * @return {SVGElement} A node that was created and appended by this function
 */

export type Annotation = (
  | SelectAnnotation
  | TextAnnotation
  | PathAnnotation
  | LineAnnotation
  | RectAnnotation
  | CircleAnnotation
) & {
  uuid?: string;
  documentId?: string;
  pageNumber?: number;
  score?: number;
  commentatorInfoView?: {
    nickName: string;
    userId: string;
  };
};

export function appendAnnotation(
  svg: SVGGElement,
  annotation: Annotation,
  visible: boolean,
  viewport: Viewport = JSON.parse(svg.parentElement?.getAttribute(PDF_ANNOTATE_VIEWPORT) || '{}')
) {
  let child: undefined | SVGPathElement | SVGRectElement | SVGGElement | SVGCircleElement;

  switch (annotation.type) {
    // case 'point':
    //   child = renderPoint(annotation);
    // break;
    case ToolBarType.Text:
      child = renderText(annotation);
      break;
    case ToolBarType.Draw:
      child = renderPath(annotation);
      break;

    case ToolBarType.Line:
    case ToolBarType.Arrow:
      child = renderTheLine(annotation);
      break;
    case ToolBarType.rect:
      child = renderRect(annotation);
      break;

    case ToolBarType.Circle:
      child = renderCircle(annotation);
      break;

    case ToolBarType.Highlight:
    case ToolBarType.AIHighlight:
    case ToolBarType.Vocabulary:
    case ToolBarType.Underline:
    case ToolBarType.select:
    case ToolBarType.hot:
      child = renderSelect(annotation);
      break;
  }

  // If no type was provided for an annotation it will result in node being null.
  // Skip appending/transforming if node doesn't exist.
  if (child) {
    // Set attributes
    child.setAttribute(PDF_ANNOTATE_ID, annotation.uuid || '');
    child.setAttribute(PDF_ANNOTATE_TYPE, annotation.type + '');
    child.setAttribute('data-pdf-annotate-score', annotation.score + '');
    child.setAttribute(
      PDF_ANNOTATE_USERID,
      annotation.commentatorInfoView?.userId || ''
    );

    child.setAttribute('aria-hidden', 'true');
    setDisplay(child, visible);
    transform(child, viewport);
    svg.appendChild(child);
  }

  return child;
}
