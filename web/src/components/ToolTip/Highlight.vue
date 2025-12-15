<template>
  <div class="highlight-container">
    <div
      class="item"
      @click.stop="addNote"
    >
      <i
        class="aiknowledge-icon icon-mark"
        :style="{ background: color || styleMap[commentState.styleId].color }"
        aria-hidden="true"
      />
      <i
        class="aiknowledge-icon icon-arrow-down"
        aria-hidden="true"
      />
    </div>

    <div class="color-options">
      <Color
        v-if="!changing"
        :color="styleMap[commentState.styleId].color"
        :handle-outside="handleOutside"
        @change="handleColorChange"
        @enterChange="handleEnterChange"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { RectOptions } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { ToolBarType, PDF_ANNOTATE_ID } from '@idea/pdf-annotate-core';
import { useCommentGlobalState } from '~/src/hooks/useNoteState';
import { currentGroupId, currentNoteInfo, store } from '~/src/store';
import { rectStyleMap, styleMap } from '@/style/select';
import Color from '@/components/Common/Color.vue';
import hexRgb from 'hex-rgb';
import { connect } from '~/src/dom/arrow';
import { connectComment } from '~/src/util/scroll';
import { useRightSideTabSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { RightSideBarType } from '../Right/TabPanel/type';
import { usePdfStore } from '~/src/stores/pdfStore';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { AnnotationSelect } from '~/src/stores/annotationStore/BaseAnnotationController';
import { computed, ref } from 'vue';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import { emitter, ANNOTATION_CREATED } from '~/src/util/eventbus';

const emit = defineEmits<{
  (event: 'createHighlight', uuid: string): void;
  (event: 'hideToolTip'): void;
}>();

const props = defineProps({
  toolTipType: {
    type: [String, Number],
    default: ToolBarType.Highlight,
  },
  pdfId: {
    type: String,
    default: '',
  },
  noteId: {
    type: String,
    default: '',
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
    default: '',
  },
  pageNumber: {
    type: Number,
    default: 1,
  },
  color: {
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
const pdfAnnotaterRef = computed(() => {
  return pdfStore.getAnnotater(props.noteId);
});

const commentState = useCommentGlobalState();
const { checkNoteTabVisible, checkSwitchGroupTab } = useRightSideTabSettings();

const addNote = async () => {
  if (props.annotateId) {
    return;
  }

  emit('hideToolTip');

  changing.value = true;

  if (props.toolTipType === ToolBarType.rect) {
    await props.createRect();
  } else {
    await addSelectNote();
  }

  changing.value = false;
};

const addSelectNote = async () => {
  const styleId = commentState.value.styleId;

  pdfViewerRef.value?.getDocumentViewer().clearSelection();

  let uuid: string;

  console.warn(props.rects);

  checkSwitchGroupTab();

  const isHighlight = currentGroupId.value === SELF_NOTEINFO_GROUPID;

  try {
    uuid = (await annotationStore.controller.addAnnotation({
      pageNumber: props.pageNumber,
      ...({
        rectStr: props.rectStr,
      } as Partial<AnnotationSelect>),
      rectangles: props.rects,
      type: ToolBarType.select,
      documentId: currentNoteInfo.value?.noteId,
      idea: '',
      styleId: +styleId,
      ...styleMap[styleId],
      isHighlight,
      commentatorInfoView: {
        avatarCdnUrl: store.state.user.userInfo?.avatarUrl ?? '',
        nickName: store.state.user.userInfo?.nickName ?? '',
        userId: store.state.user.userInfo?.id ?? '',
      },
      deleteAuthority: true,
    })) as string;
  } catch (error) {
    console.error(error);
    return;
  }

  emit('createHighlight', uuid);
  emitter.emit(ANNOTATION_CREATED);

  if (!isHighlight && checkNoteTabVisible()) {
    connectComment(uuid);
  }
};

const changing = ref(false);

const handleColorChange = async (styleId: keyof typeof styleMap) => {
  commentState.value.styleId = styleId;

  if (!props.annotateId) {
    addNote();
    return;
  }

  changing.value = true;

  const params: any = {
    ...(props.toolTipType === ToolBarType.select
      ? styleMap[styleId]
      : rectStyleMap[styleId]),
  };

  if (props.toolTipType === ToolBarType.rect) {
    params.stroke = params.color;
  }

  const patch = {
    styleId: +styleId,
    type: props.toolTipType,
    ...params,
  };

  annotationStore.controller.localPatchAnnotation(
    props.annotateId as string,
    props.pageNumber,
    patch
  );

  try {
    await annotationStore.controller.patchAnnotation(props.annotateId, patch);
  } finally {
    changing.value = false;
  }

  if (!props.isHighlightAnnotation) {
    connect(props.annotateId, true);
  }
};

const changeColor = (styleId: keyof typeof styleMap) => {
  if (props.toolTipType === ToolBarType.select) {
    const params = {
      ...styleMap[styleId],
    };

    pdfAnnotaterRef.value?.UI.commentController.setCommentToDom(
      props.annotateId as string,
      params
    );
  } else {
    const params = {
      ...rectStyleMap[styleId],
      stroke: rectStyleMap[styleId].color,
    };

    pdfAnnotaterRef.value?.UI.rectController.setRectToDom(
      props.annotateId as string,
      params
    );
  }

  if (!props.isHighlightAnnotation) {
    connect(props.annotateId as string, true);
  }
};

const handleEnterChange = async (styleId: keyof typeof styleMap) => {
  if (props.annotateId) {
    changeColor(styleId);
    return;
  }

  pdfViewerRef.value
    ?.getDocumentViewer()
    .updateSelectionColor(
      hexRgb(styleMap[styleId].color, { format: 'css', alpha: 0.3 })
    );
};

const handleOutside = () => {
  const docViewer = pdfViewerRef.value?.getDocumentViewer();
  if (props.annotateId) {
    const container = docViewer?.container!;

    const g = container.querySelector(
      `[${PDF_ANNOTATE_ID}="${props.annotateId}"]`
    )!;

    const styleId = g.getAttribute(
      'style-id'
    ) as unknown as keyof typeof styleMap;

    changeColor(styleId);
    return;
  }

  docViewer?.updateSelectionColor();
};
</script>

<style lang="less" scoped>
.highlight-container {
  position: relative;

  &:hover {
    .color-options {
      display: block;
    }
  }
}
.color-options {
  background: #393c3e;
  position: absolute;
  left: 0;
  bottom: 0;
  transform: translate(0, 100%);
  align-items: center;
  width: 100%;
  display: none;
}

.item {
  width: 32px;
  height: 32px;
  font-size: 14px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;

  &:hover {
    background: #52565a;

    .icon-arrow-down {
      display: inherit;
    }
  }

  .aiknowledge-icon {
    font-size: 14px;
    width: 16px;
    height: 16px;
    display: flex;
    justify-content: center;
    align-items: center;
    font-weight: 100;
  }

  .icon-mark {
    border-radius: 2px;
  }

  .icon-arrow-down {
    color: #ffffff;
    font-size: 6px;
    height: 2px;
    height: 6px;
    margin-top: 4px;
    position: absolute;
    bottom: 0;
    left: 50%;
    transform: translate(-50%, 0) scale(0.6);
    display: none;
  }
}
</style>
