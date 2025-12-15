/* Copyright 2021, Milkdown by Mirone. */
import { ImageOptions, image as imageNode } from '@milkdown/preset-commonmark';
import type { NodeCreator } from '@milkdown/utils';
import { NodeSelection, Plugin, PluginKey } from '@milkdown/prose/state';
import { findSelectedNodeOfType } from '@milkdown/prose';
import { EditorView } from '@milkdown/prose/view';
import { Node, NodeType } from '@milkdown/prose/model';
import interact from 'interactjs';
import { Tooltip } from 'ant-design-vue';
import { createApp } from 'vue';

const key = new PluginKey('MILKDOWN_IMAGE_RESIZE');
const keyTip = 'MILKDOWN_IMAGE_RESIZE_TIP';
const preventDrag = (e: DragEvent) => {
  e.preventDefault();
};

function resizableElement(view: EditorView, type: NodeType) {
  const selectedNode = findSelectedNodeOfType(view.state.selection, type);
  if (!selectedNode) return;

  const $wrapper = view.nodeDOM(selectedNode.pos) as HTMLElement;
  const $resizable = $wrapper.querySelector('img') as HTMLElement;
  const $tip = document.createElement('span');
  if (!$wrapper || !$resizable) return;

  const insertTip = () => {
    const flag = window.localStorage.getItem(keyTip);
    if (!flag) {
      window.localStorage.setItem(keyTip, '1');

      const app = createApp(Tooltip, {
        // title: '拖动图片边缘可调整大小',
        title: 'Drag image edges to adjust size',
        visible: true,
        placement: 'topLeft',
      });
      $tip.setAttribute('style', 'position: absolute; top: 0; left: 0;');

      $wrapper.appendChild($tip);
      app.mount($tip);
      return app;
    }
  };
  const tiper = insertTip();
  const resizer = interact($resizable).resizable({
    margin: 4,
    edges: {
      left: true,
      right: true,
      top: true,
      bottom: true, // '.resize-corner',
    },
    modifiers: [
      interact.modifiers.restrictSize({
        min: { width: 20, height: 20 },
        max: view.dom,
      }),
    ],
    inertia: true,
    listeners: {
      start: () => {
        $wrapper?.addEventListener('dragstart', preventDrag);
      },
      move: (event: Interact.ResizeEvent) => {
        if (!event.deltaRect) {
          return;
        }

        Object.assign($resizable.style, {
          width: `${event.rect.width}px`,
          height: `${event.rect.height}px`,
        });
      },
      end: () => {
        $wrapper?.removeEventListener('dragstart', preventDrag);

        const src = $resizable.getAttribute('src');
        const parsed = new URL(src!);
        const hashParams = new URLSearchParams(parsed.hash.slice(1));
        hashParams.set('width', `${$resizable.offsetWidth}px`);
        hashParams.set('height', `${$resizable.offsetHeight}px`);
        parsed.hash = `#${hashParams.toString()}`;
        view.dispatch(
          view.state.tr.setNodeAttribute(
            selectedNode!.pos,
            'src',
            parsed.toString()
          )
        );
      },
    },
  });

  return {
    destroy: () => {
      tiper?.unmount();
      $wrapper.removeChild($tip);
      resizer.unset();
    },
  };
}

export const image = <T extends object = ImageOptions>(
  node = imageNode
): NodeCreator<string, T> =>
  node.extend((original) => {
    const prosePlugins = original.prosePlugins!;

    return {
      ...original,
      // schema: (ctx) => {
      //   const schemaProps = original.schema(ctx)
      //   return {
      //     ...schemaProps,
      //     toDOM: (node) => {
      //       // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      //       // @ts-ignore
      //       const [_, attrs] = schemaProps.toDOM(node)
      //       const parsed = new URL(attrs.src)
      //       const hashParams = new URLSearchParams(parsed.hash.slice(1))

      //       return [
      //         _,
      //         {
      //           ...attrs,
      //           width: hashParams.get('width'),
      //         },
      //       ]
      //     },
      //   }
      // },
      view:
        (ctx) =>
        (node: Node, ...rest) => {
          const originalView = original.view!(ctx);
          const originalProps = originalView(node, ...rest);

          const setSize = ($span: HTMLElement, node: Node) => {
            const $img = $span.querySelector<HTMLImageElement>('img');
            if ($img) {
              const parsed = new URL(node.attrs.src);
              const hashParams = new URLSearchParams(parsed.hash.slice(1));
              const width = hashParams.get('width');
              const height = hashParams.get('height');

              Object.assign($img.style, {
                width,
                height,
              });
            }
          };
          setSize(originalProps.dom as HTMLSpanElement, node);

          return {
            ...originalProps,
            update: (updatedNode, ...rest) => {
              const result = originalProps.update!(updatedNode, ...rest);

              if (result) {
                setSize(originalProps.dom as HTMLSpanElement, updatedNode);
              }

              return result;
            },
          };
        },
      prosePlugins: (type, ctx) => [
        ...prosePlugins(type, ctx),
        new Plugin({
          key,
          view: (editorView) => {
            let instance: ReturnType<typeof resizableElement> | undefined;
            const shouldDisplay = (view: EditorView) => {
              const { selection } = view.state;

              if (!view.hasFocus()) return false;

              return (
                selection instanceof NodeSelection &&
                selection.node.type.name === 'image'
              );
            };
            const renderByView = (view: EditorView) => {
              if (shouldDisplay(view) && view.editable) {
                if (instance) return;
                instance = resizableElement(editorView, type);
              } else {
                instance?.destroy();
                instance = undefined;
              }
            };
            renderByView(editorView);

            return {
              update: (view, prevState) => {
                const isEqualSelection =
                  prevState?.doc.eq(view.state.doc) &&
                  prevState.selection.eq(view.state.selection);
                if (isEqualSelection) return;

                renderByView(view);
              },
              destroy: () => {
                instance?.destroy();
                instance = undefined;
              },
            };
          },
        }),
      ],
    };
  });
