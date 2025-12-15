<script setup lang="ts">
export interface PDFNoteInfo {
  noteId: string;
  groupId: string;
}

export interface PDFBaseInfo {
  pdfId: string;
  pdfUrl: string;
  showCopyRight: boolean;
}

export interface PDFWebviewWithNoteProps {
  pdfBaseInfo: PDFBaseInfo;
  currentNoteId: string;
  pdfNoteInfos: PDFNoteInfo[];
  clipSelecting: boolean;
  clipAction: ReturnType<typeof useClip>['clipAction'];
}
import { computed, ref, shallowRef, watch } from 'vue';
import { message } from 'ant-design-vue';
import { useI18n } from 'vue-i18n';
import { ViewerController, ViewerEvent } from '@idea/pdf-annotate-viewer';

import Left from '../Left/index.vue';
import PDFReaderWebview from './PDFReaderWebview.vue';
import PDFNoteWebview from './PDFNoteWebview.vue';
import { ToolbarLeftSideBtnEvent } from '../ToolBar/components/LeftSideBtn.vue';
import ToolbarHead from '@/components/ToolBarHead/index.vue';

import { ElementClick, reportClick } from '~/src/api/report';
import { useFullTextTranslateStore } from '~/src/stores/fullTextTranslateStore';
import { changePDFScale, goToPDFPage } from '~/src/dom/pdf';
import { connectComment, getQueryString, scrollMark } from '~/src/util/scroll';
import { ANNOTATION_ID_PARAM, PAGE_NUMBER_PARAM } from '~/src/constants';

import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import {
  currentNoteInfo,
  isOwner,
  ownNoteOrVisitSharedNote,
  selfNoteInfo,
  store,
} from '~/src/store';
import { CopyrightProps } from '../Copyright/type';
import PDFAnnotate from './PDFAnnotate.vue';
import {
  useLeftSideTabSettings,
  useRightSideTabSettings,
} from '~/src/hooks/UserSettings/useSideTabSettings';
import { RightSideBarType } from '~/src/components/Right/TabPanel/type';
import {
  usePDFWebviewScrollMode,
  usePDFWEbviewPreviewMode,
} from '@/hooks/usePDFWebviewScrollMode';
import usePDFWebviewSettings from '~/src/hooks/UserSettings/usePDFWebviewSettings';
import { SaveMode } from '@/hooks/UserSettings/useRemoteUserSettings';
import { updateDocReadStatus } from '~/src/api/setting';
import { debounce } from 'lodash';
import { NoteSubTypes } from '~/src/store/note/types';
import { useClip } from '~/src/hooks/useHeaderScreenShot';

const debouncedUpdateDocReadStatus = debounce(updateDocReadStatus, 2000, {
  trailing: true,
});
const { setSideTabSetting } = useRightSideTabSettings();

const props = defineProps<PDFWebviewWithNoteProps>();

const pdfViewInstance = shallowRef<ViewerController>();
const pdfViewFinished = ref(false);

const emit = defineEmits<{
  (
    event: 'pdfWebviewFinished',
    payload: { pdfWebview: ViewerController; pdfId: string }
  ): void;
}>();

const { setPDFWebviewInstance } = usePDFWebviewScrollMode();

const numPages = ref(0);

const onPDFNoteWebviewFinished = ({
  pdfWebview,
}: {
  pdfWebview: ViewerController;
}) => {
  numPages.value =
    pdfWebview.getDocumentViewer().getPdfViewer().pdfDocument?.numPages || 0;

  pdfViewFinished.value = true;
  pdfViewInstance.value = pdfWebview;

  setPDFWebviewInstance('main', pdfWebview);

  emit('pdfWebviewFinished', { pdfWebview, pdfId: props.pdfBaseInfo.pdfId });
};

// ---------START: 用户设置userSettings---------
const isSelfNoteInfo = computed(() => {
  const currentGroupId = props.pdfNoteInfos.find(
    (item) => item.noteId === props.currentNoteId
  )?.groupId;
  return currentGroupId === SELF_NOTEINFO_GROUPID;
});

const { pdfWebviewSettings, setScaleSetting, setCurrentPageSetting } =
  usePDFWebviewSettings(isSelfNoteInfo.value);

const { t } = useI18n();

const onPDFNoteWebviewCreated = ({
  pdfWebview,
}: {
  pdfWebview: ViewerController;
}) => {
  const checkSelf = () => isSelfNoteInfo.value && isOwner.value;

  const getSaveMode = (isSelf: boolean) =>
    isSelf ? SaveMode.remote : SaveMode.local;

  pdfWebview.addEventListener(ViewerEvent.SCALE_CHANGING, (scaleInfo) => {
    scaleInfo.presetValue = scaleInfo.presetValue || '';
    setScaleSetting(scaleInfo, getSaveMode(checkSelf()));
  });

  pdfWebview.addEventListener(ViewerEvent.PAGE_CHANGING, (pageInfo) => {
    const isSelf = checkSelf();

    setCurrentPageSetting(pageInfo.pageNumber, getSaveMode(isSelf));

    const pageCount =
      pdfWebview.getDocumentViewer().getPdfViewer().pagesCount - 1;
    const { pageNumber } = pageInfo;

    if (isSelf) {
      debouncedUpdateDocReadStatus({
        paperId: selfNoteInfo.value?.paperId ?? '',
        pdfId: selfNoteInfo.value?.pdfId ?? '',
        progress:
          pageNumber >= pageCount
            ? 100
            : Math.floor((pageNumber * 100) / pageCount),
      });
    }
  });

  if (ownNoteOrVisitSharedNote.value && isSelfNoteInfo.value) {
    pdfWebview.addEventListener(ViewerEvent.PAGES_INIT, async () => {
      // 初始化时调用一次updateReadStatus，设置进度为1%
      const isSelf = isSelfNoteInfo.value && isOwner.value;
      if (isSelf) {
        updateDocReadStatus({
          paperId: selfNoteInfo.value?.paperId ?? '',
          pdfId: selfNoteInfo.value?.pdfId ?? '',
          progress: 1,
        });
      }

      // url 带了笔记跳转参数，则不执行跳转到上次阅读位置
      const annotationId = getQueryString(ANNOTATION_ID_PARAM);
      const pageNumber = getQueryString(PAGE_NUMBER_PARAM);

      if (annotationId && pageNumber) {
        setSideTabSetting({
          tab: RightSideBarType.Note,
          subTab: NoteSubTypes.Annotation,
        });

        scrollMark(parseInt(pageNumber), annotationId, pdfWebview);
        connectComment(annotationId);

        return;
      }

      const userSettingInfo = store.state.documents.userSettingInfo;
      if (
        (userSettingInfo.scalePresetValue &&
          userSettingInfo.scalePresetValue !==
            pdfWebviewSettings.value.scalePresetValue) ||
        (userSettingInfo.scale &&
          userSettingInfo.scale !== pdfWebviewSettings.value.scale)
      ) {
        changePDFScale(
          (userSettingInfo.scalePresetValue as 'page-width') ||
            userSettingInfo.scale,
          pdfWebview
        );
      }

      if (userSettingInfo.currentPage === undefined) {
        return;
      }
      goToPDFPage(userSettingInfo.currentPage, true, pdfWebview);
      message.info(t('message.returnToLastReadingPositionTip'));
    });
  }
};
// ---------END: 用户设置userSettings---------

const copyRightProps = computed<CopyrightProps | null>(() => {
  if (!props.pdfBaseInfo.showCopyRight) {
    return null;
  }

  if (!currentNoteInfo.value) {
    return null;
  }

  const options: Omit<CopyrightProps, 'sourceMark'> = {
    crawlUrl: currentNoteInfo.value.crawlUrl || '',
    uploadUserId: currentNoteInfo.value.uploadUserId || '',
    isUserUpload: !!currentNoteInfo.value.isUserUpload,
    licenceType: currentNoteInfo.value.licenceType || '',
  };

  if (!currentNoteInfo.value.isUserUpload) {
    return {
      ...options,
      sourceMark: currentNoteInfo.value.sourceMark || '',
    };
  }

  if (isOwner.value) {
    return null;
  }

  return {
    ...options,
    sourceMark: selfNoteInfo.value.userInfo.nickName,
  };
});

const {
  sideTabSettings: leftSideTabSetting,
  setSideTabSetting: setLeftSideTabSetting,
} = useLeftSideTabSettings(
  isSelfNoteInfo.value ? SaveMode.remote : SaveMode.local
);

// ---------START: 全文翻译Webview---------
const onPDFReaderWebviewFinished = ({
  pdfWebview,
}: {
  pdfWebview: ViewerController;
}) => {
  setPDFWebviewInstance('tied', pdfWebview);
};

const fullTextTranslateStore = useFullTextTranslateStore();

const { showOriginalPDF } = usePDFWEbviewPreviewMode();

watch(
  () => showOriginalPDF.value,
  (newVal) => {
    if (!newVal) {
      setLeftSideTabSetting({ shown: false });
    }
  }
);

const showFullTextTranslate = computed(() => {
  return (
    currentNoteInfo.value?.pdfId === props.pdfBaseInfo.pdfId &&
    fullTextTranslateStore.pdfId
  );
});

// ---------END: 全文翻译Webview---------

const handleToolbarEvent = (payload: ToolbarLeftSideBtnEvent) => {
  if (payload.type === 'toolbar:leftside') {
    const shown = !leftSideTabSetting.value.shown;
    setLeftSideTabSetting({
      shown,
    });
    reportClick(ElementClick.directory, shown ? 'on' : 'off');
  }
};

const showPDFAnnotations = computed(() => {
  return currentNoteInfo.value?.showAnnotation === false ? false : true;
});
</script>

<template>
  <div class="viewer-container">
    <Left
      v-if="pdfViewFinished"
      v-show="showOriginalPDF"
      :pdf-id="pdfBaseInfo.pdfId"
      :side-tab-settings="leftSideTabSetting"
      :set-side-tab-setting="setLeftSideTabSetting"
      :pdfViewInstance="pdfViewInstance!"
    />
    <div class="main">
      <PDFReaderWebview
        v-if="showFullTextTranslate"
        :pdf-url="fullTextTranslateStore.fullTextTranslatePDFUrl"
        :scale="pdfWebviewSettings.scalePresetValue || pdfWebviewSettings.scale"
        :pdf-id="pdfBaseInfo.pdfId"
        :show-toolbar="!showOriginalPDF"
        :main-pdf-webview-settings="pdfWebviewSettings"
        @finished="onPDFReaderWebviewFinished"
      />
      <PDFNoteWebview
        v-show="showOriginalPDF"
        :pdf-url="pdfBaseInfo.pdfUrl"
        :paper-id="selfNoteInfo?.paperId || ''"
        :is-private-paper="selfNoteInfo?.isPrivatePaper ?? false"
        :pdf-webview-settings="pdfWebviewSettings"
        :copy-right-props="copyRightProps"
        :pdf-id="pdfBaseInfo.pdfId"
        :is-group-pdf="
          store.state.base.currentGroupId !== SELF_NOTEINFO_GROUPID
        "
        :left-side-tab-settings="leftSideTabSetting"
        :show-annotations="showPDFAnnotations"
        :clip-selecting="clipSelecting"
        :clip-action="clipAction"
        @toolbar-event="handleToolbarEvent"
        @finished="onPDFNoteWebviewFinished"
        @created="onPDFNoteWebviewCreated"
      />

      <template v-if="pdfViewFinished">
        <PDFAnnotate
          v-for="item in pdfNoteInfos"
          v-show="showOriginalPDF"
          :key="item.noteId"
          :pdfBaseInfo="pdfBaseInfo"
          :noteInfo="item"
          :pdfViewInstance="pdfViewInstance"
          :clip-action="clipAction"
        />
      </template>
      <ToolbarHead
        v-if="pdfViewInstance"
        :pdfViewInstance="pdfViewInstance"
        :pdfViewFinished="pdfViewFinished"
        :clip-selecting="clipSelecting"
        :clip-action="clipAction"
      />
    </div>
  </div>
</template>

<style scoped lang="less">
.viewer-container {
  height: 100%;
  flex: 1;
  display: flex;

  .main {
    flex: 1;
    height: 100%;
    position: relative;
    z-index: 1;
    display: flex;
  }
}
</style>
