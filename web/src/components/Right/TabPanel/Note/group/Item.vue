<template>
  <div
    ref="noteRef"
    class="note"
    :data-pdf-annotate-id="note.uuid"
    @dblclick.prevent.stop
  >
    <div
      class="left"
      :style="{ background: isFocus ? note.color : '#65676c' }"
    />

    <div class="right">
      <div
        v-if="
          note.type === ToolBarType.select &&
            annotationStore.activeColorMap.ref &&
            !isNotRectStr
        "
        class="ref-content"
        @click="handleContentClick"
      >
        {{ (note as AnnotationSelect).rectStr }}
      </div>

      <div
        v-if="
          note.type === ToolBarType.rect && annotationStore.activeColorMap.ref
        "
        class="ref-picture"
        @click="handleContentClick"
      >
        <img
          :src="(note as AnnotationRect).picUrl"
          alt=""
        >
      </div>

      <div
        v-show="!showReply"
        :class="['note-content', isCurrentUser ? 'current-note-content' : '']"
        @click="handleNoteClick"
      >
        <a-avatar
          class="avatar"
          :src="note.commentatorInfoView?.avatarCdnUrl"
          :size="24"
        >
          <template #icon>
            <UserOutlined />
          </template>
        </a-avatar>
        <div
          v-if="note.idea"
          class="content-container"
        >
          <div class="name">
            {{ note.commentatorInfoView?.nickName }}
          </div>
          <div
            class="content"
            v-html="inputContent"
          />
        </div>
      </div>

      <Reply
        v-show="showReply"
        ref="replyRef"
        :avatar="useInfo?.avatarUrl"
        :placeHolder="$t('teams.inputTip')"
        @onBlur="handleBlur"
      />

      <Comments
        :note="note"
        :queryTargetCommentId="queryTargetCommentId || ''"
        :positioningComment="handlePositioningComment"
      />

      <Reply
        v-show="showNoteReply"
        ref="noteReplyRef"
        :avatar="useInfo?.avatarUrl"
        :placeHolder="
          isCurrentUser
            ? $t('teams.inputTip')
            : `${$t('teams.reply')} ${note?.commentatorInfoView?.nickName}`
        "
        @onBlur="handleNoteReplyBlur"
      />
    </div>

    <div
      v-if="note.deleteAuthority"
      class="close"
      @click.stop="handleDelete"
    >
      <close-outlined class="iconguanbi" />
    </div>

    <Dot :item="note" />
  </div>
</template>

<script setup lang="ts">
import {
  AnnotationAll,
  AnnotationSelect,
  AnnotationRect,
} from '~/src/stores/annotationStore/BaseAnnotationController';
import { message, Modal } from 'ant-design-vue';
import { computed, ref, watch, onMounted } from 'vue';
import { mdi } from '@idea/aiknowledge-markdown';
import { clearArrow, connect } from '~/src/dom/arrow';
import { currentNoteInfo, useStore } from '~/src/store';
import { NoteActionTypes } from '~/src/store/note';
import Reply from './Reply.vue';
import autosize from 'autosize';
import Dot from './Dot.vue';
import { NO_ANNOTATION_ID } from '@/constants';
import {
  attrSelector,
  PDF_ANNOTATE_ID,
  ToolBarType,
} from '@idea/pdf-annotate-core';
import scrollIntoView from 'scroll-into-view-if-needed';
import { useRightSideTabSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import Comments from './Comments.vue';
import { UserOutlined, CloseOutlined } from '@ant-design/icons-vue';
import { usePdfStore } from '~/src/stores/pdfStore';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const pdfStore = usePdfStore();
const annotationStore = useAnnotationStore();

const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(currentNoteInfo.value?.pdfId);
});

const inputContent = computed(() => {
  const htmlIdea = mdi.render(props.note.idea || '');

  return htmlIdea.replace(/\n$/, '').replace(/\n/g, '<br/>');
});

const props = defineProps<{
  note: AnnotationAll;
}>();

const store = useStore();

const useInfo = computed(() => store.state.user.userInfo);

const scrollToAnnotation = () => {
  const scroller = pdfViewerRef.value?.getScrollController();
  const container = pdfViewerRef.value?.getDocumentViewer().container;
  if (!scroller || !container) {
    return;
  }

  scroller.goToPage(props.note.pageNumber);

  const annotateEle = container.querySelector(
    `[${PDF_ANNOTATE_ID}="${props.note.uuid}"]`
  );

  annotateEle &&
    scrollIntoView(annotateEle, {
      block: 'center',
      inline: 'center',
    });

  if (!isNotRectStr.value) {
    setTimeout(() => {
      connect(props.note.uuid!);
    }, 300);
  }
};

const isCurrentUser = computed(
  () => props.note.commentatorInfoView?.userId === useInfo.value?.id
);

const handleAddCommentClick = () => {
  showNoteReply.value = true;

  setTimeout(() => {
    noteReplyRef.value.inputRef.focus();
    autosize.update(noteReplyRef.value.inputRef);
  }, 100);
};

// 对笔记做评论
const showNoteReply = ref(false);
const handleContentClick = async () => {
  scrollToAnnotation();

  // 这个笔记是当前用户
  if (isCurrentUser.value) {
    // 且没有写笔记
    if (!props.note.idea) {
      handleWrite();
    }
  } else {
    // 不是则对笔记添加评论
    handleAddCommentClick();
  }
};

// 写笔记
const handleNoteReplyBlur = async (input: string) => {
  showNoteReply.value = false;

  if (input === '') {
    return;
  }

  if (input.length > 2000) {
    message.error(t('teams.limitTip'));
    return;
  }

  const params: any = {
    commentContent: input,
    groupId: store.state.base.currentGroupId,
    uuid: props.note.uuid,
    commentId: '0',
    commentedUserName: '',
    markId: props.note.uuid,
  };

  await store.dispatch(`note/${NoteActionTypes.ADD_COMMENT}`, params);
};

// 修改笔记
const handleNoteClick = () => {
  // 点击笔记如果不是本人，就代表是对笔记写评论
  if (!isCurrentUser.value) {
    handleAddCommentClick();
    return;
  }

  _handleNoteClick();
};

const _handleNoteClick = () => {
  scrollToAnnotation();
  handleWrite();
};

const showReply = ref(false);
const replyRef = ref();

const noteReplyRef = ref();

const isFocus = computed(
  () => annotationStore.currentAnnotationId === props.note.uuid
);

const isNotRectStr = computed(
  () =>
    props.note.type === ToolBarType.select &&
    (props.note as AnnotationSelect).rectStr === ''
);

const { sideTabSettings } = useRightSideTabSettings();

watch(isFocus, (val) => {
  if (val && sideTabSettings.value.shown) {
    _handleNoteClick();
  }
});

const handleWrite = () => {
  showReply.value = true;
  replyRef.value.input = replyRef.value.input || props.note.idea || '';

  setTimeout(() => {
    replyRef.value.inputRef.focus();
    autosize.update(replyRef.value.inputRef);
  }, 100);
};

const handleBlur = async (input: string) => {
  if (input.length > 2000) {
    message.error(t('teams.limitTip'));
    return;
  }

  const endReply = () => {
    showReply.value = false;
    annotationStore.currentAnnotationId = '';
  };

  if (input === '') {
    if (props.note.uuid === NO_ANNOTATION_ID) {
      annotationStore.controller.localDeleteAnnotation(
        props.note.uuid,
        props.note.pageNumber
      );
    }

    endReply();
    return;
  }

  endReply();

  if (props.note.uuid !== NO_ANNOTATION_ID) {
    await annotationStore.controller.patchAnnotation(props.note.uuid, {
      idea: input,
    });
    return;
  }

  const uuid = await annotationStore.controller.onlineSaveAnnotation(
    props.note.documentId,
    props.note.pageNumber,
    {
      ...props.note,
      idea: input,
    }
  );
  if (uuid) {
    annotationStore.controller.localPatchAnnotation(
      NO_ANNOTATION_ID,
      props.note.pageNumber,
      {
        uuid,
      }
    );
  }
};

const handleDelete = () => {
  Modal.confirm({
    title: t('teams.deleteTip'),
    onOk: async () => {
      await annotationStore.controller.deleteAnnotation(
        props.note.uuid,
        props.note.pageNumber
      );
      annotationStore.currentAnnotationId = '';
      clearArrow();
    },
    okButtonProps: {
      danger: true,
    },
    cancelButtonProps: { type: 'primary' },
    cancelText: t('viewer.cancel'),
    okText: t('viewer.delete'),
  });
};

const queryTargetNoteId = new URL(window.location.href).searchParams.get(
  'targetNoteId'
);

const queryTargetCommentId = new URL(window.location.href).searchParams.get(
  'targetCommentId'
);

const handlePositioning = (comment: Element, callback: Function) => {
  setTimeout(() => {
    scrollIntoView(comment, {
      block: 'center',
      inline: 'center',
    });

    scrollToAnnotation();

    callback();
  }, 1000);
};

const handleElementStyle = (
  parentEl: HTMLDivElement,
  targetEl: HTMLDivElement
) => {
  targetEl.style.position = 'absolute';
  targetEl.style.top = '0';
  targetEl.style.height = parentEl.offsetHeight + 'px';
  targetEl.style.backgroundColor = 'rgb(255, 255, 0, 15%)';
  targetEl.classList.add('bbox-animation');
  parentEl.appendChild(targetEl);
};

onMounted(() => {
  if (queryTargetNoteId === props.note.uuid && !queryTargetCommentId) {
    const comment = document.querySelector(
      `.note${attrSelector(PDF_ANNOTATE_ID, queryTargetNoteId)}`
    ) as HTMLDivElement;

    if (comment) {
      setTimeout(() => {
        const div = document.createElement('div');

        div.style.width = comment.offsetWidth + 'px';

        handleElementStyle(comment, div);

        handlePositioning(comment, () => {
          setTimeout(() => {
            comment.removeChild(div);
          }, 1000);
        });
      }, 300);
    }
  }
});

const handlePositioningComment = () => {
  const comment = document.querySelector(
    `.group-comment-container[data-group-comment-id="${queryTargetCommentId}"]`
  ) as HTMLDivElement;

  if (comment) {
    const div = document.createElement('div');

    div.style.width = comment.offsetWidth + 29 + 'px';

    div.style.margin = '0 -12px 0 -17px';

    div.style.padding = '0 12px 0 17px';

    handleElementStyle(comment, div);

    handlePositioning(comment, () => {
      setTimeout(() => {
        comment.removeChild(div);
      }, 2000);
    });
  }
};
</script>

<style lang="less" scoped>
.note {
  background: #393c3e;
  border-radius: 2px;
  margin-bottom: 10px;
  display: flex;
  cursor: pointer;
  position: relative;

  .left {
    flex: 0 0 4px;
    background-color: #65676c;
  }

  .right {
    padding: 16px 12px 16px 17px;
    flex: 1 1 auto;

    .ref-content {
      background: rgba(255, 255, 255, 0.05);
      mix-blend-mode: normal;

      font-family: 'Lato';
      font-style: italic;
      font-weight: 400;
      font-size: 13px;
      line-height: 20px;

      color: #a8adb3;
      margin-bottom: 8px;
      word-break: break-word;
    }

    .ref-picture {
      display: flex;
      justify-content: center;
      margin-bottom: 8px;

      img {
        opacity: 0.4;
        // max-height: 360px;
        max-width: 100%;
      }
    }

    .note-content {
      cursor: pointer;
      display: flex;

      .content-container {
        font-family:
          PingFangSC-Regular,
          PingFang SC;
        font-weight: 400;

        .name {
          font-size: 13px;
          color: rgba(255, 255, 255, 45%);
          line-height: 20px;
        }

        .content {
          font-size: 14px;
          &,
          * {
            word-break: break-all;
          }
          color: #ffffff;
          line-height: 22px;
        }
      }

      .avatar {
        flex: 0 0 24px;
        margin-right: 8px;
      }
    }

    .current-note-content {
      cursor: text;
    }
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

    .iconguanbi {
      font-size: 14px;
      color: rgba(38, 38, 38, 1);
      transform: scale(0.5);
    }
  }

  &:hover {
    .close {
      display: flex;
    }
  }
}
</style>
