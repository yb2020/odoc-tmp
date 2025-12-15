<template>
  <div>
    <a-button
      v-if="showSetting"
      type="link"
      @click="openTranslateConfigSetting"
    >
      <template #icon>
        <SettingOutlined />
      </template>
      <!-- 管理翻译渠道 -->
    </a-button>

    <settingView v-model:visible="setting" />
  </div>
</template>
<script lang="ts" setup>
import { onUnmounted } from 'vue';
import { SettingOutlined } from '@ant-design/icons-vue';

import {
  useTranslateStore,
  TranslateTabKey,
  GoogleTranslateStyle,
} from '~/src/stores/translateStore';
import settingView from '@/components/Tippy/config/TranslateSetting.vue';
import { gteElectronVersion } from '~/src/util/env';
import {
  emitter,
  CONFIG_RESET_TYPE,
  CONFIG_ADD_TYPE,
} from '@/components/Tippy/config/config';

defineProps<{
  currentTab: string;
}>();

const emit = defineEmits<{
  (event: 'addTab', payload: { title: string; type: TranslateTabKey }): void;
  (event: 'deleteTab', payload: { type: TranslateTabKey }): void;
}>();

const showSetting = gteElectronVersion('1.20.1');

const setting = defineModel('visible', { default: false });

// const setting = ref(false);

const store = useTranslateStore();
const addAddChannelListener = (e: any) => {
  const configInfo = JSON.parse(JSON.stringify(e));
  configInfo.verified = true;
  let type = TranslateTabKey.other;
  if (configInfo.txSecretId?.length > 0) {
    /// 腾讯翻译
    configInfo.type = TranslateTabKey.tencent;
    store.txConfig = configInfo;
    type = TranslateTabKey.tencent;
  }
  if (configInfo.aliAccessKeyId?.length > 0) {
    /// 阿里翻译
    configInfo.type = TranslateTabKey.ali;
    store.aliConfig = configInfo;
    type = TranslateTabKey.ali;
  }
  if (
    configInfo.googleApiKey?.length > 0 ||
    store.googleConfigVersion !== GoogleTranslateStyle.none
  ) {
    /// 谷歌翻译
    configInfo.type = TranslateTabKey.google;
    store.googleConfig = configInfo;
    type = TranslateTabKey.google;
  }
  if (configInfo.deepLKey?.length > 0) {
    /// deepl
    configInfo.type = TranslateTabKey.deepl;

    store.deeplConfig = configInfo;
    type = TranslateTabKey.deepl;
  }

  emit('addTab', {
    title: configInfo.name,
    type: configInfo.type,
  });

  // if (!titles[type]) {
  //   tabs.value.push({
  //     title: configInfo.name,
  //     type: configInfo.type,
  //   });
  // }

  // initTabs();
};

emitter.on(CONFIG_ADD_TYPE, addAddChannelListener);

const resetChannelListtener = (e: any) => {
  if (!e) {
    return;
  }

  const deleteInfo = JSON.parse(JSON.stringify(e));
  emit('deleteTab', {
    type: deleteInfo.type,
  });
  // delete titles[deleteInfo.type as TranslateTabKey];
  // currentTab.value = DEFAULT_TRANSLATE_TAB_KEY;
  // for (let i = 0; i < tabs.value.length; i++) {
  //   if (tabs.value[i].type === deleteInfo.type) {
  //     tabs.value.splice(i, 1);
  //     initTabs();
  //     break;
  //   }
  // }
};
emitter.on(CONFIG_RESET_TYPE, resetChannelListtener);

onUnmounted(() => {
  emitter.off(CONFIG_ADD_TYPE, addAddChannelListener);
  emitter.off(CONFIG_RESET_TYPE, resetChannelListtener);
});

const openTranslateConfigSetting = () => {
  setting.value = true;
};
</script>
