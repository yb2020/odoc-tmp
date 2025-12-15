import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMembershipInfo, getMembershipBaseInfo, getMembershipSubPlanInfos } from '@/api/membership'
import { defaultSubPlanInfos } from './default'
import { SubPlanInfo, MembershipInfo, BaseInfo } from './types'
import { MembershipType } from 'go-sea-proto/gen/ts/membership/MembershipApi';


export const useMembershipStore = defineStore('membership', () => {
  // 状态
  const freePlan = ref<SubPlanInfo>(defaultSubPlanInfos.find(plan => plan.type === MembershipType.FREE) || defaultSubPlanInfos[0])
  const paidPlan = ref<SubPlanInfo>(defaultSubPlanInfos.find(plan => plan.type === MembershipType.PRO) || defaultSubPlanInfos[1])
  const isLoading = ref(false)
  const error = ref<any>(null)
  const isDataLoaded = ref(false)

  // 获取订阅计划信息
  const fetchSubPlanInfos = async () => {
    const result = await getMembershipSubPlanInfos()
    if (!result || !result.subPlanInfos) {
        freePlan.value = defaultSubPlanInfos.find(plan => plan.type === MembershipType.FREE) || defaultSubPlanInfos[0]
        paidPlan.value = defaultSubPlanInfos.find(plan => plan.type === MembershipType.PRO) || defaultSubPlanInfos[1]
      return
    }

    freePlan.value = result.subPlanInfos.find(plan => plan.type === MembershipType.FREE)
    paidPlan.value = result.subPlanInfos.find(plan => plan.type === MembershipType.PRO)
    isDataLoaded.value = true
  }


  return {
    // 状态
    freePlan,
    paidPlan,
    isLoading,
    error,
    isDataLoaded,
    
    // 方法
    fetchSubPlanInfos,
  }
})
