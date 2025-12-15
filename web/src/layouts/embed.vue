<template>
  <a-config-provider :locale="isEnUS ? enUS : zhCN">
    <a-layout style="height: 100%">
      <a-layout-content :style="{ flex: 1 }">
        <slot>
          <router-view />
        </slot>
      </a-layout-content>
    </a-layout>
  </a-config-provider>
</template>
<script lang="ts" setup>
import { watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import enUS from 'ant-design-vue/es/locale/en_US';
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

const { locale } = useI18n();

// 语言管理
const { isCurrentLanguage } = useLanguage();

// 基于 proto 枚举的语言判断计算属性
const isEnUS = computed(() => isCurrentLanguage(Language.EN_US));

watch(locale, () => {
  // locale.value 现在是标准格式（'en-US', 'zh-CN'），需要转换为 dayjs 支持的格式
  const lang = isCurrentLanguage(Language.EN_US) ? 'en' : 'zh-cn';
  dayjs.locale(lang);
});
</script>
