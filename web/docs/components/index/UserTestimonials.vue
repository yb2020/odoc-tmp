<template>
  <div v-if="testimonialsData" class="user-testimonials">
    <div class="testimonials-background">
      <div class="container">
        <h2 class="main-title">{{ testimonialsData.title }}</h2>
        
        <div class="carousel-container">
          <div 
            class="carousel-track"
            :style="{ transform: `translateX(-${currentIndex * cardWidth}px)` }"
          >
            <div
              v-for="(testimonial, key) in testimonialsData.testimonials"
              :key="key"
              class="testimonial-card"
            >
              <div class="card-content">
                <p class="testimonial-text">{{ testimonial.content }}</p>
                <div class="author-info">
                  <div class="avatar-wrapper">
                    <img 
                      :src="getImageUrl(`/images/icons/${testimonial.avatar}`)" 
                      :alt="testimonial.author"
                      class="avatar"
                    />
                  </div>
                  <div class="author-details">
                    <div class="author-name">{{ testimonial.author }}</div>
                    <div class="author-source">{{ testimonial.source }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="carousel-indicators">
          <button
            v-for="(_, index) in Object.keys(testimonialsData.testimonials)"
            :key="index"
            class="indicator"
            :class="{ active: currentIndex === index }"
            @click="goToSlide(index)"
          ></button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import docsI18nInstance, { useVitePressI18n } from '../../i18n/index';
import { getImageUrl } from '../../utils/contextPath';

const { t, locale } = useVitePressI18n();

const currentIndex = ref(0);
const cardWidth = ref(400);
const autoPlayInterval = ref<NodeJS.Timeout | null>(null);

const messages = computed(() => docsI18nInstance.global.messages.value[locale.value] || {});
const testimonialsData = computed(() => (messages.value.UserTestimonials as any) || null);

const totalCards = computed(() => {
  return testimonialsData.value ? Object.keys(testimonialsData.value.testimonials).length : 0;
});

const goToSlide = (index: number) => {
  currentIndex.value = index;
  resetAutoPlay();
};

const nextSlide = () => {
  currentIndex.value = (currentIndex.value + 1) % totalCards.value;
};

const startAutoPlay = () => {
  autoPlayInterval.value = setInterval(() => {
    nextSlide();
  }, 5000);
};

const resetAutoPlay = () => {
  if (autoPlayInterval.value) {
    clearInterval(autoPlayInterval.value);
  }
  startAutoPlay();
};

onMounted(() => {
  startAutoPlay();
});

onUnmounted(() => {
  if (autoPlayInterval.value) {
    clearInterval(autoPlayInterval.value);
  }
});
</script>

<style scoped>
.user-testimonials {
  position: relative;
  width: 100vw;
  margin-left: calc(-50vw + 50%);
  background-image: url('/images/user-testimonials-bg.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.testimonials-background {
  width: 100%;
  min-height: 600px;
  position: relative;
  padding: 80px 0;
  /* background: rgba(255, 255, 255, 0.9); */
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  position: relative;
  z-index: 2;
}

.main-title {
  font-size: 48px;
  font-weight: 700;
  line-height: 1.2;
  color: var(--vp-c-text-1);
  text-align: center;
  margin: 0 0 60px 0;
}

.carousel-container {
  position: relative;
  overflow: hidden;
  width: 100%;
}

.carousel-track {
  display: flex;
  transition: transform 0.5s ease-in-out;
  gap: 24px;
}

.testimonial-card {
  flex: 0 0 400px;
  background: white;
  border-radius: 16px;
  padding: 32px;
  /* box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1); */
  border: 1px solid var(--vp-c-divider-light);
}

.card-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.testimonial-text {
  font-size: 16px;
  line-height: 1.6;
  color: var(--vp-c-text-1);
  margin: 0 0 24px 0;
  white-space: pre-line;
  flex-grow: 1;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: auto;
}

.avatar-wrapper {
  flex-shrink: 0;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.author-details {
  flex-grow: 1;
}

.author-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--vp-c-text-1);
  margin-bottom: 4px;
}

.author-source {
  font-size: 12px;
  color: var(--vp-c-text-2);
}

.carousel-indicators {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 40px;
}

.indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: none;
  background: var(--vp-c-divider);
  cursor: pointer;
  transition: all 0.3s ease;
}

.indicator.active {
  background: var(--vp-c-brand-1);
  transform: scale(1.2);
}

.indicator:hover {
  background: var(--vp-c-brand-1);
  opacity: 0.7;
}

/* Responsive */
@media (max-width: 1024px) {
  .testimonial-card {
    flex: 0 0 350px;
  }
  
  .main-title {
    font-size: 40px;
  }
}

@media (max-width: 768px) {
  .container {
    padding: 0 16px;
  }
  
  .testimonial-card {
    flex: 0 0 300px;
    padding: 24px;
  }
  
  .main-title {
    font-size: 36px;
    margin-bottom: 40px;
  }
  
  .testimonial-text {
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .testimonial-card {
    flex: 0 0 280px;
    padding: 20px;
  }
  
  .main-title {
    font-size: 28px;
  }
  
  .testimonial-text {
    font-size: 14px;
  }
  
  .avatar {
    width: 40px;
    height: 40px;
  }
}
</style>
