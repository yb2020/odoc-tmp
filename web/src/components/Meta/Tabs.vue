<template>
  <div class="update-container">
    <a-tabs v-model="activeKey">
      <a-tab-pane :key="TabKey.Manual" :tab="$t('meta.manualTab')">
        <UserEdit
          :paper-id="paperId"
          :pdf-id="pdfId"
          @submit="$emit('update:success')"
          @cancel="$emit('cancel')"
        />
      </a-tab-pane>
      <a-tab-pane :key="TabKey.English" :tab="$t('meta.englishTab')">
        <SearchDOI
          :paper-id="paperId"
          :pdf-id="pdfId"
          :page-type="pageType"
          :placeholder="$t('meta.searchPlaceholder2')"
          from="english"
          @update:success="$emit('update:success')"
        />
      </a-tab-pane>
      <!-- <a-tab-pane :key="TabKey.Chinese" :tab="$t('meta.chineseTab')">
        <SearchDOI
          :paper-id="paperId"
          :pdf-id="pdfId"
          :page-type="pageType"
          :placeholder="$t('meta.searchPlaceholder1')"
          from="chinese"
          @update:success="$emit('update:success')"
        />
      </a-tab-pane> -->
    </a-tabs>
  </div>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import SearchDOI from './SearchDOI/index.vue'
import UserEdit from './UserInput/UserEdit.vue'

enum TabKey {
  Manual = 'manual',
  English = 'english',
  // Chinese = 'chinese',
}

defineProps<{
  paperId: string
  pdfId: string
  pageType: string
}>()

defineEmits<{
  (event: 'update:success'): void
  (event: 'cancel'): void
}>()

const activeKey = ref(TabKey.English)
</script>
<style lang="less" scoped>
.update-container {
  height: 100%;
  background-color: transparent !important;
}
:deep(.ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn) {
    color: var(--site-theme-text-primary);
    text-shadow: 0 0 0.25px currentcolor;
}
:deep(.ant-tabs-tab.ant-tabs-tab-active:hover .ant-tabs-tab-btn:hover) {
    color: var(--site-theme-text-primary);
    text-shadow: 0 0 0.25px currentcolor;
}
:deep(.ant-tabs-tab) {
    color: var(--site-theme-text-secondary);
}
:deep(.ant-tabs-tab:hover) {
    color: var(--site-theme-text-secondary);
}
</style>
