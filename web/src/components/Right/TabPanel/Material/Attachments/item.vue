<template>
  <div class="attachment-item flex items-center leading-6">
    <span
      class="flex items-center flex-1 min-w-0 pr-2 text-inherit cursor-pointer"
      :class="{
        'hover:text-rp-blue-6': !uploading,
      }"
      href="javascript:void"
      @click="$emit('preview', data)"
    >
      <LoadingOutlined v-if="uploading" />
      <template v-else>
        <EyeOutlined v-if="data?.isPreview" />
        <DownloadOutlined v-else />
      </template>
      <span
        class="ml-2 block text-ellipsis overflow-hidden whitespace-nowrap"
      >{{ uploading ? file?.name : data?.name }}</span>
    </span>
    <CloseOutlined
      v-if="uploading"
      class="ml-2 text-rp-neutral-6 cursor-pointer opacity-40 hover:opacity-80"
      @click="$emit('cancel')"
    />
    <template v-else>
      <span class="text-rp-neutral-6">{{ formatBitSize(data?.size) }}</span>
      <LoadingOutlined
        v-if="removingIds?.includes(data?.id!)"
        class="ml-2"
      />
      <DeleteOutlined
        class="ml-2 text-rp-neutral-6 opacity-40 hover:opacity-80"
        @click="$emit('remove', data)"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { AttachmentInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment';
import {
  LoadingOutlined,
  DownloadOutlined,
  EyeOutlined,
  DeleteOutlined,
  CloseOutlined,
} from '@ant-design/icons-vue';
import { toRef } from 'vue';
import { formatBitSize } from '~/src/util/file';

const props = defineProps<{
  data?: AttachmentInfo;
  file?: null | File;
  uploading?: boolean;
  removingIds?: string[];
}>();

const data = toRef(props, 'data');
</script>

<style lang="postcss" scoped>
.attachment-item {
  .anticon {
    vertical-align: middle;
  }
}
</style>
