<template>
  <a-popover
    v-model:visible="visible"
    :title="null"
    :getPopupContainer="getPopupContainer"
    @visibleChange="onVisibleChange"
  >
    <template #content>
      <div class="w-80">
        <p>{{ $t('translate.aiTranslateTip1') }}</p>
        <p>{{ $t('translate.aiTranslateTip2') }}</p>
      </div>
    </template>
    <slot />
  </a-popover>
</template>
<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useTranslateStore, AiBeansType } from '~/src/stores/translateStore';
import AIBeans from '@common/components/AIBean/index.vue';

const visible = ref(false);

const getPopupContainer = (triggerNode: HTMLElement) => {
  return triggerNode.closest('.js-translate-tippy-viewer') || document.body;
};

const translateStore = useTranslateStore();

const aiTranslateConfig = computed(() => {
  return translateStore.aiTranslateConfig;
});

// translateStore.initAitranslateConfig();

const onVisibleChange = (visible: boolean) => {
  if (visible) {
    translateStore.initAitranslateConfig();
  }
};
</script>

<style scoped>
.polish-card {
  background: linear-gradient(90deg, #2173e1 0%, #6dded6 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
</style>
