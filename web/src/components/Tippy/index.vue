<template>
  <div ref="contentRef">
    <slot />
  </div>
</template>
<script lang="ts" setup>
import { onMounted, onUnmounted, ref, PropType, watch } from 'vue';
import tippy, { Instance, Placement } from 'tippy.js';
import { Nullable } from '~/src/typings/global';
import enableTippyDraggable from '~/src/dom/enableTippyDraggable';

const props = defineProps({
  triggerEle: {
    type: HTMLElement,
    default: null,
  },
  placement: {
    type: String as PropType<Placement>,
    default: 'top-end',
  },
  appendToParent: {
    type: Boolean,
    default: false,
  },
  trigger: {
    type: String,
    default: '',
  },
  showOnCreate: {
    type: Boolean,
    default: false,
  },
  maxWidth: {
    type: Number,
    default: 480,
  },
  ding: {
    type: Boolean,
    default: false,
  },
  zIndex: {
    type: Number,
    default: 999,
  },
  delay: {
    type: Array as unknown as () => NonNullable<
      Parameters<typeof tippy>[1]
    >['delay'],
    default: [0, 0],
  },
  offset: {
    type: Array as unknown as () => NonNullable<
      Parameters<typeof tippy>[1]
    >['offset'],
    default: [0, 0],
  },
  disableDraggable: {
    type: Boolean,
    default: false,
  },
});

const { triggerEle, placement, trigger, showOnCreate, maxWidth, ding } = props;

let tippyInstance: Nullable<Instance> = null;

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
    appendTo: props.appendToParent ? 'parent' : document.body,
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
  if (!props.disableDraggable && tippyInstance?.popper) {
    enableTippyDraggable(tippyInstance.popper);
  }
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
});
</script>

<style scoped lang="less"></style>
