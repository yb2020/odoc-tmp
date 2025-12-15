<template>
  <div
    v-loading="loading"
    data-loading-size="default"
    class="cite-wrap"
    :style="{
      'max-height': isClientIframe ? '170px' : 'unset',
      'min-height': isClientIframe ? '170px' : '50px',
    }"
  >
    <cite
      v-if="renderStatus === CitationRenderStatus.successful"
      :class="['cite', { bibtex: isBibTexStyle }]"
      @click="handleSelect"
      v-html="citationText"
    ></cite>

    <Result v-else-if="!loading" status="warning" title="">
      <template #subTitle>
        <div>{{ citationText }}</div>
        <div v-if="renderStatus === CitationRenderStatus.failed">
          {{ $t('home.paper.cite.warningMessage.renderError') }}
        </div>
      </template>
    </Result>
  </div>
</template>
<script lang="ts">
import { defineComponent, ref, onMounted, PropType, computed } from 'vue'
import { Result } from 'ant-design-vue'
import { PaperData } from '@/utils/citation'
import { getRenderCitation, CitationRenderStatus } from '@/utils/citation'
import { useUserStore } from '@/common/src/stores/user'
import { useI18n } from 'vue-i18n'
import { Language } from 'go-sea-proto/gen/ts/lang/Language';

export default defineComponent({
  props: {
    isClientIframe: {
      type: Boolean,
      default: false,
    },
    metaData: {
      type: Object as PropType<PaperData>,
      default: () => ({}),
    },
    fileUrl: {
      type: String,
      default: '',
    },
    isEN: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const userStore = useUserStore()

    const { t } = useI18n()

    const loading = ref<boolean>(true)

    const citationText = ref<string>('')

    const renderStatus = ref()

    onMounted(async () => {
      const res = await getRenderCitation(
        props.fileUrl,
        props.metaData,
        t,
        'html',
        props.metaData?.language || Language.EN_US
      )

      citationText.value = res?.data || ''

      renderStatus.value = res?.status

      loading.value = false
    })

    const handleSelect = (e: MouseEvent) => {
      if (!props.isEN) {
        if (!userStore.isLogin()) {
          return
        }
      }

      const text = e.target as any
      if ((document.body as any).createTextRange) {
        const range = (document.body as any).createTextRange()
        range.moveToElementText(text)
        range.select()
      } else if (window.getSelection) {
        const selection = window.getSelection()
        if (!selection) {
          return
        }
        const range = document.createRange()
        range.selectNodeContents(text)
        selection.removeAllRanges()
        selection.addRange(range)
      }
    }

    const isBibTexStyle = computed(() => props.fileUrl.includes('bibtex'))
    return {
      citationText,
      handleSelect,
      isBibTexStyle,
      renderStatus,
      CitationRenderStatus,
      loading,
    }
  },
})
</script>
<style lang="less" scoped>
.cite-wrap {
  overflow: auto;
  border-left: 1px solid #dfe6f0;
  border-right: 1px solid #dfe6f0;
  border-bottom: 1px solid #dfe6f0;
  margin-top: -20px;
  padding: 8px 16px;
  border-radius: 2px;
  margin-bottom: 8px;
  position: relative;
  min-height: 50px;
}
.cite {
  font-size: 14px;
  font-weight: 400;
  line-height: 19.6px;
  font-style: normal;
  color: #262625;
}
cite.bibtex {
  white-space: pre;
}
/deep/.ant-result {
  padding: 20px;
  .ant-result-icon {
    margin-bottom: 0;
    .anticon {
      font-size: 40px;
    }
  }
}
/deep/.csl-entry {
  font-size: 14px;
  font-weight: 400;
  line-height: 24px;
  color: #262625;
  .csl-left-margin,
  .csl-right-inline {
    display: inline;
  }
}
</style>
