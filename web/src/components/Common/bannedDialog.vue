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
        <i class="aiknowledge-icon icon-logo-readpaper" />账号封禁通知
      </div>
      <div class="main-content">
        <div>尊敬的用户:</div>
        <div>
          我们注意到您的账号存在<span class="text">{{ banInfo.banReason }}</span><span
            v-if="banInfo.banRemark"
            class="text"
          >：{{ banInfo.banRemark }}</span>，为了维护社区秩序，我们不得不对您的账号进行封禁处理。很抱歉告知您，您的账号将被封禁至<span
            class="text"
          >{{ banInfo.banEndTime }}</span>。
        </div>
        <div>
          如果您对此有任何异议，欢迎联系我们，我们会尽快为您解决问题。请加入我们官方小管家的微信：<span
            class="text"
          >readpaper888</span>，我们将在第一时间为您提供帮助。
        </div>
        <div>再次感谢您对我们的支持和理解，我们期待您的再次使用。</div>
      </div>
      <div class="logout">
        <a-button
          type="primary"
          size="large"
          @click="logout"
        >
          确认退出登录
        </a-button>
      </div>
    </div>
  </Modal>
</template>
<script setup lang="ts">
import { Modal } from 'ant-design-vue';
import { fetchPersonalLogout } from '@/api/user';
import { getDomainOrigin, IS_MOBILE } from '~/src/util/env';
import Cookies from 'js-cookie';
import { COOKIE_REFRESH_TOKEN } from '~/src/api/const';

const props = defineProps<{
  banInfo: {
    banFlag: boolean;
    banReason: string;
    banRemark?: string;
    banEndTime: string;
  };
}>();

const logout = async () => {
  const result = await fetchPersonalLogout();

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
