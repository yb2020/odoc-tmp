<template>
  <a-drawer
    v-model:visible="visible"
    placement="bottom"
    :mask="false"
    height="204px"
    :closable="false"
    :style="{
      position: 'absolute',
    }"
    :get-container="false"
  >
    <div class="textbox-container">
      <div class="inner">
        <a-textarea
          v-if="type === 'question'"
          v-model:value="input"
          :placeholder="placeholder"
          :auto-size="{ minRows: 5, maxRows: 5 }"
          :maxlength="1000"
          showCount
        />
        <div
          v-else
          class="answer-markdown-edit-wrap"
        >
          <MarkdownEditor
            :raw="input"
            :uniq-id="new Date().getTime() + ''"
            :placeholder="placeholder"
            :min-rows="5"
            :max-rows="5"
            :allow-enter="true"
            :maxLength="1000"
            :upload="upload"
            @change="handleInputChange"
          />

          <div class="tips-wrap">
            <MarkdownTip />
            <div class="count">
              {{ input.length || 0 }} / 1000
            </div>
          </div>
        </div>

        <div
          :class="[
            'btn-container',
            { 'answer-btn-container': type === 'answer' },
          ]"
        >
          <slot name="publish" />
          <a-button
            type="primary"
            :disabled="disable"
            class="btn"
            :loading="pulishLoading"
            @click="handlePublish"
          >
            发布
          </a-button>
          <a-button
            class="btn cancel"
            @click="handleCancel"
          >
            取消
          </a-button>
        </div>
      </div>
    </div>
  </a-drawer>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { message } from 'ant-design-vue';
import { MarkdownEditor } from '@idea/aiknowledge-markdown';
import '@idea/aiknowledge-markdown/dist/style.css';
import {
  IMAGE_PLACEHOLDER_UPLOADING_REGEX,
  IMAGE_PLACEHOLDER_ERROR,
} from '@idea/aiknowledge-markdown/src/components/helper';
import { uploadImage, ImageStorageType } from '~/src/api/upload';
import MarkdownTip from '@/components/Common/MarkdownTip.vue';

const { publishFn } = defineProps<{
  placeholder: string;
  publishFn: (content: string) => Promise<void>;
  type: 'question' | 'answer';
}>();

const visible = defineModel('visible', { default: false });

const input = ref('');

const disable = computed(() => input.value.trim().length < 5);

const pulishLoading = ref(false);

const handlePublish = async () => {
  input.value = input.value.replace(
    IMAGE_PLACEHOLDER_UPLOADING_REGEX,
    IMAGE_PLACEHOLDER_ERROR
  );

  if (disable.value) {
    message.info('请至少输入 5 个有效字符');
    return;
  }

  if (input.value.length > 1000) {
    message.info('最多输入 1000 个字符');
    return;
  }

  pulishLoading.value = true;

  try {
    await publishFn(input.value);

    message.success('发布成功！');
    emit('success');
    handleCancel();
  } catch (error) {}
  pulishLoading.value = false;
};

const emit = defineEmits<{
  (event: 'success'): void;
}>();

const handleCancel = () => {
  input.value = '';
  visible.value = false;
};

const handleInputChange = (data: string) => {
  input.value = data;
};

const upload = async (src: File | string) => {
  return uploadImage(src, ImageStorageType.markdown);
};
</script>

<style lang="less" scoped>
.textbox-container {
  width: 100%;
  padding: 10px;
  bottom: 0;
  background: #383a3d;
  .inner {
    overflow: hidden;
  }

  :deep(textarea) {
    color: rgba(0, 0, 0, 0.85);
    border: 1px solid #d9d9d9;
    background-color: #fff;
  }

  :deep(textarea::placeholder) {
    font-size: 14px;
    font-weight: 400;
    color: rgba(0, 0, 0, 0.25);
  }

  .btn-container {
    display: flex;
    align-items: center;
    margin-top: 24px;
    flex-direction: row-reverse;

    .btn {
      width: 93px;
      height: 32px;
    }

    .cancel {
      background: #ffffff;
      margin-right: 16px;
    }
  }
  .answer-btn-container {
    margin-top: 0;
  }
}

html[data-theme='dark'] {
  .textbox-container {
    :deep(textarea::placeholder) {
      color: #4e5969;
    }
  }
}

.answer-markdown-edit-wrap {
  .tips-wrap {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: -6px;
    .count {
      font-weight: 400;
      font-size: 12px;
      line-height: 18px;
      color: rgba(255, 255, 255, 0.3);
    }
  }
}

.ant-input-textarea-show-count::after {
  font-size: 12px;
  line-height: 18px;
  margin: 3px 0;
}
</style>
