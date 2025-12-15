<template>
  <div v-if="scenarios && activeContent" class="target-user-scenarios">
    <div class="container">
      <div class="header-section">
        <div class="title">{{ t('TargetUserScenarios.title') }}</div>
      </div>

      <div class="tabs-container">
        <button
          v-for="(tab, index) in tabsData"
          :key="tab.key"
          class="tab-button"
          :class="{ 'active-tab': activeIndex === index }"
          @click="handleTabClick(index)"
        >
          {{ tab.name }}
        </button>
      </div>
    </div>

    <div class="content-container-wrapper">
      <transition :name="transitionName" mode="out-in">
        <div :key="activeIndex" class="content-container">
          <div class="content-background" :style="{ backgroundImage: `url(${getImageUrl(activeContent.imageUrl)})` }">
            <div class="container">
              <div class="text-overlay">
                <div class="text-content">
                  <div class="quote-icon">”</div>
                  <h3 class="content-title">{{ activeContent.title }}</h3>
                  <ul class="description-list">
                    <li v-for="(item, i) in activeContent.description" :key="i" v-html="item"></li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import docsI18nInstance, { useVitePressI18n } from '../../i18n/index';
import { getImageUrl } from '../../utils/contextPath';

const { t, locale } = useVitePressI18n();

const activeIndex = ref(0);
const transitionName = ref('slide-left');

const messages = computed(() => docsI18nInstance.global.messages.value[locale.value] || {});
const scenarios = computed(() => (messages.value.TargetUserScenarios as any) || null);

const tabsData = computed(() => {
  const tabs = scenarios.value?.tabs || {};
  return Object.keys(tabs).map(key => ({
    key,
    name: tabs[key],
  }));
});

const contentData = computed(() => scenarios.value?.content || {});

const activeContent = computed(() => {
  if (!scenarios.value || !tabsData.value.length) return null;
  const activeKey = tabsData.value[activeIndex.value]?.key;
  return activeKey ? contentData.value[activeKey] : null;
});

const featureTagPositions = [
  { top: '15%', left: '45%' },
  { top: '30%', left: '60%' },
  { top: '45%', left: '35%' },
  { top: '60%', left: '55%' },
];

const getFeatureTagStyle = (index: number) => {
  return featureTagPositions[index % featureTagPositions.length];
};

const handleTabClick = (index: number) => {
  if (index > activeIndex.value) {
    transitionName.value = 'slide-left';
  } else if (index < activeIndex.value) {
    transitionName.value = 'slide-right';
  }
  activeIndex.value = index;
};
</script>

<style scoped>
.target-user-scenarios {
  position: relative;
  padding: 40px 0;
}

.target-user-scenarios::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  right: 50%;
  margin-left: -50vw;
  margin-right: -50vw;
  height: 100%;
  background-color: #f7f8fa;
  z-index: -1;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.header-section {
  text-align: center;
  margin-bottom: 40px;
}

.title {
  font-size: 2.5rem;
  font-weight: bold;
  color: var(--vp-c-text-1);
}

.tabs-container {
  display: flex;
  justify-content: center;
  margin-bottom: 40px;
  border-bottom: 1px solid var(--vp-c-divider);
  position: relative;
  z-index: 10;
}

.tab-button {
  padding: 10px 20px;
  font-size: 1rem;
  font-weight: 500;
  border: none;
  background-color: transparent;
  color: var(--vp-c-text-2);
  cursor: pointer;
  transition: all 0.3s ease;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
}

.tab-button.active-tab {
  color: var(--vp-c-brand-1);
  border-bottom-color: var(--vp-c-brand-1);
}

.content-container-wrapper {
  position: relative;
  min-height: 550px;
}

.content-container {
  width: 100%;
}

.content-background {
  min-height: 550px;
  background-size: cover;
  background-position: center;
  position: relative;
  width: 100%;
}

.content-background > .container {
  position: relative;
  height: 550px;
  display: flex;
  align-items: center;
}

.text-overlay {
  width: 65%;
  height: 100%;
  background: linear-gradient(to right, var(--vp-c-bg-soft) 70%, transparent);
  display: flex;
  align-items: center;
  padding-left: 40px;
  padding-right: 100px;
}

.text-content {
  position: relative;
}

.quote-icon {
  font-size: 6rem;
  font-weight: bold;
  color: var(--vp-c-brand-1);
  opacity: 0.15;
  position: absolute;
  top: -20px;
  left: -20px;
  line-height: 1;
}

.content-title {
  font-weight: 700;
  color: var(--vp-c-text-1);
  margin-bottom: 30px;
  line-height: 1.3;
  margin-top: 10px;
}

.description-list {
  list-style-type: disc;
  padding-left: 20px;
}

.description-list li {
  color: var(--vp-c-text-2);
  line-height: 1.6;
  margin-bottom: 12px;
}

.description-list li::marker {
  color: var(--vp-c-brand-1);
}

.feature-tags {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.feature-tag {
  position: absolute;
  background-color: rgba(255, 255, 255, 0.9);
  color: var(--vp-c-text-1);
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 500;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.feature-tag:hover {
  transform: translateY(-3px);
}

.info-card {
  position: absolute;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.6) 0%, rgba(255, 255, 255, 0.3) 100%);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.scenario-card {
  bottom: 100px;
  left: -30px;
  right: 30px;
}

.role-card {
  bottom: 20px;
  left: -30px;
  right: 30px;
}

.card-title {
  font-weight: bold;
  font-size: 1rem;
  color: var(--vp-c-text-1);
  margin-bottom: 8px;
}

.card-text {
  font-size: 0.9rem;
  color: var(--vp-c-text-2);
  line-height: 1.5;
}

/* Transitions */
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.5s cubic-bezier(0.55, 0, 0.1, 1);
}

.slide-left-enter-from {
  opacity: 0;
  transform: translateX(50px);
}

.slide-left-leave-to {
  opacity: 0;
  transform: translateX(-50px);
}

.slide-right-enter-from {
  opacity: 0;
  transform: translateX(-50px);
}

.slide-right-leave-to {
  opacity: 0;
  transform: translateX(50px);
}

@media (max-width: 768px) {
  .content-inner {
    flex-direction: column-reverse;
  }
  .title {
    font-size: 2rem;
  }
  .content-title {
    font-size: 1.5rem;
  }
}
</style>
