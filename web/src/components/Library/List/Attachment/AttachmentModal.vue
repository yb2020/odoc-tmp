<template>
  <Modal
    width="fit-content"
    :centered="true"
    :visible="!!docId"
    :header="null"
    :footer="null"
    :body-style="{
      padding: '0 0 8px',
    }"
    @cancel="emit('close')"
  >
    <List
      :data="data"
      :uploading="uploading"
      :upload-file="uploadFile"
      :removing-ids="removingIds"
      @preview="onPreview"
      @remove="onRemove"
      @upload="onUpload"
      @cancel="onCancel"
    ></List>
  </Modal>
</template>

<script lang="ts" setup>
import { computed, toRef } from 'vue'
import { Modal } from 'ant-design-vue'
import { useDocAttachments } from './useDocAttachments'
import { calculateFileMD5 } from '@/common/src/utils/md5'
import { AttachmentInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment'
import List from './AttachmentList.vue'
import { PageType } from '@/common/src/utils/report'
import { ElementName, reportElementClick } from '@/utils/report'
const props = defineProps({
  docId: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['attached', 'removed', 'close'])

const docId = toRef(props, 'docId')
const isReady = computed(() => !!docId.value)
const {
  data,
  uploading,
  uploadFile,
  removingIds,
  doUpload,
  doRemove,
  cancelUpload,
} = useDocAttachments(docId, isReady, calculateFileMD5)

const onUpload = async ({ file }: { file: File; filename: string }) => {
  await doUpload(file)
  emit('attached', docId.value, data.value?.list.length)
}

const onCancel = () => {
  cancelUpload()
}

const onPreview = (item: AttachmentInfo) => {
  reportElementClick({
    page_type: PageType.NOTE,
    type_parameter: 'none',
    element_name: ElementName.attachments,
    element_parameter: item.id,
  })

  window.open(item.url)
}

const onRemove = async (attachmentId: string) => {
  await doRemove({
    attachmentId,
  })
  emit('removed', docId.value, data.value?.list.length)
}
</script>

<style lang="postcss">
.attachments-rm-modal {
  .ant-modal-confirm-content {
    margin-left: 34px !important;
    line-height: 1;
    color: theme('colors.rp-neutral-8');
  }
}
</style>
