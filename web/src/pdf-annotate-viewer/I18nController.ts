import enUSJson from './locals/en-US.json';
import zhCNJson from './locals/zh-CN.json';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';

export default class I18nController {
  private stores;
  private lang: string;
  
  constructor(lang?: string) {
    // 保持 stores 的键与 JSON 文件名一致
    this.stores = {
      'EN_US': enUSJson,
      'ZH_CN': zhCNJson,
    } as const;
    
    // 设置默认语言为英文
    this.lang = lang || 'EN_US';
    
  }

  t(key: string, data?: {[key: string]: string | number}) {
    const arr = key.split('.');
    const message = this.stores[this.lang];
    let tmp: any = '';
    for (let i = 0; i < arr.length; i++) {
      if (i === 0) {
        tmp = message[arr[i] as keyof typeof message];
      } else {
        tmp = tmp[arr[i]];
      }
    }
    if (data && tmp) {
      Object.keys(data).forEach((k) => {
        tmp = tmp.replace(`{${k}}`, data[k]);
      });
    }
    return tmp || key
  }
}