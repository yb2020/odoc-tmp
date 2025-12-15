<template>
  <div
    ref="collectRef"
    style="display: flex; align-items: center"
    @click="showCollectDialog"
  >
    <slot />
  </div>
</template>
<script lang="ts">
/* eslint-disable vue/one-component-per-file */
import {
  computed,
  defineComponent,
  onMounted,
  PropType,
  ref,
  createApp,
} from 'vue'
import tippy, { roundArrow } from 'tippy.js'
import CollectPanel from './CollectPanel.vue'
import { collectPaper, cancelCollectPaper } from '@/common/src/api/paper'
import { ResponseError } from '@/common/src/api/type'
import { useUserStore } from '@/common/src/stores/user'
import { useVipStore } from '@/common/src/stores/vip'
import { useI18n } from 'vue-i18n'
import { getDomainOrigin } from '@/common/src/utils/env'

const ERROR_CODE_DOC_LIMIT = 2

export default defineComponent({
  props: {
    paperId: {
      type: String,
      default: '',
    },
    isCollected: {
      type: Boolean,
      default: false,
    },
    paperTitle: {
      type: String,
      default: '',
    },
    folderId: {
      type: String,
      default: '',
    },
    stopPropagation: Boolean,
    placement: {
      type: String as PropType<
        | 'top'
        | 'top-start'
        | 'top-end'
        | 'right'
        | 'right-start'
        | 'right-end'
        | 'bottom'
        | 'bottom-start'
        | 'bottom-end'
        | 'left'
        | 'left-start'
        | 'left-end'
        | 'auto'
        | 'auto-start'
        | 'auto-end'
      >,
      default: 'right-start',
    },
    collectLimitDialogReportParams: {
      type: Object, // as PropType<LimitDialogReportParams>,
      default: null,
    },
  },
  setup(props, { emit }) {
    const userStore = useUserStore()
    const vipStore = useVipStore()
    const isMobile = ref(false) // TODO
    const showCollectDialog = async (e: MouseEvent) => {
      if (props.stopPropagation) {
        e.stopPropagation()
      }
      if (!userStore.isLogin()) {
        userStore.openLogin()
        return
      }
      try {
        emit('loadingChange', true, props.isCollected)
        if (props.isCollected) {
          await cancelCollectPaper({
            paperId: props.paperId,
          })
          emit('loadingChange', false, false)
        } else {
          await collectPaper({
            paperId: props.paperId,
            paperTitle: props.paperTitle,
            folderId: props.folderId,
          } as any)
          emit('loadingChange', false, true)
          if (!isMobile.value) {
            tippyInstance?.show()
          }
        }
      } catch (error) {
        if ((error as Error)?.message === '已存在重复文献') {
          emit('loadingChange', false, true)
          if (!isMobile.value) {
            tippyInstance?.show()
          }
          return
        }
        const e = error as ResponseError
        if (e?.code === ERROR_CODE_DOC_LIMIT) {
          emit('loadingChange', false, props.isCollected)
          // 达到收藏上限
          vipStore.showVipLimitDialog(e.message, {
            exception: e.extra as any,
            leftBtn: {
              text: i18n.t('common.premium.btns.more'),
              url: `${getDomainOrigin()}/home/mine`,
            },
            reportParams: props.collectLimitDialogReportParams,
          })
          return
        }
        emit('loadingChange', false, props.isCollected)
      }
    }

    const collectRef = ref<HTMLDivElement>()

    let tippyInstance: any

    const i18n = useI18n()

    onMounted(() => {
      if (collectRef.value) {
        const app = createApp(CollectPanel, {
          title: i18n.t('paper.collectModal.title'),
          newText: i18n.t('paper.collectModal.new'),
          paperId: props.paperId,
          stopPropagation: props.stopPropagation,
          onClose: () => {
            tippyInstance?.hide()
            emit('close')
          },
        })

        const instance = app.mount(document.createElement('div')) as any

        tippyInstance = tippy(collectRef.value, {
          content: instance.$el,
          trigger: 'manual',
          hideOnClick: true,
          arrow: roundArrow,
          maxWidth: 400,
          interactive: true,
          appendTo: document.body,
          offset: [0, 8],
          placement: props.placement,
          theme: 'light',
          onShow() {
            instance.fetchClassifyData()
          },
        })
      }
    })

    return {
      showCollectDialog,
      collectRef,
    }
  },
})
</script>
