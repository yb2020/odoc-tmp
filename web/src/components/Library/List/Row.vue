<template>
  <div class="literature-list-item" type="flex" :gutter="16">
    <div v-show="collating" class="literature-list-item-check">
      <Checkbox
        v-if="collating"
        :checked="storeLibraryList.paperListCheckedMap[item.docId]"
        @change="checkItem()"
      />
    </div>
    <div
      v-for="column in storeLibraryList.paperHeadVisibleList"
      :key="column.key"
      :style="storeLibraryList.paperHeadExtra[column.key].style as StyleValue"
      class="literature-list-item-col"
      :class="['literature-list-item-col-' + column.key]"
    >
      <div
        :style="storeLibraryList.paperHeadExtra[column.key].style as StyleValue"
        :class="['literature-list-item-wrap-' + column.key]"
      >
        <div
          v-if="sortToKey[storeLibraryList.currentSortType] === column.key"
          class="sort-overlay"
        ></div>
        <div v-if="column.key === 'docName'" class="title">
          <div v-if="!isEdit" class="edit-filename">
            <ReadStatus
              :pdf-id="item.pdfId"
              :status="item.docReadingStatus"
              :percent="item.progress"
              @updated="(data) => updateReadStatus(item, data)"
            />
            <p class="mb-0 text">
              <Tag
                v-if="item.isLatestRead"
                :class="['tag-last-read']"
                >{{ $t('home.library.lastRead') }}</Tag
              >
              <span
                class="rp-pointer"
                @click="openPage(item)"
                v-html="
                  escapeHightlight(item.searchResult.hitDocName || item.docName)
                "
              >
              </span>
            </p>
            <span style="text-align: center">
              <EditOutlined
                v-if="!isEdit"
                class="iconedit"
                @click="startEditFileName()"
              />
            </span>
          </div>
          <TextArea
            v-else
            v-model:value="editContent"
            auto-size
            @pressEnter="submitFileName()"
            @blur="blurFileName()"
            @keyup.esc="cancelEditFileName()"
          />
          <div
            v-if="item.searchResult.hitNote"
            class="content"
            v-html="escapeHightlight(item.searchResult.hitNote)"
          ></div>
          <Space class="ml-9 mt-2 mb-0">
            <Tag v-if="!hasPDF" class="tag-nopdf mr-0">{{
              $t('home.paper.noPDF')
            }}</Tag>
            <span
              v-if="item.hasAttachment"
              class="tag-attached text-xs text-rp-neutral-6 font-normal"
            >
              <Divider type="vertical" />{{ $t('home.library.attachments') }}
            </span>
          </Space>
        </div>

        <ClassifyCell
          v-else-if="column.key === 'classifyInfos'"
          :classify-infos="item.classifyInfos"
          :doc-id="item.docId"
        />
        <RemarkCell
          v-else-if="column.key === 'remark'"
          :item="item"
          :doc-id="item.docId"
        />
        <Popover
          v-else-if="column.key === 'publishDate'"
          title=""
          trigger="click"
          :visible="storeLibraryList.paperListPublishEdit === item.docId"
          @visibleChange="visibleChangePublishDate"
        >
          <template #content>
            <div class="list-edit-publishDate" @click.stop>
              <div class="list-edit-publishDate-header">
                {{ $t('home.library.publishDate') }}
              </div>
              <div class="list-edit-publishDate-body">
                <InputNumber
                  v-model:value="storeLibraryList.paperListPublishYear"
                  :max="currentYear"
                  size="small"
                  :placeholder="$t('home.global.year')"
                />
                <InputNumber
                  v-model:value="storeLibraryList.paperListPublishMonth"
                  :min="1"
                  :max="12"
                  size="small"
                  :placeholder="$t('home.global.month')"
                />
                <InputNumber
                  v-model:value="storeLibraryList.paperListPublishDate"
                  :min="1"
                  :max="31"
                  size="small"
                  :placeholder="$t('home.global.day')"
                />
              </div>
              <div class="list-edit-publishDate-footer">
                <button @click="cancelPublishDate()">
                  {{ $t('home.global.cancel') }}
                </button>
                <button
                  :disabled="publishDateDisabled"
                  @click="submitPublishDate()"
                >
                  {{ $t('home.global.ok') }}
                </button>
              </div>
            </div>
          </template>
          <div
            class="literature-list-item-publishDate"
            @click="editPublishDate()"
          >
            <Rollback
              v-if="item.displayPublishDate.userEdited"
              :visible="publishDateRollbackVisible"
              :current="item.displayPublishDate.publishDate"
              :origin="item.displayPublishDate.originPublishDate"
              :width="240"
              @cancel="publishDateRollbackVisible = false"
              @ok="rollbackPublishDate()"
              @change="publishDateRollbackVisible = $event"
            />
            <span v-html="item.searchResult.hitPublishDate ? escapeHightlight(item.searchResult.hitPublishDate) : publishDate"></span>
          </div>
        </Popover>
        <div
          v-else-if="column.key === 'authors'"
          class="literature-list-item-authors"
          @click="storeLibraryList.paperListAuthorEdit = item.docId"
        >
          <Rollback
            v-if="item.displayAuthor.userEdited"
            :visible="authorRollbackVisible"
            :current="
              item.displayAuthor.authorInfos.map(getAuthorName).join('、')
            "
            :origin="
              item.displayAuthor.originAuthorInfos.map(getAuthorName).join('、')
            "
            :width="500"
            @cancel="authorRollbackVisible = false"
            @ok="rollbackAuthor()"
            @change="authorRollbackVisible = $event"
          />
          <span>{{ authors }}</span>
        </div>
        <div
          v-else-if="column.key === 'displayVenue'"
          class="literature-list-item-venue"
          :class="{
            'hover-background':
              storeLibraryList.paperListVenueEdit !== item.docId,
          }"
          :style="
            storeLibraryList.paperListVenueEdit !== item.docId
              ? {}
              : {
                  paddingLeft: 0,
                  paddingRight: 0,
                }
          "
          @click="editVenue()"
        >
          <Rollback
            v-if="item.displayVenue.userEdited"
            :visible="venueRollbackVisible"
            :current="item.displayVenue.venueInfos[0] || ''"
            :origin="item.displayVenue.originVenueInfos[0] || ''"
            :width="360"
            @cancel="venueRollbackVisible = false"
            @ok="rollbackVenue()"
            @change="venueRollbackVisible = $event"
          />

          <span v-if="storeLibraryList.paperListVenueEdit !== item.docId" v-html="item.searchResult.hitVenue ? escapeHightlight(item.searchResult.hitVenue) : (item.displayVenue.venueInfos[0] || '')"></span>
          <TextArea
            v-else
            v-model:value.trim="storeLibraryList.paperListVenueEditContent"
            :autosize="{ minRows: 1, maxRows: 5 }"
            :data-list-item-venue="item.docId"
            class="list-thin-scroll list-item-textarea"
            @click.stop
            @blur="submitVenue()"
            @pressEnter="submitVenue()"
          />
        </div>
        <div
          v-else-if="column.key === 'jcrVenuePartion'"
          class="literature-list-item-venue"
          :class="{
            'hover-background':
              storeLibraryList.paperListJcrEdit !== item.docId,
          }"
          :style="
            storeLibraryList.paperListJcrEdit !== item.docId
              ? {}
              : {
                  paddingLeft: 0,
                  paddingRight: 0,
                }
          "
          @click="editJcr()"
        >
          <Rollback
            v-if="item.jcrVenuePartion.userEdited"
            :visible="jcrRollbackVisible"
            :current="item.jcrVenuePartion.jcrVenuePartion || ''"
            :origin="item.jcrVenuePartion.originJcrVenuePartion || ''"
            :width="360"
            @cancel="jcrRollbackVisible = false"
            @ok="rollbackJcr()"
            @change="jcrRollbackVisible = $event"
          />

          <span v-if="storeLibraryList.paperListJcrEdit !== item.docId" v-html="item.searchResult.hitJcrVenuePartion ? escapeHightlight(item.searchResult.hitJcrVenuePartion) : (item.jcrVenuePartion.jcrVenuePartion || '')"></span>
          <TextArea
            v-else
            v-model:value.trim="storeLibraryList.paperListJcrEditContent"
            :autosize="{ minRows: 1, maxRows: 5 }"
            :data-list-item-jcr="item.docId"
            class="list-thin-scroll list-item-textarea"
            @click.stop
            @blur="submitJcr()"
            @pressEnter="submitJcr()"
          />
        </div>
        <div
          v-else-if="column.key === 'impactOfFactor'"
          class="literature-list-item-venue"
          :class="{
            'hover-background':
              storeLibraryList.paperListImpactFactorEdit !== item.docId,
          }"
          :style="
            storeLibraryList.paperListImpactFactorEdit !== item.docId
              ? {}
              : {
                  paddingLeft: 0,
                  paddingRight: 0,
                }
          "
          @click="editImpactFactor()"
        >
          <Rollback
            v-if="item.impactOfFactor.userEdited"
            :visible="impactFactorRollbackVisible"
            :current="item.impactOfFactor.impactOfFactor || ''"
            :origin="item.impactOfFactor.originImpactOfFactor || ''"
            :width="360"
            @cancel="impactFactorRollbackVisible = false"
            @ok="rollbackImpactFactor()"
            @change="impactFactorRollbackVisible = $event"
          />

          <span
            v-if="storeLibraryList.paperListImpactFactorEdit !== item.docId"
            >{{
              typeof item.impactOfFactor.impactOfFactor === 'number'
                ? item.impactOfFactor.impactOfFactor
                : item.impactOfFactor.impactOfFactor || ''
            }}</span
          >
          <TextArea
            v-else
            v-model:value.trim="
              storeLibraryList.paperListImpactFactorEditContent
            "
            :autosize="{ minRows: 1, maxRows: 5 }"
            :data-list-item-impact="item.docId"
            class="list-thin-scroll list-item-textarea"
            @click.stop
            @blur="submitImpactFactor()"
            @pressEnter="submitImpactFactor()"
          />
        </div>
        <div v-else-if="column.key === 'importantanceScore'">
          <Rate
            :value="item.importantanceScore || 0"
            @change="submitScore($event)"
          />
        </div>
        <div v-else-if="column.key === 'parseProgress'" class="parse-progress-cell">
          <span>{{ (item as any).parsedProgress }}</span>
          <ReloadOutlined 
            v-if="shouldShowReparseIcon(item)"
            :class="['reparse-icon', { 'reparse-spinning': isReparsing(item.docId) }]"
            @click.stop="handleReparse(item)"
          />
        </div>
        <!-- 操作列已被注释掉
        <div v-else-if="column.key === 'operation'" style="padding: 0 8px">
          <div style="display: flex; align-items: center; height: 100%">
            <Dropdown
              overlay-class-name="dropdownMenu"
              :get-popup-container="getPopupContainer"
            >
              <a class="ant-dropdown-link" @click="preventDefault">
                <EllipsisOutlined style="font-size: 16px; color: #000" />
              </a>
              <template #overlay>
                <Menu>
                  <Item
                    :class="{
                      'bg-rp-neutral-1 cursor-not-allowed': isNotMatched,
                    }"
                  >
                    <div
                      class="menu-name"
                      :class="{
                        'not-match': isNotMatched,
                      }"
                      @click="!isNotMatched && handlejump('')"
                    >
                      {{ $t('home.library.viewDetail') }}
                      <Tooltip
                        v-if="item.newPaper"
                        placement="top"
                        :title="$t('home.library.paperDetail')"
                      >
                        <i
                          aria-hidden="true"
                          class="aiknowledge-icon icon-file-done"
                        />
                      </Tooltip>
                      <Tooltip
                        v-if="isNotMatched"
                        placement="top"
                        :title="$t('home.library.noStored')"
                      >
                        <InfoCircleOutlined />
                      </Tooltip>
                    </div>
                  </Item>
                  <Item>
                    <span
                      class="menu-name"
                      @click="$emit('attach', item.docId)"
                      >{{ $t('home.library.viewAttach') }}</span
                    >
                  </Item>
                </Menu>
              </template>
            </Dropdown>
          </div>
        </div>
        -->
        <div v-else>
          {{ item[column.key] }}
        </div>
      </div>
    </div>
    <div style="flex: 0 0 14px"></div>
  </div>
</template>
<script lang="ts" setup>
import { StyleValue, computed, ref, nextTick, onUnmounted, inject, Ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import {
  PaperRepositoryStatus,
  UserDocInfo,
  DocReadingStatus,
} from 'go-sea-proto/gen/ts/doc/UserDocManage'
import { UserDocParsedStatusEnum } from 'go-sea-proto/gen/ts/doc/UserDocParsedStatus'
import { calculateParsedProgress } from '~/src/utils/pdf-upload/statusMapper'
import {
  Checkbox,
  InputNumber,
  Input,
  Rate,
  Tag,
  Space,
  Popover,
  Dropdown,
  Menu,
  Tooltip,
  Divider,
} from 'ant-design-vue'
import {
  EditOutlined,
  EllipsisOutlined,
  InfoCircleOutlined,
  ReloadOutlined,
} from '@ant-design/icons-vue'
import { liteThrottle } from '@idea/aiknowledge-special-util'
import ClassifyCell from './ClassifyCell.vue'
import RemarkCell from './RemarkCell.vue'
import ReadStatus from './ReadStatus.vue'

import { useLibraryIndex } from '@/stores/library'
import { useLibraryList, sortToKey } from '@/stores/library/list'

import Rollback from './Rollback.vue'
import { 
  getAuthors,
  updateAuthors,
  updatePublishDate,
  updateVenue,
  updateImpactFactor,
  updateJcr,
  updateScore,
  getUserDocStatusByIds,
 } from '@/api/material'
import { reParsePaper } from '@/api/parse'
import {
  escapeHightlight,
  getAuthorName,
  getOneVenue,
  getPopupContainer,
} from '../helper'
import { goPathPage, goPdfPage } from '@/common/src/utils/url'

const { Item } = Menu
const { TextArea } = Input

const currentYear = new Date().getFullYear()

export interface ArrangeEmitData {
  type: 'delete' | 'sort' | 'top'
  from: number
  to: number
}

const props = defineProps({
  collating: {
    type: Boolean,
    default: false,
  },
  item: {
    type: Object as () => Required<UserDocInfo>,
    default: () => {},
  },
  paperId: {
    type: String,
    default: '',
  },
  index: {
    type: Number,
    default: 0,
  },
  total: {
    type: Number,
    default: 0,
  },
})

const storeLibraryIndex = useLibraryIndex()

const storeLibraryList = useLibraryList()
const { t } = useI18n()

const hasPDF = computed(() => props.item.pdfId && props.item.pdfId !== '0')
const isNotMatched = computed(
  () => props.item.paperRepositoryStatus !== PaperRepositoryStatus.IN_REPOSITORY
)
const isEdit = ref(false as boolean)
const editContent = ref('' as string)

const handlejump = (param: string) => {
  if (param) {
    goPathPage(`/paper/${props.paperId}/${param}`)
  } else {
    goPathPage(`/paper/${props.paperId}`)
  }
}

const openPage = (item: UserDocInfo) => {
  if (item.pdfId && item.pdfId !== '0') {
    goPdfPage({ pdfId: item.pdfId })
  } else if (item.paperId && item.paperId !== '0') {
    // 增加 pdfId 和 noteId 参数
    goPathPage(`/note/${item.paperId}`)
  }
}

const startEditFileName = () => {
  editContent.value = props.item.docName
  isEdit.value = true
}

const cancelEditFileName = () => {
  isEdit.value = false
  editContent.value = ''
}

const submitFileName = liteThrottle(
  async () => {
    const item = storeLibraryList.paperListAll.find(
      (paper) => paper.docId === props.item.docId
    )!

    item.docName = editContent.value
    if (item.searchResult.hitDocName) {
      item.searchResult.hitDocName = ''
    }
    storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]

    const response = updateDoc({
      docId: props.item.docId,
      docName: editContent.value,
    })

    cancelEditFileName()

    if (await response) {
      storeLibraryIndex.fetchLibraryIndex()
    }
  },
  1000,
  false,
  true
)

const blurFileName = () => {
  if (!isEdit.value) {
    return
  }

  submitFileName()
}

const authors = computed(() => {
  const { authorInfos } = props.item.displayAuthor
  return authorInfos.map(getAuthorName).join(' / ')
})

const publishDate = computed(() => {
  return props.item.displayPublishDate.publishDate
})

const checkItem = () => {
  const index = storeLibraryList.paperListChecked.indexOf(props.item.docId)
  if (index === -1) {
    storeLibraryList.paperListChecked.push(props.item.docId)
  } else {
    storeLibraryList.paperListChecked.splice(index, 1)
  }
}

const editPublishDate = () => {
  const parts = props.item.displayPublishDate!.publishDate.split('-')
  const nanToUndefined = (string?: string) => {
    if (!string) {
      return undefined
    }

    return Number(string)
  }
  storeLibraryList.paperListPublishEdit = props.item.docId
  storeLibraryList.paperListPublishYear = nanToUndefined(parts.shift())
  storeLibraryList.paperListPublishMonth = nanToUndefined(parts.shift())
  storeLibraryList.paperListPublishDate = nanToUndefined(parts.shift())
}
const cancelPublishDate = () => {
  storeLibraryList.paperListPublishEdit = ''
  storeLibraryList.paperListPublishYear = currentYear
  storeLibraryList.paperListPublishMonth = undefined
  storeLibraryList.paperListPublishDate = undefined
}
const visibleChangePublishDate = (visible: boolean) => {
  if (!visible && storeLibraryList.paperListPublishEdit === props.item.docId) {
    cancelPublishDate()
  }
}
const submitPublishDate = async () => {
  await updatePublishDate({
    docId: props.item.docId,
    publishDate: [
      storeLibraryList.paperListPublishYear,
      storeLibraryList.paperListPublishMonth,
      storeLibraryList.paperListPublishDate,
    ]
      .filter(Boolean)
      .map((number) => String(number).padStart(2, '0'))
      .join('-'),
  })

  storeLibraryList.getFilesByFolderId()
  cancelPublishDate()
}
const rollbackPublishDate = async () => {
  await updatePublishDate({
    docId: props.item.docId,
  })

  storeLibraryList.getFilesByFolderId()
  publishDateRollbackVisible.value = false
}
const publishDateDisabled = computed(() => {
  if (storeLibraryList.paperListPublishDate) {
    if (
      !storeLibraryList.paperListPublishMonth ||
      !storeLibraryList.paperListPublishYear
    ) {
      return true
    }

    const date = new Date()
    date.setFullYear(storeLibraryList.paperListPublishYear)
    date.setMonth(storeLibraryList.paperListPublishMonth - 1)
    date.setDate(storeLibraryList.paperListPublishDate)

    return date.getMonth() !== storeLibraryList.paperListPublishMonth - 1
  } else if (storeLibraryList.paperListPublishMonth) {
    return !storeLibraryList.paperListPublishYear
  } else {
    return false
  }
})

const editVenue = async () => {
  storeLibraryList.paperListVenueEditContent = getOneVenue(
    props.item.displayVenue.venueInfos
  )
  storeLibraryList.paperListVenueEdit = props.item.docId
  await nextTick()
  await new Promise((resolve) => setTimeout(resolve, 300))
  const textarea = document.querySelector(
    `[data-list-item-venue="${props.item.docId}"]`
  )
  if (textarea instanceof HTMLTextAreaElement) {
    textarea.focus()
  }
}

const submitVenue = liteThrottle(
  async () => {
    if (
      getOneVenue(props.item.displayVenue.venueInfos) !==
      storeLibraryList.paperListVenueEditContent
    ) {
      await updateVenue({
        docId: props.item.docId,
        venue: storeLibraryList.paperListVenueEditContent,
      })
      storeLibraryList.venueRefreshOptionsList()
      // eslint-disable-next-line vue/no-mutating-props
      props.item.displayVenue.venueInfos.splice(
        0,
        props.item.displayVenue.venueInfos.length,
        storeLibraryList.paperListVenueEditContent
      )
      // eslint-disable-next-line vue/no-mutating-props
      props.item.displayVenue.userEdited = true
      storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
    }

    storeLibraryList.paperListVenueEdit = ''
    storeLibraryList.paperListVenueEditContent = ''
  },
  1000,
  false,
  true
)
const rollbackVenue = async () => {
  await updateVenue({
    docId: props.item.docId,
  })
  storeLibraryList.venueRefreshOptionsList()
  props.item.displayVenue!.venueInfos =
    props.item.displayVenue!.originVenueInfos
  props.item.displayVenue!.userEdited = false
  storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
  venueRollbackVisible.value = false
}

const editJcr = async () => {
  storeLibraryList.paperListJcrEditContent =
    props.item.jcrVenuePartion.jcrVenuePartion || ''
  storeLibraryList.paperListJcrEdit = props.item.docId
  await nextTick()
  await new Promise((resolve) => setTimeout(resolve, 300))
  const textarea = document.querySelector(
    `[data-list-item-jcr="${props.item.docId}"]`
  )
  if (textarea instanceof HTMLTextAreaElement) {
    textarea.focus()
  }
}

const submitJcr = liteThrottle(
  async () => {
    if (
      props.item.jcrVenuePartion.jcrVenuePartion !==
      storeLibraryList.paperListJcrEditContent
    ) {
      await updateJcr({
        docId: props.item.docId,
        jcrPartion: storeLibraryList.paperListJcrEditContent,
      })
      storeLibraryList.jcrRefreshOptionsList()
      // eslint-disable-next-line vue/no-mutating-props
      props.item.jcrVenuePartion.jcrVenuePartion =
        storeLibraryList.paperListJcrEditContent
      // eslint-disable-next-line vue/no-mutating-props
      props.item.jcrVenuePartion.userEdited = true
      storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
    }

    storeLibraryList.paperListJcrEdit = ''
    storeLibraryList.paperListJcrEditContent = ''
  },
  1000,
  false,
  true
)

const rollbackJcr = async () => {
  await updateJcr({
    docId: props.item.docId,
    jcrPartion: props.item.jcrVenuePartion.originJcrVenuePartion as string,
  })
  storeLibraryList.authorRefreshOptionsList()
  props.item.jcrVenuePartion!.jcrVenuePartion =
    props.item.jcrVenuePartion!.originJcrVenuePartion
  props.item.jcrVenuePartion!.userEdited = false
  storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
  jcrRollbackVisible.value = false
}

const editImpactFactor = async () => {
  storeLibraryList.paperListImpactFactorEditContent = String(
    props.item.impactOfFactor.impactOfFactor || ''
  )
  storeLibraryList.paperListImpactFactorEdit = props.item.docId
  await nextTick()
  await new Promise((resolve) => setTimeout(resolve, 300))
  const textarea = document.querySelector(
    `[data-list-item-impact="${props.item.docId}"]`
  )
  if (textarea instanceof HTMLTextAreaElement) {
    textarea.focus()
  }
}

const submitImpactFactor = liteThrottle(
  async () => {
    const number = Number(storeLibraryList.paperListImpactFactorEditContent)

    if (!isNaN(number)) {
      const impactOfFactor = storeLibraryList.paperListImpactFactorEditContent
        ? number
        : (null as unknown as undefined)

      if (props.item.impactOfFactor.impactOfFactor !== impactOfFactor) {
        await updateImpactFactor({
          docId: props.item.docId,
          impactOfFactor: impactOfFactor as number,
        })
        // eslint-disable-next-line vue/no-mutating-props
        props.item.impactOfFactor.impactOfFactor = impactOfFactor
        // eslint-disable-next-line vue/no-mutating-props
        props.item.impactOfFactor.userEdited = true
        storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
      }
    }

    storeLibraryList.paperListImpactFactorEdit = ''
    storeLibraryList.paperListImpactFactorEditContent = ''
  },
  1000,
  false,
  true
)

const rollbackImpactFactor = async () => {
  await updateImpactFactor({
    docId: props.item.docId,
    impactOfFactor: props.item.impactOfFactor!.originImpactOfFactor as number,
  })
  props.item.impactOfFactor!.impactOfFactor =
    props.item.impactOfFactor!.originImpactOfFactor
  props.item.impactOfFactor!.userEdited = false
  storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
  impactFactorRollbackVisible.value = false
}

const submitScore = async (score: number) => {
  const response = await updateScore({
    docId: props.item.docId,
    score,
  })

  if (response.data.data) {
    // eslint-disable-next-line vue/no-mutating-props
    props.item.importantanceScore = score
    storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
  }
}

const rollbackAuthor = async () => {
  const response = await getAuthors({
    docId: props.item.docId,
  })
  await updateAuthors({
    docId: props.item.docId,
  })
  storeLibraryList.authorRefreshOptionsList()
  props.item.displayAuthor!.authorInfos =
    response.displayAuthor!.originAuthors.map((author) => ({
      literal: author.name,
      isAuthentication: author.isAuthentication,
    }))
  props.item.displayAuthor!.userEdited = false
  storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
  storeLibraryList.paperListAuthorEdit = ''
  authorRollbackVisible.value = false
}

const authorRollbackVisible = ref(false)
const publishDateRollbackVisible = ref(false)
const venueRollbackVisible = ref(false)
const jcrRollbackVisible = ref(false)
const impactFactorRollbackVisible = ref(false)

const updateReadStatus = (
  item: UserDocInfo,
  { status, progress }: { status: DocReadingStatus; progress: number }
) => {
  item.docReadingStatus = status
  item.progress = progress
}

const preventDefault = (event: Event) => {
  event.preventDefault()
}

const getHighlightedAuthors = (item: UserDocInfo) => {
  if (!item.searchResult.hitAuthor) {
    return authors.value
  }
  return escapeHightlight(item.searchResult.hitAuthor)
}

// ========== 重新解析功能 ==========
// 从父组件注入共享的状态
const pendingReparseDocIds = inject<Ref<string[]>>('pendingReparseDocIds', ref([]))
const globalPollingTimer = inject<Ref<ReturnType<typeof setTimeout> | null>>('globalPollingTimer', ref(null))
const isFetching = inject<Ref<boolean>>('isFetching', ref(false))

const POLLING_INTERVAL = 3000 // 3秒轮询一次

// 判断是否应该显示重新解析图标
const shouldShowReparseIcon = (item: UserDocInfo) => {
  const progress = (item as any).parsedProgress
  if (!progress) return false
  
  // 提取百分比数字
  const percentMatch = progress.match(/\d+/)
  if (!percentMatch) return false
  
  const percent = parseInt(percentMatch[0], 10)
  return percent < 100
}

// 判断文档是否正在重新解析
const isReparsing = (docId: string | bigint) => {
  const id = String(docId)
  return pendingReparseDocIds.value.includes(id)
}

// 处理重新解析
const handleReparse = async (item: UserDocInfo) => {
  const pdfIdStr = String(item.pdfId || '')
  const docIdStr = String(item.docId)
  
  if (!item.pdfId || pdfIdStr === '0') {
    message.error('文档ID无效')
    return
  }

  if (isReparsing(docIdStr)) {
    message.warning('该文档正在重新解析中，请稍候')
    return
  }

  try {
    // 添加到待刷新数组（如果不存在）
    if (!pendingReparseDocIds.value.includes(docIdStr)) {
      pendingReparseDocIds.value.push(docIdStr)
    }
    
    // 调用重新解析API
    await reParsePaper({ pdfId: item.pdfId })
    message.success('重新解析请求已提交')
    
    // 只在没有正在请求且没有轮询定时器时，才立即触发第一次状态获取
    // 否则等待当前轮询周期自然完成
    if (!isFetching.value && !globalPollingTimer.value) {
      fetchAllDocStatus()
    }
  } catch (error) {
    console.error('Reparse failed:', error)
    message.error('重新解析请求失败，请稍后重试')
    // 失败时从数组中移除
    const index = pendingReparseDocIds.value.indexOf(docIdStr)
    if (index > -1) {
      pendingReparseDocIds.value.splice(index, 1)
    }
  }
}

// 批量获取所有待刷新文档的状态
const fetchAllDocStatus = async () => {
  // 如果没有待刷新的文档，直接返回
  if (pendingReparseDocIds.value.length === 0) {
    stopGlobalPolling()
    return
  }

  // 如果正在请求中，跳过本次调用
  if (isFetching.value) {
    return
  }

  try {
    isFetching.value = true
    
    // 获取当前所有待刷新的 docId
    const docIds = [...pendingReparseDocIds.value]
    
    // @ts-ignore - API 接受字符串类型的 docId 数组
    const response = await getUserDocStatusByIds({
      docIds: docIds
    })
    
    // 用于存储已完成的 docId
    const completedDocIds: string[] = []
    
    // @ts-ignore
    if (response.items && response.items.length > 0) {
      // 遍历返回的所有文档状态
      // @ts-ignore
      response.items.forEach((docStatus) => {
        const docId = String(docStatus.docId)
        
        // 查找列表中的文档
        const item = storeLibraryList.paperListAll.find(
          (paper) => String(paper.docId) === docId
        )
        
        if (item) {
          // 使用封装的方法计算进度
          const progress = calculateParsedProgress(docStatus.status, docStatus.embeddingStatus)
          
          // 本地更新 parsedProgress 字段
          ;(item as any).parsedProgress = progress
          
          // 同时更新 parsedStatus 和 embeddingStatus
          ;(item as any).parsedStatus = docStatus.status
          ;(item as any).embeddingStatus = docStatus.embeddingStatus
          
          // 如果进度达到100%，标记为已完成
          if (progress === '100%') {
            completedDocIds.push(docId)
          }
        }
      })
      
      // 触发响应式更新（一次性更新所有变化）
      storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
      
      // 从待刷新数组中移除已完成的文档
      pendingReparseDocIds.value = pendingReparseDocIds.value.filter(
        id => !completedDocIds.includes(id)
      )
    }
    
    // 根据是否还有待刷新的文档，决定是否继续轮询
    if (pendingReparseDocIds.value.length === 0) {
      // 所有文档都已完成，停止轮询
      stopGlobalPolling()
    }
  } catch (error) {
    console.error('Failed to fetch doc status:', error)
  } finally {
    isFetching.value = false
    
    // 等待 3 秒后，如果还有待刷新的文档，直接调用下一次 fetchAllDocStatus
    // 只在没有现有定时器时才设置新定时器，避免重置倒计时
    if (pendingReparseDocIds.value.length > 0 && !globalPollingTimer.value) {
      globalPollingTimer.value = setTimeout(() => {
        globalPollingTimer.value = null
        if (pendingReparseDocIds.value.length > 0) {
          fetchAllDocStatus()
        }
      }, 3000)
    }
  }
}

// 开始全局轮询
const startGlobalPolling = () => {
  // 先停止之前的轮询
  stopGlobalPolling()
  
  // 设置新的轮询定时器
  globalPollingTimer.value = setTimeout(() => {
    fetchAllDocStatus()
  }, POLLING_INTERVAL)
}

// 停止全局轮询
const stopGlobalPolling = () => {
  if (globalPollingTimer.value) {
    clearTimeout(globalPollingTimer.value)
    globalPollingTimer.value = null
  }
}

// 组件卸载时清理定时器（不清空数组，因为是共享的）
onUnmounted(() => {
  stopGlobalPolling()
})
</script>
<style lang="less" scoped>
.literature-list-item {
  display: flex;
  align-items: stretch;
  padding-left: 1px !important;
  padding-right: 1px !important;
  border-bottom: 1px solid var(--site-theme-border-color);
  background-color: var(--site-theme-background-primary);
  
  &:hover {
    background: var(--site-theme-background-hover);
  }

  .literature-list-item-check {
    padding-left: 12px;
    padding-right: 8px;
    display: flex;
    align-items: center;
  }
  .literature-list-item-col {
    position: relative;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    background-color: var(--site-theme-background-primary);
    border-bottom: 0.5px solid var(--site-theme-border-color);
    .sort-overlay {
      position: absolute;
      height: 100%;
      width: 100%;
      left: 0;
      top: 0;
      background-color: var(--site-theme-background-hover);
      opacity: 0.3;
      pointer-events: none;
    }
    > * {
      padding-left: 4px;
      padding-right: 4px;
      padding-top: 4px;
      padding-bottom: 4px;
    }
  }

  .literature-list-item-col-docName,
  .literature-list-item-wrap-docName,
  .literature-list-item-col-remark,
  .literature-list-item-wrap-remark,
  .literature-list-item-col-authors,
  .literature-list-item-wrap-authors {
    overflow: hidden;
  }
  :deep(em) {
    font-style: normal;
    color: var(--site-theme-primary-color);
    background-color: var(--site-theme-primary-color-fade);
  }
  .edit-filename {
    display: flex;
    align-items: center;
    max-width: 100%;
    .text {
      max-width: calc(100% - 50px);
      overflow: hidden;
    }
  }
  .iconedit {
    width: 25px;
    height: 25px;
    background: #dfe6f0;
    color: var(--site-theme-text-tertiary);
    margin-left: 5px;
    opacity: 0;
    display: inline-block;
    font-size: 16px;
    line-height: 25px;
  }
  .title:hover .iconedit {
    opacity: 1 !important;
  }

  .tag-last-read {
    padding: 0 3px;
  }
  .tag-nopdf {
    color: var(--site-theme-text-secondary);
    border-color: var(--site-theme-border-color);
    font-weight: normal;
    font-family:
      'Lato',
      -system-ui,
      -apple-system;
  }
  .tag-attached {
    :deep(.ant-divider) {
      border-color: var(--site-theme-border-color);
      margin-left: 0;
    }
  }
  
  .tag-last-read {
    background-color: var(--site-theme-primary-color) !important;
    color: var(--site-theme-text-inverse) !important;
    border: none;
  }
  .title {
    padding-top: 4px;
    padding-bottom: 4px;
    padding-left: 12px;
    font-size: 15px;
    // stylelint-disable-next-line
    font-family: Lato-Bold;
    font-weight: bold;
    color: var(--site-theme-text-primary);
    display: flex;
    flex-direction: column;
    height: 100%;
    width: 100%;
    justify-content: center;
    align-items: flex-start;

    .content {
      font-size: 14px;
      font-weight: 600;
      line-height: 24px;
    }
  }
  .classify {
    overflow: hidden;
  }
  .remark {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  // 解析进度列样式
  .parse-progress-cell {
    display: flex;
    align-items: center;
    gap: 8px;

    .reparse-icon {
      font-size: 16px;
      color: var(--site-theme-brand);
      cursor: pointer;
      transition: opacity 0.2s;

      &:hover {
        opacity: 0.8;
      }

      &.reparse-spinning {
        animation: rotate 1s linear infinite;
      }
    }
  }

  @keyframes rotate {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
}
.dropdownMenu {
  z-index: 99 !important;
  .ant-dropdown-menu {
    padding: 0;
    .ant-dropdown-menu-item {
      padding: 10px 48px 10px 16px;
      .menu-name {
        font-size: 14px;
        font-weight: 400;
        color: var(--site-theme-text-secondary);
        position: relative;
        &.not-match {
          color: rgba(0, 0, 0, 0.4);
        }
        .anticon,
        .icon-file-done {
          margin-left: 4px;
        }
        .anticon {
          top: 4px;
        }
        .icon-file-done {
          top: 1px;
        }
      }
      &:first-child {
        padding-top: 18px;
      }
      &:last-child {
        padding-bottom: 18px;
      }
      &:hover {
        background: var(--site-theme-background-hover);
      }
    }
  }
}

.literature-list-item-authors {
  display: flex;
  > span {
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
  }
}

.literature-list-item-venue {
  overflow: hidden;
  display: -webkit-box;
  display: flex;
  -webkit-box-orient: vertical;
  box-orient: vertical;
  -webkit-line-clamp: 5;
  line-clamp: 5;
  textarea {
    width: 100%;
    border: 1px solid var(--site-theme-border-color);
    outline: 0;
    display: block;
  }
}

.literature-list-item-authors:hover,
.literature-list-item-publishDate:hover,
.hover-background:hover {
  background-color: var(--site-theme-background-hover);
}

.list-edit-publishDate {
  .list-edit-publishDate-header {
    color: var(--site-theme-text-primary);
    font-weight: bold;
    height: 22px;
    margin-bottom: 18px;
  }
  .list-edit-publishDate-body {
    margin-bottom: 20px;
    > * {
      width: 66px;
    }
  }
  .list-edit-publishDate-footer {
    display: flex;
    justify-content: flex-end;
    button {
      margin-left: 12px;
      width: 64px;
      height: 24px;
      border: 0;
      outline: 0;
      border-radius: 2px;
      font-size: 12px;
      cursor: pointer;
      &:first-of-type {
        background: var(--site-theme-background-secondary);
        color: var(--site-theme-text-secondary);
      }
      &:last-of-type {
        background: var(--site-theme-primary-color);
        color: var(--site-theme-text-inverse);
      }
      &:disabled {
        opacity: 0.7;
      }
    }
  }
}

:deep(.ant-dropdown-link:hover) {
  > * {
    background: var(--site-theme-background-hover);
    border-radius: 2px;
  }

  * {
    color: var(--site-theme-text-primary);
  }
}
</style>

<style>
.list-item-col-rollback {
  position: absolute;
  right: 2px;
  bottom: 4px;
  height: 24px;
  width: 24px;
  display: none;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  color: var(--site-theme-text-secondary);
  background-color: var(--site-theme-background-primary);
  border-radius: 2px;
  cursor: pointer;
}

.literature-list-item-authors,
.literature-list-item-publishDate,
.literature-list-item-venue {
  position: relative;
  padding-top: 4px;
  padding-bottom: 4px;
  padding-left: 8px;
  padding-right: 8px;
  border-radius: 2px;
  min-height: 32px;
  cursor: pointer;
  &:hover {
    .list-item-col-rollback {
      display: flex;
    }
  }
}

.list-item-textarea {
  padding-left: 7px !important;
  padding-right: 7px !important;
  min-height: 32px !important;
  resize: none !important;
}
</style>
