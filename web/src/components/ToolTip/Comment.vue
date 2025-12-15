<template>
  <div class="comment-container">
    <div
      class="item"
      @click.stop="handleComment"
    >
      <i
        class="aiknowledge-icon icon-take-note"
        aria-hidden="true"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { RectOptions } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { ViewerEvent } from '@idea/pdf-annotate-viewer';
import { ToolBarType } from '@idea/pdf-annotate-core';
import { store, currentNoteInfo } from '@/store';

import { message } from 'ant-design-vue';
import { useCommentGlobalState } from '@/hooks/useNoteState';
import { styleMap } from '~/src/style/select';
import { useRightSideTabSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { RightSideBarType } from '../Right/TabPanel/type';
import { computed, nextTick, onMounted, onUnmounted } from 'vue';
import { insertNoReferenceAnnotation } from '../Right/TabPanel/Note/annotation-state';
import { usePdfStore } from '~/src/stores/pdfStore';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { AnnotationSelect } from '~/src/stores/annotationStore/BaseAnnotationController';

const props = defineProps({
  toolTipType: {
    type: [Number, String],
    default: ToolBarType.select,
  },
  annotateId: {
    type: String,
  },
  rects: {
    type: Array as () => RectOptions[],
    default: [],
  },
  rectStr: {
    type: String,
  },
  pageNumber: { type: Number, default: 1 },
  pdfId: {
    type: String,
    default: '',
  },
  isHighlightAnnotation: {
    type: Boolean,
    required: true,
  },
  createRect: {
    type: Function as unknown as () => () => Promise<string>,
    required: true,
  },
});

const pdfStore = usePdfStore();
const annotationStore = useAnnotationStore();

const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(props.pdfId);
});
const emit = defineEmits<{
  (event: 'hideToolTip'): void;
}>();

const commentState = useCommentGlobalState();

const { switchAndShowTab, activeTab, checkSwitchGroupTab } =
  useRightSideTabSettings();

const handleComment = async () => {
  if (
    JSON.stringify(props.rects || []).length > 40000 ||
    props.rectStr!.length > 2000
  ) {
    message.error('单个笔记字数最多2000字！');
    return;
  }

  pdfViewerRef.value?.getDocumentViewer().clearSelection();

  checkSwitchGroupTab();

  if (activeTab.value === RightSideBarType.Group) {
    switchAndShowTab(RightSideBarType.Group, {});
  }

  emit('hideToolTip');

  console.log('toolTipType', props.toolTipType, ToolBarType.select);

  if (props.toolTipType === ToolBarType.select) {
    await addSelectComment();
  } else if (props.annotateId) {
    popupHoverNote(props.annotateId);
  } else {
    const uuid = await props.createRect();
    await nextTick();
    annotationStore.currentAnnotationId = uuid;
  }
};

const addSelectComment = async () => {
  if (props.annotateId) {
    if (props.isHighlightAnnotation) {
      await annotationStore.controller.patchAnnotation(props.annotateId, {
        isHighlight: false,
      });
    }
    popupHoverNote(props.annotateId);
    return;
  }

  const uuid = (await annotationStore.controller.addAnnotation({
    pageNumber: props.pageNumber,
    ...({
      rectStr: props.rectStr,
    } as Partial<AnnotationSelect>),
    rectangles: props.rects,
    type: ToolBarType.select,
    documentId: currentNoteInfo.value?.noteId,
    idea: '',
    tags: [],
    styleId: +commentState.value.styleId,
    ...styleMap[commentState.value.styleId],
    commentatorInfoView: {
      avatarCdnUrl: store.state.user.userInfo?.avatarUrl ?? '',
      nickName: store.state.user.userInfo?.nickName ?? '',
      userId: store.state.user.userInfo?.id ?? '',
    },
    deleteAuthority: true,
  }))!;

  if (
    !annotationStore.activeColorMap.ref ||
    !annotationStore.activeColorMap[commentState.value.styleId]
  ) {
    if (!annotationStore.activeColorMap.ref) {
      annotationStore.activeColorMap.ref = true;
    }

    if (!annotationStore.activeColorMap[commentState.value.styleId]) {
      annotationStore.activeColorMap[commentState.value.styleId] = true;
    }
  }

  popupHoverNote(uuid);
};

const popupHoverNote = async (uuid: string) => {
  if (activeTab.value !== RightSideBarType.Group) {
    // 跨页笔记时应该取笔记数据第一个选段的页码
    // props.pageNumber是最后一个选段的页码
    const { pageNumber: page } =
      annotationStore.controller.findAnnotation(uuid);
    annotationStore.showHoverNote({
      uuid,
      page,
    });
  }

  await nextTick();
  annotationStore.currentAnnotationId = uuid;
};

const addEmptyComment = ({ pageNumber }: { pageNumber: number }) => {
  if (
    activeTab.value !== RightSideBarType.Group &&
    !annotationStore.currentAnnotationId
  ) {
    insertNoReferenceAnnotation(pageNumber, 0);
  }
};

onMounted(() => {
  pdfViewerRef.value?.addEventListener(
    ViewerEvent.EMPTY_CLICK,
    addEmptyComment
  );
});
onUnmounted(() => {
  pdfViewerRef.value?.removeEventListener(
    ViewerEvent.EMPTY_CLICK,
    addEmptyComment
  );
});
</script>

<style lang="less" scoped>
.color-options {
  background: #393c3e;
  position: absolute;
  padding: 6px 9px;
  left: 0;
  top: 0;
  transform: translate(0, -100%);
  height: 32px;
  display: flex;
  align-items: center;
  width: 100%;
}

.item {
  width: 32px;
  height: 32px;
  cursor: pointer;

  display: flex;
  align-items: center;
  justify-content: center;

  .aiknowledge-icon {
    font-size: 16px;
    font-weight: 100;
  }

  &:hover {
    background: #52565a;
  }
}
</style>
