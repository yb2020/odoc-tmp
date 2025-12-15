<template>
  <section class="error-page">
    <div class="content">
      <div class="logo-container">
        <div class="slogan">vibe reading</div>
      </div>
      <h1 class="title">
        {{ code || '404' }}
      </h1>
      <p class="message">{{ decodeURIComponent((message || 'Not Found') as string) }}</p>
      <a-button
        type="link"
        class="home-button"
        @click="handler"
      >
        {{ url ? $t('viewer.refresh') : $t('viewer.goHomePage') }}
      </a-button>
    </div>
  </section>
</template>
<script lang="ts" setup>
import { onMounted } from 'vue';
import { initTheme } from '@common/theme';

// 确保主题正确初始化
onMounted(() => {
  if (typeof window !== 'undefined') {
    // 初始化主题
    initTheme();
  }
});
const props = defineProps<{
  code?: number | string;
  message?: string;
  url?: string;
}>()

const handler = () => {
  if(props.url == undefined) {
    window.location.href = '/';
    return;
  }
  if (props.url) {
    window.location.replace(decodeURIComponent(props.url));
    return;
  }
  window.location.reload();
};
</script>

<style scoped lang="less">
/* 导入全局主题变量 */
@import '../../assets/less/theme.less';

.error-page {
  background: var(--site-theme-bg-primary);
  color: var(--site-theme-text-primary);
  min-height: 100vh;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
  
  .content {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;
    max-width: 90vw;
  }

  .logo-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 2rem;
  }

  .slogan {
    font-size: 4rem;
    font-weight: bold;
    line-height: 1.2;
    background: linear-gradient(90deg, #6889ff 0%, #8662e9 50%, #e74694 100%);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
    white-space: nowrap;
    
    @media (max-width: 768px) {
      font-size: 3.5rem;
    }
    
    @media (max-width: 480px) {
      font-size: 2.5rem;
    }
  }

  .title {
    font-size: 8rem;
    font-weight: 500;
    color: var(--site-theme-text-primary);
    line-height: 1;
    margin: 0 0 1rem;
    
    @media (max-width: 768px) {
      font-size: 6rem;
    }
    
    @media (max-width: 480px) {
      font-size: 4rem;
    }
  }

  .message {
    font-size: 1.2rem;
    margin-bottom: 2rem;
    color: var(--site-theme-text-secondary);
  }
  
  .home-button {
    font-size: 1.1rem;
    color: var(--site-theme-brand);
    transition: all 0.3s;
    
    &:hover {
      color: var(--site-theme-brand-light);
      transform: translateY(-2px);
    }
  }
}
</style>
