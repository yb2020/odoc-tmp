<template>
  <a-modal
    v-model:visible="showLogin"
    destroy-on-close
    title
    :footer="null"
    :closable="false"
    :width="showSimpleLoginDialog ? 400 : 800"
    class="loginModal"
    :mask-closable="IS_MOBILE ? true : false"
    @cancel="handleCloseModal"
  >
    <div
      v-if="!IS_MOBILE"
      class="viewer-close"
      @click="handleCloseModal"
    >
      <close-outlined style="font-size: 20px" />
    </div>

    <div
      class="login" 
      :style="{
        height: modalHeight,
      }"
    >
      <!-- 注释掉 uni-login iframe 引用 -->
      <!-- 
      <iframe
        ref="iframeRef"
        :src="`${url}/apicdn/readpaper/uni-login/?showCancel=${showCancel}${showSimpleLoginDialog ? '&isShowSimpleLogin=true': ''}`"
        :height="modalHeight"
      />
      -->
      <div class="login-placeholder">
        <p>登录功能暂时不可用</p>
        <a-button type="primary" @click="handleCloseModal">关闭</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { message as Message } from 'ant-design-vue';
import { getDomainOrigin, IS_MOBILE } from '~/src/util/env';
import { CloseOutlined } from '@ant-design/icons-vue';
import { useUserStore } from '~/src/common/src/stores/user';

defineProps<{ showCancel: boolean }>()

const userStore = useUserStore();
const url = ref('');
const iframeRef = ref();

onMounted(() => {
  url.value = getDomainOrigin();

  window.addEventListener(
    'message',
    (event) => {
      const data = event.data;

      if (data.event === 'loginSuccess') {
        Message.success(data.params.message);

        window.location.reload();

        return;
      }

      if (data.event === 'cancel') {
        showLogin.value = false;
        return;
      }
    },
    false
  );
});

const showLogin = computed({
  get() {
    return userStore.loginDialogVisible;
  },
  set(value) {
    if (value) {
      userStore.openLogin();
    } else {
      userStore.closeLogin();
    }
  },
});

// 简单登录对话框标志，在Pinia版本中不存在，使用本地状态代替
const showSimpleLoginDialog = ref(false);

const modalHeight = computed(() => {
  if (IS_MOBILE) return '480px'

  return showSimpleLoginDialog.value ? '504px' : '540px'
})

const handleCloseModal = () => {
  showLogin.value = false

  if (localStorage.getItem('loginSuccessButNoSelectChannel') === 'true') {
    localStorage.setItem('loginSuccessButNoSelectChannel', 'false')

    Message.success('登录成功')

    window.location.reload()
  }
}
</script>

<style lang="less" scoped>

.login {
  overflow: hidden;
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
