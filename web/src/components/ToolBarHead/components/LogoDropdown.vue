<template>
  <a-dropdown
    :trigger="['click']"
    placement="bottomRight"
  >
    <div class="text-base cursor-pointer">
      <MoreOutlined />
    </div>
    <template #overlay>
      <a-menu>
        <a-menu-item @click="reportClick(ElementClick.shortcut_key)">
          <a
            ref="shortcutsTippyTriggerRef"
            href="javascript:;"
            @click="toggleShortcutsHelp"
          >{{ $t('viewer.shortcut') }}</a>
        </a-menu-item>
        <!-- 旧版本反馈 -->
        <!-- <a-menu-item @click="reportClick(ElementClick.suggest)">
          <a
            ref="feedbackTippyTriggerRef"
            href="javascript:;"
          >{{
            $t('viewer.help')
          }}</a>
        </a-menu-item> -->
        <!-- 旧版本tips -->
        <!-- <a-menu-item>
          <a
            href="https://docs.qq.com/doc/DVWNURG9BU1R5ZW1T"
            target="_blank"
          >{{
            $t('viewer.tips')
          }}</a>
        </a-menu-item> -->
        <!-- 新版本反馈 -->
        <a-menu-item>
          <a 
            :href="userGuideLink" 
            target="_blank" 
            class="user-guide-link"
          >
            {{ $t('viewer.help') }}
          </a>
        </a-menu-item>
      </a-menu>
    </template>
  </a-dropdown>
  <TippyVue
    v-if="feedbackTippyTriggerRef"
    ref="tippyVueRef"
    :trigger-ele="feedbackTippyTriggerRef"
    :placement="'left-start'"
    :trigger="'click'"
  >
    <Feedback @close="close" />
  </TippyVue>

  <ShortcutsHelp
    :visible="visibleShortcutsHelp"
    @close="toggleShortcutsHelp"
  />
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { MoreOutlined } from '@ant-design/icons-vue';

import TippyVue from '../../Tippy/index.vue';
import Feedback from './Feedback.vue';
import ShortcutsHelp from './ShortcutsHelp.vue';
import { ElementClick, reportClick } from '~/src/api/report';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';

const feedbackTippyTriggerRef = ref<HTMLAnchorElement>();
const tippyVueRef = ref();

const close = () => {
  tippyVueRef.value.hide();
};

const visibleShortcutsHelp = ref(false);
const toggleShortcutsHelp = () => {
  visibleShortcutsHelp.value = !visibleShortcutsHelp.value;
};

const { locale } = useI18n();

// 语言管理
const { isCurrentLanguage } = useLanguage();

// 创建用户指南链接的计算属性
const userGuideLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide';
  }
  return '/docs/guide';
});
</script>

<style lang="less" scoped>
.icon-logo-readpaper {
  font-size: 24px;
  margin-right: 8px;
}
</style>
