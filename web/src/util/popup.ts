export function getTextRectByAspectRadio(
  txt: string,
  rate: number,
  styles: Record<string, string>,
  minWidth = 0,
  maxWidth = window.innerWidth
) {
  const div = document.createElement('div');
  div.style.position = 'fixed';
  div.style.left = '-99999px';
  div.style.top = '-99999px';
  // div.style.aspectRatio = rate;

  const pre = document.createElement('pre');
  pre.style.whiteSpace = 'normal';
  pre.style.wordBreak = 'break-word';
  pre.innerText = txt;
  Object.entries(styles).forEach(([k, v]) => pre.style.setProperty(k, v));

  div.appendChild(pre);
  document.body.appendChild(div);

  let w = maxWidth;
  let needSmaller = false;
  let needLarge = false;
  do {
    div.style.width = `${w}px`;
    div.style.height = `${w / rate}px`;
    needLarge = div.scrollHeight > div.offsetHeight;
    needSmaller = div.offsetHeight - pre.offsetHeight > 3;
    if (needSmaller) {
      maxWidth = w;
      w = (w + minWidth) / 2;
    } else if (needLarge) {
      minWidth = w;
      w = w + (maxWidth - w) / 2;
    }
  } while ((needSmaller || needLarge) && maxWidth - minWidth > 1);

  // document.removeChild(div);

  return Math.round(w);
}
