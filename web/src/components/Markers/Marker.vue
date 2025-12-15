<template>
  <span
    v-if="bbox"
    ref="spanRef"
    class="marker"
    :style="{
      width: bbox.x1 - bbox.x0 + 'px',
      height: bbox.y1 - bbox.y0 + 'px',
      top: bbox.y0 + 'px',
      left: bbox.x0 + 'px',
    }"
    @click="handleClick"
  />
  <TippyVue
    v-if="isReferenceMarker(marker) && spanRef"
    ref="tippyRef"
    :trigger-ele="spanRef"
    @onShown="onShown"
  >
    <ReferenceVue
      :paper-id="marker.paperId"
      :paper-title="marker.refRaw || marker.refContent"
      :fetch-flag="startFetch"
      @update-content="handleTippyUpdate"
    />
  </TippyVue>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue';
import {
  // ReferenceMarker,
  FigureAndTableReferenceMarker,
  PdfBBox,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import { ReferenceMarker } from 'go-sea-proto/gen/ts/pdf/PdfParse';
import { scaleMarker } from '~/src/dom/markers';
import ReferenceVue from '../Tippy/Reference/index.vue';
import TippyVue from '../Tippy/index.vue';

function isReferenceMarker(
  marker: ReferenceMarker | FigureAndTableReferenceMarker
): marker is ReferenceMarker {
  return !!(marker as ReferenceMarker).paperId;
}

const props = defineProps<{
  marker: ReferenceMarker | FigureAndTableReferenceMarker;
  viewportWidth: number;
}>();

const bbox = computed(() => scaleMarker(props.marker, props.viewportWidth));

const spanRef = ref();

const tippyRef = ref();

const startFetch = ref(false);

const handleClick = () => {
  tippyRef.value.show();
};

const handleTippyUpdate = () => {
  tippyRef.value.update();
};

const onShown = () => {
  startFetch.value = true;
};
</script>

<style scoped lang="less">
.marker {
  position: absolute;
  background: #6f77bb;
  opacity: 0.3;
  cursor: pointer;
  display: block;
  z-index: 1;
}
</style>
