<template>
  <div class="shape-color-container">
    <div
      v-if="originSize"
      class="shape-fontsize"
      @mouseleave="sizeLeave()"
    >
      <div
        v-for="size in sizeList"
        :key="size"
        :style="{
          backgroundColor: hoveringSize === size ? HIGHLIGHT : '',
          fontSize: size + 'px',
        }"
        @mouseenter="sizeEnter(size)"
        @click="clickSize(size)"
      >
        T
      </div>
    </div>
    <div
      class="shape-color"
      @mouseleave="colorLeave()"
    >
      <div
        v-for="color in colorList"
        :key="color"
        :style="{
          backgroundColor: hoveringColor === color ? HIGHLIGHT : '',
        }"
        @mouseenter="colorEnter(color)"
        @click="clickColor(color)"
      >
        <div :style="{ backgroundColor: colorMap[color] }" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { colorList, colorMap } from '@idea/pdf-annotate-core';
import { AnnotationColor } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { ref } from 'vue';
import {
  DEFAULT_FONT_SIZE,
  fontSizeList,
  useCommentGlobalState,
} from '~/src/hooks/useNoteState';

const localStorageColor = useCommentGlobalState();

interface ShapeToolbarProps {
  originColor: AnnotationColor;
  onColorEnter(color: AnnotationColor): void;
  onColorLeave(): void;
  selectColor(color: AnnotationColor): void;
  originSize?: number;
  onSizeEnter?(size: number): void;
  onSizeLeave?(): void;
  selectSize?(fontSize: number): void;
}

const props = defineProps<ShapeToolbarProps>()

let { originColor, originSize } = props;

const hoveringColor = ref(originColor);

const colorEnter = (color: AnnotationColor) => {
  hoveringColor.value = color;
  props.onColorEnter(color);
};

const colorLeave = () => {
  hoveringColor.value = originColor;
  props.onColorLeave();
};

const sizeEnter = (size: number) => {
  hoveringSize.value = size;
  props.onSizeEnter?.(size);
};

const sizeLeave = () => {
  hoveringSize.value = originSize as number;
  props.onSizeLeave?.();
};

const clickColor = (color: AnnotationColor) => {
  originColor = color;
  hoveringColor.value = color;
  props.selectColor(color);
  localStorageColor.value.shapeStyleId = color;
};

const sizeList = ref([...fontSizeList]);

const hoveringSize = ref(props.originSize ?? DEFAULT_FONT_SIZE);

const clickSize = (size: number) => {
  originSize = size;
  hoveringSize.value = size;
  localStorageColor.value.shapeFontSize = size;
  props.selectSize?.(size);
};

const HIGHLIGHT = '#484C4F';

</script>

<style lang="less" scoped>
.shape-color-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 33px;
  background-color: #323536;
  border-radius: 4px;
  .shape-color,
  .shape-fontsize {
    display: flex;

    >div {
      height: 24px;
      width: 24px;
      display: flex;
      justify-content: center;
      align-items: center;
      border-radius: 2px;
      cursor: pointer;
    }
  }

  .shape-color {
    >div {
      >div {
        height: 12px;
        width: 12px;
        border-radius: 50%;
      }
    }
  }

  .shape-fontsize {
    padding-right: 8px;
    margin-right: 8px;
    border-right: 1px solid gray;
  }
}
</style>