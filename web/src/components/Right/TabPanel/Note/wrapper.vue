<template>
  <a-tabs
    class="note-subtabs"
    :activeKey="activeSubTab || NoteSubTypes.Summary"
    @change="handleTabChange"
  >
    <a-tab-pane
      v-for="subTab in SubNoteTabs"
      :key="subTab.type"
      force-render
      :tab="`${$t(`viewer.subnotes.${subTab.i18n}`)}${
        subTabCounts[subTab.type] ? `(${subTabCounts[subTab.type]})` : ''
      }`"
    >
      <component
        :is="subTab.component"
        :tab="tab"
        :activeTab="activeTab"
        @counted="onCounted(subTab.type, $event)"
      />
    </a-tab-pane>
  </a-tabs>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import { NoteSubTypes } from '@/store/note/types';
import { isOwner } from '~/src/store';
import Summary from './Summary/index.vue';
import Vocabulary from './Vocabulary/index.vue';
import Annotations from './index.vue';
import { RightSideBarType } from '../type';
import { useRightSideTabSettings } from '~/src/hooks/UserSettings/useSideTabSettings';

const props = defineProps<{
  tab: RightSideBarType;
  activeTab: RightSideBarType;
}>();

const SubNoteTabs = [
  {
    type: NoteSubTypes.Summary,
    i18n: 'summary',
    component: Summary,
    hidden: !isOwner,
  },
  {
    type: NoteSubTypes.Vocabulary,
    i18n: 'vocabulary',
    component: Vocabulary,
    hidden: !isOwner,
  },
  {
    type: NoteSubTypes.Annotation,
    i18n: 'annotation',
    component: Annotations,
  },
].filter((x) => !x.hidden);

const { activeSubTab, setSideTabSetting } = useRightSideTabSettings();
const subTabCounts = reactive<{
  [k in NoteSubTypes]?: number;
}>({});

const onCounted = (subTab: NoteSubTypes, count: number) => {
  if (count >= 0) {
    subTabCounts[subTab] = count;
  }
};

const handleTabChange = (subTab: NoteSubTypes) => {
  setSideTabSetting({
    subTab,
  });
};
</script>

<style lang="less" scoped>
.note-subtabs {
  height: 100%;
  color: var(--site-theme-text-primary);

  :deep(.ant-tabs-nav) {
    margin: 3px 20px 0;
  }
  :deep(.ant-tabs-content) {
    height: 100%;
  }
  
  :deep(.ant-tabs-tab) {
    color: var(--site-theme-text-secondary);
    
    &:hover {
      color: var(--site-theme-brand);
    }
  }
  
  :deep(.ant-tabs-tab-active) {
    .ant-tabs-tab-btn {
      color: var(--site-theme-brand) !important;
    }
  }
  
  :deep(.ant-tabs-ink-bar) {
    background-color: var(--site-theme-brand);
  }
}
</style>
