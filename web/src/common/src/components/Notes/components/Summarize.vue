<template>
  <div class="notes-list-container">
    <NoteBreadcrumb
      :type="NoteSubTypes.Summary"
      :noteState="noteState"
      class="px-5 py-1 w-full"
      :style="{ backgroundColor: 'var(--site-theme-background-tertiary)' }"
    />

    <!-- 表格内容区域 -->
    <div class="table-content-area">
      <div class="list-table-scroll">
        <div v-if="!noteFolderMap[noteFolderSelected]?.isDoc">
          <!-- 表格头部 -->
          <div class="notes-table-header">
            <div class="header-cell title-cell">
              {{ $t('common.notes.columns.tt') }}
            </div>
            <div class="header-cell date-cell">
              {{ $t('common.notes.columns.mtime') }}
            </div>
          </div>

          <!-- 表格内容 -->
          <div class="notes-table-body">
            <div
              v-for="record in summaryList"
              :key="record.noteId"
              class="notes-table-row"
            >
              <div class="body-cell title-cell">
                <a
                  :href="`/summary.html?noteId=${record.noteId}`"
                  target="_blank"
                  class="note-link"
                  @click="onNoteClick($event, record)"
                  >{{ record.docName }}</a
                >
              </div>
              <div class="body-cell date-cell">
                {{
                  dayjs(new Date(Number(record.modifyDate))).format(
                    'YYYY-MM-DD HH:mm:ss'
                  )
                }}
              </div>
            </div>
          </div>
        </div>

        <div v-else class="editor-container">
          <SummaryEditor :note-id="noteFolderSelected" @loaded="onLoad" />
        </div>
      </div>
      <!-- 分页区域 -->
      <div
        v-if="!noteFolderMap[noteFolderSelected]?.isDoc && summaryTotal"
        class="pagination-container"
      >
        <ConfigProvider :locale="zhCN">
          <Pagination
            :current="summaryPageNumber"
            :page-size="summaryPageSize"
            size="small"
            show-size-changer
            show-quick-jumper
            :total="summaryTotal"
            :page-size-options="['10', '20', '40', '100']"
            :style="{
              userSelect: 'none',
            }"
            @change="onNoteListChange"
          />
        </ConfigProvider>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import dayjs from 'dayjs'
import { Pagination, ConfigProvider } from 'ant-design-vue'
import zhCN from 'ant-design-vue/lib/locale-provider/zh_CN'
import { ref } from 'vue'
import { useNote } from '../useNote'
import SummaryEditor from './SummaryEditor.vue'
import NoteBreadcrumb from './NoteBreadcrumb.vue'
import {
  GetSummaryResponse,
  NoteManageDocInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule'
import { NoteSubTypes } from '../types'

const props = defineProps<{
  noteState: ReturnType<typeof useNote>
}>()

const emit = defineEmits<{
  (e: 'clickNote', evt: MouseEvent, x: NoteManageDocInfo): void
  (e: 'loadedNote', x: GetSummaryResponse): void
}>()

const {
  noteFolderMap,
  noteFolderSelected,
  summaryList,
  summaryPageSize,
  summaryTotal,
  summaryPageNumber,
  fetchSummary,
} = props.noteState

const onNoteClick = (e: MouseEvent, record: NoteManageDocInfo) => {
  emit('clickNote', e, record)
}

const onNoteListChange = (current: number, pageSize: number) => {
  summaryPageSize.value = pageSize
  fetchSummary(current)
}

const onLoad = (data: GetSummaryResponse) => {
  emit('loadedNote', data)
}
</script>

<style lang="less" scoped>
// 移除错误的导入路径，直接使用内联样式

.notes-list-container {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  margin-left: 0 !important;

  // 表格内容区域
  .table-content-area {
    position: relative;
    overflow: hidden;

    .list-table-scroll {
      height: calc(100vh - 18vh); /* 为分页组件留出60px的空间 */
      overflow-y: auto;
      overflow-x: hidden;
    }

    .editor-container {
      height: 100%;
      padding: 24px;

      :deep(.milkdown-menu) {
        min-height: 50px;
      }
      :deep(.milkdown-inner) {
        min-height: calc(100vh - 200px);
      }
    }
  }

  // 分页区域
  .pagination-container {
    flex-shrink: 0; // 防止分页区域被压缩
    width: 100%;
    text-align: right;
    padding: 16px 24px;
    background-color: var(--site-theme-background);
    border-top: 1px solid var(--site-theme-border-color);
    z-index: 10;
  }
}

// 搜索栏样式 - 与Library组件保持一致
.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  background-color: var(--site-theme-background);
  border-bottom: 1px solid var(--site-theme-border-color);

  .banner {
    display: flex;
    align-items: center;
    margin-right: 20px;

    .search {
      width: 240px;
      margin-right: 16px;
    }

    .select {
      width: 120px;
      margin-right: 16px;
    }

    .divider {
      width: 1px;
      height: 24px;
      background-color: var(--site-theme-border-color);
      margin: 0 16px;
    }

    .btn {
      margin-left: 8px;
    }
  }
}

// Notes表格样式 - 与Library组件保持一致，移除列间分隔线
.notes-table-header {
  display: flex;
  background-color: var(--site-theme-background-secondary);
  border-bottom: 1px solid var(--site-theme-border-color);

  .header-cell {
    padding: 12px 16px;
    font-weight: 500;
    color: var(--site-theme-text-secondary);
    // 移除列间分隔线
    // border-right: 1px solid var(--site-theme-border-color);

    &:last-child {
      border-right: none;
    }

    &.title-cell {
      flex: 1;
    }

    &.date-cell {
      width: 200px;
      flex-shrink: 0;
    }
  }
}

.notes-table-body {
  .notes-table-row {
    display: flex;
    border-bottom: 1px solid var(--site-theme-border-color);
    transition: background-color 0.2s ease;

    &:hover {
      background-color: var(--site-theme-background-hover);
    }

    &:last-child {
      border-bottom: none;
    }

    .body-cell {
      padding: 12px 16px;
      color: var(--site-theme-text-primary);
      // 移除列间分隔线
      // border-right: 1px solid var(--site-theme-border-color);

      &:last-child {
        border-right: none;
      }

      &.title-cell {
        flex: 1;

        .note-link {
          color: var(--site-theme-primary-color);
          text-decoration: none;

          &:hover {
            text-decoration: underline;
            color: var(--site-theme-primary-color-hover);
          }
        }
      }

      &.date-cell {
        width: 200px;
        flex-shrink: 0;
        color: var(--site-theme-text-secondary);
      }
    }
  }
}

.notes-list-container {
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
      border: 1px solid var(--site-theme-border-color);
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
