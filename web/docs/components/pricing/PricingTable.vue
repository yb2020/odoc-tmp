<template>
  <div class="pricing-table-container">
    <h2 class="table-title">{{ t('Pricing.compareTitle') }}</h2>
    
    <div class="pricing-table-wrapper">
      <table class="pricing-table">
        <thead>
          <tr>
            <th></th>
            <th>
              {{ t('Pricing.free.title') }}
              <div class="price">{{ t('Pricing.free.price', [currencyToSymbol(freePlan?.base?.subInfo?.currency), toNumberString(freePlan?.base?.subInfo?.price)]) }}/{{ t('Pricing.month') }}</div>
            </th>
            <th class="highlighted-column">
              {{ t('Pricing.pro.title') }}
              <div class="price">{{ t('Pricing.pro.price', [currencyToSymbol(paidPlan?.base?.subInfo?.currency), toNumberString(paidPlan?.base?.subInfo?.price)]) }}/{{ t('Pricing.month') }}</div>
            </th>
            <th>
              <div>{{ t('Pricing.enterprise.title') }}</div>
              <div style="font-size: 12px;">{{ t('Pricing.enterprise.subtitle') }}</div>
              <div class="price">{{ t('Pricing.enterprise.price') }}</div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr class="section-header">
            <td>{{ t('Pricing.vibeReading') }}</td>
            <td></td> 
            <td></td>  
            <td></td>  
          </tr>
          <tr>
            <td>{{ t('Pricing.table.promptCredits') }}</td>
            <td> {{ t('Pricing.table.creditsPerMonth', [toNumberString(freePlan?.base?.subInfo?.credit)]) }}</td>
            <td> {{ t('Pricing.table.creditsPerMonth', [toNumberString(paidPlan?.base?.subInfo?.credit)]) }}</td>
            <td>{{ t('Pricing.enterprise.letsTalk') }}</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.addonCredits') }}</td>
            <td class="no-feature">✕</td>
            <td>{{ t('Pricing.table.addonCreditsPrice', [toNumberString(paidPlan?.base?.subAddOnCreditInfo?.price), currencyToSymbol(paidPlan?.base?.subAddOnCreditInfo?.currency), toNumberString(paidPlan?.base?.subAddOnCreditInfo?.addOnCredit)]) }}</td>
            <td>{{ t('Pricing.enterprise.letsTalk') }}</td>
          </tr>
          <tr class="section-header">
            <td>{{ t('Pricing.table.features') }}</td>
            <td></td> 
            <td></td>  
            <td></td>  
          </tr>
          <tr>
            <td>{{ t('Pricing.table.aiCopilot') }}</td>
            <td :class="freePlan?.ai?.copilot?.isEnable ? 'has-feature' : 'no-feature'">{{ freePlan?.ai?.copilot?.isEnable ? '✓' : '✕' }}</td>
            <td :class="paidPlan?.ai?.copilot?.isEnable ? 'has-feature' : 'no-feature'">{{ paidPlan?.ai?.copilot?.isEnable ? '✓' : '✕' }}</td>
            <td class="has-feature">✓</td>
          </tr>
          
          <!-- 动态模型列表 -->
          <template v-for="model in modelList" :key="model.name">
            <!-- 模型名称行 -->
            <tr>
              <td>{{ model.name }}</td>
              <td :class="getModelSupport(model, 'free') ? 'has-feature' : 'no-feature'">
                {{ getModelSupport(model, 'free') ? '✓' : '✕' }}
              </td>
              <td :class="getModelSupport(model, 'paid') ? 'has-feature' : 'no-feature'">
                {{ getModelSupport(model, 'paid') ? '✓' : '✕' }}
              </td>
              <td class="has-feature">✓</td>
            </tr>
            <!-- 模型费用行 -->
            <tr>
              <td>{{ t('Pricing.table.costCredit') }}</td>
              <td>{{ getModelCredit(model, 'free') }}</td>
              <td>{{ getModelCredit(model, 'paid') }}</td>
              <td class="has-feature">✓</td>
            </tr>
          </template>


          <tr>
            <td>{{ t('Pricing.table.docs') }} </td>
            <td class="has-feature">✓</td>
            <td class="has-feature">✓</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.docMaxCapacity') }} </td>
            <td>
              <p class="deleteTd"> {{ t('Pricing.table.docMaxStorageCapacity', [freePlan?.docs?.maxStorageCapacityOriginal]) }} </p>
              <p> <span>{{ t('Pricing.table.betaSign') }}</span> {{ t('Pricing.table.docMaxStorageCapacity', [freePlan?.docs?.maxStorageCapacity]) }}  </p>
            </td>
            <td>{{ t('Pricing.table.docMaxStorageCapacity', [paidPlan?.docs?.maxStorageCapacity]) }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.docMaxPageCount') }} </td>
            <td>
              <p class="deleteTd"> {{ t('Pricing.table.docUploadMaxPageCount', [freePlan?.docs?.docUploadMaxPageCountOriginal]) }} </p>
              <p> <span>{{ t('Pricing.table.betaSign') }}</span> {{ t('Pricing.table.docUploadMaxPageCount', [freePlan?.docs?.docUploadMaxPageCount]) }}  </p>
            </td>
            <td>{{ t('Pricing.table.docUploadMaxPageCount', [paidPlan?.docs?.docUploadMaxPageCount]) }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.docMaxSize') }} </td>
            <td>
              <p class="deleteTd"> {{ t('Pricing.table.docUploadMaxSize', [freePlan?.docs?.docUploadMaxSizeOriginal]) }} </p>
              <p> <span>{{ t('Pricing.table.betaSign') }}</span> {{ t('Pricing.table.docUploadMaxSize', [freePlan?.docs?.docUploadMaxSize]) }}  </p>
            </td>
            <td>{{ t('Pricing.table.docUploadMaxSize', [paidPlan?.docs?.docUploadMaxSize]) }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.notes') }} </td>
            <td :class="freePlan?.note?.isNoteManage ? 'has-feature' : 'no-feature'">{{ freePlan?.note?.isNoteManage ? '✓' : '✕' }}</td>
            <td :class="paidPlan?.note?.isNoteManage ? 'has-feature' : 'no-feature'">{{ paidPlan?.note?.isNoteManage ? '✓' : '✕' }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.fullTextTranslate') }} </td>
            <td :class="freePlan?.translate?.isFullTextTranslate ? 'has-feature' : 'no-feature'">{{ freePlan?.translate?.isFullTextTranslate ? '✓' : '✕' }}</td>
            <td :class="paidPlan?.translate?.isFullTextTranslate ? 'has-feature' : 'no-feature'">{{ paidPlan?.translate?.isFullTextTranslate ? '✓' : '✕' }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.fullTextTranslateCredits') }} </td>
            <td>{{ t('Pricing.table.fullTextTranslateCreditsValue', [toNumberString(freePlan?.translate?.fullTextTranslateCreditCost)]) }}</td>
            <td>{{ t('Pricing.table.fullTextTranslateCreditsValue', [toNumberString(paidPlan?.translate?.fullTextTranslateCreditCost)]) }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.fullTextTranslateMaxPageCount') }} </td>
            <td>{{ t('Pricing.table.fullTextTranslateMaxPageCountValue', [freePlan?.translate?.fullTextTranslateMaxPageCount]) }}</td>
            <td>{{ t('Pricing.table.fullTextTranslateMaxPageCountValue', [paidPlan?.translate?.fullTextTranslateMaxPageCount]) }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.wordTranslate') }} </td>
            <td :class="freePlan?.translate?.isWordTranslate ? 'has-feature' : 'no-feature'">{{ freePlan?.translate?.isWordTranslate ? '✓' : '✕' }}</td>
            <td :class="paidPlan?.translate?.isWordTranslate ? 'has-feature' : 'no-feature'">{{ paidPlan?.translate?.isWordTranslate ? '✓' : '✕' }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.wordTranslateCredits') }} </td>
            <td>{{ t('Pricing.table.wordTranslateCreditsValue', [toNumberString(freePlan?.translate?.wordTranslateCreditCost)]) }}</td>
            <td>{{ t('Pricing.table.wordTranslateCreditsValue', [toNumberString(paidPlan?.translate?.wordTranslateCreditCost)]) }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.aiTranslation') }}</td>
            <td :class="freePlan?.translate?.isAiTranslation ? 'has-feature' : 'no-feature'">{{ freePlan?.translate?.isAiTranslation ? '✓' : '✕' }}</td>
            <td :class="paidPlan?.translate?.isAiTranslation ? 'has-feature' : 'no-feature'">{{ paidPlan?.translate?.isAiTranslation ? '✓' : '✕' }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.aiTranslationCredits') }} </td>
            <td>{{ t('Pricing.table.aiTranslationCreditsValue', [toNumberString(freePlan?.translate?.aiTranslationCreditCost)]) }}</td>
            <td>{{ t('Pricing.table.aiTranslationCreditsValue', [toNumberString(paidPlan?.translate?.aiTranslationCreditCost)]) }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.ocrTranslate') }}</td>
            <td :class="freePlan?.translate?.isOcr ? 'has-feature' : 'no-feature'">{{ freePlan?.translate?.isOcr ? '✓' : '✕' }}</td>
            <td :class="paidPlan?.translate?.isOcr ? 'has-feature' : 'no-feature'">{{ paidPlan?.translate?.isOcr ? '✓' : '✕' }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr>
            <td>{{ t('Pricing.table.ocrTranslateCredits') }} </td>
            <td>{{ t('Pricing.table.ocrTranslateCreditsValue', [toNumberString(freePlan?.translate?.ocrCreditCost)]) }}</td>
            <td>{{ t('Pricing.table.ocrTranslateCreditsValue', [toNumberString(paidPlan?.translate?.ocrCreditCost)]) }}</td>
            <td class="has-feature">✓</td>
          </tr>
          <tr class="button-row">
            <td></td>
            <td>
              <a href="/workbench/recent" class="table-button">{{ t('Pricing.table.startReading') }}</a>
            </td>
            <td>
              <PaidPlanButton class="table-button primary" :buttonText="t('Pricing.pro.button')" />
            </td>
            <td>
              <a href="#" class="table-button enterprise">{{ t('Pricing.table.contactUs') }}</a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useVitePressI18n } from '../../i18n'
import { useMembershipStore } from '@/store/membership'
import { storeToRefs } from 'pinia'
import { toNumberString, currencyToSymbol } from '@/store/membership/default'
import PaidPlanButton from '@/common/src/components/Pay/PaidPlanButton.vue'
import { ref, computed } from 'vue'

// 使用项目的国际化系统
const { t } = useVitePressI18n()
// 使用会员计划store
const membershipStore = useMembershipStore()
const { freePlan, paidPlan } = storeToRefs(membershipStore)

// 从API获取模型列表数据
const modelList = computed(() => {
  // 合并免费和付费计划的模型，去重
  const freeModels = freePlan.value?.ai?.copilot?.models || []
  const paidModels = paidPlan.value?.ai?.copilot?.models || []
  
  // 使用Map去重，以key为唯一标识
  const modelMap = new Map()
  
  // 添加免费计划的模型
  freeModels.forEach(model => {
    if (model.isEnable) {
      modelMap.set(model.key, {
        key: model.key,
        name: model.name,
        isFree: model.isFree,
        freeSupported: true,
        freeCredit: model.creditCost,
        paidSupported: false,
        paidCredit: 0
      })
    }
  })
  
  // 添加付费计划的模型
  paidModels.forEach(model => {
    if (model.isEnable) {
      const existing = modelMap.get(model.key)
      if (existing) {
        // 如果已存在，更新付费信息
        existing.isFree = model.isFree,
        existing.paidSupported = true
        existing.paidCredit = model.creditCost
      } else {
        // 如果不存在，创建新的模型信息
        modelMap.set(model.key, {
          key: model.key,
          name: model.name,
          isFree: model.isFree,
          freeSupported: false,
          freeCredit: 0,
          paidSupported: true,
          paidCredit: model.creditCost
        })
      }
    }
  })
  
  return Array.from(modelMap.values())
})

// 获取模型支持状态
const getModelSupport = (model, planType) => {
  if (planType === 'free') {
    return model.freeSupported && freePlan.value?.ai?.copilot?.isEnable
  } else if (planType === 'paid') {
    return model.paidSupported && paidPlan.value?.ai?.copilot?.isEnable
  }
  return false
}

// 获取模型积分消耗
const getModelCredit = (model, planType) => {
  if (planType === 'free') {
    return model.freeSupported && !model.isFree ? `${toNumberString(model.freeCredit)}credit` : 'free'
  } else if (planType === 'paid') {
    return model.paidSupported && !model.isFree ? `${toNumberString(model.paidCredit)}credit` : 'free'
  }
  return '-'
}

</script>

<style scoped>
.pricing-table-container {
  max-width: 1200px;
  margin: 60px auto 0;
  padding: 0 20px;
}

.table-title {
  font-size: 36px;
  font-weight: 700;
  text-align: center;
  margin-bottom: 40px;
}

.pricing-table-wrapper {
  overflow-x: auto;
}

.pricing-table {
  width: 100%;
  border-collapse: collapse;
  border-spacing: 0;
  font-size: 16px;
}

.pricing-table th,
.pricing-table td {
  padding: 16px;
  text-align: center;
  border-bottom: 1px solid var(--vp-c-divider);
}

.pricing-table th {
  font-weight: 600;
  background-color: var(--vp-c-bg-soft);
  position: sticky;
  top: 0;
  z-index: 1;
}

.pricing-table th:first-child,
.pricing-table td:first-child {
  text-align: left;
  font-weight: 500;
}

.pricing-table th:first-child {
  width: 200px;
}

.price {
  font-size: 14px;
  color: var(--vp-c-text-2);
  margin-top: 4px;
}



.section-header td {
  font-weight: 600;
  background-color: var(--vp-c-bg-soft);
  text-align: left;
  padding-top: 24px;
}

.highlighted-column {
  background-color: var(--vp-c-brand-soft);
  border: 2px solid var(--vp-c-brand);
  border-bottom: 2px solid var(--vp-c-brand);
  position: relative;
}

.pricing-table td:nth-child(3) {
  background-color: var(--vp-c-bg-soft-mute);
  border-left: 2px solid var(--vp-c-brand);
  border-right: 2px solid var(--vp-c-brand);
}

.pricing-table tr:last-child td:nth-child(3) {
  border-bottom: 2px solid var(--vp-c-brand);
}

.section-header td:nth-child(3) {
  background-color: var(--vp-c-bg-soft);
  border-left: 2px solid var(--vp-c-brand);
  border-right: 2px solid var(--vp-c-brand);
}

.has-feature {
  color: #4caf50;
  font-weight: bold;
}

.no-feature {
  color: #f44336;
}

.button-row td {
  padding-top: 24px;
  padding-bottom: 24px;
  border-bottom: none;
}

.table-button {
  display: inline-block;
  width: 100%;
  max-width: 180px;
  padding: 10px 16px;
  border-radius: 6px;
  text-align: center;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
  background-color: var(--vp-c-bg-mute);
  color: var(--vp-c-text-1);
  border: 1px solid var(--vp-c-divider);
}

.table-button.primary {
  background-color: #4169e1;
  color: white;
  border: none;
}

.table-button.enterprise {
  background-color: #2c3e50;
  color: white;
  border: none;
}

.table-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

@media (max-width: 768px) {
  .table-title {
    font-size: 28px;
  }
  
  .pricing-table th,
  .pricing-table td {
    padding: 12px 8px;
    font-size: 14px;
  }
  
  .pricing-table th:first-child {
    width: 120px;
  }
  
  .table-button {
    padding: 8px 12px;
    font-size: 14px;
  }
}
.deleteTd {
  color: var(--vp-c-text-3);
  text-decoration: line-through;
}
</style>
