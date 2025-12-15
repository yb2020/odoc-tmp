<template>
  <div
    class="text-base cursor-pointer"
    @click="handleFinder"
  >
    <a-tooltip>
      <template #title>
        {{ $t('viewer.findInDocument')
        }}<component :is="getShortcutTxt(shortcut)" />
      </template>
      <search-outlined class="iconsearch" />
    </a-tooltip>
  </div>
</template>

<script lang="ts" setup>
import { togglePDFSearchViewer } from '~/src/dom/pdf';
import { SearchOutlined } from '@ant-design/icons-vue';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import { computed } from 'vue';
import { useStore } from '@/store';
import { getPlatformKey, getShortcutTxt } from '@/store/shortcuts';
import { PAGE_ROUTE_NAME } from '@/routes/type';
import { ElementClick, reportClick } from '~/src/api/report';
import useShortcuts from '~/src/hooks/useShortcuts';

const props = defineProps<{
  pdfViewInstance: ViewerController;
}>();

const handleFinder = () => {
  const isOpen = togglePDFSearchViewer(props.pdfViewInstance);
  reportClick(ElementClick.find, isOpen ? 'on' : 'off');
};

const store = useStore();
const shortcutsConfig = computed(
  () => store.state.shortcuts[PAGE_ROUTE_NAME.NOTE] || {}
);
const platformKey = getPlatformKey();
const shortcut = computed(
  () => shortcutsConfig.value.shortcuts.SEARCH.value[platformKey]
);
const opts = computed(() => ({ scope: shortcutsConfig.value.scope || 'all' }));
useShortcuts(
  shortcut,
  (event) => {
    event.preventDefault();
    handleFinder();
  },
  opts
);
</script>

<style lang="less" scoped>
.finder-container {
  .iconsearch {
    font-size: 1rem;
  }
}
</style>
