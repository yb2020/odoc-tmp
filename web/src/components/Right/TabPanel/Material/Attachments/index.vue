<template>
  <a-popover
    :visible="visible"
    overlayClassName="attachments-popup"
    trigger="hover"
    placement="left"
  >
    <template #content>
      <List
        v-if="data"
        :data="data"
        :isCollected="isCollected"
        :uploading="uploading"
        :uploadFile="uploadFile"
        :removingIds="removingIds"
        :onUpload="onUpload"
        :onBeforeUpload="onBeforeUpload"
        @close="onToggle"
        @preview="onPreview"
        @remove="onRemove"
        @cancel="onCancel"
      />
    </template>
    <!-- <div class="attachments flex">
      <span class="label pr-2 text-[13px] opacity-[0.45]">{{
        $t('info.attachments')
      }}</span>
      <div class="flex-1 min-w-0">
        <a-upload
          v-if="!count && !uploading"
          name="file"
          :maxCount="1"
          :showUploadList="false"
          :customRequest="onUpload"
        >
          <a-button
            class="!flex items-center !text-[13px] !h-6 !leading-4 !p-0 !border-0"
            :class="!isCollected ? '!text-[var(--site-theme-text-secondary)] !opacity-40' : ''"
            type="link"
            @click="onBeforeUpload"
          >
            <template #icon>
              <UploadOutlined />
            </template>
            {{ $t('info.uploadAttachments') }}
          </a-button>
        </a-upload>
        <p
          v-if="count > 0 || uploading"
          class="flex items-center text-[13px] h-6 pl-2 mb-0"
        >
          <span class="flex-1 text-ellipsis whitespace-nowrap overflow-hidden">
            <span v-if="uploading && count <= 1">{{ uploadFile?.name }}</span>
            <a
              v-else
              :title="attachments[0].name"
              @click="onPreview(attachments[0])"
            >{{ attachments[0].name }}</a>
          </span>
          <span class="flex items-center text-base">
            <LoadingOutlined
              v-if="uploading || removingIds.includes(attachments[0].id)"
              class="ml-2"
            />
            <a-upload
              v-else-if="count === 1"
              name="file"
              :maxCount="1"
              :customRequest="onUpload"
              @click="onBeforeUpload"
            >
              <UploadOutlined
                :class="!isCollected ? '!opacity-20' : ''"
                class="ml-2 cursor-pointer opacity-40 hover:opacity-80"
              />
            </a-upload>

            <CloseOutlined
              v-if="uploading"
              class="ml-2"
              @click="onCancel"
            />
            <DeleteOutlined
              v-else-if="count === 1"
              class="ml-2 cursor-pointer opacity-40 hover:opacity-80"
              @click="() => onRemove()"
            />
          </span>
        </p>
        <p
          v-if="count > 1"
          class="mt-3 mb-0 pl-2 text-[13px] leading-5 text-[rgba(255,255,255,.65)] cursor-pointer"
          @click="onToggle"
        >
          &lt; {{ $t('info.attachmentTotal', { count }) }}
        </p>
      </div>
    </div> -->
  </a-popover>
</template>

<script setup lang="ts">
import {
  LoadingOutlined,
  CloseOutlined,
  DeleteOutlined,
  UploadOutlined,
} from '@ant-design/icons-vue';
import { computed, ref, toRef, watch } from 'vue';
import { message } from 'ant-design-vue';
import List from './list.vue';
import { AttachmentInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/DocAttachment';
import $axios from '@/api/axios';
import { useDocAttachments } from '@/hooks/useDocAttachments';
import { useStore } from '@/store';
import { SELF_NOTEINFO_GROUPID } from '@/store/base/type';
import { isInElectron } from '~/src/util/env';
import { downloadUrl } from '~/src/common/src/utils/url';
import {
  ElementClick,
  getPdfIdFromUrl,
  PageType,
  reportElementClick,
} from '~/src/api/report';

const props = defineProps<{
  docId: string;
  calculateFileMD5?: (file: File) => Promise<string>;
}>();

const docId = toRef(props, 'docId');
const store = useStore();
const isCollected = computed(
  () => !!store.state.base.noteInfoMap[SELF_NOTEINFO_GROUPID].isCollected
);
const isReady = computed(() => !!docId.value);

const {
  data,
  uploadFile,
  uploading,
  removingIds,
  refresh,
  doUpload,
  cancelUpload,
  doRemove,
} = useDocAttachments($axios, docId, isReady, props.calculateFileMD5);

watch(isCollected, refresh);

const pesist = ref<undefined | true>(undefined);
const visible = computed(() =>
  (data.value?.list.length || 0) <= 1 ? false : pesist.value
);
const attachments = computed(() => data.value?.list || []);
const count = computed(() => data.value?.list.length || 0);

const onToggle = () => {
  pesist.value = pesist.value ? undefined : true;
};

const checkUpload = () => {
  if (!isCollected.value) {
    message.info('请先收藏该文献，再点击上传附件');
  }
  return isCollected.value;
};

const onBeforeUpload = (e: MouseEvent) => {
  if (!checkUpload()) {
    e.preventDefault();
    e.stopPropagation();
  }
};

const onUpload = async ({ file }: { file: File; filename: string }) => {
  if (!checkUpload()) {
    return;
  }

  await doUpload(file);
};

const onCancel = () => {
  cancelUpload();
};

const onPreview = (item: AttachmentInfo) => {
  reportElementClick({
    page_type: PageType.note,
    type_parameter: getPdfIdFromUrl(),
    element_name: ElementClick.attachments,
    element_parameter: item.id,
  });

  if (isInElectron() && !item.isPreview) {
    downloadUrl(item.url, item.name);
    return;
  }
  window.open(item.url);
};

const onRemove = (item: AttachmentInfo = attachments.value[0]) => {
  doRemove({
    attachmentId: item.id,
  });
};
</script>

<style scoped lang="less">
.label {
  color: var(--site-theme-text-secondary);
}
.attachments {
  .ant-upload {
    font-size: 1rem;
    line-height: 1;
    display: flex !important;
  }
}
.attachments-popup {
  .ant-popover-inner {
    background: var(--site-theme-bg-light);
  }
  .ant-popover-inner-content {
    padding: 0 0 8px;
    color: var(--site-theme-text-primary);
  }
}
</style>
