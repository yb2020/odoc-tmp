<template>
  <PerfectScrollbar
    :options="{
      suppressScrollX: true,
    }"
    class="catalog-tree"
  >
    <a-tree
      v-model:selectedKeys="selectedKeys"
      :default-expanded-keys="['0']"
      :tree-data="treeData"
      @select="handleSelect"
    >
      <template #title="{ title, key, pageNum, dest }">
        <a-tooltip
          :title="title"
          placement="right"
        >
          <span
            class="title"
            :class="{
              disabled: (!pageNum || parseInt(pageNum) < 1) && !dest,
            }"
          >
            <span class="text">{{ title }}</span>
          </span>
        </a-tooltip>
      </template>
    </a-tree>
  </PerfectScrollbar>
</template>
<script lang="ts" setup>
import { TreeProps } from 'ant-design-vue';
import { ref } from 'vue';
import { ViewerEvent, ViewerController } from '@idea/pdf-annotate-viewer';
import { PdfCatalogueInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import { DataNode } from 'ant-design-vue/lib/tree';
import { goBbox } from '~/src/util/bbox';
import { goToPDFPage } from '~/src/dom/pdf';

const props = defineProps<{
  catalogs: PdfCatalogueInfo;
  pdfViewInstance: ViewerController;
}>();

let isClickCalalog = false;

const handleSelect = (keys: string[], info: { node: { key: string } }) => {
  const key = keys[0] || (info.node.key ? info.node.key : '');
  if (!key) {
    return;
  }
  selectedKeys.value = [key];
  const calalog = treeMap[key] as any;
  if (calalog?.bbox) {
    isClickCalalog = true;
    goBbox(calalog.pageNum, calalog.bbox, props.pdfViewInstance);
  } else if (calalog?.pageNum) {
    isClickCalalog = true;
    goToPDFPage(calalog.pageNum, undefined, props.pdfViewInstance);
  } else if (calalog?.dest) {
    // @ts-ignore
    const { scrollController } = props.pdfViewInstance;

    isClickCalalog = true;
    scrollController.goToDestination(calalog.dest);
  } else {
    throw new Error('invalid bbox of calalog. key ' + key);
  }
};

const treeData: TreeProps['treeData'] = [];

const treeMap: Record<string, PdfCatalogueInfo> = {};

const buildTreeData = (
  calalogs: PdfCatalogueInfo[],
  treeData: DataNode[],
  treeMap: Record<string, PdfCatalogueInfo>,
  prefix: string
) => {
  calalogs.forEach((calalog: any, i) => {
    const key = prefix ? `${prefix}-${i}` : `${i}`;
    const data: DataNode = {
      title: calalog.title,
      key,
      pageNum: calalog.pageNum,
      dest: calalog.dest,
    };
    treeMap[key] = calalog;
    treeData.push(data);
    if (calalog.child?.length) {
      data.children = [];
      buildTreeData(calalog.child, data.children, treeMap, key);
    }
  });
};

buildTreeData(props.catalogs.child, treeData, treeMap, '');

props.pdfViewInstance?.addEventListener(
  ViewerEvent.PAGE_CHANGING,
  (pageInfo) => {
    if (isClickCalalog) {
      isClickCalalog = false;

      return;
    }
    const { pageNumber } = pageInfo;
    const page = treeData.find((item) => item.pageNum === pageNumber);

    if (page) {
      const { key } = page;

      if (key) {
        selectedKeys.value = [key as string];
      }
    }
  }
);

const selectedKeys = ref<string[]>();
</script>
<style lang="less" scoped>
.title {
  display: flex !important;
  flex: 1;
  padding: 0 4px;
  color: var(--site-theme-pdf-panel-text);
  &:hover {
    background-color: var(--site-theme-bg-hover);
    cursor: pointer;
  }
  &.disabled {
    background-color: var(--site-theme-bg-dark);
    &:hover {
      cursor: default;
    }
  }
  .text {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    width: 0;
  }
}
.catalog-tree {
  height: 100%;
  padding-right: 10px;
  background-color: var(--site-theme-pdf-panel);

  :deep(.ant-tree-treenode-motion) {
    width: 100%;
  }
  :deep(.ant-tree-list-holder-inner) {
    overflow: hidden;
  }
  :deep(.ant-tree-treenode) {
    width: 100% !important;
    
    &:hover {
      background-color: var(--site-theme-bg-hover);
    }
  }
  :deep(.ant-tree-node-content-wrapper) {
    flex: auto;
    display: flex;
    max-width: 100%;
    padding: 0 !important;
    background-color: transparent !important;
    .ant-tree-title {
      flex: 1;
      overflow: hidden;
    }
  }

  :deep(.ant-tree-treenode-selected .ant-tree-node-selected) {
    background-color: rgba(var(--site-theme-primary-color), 0.85) !important;
  }
  
  // 修复三角形展开/折叠图标的颜色
  :deep(.ant-tree-switcher) {
    color: var(--site-theme-pdf-panel-text);
    
    .ant-tree-switcher-icon {
      color: var(--site-theme-pdf-panel-text);
    }
  }
  
  // 修复分隔栏颜色
  :deep(.ant-tree-indent-unit) {
    &::before {
      border-color: var(--site-theme-divider-light);
    }
  }
  
  // 修复关闭按钮颜色
  :deep(.ant-tree-close-icon) {
    color: var(--site-theme-pdf-panel-text);
  }
}
</style>
