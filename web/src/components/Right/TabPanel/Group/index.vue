<template>
  <div class="group-container">
    <div class="flex items-center pr-10">
      <Switcher
        v-model:currentGroup="currentGroup"
        class="min-w-0 flex-1 px-4 py-2"
        :group-list="groupList"
      />
      <a
        class="tip"
        href="https://docs.qq.com/doc/DVWh0ZlpVWlNMUmFy"
        target="_blank"
      >
        <info-circle-outlined />{{ $t('teams.features') }}
      </a>
    </div>
    <Notes
      v-if="noteInfo"
      :tab="RightSideBarType.Group"
      :currentGroup="currentGroup"
      :active-tab="activeTab"
    />
  </div>
</template>
<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { useStore } from '~/src/store';
import Switcher, { GroupSwitcherItems } from './Switcher.vue';
import { BaseActionTypes } from '~/src/store/base';
import { $GroupProceed } from '~/src/api/group';
import Notes from './Notes.vue';
import { RightSideBarType } from '../type';
import { InfoCircleOutlined } from '@ant-design/icons-vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  groupList: GroupSwitcherItems[];
  switcherCurrentGroup: $GroupProceed;
  activeTab: RightSideBarType;
}>();

const currentGroup = ref(props.switcherCurrentGroup);

const emit = defineEmits<{
  (event: 'update:switcherCurrentGroup', group: $GroupProceed): void;
}>();

const store = useStore();

const noteInfo = computed(
  () => store.state.base.groupInfoMap[currentGroup.value.id]
);

const { t } = useI18n();

watch(
  currentGroup,
  async (newVal) => {
    await store.dispatch(`base/${BaseActionTypes.SWITCH_TO_GROUP}`, {
      groupId: newVal.id,
      t,
    });

    emit('update:switcherCurrentGroup', newVal);
  },
  {
    immediate: false,
  }
);
</script>

<style lang="less" scoped>
.group-container {
  height: 100%;
  overflow: hidden;
  position: relative;

  .tip {
    color: rgba(255, 255, 255, 0.65);
    margin: 0 4px;
    cursor: pointer;
    padding: 0 4px;

    .anticon {
      margin-right: 4px;
    }

    &:hover {
      background: rgba(255, 255, 255, 0.08);
    }
  }
}
</style>
