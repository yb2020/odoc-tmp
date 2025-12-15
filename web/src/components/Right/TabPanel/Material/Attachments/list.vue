<template>
  <section class="attachments-list">
    <div class="flex items-center h-14 px-6 py-0">
      <span class="font-bold">{{ $t('info.attachments') }}</span>
      <a-tooltip
        placement="top"
        :overlayStyle="{ maxWidth: '350px' }"
      >
        <template #title>
          <span>{{
            $t('info.attachmentSizeTip', { limitCount, limitMBSize })
          }}</span>
        </template>
        <InfoCircleOutlined class="ml-3" />
      </a-tooltip>
      <span class="flex-1" />
      <CloseOutlined
        class="opacity-40 hover:opacity-80"
        @click="onClose"
      />
    </div>
    <div class="flex items-center px-6">
      <p class="flex flex-col flex-1 m-0">
        <a-progress
          :percent="percentage"
          :show-info="false"
          size="small"
          stroke-color="#A8AFBA"
          trail-color="#F5F5F5"
        />
        <span class="text-[13px] leading-[1] mt-2">{{ totalUsedMBSize }}MB / {{ limitTotalMBSize }}MB</span>
      </p>
      <a-upload
        name="file"
        :disabled="attachments.length >= limitCount"
        :maxCount="1"
        :showUploadList="false"
        :customRequest="onUpload"
      >
        <a-button
          class="!flex !items-center ml-14"
          :class="!isCollected ? '!text-white !opacity-40' : ''"
          :disabled="attachments.length >= limitCount"
          type="primary"
          @click="onBeforeUpload"
        >
          <template #icon>
            <UploadOutlined />
          </template>
          {{ $t('info.upload') }}
        </a-button>
      </a-upload>
    </div>
    <ol class="mt-3 mb-0">
      <li
        v-if="uploading"
        class="py-2 px-6"
      >
        <Item
          :file="uploadFile"
          :uploading="uploading"
          @cancel="onCancel"
        />
      </li>
      <li
        v-for="item in attachments"
        :key="item.id"
        class="py-2 px-6 hover:bg-[rgba(255,255,255,.1)]"
      >
        <Item
          :data="item"
          :removingIds="removingIds"
          @preview="onPreview"
          @remove="onRemove"
        />
      </li>
      <li
        v-if="!attachments.length"
        class="py-24 text-center"
      >
        <img
          class="w-20 h-20 mx-auto"
          src="@/assets/images/empty.svg"
          alt="blank"
        >
        <p class="text-sm leading-5 mt-2">
          暂无附件
        </p>
      </li>
    </ol>
  </section>
</template>

<script setup lang="ts">
import {
  CloseOutlined,
  UploadOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons-vue';
import {
  AttachmentInfo,
  GetAttachmentListResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment';
import { computed, toRef } from 'vue';
import { bytes2MB } from '~/src/util/file';
import Item from './item.vue';

const props = defineProps<{
  data: GetAttachmentListResponse;
  isCollected: boolean;
  uploading: boolean;
  uploadFile: null | File;
  removingIds: string[];
  onUpload?: (x: any) => Promise<void>;
  onBeforeUpload?: (x: MouseEvent) => void;
}>();

const emit = defineEmits(['close', 'preview', 'remove', 'cancel']);

const data = toRef(props, 'data');
const attachments = computed(() => data.value?.list || []);
const limitMBSize = computed(() => bytes2MB(data.value?.sizeLimit));
const limitCount = computed(() => data.value?.amountLimit || 10);
const limitTotalMBSize = computed(() => bytes2MB(data.value?.totalSpace));
const totalUsedMBSize = computed(() => bytes2MB(data.value?.usedSpaceUsed, 1));
const percentage = computed(() => {
  if (!data.value) {
    return 0;
  }
  const { usedSpaceUsed, totalSpace } = data.value;

  return (usedSpaceUsed * 100) / totalSpace;
});

const onClose = () => emit('close');
const onPreview = (item: AttachmentInfo) => emit('preview', item);
const onRemove = (item: AttachmentInfo) => emit('remove', item);
const onCancel = () => emit('cancel');
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
}
</style>
