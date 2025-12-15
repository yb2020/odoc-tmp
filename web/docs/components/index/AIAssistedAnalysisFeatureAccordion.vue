<template>
  <div
    :id="moduleType"
    class="feature-container"
    :style="containerStyle"
  >
    <BaseAccordion
      :itemsCount="featureItems.length"
      :autoRotate="autoRotate"
      :rotationInterval="rotationInterval"
      @change="handleCarouselChange"
    >
      <template #content="{ currentIndex: activeIndex, handleItemClick, pauseAutoRotation, resumeAutoRotation }">
        <div class="header-section">
          <div class="title">
            {{ t('AIAssistedAnalysisFeature.title') }}
          </div>
          <div class="title">
            {{ t('AIAssistedAnalysisFeature.description') }}
          </div>
        </div>
        <div :class="[
          'new-introduce-item-wrap',
          { 'new-even-introduce-item-wrap': reverse }
        ]">
          <div class="text-wrap">
            <div v-for="(item, sub) in featureItems" :key="sub" class="desc-wrap">
              <button
                class="btn"
                :class="{ 'current-btn': activeIndex === sub }"
                @mouseenter="handleItemClick(sub); pauseAutoRotation()"
                @mouseleave="resumeAutoRotation()"
              >
                <img :src="item.iconUrl" class="feature-icon" alt="" />
                <div class="text-content">
                  <span class="feature-title">{{ item.title }}</span>
                  <div class="sub-desc">
                    {{ item.description }}
                  </div>
                </div>
              </button>
            </div>
            <div class="actions-section">
              <StartReadingButton :buttonText="t('AIAssistedAnalysisFeature.buttons.startReading')" />
            </div>
          </div>

          <div class="carousel-wrap" @mouseenter="pauseAutoRotation()" @mouseleave="resumeAutoRotation()">
            <a-carousel ref="carouselRef" :dots="false" :current="activeIndex">
              <div v-for="(item, num) in featureItems" :key="num">
                <img :src="item.imgUrl" alt="" class="introduce-img" />
              </div>
            </a-carousel>
          </div>
        </div>
      </template>
    </BaseAccordion>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useVitePressI18n } from '../../i18n/index';
import { getImageUrl } from '../../utils/contextPath';
import BaseAccordion from '../common/BaseAccordion.vue';
import StartReadingButton from '../common/StartReadingButton.vue';

// 使用 VitePress 特定的国际化系统
const { t, locale } = useVitePressI18n();

// 定义props，接收reverse属性和背景颜色属性
const props = defineProps({
  reverse: {
    type: Boolean,
    default: false
  },
  bgColor: {
    type: String,
    default: ''
  },
  bgGradient: {
    type: String,
    default: ''
  },
  autoRotate: {
    type: Boolean,
    default: true
  },
  rotationInterval: {
    type: Number,
    default: 5000 // 默认5秒切换一次
  }
});

// 计算容器样式，根据props设置背景颜色
const containerStyle = computed(() => {
  if (props.bgColor) {
    return { background: props.bgColor };
  } else if (props.bgGradient) {
    return { background: props.bgGradient };
  }
  return {}; // 使用默认样式
});

// 模块类型
const moduleType = ref('AIAssistedAnalysisFeature');

// 特性项目，从当前语言的翻译文件中获取
const featureItems = computed(() => {
  const features = ['summaries', 'savedQuestions', 'instantAnswers', 'followUpQuestions', 'compareModels'];
  return features.map(key => {
    const imgUrl = getImageUrl(t(`AIAssistedAnalysisFeature.features.${key}.imgUrl`));
    
    return {
      i18nKey: key,
      title: t(`AIAssistedAnalysisFeature.features.${key}.title`),
      description: t(`AIAssistedAnalysisFeature.features.${key}.description`),
      imgUrl,
      iconUrl: getImageUrl(t(`AIAssistedAnalysisFeature.features.${key}.iconUrl`))
    };
  });
});

const carouselRef = ref<any>(null);

// 处理轮播变化事件
const handleCarouselChange = (index: number) => {
  if (carouselRef.value && carouselRef.value.goTo) {
    carouselRef.value.goTo(index);
  }
};

</script>

<style scoped>
/* 使用渐变背景色并占满整行 */
.feature-container {
  width: 100vw;
  position: relative;
  left: 50%;
  right: 50%;
  margin-left: -50vw;
  margin-right: -50vw;
  margin-top: 30px;
  padding: 0;
  overflow: hidden;
}

.new-introduce-item-wrap {
  display: flex;
  max-width: 1200px;
  margin: 0 auto;
  padding: 60px 40px;
  align-items: flex-start;
  gap: 60px;
}

.new-even-introduce-item-wrap {
  flex-direction: row-reverse;
}

.text-wrap {
  flex: 1;
  max-width: 500px;
}

.new-even-introduce-item-wrap .text-wrap {
  padding-right: 0;
  padding-left: 0;
}

/* 头部区域居中样式 */
.header-section {
  text-align: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 40px 20px;
}

.title {
  font-size: 44px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--vp-c-text-1);
  font-family: var(--vp-font-family-base);
  line-height: 1.2;
}


.btn {
  display: flex;
  align-items: flex-start; /* Align items to the top */
  background: transparent;
  border: 1px solid transparent; /* Set transparent border to maintain layout */
  padding: 16px 14px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
  font-family: var(--vp-font-family-base);
  text-align: left;
  width: 100%;
  border-radius: 8px;
  min-height: 100px; /* Ensure uniform height */
}

.btn:hover,
.current-btn {
  background-image: url('/docs/images/light_selected_bg.png');
  background-size: 100% 100%;
  background-repeat: no-repeat;
}

.btn:hover .feature-title,
.current-btn .feature-title {
  color: #4A6AD7;
}

.btn:hover .sub-desc,
.current-btn .sub-desc {
  color: var(--vp-c-text-2);
}

html.dark .btn:hover,
html.dark .current-btn {
  background-image: url('/docs/images/dark_selected_bg.png');
}

html.dark .btn:hover .feature-title,
html.dark .current-btn .feature-title {
  color: #fff;
}

html.dark .btn:hover .sub-desc,
html.dark .current-btn .sub-desc {
  color: rgba(255, 255, 255, 0.8);
}

.feature-icon {
  width: 20px;
  height: 20px;
  position: relative;
  top: 2px;
  margin-right: 12px;
  flex-shrink: 0;
  /* Default icon color to match text */
  filter: invert(40%) sepia(0%) saturate(0%) hue-rotate(0deg) brightness(90%) contrast(85%);
}

.btn:hover .feature-icon,
.current-btn .feature-icon {
  /* Blue icon on hover/selection */
  filter: invert(44%) sepia(85%) saturate(1478%) hue-rotate(206deg) brightness(91%) contrast(93%);
}

/* 重写 Ant Design 按钮的默认样式 */
.btn.ant-btn {
  box-shadow: none;
  background: transparent;
  border: none;
  height: auto;
}

.btn.ant-btn:hover, 
.btn.ant-btn:focus, 
.btn.ant-btn:active {
  background: transparent;
  box-shadow: none;
  border: none;
  color: #5B7CE2;
}


.feature-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--vp-c-text-1);
}

.text-content {
  display: flex;
  flex-direction: column;
}

.sub-desc {
  margin-top: 4px;
  font-size: 14px;
  color: var(--vp-c-text-2);
  line-height: 1.5;
  font-family: var(--vp-font-family-base);
  transition: color 0.3s;
}

.carousel-wrap {
  flex: 1;
  overflow: hidden;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  background: white;
  max-width: 600px;
}

.introduce-img {
  width: 100%;
  height: auto;
  object-fit: contain;
  display: block;
}

.actions-section {
  display: flex;
  gap: 16px;
  margin-top: 10px;
}

/* 针对移动设备的响应式布局 */
@media (max-width: 768px) {
  .new-introduce-item-wrap {
    flex-direction: column;
    padding: 20px;
  }
  
  .new-even-introduce-item-wrap {
    flex-direction: column;
  }
  
  .text-wrap {
    padding-right: 0;
    margin-bottom: 20px;
  }
  
  .new-even-introduce-item-wrap .text-wrap {
    padding-left: 0;
  }
  
  .carousel-wrap {
    width: 100%;
  }
}
</style>
