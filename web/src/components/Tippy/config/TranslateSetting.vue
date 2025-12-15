<template>
  <a-modal
    v-model:visible="visible"
    centered
    :maskClosable="false"
    class="translate-modal"
    title="管理我的翻译渠道"
    width="640px"
    :zIndex="1050"
    footer=""
    :body-style="{ padding: '5px' }"
  >
    <div class="settingLeft">
      <TranslateConfig />
    </div>
    <div class="divider" />
    <div class="settingRight">
      <TranslateTips />
    </div>
  </a-modal>
</template>

<script lang="ts" setup>
import { defineComponent, reactive, provide, watch } from 'vue';
import { Modal } from 'ant-design-vue';
import TranslateConfig from './TranslateConfig.vue';
import TranslateTips from './TranslateTips.vue';
import { emitter, CONFIG_CLOSE_ACTION } from './config';

const visible = defineModel('visible', { default: false });
</script>

<style scoped lang="less">
.matedata-setting-container {
  background: white;
  width: 640px;
  z-index: 1050;
}
</style>

<style lang="less">
.translate-modal {
  background: white;
  display: flex;
  margin: 0px 16px;
  font-size: 15px;

  .ant-modal-content {
    background: white;

    .ant-modal-header {
      background: white;
      color: black;
      border-bottom: 0px solid #f0f0f0;

      .ant-modal-title {
        color: black;
        font-size: 16px;
      }
    }

    .ant-modal-close-icon {
      color: black;
    }

    .ant-modal-body {
      color: black;
      background: white;
      display: flex;
    }
  }

  .settingLeft {
    background: #f7f8fa;
    width: 35%;
    color: black;
    height: 400px;
    margin: 0px 0px 16px 16px;
    padding: 0px 5px;
    border: 1px solid #e5e6eb;
  }

  .settingRight {
    background: white;
    width: 65%;
    color: black;
    height: 400px;
    margin: 0px 16px 16px 0px;
    border-top: 1px solid #e5e6eb;
    border-right: 1px solid #e5e6eb;
    border-bottom: 1px solid #e5e6eb;
  }
}
</style>
