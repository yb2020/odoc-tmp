import { ref } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { delay } from '@idea/aiknowledge-special-util/delay';
import { bridgeAdaptor } from '@/adaptor/bridge';
import { postCollectPaper, postCancelCollectPaper } from '~/src/api/material';
import { ResponseError } from '~/src/api/type';
import { useVipStore } from '@common/stores/vip';
import { useAttachmentsStore } from '~/src/stores/attachmentsStore';
import {
  reportElementClick,
  ElementName,
  ElementClick,
  PageType,
  getPdfIdFromUrl,
} from '~/src/api/report';
import { selfNoteInfo, store } from '../store';
import { NeedVipException } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { ERROR_CODE_NEED_VIP } from '@common/api/const';
import { useI18n } from 'vue-i18n';
import { getDomainOrigin } from '@common/utils/env';

export const useCollect = (
  getPaperId: () => string,
  getDocName: () => string,
  setCollected: (collected: boolean) => void,
  report: boolean
) => {
  const attachmentsStore = useAttachmentsStore();
  const collectLoading = ref(false);
  const docId = ref('');
  const $i18n = useI18n();

  const preProcess = (): boolean => {
    if (report) {
      reportElementClick({
        page_type: PageType.note,
        type_parameter: selfNoteInfo.value?.pdfId ?? getPdfIdFromUrl(),
        element_name: ElementClick.collect,
      });
    }

    if (!store.getters['user/isLogin']) {
      bridgeAdaptor.login();
      return false;
    }

    if (collectLoading.value) {
      return false;
    }

    collectLoading.value = true;

    return true;
  };

  const postProcess = async () => {
    await delay(300);
    collectLoading.value = false;
  };

  const collectPaper = async () => {
    if (!preProcess()) {
      return false;
    }

    let result = false;

    try {
      const response = await postCollectPaper({
        paperId: getPaperId(),
        paperTitle: getDocName(),
      });
      docId.value = response.id;
      setCollected(true);
      result = true;
    } catch (error) {
      result = false;
      const responseError = error as ResponseError;

      if (responseError.code === ERROR_CODE_NEED_VIP) {
        useVipStore().showVipLimitDialog(responseError.message, {
          exception: responseError.extra as NeedVipException,
          leftBtn: {
            text: $i18n.t('common.premium.btns.more'),
            url: `${getDomainOrigin()}/home/mine`,
          },
          reportParams: {
            element_name: ElementName.upperCollectionPopup,
          },
        });
      } else {
        message.error(responseError.message);
      }
    }

    postProcess();
    return result;
  };

  const cancelCollectPaper = async () => {
    if (!preProcess()) {
      return;
    }

    try {
      if (
        !attachmentsStore.attachmentsInfo?.list.length ||
        (await new Promise((resolve) => {
          Modal.confirm({
            title: '确认要取消收藏这篇文献吗？',
            content: '文献删除后，关联的附件也会一并删除',
            onOk: () => resolve(true),
            onCancel: () => resolve(false),
            okButtonProps: {
              danger: true,
            },
            cancelButtonProps: {
              type: 'primary',
            },
          });
        }))
      ) {
        await postCancelCollectPaper({
          paperId: getPaperId(),
        });
        setCollected(false);
      }
    } catch (error) {
      const responseError = error as ResponseError;
      message.error(responseError.message);
    }

    postProcess();
  };

  return {
    collectPaper,
    cancelCollectPaper,
    collectLoading,
    docId,
  };
};
