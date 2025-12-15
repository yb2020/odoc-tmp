// 默认的订阅计划数据
export const defaultSubPlanInfos = [
  {
    ai: {
      copilot: {
        isEnable: true,
        models: [
          {
            creditCost: '1',
            isEnable: true,
            isFree: true,
            key: 'gpt-4o-mini',
            name: 'gpt 4o mini',
          },
          {
            creditCost: '1',
            isEnable: true,
            isFree: true,
            key: 'claude-3.5-sonnet',
            name: 'claude 3.5 sonnet',
          },
        ],
      },
    },
    base: {
      isEnableAddOnCredit: false,
      isEnableSubAddOnCredit: false,
      maxAddOnCreditSubCountOfMonth: 0,
      subAddOnCreditInfo: {
        addOnCredit: '0',
        credit: '0',
        originalPrice: '0',
        currency: 'usd',
        duration: 0,
        name: 'UNAVAILABLE',
        price: '0',
        type: 0,
      },
      subInfo: {
        addOnCredit: '0',
        credit: '10000',
        originalPrice: '10000',
        currency: 'usd',
        duration: 1,
        name: 'Free',
        price: '0',
        type: 1,
      },
    },
    description: 'Free Version',
    docs: {
      docUploadMaxPageCount: 100,
      docUploadMaxSize: '20',
      maxStorageCapacity: '200',
    },
    isFree: true,
    name: 'Free',
    note: {
      isNoteExtract: true,
      isNoteManage: true,
      isNotePdfDownload: true,
      isNoteSummary: true,
      isNoteWord: true,
    },
    translate: {
      aiTranslationCreditCost: '1',
      fullTextTranslateCreditCost: '1',
      isAiTranslation: true,
      isFullTextTranslate: true,
      isOcr: true,
      isWordTranslate: true,
      ocrCreditCost: '1',
      wordTranslateCreditCost: '1',
    },
    type: 1,
  },
  {
    ai: {
      copilot: {
        isEnable: true,
        models: [
          {
            creditCost: '1',
            isEnable: true,
            isFree: false,
            key: 'gpt-4o-mini',
            name: 'gpt 4o mini',
          },
          {
            creditCost: '1',
            isEnable: true,
            isFree: true,
            key: 'claude-3.5-sonnet',
            name: 'claude 3.5 sonnet',
          },
        ],
      },
    },
    base: {
      isEnableAddOnCredit: true,
      isEnableSubAddOnCredit: true,
      maxAddOnCreditSubCountOfMonth: 6,
      subAddOnCreditInfo: {
        addOnCredit: '25000',
        credit: '0',
        currency: 'usd',
        duration: 0,
        name: '$10 for 250 credits',
        price: '1000',
        type: 3,
      },
      subInfo: {
        addOnCredit: '20000',
        credit: '50000',
        currency: 'usd',
        duration: 1,
        name: 'Pro',
        price: '1200',
        type: 2,
      },
    },
    description: 'Professional Version',
    docs: {
      docUploadMaxPageCount: 100,
      docUploadMaxSize: '100',
      maxStorageCapacity: '30960',
    },
    isFree: false,
    name: 'Pro',
    note: {
      isNoteExtract: true,
      isNoteManage: true,
      isNotePdfDownload: true,
      isNoteSummary: true,
      isNoteWord: true,
    },
    translate: {
      aiTranslationCreditCost: '1',
      fullTextTranslateCreditCost: '1',
      isAiTranslation: true,
      isFullTextTranslate: true,
      isOcr: true,
      isWordTranslate: true,
      ocrCreditCost: '1',
      wordTranslateCreditCost: '1',
    },
    type: 2,
  },
]

//先放这里
export const toNumberString = (number: number) => {
  return number / 100
}

export const currencyToSymbol = (currency: string) => {
  if (!currency) return ''
  currency = currency.toUpperCase()
  switch (currency) {
    case 'USD':
      return '$'
    case 'CNY':
      return '¥'
    case 'RMB':
      return '¥'
    case 'JPY':
      return '¥'
    case 'KRW':
      return '₩'
    default:
      return currency
  }
}
