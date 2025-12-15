<script setup lang="ts">
// import { UploadOutlined } from '@ant-design/icons-vue'
import { UploadEmitEvent, usePDFUpload } from '@common/hooks/usePDFUpload';
import { UploadBizScene } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/oss/AliOSS';
import { getFileIsNeedUpload } from '@common/api/latex';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

type UploadType = 'Latex' | 'Word';

const emit = defineEmits<
  UploadEmitEvent & {
    (e: 'preclick', evt: MouseEvent): void;
  }
>();

const i18n = useI18n();
const id = `upload-${Date.now()}`;
const type = ref<UploadType>('Latex');

const onClick = (e: MouseEvent, v: UploadType) => {
  emit('preclick', e);
  if (!e.defaultPrevented) {
    type.value = v;
  }
};

const { fileList, handleUpload, beforeUpload } = usePDFUpload(emit, {
  acquirePolicyCallbackInfo: getFileIsNeedUpload,
  scene: UploadBizScene.AI_POLISH,
  i18n: i18n as any,
  type,
});
</script>
<template>
  <a-popover
    placement="bottom"
    overlayClassName="popover-upload"
  >
    <a-button
      :id="id"
      shape="round"
      size="large"
      class="!px-6 !flex items-center justify-center text-base font-bold !text-rp-neutral-6 !border-rp-neutral-4 hover:!text-rp-blue-6 hover:!border-rp-blue-6"
      @click.prevent
    >
      <template #icon>
        <img
          src="@common/../assets/images/aitools/icon-upload.svg"
          alt="Upload"
          class="w-6 mr-[10px]"
        >
      </template>{{ $t('common.text.upload') }}{{ $t('common.symbol.space')
      }}{{ $t('common.text.file') }}
    </a-button>
    <template #content>
      <a-upload
        v-model:fileList="fileList"
        name="file"
        :max-count="1"
        :multiple="false"
        :before-upload="beforeUpload"
        :custom-request="handleUpload"
        :show-upload-list="false"
      >
        <label @click="onClick($event, 'Latex')"><img
          src="@common/../assets/images/aitools/icon-latex.svg"
          alt="Latex"
          class="w-5 mr-3"
        >Latex</label>
        <label @click="onClick($event, 'Word')">
          <img
            src="@common/../assets/images/aitools/icon-word.svg"
            alt="Latex"
            class="w-5 mr-3"
          >Word</label>
      </a-upload>
    </template>
  </a-popover>
</template>

<style lang="postcss">
.popover-upload {
  .ant-popover-inner-content {
    display: flex;
    flex-direction: column;
    padding: 12px 0;
  }

  label {
    display: flex;
    align-items: center;
    justify-content: center;
    width: theme('spacing.36');
    height: theme('spacing.8');
    cursor: pointer;

    &:hover {
      background-color: theme('colors.rp-blue-1');
    }
  }
}
</style>
