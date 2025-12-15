<template>
  <teleport to="body">
    <div
      v-if="visible"
      class="cert-wrapper"
    >
      <iframe
        class="cert-iframe"
        title="个人信息认证"
        :src="`${origin}/dialog/personal?${qs}`"
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
import api from '@common/api/axios';
import { register, useCertStore, CertProps } from '@common/stores/cert';
import { getHostname, isElectronMode } from '@common/utils/env';
import { Modal, message } from 'ant-design-vue';
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { pageType } = defineProps<{ pageType: string }>();

const { t } = useI18n();
const loading = ref(true);
const error = ref();
const store = useCertStore();
const visible = computed(() => store.visible);
const sceneId = computed(() => store.sceneId);
const qs = computed(() => {
  const params = new URLSearchParams({
    pageType,
    sceneId: sceneId.value,
  });
  return params.toString();
});

const origin = isElectronMode() ? `https://${getHostname()}` : '';

// watch(visible, (cur) => {
//   if (cur) {
//     setTimeout(onSucc, 5000);
//   }
// });

register(api, (props?: CertProps) => {
  return new Promise((resolve) => {
    if (!visible.value) {
      loading.value = true;
      error.value = null;
    }
    store.showCertDialog({
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
  if (!isElectronMode() && doc?.title !== 'ReadPaper Certification') {
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
  store.hideCertDialog({
    status: 0,
  });
};

const onCancel = () => {
  store.hideCertDialog({
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
