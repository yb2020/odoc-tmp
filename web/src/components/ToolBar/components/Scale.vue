<template>
  <div
    class="js_select relative text-xs leading-3"
    @click.stop="() => (showSelect = !showSelect)"
  >
    <div class="flex items-center justify-center">
      <input
        v-model="inputValue"
        class="scale-input min-w-[50px] outline-0 border-0 bg-transparent"
        :size="inputValue.length"
        @blur="handleOptionsBlur"
        @keyup="handleKeydown"
      >
      <i
        class="aiknowledge-icon icon-arrow-down text-inherit"
        aria-hidden="true"
      />
    </div>
    <SelectOptions
      v-model:visible="showSelect"
      :options="ScaleOptions"
      class="scale-select"
      @selectChange="changeScale"
    >
      <div class="flex items-center justify-center">
        <div
          class="decrease"
          @click="handleDecrease"
        >
          <minus-circle-outlined class="icon" />
        </div>
        <div
          class="increase"
          @click="handleIncrease"
        >
          <plus-circle-outlined class="icon" />
        </div>
      </div>
    </SelectOptions>
  </div>
</template>

<script lang="ts" setup>
export type ToolbarScaleEvent =
  | {
      type: 'toolbar:scale';
      scaleValue: number | 'page-width';
    }
  | {
      type: 'toolbar:scale:increase';
    }
  | {
      type: 'toolbar:scale:decrease';
    };

import { computed, ref, watch } from 'vue';
import SelectOptions from './Select.vue';
import hotkeys from 'hotkeys-js';
import { MinusCircleOutlined, PlusCircleOutlined } from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import trim from 'lodash-es/trim';
import { isInteger } from 'lodash-es';
import { useStore } from '@/store';
import useShortcuts from '../../../hooks/useShortcuts';
import { PAGE_ROUTE_NAME } from '../../../routes/type';
import { getPlatformKey } from '../../../store/shortcuts';
import { useI18n } from 'vue-i18n';
import { ElementClick, reportClick } from '~/src/api/report';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { languageEnumToStandard } from '~/src/shared/language/service';

const props = defineProps<{
  scale: number;
  scalePresetValue: string;
}>();

const emit = defineEmits<{
  (event: 'changeScale', payload: ToolbarScaleEvent): void;
}>();

const store = useStore();
const shortcutsConfig = computed(
  () => store.state.shortcuts[PAGE_ROUTE_NAME.NOTE] || {}
);
const platformKey = getPlatformKey();
const shortcut = computed(
  () =>
    shortcutsConfig.value.shortcuts.PDF_SIZE_SELF_ADAPTION.value[platformKey]
);
const opts = computed(() => ({ scope: shortcutsConfig.value.scope || 'all' }));
const handler = () => {
  changeScale('page-width');
};

const { t, locale } = useI18n();

const ScaleOptions = computed(() => {
  return [
    {
      title: '50%',
      value: '0.5',
    },
    {
      title: '70%',
      value: '0.7',
    },
    {
      title: '100%',
      value: '1.0',
    },
    {
      title: '200%',
      value: '2.0',
    },
    {
      title: '400%',
      value: '4.0',
    },
    {
      title: t('viewer.pageWidth'),
      value: 'page-width',
      shortcut: shortcut.value,
    },
  ];
});

const getShowValue = (scale: number, presetValue: string) => {
  if (presetValue === 'page-width') {
    return t('viewer.pageWidth');
  }
  return Math.floor(scale * 100) + '%';
};

const isValidInput = () => {
  const val = trim(inputValue.value);
  const scale = parseInt(inputValue.value);

  if (val === t('viewer.pageWidth')) {
    return 'page-width';
  }

  if (isNaN(scale)) {
    message.error(t('message.scaleTip1'));
    return 0;
  }
  if (isInteger(scale)) {
    if (scale < 50 || scale > 400) {
      message.error(t('message.scaleTip2'));
      return 0;
    }
  }
  return scale;
};

const handleKeydown = (e: KeyboardEvent) => {
  if (e.keyCode == 13) {
    const scale = isValidInput();

    if (!scale) {
      return;
    }
    emit('changeScale', {
      type: 'toolbar:scale',
      scaleValue: scale === 'page-width' ? scale : scale / 100,
    });
  }
};

const handleOptionsBlur = () => {
  const scale = isValidInput();

  if (!scale) {
    inputValue.value = getShowValue(props.scale, props.scalePresetValue);
    return;
  }

  emit('changeScale', {
    type: 'toolbar:scale',
    scaleValue: scale === 'page-width' ? scale : scale / 100,
  });
};

const inputValue = ref<string>(
  getShowValue(props.scale, props.scalePresetValue)
);

watch(
  () => ({
    scale: props.scale,
    scalePresetValue: props.scalePresetValue,
  }),
  (val) => {
    inputValue.value = getShowValue(val.scale, val.scalePresetValue);
  }
);

watch(locale, () => {
  if (
    [
      t('viewer.pageWidth', languageEnumToStandard(Language.EN_US)),
      t('viewer.pageWidth', languageEnumToStandard(Language.ZH_CN)),
    ].includes(inputValue.value.trim())
  ) {
    inputValue.value = t('viewer.pageWidth');
  }
});

const changeScale = (scale: string | 'page-width') => {
  emit('changeScale', {
    type: 'toolbar:scale',
    scaleValue: scale === 'page-width' ? scale : Number(scale),
  });
  reportClick(ElementClick.size_adjust);
};

const handleDecrease = () => {
  // changePDFScaleBySteps('decrease', 1);
  emit('changeScale', { type: 'toolbar:scale:decrease' });
  reportClick(ElementClick.size_auto_narrow);
};

const handleIncrease = () => {
  emit('changeScale', { type: 'toolbar:scale:increase' });
  // changePDFScaleBySteps('increase', 1);
  reportClick(ElementClick.size_auto_large);
};

const showSelect = ref<boolean>(false);

hotkeys('ctrl+=, command+=', () => {
  handleIncrease();
  return false;
});

hotkeys('ctrl+-, command+-', () => {
  handleDecrease();
  return false;
});
useShortcuts(shortcut, handler, opts);
</script>

<style lang="less" scoped>
.decrease,
.increase {
  cursor: pointer;
  width: 36px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    background: #f5f5f5;
  }

  .icon {
    color: rgba(0, 0, 0, 64%);
  }
}

.scale-select {
  :deep(.option) {
    text-align: left;
  }
}
</style>
