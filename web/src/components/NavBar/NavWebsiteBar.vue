<template>
  <div class="nav-website-bar-container">
    <!-- <div class="sidebar-category">学术网站</div> -->
    <draggable
      v-model="websites"
      item-key="id"
      class="website-list"
      handle=".website-info"
      @end="handleDragEnd"
      animation="200"
      ghost-class="ghost"
    >
      <template #item="{ element: item }">
        <div 
          class="website-item"
          :class="{ 'active': item.id === activeKey }"
          @click="handleItemClick(item)"
        >
          <img :src="item.iconUrl || defaultIcon" class="website-icon" @error="onIconError" />
          <div class="website-info">
            <span>{{ item.name }}</span>
          </div>
          <div class="website-actions">
            <edit-outlined @click.stop="handleEdit(item)" class="icon-edit" />
            <delete-outlined @click.stop="handleDelete(item)" class="icon-delete" />
          </div>
        </div>
      </template>
    </draggable>
    <div class="add-website-button" @click="handleAdd">
      <plus-outlined />
      <span>{{ t('navbar.navwebsite.addAcademicWebsite') }}</span>
    </div>
  </div>

  <WebsiteFormModal 
    v-model:visible="isModalVisible" 
    :website-to-edit="editingWebsite" 
    @success="fetchWebsites" 
  />
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import draggable from 'vuedraggable';
import { PlusOutlined, QuestionCircleOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons-vue';
import { Modal } from 'ant-design-vue';

import { getWebsiteList, reorderWebsite, deleteWebsite } from '@/api/nav';
import { Website, WebsiteSource } from 'go-sea-proto/gen/ts/nav/Website';
import WebsiteFormModal from './WebsiteFormModal.vue';
import defaultIcon from '@/assets/images/icon-lang.svg';

const websites = ref<Website[]>([]);
const activeKey = ref<bigint | null>(null);
const isModalVisible = ref(false);
const editingWebsite = ref<Website | null>(null);
const { t } = useI18n();

const fetchWebsites = async () => {
  try {
    const res = await getWebsiteList({});
    if (res.websites) {
      websites.value = res.websites;
      if (res.websites.length > 0) {
        const activeExists = res.websites.some(w => w.id === activeKey.value);
        if (!activeExists) {
          activeKey.value = res.websites[0].id;
        }
      }
    }
  } catch (error) {
    console.error('Failed to fetch website list:', error);
  }
};

onMounted(() => {
  fetchWebsites();
});

const handleItemClick = (item: Website) => {
  activeKey.value = item.id;
  if (item.url) {
    window.open(item.url, '_blank');
  }
};

const onIconError = (e: Event) => {
  const target = e.target as HTMLImageElement;
  target.src = defaultIcon;
};

const handleDelete = (website: Website) => {
  Modal.confirm({
    title: t('navbar.navwebsite.deleteConfirmTitle'),
    content: t('navbar.navwebsite.deleteConfirmContent', { name: website.name }),
    okText: t('navbar.navwebsite.okText'),
    cancelText: t('navbar.navwebsite.cancelText'),
    onOk: async () => {
      try {
        await deleteWebsite({ id: website.id });
        if (activeKey.value === website.id) {
          activeKey.value = null;
        }
        fetchWebsites();
      } catch (error) {
        console.error('Failed to delete website:', error);
      }
    },
  });
};

const handleAdd = () => {
  editingWebsite.value = null;
  isModalVisible.value = true;
};

const handleEdit = (website: Website) => {
  editingWebsite.value = website;
  isModalVisible.value = true;
};

const handleAddSuccess = () => {
  // 添加成功后，重新获取列表以保证数据最新
  fetchWebsites();
};

const handleDragEnd = async (event: any) => {
  const { newIndex } = event;
  const movedItem = websites.value[newIndex];

  if (!movedItem) return;

  const params: any = { id: movedItem.id };

  // 仅当存在前一个元素时，才添加 prevId
  if (newIndex > 0) {
    params.prevId = websites.value[newIndex - 1].id;
  }

  // 仅当存在后一个元素时，才添加 nextId
  if (newIndex < websites.value.length - 1) {
    params.nextId = websites.value[newIndex + 1].id;
  }

  try {
    const res = await reorderWebsite(params);
    if (res.rebalanced) {
      // 如果后端进行了重平衡，则完全重新加载列表
      await fetchWebsites();
    } else if (res.updates && res.updates.length > 0) {
      // 否则，根据返回的 updates 列表进行局部更新
      res.updates.forEach(update => {
        const index = websites.value.findIndex(w => w.id === update.id);
        if (index !== -1) {
          // 直接更新本地数组中对应项的 sortOrder
          websites.value[index].sortOrder = update.sortOrder;
        }
      });
    }
  } catch (error) {
    console.error('Failed to reorder website:', error);
    // 如果排序失败，重新加载列表以恢复到一致的状态
    await fetchWebsites();
  }
};
</script>

<style lang="less" scoped>
.nav-website-bar-container {
  color: var(--site-theme-text-color);
}

.sidebar-category {
  padding: 10px 20px;
  font-size: 12px;
  color: var(--site-theme-text-secondary-color);
  text-transform: uppercase;
}

.website-list {
  display: flex;
  flex-direction: column;
}

.ghost {
  opacity: 0.5;
  background: var(--site-theme-primary-color-fade);
}

.website-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
  white-space: nowrap;
  height: 44px;

  .website-icon {
    width: 16px;
    height: 16px;
    margin-right: 10px;
    flex-shrink: 0;
  }

  .website-info {
    display: flex;
    align-items: center;
    flex-grow: 1;
    height: 100%;
    overflow: hidden; // 防止子元素溢出撑开布局

    span {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  .website-actions {
    display: none;
    gap: 8px;
  }

  &:hover {
    background-color: var(--site-theme-background-hover);
    .website-actions {
      display: flex;
    }
  }

  .icon-edit,
  .icon-delete {
    font-size: 14px;
    color: var(--site-theme-text-secondary-color);
    cursor: pointer;

    &:hover {
      color: var(--site-theme-primary-color);
    }
  }

  // &.active {
  //   background-color: var(--site-theme-primary-color-fade);
  //   color: var(--site-theme-primary-color);
  // }
}

.add-website-button {
  display: flex;
  align-items: center;
  padding: 10px 20px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
  white-space: nowrap;
  height: 44px;
  color: var(--site-theme-text-secondary-color);

  .anticon {
    margin-right: 10px;
    font-size: 16px;
  }

  .help-icon {
    margin-left: 8px;
  }

  &:hover {
    background-color: var(--site-theme-background-hover);
  }
}
</style>
