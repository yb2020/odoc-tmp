<template>
  <div class="i18n-button">
    <a
      :class="`${!isCurrentLanguage(Language.EN_US) ? 'normal' : ''}`"
      @click="switchLanguage(Language.EN_US)"
    >EN</a>
    <a-divider type="vertical" />
    <a
      :class="`${!isCurrentLanguage(Language.ZH_CN) ? 'normal' : ''}`"
      @click="switchLanguage(Language.ZH_CN)"
    >中</a>
  </div>
</template>
<script lang="ts" setup>
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '~/src/hooks/useLanguage';

const { isCurrentLanguage, changeLanguage } = useLanguage();

const switchLanguage = (lang: Language) => {
  changeLanguage(lang);
}

</script>
<style lang="less">
  .i18n-button {
    font-size: 13px;
    line-height: 20px;
    padding: 6px 16px;
    font-weight: 500;
    white-space: nowrap;
    a.normal {
      color: rgba(255, 255, 255, 0.85);
    }
    
  }

</style>