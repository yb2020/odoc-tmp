<template>
  <a-modal
    maskClosable
    wrapClassName="modal-refundbeans"
    :width="350"
    :centered="true"
    :visible="visible"
    :okButtonProps="{
      loading: isRefunding,
      disabled: !formState.reason.length,
    }"
    @ok="onRefundConfirm"
    @cancel="onRefundCancel"
  >
    <template #title>
      <a-tooltip>
        <template #title>
          <div>{{ $t('common.aibeans.refundTip') }}</div>
        </template>
        {{ $t('common.aibeans.refund') }}
        <QuestionCircleOutlined />
      </a-tooltip>
    </template>
    <p>{{ $t('common.aibeans.refundTt') }}</p>
    <a-form
      :model="formState"
      layout="vertical"
      autocomplete="off"
    >
      <a-form-item
        :label="$t('common.aibeans.refundReason')"
        name="reason"
        :rules="[{ required: true }]"
      >
        <LoadingOutlined v-if="loading" />
        <a-space
          v-else
          :size="[0, 8]"
          wrap
          class="!-mb-px"
        >
          <a-checkable-tag
            v-for="reason in allReasons"
            :key="reason"
            :checked="formState.reason.includes(reason)"
            @change="(checked: boolean) => handleTagSelect(reason, checked)"
          >
            {{ reason }}
          </a-checkable-tag>
        </a-space>
      </a-form-item>
      <a-form-item
        :label="$t('common.aibeans.refundDesc')"
        name="supplementaryExplanation"
      >
        <a-textarea
          v-model:value="formState.supplementaryExplanation"
          class="!border-rp-neutral-4 resize-none !text-rp-neutral-10"
        />
      </a-form-item>
      <a-form-item
        :label="$t('common.aibeans.refundImage')"
        name="screenshotUrl"
      >
        <a-upload
          :max-count="1"
          :multiple="false"
          :customRequest="handleUpload"
          :show-upload-list="true"
          :file-list="fileList"
          @remove="handleRemoveFile"
        >
          <a-button>
            <upload-outlined />
            <template v-if="formState.screenshotUrl">
              {{
                $t('common.text.re')
              }}
            </template>
            {{ $t('common.text.upload') }}
          </a-button>
        </a-upload>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { BizType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/polish/PolishFeedbackInfo';
import {
  LoadingOutlined,
  QuestionCircleOutlined,
  UploadOutlined,
} from '@ant-design/icons-vue';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { standardToLanguageEnum } from '../../../../shared/language/service';
import { useAIBeansRefundReasons } from '../../hooks/useAIBeans';
import { RefundReasonScene, refundAiBeans } from '../../api/aibeans';
import { UploadProps } from 'ant-design-vue';
import { uploadImageFile, ImageStorageType } from '../../api/upload';
import { useRequest } from 'ahooks-vue';

const visible = defineModel('visible', { default: false });

const props = defineProps<{
  tid: string;
  ttype?: BizType;
  scene: RefundReasonScene;
}>();
const emit = defineEmits(['cancel', 'ok']);

const formState = ref<
  Omit<
    /*RefundAiBeanReq*/ {
      reason: string[];
      supplementaryExplanation?: string;
      screenshotUrl?: string;
      screenshotName?: string;
    },
    'answerId'
  >
>({
  reason: [],
});
const fileList = computed(() => {
  const arr: UploadProps['fileList'] = [];

  if (formState.value.screenshotUrl) {
    arr.push({
      uid: '1',
      name: formState.value.screenshotName || 'Unknown',
      status: 'done',
      url: formState.value.screenshotUrl,
    });
  }

  return arr;
});

const { locale } = useI18n();
const { data, loading, refresh } = useAIBeansRefundReasons(props.scene);

// 判断当前是否为中文语言
const isChineseLanguage = computed(() => {
  const langEnum = standardToLanguageEnum(locale.value);
  return langEnum === Language.ZH_CN;
});

const allReasons = computed(
  () =>
    data.value?.map((x) => {
      return isChineseLanguage.value ? x.text : x.textEn ?? x.text;
    }) ?? []
);

watch(visible, (v) => {
  if (v && !data.value) {
    refresh();
  }
});

const handleTagSelect = (reason: string, checked: boolean) => {
  let reasons = formState.value.reason;
  if (checked) {
    if (!reasons.includes(reason)) {
      reasons.push(reason);
    }
  } else {
    const i = reasons.indexOf(reason);
    reasons.splice(i, 1);
  }
  // triggerRef(formState)
};

const handleUpload: UploadProps['customRequest'] = async ({ file }) => {
  const url = await uploadImageFile(file as File, ImageStorageType.screenshot);

  formState.value.screenshotName = (file as File).name;
  formState.value.screenshotUrl = url;
};

const handleRemoveFile = () => {
  formState.value.screenshotName = '';
  formState.value.screenshotUrl = '';
};

const onRefundCancel = () => {
  visible.value = false;

  emit('cancel');
};

const { run: onRefundConfirm, loading: isRefunding } = useRequest(
  async () => {
    if (!formState.value.reason?.length) {
      return;
    }

    const k = props.scene === RefundReasonScene.Copilot ? 'answerId' : 'taskId';
    const res = await refundAiBeans(props.scene, {
      ...formState.value,
      [`${k}`]: props.tid,
      bizType: props.ttype,
    });

    visible.value = false;

    emit('ok', res?.optimizeId);
  },
  {
    manual: true,
  }
);
</script>

<style lang="less">
.modal-refundbeans {
  .ant-modal-close-x {
    @apply h-12;
    line-height: theme('spacing.12');
  }

  .ant-modal-header {
    @apply py-3;
    @apply mt-px;
  }

  .ant-modal-body {
    @apply pb-0;
    @apply pt-4;
    line-height: 1.375rem;
  }

  .ant-modal-footer {
    @apply p-4;
    @apply -mt-4;
    @apply border-t-0;
  }

  .ant-form-item {
    @apply mb-4;
  }

  .ant-space-item {
    height: 1.375rem;
    line-height: 1.375rem;
  }

  .ant-form-vertical .ant-form-item-label {
    line-height: 1.375rem;

    &
      > label.ant-form-item-required:not(
        .ant-form-item-required-mark-optional
      )::before {
      position: absolute;
      left: 100%;
    }
  }

  .ant-tag-checkable {
    color: #1d2129;
    @apply bg-rp-neutral-2;
    @apply !transition-none;

    &:active,
    &.ant-tag-checkable-checked {
      @apply bg-rp-blue-1;
      @apply text-rp-blue-6;
    }
  }

  .ant-upload-list-item-card-actions-btn.ant-btn-sm {
    @apply h-6;
  }

  .ant-upload-list-item,
  .ant-upload-list-item-info,
  .ant-upload-list-item-info > span {
    height: auto;
  }
}
</style>
