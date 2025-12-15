<template>
  <div
    v-if="!isValidApplyStatus"
    class="h-full flex items-center justify-center w-full"
  >
    <a-spin />
  </div>
  <Welcome
    v-else-if="!hasPermission"
    :from-client="fromClient"
  />
  <Recent
    v-else
    v-bind="props"
  />
</template>

<script setup lang="ts">
import useApplyPermission from '@common/hooks/aitools/useApplyPermission';
import Welcome from './Welcome/index.vue';
import Recent from './Recent.vue';
import { SimpleProjectInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { onActivated } from 'vue';
import {
  reportModuleImpression,
  PageType,
  ModuleType,
} from '@common/utils/report';

const props = defineProps<{
  loading: boolean;
  projects: SimpleProjectInfo[];
  fromClient?: boolean;
  openInNewTab?: boolean;
  precheck: (e: MouseEvent) => void;
  isModuleImpression?: boolean;
}>();

const { hasPermission, isValidApplyStatus } = useApplyPermission();

if (props.isModuleImpression) {
  onActivated(() => {
    // 上报模块曝光
    reportModuleImpression({
      page_type: PageType.POLISH,
      module_type: ModuleType.AI_POLISH_MENTOR,
    });
  });
}
</script>
