<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import {
  store,
  currentGroupId,
  currentNoteInfo,
  selfNoteInfo,
  pdfStatusInfo,
  isOwner,
} from '~/src/store';
import Right from '../Right/index.vue';
import NoteInfo from '../Private/NoteInfo.vue';
import FinderVue from '../Right/Finder.vue';
import PDFWebviewWithNote, {
  PDFWebviewWithNoteProps,
} from './PDFWebviewWithNote.vue';

import { ViewerController } from '@idea/pdf-annotate-viewer';
import NoPaper from './NoPaper.vue';
import { useBaseStore } from '~/src/stores/baseStore';
import { useSyncAnnotationVisible } from '~/src/stores/annotationStore';
import { useEnvStore } from '~/src/stores/envStore';
import { usePdfStore } from '~/src/stores/pdfStore';
import { useClip } from '~/src/hooks/useHeaderScreenShot';

const finderRef = ref();

useSyncAnnotationVisible();

type PDFWebviewWithNoteParams = Omit<
  PDFWebviewWithNoteProps,
  'clipSelecting' | 'clipAction'
>;

const pdfWebviewWithNoteList = computed<PDFWebviewWithNoteParams[]>(() => {
  /**
   * 根据pdfId聚合一下数据
   */
  const pdfIdMap: Record<string, PDFWebviewWithNoteParams> = {};
  const { noteInfoMap } = store.state.base;
  Object.keys(noteInfoMap).forEach((groupId) => {
    const { noteId, pdfId, pdfUrl } = noteInfoMap[groupId];
    if (!(pdfId in pdfIdMap)) {
      pdfIdMap[pdfId] = {
        pdfBaseInfo: {
          pdfId,
          pdfUrl,
          showCopyRight: envStore.viewerConfig.PDFSource !== false,
        },
        currentNoteId: noteId,
        pdfNoteInfos: [],
      };
    }

    pdfIdMap[pdfId].pdfNoteInfos.push({
      noteId,
      groupId,
    });
  });

  return Object.values(pdfIdMap);
});

const pdfStore = usePdfStore();

const onPdfWebviewFinished = ({
  pdfId,
  pdfWebview,
}: {
  pdfId: string;
  pdfWebview: ViewerController;
}) => {
  if (currentNoteInfo.value && currentNoteInfo.value.pdfId === pdfId) {
    pdfStore.setViewer(pdfId, pdfWebview);
  }

  if (finderRef.value) {
    pdfWebview.getFinderViewer().setWrapper(finderRef.value.getFinderWrapper());
  }
};

watch(
  () => currentNoteInfo.value?.pdfId,
  (newVal) => {
    if (!newVal) {
      return;
    }

    const pdfWebviewController = pdfStore.getViewer(newVal);

    if (pdfWebviewController) {
      pdfStore.setViewer(newVal, pdfWebviewController);
    }
  }
);

const baseStore = useBaseStore();
baseStore.getRedDotInfo();

const envStore = useEnvStore();

const { clipSelecting, clipAction } = useClip();

clipAction.onMouseDownClearClip();
</script>

<template>
  <div class="wrap">
    <div class="viewer">
      <PDFWebviewWithNote
        v-for="item in currentGroupId ? pdfWebviewWithNoteList : []"
        v-show="item.pdfBaseInfo.pdfId === currentNoteInfo?.pdfId"
        :key="item.pdfBaseInfo.pdfId"
        :pdf-base-info="item.pdfBaseInfo"
        :pdf-note-infos="item.pdfNoteInfos"
        :current-note-id="currentNoteInfo?.noteId || ''"
        :current-pdf-id="currentNoteInfo?.pdfId || ''"
        :clip-selecting="clipSelecting"
        :clip-action="clipAction"
        @pdf-webview-finished="onPdfWebviewFinished"
      />
      <NoPaper v-if="currentGroupId && !currentNoteInfo" />
    </div>

    <FinderVue ref="finderRef" />
    <Right
      v-if="isOwner || pdfStatusInfo.hasPdfAccessFlag"
      :is-owner="isOwner"
      :paper-id="selfNoteInfo!.paperId"
      :pdf-id="selfNoteInfo!.pdfId"
      :pdf-status-info="pdfStatusInfo"
      :is-private-paper="selfNoteInfo!.isPrivatePaper"
      :clip-selecting="clipSelecting"
      :clip-action="clipAction"
    />
    <NoteInfo v-if="!isOwner" />
  </div>
</template>

<style scoped lang="less">
.wrap {
  height: 100%;
  width: 100%;
  display: flex;

  .viewer {
    flex: 1;
    height: 100%;
    position: relative;
    background-color: var(--site-theme-pdf-collapsed);
  }
}
</style>
