<template>
  <template v-if="tab === TabKeyType.zhpolish">
    <LoadingOutlined v-if="polishEnableLoading" />
    <template v-else-if="polishEnable">
      <span class="text-rp-neutral-8">{{ $t(`common.aibeans.freeuse`) }}{{ $t('common.symbol.colon') }}</span>
      <span class="polish-card">{{ $t(`common.aibeans.cardBetaPolishAct`) }}
      </span>
    </template>
    <slot v-else />
  </template>
  <slot v-else />
</template>
<script setup lang="ts">
import { LoadingOutlined } from '@ant-design/icons-vue';
import { getTrialBeta } from '@common/api/revise';
import { useRequest } from 'ahooks-vue';
import { TabKeyType } from './type';
import { ref } from 'vue';
import { TrialFeature } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/polish/TrialCount';

const props = defineProps<{
  tab: TabKeyType;
}>();

const { data: polishEnable, loading: polishEnableLoading } = useRequest(
  async () => {
    if (props.tab === TabKeyType.zhpolish) {
      const res = await getTrialBeta({
        feature: TrialFeature.ZH_POLISH_REWRITE,
      });

      return !!res?.enabled;
    }
  },
  {}
);
</script>

<style scoped>
.polish-card {
  background: linear-gradient(90deg, #2173e1 0%, #6dded6 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
</style>
