<template>
  <div class="private-container">
    <div class="private-page">
      <div class="mask">
        <img class="icon" :src="Copyright" alt="icon" />
        <div class="title">该PDF属于有版权争议的PDF，无法公开查看</div>
        <div class="detail" v-if="selfNoteInfo">
          当前上传的版本属于【{{
            selfNoteInfo?.userInfo.showName || selfNoteInfo?.userInfo.nickName
          }}】，仅用户本人可以查看以及记录笔记
        </div>
        <div class="detail">您可以通过上传同一篇论文的PDF来获取阅读权限</div>
        <div v-if="IS_MOBILE" class="index" @click="handleGoIndex">
          去首页逛逛
          <ArrowRightOutlined />
        </div>
        <div v-else class="upload">
          <a-button
            type="primary"
            ref="uploadBtnRef"
            v-if="
              !pdfStatusInfo.hasPdfAccessFlag &&
              pdfStatusInfo.noteUserStatus === UserStatusEnum.TOURIST
            "
            @click="handleClick"
          >
            上传论文
          </a-button>

          <template v-else>
            <a-upload
              name="file"
              :multiple="false"
              :customRequest="handleUpload"
              :showUploadList="false"
            >
              <a-button type="primary" ref="uploadBtnRef">上传论文</a-button>
            </a-upload>
            <a-progress
              v-show="visible"
              :percent="percent"
              :status="progressStatus"
            />
          </template>
        </div>
      </div>
    </div>
    <NoteInfo />
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { store, pdfStatusInfo, selfNoteInfo, isOwner } from '@/store';
import Copyright from '@/assets/images/annotate-view/copyright.png';
import { ArrowRightOutlined } from '@ant-design/icons-vue';
import NoteInfo from './NoteInfo.vue';
import {
  isErrorStatus,
  createItem,
  Item,
  uploadItem,
} from '@/utils/pdf-upload/index.js';
import api from '~/src/api/axios';
import {
  CreateDocScene,
  UserDocCreateStatus,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/userCenter/UserDoc';
import { UserStatusEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';
import { Modal, ModalFuncProps } from 'ant-design-vue';
import {
  ElementName,
  PageType,
  reportPdfUploadSuccess,
} from '~/src/api/report';
import { getDomainOrigin, IS_MOBILE } from '~/src/util/env';
import { getHostname, isInElectron } from '../../util/env';
import { goPathPage } from '~/src/common/src/utils/url';
import { isEmptyPaperId } from '~/src/api/helper';
import { useVipStore } from '@common/stores/vip';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '~/src/common/src/stores/user';

const handleGoIndex = () => {
  goPathPage(getDomainOrigin());
};

const vipStore = useVipStore();

const $i18n = useI18n();

const visible = ref(false);

const uploadBtnRef = ref();

const uploading = ref<Item | null>(null);
let cancelItem: Awaited<ReturnType<typeof uploadItem>>['cancelItem'] | null =
  null;

const handleUpload = async (info: { file: File & { uid: string } }) => {
  if (uploading.value) {
    uploading.value = null;
    if (cancelItem) {
      cancelItem();
      cancelItem = null;
    }
  }

  const paperId = isEmptyPaperId(store.state.base.paperId)
    ? ''
    : store.state.base.paperId;

  uploading.value = createItem(
    {
      groupId: '',
      classifyId: '',
      folderId: '',
      paperId,
      scene: paperId
        ? CreateDocScene.SCENE_VERIFY
        : CreateDocScene.SCENE_DEFAULT,
    },
    {}
  );
  const result = await uploadItem(uploading.value, info.file, () => {});
  cancelItem = result.cancelItem;
  visible.value = true;

  await result.promise;

  const viewMyNote = () => {
    window.location.replace(
      `${getDomainOrigin()}/pdf-annotate/note?pdfId=${uploading.value?.docInfo
        ?.pdfId}`
    );
  };

  const successModal: ModalFuncProps = {
    title: '解锁成功',
    content: '刷新页面可以继续阅读（该文献已自动帮您上传到个人中心）',
    okText: '刷新页面',
  };

  const failedModal: ModalFuncProps = {
    title: '解锁失败',
    content:
      '上传的文献与当前浏览文献不匹配（上传的文献已自动绑定帮您上传到 个人中心）',
    cancelText: '取消',
    okText: '查看上传文献',
    onOk: viewMyNote,
  };

  const retryModal: ModalFuncProps = {
    title: '上传失败',
    content: uploading.value?.message || '',
    onOk() {
      uploadBtnRef.value.$el.click();
    },
    okText: '重新上传',
    maskClosable: true,
  };

  if (uploading.value.status === 'limited') {
    // 达到上限
    vipStore.showVipLimitDialog(uploading.value.message || '达到上传上限', {
      exception: uploading.value.exception,
      leftBtn: {
        text: $i18n.t('common.premium.btns.more'),
        url: `${getDomainOrigin()}/home/mine`,
      },
      reportParams: {
        page_type: PageType.note,
        element_name: ElementName.upperCollectionPopup,
      },
    });
  } else if (isErrorStatus(uploading.value.status)) {
    Modal.error(retryModal);
  } else if (
    uploading.value?.status === UserDocCreateStatus.SELF_OTHER_VERSION
  ) {
    Modal.confirm(failedModal);
  } else if (
    uploading.value?.status === UserDocCreateStatus.FINISH ||
    uploading.value?.status === 'repeat'
  ) {
    if (isOwner.value) {
      Modal.success({
        ...successModal,
        onOk: viewMyNote,
      });
    } else {
      await store.dispatch(`base/${BaseActionTypes.GET_STATUSINFO}`, {
        noteId: selfNoteInfo.value!.noteId,
        pdfId: selfNoteInfo.value!.pdfId,
      });
      if (pdfStatusInfo.value.hasPdfAccessFlag) {
        Modal.success({
          ...successModal,
          onOk() {
            window.location.reload();
          },
        });
      } else {
        Modal.confirm(failedModal);
      }
    }
  }
  visible.value = false;

  if (
    [
      UserDocCreateStatus.FINISH,
      UserDocCreateStatus.CONFLICT,
      UserDocCreateStatus.GROUP_OTHER_VERSION,
      UserDocCreateStatus.SELF_OTHER_VERSION,
    ].includes(uploading.value.status as UserDocCreateStatus)
  ) {
    reportPdfUploadSuccess(uploading.value.docInfo?.pdfId || '');
  }
};

const percent = computed(() => {
  if (!uploading.value) {
    return 0;
  }

  if (uploading.value.status === 'waiting') {
    return 20;
  }

  if (uploading.value.status === UserDocCreateStatus.UPLOAD) {
    return 40;
  }

  if (uploading.value.status === UserDocCreateStatus.PARSE) {
    return 60;
  }

  if (uploading.value.status === UserDocCreateStatus.MATCH) {
    return 80;
  }

  if (uploading.value.status === UserDocCreateStatus.FINISH) {
    return 100;
  }

  return 0;
});

const progressStatus = computed(() => {
  if (uploading.value) {
    if (isErrorStatus(uploading.value.status)) {
      return 'exception';
    } else if (uploading.value.status === UserDocCreateStatus.FINISH) {
      return 'success';
    }
  }

  return 'active';
});

const userStore = useUserStore();

const handleClick = () => {
  if (!userStore.isLogin()) {
    userStore.openLogin();
    return;
  }
}
</script>

<style lang="less" scoped>
.private-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.private-page {
  position: relative;
  width: 50%;
  height: 100%;
  .mask {
    background: var(--site-theme-bg-primary);

    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    width: 100%;
    height: 100%;
  }

  .icon {
    width: 64px;
    height: 64px;
  }

  .title {
    font-size: 18px;
    margin-top: 34px;
    font-size: 18px;
    font-weight: 600;
    color: var(--site-theme-text-primary);
    line-height: 25px;
  }

  .detail {
    margin-top: 10px;
    width: 434px;

    font-size: 14px;
    font-weight: 600;
    color: var(--site-theme-text-secondary);
    line-height: 20px;
    text-align: center;
  }

  .upload {
    margin-top: 32px;
    width: 100%;
    text-align: center;
    padding: 0 20px;
    :deep(.ant-progress) {
      margin-top: 16px;
    }
    :deep(.ant-progress-inner) {
      background-color: var(--site-theme-bg-soft);
    }
    :deep(.ant-progress-text) {
      color: var(--site-theme-text-primary);
    }
  }
}

.progress-container {
  display: flex;
  justify-content: center;
}

.mobile-viewport {
  .private-page {
    width: 100%;
    .mask {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
    }

    .icon {
      width: 64rpx;
      height: 64rpx;
      margin-top: 0;
    }

    .title {
      font-size: 16rpx;
      line-height: 26rpx;
    }

    .detail {
      font-size: 13rpx;
      color: var(--site-theme-text-secondary);
      line-height: 22rpx;
      width: 209rpx;
    }

    .index {
      font-size: 14rpx;
      font-family:
        PingFangSC-Semibold,
        PingFang SC;
      font-weight: 600;
      color: var(--site-theme-primary);
      margin-left: 10rpx;
      margin-top: 24rpx;
    }
  }
}
</style>
