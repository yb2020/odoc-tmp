<template>
  <div
    :class="[PDF_ANNOTATE_EDIT_OVERLAY]"
    :style="{
      top: top + 'px',
      left: left + 'px',
      width: width + 'px',
      height: height + 'px',
      ...style,
    }"
    @click.stop
  >
    <p
      v-if="annotation && 'rectStr' in annotation"
      ref="text"
      class="m-0 opacity-0"
    >
      {{ annotation.rectStr }}
    </p>
    <div
      v-if="
        isOwner && annotationStore.activeOverlayPageNumber === rectPageNumber
      "
      class="close"
      @mousedown.stop="handleClick"
    >
      <close-outlined class="iconguanbi" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { Modal, message } from 'ant-design-vue';
import { clearArrow } from '@/dom/arrow';
import { isOwner } from '~/src/store';
import { CloseOutlined } from '@ant-design/icons-vue';

import { RightSideBarType } from '~/src/components/Right/TabPanel/type';
import { useRightSideTabSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { PDF_ANNOTATE_EDIT_OVERLAY } from '~/src/constants';
import { useI18n } from 'vue-i18n';

const annotationStore = useAnnotationStore();

const props = defineProps<{
  annotatePageNumber: number;
  rectPageNumber: number;
  top: number;
  left: number;
  width: number;
  height: number;
  annotateId: string;
  style: Object;
}>();

const { activeTab } = useRightSideTabSettings();

const { t } = useI18n();
const text = ref<HTMLElement>();

const annotation = computed(() =>
  annotationStore.pageMap[props.annotatePageNumber]?.find(
    (item) => item.uuid === props.annotateId
  )
);

const handleClick = async () => {
  if (activeTab.value === RightSideBarType.Group) {
    const { deleteAuthority } = annotation.value || {};

    if (!deleteAuthority) {
      message.error('您无权限删除该标记~');
      return;
    }
  }

  const onOk = async () => {
    await annotationStore.controller.deleteAnnotation(
      props.annotateId,
      Number(props.annotatePageNumber)
    );

    clearArrow();

    annotationStore.delHoverNote(props.annotateId);
  };

  if (annotation.value?.isHighlight) {
    document.body.click();
    return onOk();
  }

  Modal.confirm({
    title: t('message.confirmToDeleteAnnotationTip'),
    onOk,
    maskClosable: true,
    okButtonProps: {
      danger: true,
    },
    cancelButtonProps: { type: 'primary' },
    okText: t('viewer.delete'),
  });
};

onMounted(() => {
  setTimeout(() => {
    if (text.value?.offsetWidth) {
      const range = document.createRange();
      range.selectNode(text.value);
      const selection = window.getSelection();
      selection?.removeAllRanges();
      selection?.addRange(range);
    }
  }, 60);
});
</script>

<style lang="less" scoped>
.pdf-annotate-edit-overlay {
  box-sizing: content-box;
  position: absolute;
  border: 2px solid #00bfff;
  border-radius: 1px;
  z-index: 99;

  .close {
    position: absolute;
    right: 0;
    top: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    transform: translate(100%, -3px);
    cursor: pointer;

    width: 20px;
    height: 20px;
    background: #dadfe6;
    box-shadow: 0px 2px 4px 0px rgba(0, 0, 0, 0.2);
    border-radius: 2px;
    border: 1px solid #d3d6d8;

    .iconguanbi {
      font-size: 12px;
      color: #000000;
      // transform: scale(0.5);
    }
  }
}
</style>
