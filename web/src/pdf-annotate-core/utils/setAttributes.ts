// Don't convert these attributes from camelCase to hyphenated-attributes
const BLACKLIST = [
  'viewBox',
  'markerWidth',
  'markerHeight',
  'markerUnits',
  'refX',
  'refY',
];

const keyCase = (key: string) => {
  if (BLACKLIST.includes(key)) {
    return key
  }

  return key.replace(/[A-Z]/g, (match) => '-' + match.toLowerCase());
};

/**
 * Set attributes for a node from a map
 *
 * @param {Node} node The node to set attributes on
 * @param {Object} attributes The map of key/value pairs to use for attributes
 */
export default function setAttributes(
  node: Element,
  attributes: Record<string, any>
) {
  Object.keys(attributes).forEach((key) => {
    try {
      const attr = attributes[key];

      if (typeof attr === 'object') {
        node.setAttribute(keyCase(key), JSON.stringify(attr));
      } else {
        node.setAttribute(keyCase(key), attr);
      }
    } catch (err) {
      console.error(err);
    }
  });
}

export function numberToPx(number: number) {
  return typeof number !== 'number'
    ? ''
    : number === 0
      ? '0'
      : `${number}px`;
}

export function pxToNumber(px: string) {
  return typeof px !== 'string'
    ? 0
    : px === '0'
      ? 0
      : Number(px.replace('px', ''));
}
