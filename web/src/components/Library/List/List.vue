<template>
  <div class="library-list-container">
    <!-- 表格内容区域 -->
    <div class="table-content-area">
      <div class="list-table-scroll">
        <ListHead v-if="!storeLibraryList.paperListEmpty" :collating="collating"></ListHead>
        <ListCommonTable v-if="
            !storeLibraryList.paperListLoading && storeLibraryList.paperListTotal
          " :list="storeLibraryList.paperListAll" :collating="collating"
          :paper-head-visible-list="storeLibraryList.paperHeadVisibleList" @addButtonClick="handleClick"
          @refresh-classiy-author-venue-paper-list="
            storeLibraryList.refreshClassiyAuthorVenuePaperList()
          " />
      </div>
      <Loading v-if="storeLibraryList.paperListLoading" spin-size="default" class="list-table-loading" />
      <EmptyComponent v-else-if="!storeLibraryList.paperListTotal" class="empty" :desc="$t('home.library.noCollection')"
        :button="
          !storeLibraryList.searchInput && storeLibraryList.paperListEmpty
            ? $t('home.library.startAdd')
            : ''
        " @buttonClick="handleClick" />
    </div>
    
    <!-- 分页区域 -->
    <div v-if="storeLibraryList.paperListTotal" class="pagination-container">
      <ConfigProvider :locale="isWebEN ? enUS : zhCN">
        <Pagination :current="storeLibraryList.paperListPageNumber" :page-size="storeLibraryList.paperListPageSize"
          size="small" show-size-changer show-quick-jumper :total="storeLibraryList.paperListTotal"
          :page-size-options="['10', '20', '40', '100']" :style="{
            userSelect: 'none',
            pointerEvents: storeLibraryList.paperListLoading ? 'none' : '',
          }" @change="onChange" />
      </ConfigProvider>
    </div>
    
    <AuthorEdit v-if="storeLibraryList.paperListAuthorEdit" />
  </div>
</template>
<script lang="ts" setup>
import lodash from 'lodash'
import { Pagination, ConfigProvider } from 'ant-design-vue'
import { watch, onMounted } from 'vue'
import enUS from 'ant-design-vue/lib/locale-provider/en_US'
import zhCN from 'ant-design-vue/lib/locale-provider/zh_CN'
import ListCommonTable from './ListCommonTable.vue'
import ListHead from '../Head/index.vue'
import EmptyComponent from './Empty.vue'
import AuthorEdit from './AuthorEdit.vue'
import { useLibraryIndex } from '@/stores/library'
import { useLibraryList } from '@/stores/library/list'
import Loading from '../../Common/loading.vue'
import { useLanguage } from '@/hooks/useLanguage'
import { useUserStore } from '@/common/src/stores/user'

defineProps<{
  collating: boolean
}>()

const emit = defineEmits<{
  (event: 'buttonClick'): void
}>()

const { isEnUS } = useLanguage()
const isWebEN = isEnUS // 保持向后兼容的命名

const storeLibraryIndex = useLibraryIndex()

const storeLibraryList = useLibraryList()

const userStore = useUserStore()

;(window as any).storeLibraryList = storeLibraryList

const refresh = () => {
  storeLibraryList.refreshClassiyAuthorVenuePaperList(true)
}
onMounted(refresh)
watch(() => storeLibraryIndex.rawFolderId, refresh)

const handleClick = () => {
  emit('buttonClick')
}

const onChange = (number: number, size: number) => {
  console.warn({ number, size })
  storeLibraryList.paperListChecked = []
  storeLibraryList.getFilesByFolderId(number, size)
}

const showSizeChange = (number: number, size: number) => {
  // storeLibraryList.paperListChecked = []
  // storeLibraryList.getFilesByFolderId(number, size)
}

const search = lodash.debounce(() => {
  if (userStore.userInfo) {
    storeLibraryList.getFilesByFolderId(1)
  }
}, 400)

watch(() => storeLibraryList.searchInput, search)
</script>
<style lang="less" scoped>
@import '../Menu/Menu.less';

.library-list-container {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  margin-left: 0 !important;
  
  // 表格内容区域
  .table-content-area {
    flex: 1;
    min-height: 0; // 重要：允许flex子元素缩小
    position: relative;
    overflow: hidden;
    
    .list-table-scroll {
      height: 100%;
      overflow: auto;
      
      /* 自定义滚动条样式 - Webkit浏览器 */
      &::-webkit-scrollbar {
        width: 8px;
        height: 8px;
      }
      
      &::-webkit-scrollbar-track {
        background: transparent; /* 透明背景 */
        border-radius: 4px;
      }
      
      &::-webkit-scrollbar-thumb {
        background: rgba(0, 0, 0, 0.2); /* 半透明滚动条 */
        border-radius: 4px;
        transition: background 0.2s ease;
        
        &:hover {
          background: rgba(0, 0, 0, 0.4); /* 悬停时稍微深一点 */
        }
      }
      
      &::-webkit-scrollbar-corner {
        background: transparent; /* 滚动条交汇处透明 */
      }
      
      /* Firefox 滚动条样式 */
      scrollbar-width: thin;
      scrollbar-color: rgba(0, 0, 0, 0.2) transparent;
    }
    
    .list-table-loading {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      z-index: 100;
    }
    
    .empty {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      width: 100%;
    }
  }
  
  // 分页区域
  .pagination-container {
    flex-shrink: 0; // 防止分页区域被压缩
    width: 100%;
    text-align: right;
    padding: 16px 0;
    background-color: var(--site-theme-background);
    border-top: 1px solid var(--site-theme-border-color);
    z-index: 10;
  }

  :deep(.ant-pagination-options) {
    .ant-select-selector {
      color: var(--site-theme-text-primary) !important;
    }

    .ant-select-selection-item {
      color: var(--site-theme-text-primary) !important;
    }
  }

  /* 修复下拉菜单样式，使用主题变量 */
  :deep(.ant-select-dropdown) {
    background-color: var(--site-theme-background-primary);

    .ant-select-item {
      color: var(--site-theme-text-primary) !important;
      background-color: var(--site-theme-background-primary);
    }

    .ant-select-item-option-active {
      background-color: var(--site-theme-background-secondary);
      color: var(--site-theme-text-primary) !important;
    }

    .ant-select-item-option-selected {
      color: var(--site-theme-text-primary) !important;
      background-color: var(--site-theme-background-secondary);
      font-weight: bold;
    }
  }

  /* 修复跳转输入框样式 */
  :deep(.ant-pagination-options-quick-jumper) {
    color: var(--site-theme-text-primary) !important;

    input {
      color: var(--site-theme-text-primary) !important;
      background-color: var(--site-theme-background-primary) !important;
      border: 1px solid var(--site-theme-border-color);
      border-radius: 2px;
    }
  }

  /* 修复翻页按钮样式 */
  :deep(.ant-pagination-prev),
  :deep(.ant-pagination-next) {
    .ant-pagination-item-link {
      color: var(--site-theme-text-primary) !important;
      background-color: var(--site-theme-background-primary) !important;

      &:hover {
        border-color: var(--site-theme-primary-color);
        color: var(--site-theme-primary-color) !important;
      }
    }
  }

  /* 修复页码样式 */
  :deep(.ant-pagination-item) {
    background-color: var(--site-theme-background-primary) !important;
    border-color: var(--site-theme-border-color) !important;

    a {
      color: var(--site-theme-text-primary) !important;
    }

    &:hover {
      border-color: var(--site-theme-primary-color);

      a {
        color: var(--site-theme-primary-color) !important;
      }
    }

    &.ant-pagination-item-active {
      border-color: var(--site-theme-primary-color);
      background-color: var(--site-theme-primary-color) !important;

      a {
        color: #ffffff !important;
      }

      &:hover {
        a {
          color: #ffffff !important;
        }
      }
    }
  }
}
</style>




