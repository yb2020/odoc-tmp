import {
  editorViewOptionsCtx,
  parserCtx,
  schemaCtx,
  serializerCtx,
} from '@milkdown/core';
import { DOMParser, Slice } from '@milkdown/prose/model';
import { Plugin, PluginKey } from '@milkdown/prose/state';
import { createPlugin } from '@milkdown/utils';
import { parse, HtmlGenerator } from 'latex.js';

export const LatexKey = new PluginKey('MILKDOWN_LATEX');

export const latexPlugin = createPlugin(() => {
  return {
    prosePlugins: (_, ctx) => {
      const schema = ctx.get(schemaCtx);

      // Set editable props for https://github.com/Saul-Mirone/milkdown/issues/190
      ctx.update(editorViewOptionsCtx, (prev) => ({
        ...prev,
        editable: prev.editable ?? (() => true),
      }));

      const plugin = new Plugin({
        key: LatexKey,
        props: {
          handlePaste: (view, event, originalSlice) => {
            const parser = ctx.get(parserCtx);
            const serializer = ctx.get(serializerCtx);
            const editable = view.props.editable?.(view.state);
            const { clipboardData } = event;
            if (!editable || !clipboardData) {
              return false;
            }

            const currentNode = view.state.selection.$from.node();
            if (currentNode.type.spec.code) {
              return false;
            }

            let text = clipboardData.getData('text/plain');
            if (!/\\documentclass/.test(text)) {
              return false;
            }
            const generator = new HtmlGenerator({ hyphenate: false });
            const doc = parse(text, { generator });
            const html: DocumentFragment = doc.domFragment();
            const node = DOMParser.fromSchema(schema).parse(html);
            text = serializer(node);

            const slice = parser(text);
            if (!slice || typeof slice === 'string') return false;

            view.dispatch(
              view.state.tr.replaceSelection(
                new Slice(
                  slice.content,
                  originalSlice.openStart,
                  originalSlice.openEnd
                )
              )
            );

            return true;
          },
        },
      });

      return [plugin];
    },
  };
});

export const latex = latexPlugin();
