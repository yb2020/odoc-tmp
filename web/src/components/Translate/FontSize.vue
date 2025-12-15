<template>
  <div
    class="flex items-center relative"
    @click="showSelect = true"
  >
    <div class="text-sm text-rp-neutral-8 cursor-pointer">
      {{ $t('translate.fontSize') }}<span class="ml-0.5">{{ fontSize }}</span><CaretDownOutlined />
    </div>
    <SelectOptions
      v-model:visible="showSelect"
      class="w-full top-6 h-[180px] z-10"
      :options="[]"
      @selectChange="changePage"
    >
      <template #options>
        <div class="bg-white text-rp-neutral-10">
          <div
            v-for="item in pageOptions"
            :key="item.value"
            :class="[
              'p-1 hover:bg-rp-neutral-2',
              fontSize === item.value ? 'bg-rp-neutral-2' : '',
            ]"
            :style="{ fontSize: item.value + 'px' }"
            @click="changePage(item.value)"
          >
            {{ item.title }}
          </div>
        </div>
      </template>
    </SelectOptions>
  </div>
</template>
<script lang="ts" setup>
import SelectOptions from '@/components/ToolBar/components/Select.vue';
import { ref } from 'vue';
import { useTranslateFontSize } from '@/hooks/useTranslation';
import { CaretDownOutlined } from '@ant-design/icons-vue';

const { fontSize } = useTranslateFontSize();

const showSelect = ref(false);

const changePage = (value: string) => {
  fontSize.value = value;
};
const pageOptions = [
  { title: '12', value: '12' },
  { title: '14', value: '14' },
  { title: '16', value: '16' },
  { title: '18', value: '18' },
  { title: '20', value: '20' },
  { title: '22', value: '22' },
  { title: '24', value: '24' },
  { title: '26', value: '26' },
];
</script>
