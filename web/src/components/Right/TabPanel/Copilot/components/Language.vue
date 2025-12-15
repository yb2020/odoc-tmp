<template>
  <a-select
    :value="userSettings.copilotLanguage"
    size="small"
    class="copilot-language"
    @change="handleLanguageChange"
  >
    <a-select-option :value="LangType.CHINESE">
      {{
        $t('aiCopilot.useChinese')
      }}
    </a-select-option>
    <a-select-option :value="LangType.ENGLISH">
      {{
        $t('aiCopilot.useEnglish')
      }}
    </a-select-option>
  </a-select>
</template>
<script setup lang="ts">
import { SelectProps } from 'ant-design-vue';
import useRemoteSettings from '~/src/hooks/UserSettings/useRemoteUserSettings';
import { LangType } from '~/src/stores/copilotType';

const { userSettings, saveRemoteUserSettings } = useRemoteSettings();

const handleLanguageChange: SelectProps['onChange'] = (value) => {
  saveRemoteUserSettings({
    copilotLanguage: value as LangType,
  });
};
</script>

<style lang="less" scoped>
.copilot-language {
  :deep(.ant-select-selector) {
    background-color: var(--site-theme-bg-secondary) !important;
    border-color: var(--site-theme-border) !important;
    color: var(--site-theme-text-primary) !important;
  }
  
  :deep(.ant-select-arrow) {
    color: var(--site-theme-text-secondary) !important;
  }
  
  :deep(.ant-select-selection-item) {
    color: var(--site-theme-text-primary) !important;
  }
  
  &:hover {
    :deep(.ant-select-selector) {
      border-color: var(--site-theme-brand) !important;
    }
  }
}
</style>
<style lang="less">
.ant-select-dropdown {
  background-color: var(--site-theme-bg-secondary) !important;
  
  .ant-select-item {
    color: var(--site-theme-text-secondary) !important;
    
    &-option-active:not(.ant-select-item-option-disabled) {
      background-color: var(--site-theme-bg-mute) !important;
    }
    
    &-option-selected:not(.ant-select-item-option-disabled) {
      background-color: var(--site-theme-brand-light) !important;
      color: var(--site-theme-text-primary) !important;
    }
  }
}
</style>
