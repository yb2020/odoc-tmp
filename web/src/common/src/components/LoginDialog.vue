<template>
  <a-modal
    v-model:visible="visible"
    destroy-on-close
    title
    :footer="null"
    :closable="false"
    :width="800"
    class="loginModal"
    :mask-closable="IS_MOBILE ? true : false"
    @cancel="closeLoginDialog"
  >
    <div
      v-if="!IS_MOBILE"
      class="viewer-close"
      @click="closeLoginDialog"
    >
      <close-outlined style="font-size: 20px" />
    </div>
    <div class="login">
      <!-- 注释掉 uni-login iframe 引用 -->
      <!-- 
      <iframe
        ref="iframeRef"
        :src="`${url}/apicdn/readpaper/uni-login/?showCancel=${showCancel}${
          isWebEN ? '&lang=web_en' : ''
        }`"
        :height="IS_MOBILE ? 480 : 540"
      />
      -->
      <div class="login-placeholder">
        <p>登录功能暂时不可用</p>
        <a-button type="primary" @click="closeLoginDialog">关闭</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { message as Message } from 'ant-design-vue';
import { getDomainOrigin, IS_MOBILE } from '@common/utils/env';
import { CloseOutlined } from '@ant-design/icons-vue';
import { useUserStore } from '@common/stores/user';
import { computed } from 'vue';
import { onUnmounted } from 'vue';
import { watch } from 'vue';

defineProps<{ showCancel: boolean; isWebEN?: boolean }>();

const userStore = useUserStore();

const url = ref('');

const iframeRef = ref();

const messageHandler = (event: MessageEvent) => {
  const data = event.data as { event: string; params: { message: string } };

  if (data.event === 'loginSuccess') {
    Message.success(data.params.message);
    userStore.getUserInfo().then(() => {
      userStore.closeLogin();
      // 如果需要刷新页面，可以在获取用户信息后刷新
      // window.location.reload();
    });
    return;
  }

  if (data.event === 'cancel') {
    closeLoginDialog();
    return;
  }
};

onMounted(() => {
  url.value = getDomainOrigin();

  window.addEventListener('message', messageHandler, false);
});

onUnmounted(() => {
  window.removeEventListener('message', messageHandler, false);
});

const visible = ref(false);

watch(
  () => userStore.loginDialogVisible,
  (val) => {
    visible.value = val;
  },
  {
    immediate: true,
  }
);

const closeLoginDialog = () => {
  userStore.closeLogin();
};
</script>
<style lang="less" scoped>
.login {
  overflow: hidden;
  height: 540px;
  margin: -24px;
  background-color: #fff;
  iframe {
    border: none;
    width: 100%;
  }
}
.viewer-close {
  color: #fff;
  position: absolute;
  right: -48px;
  top: 0;
  width: 48px;
  height: 48px;
  font-size: 24px;
  line-height: 40px;
  text-align: center;
  cursor: pointer;
}
.mobile-viewport {
  .login {
    height: 480px;
  }
}
.login-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 20px;
  
  p {
    font-size: 16px;
    color: #666;
  }
}
</style>
