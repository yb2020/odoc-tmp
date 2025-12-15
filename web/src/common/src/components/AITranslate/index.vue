<script setup lang="ts">
import { ref, nextTick, watch, onActivated } from 'vue';
import Tip from '@common/components/AITools/Tip.vue';
import Panel from '@common/components/AITools/Panel.vue';
import TextInput from '@common/components/AITools/TextInput.vue';
import { TabKeyType, TextModeType } from '@common/components/AITools/type';
import TextMode from '@common/components/AITools/TextMode.vue';
import Sentences from '@common/components/AITools/Sentences/index.vue';
import Paragraph from '@common/components/AITools/Paragraph.vue';
import { useI18n } from 'vue-i18n';
import { useLocalStorage } from '@vueuse/core';
import {
  reportModuleImpression,
  PageType,
  ModuleType,
} from '@common/utils/report';

const props = defineProps<{
  allowed: boolean;
}>();

const emit = defineEmits<{
  (e: 'intercepted'): void;
  (e: 'started'): void;
}>();

const sentencesRef = ref<InstanceType<typeof Sentences>>();

const paragraphRef = ref<InstanceType<typeof Paragraph>>();

const textInputRef = ref<InstanceType<typeof TextInput>>();

const isInSentencesProgress = ref(false);

const textMode = useLocalStorage<TextModeType>(
  'polish/ai-translate-client-textmode',
  TextModeType.PARAGRAPH
);

const onClear = () => {
  if (textMode.value === TextModeType.PARAGRAPH) {
    paragraphRef.value?.clear();
  } else {
    sentencesRef.value?.clear();
  }
};

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

watch(
  () => textMode.value,
  () => {
    textInputRef.value?.setButtonType('start');
    paragraphRef.value?.clear();
    sentencesRef.value?.clear();
  }
);

onActivated(() => {
  // 上报模块曝光
  reportModuleImpression({
    page_type: PageType.POLISH,
    module_type: ModuleType.ZH_TO_EN,
  });
});
</script>
<template>
  <div class="flex flex-col h-full">
    <Tip
      class="p-4"
      :text="$t('common.aitools.aiTranslateTip')"
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
          <div class="text-base font-medium pr-4 text-rp-neutral-10 leading-7">
            {{ $t('common.aitools.aiTranslate') }}
          </div>
          <TextMode v-model:mode="textMode" />
        </div>
      </template>
      <div class="flex flex-1 h-0">
        <TextInput
          v-show="
            !(textMode === TextModeType.SENTENCE && isInSentencesProgress)
          "
          ref="textInputRef"
          :placeholder="$t('common.aitools.aiTranslatePlaceholder')"
          :type="TabKeyType.translate"
          :start="startBuild"
          :is-chinese="true"
          @clear="onClear"
        />
        <Paragraph
          v-if="textMode === TextModeType.PARAGRAPH"
          ref="paragraphRef"
          :type="TabKeyType.translate"
          class="flex-1"
          :is-diff="false"
        />
      </div>
    </Panel>
    <Sentences
      v-if="textMode === TextModeType.SENTENCE"
      ref="sentencesRef"
      v-model:isInProcess="isInSentencesProgress"
      class="flex-1 h-0"
      :is-diff="false"
      :type="TabKeyType.translate"
    />
  </div>
</template>
