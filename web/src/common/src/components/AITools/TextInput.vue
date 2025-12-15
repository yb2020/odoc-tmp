<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import wordsCount from 'words-count';
import { TabKeyType } from './type';
import { SnippetsOutlined } from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  start: (text: string, words: number) => Promise<boolean>;
  placeholder: string;
  type: TabKeyType;
}>();

const textValue = ref('');
const textCount = computed(() =>
  props.type === TabKeyType.translate || props.type === TabKeyType.zhpolish
    ? textValue.value.trim().length
    : wordsCount(textValue.value)
);
const buttonType = ref<'start' | 'retry'>('start');
const buttonLoading = ref(false);

const startButtonDisabled = computed(() => {
  return !textValue.value.trim();
});

const emit = defineEmits<{
  (event: 'clear'): void;
}>();

watch(
  () => textValue.value,
  () => {
    if (!textValue.value.trim()) {
      buttonType.value = 'start';
    }
  }
);

const startBuild = async () => {
  buttonLoading.value = true;
  try {
    await props.start(textValue.value.trim(), textCount.value);
    buttonType.value = 'retry';
  } catch (error) {
  } finally {
    buttonLoading.value = false;
  }
};

const handleClear = () => {
  textValue.value = '';
  buttonType.value = 'start';
  emit('clear');
};

defineExpose({
  setButtonType(type: 'start' | 'retry') {
    buttonType.value = type;
  },
});

const { t } = useI18n();

const handlePaste = (e: MouseEvent) => {
  e.preventDefault();
  navigator.clipboard
    .readText()
    .then((text) => {
      textValue.value = text;
    })
    .catch(() => {
      message.warn(t('common.aitools.pasteErrorTip'));
    });
};
</script>
<template>
  <div class="flex-1 p-6 border-r border-rp-neutral-3">
    <div class="flex flex-col h-full">
      <div class="flex-1 relative">
        <a-textarea
          v-model:value="textValue"
          class="!h-full !border-none !px-0"
          :placeholder="placeholder"
        />
        <div
          v-if="!textValue"
          class="absolute cursor-pointer px-5 py-3 text-base text-rp-blue-6 flex flex-col justify-center items-center top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-xl border border-rp-blue-6"
          @click="handlePaste"
        >
          <SnippetsOutlined class="text-xl" />
          <div class="pt-1 text-sm">
            {{ $t('common.aitools.pasteText') }}
          </div>
        </div>
      </div>
      <div class="flex justify-between items-center mt-4">
        <span class="text-rp-neutral-6">{{ textCount }}
          {{ $t('common.text.words', textCount > 1 ? 2 : 1) }}</span>
        <a-button
          v-if="buttonType === 'start'"
          type="primary"
          shape="round"
          class=""
          :loading="buttonLoading"
          :disabled="startButtonDisabled"
          @click="startBuild"
        >
          {{
            $t(
              type !== TabKeyType.translate
                ? 'common.aitools.generate'
                : 'common.aitools.translate'
            )
          }}
        </a-button>
        <div
          v-else
          class="space-x-6"
        >
          <a-button
            shape="round"
            class="!border-rp-neutral-4 !text-black"
            @click="handleClear"
          >
            {{ $t('common.aitools.clear') }}
          </a-button>
          <a-button
            shape="round"
            type="default"
            :loading="buttonLoading"
            :disabled="startButtonDisabled"
            @click="startBuild"
          >
            {{
              $t(
                type !== TabKeyType.translate
                  ? 'common.aitools.regenerate'
                  : 'common.aitools.retranslate'
              )
            }}
          </a-button>
        </div>
      </div>
    </div>
  </div>
</template>
