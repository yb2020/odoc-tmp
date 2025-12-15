<template>
  <div class="dot-wrapper">
    <div
      class="dot"
      :style="{ background: `${item.fill}` }"
      @click.stop="() => {}"
      @mouseenter="handleMouseEnter"
      @mouseleave="handleMouseLeave"
    >
      <div
        v-if="isOwner"
        ref="optionsRef"
        :class="['dot-option-list', isBottom ? 'bottom-options' : '']"
      >
        <div
          v-for="(option, key) in styleMap"
          :key="option.color"
          class="dot-option"
          @click.stop="handleClick(key)"
        >
          <div
            class="dot-option-dot"
            :style="{ backgroundColor: option.color }"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { rectStyleMap, styleMap } from '@/style/select';
import { isOwner } from '~/src/store';
import { connect } from '@/dom/arrow';

import { useCommentGlobalState } from '@/hooks/useNoteState';
import { ToolBarType } from '@idea/pdf-annotate-core';
import {
  AnnotationAll,
  AnnotationSelect,
} from '~/src/stores/annotationStore/BaseAnnotationController';
import { useAnnotationStore } from '~/src/stores/annotationStore';

const annotationStore = useAnnotationStore();

const props = defineProps<{
  item: AnnotationAll;
}>();

const localStorageColor = useCommentGlobalState();

const handleClick = async (key: keyof typeof styleMap) => {
  localStorageColor.value.styleId = key;

  const params = {
    ...(props.item.type === ToolBarType.select
      ? styleMap[key]
      : rectStyleMap[key]),
  };

  if (props.item.type === ToolBarType.rect) {
    (params as any).stroke = params.color;
  }

  await annotationStore.controller.patchAnnotation(props.item.uuid, {
    styleId: +key,
    ...params,
  });

  if (
    props.item.type === ToolBarType.select &&
    !(props.item as AnnotationSelect).rectStr
  ) {
    return;
  }

  connect(props.item.uuid, true);
};

const isBottom = ref(false);

const optionsRef = ref<HTMLDivElement>();

const handleMouseEnter = () => {
  if (!optionsRef.value) {
    return;
  }

  const { bottom } = optionsRef.value.getBoundingClientRect();

  isBottom.value = bottom > window.innerHeight;
};

const handleMouseLeave = () => {
  isBottom.value = false;
};
</script>

<style lang="postcss" scoped>
.dot {
  width: 10px;
  height: 10px;
  background: var(--site-theme-brand);
  border-radius: 5px;
}

.dot-wrapper {
  position: relative;

  &:hover {
    .dot-option-list {
      display: block;
    }
  }

  .bottom-options {
    top: 0;
    transform: translate(-50%, -100%);
    bottom: auto;
  }
}
.dot-option-list {
  position: absolute;
  background: var(--site-theme-pdf-panel-secondary);
  border-radius: 1px;
  border: 1px solid var(--site-theme-border);
  bottom: 0px;
  left: 50%;
  z-index: 3;
  width: 24px;
  padding: 10px 0;
  transform: translate(-50%, 100%);
  display: none;

  .dot-option {
    width: 100%;
    height: 22px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;

    &:hover {
      background: var(--site-theme-bg-hover);
    }

    &-dot {
      border-radius: 5px;
      width: 10px;
      height: 10px;
    }
  }
}
</style>
