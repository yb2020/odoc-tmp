<template>
  <Intersect @enter="handleEnter(paperData.id)">
    <li :class="['paper-item', { 'paper-item-small': size === 'small' }]">
      <h3
        class="title rp-pointer"
        @click="goPaperPage(paperData.id)"
        v-html="paperData.title"
      ></h3>
      <a
        class="info rp-pointer"
        :href="`${getPaperDetailLink(paperData.id)}`"
        target="_blank"
        @click="goPaperPage(paperData.id, $event)"
      >
        <div
          v-if="!isShowUserDoc"
          class="rp-line-ellipsis authors"
          style="display: flex"
        >
          <div
            v-if="authorNames && authorNames.length > 0"
            :class="{
              block: (venues && venues.length > 0) || journalOrConference,
            }"
          >
            <span
              v-for="(item, index) in authorNames"
              :key="`${index}`"
              class="author"
              v-html="item.type === 'author' ? item.value.name : '...'"
            >
            </span>
          </div>

          <div
            v-if="venues && venues.length > 0"
            ref="journalOrConferenceRef"
            class="venues-wrap"
          >
            <span v-for="(venue, index) in venues" :key="index" class="block">
              {{ venue }}
            </span>
          </div>

          <span
            v-else-if="journalOrConference"
            ref="journalOrConferenceRef"
            v-html="xss(journalOrConference)"
          >
          </span>

          <span
            v-if="publishDate"
            :class="[
              'date',
              {
                'time-block':
                  (venues && venues.length > 0) ||
                  journalOrConference ||
                  (authorNames && authorNames.length > 0),
              },
            ]"
          >
            {{ publishDate }}</span
          >
        </div>

        <div v-else class="rp-line-ellipsis authors">
          <span
            v-for="(item, index) in showAuthorNames"
            :key="`${index}`"
            class="author"
            v-html="item.type === 'author' ? item.value.name : '...'"
          >
          </span>

          <div style="display: flex">
            <div
              v-if="venues && venues.length > 0"
              ref="journalOrConferenceRef"
              class="venues-wrap"
            >
              <span v-for="(venue, index) in venues" :key="index" class="block">
                {{ venue }}
              </span>
            </div>

            <span
              v-if="paperData && paperData.publishDate"
              :class="['date', { 'time-block': venues && venues.length > 0 }]"
            >
              {{ paperData.publishDate }}
            </span>
          </div>
        </div>

        <p
          v-if="paperData.summary"
          class="summary"
          v-html="xss(paperData.summary)"
        ></p>
      </a>

      <slot name="other" />

      <div class="extra">
        <div class="left">
          <span
            v-if="
              paperData.pdfId &&
              paperData.pdfId !== '0' &&
              pageType === PageType.newPaper
            "
            class="new-tag"
            @click.stop="gotoPdfPage(paperData.pdfId)"
            ><a-icon type="file-pdf" theme="filled" class="file-pdf-icon" />{{
              $t('home.searchPage.paperItem.readPDF')
            }}</span
          >
          <span
            v-else-if="
              paperData.pdfId && paperData.pdfId !== '0' && isShowHasPdfTag
            "
            :class="['tag', { 'rp-pointer': pageType === PageType.newPaper }]"
            @click.stop="gotoPdfPage(paperData.pdfId)"
            >{{ $t('home.searchPage.paperItem.hasPdf') }}</span
          >

          <span
            v-if="
              paperData.venueTags &&
              paperData.venueTags.length > 0 &&
              isShowVenueTagging
            "
          >
            <span
              v-for="item in paperData.venueTags"
              :key="item"
              class="venue-tagging-tag"
            >
              {{ item }}
            </span>
          </span>

          <div
            v-if="
              paperData.recommendItem &&
              paperData.recommendItem.reason &&
              paperData.recommendItem.reason !== RecommendReason.hot
            "
          >
            <img
              v-if="!paperData.recommendItem.qtitle"
              :src="
                recommendSrc[paperData.recommendItem.reason][isWebEN ? 1 : 0]
              "
              alt=""
              class="recommend-img"
            />
            <a-popover
              v-else
              placement="bottomLeft"
              overlay-class-name="recommend-popover"
              :get-popup-container="(triggerNode: HTMLDivElement) => triggerNode.parentNode"
            >
              <template #content>
                <i
                  class="aiknowledge-icon icon-ai"
                  aria-hidden="true"
                  style="margin-right: 8px; font-size: 16px"
                ></i>
                {{ $t('home.searchPage.paperItem.recommend1') }}
                <span
                  class="recommend-title rp-pointer"
                  :style="{ 'max-width': isMobile ? '140px' : '300px' }"
                  @click.stop="
                    goPathPage(`/paper/${paperData.recommendItem.qid}`)
                  "
                  >{{ paperData.recommendItem.qtitle }} </span
                >{{ $t('home.searchPage.paperItem.recommend2') }}
              </template>
              <img
                :src="
                  recommendSrc[paperData.recommendItem.reason][isWebEN ? 1 : 0]
                "
                alt=""
                class="recommend-img rp-pointer"
              />
            </a-popover>
          </div>
          <div class="leftContent">
            <slot name="left"></slot>
          </div>
        </div>
        <div class="right">
          <div class="tip">
            <slot name="extra" />
          </div>
          <div v-if="isCollectPaper" class="rp-pointer collect">
            <CollectDialog
              :paper-id="paperData.id"
              :is-collected="isCollected"
              :collect-limit-dialog-report-params="
                collectLimitDialogReportParams
              "
              @loadingChange="handleCollect"
            >
              <a-icon v-if="isCollectLoading" type="loading" />
              {{
                isCollected
                  ? $t('home.searchPage.paperItem.collected')
                  : $t('home.searchPage.paperItem.collect')
              }}<i class="aiknowledge-icon icon-bookmark" aria-hidden="true" />
            </CollectDialog>
          </div>
        </div>
      </div>
    </li>
  </Intersect>
</template>
<script lang="ts">
import { defineComponent, PropType, ref, computed, onMounted } from 'vue'
import moment from 'moment'
import Intersect from 'vue-intersect'
import xss from 'xss'
import { $PaperAuthor, $PaperDetail, RecommendReason } from '@/common/src/api/paper'
import CollectDialog from './CollectDialog.vue'
import useBehaviorReport, {
  ReportActionType,
  ReportItemType,
  ReportSceneId,
} from '@/hooks/useBehaviorReport'
import { goPathPage, goPdfPage } from '@/common/src/utils/url'
import recommendRef from '@/assets/images/paper/recommend-ref.png'
import recommendClassic from '@/assets/images/paper/recommend-classic.png'
import recommendNew from '@/assets/images/paper/recommend-new.png'
import recommendRefEn from '@/assets/images/paper/recommend-ref-en.png'
import recommendClassicEn from '@/assets/images/paper/recommend-classic-en.png'
import recommendNewEn from '@/assets/images/paper/recommend-new-en.png'
import {
  EventCode,
  PageType,
  reportPaperItem,
  reportElementClick,
} from '@/utils/report'

import { useUserStore } from '@/common/src/stores/user'
import { useLanguage } from '@/hooks/useLanguage'

export const LOCALSTORAGE_CURRENT_SEARCH_ID = 'currentSearchId'

export default defineComponent({
  components: { CollectDialog, Intersect },
  props: {
    paperData: {
      type: Object as PropType<$PaperDetail>,
      default: () => ({}),
    },
    sceneId: {
      type: String as PropType<ReportSceneId>,
      default: ReportSceneId.unknown,
    },
    isCollectPaper: {
      type: Boolean,
      default: true,
    },
    size: {
      type: String as PropType<'default' | 'small'>,
      default: 'default',
    },
    keyword: {
      type: String,
      default: '',
    },
    index: {
      type: Number,
      default: -1,
    },
    isShowUserDoc: {
      type: Boolean,
      default: false,
    },
    pageType: {
      type: String, // as () => PageType,
      default: '',
    },
    moduleType: {
      type: String,
      default: '',
    },
    isShowVenueTagging: {
      type: Boolean,
      default: true,
    },
    collectLimitDialogReportParams: {
      type: Object, // as PropType<LimitDialogReportParams>,
      default: null,
    },
    isShowHasPdfTag: {
      type: Boolean,
      default: true,
    },
    authorLimitNumber: {
      type: Number,
      default: 4,
    },
    typeParameter: {
      type: String,
      default: 'none',
    },
    hasPdfElementName: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const userStore = useUserStore()
    const authorNames = computed(() => {
      const list = props.paperData.authorList

      if (!list || !list.length) {
        return []
      }
      if (list.length <= 4) {
        return list.map((item) => {
          return {
            type: 'author',
            value: item,
          }
        })
      }
      const first = list.findIndex((item) => /<em>/.test(item.name))
      const second =
        first === -1
          ? -1
          : list
              .slice(first + 1)
              .findIndex((item) => /<em>/.test(item.name)) === -1
          ? list.slice(first + 1).findIndex((item) => /<em>/.test(item.name)) +
            first +
            2
          : list.slice(first + 1).findIndex((item) => /<em>/.test(item.name)) +
            first +
            1

      const third =
        second === -1
          ? -1
          : list
              .slice(second + 1)
              .findIndex((item) => /<em>/.test(item.name)) === -1
          ? list.slice(second + 1).findIndex((item) => /<em>/.test(item.name)) +
            second +
            2
          : list.slice(second + 1).findIndex((item) => /<em>/.test(item.name)) +
            second +
            1

      let result: (
        | { type: 'author'; value: $PaperAuthor }
        | { type: 'ellipsis' }
      )[] = []
      if (third <= 2) {
        result = list.slice(0, 3).map((item) => {
          return {
            type: 'author',
            value: item,
          }
        })
        result.push({
          type: 'ellipsis',
        })
        result.push({
          type: 'author',
          value: list[list.length - 1],
        })
        return result
      }

      if (first > 0) {
        result.push({
          type: 'ellipsis',
        })
      }
      result.push({
        type: 'author',
        value: list[first],
      })
      if (second > first + 1 && second <= list.length - 1) {
        result.push({
          type: 'ellipsis',
        })
      }

      if (second <= list.length - 1) {
        result.push({
          type: 'author',
          value: list[second],
        })
      }

      if (third > second + 1 && third <= list.length - 1) {
        result.push({
          type: 'ellipsis',
        })
      }

      if (third <= list.length - 1) {
        result.push({
          type: 'author',
          value: list[third],
        })
      }

      if (third < list.length - 1) {
        result.push({
          type: 'ellipsis',
        })
      }

      return result
    })

    const otherInfo = {
      query: props.keyword,
      search_id: localStorage.getItem(LOCALSTORAGE_CURRENT_SEARCH_ID),
      rank: props.index,
    }

    const reportHomeOrSearchPaperItem = (
      event:
        | EventCode.readpaperPaperItemClick
        | EventCode.readpaperPaperItemImpression,
      id: string
    ) => {
      if (!props.pageType || !props.moduleType) {
        return
      }

      let sceneId = ''
      if ([PageType.library, PageType.subjectPage].includes(props.pageType)) {
        sceneId = props.paperData.recommendId || ''
      } else if (props.pageType === PageType.search) {
        sceneId = localStorage.getItem(LOCALSTORAGE_CURRENT_SEARCH_ID) || ''
      }

      return reportPaperItem({
        event_code: event,
        type_parameter: props.typeParameter,
        page_type: props.pageType,
        module_type: props.moduleType,
        paper_id: id,
        scene_id: sceneId,
        order_num: props.index,
        subject_id: 'none',
      })
    }

    const reportUserBehavior = useBehaviorReport()
    const goPaperPage = (id: string, e?: Event) => {
      if (e) e.preventDefault()

      if (!id || id === '0') return

      emit('reportFeedClick')

      reportHomeOrSearchPaperItem(
        EventCode.readpaperPaperItemClick,
        props.paperData.id
      )

      reportUserBehavior({
        itemId: props.paperData.id,
        itemType: ReportItemType.paper,
        actionType: ReportActionType.click,
        sceneId: props.sceneId,
        otherInfo: JSON.stringify(otherInfo),
      })
      goPathPage(`/paper/${id}`)
    }

    let hasReport = false

    const handleEnter = (paperId: string) => {
      if (props.sceneId === ReportSceneId.unknown || hasReport) {
        return
      }

      hasReport = true
      reportHomeOrSearchPaperItem(
        EventCode.readpaperPaperItemImpression,
        paperId
      )
      reportUserBehavior(
        {
          itemId: paperId,
          itemType: ReportItemType.paper,
          actionType: ReportActionType.show,
          sceneId: props.sceneId,
          otherInfo: JSON.stringify(otherInfo),
        },
        undefined,
        1500
      )
    }

    const isCollectLoading = ref<boolean>(false)

    const isCollected = ref<boolean>(props.paperData.isCollected)

    const handleCollect = (isLoading: boolean, result: boolean) => {
      isCollectLoading.value = isLoading
      if (!isLoading) {
        isCollected.value = result
      }
    }

    const publishDate = computed<string>(() => {
      if (props.paperData.publishDate) {
        const publishDate = moment(props.paperData.publishDate - 0).locale('en')
        const time = publishDate.format('MMM YYYY')
        return time || ''
      }
      return ''
    })

    const journalOrConference = computed<string>(() => {
      return props.paperData.journal || props.paperData.conferenceInfo || ''
    })

    const journalOrConferenceRef = ref()

    onMounted(() => {
      if (journalOrConferenceRef.value) {
        const dom: any = journalOrConferenceRef.value
        const date = dom.nextElementSibling
        if (!date || !dom.offsetParent) {
          return
        }
        if (
          dom.offsetLeft + dom.offsetWidth >
          dom.offsetParent.offsetWidth - date.offsetWidth
        ) {
          date.className = 'date fixed'
        }
      }
    })

    const showAuthorNames = computed(() => {
      const list = props.paperData.authorList

      if (!list || !list.length) {
        return []
      }

      let result: (
        | { type: 'author'; value: $PaperAuthor }
        | { type: 'ellipsis' }
      )[] = []

      if (list.length <= props.authorLimitNumber)
        return list.map((item) => {
          return {
            type: 'author',
            value: item,
          }
        })
      else {
        result = list.slice(0, props.authorLimitNumber - 1).map((item) => {
          return {
            type: 'author',
            value: item,
          }
        })
        result.push(
          { type: 'ellipsis' },
          { type: 'author', value: list[list.length - 1] }
        )

        return result
      }
    })

    const venues = computed(() => {
      if (!props.paperData.venues || !props.paperData.venues.length) {
        return []
      }

      return props.paperData.venues.filter((item) => item != null)
    })

    const isMobile = false // TODO

    const recommendSrc = {
      [RecommendReason.ref]: [recommendRef, recommendRefEn],
      [RecommendReason.classic]: [recommendClassic, recommendClassicEn],
      [RecommendReason.new]: [recommendNew, recommendNewEn],
    }

    const gotoPdfPage = (pdfId: string) => {
      if (props.pageType !== PageType.newPaper) return

      reportElementClick({
        page_type: props.pageType,
        type_parameter: props.typeParameter,
        element_name: props.hasPdfElementName,
        status: 'none',
      })

      if (!userStore.isLogin()) {
        return
      }
      goPdfPage({ pdfId })
    }

    const { isEnUS } = useLanguage()
    const isWebEN = isEnUS // 保持向后兼容的命名
    const getPaperDetailLink = (id: string) => {
      return `${window.location.origin}/paper/${id}`
    }

    return {
      authorNames,
      goPaperPage,
      handleEnter,
      publishDate,
      handleCollect,
      isCollectLoading,
      isCollected,
      journalOrConference,
      journalOrConferenceRef,
      xss,
      showAuthorNames,
      venues,
      isMobile,
      recommendSrc,
      goPathPage,
      RecommendReason,
      getPaperDetailLink,
      gotoPdfPage,
      PageType,
      isWebEN,
    }
  },
})
</script>
<style lang="less" scoped>
@import './semantic.less';
.paper-item {
  /deep/em {
    font-style: normal;
    color: #1f71e0;
  }
  .title {
    font-size: 20px;
    // stylelint-disable-next-line
    font-family: Lato-Black, Lato;
    line-height: 1.4;
    font-weight: 900;
    color: #262625;
    margin-bottom: 10px;
  }
  .info {
    color: rgba(0, 0, 0, 0.45);
    font-size: 13px;
    line-height: 18px;
    margin-bottom: 16px;
    display: block;
  }
  .authors {
    font-size: 13px;
    font-weight: 400;
    color: #62738c;
    line-height: 18px;
    position: relative;
    .author + .author::before {
      content: '/';
      padding: 0 8px;
      opacity: 0.5;
    }
    .num::before {
      content: '/';
      padding: 0 8px;
      opacity: 0.5;
    }
    .last {
      padding-left: 8px;
    }
    .block {
      &::after {
        content: '\00b7';
        padding: 0 6px;
        color: #999;
      }
    }
    .venues-wrap {
      display: flex;
      .block:last-child {
        &::after {
          content: '';
          margin-left: -12px;
        }
      }
    }
    .time-block {
      &::before {
        content: '\00b7';
        padding: 0 6px;
        color: #999;
      }
    }
    .date {
      &.fixed {
        position: absolute;
        right: 0;
        background: #fff;
        &::before {
          content: '...';
          padding-right: 4px;
        }
      }
    }
  }
  .summary {
    display: -webkit-box;
    overflow: hidden;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    font-size: 15px;
    // stylelint-disable-next-line
    font-family: Lato-Regular, Lato;
    font-weight: 400;
    color: #262625;
    line-height: 22px;
    margin-top: 5px;
  }

  & + & {
    margin-top: 32px;
  }

  .extra {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    font-size: 12px;
    font-weight: 400;
    line-height: 18px;
    .left {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      .leftContent {
        font-size: 13px;
        font-weight: 400;
        color: #73716f;
        line-height: 22px;
        white-space: nowrap;
      }
    }
    .tag {
      color: #62ac06;
      padding: 1px 6px;
      background: rgba(98, 172, 6, 0.1);
      border-radius: 11px;
      margin-right: 6px;
    }
    .new-tag {
      background: #1f71e0;
      color: #fff;
      padding: 1px 6px;
      border-radius: 11px;
      margin-right: 6px;
      cursor: pointer;
      .file-pdf-icon {
        margin-right: 4px;
      }
    }
    .venue-tagging-tag {
      display: inline-block;
      color: #4e5969;
      padding: 1px 6px;
      background: #f0f2f5;
      border-radius: 11px;
      margin-right: 6px;
    }
    .right {
      display: flex;
      justify-content: flex-end;
      align-items: center;
      color: #86919c;
      font-size: 13px;
      line-height: 20px;
      .collect {
        font-size: 14px;
        font-weight: 400;
        color: #86919c;
        line-height: 24px;
        flex-shrink: 0;
        display: flex;
        margin-left: 16px;
        .icon-bookmark {
          margin-left: 7px;
          font-size: 13px;
          height: 13px;
          line-height: 13px;
        }
        .collectionCount {
          margin-right: 7px;
        }
      }
    }
    .recommend-img {
      height: 20px;
    }
  }

  &-small {
    .title {
      font-size: 14px;
      line-height: 20px;
      margin-bottom: 0;
    }
    .info {
      margin-bottom: 8px;
    }
  }
}

.mobile-viewport {
  .paper-item {
    .title {
      font-size: 18rpx;
      line-height: 24rpx;
    }
  }
  .paper-item + .paper-item {
    margin-top: 24rpx;
  }
}
</style>
