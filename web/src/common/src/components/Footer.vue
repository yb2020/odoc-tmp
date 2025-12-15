<template>
  <footer class="footer w-[1192px]">
    <div
      class="flex justify-between pt-10 pb-8 border-b border-solid border-rp-white-a-3"
    >
      <img
        :src="logo || LOGO"
        alt=""
        class="h-[42px]"
      >
      <div class="flex gap-[120px]">
        <div
          v-for="item in FooterData"
          :key="item.title"
          class="item"
        >
          <div
            class="title mb-6 font-medium text-[20px] leading-[30px] text-rp-neutral-10"
          >
            {{ $t(`common.footer.${item.i18nKey}.title`) }}
          </div>
          <div
            v-for="detail in item.detailList"
            :key="detail.subtitle"
            class="subtitle text-sm text-rp-neutral-8 mt-[10px]"
          >
            <a
              v-if="detail.icon !== 'icon-wechat'"
              :href="detail.url"
              target="_blank"
              class="flex items-center"
            ><i
               v-if="detail.icon && detail.icon !== 'zhihu'"
               :class="['aiknowledge-icon', detail.icon]"
               aria-hidden="true"
             />
              <ZhihuCircleFilled
                v-if="detail.icon && detail.icon === 'zhihu'"
                class="aiknowledge-icon icon-zhihu"
              />
              {{
                $t(
                  `common.footer.${item.i18nKey}.detailList.${detail.i18nKey}.subtitle`
                )
              }}
            </a>

            <a-popover
              v-else
              :title="null"
              placement="top"
              arrow-point-at-center
            >
              <template #content>
                <img
                  :src="detail.url"
                  alt=""
                  class="w-[100px]"
                >
              </template>
              <span class="flex items-center">
                <a-icon
                  type="wechat"
                  class="aiknowledge-icon icon-wechat"
                />{{
                  $t(
                    `common.footer.${item.i18nKey}.detailList.${detail.i18nKey}.subtitle`
                  )
                }}
              </span>
            </a-popover>
          </div>
        </div>
      </div>
    </div>
    <div
      class="font-[Noto Sans SC] py-4 flex justify-between text-xs text-rp-neutral-6 leading-[22px]"
    >
      <div class="flex items-center">
        <a
          href="http://beian.miit.gov.cn/"
          target="_blank"
        >
          备案号：{{ siteOptions.beianhao }}</a>
        <a
          target="_blank"
          href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=44030402006181"
          class="ml-2"
        ><img
          class="inline"
          src="@common/../assets/images/beian.png"
        >&nbsp;粤公网安备 44030402006181号
        </a>
      </div>

      <div class="flex items-center gap-2">
        <span> &nbsp;Copyright ©{{ new Date().getFullYear() }} </span>
        <span v-html="siteOptions.copyright" />
      </div>
    </div>
  </footer>
</template>
<script lang="ts" setup>
import { useSeoStore } from '@common/stores/seo';
import { ZhihuCircleFilled } from '@ant-design/icons-vue';
import LOGO from '~common/assets/images/logo-ft.png';

defineProps<{
  logo?: string;
}>();

const seoStore = useSeoStore();
const { siteOptions } = seoStore;

const FooterData = [
  {
    i18nKey: 'product',
    title: '产品',
    detailList: [
      {
        i18nKey: 'feature',
        subtitle: '功能概览',
        url: 'https://readpaper.com/help#tab=feature',
        icon: '',
      },
      {
        i18nKey: 'tips',
        subtitle: '使用技巧',
        url: 'https://readpaper.com/help#tab=useSkill',
        icon: '',
      },
      {
        i18nKey: 'agreement',
        subtitle: '用户协议',
        url: 'https://nuxt.cdn.readpaper.com/privacy/ReadPaper%E7%94%A8%E6%88%B7%E5%8D%8F%E8%AE%AE.pdf',
        icon: '',
      },
      {
        i18nKey: 'privacy',
        subtitle: '隐私政策',
        url: 'https://nuxt.cdn.readpaper.com/privacy/%E9%9A%90%E7%A7%81%E6%94%BF%E7%AD%96.pdf',
        icon: '',
      },
    ],
  },
  {
    i18nKey: 'team',
    title: '团队',
    detailList: [
      {
        i18nKey: 'about',
        subtitle: '关于我们',
        url: 'https://readpaper.com/help',
        icon: '',
      },
      {
        i18nKey: 'log',
        subtitle: '更新日志',
        url: 'https://readpaper.com/help#tab=changelog',
        icon: '',
      },
      {
        i18nKey: 'contact',
        subtitle: '联系我们',
        url: 'https://readpaper.com/help#tab=feedback',
        icon: '',
      },
      {
        i18nKey: 'feedback',
        subtitle: '建议反馈',
        url: 'https://support.qq.com/product/367725',
        icon: '',
      },
    ],
  },
  {
    i18nKey: 'media',
    title: '社交媒体',
    detailList: [
      {
        i18nKey: 'bilibili',
        subtitle: 'Bilibili',
        url: 'https://space.bilibili.com/1706874133',
        icon: 'icon-bilibili',
      },
      {
        i18nKey: 'zhihu',
        subtitle: '知乎',
        url: 'https://www.zhihu.com/people/readpaperlun-wen-yue-du',
        icon: 'zhihu',
      },
      {
        i18nKey: 'weChat',
        subtitle: '微信公众号',
        url: 'https://readpaper.com/doc/assets/img/Qrcode.png',
        icon: 'icon-wechat',
      },
    ],
  },
];
</script>
<style lang="less" scoped>
.footer {
  .aiknowledge-icon {
    font-size: 16px;
    height: 16px;
    line-height: 16px;
    margin-right: 5px;
  }
  .icon-bilibili {
    color: #00a5e5;
  }
  .icon-zhihu {
    color: #0064df;
  }
  .icon-wechat {
    color: #3da611;
  }
}
</style>
