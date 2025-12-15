<template>
  <Copilot
    :page="PAGE_ROUTE_NAME.CHAT"
    :clip-selecting="clipSelecting"
    :clip-action="clipAction"
  />
</template>

<script lang="ts" setup>
import { defineAsyncComponent } from 'vue';
// 使用defineAsyncComponent正确处理懒加载组件
const Copilot = defineAsyncComponent(() => import('@/components/Right/TabPanel/Copilot/index.vue'));
import { useClip } from '~/src/hooks/useHeaderScreenShot';
import { PAGE_ROUTE_NAME } from '~/src/routes/type';

const { clipAction, clipSelecting } = useClip();
</script>
