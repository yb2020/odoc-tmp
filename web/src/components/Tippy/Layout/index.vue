<template>
  <div
    ref="wrapperRef"
    @click.stop
    :class="`tippy-viewer js-${group || 'default'}-tippy-viewer`"
    :style="realStyle"
  >
    <div class="title">
      <slot
        v-if="hasTitleSlot"
        name="title"
      />
      <div v-else>
        {{ title }}
      </div>
      <div class="close flex items-center">
        <a-tooltip
          v-if="hasLock"
          :title="$t('translate.lockTip')"
        >
          <img
            src="@/assets/images/lock-to-right.svg"
            class="mr-7 w-4 h-4"
            @click="handleLock"
          >
        </a-tooltip>
        <i
          v-if="!noDing"
          :class="[
            'aiknowledge-icon',
            props.isDing?.value ? 'icon-pinned-fill' : 'icon-pin',
          ]"
          aria-hidden="true"
          @click.stop="handleDing"
        />

        <CloseOutlined @click.stop="tippyHandler('close')" />
      </div>
    </div>
    <div class="inner">
      <slot />
    </div>
  </div>
</template>
<script lang="ts" setup>
import { ref, useSlots } from 'vue';
import { CloseOutlined } from '@ant-design/icons-vue';

const props = defineProps<{
  title: string;
  group?: string;
  style: Record<string, string | number>;
  noDing?: boolean;
  hasLock?: boolean;
  isDing?: import('vue').Ref<boolean>; // 新增：从父组件接收响应式钉住状态
  tippyHandler: (event: 'ding' | 'close' | 'unding' | 'lock') => void;
}>();

const hasTitleSlot = !!useSlots().title;

const wrapperRef = ref<HTMLDivElement>();

const realStyle = { ...props.style };

const handleDing = () => {
  // 直接通知父组件切换状态，不维护本地状态
  props.tippyHandler(props.isDing?.value ? 'unding' : 'ding');
};

defineExpose({
  getDing: () => props.isDing?.value || false,
});

const handleLock = () => {
  props.tippyHandler('lock');
};
</script>
<style lang="less" scoped>
.tippy-viewer {
  border-radius: 4px 4px 0px 0px;
  width: 480px;
  font-size: 20px;
  font-family: sans-serif;
  height: auto;
  touch-action: none;
  box-sizing: border-box;
  font-size: 14px;
  line-height: 18px;
  background-color: #fff;
  color: rgba(0, 0, 0, 0.85);
  display: flex;
  flex-direction: column;
  .title {
    z-index: 1;
    padding: 7px 0 7px 16px;
    display: flex;
    justify-content: space-between;
    color: #fff;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    background: #1f71e0;
    border-radius: 4px 4px 0px 0px;

    .close {
      padding: 0 10px;
      cursor: pointer;
      display: flex;
      align-items: center;
      .aiknowledge-icon {
        margin-right: 27px;
        font-size: 16px;
      }
    }
  }

  .inner {
    margin-top: 32px;
    height: calc(100% - 32px);
  }
}
</style>
