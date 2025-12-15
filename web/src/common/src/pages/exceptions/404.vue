<script setup lang="ts">
import { useRouter } from 'vue-router';
const router = useRouter();
const goBack = () => {
  router.replace({ path: '/' });
};
</script>

<template>
  <div class="error-404-page w-full h-full">
    <div
    >vibe reading - odoc.ai</div>
    <div class="desc">
      {{ $t('common.html.pageTitle.pageError') }}
    </div>
    <a-button
      type="primary"
      class="goHomeBtn"
      @click="goBack"
    >
      {{ $t('common.html.pageTitle.gobackHome') }}
    </a-button>
  </div>
</template>

<style scoped lang="less">
.error-404-page {
  background: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  .desc {
    font-size: 22px;
    font-weight: 600;
    color: #262625;
    line-height: 32px;
    margin: 22px 0 24px;
  }
  .goHomeBtn {
    color: rgba(255, 255, 255, 0.85);
  }
}
</style>
