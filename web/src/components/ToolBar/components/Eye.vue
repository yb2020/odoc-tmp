<template>
  <a-tooltip
    v-if="envStore.viewerConfig.otherAnnotation !== false"
    overlayClassName="eye_tips"
    :overlayStyle="isOwner ? {} : { display: 'none' }"
  >
    <template #title>
      <div class="box_list">
        <div class="box_item">
          <a-checkbox
            v-if="isHotDisable"
            :class="{ check: annotationStore.hotVisible }"
            :checked="annotationStore.hotVisible"
            @change="onHotChange"
          >
            {{ $t('viewer.otherAnnotation') }}
          </a-checkbox>
          <a-tooltip
            v-else
            placement="topLeft"
          >
            <template #title>
              {{ $t('message.annotionsTip3') }}
            </template>
            <a-checkbox
              :checked="annotationStore.hotVisible"
              disabled
            >
              {{ $t('viewer.otherAnnotation') }}
            </a-checkbox>
          </a-tooltip>
          <a-tooltip placement="topLeft">
            <template #title>
              ①{{ $t('message.annotionsTip1') }}<br>
              ②{{ $t('message.annotionsTip2') }}
            </template>
            <question-circle-outlined />
          </a-tooltip>
        </div>
        <div class="box_item">
          <a-checkbox
            :class="{ check: annotationStore.personVisible }"
            :checked="annotationStore.personVisible"
            @change="onMineChange"
          >
            {{ $t('viewer.myAnnotation') }}
          </a-checkbox>
        </div>
      </div>
    </template>
    <div
      class="eye-container"
      @click="isOwner ? clickEye() : onMineChange()"
    >
      <eye-outlined
        v-if="isVisible"
        class="icon"
      />
      <eye-invisible-outlined
        v-else
        class="icon"
      />
    </div>
  </a-tooltip>
  <div
    v-else
    class="eye-container"
    @click="isOwner ? clickEye() : onMineChange()"
  >
    <eye-outlined
      v-if="isVisible"
      class="icon"
    />
    <eye-invisible-outlined
      v-else
      class="icon"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, watch } from 'vue';
import { isOwner } from '@/store';
import {
  EyeOutlined,
  EyeInvisibleOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons-vue';

import { setDisplay } from '@idea/pdf-annotate-core';

import { hideAvatar, showAvatar } from '~/src/dom/avatar';
import { setAnnotateThumbnails } from '~/src/dom/thumbnails';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { useEnvStore } from '~/src/stores/envStore';
import { usePdfStore } from '~/src/stores/pdfStore';

const props = defineProps<{
  pdfId: string;
  noteId?: string;
  currentPage: number;
}>();

const pdfStore = usePdfStore();
const annotationStore = useAnnotationStore();

const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(props.pdfId);
});
const pdfAnnotaterRef = computed(() => {
  return pdfStore.getAnnotater(props.noteId ?? '');
});

const isHotDisable = computed(
  () => annotationStore.pageHotMap[props.currentPage]?.length > 0
);

const isVisible = computed(() => {
  return isOwner.value
    ? annotationStore.personVisible || annotationStore.hotVisible
    : annotationStore.personVisible;
});

watch(isVisible, (value) => {
  if (isVisible.value) {
    setAnnotateThumbnails(pdfViewerRef.value, pdfAnnotaterRef.value);
  }
});

const clickEye = () => {
  if (annotationStore.hotVisible === annotationStore.personVisible) {
    annotationStore.hotVisible = !annotationStore.hotVisible;
    annotationStore.personVisible = !annotationStore.personVisible;
  } else {
    annotationStore.hotVisible = true;
    annotationStore.personVisible = true;
  }
};

const onMineChange = () => {
  annotationStore.personVisible = !annotationStore.personVisible;
};
const onHotChange = () => {
  annotationStore.hotVisible = !annotationStore.hotVisible;
};

const envStore = useEnvStore();
</script>

<style lang="less" scoped>
.eye-container {
  cursor: pointer;

  .icon {
    font-size: 16px;
    color: rgba(0, 0, 0, 64%);
  }
}
</style>

<style lang="less">
.eye_tips {
  .ant-tooltip-arrow {
    display: none !important;
  }
  .ant-tooltip-inner {
    background: #fff;
    .ant-checkbox-wrapper {
      color: rgba(0, 0, 0, 85%);
      &.check {
        color: #1890ff;
        .ant-checkbox-inner {
          border-color: #1f71e0 !important;
        }
      }
      &.ant-checkbox-wrapper-disabled {
        span {
          color: rgba(0, 0, 0, 0.25);
        }
        .ant-checkbox-inner {
          background: #f5f5f5;
          border-color: #d9d9d9 !important;
          &:after {
            border-color: #00000040;
          }
        }
      }
      .ant-checkbox-inner {
        border-color: #d9d9d9;
      }
    }
  }
  .box_list {
    user-select: none;
    .box_item {
      padding: 4px;
      color: #1d2229;
      .anticon {
        cursor: pointer;
      }
    }
  }
}
</style>
