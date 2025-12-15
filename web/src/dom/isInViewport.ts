import { debounce } from 'lodash-es';
import isElementVisible from 'element-visible';

export type Listener = (x: Element) => void;

export interface ListenerOpts {
  threshold?: number;
  listener: Listener;
}

export function isInDOMTree(element?: null | Node) {
  do {
    if (element === document) {
      return true;
    }
  } while ((element = element?.parentNode));

  return false;
}

export class ViewObserver {
  static isInited = 0;

  static THRESHOLD = 1;

  static triggered = new Set();

  static listeners = new Map<Element[], ListenerOpts>();

  static initViewportEvents() {
    if (ViewObserver.isInited) {
      return;
    }
    ViewObserver.isInited = 1;
    window.addEventListener('scroll', ViewObserver.checkAllElements);
    window.addEventListener('resize', ViewObserver.checkAllElements);
  }

  static checkAllElements = debounce(() => {
    const { listeners } = ViewObserver;
    const entries = listeners.entries();

    let item = entries.next();
    while (!item.done) {
      const [elements, { listener, threshold }] = item.value;
      ViewObserver.checkElements(elements, listener, threshold);
      item = entries.next();
    }
  }, 300);

  static checkElements = (
    elements: Element[],
    listener: (x: Element) => void,
    threshold = ViewObserver.THRESHOLD
  ) => {
    const { triggered } = ViewObserver;
    const allTriggered = elements.every((element) => {
      let hasTriggered = triggered.has(element);
      if (isInDOMTree(element) && isElementVisible(element, threshold) && !hasTriggered) {
        triggered.add(element);
        listener(element);
        hasTriggered = true;
      }
      return hasTriggered;
    });
    if (allTriggered) {
      ViewObserver.listeners.delete(elements);
    }
  };

  static watchElements(
    elements: Element[],
    listenerOpts: Listener | ListenerOpts,
    scrollEl?: Element
  ) {
    const { listeners } = ViewObserver;
    if (typeof listenerOpts === 'function') {
      listenerOpts = { listener: listenerOpts };
    }
    const { listener, threshold } = listenerOpts;

    listeners.set(elements, listenerOpts);

    ViewObserver.initViewportEvents();
    if (scrollEl) {
      scrollEl.addEventListener(
        'scroll',
        debounce(() => {
          ViewObserver.checkElements(elements, listener, threshold);
        }, 300)
      );
    }
    ViewObserver.checkElements(elements, listener, threshold);
  }
}
