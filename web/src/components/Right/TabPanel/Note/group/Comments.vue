<template>
  <div class="comments-container">
    <Comment
      v-for="comment in showComments"
      :comment="comment"
      :showCommentReply="handleShowCommentReply"
      :note="note"
    />
    <Reply
      v-show="showCommentReply"
      ref="replyRef"
      :avatar="useInfo?.avatarUrl"
      :placeHolder="`${$t('teams.reply')} ${
        commentUser?.commentatorInfoView?.nickName
      }`"
      @onBlur="handleBlur"
    />
    <div
      v-if="comments?.length > 2"
      class="last-comments"
      @click="handleLast"
    >
      <div>
        {{ showLast ? $t('viewer.collapse') : $t('viewer.expand') }}
        {{ comments?.length - 2 }}
        {{ $t('teams.comments', comments?.length - 2) }}
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue';
import { AnnotationAll } from '~/src/stores/annotationStore/BaseAnnotationController';
import { useStore } from '@/store';
import Comment from './Comment.vue';
import { GroupComment } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/GroupNote';
import Reply from './Reply.vue';
import { message } from 'ant-design-vue';
import { NoteActionTypes } from '~/src/store/note';
import autosize from 'autosize';

const props = defineProps<{
  note: AnnotationAll;
  queryTargetCommentId: string;
  positioningComment: () => void;
}>();

const store = useStore();

const useInfo = computed(() => store.state.user.userInfo);

const comments = computed(() => store.state.note.comments[props.note.uuid!]);

const showLast = ref(false);

const showComments = computed(() =>
  showLast.value ? comments.value : (comments.value || []).slice(0, 2)
);

const handleLast = () => {
  showLast.value = !showLast.value;
};

const showCommentReply = ref(false);

const commentUser = ref<GroupComment | null>(null);

const handleShowCommentReply = (comment: GroupComment) => {
  showCommentReply.value = true;
  commentUser.value = comment;

  setTimeout(() => {
    replyRef.value.inputRef.focus();
    autosize.update(replyRef.value.inputRef);
  }, 100);
};

const replyRef = ref();

const handleBlur = async (input: string) => {
  showCommentReply.value = false;

  if (input === '') {
    return;
  }

  if (input.length > 2000) {
    message.error('单个评论字数最多2000字！');
    return;
  }

  const params: any = {
    commentContent: input,
    groupId: store.state.base.currentGroupId,
    commentId: commentUser.value?.commentId,
    commentedUserName: commentUser.value?.commentatorInfoView?.nickName,
    markId: props.note.uuid,
  };

  await store.dispatch(`note/${NoteActionTypes.ADD_COMMENT}`, params);
};

const isPositioned = ref<boolean>(false);

watch(
  comments,
  (newVal) => {
    if (newVal && props.queryTargetCommentId && !isPositioned.value) {
      const findIndex = newVal.findIndex(
        (item: any) => item.commentId === props.queryTargetCommentId
      );

      if (findIndex > 1) {
        showLast.value = true;
      }

      if (findIndex !== -1) {
        setTimeout(() => {
          isPositioned.value = true;

          props.positioningComment();
        }, 300);
      }
    }
  },
  {
    immediate: true,
  }
);
</script>

<style lang="less" scoped>
.last-comments {
  height: 22px;

  font-size: 13px;
  font-family: PingFangSC-Regular, PingFang SC;
  font-weight: 400;
  color: rgba(255, 255, 255, 65%);
  line-height: 18px;
  padding: 2px 8px;
  display: flex;
  align-items: center;
  justify-content: center;

  div {
    background: rgba(255, 255, 255, 4%);
    padding: 2px 8px;
  }
}
</style>
