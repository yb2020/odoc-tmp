import { defineStore } from 'pinia';
import {
  $SiteOptionData,
  TDKOptions,
  getTKDOptions,
  getWebSiteInfo,
} from '@common/api/seo';
import { SeoPageType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/seo/SEOCommon';

type SeoState = {
  tdkOptions: {
    default: TDKOptions;
  };
  siteOptions: Partial<$SiteOptionData>;
};

export const useSeoStore = defineStore('seo', {
  state: (): SeoState => ({
    tdkOptions: {
      default: {
        tplTitle: '',
        tplKeywords: '',
        tplDescription: '',
      },
    },
    siteOptions: {},
  }),
  actions: {
    async getDefaultTDK() {
      const defaultTDK = await getTKDOptions({
        type: SeoPageType.HOME,
      });
      if (defaultTDK) {
        this.tdkOptions.default = defaultTDK;
      } else {
        console.log('获取默认TDK失败');
      }
    },
    async getWebOptions() {
      const res = await getWebSiteInfo();
      if (res) {
        this.siteOptions = res;
      } else {
        console.log('获取网站配置失败');
      }
    },
  },
});
