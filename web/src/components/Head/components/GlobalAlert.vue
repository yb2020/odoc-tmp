<template>
  <a-alert
    v-if="!!message"
    :message="message"
    type="warning"
    class="text-center"
    closable
    @close="onClose"
  />
</template>
<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue';
import { useLocalStorage } from '@vueuse/core';
import { getGlogalMessage } from '@/api/message';
import { useLanguage } from '@/hooks/useLanguage';
import { MESSAGE } from '@/common/src/constants/storage-keys';

export default defineComponent({
  setup() {
    const showMessage = useLocalStorage(MESSAGE.GLOBAL_ALERT, '');
    const message = ref('');
    const { isEnUS } = useLanguage();
    const isWebEN = isEnUS; // 保持向后兼容的命名
    let lsKey = '';
    onMounted(async () => {
      try {
        const res = await getGlogalMessage() as any;
        if (res && res.show && showMessage.value !== `hidden-${res.lsKey}`) {
          showMessage.value = `show-${res.lsKey}`;
          message.value = res.message?.[isWebEN.value ? 'en' : 'zh'] || '';
          lsKey = res.lsKey;
        }
      } catch (error) {
        console.error('获取全局消息失败:', error);
      }
    });
    return {
      showMessage,
      message,
      onClose: () => {
        showMessage.value = `hidden-${lsKey}`;
      },
    };
  },
});
</script>
<style lang="less" scoped>
.ant-alert-warning {
  background-color: #fffbe6;
  border: 1px solid #ffe58f;
  :deep(.ant-alert-message) {
    color: var(--rp-theme-fg-000000d9, #000000d9);
  }
  :deep(.ant-alert-close-icon .anticon-close) {
    color: var(--rp-theme-fg-000000d9, #000000d9);
  }
}
</style>
