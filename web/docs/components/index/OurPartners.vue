<template>
  <div v-if="partnersData" class="our-partners">
    <div class="partners-background">
      <div class="container">
        <h2 class="main-title">{{ partnersData.title }}</h2>
        
        <div class="partners-grid">
          <div
            v-for="(partner, key) in partnersData.partners"
            :key="key"
            class="partner-logo"
          >
            <img 
              :src="getImageUrl(`/images/icons/${partner.logo}`)" 
              :alt="partner.name"
              class="logo-image"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import docsI18nInstance, { useVitePressI18n } from '../../i18n/index';
import { getImageUrl } from '../../utils/contextPath';

const { locale } = useVitePressI18n();

const messages = computed(() => docsI18nInstance.global.messages.value[locale.value] || {});
const partnersData = computed(() => (messages.value.OurPartners as any) || null);
</script>

<style scoped>
.our-partners {
  position: relative;
  width: 100vw;
  margin-left: calc(-50vw + 50%);
  background-color: var(--vp-c-bg-soft);
}

.partners-background {
  width: 100%;
  padding: 80px 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.main-title {
  font-size: 48px;
  font-weight: 700;
  line-height: 1.2;
  color: var(--vp-c-text-1);
  text-align: center;
  margin: 0 0 60px 0;
  border: none;
}

.partners-grid {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.partner-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.8;
  transition: all 0.3s ease;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 8px 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

html.dark .partner-logo {
  background-color: rgba(255, 255, 255, 0.9);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.partner-logo:hover {
  opacity: 1;
}

.logo-image {
  height: 40px;
  width: 120px;
  object-fit: contain;
  filter: grayscale(100%);
  transition: filter 0.3s ease;
}

html.dark .logo-image {
  /* 在白色背景上使用正常灰度滤镜 */
  filter: grayscale(100%);
  -webkit-filter: grayscale(100%);
}

.partner-logo:hover .logo-image {
  filter: grayscale(0%);
}

html.dark .partner-logo:hover .logo-image {
  filter: grayscale(0%);
  -webkit-filter: grayscale(0%);
}

/* 深色模式下的logo处理 */
.dark .logo-image {
  filter: grayscale(100%) brightness(0) invert(1);
}

.dark .partner-logo:hover .logo-image {
  filter: grayscale(0%) brightness(1) invert(0);
}

/* Responsive */
@media (max-width: 1024px) {
  .main-title {
    font-size: 40px;
  }
  
  .partners-grid {
    gap: 40px;
  }
  
  .logo-image {
    height: 36px;
    width: 100px;
  }
}

@media (max-width: 768px) {
  .container {
    padding: 0 16px;
  }
  
  .main-title {
    font-size: 36px;
    margin-bottom: 40px;
  }
  
  .partners-grid {
    gap: 30px;
  }
  
  .logo-image {
    height: 32px;
    width: 90px;
  }
}

@media (max-width: 480px) {
  .main-title {
    font-size: 28px;
  }
  
  .partners-grid {
    gap: 24px;
  }
  
  .logo-image {
    height: 28px;
    width: 80px;
  }
}
</style>
