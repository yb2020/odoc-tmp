<template>
  <div ref="contentRef">
    <slot />
  </div>
</template>
<script lang="ts" setup>
import { onMounted, onUnmounted, ref, watch } from 'vue';
import tippy, { Instance, Props } from 'tippy.js';

const props = defineProps<
  Partial<Props> & { ding?: boolean; triggerEle: HTMLDivElement }
>();

const { triggerEle, placement, trigger, showOnCreate, maxWidth } = props;

let tippyInstance: Instance | null = null;

const contentRef = ref();

const emit = defineEmits(['onShown', 'onHide']);

onMounted(() => {
  tippyInstance = tippy(triggerEle, {
    content: contentRef.value,
    trigger,
    arrow: false,
    theme: 'ref-paper',
    placement,
    maxWidth,
    interactive: true,
    appendTo: props.appendTo || document.body,
    hideOnClick: false,
    showOnCreate,
    zIndex: props.zIndex,
    delay: props.delay,
    offset: props.offset,
    onShown() {
      emit('onShown');
    },
    onHide(instance) {
      emit('onHide');
      instance.popper.removeAttribute('data-x');
      instance.popper.removeAttribute('data-y');
    },

    onClickOutside(instance) {
      /**
       * 这里要用props.ding来使用ding的值，不能直接用解构后的ding
       */
      if (!props.ding) {
        instance?.hide();
      }
    },
  });
});

onUnmounted(() => {
  tippyInstance?.destroy();
});

const show = () => {
  tippyInstance?.show();
};

const hide = () => {
  tippyInstance?.hide();
};

const update = () => {
  tippyInstance?.setContent(contentRef.value);
};

watch(
  () => props.trigger,
  (newVal) => {
    tippyInstance?.setProps({
      trigger: newVal,
    });
  }
);

defineExpose({
  contentRef,
  show,
  update,
  hide,
  getTippyInstance() {
    if (!tippyInstance) {
      throw new Error('tippyInstance is null');
    }
    return tippyInstance;
  },
});
</script>

<style scoped lang="less"></style>
