<template>
  <LoadingOutlined v-if="loading" />
  <template v-else-if="isUnlimitedCard">
    <span class="text-rp-neutral-8">{{ $t(`common.aibeans.freeuse`) }}{{ $t('common.symbol.colon') }}</span>
    <span class="polish-card">{{ $t(`common.aibeans.cardUnlimitedPolishAct`)
    }}{{ $t('common.symbol.space')
    }}{{ $t(`common.aibeans.freeuseTip`) }}</span>
  </template>
  <slot v-else />
</template>

<script setup lang="ts">
import { LoadingOutlined } from '@ant-design/icons-vue';
import { ActivityTypeEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipActivitiesInfo';
import { computed } from 'vue';
import { useActivityCards } from '../../hooks/useActivityCards';

const { cards, loading } = useActivityCards();
const isUnlimitedCard = computed(
  () => cards.value[ActivityTypeEnum.UN_LIMIT_POLISH_CARD]
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
