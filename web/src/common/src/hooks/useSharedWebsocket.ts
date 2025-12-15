import {
  DetailMessageCount,
  OfficialMessage,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/operator/admin/msg/SystemMsgInfo';
import { Ref, computed, ref, watch } from 'vue';
import { IS_MOBILE as isMobile, isInElectron } from '@common/utils/env';
import { useUserStore } from '../stores/user';
import { createSharedComposable } from '@vueuse/core';

const initSharedWorker = async (uid: string, callback: (x: any) => void) => {
  try {
    const SharedWorker = await import('@okikio/sharedworker');
    const sw = new SharedWorker.SharedWorkerPolyfill(
      window.location.origin + '/rp-worker.js'
    );
    sw.onmessage = (e: MessageEvent<any>) => {
      console.info('[ws]', e.data);
      try {
        const json = JSON.parse(e.data);
        callback(json);
      } catch (error) {}
    };
    sw.onerror = (err: ErrorEvent) => {
      console.error('[ws]', err);
      sw.close();
    };

    sw.postMessage(
      JSON.stringify({
        scene: 'ws-url',
        content: `${
          !['dev', 'uat'].includes(import.meta.env.VITE_API_ENV)
            ? 'wss://conn.readpaper.com'
            : `ws://conn.${import.meta.env.VITE_API_ENV}.aiteam.cc`
        }/ws?uid=${uid}&appId=microService-app-aiKnowledge`,
      })
    );
  } catch (error) {
    console.error(error);
  }
};

export type Msg = {
  scene: string | 'user_center' | 'message_center' | 'OFFICIAL_MESSAGE';
  content: number;
  officialMessage?: OfficialMessage;
  detailMessageCounts?: DetailMessageCount[];
};

function useNotificationWebsocket(uid?: Ref<string>) {
  const userStore = useUserStore();
  const isEnabled = computed(
    () => !isInElectron() && !isMobile && (uid?.value || userStore.isLogin())
  );

  const sw = ref();
  const msg = ref<Msg>();
  const update = (data: Msg) => {
    msg.value = data;
  };

  watch(
    isEnabled,
    async () => {
      if (isEnabled.value) {
        const id = uid?.value || userStore.userInfo?.id || '';
        sw.value = await initSharedWorker(id, update);
      } else {
        sw.value?.close();
      }
    },
    {
      immediate: true,
    }
  );

  return {
    msg,
  };
}

export default createSharedComposable(useNotificationWebsocket);
