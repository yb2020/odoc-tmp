<template>
  <Layout
    :title="$t('glossary.title')"
    noDing
    :tippyHandler="onTippyHandler"
    group="glossary"
    :style="{ width: '600px' }"
  >
    <div
      ref="containerRef"
      class="py-3"
    >
      <div class="flex justify-between items-center px-3">
        <div ref="createGlossaryButtonRef">
          <a-button
            v-if="checkedItems.length"
            type="default"
            size="small"
            class="mr-8"
            @click="onDelete"
          >
            {{ $t('viewer.delete') }}
          </a-button>
          <a-button
            v-else
            type="default"
            size="small"
            class="mr-8"
            @click="onCreate"
          >
            {{ $t('glossary.addTitle') }}
          </a-button>
        </div>
        <a-input-search
          v-model:value="searchValue"
          class="flex-1 glossary-search"
          :placeholder="$t('glossary.search')"
          size="small"
          @search="onSearch(false)"
          @change="onSearch(true)"
        />
      </div>
      <div
        class="flex"
        :style="{ 'min-height': '240px' }"
      >
        <div
          v-if="loading && !list?.length"
          class="flex-1 flex items-center justify-center"
        >
          <loading-outlined />
        </div>
        <div
          v-else-if="error"
          class="flex-1 flex flex-col items-center justify-center"
        >
          <div class="text-[#4e5969]">
            {{ error.message }}
          </div>
          <a-button
            type="link"
            class="mt-2"
            @click="run"
          >
            {{
              $t('glossary.table.reload')
            }}
          </a-button>
        </div>
        <div
          v-else-if="total === 0"
          class="flex flex-1 items-center justify-center text-[#4e5969]"
        >
          {{ $t('message.noDataTip') }}
        </div>
        <div
          v-else-if="list?.length"
          class="flex-1 flex flex-col justify-between"
        >
          <div class="mt-4 glossary-table">
            <a-row class="w-full bg-[#f5f7fa] p-2">
              <a-col
                class="px-1"
                :span="2"
              >
                <a-checkbox
                  class="rp-light-theme"
                  size="small"
                  :checked="isCheckedAll"
                  @change="onCheckedAll"
                />
              </a-col>
              <a-col
                class="px-1"
                :span="6"
              >
                {{
                  $t('glossary.table.original')
                }}
              </a-col>
              <a-col
                class="px-1"
                :span="6"
              >
                {{
                  $t('translate.translation')
                }}
              </a-col>
              <a-col
                class="px-1"
                :span="5"
              >
                {{
                  $t('glossary.table.isCaseSensitive')
                }}
              </a-col>
              <a-col
                class="px-1"
                :span="5"
              >
                {{
                  $t('glossary.table.operation')
                }}
              </a-col>
            </a-row>
            <a-row
              v-for="(item, idx) in list"
              :key="item.id"
              class="w-full p-2 text-sm flex items-center border-l-0 border-t-0 border-r-0 border-b border-solid border-[var(--rp-theme-bg-f5f7fa,#f5f7fa)]"
            >
              <a-col
                class="px-1"
                :span="2"
              >
                <a-checkbox
                  class="rp-light-theme"
                  :checked="checkedItems.includes(item.id)"
                  @change="onItemCheckedChange(item.id)"
                />
              </a-col>
              <a-col
                class="px-1 overflow-hidden text-ellipsis"
                :span="6"
              >
                {{
                  item.originalText
                }}
              </a-col>
              <a-col
                v-if="!item.ignored"
                class="px-1 overflow-hidden text-ellipsis"
                :span="6"
              >
                {{ item.translationText }}
              </a-col>
              <a-col
                v-else
                class="px-1 text-[rgba(0,0,0,0.5)]"
                :span="6"
              >
                /
              </a-col>
              <a-col
                class="px-1"
                :span="5"
              >
                {{
                  item.matchCase
                    ? $t('glossary.table.yes')
                    : $t('glossary.table.no')
                }}
              </a-col>
              <a-col
                class="px-1 space-x-2"
                :span="5"
              >
                <a @click="onEdit(item)">{{ $t('glossary.table.edit') }}</a>
                <a-popconfirm
                  :title="$t('glossary.table.deleteTip')"
                  :ok-text="$t('viewer.delete')"
                  :cancel-text="$t('viewer.cancel')"
                  :getPopupContainer="getPopoverContainer"
                  placement="topRight"
                  @confirm="onConfirmDelete([item.id])"
                >
                  <a href="#">{{ $t('viewer.delete') }}</a>
                </a-popconfirm>
              </a-col>
            </a-row>
          </div>
          <div class="mt-4 text-right px-4">
            <a-pagination
              v-model:current="currentPage"
              class="glossary-pagination"
              size="small"
              :total="total"
              :showSizeChanger="false"
              @change="onPageChange"
            />
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>
<script setup lang="ts">
import Layout from '@/components/Tippy/Layout/index.vue';
import { GlossaryItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/GlossaryManage';
import { useRequest } from 'ahooks-vue';
import { computed, ref } from 'vue';
import {
  deleteGlossaryItems,
  fethGlossrayListByPage,
} from '~/src/api/glossary';
import { LoadingOutlined } from '@ant-design/icons-vue';
import { Modal, message } from 'ant-design-vue';
import { createGlossaryCreateEditTippyVue } from '~/src/dom/tippy';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  close: () => void;
}>();

const { t } = useI18n();

const onTippyHandler = (event: 'ding' | 'close' | 'unding' | 'lock') => {
  if (event === 'close') {
    props.close();
  }
};

const searchValue = ref('');
const onSearch = (fromChange?: boolean) => {
  const value = searchValue.value.trim();
  if (fromChange) {
    !value && run();
    return;
  }
  if (value) {
    currentPage.value = 1;
    run();
  }
};

const containerRef = ref<HTMLElement | null>(null);
const getPopoverContainer = () => {
  return containerRef.value || document.body;
};

const currentPage = ref(1);
const total = ref(-1);
const list = ref<GlossaryItem[]>([]);
const { error, loading, run } = useRequest(
  () => {
    return fethGlossrayListByPage({
      currentPage: currentPage.value,
      pageSize: 10,
      searchText: searchValue.value.trim()
        ? searchValue.value.trim()
        : undefined,
    });
  },
  {
    onSuccess: (data) => {
      total.value = data.total;
      list.value = data.items;
      checkedItems.value = [];
    },
  }
);

const onPageChange = (page: number) => {
  currentPage.value = page;
  run();
};

const onConfirmDelete = async (ids: string[]) => {
  console.log('delete');
  await deleteGlossaryItems({ ids });
  await run();
  message.success(t('glossary.table.deleteSuccessTip'));
};

const checkedItems = ref<string[]>([]);
const isCheckedAll = computed(() => {
  return checkedItems.value.length === list.value.length;
});
const onItemCheckedChange = (id: string) => {
  if (checkedItems.value.includes(id)) {
    checkedItems.value = checkedItems.value.filter((item) => item !== id);
  } else {
    checkedItems.value = [...checkedItems.value, id];
  }
};

const onCheckedAll = () => {
  if (!isCheckedAll.value) {
    checkedItems.value = list.value.map((item) => item.id);
  } else {
    checkedItems.value = [];
  }
};

const createGlossaryButtonRef = ref<HTMLElement | null>(null);
const onCreate = () => {
  if (!createGlossaryButtonRef.value) {
    return;
  }
  createGlossaryCreateEditTippyVue({
    triggerEle: createGlossaryButtonRef.value,
    item: null,
    refresh: run,
  });
};

const onEdit = (item: GlossaryItem) => {
  if (!createGlossaryButtonRef.value) {
    return;
  }
  createGlossaryCreateEditTippyVue({
    triggerEle: createGlossaryButtonRef.value,
    item,
    refresh: run,
  });
};

const onDelete = () => {
  Modal.confirm({
    title: t(
      'glossary.table.deleteTip',
      checkedItems.value.length === 1 ? 1 : 2
    ),
    okText: t('viewer.delete'),
    cancelText: t('viewer.cancel'),
    onOk: async () => {
      await onConfirmDelete(checkedItems.value);
    },
  });
};
</script>
<style lang="less">
.glossary-search {
  .ant-input.ant-input-sm {
    border-color: rgb(217, 217, 217);
    color: rgba(0, 0, 0, 0.65) !important;
    &::placeholder {
      color: #aaa;
    }
  }
  .ant-input-group-addon .ant-input-search-button {
    border-color: rgb(217, 217, 217) !important;
    color: #4e5969 !important;
    height: 26px;
  }
}
.glossary-pagination {
  &.ant-pagination {
    .ant-pagination-prev,
    .ant-pagination-next {
      .ant-pagination-item-link {
        color: #4e5969;
      }
    }
    .ant-pagination-item:not(.ant-pagination-item-active) a {
      color: #4e5969;
    }
    .ant-pagination-item-link .ant-pagination-item-container {
      .ant-pagination-item-ellipsis {
        color: #4e5969;
      }
    }
  }
}
</style>
