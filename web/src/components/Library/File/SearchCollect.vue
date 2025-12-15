<template>
  <Button
    :loading="isLoading"
    class="btn"
    :disabled="currentCollected || !userStore.userInfo"
    :type="currentCollected ? 'default' : 'primary'"
    @click="handleClick"
    >{{
      !userStore.userInfo
        ? $t('common.tips.loginFirst')
        : currentCollected
          ? $t('home.upload.added')
          : $t('home.upload.add')
    }}</Button
  >
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import { Button } from 'ant-design-vue'
import {
  ERROR_CODE_DOC_LIMIT,
  REQUEST_SERVICE_NAME_APP,
} from '@/common/src/api/const'
import { ResponseError, SuccessResponse } from '@/common/src/api/type'
import {
  CollectPaperRequest,
  MyCollectedDocInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/userCenter/UserDoc'
import api from '@/common/src/api/axios'
import { useVipStore } from '@/common/src/stores/vip'
import { ElementName, PageType } from '@/utils/report'
import { useUserStore } from '@/common/src/stores/user'
import { NeedVipException } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo'
import { getDomainOrigin } from '@/common/src/utils/env'
import { useI18n } from 'vue-i18n'
// import { LimitDialogReportParams } from '~/store/vip'

const i18n = useI18n()
const vipStore = useVipStore()
const userStore = useUserStore()

const props = defineProps({
  paperId: {
    type: String,
    default: '',
  },
  isCollected: {
    type: Boolean,
    default: false,
  },
  folderId: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['addSuccess'])
const currentCollected = ref<boolean>(props.isCollected)
const isLoading = ref<boolean>(false)

const collectPaper = async (param: Partial<CollectPaperRequest>) => {
  const res = await api.post<SuccessResponse<MyCollectedDocInfo>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/collectPaper`,
    param
  )
  return res.data.data || {}
}

const handleClick = async () => {
  isLoading.value = true
  try {
    await collectPaper({
      paperId: props.paperId,
      paperTitle: '',
      folderId: props.folderId || '0',
    } as any)

    currentCollected.value = true

    isLoading.value = false
    emit('addSuccess')
  } catch (error) {
    const e = error as ResponseError
    if (e?.code === ERROR_CODE_DOC_LIMIT) {
      // 达到收藏上限
      vipStore.showVipLimitDialog(e.message, {
        exception: e.extra as NeedVipException,
        leftBtn: {
          text: i18n.t('common.premium.btns.more'),
          url: `${getDomainOrigin()}/home/mine`,
        },
        reportParams: {
          page_type: PageType.library,
          element_name: ElementName.upperCollectionPopup,
        },
      })
    }
    if ((error as Error)?.message === '已存在重复文献') {
      currentCollected.value = true
    }
    isLoading.value = false
  }
}
</script>
<style lang="less" scoped>
.uploader-search-input {
  .input {
    margin-top: 12px;
    margin-bottom: 24px;
  }
  .item {
    display: flex;
    justify-content: space-between;
    .paper {
      margin-right: 36px;
      overflow: hidden;
    }
    .btn {
      padding: 0 16px;
    }
  }
  .item + .item {
    margin-top: 14px;
  }
}
</style>
