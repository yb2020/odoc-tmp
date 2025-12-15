<template>
  <div v-loading="loading" class="cite-dialog" data-loading-size="default">
    <div v-if="fetchState && fetchState.error" class="error" @click="fetch">
      <ReloadOutlined /> {{ fetchState.error.message }}
    </div>

    <div
      v-else-if="
        metaData &&
        showTabCitationStyleList &&
        showTabCitationStyleList.length &&
        !loading
      "
      class="content"
    >
      <a-tabs v-model="activeKey" :animated="false" @change="onTabChange">
        <a-tab-pane
          v-for="citationStyle in showTabCitationStyleList"
          :key="citationStyle.id"
          :tab="
            citationStyle.customDefineTitle ||
            citationStyle.shortTitle ||
            citationStyle.title
          "
        >
          <CiteContent
            :is-client-iframe="isClientIframe"
            :meta-data="metaData"
            :file-url="citationStyle.fileUrl"
            :is-e-n="isWebEN"
          ></CiteContent>
        </a-tab-pane>
        <span
          slot="tabBarExtraContent"
          class="manage"
          @click="handleManageCitationStyle"
        >
          <SettingOutlined class="setting-icon" />
          {{ $t('home.paper.cite.manage') }}
        </span>
      </a-tabs>

      <div class="footer">
        <div>
          <span class="text">{{ $t('home.paper.cite.export') }}</span>
          <a-button
            v-if="!authenticated && !isWebEN"
            type="primary"
            class="btn"
            @click="handleLogin"
          >
            {{ CitationStyle.BIBTEX }}
          </a-button>
          <a
            v-else
            class="btn"
            :href="bibTextBase64"
            download="citation.bib"
            @click="
              handleReportPopupPaperReferenceClick(
                PopupPaperReferenceElementName.bibTex
              )
            "
          >
            <a-button type="primary" :disabled="!bibTextBase64">
              {{ CitationStyle.BIBTEX }}
            </a-button>
          </a>

          <a-button
            v-if="!authenticated && !isWebEN"
            type="primary"
            class="btn"
            @click="handleLogin"
          >
            {{ CitationStyle.ENDNOTE }}
          </a-button>
          <a
            v-else
            class="btn"
            :href="endnoteTextBase64"
            download="citation.enw"
            @click="
              handleReportPopupPaperReferenceClick(
                PopupPaperReferenceElementName.endNote
              )
            "
          >
            <a-button type="primary" :disabled="!endnoteTextBase64">
              {{ CitationStyle.ENDNOTE }}
            </a-button>
          </a>
          <a-button
            v-if="authenticated"
            class="btn copy-btn"
            :data-clipboard-text="copyText"
            :disabled="!copyText"
            @click="handleCopy"
            ><CopyOutlined />{{ $t('home.paper.cite.copy') }}</a-button
          >
        </div>
        <a-button
          v-if="!authenticated"
          class="login-copy-btn"
          @click="handleLogin"
          ><LoginOutlined />{{ $t('home.paper.cite.loginToCopy') }}</a-button
        >
        <CiteUpdate
          v-if="authenticated && hasUpdateBtn"
          :paper-id="paperId"
          :pdf-id="pdfId"
          :trigger-iframe="triggerIframe"
          :get-literature-format="getLiteratureFormat"
          :page-type="pageType"
          @update:success="onUpdateSuccess"
        />
      </div>
    </div>

    <div v-else-if="!loading && isShowEmpty" class="empty">
      {{ $t('home.paper.cite.empty') }}
    </div>
  </div>
</template>
<script lang="ts">
import {
  computed,
  defineComponent,
  ref,
  watch,
  onMounted,
  onUnmounted,
} from 'vue'
import { encode } from 'js-base64'
import copyTextToClipboard from 'copy-to-clipboard-ultralight'
import { message } from 'ant-design-vue'
import {
  CopyOutlined,
  LoginOutlined,
  ReloadOutlined,
  SettingOutlined,
} from '@ant-design/icons-vue'
import { PaperData } from '@/utils/citation'
import { MyCslItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/CSL'
import CiteContent from './citeContent.vue'
import CiteUpdate from './CiteUpdate.vue'
import {
  CitationStyle,
  getRenderCitation,
  getEndNoteCitation,
} from '@/utils/citation'
import {
  getMyCslList,
  getDefaultCslList,
  getDocMetaInfo,
} from '@common/api/citation'
import {
  PageType,
  reportPopupPaperReferenceClick,
  PopupPaperReferenceElementName,
  reportElementImpression,
  EventCode,
  ElementName,
} from '@common/utils/report'
import useFetch from '@common/hooks/useFetch'
import { useUserStore } from '@/common/src/stores/user'
import { useLanguage } from '@/hooks/useLanguage'
import { useI18n } from 'vue-i18n'
import { Language } from 'go-sea-proto/gen/ts/lang/Language'

export default defineComponent({
  components: {
    CiteContent,
    CiteUpdate,
    ReloadOutlined,
    SettingOutlined,
    CopyOutlined,
    LoginOutlined,
  },
  props: {
    paperId: {
      type: String,
      default: '',
    },
    isClientIframe: {
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
    hasUpdateBtn: {
      type: Boolean,
      default: false,
    },
    triggerIframe: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:success', 'clickManageCitationStyle'],
  fetchOnServer: false,
  setup(props, { emit }) {
    const { isEnUS } = useLanguage()
    const isWebEN = isEnUS // 保持向后兼容的命名
    const { t } = useI18n()
    const { getCurrentLanguage } = useLanguage()

    const userStore = useUserStore()

    const authenticated = computed(() => userStore.isLogin())

    const defaultCitationStyleList = ref<MyCslItem[]>([])

    const metaData = ref<PaperData>()

    const personalCitationStyleList = ref<MyCslItem[]>([])

    const showTabCitationStyleList = ref<MyCslItem[]>([])

    const loading = computed(() => fetchState.pending && !fetchState.error)

    const activeKey = ref<string>('')

    const activeFileUrl = computed(() => {
      if (!activeKey.value) return ''

      const find = showTabCitationStyleList.value.find((item) => {
        return item.id === activeKey.value
      })

      if (find) return find.fileUrl

      return ''
    })

    const bibTexFileUrl = ref<string>('')

    const isShowEmpty = ref<boolean>(false)

    useFetch(async () => {
      const res = await getDefaultCslList()

      defaultCitationStyleList.value = res.list || []

      bibTexFileUrl.value = defaultCitationStyleList.value[0].fileUrl

      if (!authenticated.value) {
        // 国际版的去掉GB/T 7714
        showTabCitationStyleList.value = isWebEN.value
          ? defaultCitationStyleList.value.filter(
              (item) => item.customDefineTitle !== 'GB/T 7714'
            )
          : defaultCitationStyleList.value

        activeKey.value = showTabCitationStyleList.value[0].id
      }
    })

    const { fetchState, fetch } = useFetch(async () => {
      const paperId =
        props.paperId && props.paperId !== '0' ? props.paperId : undefined

      const pdfId = props.pdfId && props.pdfId !== '0' ? props.pdfId : undefined

      const res = await getDocMetaInfo({
        paperId,
        pdfId,
      })

      if (res) {
        metaData.value = {
          ...res,
          // 判断字段是否符合中英文，否则降级英文
          language:
            res.language === Language.EN_US || res.language === Language.ZH_CN
              ? res.language
              : Language.EN_US,
          paperId: props.paperId && props.paperId !== '0' ? props.paperId : '',
          pdfId: props.pdfId && props.pdfId !== '0' ? props.pdfId : '',
          docName: res.title,
          authors: res?.authorList?.map((x) => x.literal),
          displayPublishDate: {
            publishDate: res.publishTimestamp,
            originPublishDate: res.publishTimestamp,
            rollbackEnable: false,
          },
          venues: res.containerTitle || [],
          primaryVenue: '',
          authorList: (res.authorList as any) || [],
          publishDate: res.publishTimestamp || '',
          doi: res.doi || undefined,
          url: res.url || undefined,
          issue: res.issue || undefined,
          volume: res.volume || undefined,
          docType: res.docType || 'article-journal',
        }
      } else {
        isShowEmpty.value = true
      }
    })

    const { fetch: fetchMyCslList } = useFetch(async () => {
      if (!authenticated.value) return

      const res = await getMyCslList()

      personalCitationStyleList.value = res.list || []

      const partialPersonalCitationStyleList =
        personalCitationStyleList.value.slice(0, 5)

      // 国际版的去掉GB/T 7714
      showTabCitationStyleList.value = isWebEN.value
        ? partialPersonalCitationStyleList.filter(
            (item) => item.customDefineTitle !== 'GB/T 7714'
          )
        : partialPersonalCitationStyleList

      const findIndex = showTabCitationStyleList.value.findIndex((item) => {
        return item.id === activeKey.value
      })

      if (!activeKey.value || findIndex === -1) {
        activeKey.value = showTabCitationStyleList.value[0]?.id
      }

      if (!showTabCitationStyleList.value?.length) {
        isShowEmpty.value = true
      }
    })

    const onTabChange = (key: string) => {
      activeKey.value = key
    }

    const bibText = ref<string>('')

    const copyText = ref<string>('')

    const bibTextBase64 = computed(() => {
      if (!bibText.value) return ''

      return `data:application/octet-stream;charset=utf-16le;base64,${encode(
        bibText.value
      )}`
    })

    const endnoteTextBase64 = computed(() => {
      if (!metaData.value) return ''

      return `data:application/octet-stream;charset=utf-16le;base64,${encode(
        getEndNoteCitation(metaData.value)
      )}`
    })

    const handleCopy = () => {
      copyTextToClipboard(copyText.value)
      message.success(t('home.paper.cite.copySuccessful') as string)
      handleReportPopupPaperReferenceClick(PopupPaperReferenceElementName.copy)
    }

    const handleManageCitationStyle = () => {
      if (!authenticated.value) {
        userStore.openLogin()

        return
      }

      emit('clickManageCitationStyle')
    }

    const handleRenderCitation = async (style: string, metaData: PaperData) => {
      const res = await getRenderCitation(
        style,
        metaData,
        t,
        'text',
        typeof metaData.language === 'string' ? metaData.language : Language.EN_US
      )
      if (res?.status === 1) return res?.data

      return ''
    }

    watch(
      () => metaData.value,
      async (newVal) => {
        if (newVal) {
          if (bibTexFileUrl.value) {
            bibText.value = await handleRenderCitation(
              bibTexFileUrl.value,
              newVal
            )
          }

          if (activeFileUrl.value) {
            copyText.value = await handleRenderCitation(
              activeFileUrl.value,
              newVal
            )
          }
        }
      }
    )

    const handleChangeStyleSuccess = () => {
      fetchMyCslList()
    }

    watch(
      () => activeFileUrl.value,
      async (newVal) => {
        if (newVal && metaData.value) {
          copyText.value = await handleRenderCitation(
            activeFileUrl.value,
            metaData.value
          )
        }
      }
    )

    const getLiteratureFormat = () => {
      return showTabCitationStyleList.value
        .find((item) => item.id === activeKey.value)
        ?.title.toLowerCase()
    }

    const handleReportPopupPaperReferenceClick = (element: string) => {
      const format = getLiteratureFormat()
      if (typeof format === 'string') {
        reportPopupPaperReferenceClick({
          page_type: props.pageType,
          type_parameter: props.typeParameter,
          literature_format: format,
          element_name: element,
        })
      }
    }

    const handleSelectCopy = (event: any) => {
      const flag =
        event.target.className.includes('csl-entry') ||
        event.target?.parentNode?.className?.includes('csl-entry')

      if (flag) {
        handleReportPopupPaperReferenceClick(
          PopupPaperReferenceElementName.selectCopy
        )
      }
    }

    onMounted(() => {
      window.addEventListener('copy', handleSelectCopy)

      reportElementImpression(EventCode.readpaperElementImpression, {
        page_type: props.pageType,
        element_name: ElementName.popupReference,
        type_parameter: props.typeParameter,
        element_parameter: 'none',
      })
    })

    onUnmounted(() => {
      window.removeEventListener('copy', handleSelectCopy)
    })

    const handleLogin = () => {
      userStore.openLogin()
    }

    const onUpdateSuccess = () => {
      fetch()
      emit('update:success')
    }
    return {
      CitationStyle,
      bibText,
      copyText,
      bibTextBase64,
      endnoteTextBase64,
      handleCopy,
      activeKey,
      onTabChange,
      defaultCitationStyleList,
      personalCitationStyleList,
      showTabCitationStyleList,
      fetchState,
      fetch,
      loading,
      handleManageCitationStyle,
      fetchMyCslList,
      handleChangeStyleSuccess,
      PopupPaperReferenceElementName,
      handleReportPopupPaperReferenceClick,
      metaData,
      activeFileUrl,
      bibTexFileUrl,
      isShowEmpty,
      isWebEN,
      authenticated,
      handleLogin,
      onUpdateSuccess,
      getLiteratureFormat,
    }
  },
})
</script>
<style lang="less" scoped>
@import './scrollbar.less';

.cite-dialog {
  background-color: var(--site-theme-background-primary);
  position: relative;
  min-height: 100px;
  .error {
    text-align: center;
    padding-top: 38px;
    cursor: pointer;
    color: var(--site-theme-text-primary);
  }
  .empty {
    text-align: center;
    padding-top: 38px;
    color: var(--site-theme-text-secondary);
  }
  :deep(.ant-tabs-nav-wrap) {
    margin-left: 16px;
  }
  .copy-btn {
    font-size: 13px;
  }
  .footer {
    // margin-top: 24px;
    padding-top: 16px;
    display: flex;
    justify-content: space-between;
    .text {
      margin-right: 16px;
      color: var(--site-theme-text-tertiary);
    }
    .btn + .btn {
      margin-left: 16px;
    }
  }
  .manage {
    display: inline-block;
    padding: 2px 12px;
    background: var(--site-theme-background-secondary);
    border-radius: 2px;
    font-family: 'Noto Sans SC';
    font-style: normal;
    font-weight: 400;
    font-size: 12px;
    line-height: 20px;
    color: var(--site-theme-text-secondary);
    cursor: pointer;
    .setting-icon {
      margin-right: 6px;
    }
  }
  .login-copy-btn {
    background: var(--site-theme-primary-color);
    color: #fff;
    border-color: var(--site-theme-primary-color);
    &:hover {
      background: var(--site-theme-primary-color);
      color: #fff;
      border-color: var(--site-theme-primary-color);
    }
  }
}
:deep(.ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn) {
    color: var(--site-theme-text-primary);
    text-shadow: 0 0 0.25px currentcolor;
}
:deep(.ant-tabs-tab) {
    color: var(--site-theme-text-secondary);
}
</style>
