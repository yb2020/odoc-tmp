/* eslint-disable @typescript-eslint/ban-ts-comment */
import IconPremium from '~common/assets/images/premium/icon-premium.svg?skipsvgo';
import IconDocLimit from '~common/assets/images/vip/icon-file.svg';
import IconDocExport from '~common/assets/images/vip/icon-pdfnote.svg';
import IconFullTrans from '~common/assets/images/vip/icon-full-translate.svg';
import IconAIBeans from '~common/assets/images/vip/icon-aibeans.svg';
import IconMaster from '~common/assets/images/vip/icon-master.svg';
import IconReviewer from '~common/assets/images/vip/icon-reviewer.svg';
import IconGroup from '~common/assets/images/vip/icon-group.svg';
// import IconRocket from '~common/assets/images/vip/icon-rocket.svg'
import IconDocLimitH5 from '~common/assets/images/vip/h5/icon-file.svg';
import IconDocExportH5 from '~common/assets/images/vip/h5/icon-pdfnote.svg';
import IconFullTransH5 from '~common/assets/images/vip/h5/icon-full-translate.svg';
import IconAIBeansH5 from '~common/assets/images/vip/h5/icon-aibeans.svg';
import IconOCRTransH5 from '~common/assets/images/vip/h5/icon-ocr.svg';
import IconAICopilotH5 from '~common/assets/images/vip/h5/icon-copilot.svg';
import IconAIPolishH5 from '~common/assets/images/vip/h5/icon-polish.svg';
import IconGroupH5 from '~common/assets/images/vip/h5/icon-group.svg';
import IconCalendarH5 from '~common/assets/images/vip/h5/icon-calendar.svg';
import IconPeopleH5 from '~common/assets/images/vip/h5/icon-people.svg';
// import IconNewH5 from '~common/assets/images/vip/h5/icon-new.svg'
import BadgeFree from '~common/assets/images/premium/badge-free.svg';
import BadgeStd from '~common/assets/images/premium/badge-std.svg';
import BadgePro from '~common/assets/images/premium/badge-pro.svg';
import BadgeEnt from '~common/assets/images/premium/badge-ent.svg';
import {
  VipPrivilege,
  VipType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/user/VipUserInterface';
import {
  AiBeanCountResponse,
  VipPayType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { isNumber } from 'lodash-es';
import { CloseCircleFilled, CheckOutlined } from '@ant-design/icons-vue';
import { formatSize } from '@common/utils/file';

export type VipTypePayable = Exclude<
  VipType,
  VipType.FREE | VipType.UNRECOGNIZED
>;

export const VipType2PayType = {
  [VipType.STANDARD]: VipPayType.STANDARD,
  [VipType.PROFESSIONAL]: VipPayType.PROFESSIONAL,
  [VipType.OUTSTANDING]: VipPayType.OUTSTANDING,
};

export const VipType2ElementName = {
  [VipType.STANDARD]: 'std_pay',
  [VipType.PROFESSIONAL]: 'pro_pay',
  [VipType.ENTERPRISE]: 'ent_pay',
  [VipType.OUTSTANDING]: 'ost_pay',
};

export enum PrivilegeTypes {
  DOC_LIMIT = 'limit',
  DOC_EXPORT = 'export',
  DOC_ATTACHEMNTS = 'attachment',
  TRANS_FULLTEXT = 'fulltext',
  TRANS_SELECTTEXT = 'selecttext',
  TRANS_OCR = 'ocr',
  AI_BEANS = 'beans',
  AI_COPILOT = 'copilot',
  AI_REVISE = 'revise',
  AI_MASTER = 'master',
  AI_REVIEW = 'review',
  AI_TRANSLATE = 'ai_translate',
  GROUP_SHARE = 'gshare_limit',
  GROUP_NOTES = 'gnote_limit',
  GROUP_LIMIT = 'glimit',
  GROUP_NUM_LIMIT = 'gnum_limit',
  OTHER_READTIMES = 'readtimes',
  OTHER_EXPERIENCE = 'experience',
}

export const PrivilegeTypes2Field = {
  [PrivilegeTypes.DOC_LIMIT]: 'documentCountLimit',
  [PrivilegeTypes.DOC_EXPORT]: 'exportEnable',
  [PrivilegeTypes.DOC_ATTACHEMNTS]: 'docAttachmentTotalSpace',
  [PrivilegeTypes.TRANS_FULLTEXT]: [
    'fullTextTranslationCountLimit',
    'fullTextTranslatePageCountLimit',
  ],
  [PrivilegeTypes.TRANS_SELECTTEXT]: 'wordSelectionTranslateCountLimit',
  [PrivilegeTypes.TRANS_OCR]: 'ocrTranslateCountLimit',
  [PrivilegeTypes.AI_BEANS]: 'aiBeanCountLimit',
  [PrivilegeTypes.AI_COPILOT]: [
    'readingEnable',
    'aiReadingDocumentPageCountLimit',
  ],
  [PrivilegeTypes.AI_REVISE]: 'polishEnable',
  [PrivilegeTypes.AI_MASTER]: 'polishEnable',
  [PrivilegeTypes.AI_REVIEW]: 'polishEnable',
  [PrivilegeTypes.AI_TRANSLATE]: 'aiTranslateEnable',
  [PrivilegeTypes.GROUP_SHARE]: 'groupDocumentCountLimit',
  [PrivilegeTypes.GROUP_NOTES]: 'groupDocumentNoteCountLimit',
  [PrivilegeTypes.GROUP_LIMIT]: [
    'groupCreateCountLimit',
    'groupJoinCountLimit',
  ],
  [PrivilegeTypes.GROUP_NUM_LIMIT]: 'groupMemberCountLimit',
  [PrivilegeTypes.OTHER_READTIMES]: 'learningStatisticsEnable',
  [PrivilegeTypes.OTHER_EXPERIENCE]: 'newFeatureEnable',
};

/**
 * @deprecated
 */
export const PrivilegeTypes2BeansCostField = {
  [PrivilegeTypes.AI_COPILOT]: 'aiBeanCountCostEveryQuesion',
  [PrivilegeTypes.AI_REVISE]: 'aiBeanCountCostEveryPolish',
  [PrivilegeTypes.AI_MASTER]: 'aiBeanCountCostEveryAbstractGeneration',
  [PrivilegeTypes.AI_REVIEW]: 'aiBeanCountCostEveryReview',
  [PrivilegeTypes.AI_TRANSLATE]: 'aiTranslateConsumeBeanCount',
};

export const PrivilegeTypes2CostField = {
  [PrivilegeTypes.AI_COPILOT]: 'readingAskQuestion',
  [PrivilegeTypes.AI_REVISE]: 'polishZH',
  [PrivilegeTypes.AI_MASTER]: 'polishAbstractGenerate',
  [PrivilegeTypes.AI_REVIEW]: 'polishReviewAll',
  [PrivilegeTypes.AI_TRANSLATE]: 'aiTranslateConsumeBeanCount',
};

export const PrivilegeTypes2Txt: Record<string, any[]> = {
  [PrivilegeTypes.GROUP_SHARE]: [100],
  [PrivilegeTypes.GROUP_LIMIT]: [3],
  [PrivilegeTypes.GROUP_NUM_LIMIT]: [20],
};

export enum PrivilegeRenderTypes {
  COUNT = 'count',
  BETA = 'beta',
  RADIO = 'radio',
  BYTES = 'bytes',
  PAGES = 'pages',
  PAGES_PERTIME = 'pagespertime',
  BEANS_PERTIME = 'beanspertime',
}

export enum PrivilegeCountUnits {
  ONE = 'one',
  PIECE = 'piece',
  PERSON = 'user',
  PERDAY = 'pertimesday',
  PERWEEK = 'pertimesweek',
  PERBEANSWEEK = 'perbeansweek',
  TPL_DOC = 'tpl_doc',
  TPL_GROUP = 'tpl_group',
}

export interface PrivilegeConifg {
  key: string;
  title: string;
  icon: string;
  iconH5?: string;
  tip: 1 | 0;
  type: PrivilegeRenderTypes | string;
  typeTip?: PrivilegeRenderTypes | string;
  typeDesc?: PrivilegeRenderTypes | string;
  unit?: PrivilegeCountUnits;
}

/* eslint-disable no-template-curly-in-string */
export const PremiumPrivileges = [
  {
    prefix: 'common.premium.privilege.management',
    key: 'tt',
    title: '文献管理',
    privileges: [
      {
        key: PrivilegeTypes.DOC_LIMIT,
        title: '文献上限',
        icon: IconDocLimit,
        iconH5: IconDocLimitH5,
        type: PrivilegeRenderTypes.COUNT,
      },
      {
        key: PrivilegeTypes.DOC_EXPORT,
        title: 'PDF带笔记导出',
        icon: IconDocExport,
        iconH5: IconDocExportH5,
        type: PrivilegeRenderTypes.RADIO,
      },
      // {
      //   key: PrivilegeTypes.DOC_ATTACHEMNTS,
      //   title: '文献附件',
      //   type: PrivilegeRenderTypes.BYTES,
      // },
    ] as PrivilegeConifg[],
  },
  {
    prefix: 'common.premium.privilege.translation',
    key: 'tt',
    title: '文献翻译',
    privileges: [
      {
        key: PrivilegeTypes.TRANS_FULLTEXT,
        unit: PrivilegeCountUnits.PERWEEK,
        title: '全文翻译',
        icon: IconFullTrans,
        iconH5: IconFullTransH5,
        type: PrivilegeRenderTypes.COUNT,
        typeDesc: PrivilegeRenderTypes.PAGES_PERTIME,
      },
      // {
      //   key: PrivilegeTypes.TRANS_SELECTTEXT,
      //   unit: PrivilegeCountUnits.PERDAY,
      //   title: '划词翻译',
      //   type: PrivilegeRenderTypes.COUNT,
      // },
      {
        key: PrivilegeTypes.TRANS_OCR,
        unit: PrivilegeCountUnits.PERWEEK,
        title: 'OCR翻译',
        icon: 'icon-translate',
        iconH5: IconOCRTransH5,
        type: PrivilegeRenderTypes.COUNT,
      },
    ] as PrivilegeConifg[],
  },
  {
    prefix: 'common.premium.privilege.ai',
    key: 'tt',
    title: 'AI功能',
    privileges: [
      {
        key: PrivilegeTypes.AI_BEANS,
        unit: PrivilegeCountUnits.PERBEANSWEEK,
        title: 'AI豆',
        icon: IconAIBeans,
        iconH5: IconAIBeansH5,
        tip: 1,
        type: PrivilegeRenderTypes.COUNT,
      },
      {
        key: PrivilegeTypes.AI_TRANSLATE,
        title: 'AI翻译',
        icon: '',
        iconH5: IconAICopilotH5,
        type: PrivilegeRenderTypes.RADIO,
        typeTip: PrivilegeRenderTypes.BEANS_PERTIME,
      },
      {
        key: PrivilegeTypes.AI_COPILOT,
        title: 'AI辅读',
        icon: '',
        iconH5: IconAICopilotH5,
        type: PrivilegeRenderTypes.RADIO,
        typeTip: PrivilegeRenderTypes.BEANS_PERTIME,
        typeDesc: PrivilegeRenderTypes.PAGES,
      },
      {
        key: PrivilegeTypes.AI_REVISE,
        title: 'AI润色',
        icon: 'icon-polish-label',
        iconH5: IconAIPolishH5,
        type: PrivilegeRenderTypes.RADIO,
        typeTip: PrivilegeRenderTypes.BEANS_PERTIME,
      },
      {
        key: PrivilegeTypes.AI_MASTER,
        title: '标题/摘要神器',
        icon: IconMaster,
        type: PrivilegeRenderTypes.RADIO,
        typeTip: PrivilegeRenderTypes.BEANS_PERTIME,
      },
      {
        key: PrivilegeTypes.AI_REVIEW,
        title: 'AI审稿',
        icon: IconReviewer,
        type: PrivilegeRenderTypes.RADIO,
        typeTip: PrivilegeRenderTypes.BEANS_PERTIME,
      },
    ] as PrivilegeConifg[],
  },
  {
    prefix: 'common.premium.privilege.group',
    key: 'tt',
    title: '团队功能',
    privileges: [
      {
        key: PrivilegeTypes.GROUP_SHARE,
        title: '小组文献上传上限',
        icon: IconDocLimit,
        iconH5: IconDocLimitH5,
        unit: PrivilegeCountUnits.PIECE,
        type: PrivilegeRenderTypes.COUNT,
      },
      // {
      //   key: PrivilegeTypes.GROUP_NOTES,
      //   title: '小组文献笔记上限',
      //   unit: PrivilegeCountUnits.TPL_DOC,
      //   type: PrivilegeRenderTypes.COUNT,
      // },
      {
        key: PrivilegeTypes.GROUP_LIMIT,
        title: '小组创建',
        icon: IconGroup,
        iconH5: IconGroupH5,
        // unit: PrivilegeCountUnits.TPL_GROUP,
        unit: PrivilegeCountUnits.ONE,
        type: PrivilegeRenderTypes.COUNT,
      },
      {
        key: PrivilegeTypes.GROUP_NUM_LIMIT,
        title: '单个小组最大人数',
        icon: IconGroup,
        iconH5: IconPeopleH5,
        unit: PrivilegeCountUnits.PERSON,
        type: PrivilegeRenderTypes.COUNT,
      },
    ] as PrivilegeConifg[],
  },
  {
    prefix: 'common.premium.privilege.others',
    key: 'tt',
    title: '其他功能',
    privileges: [
      {
        key: PrivilegeTypes.OTHER_READTIMES,
        title: '阅读时长统计',
        icon: 'calendar',
        iconH5: IconCalendarH5,
        type: PrivilegeRenderTypes.RADIO,
      },
      // {
      //   key: PrivilegeTypes.OTHER_EXPERIENCE,
      //   title: '新功能优先内测',
      //   icon: IconRocket,
      //   iconH5: IconNewH5,
      //   type: PrivilegeRenderTypes.RADIO,
      // },
    ] as PrivilegeConifg[],
  },
];

// export const PrivilegeType2Config = keyBy(
//   // eslint-disable-next-line @typescript-eslint/no-unused-vars
//   PremiumPrivileges.reduce((arr, { title, key, privileges, ...rest }) => {
//     arr.push(...privileges.map((y) => ({ ...y, ...rest })))

//     return arr
//   }, [] as PrivilegeConifg[]),
//   'key',
// )

export const PremiumCorePrivileges = {
  [VipType.STANDARD]: {
    privileges: [
      'premium.privilege.ai.copilot',
      {
        key: 'premium.privilege.translation.fulltext',
        label: 'premium.versionCore.fulltext',
      },
      'premium.privilege.translation.ocr',
      'premium.versionCore.beans',
    ],
  },
  [VipType.PROFESSIONAL]: {
    privileges: [
      'premium.versionCore.contain',
      {
        key: 'premium.versionCore.aireviseBilingual',
        label: 'premium.versionCore.eureka',
      },
      'premium.privilege.ai.review',
      'premium.versionCore.beans',
    ],
  },
  [VipType.ENTERPRISE]: {
    privileges: [
      'premium.versionCore.contain',
      'premium.versionCore.customGroup',
      'premium.versionCore.customFeature',
      'premium.versionCore.customApi',
    ],
  },
  [VipType.OUTSTANDING]: {
    privileges: [
      'premium.versionCore.contain',
      'premium.versionCore.aicopilotUnlimited',
      {
        key: 'premium.versionCore.aireviseUnlimited',
        label: 'premium.versionCore.eureka',
      },
      'premium.versionCore.doclimit',
    ],
  },
};

export const PremiumVipPreferences = {
  [VipType.FREE]: {
    key: 'free',
    icon: null,
    badge: BadgeFree,
    color: '#A8AFBA',
    sloganI18n: '',
    slogan: '',
  },
  [VipType.STANDARD]: {
    key: 'standard',
    icon: IconPremium,
    badge: BadgeStd,
    color: '#6AA9EC',
    sloganI18n: 'common.premium.wordings.slogan.standard',
    slogan: '有AI，轻松读论文',
  },
  [VipType.PROFESSIONAL]: {
    key: 'profession',
    icon: IconPremium,
    badge: BadgePro,
    color: '#5A77EB',
    sloganI18n: 'common.premium.wordings.slogan.professional',
    slogan: '有AI，轻松做科研',
  },
  [VipType.ENTERPRISE]: {
    key: 'enterprise',
    icon: IconPremium,
    badge: BadgeEnt,
    color: '#24304B',
    sloganI18n: 'common.premium.wordings.slogan.enterprise',
    slogan: '有AI，轻松带学生',
  },
  [VipType.OUTSTANDING]: {
    key: 'outstanding',
    icon: IconPremium,
    badge: BadgeEnt,
    color: '#24304B',
    sloganI18n: 'common.premium.wordings.slogan.enterprise',
    slogan: '有AI，轻松带学生',
  },
  [VipType.UNRECOGNIZED]: {} as any,
};

export const DefaultVipList: VipPrivilege[] = [
  {
    vipType: VipType.FREE,
    aiTranslateEnable: false,
    paySubject: '',
    payTotalAmount: '0',
    exportEnable: false,
    documentCountLimit: 200,
    docAttachmentTotalSpace: 0,
    singleDocumentSizeLimit: (100 * 1024 * 1024) / 8,
    fullTextTranslationCountLimit: 3,
    fullTextTranslatePageCountLimit: 30,
    wordSelectionTranslateCountLimit: 50,
    ocrTranslateCountLimit: 0,
    groupCreateCountLimit: 1,
    groupJoinCountLimit: 5,
    groupDocumentCountLimit: 20,
    groupDocumentNoteCountLimit: 10,
    groupMemberCountLimit: 5,
    aiReadingDocumentPageCountLimit: 0,
    aiBeanCountLimit: 0,
    duplicatePaperCheckingCountLimit: 0,
    readingEnable: false,
    polishEnable: false,
    originalPayTotalAmount: '',
    costLabelUrl: '',
    costLabelTip: '',
    disseminateText: '',
    disseminateTextes: [],
    learningStatisticsEnable: false,
    newFeatureEnable: false,
    aiTranslateConsumeBeanCount: 1,
    groupBuyPayTotalAmount: '',
    discount: '80',
  },
  {
    vipType: VipType.STANDARD,
    aiTranslateEnable: true,
    paySubject: '',
    payTotalAmount: '9900',
    exportEnable: true,
    documentCountLimit: 1000,
    docAttachmentTotalSpace: (500 * 1024 * 1024) / 8,
    singleDocumentSizeLimit: (100 * 1024 * 1024) / 8,
    fullTextTranslationCountLimit: 10,
    fullTextTranslatePageCountLimit: 50,
    wordSelectionTranslateCountLimit: -1,
    ocrTranslateCountLimit: 100,
    groupCreateCountLimit: 10,
    groupJoinCountLimit: 20,
    groupDocumentCountLimit: 100,
    groupDocumentNoteCountLimit: 100,
    groupMemberCountLimit: 20,
    aiReadingDocumentPageCountLimit: 50,
    aiBeanCountLimit: 200,
    duplicatePaperCheckingCountLimit: 1,
    readingEnable: true,
    polishEnable: false,
    originalPayTotalAmount: '19900',
    costLabelUrl: '',
    costLabelTip: '',
    disseminateText: '',
    disseminateTextes: [],
    learningStatisticsEnable: true,
    newFeatureEnable: true,
    aiTranslateConsumeBeanCount: 1,
    groupBuyPayTotalAmount: '',
    discount: '80',
  },
  {
    vipType: VipType.PROFESSIONAL,
    aiTranslateEnable: true,
    paySubject: '',
    payTotalAmount: '36000',
    exportEnable: true,
    documentCountLimit: 2000,
    docAttachmentTotalSpace: (2 * 1024 * 1024 * 1024) / 8,
    singleDocumentSizeLimit: (100 * 1024 * 1024) / 8,
    fullTextTranslationCountLimit: 30,
    fullTextTranslatePageCountLimit: 100,
    wordSelectionTranslateCountLimit: -1,
    ocrTranslateCountLimit: 300,
    groupCreateCountLimit: 20,
    groupJoinCountLimit: 50,
    groupDocumentCountLimit: 1000,
    groupDocumentNoteCountLimit: -1,
    groupMemberCountLimit: 50,
    aiReadingDocumentPageCountLimit: 100,
    aiBeanCountLimit: 600,
    duplicatePaperCheckingCountLimit: 1,
    readingEnable: true,
    polishEnable: true,
    originalPayTotalAmount: '79800',
    costLabelUrl: '',
    costLabelTip: '',
    disseminateText: '',
    disseminateTextes: [],
    learningStatisticsEnable: true,
    newFeatureEnable: true,
    aiTranslateConsumeBeanCount: 1,
    groupBuyPayTotalAmount: '',
    discount: '80',
  },
  {
    vipType: VipType.ENTERPRISE,
    aiTranslateEnable: true,
    paySubject: '',
    payTotalAmount: '120000',
    exportEnable: true,
    documentCountLimit: 2000,
    docAttachmentTotalSpace: (2 * 1024 * 1024 * 1024) / 8,
    singleDocumentSizeLimit: (100 * 1024 * 1024) / 8,
    fullTextTranslationCountLimit: 30,
    fullTextTranslatePageCountLimit: 100,
    wordSelectionTranslateCountLimit: -1,
    ocrTranslateCountLimit: 300,
    groupCreateCountLimit: 20,
    groupJoinCountLimit: 50,
    groupDocumentCountLimit: 1000,
    groupDocumentNoteCountLimit: -1,
    groupMemberCountLimit: 50,
    aiReadingDocumentPageCountLimit: 100,
    aiBeanCountLimit: 600,
    duplicatePaperCheckingCountLimit: 1,
    readingEnable: true,
    polishEnable: true,
    originalPayTotalAmount: '150000',
    costLabelUrl: '',
    costLabelTip: '',
    disseminateText: '',
    disseminateTextes: [],
    learningStatisticsEnable: true,
    newFeatureEnable: true,
    aiTranslateConsumeBeanCount: 1,
    groupBuyPayTotalAmount: '',
    discount: '80',
  },
  {
    vipType: VipType.OUTSTANDING,
    aiTranslateEnable: true,
    paySubject: '',
    payTotalAmount: '99900',
    exportEnable: true,
    documentCountLimit: 5000,
    docAttachmentTotalSpace: (2 * 1024 * 1024 * 1024) / 8,
    singleDocumentSizeLimit: (100 * 1024 * 1024) / 8,
    fullTextTranslationCountLimit: -1,
    fullTextTranslatePageCountLimit: 100,
    wordSelectionTranslateCountLimit: -1,
    ocrTranslateCountLimit: -1,
    groupCreateCountLimit: 100,
    groupJoinCountLimit: 100,
    groupDocumentCountLimit: 2000,
    groupDocumentNoteCountLimit: -1,
    groupMemberCountLimit: 100,
    aiReadingDocumentPageCountLimit: 100,
    aiBeanCountLimit: 1500,
    duplicatePaperCheckingCountLimit: 1,
    readingEnable: true,
    polishEnable: true,
    originalPayTotalAmount: '160000',
    costLabelUrl: '',
    costLabelTip: '',
    disseminateText: '',
    disseminateTextes: [],
    learningStatisticsEnable: true,
    newFeatureEnable: true,
    polishZH: 0,
    polishEN: 0,
    polishZh2En: 0,
    polishTitleGenerate: 0,
    polishAbstractGenerate: 0,
    readingAskQuestion: 0,
    polishAbstractRewrite: 0,
    polishAbstractTruingFramework: 0,
    aiTranslateConsumeBeanCount: -1,
    groupBuyPayTotalAmount: '',
    discount: '80',
  },
].map((x) => {
  return {
    polishZH: 2,
    polishEN: 2,
    polishZh2En: 2,
    polishTitleGenerate: 6,
    polishAbstractGenerate: 8,
    readingAskQuestion: 8,
    polishAbstractRewrite: 8,
    polishAbstractTruingFramework: 0,
    polishReviewAll: 30,
    ...x,
  };
});

export const StepwiseVipTypes = [
  VipType.FREE,
  VipType.STANDARD,
  VipType.PROFESSIONAL,
  // [VipType.PROFESSIONAL, VipType.ENTERPRISE],
  VipType.OUTSTANDING,
];

export const compareVipType = (a: VipType, b: VipType) => {
  const aIdx = StepwiseVipTypes.findIndex((v) => [v].flat().includes(a));
  const bIdx = StepwiseVipTypes.findIndex((v) => [v].flat().includes(b));

  if (aIdx >= 0 && bIdx >= 0) {
    return aIdx > bIdx ? 1 : aIdx < bIdx ? -1 : 0;
  }

  return Number.NaN;
};

export const getLowerVipType = (vipType: VipType) => {
  const idx = StepwiseVipTypes.findIndex((v) => [v].flat().includes(vipType));
  if (idx > 0) {
    return StepwiseVipTypes[idx - 1];
  }

  return undefined;
};

export const getHigherVipType = (vipType: VipType) => {
  const idx = StepwiseVipTypes.findIndex((v) => [v].flat().includes(vipType));
  if (idx >= 0) {
    return StepwiseVipTypes[idx + 1];
  }

  return undefined;
};

export const getPrivilegeTip = (
  context: { i18n: any },
  config: PrivilegeConifg,
  costs: AiBeanCountResponse,
  vipPrivilege?: VipPrivilege
) => {
  const { i18n } = context;
  const { key, typeTip } = config;
  const field =
    PrivilegeTypes2BeansCostField[
      key as keyof typeof PrivilegeTypes2BeansCostField
    ];
  const fields = Array.isArray(field) ? (field as string[]) : [field];
  // @ts-ignore
  let values = fields.map((k) => costs[k]) as string[];

  if (values[0] === undefined) {
    values = [vipPrivilege?.[field as keyof VipPrivilege]] as string[];
  }

  console.log('getPrivilegeTip', key, field, typeTip, values, costs);

  switch (typeTip) {
    case PrivilegeRenderTypes.BEANS_PERTIME:
      return i18n.t('common.premium.wordings.beanspertime', values);
    default:
      return `${values[0]}`;
  }
};

export const getPrivilegeTxt = (
  context: { i18n: any },
  config: PrivilegeConifg,
  vipPrivilege?: VipPrivilege,
  useIcon = true
) => {
  const { i18n } = context;
  const { type, key } = config;
  const field = PrivilegeTypes2Field[key as PrivilegeTypes];
  // if (key === PrivilegeTypes.AI_TRANSLATE) {
  //   debugger
  // }
  const fieldCost =
    PrivilegeTypes2CostField[key as keyof typeof PrivilegeTypes2CostField];
  const fields = Array.isArray(field) ? field : [field];
  const resUnlimited = i18n.t(`common.premium.wordings.unlimit`);
  const resAble = useIcon ? [CheckOutlined, {}] : '✔️'; // i18n.t('common.premium.wordings.able')
  if (useIcon) {
    resAble.toString = () => '✔️';
  }
  const resNull = [
    CloseCircleFilled,
    {
      class: '!text-rp-neutral-3',
    },
  ] as const;
  resNull.toString = () => '';
  const values =
    (vipPrivilege?.vipType !== VipType.OUTSTANDING ? PrivilegeTypes2Txt : {})[
      key
      // @ts-ignore
    ] || fields.map((k) => vipPrivilege?.[k]);
  // @ts-ignore
  const valueCost = fieldCost ? Number(vipPrivilege?.[fieldCost]) : null;

  switch (type) {
    case PrivilegeRenderTypes.RADIO:
      return values[0]
        ? isNumber(valueCost) && valueCost <= 0
          ? i18n.t('common.aibeans.freeuse')
          : resAble
        : resNull;
    case PrivilegeRenderTypes.BETA:
      return values[0]
        ? isNumber(valueCost) && valueCost <= 0
          ? i18n.t('common.aibeans.freeuse')
          : (i18n.t('common.premium.wordings.gray') as string)
        : resNull;
    case PrivilegeRenderTypes.BYTES:
      return values[0] === -1
        ? resUnlimited
        : values[0] > 0
          ? formatSize(values[0])
          : resNull;
    case PrivilegeRenderTypes.COUNT:
      return values[0] === -1
        ? resUnlimited
        : values[0] > 0
          ? getPrivilegeUnitTxt(context, config, values)
          : resNull;
    default:
      return `${values[0]}`;
  }
};

export const getPrivilegeDesc = (
  context: { i18n: any },
  config: PrivilegeConifg,
  vipPrivilege?: VipPrivilege,
  short = false
) => {
  const { i18n } = context;
  const { typeDesc, key } = config;
  const field = PrivilegeTypes2Field[key as PrivilegeTypes];
  const fields = Array.isArray(field) ? field : [field];

  // @ts-ignore
  const values = fields.map((k) => vipPrivilege?.[k]);
  if (!values[0]) {
    return '';
  }
  if (short) {
    return i18n.t(`common.premium.units.tpl_pagesshort`, values) as string;
  }
  switch (typeDesc) {
    case PrivilegeRenderTypes.PAGES:
      return i18n.t(`common.premium.units.tpl_pages`, values) as string;
    case PrivilegeRenderTypes.PAGES_PERTIME:
      return i18n.t(`common.premium.units.tpl_pagespertime`, values) as string;
    default:
      return `${values[0]}`;
  }
};

export const getPrivilegeUnitTxt = (
  context: { i18n: any },
  config: PrivilegeConifg,
  values: any[]
): string => {
  const { i18n } = context;
  switch (config.unit) {
    case PrivilegeCountUnits.PERBEANSWEEK:
      return `${i18n.t('common.premium.wordings.fillUp')} ${values[0]} ${
        i18n.t(`common.premium.units.${config.unit}`, values[0] > 1 ? 2 : 1) ??
        ''
      }`;
    case PrivilegeCountUnits.ONE:
    case PrivilegeCountUnits.PERSON:
    case PrivilegeCountUnits.PIECE:
    case PrivilegeCountUnits.PERDAY:
    case PrivilegeCountUnits.PERWEEK:
      return `${values[0]} ${
        i18n.t(`common.premium.units.${config.unit}`, values[0] > 1 ? 2 : 1) ??
        ''
      }`;
    case PrivilegeCountUnits.TPL_DOC:
      return `${i18n.t('common.premium.units.tpl_doc', values)}`;
    case PrivilegeCountUnits.TPL_GROUP:
      return `${i18n.t('common.premium.units.tpl_group', values)}`;
    default:
      return `${values[0]}`;
  }
};
