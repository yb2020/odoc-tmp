<template>
  <div
    class="fullpage-container"
    @click="handleClick"
  >
    <a-tooltip>
      <template #title>
        {{ $t('viewer.fullscreenMode') }}<component :is="getShortcutTxt(shortcut)" />
      </template>
      <FullscreenExitOutlined
        v-if="isFullPage"
        class="icon"
      />
      <FullscreenOutlined
        v-else
        class="icon"
      />
    </a-tooltip>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { useStore } from '@/store';
import { DocumentsMutationTypes } from '~/src/store/documents';
import screenfull from 'screenfull';
import { FullscreenExitOutlined, FullscreenOutlined } from '@ant-design/icons-vue';
import { PAGE_ROUTE_NAME } from '../../../routes/type';
import useShortcuts from '../../../hooks/useShortcuts';
import { getPlatformKey, getShortcutTxt } from '../../../store/shortcuts';
import { ElementClick, reportClick } from '~/src/api/report';

export default defineComponent({
  components: {
    FullscreenExitOutlined,
    FullscreenOutlined,
  },
  setup() {
    const store = useStore();
    const isFullPage = computed(() => store.state.documents.isFullPage);

    const handleClick = () => {
      if (screenfull.isEnabled) {
        screenfull
          .toggle()
          .then(() => {
            store.commit(`documents/${DocumentsMutationTypes.SET_FULL_PAGE}`, !isFullPage.value);
            
          })
          .catch((e) => {
            message.info('请聚焦页面后使用全屏快捷键~');
          });
      }
    };

    onMounted(() => {
      if (screenfull.isEnabled) {
        screenfull.on('change', () => {
          reportClick(ElementClick.full_screen, isFullPage.value ? 'on' : 'off')
          if (!(screenfull as any).isFullscreen) {
            store.commit(`documents/${DocumentsMutationTypes.SET_FULL_PAGE}`, false);
          }
        });
      }
    });

    const shortcutsConfig = computed(() => store.state.shortcuts[PAGE_ROUTE_NAME.NOTE] || {});
    const platformKey = getPlatformKey();
    const shortcut = computed(() => shortcutsConfig.value.shortcuts.FULLSCREEN.value[platformKey]);
    const opts = computed(() => ({ scope: shortcutsConfig.value.scope || 'all' }));
    const handler = (e: KeyboardEvent) => {
      
      e.preventDefault();
      handleClick();
    };
    useShortcuts(shortcut, handler, opts);

    return {
      handleClick,
      isFullPage,
      shortcut,
      getShortcutTxt,
    };
  },
});
</script>

<style lang="less" scoped>
.fullpage-container {
  display: flex;
  align-items: center;
  justify-content: center;

  width: 36px;
  height: 36px;
  border-radius: 2px;
  cursor: pointer;

  &:hover {
    background: #f5f5f5;
  }

  .icon {
    color: rgba(0, 0, 0, 64%);
    font-size: 16px;
  }
}
</style>
