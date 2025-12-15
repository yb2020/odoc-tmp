<template>
  <div
    class="note-breadcrumb-container"
    :class="{
      'note-breadcrumb-fadeout': horizonal.fadeOut,
    }"
  >
    <div
      v-if="horizonal.fadeOut"
      class="note-breadcrumb-head"
    />
    <a-breadcrumb :style="{ fontSize: '14px' }">
      <a-breadcrumb-item
        v-for="(item, index) in list"
        :key="index"
        @click.native="clickNoteFolder([item.key])"
      >
        {{
          item.key === noteAllFolder.key
            ? `${$t(`common.text.all`)}${$t(
              `common.notes.${NoteSubType2I18nKey[type]}`,
              2
            )}`
            : item.title
        }}
      </a-breadcrumb-item>
    </a-breadcrumb>
    <div
      v-if="horizonal.fadeOut"
      class="note-breadcrumb-tail"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { useNote, getCharWidth } from '../useNote';
import { NoteBreadcrumb, NoteSubTypes, NoteSubType2I18nKey } from '../types';

const props = defineProps<{
  type: NoteSubTypes;
  noteState: ReturnType<typeof useNote>;
}>();
const sum = (a: number, b: number) => a + b;

const { noteAllFolder, noteBreadcrumbList, clickNoteFolder } = props.noteState;

const list = computed(() => noteBreadcrumbList.value?.slice());
const SEPERATOR_WIDTH = 20;
const FONT_SIZE = 14;
const MAX_TEXT_WIDTH = 380;
const FADEOUT_WIDTH = 20;
const horizonal = computed(() => {
  const width = noteBreadcrumbList.value
    ?.map((bread: NoteBreadcrumb) =>
      bread.title

        .split('')
        .map((char) => getCharWidth(char, FONT_SIZE))
        .reduce(sum, SEPERATOR_WIDTH)
    )
    .reduce(sum, -SEPERATOR_WIDTH + FADEOUT_WIDTH);

  if (width <= MAX_TEXT_WIDTH) {
    return {
      flexBasis: `${width}px`,
      fadeOut: false,
    };
  }

  return {
    flexBasis: `${Math.min(width, MAX_TEXT_WIDTH) + FADEOUT_WIDTH}px`,
    fadeOut: true,
  };
});
</script>

<style lang="less" scoped>
.note-breadcrumb-container {
  position: relative;
  flex-grow: 0;
  flex-shrink: 0;
  display: inline-block;
  overflow: hidden;
  background-color: var(--site-theme-bg-secondary);
  .note-breadcrumb-head,
  .note-breadcrumb-tail {
    position: absolute;
    top: 0;
    height: 100%;
    width: 20px;
    z-index: 10;
  }
  .note-breadcrumb-head {
    left: 0;
    background: linear-gradient(to right, var(--site-theme-bg-secondary), transparent);
  }
  .note-breadcrumb-tail {
    right: 0;
    background: linear-gradient(to right, transparent, var(--site-theme-bg-secondary));
  }
}
</style>
<style lang="less">
.note-breadcrumb-container {
  .ant-breadcrumb {
    overflow-x: auto;
    white-space: nowrap;
    text-align: left;
    padding-left: 0;
    padding-top: 0 !important;
    padding-bottom: 0 !important;
    &::-webkit-scrollbar {
      display: none;
    }
    
    .ant-breadcrumb-separator {
      color: var(--site-theme-text-tertiary);
    }
    
    > * {
      cursor: pointer;
      display: inline-block;
      > span {
        display: inline-block;
        color: var(--site-theme-text-secondary);
        
        &:hover {
          color: var(--site-theme-primary-color);
        }
      }
      &:last-of-type {
        > span {
          color: var(--site-theme-text-primary);
          font-weight: 500;
        }
      }
    }
  }
}
</style>
