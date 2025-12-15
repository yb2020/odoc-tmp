<template>
  <div
    ref="elRef"
    class="milkdown-wrapper"
    :class="{
      'milkdown-wrapper--disabled': disabled,
    }"
  />
</template>

<script lang="ts" setup>
import { onMounted, ref, toRef, watch } from 'vue';
import {
  Editor,
  defaultValueCtx,
  rootCtx,
  editorViewOptionsCtx,
  MilkdownPlugin,
  editorViewCtx,
} from '@milkdown/core';
import { AtomList } from '@milkdown/utils';
import { nordLight } from '@milkdown/theme-nord';
import { heading, image } from '@milkdown/preset-commonmark';
import { image as imageResizable } from './plugins/image';
import {
  imagePickerPreset,
  imagePickerView,
  ImageOptions,
} from 'milkdown-plugin-image-picker';
import { gfm } from '@milkdown/preset-gfm';
import { listener, listenerCtx } from '@milkdown/plugin-listener';
import { menu, menuPlugin } from '@milkdown/plugin-menu';
import { upload, uploadPlugin } from '@milkdown/plugin-upload';
import { placeholderCtx } from 'milkdown-plugin-placeholder';
import {
  highlightCtx,
  HighlightCmd,
  HighlightPayload,
} from './plugins/highlight';

import { replaceAll, getHTML, getMarkdown, callCommand } from '@milkdown/utils';
import { elementOffsetTop } from '@common/utils/offset';
import { getMenuConfig, getDefaultPlugins } from './utils';

import 'material-icons/iconfont/outlined.css';
import 'katex/dist/katex.min.css';

export interface UpdatePayload {
  text: string;
  markdown: string;
  count: number;
  html: string;
}

const props = defineProps<{
  lang?: string;
  disabled?: boolean;
  modelValue: string;
  hlReplace?: HighlightPayload;
  width?: string;
  placeholder: string;
  plugins?: MilkdownPlugin[];
  ignoreMenuKeys?: string[];
  defaultPlugins?: Record<string, MilkdownPlugin | AtomList>;
  shouldFixHeadingSelector?: boolean;
  uploader?: (files: FileList) => Promise<
    Array<{
      src: string;
      alt: string;
    }>
  >;
}>();
const emit = defineEmits<{
  (event: 'update', payload: UpdatePayload): void;
  (event: 'created', editor: Editor): void;
}>();

const modelValue = defineModel('modelValue', { default: '' });
const disabled = toRef(props, 'disabled');
const hlReplace = toRef(props, 'hlReplace');
const width = toRef(props, 'width');
const uploader = toRef(props, 'uploader');

const elRef = ref<HTMLDivElement>();
const editorRef = ref<null | Editor>();
const count = ref(0);

const updateMark = () => {
  editorRef.value?.action((ctx) => {
    const md = getMarkdown()(ctx);
    if (md !== modelValue.value) {
      const fn = replaceAll(modelValue.value);
      fn(ctx);
    }
  });
};

watch(modelValue, updateMark);

watch(hlReplace, () => {
  editorRef.value?.action((ctx) => {
    callCommand(HighlightCmd, hlReplace.value)(ctx);
  });
});

watch(
  () => props.placeholder,
  (newVal) => {
    // TODO placeholder插件提供一个更新placeholder的方法
    editorRef.value?.action((ctx) => {
      const view = ctx.get(editorViewCtx);
      ctx.update(placeholderCtx, () => {
        return newVal;
      });
      view.updateState(view.state);
    });
  }
);

watch(width, (w) => {
  if (w) {
    const el = elRef.value?.querySelector('.milkdown') as null | HTMLDivElement;
    if (el) {
      el.style.width = w;
    }
  }
});

const fixHeadingSelector = () => {
  if (!props.shouldFixHeadingSelector) {
    return;
  }
  // menu plugin的定位错误，强行更正
  // node_modules/@milkdown/plugin-menu/src/select.ts
  const btnHeading = elRef.value?.querySelector(
    '.milkdown-menu .menu-selector'
  );
  const selectHeading = elRef.value?.querySelector(
    '.milkdown-menu .menu-selector-list'
  ) as HTMLElement;
  if (btnHeading instanceof HTMLElement) {
    btnHeading.addEventListener('mousedown', () => {
      selectHeading.style.left = `${btnHeading.getBoundingClientRect()
        ?.left}px`;
    });
  }
};

const fixScrolling = () => {
  const elMilkdown = elRef.value?.querySelector('.milkdown') as HTMLDivElement;
  let elInner = elRef.value?.querySelector('.milkdown-inner');
  if (elInner || !elMilkdown) {
    return;
  }
  elInner = document.createElement('div');
  elInner.classList.add('milkdown-inner');
  elMilkdown?.replaceWith(elInner);
  elInner.appendChild(elMilkdown);
};

const scrollIntoHighlightView = (el: HTMLElement) => {
  const scrollEl = elRef.value?.querySelector(
    '.milkdown-inner'
  ) as null | HTMLElement;

  if (scrollEl) {
    const top = elementOffsetTop(el, scrollEl);
    scrollEl.scrollTop = Math.max(0, top - scrollEl.offsetHeight / 4);
  }
};

watch(disabled, (flag: boolean) => {
  if (!flag) {
    setTimeout(fixScrolling, 2000);
  }
});

onMounted(() => {
  const imagePlugin = imageResizable<ImageOptions>(imagePickerPreset())({
    uploader: uploader.value,
  });

  const editor = Editor.make()
    .config((ctx) => {
      ctx.set(rootCtx, elRef.value as Element);
      // 注意这里set一定要是字符串
      // @notice 这里不设置modelVale，因为要触发update事件传递text, html以便外层初始化
      ctx.set(defaultValueCtx, '');
      ctx.set(placeholderCtx, props.placeholder);

      ctx.update(editorViewOptionsCtx, (prev) => ({
        ...prev,
        editable: () => !disabled.value,
      }));

      let text = '';
      ctx.get(listenerCtx).markdownUpdated((_ctx, markdown, _prevMarkdown) => {
        modelValue.value = markdown;
        count.value = [...text].length;
        // 此时dom还没有更新，所以要等到下一个tick
        setTimeout(() => {
          const fn = getHTML();
          emit('update', {
            text,
            markdown,
            count: count.value,
            html: fn(_ctx),
          });
        }, 0);
      });

      ctx.get(listenerCtx).updated((_ctx, node) => {
        text = node
          .textBetween(0, node.content.size, '\n')
          .replaceAll('\u00a0', ' ');
      });
    })
    .use(nordLight.override(imagePickerView))
    .use(
      gfm
        .configure(heading, {
          displayHashtag: false,
        })
        .replace(image, imagePlugin)
    )
    .use(
      upload.configure(uploadPlugin, {
        uploader: async (files, schema) => {
          const images = (await uploader.value?.(files)) ?? [];

          return images.map(({ src, alt }) => {
            return schema.nodes.image.createAndFill({
              src,
              alt,
            })!;
          });
        },
      })
    )
    .use(
      menu.configure(menuPlugin, {
        config: getMenuConfig(props.lang, props.ignoreMenuKeys),
      })
    )
    .use(listener);

  if (props.plugins?.length) {
    editor.use(props.plugins);
  }

  getDefaultPlugins(props.defaultPlugins).forEach((p) => editor.use(p));

  editor.create().then((editor) => {
    editorRef.value = editor;

    fixHeadingSelector();
    fixScrolling();
    updateMark();

    editorRef.value.action((ctx) => {
      ctx.set(highlightCtx, {
        onHighlight: scrollIntoHighlightView,
      });
    });

    emit('created', editor);
  });
});

defineExpose({
  getHTML: () => {
    return new Promise((resolve, reject) => {
      const timeout = window.setTimeout(() => {
        reject(Error('timeout'));
      }, 3000);
      window.setTimeout(() => {
        editorRef.value?.action((ctx) => {
          const fn = getHTML();
          window.clearTimeout(timeout);
          resolve(fn(ctx));
        });
      }, 0);
    });
  },
  forceUpdatMark: (md: string, cb: () => void) => {
    editorRef.value?.action((ctx) => {
      const fn = replaceAll(md);
      fn(ctx);
      cb();
    });
  },
});
</script>

<style lang="less" scoped>
@h1Size: 21px;
@h2Size: 16px;
@h3Size: 14px;
@h4Size: 12px;
@pSize: 14px;

.disabled() {
  opacity: 0.5;
  pointer-events: none;
}

.milkdown-wrapper {
  :deep(.milkdown) {
    font-family: inherit;
    color: inherit;

    .editor {
      padding: 24px 0;

      & > * {
        margin: 10px 0;
      }

      .list-item {
        margin: 0;

        & > .list-item_label {
          font-size: 12px;
          align-items: center;
          color: inherit;
        }
      }
    }

    .ProseMirror[data-placeholder]::before {
      margin: 10px 0;
      font-size: 14px;
    }

    h1 {
      font-size: @h1Size;
      margin: 10px 0;
    }

    h2 {
      font-size: @h2Size;
      margin: 10px 0;
    }

    h3 {
      font-size: @h3Size;
    }

    h4 {
      font-size: @h4Size;
    }

    h5,
    h6 {
      font-size: unset;
    }

    p {
      font-size: @pSize;
    }

    ol,
    ul {
      padding-left: 20px;
    }

    blockquote {
      line-height: 24px;
      padding-left: 8px;
      padding-right: 8px;
      padding-top: 10px;
      padding-bottom: 10px;
      border-left-width: 4px;
      border-left-style: solid;
      border-left-color: #c9cdd4;

      p {
        margin: 0;
      }
    }

    a {
      color: #0d3e81;
    }
  }

  :deep(.milkdown-menu) {
    justify-content: space-around;

    .divider {
      margin: 12px 2px;
    }

    .button {
      margin: 8px 0;
      width: 27px;
      height: 27px;

      &[disabled] {
        display: flex;
        .disabled();
      }
    }

    .menu-selector {
      width: 104px;
      padding: 4px 8px;
      margin: 8px 0;

      &[disabled] {
        display: flex;
        .disabled();
      }

      &-wrapper {
        border: none;

        &.disabled {
          display: block;
          .disabled();

          & + .divider {
            display: block;
          }
        }
      }
    }

    .menu-selector-list {
      width: 180px;
    }

    .menu-selector-list-item {
      &[data-id='1'] {
        font-size: @h1Size;
      }

      &[data-id='2'] {
        font-size: @h2Size;
      }

      &[data-id='3'] {
        font-size: @h3Size;
      }

      &[data-id='4'] {
        font-size: @pSize;
      }
    }
  }
}
</style>
