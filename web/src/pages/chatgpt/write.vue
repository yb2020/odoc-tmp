<template>
  <a-row
    type="flex"
    :gutter="24"
    class="write-container"
  >
    <a-col :flex="1">
      <div class="write-container-text">
        <a-textarea
          v-model:value="textValue"
          placeholder="在这里填写你的论文片段/全文"
          :rows="40"
        />
      </div>
    </a-col>
    <a-col flex="40%">
      <div class="write-container-chat">
        <Copilot
          :page="PAGE_ROUTE_NAME.WRITE"
          :text-value="textValue"
          :clip-selecting="clipSelecting"
          :clip-action="clipAction"
        />
      </div>
    </a-col>
  </a-row>
</template>

<script lang="ts" setup>
import { defineAsyncComponent } from 'vue';
// 使用defineAsyncComponent正确处理懒加载组件
const Copilot = defineAsyncComponent(() => import('@/components/Right/TabPanel/Copilot/index.vue'));
import { ref } from 'vue';
import { useClip } from '~/src/hooks/useHeaderScreenShot';
import { PAGE_ROUTE_NAME } from '~/src/routes/type';

const textValue = ref('');

const { clipAction, clipSelecting } = useClip();
</script>
<style lang="less">
.write-container {
  height: 100%;
  flex-flow: row nowrap !important;
  &-text {
    // background-color: #383a3d;
    padding: 24px;
    overflow: auto;
    height: calc(100vh - 56px);
  }
  &-chat {
    height: calc(100vh - 56px);
    border-left: 2px solid #383a3d;
  }
}
</style>
