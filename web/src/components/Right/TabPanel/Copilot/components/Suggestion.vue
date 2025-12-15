<script lang="ts" setup>
import { computed, ref, watch } from 'vue';
import { DownOutlined, UpOutlined } from '@ant-design/icons-vue';
import useRemoteSettings from '~/src/hooks/UserSettings/useRemoteUserSettings';
import Custom from './Custom.vue';
import { useLocalStorage } from '@vueuse/core';
import { PDF_READER } from '~/src/common/src/constants/storage-keys';

enum TabKey {
  DEFAULT = '1',
  CUSTOM = '2',
}

const { userSettings } = useRemoteSettings();
const suggestionVisible = ref(false);
const activeKey = useLocalStorage<TabKey>(
  PDF_READER.COPILOT_SUGGESTION_TAB,
  TabKey.CUSTOM
);

// const suggestionQuestionList = ref<QuestionItem[]>([]);

const emit = defineEmits<{
  (event: 'suggestion:fill', question: string): void;
}>();

const domRef = ref<HTMLDivElement>();
const getContainer = () => domRef.value || document.body;

const handleTab = (key: TabKey) => {
  activeKey.value = key;
};

const handleSuggestionClick = (question: string) => {
  emit('suggestion:fill', question);
  suggestionVisible.value = false;
};
</script>
<template>
  <div
    ref="domRef"
    class="copilot-suggestions"
  >
    <a-dropdown
      v-model:visible="suggestionVisible"
      :trigger="['hover']"
      :get-popup-container="getContainer"
      :destroy-popup-on-hide="true"
      overlayClassName="copilot-suggestions-overlay"
    >
      <div
        class="copilot-suggestions-tips"
        @click.prevent
      >
        {{ $t('aiCopilot.suggestions') }}
        <UpOutlined
          v-if="suggestionVisible"
          class="copilot-suggestions-arrow"
        />
        <DownOutlined
          v-else
          class="copilot-suggestions-arrow"
        />
      </div>
      <template #overlay>
        <div class="copilot-suggestions-wrap">
          <div class="copilot-suggestions-tabs flex space-x-4 px-4">
            <!-- 
            <div
              :class="[
                'text-base',
                'tab',
                'cursor-pointer',
                activeKey === TabKey.DEFAULT ? 'active' : '',
              ]"
              @click="handleTab(TabKey.DEFAULT)"
            >
              {{ $t('aiCopilot.defaultTab') }}
            </div>
            -->
            <div
              :class="[
                'text-base',
                'tab',
                'cursor-pointer',
                activeKey === TabKey.CUSTOM ? 'active' : '',
              ]"
              @click="handleTab(TabKey.CUSTOM)"
            >
              {{ $t('aiCopilot.customTab') }}
            </div>
          </div>
          <KeepAlive>
            <Custom
              v-if="activeKey === TabKey.CUSTOM"
              @suggestion:fill="handleSuggestionClick"
            />
          </KeepAlive>
        </div>
      </template>
    </a-dropdown>
  </div>
</template>
<style lang="less" scoped>
.copilot-suggestions {
  &-tips {
    background-color: var(--site-theme-bg-mute);
    margin: 0 10px;
    padding: 0 20px 0 12px;
    border-radius: 24px;
    display: -webkit-box;
    overflow: hidden;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    cursor: pointer;
    position: relative;
    opacity: 0.85;
    line-height: 36px;
    color: var(--site-theme-text-primary);
    // max-width: 320px;
  }
  &-arrow {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 14px !important;
    color: var(--site-theme-text-primary);
  }
  &-wrap {
    background-color: var(--site-theme-bg-primary);
  }
  &-tabs {
    color: var(--site-theme-text-tertiary);
    border-bottom: 1px solid var(--site-theme-border);

    .tab {
      padding: 6px 2px 10px;
    }
    .active {
      color: var(--site-theme-text-primary);
      font-weight: 500;
      border-bottom: 2px solid var(--site-theme-brand);
    }
  }
}
</style>
<style lang="less">
.copilot-suggestions {
  .ant-dropdown.copilot-suggestions-overlay {
    width: 100%;
    padding: 0 10px;
    
    .ant-tooltip {
      .ant-tooltip-inner {
        background-color: var(--site-theme-bg-secondary);
        color: var(--site-theme-text-primary);
      }
      .ant-tooltip-arrow-content {
        background-color: var(--site-theme-bg-secondary);
      }
    }
  }
}
</style>
