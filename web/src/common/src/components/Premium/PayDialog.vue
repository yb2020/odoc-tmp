<template>
  <component :is="inline ? 'div' : Modal" class="vip-trigger-modal" v-model:visible="visible" :width="702"
    :header="null" :footer="null">
    <div class="vip-trigger-inner relative" :class="{
        switchable,
      }">
      <div class="flex items-center p-4 text-white">
        <a-avatar :size="40" :src="userInfo?.avatarUrl"></a-avatar>
        <div class="ml-2 flex-1">
          <p class="text-base leading-22px font-medium m-0">
            {{ userInfo?.nickName }}
          </p>
          <p v-if="!isGroupBuy" class="text-sm leading-18px m-0">
            <template v-if="activeRole">
              <span v-for="role in [activeRole]" :key="role.vipType">
                {{
                $t(
                `common.premium.versions.${
                PremiumVipPreferences[role.vipType].key
                }`
                )
                }}{{ $t('common.premium.wordings.left') }}
                <span class="text-rp-red-6">{{ role.leftDays }}
                  {{ i18n.t('common.premium.units.day', role.leftDays)
                  }}<a-tooltip v-if="role.seniorIdx >= 0" class="ml-1">
                    <template #title>
                      {{
                      $t('common.premium.wordings.lefttip', [
                      $t(
                      `common.premium.versions.${
                      PremiumVipPreferences[
                      userRoles[role.seniorIdx].vipType
                      ].key
                      }`
                      ),
                      ])
                      }}
                    </template>
                    <InfoCircleOutlined />
                  </a-tooltip></span>
              </span>
            </template>
            <template v-else>{{ $t('common.premium.wordings.none') }}{{ vipTxt }}</template>
          </p>
        </div>
        <p v-if="licensee" class="mr-12 mb-0" :style="{ flex: '2 1 0%' }">
          ReadPaper {{ $t('common.pay.authorize', [licensee]) }}
        </p>
      </div>
      <!-- 团购tabs -->
      <a-tabs v-if="isGroupBuy" :active-key="groupTableKey" :animated="false" @change="onSwitch">
        <!-- 添加团购 tab-pane -->
        <a-tab-pane key="group-buy">
          <template #tab>
            <div class="w-full h-14 text-center" :class="{
                'bg-tab--first': true
              }" :style="{
                padding: '6px 0',
                color: '#fff'
              }">
              <p class="m-0 text-base leading-26px font-medium flex items-center justify-center">
                {{ $t('common.premium.versions.groupBuy') }}
              </p>
              <p class="text-xs">
                {{ $t('common.premium.versionTips.groupBuy') }}
              </p>
            </div>
          </template>
          <div class="p-4 bg-white">
            <!-- 团购内容区域 -->
            <a-row type="flex">
              <!-- 右侧支付区域 -->
              <a-col :style="{
                width: '40%',
                display: 'flex',
                }" class="groupbuy-trigger-right relative flex items-center justify-center">
                <div class="pl-6 min-w-0 flex items-center">
                  <template v-if="!enabled">
                    <div class="flex-1 text-rp-neutral-8 text-base text-center font-medium">
                      <img class="inline w-20 mb-2" src="@common/../assets/images/beans/pay-disabled.svg"
                        alt="disabled" />
                      <p class="m-0">{{ $t('common.premium.disabled') }}</p>
                    </div>
                  </template>
                  <template v-else>
                    <a-spin :spinning="loading">
                      <div class="relative mr-6 pb-8">
                        <div class="vip-trigger-code border border-solid border-rp-neutral-3" :style="{
                            width: '200px',
                            height: '200px',
                          }">
                          <QrcodeVue v-if="orderId" :value="url" :size="180" />
                          <div v-else
                            class="w-full h-full bg-rp-neutral-1 p-2 flex flex-col justify-center text-center">
                            <span v-if="error">{{ error.message }}</span>
                          </div>
                        </div>
                        <p v-if="waiting"
                          class="absolute w-full bottom-0 left-0 m-0 text-base text-center text-rp-red-6">
                          {{ $t('common.pay.scanned') }}
                        </p>
                      </div>
                    </a-spin>
                  </template>
                </div>
              </a-col>
              <!-- 左侧说明区域 -->
              <a-col :style="{
                width: '60%',
                background: '#FFFFFF',
                }" class="groupbuy-trigger-left p-4 rounded-sm text-white">
                
                <!-- 商品列表信息 -->
                <div class="group-buy-list mt-4 mb-10 text-black">
                  <h3 class="text-[24px] text-black">
                    {{ $t('common.premium.wordings.groupBuyTitle') }}
                  </h3>
                  <div v-for="(item, index) in groupProductItem" :key="index" class="group-buy-item">
                    <div class="group-buy-item flex justify-between items-center mt-1 mb-1">
                      <!-- <div class="item-header">
                        {{ item.name }}
                      </div> -->
                      <div class="item-content flex items-center gap-3 text-[16px]">
                        <span class="item-price">{{ item.name }}</span>
                        <span class="item-price">¥{{ item.price }}</span>
                        <span class="item-count ml-4">x{{ item.count }}</span>
                      </div>
                    </div>

                  </div>
                </div>

                <!-- 价格信息 -->
                <div class="mt-4">
                  <p class="flex items-end mb-3 whitespace-nowrap text-ellipsis overflow-hidden">
                    <span class="mr-1 text-4xl text-rp-red-6 font-black">
                      &yen;
                      <LoadingOutlined v-if="!price" />
                      <template v-else>{{ price }}</template>
                    </span>
                    <!-- <span v-if="discount > 0" class="text-sm text-rp-neutral-6 ml-2">
                      {{ $t('common.premium.wordings.save', [formatPrice(discount)]) }}
                    </span> -->
                    <!-- <span class="text-sm text-rp-neutral-6 text-[#000] ml-2">
                      {{ $t('common.premium.units.peruseryear') }}
                    </span> -->
                    <span v-if="discount > 0" class="text-sm text-rp-neutral-6 ml-2">
                      {{ $t('common.premium.wordings.save2', [formatPrice(discount)]) }}
                    </span>
                  </p>
                </div>
                <!-- Support 和合同信息 -->
                <div class="mt-4 text-black">
                  <Support class="!justify-normal" />
                  <p v-if="contract" class="mt-4 mb-0 text-rp-neutral-8">
                    {{ $t('common.premium.contracttip') }}
                    <a :href="contract" target="_blank">
                      {{ $t('common.premium.contract') }}
                    </a>
                  </p>
                </div>

                <!-- 帮助提示 -->
                <Help>
                  <p class="absolute right-0 bottom-0 m-0 text-rp-blue-6 cursor-help">
                    {{ $t('common.pay.supportTip') }}
                  </p>
                </Help>
              </a-col>
            </a-row>

          </div>
        </a-tab-pane>
      </a-tabs>

      <!-- 原来的内容tabs -->
      <a-tabs v-else :active-key="vipType" :animated="false" @change="onSwitch">
        <!-- 原有的会员类型 tabs -->
        <a-tab-pane v-for="(item, i) in vipList" :key="item.vipType">
          <template #tab>
            <div class="w-full h-14 text-center" :class="{
                [`bg-tab--${
                  i === 0
                    ? 'first'
                    : i === vipList.length - 1
                      ? 'last'
                      : 'middle'
                }`]: item.vipType === vipType,
              }" :style="{
                padding: '6px 0',
                color:
                  item.vipType == vipType
                    ? PremiumVipPreferences[item.vipType].color
                    : '#fff',
              }">
              <p class="m-0 text-base leading-26px font-medium flex items-center justify-center">
                <component :is="PremiumVipPreferences[item.vipType].icon" :class="{
                    [PremiumVipPreferences[item.vipType].key]:
                      item.vipType == vipType,
                  }" style="margin-right: 6px" />{{
                $t(
                `common.premium.versions.${
                PremiumVipPreferences[item.vipType].key
                }`
                )
                }}
              </p>
              <p class="text-xs">
                {{
                $t(
                `common.premium.versionTips.${
                PremiumVipPreferences[item.vipType].key
                }`
                )
                }}
              </p>
            </div>
          </template>
          <div class="p-4 bg-white">
            <a-row type="flex">
              <a-col v-if="visiblePrivilege" :style="{
                  width: '35%',
                  background: `${PremiumVipPreferences[item.vipType].color}`,
                }" class="vip-trigger-left p-4 rounded-sm text-white">
                <h4 class="text-base text-white">
                  {{ $t('common.premium.wordings.enhanced') }}
                </h4>
                <hr class="h-px border-0" />
                <ul class="flex flex-col gap-2 mt-3 mb-1 text-xs" style="line-height: 1.125rem">
                  <template v-if="vipType === VipType.ENTERPRISE">
                    <li v-for="tip in $t('common.premium.customtips')" :key="`${tip}`">
                      {{ tip }}
                    </li>
                  </template>
                  <template v-else>
                    <li v-for="privilege in privileges.slice(0, 12)" :key="`${privilege.name}`"
                      class="flex items-center justify-between">
                      <span>{{ privilege.name }}</span>
                      <span v-if="typeof privilege.value === 'string'">{{
                        privilege.value
                        }}</span>
                      <span v-else>-</span>
                    </li>
                    <li>
                      <a target="_blank" href="/vip" class="underline">{{
                        $t('common.premium.wordings.viewdiff')
                        }}</a>
                    </li>
                  </template>
                </ul>
              </a-col>
              <a-col :style="{
                  width: visiblePrivilege ? '65%' : '100%',
                  display: 'flex',
                }" class="vip-trigger-right relative flex items-center justify-center">
                <div class="pl-6 min-w-0 flex items-center">
                  <template v-if="!enabled">
                    <div class="flex-1 text-rp-neutral-8 text-base text-center font-medium">
                      <img class="inline w-20 mb-2" src="@common/../assets/images/beans/pay-disabled.svg"
                        alt="disabled" />
                      <p class="m-0">{{ $t('common.premium.disabled') }}</p>
                    </div>
                  </template>
                  <template v-else-if="vipType === VipType.ENTERPRISE">
                    <div class="flex-1 flex flex-col items-center text-base text-center font-medium">
                      <p class="mb-6">
                        {{ $t('common.premium.wordings.contact') }}
                      </p>
                      <div class="border border-solid border-rp-neutral-3 flex items-center justify-center" :style="{
                          width: '200px',
                          height: '200px',
                        }">
                        <QrcodeVue :value="entcode" :size="180" />
                      </div>
                    </div>
                  </template>
                  <template v-else-if="!planable">
                    <div class="flex-1 text-base text-center font-medium">
                      <img class="inline w-20 mb-2" src="@common/../assets/images/beans/pay-disabled.svg"
                        alt="disabled" />
                      <p v-if="vipLeftDays > maxDays" class="m-0">
                        {{
                        $t('common.premium.wordings.unpayable2', [
                        vipTxt,
                        maxDays,
                        ])
                        }}
                      </p>
                      <p v-else class="m-0">
                        {{ $t('common.premium.wordings.unpayable', [roleTxt]) }}
                        <a href="javascript:;" class="ml-1 text-rp-blue-6 underline"
                          @click.prevent="onSwitch(roleType)">
                          {{ $t('common.premium.btns.renew') }} {{ roleTxt }}</a>
                      </p>
                    </div>
                  </template>
                  <template v-else>
                    <a-spin :spinning="loading">
                      <div class="relative mr-6 pb-8">
                        <!-- <Support v-if="vipUpgrading" class="justify-center" /> -->
                        <div class="vip-trigger-code border border-solid border-rp-neutral-3" :style="{
                            width: '200px',
                            height: '200px',
                          }">
                          <QrcodeVue v-if="orderId" :value="url" :size="180" />
                          <div v-else
                            class="w-full h-full bg-rp-neutral-1 p-2 flex flex-col justify-center text-center">
                            <span v-if="error">{{ error.message }}</span>
                          </div>
                        </div>
                        <p v-if="waiting"
                          class="absolute w-full bottom-0 left-0 m-0 text-base text-center text-rp-red-6">
                          {{ $t('common.pay.scanned') }}
                        </p>
                        <p class="absolute left-0 w-max mt-8 h-4 text-rp-neutral-6 text-xs" style="margin-left: -8%">
                          <span class="block" style="transform: scale(0.833333)"><template
                              v-for="x in privilegesOnlyPro">{{ x.name }}{{ !x.isLast ? '、' : '' }}</template>{{
                            $t('common.premium.wordings.buytip')
                            }}</span>
                        </p>
                      </div>
                    </a-spin>
                    <div class="pb-8 min-w-0 flex-1 flex flex-col justify-center">
                      <img v-if="vipData?.costLabelUrl" :src="vipData.costLabelUrl" alt="label" class="mb-4" :style="{
                          width: 'fit-content',
                        }" />
                      <p v-if="vipUpgrading" class="mb-1 text-xs text-rp-neutral-8">
                        {{ $t('common.premium.wordings.pflprice', [vipPrice])
                        }}{{ $t('common.symbol.comma')
                        }}{{
                        lowerRolesIn365Days
                        .map((x) => {
                        return $t(
                        `common.premium.wordings.patchprice${
                        x.deductDays < x.leftDays ? 'full' : '' }`, [ $t( `common.premium.versions.${
                          PremiumVipPreferences[x.vipType].key }` ), x.deductDays, ] ) })
                          .join($t('common.symbol.comma')) }}{{ $t('common.symbol.comma') }}{{
                          $t('common.premium.wordings.pchprice') }}{{ $t('common.symbol.colon') }} </p>
                          <p class="flex items-end mb-3 whitespace-nowrap text-ellipsis overflow-hidden">
                            <span class="mr-1 text-4xl text-rp-red-6 font-black">&yen;
                              <LoadingOutlined v-if="!price" /><template v-else>{{ price }}</template>
                            </span>{{ $t('common.premium.units.peryear') }}
                            <span v-if="discount > 0" class="text-sm text-rp-neutral-6 ml-2">{{
                              $t('common.premium.wordings.save', [
                              formatPrice(discount),
                              ])
                              }}</span>
                          </p>
                          <!-- <Support v-if="!vipUpgrading" /> -->
                          <Support class="!justify-normal" />
                          <p v-if="contract" class="mt-4 mb-0 text-rp-neutral-8">
                            {{ $t('common.premium.contracttip')
                            }}<a :href="contract" target="_blank">{{
                              $t('common.premium.contract')
                              }}</a>
                          </p>
                    </div>
                    <Help>
                      <p class="absolute right-0 bottom-0 m-0 text-rp-blue-6 cursor-help">
                        {{ $t('common.pay.supportTip') }}
                      </p>
                    </Help>
                  </template>
                </div>
              </a-col>
            </a-row>
          </div>
        </a-tab-pane>
      </a-tabs>
    </div>
  </component>
</template>

<script lang="ts" setup>
import _, { mapKeys, snakeCase } from 'lodash'
import { computed, watch, ref } from 'vue'
import QrcodeVue from 'qrcode.vue'
import { useRequest } from 'ahooks-vue'
import { message, Modal } from 'ant-design-vue'
import { LoadingOutlined, InfoCircleOutlined } from '@ant-design/icons-vue'
import {
  VipRole,
  VipType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/user/VipUserInterface'
import { GetScanQRCodeResp, VipPayType, GroupBuyItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo'
import {
  getPayOrigin,
  getQrCodeParams,
  GetGroupBuyScanQRCode,
  PayStatus,
  VipType2ProductType,
} from '@common/api/vipPay'
import { formatPrice, isWaiting } from '@common/utils/pay'
import {
  PremiumVipPreferences,
  PremiumPrivileges,
  DefaultVipList,
  getPrivilegeTxt,
  PrivilegeTypes2Txt,
  VipTypePayable,
  VipType2PayType,
  VipType2ElementName,
  getPrivilegeDesc,
} from '@common/components/Premium/types'
import Help from '@common/components/Help.vue'
import Support from '@common/components/Pay/Support.vue'
import { useOrderInfo } from '@common/hooks/useOrderInfo'
import usePayConfig from '@common/hooks/useVipPayConfig'
import { compareVipType } from '@common/components/Premium/types'
import { EventCode, PageType, reportEvent } from '@common/utils/report'
import { useVipStore } from '../../stores/vip'
import { useUserStore } from '../../stores/user'
import { useI18n } from 'vue-i18n'
//import { GroupProductItem } from '~common/src/components/Premium/gropbuy' // 引入 gropbuy

defineProps<{
  inline?: boolean
  switchable?: boolean
  visiblePrivilege?: boolean
}>()

const groupTableKey = "group-buy" // 添加团购参数


const emit = defineEmits(['paysucc', 'paycancel'])

const PrivilegesConfig = _.flatten(
  Object.values(PremiumPrivileges).map((x) =>
    x.privileges
      .map((y) => ({
        ...y,
        prefix: x.prefix,
      }))
      .filter((y) => !(y.key in PrivilegeTypes2Txt))
  )
)

const i18n = useI18n()
const context = { i18n }
const userStore = useUserStore()
const vipStore = useVipStore()
const visible = ref(false)
const props = computed(() => vipStore.payDialogProps || {})
const userInfo = computed(() => userStore.userInfo)
const userRoles = computed(() => vipStore.roles)
const seniorRoles = computed(() => vipStore.seniorRoles)
const activeRole = computed(() =>
  seniorRoles.value.find((x) => x.vipType === vipType.value)
)
const lowerRolesIn365Days = computed(() => {
  let days = 0

  return seniorRoles.value
    .filter(
      (x) =>
        x.vipType >= VipType.STANDARD &&
        compareVipType(x.vipType, vipType.value) === -1
    )
    .reduce(
      (arr, x) => {
        if (days < 365) {
          days += x.leftDays
          arr.push({
            ...x,
            deductDays: days > 365 ? 365 - (days - x.leftDays) : x.leftDays,
          })
        }

        return arr
      },
      [] as Array<
        VipRole & {
          leftDays: number
          deductDays: number
        }
      >
    )
})
const roleType = computed(
  () => (userRoles.value[0]?.vipType ?? VipType.STANDARD) as VipTypePayable
)
const roleTxt = computed(() =>
  i18n.t(`common.premium.versions.${PremiumVipPreferences[roleType.value].key}`)
)
const maxDays = computed(() => vipConfig.value?.maxDays ?? 360)
const allowPlanLowVipDays = computed(
  () => vipConfig.value?.allowPlanLowVipDays ?? 90
)
const planable = computed(
  () =>
    compareVipType(vipType.value, roleType.value) === 1 ||
    (compareVipType(vipType.value, roleType.value) === 0 &&
      vipLeftDays.value <= maxDays.value) ||
    (compareVipType(vipType.value, roleType.value) === -1 &&
      seniorRoles.value
        ?.filter((x) => compareVipType(x.vipType, vipType.value) === 1)
        .every((x) => x.leftDays <= allowPlanLowVipDays.value) &&
      (activeRole.value?.leftDays ?? 0) <= maxDays.value)
)

const { data: vipConfig } = usePayConfig(0)
const vipType = ref<VipTypePayable>(props.value.needVipType ?? VipType.STANDARD)

const vipList = computed(() => {
  return (vipConfig.value?.vipPayPrivilege ?? DefaultVipList).filter(
    (x) => x.vipType !== VipType.FREE
  )
})
const vipData = computed(() =>
  vipList.value.find((x) => x.vipType === vipType.value)
)
const vipPrice = computed(() => {
  return vipData.value ? formatPrice(+vipData.value.payTotalAmount) : ''
})
const vipLeftDays = computed(() => activeRole.value?.leftDays ?? 0)
const vipTxt = computed(() =>
  i18n.t(`common.premium.versions.${PremiumVipPreferences[vipType.value].key}`)
)
const vipUpgrading = computed(() =>
  lowerRolesIn365Days.value.some((x) => x.leftDays > 0)
)

const privileges = computed(() => {
  return PrivilegesConfig.map((x) => {
    const txt = getPrivilegeTxt(context, x, vipData.value)
    return {
      name: i18n.t(`${x.prefix}.${x.key}`),
      value:
        typeof txt === 'string'
          ? `${txt}${
              x.typeDesc
                ? `(${getPrivilegeDesc(context, x, vipData.value, true)})`
                : ''
            }`
          : txt,
    }
  })
})
const privilegesOnlyPro = computed(() => {
  const all = vipConfig.value?.vipPayPrivilege || DefaultVipList
  const stdPrivileges = all.find((x) => x.vipType === VipType.STANDARD)
  const proPrivileges = all.find((x) => x.vipType === VipType.PROFESSIONAL)
  return PrivilegesConfig.filter((x) => {
    const stdFlag = getPrivilegeTxt(context, x, stdPrivileges)
    const proFlag = getPrivilegeTxt(context, x, proPrivileges)

    return !stdFlag.toString() && proFlag.toString()
  }).map((x, i, arr) => {
    return {
      name: i18n.t(`${x.prefix}.${x.key}`),
      isLast: i === arr.length - 1,
    }
  })
})
const enabled = computed(() => vipConfig.value?.paySwitch ?? true)
const licensee = computed(() => vipConfig.value?.licensee ?? '')
const contract = computed(() => vipConfig.value?.contract)
const entcode = computed(() => vipConfig.value?.enterpriseQRCode ?? '')

const isGroupBuy = computed(() => vipStore.payDialogProps.isGroupBuy || false) // 添加团购参数
const groupProductItem = computed(() => vipStore.payDialogProps.groupProductItem || []) // 添加团购参数
const groupBuyItemList = computed(() => {
  const items = vipStore.payDialogProps.groupProductItem || []
  return items.map(item => ({
    buyVipType: item.vipType,
    vipAmount: item.count
  })) as GroupBuyItem[]
})


const {
  data: order,
  loading,
  error,
  run: onPlaceOrder,
} = useRequest(
  async () => {
    // const res =
    //   vipType.value === VipType.ENTERPRISE
    //     ? ({} as GetScanQRCodeResp)
    //     : await getQrCodeParams({
    //         vipPayType: VipType2PayType[vipType.value],
    //       })
    //debugger
    console.log("isGrpouBuy", isGroupBuy.value)
    if(isGroupBuy.value){
      const res = await GetGroupBuyScanQRCode({
            vipPayType: VipPayType.GROUP_PAY,
            groupBuyItem: groupBuyItemList.value
          })

      return {
        ...res,
        vipType: vipType.value,
      }
    }else{
      const res =
      vipType.value === VipType.ENTERPRISE
        ? ({} as GetScanQRCodeResp)
        : await getQrCodeParams({
            vipPayType: VipType2PayType[vipType.value],
          })
      return {
        ...res,
        vipType: vipType.value,
      }
    }
    
  },
  {
    manual: true,
  }
)

const orderId = computed(() => order.value?.preOrderId)
const { data: orderInfo } = useOrderInfo(orderId)
const price = computed(() => {
  const s = orderInfo.value?.payInfo?.vipPrivilege?.payTotalAmount
  const n = parseInt(s ?? '', 10)

  return Number.isNaN(n) ? s : formatPrice(n)
})
const discount = computed(() => {
  const s = orderInfo.value?.payInfo?.vipPrivilege?.originalPayTotalAmount
  const s2 = orderInfo.value?.payInfo?.vipPrivilege?.payTotalAmount
  const n = parseInt(s ?? '', 10)
  const p = parseInt(s2 ?? '', 10)

  return Number.isNaN(n) || Number.isNaN(p) ? 0 : n - p
})
const url = computed(() => {
  const params = new URLSearchParams({
    id: orderId.value || '',
    //type: VipType2ProductType[vipType.value],// 旧的赋值方式
    type: isGroupBuy ? "groupbuy" : VipType2ProductType[vipType.value], // 团购时赋值groupbuy，之前的方式保留
  })
  if (vipConfig.value?.env) {
    params.append('env', vipConfig.value?.env)
  }
  const origin = vipStore.payOrigin || getPayOrigin()

  return `${origin}/pay?${params.toString()}`
})
const waiting = computed(() => isWaiting(orderInfo.value?.status))

const onSwitch = (type: VipTypePayable) => {
  vipType.value = type
}

const onCancel = () => {
  emit('paycancel')
  order.value = undefined
  orderInfo.value = undefined
  props.value.onPayCancel?.()
  vipStore.hideVipPayDialog()
}

watch(
  () => vipStore.showPayDialog,
  (v) => {
    if (v) {
      visible.value = true
    }
  },
  {
    immediate: true,
  }
)

watch(
  visible,
  () => {
    if (visible.value) {
      if (!userRoles.value.length) {
        vipStore.fetchVipProfile()
      }
    } else {
      onCancel()
    }
  },
  {
    immediate: true,
  }
)

watch(
  visible,
  () => {
    if (visible.value) {
      const defaultVipType =
        vipStore.role.vipType !== VipType.FREE
          ? (vipStore.role.vipType as VipTypePayable)
          : VipType.STANDARD
      vipType.value = props.value.needVipType ?? defaultVipType
    }
  },
  {
    immediate: true,
  }
)

watch(
  [visible, vipType],
  () => {
    if (
      visible.value &&
      enabled.value &&
      !loading.value &&
      vipType.value !== VipType.ENTERPRISE &&
      order.value?.vipType !== vipType.value
    ) {
      order.value = undefined
      orderInfo.value = undefined
      onPlaceOrder()
    }
  },
  {
    flush: 'post',
  }
)

watch(orderId, () => {
  const {
    payDialogProps: { reportParams },
  } = vipStore
  reportEvent(EventCode.readpaperVipPayPopupImpression, {
    page_type: PageType.premium,
    element_name: VipType2ElementName[vipType.value],
    pay_code_id: orderId.value,
    ...(!vipStore.payByDialog
      ? {}
      : mapKeys(reportParams, (_, k) => snakeCase(k))),
  })
})

watch(
  () => orderInfo.value?.status,
  (payStatus) => {
    if (payStatus === PayStatus.PAY_SUCCESS) {
      message.success(i18n.t('common.premium.succtip') as string)
      emit('paysucc', payStatus)
      props.value.onPaySucc?.(payStatus)
      vipStore.hideVipPayDialog(payStatus)
    } else if (payStatus === PayStatus.PAY_TIMEOUT) {
      onPlaceOrder()
    }
  }
)
</script>

<style lang="less">
.vip-trigger-modal {
  border-radius: 4px;

  .vip-trigger-inner {
    background-color: var(--site-theme-bg-dark);
    background-image: url('~common/assets/images/premium/bg-dialog.png');
    background-repeat: no-repeat;
    background-size: 100%;

    .ant-tabs-nav {
      display: none;
    }
  }

  .vip-trigger-inner.switchable {
    .ant-tabs-nav {
      display: block;
    }
  }

  .leading-18px {
    line-height: 18px;
  }

  .leading-22px {
    line-height: 22px;
  }

  .leading-26px {
    line-height: 26px;
  }

  .ant-modal-wrap {
    z-index: 1010;
  }

  .ant-modal-content {
    border-radius: 4px;
  }

  .ant-modal-body {
    padding: 0;
  }

  .ant-modal-close-x {
    color: var(--site-theme-text-inverse);
    font-size: 16px;
    line-height: 48px;
    width: 64px;

    &:hover {
      color: var(--site-theme-primary);
    }
  }

  .ant-tabs-bar {
    margin: 0;
  }

  .ant-tabs-nav {
    width: 100%;
    display: block;
    background: var(--site-theme-bg-translucent);
    backdrop-filter: blur(2px);

    &>div {
      display: flex;
    }
  }

  .ant-tabs-nav .ant-tabs-tab {
    flex: 1;
    margin: 0;
    padding: 0;

    .bg-tab--first {
      background: url('~common/assets/images/premium/bg-version-tab-first.svg') no-repeat;
    }

    .bg-tab--middle {
      background: url('~common/assets/images/premium/bg-version-tab-middle.svg') no-repeat;
    }

    .bg-tab--last {
      background: url('~common/assets/images/premium/bg-version-tab-last.svg') no-repeat;
    }
  }

  .ant-tabs-nav .ant-tabs-tab::before {
    width: 0;
    transition: none;
  }

  .ant-tabs-tab:not(.ant-tabs-tab-active)+.ant-tabs-tab:not(.ant-tabs-tab-active) {
    &::before {
      content: '';
      width: 1px;
      height: 24px;
      background: var(--site-theme-divider);
      position: absolute;
      left: 0;
      top: 50%;
      margin-top: -12px;
    }
  }

  .ant-tabs-ink-bar {
    display: none !important;
  }

  .ant-tabs .ant-tabs-tabpane {
    transition: none;
  }

  .vip-trigger-left,
  .vip-trigger-right {
    min-height: 412px;
  }

  .vip-trigger-left {
    padding: 16px;
    background: url('~common/assets/images/premium/bg-version-left.svg') no-repeat;
    background-position: top -14px right 16px;

    &>hr {
      margin: 10px 0 0;
      background: var(--site-theme-divider-light);
    }

    &>ul {
      gap: 8.5px;
    }
  }

  .vip-trigger-right {
    canvas {
      display: block;
    }
  }

  .vip-trigger-code {
    padding: 9px;
  }


  .groupbuy-trigger-left,
  .groupbuy-trigger-right {
    min-height: 412px;
  }

  .groupbuy-trigger-left {
    padding: 16px;
    background: url('~common/assets/images/premium/bg-version-left.svg') no-repeat;
    background-position: top -14px right 16px;

    &>hr {
      margin: 10px 0 0;
      background: var(--site-theme-divider-light);
    }

    &>ul {
      gap: 8.5px;
    }
  }

  .groupbuy-trigger-right {
    canvas {
      display: block;
    }
  }

  .groupbuy-trigger-code {
    padding: 9px;
  }
}
</style>
