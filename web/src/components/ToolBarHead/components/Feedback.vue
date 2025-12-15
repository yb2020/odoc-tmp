<template>
  <Layout
    v-if="isWebEN"
    :title="$t('viewer.help')"
    :tippy-handler="tippyHandler"
    :no-ding="true"
    :style="{ width: '400px' }"
  >
    <div class="feedback-wrap__email js-interact-drag-ignore">
      <div class="tip">
        Please contact us at
      </div>
      <div class="email">
        {{ envStore.viewerConfig.feedbackEmail }}
      </div>
      <a-button
        size="large"
        type="primary"
        @click.stop="copyToPaste(envStore.viewerConfig.feedbackEmail!)"
      >
        <copy-outlined /> Copy email
      </a-button>
    </div>
  </Layout>
  <!-- <Layout
    v-else
    title="帮助"
    :tippy-handler="tippyHandler"
    :no-ding="true"
    :style="{ width: '400px' }"
  >
    <div class="feedback-wrap">
       <img src="https://readpaper.com/doc/assets/img/Qrcode.png">
      <div class="tip1">
        微信扫码添加小管家
      </div>
      <div class="tip2">
        您有任何使用问题、新功能建议<br>都可以向小管家吐槽或留言反馈
      </div>
      <div class="text-xs text-rp-neutral-6 my-2">
        Revision: {{ revision || '' }}
      </div>
      <a-button
        type="primary"
        @click.stop="goToFeedbackPage"
      >
        <edit-outlined /> 我要反馈
      </a-button>
      <a :href="userGuideLink" target="_blank" class="user-guide-link">{{ $t('home.upload.errorDetail.guide') }}</a>
    </div>
  </Layout> -->
</template>

<script lang="ts" setup>
import Layout from '../../Tippy/Layout/index.vue';
import { EditOutlined, CopyOutlined } from '@ant-design/icons-vue';
import { goPathPage } from '~/src/common/src/utils/url';
import { getDomainOrigin } from '~/src/util/env';
import { useEnvStore } from '~/src/stores/envStore';
import { copyToPaste } from '@/util/copy';
import { useI18n } from 'vue-i18n'
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';
import { computed } from 'vue';

declare global {
  interface Window {
    __REVISION__: string;
  }
}

const emit = defineEmits<{
  (event: 'close'): void;
}>();

const tippyHandler = (event: 'close' | 'ding' | 'unding' | 'lock') => {
  emit('close');
};

const goToFeedbackPage = () => {
  goPathPage(`${getDomainOrigin()}/feedback`);
};

const envStore = useEnvStore();
const revision = import.meta.env.REVISION || '';

const { isEnUS } = useLanguage();
const isWebEN = isEnUS; // 保持向后兼容的命名

const { locale } = useI18n();

// 语言管理
const { isCurrentLanguage } = useLanguage();

const userGuideLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide';
  }
  return '/docs/guide';
});
</script>

<style lang="less" scoped>
.feedback-wrap {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 48px 0;
  line-height: 24px;
  img {
    width: 200px;
  }
  .tip1 {
    margin-top: 8px;
    color: #73716f;
    font-size: 13px;
  }

  .tip2 {
    margin-top: 8px;
    font-size: 14px;
    color: #262625;
  }
  &__email {
    text-align: center;
    padding: 24px 40px;
    .tip {
      font-size: 20px;
      line-height: 30px;
      margin-bottom: 22px;
    }
    .email {
      background-color: #f7f8fa;
      border-radius: 2px;
      color: #1f71e0;
      font-weight: 500;
      font-size: 16px;
      line-height: 26px;
      padding: 8px 16px;
      margin-bottom: 32px;
    }
  }
}
.user-guide-link {
  display: block;
  color: #1f71e0;
  text-decoration: none;
  padding: 8px 0;
  
  &:hover {
    color: #0d5cbf;
    text-decoration: underline;
  }
}
</style>
