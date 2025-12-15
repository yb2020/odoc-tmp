<template>
  <div class="time-group">
    <div class="time-label" @click="toggleExpand">
      <span class="expand-icon">
        <down-outlined v-if="isExpanded" />
        <right-outlined v-else />
      </span>
      {{ title }}
    </div>
    <div class="document-list" v-show="isExpanded">
      <div 
        v-for="doc in documents" 
        :key="doc.docId" 
        class="document-item"
        @click="openPage(doc)"
      >     
        <div class="document-icon">
          <file-pdf-outlined style="color: #ff4d4f; font-size: 24px;" />
        </div>
        <div class="document-info">
          <div class="document-title">{{ doc.docName }}</div>
          <div class="document-meta">
            <template v-if="doc.venue">{{ t('recentReading.showAttributes.venue') }}: {{ doc.venue }}</template>
            <template v-if="doc.publishDate">{{ doc.venue ? ' | ' : '' }}{{ t('recentReading.showAttributes.publishDate') }}: {{ doc.publishDate }}</template>
            <template v-if="doc.remark">{{ (doc.venue || doc.publishDate) ? ' | ' : '' }}{{ t('recentReading.showAttributes.remark') }}: {{ doc.remark }}</template>
          </div>
        </div>
        <div class="document-date">{{ formatTime(doc.lastReadTime) }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { FilePdfOutlined, DownOutlined, RightOutlined } from '@ant-design/icons-vue';
import { useI18n } from 'vue-i18n';
import { goPathPage, goPdfPage } from '@/common/src/utils/url'
import { ref, computed } from 'vue';

const { t } = useI18n();

// 接收属性
const props = defineProps({
  title: {
    type: String,
    required: true
  },
  documents: {
    type: Array,
    required: true
  },
  defaultExpanded: {
    type: Boolean,
    default: true
  }
});

// 根据是否有数据决定初始展开状态
const hasDocuments = computed(() => props.documents && props.documents.length > 0);

// 展开/收缩状态 - 如果没有文档则默认收起
const isExpanded = ref(hasDocuments.value ? props.defaultExpanded : false);

// 切换展开/收缩状态
const toggleExpand = () => {
  isExpanded.value = !isExpanded.value;
};

const openPage = (doc) => {
  if (doc.pdfId && doc.pdfId !== '0') {
    goPdfPage({ pdfId: doc.pdfId });
  } else if (doc.paperId && doc.paperId !== '0') {
    goPathPage(`/note/${doc.paperId}`);
  }
};

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '';
  
  const date = new Date(parseInt(timestamp) * 1000);
  const month = (date.getMonth() + 1).toString().padStart(2, '0');
  const day = date.getDate().toString().padStart(2, '0');
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  
  return `${month}/${day} ${hours}:${minutes}`;
};
</script>

<style lang="less" scoped>
.time-group {
  margin-bottom: 30px;

  .time-label {
    font-size: 14px;
    color: var(--site-theme-text-secondary-color);
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    cursor: pointer;
    user-select: none;

    .expand-icon {
      margin-right: 8px;
      display: flex;
      align-items: center;
      font-size: 12px;
      transition: transform 0.3s;
    }

    &::after {
      content: '';
      flex: 1;
      height: 1px;
      background-color: var(--site-theme-border-color);
      margin-left: 10px;
    }

    &:hover {
      color: var(--site-theme-primary-color);
    }
  }
}

.document-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  transition: max-height 0.3s ease-in-out, opacity 0.3s ease-in-out;
}

.document-item {
  display: flex;
  align-items: center;
  padding: 12px 15px;
  background-color: var(--site-theme-background-secondary);
  border-radius: 4px;
  transition: all 0.3s;
  cursor: pointer;

  &:hover {
    background-color: var(--site-theme-background-hover);
  }

  .document-icon {
    font-size: 24px;
    margin-right: 15px;
  }

  .document-info {
    flex: 1;

    .document-title {
      font-size: 14px;
      margin-bottom: 5px;
      color: var(--site-theme-text-color);
    }

    .document-meta {
      font-size: 12px;
      color: var(--site-theme-text-secondary-color);
    }
  }

  .document-date {
    font-size: 12px;
    color: var(--site-theme-text-secondary-color);
  }
}

.empty-group-message {
  font-size: 14px;
  color: var(--site-theme-text-secondary-color);
  padding: 12px 15px;
  background-color: var(--site-theme-background-secondary);
  border-radius: 4px;
}
</style>
