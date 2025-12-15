<template>
  <div
    v-if="!onlyAddToParse"
    class="divide"
  >
    <span class="title">{{ $t('translate.translation') }}</span>
    <span class="line" />
    <CopyOutlined
      class="text-2xl !text-rp-neutral-8 ml-3 cursor-pointer"
      @click="handleCopy"
    />
    <span
      v-if="allowToAddNote && !store.multiSegment"
      class="icon"
    >
      <a-tooltip
        v-if="!store.isExistingAnnotation"
        :getPopupContainer="getPopupContainer"
        placement="left"
      >
        <template #title>{{ $t('translate.addToAnnotation') }}</template>
        <i
          class="aiknowledge-icon text-2xl icon-add-to-note"
          @click="addToNote(false)"
        />
      </a-tooltip>
      <a-tooltip
        v-if="isPhrase"
        :getPopupContainer="getPopupContainer"
        placement="left"
      >
        <template #title>{{ $t('translate.addToWordPhrase') }}</template>
        <i
          class="aiknowledge-icon text-2xl ml-2 icon-add-to-phrase"
          @click="addToNote(true)"
        />
      </a-tooltip>
    </span>
  </div>
  <div v-else>
    <a-tooltip
      :getPopupContainer="getPopupContainer"
      placement="left"
    >
      <template #title>
        {{ $t('translate.addToWordPhrase') }}
      </template>
      <i
        class="aiknowledge-icon text-xl ml-2 icon-add-to-phrase text-rp-neutral-4 cursor-pointer"
        @click="addToNote(true)"
      />
    </a-tooltip>
  </div>
</template>
<script setup lang="ts">
import { CopyOutlined } from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { UniTranslateResp } from '~/src/api/translate';
import { useTranslateStore } from '~/src/stores/translateStore';
import { copyToPaste } from '~/src/util/copy';

const props = defineProps<{
  translatedContent: string;
  translatedData?: UniTranslateResp;
  input: string;
  onlyAddToParse?: boolean;
  addToNoteHandler: (
    isPhrase: boolean,
    phrase: string,
    translation: string,
    translationRes: UniTranslateResp
  ) => void;
}>();

const isPhrase = computed(() => props.input.length < 60);

const store = useTranslateStore();

const { t } = useI18n();

const addToNote = (isAddToPhrase: boolean) => {
  const translatedData = props.translatedData;
  console.log('addToNote', props.translatedData);
  if (translatedData) {
    props.addToNoteHandler(
      isAddToPhrase,
      props.input,
      !isAddToPhrase && translatedData.targetResp?.length
        ? translatedData.targetResp[0].targetContent?.join(' ')
        : props.translatedContent,
      translatedData
    );
    if (isAddToPhrase) {
      message.success(t('translate.addToWordPhraseSuccessTip'));
    }
  }
};

const translateStore = useTranslateStore();

const allowToAddNote = computed(() => translateStore.allowToAddNote);

const handleCopy = () => {
  const translatedData = props.translatedData;
  let content = props.translatedContent;
  if (translatedData?.targetResp?.length) {
    content = translatedData.targetResp
      .map((item) => item.targetContent.join(' '))
      .join(' ');
  }
  copyToPaste(content);
};

const getPopupContainer = (triggerNode: HTMLElement) => {
  return triggerNode.closest('.js-translate-tippy-viewer') || document.body;
};
</script>
<style scoped lang="less">
.divide {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  height: 22px;
  line-height: 22px;

  .title {
    font-weight: 600;
    color: #1d2129;
    margin-right: 12px;
  }

  .line {
    height: 1px;
    background: #e4e7ed;
    flex: 1;
  }

  .icon {
    color: #4e5969;
    margin-left: 12px;
    cursor: pointer;
  }
}
</style>
