<template>
  <a-popover
    :visible="true"
    :overlay-style="{
      display: visibility === 'visible' && msgCount > 0 ? 'block' : 'none',
    }"
    overlay-class-name="official-popup"
    :align="
      align || {
        points: ['tr', 'br'],
        offset: [56, 8],
      }
    "
  >
    <template #content>
      <iframe
        v-if="isLogin"
        ref="msgFrame"
        :src="`https://${getWebHost()}/dialog/official?ws=1`"
        width="360"
        :height="height"
        class="official-iframe"
      />
    </template>
    <slot />
  </a-popover>
</template>

<script lang="ts" setup>
import { useDocumentVisibility } from '@vueuse/core';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useUserStore } from '@common/stores/user';
import { getWebHost } from '../utils/env';

const props = defineProps<{
  align?: object;
  logined?: boolean;
}>();

const visibility = useDocumentVisibility();
const msgCount = ref(0);
const msgFrame = ref<HTMLIFrameElement>();
const height = ref(0);

const userStore = useUserStore();
const isLogin = computed(() => userStore.isLogin() || props.logined);

const calHeight = () => {
  if (!msgFrame.value) {
    return;
  }
  const { contentDocument } = msgFrame.value;
  const { scrollHeight = 0 } = contentDocument?.body || {};
  console.debug('[official iframe]scroll height', scrollHeight);
  if (scrollHeight !== 0) {
    height.value = scrollHeight;
  } else {
    // 重复计算到有效为止
    setTimeout(calHeight, 16);
  }
};

const onMessage = (e: MessageEvent) => {
  if (e.data.event === 'official_msg') {
    msgCount.value = e.data.params?.msgCount || 0;
    if (msgCount.value > 0) {
      calHeight();
    }
  }
};

onMounted(() => {
  window.addEventListener('message', onMessage);
});
onUnmounted(() => {
  window.removeEventListener('message', onMessage);
});
</script>

<style lang="less">
.official-popup {
  .ant-popover-content {
    .ant-popover-arrow {
      // 24px是右侧间距
      right: calc(8px + 24px + 32px);
    }
    .ant-popover-arrow-content {
      background-color: rgba(255, 255, 255, 0.85);
    }

    .ant-popover-inner .ant-popover-inner-content {
      padding: 0 !important;
      display: flex;
    }
  }
}
.official-iframe {
  border: 0;
  min-height: 160px;
}
</style>
