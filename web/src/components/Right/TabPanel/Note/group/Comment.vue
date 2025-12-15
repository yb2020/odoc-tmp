<template>
  <div
    class="group-comment-container"
    :data-group-comment-id="props.comment.commentId"
  >
    <div
      v-show="!showReply"
      :class="['comment', isCurrentUser ? 'current-comment' : '']"
      @click="handleCommentClick"
    >
      <a-avatar
        class="avatar"
        :src="comment.commentatorInfoView?.avatarCdnUrl"
        :size="24"
      >
        <template #icon>
          <UserOutlined />
        </template>
      </a-avatar>
      <div class="content-container">
        <div class="name">
          {{ comment.commentatorInfoView?.nickName }}
          <template v-if="comment.commentedUserName">
            <span class="split">{{ $t('teams.reply') }}</span>
            {{ comment.commentedUserName }}
          </template>
        </div>
        <div
          class="content"
          v-html="inputContent"
        />
      </div>
      <div
        v-if="comment.deleteAuthority"
        class="close"
        @click.stop="handleDelete"
      >
        <close-outlined class="iconguanbi" />
      </div>
    </div>

    <Reply
      v-show="showReply"
      ref="replyRef"
      :avatar="userInfo?.avatarUrl"
      @onBlur="handleBlur"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
import { useStore } from '@/store';
import { mdi } from '@idea/aiknowledge-markdown';
import { GroupComment } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/GroupNote';
import Reply from './Reply.vue';
import { NoteActionTypes } from '~/src/store/note';
import autosize from 'autosize';
import { Modal } from 'ant-design-vue';
import { UserOutlined, CloseOutlined } from '@ant-design/icons-vue';
import { AnnotationAll } from '~/src/stores/annotationStore/BaseAnnotationController';

const props = defineProps<{
  comment: GroupComment;
  showCommentReply: (comment: GroupComment) => void;
  note: AnnotationAll;
}>();
const store = useStore();
const userInfo = computed(() => store.state.user.userInfo);

const replyRef = ref();

const showReply = ref(false);

const isCurrentUser = computed(
  () => props.comment.commentatorInfoView?.userId === userInfo.value?.id
);

const inputContent = computed(() => {
  const htmlIdea = mdi.render(props.comment.comment || '');

  return htmlIdea.replace(/\n$/, '').replace(/\n/g, '<br/>');
});

const handleCommentClick = () => {
  if (isCurrentUser.value) {
    showReply.value = true;

    replyRef.value.input = replyRef.value.input || props.comment.comment || '';

    setTimeout(() => {
      replyRef.value.inputRef.focus();
      autosize.update(replyRef.value.inputRef);
    }, 100);
  } else {
    props.showCommentReply(props.comment);
  }
};

const handleBlur = async (input: string) => {
  showReply.value = false;

  await store.dispatch(`note/${NoteActionTypes.UPDATE_COMMENT}`, {
    commentId: props.comment.commentId,
    commentContent: input,
    markId: props.note.uuid,
  });
};

const handleDelete = async () => {
  Modal.confirm({
    title: '确定删除评论？',
    onOk: async () => {
      await store.dispatch(`note/${NoteActionTypes.DELETE_COMMENT}`, {
        commentId: props.comment.commentId,
        markId: props.note.uuid,
      });
    },
    okButtonProps: {
      danger: true,
    },
    cancelButtonProps: { type: 'primary' },
    cancelText: '取消',
    okText: '删除',
  });
};
</script>

<style lang="less" scoped>
.group-comment-container {
  word-break: break-all;
  position: relative;
}

.comment {
  cursor: pointer;
  display: flex;
  padding: 8px 0;
  position: relative;

  .content-container {
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;

    .name {
      font-size: 13px;
      color: rgba(255, 255, 255, 45%);
      line-height: 20px;

      .split {
        font-weight: 400;
        color: #ffffff;
        line-height: 20px;
      }
    }

    .content {
      font-size: 14px;
      color: #ffffff;
      line-height: 22px;
    }
  }

  .avatar {
    flex: 0 0 24px;
    margin-right: 8px;
  }

  .close {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: #76797d;

    display: none;
    align-items: center;
    justify-content: center;
    position: absolute;
    right: 0;
    top: 0;
    transform: translate(50%, -50%);
    z-index: 2;
    cursor: pointer;

    .iconguanbi {
      font-size: 8px;
      color: rgba(38, 38, 38, 1);
      transform: scale(0.5);
    }
  }

  &:hover {
    background: rgba(255, 255, 255, 0.04);
    border-radius: 0px 2px 2px 0px;

    .close {
      display: flex;
    }
  }
}

.current-comment {
  cursor: text;
}
</style>
