<template>
  <a-dropdown
    v-model:visible="visible"
    :trigger="['click']"
    :overlayStyle="{ width: '140px' }"
    placement="bottomRight"
  >
    <slot />
    <template #overlay>
      <a-menu @click="handleMenuClick">
        <VueDraggableNext
          :list="tabs"
          @change="onChange"
        >
          <a-menu-item
            v-for="item in tabs"
            :key="item.key"
          >
            <a
              href="javascript:;"
              class="right-pane-more-item"
            ><span> <HolderOutlined />{{ $t(item.title) }} </span>
              <check-outlined v-if="item.shown" />
            </a>
          </a-menu-item>
        </VueDraggableNext>
      </a-menu>
    </template>
  </a-dropdown>
</template>
<script setup lang="ts">
import { CheckOutlined, HolderOutlined } from '@ant-design/icons-vue';
import { ref, watch } from 'vue';
import type { MenuProps } from 'ant-design-vue';
import {
  RightSideTabItem,
  useRightSideTabSettings,
} from '~/src/hooks/UserSettings/useSideTabSettings';
import { RightSideBarType } from '../type';
import { VueDraggableNext } from 'vue-draggable-next';

const props = defineProps<{
  allTabs: RightSideTabItem[];
}>();

const { toggleTab, updateRightTabBars } = useRightSideTabSettings();

const visible = ref(false);

const handleMenuClick: MenuProps['onClick'] = (e) => {
  toggleTab(e.key as RightSideBarType);
};

const tabs = ref(props.allTabs.filter((tab) => tab.sortable));

watch(
  () => props.allTabs,
  (newVal) => {
    tabs.value = newVal.filter((tab) => tab.sortable);
  }
);

const onChange = () => {
  updateRightTabBars(tabs.value);
};
</script>
<style lang="less" scoped>
.right-pane-more-item {
  display: flex;
  width: 100%;
  justify-content: space-between;
  align-items: center;

  .anticon-holder {
    color: #969ca3;
    margin-right: 10px;
  }

  .anticon-check {
    color: rgba(255, 255, 255, 0.65);
  }
}
</style>
