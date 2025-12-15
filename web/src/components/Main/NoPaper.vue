<template>
  <div class="no-paper-wrap">
    <div class="tip">
      {{ $t('teams.emptyTip', { team: groupName || $t('teams.thisTeam') }) }}
    </div>
    <a-button
      type="primary"
      :loading="loading"
      @click="addToGroup"
    >
      <plus-outlined />{{ $t('teams.addToTeam') }}
    </a-button>
  </div>
</template>
<script setup lang="ts">
import { PlusOutlined } from '@ant-design/icons-vue';
import { Modal } from 'ant-design-vue';
import { createVNode, ref, computed } from 'vue';
import { addPaperToGroup } from '~/src/api/group';
import { store, selfNoteInfo } from '~/src/store';
import { BaseActionTypes } from '~/src/store/base';
import useGroupSettings from '~/src/hooks/UserSettings/useGroupSettings';
import { groupAnnotationController } from '~/src/stores/annotationStore';
import { useI18n } from 'vue-i18n';

const { groupSettings } = useGroupSettings();
const groupInfoMap = computed(() => store.state.base.groupInfoMap);
const groupName = computed(
  () => groupInfoMap.value[groupSettings.value.currentGroupId].name
);

const loading = ref(false);

const { t } = useI18n();

const addToGroup = async () => {
  const currentGroup =
    store.state.base.groupInfoMap[store.state.base.currentGroupId];
  Modal.confirm({
    title: t('teams.addToTeam'),
    content: createVNode('span', {}, [
      t('teams.tip1'),
      createVNode(
        'span',
        { style: { color: '#1f71e0' } },
        `${selfNoteInfo.value?.paperTitle ?? selfNoteInfo.value?.docName}`
      ),
      t('teams.tip2'),
      createVNode(
        'span',
        { style: { color: '#1f71e0' } },
        `${currentGroup?.name}`
      ),
    ]),
    onOk: async () => {
      loading.value = true;

      groupAnnotationController.clearGroupAnnotation();

      try {
        await addPaperToGroup({
          paperId: selfNoteInfo.value!.paperId,
          pdfId: selfNoteInfo.value!.pdfId,
          paperTitle: selfNoteInfo.value!.paperTitle,
          groupId: store.state.base.currentGroupId,
        });

        store.dispatch(`base/${BaseActionTypes.SWITCH_TO_GROUP}`, {
          groupId: store.state.base.currentGroupId,
          t,
        });
      } catch (error) {}
      loading.value = false;
    },
    cancelText: t('viewer.cancel'),
    okText: t('viewer.confirm'),
    okType: 'danger',
  });
  return;
};
</script>
<style lang="less">
.no-paper-wrap {
  flex: 1;
  background-color: #fff;
  display: flex;
  height: 100%;
  margin: 0 15%;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
  .tip {
    color: #1d2129;
    margin-bottom: 15px;
  }
  .question {
    position: absolute;
    top: 30px;
    right: 20px;
    color: #1d2129;
    font-size: 18px;
    cursor: pointer;
  }
}
</style>
