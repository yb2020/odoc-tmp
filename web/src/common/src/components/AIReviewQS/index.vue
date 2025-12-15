<template>
  <div class="flex flex-col h-full overflow-auto">
    <Tip
      class="p-4"
      text=""
    >
      <template
        v-if="vipStore.role.vipType !== VipType.FREE"
        #beans
      >
        <slot name="beans" />
      </template>
    </Tip>
    <TaskList
      ref="$list"
      class="px-4"
    >
      <template #upload>
        <a-button
          class="!w-32 !h-10 !text-base font-medium"
          type="primary"
          shape="round"
          @click="visible = true"
        >
          {{ $t('common.text.upload') }}{{ $t('common.symbol.space')
          }}{{ $t('common.text.paper') }}
        </a-button>
      </template>
      <section class="px-6 text-center">
        <h2 class="text-5xl mt-8 font-normal">
          {{ $t('common.aiReviewerQS.tt') }}
        </h2>
        <p
          class="text-3xl mt-12"
          v-html="$t('common.aiReviewerQS.desc')"
        />
        <a-button
          class="my-8 !w-48 !h-14 !text-xl font-medium"
          type="primary"
          shape="round"
          @click="visible = true"
        >
          {{ $t('common.text.upload') }}{{ $t('common.symbol.space')
          }}{{ $t('common.text.paper') }}
        </a-button>
        <div>
          <p>
            {{ $t('common.aiReviewerQS.demo') }}（<a
              v-if="data?.exampleFileUrl"
              :href="data.exampleFileUrl"
              class="underline"
            >
              {{ $t('common.aiReviewerQS.demotip') }}</a>）
          </p>
          <img
            class="bg-white w-full min-h-[600px]"
            :src="data?.examplePicUrl"
            alt="Demo"
          >
        </div>
      </section>
    </TaskList>
  </div>
  <TaskForm
    v-model:visible="visible"
    :data="data!"
    @created="onCreated"
  />
</template>

<script setup lang="ts">
import { useRequest } from 'ahooks-vue';
import { onActivated, ref } from 'vue';
import Tip from '@common/components/AITools/Tip.vue';
import TaskList from './TaskList.vue';
import TaskForm from './TaskForm.vue';
import { getQSReviewConfig } from '@common/api/review';
import { useVipStore, VipType } from '@common/stores/vip';
import {
  PageType,
  ModuleType,
  reportModuleImpression,
} from '../../utils/report';

const $list = ref();
const visible = ref(false);

const vipStore = useVipStore();

const { data } = useRequest(async () => {
  const res = await getQSReviewConfig();

  return {
    // @ts-ignore
    currentAiBeanPrice: 500,
    // @ts-ignore
    originalAiBeanPrice: 2880,
    // @ts-ignore
    aiBeanAmountOneYuan: 10,
    ...res,
  };
}, {});

const onCreated = () => {
  $list.value?.refresh();
};

onActivated(() => {
  // 上报模块曝光
  reportModuleImpression({
    page_type: PageType.POLISH,
    module_type: ModuleType.AI_MENTOR_GRA,
  });
});
</script>

<style scoped></style>
