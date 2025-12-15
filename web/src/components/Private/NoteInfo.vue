<template>
  <div
    v-if="selfNoteInfo && !ignored"
    class="auth"
    :style="{ height: `${FOOTER_NOTE_INFO_HEIGHT}px` }"
  >
    <CloseOutlined
      v-if="pdfStatusInfo.noteUserStatus === UserStatusEnum.GUEST && !IS_MOBILE"
      class="close-button"
      @click="ignored = true"
    />
    <div class="auth-content">
      <div class="user">
        <a-avatar
          v-if="selfNoteInfo.userInfo.avatarUrl"
          :src="selfNoteInfo.userInfo.avatarUrl"
          class="avatar"
        />
        <!-- <a-avatar   暂时移除了点击事件
          v-if="selfNoteInfo.userInfo.avatarUrl"
          :src="selfNoteInfo.userInfo.avatarUrl"
          class="avatar"
          @click="visitUser()"
        /> -->
        <div class="title">
          <div class="name">
            <div class="username">
              {{
                selfNoteInfo.userInfo.showName || selfNoteInfo.userInfo.nickName
              }}
            </div>
            <!-- <div
              v-if="selfNoteInfo.userInfo.tags"
              class="tags"
            >
               <div
                v-for="tag in selfNoteInfo.userInfo.tags.split('|')"
                class="tag"
              >
                {{ tag }}
              </div>
            </div> -->
            <div class="text">
              &nbsp;{{ $t('forbidden.using')
              }}<span
                class="cursor-pointer"
                @click="gotoNewPage"
              >ODOC.AI</span>{{ $t('forbidden.forReading') }}
            </div>
          </div>
        </div>

        <a-dropdown placement="bottomLeft">
          <MoreOutlined class="more-icon" />
          <template #overlay>
            <a-menu>
              <a-menu-item @click="handleUgcReport">
                {{ $t('source.report.button') }}
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>

        <a-button
          type="primary"
          shape="round"
          size="large"
          class="own"
          @click="handleClick"
        >
          {{ $t('forbidden.startReading') }}
        </a-button>
      </div>

      <div class="bottom">
        <div
          v-if="isEmptyPaperId(selfNoteInfo?.paperId ?? '')"
          class="no-paper-title"
        >
          {{ $t('message.noTitle') }} 
        </div>
        <div
          v-else
          class="paper-title"
          :style="{ cursor: selfNoteInfo.isPrivatePaper ? '' : 'pointer' }"
          @click="handleTitle"
        >
          {{ selfNoteInfo?.paperTitle }}
        </div>
        <!-- <div
          class="search"
          @click="handleSearch"
        >
          {{ $t('message.searchOtherPapersTip') }}
        </div> -->
        <div
          v-if="IS_MOBILE"
          class="h5"
        >
          {{ $t('message.searchOtherPapersH5Tip') }}
        </div>
      </div>
    </div>

    <UgcReportDialog />
  </div>
</template>

<script lang="ts" setup>
import { store, selfNoteInfo, pdfStatusInfo } from '@/store';
import { UserStatusEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';
import { getOwnerPaperNoteBaseInfo } from '~/src/api/base';
import { message } from 'ant-design-vue';
import { CloseOutlined, MoreOutlined } from '@ant-design/icons-vue';
import {
  getDomainOrigin,
  getPdfAnnotateNoteUrl,
  IS_MOBILE,
} from '~/src/util/env';
import { bridgeAdaptor } from '~/src/adaptor/bridge';
import { goPathPage, goPersonPage } from '~/src/common/src/utils/url';
import { isEmptyPaperId } from '~/src/api/helper';
import {
  FOOTER_NOTE_INFO_HEIGHT,
  ignored,
} from '~/src/hooks/UserSettings/useFooterNoteInfo';
import {
  ElementClick,
  getPdfIdFromUrl,
  PageType,
  reportElementClick,
} from '~/src/api/report';
import { useUgcReportStore } from '~/src/stores/ugcReport';
import UgcReportDialog from '@/components/ugcReport/dialog.vue';
import {
  ReportType,
  UgcReportSubType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/UgcReport';

const gotoNewPage = () => {
  goPathPage(`${getDomainOrigin()}/new`);
};

const reportButtonClick = (name: ElementClick) => {
  reportElementClick({
    page_type: PageType.note,
    type_parameter: selfNoteInfo.value?.pdfId ?? getPdfIdFromUrl(),
    element_name: name,
  });
};

const visitUser = () => {
  reportButtonClick(ElementClick.head_portrait);

  goPersonPage(selfNoteInfo.value?.userInfo.id ?? '');
};

const handleClick = async () => {
  reportButtonClick(ElementClick.my_note);

  if (!store.getters['user/isLogin']) {
    bridgeAdaptor.login();
    return;
  }

  if (!pdfStatusInfo.value.hasPdfAccessFlag) {
    message.error('暂无该PDF权限');
    return;
  }

  if (!selfNoteInfo.value) {
    return;
  }

  const pdfId = pdfStatusInfo.value.authPdfId || selfNoteInfo.value?.pdfId;
  if (!pdfId) {
    return;
  }

  try {
    const { noteId } = await getOwnerPaperNoteBaseInfo({
      pdfId,
    });

    goPathPage(String(getPdfAnnotateNoteUrl({ noteId })));
  } catch (err) {
    message.error((err as Error).message);
  }
};

const handleTitle = () => {
  reportButtonClick(ElementClick.paper_link);

  if (selfNoteInfo.value && !selfNoteInfo.value.isPrivatePaper) {
    goPathPage(`${getDomainOrigin()}/paper/${selfNoteInfo.value.paperId}`);
  }
};

const handleSearch = () => {
  reportButtonClick(ElementClick.search_paper);

  if (!selfNoteInfo.value) {
    return;
  }

  const docName =
    selfNoteInfo.value?.paperTitle || selfNoteInfo.value?.docName || '';
  goPathPage(`${getDomainOrigin()}/search/${encodeURIComponent(docName)}`);
};

const ugcReportStore = useUgcReportStore();

const handleUgcReport = () => {
  if (!store.getters['user/isLogin']) {
    bridgeAdaptor.login();
    return;
  }

  const reportParams = {
    reportType: ReportType.USER_GENERATED_CONTENT,
    ugcReportSubType: UgcReportSubType.PERSONAL_OPEN_NOTE,
    contentId: selfNoteInfo.value?.noteId,
    pageUrl: window.location.href,
  };

  if (!IS_MOBILE) {
    ugcReportStore.showUgcReportDialog(reportParams);

    return;
  }

  const url = new URL(`${window.location.origin}/mobileUgcReport`);

  for (const key of Object.keys(reportParams)) {
    url.searchParams.set(key, (reportParams as any)[key]);
  }

  window.location.href = url.toString();
};
</script>

<style lang="less" scoped>
.auth {
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100%;
  background: rgba(245, 245, 245, 0.8);

  box-shadow: 0px -4px 13px rgba(0, 0, 0, 0.12);

  backdrop-filter: blur(8px);

  display: flex;
  flex-direction: column;
  padding-top: 16px;
  align-items: center;
  z-index: 200;

  .close-button {
    position: absolute;
    right: 26px;
    z-index: 99999;
    color: #4e5969;
    top: 18px;
    font-size: 17px;
    cursor: pointer;
  }

  .auth-content {
    width: 758px;
  }

  .user {
    display: flex;
    position: relative;
    height: 48px;

    .avatar {
      border: 1px solid #ffffff;
      margin-right: 20px;
      cursor: pointer;
      height: 40px;
      width: 40px;
      box-sizing: border-box;
    }

    .title {
      display: flex;
      height: 40px;
      align-items: center;
      max-width: calc(100% - 210px);

      .name {
        font-size: 16px;
        font-weight: 600;
        color: rgba(0, 0, 0, 84%);
        display: flex;
        overflow: hidden;
        .username {
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }
        .text {
          flex-shrink: 0;
          span {
            color: #1f71e0;
          }
        }
      }

      .tags {
        display: flex;
        margin-left: 8px;
        .tag {
          height: 24px;
          background: #e1eaf5;
          border-radius: 11px;
          padding: 3px 8px;

          font-size: 12px;
          font-weight: 400;
          color: #1f71e0;
          margin-right: 5px;
          display: flex;
          align-items: center;
          justify-content: center;
          &:last-child {
            margin-right: 8px;
          }
        }
      }
    }

    .own {
      position: absolute;
      right: 0;
      height: 48px;
      border-radius: 2px;
    }
  }

  .bottom {
    display: flex;
    justify-content: space-between;
    margin-top: 16px;
    line-height: 22px;
    font-size: 14px;

    .no-paper-title,
    .paper-title {
      width: 604px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    .no-paper-title {
      color: #4e5969 !important;
    }
    .paper-title {
      font-family: 'Lato';
      font-weight: bold;
      color: #1f71e0;
    }

    .search {
      margin-left: 48px;
      color: rgba(0, 0, 0, 86%);
      text-decoration: underline;
      width: 106px;
      display: flex;
      justify-content: center;
      cursor: pointer;
    }
  }

  .more-icon {
    transform: rotate(90deg);
    font-size: 16px;
    margin-left: 16px;
    margin-top: 12px;
    color: #000;
    height: 16px;
  }
}

.mobile-viewport {
  .auth {
    height: auto;
    background: rgba(245, 245, 245, 0.85);
    backdrop-filter: blur(2rpx);
    padding: 20rpx 16rpx 25rpx 16rpx;
    justify-content: flex-start;
    align-items: flex-start;
    // z-index: 0;
    .auth-content {
      width: 100%;
    }

    .user {
      width: 100%;
      .avatar {
        width: 40rpx;
        height: 40rpx;
        border: 2rpx solid #ffffff;
        margin-right: 12rpx;
        flex-shrink: 0;
      }

      .own {
        display: none;
      }

      .title {
        max-width: calc(100% - 74rpx);
        .name {
          font-size: 16rpx;
          line-height: 22rpx;
          overflow: hidden;
        }

        .tags {
          .tag {
            height: 22rpx;
            background: #fafafa;
            border-radius: 11rpx;
            padding: 3rpx 8rpx;

            font-size: 12rpx;
            margin-right: 5rpx;
          }
        }
      }
    }

    .bottom {
      width: 100%;
      margin-top: 16rpx;
      flex-direction: column;
      .paper-title {
        width: auto;
        font-size: 14rpx;
        font-family: Roboto-Bold, Roboto;
        font-weight: bold;
        color: #1f71e0;
        line-height: 22rpx;
        margin-right: 56rpx;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        text-overflow: ellipsis;
        word-break: break-all;
      }

      .search {
        display: none;
      }

      .h5 {
        font-size: 12rpx;
        font-weight: 400;
        color: #757980;
        line-height: 18rpx;
        margin-top: 16rpx;
      }
    }
  }
}
</style>
