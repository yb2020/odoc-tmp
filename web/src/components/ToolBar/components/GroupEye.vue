<template>
  <a-tooltip overlayClassName="eye_tips">
    <template #title>
      <div class="box_list">
        <div class="box_item">
          <a-checkbox
            :class="{ check: annotationStore.groupOtherVisible }"
            :checked="annotationStore.groupOtherVisible"
            @change="onGroupOtherChange"
          >
            {{ $t('viewer.teamAnnotation') }}
          </a-checkbox>
        </div>
        <div class="box_item">
          <a-checkbox
            :class="{ check: annotationStore.groupSelfVisible }"
            :checked="annotationStore.groupSelfVisible"
            @change="onGroupMineChange"
          >
            {{ $t('viewer.myAnnotation') }}
          </a-checkbox>
        </div>
        <div class="box_item">
          <a-checkbox
            :class="{ check: annotationStore.groupImageVisible }"
            :checked="annotationStore.groupImageVisible"
            @change="onGroupImageChange"
          >
            {{ $t('viewer.showAvatar') }}
          </a-checkbox>
        </div>
      </div>
    </template>
    <div
      class="eye-container"
      @click="handleEye"
    >
      <eye-outlined
        v-if="
          annotationStore.groupOtherVisible || annotationStore.groupSelfVisible
        "
        class="icon"
      />
      <eye-invisible-outlined
        v-else
        class="icon"
      />
    </div>
  </a-tooltip>
</template>

<script lang="ts" setup>
import { computed, watch } from 'vue';
import { EyeOutlined, EyeInvisibleOutlined } from '@ant-design/icons-vue';
import { setDisplay } from '@idea/pdf-annotate-core';

import { useStore } from '@/store';

import { toggleAvatar } from '~/src/dom/avatar';
import { setAnnotateThumbnails } from '~/src/dom/thumbnails';

import { usePdfStore } from '~/src/stores/pdfStore';
import { useAnnotationStore } from '~/src/stores/annotationStore';

const props = defineProps<{
  pdfId: string;
  noteId?: string;
}>();

const pdfStore = usePdfStore();
const annotationStore = useAnnotationStore();

const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(props.pdfId);
});
const pdfAnnotaterRef = computed(() => {
  return pdfStore.getAnnotater(props.noteId ?? '');
});

const isHidden = computed(
  () => !annotationStore.groupOtherVisible && !annotationStore.groupSelfVisible
);

watch(isHidden, () => {
  if (!isHidden.value) {
    setAnnotateThumbnails(pdfViewerRef.value, pdfAnnotaterRef.value);
  }
});

const handleEye = () => {
  if (annotationStore.groupOtherVisible === annotationStore.groupSelfVisible) {
    annotationStore.groupOtherVisible = !annotationStore.groupOtherVisible;
    annotationStore.groupSelfVisible = !annotationStore.groupSelfVisible;
    annotationStore.groupImageVisible = annotationStore.groupSelfVisible;
  } else {
    annotationStore.groupOtherVisible = true;
    annotationStore.groupSelfVisible = true;
    annotationStore.groupImageVisible = true;
  }

  handleIsHiddenChange();
};

const store = useStore();
const userId = computed(() => store.state.user.userInfo?.id);

const handleIsHiddenChange = () => {
  toggleAvatar(
    userId.value as string,
    annotationStore.groupOtherVisible,
    annotationStore.groupSelfVisible,
    annotationStore.groupImageVisible
  );
};

const onGroupMineChange = () => {
  annotationStore.groupSelfVisible = !annotationStore.groupSelfVisible;

  handleIsHiddenChange();
};
const onGroupOtherChange = () => {
  annotationStore.groupOtherVisible = !annotationStore.groupOtherVisible;

  handleIsHiddenChange();
};
const onGroupImageChange = () => {
  annotationStore.groupImageVisible = !annotationStore.groupImageVisible;

  handleIsHiddenChange();
};
</script>

<style lang="less" scoped>
.eye-container {
  width: 36px;
  height: 36px;
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;

  &:hover {
    background: #f5f5f5;
  }

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
