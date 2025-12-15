<template>
  <div class="recent-reading">
    <div class="header">
      <div class="header-left">
        <a-button type="primary" class="upload-btn" @click="handleUploadClick">
          <plus-outlined />
          {{ t('recentReading.importDocument') }}
        </a-button>
      </div>
    </div>

    <div class="document-section">
      <div class="section-header">
        <h2>{{ t('recentReading.recentlyRead') }}</h2>
        <!-- <div class="section-actions">
          <a-button type="text" class="action-btn">
            <appstore-outlined />
          </a-button>
          <a-button type="text" class="action-btn">
            <sort-ascending-outlined />
          </a-button>
        </div> -->
      </div>

      <div class="time-groups">
        <!-- 加载中状态 -->
        <div v-if="loading" class="loading-container">
          <a-spin />
        </div>
        <!-- 使用DocumentGroup组件显示各时间段的文档 -->
        <template v-else>
          <document-group 
            :title="t('recentReading.today')"
            :documents="todayDocs"
          />
          
          <document-group 
            :title="t('recentReading.yesterday')"
            :documents="yesterdayDocs"
          />
          
          <document-group 
            :title="t('recentReading.withinSevenDays')"
            :documents="weekDocs"
          />
          
          <document-group 
            :title="t('recentReading.earlierThisMonth')"
            :documents="earlierDocs"
          />
        </template>
      </div>
    </div>

    <!-- 添加文件上传组件 -->
    <FileUploader
      :visible="storeLibraryIndex.uploaderVisible"
      :folder-id="storeLibraryIndex.rawFolderId || '0'"
      :selected-key="storeLibraryIndex.selectedKey"
      :bread-crumb-list="storeLibraryIndex.breadCrumbList || []"
      @close="storeLibraryIndex.uploaderVisible = false"
      @addSuccess="handleAddSuccess"
      @refreshList="handleAddSuccess"
    />
  </div>
</template>

<script lang="ts" setup>
import {
  PlusOutlined,
  AppstoreOutlined,
  SortAscendingOutlined,
} from '@ant-design/icons-vue';
import { useI18n } from 'vue-i18n';
import { computed, onMounted, ref, defineAsyncComponent } from 'vue';
import { ThemeType } from '@/theme';
import { getLatestReadDocListNew } from '@/api/document';
import { message } from 'ant-design-vue';
import { useLibraryIndex } from '@/stores/library';

// 使用defineAsyncComponent正确处理懒加载组件
const DocumentGroup = defineAsyncComponent(() => import('@/components/RecentReading/DocumentGroup.vue'));
const FileUploader = defineAsyncComponent(() => import('@/components/Library/File/Uploader.vue'));

// 使用i18n
const { t } = useI18n();

// 获取library store
const storeLibraryIndex = useLibraryIndex();

// 处理上传按钮点击事件
const handleUploadClick = () => {
  // 设置上传组件可见
  storeLibraryIndex.uploaderVisible = true;
};

// 处理文件上传成功
const handleAddSuccess = async () => {
  // 刷新最近阅读列表
  try {
    loading.value = true;
    const result = await getLatestReadDocListNew({});
    
    // 根据时间分类文档
    if (result && result.docInfos) {
      categorizeDocsByTime(result.docInfos, result.currentTime);
    }
  } catch (error) {
    message.error(t('recentReading.fetchError'));
  } finally {
    loading.value = false;
  }
};

// 获取当前主题
const currentTheme = computed(() => {
  if (typeof localStorage !== 'undefined') {
    return localStorage.getItem('theme') || ThemeType.LIGHT;
  }
  return ThemeType.LIGHT;
});

// 判断是否为暗黑模式
const isDarkMode = computed(() => currentTheme.value === ThemeType.DARK);

// 定义文档分组
const todayDocs = ref([]);
const yesterdayDocs = ref([]);
const weekDocs = ref([]);
const earlierDocs = ref([]);
const loading = ref(true);

// 分组文档函数
const categorizeDocsByTime = (docList, serverTime) => {
  // 清空之前的分类
  todayDocs.value = [];
  yesterdayDocs.value = [];
  weekDocs.value = [];
  earlierDocs.value = [];
  
  // 计算时间界限（秒级时间戳）
  const currentTime = parseInt(serverTime);
  const oneDaySeconds = 24 * 60 * 60;
  
  // 今天开始时间（当天0点）
  const todayStart = currentTime - (currentTime % oneDaySeconds);
  // 昨天开始时间
  const yesterdayStart = todayStart - oneDaySeconds;
  // 一周前开始时间
  const weekStart = todayStart - (7 * oneDaySeconds);
  // 一个月前开始时间（近似为30天）
  const monthStart = todayStart - (30 * oneDaySeconds);
  
  // 分类文档
  docList.forEach(doc => {
    const readTime = parseInt(doc.lastReadTime);
    
    if (readTime >= todayStart) {
      // 今天阅读的文档
      todayDocs.value.push(doc);
    } else if (readTime >= yesterdayStart) {
      // 昨天阅读的文档
      yesterdayDocs.value.push(doc);
    } else if (readTime >= weekStart) {
      // 七天内阅读的文档
      weekDocs.value.push(doc);
    } else if (readTime >= monthStart) {
      // 更早阅读的文档
      earlierDocs.value.push(doc);
    }
  });
};

// 在组件加载时调用API获取数据
onMounted(async () => {
  try {
    loading.value = true;
    const result = await getLatestReadDocListNew({});
    console.log('最近阅读文档数据:', result);
    
    // 根据时间分类文档
    if (result && result.docInfos) {
      categorizeDocsByTime(result.docInfos, result.currentTime);
      console.log('今天的文档:', todayDocs.value);
      console.log('昨天的文档:', yesterdayDocs.value);
      console.log('七天内的文档:', weekDocs.value);
      console.log('本月更早的文档:', earlierDocs.value);
    }
  } catch (error) {
    console.error('获取最近阅读文档失败:', error);
    message.error(t('recentReading.fetchError'));
  } finally {
    loading.value = false;
  }
});
</script>

<style lang="less" scoped>
.recent-reading {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  &-left {
    display: flex;
    align-items: center;
    gap: 15px;
  }
}

.upload-btn {
  background-color: var(--site-theme-primary-color);
  border-color: var(--site-theme-primary-color);

  &:hover {
    background-color: var(--site-theme-primary-color-hover);
    border-color: var(--site-theme-primary-color-hover);
  }
}

.document-section {
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;

    h2 {
      font-size: 18px;
      margin: 0;
      color: var(--site-theme-text-color);
    }

    .section-actions {
      display: flex;
      gap: 10px;
    }
  }
}

.time-groups {
  /* 增加内下边距，避免“本月更早”等底部分组被遮挡或显示不全 */
  padding-bottom: 64px;
  .time-group {
    margin-bottom: 30px;

    .time-label {
      font-size: 14px;
      color: var(--site-theme-text-secondary-color);
      margin-bottom: 10px;
      display: flex;
      align-items: center;

      &::after {
        content: '';
        flex: 1;
        height: 1px;
        background-color: var(--site-theme-border-color);
        margin-left: 10px;
      }
    }
  }
}

.document-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
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

:deep(.ant-input-search) {
  .ant-input {
    background-color: var(--site-theme-background-secondary);
    border-color: var(--site-theme-border-color);
    color: var(--site-theme-text-color);

    &:focus {
      border-color: var(--site-theme-primary-color);
    }
  }

  .ant-input-search-button {
    background-color: var(--site-theme-primary-color);
    border-color: var(--site-theme-primary-color);

    &:hover {
      background-color: var(--site-theme-primary-color-hover);
      border-color: var(--site-theme-primary-color-hover);
    }
  }
}

:deep(.ant-checkbox-wrapper) {
  color: var(--site-theme-text-color);

  .ant-checkbox-checked .ant-checkbox-inner {
    background-color: var(--site-theme-primary-color);
    border-color: var(--site-theme-primary-color);
  }
}

:deep(.ant-btn-text) {
  color: var(--site-theme-text-color);

  &:hover {
    background-color: var(--site-theme-background-hover);
    color: var(--site-theme-primary-color);
  }
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;

  p {
    margin-top: 10px;
    color: var(--site-theme-text-secondary-color);
  }
}
</style>