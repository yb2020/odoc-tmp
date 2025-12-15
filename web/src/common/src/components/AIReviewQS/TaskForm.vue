<template>
  <a-modal
    v-model:visible="visible"
    :width="600"
    :title="$t('common.aiReviewerQS.createTt')"
    :footer="null"
    @cancel="onCancel"
  >
    <p>{{ $t('common.aiReviewerQS.createTip') }}</p>
    <a-form
      ref="form"
      layout="vertical"
      :model="formState"
    >
      <a-form-item
        required
        name="fileName"
        :label="`${$t('common.text.paper')}(${$t(
          'common.aiReviewerQS.labelTipFile'
        )})`"
        :rules="[
          {
            validator: fileValidator,
          },
        ]"
      >
        <a-upload
          :maxCount="1"
          :multiple="false"
          :showUploadList="false"
          :customRequest="onUpload"
        >
          <div class="flex items-center gap-2">
            <a-button shape="round">
              <UploadOutlined />
              <template v-if="formState.fileName">
                {{
                  $t('common.text.re')
                }}
              </template>{{ $t('common.text.upload') }}
            </a-button>
            <LoadingOutlined v-if="uploading" />
            <span
              v-else-if="formState.fileName"
              class="text-rp-neutral-8"
            >{{
              formState.fileName
            }}</span>
          </div>
        </a-upload>
      </a-form-item>
      <a-form-item
        required
        name="major"
        :label="$t('common.aiReviewerQS.labelMajor')"
        :rules="[{ required: true }, { max: 20 }]"
      >
        <a-input v-model:value="formState.major" />
      </a-form-item>
      <a-form-item
        required
        name="kind"
        :label="$t('common.aiReviewerQS.labelKind')"
      >
        <a-select v-model:value="formState.kind">
          <a-select-option
            v-for="item in OptionsKind"
            :key="item.value"
            :value="item.value"
            :disabled="item.disabled"
          >
            {{ $t(`common.aitools.paperKind.${item.label}`) }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item
        required
        name="lang"
        :label="$t('common.aiReviewerQS.labelLang')"
      >
        <a-select v-model:value="formState.lang">
          <a-select-option :value="LangType.ZH">
            {{
              $t('common.text.chinese')
            }}
          </a-select-option>
          <a-select-option
            :value="LangType.EN"
            disabled
          >
            {{
              $t('common.text.english')
            }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item
        name="agreed"
        :rules="[
          {
            validator: agreedValidator,
            trigger: 'change',
          },
        ]"
      >
        <p class="mb-0">
          <span class="mr-1"><a-checkbox
            v-model:checked="formState.agreed"
            @change="form.validate(['agreed'])"
          /></span>
          {{ $t('common.aiReviewerQS.labelRules') }}
        </p>
        <ul class="mb-0 text-rp-neutral-6">
          <li v-for="item in $t('common.aiReviewerQS.createRules').split('\n')">
            {{ item }}
          </li>
        </ul>
      </a-form-item>

      <div class="text-center">
        <p class="mb-2 flex items-center justify-center gap-2">
          <span class="text-3xl font-medium">{{ data?.currentAiBeanPrice || 500 }}
            {{ $t('common.aibeans.name') }}/{{
              $t('common.premium.units.piece')
            }}</span>
          <span class="flex flex-col">
            <span
              class="bg-rp-red-6 text-xs tracking-tighter text-white py-1 px-2 rounded-tl-lg rounded-br-lg"
            >内测优惠</span>
            <span class="line-through text-xs text-rp-neutral-6">{{ data?.originalAiBeanPrice || 2880 }}
              {{ $t('common.aibeans.name') }}</span>
          </span>
        </p>
        <a-button
          type="primary"
          size="large"
          shape="round"
          :loading="isSaving"
          @click="onSaveTask"
        >
          {{ $t('common.aiReviewerQS.createBtn') }}
        </a-button>
        <p class="text-sm text-rp-neutral-8 mt-2 mb-0">
          {{
            $t('common.aiReviewerQS.beansTip', [
              data?.aiBeanAmountOneYuan || 10,
            ])
          }}
        </p>
      </div>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { LoadingOutlined, UploadOutlined } from '@ant-design/icons-vue';
import {
  GetConfigInfoResponse,
  PaperType as PaperKinds,
  LangType,
  SaveTaskResponse,
  TaskStatus,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/review/AiReviewPaper';
import { computed, ref, watch } from 'vue';
import {
  uploadQSReviewFile,
  saveQSReviewTask,
  cancelQSReviewTask,
  getQSReviewTask,
} from '@common/api/review';
import { useRequest } from 'ahooks-vue';
// import { useEventListener } from '@vueuse/core'
import {
  ERROR_CODE_BEANS_CASH_DEDUCTION,
  REQUEST_APPID,
} from '../../api/const';
import { BeanScenes, useAIBeans } from '../../hooks/useAIBeans';
import { useI18n } from 'vue-i18n';
import { ResponseError } from '../../api/type';
import { message } from 'ant-design-vue';
import { isDev } from '../../utils/env';
import { identity, pickBy } from 'lodash';

const DEFAULT_STATE = {
  fileName: '',
  taskId: '',
  major: '',
  kind: PaperKinds.FULL_TIME_EDUCATION_UNDERGRADUATE_THESIS,
  lang: LangType.ZH,
  agreed: false,
};

const props = defineProps<{
  data?: GetConfigInfoResponse;
}>();

const emit = defineEmits<{
  (e: 'created', taskId: string): void;
}>();

const { t, locale } = useI18n();

const visible = defineModel('visible', {
  type: Boolean,
  default: false,
});

const OptionsKind = [
  {
    label: 'fullTimeUndergraduate',
    value: PaperKinds.FULL_TIME_EDUCATION_UNDERGRADUATE_THESIS,
  },
  {
    label: 'partTime',
    value: PaperKinds.ADULT_EDUCATION_UNDERGRADUATE_THESIS,
  },
  {
    label: 'fullTimeMaster',
    value: PaperKinds.GRADUATE_THESIS,
    disabled: true,
  },
  {
    label: 'fullTimeDoctor',
    value: PaperKinds.DOCTORAL_THESIS,
    disabled: true,
  },
];

const form = ref();
const formState = ref({
  ...DEFAULT_STATE,
});
const isSaving = ref(false);
const saved = ref(false);
const fileValidator = async (_: unknown, v: string) => {
  const msg =
    error.value?.message ??
    (!formState.value.taskId && formState.value.fileName
      ? `${t('common.text.please')}${t('common.symbol.space')}${t(
          'common.text.upload'
        )}`
      : '');

  return msg ? Promise.reject(msg) : Promise.resolve();
};
const agreedValidator = async (_: unknown, v: boolean) => {
  return v
    ? Promise.resolve()
    : Promise.reject(t('common.aiReviewerQS.errorRules'));
};
const msgHandler = ref();

watch(msgHandler, (_, prevHandler) => {
  if (prevHandler) {
    window.removeEventListener('message', prevHandler);
  }
});

const {
  run: onUpload,
  loading: uploading,
  error,
} = useRequest(
  async ({ file, filename }: { file: File; filename: string }) => {
    const formdata = new FormData();
    formdata.append('appId', REQUEST_APPID);
    formdata.append('file', file);
    formdata.append('fileName', filename);
    formState.value.fileName = file.name;
    const res = await uploadQSReviewFile(formdata);

    if (!res.isSuccess) {
      throw new Error(res.msg);
    }

    if (!res.taskId) {
      throw new Error('Error: taskId is empty!');
    }

    formState.value = {
      ...formState.value,
      taskId: res.taskId,
    };

    return res;
  },
  {
    manual: true,
  }
);
watch([formState, error], () => {
  form.value.validate(['fileName']);
});

const { beans, consumeBeans } = useAIBeans();
const { run: doSaveTask } = useRequest(
  async () => {
    const { taskId, major, kind, lang } = formState.value;
    const res = await saveQSReviewTask({
      taskId,
      major,
      paperType: kind,
      langType: lang,
    })
      .then(() => true)
      .catch((e: unknown) => {
        const err = e as ResponseError;
        // 需要额外付款
        if (err?.code === ERROR_CODE_BEANS_CASH_DEDUCTION) {
          return new Promise((resolve) => {
            msgHandler.value = (
              event: MessageEvent<{ type?: string; source?: string }>
            ) => {
              if (event.data?.source !== 'qsbuy') {
                return;
              }
              resolve(event.data?.type === 'paySucc');
              window.removeEventListener('message', msgHandler.value);
            };
            window.addEventListener('message', msgHandler.value);

            const params = new URLSearchParams(
              pickBy(
                {
                  id: taskId,
                  resId: (err.extra as SaveTaskResponse).reservedResourceId!,
                  lang: locale.value,
                },
                identity
              )
            );
            window.open(
              `/${
                isDev ? 'readpaper-ai' : 'home'
              }/qsbuy.html?${params.toString()}`,
              '_blank'
            );
          });
        } else {
          message.error(err?.message || 'Unknown Error');
        }

        return false;
      });

    if (res) {
      message.success(t('common.aiReviewerQS.createSuccess'));
      emit('created', taskId);
      saved.value = true;
      onCancel();
    }

    return res;
  },
  {
    manual: true,
  }
);
const onSaveTask = async () => {
  try {
    isSaving.value = true;
    const isValid = await form.value.validate();
    if (isValid) {
      let flag = true;
      if (beans.value > props.data!.currentAiBeanPrice) {
        flag = await consumeBeans(
          BeanScenes.REVIEWQS,
          props.data!.currentAiBeanPrice
        );
      } else {
        beans.value -= props.data!.currentAiBeanPrice;
      }

      if (flag) {
        await doSaveTask();
      }
    }
  } finally {
    isSaving.value = false;
  }
};

useRequest(
  async () => {
    const { taskId } = formState.value;
    if (!taskId || saved.value) {
      return;
    }
    const res = await getQSReviewTask({
      taskId,
    });

    if (
      [
        TaskStatus.REVIEWING,
        TaskStatus.SUCCESS,
        TaskStatus.FAIL,
        TaskStatus.CONSUME_BEAN_FAIL,
        TaskStatus.CANCEL,
      ].includes(res.taskStatus)
    ) {
      message.success(t('common.aiReviewerQS.createSuccess'));
      emit('created', taskId);
      saved.value = true;
      onCancel();
    }

    return res;
  },
  {
    // @ts-ignore
    ready: computed(() => !!formState.value.taskId),
    pollingInterval: 3000,
    pollingWhenHidden: false,
    pollingSinceLastFinished: true,
  }
);

// const { run: onCancelTask } = useRequest(
//   async () => {
//     if (!formState.value.taskId) {
//       return
//     }
//     await cancelQSReviewTask({
//       taskId: formState.value.taskId,
//     })
//   },
//   {
//     manual: true,
//   }
// )

const onCancel = () => {
  if (saved.value) {
    formState.value = { ...DEFAULT_STATE };
  }
  saved.value = false;
  visible.value = false;
};

// useEventListener(window, 'beforeunload', () => {
//   onCancel()
// })
</script>
