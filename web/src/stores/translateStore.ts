// store

import { defineStore } from 'pinia';
import { retryable } from '@idea/aiknowledge-special-util/retryable';
import { isInElectron } from '~/src/util/env';
import {
  getCustomerTranslateTabList,
  getTranslateTabList,
} from '../api/translate';
import { TranslateTabKey, UniTranslateResp } from '~/src/api/translate';
import {
  TxConfig,
  AliConfig,
  GoogleConfig,
  DeepLConfig,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/CustomTranslateInterface';
import { PageSelectText } from '@idea/pdf-annotate-viewer';
import { Nullable } from '../typings/global';
// 移除未使用的导入，因为我们已经取消了会员权限检查

export const enum AllowGoogleTranslate {
  unknown,
  yes,
  no,
}

export const enum GoogleTranslateStyle {
  none,
  default,
  custom,
}

export { TranslateTabKey };

// export const DEFAULT_TRANSLATE_TAB_KEY = TranslateTabKey.idea;
export const DEFAULT_TRANSLATE_TAB_KEY = TranslateTabKey.google; //默认翻译设置为谷歌

export enum AiBeansType {
  unknown,
  none, // 无
  ready, // 准备
  free, // 免豆
  trial, // 试用
  standard, // 标准
}

export interface TranslateState {
  pdfId: string;
  content: {
    origin: string;
    ocrTranslate: UniTranslateResp | null;
    ocrChannel?: TranslateTabKey;
  };
  multiSegment: boolean;
  isExistingAnnotation: boolean;
  allowToAddNote: boolean;
  tabs: TranslateTabKey[];
  cusTabs: {
    name: string;
    type: TranslateTabKey;
    id: string;
    verified: boolean;
  }[];
  txConfig: TxConfig;
  aliConfig: AliConfig;
  googleConfig: GoogleConfig;
  deeplConfig: DeepLConfig;
  googleConfigVersion: GoogleTranslateStyle;
  extraInfo: {
    selections?: Nullable<PageSelectText[]>;
    ocr?: {
      text: TranslateState['content'];
      addOcrNote(translation: string): Promise<void>;
    };
    translateData?: UniTranslateResp | null;
  };
  aiTranslateConfig: {
    aiBeansType: AiBeansType;
    trialCount: number;
    cost: number;
    inBeta: boolean;
  };
}

export const LSKEY_CURRENT_TRANSLATE_TAB = 'pdf-annotate/2.0/translateTab';
export const LSKEY_CURRENT_TRANSLATE_TAB_RESET_TIME =
  'pdf-annotate/2.0/translateTabResetTime';

export const resetLSCurrentTranslateTabKey = () => {
  try {
    const currentTab = localStorage.getItem(LSKEY_CURRENT_TRANSLATE_TAB);
    /**
     * 1、凡是切换为自定义翻译引擎的，记录原本的选项，不再调整回idea
     * 2、切换为有道/百度等的，依然默认是idea
     */

    if (
      currentTab === TranslateTabKey.idea ||
      (currentTab !== TranslateTabKey.youdao &&
        currentTab !== TranslateTabKey.baidu)
    ) {
      return;
    }
    const currentTabResetTime = parseInt(
      localStorage.getItem(LSKEY_CURRENT_TRANSLATE_TAB_RESET_TIME) || '0',
      10
    );
    // 每天凌晨2:00将currentTab重置为默认值
    const now = new Date();
    // 今天凌晨2:00
    const today = new Date(
      now.getFullYear(),
      now.getMonth(),
      now.getDate(),
      2,
      0,
      0
    );
    if (now.getTime() < today.getTime()) {
      today.setDate(today.getDate() - 1);
    }
    if (currentTabResetTime !== today.getTime()) {
      localStorage.setItem(
        LSKEY_CURRENT_TRANSLATE_TAB_RESET_TIME,
        today.getTime().toString()
      );
      localStorage.setItem(LSKEY_CURRENT_TRANSLATE_TAB, TranslateTabKey.idea);
    }
  } catch (error) {
    console.error(error);
  }
};

export const useTranslateStore = defineStore('translate', {
  state: (): TranslateState => ({
    pdfId: '',
    content: {
      origin: '',
      ocrTranslate: null,
    },
    multiSegment: false,
    isExistingAnnotation: false,
    allowToAddNote: false,
    tabs: [DEFAULT_TRANSLATE_TAB_KEY],
    googleConfigVersion: GoogleTranslateStyle.none,
    cusTabs: [],
    aliConfig: <AliConfig>{
      id: '',
      verified: false,
      name: '阿里云翻译',
      aliAccessKeyId: '',
      aliAccessKeySecret: '',
      createDate: '',
      aliInterfaceVersion: 'general',
    },
    txConfig: <TxConfig>{
      id: '',
      verified: false,
      name: '腾讯翻译君',
      txSecretId: '',
      txSecretKey: '',
      createDate: '',
    },
    googleConfig: <GoogleConfig>{
      id: '',
      verified: false,
      name: '谷歌翻译',
      googleApiKey: '',
      createDate: '',
    },
    deeplConfig: <DeepLConfig>{
      id: '',
      verified: false,
      name: 'DeepL翻译',
      deepLApi: 'Free',
      deepLKey: '',
      deepLFormality: 'default',
    },
    extraInfo: {},
    aiTranslateConfig: {
      aiBeansType: AiBeansType.unknown,
      trialCount: 0,
      cost: 0,
      inBeta: false,
    },
  }),
  getters: {
    accessToAiTranslate: () => {
      // 移除会员权限检查，改为积分制，所有用户都可以使用AI翻译功能
      return true;
    },
  },
  actions: {
    setContent(payload: TranslateState['content']) {
      this.content = payload;
    },
    setPdfId(pdfId: string) {
      this.pdfId = pdfId;
    },
    enableAllowToAddNote(enable: boolean) {
      this.allowToAddNote = enable;
    },

    async getTabs() {
      const tabs = await retryable(getTranslateTabList, 2, 200);

      if (tabs.tabs?.length) {
        this.tabs = tabs.tabs as TranslateTabKey[];
      }

      if (isInElectron()) {
        const customerTabs = await retryable(
          getCustomerTranslateTabList,
          2,
          200
        );

        this.txConfig.verified = false;
        this.aliConfig.verified = false;
        this.googleConfig.verified = false;
        this.deeplConfig.verified = false;
        this.googleConfigVersion = GoogleTranslateStyle.none;
        this.cusTabs = [];
        if (customerTabs.aliConfig) {
          this.aliConfig = customerTabs.aliConfig;
          const isExist = this.cusTabs.some(
            (item) => item.type === TranslateTabKey.ali
          );
          if (!this.tabs.includes(TranslateTabKey.ali)) {
            this.tabs.push(TranslateTabKey.ali);
          }
          if (!isExist) {
            this.cusTabs.push({
              name: this.aliConfig.name,
              type: TranslateTabKey.ali,
              verified: this.aliConfig.verified,
              id: this.aliConfig.id?.toString() as string,
            });
          }
        }

        if (customerTabs.txConfig) {
          this.txConfig = customerTabs.txConfig;

          if (!this.tabs.includes(TranslateTabKey.tencent)) {
            this.tabs.push(TranslateTabKey.tencent);
          }
          const isExist = this.cusTabs.some(
            (item) => item.type === TranslateTabKey.tencent
          );

          if (!isExist) {
            this.cusTabs.push({
              name: this.txConfig.name,
              type: TranslateTabKey.tencent,
              verified: this.txConfig.verified,
              id: this.txConfig.id?.toString() as string,
            });
          }
        }
        if (customerTabs.googleConfig) {
          this.googleConfig = customerTabs.googleConfig;
          this.googleConfigVersion =
            customerTabs.googleConfig.googleApiKey?.length === 0
              ? GoogleTranslateStyle.default
              : GoogleTranslateStyle.custom;

          if (!this.tabs.includes(TranslateTabKey.google)) {
            this.tabs.push(TranslateTabKey.google);
          }
          const isExist = this.cusTabs.some(
            (item) => item.type === TranslateTabKey.google
          );

          if (!isExist) {
            this.cusTabs.push({
              name: this.googleConfig.name,
              type: TranslateTabKey.google,
              verified: this.googleConfig.verified,
              id: this.googleConfig.id?.toString() as string,
            });
          }
        }

        if (customerTabs.deepLConfig) {
          this.deeplConfig = customerTabs.deepLConfig;

          if (!this.tabs.includes(TranslateTabKey.deepl)) {
            this.tabs.push(TranslateTabKey.deepl);
          }

          const isExist = this.cusTabs.some(
            (item) => item.type === TranslateTabKey.deepl
          );

          if (!isExist) {
            this.cusTabs.push({
              name: this.deeplConfig.name,
              type: TranslateTabKey.deepl,
              verified: this.deeplConfig.verified,
              id: this.deeplConfig.id?.toString() as string,
            });
          }
        }
      }
      return this.tabs;
    },

    setExtraInfo(info: TranslateState['extraInfo']) {
      this.extraInfo = info;
    },

    async initAitranslateConfig() {
      const config: TranslateState['aiTranslateConfig'] = {
        aiBeansType: AiBeansType.standard,
        trialCount: 0,
        cost: 0,
        inBeta: false,
      };
      if (config.aiBeansType === AiBeansType.trial) {
        this.aiTranslateConfig.aiBeansType =
          config.trialCount > 0 ? AiBeansType.trial : AiBeansType.none;
      } else {
        this.aiTranslateConfig.aiBeansType = config.aiBeansType;
      }
      this.aiTranslateConfig.trialCount = config.trialCount;
      this.aiTranslateConfig.cost = config.cost;
      this.aiTranslateConfig.inBeta = config.inBeta;
    },
  },
});
