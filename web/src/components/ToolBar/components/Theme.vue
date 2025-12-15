<template>
  <a-dropdown :trigger="['click']" placement="top" overlayClassName="theme-overlay">
    <div class="relative" :class="[
        'theme-container',
        { 'red-dot-container': baseStore.colorRedDot },
      ]" @click="handleClickThemeIcon">
      <a-tooltip>
        <template #title>
          {{ $t('viewer.readMode') }}
        </template>
        <a-badge :dot="baseStore.colorRedDot" color="#FAAD14">
          <img :src="colorIconImgUrl" class="w-4" alt="">
        </a-badge>
      </a-tooltip>
    </div>

    <template #overlay>
      <a-menu>
        <div class="theme-popover">
          <div class="title">
            {{ $t('viewer.theme') }}
          </div>
          <a-row type="flex" :gutter="[12, 12]">
            <a-col v-for="item in themeConfigs" :key="item.value" flex="1" @click="onChangeTheme(item.value as any)">
              <div :class="[
                  'color-box',
                  'no-theme',
                  curTheme === item.value ? 'current' : '',
                ]" :style="{ background: item.bg }" />
              <a-radio :value="item.value" :checked="curTheme === item.value">
                {{
                $t(item.i18n)
                }}
              </a-radio>
            </a-col>
          </a-row>
        </div>
      </a-menu>
    </template>
  </a-dropdown>
</template>

<script lang="ts" setup>
import { computed, watch } from 'vue';
import { ThemeType, changeReaderTheme } from '@/theme';
import { useLocalStorage } from '@vueuse/core';
import { message } from 'ant-design-vue';
import { useFullTextTranslateStore } from '~/src/stores/fullTextTranslateStore';
import { useBaseStore } from '~/src/stores/baseStore';
import colorIconImgUrl from '@/assets/images/color_icon.png';
import { clearRedDot } from '~/src/api/user';
import { FunctionTypeEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/RedDot';
import { PageType, reportElementClick } from '~/src/api/report';
import { useI18n } from 'vue-i18n';
import { THEME } from '~/src/common/src/constants/storage-keys';

const props = defineProps<{
  onChangeTheme: (x: ThemeType | 'default') => void;
}>();

const themeConfigs = computed(() => {
  return [
    {
      value: 'default',
      name: '原生',
      i18n: 'viewer.white',
      bg: '#fff',
    },
    {
      value: ThemeType.dark,
      name: '深色',
      i18n: 'viewer.dark',
      bg: '#282B2E',
    },
    {
      value: ThemeType.beige,
      name: '米色',
      i18n: 'viewer.beige',
      bg: '#F0E6D3',
    },
    {
      value: ThemeType.green,
      name: '绿色',
      i18n: 'viewer.green',
      bg: '#E8EDE4',
    },
  ];
});

const curTheme = useLocalStorage(THEME.SITE, 'default');

const curBg = computed(() => {
  const cur = themeConfigs.value.find((item) => item.value === curTheme.value);
  return cur?.bg || '#fff';
});

watch(
  () => curTheme.value,
  (newVal, oldVal) => {
    try {
      changeReaderTheme(newVal as 'default' | ThemeType);
      props.onChangeTheme?.(newVal as 'default' | ThemeType);
    } catch (error) {}
  }
);

const fullTextTranslateStore = useFullTextTranslateStore();

watch(
  () => fullTextTranslateStore.pdfId,
  (newVal) => {
    if (newVal && curTheme.value === ThemeType.dark) {
      message.error('全文翻译暂不支持深色模式');
      onChangeTheme('default');
    }
  }
);

const { t } = useI18n();

const onChangeTheme = (val: 'default' | ThemeType) => {
  if (val === curTheme.value) {
    return;
  }

  if (val === ThemeType.dark && fullTextTranslateStore.pdfId) {
    message.error('全文翻译暂不支持深色模式');
    return;
  }

  const tips = {
    default: t('message.turnDownReadModeTip'),
    [ThemeType.beige]: t('message.turnUpReadModeTip', {
      mode: t('viewer.beige'),
    }),
    [ThemeType.green]: t('message.turnUpReadModeTip', {
      mode: t('viewer.green'),
    }),
    [ThemeType.dark]: t('message.turnUpReadModeTip', {
      mode: t('viewer.dark'),
    }),
  };
  message.success(tips[val]);
  curTheme.value = val;
};

const baseStore = useBaseStore();

const handleClickThemeIcon = () => {
  reportElementClick({
    page_type: PageType.note,
    type_parameter: 'none',
    element_name: 'color_switch' as any,
  });

  if (baseStore.colorRedDot) {
    baseStore.setColorRedDot(false);

    clearRedDot({ functionType: FunctionTypeEnum.COLOR_PATTERN });
  }
};
</script>

<style lang="less" scoped>
.theme-container {
  .ant-badge {
    :deep(.ant-badge-dot) {
      box-shadow: none;
    }
  }
}

.theme-overlay {
  .ant-dropdown-menu {
    padding: 0;
    box-shadow: 0px 8px 12px 0px rgba(0, 0, 0, 0.16);
  }
}

.theme-popover {
  width: 428px;
  padding: 8px 12px;
  margin-bottom: 10px;
  color: #4e5969;
  border-radius: 4px;
  background-color: #fff;

  .title {
    font-weight: 600;
  }

  :deep(.ant-radio-wrapper) {
    color: #4e5969;

    .ant-radio-inner {
      box-shadow: none;
      border: 1px solid #c9cdd4;
    }
  }

  .color-box {
    height: 56px;
    border: 2px solid #c9cdd4;
    border-radius: 2px;
    margin: 16px 0 8px;
    cursor: pointer;

    &.current {
      border-color: #1f71e0;
    }
  }
}

// html[data-theme="dark"] {
//   .theme-popover {
//     .title {
//       color: #C9CDD4;
//     }
//   }
// }

.red-dot-container {
  background: rgba(255, 255, 255, 0.08);
}
</style>
