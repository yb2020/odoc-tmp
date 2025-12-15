import { ToolBarType } from '../constants';
import normalizeColor from '../utils/normalizeColor';

import setAttributes from '../utils/setAttributes';

/**
 * Create SVGTextElement from an annotation definition.
 * This is used for anntations of type `textbox`.
 *
 * @param {Object} a The annotation definition
 * @return {SVGTextElement} A text to be rendered
 */
export interface TextAnnotation {
  type: ToolBarType.Text;
  content: string;
  size: string;
  x: number;
  y: number;
  color: string;
}

export default function renderText(a: TextAnnotation) {
  const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');

  if (typeof a.content !== 'string') {
    return group;
  }

  const contentArray = a.content.split('\n').filter((str) => !!str);

  let y = a.y + parseInt(a.size, 10);

  setAttributes(group, {
    x: a.x,
    y,
    fontSize: a.size,
  });

  for (const str of contentArray) {
    const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');

    setAttributes(text, {
      x: a.x,
      y,
      fill: normalizeColor(a.color || '#000'),
      fontSize: a.size,
    });

    text.innerHTML = str;

    group.appendChild(text);

    y += parseInt(a.size, 10);
  }

  return group;
}
