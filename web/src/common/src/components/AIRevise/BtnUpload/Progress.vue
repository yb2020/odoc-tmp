<script setup lang="ts">
import { ref } from 'vue';
import ErrorResult from '@common/components/ErrorResult.vue';
import { EllipsisOutlined } from '@ant-design/icons-vue';

defineProps<{ error: Error | null; progress: number }>();

const emit = defineEmits<{
  (event: 'cancel'): void;
}>();

const handleCancel = () => {
  emit('cancel');
};

const visible = ref(true);
</script>
<template>
  <a-modal
    v-model:visible="visible"
    class="rp-modal-rounded"
    :title="$t('common.text.progress')"
    :footer="null"
    :closable="!!error"
    :maskClosable="false"
    @cancel="handleCancel"
  >
    <ErrorResult
      v-if="error"
      :error="error"
    />
    <div
      v-else
      class="flex-1 flex flex-col items-center justify-center p-4 h-full"
    >
      <div class="text-base font-medium mb-4">
        {{ $t('common.tips.uploading') }}<EllipsisOutlined />
      </div>
      <a-progress
        :percent="progress"
        status="active"
      />
    </div>
  </a-modal>
</template>
