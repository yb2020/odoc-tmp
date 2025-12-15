<template>
  <div
    class="question"
    @click="handleQuestionClick"
  >
    <div class="title-container">
      <div class="title">
        {{ question.title }}
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
    <div class="user">
      提问者：
      <span
        @click="handleUser"
      >{{ isCurrentUser ? '我' : question.userInfo?.userInfo?.realName || 'unknown' }}</span>
      <span
        v-if="isShowAuthorLabel"
        class="label"
      >作者</span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Modal, message } from 'ant-design-vue';
import { computed, ref, nextTick, createVNode } from 'vue';
import { useStore } from '@/store';
import { deleteQuestion } from '~/src/api/question';
import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
import { PaperQuestion } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/QuestionAnswer';
import { getDomainOrigin } from '~/src/util/env';
import { goPathPage } from '~/src/common/src/utils/url';

const props = defineProps<{ question: PaperQuestion, isAllowedDelete: boolean }>()

const emit = defineEmits<{
  (event: 'delete'): void,
  (event: 'questionClick'): void,
}>()

const handleDelete = () => {
  Modal.confirm({
    title: '确定要删除问题？',
    content: '',
    icon: createVNode(ExclamationCircleOutlined),
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      try {
        await deleteQuestion({
          questionId: props.question.questionId,
        });

        message.success('删除成功！');

        emit('delete');
      } catch (err) {
        message.error('删除失败，请稍后再试');
      }
    },
    onCancel() { },
  });
};

const store = useStore();
const user = computed(() => store.state.user);

const isCurrentUser = computed(() => {
  const item = props.question;
  const userId = user.value?.userInfo?.id;
  const questionUserId = item?.userInfo?.userInfo?.userId;

  if (!userId || !questionUserId) {
    return false;
  }

  return userId === questionUserId;
});

const handleUser = () => {
  const userId = props.question?.userInfo?.userInfo?.userId;

  if (!userId) return;

  if (isCurrentUser.value) return;

  goPathPage(`${getDomainOrigin()}/user/visitor/${userId}`);
};

const handleQuestionClick = () => {
  emit('questionClick');
};

const isShowAuthorLabel = computed(() => props.question?.userInfo?.isAuthentication && props.question?.userInfo?.isPaperAuthor)


</script>

<style lang="less" scoped>
.question {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
  // border: 1px solid #e4e7ed;
  padding: 10px;
  flex: 0 0 auto;
  margin: 10px 10px 0;
  cursor: pointer;

  .title-container {
    font-size: 14px;
    color: #f0f2f5;
    display: flex;
    justify-content: space-between;
    position: relative;

    // .title {
    //   display: -webkit-box;
    //   -webkit-line-clamp: 3;
    //   -webkit-box-orient: vertical;
    //   overflow: hidden;
    //   text-overflow: ellipsis;
    //   word-break: break-all;
    // }

    .dot {
      flex: 0 0 auto;
      color: #a8adb3;
      font-size: 16px;
      margin-left: 5px;
      height: 24px;
    }

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
      cursor: pointer;

      .option {
        width: 66px;
        height: 24px;
        display: flex;
        align-items: center;
        justify-content: center;

        font-size: 14px;
        font-weight: 400;
        color: #a8adb3;

        &:hover {
          background: #f2f2f2;
        }
      }
    }
  }

  .user {
    font-size: 13px;
    font-weight: 400;
    color: #a8adb3;
    margin-top: 10px;
    .label {
      color: #428DE6;
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
}
</style>
