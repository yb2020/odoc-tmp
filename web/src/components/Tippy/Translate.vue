<template>
  <Layout
    ref="layoutRef"
    :style="style"
    :title="''"
    :tippyHandler="tippyHandler"
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
    <PerfectScrollbar
      v-if="shown"
      class="translate-content"
    >
      <div class="tabs">
        <div
          v-for="tab in tabs"
          :key="tab.type"
          :class="['tab', currentTab === tab.type ? 'active' : '']"
          @click="changeTab(tab.type)"
        >
          {{ tab.title }}
        </div>
        <div class="tabSettingIcon">
          <a-button
            v-if="showSetting"
            type="link"
            class="tabSetting"
            @click="openTranslateConfigSetting"
          >
            <template #icon>
              <SettingOutlined />
            </template>
            <!-- 管理翻译渠道 -->
          </a-button>
        </div>

        <settingView
          :visible="setting"
          @update:visible="visibleChange"
        />
      </div>
      <div class="flex justify-end items-center gap-3 mb-3">
        <FontSizeSelect />
        <GlossaryButton
          v-if="showGlossaryButton"
          :getPopupContainer="getPopupContainer"
          @change="onGlossaryCheckedChange"
        />
      </div>
      <div ref="textareaRef">
        <div
          v-if="!editing"
          class="text js-interact-drag-ignore"
          @click="showEdit"
        >
          <div
            class="ellipsis"
            :style="{ fontSize: fontSize + 'px', lineHeight: '1.4' }"
          >
            {{ input }}
          </div>
        </div>
        <a-textarea
          v-else
          v-model:value="input"
          class="js-interact-drag-ignore"
          :placeholder="$t('translate.placeholder')"
          :auto-size="{ minRows: 1, maxRows: 6 }"
          :style="{
            color: 'rgba(0,0,0,.85)',
            border: '1px solid #f5f7fa',
            background: '#fefefe',
            boxShadow: '0 4px 16px 0 rgba(12, 53, 115, 0.2)',
            fontSize: '14px',
          }"
          @blur="onTranslate"
          @keydown="handleKeyDown"
        />
        <div class="divide">
          <span class="title">{{ $t('translate.translation') }}</span>
          <span class="line" />
          <CopyOutlined
            class="text-2xl !text-rp-neutral-8 ml-3 cursor-pointer"
            @click="handleCopy"
          />
          <span
            v-if="allowToAddNote && !store.multiSegment"
            class="icon"
          >
            <a-tooltip
              v-if="!store.isExistingAnnotation"
              :getPopupContainer="getPopupContainer"
              placement="left"
            >
              <template #title>{{ $t('translate.addToAnnotation') }}</template>
              <i
                class="aiknowledge-icon text-2xl icon-add-to-note"
                @click="addToNote(false)"
              />
            </a-tooltip>
            <a-tooltip
              v-if="isPhrase"
              :getPopupContainer="getPopupContainer"
              placement="left"
            >
              <template #title>{{ $t('translate.addToWordPhrase') }}</template>
              <i
                class="aiknowledge-icon text-2xl ml-2 icon-add-to-phrase"
                @click="addToNote(true)"
              />
            </a-tooltip>
          </span>
        </div>
        <a-spin :spinning="fetchState.pending">
          <div
            v-if="fetchState.error"
            class="error"
            @click="fetch"
          >
            <redo-outlined /> {{ fetchState.error.message }}
          </div>
          <template v-else>
            <div
              v-if="translatedData?.targetResp?.length"
              class="translate-word js-interact-drag-ignore"
            >
              <div>
                <WordPronunciation
                  v-if="translatedData.britishSymbol"
                  prefix="英"
                  :title="translatedData.britishSymbol"
                  :type="translatedData.britishFormat"
                  :audio="translatedData.britishPronunciation"
                />
                <WordPronunciation
                  v-if="translatedData.americaSymbol"
                  prefix="美"
                  :title="translatedData.americaSymbol"
                  :type="translatedData.americaFormat"
                  :audio="translatedData.americaPronunciation"
                />
              </div>

              <ul
                v-if="translatedData.targetResp?.length"
                class="translate-part-list"
                :style="{ fontSize: fontSize + 'px', lineHeight: '1.4' }"
              >
                <li
                  v-for="item in translatedData.targetResp"
                  :key="item.part"
                  class="translate-part-item"
                >
                  <span class="translate-part">{{ item.part }}</span>
                  <div class="translate-part-txt">
                    <p
                      v-for="text in item.targetContent"
                      :key="text"
                    >
                      {{ text }};
                    </p>
                  </div>
                </li>
              </ul>
            </div>
            <div
              v-else
              class="translate-text js-interact-drag-ignore"
              :style="{ fontSize: fontSize + 'px', lineHeight: '1.4' }"
            >
              {{ translatedContent }}
            </div>
          </template>
        </a-spin>
        <div
          v-if="feedbackVisible && !fetchState.error"
          class="feedback"
        >
          <a-tooltip
            :getPopupContainer="getPopupContainer"
            placement="right"
          >
            <template #title>
              {{ $t('translate.reportTip') }}
            </template>
            <DislikeOutlined @click="submitFeedback" />
          </a-tooltip>
        </div>
      </div>
    </PerfectScrollbar>
  </Layout>
</template>
<script lang="ts" setup>
import { computed, nextTick, onMounted, ref, watch, onUnmounted } from 'vue';
import {
  QuestionCircleOutlined,
  SettingOutlined,
  CopyOutlined,
} from '@ant-design/icons-vue';
import useFetch from '~/src/hooks/useFetch';
import { correctTranslate, UniTranslateResp } from '@/api/translate';
import Layout from './Layout/index.vue';
import WordPronunciation from '@common/components/Notes/components/WordPronunciation.vue';
import trim from 'lodash-es/trim';

import {
  useTranslateStore,
  TranslateTabKey,
  DEFAULT_TRANSLATE_TAB_KEY,
  LSKEY_CURRENT_TRANSLATE_TAB,
  GoogleTranslateStyle,
} from '~/src/stores/translateStore';
import { RedoOutlined, DislikeOutlined } from '@ant-design/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import reporter from '@idea/aiknowledge-report';
import { PageType, EventCode, reportTranslateClick } from '~/src/api/report';
import { message } from 'ant-design-vue';
import { getPlatformKey } from '~/src/store/shortcuts';
import { useVipStore } from '@common/stores/vip';
import settingView from './config/TranslateSetting.vue';
import { gteElectronVersion } from '~/src/util/env';
import { emitter, CONFIG_RESET_TYPE, CONFIG_ADD_TYPE } from './config/config';
import { useI18n } from 'vue-i18n';
import {
  useTranslateApi,
  useTranslateFontSize,
} from '~/src/hooks/useTranslation';
import { JS_IGNORE_MOUSE_OUTSIDE } from '@idea/pdf-annotate-viewer';
import GlossaryButton from '@/components/Translate/Glossary/Button.vue';
import FontSizeSelect from '@/components/Translate/FontSize.vue';
import copyTextToClipboard from 'copy-text-to-clipboard';
import { copyToPaste } from '~/src/util/copy';
import { useUserStore } from '~/src/common/src/stores/user';

const props = defineProps<{
  pdfId: string;
  width: number;
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

const visibleChange = (visible: boolean) => {
  setting.value = visible;
};

const store = useTranslateStore();
const vipStore = useVipStore();

const translatedContent = computed(
  () => translatedData?.value?.targetContent || ''
);

const input = ref(store.content.origin);

const preInput = ref(store.content.origin);

const editing = ref(false);

const setting = ref(false);

const isPhrase = computed(() => input.value.length < 60);

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

const textareaRef = ref();

const reportTranslate = () =>
  reportTranslateClick({
    page_type: PageType.note,
    type_parameter: props.pdfId || '',
    tran_content: input.value,
    sources: currentTab.value,
  });

const titles: { [k in TranslateTabKey]?: string } = {
  [TranslateTabKey.youdao]: 'translate.youdao',
  [TranslateTabKey.baidu]: 'translate.baidu',
  [TranslateTabKey.idea]: 'IDEA',
};

const { t } = useI18n();

const tabs = computed(() => {
  return store.tabs.map((tab: any) => {
    if (!titles[tab as TranslateTabKey]) {
      const cusTab = store.cusTabs.find((item: any) => item.type === tab);
      if (cusTab) {
        titles[tab as TranslateTabKey] = cusTab.name;
      }
    }
    return {
      title: [TranslateTabKey.youdao, TranslateTabKey.baidu].includes(tab)
        ? t(titles[tab as TranslateTabKey] || '')
        : titles[tab as TranslateTabKey],
      type: tab,
    };
  });
});

const currentTab = useLocalStorage<string>(
  LSKEY_CURRENT_TRANSLATE_TAB,
  DEFAULT_TRANSLATE_TAB_KEY
);

const initTabs = async () => {
  try {
    const tabs = await store.getTabs();
    if (currentTab.value !== DEFAULT_TRANSLATE_TAB_KEY) {
      if (!tabs.includes(currentTab.value as TranslateTabKey)) {
        changeTab(DEFAULT_TRANSLATE_TAB_KEY);
      }
    }
  } catch (error) {
    console.error('getTranslateTabList error');
  }
};

initTabs();

const addAddChannelListener = (e: any) => {
  const configInfo = JSON.parse(JSON.stringify(e));
  configInfo.verified = true;
  let type = TranslateTabKey.other;
  if (configInfo.txSecretId?.length > 0) {
    /// 腾讯翻译
    configInfo.type = TranslateTabKey.tencent;
    store.txConfig = configInfo;
    type = TranslateTabKey.tencent;
  }
  if (configInfo.aliAccessKeyId?.length > 0) {
    /// 阿里翻译
    configInfo.type = TranslateTabKey.ali;
    store.aliConfig = configInfo;
    type = TranslateTabKey.ali;
  }
  if (
    configInfo.googleApiKey?.length > 0 ||
    store.googleConfigVersion !== GoogleTranslateStyle.none
  ) {
    /// 谷歌翻译
    configInfo.type = TranslateTabKey.google;
    store.googleConfig = configInfo;
    type = TranslateTabKey.google;
  }
  if (configInfo.deepLKey?.length > 0) {
    /// deepl
    configInfo.type = TranslateTabKey.deepl;

    store.deeplConfig = configInfo;
    type = TranslateTabKey.deepl;
  }

  if (!titles[type]) {
    tabs.value.push({
      title: configInfo.name,
      type: configInfo.type,
    });
  }

  initTabs();
};

emitter.on(CONFIG_ADD_TYPE, addAddChannelListener);

const resetChannelListtener = (e: any) => {
  if (!e) {
    return;
  }

  const deleteInfo = JSON.parse(JSON.stringify(e));
  delete titles[deleteInfo.type as TranslateTabKey];
  currentTab.value = DEFAULT_TRANSLATE_TAB_KEY;
  for (let i = 0; i < tabs.value.length; i++) {
    if (tabs.value[i].type === deleteInfo.type) {
      tabs.value.splice(i, 1);
      initTabs();
      break;
    }
  }
};
emitter.on(CONFIG_RESET_TYPE, resetChannelListtener);

onUnmounted(() => {
  emitter.off(CONFIG_ADD_TYPE, addAddChannelListener);
  emitter.off(CONFIG_RESET_TYPE, resetChannelListtener);
});

const openTranslateConfigSetting = () => {
  setting.value = !setting.value;
};

const changeTab = (type: string) => {
  currentTab.value = type;
  fetch();
  reportTranslate();
};

const { data: translatedData, run: fetchTranslate } = useTranslateApi();
const { fetchState, fetch } = useFetch(async () => {
  const value = trim(input.value);

  await fetchTranslate(currentTab.value as TranslateTabKey, value, {
    onVipLimit() {
      props.tippyHandler('close');
    },
  });

  editing.value = false;
  feedbackVisible.value = true;

  props.fixPlacement();

  // 刷新积分
  const userStore = useUserStore();
  userStore.refreshUserCredits()
}, false);

const onTranslate = () => {
  if (input.value === preInput.value) {
    editing.value = false;
    return;
  }

  fetch();

  preInput.value = input.value;
};

watch(
  () => store.content,
  (newVal) => {
    input.value = newVal.origin;
    preInput.value = newVal.origin;

    if (newVal.ocrTranslate) {
      editing.value = false;
      translatedData.value = newVal.ocrTranslate;
      props.fixPlacement();
    } else {
      fetch();
    }

    reportTranslate();
  }
);

const showEdit = () => {
  editing.value = true;
  nextTick(() => {
    const div = textareaRef.value as HTMLDivElement;
    div.querySelector('textarea')?.focus();
  });
};
/**
 * 如果直接展示perfectscrollbar，会导致在tippy里面bar在左边而不是右边
 */
const shown = ref(false);

const showSetting = gteElectronVersion('1.20.1');

onMounted(() => {
  shown.value = true;
});

/*
 * 设置输入域(input/textarea)光标的位置
 * @param {HTMLInputElement/HTMLTextAreaElement} elem
 * @param {Number} index
 */
function setCursorPosition(elem: HTMLTextAreaElement, index: number) {
  const val = elem.value;
  const len = val.length; // 超过文本长度直接返回

  if (len < index) return;

  setTimeout(function () {
    elem.focus();

    if (elem.setSelectionRange) {
      // 标准浏览器
      elem.setSelectionRange(index, index);
    }
  }, 10);
}

const handleKeyDown = (e: KeyboardEvent) => {
  const target: any = e.target;

  if (e.keyCode == 13 && (e.ctrlKey || e.metaKey)) {
    const selectionStart = target.selectionStart;

    input.value =
      input.value.slice(0, selectionStart) +
      '\n' +
      input.value.slice(selectionStart);

    setCursorPosition(target, selectionStart + 1);
  } else if (e.keyCode == 13) {
    e.preventDefault();

    onTranslate();
  }
};

const addToNote = (isAddToPhrase: boolean) => {
  console.log('addToNote', translatedData.value);
  if (translatedData.value) {
    props.addToNoteHandler(
      isAddToPhrase,
      input.value,
      !isAddToPhrase && translatedData.value.targetResp?.length
        ? translatedData.value.targetResp[0].targetContent?.join(' ')
        : translatedContent.value,
      translatedData.value
    );
    if (isAddToPhrase) {
      message.success(t('translate.addToWordPhraseSuccessTip'));
    }
  }
};

const translateStore = useTranslateStore();

const allowToAddNote = computed(() => translateStore.allowToAddNote);

const feedbackVisible = ref(true);

const layoutRef = ref<InstanceType<typeof Layout>>();

const getPopupContainer = (triggerNode: HTMLElement) => {
  return triggerNode.closest('.js-translate-tippy-viewer') || document.body;
};

const submitFeedback = async () => {
  const { requestId = '' } = translatedData.value || {};
  if (requestId) {
    await correctTranslate({
      requestId,
      sources: currentTab.value,
    });
  }
  reporter.report(
    {
      event_code: EventCode.readpaperPopupTranslateFeedbackClick,
    },
    {
      sources: currentTab.value,
      tran_content: trim(input.value),
      tran_result: translatedContent.value,
    }
  );
  feedbackVisible.value = false;
  message.success(t('translate.reportSuccessTip'));
};

const CtrlOrCommand = getPlatformKey() === 'win32' ? 'Ctrl' : 'Command';

const showGlossaryButton = computed(() => {
  return [
    TranslateTabKey.idea,
    TranslateTabKey.youdao,
    TranslateTabKey.baidu,
  ].includes(currentTab.value as TranslateTabKey);
});

const onGlossaryCheckedChange = (checked: boolean) => {
  fetch();
};

const { fontSize } = useTranslateFontSize();

const handleCopy = () => {
  let content = translatedContent.value;
  if (translatedData.value?.targetResp?.length) {
    content = translatedData.value.targetResp
      .map((item) => item.targetContent.join(' '))
      .join(' ');
  }
  copyToPaste(content);
};
</script>
<style lang="less" scoped>
.translate-ctrl {
  position: absolute;
  top: 9px;
  left: 54px;
  z-index: 9;
  color: #fff;
  opacity: 0.9;
}

.translate-content {
  padding: 8px 16px 16px;
  line-height: 24px;
  height: 100%;
  max-height: calc(100vh - 32px);

  // overflow: auto;
  .text {
    background: #f5f7fa;
    border-radius: 4px;
    font-size: 14px;
    font-family: Lato-Regular, Lato;
    font-weight: 400;
    color: #1d2129;
    line-height: 22px;
    padding: 8px;
    cursor: text;

    .ellipsis {
      text-overflow: ellipsis;
      overflow: hidden;
      -webkit-line-clamp: 6;
      -webkit-box-orient: vertical;
      display: -webkit-box;
    }
  }

  .translate-part-txt,
  .translate-text {
    // font-size: 16px;
    font-family:
      PingFangSC-Regular,
      PingFang SC;
    font-weight: 400;
    color: #1d2129;
    // line-height: 26px;
    cursor: text;
  }

  .translate-text {
    min-height: 40px;
    margin-top: 8px;
    padding: 4px 0 8px;
  }

  .divide {
    display: flex;
    justify-content: space-between;
    margin-top: 16px;
    align-items: center;
    font-size: 14px;
    height: 22px;
    line-height: 22px;

    .title {
      font-weight: 600;
      color: #1d2129;
      margin-right: 12px;
    }

    .line {
      height: 1px;
      background: #e4e7ed;
      flex: 1;
    }

    .icon {
      color: #4e5969;
      margin-left: 12px;
      cursor: pointer;
    }
  }

  .error {
    text-align: center;
    color: #43464a;
    cursor: pointer;
  }

  :deep(.ant-spin-container) {
    &::after {
      background-color: #fff;
    }
  }

  .tabs {
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

  .feedback {
    margin-top: 6px;
    color: #86919c;
    font-size: 14px;
    line-height: 16px;
  }

  .translate-word {
    margin-top: 14px;

    .pronunciation-item + .pronunciation-item {
      margin-top: 8px;
    }
  }

  .translate-part-list {
    padding: 0;
    margin: 12px 0 14px;
  }

  .translate-part-item {
    display: flex;
    // font-size: 16px;
    // line-height: 26px;
    color: #1d2229;

    .translate-part {
      // width: 28px;
      margin-right: 4px;
    }

    .translate-part-txt {
      margin-left: 12px;

      p {
        margin: 0;
      }
    }
  }

  .translate-part-item + .translate-part-item {
    margin-top: 8px;
  }
}
</style>
