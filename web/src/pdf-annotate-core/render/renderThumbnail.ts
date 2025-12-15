/*
 * Created Date: March 24th 2022, 10:59:45 am
 * Author: zhoupengcheng
 * -----
 * Last Modified: March 28th 2022, 11:29:39 am
 */
export const renderThumbnail = (svg: SVGElement) =>
  new Promise<string>((resolve) => {
    const s = new XMLSerializer().serializeToString(svg);
    const _s = unescape(encodeURIComponent(s));

    const encodedData = window.btoa(_s);

    const rect = svg.getBoundingClientRect();

    const img = new Image();
    img.crossOrigin = 'Anonymous';

    const pbx = document.createElement('img');

    pbx.style.width = rect.width + 'px';
    pbx.style.height = rect.height + 'px';

    pbx.src = 'data:image/svg+xml;base64,' + encodedData;

    pbx.onload = () => {
      const canvas = document.createElement('canvas');

      const context = canvas.getContext('2d')!;

      canvas.width = rect.width;

      canvas.height = rect.height;

      canvas.style.position = 'absolute';

      context.drawImage(pbx, 0, 0);

      resolve(canvas.toDataURL('image/png', 1));
    };
  });
