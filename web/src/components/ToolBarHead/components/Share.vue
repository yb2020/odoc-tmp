<template>
  <!-- <a-popconfirm
      title="笔记分享功能内测中"
      placement="bottomLeft"
      ok-text="我知道了"
      :showCancel="false"
      :visible="tipsVisible"
      @confirm="hideTips()"
    >
      
    </a-popconfirm> -->
  <div
    class="text-base cursor-pointer"
    @click="visible = true"
  >
    <a-tooltip :title="$t('toolbar.share')">
      <i class="aiknowledge-icon icon-share-square-fill" />
      <!-- {{ $t('viewer.share') }} -->
    </a-tooltip>
  </div>
  <Modal
    v-model:visible="visible"
    destroy-on-close
    :footer="null"
    :closable="false"
    :width="344"
    wrap-class-name="share-modal-wrap"
  >
    <div class="share-modal-content">
      <div class="share-modal-title">
        <span>{{ $t('viewer.share') }}</span>
        <CloseOutlined @click="visible = false" />
      </div>
      <!-- <div class="share-modal-tips">
        {{ $t('viewer.shareDialog.title') }}
      </div> -->
      <!-- <div class="share-modal-switch">
        <div>
          {{ $t('viewer.shareDialog.public') }}
          <Tooltip
            :title="$t('viewer.shareDialog.publicTip')"
            placement="topLeft"
            overlay-class-name="metadata-rollback-tooltip"
          >
            <QuestionCircleOutlined />
          </Tooltip>
          ：
        </div>
        <Switch
          :checked="pdfStatusInfo.noteOpenAccessFlag"
          :style="{
            cursor: loading ? 'wait' : 'pointer',
          }"
          @click="shareSwitch()"
        />
      </div> -->
      <div class="share-modal-link">
        <div>{{ $t('viewer.shareDialog.link') }}</div>
      </div>
      <div class="share-modal-href">
        {{ getHref() }}
      </div>
      <div
        v-if="!pdfStatusInfo.authPdfId"
        class="share-modal-pdf"
      >
        <ExclamationCircleFilled />
        {{ $t('viewer.shareDialog.privateTip') }}
      </div>
      <!-- <div
        class="share-modal-copy"
        :class="{
          disabled: !pdfStatusInfo.noteOpenAccessFlag,
        }"
        @click="pdfStatusInfo.noteOpenAccessFlag && copyLink()"
      >
        {{ $t('viewer.shareDialog.copyLink') }}
      </div> -->
      <div
        class="share-modal-copy"
        :class="{
          disabled: false,
        }"
        @click="copyLink()"
      >
        {{ $t('viewer.shareDialog.copyLink') }}
      </div>
    </div>
  </Modal>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { Modal, Tooltip, Switch, message } from 'ant-design-vue';
import {
  QuestionCircleOutlined,
  ExclamationCircleFilled,
  CloseOutlined,
} from '@ant-design/icons-vue';
import copyTextToClipboard from 'copy-text-to-clipboard';
import { store, pdfStatusInfo, selfNoteInfo } from '~/src/store';
import { closeAccess, openAccess } from '~/src/api/base';
import { BaseMutationTypes } from '~/src/store/base';
import { GetPdfStatusInfoResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';
// import { useNewFeature } from '~/src/hooks/useNewFeature';
import { ElementClick, PageType, reportElementClick } from '~/src/api/report';
import { getPdfAnnotateNoteUrl } from '~/src/util/env';

const getHref = (): string => {
  const { searchParams } = new URL(window.location.href);

  return String(
    getPdfAnnotateNoteUrl({
      pdfId: selfNoteInfo.value?.pdfId ?? searchParams.get('pdfId') ?? '',
      noteId: selfNoteInfo.value?.noteId ?? searchParams.get('noteId') ?? '',
    })
  );
};

// const { tipsVisible, hideTips } = useNewFeature('share-note', 2022, 12, 31);

const visible = ref(false);
const loading = ref(false);
const shareSwitch = async () => {
  if (loading.value) {
    return;
  }

  loading.value = true;

  const params = {
    noteId: selfNoteInfo.value!.noteId,
  };
  const request = pdfStatusInfo.value.noteOpenAccessFlag
    ? closeAccess(params)
    : openAccess(params);

  try {
    await request;
  } catch (error) {
    const operation = pdfStatusInfo.value.noteOpenAccessFlag
      ? '关闭分享笔记'
      : '公开笔记';
    message.error(operation + '失败，请稍后再试');
    loading.value = false;
    return;
  }

  store.commit(`base/${BaseMutationTypes.SET_STATUS_INFO}`, {
    ...pdfStatusInfo.value,
    noteOpenAccessFlag: !pdfStatusInfo.value.noteOpenAccessFlag,
  } as GetPdfStatusInfoResponse);

  loading.value = false;
};

const copyLink = () => {
  const success = copyTextToClipboard(getHref());
  if (success) {
    message.success('已复制到剪贴板');
  } else {
    message.warn('复制失败，请手动复制以上链接');
  }

  reportElementClick({
    page_type: PageType.note,
    type_parameter: selfNoteInfo.value?.pdfId ?? '',
    element_name: ElementClick.copy_note_link,
  });
};
</script>

<style lang="less">
.share-modal-wrap {
  .ant-modal-content {
    box-shadow:
      0 3px 6px -4px rgb(0 0 0 / 12%),
      0 6px 16px 0 rgb(0 0 0 / 8%),
      0 9px 28px 8px rgb(0 0 0 / 5%);
  }
  .ant-modal-body {
    padding-top: 14px;
    padding-bottom: 16px;
    padding-left: 16px;
    padding-right: 16px;
    background-color: white !important;
  }

  .ant-switch {
    &[aria-checked='false'] {
      background-color: rgba(0, 0, 0, 0.25);
    }
    .ant-switch-handle::before {
      box-shadow: 0 2px 4px 0 rgb(0 35 11 / 20%);
    }
  }
}

.share-modal-content {
  color: #1d2229;
  font-size: 14px;
  .share-modal-title {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 26px;
    > span:first-of-type {
      font-size: 16px;
      font-weight: bold;
    }
    > *:last-of-type {
      color: rgba(0, 0, 0, 0.45);
    }
  }
  .share-modal-tips {
    margin-top: 24px;
    height: 22px;
  }
  .share-modal-switch,
  .share-modal-link {
    margin-top: 16px;
    height: 22px;
    position: relative;
    > div:first-of-type {
      width: 70px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
  .share-modal-switch {
    display: flex;
    align-items: center;
    > div:first-of-type {
      > i {
        position: absolute;
        left: 34px;
      }
    }
    > button {
      margin-left: 24px;
    }
  }
  .share-modal-href {
    margin-top: 8px;
    padding: 8px;
    word-break: break-all;
    background: #f7f8fa;
  }
  .share-modal-pdf {
    margin-top: 16px;
    line-height: 18px;
    color: #86919c;
    font-size: 12px;
  }
  .share-modal-copy {
    margin-top: 16px;
    color: white;
    background: #1f71e0;
    border-radius: 2px;
    height: 40px;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    &:hover {
      opacity: 0.85;
    }
    &.disabled {
      cursor: not-allowed;
      opacity: 0.5;
    }
  }
}
</style>
