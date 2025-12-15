<script setup lang="ts">
import { SwapOutlined } from '@ant-design/icons-vue';
import { TextModeType } from './type';

const textMode = defineModel('mode', {
  default: TextModeType.PARAGRAPH,
});

const handleChangeTextMode = () => {
  textMode.value =
    textMode.value === TextModeType.PARAGRAPH
      ? TextModeType.SENTENCE
      : TextModeType.PARAGRAPH;
};
</script>
<template>
  <span
    class="cursor-pointer text-rp-blue-6"
    @click="handleChangeTextMode"
  ><SwapOutlined class="mr-1" />{{
    textMode === TextModeType.PARAGRAPH
      ? $t('common.aitools.sentenceMode')
      : $t('common.aitools.paragraphMode')
  }}</span>
</template>
