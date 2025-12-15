<template>
  <a-modal
    v-model:visible="visible"
    :footer="null"
    :closable="false"
  >
    <template #title>
      <div class="title">
        <exclamation-circle-outlined
          :style="{ color: '#FF8C19', fontSize: '21px' }"
          class="icon"
        />
        <div>{{ title || 'unknown error' }}</div>
      </div>
    </template>
    <div class="error-container">
      <a :href="userGuideLink" target="_blank" class="user-guide-link">{{ $t('viewer.pdfError.guide') }}</a>
      <div class="error">
        <div>{{ $t('viewer.pdfError.errorLog') }}</div>
        {{ log }}
      </div>

      <a-button
        type="primary"
        class="btn"
        @click="handleOk"
      >
        {{ $t('viewer.pdfError.confirm') }}
      </a-button>
    </div>
  </a-modal>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

defineProps<{ title: string; log: string }>();

const visible = defineModel('visible', { default: false });
const { isCurrentLanguage } = useLanguage()

// 创建用户指南链接的计算属性
const userGuideLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide'
  }
  return '/docs/guide'
})

const handleOk = () => {
  visible.value = false;
};


</script>

<style lang="less" scoped>
.title {
  font-size: 14px;
  font-weight: 500;
  // color: #262625;

  display: flex;
  align-items: center;

  .icon {
    margin-right: 13px;
  }
}
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 10px;

  .btn {
    margin-top: 16px;
  }
}

.qrcode {
  width: 128px;
  height: 128px;
}

.error {
  font-size: 14px;
  font-weight: 400;
  // color: #262625;
  margin-top: 16px;
  word-break: break-all;
}
</style>
