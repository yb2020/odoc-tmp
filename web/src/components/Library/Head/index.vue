<template>
  <div class="literature-list-mode-table">
    <div v-if="collating" class="literature-list-item-check" style="background: var(--site-theme-background-secondary)">
      <Checkbox :checked="paperListAllChecked" :indeterminate="paperListIndeterminate" @change="toggleChecked()" />
    </div>
    <div v-for="item in storeLibraryList.paperHeadVisibleList" :key="item.key"
      :style="storeLibraryList.paperHeadExtra[item.key].style as StyleValue" class="literature-list-head">
      <div v-if="sortToKey[storeLibraryList.currentSortType] === item.key" class="sort-overlay"></div>
      <div v-if="item.key === 'docName'" class="literature-list-head-docName"
        @click="sortBy(UserDocListSortType.DOC_NAME)">
        {{ $t(TableHeadName[item.key]) }}
        <div class="literature-list-sort">
          <CaretUpOutlined :style="{
              color:
                sortToKey[storeLibraryList.currentSortType] === item.key &&
                storeLibraryList.paperListLocalSortDirection === 1
                  ? 'var(--site-theme-primary-color)'
                  : 'var(--site-theme-text-tertiary)',
            }" />
          <CaretDownOutlined :style="{
              color:
                sortToKey[storeLibraryList.currentSortType] === item.key &&
                storeLibraryList.paperListLocalSortDirection === -1
                  ? 'var(--site-theme-primary-color)'
                  : 'var(--site-theme-text-tertiary)',
            }" />
        </div>
      </div>
      <Classify v-else-if="item.key === 'classifyInfos'" />
      <HeadFilter v-else-if="item.key === 'authors'" :text="$t(TableHeadName[item.key])" :use-keyword="true"
        :search-input="storeLibraryList.searchInput" :option-list="storeLibraryList.authorOptionList"
        :checked="storeLibraryList.authorChecked" :empty-checked="storeLibraryList.authorEmptyChecked"
        :all-checked="storeLibraryList.authorAllChecked" :indeterminate="storeLibraryList.authorIndeterminate"
        @toggle="storeLibraryList.authorToggle($event)" @empty-toggle="storeLibraryList.authorEmptyToggle()"
        @all-toggle="storeLibraryList.authorAllToggle()" />
      <HeadFilter v-else-if="item.key === 'displayVenue'" :text="$t(TableHeadName[item.key])"
        :search-input="storeLibraryList.searchInput" :option-list="storeLibraryList.venueOptionList"
        :checked="storeLibraryList.venueChecked" :empty-checked="storeLibraryList.venueEmptyChecked"
        :all-checked="storeLibraryList.venueAllChecked" :indeterminate="storeLibraryList.venueIndeterminate"
        @toggle="storeLibraryList.venueToggle($event)" @empty-toggle="storeLibraryList.venueEmptyToggle()"
        @all-toggle="storeLibraryList.venueAllToggle()" />
      <div v-else-if="item.key === 'remark' && item.visible" :key="`remark-${item.key}`" flex="150px">
        {{ $t(TableHeadName[item.key]) }}
      </div>
      <div v-else-if="item.key === 'parseProgress'" :key="`parseProgress-${item.key}`" flex="150px">
        {{ $t(TableHeadName[item.key]) }}
      </div>
      <div v-else-if="item.key === 'publishDate'" class="literature-list-head-publishDate"
        @click="sortBy(UserDocListSortType.PUBLISH_DATE)">
        {{ $t(TableHeadName[item.key]) }}
        <div class="literature-list-sort">
          <CaretUpOutlined :style="{
              color:
                sortToKey[storeLibraryList.currentSortType] === item.key &&
                storeLibraryList.paperListLocalSortDirection === 1
                  ? 'var(--site-theme-primary-color)'
                  : 'var(--site-theme-text-tertiary)',
            }" />
          <CaretDownOutlined :style="{
              color:
                sortToKey[storeLibraryList.currentSortType] === item.key &&
                storeLibraryList.paperListLocalSortDirection === -1
                  ? 'var(--site-theme-primary-color)'
                  : 'var(--site-theme-text-tertiary)',
            }" />
        </div>
      </div>
      <HeadFilter v-else-if="item.key === 'jcrVenuePartion'" :text="$t(TableHeadName[item.key])"
        :search-input="storeLibraryList.searchInput" :option-list="storeLibraryList.jcrOptionList"
        :checked="storeLibraryList.jcrChecked" :empty-checked="storeLibraryList.jcrEmptyChecked"
        :all-checked="storeLibraryList.jcrAllChecked" :indeterminate="storeLibraryList.jcrIndeterminate"
        @toggle="storeLibraryList.jcrToggle($event)" @empty-toggle="storeLibraryList.jcrEmptyToggle()"
        @all-toggle="storeLibraryList.jcrAllToggle()" />
      <div v-if="item.key === 'impactOfFactor'" class="literature-list-head-docName"
        @click="sortBy(UserDocListSortType.IMPACT_OF_FACTOR)">
        <Dropdown v-model:visible="factorDropdownVisible">
          <div style="display: flex">
            {{ $t(TableHeadName[item.key]) }}
            {{ factorView }}
            <div class="literature-list-sort">
              <CaretUpOutlined :style="{
                  color:
                    sortToKey[storeLibraryList.currentSortType] === item.key &&
                    storeLibraryList.paperListLocalSortDirection === 1
                      ? 'var(--site-theme-primary-color)'
                      : 'var(--site-theme-text-tertiary)',
                }" />
              <CaretDownOutlined :style="{
                  color:
                    sortToKey[storeLibraryList.currentSortType] === item.key &&
                    storeLibraryList.paperListLocalSortDirection === -1
                      ? 'var(--site-theme-primary-color)'
                      : 'var(--site-theme-text-tertiary)',
                }" />
            </div>
          </div>
          <template #overlay>
            <div class="literature-list-head-impact-factor-filter">
              <div @click="factorToggleLimit()">
                <Checkbox :checked="storeLibraryList.paperListImpactFactorNoLimit" />
                <span>{{ $t('home.library.noLimit') }}</span>
              </div>
              <Divider />
              <div class="impact-dropdown-2">
                <InputNumber v-model:value="factorMin" :min="0" size="small" />
                -
                <InputNumber v-model:value="factorMax" :min="0" size="small" />
              </div>
              <div class="impact-dropdown-3">
                <Button size="small" @click="factorReset()">{{
                  $t('home.global.clear')
                  }}</Button>
                <Button type="primary" size="small" @click="factorSearch()">
                  {{ $t('home.global.ok') }}
                </Button>
              </div>
            </div>
          </template>
        </Dropdown>
      </div>
      <div v-if="item.key === 'importantanceScore'" class="literature-list-head-docName"
        @click="sortBy(UserDocListSortType.IMPORTANCE_SCORE)">
        {{ $t(TableHeadName[item.key]) }}
        <div class="literature-list-sort">
          <CaretUpOutlined :style="{
              color:
                sortToKey[storeLibraryList.currentSortType] === item.key &&
                storeLibraryList.paperListLocalSortDirection === 1
                  ? 'var(--site-theme-primary-color)'
                  : 'var(--site-theme-text-tertiary)',
            }" />
          <CaretDownOutlined :style="{
              color:
                sortToKey[storeLibraryList.currentSortType] === item.key &&
                storeLibraryList.paperListLocalSortDirection === -1
                  ? 'var(--site-theme-primary-color)'
                  : 'var(--site-theme-text-tertiary)',
            }" />
        </div>
      </div>
      <!-- 操作列已被注释掉
      <div v-else-if="item.key === 'operation'" class="quickAction" :flex="
          storeLibraryList.currentSortType !== UserDocListSortType.LAST_ADD &&
          collating &&
          storeLibraryList.searchInput
            ? '100px'
            : '245px'
        ">
        {{ $t('home.library.quickActions') }}
      </div>
      -->
      <div class="literature-list-head-resize" @mousedown="resizeStart($event, item.key)">
        <div></div>
      </div>
    </div>
    <Config />
  </div>
</template>
<script lang="ts" setup>
import { StyleValue, computed, ref } from 'vue'
import { UserDocListSortType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/UserDocManage'
import {
  message,
  Button,
  Dropdown,
  InputNumber,
  Divider,
  Checkbox,
} from 'ant-design-vue'
import { CaretUpOutlined, CaretDownOutlined } from '@ant-design/icons-vue'
import { liteThrottle } from '@idea/aiknowledge-special-util'
import {
  useLibraryList,
  TableHeadName,
  TableHeadKey,
  defaultTableHeadWidth,
  sortToKey,
} from '../../../stores/library/list'
import Classify from './Classify.vue'
import HeadFilter from './Filter.vue'
import Config from './Config.vue'

defineProps({
  collating: {
    type: Boolean,
    default: false,
  },
})

const storeLibraryList = useLibraryList()

const { getFilesByFolderId } = storeLibraryList

const factorDropdownVisible = ref(false)
const factorMin = ref<number | undefined>(undefined)
const factorMax = ref<number | undefined>(undefined)
const factorView = computed(() => {
  const undefinedToEmptyString = (value: number | undefined) => {
    if (value === undefined) {
      return ''
    }

    return value
  }

  if (
    storeLibraryList.paperListImpactFactorMin !== undefined ||
    storeLibraryList.paperListImpactFactorMax !== undefined
  ) {
    return `: ${undefinedToEmptyString(
      storeLibraryList.paperListImpactFactorMin
    )} - ${undefinedToEmptyString(storeLibraryList.paperListImpactFactorMax)}`
  }

  if (storeLibraryList.paperListImpactFactorNoLimit) {
    return ': 不限'
  }

  return ''
})

const factorSearch = () => {
  if (
    factorMin.value !== undefined &&
    factorMax.value !== undefined &&
    factorMin.value > factorMax.value
  ) {
    message.warn('最小值不能大于最大值')
    return
  }

  storeLibraryList.paperListImpactFactorMin = factorMin.value
  storeLibraryList.paperListImpactFactorMax = factorMax.value

  getFilesByFolderId()
}

const factorReset = () => {
  factorMin.value = undefined
  factorMax.value = undefined
  storeLibraryList.paperListImpactFactorMin = undefined
  storeLibraryList.paperListImpactFactorMax = undefined
  storeLibraryList.paperListImpactFactorNoLimit = false
  getFilesByFolderId()
}

const factorToggleLimit = () => {
  storeLibraryList.paperListImpactFactorNoLimit =
    !storeLibraryList.paperListImpactFactorNoLimit

  if (
    storeLibraryList.paperListImpactFactorMin === undefined &&
    storeLibraryList.paperListImpactFactorMax === undefined
  ) {
    getFilesByFolderId()
  }
}

let resizingIndex = NaN
let resizeOriginX = NaN
let resizeOriginWidth = NaN
const resizeStart = (event: MouseEvent, key: TableHeadKey) => {
  resizingIndex = storeLibraryList.paperHeadList.findIndex(
    (item) => item.key === key
  )
  resizeOriginX = event.clientX
  resizeOriginWidth = storeLibraryList.paperHeadList[resizingIndex].width
  document.addEventListener('mousemove', wrapedResizeMove, { passive: true })
  document.addEventListener('mouseup', () => {
    document.removeEventListener('mousemove', wrapedResizeMove)
  })
}

const resizeMove = (event: MouseEvent) => {
  const width = event.clientX - resizeOriginX + resizeOriginWidth
  const head = storeLibraryList.paperHeadList[resizingIndex]
  const minWidth = defaultTableHeadWidth[head.key]
  head.width = Math.max(width, minWidth)
  storeLibraryList.paperHeadSync()
}

const wrapedResizeMove = liteThrottle(resizeMove, 50)

const sortBy = (type: UserDocListSortType) => {
  if (type !== storeLibraryList.currentSortType) {
    storeLibraryList.currentSortType = type
    storeLibraryList.paperListLocalSortDirection = 1
  } else if (storeLibraryList.paperListLocalSortDirection === 1) {
    storeLibraryList.paperListLocalSortDirection = -1
  } else {
    storeLibraryList.currentSortType = storeLibraryList.dropdownSortType
    storeLibraryList.paperListLocalSortDirection = 1
  }

  getFilesByFolderId()
}

const paperListAllChecked = computed(() => {
  const values = Object.values(storeLibraryList.paperListCheckedMap)
  return Boolean(values.length) && values.every(Boolean)
})
const paperListIndeterminate = computed(() => {
  return (
    !paperListAllChecked.value &&
    Object.values(storeLibraryList.paperListCheckedMap).some(Boolean)
  )
})

const toggleChecked = () => {
  if (paperListAllChecked.value && !paperListIndeterminate.value) {
    storeLibraryList.paperListChecked = []
  } else {
    storeLibraryList.paperListChecked = storeLibraryList.paperListAll.map(
      (item) => item.docId
    )
  }
}
</script>
<style lang="less" scoped>
.literature-select-dropdown-item {
  padding: 9px 16px;
  cursor: pointer;

  .checkbox-wrap {
    display: flex;
    align-items: center;

    .checkbox {
      margin-right: 8px;
    }
  }
}

.literature-list-mode-table {
  display: flex;
  flex-direction: row;
  position: sticky;
  margin-left: 1px;
  top: 0;
  justify-content: flex-start;
  height: 40px;
  border-bottom: 1px solid var(--site-theme-border-color);
  color: var(--site-theme-text-primary);
  user-select: none;
  z-index: 2;
  background-color: var(--site-theme-background-secondary);

  .literature-list-item-check {
    padding-left: 12px;
    padding-right: 8px;
    display: flex;
    align-items: center;
    background-color: var(--site-theme-background-secondary);
  }

  .dot {
    position: sticky;
    right: 0;
    top: 8px;
    display: flex;
    align-items: center;
    background-color: var(--site-theme-background-secondary);

    .iconmore {
      font-size: 22px;
      color: var(--site-theme-text-primary);
    }
  }

  .empty {
    margin-top: 158px;
    margin-bottom: 260px;
  }

  .classify-link {
    padding: 0;
  }
}

.literature-select-dropdown-wrap {
  max-height: 300px;
  overflow: auto;
}

/deep/.ant-select-dropdown.ant-select-dropdown--single.ant-select-dropdown--empty {
  left: -101px !important;
}

.literature-list-head {
  position: relative;
  background: var(--site-theme-background-secondary);
  display: flex;
  align-items: center;

  >* {
    padding-left: 12px;
    flex: 1 1 100%;
  }

  .literature-list-head-resize {
    position: absolute;
    top: 0;
    right: 0;
    height: 100%;
    width: 6px;
    cursor: ew-resize;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    padding-left: 0 !important;

    div {
      height: 20px;
      width: 1px;
      background-color: var(--site-theme-border-color);
    }
  }

  .literature-list-head-docName,
  .literature-list-head-publishDate {
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
    align-items: center;
    cursor: pointer;
    margin-right: 8px;
  }

  .literature-list-sort {
    margin-left: 6px;
    display: flex;
    flex-direction: column;
    justify-content: center;

    >* {
      font-size: 12px;

      &:first-of-type {
        margin-bottom: -5px;
      }
    }
  }
}

.sort-overlay {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 100%;
  background-color: var(--site-theme-text-primary);
  opacity: 0.03;
  pointer-events: none;
}
</style>
<style lang="less">
.literature-list-head-impact-factor-filter {
  /* shadow-1 */
  box-shadow: 2px 6px 16px 0 rgba(0, 0, 0, 0.12);
  background-color: var(--site-theme-background-primary);
  padding: 8px 0;

  >div {

    &:first-of-type,
    &.impact-dropdown-2 {
      margin: 0 16px;
    }

    &:first-of-type {
      cursor: pointer;

      >span {
        margin-left: 8px;
        color: var(--site-theme-text-primary);
      }
    }

    &.ant-divider {
      margin: 8px 0;
      border-color: var(--site-theme-border-color);
    }

    &.impact-dropdown-2 {
      .ant-input-number {
        width: 54px;

        &:first-of-type {
          margin-right: 4px;
        }

        &:last-of-type {
          margin-left: 4px;
        }
      }
    }

    &.impact-dropdown-3 {
      margin-top: 8px;
      padding-right: 16px;
      padding-left: 16px;
      display: flex;
      justify-content: space-between;
      width: 100%;

      >*:first-of-type {
        background: var(--site-theme-background-secondary);
        color: var(--site-theme-text-secondary);
        border: 0;
        margin-right: 4px;
      }
    }
  }

  .impact-factor-search-icon {
    color: var(--site-theme-primary-color);
  }
}
</style>
