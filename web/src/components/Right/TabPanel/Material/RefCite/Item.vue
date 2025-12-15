<template>
  <div class="paper-item">
    <div
      ref="triggerRef"
      class="item"
      @click.stop="openPaperPage"
    >
      <div class="index">
        [{{ index }}]
      </div>
      <div class="title">
        {{ title }}
      </div>
    </div>
    <TippyVue
      v-if="triggerRef"
      ref="tippyRef"
      :trigger-ele="triggerRef"
      :placement="'left-start'"
      :trigger="'mouseenter'"
      @onShown="onShown"
    >
      <ReferenceVue
        :paper-id="paperId"
        :paper-title="title"
        :fetch-flag="startFetch"
        @update-content="handleTippyUpdate"
        @hide="handleTippyHide"
      />
    </TippyVue>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue';
import ReferenceVue from '../../../../Tippy/Reference/index.vue';
import TippyVue from '../../../../Tippy/index.vue';
import { getDomainOrigin } from '~/src/util/env';
import { goPathPage } from '~/src/common/src/utils/url';

const { paperId } =
  defineProps<{ index: number; title: string; paperId: string }>();

const triggerRef = ref();

const tippyRef = ref();

const startFetch = ref(false);

const openPaperPage = (e: any) => {
  return;
  // 旧方式调整链接
  // if (!paperId) {
  //   return;
  // }

  // console.log(e);
  // goPathPage(`${getDomainOrigin()}/paper/${paperId}`);
};

const onShown = () => {
  startFetch.value = true;
};

const handleTippyUpdate = () => {
  tippyRef.value.update();
};

const handleTippyHide = () => {
  tippyRef.value.hide();
};
</script>

<style scoped lang="less">
.item {
  display: flex;
  .index {
    margin-right: 6px;
    color: #86919c;
    cursor: pointer;
  }

  .title {
    display: -webkit-box;
    overflow: hidden;
    -webkit-line-clamp: 5;
    -webkit-box-orient: vertical;
    cursor: pointer;
    color: #1d2229;
  }
}

.paper-item {
  padding: 8px 20px;

  &:hover {
    background-color: #f0f2f5;
  }
}

// .paper-item + .paper-item {
//   margin-top: 15px;
// }
</style>
