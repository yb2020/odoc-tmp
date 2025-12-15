<template>
  <PerfectScrollbar
    v-if="shown"
    class="translate-content"
  >
    <VipForbidden :current-tab="currentTab as TranslateTabKey">
      <div v-if="store.content.origin">
        <div class="flex justify-end items-center gap-3 mb-3">
          <FontSizeSelect />
          <GlossaryButton
            v-if="showGlossaryButton"
            :getPopupContainer="getPopupContainer"
            @change="onGlossaryCheckedChange"
          />
        </div>
        <div>
          <InputOrigin
            ref="inputOriginRef"
            :onTranslate="onTranslate"
            :fontSize="fontSize"
            :glossary-list="translatedData?.glossaryList"
            :add-to-note-handler="addToNoteHandler"
            @translate="onTranslate"
          />
          <ExtraIcons
            class="mt-4"
            :input="input"
            :translated-content="translatedContent"
            :translated-data="translatedData"
            :add-to-note-handler="addToNoteHandler"
          />
          <a-spin :spinning="fetchState.pending && !isUsingSSE">
            <div
              v-if="fetchState.error || sseError"
              class="error py-3"
              @click="fetch"
            >
              <redo-outlined /> {{ (fetchState.error || sseError)?.message }}
            </div>
            <template v-else-if="translatedData">
              <WordResult
                v-if="translatedData?.targetResp?.length"
                :translatedData="translatedData"
                :fontSize="fontSize"
              />
              <div
                v-else
                class="translate-text js-interact-drag-ignore"
                :style="{ fontSize: fontSize + 'px', lineHeight: '1.4' }"
              >
                {{ translatedContent }}
              </div>
            </template>
            <template v-else-if="isUsingSSE && !showSSEResult">
              <Typing
                :text="sseText"
                :is-pending="sseLoading"
                class="text-black mt-3"
                @typing:finished="onTypingFinished"
              />
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
      </div>
      <Empty v-else />
    </VipForbidden>
  </PerfectScrollbar>
</template>
<script lang="ts" setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue';
import useFetch from '~/src/hooks/useFetch';
import { correctTranslate, UniTranslateResp } from '@/api/translate';
import trim from 'lodash-es/trim';

import {
  useTranslateStore,
  TranslateTabKey,
} from '~/src/stores/translateStore';
import { RedoOutlined, DislikeOutlined } from '@ant-design/icons-vue';
import reporter from '@idea/aiknowledge-report';
import { PageType, EventCode, reportTranslateClick } from '~/src/api/report';
import { message } from 'ant-design-vue';
import { useI18n } from 'vue-i18n';
import {
  useTranslateApi,
  useTranslateFontSize,
  useTranslateSSE,
} from '~/src/hooks/useTranslation';
import GlossaryButton from '@/components/Translate/Glossary/Button.vue';
import FontSizeSelect from '@/components/Translate/FontSize.vue';
import ExtraIcons from './ExtraIcons.vue';
import Empty from './Empty.vue';
import InputOrigin from './InputOrigin.vue';
import WordResult from './WordResult.vue';
import VipForbidden from './VipForbidden.vue';
import Typing from '@common/components/TypingTxt/index.vue';
import { useUserStore } from '~/src/common/src/stores/user';
// import { useAIBeans } from '@common/hooks/useAIBeans'; //去除积分制后，不再需要刷新积分

const props = defineProps<{
  pdfId: string;
  // width: number;
  tippyHandler: (event: 'ding' | 'close' | 'unding' | 'lock') => void;
  addToNoteHandler: (
    isPhrase: boolean,
    phrase: string,
    translation: string,
    translationRes: UniTranslateResp
  ) => void;
  fixPlacement(): void;
  currentTab: string;
  updateCurrentTab: (tab: TranslateTabKey) => void;
}>();

const inputOriginRef = ref<InstanceType<typeof InputOrigin>>();

const input = computed<string>(
  () => inputOriginRef.value?.getOriginInput() || ''
);

const store = useTranslateStore();

store.initAitranslateConfig();

const translatedContent = computed(
  () => translatedData?.value?.targetContent || ''
);

const editing = ref(false);

const reportTranslate = () =>
  reportTranslateClick({
    page_type: PageType.note,
    type_parameter: props.pdfId || '',
    tran_content: input.value,
    sources: props.currentTab,
  });

const { t } = useI18n();

watch(
  () => props.currentTab,
  (newTab, oldTab) => {
    // 重置状态，确保切换标签时不会保留旧的错误状态
    translatedData.value = undefined;
    sseError.value = null;
    showSSEResult.value = false;
    sseText.value = '';
    
    // 重新获取翻译结果
    fetch();
    reportTranslate();
  }
);

const isUsingSSE = computed(() => {
  return props.currentTab === TranslateTabKey.ai;
});

// const { refresh: refreshAIBeans } = useAIBeans(); //去除积分制后，不再需要刷新积分

const showSSEResult = ref(false);
const { data: translatedData, run: fetchTranslate } = useTranslateApi();
const {
  data: sseTranslatedData,
  sseText,
  loading: sseLoading,
  startPollingTranslate,
  error: sseError,
} = useTranslateSSE();
const { fetchState, fetch } = useFetch(async () => {
  const value = trim(input.value);

  if (!value) {
    return;
  }

  translatedData.value = undefined;
  feedbackVisible.value = false;

  if (isUsingSSE.value) {
    showSSEResult.value = false;
    await startPollingTranslate(value).promise;
    //refreshAIBeans(); //去除积分制后，不再需要刷新积分
  } else {
    // 移除会员权限检查回调，因为已经改为积分制
    await fetchTranslate(props.currentTab as TranslateTabKey, value);
    // 刷新积分
    const userStore = useUserStore();
    userStore.refreshUserCredits()
  }

  editing.value = false;
  feedbackVisible.value = true;
  props.fixPlacement();
}, false);
const onTypingFinished = () => {
  if (!sseTranslatedData.value) {
    return;
  }
  showSSEResult.value = true;
  translatedData.value = sseTranslatedData.value! || {};
  console.log('onTypingFinished', sseTranslatedData.value);
  // translatedData.value.glossaryList = [
  //   {
  //     originalText: '111',
  //     translationText: '222',
  //     start: 2,
  //     end: 4,
  //   },
  // ];
  // translatedData.value.targetContent = '333';
  // 刷新积分（AI翻译）
  const userStore = useUserStore();
  userStore.refreshUserCredits()
};

const onTranslate = () => {
  fetch();
};

watch(
  () => store.content,
  (newVal) => {
    inputOriginRef.value?.updateOriginInput(newVal.origin);

    console.log('newVal', newVal, store.extraInfo.translateData);

    if (newVal.ocrTranslate) {
      if (newVal.ocrChannel === TranslateTabKey.ai) {
        fetch();
      } else {
        translatedData.value = newVal.ocrTranslate;
        props.fixPlacement();
      }
      if (newVal.ocrChannel) {
        props.updateCurrentTab(newVal.ocrChannel);
      }
    } else if (store.extraInfo.translateData && newVal.origin) {
      translatedData.value = store.extraInfo.translateData;
      store.extraInfo.translateData = null;
    } else if (newVal.origin) {
      fetch();
    }
    reportTranslate();
  },
  {
    immediate: true,
  }
);

/**
 * 如果直接展示perfectscrollbar，会导致在tippy里面bar在左边而不是右边
 */
const shown = ref(false);

onMounted(() => {
  shown.value = true;
});

const feedbackVisible = ref(true);

const getPopupContainer = (triggerNode: HTMLElement) => {
  return triggerNode.closest('.js-translate-tippy-viewer') || document.body;
};

const submitFeedback = async () => {
  const { requestId = '' } = translatedData.value || {};
  if (requestId) {
    await correctTranslate({
      requestId,
      sources: props.currentTab,
    });
  }
  reporter.report(
    {
      event_code: EventCode.readpaperPopupTranslateFeedbackClick,
    },
    {
      sources: props.currentTab,
      tran_content: trim(input.value),
      tran_result: translatedContent.value,
    }
  );
  feedbackVisible.value = false;
  message.success(t('translate.reportSuccessTip'));
};

const showGlossaryButton = computed(() => {
  return [
    TranslateTabKey.idea,
    TranslateTabKey.youdao,
    TranslateTabKey.baidu,
    TranslateTabKey.ai,
    TranslateTabKey.google
  ].includes(props.currentTab as TranslateTabKey);
});

const onGlossaryCheckedChange = () => {
  fetch();
};

const { fontSize } = useTranslateFontSize();

defineExpose({
  getTranslateData: () => {
    return translatedData.value;
  },
  storeTranslateData: () => {
    store.extraInfo.translateData = translatedData.value || null;
  },
});
</script>
<style lang="less" scoped>
.translate-content {
  padding: 0px 16px 16px;
  line-height: 24px;
  height: 100%;
  max-height: calc(100vh - 32px);
  background-color: white;
  // overflow: auto;

  .translate-text {
    font-weight: 400;
    color: #1d2129;
    cursor: text;
  }

  .translate-text {
    min-height: 40px;
    margin-top: 8px;
    padding: 4px 0 8px;
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

  .feedback {
    margin-top: 6px;
    color: #86919c;
    font-size: 14px;
    line-height: 16px;
  }
}
</style>
