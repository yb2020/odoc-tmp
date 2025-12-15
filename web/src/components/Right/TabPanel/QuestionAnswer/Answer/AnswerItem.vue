<template>
  <div
    class="item"
    @click="handleItemClick"
  >
    <div class="avatar">
      <a-avatar :src="item.answerAuthor?.avatarUrl">
        <template #icon>
          <UserOutlined />
        </template>
      </a-avatar>
    </div>

    <div class="text-container">
      <div class="name">
        <div class="reply-user">
          <span
            class="user-click"
            @click.stop="handleUser('answer')"
          >{{
            answerName
          }}</span>
          <span
            v-if="isShowAnswerAuthorLabel"
            class="label"
          >作者</span>
          <template v-if="item.replyUser">
            <span class="reply-text">{{ $t('teams.reply') }}</span>
            <span
              class="user-click"
              @click.stop="handleUser('reply')"
            >{{
              isCurrentReplyUser ? '我' : item.replyUser.nickName
            }}</span>
            <span
              v-if="isShowReplyAuthorLabel"
              class="label"
            >作者</span>
          </template>
        </div>
        <a-dropdown placement="bottomRight">
          <span
            v-if="isCurrentUser || isAllowedDelete"
            class="dot"
            @click.prevent
            @click.stop
          >
            ···
          </span>
          <template #overlay>
            <a-menu>
              <a-menu-item @click.stop="handleDelete">
                <a href="javascript:;">删除</a>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>

      <div class="text-sub-container">
        <MarkdownViewer
          :raw="item.content || ''"
          class="text"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Modal, message } from 'ant-design-vue';
import { computed, createVNode } from 'vue';
import { useStore } from '@/store';
import { deleteAnswer } from '~/src/api/question';
import { UserOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue';
import { PaperAnswer } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/QuestionAnswer';
import { getDomainOrigin } from '~/src/util/env';
import { goPathPage } from '~/src/common/src/utils/url';
import { MarkdownViewer } from '@idea/aiknowledge-markdown';

const props = defineProps<{
  item: PaperAnswer;
  isAllowedDelete: boolean;
}>();

const emit = defineEmits<{
  (event: 'itemClick', info: { item: PaperAnswer; answerName: string }): void;
  (event: 'delete'): void;
}>();

const store = useStore();

const user = computed(() => store.state.user);

const isCurrentUser = computed(() => {
  const item = props.item;
  const userId = user.value?.userInfo?.id;
  const answerUserId = item?.answerAuthor?.id;
  if (!userId || !answerUserId) {
    return false;
  }

  return userId === answerUserId;
});

const isCurrentReplyUser = computed(() => {
  const item = props.item;
  const userId = user.value?.userInfo?.id;
  const replyUserId = item?.replyUser?.id;

  if (!userId || !replyUserId) {
    return false;
  }

  return userId === replyUserId;
});

const answerName = computed(() =>
  isCurrentUser.value ? '我' : props.item?.answerAuthor?.nickName || ''
);

const handleDelete = () => {
  Modal.confirm({
    title: '确定要删除回答？',
    content: '',
    icon: createVNode(ExclamationCircleOutlined),
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      try {
        await deleteAnswer({
          answerId: props.item.answerId,
        });
        emit('delete');
        message.success('删除成功！');
      } catch (err) {
        message.error('删除失败！请稍后再试');
      }
    },
    onCancel() {},
  });
};

const handleItemClick = () => {
  emit('itemClick', { item: props.item, answerName: answerName.value });
};

const handleUser = (type: 'answer' | 'reply') => {
  if (type === 'answer') {
    if (isCurrentUser.value) return;

    const userId = props.item?.answerAuthor?.id;

    if (!userId) return;

    goPathPage(`${getDomainOrigin()}/user/visitor/${userId}`);

    return;
  }

  if (type === 'reply') {
    if (isCurrentReplyUser.value) return;

    const userId = props.item?.replyUser?.id;

    if (!userId) return;

    goPathPage(`${getDomainOrigin()}/user/visitor/${userId}`);

    return;
  }
};

const isShowAnswerAuthorLabel = computed(
  () =>
    props.item?.answerAuthor?.isAuthentication &&
    props.item?.answerAuthor?.isPaperAuthor
);

const isShowReplyAuthorLabel = computed(
  () =>
    props.item?.replyUser?.isAuthentication &&
    props.item?.replyUser?.isPaperAuthor
);
</script>

<style lang="less" scoped>
.item {
  border-radius: 4px;
  display: flex;
  padding: 10px;
  cursor: pointer;

  &:hover {
    background: rgba(255, 255, 255, 0.05);
  }

  .avatar {
    height: 40px;
    flex: 0 0 40px;
    margin-right: 4px;

    img {
      width: 100%;
      height: 100%;
      border-radius: 50%;
    }
  }

  .text-container {
    flex: 1 1 auto;
    background: rgba(255, 255, 255, 0.05);
    margin-left: 12px;
  }

  .name {
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    position: relative;

    font-size: 14px;
    font-weight: 600;
    color: #a8adb3;
    padding: 0 14px;

    .options {
      position: absolute;
      bottom: 0;
      right: 0;
      transform: translate(0, 100%);

      background: #ffffff;
      box-shadow: 0px 4px 12px 0px rgba(0, 0, 0, 0.2);
      border-radius: 4px;
      border: 1px solid #dcdfe6;
      outline: none;

      .option {
        width: 66px;
        height: 24px;
        background: rgba(242, 242, 242, 74%);
        display: flex;
        align-items: center;
        justify-content: center;

        font-size: 14px;
        font-weight: 400;
        color: #000000;
        cursor: pointer;

        &:hover {
          background: #f2f2f2;
        }
      }
    }

    .dot {
      font-size: 16px;
      cursor: pointer;
    }

    .reply-user {
      .reply-text {
        font-size: 14px;
        font-weight: 400;
        color: #f0f2f5;
        margin: 0 4px;
      }

      .user-click {
        cursor: pointer;
      }
    }
  }

  .text-sub-container {
    padding: 10px;
  }

  .text {
    font-size: 14px;
    font-weight: 400;
    color: #f0f2f5;

    display: -webkit-box;
    -webkit-line-clamp: 100;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
  }
  .label {
    color: #428de6;
    background: rgba(0, 0, 0, 0.15);
    padding: 2px 8px;
    font-size: 12px;
    border-radius: 10px;
    font-weight: 400;
    line-height: 18px;
    margin-left: 5px;
    word-break: keep-all;
  }
}
</style>
