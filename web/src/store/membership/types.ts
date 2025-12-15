// 模型定义
export interface Model {
  creditCost: string;
  isEnable: boolean;
  isFree: boolean;
  key: string;
  name: string;
}

// AI功能定义
export interface AiFeature {
  copilot: {
    isEnable: boolean;
    models: Model[];
  };
}

// 订阅信息定义
export interface SubscriptionInfo {
  addOnCredit: string;
  credit: string;
  currency: string;
  duration: number;
  name: string;
  price: string;
  type: number;
}

// 基础信息定义
export interface BaseFeature {
  isEnableAddOnCredit: boolean;
  isEnableSubAddOnCredit: boolean;
  maxAddOnCreditSubCountOfMonth: number;
  subAddOnCreditInfo: SubscriptionInfo;
  subInfo: SubscriptionInfo;
}

// 文档功能定义
export interface DocFeature {
  docUploadMaxPageCount: number;
  docUploadMaxSize: string;
  maxStorageCapacity: string;
}

// 笔记功能定义
export interface NoteFeature {
  isNoteExtract: boolean;
  isNoteManage: boolean;
  isNotePdfDownload: boolean;
  isNoteSummary: boolean;
  isNoteWord: boolean;
}

// 翻译功能定义
export interface TranslateFeature {
  aiTranslationCreditCost: string;
  fullTextTranslateCreditCost: string;
  isAiTranslation: boolean;
  isFullTextTranslate: boolean;
  isOcr: boolean;
  isWordTranslate: boolean;
  ocrCreditCost: string;
  wordTranslateCreditCost: string;
}

// 订阅计划信息
export interface SubPlanInfo {
  ai: AiFeature;
  base: BaseFeature;
  description: string;
  docs: DocFeature;
  isFree: boolean;
  name: string;
  note: NoteFeature;
  translate: TranslateFeature;
  type: number;
}

// 会员信息
export interface MembershipInfo {
  [key: string]: any;
}

// 基础信息
export interface BaseInfo {
  [key: string]: any;
}