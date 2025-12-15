/* Copyright 2021, Milkdown by Mirone. */
import { Ctx, createCmd, createCmdKey, createSlice } from '@milkdown/core';
import { Fragment, Node, Slice } from '@milkdown/prose/model';
import {
  Command,
  Plugin,
  PluginKey,
  TextSelection,
} from '@milkdown/prose/state';
import { Decoration, DecorationSet } from '@milkdown/prose/view';
import { createPlugin } from '@milkdown/utils';

export interface HighlightPayload {
  segment: string;
  replacer?: string;
  start?: number;
}

export interface HighlighCtx {
  ignoredStrs?: string[];
  onHighlight?: (el: HTMLElement) => void;
}

export const HighlightPluginKey = new PluginKey('MILKDOWN_HIGHLIGHT');
export const HighlightCmd = createCmdKey<HighlightPayload>('Highlight');
export const highlightCtx = createSlice<null | HighlighCtx, 'HighlighCtx'>(
  null,
  'HighlighCtx'
);

function findRangeInNode(
  root: Node,
  query: string,
  start = 0,
  ignoredReg?: RegExp
) {
  let startIdx = -1;
  let endIdx = -1;

  // https://prosemirror.net/docs/guide/#doc
  // 纯文本匹配的开始位置一定比prosemirror的indexing要小
  let i = start;
  let j = start + query.length;
  let lastInclude = false;
  while (j <= root.nodeSize - 2) {
    let text = root.textBetween(i, j, '\n').replaceAll('\u00a0', ' ');
    if (ignoredReg) {
      text = text.replaceAll(ignoredReg, '');
    }
    const isFullMatch = text.startsWith(query);
    if (
      isFullMatch ||
      // 有可能出现最后一个字符不匹配的情况，比如加了个标点符号.
      (query.length === text.length && lastInclude)
    ) {
      startIdx = i;
      endIdx = j - (isFullMatch ? 0 : 1);
      break;
    }
    if (query.includes(text)) {
      lastInclude = true;
      j++;
    } else {
      lastInclude = false;
      i++;
      j++;
    }
  }

  return [startIdx, endIdx];
}

function findClosestParentHTMLEl(node: Element) {
  if (node && !(node instanceof HTMLElement)) {
    return node.parentElement;
  }

  return node;
}

class HighlightState {
  public time = 0;

  constructor(
    public query: string,
    public bgnIdx: number,
    public endIdx: number
  ) {}
}

const dispatchHighlight: (ctx: Ctx, payload?: HighlightPayload) => Command =
  (ctx, payload) => (state, dispatch, view) => {
    const hlState: HighlightState = HighlightPluginKey.getState(state);
    if (!payload || !hlState) return false;
    const { segment, replacer, start } = payload;
    const { ignoredStrs, onHighlight } = ctx.get(highlightCtx) || {};
    const ignoredReg = ignoredStrs?.length
      ? new RegExp(`(${ignoredStrs.join('|')})`, 'g')
      : undefined;
    if (dispatch) {
      hlState.query = segment ?? '';
      [hlState.bgnIdx, hlState.endIdx] = findRangeInNode(
        state.tr.doc,
        hlState.query,
        start,
        ignoredReg
      );
      if (hlState.bgnIdx === -1 && hlState.endIdx === -1) {
        return false;
      }
      let { tr } = state;
      const { schema } = state;
      if (replacer) {
        // 这里不考虑跨段的情况，如出现则变为一段
        let i = hlState.bgnIdx;
        let node = tr.doc.nodeAt(i);
        while (!node?.isText && i < hlState.endIdx) {
          i++;
          node = tr.doc.nodeAt(i);
        }

        const fragment = Fragment.from(schema.text(replacer));
        const slice = new Slice(fragment, 0, 0);
        tr = tr.replace(i, hlState.endIdx, slice);
        hlState.endIdx = i + fragment.size;
      }
      const posEnd = tr.doc.resolve(hlState.endIdx);
      const selection = new TextSelection(posEnd, posEnd);

      hlState.time = tr.time;
      dispatch(
        tr
          .setSelection(selection)
          // 不支持offset
          .scrollIntoView()
      );

      const node = view?.domAtPos(
        hlState.bgnIdx + (hlState.endIdx - hlState.bgnIdx) / 2
      )?.node as Element;
      const el = findClosestParentHTMLEl(node);
      if (el) {
        onHighlight?.(el);
      }
    }
    return true;
  };

export const highlightPlugin = createPlugin(() => {
  return {
    injectSlices: [highlightCtx],
    commands: (_, ctx) => [
      createCmd(HighlightCmd, (payload?: HighlightPayload) => {
        return dispatchHighlight(ctx, payload);
      }),
    ],
    prosePlugins: (_, ctx) => {
      ctx.set(highlightCtx, {});
      const plugin = new Plugin({
        key: HighlightPluginKey,
        state: {
          init() {
            return new HighlightState('', -1, -1);
          },
          apply(tr, hlState) {
            // 用time判断是否来源自身的变更
            if (tr.docChanged && hlState.time !== tr.time) {
              // 清空高亮
              return new HighlightState('', -1, -1);
            }

            return hlState;
          },
        },
        props: {
          decorations(state) {
            const hlState = this.getState(state);
            if (!hlState || hlState.bgnIdx === -1) {
              return;
            }

            return DecorationSet.create(state.doc, [
              Decoration.inline(hlState.bgnIdx, hlState.endIdx, {
                style: 'background: yellow',
              }),
            ]);
          },
        },
      });

      return [plugin];
    },
  };
});

export const highlight = highlightPlugin();
