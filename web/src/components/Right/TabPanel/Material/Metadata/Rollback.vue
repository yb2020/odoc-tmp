<template>
  <Popover
    :visible="visible"
    title=""
    trigger="click"
    :destroy-tooltip-on-hide="true"
    overlay-class-name="metadata-popover"
    placement="left"
  >
    <template #content>
      <div
        class="metadata-rollback-confirm"
        :style="{ width }"
      >
        <div class="metadata-rollback-confirm-header">
          {{ $t('message.resetData.title') }}
        </div>
        <div class="metadata-rollback-confirm-body">
          <div>
            <div>{{ $t('message.resetData.content1') }}</div>
            <div class="thin-scroll">
              {{ currentInfo }}
            </div>
          </div>
          <div>
            <div>{{ $t('message.resetData.content2') }}</div>
            <div class="thin-scroll">
              {{ originInfo }}
            </div>
          </div>
        </div>
        <div class="metadata-rollback-confirm-footer">
          <button @click.stop="visible = false">
            {{ $t('message.resetData.cancelText') }}
          </button>
          <button @click.stop="emit('click')">
            {{ $t('message.resetData.confirmText') }}
          </button>
        </div>
      </div>
    </template>
    <Tooltip
      :title="$t('message.rollbackDataTip')"
      placement="topRight"
      overlay-class-name="metadata-rollback-tooltip"
    >
      <div
        class="metadata-rollback"
        @click.stop="visible = true"
      >
        <ReloadOutlined />
      </div>
    </Tooltip>
  </Popover>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { Tooltip, Popover } from 'ant-design-vue';
import { ReloadOutlined } from '@ant-design/icons-vue';

defineProps<{
  width: string;
  currentInfo: string;
  originInfo: string;
}>();

const emit = defineEmits<{
  (event: 'click'): void;
}>();

const visible = ref(false);
</script>

<style lang="less">
@import './style.less';

.metadata-rollback-tooltip {
  max-width: 999px !important;
  .ant-tooltip-inner {
    background-color: #ffffff;
    font-size: 12px;
    color: #4e5969;
  }
  .ant-tooltip-arrow-content {
    background-color: #ffffff;
  }
}
.metadata-rollback {
  position: absolute;
  top: 4px;
  right: 2px;
  height: 24px;
  width: 24px;
  display: flex;
  visibility: hidden;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
  background-color: #5f666d;
  border-radius: 2px;
  cursor: pointer;
}

.metadata-rollback-confirm {
  .metadata-rollback-confirm-header {
    color: #1d2229;
    font-weight: bold;
    height: 22px;
    margin-bottom: 18px;
  }
  .metadata-rollback-confirm-footer {
    display: flex;
    justify-content: flex-end;
    button {
      margin-left: 12px;
      width: 64px;
      height: 24px;
      border: 0;
      outline: 0;
      border-radius: 2px;
      font-size: 12px;
      cursor: pointer;
      &:first-of-type {
        background: #f0f2f5;
        color: #4e5969;
      }
      &:last-of-type {
        background: #1f71e0;
        color: #fff;
      }
      &:disabled {
        opacity: 0.7;
      }
    }
  }
}

.metadata-rollback-confirm-body {
  display: flex;
  flex-direction: column;
  padding-bottom: 12px;
  > div {
    display: flex;
    flex-direction: row;
    margin-bottom: 4px;
    > div:first-of-type {
      flex: 0 0 100px;
    }
    > div:last-of-type {
      flex: 1 1 100%;
      max-height: 300px;
      overflow: auto;
    }
  }
}
</style>
