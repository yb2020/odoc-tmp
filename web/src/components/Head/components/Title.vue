<template>
  <div class="head-container">
    <!-- <a-tooltip placement="bottom" :overlayStyle="{ 'max-width': '480px' }"> -->
      <!-- <template #title> -->
        <div class="full-title">
          {{ docNameDecoded }}
        </div>
      <!-- </template> -->
      <!-- <span v-if="!isOwner && isEmptyPaperId(paperId)" class="no-title">{{ defaultTitle }}</span> -->
      <!-- <span v-else @click="goPaperDetails" class="title">{{ docNameDecoded }}</span> -->
    <!-- </a-tooltip> -->

    <!-- <a-spin v-if="isOpenPaper" :spinning="collectLoading" size="small">
      <div v-if="!isCollected" class="collect-btn" @click="collectPaper">
        <i class="aiknowledge-icon icon-bookmark-fill" aria-hidden="true" />
        {{ $t('viewer.save') }}
      </div>
      <div v-else-if="!(docId || selfNoteInfo.userDocId)" class="collect-btn collected-btn" @click="cancelCollectPaper">
        <i class="aiknowledge-icon icon-bookmark-fill" aria-hidden="true" />
        {{ $t('viewer.saved') }}
      </div>
      <FolderMenu v-else :doc-id="docId || selfNoteInfo.userDocId">
        <div class="collect-btn collected-btn" @click="cancelCollectPaper">
          <i class="aiknowledge-icon icon-bookmark-fill" aria-hidden="true" />
          {{ $t('viewer.saved') }}
        </div>
      </FolderMenu>
    </a-spin> -->
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted } from 'vue';
import { selfNoteInfo, isOwner, store } from '~/src/store';
import { message, Modal } from 'ant-design-vue';
import { getDomainOrigin } from '~/src/util/env';
import FolderMenu from '@/components/Folder/FolderMenu.vue';

import { goPathPage } from '~/src/common/src/utils/url';
import { checkOpenPaper, isEmptyPaperId } from '~/src/api/helper';
import { useRoute, useRouter } from 'vue-router';
import { PAGE_ROUTE_NAME } from '~/src/routes/type';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { useI18n } from 'vue-i18n';

import { useCollect } from '~/src/hooks/useCollect';
import { BaseMutationTypes } from '~/src/store/base';

const annotationStore = useAnnotationStore();

const docName = computed(
  () => selfNoteInfo.value?.docName || selfNoteInfo.value?.paperTitle || ''
);
const docNameDecoded = computed(() => decodeURIComponent(docName.value));
const paperId = computed(() => selfNoteInfo.value?.paperId ?? '');
const isOpenPaper = computed(() =>
  checkOpenPaper(paperId.value, selfNoteInfo.value?.isPrivatePaper ?? false)
);

const isCollected = computed(() => !!selfNoteInfo.value?.isCollected);

const { collectPaper, cancelCollectPaper, collectLoading, docId } = useCollect(
  () => paperId.value,
  () => docName.value,
  (collected) => {
    store.commit(`base/${BaseMutationTypes.SET_COLLECTED}`, collected);
  },
  true
);

const goPaperDetails = () => {
  if (!isOpenPaper.value) {
    message.error('该文献暂未入库');
    return;
  }

  goPathPage(`${getDomainOrigin()}/paper/${paperId.value}`);
  return;
};

const openModal = () => {
  Modal.confirm({
    title: t('message.remindToSavePaper.title'),
    content: t('message.remindToSavePaper.content'),
    onOk: async () => {
      window.onbeforeunload = () => {};
      window.open('about:blank', '_self')!.close();
    },
    cancelButtonProps: { type: 'primary' },
    cancelText: t('viewer.save'),
    okText: t('message.remindToSavePaper.cancelText'),
    onCancel: () => {
      window.onbeforeunload = () => {};
      collectPaper();
    },
  });
};

onMounted(() => {
  // window.onbeforeunload = function (e) {
  //   e = e || window.event;
  //   const noteShowModal =
  //     isCollected.value || Object.keys(annotationStore.pageMap).length === 0;
  //   if (!noteShowModal) {
  //     if (e) {
  //       e.returnValue = false;
  //     }
  //     e.preventDefault();
  //     openModal();
  //     return false;
  //   }
  // };
});

const route = useRoute();

const { t } = useI18n();

const router = useRouter();

// 判断当前是否在PDF页面
const isNotePage = computed(() => {
  const route = router.currentRoute.value;
  return route.path.startsWith('/note');
});


const defaultTitle = computed(() => {
  if (route.name === PAGE_ROUTE_NAME.CHAT) {
    return 'ReadPaper chatGPT 领域问答版本';
  }
  if (route.name === PAGE_ROUTE_NAME.WRITE) {
    return 'ReadPaper chatGPT 写论文版本';
  }
  if (isNotePage.value) {
    return t('message.noTitle');
  }
  return '';
});
</script>

<style lang="less" scoped>
.full-title {
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

.head-container {
  display: flex;
  align-items: center;
  z-index: 2;
  cursor: pointer;
  flex: 1;
  width: 100%;
  justify-content: center;
  .no-title,
  .title {
    font-size: 14px;
    line-height: 20px;
    max-width: calc(100vw - 680px);
    text-align: center;
  }
  .no-title {
    color: var(--site-theme-text-tertiary);
  }
  .title {
    font-family: Lato-Bold, Lato;
    font-weight: bold;
    color: var(--site-theme-text-primary);
  }

  .collect-btn {
    width: 94px;
    height: 32px;
    background: var(--site-theme-brand);
    border-radius: 2px;

    display: flex;
    align-items: center;
    justify-content: center;

    font-size: 13px;
    font-weight: 400;
    color: #ffffff;
    line-height: 18px;

    margin-left: 10px;

    .icon-bookmark-fill {
      font-size: 15px;
      color: #fff;
      margin-right: 4px;
      height: 15px;
      line-height: 15px;
    }
  }

  .collected-btn {
    background: transparent;
    color: #b4b9bf;

    .icon-bookmark-fill {
      color: #b4b9bf;
    }
  }
}

.mobile-viewport {
  .head-container {
    .title {
      display: none;
    }

    .collect-btn {
      position: fixed;
      width: 56rpx;
      height: 56rpx;
      background: #ffffff;
      box-shadow: 0px 4rpx 6rpx 0px rgba(0, 0, 0, 0.3);
      bottom: 82rpx;
      right: 12rpx;
      border-radius: 50%;
      color: rgba(51, 51, 51, 1);
      flex-direction: column;

      font-size: 10rpx;
      color: #333333;
      line-height: 14rpx;

      z-index: 9;
      .icon-bookmark-fill {
        color: rgba(51, 51, 51, 1);
        margin-bottom: 3rpx;
        margin-right: 0;
      }
    }

    .collected-btn {
      color: #a19f9d;

      .icon-bookmark-fill {
        color: rgba(161, 159, 157, 1);
      }
    }
  }
}
</style>
