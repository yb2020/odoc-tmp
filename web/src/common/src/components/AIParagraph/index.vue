<script setup lang="ts">
import { nextTick, onActivated, ref, watch } from 'vue';
import Tip from '../AITools/Tip.vue';
import Panel from '../AITools/Panel.vue';
import Mode from './components/Mode.vue';
import TextMode from '../AITools/TextMode.vue';
import Paragraph from '../AITools/Paragraph.vue';
import TextInput from '../AITools/TextInput.vue';
import Sentences from '../AITools/Sentences/index.vue';
import { TabKeyType, TextModeType } from '../AITools/type';
import { ParagraphMode } from '../AIParagraph/type';
import { useLocalStorage } from '@vueuse/core';
import { useI18n } from 'vue-i18n';
import {
  reportModuleImpression,
  PageType,
  ModuleType,
} from '@common/utils/report';

const props = defineProps<{
  allowed: boolean;
  type: TabKeyType.polish | TabKeyType.zhpolish;
}>();

const emit = defineEmits<{
  (e: 'intercepted'): void;
  (e: 'started'): void;
}>();

const textMode = useLocalStorage<TextModeType>(
  'polish/ai-paragraph-client-textmode' +
    (props.type === TabKeyType.zhpolish ? '-zh' : ''),
  TextModeType.PARAGRAPH
);

const polishMode = ref(ParagraphMode.improve);

const paragraphRef = ref<InstanceType<typeof Paragraph>>();

const sentencesRef = ref<InstanceType<typeof Sentences>>();

const textInputRef = ref<InstanceType<typeof TextInput>>();

const isInSentencesProgress = ref(false);

const { t } = useI18n();

const startBuild = async (textValue: string, textCount: number) => {
  if (!props.allowed) {
    emit('intercepted');
    return false;
  }
  // if (textCount > 300) {
  //   message.error(t('revise.client.wordsLimitTip'))
  //   return false
  // }

  if (textMode.value === TextModeType.PARAGRAPH) {
    await paragraphRef.value?.startParagraph(textValue);
  } else {
    isInSentencesProgress.value = true;
    await nextTick();
    await sentencesRef.value?.startSentences(textValue);
  }

  emit('started');
  return true;
};

const onClear = () => {
  if (textMode.value === TextModeType.PARAGRAPH) {
    paragraphRef.value?.clear();
  } else {
    sentencesRef.value?.clear();
  }
};

watch(
  () => textMode.value,
  () => {
    paragraphRef.value?.clear();
    sentencesRef.value?.clear();
    setTimeout(() => {
      textInputRef.value?.setButtonType('start');
    }, 0);
  }
);

watch(
  () => polishMode.value,
  () => {
    paragraphRef.value?.clear();
    sentencesRef.value?.clear();
    setTimeout(() => {
      textInputRef.value?.setButtonType('start');
    }, 0);
  }
);

onActivated(() => {
  // 上报模块曝光
  reportModuleImpression({
    page_type: PageType.POLISH,
    module_type:
      props.type === TabKeyType.zhpolish
        ? ModuleType.POLISH_REWRITE_ZH
        : ModuleType.POLISH_REWRITE,
  });
});
</script>
<template>
  <div class="flex flex-col h-full">
    <Tip
      class="p-4"
      :text="
        props.type === TabKeyType.zhpolish
          ? $t('common.aitools.zhAiPolishTip')
          : $t('common.aitools.aiPolishTip')
      "
    >
      <template #beans>
        <slot name="beans" />
      </template>
    </Tip>
    <Panel
      v-show="!isInSentencesProgress"
      class="flex-1 h-0 m-4 mt-0"
    >
      <template #title>
        <div class="px-6 py-3 flex justify-between">
          <Mode
            v-model:mode="polishMode"
            :disabled-mode="
              textMode === TextModeType.SENTENCE
                ? [
                  ParagraphMode.expand,
                  ParagraphMode.shorten,
                  ParagraphMode.simple,
                  ParagraphMode.standard,
                ]
                : []
            "
          />
          <TextMode v-model:mode="textMode" />
        </div>
      </template>
      <div class="flex flex-1 h-0">
        <TextInput
          v-show="
            !(textMode === TextModeType.SENTENCE && isInSentencesProgress)
          "
          ref="textInputRef"
          :placeholder="
            type === TabKeyType.zhpolish
              ? $t('common.aitools.aiZhPolishPlaceholder')
              : $t('common.aitools.aiPolishPlaceholder')
          "
          :type="type"
          :start="startBuild"
          @clear="onClear"
        />
        <Paragraph
          v-if="textMode === TextModeType.PARAGRAPH"
          ref="paragraphRef"
          :type="type"
          :mode-value="polishMode"
          :is-diff="true"
          class="flex-1"
        />
      </div>
    </Panel>
    <Sentences
      v-if="textMode === TextModeType.SENTENCE"
      ref="sentencesRef"
      v-model:isInProcess="isInSentencesProgress"
      :mode-value="polishMode"
      class="flex-1 h-0"
      :is-diff="true"
      :type="type"
    />
  </div>
</template>
