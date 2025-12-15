<template>
  <div class="flex items-center justify-end text-[#4e5969] glossary">
    <a-switch
      v-model:checked="glossaryChecked"
      size="small"
      class="glossary-button"
    />
    <a-tooltip
      :getPopupContainer="getPopupContainer"
      :title="$t('glossary.tooltip')"
    >
      <div
        ref="triggerRef"
        class="cursor-pointer text-sm hover:text-[#1f71e0]"
        @click="onOpenManageTable"
      >
        <span class="mx-1">{{ $t('glossary.terms') }}</span>
        <RightOutlined class="text-[rgba(255,255,255,.45)]" />
      </div>
      <a-dropdown size="small" :getPopupContainer="getPopupContainer" overlayClassName="glossary-overlay">
        
        <template #overlay>
          <a-menu>
            <a-menu-item>
              <a href="javascript:;">个人术语库</a>
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </a-tooltip>
  </div>
</template>
<script setup lang="ts">
import { ref, watch } from 'vue';
import { RightOutlined } from '@ant-design/icons-vue';
import { createGlossaryTableTippyVue } from '@/dom/tippy';
import { useGlossary } from '~/src/hooks/useGlossary';

defineProps<{
  getPopupContainer: (triggerNode: HTMLElement) => Element;
}>();

const { glossaryChecked } = useGlossary();
const triggerRef = ref<HTMLElement | null>(null);

const onOpenManageTable = () => {
  createGlossaryTableTippyVue({
    triggerEle: triggerRef.value!,
  });
};

const emit = defineEmits<{
  (event: 'change', checked: boolean): void;
}>();

watch(glossaryChecked, (checked) => {
  emit('change', checked);
});
</script>
<style less scoped>
.glossary-button.ant-switch:not(.ant-switch-checked) {
  background-color: rgb(239, 239, 239);
}
</style>
<!-- <style less>
.glossary-overlay {
  .ant-dropdown-menu {
    background-color: #fff;
    .ant-dropdown-menu-title-content {
      color: #000;
    }
  }
}
</style> -->
