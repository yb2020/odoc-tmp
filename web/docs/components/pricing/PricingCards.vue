<template>
  <div class="pricing-container">
    <div class="gradient-title">{{ t('Pricing.title') }}</div>
    <!-- <div class="beta-banner">
      <span class="beta-text">{{ t('Pricing.betaTitle') }}</span>
    </div> -->
    <p class="subtitle">{{ t('Pricing.subtitle') }}</p>

    <div class="pricing-promo">
      <a href="#pricing-plans" class="promo-button">{{ t('Pricing.promo') }} →</a>
    </div>
    <div class="beta-banner">
      <span>{{ t('Pricing.betaTitle') }}</span>
    </div>

    <div class="pricing-plans" id="pricing-plans">
      <div class="pricing-plan">
        <h2>{{ t('Pricing.free.title') }}</h2>
        <div class="price">
          <span class="amount">{{ t('Pricing.free.price', [currencyToSymbol(freePlan?.base?.subInfo?.currency), toNumberString(freePlan?.base?.subInfo?.price||0)]) }}</span>
          <span class="period">{{ t('Pricing.free.period') }}</span>
        </div>
        <a href="/workbench/recent" class="plan-button">{{ t('Pricing.free.button') }}</a>
        <ul class="plan-features">
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.free.features.trial') }}</li>
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.free.features.copilot') }}</li>
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.free.features.aiTranslation') }}</li>
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.free.features.docmentCapacity') }}</li>
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.free.features.fullTextTranslation') }}</li>
          <li class="featureDelete"><span class="check">✓</span> {{ t('Pricing.free.features.credits', [toNumberString(freePlan?.base?.subInfo?.originalCredit||2500)]) }}</li>
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.free.features.betaCredits', [toNumberString(freePlan?.base?.subInfo?.credit||2500)]) }}</li>
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.free.features.addonCredits') }}</li>
        </ul>
      </div>
      
      <div class="pricing-plan featured">
        <h2>{{ t('Pricing.pro.title') }}</h2>
        <div class="price-group"> <!-- 新增的容器 -->
          <div class="pro-price">
            <span class="amount">{{ t('Pricing.pro.price', [currencyToSymbol(paidPlan?.base?.subInfo?.currency), toNumberString(paidPlan?.base?.subInfo?.originalPrice||1000)]) }}</span>
            <span class="period">{{ t('Pricing.pro.period') }}</span>
          </div>
          <div class="beta-price">
            <span class="beta" style="margin: 0 4px; font-weight: bold; font-size: 20px;">{{ t('Pricing.pro.betaPrice') }}</span>
            <span class="amount">{{ t('Pricing.pro.price', [currencyToSymbol(paidPlan?.base?.subInfo?.currency), toNumberString(paidPlan?.base?.subInfo?.price||1000)]) }}</span>
            <span class="period">{{ t('Pricing.pro.period') }}</span>
          </div>
        </div>
        <PaidPlanButton class="plan-button primary" :buttonText="t('Pricing.pro.button')"/>
        <ul class="plan-features">
          <li class="feature highlight">{{ t('Pricing.pro.features.everything') }}</li>
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.pro.features.copilot') }}</li>
          <li class="featureDelete"><span class="check">✓</span> {{ t('Pricing.pro.features.credits', [toNumberString(paidPlan?.base?.subInfo?.originalCredit||50000)]) }}</li>
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.pro.features.betaCredits', [toNumberString(paidPlan?.base?.subInfo?.credit||50000)]) }}</li>
          <li class="feature"><span class="check">✓</span> {{ t('Pricing.pro.features.addonCredits', [currencyToSymbol(paidPlan?.base?.subAddOnCreditInfo?.currency), toNumberString(paidPlan?.base?.subAddOnCreditInfo?.price||1000),toNumberString(paidPlan?.base?.subAddOnCreditInfo?.addOnCredit||25000)]) }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useVitePressI18n } from '../../i18n'
import { onMounted } from 'vue'
import { useMembershipStore } from '@/store/membership'
import { toNumberString, currencyToSymbol } from '@/store/membership/default'
import PaidPlanButton from '@/common/src/components/Pay/PaidPlanButton.vue'
import { storeToRefs } from 'pinia'

// 使用项目的国际化系统
const { t } = useVitePressI18n()

// 使用会员计划store
const membershipStore = useMembershipStore()
const { freePlan, paidPlan } = storeToRefs(membershipStore)

// 在组件挂载时获取数据
onMounted(async () => {
  if (!membershipStore.isDataLoaded) await membershipStore.fetchSubPlanInfos()
})
</script>

<style scoped>
.pricing-container {
  max-width: 1200px;
  margin: 0 auto;
  margin-top: 80px;
  padding: 0 20px;
}

.gradient-title {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 16px;
  background-image: linear-gradient(90deg, var(--vp-c-brand), var(--vp-c-brand-light));
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  color: var(--vp-c-brand); /* Fallback */
  text-align: center;
}

/* Beta discount banner */
.beta-banner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin: 12px auto 24px;
  padding: 10px 14px;
  border-radius: 10px;
  width: fit-content;
  max-width: 100%;
  text-align: center;
  color: var(--vp-c-red, #e11d48);
}

.subtitle {
  font-size: 20px;
  color: var(--vp-c-text-2);
  text-align: center;
  margin-bottom: 40px;
}

.pricing-promo {
  display: flex;
  justify-content: center;
  margin: 30px 0;
}

.promo-button {
  display: inline-flex;
  align-items: center;
  background-color: var(--vp-c-brand);
  color: var(--vp-c-text-fg-ffffff);
  padding: 12px 24px;
  border-radius: 24px;
  text-decoration: none;
  font-size: 16px;
  transition: all 0.3s ease;
}

.promo-button:hover {
  background-color: var(--vp-c-brand-light);
  transform: translateY(-2px);
}

.pricing-plans {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
  margin-bottom: 60px;
}

.pricing-plan {
  background-color: var(--vp-c-bg-soft);
  border-radius: 12px;
  padding: 32px 24px;
  border: 1px solid var(--vp-c-divider);
  transition: all 0.3s ease;
}

.pricing-plan:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.pricing-plan.featured {
  border-color: var(--vp-c-brand);
}

.pricing-plan h2 {
  font-size: 24px;
  font-weight: 600;
  margin-top: 0;
  margin-bottom: 16px;
}

.price {
  margin-bottom: 24px;
}

.pro-price .amount {
  font-size: 36px;
  line-height: 1;
  text-decoration: line-through;
  color: var(--vp-c-text-3);
}
.pro-price .period {
  font-size: 14px;
  color: var(--vp-c-text-3);
  margin-left: 4px;
}

.beta-price .amount {
  font-size: 36px;
  font-weight: 700;
  line-height: 1;
}
.beta-price .period {
  font-size: 14px;
  color: var(--vp-c-text-2);
  margin-left: 4px;
}
  

.pricing-plan.featured .price-group {
  display: flex;
  justify-content: flex-start; /* 将第一个价格靠左对齐 */
  align-items: center; /* 垂直居中 */
  margin-bottom: 24px; /* 为整个价格组设置底部间距 */
  gap: 30px; /* 两个价格块之间的间距，根据需要调整，现在设置为15px，你可以尝试更小的值如 10px */
}


.price .amount {
  font-size: 36px;
  font-weight: 700;
  line-height: 1;
}

.price .period {
  font-size: 14px;
  color: var(--vp-c-text-2);
  margin-left: 4px;
}

.plan-button {
  display: block;
  width: 100%;
  padding: 12px;
  border-radius: 6px;
  font-weight: 500;
  text-align: center;
  text-decoration: none;
  transition: all 0.3s ease;
  background-color: var(--vp-c-bg-mute);
  color: var(--vp-c-text-1);
  border: 1px solid var(--vp-c-divider);
  margin-bottom: 24px;
}

.plan-button.primary {
  background-color: var(--vp-c-brand);
  color: white;
  border: none;
}

.plan-button:hover {
  transform: translateY(-2px);
}

.plan-button.primary:hover {
  background-color: var(--vp-c-brand-dark);
}

.plan-features {
  list-style: none;
  padding: 0;
  margin: 0;
}

.feature {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;
  line-height: 1.4;
}

.feature.highlight {
  font-weight: 500;
  margin-bottom: 20px;
  color: var(--vp-c-brand);
}

.feature .check {
  color: var(--vp-c-green);
  margin-right: 10px;
  font-weight: bold;
}
/** 删除的样式 */
.featureDelete {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;
  line-height: 1.4;
  color: var(--vp-c-text-3);
  text-decoration: line-through; /* 添加删除线 */
}

.featureDelete.highlight {
  font-weight: 500;
  margin-bottom: 20px;
}

.featureDelete .check {
  color: var(--vp-c-divider);
  margin-right: 10px;
  font-weight: bold;
}

@media (max-width: 768px) {
  .pricing-plans {
    grid-template-columns: 1fr;
  }
  
  .gradient-title {
    font-size: 36px;
  }
  
  .subtitle {
    font-size: 18px;
  }
  
  /* Responsive sizes for beta banner */
  .beta-banner {
    gap: 8px;
    padding: 8px 12px;
    margin: 10px auto 20px;
  }
  .beta-badge {
    font-size: 8px;
    padding: 2px 7px;
  }
  .beta-text {
    font-size: 10px;
  }
}
</style>
