import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRequest } from 'ahooks-vue';
// 移除未使用的导入
import {
  REQUEST_SERVICE_NAME_TRANSLATE,
} from '@common/api/const';
import { selfNoteInfo } from '../store';
import {
  useTranslateStore,
  TranslateTabKey,
  DEFAULT_TRANSLATE_TAB_KEY,
  LSKEY_CURRENT_TRANSLATE_TAB,
} from '../stores/translateStore';
// 移除未使用的导入
import { UniTranslateResp, fetchTranslate } from '../api/translate';
import { useGlossary } from './useGlossary';
import { createSharedComposable, useLocalStorage } from '@vueuse/core';
import { useDataSSE } from '@common/hooks/aitools/useClientSSE';
import { AiTranslateRequest } from 'go-sea-proto/gen/ts/translate/TextTranslate';
import { BeanScenes } from '@common/hooks/useAIBeans';
import { PDF_READER } from '../common/src/constants/storage-keys';
import { useUserStore } from '../common/src/stores/user';

export function useTranslateApi() {
  const { t } = useI18n();
  const store = useTranslateStore();
  // 移除未使用的vipStore变量

  const pdfId = computed(() => selfNoteInfo.value.pdfId);

  const { glossaryChecked } = useGlossary();

  return useRequest(
    async (
      channel: TranslateTabKey,
      content: string
      // 移除会员权限检查相关参数
    ) => {
      if (!content) {
        throw Error(t('translate.placeholder'));
      }

      let param;
      let res;

      if (!store.tabs.find((tab) => tab === channel) && !store.cusTabs.length) {
        await store.getTabs();
      }

      switch (channel) {
        case TranslateTabKey.tencent:
          param = store.txConfig;
          break;
        case TranslateTabKey.google:
          param = store.googleConfig;
          break;
        case TranslateTabKey.ali:
          param = store.aliConfig;
          break;
        case TranslateTabKey.deepl:
          param = store.deeplConfig;
          break;
        default:
          break;
      }

      try {
        // 移除会员权限检查，改为积分制
        res = await fetchTranslate({
          type: channel,
          content,
          pdfId: pdfId.value,
          param: param,
          useGlossary: glossaryChecked.value,
        });
        // 刷新积分
        const userStore = useUserStore();
        userStore.refreshUserCredits();
      } catch (err) {
        // 移除会员权限检查相关代码，只保留基本错误处理
        console.error('Translation error:', err);
      }

      return res;
    },
    {
      manual: true,
    }
  );
}

export const useTranslateSSE = () => {
  const translateData = ref<UniTranslateResp | null>(null);
  const sseText = ref<string>('');

  const { glossaryChecked } = useGlossary();

  const { startPollingData, ...rest } = useDataSSE<
    AiTranslateRequest & { text: string },
    UniTranslateResp | null,
    UniTranslateResp
  >(
    `/api/text/translate/completions`,
    BeanScenes.AI_TRANSLATE,
    translateData,
    (d) => {
      console.log('useTranslateSSE', d);
      if (d === undefined) {
        translateData.value = null;
        return;
      }
      if (d.targetResp || d.targetContent) {
        translateData.value = d;
        translateData.value.targetContent = Array.isArray(d.targetContent)
          ? d.targetContent[0]
          : d.targetContent;
        sseText.value += d.sseData ?? ' ';
      } else {
        sseText.value += d.sseData ?? '';
      }
    }
  );

  const startPollingTranslate = (content: string) => {
    sseText.value = '';
    translateData.value = null;
    return startPollingData({
      content,
      text: content,
      useGlossary: glossaryChecked.value,
      pdfId: selfNoteInfo.value.pdfId,
    });
  };

  return {
    data: translateData,
    sseText,
    startPollingTranslate,
    ...rest,
  };
};

export const LSKeyForTranslateFontSize = 'pdf-annotate/2.0/translateFontSize';

export const useTranslateFontSize = createSharedComposable(() => {
  const fontSize = useLocalStorage(LSKeyForTranslateFontSize, '16');

  return {
    fontSize,
  };
});

export const LSKeyForTranslateLock = PDF_READER.TRANSLATE_LOCK;
export const useTranslateLock = createSharedComposable(() => {
  const translateLock = useLocalStorage(LSKeyForTranslateLock, false);
  const showTranslateTabInRight = computed(() => {
    return !!translateLock.value;
  });

  return {
    translateLock,
    showTranslateTabInRight,
  };
});

export const useTranslateTabs = createSharedComposable(() => {
  const currentTab = useLocalStorage<string>(
    LSKEY_CURRENT_TRANSLATE_TAB,
    DEFAULT_TRANSLATE_TAB_KEY
  );
  console.log("use translate tabs", currentTab.value)
  const { t } = useI18n();
  const store = useTranslateStore();
  const titles: { [k in TranslateTabKey]?: string } = {
    [TranslateTabKey.youdao]: 'translate.youdao',
    [TranslateTabKey.baidu]: 'translate.baidu',
    [TranslateTabKey.idea]: 'IDEA',
    [TranslateTabKey.ai]: 'translate.ai',
    [TranslateTabKey.google]: 'Google',
  };

  const changeTab = (type: string) => {
    currentTab.value = type;
  };

  const tabs = computed(() => {
    return store.tabs.map((tab: any) => {
      if (!titles[tab as TranslateTabKey]) {
        const cusTab = store.cusTabs.find((item: any) => item.type === tab);
        if (cusTab) {
          titles[tab as TranslateTabKey] = cusTab.name;
        }
      }
      return {
        title: [
          TranslateTabKey.youdao,
          TranslateTabKey.baidu,
          TranslateTabKey.ai,
          TranslateTabKey.google,
        ].includes(tab)
          ? t(titles[tab as TranslateTabKey] || '')
          : titles[tab as TranslateTabKey],
        type: tab as TranslateTabKey,
      };
    });
  });

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

  const addSettingTab = (payload: { title: string; type: TranslateTabKey }) => {
    if (!titles[payload.type]) {
      tabs.value.push({
        title: payload.title,
        type: payload.type,
      });
    }

    initTabs();
  };

  const deleteSettingTab = ({ type }: { type: TranslateTabKey }) => {
    delete titles[type];
    currentTab.value = DEFAULT_TRANSLATE_TAB_KEY;
    for (let i = 0; i < tabs.value.length; i++) {
      if (tabs.value[i].type === type) {
        tabs.value.splice(i, 1);
        initTabs();
        break;
      }
    }
  };

  return {
    currentTab,
    tabs,
    initTabs,
    addSettingTab,
    deleteSettingTab,
    changeTab,
  };
});
