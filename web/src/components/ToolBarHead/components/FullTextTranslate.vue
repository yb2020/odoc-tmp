<template>
  <a-dropdown
    destroyPopupOnHide
    :trigger="['hover']"
  >
    <template
      v-if="
        data.status !== FullTranslateFlowStatus.TRANSLATE_FINISHED 
        || fullTextTranslateStore.pdfId
      "
      #overlay
    >
      <a-menu>
        <FullTextTranslatePopover
          :status="data.status"
          :progressPercent="data.progressPercent || ''"
          :pdf-info="selfPdfInfo"
          :translate-loading="translateLoading"
          @translate="startTranslate"
        />
      </a-menu>
    </template>
    <button
      class="relative h-8 !p-2 rounded-sm bg-transparent border border-solid cursor-pointer full-text-translate-btn"
      :class="[
        finished ? 'finished' : 'unfinished',
        {
          'active': !closed,
        },
      ]"
      @click="handleFullTextTranslate"
    >
      <span>{{ $t('viewer.documentTranslate') }}</span>
      <span
        v-if="status === FullTranslateFlowStatus.TRANSLATE_FINISHED"
        class="absolute h-0.5 bottom-0 left-0 progress-bar"
        :style="{
          width: `${data?.progressPercent ?? 1}%`,
        }"
      />
    </button>
  </a-dropdown>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue';
import { message } from 'ant-design-vue';
import {
  FullTranslateFlowStatus,
  GetTranslateStatusReq,
  GetTranslateStatusResponse,
} from 'go-sea-proto/gen/ts/translate/FullTextTranslate';
import { useFullTextTranslateStore } from '~/src/stores/fullTextTranslateStore';
import {
  getTranslateStatus,
  postTranslate,
} from '@/api/fullTextTranslate';
import FullTextTranslatePopover from './FullTextTranslatePopover.vue';
import { currentNoteInfo, selfNoteInfo } from '~/src/store';
import { polling } from '@/util/polling';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@common/stores/user';


const { t } = useI18n();

const props = defineProps<{
  pdfViewFinished?: boolean;
}>();
const fullTextTranslateStore = useFullTextTranslateStore();
const userStore = useUserStore();

const fetchTranslateStatus = (params: GetTranslateStatusReq) => {
  return polling<GetTranslateStatusReq, GetTranslateStatusResponse>({
    fn: getTranslateStatus,
    maxAttempts: 30,
    interval: 5000,
    params,
    validate: (res) => {
      data.value = res;
      if (res.status === FullTranslateFlowStatus.TRANSLATING) {
        return false;
      }
      return true;
    },
  });
};

// let fetchHandler: null | { cancel: () => void } = null;

const data = ref<GetTranslateStatusResponse>({
  status: FullTranslateFlowStatus.WITHOUT_TRANSLATE_HISTORY,
});
const status = computed(() => data.value.status);
const finished = computed(() => data.value.status === FullTranslateFlowStatus.TRANSLATE_FINISHED && data.value.translationFileUrl ); // 确保有翻译文件URL
const closed = computed(() => !fullTextTranslateStore.pdfId);

const selfPdfInfo = computed(() => {
  return {
    pdfId: selfNoteInfo.value?.pdfId || '',
    pdfUrl: selfNoteInfo.value?.pdfUrl || '',
  };
});

const route = useRoute();
const startToPolling = async (
  pdfInfo: { pdfId: string; pdfUrl: string },
  noErrorTip: boolean
) => {
  const res = fetchTranslateStatus({
    pdfId: pdfInfo.pdfId,
    needTranslateFileUrl: pdfInfo.pdfUrl,
  });
  // fetchHandler = {
  //   cancel: res.cancel,
  // };
  try {
    const result = await res.request;
    data.value = result;
    if (result.status === FullTranslateFlowStatus.TRANSLATE_FINISHED) {
      fullTextTranslateStore.setFullTextTranslatePDFUrl(
        result.translationFileUrl!
      );
      if (result.alignment) {
        fullTextTranslateStore.setFullTextTranslateAlignment(result.alignment);
      }
      if (route.query.fullText === 'open') {
        // 打开全文翻译
        handleFullTextTranslate();
      }
    } else if (result.status === FullTranslateFlowStatus.TRANSLATE_FAIL && !noErrorTip) {
      message.error(t('viewer.fullTextTranslateError.msg'));
      //刷新用户积分
      userStore.refreshUserCredits()
    }
  } catch (error) {
    data.value = {
      status: FullTranslateFlowStatus.WITHOUT_TRANSLATE_HISTORY,
    };
    //刷新用户积分
    userStore.refreshUserCredits()
  }
};

startToPolling(selfPdfInfo.value, true);

watch(
  () => props.pdfViewFinished,
  () => {
    if (props.pdfViewFinished && route.query.fullText === 'open') {
      // 打开全文翻译
      handleFullTextTranslate();
    }
  },
  { immediate: true }
);

// watch(currentPdfInfo, (newVal, oldVal) => {
//   if (isEqual(newVal, oldVal)) {
//     return
//   }
//   // 取消轮询
//   fetchHandler?.cancel()
//   fullTextTranslateStore.setFullTextTranslatePDFUrl('')
//   fullTextTranslateStore.toggleFullTextTranslate('')
//   if (newVal.pdfId && newVal.pdfUrl) {
//     currentFullTextStatus.value = {
//       status: TranslateStatus.unknown,
//     }
//     startToPolling(newVal)
//   }
// })

const handleFullTextTranslate = async () => {
  if (data.value.status === FullTranslateFlowStatus.TRANSLATE_FINISHED) {
    // 在线查看翻译PDF前，重新获取最新的翻译文件URL（链接可能过期）
    try {
      const pdfId = currentNoteInfo.value?.pdfId || '';
      const pdfUrl = currentNoteInfo.value?.pdfUrl || '';

      if (pdfId && pdfUrl) {
        console.log('获取最新翻译PDF地址...');
        const latest = await getTranslateStatus({
          pdfId,
          needTranslateFileUrl: pdfUrl,
        });
        
        if (latest?.translationFileUrl) {
          console.log('更新翻译PDF地址:', latest.translationFileUrl);
          fullTextTranslateStore.setFullTextTranslatePDFUrl(latest.translationFileUrl);
        }
      }
    } catch (e) {
      console.error('获取最新翻译文件链接失败：', e);
      // 如果获取失败，仍然使用现有的URL继续
    }

    // 切换全文翻译显示状态
    fullTextTranslateStore.toggleFullTextTranslate(
      currentNoteInfo.value?.pdfId
    );
  }
};

const translateLoading = ref(false);

const startTranslate = async () => {
  translateLoading.value = true;

  data.value.status = FullTranslateFlowStatus.TRANSLATING;
  try {
    const postTranslateData = await postTranslate({
      pdfId: selfPdfInfo.value.pdfId,
      needTranslateFileUrl: selfPdfInfo.value.pdfUrl,
    });

    if (postTranslateData.errorCode === '0') {
       
    } else {
      data.value.status = FullTranslateFlowStatus.TRANSLATE_FAIL;
      translateLoading.value = false;
      message.error(postTranslateData.message || '全文翻译失败');
      return;
    }
    // 翻译请求不管成功失败，都需要刷新用户积分
    userStore.refreshUserCredits()
    setTimeout(() => {
      startToPolling(selfPdfInfo.value, false);
    }, 5000);
  } catch (error) {
    data.value.status = FullTranslateFlowStatus.TRANSLATE_FAIL;
  }
  translateLoading.value = false;
};
</script>

<style scoped>
.full-text-translate-btn {
  color: var(--site-theme-text-color, #000000);
  border-color: var(--site-theme-divider, #d9d9d9);
  background-color: transparent;
}

.full-text-translate-btn:hover {
  background-color: var(--site-theme-background-hover, rgba(0, 0, 0, 0.05)) !important;
}

.full-text-translate-btn.finished {
  border-color: var(--site-theme-primary-color, #52c41a) !important;
}

.full-text-translate-btn.unfinished {
  border-color: var(--site-theme-divider, #d9d9d9) !important;
}

.full-text-translate-btn.active {
  background-color: var(--site-theme-background-hover, rgba(0, 0, 0, 0.05)) !important;
}

.progress-bar {
  background-color: var(--site-theme-primary-color, #1890ff);
}
</style>
