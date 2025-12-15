<template>
  <Popover
    :visible="visible"
    title=""
    trigger="click"
    :destroy-tooltip-on-hide="true"
    @visibleChange="$emit('change', $event)"
  >
    <template #content>
      <div class="list-edit-rollback" :style="{ width: width + 'px' }">
        <div class="list-edit-rollback-header">
          {{ $t('home.library.resetData') }}
        </div>
        <div class="list-edit-rollback-body">
          <div>
            <div>{{ $t('home.library.currentData') }}</div>
            <div>
              {{ current }}
            </div>
          </div>
          <div>
            <div>{{ $t('home.library.isReset') }}</div>
            <div>
              {{ origin }}
            </div>
          </div>
        </div>
        <div class="list-edit-rollback-footer">
          <button @click="$emit('cancel')">
            {{ $t('home.global.cancel') }}
          </button>
          <button @click="$emit('ok')">{{ $t('home.global.ok') }}</button>
        </div>
      </div>
    </template>
    <div class="list-item-col-rollback" @click.stop>
      <ReloadOutlined />
    </div>
  </Popover>
</template>

<script lang="ts" setup>
import { Popover } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'

defineProps({
  visible: {
    type: Boolean,
    required: true,
  },
  current: {
    type: [String, Number],
    default: '',
  },
  origin: {
    type: [String, Number],
    default: '',
  },
  width: {
    type: Number,
    default: 240,
  },
})

defineEmits(['cancel', 'ok', 'change'])
</script>

<style lang="less">
.list-edit-rollback {
  .list-edit-rollback-header {
    color: #1d2229;
    font-weight: bold;
    height: 22px;
    margin-bottom: 18px;
  }
  .list-edit-rollback-body {
    margin-bottom: 20px;
  }
  .list-edit-rollback-footer {
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
        color: white;
      }
      &:disabled {
        opacity: 0.7;
      }
    }
  }
}

.list-edit-rollback-body {
  display: flex;
  flex-direction: column;
  padding-bottom: 12px;
  > div {
    display: flex;
    flex-direction: row;
    > div:first-of-type {
      flex: 0 0 100px;
    }
    > div:last-of-type {
      flex: 1 1 100%;
    }
  }
}
</style>
