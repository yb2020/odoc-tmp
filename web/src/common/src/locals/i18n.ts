import { createI18n, I18nOptions, LocaleMessages } from 'vue-i18n';
import forIn from 'lodash-es/forIn';
import merge from 'lodash-es/merge';
import { isInOverseaseElectron } from '../utils/env';
import commonZhJSON from './files/zh-CN.json';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';

export type I18nCommonMessageType = typeof commonZhJSON;

// 定义语言类型映射，将 Language 枚举映射到标准格式语言代码（RFC 5646）
export const LanguageCodeMap = {
  [Language.EN_US]: 'en-US',  // 0 → 'en-US'
  [Language.ZH_CN]: 'zh-CN',  // 1 → 'zh-CN'
};

// 反向映射，使用标准格式语言代码作为键
export const CodeToLanguageMap = {
  'en-US': Language.EN_US,  // 'en-US' → 0
  'zh-CN': Language.ZH_CN,  // 'zh-CN' → 1
};
export function createMessages<T extends I18nOptions['messages']>(
  modules: Record<string, any>
) {
  const messages: { [x: string]: T } = {};
  forIn(modules, (value, key) => {
    const idx = key.lastIndexOf('/');
    const lang = key.substring(idx + 1, key.length - 5);
    const prev = messages[lang] || {};
    messages[lang] = merge(prev, value.default as T);
  });
  delete messages.lang;
  return messages;
}

const currentLocal = 'en-US'; // 统一默认为英文

const numberFormats: I18nOptions['numberFormats'] = {
  'en-US': {
    integer: {
      style: 'decimal',
      useGrouping: true,
    },
    percent: {
      style: 'percent',
      useGrouping: false,
    },
  },
  'zh-CN': {
    integer: {
      style: 'decimal',
      useGrouping: false,
    },
    percent: {
      style: 'percent',
      useGrouping: false,
    },
  },
} as const;

const datetimeFormats: I18nOptions['datetimeFormats'] = {
  'en-US': {
    short: {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    },
  },
  'zh-CN': {
    short: {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    },
  },
} as const;

function createVueI18n<T extends LocaleMessages<any>>(
  modules: Record<string, any>
) {
  const commonMessages = createMessages<I18nCommonMessageType>(
    import.meta.glob('./files/*.json', { eager: true })
  );
  const localMessages = createMessages<T>(modules);

  const messages = merge({}, commonMessages, localMessages) as Record<
    string,
    I18nCommonMessageType & T
  >;
  const i18n = createI18n({
    locale: currentLocal,
    legacy: false,
    globalInjection: true,
    fallbackLocale: 'en-US',
    messages,
    numberFormats,
    datetimeFormats,
  });
  return i18n;
}

export default createVueI18n;

export const createCommonVueI18n = () => {
  const i18n = createVueI18n(
    import.meta.glob('./files/*.json', { eager: true })
  );
  return i18n;
};

export const setGlobalLocale = (
  i18n: ReturnType<typeof createVueI18n>,
  locale?: string
) => {
  locale =
    locale ||
    (/^zh/.test(navigator.language)
      ? 'zh-CN'
      : 'en-US');
  i18n.global.locale.value = locale;
};
