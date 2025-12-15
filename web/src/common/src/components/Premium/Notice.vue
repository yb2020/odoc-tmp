<template>
  <div
    v-show="barVisible"
    :class="{ 'notice-mobile': isMobile }"
  >
    <!-- <div class="notice">
      <div style="display: flex; align-items: center" @click="clickNotice()">
        <a-icon type="sound" theme="filled" class="iconsounder_fill" />
        <div class="notice-msg">
          尊敬的用户您好，由于产品功能的迭代升级，我们的VIP规则有相应调整，请点击查看详情。
        </div>
      </div>
      <div @click="close">
        <i class="aiknowledge-icon icon-nav-close" aria-hidden="true" />
      </div>
    </div> -->
    <!-- <Modal
      :visible="false"
      width="480px"
      :closable="false"
      dialog-class="folder-upgrade"
      centered
      :mask-closable="false"
    >
      <div class="title">
        <i
          class="aiknowledge-icon icon-upgrade-circle-fill"
          aria-hidden="true"
        />
        <span>小组讨论功能上线！</span>
      </div>
      <div class="dialog dialog-part1">
        <div class="sub-title">小组讨论功能上线</div>
        <div>
          ①小组成员可以对小组内的论文进行标注及回复（在PDF阅读界面右侧-小组讨论tab-选择小组后可以开启讨论）
        </div>
        <div>②开启讨论的方式是直接往对应小组里传论文</div>
        <div>
          ③若是有版权争议的论文，则依然是只有已经自己上传过的小组成员才能讨论
        </div>
      </div>
      <div class="dialog dialog-part2">
        <div class="sub-title">另外为了保证小组讨论功能的稳定，我们决定：</div>
        <div>①限制单个用户创建及加入小组的上限，现阶段最多创建3个、加入5个</div>
        <div>
          ②已经超过上限的用户不会被强制解散/退出小组，但无法再新创建或者新加入小组
        </div>
        <div>③后续随功能的逐步稳定，我们会考虑重新放开限制</div>
      </div>

      <div class="dialog dialog-part3">
        相关功能演示可以<a
          href="https://docs.qq.com/doc/DVWh0ZlpVWlNMUmFy?u=eda27474cce44903bae790026a045e3c"
          target="_blank"
          style="color: #1f71e0; cursor: pointer"
          >看这里</a
        >
      </div>

      <span slot="footer" class="dialog-footer">
        <a-button class="btn-close" @click="close()">知道了</a-button>
      </span>
    </Modal> -->
  </div>
</template>
<script lang="ts">
import { computed, defineComponent, onMounted, ref } from 'vue';
import { Modal } from 'ant-design-vue';
import { useLocalStorage } from '@vueuse/core';
import { postVipAnnouncementVisible } from '~/src/api/user';

export const CURRENT_NOTICE_KEY = 'notice';
export const CURRENT_NOTICE_ID = '231208';

export default defineComponent({
  setup(_, { emit }) {
    const isMobile = false;
    const modalVisible = ref(false);
    const close = () => {
      modalVisible.value = false;
      read.value = CURRENT_NOTICE_ID;
      emit('closeNotice');
    };
    const download = () => {
      window.location.href = '/download';
    };
    const clickNotice = () => {
      window.open('/user/message?type=vipAnnouncement', '_blank');
    };

    const vipVisible = ref(false);

    const read = useLocalStorage(CURRENT_NOTICE_KEY, '');

    const isRead = computed(() => {
      return read.value === CURRENT_NOTICE_ID;
    });

    const barVisible = computed(() => {
      return !isRead.value && vipVisible.value;
    });

    onMounted(async () => {
      if (isRead.value) {
        return;
      }

      try {
        const res = await postVipAnnouncementVisible();
        vipVisible.value = res;
      } catch (error) {
        console.log(error);
      }
    });
    return {
      isMobile,
      barVisible,
      modalVisible,
      close,
      download,
      clickNotice,
    };
  },
});
</script>
<style lang="less" scoped>
.notice-mobile {
  .notice {
    width: 100% !important;
    height: auto !important;
    .notice-msg {
      width: 100% !important;
      height: auto !important;
      line-height: initial !important;
    }
    .iconsounder_fill {
      margin-left: 10px !important;
    }
  }
}
.notice {
  width: 100%;
  height: 36px;
  flex: 0 0 36px;
  background: #faf0e0;
  border-radius: 2px;
  display: flex;
  justify-content: space-between;
  cursor: pointer;
  .iconsounder_fill {
    margin-left: 20px;
    color: #e68b39;
  }
  .icon-nav-close {
    margin-left: auto;
    margin-right: 10px;
    margin-top: 5px;
    color: #e68b39;
    display: inline-block;
    cursor: pointer;
    font-size: 16px;
  }

  .notice-msg {
    width: 1040px;
    height: 36px;
    line-height: 36px;
    margin-left: 5px;
    font-size: 14px;
    font-weight: bold;
    color: #e68b39;
  }
}
.folder-upgrade {
  .ant-modal-close {
    height: 48px;
    line-height: 48px;
  }
  .dialog {
    width: 420px;
    font-size: 14px;
    color: #262625;
    margin: 0 auto;
  }
  .dialog-part1 {
    line-height: 24px;
  }
  .dialog-part2 {
    line-height: 24px;
    margin-top: 12px;
  }
  .dialog-part3 {
    padding-top: 20px;
    padding-bottom: 20px;
    color: #73716f;
    img {
      width: 20px;
      height: 20px;
      margin-top: -2px;
      margin-right: 2px;
    }
  }

  .title {
    height: 60px;
    line-height: 60px;
    font-size: 18px;
    font-weight: 600;
    color: #262625;
    margin-left: 10px;
    display: flex;
    align-items: center;
    .icon-upgrade-circle-fill {
      font-size: 18px;
      color: #1f71e0;
      height: 18px;
      line-height: 18px;
    }
    span {
      margin-left: 2px;
    }
  }
  .sub-title {
    font-weight: 600;
  }

  .btn-download {
    width: 118px;
    height: 32px;
    border-radius: 2px;
    text-align: center;
    color: #1f71e0;
    border: 1px solid #1f71e0;
  }
  .btn-close {
    width: 90px;
    height: 32px;
    border-radius: 2px;
    text-align: center;
    color: #fff;
    background: #1f71e0;
  }
}
:deep(.ant-modal-footer) {
  border-top: none;
  padding-bottom: 20px;
}
</style>
