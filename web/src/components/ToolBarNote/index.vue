<template>
  <!-- <CollapsedBar
    v-if="isOwner && showOriginalPDF"
    class="toolbar-note text-base leading-4"
    :collapsed="!toolbarSettings.toolBarNoteVisible"
    @toggled="
      (v) =>
        setToolbarSettings({
          toolBarNoteVisible: !v,
        })
    "
  >
    <template #icon-expand>
      <DoubleLeftOutlined />
    </template>
    <ScreenShot
      :clip-selecting="clipSelecting"
      :cancel-clip="cancelClip"
      :handle-screen-shot="handleScreenShot"
    />
    <span
      v-if="isSelfNoteTab"
      :style="{
        backgroundColor:
          shapeState.creating.value === ShapeType.rectangle ? '#414548' : '',
      }"
      @click="createRect()"
    >
      <a-tooltip
        placement="left"
        :title="$t('toolbar.rectangle')"
      >
        <i
          class="aiknowledge-icon icon-shape-rect"
          aria-hidden="true"
        />
      </a-tooltip>
    </span>
    <span
      v-if="isSelfNoteTab"
      :style="{
        backgroundColor:
          shapeState.creating.value === ShapeType.ellipse ? '#414548' : '',
      }"
      @click="createEllipse()"
    >
      <a-tooltip
        placement="left"
        :title="$t('toolbar.circle')"
      >
        <img
          src="@/assets/images/shape-circle.svg"
          style="height: 1rem; width: auto"
        >
      </a-tooltip>
    </span>
    <span
      v-if="isSelfNoteTab"
      :style="{
        backgroundColor:
          shapeState.creating.value === ShapeType.arrow ? '#414548' : '',
      }"
      @click="createShape(ShapeType.arrow)"
    >
      <a-tooltip
        placement="left"
        :title="$t('toolbar.arrow')"
      >
        <i
          class="aiknowledge-icon icon-shape-arrow"
          aria-hidden="true"
        />
      </a-tooltip>
    </span>
    <span
      v-if="isSelfNoteTab"
      :style="{
        backgroundColor: shapeState.creating.value === 'text' ? '#414548' : '',
      }"
      @click="createShape('text')"
    >
      <a-tooltip
        placement="left"
        :title="$t('toolbar.text')"
      >
        <i
          class="aiknowledge-icon icon-shape-text"
          aria-hidden="true"
        />
      </a-tooltip>
    </span>
    <FontSizeOutlined v-if="false" />
  </CollapsedBar> -->
</template>

<script setup lang="ts">
import { DoubleLeftOutlined, FontSizeOutlined } from '@ant-design/icons-vue';
import CollapsedBar from '@/components/Common/CollapsedBar.vue';
import { currentGroupId, currentNoteInfo, isOwner } from '~/src/store';
import { usePDFWEbviewPreviewMode } from '~/src/hooks/usePDFWebviewScrollMode';
import { useToolBarSettings } from '~/src/hooks/UserSettings/useToolBarSettings';
import ScreenShot from './components/ScreenShot.vue';
import { cancelShape, shapeState } from '../Main/mouseCore';
import { ShapeType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { useI18n } from 'vue-i18n';
import {
  executeAndSetBeforeChangeTab,
  removeBeforeChangeTab,
} from '../../hooks/UserSettings/useSideTabSettings';
import { ElementClick, reportClick } from '~/src/api/report';
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import {
  ANNOTATION_SCREENSHOT,
  ScreenShotPayload,
  useClip,
} from '~/src/hooks/useHeaderScreenShot';
import { AnnotationRect } from '~/src/stores/annotationStore/BaseAnnotationController';
import { usePdfStore } from '~/src/stores/pdfStore';

const props = defineProps<{
  clipSelecting: boolean;
  clipAction: ReturnType<typeof useClip>['clipAction'];
}>();

const pdfStore = usePdfStore();
const pdfAnnotater = computed(() => {
  return pdfStore.getAnnotater(currentNoteInfo.value?.noteId);
});

const { showOriginalPDF } = usePDFWEbviewPreviewMode();
const { toolbarSettings, setToolbarSettings } = useToolBarSettings();

const isSelfNoteTab = computed(() => {
  return currentGroupId.value === SELF_NOTEINFO_GROUPID;
});

const { t } = useI18n();
shapeState.deleteTitle = t('message.confirmToDeleteNoteTip');
shapeState.deleteOk = t('viewer.delete');

const createRect = () => {
  createShape(ShapeType.rectangle);
};

const createEllipse = () => {
  createShape(ShapeType.ellipse);
};

const createShape = (type: ShapeType | 'text') => {
  if (props.clipSelecting) {
    cancelScreenShot();
  }

  if (shapeState.creating.value === type) {
    cancelShape();
  } else if (shapeState.creating.value === ShapeType.UNRECOGNIZED) {
    shapeState.creating.value = type;
  } else {
    cancelShape();
    shapeState.creating.value = type;
  }
};

const cancelClip = () => {
  props.clipAction.cancelCut();
  removeBeforeChangeTab(cancelClip);
};

const handleScreenShot = () => {
  if (shapeState.creating.value !== ShapeType.UNRECOGNIZED) {
    cancelShape();
  }

  if (props.clipSelecting) {
    cancelScreenShot();
  } else {
    startScreehShot();
  }
};

const startScreehShot = async () => {
  executeAndSetBeforeChangeTab(cancelClip);

  await nextTick();

  reportClick(ElementClick.screenshots, 'on');

  const payload = await props.clipAction.init(true);

  if (!payload) {
    return;
  }

  const { newRectElement, newRectAnnotation } = payload;

  pdfAnnotater.value?.UI.emit(ANNOTATION_SCREENSHOT, {
    rect: newRectElement?.getBoundingClientRect(),
    pageNum: (newRectAnnotation as AnnotationRect).pageNumber,
  } as ScreenShotPayload);
};

const cancelScreenShot = () => {
  cancelClip();
  reportClick(ElementClick.screenshots, 'off');
};

const bindEsc = (event: KeyboardEvent) => {
  if (event.key !== 'Escape') {
    return;
  }

  if (props.clipSelecting) {
    cancelScreenShot();
  } else if (shapeState.creating.value !== ShapeType.UNRECOGNIZED) {
    cancelShape();
  }
};

onMounted(() => {
  document.addEventListener('keydown', bindEsc, { passive: true });
});
onUnmounted(() => {
  document.removeEventListener('keydown', bindEsc);
});
</script>

<style lang="postcss">
.toolbar-note.collapsed-icon,
.toolbar-note.collapsed-bar {
  bottom: 30px;

  background: #2f3337;
}

.toolbar-note.collapsed-icon {
  right: 2px;
  .anticon-double-left {
    transform: rotate(90deg);
  }
}

.toolbar-note.collapsed-bar {
  flex-direction: column;
  left: auto;
  width: fit-content;
  right: 12px;
  padding: 8px 4px;

  .btn-collapse {
    color: inherit;
    transform: rotate(-90deg);
    order: 99;
  }

  & > * {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    padding: 10px;
    line-height: 1;
    border-radius: 2px;
    cursor: pointer;

    &:hover {
      background-color: theme('colors.rp-dark-4');
    }
  }
}
</style>
