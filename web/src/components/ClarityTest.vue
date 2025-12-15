<template>
  <div class="clarity-test-container">
    <div class="test-section">
      <h2>Microsoft Clarity 集成测试</h2>
      <p>此组件用于测试 Microsoft Clarity 的各项功能</p>
      
      <div class="button-group">
        <a-button type="primary" @click="testIdentifyUser">
          识别用户
        </a-button>
        
        <a-button @click="testSetTag">
          设置标签
        </a-button>
        
        <a-button @click="testTrackEvent">
          记录事件
        </a-button>
        
        <a-button @click="checkStatus">
          检查状态
        </a-button>
      </div>
      
      <div v-if="status" class="status-info">
        <h3>Clarity 状态</h3>
        <p><strong>初始化状态:</strong> {{ status.initialized ? '已初始化' : '未初始化' }}</p>
        <p><strong>项目 ID:</strong> {{ status.projectId }}</p>
        <p><strong>最后操作:</strong> {{ lastAction }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { clarityService, identifyUser, setClarityTag, trackEvent } from '@/utils/clarity';

const status = ref<{
  initialized: boolean;
  projectId: string;
} | null>(null);

const lastAction = ref<string>('无');

const testIdentifyUser = () => {
  const userId = `test_user_${Date.now()}`;
  identifyUser(userId);
  lastAction.value = `识别用户: ${userId}`;
  message.success('用户识别完成');
};

const testSetTag = () => {
  const tagKey = 'test_tag';
  const tagValue = `test_value_${Date.now()}`;
  setClarityTag(tagKey, [tagValue]);
  lastAction.value = `设置标签: ${tagKey} = ${tagValue}`;
  message.success('标签设置完成');
};

const testTrackEvent = () => {
  const eventName = 'test_button_click';
  const properties = {
    timestamp: new Date().toISOString(),
    page: 'clarity_test',
    action: 'button_click'
  };
  trackEvent(eventName, properties);
  lastAction.value = `记录事件: ${eventName}`;
  message.success('事件记录完成');
};

const checkStatus = () => {
  status.value = {
    initialized: clarityService.initialized,
    projectId: clarityService.id
  };
  lastAction.value = '检查状态';
  message.info('状态已更新');
};

// 组件挂载时自动检查状态
checkStatus();
</script>

<style scoped>
.clarity-test-container {
  padding: 24px;
  max-width: 600px;
  margin: 0 auto;
}

.test-section {
  background: #f8f9fa;
  padding: 24px;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.button-group {
  display: flex;
  gap: 12px;
  margin: 24px 0;
  flex-wrap: wrap;
}

.status-info {
  margin-top: 24px;
  padding: 16px;
  background: #fff;
  border-radius: 6px;
  border: 1px solid #d1d5db;
}

.status-info h3 {
  margin-top: 0;
  color: #374151;
}

.status-info p {
  margin: 8px 0;
  color: #6b7280;
}
</style>
