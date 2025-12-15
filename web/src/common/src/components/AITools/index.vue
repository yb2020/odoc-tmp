<script setup lang="ts">
import { computed, provide, ref, shallowRef, watch } from 'vue';
import TranslateIcon from '~common/assets/images/aitools/translate-icon.svg';
import TranslateActiveIcon from '~common/assets/images/aitools/translate-icon__active.svg';
import PolishIcon from '~common/assets/images/aitools/polish-icon.svg';
import PolishActiveIcon from '~common/assets/images/aitools/polish-icon__active.svg';
import ZhPolishIcon from '~common/assets/images/aitools/zhpolish-icon.svg';
import ZhPolishActiveIcon from '~common/assets/images/aitools/zhpolish-icon__active.svg';
import Drawer from '@common/components/Drawer/index.vue';
import LoginFirst from '@common/components/LoginFirst/index.vue';
import AiBeanAuth from '@common/components/AIBean/Auth.vue';
import ApplyModal from '@common/components/AITools/Apply/index.vue';
import RecentPage from '@common/components/AIRevise/index.vue';
import ReviewQSPage from '@common/components/AIReviewQS/index.vue';
import PolishPage from '../AIParagraph/index.vue';
import TranslatePage from '../AITranslate/index.vue';
import UnlimitedCard from './UnlimitedCard.vue';
import BetaCard from './BetaCard.vue';

import '@idea/aiknowledge-icon/dist/css/iconfont.css';
import { useLocalStorage } from '@vueuse/core';
import { useRouter } from 'vue-router';
import { useRequest } from 'ahooks-vue';
import { VipType, useVipStore } from '@common/stores/vip';
import useApplyPermission from '@common/hooks/aitools/useApplyPermission';
import { getTrialCount } from '@common/api/revise';
import { TrialFeature } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/polish/TrialCount';
import { debounce } from 'lodash';
import { TabKeyType } from './type';
import { compareVipType } from '../Premium/types';

const props = defineProps<{
  fromClient?: boolean;
  openInNewTab?: boolean;
  isPermissionInited?: boolean;
  isModuleImpression?: boolean;
  isMalong?: boolean;
}>();

const router = useRouter();
const vipStore = useVipStore();
const { hasPermission, initPermission, checkPermission } = useApplyPermission();
const {
  data: trialCounts,
  refresh: refreshTrial,
  loading: trialLoading,
} = useRequest(getTrialCount, {});
const refreshTrialCount = debounce(() => {
  if (activeTabTrialCount.value > 0) {
    refreshTrial();
  } else {
    activeTabRefundable.value = true;
  }
}, 800);

if (!props.isPermissionInited) {
  initPermission();
}

const TAB_LIST = props.isMalong
  ? [
      {
        key: TabKeyType.zhpolish,
        label: '中文润色改写',
        i18nKey: 'common.aitools.zhAiPolish',
        component: PolishPage,
        badge: 'new',
        icon: ZhPolishIcon,
        activeIcon: ZhPolishActiveIcon,
      },
      {
        key: TabKeyType.polish,
        label: '英文润色改写',
        i18nKey: 'common.aitools.enAiPolish',
        component: PolishPage,
        icon: PolishIcon,
        activeIcon: PolishActiveIcon,
      },
    ]
  : [
      {
        key: TabKeyType.zhpolish,
        label: '中文润色改写',
        i18nKey: 'common.aitools.zhAiPolish',
        component: PolishPage,
        badge: 'new',
        icon: ZhPolishIcon,
        activeIcon: ZhPolishActiveIcon,
      },
      {
        key: TabKeyType.polish,
        label: '英文润色改写',
        i18nKey: 'common.aitools.enAiPolish',
        component: PolishPage,
        icon: PolishIcon,
        activeIcon: PolishActiveIcon,
      },
      {
        key: TabKeyType.translate,
        label: '中译英',
        i18nKey: 'common.aitools.aiTranslate',
        component: TranslatePage,
        icon: TranslateIcon,
        activeIcon: TranslateActiveIcon,
      },
      {
        key: TabKeyType.reviewer,
        label: 'AI论文导师',
        i18nKey: 'common.aitools.aiReviewer',
        component: RecentPage,
        icon: 'icon-reviewer-label',
        activeIcon: '',
      },
      {
        key: TabKeyType.reviewerqs,
        label: 'AI毕业论文审稿',
        badge: 'new',
        i18nKey: 'common.aitools.aiReviewerQS',
        component: ReviewQSPage,
        icon: 'icon-reviewer-label',
        activeIcon: '',
      },
    ];

const activeTabKey = useLocalStorage<TabKeyType>(
  'polish/ai-paragraph-client-tab',
  TAB_LIST[0].key
);

const isAboveProfessional = computed(
  () =>
    compareVipType(vipStore.role.vipType, VipType.PROFESSIONAL) >= 0 ||
    vipStore.role.vipType === VipType.ENTERPRISE
);
const activeTabTrialCount = computed(() => {
  const count =
    trialCounts.value?.trialFeatureAndCounts.find(
      (x) =>
        activeTabKey.value ===
        (
          {
            [TrialFeature.POLISH_REWRITE]: TabKeyType.polish,
            [TrialFeature.ZH_TO_EN_TRANSLATION]: TabKeyType.translate,
            [TrialFeature.UNRECOGNIZED]: 'unrecoginized',
            [TrialFeature.ZH_POLISH_REWRITE]: TabKeyType.zhpolish,
            [TrialFeature.AI_TRANSLATE]: TabKeyType.translate,
          } as Record<string, string>
        )[x.feature]
    )?.trialCount ?? 0;

  if (activeTabKey.value === TabKeyType.zhpolish && isAboveProfessional.value) {
    // 中文润色，同时已经是专业会员，不显示试用次数
    return 0;
  }

  if (count > 0) {
    activeTabRefundable.value = false;
  }

  return count;
});
const activeTabRefundable = ref(true);

provide('refundable', activeTabRefundable);

const componentId = shallowRef(
  (TAB_LIST.find((tab) => tab.key === activeTabKey.value) || TAB_LIST[0])
    .component
);

watch(
  () => router.currentRoute.value.query.tab,
  (tab) => {
    if (tab) {
      const tabItem = TAB_LIST.find((item) => item.key === tab);
      if (tabItem) {
        activeTabKey.value = tabItem.key;
        componentId.value = tabItem.component;
      }
    }
  },
  {
    immediate: true,
  }
);

const handleChangeTab = (tab: (typeof TAB_LIST)[0]) => {
  router.push({
    query: {
      tab: tab.key,
    },
  });
};

const isActiveTab = (tab: (typeof TAB_LIST)[0]) => {
  return activeTabKey.value === tab.key;
};
</script>
<template>
  <div class="flex h-full bg-[#F2F4F7]">
    <Drawer
      class="bg-[#F5F4F2] border border-r border-rp-neutral-3"
      :initial-width="256"
      :max-width="256"
      :min-width="256"
      placement="left"
    >
      <div class="px-1 py-2">
        <div
          v-for="tab in TAB_LIST"
          :key="tab.key"
          :class="[
            'item px-4 py-3 text-rp-neutral-10 font-medium text-sm space-x-2 cursor-pointer flex items-center rounded',
            isActiveTab(tab) ? 'active' : '',
          ]"
          @click="handleChangeTab(tab)"
        >
          <i
            v-if="/^icon-/.test(tab.icon)"
            :class="['aiknowledge-icon text-rp-neutral-6', tab.icon]"
          />
          <img
            v-else
            :src="isActiveTab(tab) ? tab.activeIcon : tab.icon"
          >
          <span>{{ $t(tab.i18nKey) }}</span>
          <div
            v-if="tab.badge"
            class="px-2 py-0.5 rounded-[2px] bg-rp-blue-2 text-rp-blue-6 text-xs font-normal"
          >
            <span>{{ tab.badge }}</span>
          </div>
          <!-- <div class="absolute right-0 top-1">
            <a-badge-ribbon
              v-if="tab.badge"
              :text="tab.badge"
              color="red"
            ></a-badge-ribbon>
          </div> -->
        </div>
      </div>
    </Drawer>
    <div class="flex-1 relative">
      <LoginFirst>
        <template #blank>
          <img
            class="w-20"
            src="@common/../assets/images/aitools/empty.svg"
            alt="empty"
          >
          <p>
            {{ $t('common.aitools.loginFirst') }}
          </p>
        </template>
        <keep-alive>
          <component
            :is="componentId"
            :key="activeTabKey"
            :allowed="hasPermission"
            :fromClient="fromClient"
            :openInNewTab="openInNewTab"
            :isModuleImpression="isModuleImpression"
            :type="activeTabKey"
            @intercepted="checkPermission"
            @started="refreshTrialCount"
          >
            <template
              v-if="
                vipStore.enabled &&
                  !trialLoading &&
                  activeTabTrialCount !== undefined &&
                  activeTabTrialCount <= 0
              "
              #beans
            >
              <div class="flex items-center whitespace-nowrap">
                <!-- <template v-if="activeTabTrialCount > 0">
                  <span class="text-rp-neutral-8"
                    >{{ $t(`common.aitools.trialBalance`)
                    }}{{ $t('common.symbol.colon') }}</span
                  >{{ activeTabTrialCount }}
                </template> -->
                <BetaCard :tab="activeTabKey">
                  <UnlimitedCard>
                    <span class="text-rp-neutral-8">{{ $t(`common.aitools.aiBeanBalance`)
                    }}{{ $t('common.symbol.colon') }}</span>
                    <AiBeanAuth :lowest-vip-type="VipType.PROFESSIONAL" />
                  </UnlimitedCard>
                </BetaCard>
              </div>
            </template>
          </component>
        </keep-alive>
        <ApplyModal />
      </LoginFirst>
    </div>
  </div>
</template>
<style lang="less" scoped>
.active {
  background-color: #438de6;
  color: #fff;
  .aiknowledge-icon {
    color: #fff;
  }
}
.item {
  transition: all 0.3s;
  &:not(.active):hover {
    background-color: theme('colors.rp-neutral-3');
  }
}
</style>
