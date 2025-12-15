<template>
  <div
    ref="el"
    @click="handleQuickAsk"
  >
    <MarkdownViewer
      :text="formatedText"
      :cursor="cursor"
    />
  </div>
</template>
<script setup lang="ts">
import { computed, onUnmounted, onUpdated, ref, watch } from 'vue';
import { Typed } from './Typed';
import MarkdownViewer from './MarkdownText.vue';
import { formatQuestion } from './tool';
import lodash from 'lodash';

const props = defineProps<{
  text: string;
  scrollToBottom: () => void;
  isPending: boolean;
  anchors?: string[];
}>();

const showText = ref('');

const isTyping = ref(true);

const cursor = computed(() => {
  return isTyping.value || props.isPending
    ? `<span class="typed-cursor ${
        showText.value ? 'ml-1' : ''
      }" style="vertical-align: sub"><span style="display: inline-block; width: 8px; height: 16px; line-height: 16px; background-color: #1f71e0;" /></span>`
    : '';
});

const formatedText = computed(() => {
  return formatQuestion(showText.value, props.anchors);

  // return (
  //   answer +
  //   (isTyping.value || props.isPending
  //     ? `<span class="typed-cursor ${
  //         showText.value ? 'ml-1' : ''
  //       }" style="vertical-align: sub"><span style="display: inline-block; width: 8px; height: 16px; line-height: 16px; background-color: #1f71e0;" /></span>`
  //     : '')
  // );
});

const typed: Typed = new Typed({
  strings: [
    {
      string: props.text,
      typingTime: 200,
    },
  ],
  typeSpeed: 5, //基础打字速度，即每个字符之间的输入延迟大约是 30 毫秒
});
const handleTyping = lodash.throttle((text) => {
  showText.value = text.fullString;
  props.scrollToBottom();
}, 50);

typed
  .onTyping((text) => {
    handleTyping(text);
    isTyping.value = true;
  })
  .onStop(() => {
    // onStop的时候，要取消掉throttle, 并且立即执行一次，保证最后一次的数据能被更新
    handleTyping.flush();
    isTyping.value = false;
  });

const el = ref<HTMLElement>();

onUnmounted(() => {
  typed.destroy();
});

watch(
  () => props.text,
  (newVal, oldVal) => {
    if (newVal === '') {
      typed.reset();
      showText.value = '';
      isTyping.value = true;
    }
    const diff = newVal.slice(oldVal.length);
    if (diff.length) {
      typed.flushStrings({
        string: diff,
        typingTime: 2000,
      });
    }
  }
);
const handleTextShowComplete = lodash.debounce(() => {
  if (!isTyping.value && !props.isPending) {
    emit('textShowComplete');
  }
}, 300);

watch(
  () => [isTyping.value, props.isPending],
  () => {
    handleTextShowComplete();
  },
  { immediate: true }
);



const emit = defineEmits<{
  (event: 'quickAskQuestion', e: MouseEvent): void;
  (event: 'textShowComplete'): void;
}>();

const handleQuickAsk = (e: MouseEvent) => {
  emit('quickAskQuestion', e);
};
</script>
