<template>
  <component
    :is="!inline ? Popover : 'div'"
    :placement="placement || 'top'"
  >
    <template #content>
      <a :href="userGuideLink" target="_blank" class="user-guide-link">{{ $t('home.upload.errorDetail.guide') }}</a>
      <!-- <div class="wrapper">
        <p class="text-rp-neutral-3 mb-3 text-center">
          {{ $t('common.tips.wechatContact') }}
        </p>
        <img
          class="code"
          src="https://readpaper.com/doc/assets/img/Qrcode.png"
          alt="ReadPaper Wechat"
        >
      </div> -->
    </template>
    <slot />
  </component>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { Popover } from 'ant-design-vue';
import { useI18n } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

defineProps<{
  inline?: boolean;
  placement?: string;
}>();

const { locale } = useI18n();

// 语言管理
const { isCurrentLanguage } = useLanguage();

// 创建用户指南链接的计算属性
const userGuideLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide';
  }
  return '/docs/guide';
});
</script>

<style scoped>
.wrapper {
  width: 168px;
  text-align: center;
}

.code {
  width: 100%;
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
