<template>
  <a-spin v-if="loading && !list.length" />
  <div
    v-else-if="list.length"
    class="task-list"
  >
    <div class="mb-4">
      <slot name="upload" />
      <a-button
        type="link"
        :loading="loading"
        @click="refresh"
      >
        <template
          v-if="!loading"
          #icon
        >
          <ReloadOutlined />
        </template>{{ $t('common.text.refresh') }}
      </a-button>
    </div>
    <a-table
      :columns="columns"
      :data-source="list"
      :pagination="false"
      :loading="loading"
    >
      <template
        #bodyCell="{ column, record }: { column: any; record: TaskInfo }"
      >
        <template v-if="column.key === 'id'">
          <span>{{ record.id }}</span>
        </template>
        <template v-else-if="column.key === 'createDate'">
          <span>{{
            dayjs(+record.createDate).format('YYYY-MM-DD HH:mm:ss')
          }}</span>
        </template>
        <template v-else-if="column.key === 'fileName'">
          <a
            class="text-rp-blue-7 hover:text-rp-blue-6"
            :href="record.originalFileUrl"
            target="_blank"
          >{{ record.fileName }}</a>
        </template>
        <template v-else-if="column.key === 'major'">
          <span>{{ record.major }}</span>
        </template>
        <template v-else-if="column.key === 'paperType'">
          <span>{{
            !kind2Key[record.paperType]
              ? 'Unknown'
              : $t(`common.aitools.paperKind.${kind2Key[record.paperType]}`)
          }}</span>
        </template>
        <template v-else-if="column.key === 'langType'">
          <span>{{ lang2Txt[record.langType] }}</span>
        </template>
        <template v-else-if="column.key === 'taskStatus'">
          <span>{{ status2Txt[record.taskStatus] }}</span><template
            v-if="
              [TaskStatus.CONSUME_BEAN_FAIL, TaskStatus.FAIL].includes(
                record.taskStatus
              )
            "
          >
            {{ $t('common.symbol.comma')
            }}{{ $t('common.aiReviewerQS.failTip') }}
          </template>
          <Help
            v-if="
              [TaskStatus.CONSUME_BEAN_FAIL, TaskStatus.FAIL].includes(
                record.taskStatus
              )
            "
          >
            <span class="text-rp-gold-7 ml-1"><QuestionCircleOutlined /></span>
          </Help>
        </template>
        <template v-else-if="column.key === 'resultFileUrl'">
          <span v-if="!record.resultFileUrl">暂无</span>
          <a
            v-else
            class="text-rp-blue-7 hover:text-rp-blue-6"
            :href="record.resultFileUrl"
            target="_blank"
            @click="onDownload(record)"
          >{{ $t('common.text.download') }}</a>
        </template>
      </template>
      <!-- <template #footer></template> -->
    </a-table>
    <p class="my-4 text-rp-neutral-6 text-sm">
      <span class="mr-2"><InfoCircleOutlined /></span>{{ $t('common.aiReviewerQS.tasksTip') }}
    </p>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
import {
  ReloadOutlined,
  InfoCircleOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons-vue';
import { useRequest } from 'ahooks-vue';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import dayjs from 'dayjs';
import { getQSReviewTaskList } from '@common/api/review';
import {
  PaperType,
  TaskInfo,
  TaskStatus,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/review/AiReviewPaper';
import Help from '@common/components/Help.vue';
import { useAIReviewQSWordings } from './consts';
import {
  PageType,
  ModuleType,
  EventCode,
  reportEvent,
} from '../../utils/report';

const { t } = useI18n();

const columns = computed(() => [
  {
    title: t('common.aiReviewerQS.labelTask'),
    key: 'id',
  },
  {
    title: t('common.aiReviewerQS.createTime'),
    key: 'createDate',
  },
  {
    title: t('common.aiReviewerQS.labelFile'),
    key: 'fileName',
  },
  {
    title: t('common.aiReviewerQS.labelMajor'),
    key: 'major',
  },
  {
    title: t('common.aiReviewerQS.labelKind'),
    key: 'paperType',
  },
  {
    title: t('common.aiReviewerQS.labelLang'),
    key: 'langType',
  },
  {
    title: t('common.aiReviewerQS.labelStatus'),
    key: 'taskStatus',
  },
  {
    title: t('common.aiReviewerQS.labelResult'),
    key: 'resultFileUrl',
  },
]);

const kind2Key = {
  [PaperType.FULL_TIME_EDUCATION_UNDERGRADUATE_THESIS]: 'fullTimeUndergraduate',
  [PaperType.GRADUATE_THESIS]: 'fullTimeMaster',
  [PaperType.DOCTORAL_THESIS]: 'fullTimeDoctor',
  [PaperType.ADULT_EDUCATION_UNDERGRADUATE_THESIS]: 'partTime',
  [PaperType.UNRECOGNIZED]: '',
};

const { lang2Txt, status2Txt } = useAIReviewQSWordings();

const { data, loading, refresh } = useRequest(async () => {
  const res = await getQSReviewTaskList();

  return res;
}, {});
const list = computed(
  () =>
    data.value?.taskList ||
    [
      // {
      //   id: '1',
      //   createDate: '2024-01-01 10:00:00',
      //   fileName: '研究.doc',
      //   major: '软件工程',
      //   paperType: PaperType.FULL_TIME_EDUCATION_UNDERGRADUATE_THESIS,
      //   langType: LangType.ZH,
      //   taskStatus: TaskStatus.SUCCESS,
      //   resultFileUrl: '',
      //   originalFileUrl: '',
      // },
    ]
);

const onDownload = (x: TaskInfo) => {
  reportEvent(EventCode.readpaper_ai_polish_result_feedback_click, {
    element_name: 'down_load',
    polish_type: ModuleType.AI_MENTOR_GRA,
    task_id: x.id,
    page_type: PageType.POLISH,
  });
};

defineExpose({
  refresh,
});
</script>

<style lang="less" scoped>
.task-list {
  :deep(*) {
    .ant-table {
      .ant-table-footer {
        background: none;
      }
    }
  }
}
</style>
