<template>
  <div
    class="aibeans-way w-[295px] pb-4 rounded-lg border border-solid border-rp-neutral-3"
  >
    <h3
      class="mb-4 relative h-12 flex justify-center items-center text-base font-medium"
    >
      <span class="relative z-10 inline-flex items-center text-white">
        <slot name="title">
          <IconBean
            class="mr-2"
            fill="white"
          />
          {{ $t('common.aibeans.buyPkg') }}
        </slot>
      </span>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 295 48"
        class="w-full absolute top-0 z-0"
        :fill="fill || '#47CCC7'"
      >
        <path
          d="M0 4C0 1.79086 1.79086 0 4 0H291C293.209 0 295 1.79086 295 4V30.6404C295 32.706 293.427 34.432 291.37 34.6232L147.87 47.9656C147.624 47.9885 147.376 47.9885 147.13 47.9656L3.62968 34.6232C1.57289 34.432 0 32.706 0 30.6404V4Z"
        />
      </svg>
    </h3>
    <p class="mb-4 flex justify-center items-baseline h-10">
      <slot name="beans">
        <span class="text-rp-red-6 text-[40px] font-bold leading-[1]">{{ formatPrice(price ?? 3000) }}{{ $t('common.units.yen') }}</span><span class="text-sm leading-[1]">{{ $t('common.aibeans.perPkg') }}({{
          $t('common.aibeans.tipPkg', [beans])
        }})</span>
      </slot>
    </p>
    <p class="mb-4 text-sm text-center py-px">
      <slot name="desc">
        <span class="text-rp-neutral-6 text-xs">({{ $t('common.aibeans.tipPkg2', [expireDays]) }})</span>
      </slot>
    </p>
    <slot name="qrcode">
      <div class="w-40 h-40">
        <a-spin />
      </div>
    </slot>

    <PaySupport class="mt-4 pt-px" />
  </div>
</template>

<script setup lang="ts">
import IconBean from './Icon.vue';
import PaySupport from '../Pay/Support.vue';
import { formatPrice } from '../../utils/pay';

defineProps<{
  fill?: string;
  beans?: number;
  price?: number;
  expireDays?: number;
}>();
</script>
