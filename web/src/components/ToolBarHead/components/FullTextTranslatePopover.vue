<template>
  <div
    ref="fulltextPopoverRef"
    class="fulltext-translate-popover"
  >
    <div v-if="status === FullTranslateFlowStatus.TRANSLATE_FINISHED">
      <div class="switcher">
        <span>{{ $t('viewer.fullTextScroll') }}</span>
        <a-switch
          :checked="
            fullTextTranslateStore.scrollMode === PDFWebviewScrollMode.lock
          "
          @change="handleScrollMode"
        />
      </div>
      <div class="switcher">
        <span>{{ $t('viewer.fullTextOriginal') }}</span>
        <a-switch
          :checked="
            fullTextTranslateStore.previewMode ===
              PDFWebviewPreviewMode.withOriginalPDF
          "
          @change="handlePreviewMode"
        />
      </div>
      <div class="switcher">
        <span>{{ $t('viewer.exportTranslatedPDF') }}</span>
        <a-button
          type="primary"
          size="noraml"
          @click="exportTranslatedPDF"
        >
          {{
            $t('viewer.exportTranslatedPDFButtonText')
          }}
        </a-button>
      </div>
    </div>
    <a-spin
      v-else
      :spinning="pending"
    >
      <div class="content">
        <div class="title">
          <a 
            :href="userGuideLink" 
            target="_blank" 
            class="user-guide-link"
          >
            {{ $t('viewer.useCase') }}
          </a>
        </div>
        <div
          v-if="fetchError"
          class="error"
          @click="initRightInfoPopover"
        >
          {{ fetchError.message }}
          <redo-outlined />
        </div>
        <div
          v-if="rightInfo"
          class="rights"
        >
          <div
            v-if="rightInfo.ruleDesc"
            class="rules"
          >
            <div class="info">
              {{ $t('viewer.rule') }}
            </div>
            <!-- <div v-html="rightInfo.ruleDesc" /> -->
            <div class="additional-rules">
              <div class="rule-item">1. {{ $t('viewer.rule1') }}</div>
              <div class="rule-item">2. {{ $t('viewer.rule2') }}</div>
              <div class="rule-item">3. {{ $t('viewer.rule3') }}</div>
            </div>
          </div>
          <div class="bottom">
            <div>
              {{
                $t('viewer.fullTextTip1', {
                  points: userStore.getTotalCredits()
                })
              }}
            </div>
            <div
              v-if="status === FullTranslateFlowStatus.TRANSLATING"
              class="doing"
            >
              {{ $t('viewer.fullTextTip2', { progress: formatedWaiting }) }}
            </div>
            <div v-else class="translate-action-wrapper">
              <div v-if="hasRight">
                <a-button
                  :loading="translateLoading"
                  type="primary"
                  @click="handleTranslate"
                  :disabled="!isReadyForTranslation"
                >
                  {{ $t('viewer.startTranslation') }}
                </a-button>
                <div v-if="!isReadyForTranslation" class="unparsed-tip">
                  {{ $t('viewer.unparsedDocumentTip') }}
                </div>
              </div>
              <Trigger
                v-else-if="vipStore.seniorRole"
                visible
                :need-vip-type="vipStore.seniorRole"
                :report-params="{
                  pageType: PageType.note,
                  elementName:
                    vipStore.role.vipType > VipType.FREE
                      ? ElementName.upperPaperTranslateLimitPopup
                      : ElementName.upperPaperTranslatePopup,
                }"
                :btn-txt="`${$t('common.premium.btns.get')} ${$t(
                  `common.premium.${
                    vipStore.seniorRole <
                    StepwiseVipTypes[StepwiseVipTypes.length - 1]
                      ? 'senior'
                      : `versions.${
                        PremiumVipPreferences[vipStore.seniorRole].key
                      }`
                  }`
                )}`"
              />
              <div
                v-else
                class="error"
              >
                {{ $t('viewer.fullTextTip3') }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </a-spin>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useDocStore } from '~/src/stores/docStore';
import { storeToRefs } from 'pinia';
import { UserDocParsedStatusEnum } from 'go-sea-proto/gen/ts/doc/UserDocParsedStatus';
import { GetFullTextTranslateRightInfoResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/translate/FullTextTranslate';
import {
  useFullTextTranslateStore,
  PDFWebviewScrollMode,
  PDFWebviewPreviewMode,
} from '~/src/stores/fullTextTranslateStore';
import {
  useVipStore,
  PremiumVipPreferences,
  VipType,
} from '@common/stores/vip';
import {
  getRightInfo,
  reportTranslateExposure,
  getTranslateStatus,
} from '@/api/fullTextTranslate';
import { ResponseError } from '~/src/api/type';
import {
  PageType,
  ElementClick,
  reportElementClick,
  ElementName,
} from '~/src/api/report';
import { RedoOutlined } from '@ant-design/icons-vue';
import { getHostname, isInElectron } from '~/src/util/env';
import Trigger from '@common/components/Premium/Trigger.vue';
import { StepwiseVipTypes } from '@common/components/Premium/types';
import { currentNoteInfo } from '~/src/store';
import { FullTranslateFlowStatus } from 'go-sea-proto/gen/ts/translate/FullTextTranslate';
import { useI18n } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useLanguage } from '@/hooks/useLanguage';
import { useUserStore } from '@common/stores/user'

const { locale } = useI18n();

// 语言管理
const { isCurrentLanguage } = useLanguage();

// 创建用户指南链接的计算属性
const userGuideLink = computed(() => {
  if (isCurrentLanguage(Language.ZH_CN)) {
    return '/docs/zh/guide';
  }
  return '/docs/guide';
});

const props = defineProps<{
  status: FullTranslateFlowStatus;
  progressPercent: string;
  translateLoading: boolean;
}>();

const emit = defineEmits<{ (event: 'translate'): void }>();

const vipStore = useVipStore();
const userStore = useUserStore();

const rightInfo = ref<GetFullTextTranslateRightInfoResponse>();
const pending = ref(false);
const fetchError = ref<null | ResponseError>(null);

const remainCount = computed(() => rightInfo.value?.remainCount ?? 0);
const hasRight = computed(
  () => remainCount.value === -1 || remainCount.value > 0
);

const docStore = useDocStore();
const { docInfo } = storeToRefs(docStore);
const isReadyForTranslation = computed(() => {
  return (docInfo.value as any)?.parsedStatus === UserDocParsedStatusEnum.CONTENT_DATA_PARSED;
});

const initRightInfoPopover = async () => {
  
  pending.value = true;
  fetchError.value = null;
  try {
    const res = await getRightInfo();
    rightInfo.value = res;
  } catch (error) {
    fetchError.value = error as ResponseError;
  }
  pending.value = false;
};
initRightInfoPopover();

const formatedWaiting = computed(() => {
  if (!props.progressPercent) {
    return '0%';
  }

  return Math.round(parseFloat(props.progressPercent)) + '%';
});

const fullTextTranslateStore = useFullTextTranslateStore();

const handleScrollMode = (checked: boolean) => {
  fullTextTranslateStore.setScrollMode(
    checked ? PDFWebviewScrollMode.lock : PDFWebviewScrollMode.unlock
  );
};

const handlePreviewMode = (checked: boolean) => {
  fullTextTranslateStore.setPreviewMode(
    checked
      ? PDFWebviewPreviewMode.withOriginalPDF
      : PDFWebviewPreviewMode.onlyTranslatePDF
  );
  // 显示原文按钮点击上报
  reportElementClick({
    page_type: PageType.note,
    type_parameter: currentNoteInfo.value?.pdfId,
    element_name: ElementClick.original_text,
    status: checked ? 'on' : 'off',
  });
};

// 定义一个变量来跟踪下载状态
let isDownloading = false;

const exportTranslatedPDF = async () => {
  // 如果已经在下载，则直接返回
  if (isDownloading) {
    console.warn('下载已在进行中，请稍候...');
    return;
  }

  // 设置标志位为 true，表示正在下载
  isDownloading = true;

  const fullTextTranslateStore = useFullTextTranslateStore();

  // 每次下载前重新获取最新的翻译文件URL（链接可能过期）
  let fileUrl = '';
  try {
    const pdfId = currentNoteInfo.value?.pdfId || '';
    const pdfUrl = currentNoteInfo.value?.pdfUrl || '';

    if (!pdfId || !pdfUrl) {
      console.error('无法获取 pdfId 或 pdfUrl，无法刷新下载链接');
    } else {
      const latest = await getTranslateStatus({
        pdfId,
        needTranslateFileUrl: pdfUrl,
      });
      if (latest?.translationFileUrl) {
        fullTextTranslateStore.setFullTextTranslatePDFUrl(latest.translationFileUrl);
      }
    }
  } catch (e) {
    console.error('刷新翻译文件链接失败：', e);
  }

  // 如果调用失败或未返回，回退使用现有的 store 值
  fileUrl = fullTextTranslateStore.fullTextTranslatePDFUrl;

  if (!fileUrl) {
    console.error('未找到文件 URL！');
    isDownloading = false; // 重置标志位
    return;
  }

  try {
    // 创建一个 <a> 标签
    const link = document.createElement('a');
    link.target = '_blank';
    link.href = fileUrl;

    // 设置下载属性，指定下载的文件名
    link.download = 'translated.pdf';

    // 触发点击事件来启动下载
    link.click();

    // 下载完成，移除 <a> 标签
    link.remove();
  } catch (error) {
    console.error('下载过程中出现错误：', error);
  } finally {
    // 重置标志位，以允许下次下载
    isDownloading = false;
  }
};

const handleTranslate = async () => {
  emit('translate');
};

const fulltextPopoverRef = ref<HTMLDivElement>();

onMounted(() => {
  fulltextPopoverRef.value?.addEventListener('click', (e) => {
    const anchor = (e.target as HTMLElement).closest('.js-go-invite-page');
    if (anchor) {
      const path = anchor.getAttribute('data-url') || '/invite';
      const url = isInElectron() ? `https://${getHostname()}${path}` : path;
      window.open(url);
    }
  });
});

watch([rightInfo, vipStore], () => {
  if (!rightInfo.value) {
    return;
  }

  if (!hasRight.value) {
    return;
  }

  if (!vipStore.seniorRole) {
    reportTranslateExposure(true);
    return;
  }

  if (props.status === FullTranslateFlowStatus.TRANSLATING) {
    return;
  }

  reportTranslateExposure(false);
});
</script>

<style lang="less" scoped>
.fulltext-translate-popover {
  position: absolute;
  background-color: #fff;
  padding: 20px 24px;
  color: #1d2229;
  border-radius: 4px;
  box-shadow:
    0px 3px 6px rgba(0, 0, 0, 0.12),
    0px 6px 16px rgba(0, 0, 0, 0.08),
    0px 9px 28px rgba(0, 0, 0, 0.05);

  .switcher {
    width: 264px;
    display: flex;
    justify-content: space-between;
  }

  .switcher + .switcher {
    margin-top: 8px;
  }

  .content {
    min-height: 196px;
    width: 420px;

    & > .error {
      text-align: center;
      padding-top: 60px;
      cursor: pointer;
    }
  }

  .title {
    font-size: 18px;
    line-height: 28px;
    margin-bottom: 16px;
    a {
      text-decoration: underline;
    }
  }
  .additional-rules {
        margin-top: 16px;
        
        .rule-item {
          line-height: 24px;
          margin: 6px 0;
          color: #666;
          
          .rule-link {
            color: #1f71e0;
            text-decoration: underline;
            
            &:hover {
              color: #0d5aa7;
            }
          }
        }
    }

  .rights {
    .rules {
      :deep(p) {
        line-height: 24px;
        margin: 0;
      }

      :deep(.label) {
        color: #1f71e0;
        border: 1px solid #1f71e0;
        border-radius: 2px;
        padding: 2px 8px;
        margin: 0 2px;
        cursor: pointer;
      }
    }
  }

  .bottom {
    margin-top: 24px;
    display: flex;
    justify-content: space-between;

    .doing {
      color: #3da611;
    }

    .error {
      color: #e66045;
    }

    .translate-action-wrapper {
      display: flex;
      flex-direction: column;
      align-items: flex-end;
      gap: 8px;
    }

    .unparsed-tip {
      color: #8a8a8a; /* Gray color */
      text-align: right;
      font-size: 12px;
      margin-top: 4px;
    }
  }
}

:deep(.ant-spin-container)::after {
  background: transparent;
}
</style>
