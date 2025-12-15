<template>
  <div v-if="differData" class="differ-with-chatgpt">
    <div class="container">
      <div class="content-layout">
        <div class="left-column">
          <h2 class="main-title">{{ differData.title }}</h2>
        </div>
        
        <div class="right-column">
          <div class="sections-list">
            <div
              v-for="(section, key) in differData.sections"
              :key="key"
              class="section-item"
              :class="{ 'is-expanded': expandedItems.includes(String(key)) }"
            >
              <button
                class="section-header"
                @click="toggleSection(String(key))"
              >
                <h3 class="section-title">{{ section.title }}</h3>
                <span class="expand-icon">+</span>
              </button>
              
              <div class="section-content-wrapper">
                <p class="section-content">{{ section.content }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import docsI18nInstance, { useVitePressI18n } from '../../i18n/index';

const { t, locale } = useVitePressI18n();

const expandedItems = ref<string[]>([]);

const messages = computed(() => docsI18nInstance.global.messages.value[locale.value] || {});
const differData = computed(() => (messages.value.DifferWithChatGPT as any) || null);

const toggleSection = (key: string) => {
  const index = expandedItems.value.indexOf(key);
  if (index > -1) {
    expandedItems.value.splice(index, 1);
  } else {
    expandedItems.value.push(key);
  }
};
</script>

<style scoped>
.differ-with-chatgpt {
  position: relative;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.content-layout {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 80px;
  align-items: start;
}

.left-column {
  position: sticky;
  top: 100px;
}

.main-title {
  font-size: 48px;
  font-weight: 700;
  line-height: 1.2;
  color: var(--vp-c-text-1);
  margin: 0;
  border: none;
}
.sections-list {
  display: flex;
  flex-direction: column;
}

.section-item {
  border-bottom: 1px solid var(--vp-c-divider);
  transition: all 0.3s ease;
}

.section-item:last-child {
  border-bottom: none;
}

.section-header {
  width: 100%;
  background: none;
  border: none;
  padding: 30px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
}


.section-title {
  font-size: 18px;
  font-weight: 400;
  color: var(--vp-c-text-1);
  margin: 0;
  line-height: 1.3;
}

.expand-icon {
  font-size: 24px;
  font-weight: 300;
  color: var(--vp-c-brand-1);
  transition: transform 0.3s ease;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.section-item.is-expanded .expand-icon {
  transform: rotate(45deg);
}

.section-content-wrapper {
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease;
}

.section-item.is-expanded .section-content-wrapper {
  max-height: 300px;
}

.section-content {
  font-size: 16px;
  line-height: 1.7;
  color: var(--vp-c-text-2);
  margin: 0;
  padding: 0 0 40px 0;
}

/* Responsive */
@media (max-width: 1024px) {
  .content-layout {
    gap: 60px;
  }
  
  .main-title {
    font-size: 40px;
  }
  
  .section-title {
    font-size: 22px;
  }
}

@media (max-width: 768px) {
  .container {
    padding: 0 16px;
  }
  
  .content-layout {
    grid-template-columns: 1fr;
    gap: 40px;
  }
  
  .left-column {
    position: static;
    text-align: center;
  }
  
  .main-title {
    font-size: 36px;
  }
  
  .right-column {
    min-height: auto;
  }
  
  .section-item {
    padding: 32px 0;
  }
  
  .section-title {
    font-size: 20px;
  }
  
  .section-content {
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .main-title {
    font-size: 28px;
  }
  
  .section-item {
    padding: 24px 0;
  }
  
  .section-title {
    font-size: 18px;
  }
  
  .section-content {
    font-size: 14px;
  }
}
</style>
