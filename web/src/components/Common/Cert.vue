<template>
  <teleport to="body">
    <div
      v-if="visible"
      class="cert-wrapper"
    >
      <iframe
        class="cert-iframe"
        title="个人信息认证"
        :src="`${origin}/dialog/personal?pageType=${PageType.note}&sceneId=${sceneId}`"
        :style="{ display: error || loading ? 'none' : 'block' }"
        @load="onLoad"
      />
      <a-spin
        v-if="loading"
        class="cert-loading"
      />
    </div>
  </teleport>
</template>
<script lang="ts" setup>
import api from '@/api/axios';
import { register } from '@/store/cert';
import { PageType } from '@/api/report';
import { Modal, message } from 'ant-design-vue';
import { computed, onMounted, ref } from 'vue';
import { useStore } from '../../store';
import { CertProps } from '../../store/cert';
import { useI18n } from 'vue-i18n';
import { getHostname, isInElectron } from '../../util/env';

const { t } = useI18n();
const store = useStore();
const loading = ref(true);
const error = ref();
const visible = computed(() => store.state.cert.visible);
const sceneId = computed(() => store.state.cert.sceneId);

const origin = isInElectron() ? `https://${getHostname()}` : '';

// watch(visible, (cur) => {
//   if (cur) {
//     setTimeout(onSucc, 5000);
//   }
// });

register(api, (props?: CertProps) => {
  return new Promise((resolve) => {
    const { result } = store.state.cert;
    if (result) {
      return resolve(result);
    }

    if (!visible.value) {
      loading.value = true;
      error.value = null;
    }
    store.commit('cert/SHOW_CERT_DIALOG', {
      sceneId: props?.sceneId,
      callback: resolve,
    });
  });
});

const onLoad = (e: Event) => {
  loading.value = false;
  const iframe = e.target as HTMLIFrameElement;
  const doc = iframe.contentDocument;
  // Electron内跨域了
  if (!isInElectron() && doc?.title !== 'ReadPaper Certification') {
    error.value = new Error('Failed to load');
    Modal.error({
      title: t('viewer.certPageLoadFail'),
      content: error.value.message || error.value.name,
      onOk: onCancel,
      onCancel,
      cancelText: '',
      zIndex: 10001,
    });
  }
};

const onSucc = () => {
  store.commit('cert/HIDE_CERT_DIALOG', {
    status: 0,
  });
};

const onCancel = () => {
  store.commit('cert/HIDE_CERT_DIALOG', {
    status: 1,
  });
};

const onError = (msg: string) => {
  message.error(msg);
};

onMounted(() => {
  const handler = (e: MessageEvent) => {
    const { event } = e.data || {};
    if (event === 'succ') {
      onSucc();
    } else if (event === 'cancel') {
      onCancel();
    } else if (event === 'error') {
      onError(e.data);
    }
  };
  window.addEventListener('message', handler);
  return () => {
    window.removeEventListener('message', handler);
  };
});
</script>

<style lang="less" scoped>
.cert-wrapper {
  position: fixed;
  // 得比 tippy 组件高
  z-index: 10000;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background: transparent;
}
.cert-iframe {
  border: none;
  width: 100%;
  height: 100%;
  user-select: none;
}
.cert-loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>
