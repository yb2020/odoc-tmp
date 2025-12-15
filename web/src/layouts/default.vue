<template>
  <a-config-provider :locale="isEnUS ? enUS : zhCN">
    <a-layout style="height: 100%">
      <LayoutHeader />
      <a-layout-content :style="{ flex: 1 }">
        <router-view />
      </a-layout-content>
    </a-layout>
    <!-- 移除登录弹窗，改为使用路由导航到登录页面 -->
    <CertDialog v-if="!bannedDialogVisible" />
    <BannedDialog
      v-if="banInfo && bannedDialogVisible"
      :banInfo="banInfo"
    />
  </a-config-provider>
</template>
<script lang="ts" setup>
import LayoutHeader from '@/components/Head/index.vue';
// 移除 Login 组件的导入
import CertDialog from '@/components/Common/Cert.vue';
import BannedDialog from '@/components/Common/bannedDialog.vue';

import enUS from 'ant-design-vue/es/locale/en_US';
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import dayjs from 'dayjs';

import 'dayjs/locale/zh-cn';
import { computed } from 'vue';
import { store } from '~/src/store';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

dayjs.locale('zh-cn');

// 语言管理
const { isCurrentLanguage } = useLanguage();

const banInfo = computed(() => store.state.user.userInfo?.banInfo);

const bannedDialogVisible = computed(() => banInfo.value?.banFlag);

// 基于 proto 枚举的语言判断计算属性
const isEnUS = computed(() => isCurrentLanguage(Language.EN_US));
</script>

<style lang="less">
body {
  background-color: #282b2e;
}
</style>
