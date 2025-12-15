import interact from 'interactjs';
import trim from 'lodash-es/trim';
import merge from 'lodash-es/merge';
import flow from 'lodash-es/flow';
import { ResizableOptions } from '@interactjs/actions/resize/plugin';

export default function enableInteractResizable(
  target: HTMLElement | SVGElement | string,
  options: Partial<Interact.OrBoolean<ResizableOptions>> = {}
) {
  return interact(target).resizable(
    merge(
      {
        margin: 10,
        // resize from all edges and corners
        edges: { left: true, right: true, bottom: true, top: true },

        modifiers: [
          // keep the edges inside the parent
          interact.modifiers.restrictEdges({
            outer: 'body',
          }),

          // minimum size
          interact.modifiers.restrictSize({
            min: { width: 200, height: 200 },
          }),

          // interact.modifiers.snapSize({
          //   targets: [
          //     { width: 10 },
          //     // interact.snappers.grid({ width: 100, height: 100 }),
          //   ],
          // }),
        ],
        inertia: true,
      },
      options,
      {
        listeners: {
          move: flow(
            (event: Interact.ResizeEvent) => {
              if (!event.deltaRect) {
                return;
              }

              const target = event.target as HTMLElement;
              // update the element's style
              target.style.width = event.rect.width + 'px';

              if (
                (event.edges as Interact.EdgeOptions).bottom ||
                (event.edges as Interact.EdgeOptions).top
              ) {
                target.style.height = event.rect.height + 'px';
                target.style.maxHeight = event.rect.height + 'px';
              }

              const tippyTarget = target.closest(
                `[data-tippy-root]`
              ) as HTMLElement;

              if (!tippyTarget) {
                return;
              }

              let x = parseFloat(tippyTarget.getAttribute('data-x') || '');
              let y = parseFloat(tippyTarget.getAttribute('data-y') || '');

              if (isNaN(x) || isNaN(y)) {
                const computedStyle = window.getComputedStyle(
                  tippyTarget,
                  null
                );
                const transform = computedStyle.transform.substring(6);
                if (transform) {
                  x = parseFloat(transform.split(',')[4]);
                  y = parseFloat(transform.split(',')[5]);
                } else {
                  x = 0;
                  y = 0;
                }
              }

              const insets = trim(tippyTarget.style.inset).split(/\s+/);

              const isLeftFixed = insets[3] !== 'auto';
              const isRightFixed = insets[1] !== 'auto';
              const isTopFixed = insets[0] !== 'auto';
              const isBottomFixed = insets[2] !== 'auto';

              if (isRightFixed && event.deltaRect.right) {
                // translate when resizing from right edges
                x += event.deltaRect.right;
              } else if (isLeftFixed && event.deltaRect.left) {
                // translate when resizing from left edges
                x += event.deltaRect.left;
              } else if (isTopFixed && event.deltaRect.top) {
                // translate when resizing from top edges
                y += event.deltaRect.top;
              } else if (isBottomFixed && event.deltaRect.bottom) {
                // translate when resizing from bottom edges
                y += event.deltaRect.bottom;
              }

              tippyTarget.style.transform =
                'translate(' + x + 'px,' + y + 'px)';

              tippyTarget.setAttribute('data-x', `${x}`);
              tippyTarget.setAttribute('data-y', `${y}`);

              return event;
            },
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            (event: Interact.ResizeEvent) => options.listeners?.move?.(event)
          ),
        },
      }
    )
  );
}
