<template>
  <section class="attachments-list">
    <div class="flex items-center h-14 px-6 py-0">
      <span class="font-bold">{{ $t('home.library.attachment') }}</span>
      <Tooltip placement="top" :overlay-style="{ maxWidth: '350px' }">
        <template #title>
          <span
            >{{
              $t('home.library.uploadLimited', {
                count: limitCount,
                size: limitMBSize,
              })
            }}
          </span>
        </template>
        <InfoCircleOutlined class="ml-3" />
      </Tooltip>
      <span class="flex-1" />
      <CloseOutlined
        v-if="closable"
        class="opacity-40 hover:opacity-80"
        @click="onClose"
      />
    </div>
    <div class="flex items-center px-6">
      <p class="flex flex-col flex-1 m-0">
        <Progress
          :percent="percentage"
          :show-info="false"
          size="small"
          stroke-color="#A8AFBA"
          trail-color="#F5F5F5"
        />
        <span class="text-[13px] leading-[1] mt-2"
          >{{ totalUsedMBSize }}MB / {{ limitTotalMBSize }}MB</span
        >
      </p>
      <Upload
        name="file"
        :disabled="attachments.length >= limitCount"
        :max-count="1"
        :show-upload-list="false"
        :custom-request="onUpload"
      >
        <Button
          :disabled="attachments.length >= limitCount"
          class="btn-upload flex items-center ml-14"
          type="primary"
        >
          <UploadOutlined />{{ $t('home.upload.upload') }}
        </Button>
      </Upload>
    </div>
    <ol class="p-0 mt-3">
      <li v-if="uploading" class="py-2 px-6">
        <Item
          :file="uploadFile"
          :uploading="uploading"
          @cancel="onCancel"
        ></Item>
      </li>
      <li
        v-for="item in attachments"
        :key="item.id"
        class="py-2 px-6 hover:bg-[rgba(255,255,255,.1)]"
      >
        <Item
          :data="item"
          :removing-ids="removingIds"
          @preview="onPreview"
          @remove="onRemove"
        ></Item>
      </li>
      <li v-if="!attachments.length" class="py-24 text-center">
        <img class="w-20 h-20 mx-auto" :src="GROUP_EMPTY" alt="blank" />
        <p class="text-sm leading-5 mt-2">
          {{ $t('home.library.noAttachment') }}
        </p>
      </li>
    </ol>
  </section>
</template>

<script lang="ts" setup>
import {
  AttachmentInfo,
  GetAttachmentListResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment'
import { computed, PropType, toRef } from 'vue'
import { Upload, Button, Tooltip, Progress } from 'ant-design-vue'
import {
  CloseOutlined,
  InfoCircleOutlined,
  UploadOutlined,
} from '@ant-design/icons-vue'
import { bytes2MB } from '@/common/src/utils/file'
import GROUP_EMPTY from '@/assets/images/group-empty.svg'
import Item from './AttachmentItem.vue'

const props = defineProps({
  data: {
    type: Object as PropType<GetAttachmentListResponse>,
    default: null,
  },
  uploading: {
    type: Boolean,
    default: false,
  },
  uploadFile: {
    type: Object as PropType<null | File>,
    default: null,
  },
  removingIds: {
    type: Array as PropType<string[]>,
    default: () => [],
  },
  closable: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['preview', 'remove', 'upload', 'cancel', 'close'])

const data = toRef(props, 'data')
const attachments = computed(() => data.value?.list || [])
const limitMBSize = computed(() => bytes2MB(data.value?.sizeLimit))
const limitCount = computed(() => data.value?.amountLimit || 10)
const limitTotalMBSize = computed(() => bytes2MB(data.value?.totalSpace))
const totalUsedMBSize = computed(() => bytes2MB(data.value?.usedSpaceUsed, 1))
const percentage = computed(() => {
  if (!data.value) {
    return 0
  }
  const { usedSpaceUsed, totalSpace } = data.value

  return (usedSpaceUsed * 100) / totalSpace
})

const onClose = () => emit('close')
const onPreview = (item: AttachmentInfo) => emit('preview', item)
const onRemove = (item: AttachmentInfo) => emit('remove', item.id)
const onCancel = () => emit('cancel')
const onUpload = (data: any) => emit('upload', data)
</script>

<style lang="postcss">
.attachments-list {
  width: 336px;

  .ant-progress,
  .ant-progress-outer {
    height: 6px;
    line-height: 6px;
  }
  .ant-popover-inner-content {
    padding: 12px 0;
  }

  .btn-upload {
    display: flex;
  }
}
</style>
