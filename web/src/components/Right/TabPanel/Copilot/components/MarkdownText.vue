<template>
  <div
    class="rp-markdown"
    v-html="markdownText"
  />
</template>
<script setup lang="ts">
import { computed } from 'vue';
import { md } from './tool';

const props = defineProps<{
  text: string;
  cursor?: string;
}>();

const markdownText = computed(() => {
  if ((props.text.match(/```/g)?.length || 0) % 2 === 1) {
    return md.render(renderDeepThink(props.text)) + props.cursor ?? '';
  }
  return md.render(renderDeepThink(props.text) + (props.cursor || ''));
});

const renderDeepThink = (content: string) => {
  // 1. 当检测到props.text中有<think>打头时，把<think>换成<div class='reasoner_think'>
  if (content.includes('<think>') && content.includes('</think>')) {
    content = content.replace('<think>', "<div class='reasoner_think'>");
    content = content.replace('</think>', '</div>');
  } else if (content.includes('<think>') && !content.includes('</think>')) {
    content = content.replace('<think>', "<div class='reasoner_think'>");
  }
  // console.log('content->', content);

  return content;
};
</script>
<style lang="less">
.rp-markdown {
  p:last-child {
    margin-bottom: 0;
  }
}
.reasoner_think {
  color: #979797;
  padding: 10px 0 10px 13px;
  position: relative;
  border-left: 2px solid #979797;
  height: calc(100% - 10px);
}
</style>
