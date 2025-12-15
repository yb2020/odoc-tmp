import EventEmitter from 'events';
import { findSVGAtPoint, findAnnotationListAtPoint } from './utils';
import { ToolBarType } from '../constants';

export const ANNOTATION_CLICK = 'annotation:click';

export const useAnnotationMouseDown = (emitter: EventEmitter) => {
  let clickNode: SVGGElement | SVGRectElement | undefined;

  return (element: SVGSVGElement | HTMLDivElement) => {
    (element as HTMLDivElement).addEventListener('pointerdown', (e) => {
      console.warn('element pointerdown');

      element.addEventListener('pointerup', onMouseUp, { once: true });
      setTimeout(() => {
        element.removeEventListener('pointerup', onMouseUp);
      }, 300);

      function onMouseUp() {
        if (!findSVGAtPoint(e.clientX, e.clientY)) {
          return;
        }

        // 默认取最后一个层级最高
        const target = findAnnotationListAtPoint(e.clientX, e.clientY)
          ?.filter((item) => {
            const type = Number(item.getAttribute('type'));
            return type !== ToolBarType.hot;
          })
          .pop();
        // Emit annotation:blur if clickNode is no longer clicked
        if (clickNode && clickNode !== target) {
          emitter.emit('annotation:blur', clickNode, e);
        }

        // Emit annotation:click if target was clicked
        if (target) {
          emitter.emit(ANNOTATION_CLICK, target, e);
        }

        clickNode = target;
      }
    });
  };
};

export const ANNOTATION_MOUSEOVER = 'annotation:mouseover';

export const useAnnotationMouseMove = (emitter: EventEmitter) => {
  let mouseOverNode: ReturnType<typeof findAnnotationListAtPoint>;

  return (e: MouseEvent) => {
    const targetList = findAnnotationListAtPoint(e.clientX, e.clientY);

    if (targetList.length === 0 && mouseOverNode?.length) {
      // Emit annotation:mouseout if target was mouseout'd
      emitter.emit('annotation:mouseout', mouseOverNode, e);
    } else {
      const outTargets = mouseOverNode?.filter((x) => !targetList.includes(x));

      if (outTargets?.length) {
        emitter.emit('annotation:mouseout', outTargets, e);
      }

      // Emit annotation:mouseover if target was mouseover'd
      emitter.emit(ANNOTATION_MOUSEOVER, targetList, e);
    }

    mouseOverNode = targetList;
  };
};
