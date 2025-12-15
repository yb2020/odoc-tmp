<template>
  <CollapsedBar class="toolbar-hd" :collapsed="!toolbarSettings.toolBarHeadVisible" @toggled="
      (v) =>
        setToolbarSettings({
          toolBarHeadVisible: !v,
        })
    ">
    <!-- Components from ToolBarNote -->
    <ScreenShot v-if="isOwner && showOriginalPDF" :clip-selecting="clipSelecting" :cancel-clip="cancelClip"
      :handle-screen-shot="handleScreenShot" />
    <span v-if="isOwner && showOriginalPDF && isSelfNote" :style="{
        backgroundColor:
          shapeState.creating.value === ShapeType.rectangle ? getActiveBackgroundColor() : '',
      }" @click="createRect()">
      <a-tooltip placement="top" :title="$t('toolbar.rectangle')">
        <i class="aiknowledge-icon icon-shape-rect" aria-hidden="true" />
      </a-tooltip>
    </span>
    <span v-if="isOwner && showOriginalPDF && isSelfNote" :style="{
        backgroundColor:
          shapeState.creating.value === ShapeType.ellipse ? getActiveBackgroundColor() : '',
      }" @click="createEllipse()">
      <a-tooltip placement="top" :title="$t('toolbar.circle')">
        <i class="aiknowledge-icon icon-shape-icon" aria-hidden="true" />
        <svg width="15" height="15" viewBox="0 0 32 32" class="shape-icon" xmlns="http://www.w3.org/2000/svg">
          <g clip-path="url(#clip0_14738_33544)">
            <path fill-rule="evenodd" clip-rule="evenodd" d="M16 2.6C8.59897 2.6 2.6 8.59897 2.6 16C2.6 23.401 8.59897 29.4 16 29.4C23.401 29.4 29.4 23.401 29.4 16C29.4 8.59897 23.401 2.6 16 2.6ZM0 16C0 7.16303 7.16303 0 16 0C24.837 0 32 7.16303 32 16C32 24.837 24.837 32 16 32C7.16303 32 0 24.837 0 16Z"/>
          </g>
        </svg>
      </a-tooltip>
    </span>
    <span v-if="isOwner && showOriginalPDF && isSelfNote" :style="{
        backgroundColor:
          shapeState.creating.value === ShapeType.arrow ? getActiveBackgroundColor() : '',
      }" @click="createShape(ShapeType.arrow)">
      <a-tooltip placement="top" :title="$t('toolbar.arrow')">
        <i class="aiknowledge-icon icon-shape-arrow" aria-hidden="true" />
      </a-tooltip>
    </span>
    <span v-if="isOwner && showOriginalPDF && isSelfNote" :style="{
        backgroundColor: shapeState.creating.value === 'text' ? getActiveBackgroundColor() : '',
      }" @click="createShape('text')">
      <a-tooltip placement="top" :title="$t('toolbar.text')">
        <i class="aiknowledge-icon icon-shape-text" aria-hidden="true" />
      </a-tooltip>
    </span>
    <FullTextTranslate v-if="
        isOwner && envStore.viewerConfig.headerFulltextTranslateButton !== false
      " :pdfViewFinished="pdfViewFinished" class="ml-2" />
    <AIHighlighter v-if="
        isOwner && envStore.viewerConfig.headerAIHighlighterButton !== false
      " class="mx-2" />
    <Share v-if="isSelfNote && envStore.viewerConfig.headerShareButton !== false" />
    <Export v-if="isSelfNote && !enableFullTextTranslate" :pdf-id="selfNoteInfo?.pdfId ?? ''"
      :note-id="selfNoteInfo?.noteId ?? ''" :with-note="envStore.viewerConfig.exportPDFWithNoteButton !== false" />
    <Finder v-if="ownNoteOrVisitSharedNote && showOriginalPDF" :pdfViewInstance="pdfViewInstance" />
    <LogoDropdown />
  </CollapsedBar>

  <!-- Add the second toolbar for note tools with its own visibility control -->
  <!-- <CollapsedBar v-if="isOwner && showOriginalPDF" class="toolbar-note text-base leading-4"
    :collapsed="!toolbarSettings.toolBarNoteVisible" @toggled="
      (v) =>
        setToolbarSettings({
          toolBarNoteVisible: !v,
        })
    ">
    <template #icon-expand>
      <DoubleLeftOutlined />
    </template>
  </CollapsedBar> -->
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { ViewerController } from '@idea/pdf-annotate-viewer'
import {
  store,
  isOwner,
  selfNoteInfo,
  ownNoteOrVisitSharedNote,
  currentGroupId,
  currentNoteInfo,
} from '~/src/store'
import { useEnvStore } from '~/src/stores/envStore'
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type'
import { usePDFWEbviewPreviewMode } from '~/src/hooks/usePDFWebviewScrollMode'
import { useToolBarSettings } from '~/src/hooks/UserSettings/useToolBarSettings'

import CollapsedBar from '@/components/Common/CollapsedBar.vue'
import Finder from './components/Finder.vue'
import LogoDropdown from './components/LogoDropdown.vue'
import FullTextTranslate from './components/FullTextTranslate.vue'
import AIHighlighter from './components/AIHighlighter.vue'
import Share from './components/Share.vue'
import Export from './components/Export.vue'

// Import components from ToolBarNote
import ScreenShot from '../ToolBarNote/components/ScreenShot.vue'
import { DoubleLeftOutlined } from '@ant-design/icons-vue'
import { cancelShape, shapeState } from '../Main/mouseCore'
import { ShapeType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common'
import { useI18n } from 'vue-i18n'
import {
  executeAndSetBeforeChangeTab,
  removeBeforeChangeTab,
} from '../../hooks/UserSettings/useSideTabSettings'
import { ElementClick, reportClick } from '~/src/api/report'
import {
  ANNOTATION_SCREENSHOT,
  ScreenShotPayload,
  useClip,
} from '~/src/hooks/useHeaderScreenShot'
import { AnnotationRect } from '~/src/stores/annotationStore/BaseAnnotationController'
import { usePdfStore } from '~/src/stores/pdfStore'

const props = defineProps<{
  pdfViewInstance: ViewerController
  pdfViewFinished?: boolean
  clipSelecting?: boolean
  clipAction?: ReturnType<typeof useClip>['clipAction']
}>()

const envStore = useEnvStore()
const { showOriginalPDF, enableFullTextTranslate } = usePDFWEbviewPreviewMode()
const { toolbarSettings, setToolbarSettings } = useToolBarSettings()

const isSelfNote = computed(() => {
  return (
    isOwner.value && store.state.base.currentGroupId === SELF_NOTEINFO_GROUPID
  )
})

// 获取主题相关的颜色
const getActiveBackgroundColor = () => {
  // 从 CSS 变量中获取主题颜色，如果不可用则使用默认值
  const style = getComputedStyle(document.documentElement)
  const hoverColor = style.getPropertyValue('--site-theme-background-hover').trim()
  return hoverColor || '#414548'
}

// Add code from ToolBarNote
const pdfStore = usePdfStore()
const pdfAnnotater = computed(() => {
  return pdfStore.getAnnotater(currentNoteInfo.value?.noteId)
})

const { t } = useI18n()
shapeState.deleteTitle = t('message.confirmToDeleteNoteTip')
shapeState.deleteOk = t('viewer.delete')

const createRect = () => {
  createShape(ShapeType.rectangle)
}

const createEllipse = () => {
  createShape(ShapeType.ellipse)
}

const createShape = (type: ShapeType | 'text') => {
  if (props.clipSelecting) {
    cancelScreenShot()
  }

  if (shapeState.creating.value === type) {
    cancelShape()
  } else if (shapeState.creating.value === ShapeType.UNRECOGNIZED) {
    shapeState.creating.value = type
  } else {
    cancelShape()
    shapeState.creating.value = type
  }
}

const cancelClip = () => {
  props.clipAction?.cancelCut()
  removeBeforeChangeTab(cancelClip)
}

const handleScreenShot = () => {
  if (shapeState.creating.value !== ShapeType.UNRECOGNIZED) {
    cancelShape()
  }

  if (props.clipSelecting) {
    cancelScreenShot()
  } else {
    startScreehShot()
  }
}

const startScreehShot = async () => {
  executeAndSetBeforeChangeTab(cancelClip)

  await nextTick()

  reportClick(ElementClick.screenshots, 'on')

  const payload = await props.clipAction?.init(true)

  if (!payload) {
    return
  }

  const { newRectElement, newRectAnnotation } = payload

  pdfAnnotater.value?.UI.emit(ANNOTATION_SCREENSHOT, {
    rect: newRectElement?.getBoundingClientRect(),
    pageNum: (newRectAnnotation as AnnotationRect).pageNumber,
  } as ScreenShotPayload)
}

const cancelScreenShot = () => {
  cancelClip()
  reportClick(ElementClick.screenshots, 'off')
}

const bindEsc = (event: KeyboardEvent) => {
  if (event.key !== 'Escape') {
    return
  }

  if (props.clipSelecting) {
    cancelScreenShot()
  } else if (shapeState.creating.value !== ShapeType.UNRECOGNIZED) {
    cancelShape()
  }
}

onMounted(() => {
  document.addEventListener('keydown', bindEsc, { passive: true })
})
onUnmounted(() => {
  document.removeEventListener('keydown', bindEsc)
})
</script>

<style lang="postcss">
.toolbar-hd.collapsed-icon,
.toolbar-hd.collapsed-bar {
  left: 0;
  background: var(--site-theme-bg-secondary, #2f3337);
}

.toolbar-hd.collapsed-bar {
  width: max-content;

  .btn-collapse {
    color: inherit;
  }

  &>* {
    flex: 1 0 auto;
    display: flex;
    align-items: center;
    justify-content: center;
    height: theme('spacing.8');
    padding: theme('spacing.2') 10px;
    line-height: 1;
    border-radius: 2px;

    &:hover {
      background-color: var(--site-theme-background-hover, theme('colors.rp-dark-4'));
    }
  }
}

/* Styles from ToolBarNote */
.toolbar-note.collapsed-icon,
.toolbar-note.collapsed-bar {
  bottom: 30px;
  background: var(--site-theme-bg-secondary, #2f3337);
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

  &>* {
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
      background-color: var(--site-theme-background-hover, theme('colors.rp-dark-4'));
    }
  }
}

/* 形状图标主题适配 */
.shape-icon {
  fill: var(--site-theme-text-color);
}

</style>
