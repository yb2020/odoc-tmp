export function elementOffsetTop(el: HTMLElement, target?: HTMLElement) {
  let elem: null | HTMLElement = el;
  let offsetTop = 0;

  do {
    if (!isNaN(el.offsetTop)) {
      offsetTop += elem.offsetTop;
    }
    elem = elem?.offsetParent as HTMLElement;
  } while (
    elem &&
    elem.parentElement !== target &&
    elem.parentElement !== document.scrollingElement
  );

  return offsetTop;
}
