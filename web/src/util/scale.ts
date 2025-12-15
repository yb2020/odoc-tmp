import { PageViewport } from '@idea/pdfjs-dist';

export function scaleDownRect(
  rect: Record<string, any>,
  viewport: PageViewport
) {
  const result = {} as any;

  Object.keys(rect).forEach((key) => {
    if (typeof rect[key] !== 'number') {
      return;
    }

    result[key as string] = rect[key] / viewport.scale;
  });

  return result;
}
