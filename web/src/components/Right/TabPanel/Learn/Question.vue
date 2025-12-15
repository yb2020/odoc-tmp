<template>
  <div class="questionItem">
    <div class="title">
      <span class="index">Q{{ index + 1 }}</span>
      <a-tooltip
        v-if="data.answerCount > 0"
        placement="top"
        :arrow-point-at-center="true"
      >
        <template #title>
          {{
            $t('tasks.answersTotal', { count: data.answerCount })
          }}
        </template>
        <span
          class="detail-title"
          @click="goToQuestionsDetailPage"
        >{{
          questionTitle
        }}</span>
      </a-tooltip>
      <span
        v-else
        class="detail-title"
        @click="goToQuestionsDetailPage"
      >{{
        questionTitle
      }}</span>
      <edit-outlined
        v-show="!isShowInput"
        class="iconedit"
        @click="handleClick"
      />
    </div>

    <div
      v-if="data.answer || isShowInput"
      class="addAnswer"
    >
      <IdeaMarkdown
        ref="ideaMarkdownRef"
        :raw="data.answer || ''"
        :uniq-id="data.questionId"
        :editing="isShowInput"
        :upload="upload"
        :placeholder="$t('tasks.inputPlaceholder')"
        :minRows="2"
        :maxRows="50"
        :allow-enter="true"
        :class="{
          'part-answer': showMoreAnswer === 'close',
          'all-answer': showMoreAnswer === 'show',
        }"
        @blur="handlePublish"
      />

      <span
        v-if="showMoreAnswer === 'close'"
        class="btn"
        @click="showMoreAnswer = 'show'"
      >{{ $t('viewer.expand') }}</span>
      <span
        v-if="showMoreAnswer === 'show'"
        class="btn"
        @click="showMoreAnswer = 'close'"
      >{{ $t('viewer.collapse') }}</span>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { ref, computed, onMounted, watch } from 'vue';
import { mdi } from '@idea/aiknowledge-markdown';
import { FastReadingQuestionInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/paper/PaperDetail';
import { message } from 'ant-design-vue';
import { saveAnswer, updateAnswer } from '@/api/learn';
import { PageType, reportTenQListClick } from '~/src/api/report';
import { currentNoteInfo } from '~/src/store';
import { EditOutlined } from '@ant-design/icons-vue';
import { uploadImage, ImageStorageType } from '~/src/api/upload';
import { IdeaMarkdown } from '@idea/aiknowledge-markdown';
import '@idea/aiknowledge-markdown/dist/style.css';
import { goPathPage } from '~/src/common/src/utils/url';
import { getDomainOrigin } from '~/src/util/env';
import { useI18n } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '~/src/hooks/useLanguage';

const props = defineProps<{
  data: FastReadingQuestionInfo;
  paperId: string;
  index: number;
  width: number;
}>();

const emit = defineEmits<{
  (event: 'success', option: any): void;
}>();

const { locale } = useI18n();
const { isCurrentLanguage } = useLanguage();

const questionTitle = computed(() => {
  const title = props.data.questionTitle.split('&&');

  return title[isCurrentLanguage(Language.EN_US) ? 0 : 1];
});

const isShowInput = ref<boolean>(false);

const prevInput = computed(() => (props.data && props.data.answer) || '');

const handleClick = () => {
  isShowInput.value = true;

  showMoreAnswer.value = '';
};

const handlePublish = async (input: string) => {
  if (prevInput.value === input.trim()) {
    isShowInput.value = false;

    handleIsShowMore();

    return;
  }

  if (input.trim().length > 0 && input.trim().length < 5) {
    message.info('请至少输入 5 个有效字符');
    return;
  }

  const request = props.data.answerId ? updateAnswer : saveAnswer;

  await request({
    paperId: props.paperId,
    id: props.data.answerId,
    classicQuestionId: props.data.questionId,
    answer: input.trim(),
    htmlAnswer: mdi.render(input.trim()),
  });

  reportTenQListClick({
    type_parameter: currentNoteInfo.value?.pdfId || '',
    question_id: props.data.questionId,
    page_type: PageType.note,
  });

  message.success(props.data.answerId ? '修改回答成功' : '回答成功');

  await emit('success', () => {
    (isShowInput.value = false), handleIsShowMore();
  });
};

const upload = async (src: File | string) => {
  return uploadImage(src, ImageStorageType.markdown);
};

const goToQuestionsDetailPage = () => {
  goPathPage(
    `${getDomainOrigin()}/paper/${props.paperId}/questions-detail?questionId=${
      props.data.questionId
    }`
  );
};

const ideaMarkdownRef = ref();

const showMoreAnswer = ref<string>('');

const handleIsShowMore = () => {
  setTimeout(() => {
    if (ideaMarkdownRef.value) {
      const dom: any = ideaMarkdownRef.value.$el;

      if (dom.offsetHeight - 6 > 66) {
        showMoreAnswer.value = 'close';
      }
    }
  }, 200);
};

onMounted(() => {
  handleIsShowMore();
});

watch(
  () => props.width,
  () => {
    showMoreAnswer.value = '';

    handleIsShowMore();
  }
);
</script>
<style lang="less" scoped>
.questionItem {
  padding: 8px 16px;
  font-weight: 400;
  font-size: 14px;
  line-height: 22px;
  color: rgba(255, 255, 255, 0.65);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  &:last-child {
    border-bottom: none;
  }
  .title {
    .index {
      margin-right: 5px;
    }
    .iconedit {
      display: none;
      cursor: pointer;
      margin-left: 5px;
    }
    .detail-title {
      cursor: pointer;
      &:hover {
        text-decoration: underline;
      }
    }
  }
  &:hover {
    color: #fff;
    background: rgba(255, 255, 255, 0.08);
    .iconedit {
      display: inline-block;
    }
  }

  .addAnswer {
    margin-top: 11px;
    .part-answer {
      display: -webkit-box;
      -webkit-line-clamp: 3;
      -webkit-box-orient: vertical;
      overflow: hidden;
      text-overflow: ellipsis;
      line-height: 22px;
      padding-bottom: 0;
    }
    .all-answer {
      -webkit-line-clamp: 1000;
      padding-bottom: 0;
    }
    .btn {
      font-size: 13px;
      font-weight: 400;
      color: #1f71e0;
      cursor: pointer;
      display: inline-block;
      margin-top: 6px;
    }
  }
}
</style>
