<template>
  <div class="reference-info">
    <div
      class="title"
      @click="handleTitleClick"
    >
      {{ paperData.title }}
    </div>
    <div class="authors info">
      <div
        v-if="authorList.length"
        ref="authorsRef"
        class="authors"
        style="position: absolute; visibility: hidden"
      >
        <span
          v-for="(author, index) in calcAuthorsList"
          :key="`pre-${index}`"
          class="author"
        >
          <span>{{ author.name }}</span>
        </span>
      </div>
      <div
        v-if="authorList.length"
        class="authors block"
      >
        <span
          v-for="(author, index) in authorList[0]"
          :key="`pre-${index}`"
          class="author"
        >
          <span>{{ author.name }}</span>
        </span>
        <span
          v-for="(author, index) in authorList[1] && allAuthors
            ? authorList[1]
            : []"
          :key="`mid-${index}`"
          class="author"
        >
          <span>{{ author.name }}</span>
        </span>
        <span
          v-if="authorList[1]"
          class="num"
          @click="showMore"
        >
          {{ !allAuthors ? '...+' + authorList[1].length : '收起' }}
        </span>
        <span
          v-if="authorList[2]"
          class="last"
        >
          <span>{{ authorList[2][0].name }}</span>
        </span>
      </div>
    </div>

    <div class="venues-date">
      <span
        v-for="(item, index) in venues"
        :key="index"
        class="block"
      >
        {{ item }}
      </span>

      <span
        v-if="paperData.publishDate"
        class="block"
      >
        {{ paperData.publishDate }}
      </span>
    </div>

    <PerfectScrollbar
      class="summary"
      :options="{
        suppressScrollX: true,
      }"
    >
      <div
        v-if="summary"
        class="js-interact-drag-ignore"
      >
        {{ summary }}
        <span
          v-if="paperData.originalAbstract?.length > 300"
          class="btn"
        >
          <ellipsis-outlined
            v-show="isToggle"
            class="aiknowledge-icon icongengduo"
          />
          <i
            :class="[
              'aiknowledge-icon',
              isToggle ? 'icon-arrow-down' : 'icon-arrow-up',
            ]"
            aria-hidden="true"
            @click="handleToggle"
          />
        </span>
      </div>
    </PerfectScrollbar>

    <div class="tag-wrap">
      <span
        v-if="paperData.pdfId && paperData.pdfId !== '0'"
        class="tag"
      >{{
        $t('viewer.hasPDF')
      }}</span>

      <span v-if="paperData.venueTags && paperData.venueTags.length > 0">
        <span
          v-for="item in paperData.venueTags"
          :key="item"
          class="venue-tagging-tag"
        >
          {{ item }}
        </span>
      </span>
    </div>
    <div class="bottom">
      <div class="text-xs">
        <span v-if="paperData.noteCount > 10">{{ paperData.noteCount || 0 }} {{ $t('viewer.noteCount') }} ·
        </span>
        <span>{{ $n(paperData.citationCount || 0, 'integer') }}
          {{ $t('viewer.citations') }}</span>
      </div>
      <div class="flex">
        <div
          v-if="showCiteBtn"
          class="cursor-pointer w-13 h-8 text-rp-blue-6 text-sm font-medium flex items-center"
          @click="citeDialogVisible = true"
        >
          <i
            class="aiknowledge-icon icon-cite mr-1 text-base"
            aria-hidden="true"
          />
          {{ $t('viewer.quoteLabel') }}
        </div>
        <template v-if="!noCollect">
          <div
            v-if="!paperData.isCollected"
            class="collect"
            @click="collectPaper"
          >
            <a-spin
              :spinning="collectLoading"
              size="small"
            >
              <div style="display: flex; align-items: center">
                <i
                  class="aiknowledge-icon icon-bookmark"
                  aria-hidden="true"
                />
                {{ $t('viewer.savePaper') }}
              </div>
            </a-spin>
          </div>
          <FolderMenu
            v-else
            :doc-id="docId || paperData.userDocId"
            :get-popup-container="getPopupContainer"
          >
            <div
              class="collect collected-btn"
              @click="cancelCollectPaper"
            >
              <a-spin
                :spinning="collectLoading"
                size="small"
              >
                <div style="display: flex; align-items: center">
                  <i
                    class="aiknowledge-icon icon-bookmark"
                    aria-hidden="true"
                  />
                  {{ $t('viewer.saved') }}
                </div>
              </a-spin>
            </div>
          </FolderMenu>
        </template>
      </div>
    </div>
    <CiteModal
      v-if="citeDialogVisible"
      v-model:citeDialogVisible="citeDialogVisible"
      :paperData="paperData"
      @update:success="onUpdateSuccess"
    />
  </div>
</template>
<script setup lang="ts">
import { PaperDetailInfo } from 'go-sea-proto/gen/ts/paper/Paper'
import { ref, computed, onMounted, watch } from 'vue';
import { getDomainOrigin } from '~/src/util/env';
import FolderMenu from '@/components/Folder/FolderMenu.vue';
import { goPathPage } from '~/src/common/src/utils/url';
import { EllipsisOutlined } from '@ant-design/icons-vue';
import { useCollect } from '~/src/hooks/useCollect';
import {
  reportPaperItemClick,
  PageType,
  getPdfIdFromUrl,
  reportPopupNoteCiteFigureImpression,
} from '~/src/api/report';
import { ReferenceMarker } from 'go-sea-proto/gen/ts/pdf/PdfParse';

import { selfNoteInfo } from '@/store';
import CiteModal from './citeModal.vue';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { reportElementClick, ElementClick } from '@/api/report';

const props = defineProps<{
  paperData: PaperDetailInfo;
  showCiteBtn?: boolean;
  noCollect?: boolean;
  marker?: ReferenceMarker;
  paperDataReload: () => void | Promise<any>;
}>();

const emit = defineEmits<{
  (event: 'hide'): void;
}>();

const annotationStore = useAnnotationStore();
const citeDialogVisible = ref(false);
const isToggle = ref(true);
const handleToggle = () => {
  isToggle.value = !isToggle.value;
};

const onUpdateSuccess = async () => {
  await props.paperDataReload();
};

watch(
  () => citeDialogVisible.value,
  (value) => {
    annotationStore.showReferenceTippy = value;
  }
);

const summary = computed(() => {
  const originalAbstract = props.paperData.originalAbstract;
  if (!originalAbstract) {
    return '';
  }
  if (isToggle.value) {
    return originalAbstract.slice(0, 300);
  }
  return originalAbstract;
});

const handleTitleClick = () => {
  reportPaperItemClick({
    page_type: PageType.note,
    type_parameter: selfNoteInfo.value?.pdfId ?? getPdfIdFromUrl(),
    module_type: 'cite_popup_paper',
    paper_id: props.paperData.paperId,
    scene_id: '',
    ref_content: props?.marker?.refContent || '',
    subject_id: '',
  });

  goPathPage(`${getDomainOrigin()}/paper/${props.paperData.paperId}`);
};

const { collectLoading, docId, collectPaper, cancelCollectPaper } = useCollect(
  () => props.paperData.paperId,
  () => props.paperData.title,
  (collected) => {
    props.paperData.isCollected = collected;
    reportElementClick({
      page_type: PageType.note,
      element_name: ElementClick.save_reference_paper,
      type_parameter: 'none',
      status: collected ? 'on' : 'off',
    });
  },
  false
);

const showAllAuthors = ref<boolean>(false);

const authorList = computed(() => {
  if (!props.paperData.authorList || !props.paperData.authorList.length) {
    return [];
  }
  if (props.paperData.authorList.length <= 4 || showAllAuthors.value) {
    return [props.paperData.authorList];
  }
  const list = props.paperData.authorList;
  const num = 2;
  return [
    list.slice(0, num),
    list.slice(num, list.length - 1),
    list.slice(list.length - 1),
  ];
});

const allAuthors = ref<boolean>(false);
const showMore = () => {
  allAuthors.value = !allAuthors.value;
};

const calcAuthorsList = computed(
  () => props.paperData.authorList?.slice(0, 15) || []
);

const authorsRef = ref<HTMLDivElement>();

const calcShowAllAuthors = () => {
  if (authorsRef.value && props.paperData.authorList?.length > 4) {
    const parentNode = authorsRef.value.parentElement;
    if (!parentNode) {
      return;
    }
    const children = authorsRef.value.children;
    let i = 0;
    if (children && children.length) {
      const maxWidth = parentNode.offsetWidth;
      let width = 0;
      for (; i < children.length - 1; i += 1) {
        width += (children[i] as HTMLElement).offsetWidth;
        if (width >= maxWidth) {
          break;
        }
      }
      if (props.paperData.authorList?.length <= i + 1) {
        showAllAuthors.value = true;
      }
    }
  }
};

onMounted(() => {
  reportPopupNoteCiteFigureImpression({
    page_type: PageType.note,
    type_parameter: selfNoteInfo.value?.pdfId ?? getPdfIdFromUrl(),
    popup_type: 'paper',
    popup_id: props.paperData.paperId,
    order_num: props?.marker?.refContent || '',
  });

  calcShowAllAuthors();
});

const venues = computed(() => {
  if (!props.paperData.venues || !props.paperData.venues.length) {
    return [];
  }

  return props.paperData.venues.filter((item) => item != null);
});

const getPopupContainer = () =>
  document.getElementById('peak') as HTMLDivElement;
</script>

<style scoped lang="less">
.reference-info {
  padding: 20px 24px 0;
  .title {
    line-height: 22px;
    cursor: pointer;
    font-size: 14px;
    font-weight: bold;
    color: #1f71e0;
    font-family: Roboto-Bold, Roboto;

    display: -webkit-box;
    -webkit-line-clamp: 5;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
  }

  .info {
    overflow: hidden;
    text-overflow: ellipsis;
    margin-top: 5px;
    font-size: 12px;
    font-weight: 400;
    color: #73716f;
    line-height: 18px;
  }

  .authors {
    .author + .author::before {
      content: '/';
      padding: 0 8px;
    }

    .num {
      cursor: pointer;
      font-size: 13px;
      color: #1f71e0;
      padding: 0 8px;
    }
    .num::after {
      content: '/';
      padding-left: 8px;
    }
    .num::before {
      content: '/';
      padding-right: 8px;
    }
  }
  .summary {
    cursor: text;
    margin-top: 17px;
    font-size: 13px;
    font-weight: 400;
    color: #262625;
    line-height: 20px;
    max-height: 150px;
    overflow-y: auto;
    .btn {
      display: inline-flex;
      align-items: center;

      .aiknowledge-icon {
        font-size: 14px;
        color: rgba(0, 0, 0, 64%);
      }

      .icongengduo {
        margin-right: 3px;
      }

      .icon-arrow-down,
      .icon-arrow-up {
        cursor: pointer;
      }
    }
  }

  .bottom {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 48px;
    background: #ffffff;
    border-radius: 4px;
    border-top: 1px solid #e4e7ed;
    margin-top: 16px;
    color: #86919c;
    font-family: Noto Sans SC;
    .collect {
      width: 94px;
      height: 32px;
      background: #1f71e0;
      border-radius: 2px;

      display: flex;
      align-items: center;
      justify-content: center;

      font-size: 13px;
      font-weight: 400;
      color: #ffffff;
      line-height: 18px;

      margin-left: 16px;
      cursor: pointer;

      .icon-bookmark {
        margin-right: 5px;
        font-size: 16px;
      }
    }

    .collected-btn {
      background: transparent;
      font-weight: 600;
      color: #919fb5;

      .icon-bookmark {
        color: #919fb5;
      }
    }
  }

  .venues-date {
    font-size: 12px;
    font-weight: 400;
    color: #73716f;
    line-height: 18px;
    word-break: break-word;
    .block::before {
      content: '\00b7';
      padding: 0 6px;
      color: #999;
    }
    .block:first-child {
      &::before {
        content: '';
        margin-left: -12px;
      }
    }
  }

  .tag-wrap {
    font-weight: 400;
    font-size: 12px;
    line-height: 18px;
    margin-top: 12px;
    .tag {
      display: inline-block;
      color: #62ac06;
      padding: 1px 6px;
      background: rgba(98, 172, 6, 0.1);
      border-radius: 11px;
      margin-right: 6px;
    }
    .venue-tagging-tag {
      display: inline-block;
      color: #4e5969;
      padding: 1px 6px;
      background: #f0f2f5;
      border-radius: 11px;
      margin-right: 6px;
    }
  }
}

html[data-theme='dark'] {
  .reference-info {
    .tag-wrap {
      .tag {
        background-color: #e9f7c6;
      }
    }
  }
}
</style>
