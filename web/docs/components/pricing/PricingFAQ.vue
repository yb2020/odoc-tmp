<template>
  <div class="pricing-faq-container">
    <h2 class="faq-title">{{ t('Pricing.faq.title') }}</h2>
    
    <div class="faq-list">
      <div 
        v-for="(item, index) in faqItems" 
        :key="index" 
        class="faq-item"
        :class="{ 'active': activeIndex === index }"
      >
        <div class="faq-question" @click="toggleFaq(index)">
          <span>{{ item.question }}</span>
          <span class="toggle-icon">{{ activeIndex === index ? '−' : '+' }}</span>
        </div>
        <div class="faq-answer" v-show="activeIndex === index">
          <p>{{ item.answer }}</p>
        </div>
      </div>
    </div>
    
    <div class="faq-support">
      <p>{{ t('Pricing.faq.stillHaveQuestions') }} <a :href="supportLink" class="support-link">{{ t('Pricing.faq.support') }}</a> {{ t('Pricing.faq.team') }}.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useVitePressI18n } from '../../i18n'

// 使用项目的国际化系统
const { t } = useVitePressI18n()

// 控制当前打开的FAQ项
const activeIndex = ref(null)

// 切换FAQ项的打开/关闭状态
const toggleFaq = (index) => {
  activeIndex.value = activeIndex.value === index ? null : index
}

// FAQ项目列表
const faqItems = computed(() => [
  {
    question: t('Pricing.faq.questions.promptCredit.question'),
    answer: t('Pricing.faq.questions.promptCredit.answer')
  },
  {
    question: t('Pricing.faq.questions.runOutCredits.question'),
    answer: t('Pricing.faq.questions.runOutCredits.answer')
  },
  {
    question: t('Pricing.faq.questions.studentPricing.question'),
    answer: t('Pricing.faq.questions.studentPricing.answer')
  },
  {
    question: t('Pricing.faq.questions.earlyAdopter.question'),
    answer: t('Pricing.faq.questions.earlyAdopter.answer')
  },
  {
    question: t('Pricing.faq.questions.autoRefills.question'),
    answer: t('Pricing.faq.questions.autoRefills.answer')
  },
  {
    question: t('Pricing.faq.questions.subscription.question'),
    answer: t('Pricing.faq.questions.subscription.answer')
  },
  {
    question: t('Pricing.faq.questions.subscriptionEnd.question'),
    answer: t('Pricing.faq.questions.subscriptionEnd.answer')
  },
  {
    question: t('Pricing.faq.questions.permanentPurchase.question'),
    answer: t('Pricing.faq.questions.permanentPurchase.answer')
  }
])

// 支持链接
const supportLink = '#'
</script>

<style scoped>
.pricing-faq-container {
  max-width: 1000px;
  margin: 80px auto 0;
  padding: 0 20px;
}

.faq-title {
  font-size: 36px;
  font-weight: 700;
  text-align: center;
  margin-bottom: 40px;
}

.faq-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.faq-item {
  border-bottom: 1px solid var(--vp-c-divider);
}

.faq-question {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
  font-size: 18px;
  font-weight: 500;
  cursor: pointer;
  transition: color 0.3s;
}

.faq-question:hover {
  color: var(--vp-c-brand);
}

.toggle-icon {
  font-size: 24px;
  font-weight: 300;
}

.faq-answer {
  padding: 0 0 20px;
  color: var(--vp-c-text-2);
  line-height: 1.6;
}

.faq-support {
  margin-top: 60px;
  text-align: center;
  font-size: 18px;
}

.support-link {
  color: var(--vp-c-brand);
  text-decoration: none;
  font-weight: 500;
  transition: opacity 0.2s;
}

.support-link:hover {
  opacity: 0.8;
}

@media (max-width: 768px) {
  .faq-title {
    font-size: 28px;
    margin-bottom: 30px;
  }
  
  .faq-question {
    font-size: 16px;
    padding: 16px 0;
  }
  
  .faq-support {
    margin-top: 40px;
    font-size: 16px;
  }
}
</style>
