


<template>
  <div class="error-container">
    <div
      v-if="imgUrl"
      class="icon"
    >
      <img
        :src="imgUrl"
        alt=""
      >
    </div>
    <div class="detail">
      {{ message }}
    </div>
    <slot />
  </div>
</template>

<script lang="ts" setup>
import Initial from '@/assets/images/initial.png';
import NoPaper from '@/assets/images/no_paper.png';

defineProps<{ message: string, imgUrl?: string }>()


</script>

<style lang="less" scoped>
.error-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  padding: 20px 0;

  .icon {
    width: 64px;
    height: 64px;

    img {
      width: 100%;
      height: 100%;
    }
  }

  .detail {
    font-size: 13px;
    font-weight: 400;
    color: #929599;
    margin-top: 10px;
  }
}

</style>
