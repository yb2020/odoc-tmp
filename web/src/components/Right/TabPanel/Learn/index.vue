<template>
  <PerfectScrollbar class="learn-wrap">
    <ErrorVue
      v-if="!isOpenPaper"
      :style="{ height: '100%' }"
      :img-url="noPaperImgUrl"
      :message="$t('tasks.errorMsg')"
    />

    <div
      v-if="isOpenPaper"
      :class="{ center: !qaInfo }"
    >
      <a-spin :spinning="fetchState.pending">
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

        <div v-else-if="qaInfo && qaInfo.length > 0">
          <div class="info">
            <a-tooltip
              :title="$t('tasks.tip')"
              placement="bottomLeft"
              :arrow-point-at-center="true"
              :overlayStyle="{ 'max-width': '360px' }"
            >
              <info-circle-outlined style="margin-right: 5px" />
            </a-tooltip>
            {{ $t('tasks.tipLabel') }}
          </div>

          <div class="firstStage">
            <div class="header">
              {{ $t('tasks.label1') }}
            </div>
            <div class="content">
              <div class="title flex">
                <span
                  class="abstract"
                  @click="goToAbstractPage"
                >{{
                  $t('tasks.abstract')
                }}</span>
                <edit-outlined
                  v-show="
                    translationInfo &&
                      translationInfo.translation &&
                      !isShowInput
                  "
                  class="iconedit"
                  @click="handleClick"
                />
              </div>
              <div v-show="!isShowInput">
                <div
                  v-if="translationInfo && translationInfo.translation"
                  class="translation-wrap"
                >
                  <p
                    ref="translationRef"
                    :class="{
                      'part-translation': showMoreTranslation === 'close',
                      'all-translation': showMoreTranslation === 'show',
                    }"
                    v-html="htmlTranslation"
                  />
                  <span
                    v-if="showMoreTranslation === 'close'"
                    class="btn"
                    @click="showMoreTranslation = 'show'"
                  >{{ $t('viewer.expand') }}</span>
                  <span
                    v-if="showMoreTranslation === 'show'"
                    class="btn"
                    @click="showMoreTranslation = 'close'"
                  >{{ $t('viewer.collapse') }}</span>
                </div>

                <div
                  v-else
                  class="add-translate-btn-wrap"
                >
                  <a-button
                    class="btn"
                    @click="handleClick"
                  >
                    <plus-outlined />{{ $t('tasks.abstractBtn') }}
                  </a-button>
                </div>
              </div>
              <div
                v-show="isShowInput"
                style="margin-top: 8px"
              >
                <a-textarea
                  ref="inputRef"
                  v-model:value="input"
                  :placeholder="$t('tasks.inputPlaceholder')"
                  :auto-size="{ minRows: 3, maxRows: 50 }"
                  class="input-textarea"
                  @blur="handleTranslate"
                />
              </div>
            </div>
          </div>

          <div class="secondStage">
            <div class="header flex">
              <span>{{ $t('tasks.label2') }}</span>
              <span class="answer-num">{{ $t('tasks.progress') }}{{ answerCount }}/10</span>
            </div>
            <div class="content">
              <Question
                v-for="(item, index) in showQaInfo"
                :key="item.questionId"
                :data="item"
                :index="index"
                :paper-id="paperId"
                :width="sideTabSettings.width"
                @success="handleAnswerSuccess"
              />

              <div
                v-if="!showMoreQuestion"
                class="btn"
                @click="showMoreQuestion = true"
              >
                {{ $t('tasks.expandTenQuestions')
                }}<down-outlined style="margin-left: 4px" />
              </div>

              <div
                v-else
                class="btn"
                @click="showMoreQuestion = false"
              >
                {{ $t('tasks.collapseQuestions')
                }}<up-outlined style="margin-left: 4px" />
              </div>
            </div>
          </div>

          <Comment :paperId="paperId" />
        </div>
      </a-spin>
    </div>
  </PerfectScrollbar>
</template>
<script lang="ts" setup>
import { computed, ref, nextTick, onMounted, watch } from 'vue';
import {
  InfoCircleOutlined,
  PlusOutlined,
  EditOutlined,
  DownOutlined,
  UpOutlined,
} from '@ant-design/icons-vue';
import useFetch from '~/src/hooks/useFetch';
import ErrorVue from '../QuestionAnswer/Error.vue';
import {
  FastReadingQuestionInfo,
  BaseTranslationInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/PaperDetail';
import {
  getMyLearningTaskInfo,
  savePaperTranslation,
  updatePaperTranslation,
} from '@/api/learn';
import Question from './Question.vue';
import lodash from 'lodash-es';
import { message } from 'ant-design-vue';
import noPaperImgUrl from '@/assets/images/no_paper.png';
import { goPathPage } from '~/src/common/src/utils/url';
import { getDomainOrigin } from '~/src/util/env';
import Comment from './Comment.vue';
import { useRightSideTabSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { checkOpenPaper } from '~/src/api/helper';
import { useI18n } from 'vue-i18n';

const props = defineProps<{ paperId: string; isPrivatePaper: boolean }>();

const qaInfo = ref<FastReadingQuestionInfo[]>();

const translationInfo = ref<BaseTranslationInfo>();

const { t } = useI18n();

const { fetch, fetchState } = useFetch(async () => {
  const res = await getMyLearningTaskInfo({ paperId: props.paperId });

  qaInfo.value = res?.qaInfo || [];

  translationInfo.value = res?.translationInfo;
});

const inputRef = ref<HTMLElement>();

const input = ref<string>('');

const prevInput = computed(
  () => (translationInfo.value && translationInfo.value.translation) || ''
);

const isShowInput = ref<boolean>(false);

nextTick(() => {
  inputRef.value?.focus();
});

const isOpenPaper = computed(() =>
  checkOpenPaper(props.paperId, props.isPrivatePaper)
);

const htmlTranslation = computed(() => {
  return (
    translationInfo.value &&
    translationInfo.value.translation &&
    lodash.escape(translationInfo.value.translation).replace(/\n/g, '<br/>')
  );
});

const handleClick = () => {
  isShowInput.value = true;

  nextTick(() => {
    inputRef.value?.focus();
  });

  if (translationInfo.value && translationInfo.value.translationId) {
    input.value = translationInfo.value.translation;
  }
};

const handleTranslate = async () => {
  if (prevInput.value === input.value.trim()) {
    isShowInput.value = false;
    return;
  }

  if (input.value.trim().length > 0 && input.value.trim().length < 5) {
    message.info(t('tasks.inputLimitTip'));
    return;
  }

  const translationId =
    translationInfo.value && translationInfo.value.translationId;

  const request = translationId ? updatePaperTranslation : savePaperTranslation;

  await request({
    paperId: props.paperId,
    id: translationId || '',
    translation: input.value.trim(),
  });

  message.success(
    translationId
      ? t('tasks.modifyTranslateTip')
      : t('tasks.translateSuccessTip')
  );

  await fetch();

  isShowInput.value = false;

  showMoreTranslation.value = '';

  handleIsShowMore();
};

const answerCount = computed(() =>
  qaInfo.value?.reduce((a: any, c: any) => {
    if (c.answer) {
      return a + 1;
    } else {
      return a;
    }
  }, 0)
);

const handleAnswerSuccess = async (callback: Function) => {
  await fetch();

  callback();
};

const goToAbstractPage = () => {
  goPathPage(`${getDomainOrigin()}/paper/${props.paperId}/abstract`);
};

const translationRef = ref();

const showMoreTranslation = ref<string>('');

const handleIsShowMore = () => {
  setTimeout(() => {
    if (translationRef.value) {
      const dom: any = translationRef.value;

      if (dom.offsetHeight > 240) {
        showMoreTranslation.value = 'close';
      }
    }
  }, 200);
};

onMounted(() => {
  handleIsShowMore();
});

const showMoreQuestion = ref<boolean>(false);

const showQaInfo = computed(() => {
  if (!qaInfo.value) return [];

  return showMoreQuestion.value ? qaInfo.value : qaInfo.value.slice(0, 5);
});

const { sideTabSettings } = useRightSideTabSettings();

watch(
  () => sideTabSettings.value.width,
  () => {
    showMoreTranslation.value = '';

    handleIsShowMore();
  }
);
</script>

<style lang="less" scoped>
.learn-wrap {
  height: 100%;
  padding: 0 16px;
  .info {
    font-weight: 400;
    font-size: 12px;
    line-height: 16px;
    color: rgba(255, 255, 255, 0.65);
    margin-bottom: 12px;
  }
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
  :deep(textarea)::placeholder {
    color: #c9cdd4;
  }
  .firstStage {
    margin-bottom: 16px;
    .content {
      padding: 8px 16px;
      background: #484a4d;
      border-radius: 2px;
      .title {
        font-weight: 400;
        font-size: 14px;
        line-height: 22px;
        color: rgba(255, 255, 255, 0.85);
        margin-top: 4px;
        .iconedit {
          cursor: pointer;
        }
        .abstract {
          cursor: pointer;
          &:hover {
            text-decoration: underline;
          }
        }
      }
      .add-translate-btn-wrap {
        display: flex;
        justify-content: center;
        margin: 16px 0 8px;
        .btn {
          color: rgba(255, 255, 255, 0.85);
          opacity: 0.85;
          border: 1px solid rgba(255, 255, 255, 0.3);
          border-radius: 2px;
          width: 216px;
        }
      }
      .translation-wrap {
        margin-top: 8px;
        word-break: break-word;
        p {
          margin-bottom: 0;
        }
        .part-translation {
          display: -webkit-box;
          -webkit-line-clamp: 10;
          -webkit-box-orient: vertical;
          overflow: hidden;
          text-overflow: ellipsis;
          line-height: 24px;
        }
        .all-translation {
          -webkit-line-clamp: 1000;
        }
        .btn {
          font-size: 13px;
          font-weight: 400;
          color: #1f71e0;
          cursor: pointer;
          display: inline-block;
          margin-top: 4px;
        }
      }
    }
  }
  .secondStage {
    .answer-num {
      font-weight: 400;
      font-size: 12px;
      line-height: 18px;
      color: rgba(255, 255, 255, 0.65);
    }
    .content {
      background: #484a4d;
      border-radius: 2px;
      padding: 8px 0;
      .btn {
        padding: 10px 0 4px 16px;
        font-weight: 400;
        font-size: 13px;
        line-height: 20px;
        color: rgba(255, 255, 255, 0.65);
        cursor: pointer;
      }
    }
  }
  .center {
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .input-textarea {
    color: #4e5969;
    background: #fff;
  }
}
</style>
