<template>
  <a-modal
    :visible="visible && unpolishable"
    :closable="closable ?? false"
    :footer="null"
    :getContainer="() => container"
    @cancel="visible = false"
  >
    <div class="text-rp-neutral-10">
      <h3 class="text-base mb-2">
        {{ $t('common.aitools.vipLimit.ruleTt') }}
      </h3>
      <ol class="text-sm leading-[22px]">
        <li v-for="x in $t('common.aitools.vipLimit.ruleTxt').split('\n')">
          {{ x }}
        </li>
      </ol>
      <h3 class="text-base mb-2">
        {{ $t('common.aitools.vipLimit.exampleTt') }}
      </h3>
      <p>
        {{ $t('common.aitools.vipLimit.exampleTxt')
        }}<a
          class="text-rp-blue-6 underline"
          href="javascript:;"
          @click="openBeginnerGuide(isCurrentLanguage(Language.EN_US))"
        >{{ $t('common.aitools.guide') }}</a>
      </p>
      <div class="flex flex-row-reverse">
        <Trigger
          visible
          :btnProps="{ type: 'primary' }"
          :btnTxt="`${$t('common.premium.btns.get')} ${$t(
            `common.premium.versions.${vipPrefs.key}`
          )}`"
          :needVipType="vipType"
          :report-params="{
            pageType: PageType.POLISH,
            elementName: ElementName.upper_ai_polish_mentor,
          }"
        />
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import {
  PremiumVipPreferences,
  useVipStore,
  VipType,
} from '@common/stores/vip';
import Trigger from '@common/components/Premium/Trigger.vue';
import { openBeginnerGuide } from '@common/components/AIRevise/utils';
import { computed, onMounted } from 'vue';
import { PageType, ElementName } from '@common/utils/report';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

const visible = defineModel('visible', { default: false });

// 语言管理
const { isCurrentLanguage } = useLanguage();

const props = defineProps<{
  closable?: boolean;
  container?: HTMLElement;
}>();

const container = computed(() => props.container || document.body);

const vipType = VipType.PROFESSIONAL as const;
const vipPrefs = PremiumVipPreferences[vipType];
const vipStore = useVipStore();

const unpolishable = computed(() => {
  return (
    Object.keys(vipStore.privileges).length > 0 &&
    !vipStore.privileges.polishEnable &&
    vipStore.inited &&
    vipStore.enabled
  );
});

onMounted(() => {
  vipStore.fetchVipProfile();
});
</script>
