<template>
  <div
    v-show="dir"
    class="back-to-scroll"
    @click="goBack"
  >
    <div class="mask" />
    <div class="text">
      <i
        :class="[
          'aiknowledge-icon',
          { 'icon-arrow-down': dir === 'down', 'icon-arrow-up': dir === 'up' },
        ]"
        aria-hidden="true"
      />
      {{ $t('viewer.goBack') }}
    </div>
  </div>
</template>
<script lang="ts" setup>
import { ViewerController } from '~/../../packages/pdf-annotate-viewer/typing';
import { scrollToPDFLastPosition } from '~/src/dom/pdf';

const props = defineProps<{
  dir: 'up' | 'down' | false;
  pdfViewInstance: ViewerController;
}>();

const goBack = () => {
  scrollToPDFLastPosition(props.pdfViewInstance);
};
</script>
<style lang="less" scoped>
.back-to-scroll {
  position: absolute;
  bottom: 115px;
  width: 156px;
  height: 40px;
  border-radius: 20px;
  overflow: hidden;
  opacity: 0.8;
  left: calc(50% - 78px);
  cursor: pointer;
  z-index: 99;
  .mask {
    background: #677080;
    width: 100%;
    height: 100%;
  }

  .text {
    color: #fff;
    font-size: 14px;
    line-height: 40px;
    position: absolute;
    top: 0;
    width: 100%;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;
    .aiknowledge-icon {
      font-size: 16px;
      margin-right: 5px;
    }
  }
}
</style>
