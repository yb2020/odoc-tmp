<template>
  <div class="flex flex-col h-full">
    <a-tabs
      :activeKey="currentTab"
      :centered="false"
      class="h-fullss"
      :destroyInactiveTabPane="true"
      @change="changeTab"
    >
      <template #rightExtra>
        <div class="flex items-center -mr-4">
          <ScreenShot
            :clip-selecting="clipSelecting"
            :clip-action="clipAction"
            :pdfAnnotater="pdfAnnotater"
          />
          <a-tooltip :title="$t('translate.unlockTip')">
            <!-- <img
              src="@/assets/images/unlock.svg"
              class="mr-1 w-4 h-4 cursor-pointer"
              @click="handleLock"
            > -->
            <svg class="mr-1 w-4 h-4 cursor-pointer shape-icon" @click="handleLock" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none">
              <g clip-path="url(#clip0_16244_33193)">
              <path fill-rule="evenodd" clip-rule="evenodd" d="M14.5711 0.854909C14.8872 0.854909 15.1426 1.11027 15.1426 1.42634V14.5692C15.1426 14.8853 14.8872 15.1406 14.5711 15.1406H8.14258C8.06401 15.1406 7.99972 15.0763 7.99972 14.9978V13.9978C7.99972 13.9192 8.06401 13.8549 8.14258 13.8549H10V2.14062H2.14258V7.85491C2.14258 7.93348 2.07829 7.99777 1.99972 7.99777H0.999719C0.921147 7.99777 0.856862 7.93348 0.856862 7.85491V1.42634C0.856862 1.11027 1.11222 0.854909 1.42829 0.854909H14.5711ZM11.4 2.14062V13.8549H13.8569V2.14062H11.4ZM3.37703 13.5884L4.30917 14.5205C4.32805 14.5395 4.34123 14.5634 4.34721 14.5895C4.35318 14.6156 4.35171 14.6428 4.34297 14.6681C4.33422 14.6934 4.31855 14.7158 4.29774 14.7326C4.27693 14.7495 4.25181 14.7601 4.22524 14.7634L1.02167 15.1384C0.930598 15.1491 0.852026 15.0723 0.862741 14.9795L1.23774 11.7759C1.25203 11.658 1.39667 11.608 1.4806 11.692L2.41631 12.6277L6.99131 8.05268C7.04667 7.99732 7.13774 7.99732 7.1931 8.05268L7.95024 8.80982C8.0056 8.86518 8.0056 8.95625 7.95024 9.01161L3.37703 13.5884Z" fill-opacity="0.65"/>
              </g>
              <defs>
              <clipPath id="clip0_16244_33193">
              <rect width="16" height="16" transform="matrix(-1 0 0 -1 16 16)"/>
              </clipPath>
              </defs>
            </svg>
          </a-tooltip>
          <TranslateSettingButton
            v-model:visible="setting"
            :current-tab="currentTab"
            @add-tab="addSettingTab"
            @delete-tab="deleteSettingTab"
          />
        </div>
      </template>
      <a-tab-pane
        v-for="tab in tabs"
        :key="tab.type"
        :forceRender="true"
      >
        <template #tab>
          <AITranslateTab v-if="tab.type === TranslateTabKey.ai">
            {{ tab.title }}
          </AITranslateTab>
          <template v-else>
            {{ tab.title }}
          </template>
        </template>
      </a-tab-pane>
    </a-tabs>
    <div class="flex-1 overflow-hidden">
      <TranslatePanel
        ref="translatePanelRef"
        :pdf-id="translateStore.pdfId"
        :tippy-handler="tippyHandler"
        :add-to-note-handler="addToNoteHandler"
        :fix-placement="() => {}"
        :current-tab="currentTab"
        :update-current-tab="changeTab"
        class="h-full !pt-4"
      />
    </div>
  </div>
</template>
<script setup lang="ts">
import TranslatePanel from '@/components/Tippy/Translate/Panel.vue';
import { UniTranslateResp } from '~/src/api/translate';
import { createTranslateTippyVue } from '~/src/dom/tippy';
import useCreateNote from '~/src/hooks/note/useCreateNote';
import { useTranslateLock, useTranslateTabs } from '~/src/hooks/useTranslation';
import { currentNoteInfo, selfNoteInfo } from '~/src/store';
import { getPlatformKey } from '~/src/store/shortcuts';
import { usePdfStore } from '~/src/stores/pdfStore';
import {
  useTranslateStore,
  TranslateTabKey,
} from '~/src/stores/translateStore';
import TranslateSettingButton from '@/components/Tippy/Translate/Setting.vue';
import ScreenShot from './ScreenShot.vue';
import { computed, nextTick, ref } from 'vue';
import { ElementClick, PageType, reportElementClick } from '~/src/api/report';
import { useClip } from '~/src/hooks/useHeaderScreenShot';
import AITranslateTab from '@/components/Tippy/Translate/AITranslateTab.vue';

defineProps<{
  clipSelecting: boolean;
  clipAction: ReturnType<typeof useClip>['clipAction'];
}>();

const translateStore = useTranslateStore();

const pdfStore = usePdfStore();
const pdfViewer = pdfStore.getViewer(currentNoteInfo.value.pdfId);
const pdfAnnotater = computed(() =>
  pdfStore.getAnnotater(currentNoteInfo.value.noteId)
);

const setting = ref(false);
const {
  currentTab,
  tabs,
  initTabs,
  changeTab,
  addSettingTab,
  deleteSettingTab,
} = useTranslateTabs();

initTabs();

const { add: addNote, addWord } = useCreateNote({
  pdfId: translateStore.pdfId,
  noteId: selfNoteInfo.value?.noteId || '',
  pdfViewer,
});

const addToNoteHandler = (
  isPhrase: boolean,
  phrase: string,
  translation: string,
  translationRes: UniTranslateResp
) => {
  if (isPhrase) {
    addWord(
      { info: translateStore.extraInfo.selections ?? null },
      phrase,
      translationRes
    );
  } else if (translateStore?.content.ocrTranslate) {
    translateStore.extraInfo.ocr?.addOcrNote(translation);
  } else {
    addNote({ info: translateStore.extraInfo.selections ?? null }, translation);
  }
};

const tippyHandler = () => {};

const translatePanelRef = ref<InstanceType<typeof TranslatePanel>>();
const { translateLock } = useTranslateLock();
const handleLock = async () => {
  const translateData = translatePanelRef.value?.getTranslateData() || null;
  await nextTick();
  translateLock.value = false;
  createTranslateTippyVue.show({
    pdfId: translateStore.pdfId,
    triggerEle:
      document.querySelector('.right-panel.js-drawer') || document.body,
    isExistingAnnotation: false,
    pageTexts: translateStore.extraInfo.selections || undefined,
    ocr: translateStore.extraInfo.ocr,
    translatedData: translateData,
  });
  reportElementClick({
    element_name: ElementClick.trans_box_release,
    page_type: PageType.note,
    type_parameter: 'none',
  });
};
</script>
<style lang="less" scoped>
/* 形状图标主题适配 */
.shape-icon {
  fill: var(--site-theme-text-color);
}
</style>
