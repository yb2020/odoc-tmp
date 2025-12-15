import interact from 'interactjs';
import merge from 'lodash-es/merge';
import { DraggableOptions } from '@interactjs/actions/drag/plugin';

export default function enableInteractDraggable(
  target: HTMLElement | SVGElement | string,
  options: Partial<Interact.OrBoolean<DraggableOptions>> = {}
) {
  const instance = interact(target).draggable(
    merge(
      {
        ignoreFrom: '.ps__rail-y,.ps__rail-x,.js-interact-drag-ignore',
        listeners: {
          move: function dragMoveListener(event: Interact.DragEvent) {
            const target = event.target as HTMLElement;
            let x = parseFloat(target.getAttribute('data-x') || '');
            let y = parseFloat(target.getAttribute('data-y') || '');
            if (isNaN(x) || isNaN(y)) {
              const computedStyle = window.getComputedStyle(target, null);
              const transform = computedStyle.transform.substring(6);
              if (transform) {
                x = parseFloat(transform.split(',')[4]);
                y = parseFloat(transform.split(',')[5]);
              } else {
                x = 0;
                y = 0;
              }
            }
            // keep the dragged position in the data-interact-x/data-interact-y attributes
            x += event.dx;
            y += event.dy;

            // translate the element
            target.style.transform = 'translate(' + x + 'px, ' + y + 'px)';

            // update the posiion attributes
            target.setAttribute('data-x', `${x}`);
            target.setAttribute('data-y', `${y}`);
          },
        },
        inertia: true,
        modifiers: [
          interact.modifiers.restrictRect({
            restriction: 'body',
            endOnly: true,
          }),
        ],
      },
      options
    )
  );

  instance.preventDefault();

  return instance;
}
