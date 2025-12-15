<template>
  <div v-html="showTextWithCursor" />
</template>
<script setup lang="ts">
import { computed, onUnmounted, onUpdated, ref, watch } from 'vue';
import { Typed } from './Typed';

const props = defineProps<{
  text: string;
  scrollToBottom?: () => void;
  isPending: boolean;
}>();

const showText = ref('');

const isTyping = ref(true);

const showTextWithCursor = computed(() => {
  return (
    showText.value +
    (isTyping.value
      ? `<span class="typed-cursor ${showText.value ? 'ml-1' : ''}">_</span>`
      : '')
  );
});

const emit = defineEmits<{
  (event: 'typing:finished'): void;
}>();

const typed: Typed = new Typed({
  strings: [
    {
      string: props.text,
      typingTime: 2000,
    },
  ],
  typeSpeed: 30,
});
typed
  .onTyping((text) => {
    showText.value = text.fullString;
    isTyping.value = true;
  })
  .onStop(() => {
    if (!props.isPending) {
      isTyping.value = false;
      emit('typing:finished');
    }
  });

onUnmounted(() => {
  typed.destroy();
});

watch(
  () => props.text,
  (newVal, oldVal) => {
    const diff = newVal.slice(oldVal.length);
    if (diff.length) {
      typed.flushStrings({
        string: diff,
        typingTime: 2000,
      });
    }
  }
);

watch(
  () => props.isPending,
  (newVal) => {
    if (!newVal) {
      typed.flushStrings({
        string: '',
        typingTime: 2000,
      });
    }
  }
);

onUpdated(() => {
  props.scrollToBottom?.();
});
</script>
