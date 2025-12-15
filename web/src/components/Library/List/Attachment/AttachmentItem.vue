<template>
  <div class="attachment-item flex items-center leading-6">
    <span
      :class="`
        flex
        items-center
        flex-1
        min-w-0
        pr-2
        text-inherit
        cursor-pointer
        ${!uploading ? 'hover:text-rp-blue-6' : ''}
      `"
      @click="$emit('preview', data)"
    >
      <LoadingOutlined v-if="uploading" />
      <template v-else>
        <EyeOutlined v-if="data && data.isPreview" />
        <DownloadOutlined v-else />
      </template>
      <span
        class="ml-2 block text-ellipsis overflow-hidden whitespace-nowrap"
        >{{ uploading ? file && file.name : data && data.name }}</span
      >
    </span>
    <CloseOutlined
      v-if="uploading"
      class="ml-2 text-rp-neutral-6 cursor-pointer opacity-40 hover:opacity-80"
      @click="$emit('cancel')"
    />
    <template v-else>
      <span class="text-rp-neutral-6">{{
        formatBitSize(data && data.size)
      }}</span>
      <LoadingOutlined
        v-if="removingIds.includes(data && data.id)"
        class="ml-2"
      />
      <DeleteOutlined
        class="ml-2 text-rp-neutral-6 opacity-40 hover:opacity-80"
        @click="$emit('remove', data)"
      />
    </template>
  </div>
</template>

<script lang="ts" setup>
import { AttachmentInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment'
import { PropType } from 'vue'
import {
  CloseOutlined,
  LoadingOutlined,
  DeleteOutlined,
  EyeOutlined,
  DownloadOutlined,
} from '@ant-design/icons-vue'
import { formatBitSize } from '@/common/src/utils/file'

defineProps({
  data: {
    type: Object as PropType<AttachmentInfo>,
    default: null,
  },
  file: {
    type: Object as PropType<File | null>,
    default: null,
  },
  uploading: {
    type: Boolean,
    default: false,
  },
  removingIds: {
    type: Array as PropType<string[]>,
    default: () => [],
  },
})
</script>

<style lang="postcss">
.attachment-item {
  .anticon {
    vertical-align: middle;
  }
}
</style>
