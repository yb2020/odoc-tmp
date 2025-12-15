<template>
  <Exception
    v-if="baseStore.pageError"
    :code="baseStore.pageError.code"
    :message="baseStore.pageError.message"
    :url="href"
  />
  <Private v-else-if="!pdfStatusInfo.hasPdfAccessFlag" />
  <MobileView v-else-if="IS_MOBILE" />
  <Main v-else />
  <LimitDialog />
  <OcrModal />
  <PayDialog />
</template>

<script lang="ts" setup>
// 使用defineAsyncComponent正确处理懒加载组件
const Main = defineAsyncComponent(() => import('@/components/Main/index.vue'));
const Private = defineAsyncComponent(() => import('@/components/Private/index.vue'));

import {
  getPdfIdFromUrl,
} from '@/api/report';
import { computed, onMounted, ref, watch, defineAsyncComponent } from 'vue';
import {
  PageType,
  useReportVisitDuration,
} from '@common/utils/report';
import { pdfStatusInfo, selfNoteInfo } from '../store';
import MobileView from '../components/MobileView/index.vue';
import { IS_MOBILE } from '../util/env';
import { UserStatusEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';
import OcrModal from '../components/ToolTip/OcrModal.vue';
import LimitDialog from '@common/components/Premium/LimitDialog.vue';
import PayDialog from '@common/components/Premium/PayDialog.vue';
import Exception from './exception/exception.vue';
import { useBaseStore } from '../stores/baseStore';

const baseStore = useBaseStore();
const { href } = window.location;

if (IS_MOBILE) {
  document.body.classList.add('mobile-viewport');
}

useReportVisitDuration(
  () => String(selfNoteInfo.value?.userInfo?.id ?? ''),
  () => {
    const pdfId = selfNoteInfo.value?.pdfId ?? getPdfIdFromUrl();
    return {
      page_type: PageType.NOTE,
      type_parameter: String(pdfId),
    };
  },
  () => {
    return UserStatusEnum.TOURIST === pdfStatusInfo.value.noteUserStatus;
  }
);
</script>
