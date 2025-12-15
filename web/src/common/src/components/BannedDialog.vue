<template>
  <Modal
    :visible="true"
    destroy-on-close
    :title="null"
    :footer="null"
    :closable="false"
    :maskClosable="false"
    centered
    :width="IS_MOBILE ? 327 : 482"
    wrap-class-name="banned-modal-wrap"
  >
    <div class="modal-content-wrap">
      <div class="header">
        <i class="aiknowledge-icon icon-logo-readpaper" />{{
          $t('common.account.banned.tt')
        }}
      </div>
      <div class="main-content">
        <div>
          {{ $t('common.account.banned.head') }}{{ $t('common.symbol.colon') }}
        </div>
        <div>
          {{ $t('common.account.banned.msg1')
          }}<span class="text">{{ banInfo.banReason }}</span><span
            v-if="banInfo.banRemark"
            class="text"
          >{{ $t('common.symbol.comma') }}{{ banInfo.banRemark }}</span>{{ $t('common.symbol.comma') }}{{ $t('common.account.banned.msg2')
          }}<span class="text">{{ banInfo.banEndTime }}</span>{{ $t('common.symbol.period') }}
        </div>
        <div>
          {{ $t('common.account.banned.contact1') }}
          <span class="text">readpaper888</span>
          {{ $t('common.account.banned.contact2') }}
        </div>
        <div>{{ $t('common.account.banned.tail') }}</div>
      </div>
      <div class="logout">
        <a-button
          type="primary"
          size="large"
          @click="logout"
        >
          {{
            $t('common.account.banned.btn')
          }}
        </a-button>
      </div>
    </div>
  </Modal>
</template>
<script setup lang="ts">
import { Modal } from 'ant-design-vue';
import { getDomainOrigin, IS_MOBILE } from '@common/utils/env';
import { COOKIE_REFRESH_TOKEN } from '@common/api/const';
import Cookies from 'js-cookie';
import { doLogout } from '~/src/api/user';

const props = defineProps<{
  banInfo: {
    banFlag: boolean;
    banReason: string;
    banRemark?: string;
    banEndTime: string;
  };
}>();

const logout = async () => {
  const result = await doLogout();

  if (result) {
    Cookies.remove(COOKIE_REFRESH_TOKEN);
    window.location.href = getDomainOrigin();
  }
};
</script>

<style lang="less" scoped>
.modal-content-wrap {
  padding: 8px 16px;
  .header {
    font-weight: 600;
    font-size: 24px;
    line-height: 36px;
    color: #1d2229;
    display: flex;
    align-items: center;
    justify-content: center;
    .icon-logo-readpaper {
      font-size: 32px;
      margin-right: 12px;
      color: #1f71e0;
      height: 32px;
      line-height: 32px;
    }
  }
  .main-content {
    font-family: 'Noto Sans SC';
    font-weight: 400;
    font-size: 15px;
    line-height: 24px;
    color: #1d2229;
    margin: 32px 0;
    div + div {
      margin-top: 24px;
      .text {
        color: #1f71e0;
      }
    }
  }
  .logout {
    text-align: center;
    button {
      width: 280px;
    }
  }
}
.mobile-viewport {
  .modal-content-wrap {
    padding: 8px 0px 0px;
    .main-content {
      margin: 24px 0;
      div + div {
        margin-top: 16px;
      }
    }
  }
}
</style>
<style lang="less">
.banned-modal-wrap {
  .ant-modal-body {
    background: #fff;
  }
}
</style>
