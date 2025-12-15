<template>
  <div @click="handleClick">
    <slot />
    <Modal
      v-model:visible="dialogVisible"
      :title="`${$t('home.paper.cite.title')}`"
      destroy-on-close
      width="560px"
      class="cite-dialog"
      :footer="null"
    >
      <CitePaneVue
        ref="citePaneVueRef"
        :paper-id="paperId"
        :page-type="pageType"
        :type-parameter="typeParameter"
        :pdf-id="pdfId"
        :has-update-btn="true"
        @update:success="onUpdateSuccess"
        @clickManageCitationStyle="handleClickManageCitationStyle"
      />
    </Modal>

    <Modal
      v-model:visible="manageCitationStyleVisible"
      destroy-on-close
      :title="`${$t('home.paper.citeManage.title')}`"
      :footer="null"
      :width="560"
      wrap-class-name="manage-citation-style-modal"
      @cancel="handleCancel"
    >
      <div class="modal-content">
        <ManageCitationStyle></ManageCitationStyle>
      </div>
    </Modal>
  </div>
</template>
<script lang="ts">
import { defineComponent, ref } from 'vue'
import CitePaneVue from './CitePane.vue'
import ManageCitationStyle from './manage.vue'
import { initClipBoradButton } from '@/common/src/utils/clipboard'
import { PageType, reportElementClick } from '@/utils/report'
import { Modal } from 'ant-design-vue'

initClipBoradButton('.js-cite-copy-btn')

export default defineComponent({
  components: { CitePaneVue, ManageCitationStyle, Modal },
  props: {
    paperId: {
      type: String,
      default: '',
    },
    visible: {
      type: Boolean,
      default: false,
    },
    pageType: {
      type: String, // TODO as () => PageType,
      default: '',
    },
    typeParameter: {
      type: String,
      default: 'none',
    },
    pdfId: {
      type: String,
      default: '',
    },
    elementName: {
      type: String,
      default: 'cite',
    },
  },
  emits: ['update:success'],
  setup(props, { emit }) {
    const dialogVisible = ref<boolean>(props.visible)

    const citePaneVueRef = ref()

    const manageCitationStyleVisible = ref<boolean>(false)

    const handleClickManageCitationStyle = () => {
      manageCitationStyleVisible.value = true
    }

    const handleCancel = () => {
      manageCitationStyleVisible.value = false

      citePaneVueRef.value?.fetchMyCslList()
    }

    const onUpdateSuccess = () => {
      emit('update:success')
    }

    const handleClick = () => {
      dialogVisible.value = true

      reportElementClick({
        page_type: props.pageType,
        type_parameter: props.typeParameter,
        element_name: props.elementName,
        status: 'none',
      })
    }

    return {
      dialogVisible,
      citePaneVueRef,
      manageCitationStyleVisible,
      handleClickManageCitationStyle,
      handleCancel,
      onUpdateSuccess,
      handleClick,
    }
  },
})
</script>
<style lang="less">
.cite-dialog,
.manage-citation-style-modal {
  .ant-modal-header {
    border: none !important;
    padding: 16px 20px;
    .ant-modal-title {
      font-family: 'Noto Sans SC';
      font-weight: 600;
      font-size: 15px;
      line-height: 24px;
      color: #1d2229;
    }
  }
  .ant-modal-close {
    .ant-modal-close-x {
      height: 56px;
      line-height: 56px;
      font-size: 16px;
    }
  }
}
.manage-citation-style-modal {
  .ant-modal-content {
    margin: 32px 0 0 24px;
    width: 100%;
    .modal-content {
      margin-top: -20px;
    }
  }
}
</style>
