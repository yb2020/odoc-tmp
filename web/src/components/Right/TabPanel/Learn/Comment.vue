<template>
  <div class="threeStage">
    <div class="header flex">
      <span>{{ $t('tasks.label3') }}</span>
      <span class="answer-num">{{ $t('tasks.progress') }}{{ isCommented ? 1 : 0 }}/1</span>
    </div>
    <div class="content">
      <a-spin
        :spinning="fetchState.pending"
        wrapperClassName="loading"
      >
        <ErrorVue
          v-if="fetchState.error"
          :message="fetchState.error.message"
        >
          <a-button
            type="text"
            @click="fetch"
          >
            {{
              $t('tasks.retry')
            }}
          </a-button>
        </ErrorVue>

        <div v-else-if="!fetchState.pending">
          <div
            v-if="!isCommented"
            class="empty"
          >
            <a-radio-group v-model:value="myCommentData.paperCommentLevel">
              <a-radio
                v-for="option in COMMENT_OPTIONS"
                :key="option.value"
                :value="option.value"
                :style="radioStyle"
                @click="showModal(option.value)"
              >
                {{ $t(option.i18n) }}
              </a-radio>
            </a-radio-group>
          </div>

          <div
            v-else
            class="comment"
          >
            <div class="level">
              <a-radio :checked="true">
                {{ myCommentLabel }}
              </a-radio>
              <edit-outlined
                class="iconedit"
                @click="showModal(myCommentData.paperCommentLevel)"
              />
            </div>

            <div
              v-if="dimensionScore && dimensionScore.length > 0"
              class="score"
            >
              <div
                v-for="dimension in dimensionScore"
                :key="dimension.paperCommentDimension"
                class="item"
              >
                {{ getDimension(dimension.paperCommentDimension) }}：{{
                  dimension.score
                }}
              </div>
            </div>

            <MarkdownViewer
              v-if="myCommentData.commentContent"
              :raw="myCommentData.commentContent"
              class="comment-content"
            />

            <div
              v-if="myCommentData.anonymous"
              class="anonymous"
            >
              {{ $t('tasks.anonymous') }}
            </div>
          </div>

          <a-modal
            v-model:visible="isShowModal"
            :width="640"
            :centered="true"
            :title="$t('tasks.evaluation')"
            footer=""
            class="comments-modal-wrap"
            :bodyStyle="{ height: iframeHeight + 'px' }"
            @cancel="hideModal"
          >
            <iframe
              ref="commentsRef"
              :src="`${getDomainOrigin()}/comment`"
              title="comments"
              width="640"
              :onload="handleLoad"
              class="comments-iframe"
              :height="iframeHeight"
            />
          </a-modal>
        </div>
      </a-spin>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted } from 'vue';
import useFetch from '~/src/hooks/useFetch';
import { PaperCommentLevel } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/PaperCommentView';
import {
  getCommentAggregation,
  updatePaperComment,
  createPaperComment,
} from '@/api/learn';
import { getDomainOrigin } from '~/src/util/env';
import { EditOutlined } from '@ant-design/icons-vue';
import {
  COMMENT_OPTIONS,
  DEFAULT_COMMENT,
  DIMENSION_NAME_MAP,
} from './constants';
import { MarkdownViewer } from '@idea/aiknowledge-markdown';
import ErrorVue from '../QuestionAnswer/Error.vue';
import { useI18n } from 'vue-i18n';
interface CommentDimensionType {
  0: string;
  1: string;
  2: string;
  3: string;
}

const { t } = useI18n();

const props = defineProps<{ paperId: string }>();

const myCommentData = ref<any>({
  paperCommentLevel: -1,
});

const isCommented = computed(() => myCommentData.value.paperCommentId);

const radioStyle = ref({
  display: 'flex',
  fontFamily: 'Noto Sans SC',
  fontWeight: '400',
  fontSize: '14px',
  lineHeight: '22px',
  color: 'rgba(255, 255, 255, 0.65)',
  marginBottom: '12px',
});

const { fetch, fetchState } = useFetch(async () => {
  const { myPaperComment } = await getCommentAggregation({
    paperId: props.paperId,
  });

  myCommentData.value = myPaperComment || { ...DEFAULT_COMMENT };

  const commentAnonymous = localStorage.getItem('commentAnonymous');

  if (!myCommentData.value.paperCommentId) {
    myCommentData.value.anonymous = commentAnonymous !== '0';
  }
});

const isShowModal = ref<boolean>(false);

const showModal = (key: any) => {
  isShowModal.value = true;

  myCommentData.value.paperCommentLevel = key;

  if (commentsRef.value) {
    handleLoad();
  }
};

const commentsRef = ref();

const handleLoad = () => {
  commentsRef.value.contentWindow.postMessage(
    {
      event: 'openCommentsPage',
      params: {
        isShowModal: isShowModal.value,
        myCommentData: JSON.stringify(myCommentData.value),
      },
    },
    '*'
  );
};

const updateComment = async () => {
  const {
    paperCommentLevel,
    commentContent,
    anonymous,
    paperCommentDimensionScore,
    paperCommentId,
  } = myCommentData.value;

  if (paperCommentId) {
    await updatePaperComment({
      paperCommentLevel,
      paperCommentContent: commentContent,
      anonymous,
      paperCommentDimensions: paperCommentDimensionScore,
      paperId: props.paperId,
      paperCommentId,
    });
  } else {
    const { paperCommentId } = await createPaperComment({
      paperCommentLevel,
      paperCommentContent: commentContent,
      anonymous,
      paperCommentDimensions: paperCommentDimensionScore,
      paperId: props.paperId,
    });

    myCommentData.value.paperCommentId = paperCommentId;
  }

  fetch();
};

const hideModal = () => {
  isShowModal.value = false;
  if (!isCommented.value) {
    myCommentData.value.paperCommentLevel = -1;
  }
};

const iframeHeight = ref<number>(620);

onMounted(() => {
  window.addEventListener(
    'message',
    async (event) => {
      const data = event.data;

      if (data.event === 'setMyCommentData') {
        myCommentData.value = JSON.parse(data.params.myCommentData);

        if (data.params.isRequest) {
          await updateComment();
        }

        hideModal();

        return;
      }

      if (data.event === 'hideModal') {
        hideModal();

        return;
      }

      if (data.event === 'handleIsHasCommentTips') {
        iframeHeight.value = data.params.IsHasCommentTips ? 676 : 620;

        return;
      }
    },
    false
  );
});

const myCommentLabel = computed(() => {
  const findIndex = COMMENT_OPTIONS.findIndex(
    (item) => item.value === myCommentData.value.paperCommentLevel
  );

  return t(COMMENT_OPTIONS[findIndex]?.i18n || '');
});

const dimensionScore = computed(() => {
  return myCommentData.value?.paperCommentDimensionScore?.filter(
    (item: any) => item.score > 0
  );
});

const getDimension = (dimension: keyof CommentDimensionType) => {
  return t(DIMENSION_NAME_MAP[dimension]);
};
</script>

<style lang="less" scoped>
.threeStage {
  margin: 16px 0 20px;
  .header {
    font-weight: 500;
    font-size: 14px;
    line-height: 22px;
    color: rgba(255, 255, 255, 0.85);
    margin-bottom: 9px;
  }
  .flex {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .answer-num {
    font-weight: 400;
    font-size: 12px;
    line-height: 18px;
    color: rgba(255, 255, 255, 0.65);
  }
  .content {
    background: #484a4d;
    border-radius: 2px;
    padding: 12px;
    :deep(.ant-radio-inner) {
      border: 1px solid rgba(255, 255, 255, 0.25);
    }
    display: flex;
    flex-direction: column-reverse;
    align-items: center;
    .loading {
      width: 100%;
    }
    .empty {
      :deep(.ant-radio-input):focus + .ant-radio-inner {
        box-shadow: none;
      }
    }
    .comment {
      .level {
        display: flex;
        justify-content: space-between;
        align-items: center;
        .iconedit {
          cursor: pointer;
        }
      }
      .score {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        margin-top: 8px;
        .item {
          font-size: 13px;
          line-height: 20px;
          color: rgba(255, 255, 255, 0.65);
          padding: 0 8px;
          border: 1px solid rgba(255, 255, 255, 0.25);
          border-radius: 2px;
          margin: 0 4px 4px 0;
        }
      }
      .comment-content {
        font-size: 14px;
        line-height: 24px;
        color: rgba(255, 255, 255, 0.85);
        margin-top: 8px;
        padding-bottom: 0;
      }
      .anonymous {
        font-size: 13px;
        line-height: 20px;
        color: rgba(255, 255, 255, 0.65);
        margin-top: 8px;
      }
    }
  }
}
.comments-iframe {
  border-width: 0;
}
</style>

<style lang="less">
.comments-modal-wrap {
  .ant-modal-content {
    .ant-modal-body {
      padding: 0;
    }
  }
}
</style>
