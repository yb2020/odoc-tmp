<template>
  <div
    class="dot"
    :style="{ background: `${item.fill}` }"
    @click.stop="() => {}"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div
      v-if="isOwner && isCurrentUser"
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
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { rectStyleMap, styleMap } from '@/style/select';
import { store, isOwner } from '~/src/store';
import { connect } from '@/dom/arrow';

import { useCommentGlobalState } from '@/hooks/useNoteState';
import { ToolBarType } from '@idea/pdf-annotate-core';
import {
  AnnotationAll,
  AnnotationSelect,
} from '~/src/stores/annotationStore/BaseAnnotationController';
import { useAnnotationStore } from '~/src/stores/annotationStore';

const annotationStore = useAnnotationStore();

const props = defineProps({
  item: { type: Object as () => AnnotationAll, required: true },
});

const useInfo = computed(() => store.state.user.userInfo);

const isCurrentUser = computed(
  () => props.item.commentatorInfoView?.userId === useInfo.value?.id
);

const localStorageColor = useCommentGlobalState();

const handleClick = async (key: keyof typeof styleMap) => {
  localStorageColor.value.styleId = key;

  const params: any = {
    ...(props.item.type === ToolBarType.select
      ? styleMap[key]
      : rectStyleMap[key]),
  };

  if (props.item.type === ToolBarType.rect) {
    params.stroke = params.color;
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

  connect(props.item.uuid);
};

const isBottom = ref(false);

const optionsRef = ref<HTMLDivElement>();

const handleMouseEnter = () => {
  const { bottom } = optionsRef.value!.getBoundingClientRect();

  isBottom.value = bottom > window.innerHeight;
};

const handleMouseLeave = () => {
  isBottom.value = false;
};
</script>

<style lang="less" scoped>
.dot {
  width: 10px;
  height: 10px;
  background: #447ac7;
  border-radius: 5px;
  position: absolute;

  right: 12px;
  bottom: 12px;

  &:hover {
    .dot-option-list {
      display: block;
    }
  }

  .dot-option-list {
    position: absolute;
    background: #282b2e;
    border-radius: 1px;
    border: 1px solid #797979;
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

      &-dot {
        border-radius: 5px;
        width: 10px;
        height: 10px;
      }
    }
  }

  .bottom-options {
    top: 0;
    transform: translate(-50%, -100%);
    bottom: auto;
  }
}
</style>
