<script lang="ts">
import { defineComponent, onMounted, onUnmounted, ref } from 'vue'
import { postMessage } from '@/utils/message'
import { reportPopupPaperReferenceClick } from '@/utils/report'
// Tabs组件在模板中使用，不要删除此导入
import Tabs from '@/components/Meta/Tabs.vue' // eslint-disable-line

export default defineComponent({
  components: {
    Tabs,
  },
  props: {
    paperId: {
      type: String,
      default: '',
    },
    pdfId: {
      type: String,
      default: '',
    },
    triggerIframe: {
      type: Boolean,
      default: false,
    },
    pageType: {
      type: String as () => string, // PageType,
      default: '',
    },
    getLiteratureFormat: {
      type: Function as unknown as () => () => string | undefined,
      required: true,
    },
  },
  emits: ['update:success'],
  setup(props, { emit }) {
    const loading = ref(false)
    const showUpdateDialog = ref(false)

    const handleOpenUpdateDialog = () => {
      if (loading.value) {
        return
      }

      reportPopupPaperReferenceClick({
        page_type: props.pageType,
        literature_format: props.getLiteratureFormat() || '',
        element_name: 'renew',
      })

      if (props.triggerIframe) {
        loading.value = true
        postMessage({
          event: 'clickUpdateCitation',
          params: {
            paperId: props.paperId,
            pdfId: props.pdfId,
          },
        })
        setTimeout(() => {
          loading.value = false
        }, 10000)
        return
      }

      // 直接显示对话框，不再加载外部库
      showUpdateDialog.value = true
    }
    
    const handleCloseDialog = () => {
      showUpdateDialog.value = false
    }
    
    const handleUpdateSuccess = () => {
      showUpdateDialog.value = false
      emit('update:success')
    }

    const onTopWindowMessage = (event: MessageEvent) => {
      const { data } = event
      if (data.event === 'stopLoadingClickUpdateCitation') {
        loading.value = false
      }
    }

    onUnmounted(() => {
      window.removeEventListener('message', onTopWindowMessage)
    })

    onMounted(() => {
      if (props.triggerIframe) {
        window.addEventListener('message', onTopWindowMessage)
      }
    })

    return {
      loading,
      showUpdateDialog,
      handleOpenUpdateDialog,
      handleCloseDialog,
      handleUpdateSuccess,
    }
  },
})
</script>
<template>
  <div>
    <a-button 
      :loading="loading" 
      type="link" 
      @click="handleOpenUpdateDialog()"
    >
      {{ $t('home.paper.cite.update') }}
    </a-button>
    
    <a-modal
      v-model:visible="showUpdateDialog"
      :title="$t('meta.title')"
      :footer="null"
      :width="800"
      @cancel="handleCloseDialog"
    >
      <Tabs
        v-if="showUpdateDialog"
        :paper-id="paperId"
        :pdf-id="pdfId"
        :page-type="pageType"
        @update:success="handleUpdateSuccess"
        @cancel="handleCloseDialog"
      />
    </a-modal>
  </div>
</template>
<style lang="less" scoped>
.ant-btn-link {
  color: #1f71e0;
}
</style>
