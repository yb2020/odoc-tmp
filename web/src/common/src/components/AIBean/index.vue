<template>
  <span
    class="aibeans text-sm cursor-pointer inline-flex items-center"
    @click="onClick($event)"
  >
    <IconBean class="mr-1" />
    <u class="underline">
      <LoadingOutlined v-if="loading && !beans" />
      <!-- <RollingText
        v-else
        :duration="1"
        :start-num="beansLast"
        :target-num="beans"
        :direction="beans - beansLast >= 0 ? 'up' : 'down'"
      /> -->
      <Roller
        v-else
        class="!flex-nowrap"
        :value="`${beans}`"
      />
    </u>
  </span>
</template>

<script setup lang="ts">
import { getCurrentInstance } from 'vue';
// import aninejs, { AnimeInstance } from 'animejs'
import { LoadingOutlined } from '@ant-design/icons-vue';
// import { RollingText } from 'vant'
// import 'vant/lib/rolling-text/index.css'
import { Roller } from 'vue-roller';
import 'vue-roller/dist/style.css';
// import { usePrevious } from '@vueuse/core'
import { useAIBeans, useAIBeansBuy } from '../../hooks/useAIBeans';
import IconBean from './Icon.vue';

const emit = defineEmits<{
  (e: 'click', evt: MouseEvent): void;
}>();

const appCtx = getCurrentInstance()?.appContext;

const { beans, loading } = useAIBeans();
// const beansLast = usePrevious(beans.value, 0)
const { showBuyDialog } = useAIBeansBuy(appCtx);

// let animation: AnimeInstance | null = null
// watch(beans, (curr, prev) => {
//   if (curr !== prev) {
//     animation?.pause()
//     animation = aninejs({
//       targets: beansNum,
//       value: curr,
//       round: 1,
//       easing: 'easeOutExpo',
//       duration: 500,
//     })
//     animation.play()
//   }
// })

const onClick = (e: MouseEvent) => {
  emit('click', e);
  if (!e.defaultPrevented) {
    showBuyDialog();
  }
};
</script>

<style lang="less">
.aibeans-modal.ant-modal-wrap {
  .ant-modal,
  .ant-modal-body {
    padding: 0;
  }
  .ant-modal-confirm-content {
    margin-top: 0;
  }
  .ant-modal-confirm-btns {
    display: none;
  }
  .ant-modal-close {
    margin: 27px 24px 0 0;
  }
  .ant-modal-close-x {
    width: 16px;
    height: 16px;
    line-height: 1;
    color: theme('colors.rp-neutral-8');
    &:hover {
      color: theme('colors.rp-neutral-10');
    }
  }
}

.aibeans-confirm.ant-modal-wrap {
  .ant-modal-confirm-body {
    min-height: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .ant-modal-confirm-btns {
    float: none;
    text-align: center;
    margin-top: 0;

    .ant-btn-primary {
      color: #fff;
      width: 168px;
    }
    .ant-btn:not(.ant-btn-primary) {
      display: none;
    }
  }
}
</style>

<style scoped lang="less">
.aibeans {
  // .van-rolling-text-item {
  //   width: auto;
  // }
  :deep(.roller-item) {
    // 需要覆盖style
    width: 8px !important;
  }
}
</style>
