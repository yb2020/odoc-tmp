<template>
  <div class="group-switcher">
    <a-dropdown
      overlayClassName="group-switcher-dropdown"
      :overlayStyle="{ maxHeight: '300px', overflow: 'hidden' }"
    >
      <a
        class="group-name flex"
        @click.prevent
      >
        <span class="name">{{ currentGroup.name }}</span>
        <DownOutlined />
      </a>
      <template #overlay>
        <a-menu>
          <PerfectScrollbar
            :options="{
              suppressScrollX: true,
            }"
            style="max-height: 300px"
          >
            <div
              v-for="(items, index) in groupList"
              :key="items.type"
            >
              <div class="groups">
                {{ $t(items.groupsNameI18n) }}
              </div>
              <a-menu-item
                v-for="item in items.groups"
                :key="item.id"
              >
                <a
                  href="javascript:;"
                  @click="handleSelect(item)"
                >{{
                  item.name
                }}</a>
              </a-menu-item>
              <a-menu-divider v-if="index < groupList.length" />
            </div>
          </PerfectScrollbar>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>
<script setup lang="ts">
import { DownOutlined } from '@ant-design/icons-vue';
import { $GroupProceed } from '~/src/api/group';

export interface GroupSwitcherItems {
  groupsName: string;
  groupsNameI18n: string;
  groups: $GroupProceed[];
  type: string;
}

defineProps<{
  groupList: GroupSwitcherItems[];
  currentGroup: $GroupProceed;
}>();

const emit = defineEmits<{
  (event: 'update:currentGroup', current: $GroupProceed): void;
}>();

const handleSelect = (item: $GroupProceed) => {
  emit('update:currentGroup', item);
};
</script>
<style lang="less" scoped>
.group-switcher {
  .group-name {
    color: rgba(255, 255, 255, 0.65);
    overflow: hidden;
    .name {
      white-space: nowrap;
      text-overflow: ellipsis;
      overflow: hidden;
    }
    .anticon-down {
      margin-left: 6px;
      vertical-align: top !important;
      line-height: 24px !important;
    }
  }
}
</style>
<style lang="less">
.group-switcher-dropdown {
  .groups {
    color: rgba(255, 255, 255, 0.3) !important;
    font-size: 12px;
    padding: 5px 12px;
  }
  .ant-dropdown-menu-title-content {
    font-size: 13px !important;
    color: rgba(255, 255, 255, 0.65);
    &:hover {
      color: inherit;
    }
  }
}
</style>
