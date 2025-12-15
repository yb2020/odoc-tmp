<script setup lang="ts">
import { onMounted, ref } from 'vue';
import * as jsdiff from 'diff';

const props = defineProps<{
  originalText: string;
  modifiedText: string;
  diffChars?: boolean;
}>();

const diffTxtHtml = ref('');

const diff = () => {
  const diffs = props.diffChars
    ? jsdiff.diffChars(props.originalText, props.modifiedText)
    : jsdiff.diffWords(props.originalText, props.modifiedText);
  let html = '';
  diffs.map((diff: any) => {
    if (diff.added) {
      html += `<span class="added">${diff.value}</span>`;
    } else if (diff.removed) {
      html += `<span class="removed">${diff.value}</span>`;
    } else {
      html += diff.value;
    }
  });
  diffTxtHtml.value = html;
};

onMounted(() => {
  diff();
});
</script>
<template>
  <div
    class="wrapper-dw"
    v-html="diffTxtHtml"
  />
</template>
<style lang="less" scoped>
.wrapper-dw:deep(.added) {
  color: theme('colors.rp-blue-6');
}

.wrapper-dw:deep(.removed) {
  color: theme('colors.rp-red-6');
  text-decoration: line-through;
}
</style>
