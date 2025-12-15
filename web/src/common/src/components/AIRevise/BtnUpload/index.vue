<script setup lang="ts">
import { ref } from 'vue';
import Progress from './Progress.vue';
import Btn from './Btn.vue';
import { ObjectStoreInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/oss/AliOSS';
import { parseFile } from '@common/api/latex';
import { createProject } from '@common/api/revise';
import { ResponseError } from '@common/api/type';
import { ERROR_CODE_NEED_VIP } from '@common/api/const';

const btnProps = {
  size: 'large',
  class: '!rounded-2xl',
};

const importBucketError = ref<null | Error>(null);
const importBucketProgress = ref(-1);

const onImportBucketError = (err: Error) => {
  if (err instanceof ResponseError && err.code === ERROR_CODE_NEED_VIP) {
    onProgressCancel();
    emit('errvip', err);
    return;
  }

  importBucketError.value = err;
};

const props = defineProps<{
  alwaysNew?: boolean;
  curProjectId?: string;
  curVersionId?: string;
}>();

const emit = defineEmits<{
  (event: 'preclick', e: MouseEvent): void;
  (event: 'import:success', x: string): void;
  (event: 'errvip', err: ResponseError): void;
}>();

const entry = ref<string>();
const entryList = ref<string[]>([]);
const entryResolve = ref<(v: unknown) => void>();

const ensureParseFile = async (
  bucket: ObjectStoreInfo,
  curProjectId?: string,
  curVersionId?: string,
  resolved = false
) => {
  if (props.alwaysNew) {
    if (!curProjectId) {
      ({ projectId: curProjectId } = await createProject());
    }
  } else {
    curProjectId = props.curProjectId!;
    curVersionId = props.curVersionId!;
  }

  const res = await parseFile({
    pageId: curVersionId,
    projectId: curProjectId,
    ...bucket,
    latexName: entry.value,
    chooseLatex: resolved ? !entry.value : undefined,
  });

  if (res.latexNames?.length && typeof entry.value === 'undefined') {
    const isContinue = await new Promise((resolve) => {
      entryResolve.value = resolve;
      entryList.value = res.latexNames;
    });
    if (isContinue === false) {
      return false;
    }
    await ensureParseFile(bucket, curProjectId, curVersionId, true);
  }

  return {
    curProjectId,
    curVersionId,
  };
};

const onCancelEntry = async (e: PointerEvent) => {
  const isSkip = (e.currentTarget as HTMLElement)?.classList.contains(
    'ant-btn'
  );

  entryResolve.value?.(isSkip ? '' : false);

  entryList.value = [];
};

const onConfirmEntry = async () => {
  if (!entry.value) {
    return;
  }
  entryResolve.value?.(entry.value);

  entryList.value = [];
};

const onImportBucketProgress = async (
  progress: number,
  bucket?: ObjectStoreInfo
) => {
  if (bucket) {
    importBucketProgress.value = 99;
    try {
      const res = await ensureParseFile(bucket);
      onProgressCancel();

      if (res) {
        emit('import:success', res.curProjectId);
      }
    } catch (error) {
      onImportBucketError(error as Error);
    }
  } else {
    importBucketProgress.value = progress;
  }
};

const onProgressCancel = () => {
  importBucketProgress.value = -1;
  importBucketError.value = null;
};
</script>
<template>
  <Btn
    @preclick="emit('preclick', $event)"
    @import:bucket:error="onImportBucketError"
    @import:bucket:progress="onImportBucketProgress"
  />
  <Progress
    v-if="importBucketProgress >= 0"
    :error="importBucketError"
    :progress="importBucketProgress"
    @cancel="onProgressCancel"
  />
  <a-modal
    v-if="entryList.length > 0"
    class="modal-latex-picker"
    visible
    :width="560"
    :header="null"
    :cancelText="$t('common.aitools.latexEntryPick.skip')"
    :okButtonProps="btnProps"
    :cancelButtonProps="btnProps"
    @ok="onConfirmEntry"
    @cancel="onCancelEntry"
  >
    <h3 class="text-2xl text-center mb-10">
      {{ $t('common.aitools.latexEntryPick.tt') }}
    </h3>
    <p class="text-base mb-6">
      {{ $t('common.aitools.latexEntryPick.tip') }}
    </p>
    <a-radio-group
      v-model:value="entry"
      class="!flex flex-col"
    >
      <a-radio
        v-for="item in entryList"
        :value="item"
      >
        {{ item }}
      </a-radio>
    </a-radio-group>
  </a-modal>
</template>

<style lang="postcss">
.modal-latex-picker {
  .ant-modal-body {
    padding: 20px theme('spacing.12') theme('spacing.10');
  }
  .ant-modal-footer {
    border: 0;
    padding: 0 theme('spacing.12') theme('spacing.10');
  }
  .ant-btn:not(.ant-btn-primary) {
    color: theme('colors.rp-neutral-8');
    border-color: transparent;
    background-color: theme('colors.rp-neutral-2');

    &:hover {
      color: theme('colors.rp-neutral-10');
    }
  }
}
</style>
