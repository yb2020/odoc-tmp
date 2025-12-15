<template>
  <Notes
    v-if="currentNoteInfo"
    :tab="RightSideBarType.Group"
    :activeTab="activeTab"
  />
</template>

<script lang="ts" setup>
import { currentNoteInfo, store } from '~/src/store';
import { onUnmounted } from 'vue';
import { RightSideBarType } from '../type';
import { NoteActionTypes } from '~/src/store/note';
import { OperationType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/GroupNote';
import Notes from '../Note/index.vue';
import { $GroupProceed } from '~/src/api/group';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import { groupAnnotationController } from '~/src/stores/annotationStore';

const props = defineProps<{
  activeTab: RightSideBarType;
  tab: string;
  currentGroup: $GroupProceed;
}>();

const clearTimer = (timer: NodeJS.Timer | null) => {
  groupAnnotationController.clearGroupAnnotation();

  timer && clearInterval(timer);

  timer = null;
};

let groupTimer: NodeJS.Timer | null = null;

const _getGroupNotes = async (time?: string) => {
  try {
    await groupAnnotationController.syncGroupAnnotation(time);
  } catch (err) {
    console.error(err);
  }
};

const getGroupNotes = async (time?: string) => {
  clearTimer(groupTimer);

  groupTimer = setInterval(async () => {
    _getGroupNotes(time);
  }, 5000);
};

let commentTimer: NodeJS.Timer | null = null;

const _getGroupNoteComments = async () => {
  const modifiedTime = store.state.note.commentModifiedTime;

  const params: any = {
    commentModifiedTime: modifiedTime,
    operationTypes:
      modifiedTime === '0'
        ? [OperationType[2]]
        : [OperationType[0], OperationType[1], OperationType[2]],
    groupNoteId: currentNoteInfo.value?.noteId,
  };

  try {
    await store.dispatch(`note/${NoteActionTypes.GET_GROUP_COMMENTS}`, params);
  } catch (err) {
    console.error(err);
  }
};

const getGroupNoteComments = async () => {
  clearTimer(commentTimer);

  commentTimer = setInterval(async () => {
    _getGroupNoteComments();
  }, 5000);
};

store.watch(
  (state) => state.base.currentGroupId,
  (val) => {
    if (val === SELF_NOTEINFO_GROUPID || !currentNoteInfo.value) {
      clearTimer(groupTimer);
      clearTimer(commentTimer);
      return;
    }

    _getGroupNotes('0');
    _getGroupNoteComments();

    getGroupNotes();
    getGroupNoteComments();
  },
  { immediate: true }
);

onUnmounted(() => {
  clearTimer(groupTimer);
  clearTimer(commentTimer);
});
</script>

<style></style>
