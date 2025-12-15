<template>
  <FileUploaderModal
    v-bind="$props"
    @addSuccess="$emit('addSuccess')"
    @close="$emit('close')"
    @refreshList="$emit('refreshList')"
  >
    <template #addButton="slotProps">
      <SearchCollectBtn
        :paper-id="slotProps.paperId"
        :folder-id="folderId"
        :is-collected="slotProps.isCollected"
        @addSuccess="$emit('addSuccess')"
      />
    </template>
  </FileUploaderModal>
</template>
<script lang="ts" setup>
import SearchCollectBtn from './SearchCollect.vue'
import FileUploaderModal from './Modal.vue'
import { BreadCrumb } from '@/stores/library'
import { PropType } from 'vue'
import { LimitDialogReportParams } from '../helper'

defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  folderId: {
    type: String,
    default: '',
  },
  collectLimitDialogReportParams: {
    type: Object as PropType<LimitDialogReportParams>,
    default: null,
  },
  selectedKey: {
    type: String,
    required: true,
  },
  breadCrumbList: {
    type: Array as () => BreadCrumb[],
    required: true,
  },
})

defineEmits(['addSuccess', 'close', 'refreshList'])
</script>
<style lang="less" scoped>
.file-uploader {
  &-section {
    .title {
      font-weight: 600;
      color: #262625;
      line-height: 22px;
      margin-top: 4px;
      margin-bottom: 12px;
    }
  }
  &-drag {
    margin-bottom: 40px;
  }
}
</style>
