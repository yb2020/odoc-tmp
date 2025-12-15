<template>
  <Layout
    ref="layoutRef"
    :style="style"
    :title="''"
    :isDing="isDing"
    :tippyHandler="onTippyHandler"
    :has-lock="true"
    group="translate"
    :class="[JS_IGNORE_MOUSE_OUTSIDE]"
  >
    <template #title>
      <div>
        {{ $t('translate.translate') }}
        <a-tooltip placement="topLeft">
          <template #title>
            {{
              $t('translate.translateTip', { cmd: CtrlOrCommand })
            }}
          </template>
          <question-circle-outlined />
        </a-tooltip>
      </div>
    </template>
    <div>
      <div class="flex justify-between w-full items-center">
        <div class="tabs flex-1">
          <div
            v-for="tab in tabs"
            :key="tab.type"
            :class="['tab', currentTab === tab.type ? 'active' : '']"
            @click="changeTab(tab.type)"
          >
            <AITranslateTab v-if="tab.type === TranslateTabKey.ai">
              {{ tab.title }}
            </AITranslateTab>
            <template v-else>
              {{ tab.title }}
            </template>
          </div>
        </div>
        <TranslateSettingButton
          v-model:visible="setting"
          :current-tab="currentTab"
          @add-tab="addSettingTab"
          @delete-tab="deleteSettingTab"
        />
      </div>
      <Panel
        ref="translatePanelRef"
        :pdf-id="pdfId"
        :add-to-note-handler="addToNoteHandler"
        :tippy-handler="tippyHandler"
        :fix-placement="fixPlacement"
        :current-tab="currentTab"
        :update-current-tab="changeTab"
      />
    </div>
  </Layout>
</template>
<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import { UniTranslateResp } from '@/api/translate';
import Layout from '@/components/Tippy/Layout/index.vue';

import {
  TranslateTabKey,
  useTranslateStore,
} from '~/src/stores/translateStore';
import { getPlatformKey } from '~/src/store/shortcuts';
import { useI18n } from 'vue-i18n';
import { JS_IGNORE_MOUSE_OUTSIDE } from '@idea/pdf-annotate-viewer';
import Panel from './Panel.vue';
import TranslateSettingButton from './Setting.vue';
import { useTranslateTabs } from '~/src/hooks/useTranslation';
import AITranslateTab from '@/components/Tippy/Translate/AITranslateTab.vue';

const props = defineProps<{
  pdfId: string;
  width: number;
  isDing: import('vue').Ref<boolean>; // 新增：响应式的钉住状态
  tippyHandler: (event: 'ding' | 'close' | 'unding' | 'lock') => void;
  addToNoteHandler: (
    isPhrase: boolean,
    phrase: string,
    translation: string,
    translationRes: UniTranslateResp
  ) => void;
  fixPlacement(): void;
}>();

const style = computed(() => {
  return {
    width: props.width + 'px',
  };
});

const store = useTranslateStore();

const setting = ref(false);

watch(
  () => setting.value,
  (newVal) => {
    if (newVal) {
      props.tippyHandler('ding');
    } else if (!layoutRef.value?.getDing()) {
      props.tippyHandler('unding');
    }
  }
);

const {
  currentTab,
  tabs,
  initTabs,
  addSettingTab,
  deleteSettingTab,
  changeTab,
} = useTranslateTabs();
initTabs();

/**
 * 如果直接展示perfectscrollbar，会导致在tippy里面bar在左边而不是右边
 */
const shown = ref(false);

onMounted(() => {
  shown.value = true;
});

const layoutRef = ref<InstanceType<typeof Layout>>();

const CtrlOrCommand = getPlatformKey() === 'win32' ? 'Ctrl' : 'Command';

const translatePanelRef = ref<InstanceType<typeof Panel>>();
const onTippyHandler = async (evt: 'ding' | 'close' | 'unding' | 'lock') => {
  if (evt === 'lock') {
    translatePanelRef.value?.storeTranslateData();
  }
  props.tippyHandler(evt);
};
</script>
<style lang="less" scoped>
.tabs {
  padding: 8px 16px 0px;
  display: flex;
  border-bottom: 1px solid #e9ebf0;
  margin-bottom: 8px;
  flex-direction: row;
  flex-wrap: wrap;
  width: 100%;

  .tabSettingIcon {
    flex: 1;
    display: flex;
    justify-content: flex-end;

    .tabSetting {
      padding: 5px 0px;
      flex: end;
      color: black;
      cursor: pointer;
      font-size: 12px;
      margin-right: 5px;
      display: flex;
      justify-content: flex-end;
    }
  }

  .tab {
    // flex: 1;

    font-size: 14px;
    font-weight: 400;
    color: #4e5969;
    line-height: 20px;
    padding: 6px 0;
    margin-right: 24px;
    cursor: pointer;

    &.active {
      font-weight: 500;
      color: #176ae5;
      border-bottom: 2px solid #176ae5;
    }
  }
}
.translate-ctrl {
  position: absolute;
  top: 9px;
  left: 54px;
  z-index: 9;
  color: #fff;
  opacity: 0.9;
}
</style>
