<template>
  <div class="right-panel js-drawer">
    <a-tabs
      :activeKey="activeTab"
      :centered="false"
      class="tabs"
      type="card"
      @change="handleChange"
    >
      <template #rightExtra>
        <more :all-tabs="tabs">
          <more-outlined class="right-icon" />
        </more>
        <close-outlined
          class="right-icon"
          @click="handleClose"
        />
      </template>
      <a-tab-pane
        v-for="tab in showTabList"
        :key="tab.key"
        :tab="$t(tab.title)"
        :force-render="tab.key === RightSideBarType.Note"
      >
        <material-vue
          v-if="tab.key === RightSideBarType.Matirial"
          :paper-id="paperId"
          :pdf-id="pdfId"
        />
        <component
          :is="isOwner ? NotesVue : NoteVue"
          v-else-if="tab.key === RightSideBarType.Note"
          :tab="RightSideBarType.Note"
          :active-tab="activeTab"
        />
        <question-answer-vue
          v-else-if="tab.key === RightSideBarType.Question"
          :paper-id="paperId"
          :is-private-paper="isPrivatePaper"
        />
        <group-vue
          v-else-if="tab.key === RightSideBarType.Group && currentGroup"
          v-model:switcherCurrentGroup="currentGroup"
          :groupList="groupList"
          :activeTab="activeTab"
        />
        <learn-vue
          v-else-if="tab.key === RightSideBarType.Learn"
          :paper-id="paperId"
          :is-private-paper="isPrivatePaper"
        />
        <copilot-vue
          v-else-if="tab.key === RightSideBarType.Copilot"
          :clip-selecting="clipSelecting"
          :clip-action="clipAction"
        />
        <translate-vue
          v-else-if="tab.key === RightSideBarType.Translate"
          :clip-selecting="clipSelecting"
          :clip-action="clipAction"
        />
      </a-tab-pane>
    </a-tabs>
  </div>
</template>
<script setup lang="ts">
import { RightSideBarType } from './type';
import MaterialVue from './Material/index.vue';
import NotesVue from './Note/wrapper.vue';
import NoteVue from './Note/index.vue';
import QuestionAnswerVue from './QuestionAnswer/index.vue';
import { CloseOutlined, MoreOutlined } from '@ant-design/icons-vue';
import { computed, ref, watch } from 'vue';
import { useInitGroupTabPane } from '~/src/hooks/GroupTabPane/useInitGroupTabPane';
import GroupVue from './Group/index.vue';
import { isOwner, selfNoteInfo, store } from '~/src/store';
import { BaseActionTypes, BaseMutationTypes } from '~/src/store/base';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import { $GroupProceed } from '~/src/api/group';
import More from './More/index.vue';
import { useRightSideTabSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import useGroupSettings from '~/src/hooks/UserSettings/useGroupSettings';
import { reportRightTab } from '~/src/api/report';
import { clearArrow } from '~/src/dom/arrow';
import { destroyEditOverlay, hideTooltip } from '~/src/dom/tooltip';
import LearnVue from './Learn/index.vue';
import { isEmptyPaperId } from '~/src/api/helper';
import CopilotVue from './Copilot/index.vue';
import { usePdfStore } from '~/src/stores/pdfStore';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { useI18n } from 'vue-i18n';
import { useClip } from '~/src/hooks/useHeaderScreenShot';
import TranslateVue from './Translate/index.vue';

const props = defineProps<{
  paperId: string;
  pdfId: string;
  isPrivatePaper: boolean;
  clipSelecting: boolean;
  clipAction: ReturnType<typeof useClip>['clipAction'];
}>();

const { t } = useI18n();

const emit = defineEmits<{
  (event: 'visibleChange', val: boolean): void;
}>();

const pdfStore = usePdfStore();
const annotationStore = useAnnotationStore();

const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(props.pdfId);
});

const handleClose = () => {
  emit('visibleChange', false);
};

const { activeTab, activeSubTab, setSideTabSetting, tabs, toggleDisabledTab } =
  useRightSideTabSettings();

const handleChange = async (tab: RightSideBarType) => {
  refreshNoteList();

  setSideTabSetting({ tab });
};

const currentGroup = ref<$GroupProceed>();
const { groupSettings, setCurrentGroupId } = useGroupSettings();
const { createdGroupList, joinedGroupList, getGroupList } =
  useInitGroupTabPane();

(async () => {
  if (!isOwner.value) {
    await showSelfNoteList();
    ifNoteTabGetNoteList();
    return;
  }

  if (isEmptyPaperId(props.paperId)) {
    disableGroupTabAndShowSelfNoteList();
    return;
  }

  const allGroupList = await getGroupList();
  if (!allGroupList) {
    disableGroupTabAndShowSelfNoteList();
    return;
  }

  const group =
    allGroupList.find(
      (group) => group.id === groupSettings.value.currentGroupId
    ) || allGroupList[0];

  if (activeTab.value !== RightSideBarType.Group) {
    currentGroup.value = group;
    await showSelfNoteList();
    ifNoteTabGetNoteListAndTagList();
    return;
  }

  await store.dispatch(`base/${BaseActionTypes.SWITCH_TO_GROUP}`, {
    groupId: group.id,
    t,
  });
  currentGroup.value = group;

  ifNoteTabGetNoteListAndTagList();

  reportRightTab(activeTab.value);

  async function showSelfNoteList() {
    store.commit(`base/${BaseMutationTypes.SET_NOTE_INFO}`, {
      groupId: SELF_NOTEINFO_GROUPID,
      noteInfo: selfNoteInfo.value,
    });

    reportRightTab(activeTab.value, activeSubTab.value);

    await store.dispatch(`base/${BaseActionTypes.SWITCH_TO_GROUP}`, {
      groupId: SELF_NOTEINFO_GROUPID,
      t,
    });
  }

  async function disableGroupTabAndShowSelfNoteList() {
    toggleDisabledTab(RightSideBarType.Group, true);
    await showSelfNoteList();
    ifNoteTabGetNoteListAndTagList();
  }
})();

const showTabList = computed(() => tabs.value.filter((tab) => tab.shown));

watch(currentGroup, (newVal, oldVal) => {
  if (newVal && newVal.id !== oldVal?.id) {
    setCurrentGroupId(newVal.id);
  }
});

watch([activeTab, activeSubTab], () => {
  console.log('activeTab', activeTab.value, activeSubTab.value);
  clearArrow();
  destroyEditOverlay(pdfViewerRef.value);
  hideTooltip();

  if (activeTab.value === RightSideBarType.Translate) {
    // do nothing
  } else if (activeTab.value !== RightSideBarType.Group) {
    store.dispatch(`base/${BaseActionTypes.SWITCH_TO_GROUP}`, {
      groupId: SELF_NOTEINFO_GROUPID,
      t,
    });
  } else if (currentGroup.value) {
    store.dispatch(`base/${BaseActionTypes.SWITCH_TO_GROUP}`, {
      groupId: currentGroup.value.id,
      t,
    });
  }

  reportRightTab(activeTab.value, activeSubTab.value);

  refreshNoteList();
});

const groupList = computed(() => {
  const list = [];
  if (createdGroupList.value?.length) {
    list.push({
      groupsName: '我创建的群组',
      groupsNameI18n: 'teams.createdTeams',
      groups: createdGroupList.value,
      type: 'created',
    });
  }
  if (joinedGroupList.value?.length) {
    list.push({
      groupsName: '我加入的群组',
      groupsNameI18n: 'teams.joinedTeams',
      groups: joinedGroupList.value,
      type: 'joined',
    });
  }
  return list;
});

function onlyGetNoteList() {
  annotationStore.controller.loadAnnotationMap();
}

function ifNoteTabGetNoteList() {
  if (activeTab.value === RightSideBarType.Note) {
    onlyGetNoteList();
  }
}

function ifNoteTabGetNoteListAndTagList() {
  if (activeTab.value === RightSideBarType.Note) {
    onlyGetNoteList();
    annotationStore.refreshTagList();
  }
}

function refreshNoteList() {
  if (isOwner.value) {
    ifNoteTabGetNoteListAndTagList();
  } else {
    ifNoteTabGetNoteList();
  }
}
</script>
<style lang="less" scoped>
.right-panel {
  display: flex;
  height: 100%;
  background-color: var(--site-theme-pdf-panel);

  .right-panel-tab {
    width: 100%;
  }

  :deep(.ant-tabs-content) {
    height: 100%;
  }

  .tabs.ant-tabs.ant-tabs-card {
    width: 100%;

    :deep(> .ant-tabs-nav) {
      margin: 0;
      border-bottom: 0.5px solid var(--site-theme-divider);
      .ant-tabs-tab {
        border: 0;
        margin: 0;
        padding: 8px 22px;
        border-right: 1px solid var(--site-theme-divider);
        background-color: var(--site-theme-pdf-panel);
        color: var(--site-theme-pdf-panel-text);
        &.ant-tabs-tab-active {
          border-right: none;
          background-color: var(--site-theme-pdf-tab-active-bg);
          color: var(--site-theme-pdf-tab-active-text) !important;
          .ant-tabs-tab-btn {
            color: var(--site-theme-pdf-tab-active-text) !important;
            text-shadow: 0 0 0.25px currentcolor;
          }
          
          /* 确保最高优先级，防止被其他样式覆盖 */
          &.ant-tabs-tab-active .ant-tabs-tab-btn,
          &.ant-tabs-tab-active > div,
          &.ant-tabs-tab-active * {
            color: var(--site-theme-pdf-tab-active-text) !important;
          }
        }
      }
    }
  }

  .right-icon {
    font-size: 16px;
    padding: 4px;
    cursor: pointer;
    margin-left: 12px;
    color: var(--site-theme-pdf-panel-text);
    &:hover {
      background-color: var(--site-theme-bg-hover);
    }
  }
}
</style>
