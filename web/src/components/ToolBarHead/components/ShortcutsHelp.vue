<template>
  <a-modal
    centered
    :visible="visible"
    :mask="false"
    :footer="null"
    :title="$t('viewer.shortcut')"
    :width="600"
    @cancel="close"
  >
    <div class="shortcuts-help">
      <a-table
        :pagination="false"
        table-layout="fixed"
        :dataSource="allShortcuts"
        :columns="[
          {
            title: '',
            dataIndex: 'name',
            width: 195,
          },
          {
            title: 'Windows',
            dataIndex: ['value', 'win32'],
          },
          {
            title: 'MacOS',
            dataIndex: ['value', 'darwin'],
          },
        ]"
      >
        <template #bodyCell="{ column, text, record }">
          <template v-if="'' === column.title">
            <span class="shortcuts-name">
              <component
                :is="record.icon"
                v-bind="record.iconAttrs"
                class="shortcuts-icon"
              />{{
                $t(record.i18n)
              }}
            </span>
          </template>
          <template v-if="['Windows', 'MacOS'].includes(column.title)">
            <ul class="shortcuts-key-list">
              <template v-for="item in text.split('+')">
                <li
                  v-if="isText(item)"
                  class="shortcuts-key-text"
                >
                  {{ $t(item) }}
                </li>
                <li
                  v-else
                  class="shortcuts-key-item"
                >
                  {{ item in shortcutTxtMap ? shortcutTxtMap[item] : item }}
                </li>
              </template>
            </ul>
          </template>
        </template>
      </a-table>
    </div>
  </a-modal>
</template>

<script lang="ts" setup>
import { useStore } from '@/store';
import { computed } from 'vue';
import { PAGE_ROUTE_NAME } from '../../../routes/type';
import { shortcutTxtMap } from '../../../store/shortcuts';

const props = defineProps<{
  visible: boolean;
}>();
const emit = defineEmits<{ (event: 'close'): void }>();

const close = () => emit('close');

const store = useStore();
const allShortcuts = computed(() => {
  const { shortcuts } = store.state.shortcuts[PAGE_ROUTE_NAME.NOTE];

  return Object.values(shortcuts).sort((a, b) => Number(a.order) - Number(b.order));
});

const isText = (string: string) => {
  return /^viewer\./.test(string)
};
</script>

<style lang="less" scoped>
.shortcuts {
  &-help:deep(.ant-table) {
    thead > tr > th {
      background-color: transparent;
      &::before {
        display: none;
      }
    }
    tbody > tr > td {
      background-color: transparent !important;
    }
    .ant-table-cell {
      padding: 12px 16px;

      &:first-child {
        padding: 12px 0;
      }
    }
  }

  &-icon {
    width: 18px;
    height: 18px;
    font-size: 18px;
    margin-right: 8px;
  }

  &-name {
    display: flex;
    align-items: center;
    line-height: 22px;
  }

  &-key-list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    align-items: center;
    font-size: 13px;
    line-height: 1;
  }

  &-key-item {
    padding: 5.5px 7px;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 1px;
    margin-right: 4px;

    &:last-child {
      margin-right: 0px;
    }
  }

  &-key-text {
    margin-right: 4px;
  }
}
</style>
