<template>
  <div class="initial-container">
    <div class="icon">
      <img
        :src="img"
        alt=""
      >
    </div>
    <div class="detail">
      {{ message || $t('message.noDataTip') }}
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent } from 'vue';
import InitialImg from '@/assets/images/initial.png';

export default defineComponent({
  props: {
    message: {
      type: String,
      default: '',
    },
    img: {
      type: String,
      default: InitialImg,
    },
  },
});
</script>
<style lang="less" scoped>
.initial-container {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  height: 100%;
  position: relative;
  user-select: none;

  .icon {
    width: 64px;
    height: 64px;
  }

  .detail {
    font-size: 13px;
    font-weight: 400;
    color: var(--site-theme-text-secondary);
    margin-top: 10px;
  }
}
</style>
