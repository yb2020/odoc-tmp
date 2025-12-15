<template>
  <a-modal
    v-model:visible="applyingModalVisible"
    :footer="null"
    :width="560"
    :closable="false"
    class="rp-modal-rounded"
  >
    <div
      v-if="applyStatus === ApplyQuotaStatus.APPLYING"
      class="text-center px-12 pb-5"
    >
      <h1 class="mb-10 font-medium text-2xl">
        {{ $t('common.aitools.welcome.applied') }}
      </h1>
      <p class="font-medium text-lg">
        {{ $t('common.aitools.welcome.waitState') }}
        <br>
        {{ $t('common.aitools.welcome.waitTip') }}
      </p>
      <p class="mb-10 font-medium text-lg">
        {{ $t('common.aitools.welcome.waitDesc') }}
      </p>
      <div>
        <a-button
          type="primary"
          size="large"
          class="w-52 !rounded-[20px]"
          @click="applyingModalVisible = false"
        >
          {{ $t('common.text.ok') }}
        </a-button>
      </div>
    </div>
    <div
      v-else
      class="text-center px-12 pb-5 mt-6"
    >
      <!-- <h1 class="mb-10 font-medium text-2xl">申请内测</h1> -->
      <p class="font-medium text-lg">
        {{ $t('common.aitools.welcome.tip') }}
      </p>
      <p class="font-medium text-lg">
        {{ $t('common.aitools.welcome.desc') }}
      </p>
      <div>
        <a-button
          type="primary"
          size="large"
          class="w-52 !rounded-[20px]"
          @click="handleApply"
        >
          {{ $t('common.aitools.welcome.apply') }}
        </a-button>
      </div>
    </div>
  </a-modal>
</template>
<script setup lang="ts">
import { ApplyQuotaStatus } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import useApplyPermission from '@common/hooks/aitools/useApplyPermission';
// import { watch } from 'vue'

// const props = defineProps<{
//   visible?: boolean
// }>()

const { applyStatus, applyingModalVisible, applyPermission } =
  useApplyPermission();

const handleApply = async () => {
  await applyPermission();
};

// watch(
//   () => props.visible,
//   (v) => {
//     applyingModalVisible.value = v
//   }
// )
</script>
