<template>
  <div
    class="left-source js-copyright js-source"
    :style="{
      transformOrigin: 'left',
    }"
  >
    <CopyRightVue 
      v-if="copyright" 
      :crawl-url="copyright.crawlUrl" 
      :licence-type="copyright.licenceType"
      :source-mark="copyright.sourceMark"
      :is-user-upload="copyright.isUserUpload"
      :upload-user-id="copyright.uploadUserId"
    />
  </div>
  <div
    class="right-source js-group-intro js-source js-source-right"
    :style="{
      transformOrigin: 'right',
    }"
  >
    <GroupPdfToolTip v-show="isGroupPdf" />
  </div>
  <div 
    v-if="checkOpenPaper(paperId, isPrivatePaper)" 
    class="right-source js-pdf-versions js-source-right" 
    :style="{
      transformOrigin: 'right',
      display: 'flex',
    }" 
  >
    <VersionSwitcher
      v-show="!isGroupPdf"
      :paper-id="paperId"
    />
  </div>
</template>

<script lang="ts" setup>
import GroupPdfToolTip from '../Common/GroupPdfToolTip.vue';
import { Nullable } from '~/src/typings/global';
import CopyRightVue, { CopyrightProps } from '../Copyright/index.vue'
import VersionSwitcher from './VersionSwitcher.vue';
import { checkOpenPaper } from '~/src/api/helper';

interface PdfRightMarkerProps {
  copyright: Nullable<CopyrightProps>;
  isGroupPdf?: boolean;
  paperId: string;
  isPrivatePaper: boolean;
}

defineProps<PdfRightMarkerProps>()


</script>

<style lang="less" scoped>
.left-source {
  position: absolute;
  top: 8px;
  left: 50%;
  display: flex;
  padding-left: 16px;
  display: none;
  z-index: 99;
}

.right-source {
  position: absolute;
  top: 8px;
  right: 50%;
  padding-right: 16px;
  display: none;
  z-index: 99;
}

.mobile-viewport {
  .left-source {
    left: 0;
    transform: none !important;
    right: auto;
    height: 24rpx;
    justify-content: flex-end;
    width: 100%;
    top: 0;
    border-bottom: 0.53333vw solid #e1e3e6;
    background: #ffffff;
  }
}
</style>
