<template>
  <div class="classify-cell" @click="handleToggleInput(true)">
    <span v-if="!showSelectInput" class="tag">
      <Tag
        v-for="tag in currentTags"
        :key="tag.key"
        :closable="false"
        class="classify-tag"
        >{{ tag.label }}</Tag
      >

      <span
        v-if="currentTags.length === 0"
        v-show="!showSelectInput"
        aria-hidden="true"
        class="classify-add"
      >
        <i class="aiknowledge-icon icon-add" aria-hidden="true"
      /></span>
    </span>
    <Select
      v-if="showSelectInput"
      ref="selectRef"
      :value="currentTags as any"
      mode="multiple"
      style="width: 100%; margin: 0 !important"
      :placeholder="$t('home.library.selectTag')"
      label-in-value
      autofocus
      :get-popup-container="getPopupContainer"
      default-open
      class="select"
      option-label-prop="label"
      :default-active-first-option="false"
      :not-found-content="null"
      :filter-option="false"
      @select="handleSelect"
      @deselect="handleDeselect as any"
      @search="handleSearch"
      @input-key-down="handleKeydown"
    >
      <SelectOption
        v-for="item in filteredOptions"
        :key="item.key"
        :value="item.key"
        :label="item.label"
      >
        <Tag style="cursor: pointer" class="classify-tag" :closable="false">{{
          item.label
        }}</Tag>
      </SelectOption>
    </Select>
  </div>
</template>
<script lang="ts" setup>
import { UserDocClassifyInfo } from 'go-sea-proto/gen/ts/doc/UserDocManage'
import { computed, PropType, ref, watch } from 'vue'
import { Select, SelectOption } from 'ant-design-vue'
import { delay } from '@/common/src/utils/aiknowledge-special-util'
import _ from 'lodash'

import { onClickOutside } from '@vueuse/core'
import Tag from '../Head/Tag.vue'
import {
  addClassify,
  attachDocToClassify,
  removeDocFromClassify,
} from '@/api/material'
import { getPopupContainer } from '../helper'

import { useLibraryList } from '@/stores/library/list'
import { useClassify } from '@/stores/classify'
import { DefaultOptionType } from 'ant-design-vue/lib/select'

interface ClassifyLabel {
  key: string
  label: string
}

const storeLibraryList = useLibraryList()
const storeClassify = useClassify()

const props = defineProps({
  classifyInfos: {
    type: Array as PropType<UserDocClassifyInfo[]>,
    default: () => [],
  },
  docId: {
    type: String,
    default: '',
  },
})

const allClassifyList = computed<ClassifyLabel[] | null>(() => {
  return storeClassify.classifyList.map((item) => {
    return {
      key: item.classifyId,
      label: item.classifyName,
    }
  })
})
const currentTags = ref<ClassifyLabel[]>(
  (props.classifyInfos || []).map((item) => {
    return {
      key: item.classifyId,
      label: item.classifyName,
    }
  })
)

const updateOfflineClassifyInfos = (tags: ClassifyLabel[]) => {
  storeLibraryList.paperListClassifyRefresh()
  const doc = storeLibraryList.paperListAll.find(
    (paper) => paper.docId === props.docId
  )!
  doc.classifyInfos = tags.map((tag) => {
    const exists = doc.classifyInfos.find(
      (info: { classifyId: string }) => info.classifyId === tag.key
    )

    if (exists) {
      return exists
    }

    return {
      classifyId: tag.key,
      classifyName: tag.label,
      docId: props.docId,
    }
  })
  storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
}

const showSelectInput = ref<boolean>(false)
const filteredOptions = computed<ClassifyLabel[]>(() => {
  if (!allClassifyList.value?.length) {
    return []
  }
  const list = _.pullAllBy([...allClassifyList.value], currentTags.value, 'key')
  if (!inputValue.value) {
    return list
  }
  return list.filter((item) => {
    return item.label.includes(inputValue.value)
  })
})
const inputValue = ref<string>('')

const handleToggleInput = async (show: boolean) => {
  if (show) {
    await storeClassify.initClassifyList()
    window.setTimeout(() => {
      if (inputRef.value) {
        inputRef.value.focus()
      }
    }, 200)
  }

  showSelectInput.value = show
}

const handleEnter = async () => {
  const value = String(inputValue.value).trim()
  if (!value) {
    handleToggleInput(false)
    return
  }
  const findInCurrent = currentTags.value.find((item) => item.label === value)
  if (findInCurrent) {
    handleToggleInput(false)
    inputValue.value = ''
    return
  }
  let findInAll = filteredOptions.value.find((item) => item.label === value)
  if (!findInAll) {
    try {
      const res = await addClassify({ classifyName: value })
      findInAll = {
        key: res.classifyId,
        label: res.classifyName,
      }
      storeClassify.refreshClassifyList()
    } catch (error) {
      return
    }
  }
  if (findInAll) {
    const result = await attachDocToClassify({
      docId: props.docId,
      classifyId: findInAll.key,
    })
    if (result) {
      updateOfflineClassifyInfos([...currentTags.value, findInAll])
      inputValue.value = ''
      handleToggleInput(false)
    }
  }
}

const inputRef = ref<HTMLElement>()

watch(
  () => [...(props.classifyInfos || [])],
  (newValue, oldValue) => {
    const a = newValue || []
    const b = oldValue || []
    if (a.length === b.length && !_.differenceBy(a, b, 'classifyId').length) {
      return
    }
    currentTags.value = a.map((item) => {
      return {
        key: item.classifyId,
        label: item.classifyName,
      }
    })
  }
)

let latestSelectTiming = 0

const handleSelect = async (...[, value]: [unknown, DefaultOptionType]) => {
  console.log('handleSelect', value)
  latestSelectTiming = Date.now()
  const tag: ClassifyLabel = {
    key: value.value as string,
    label: value.label as string,
  }
  const result = await attachDocToClassify({
    docId: props.docId,
    classifyId: tag.key,
  })
  if (result) {
    updateOfflineClassifyInfos([...currentTags.value, tag])
  }
}
const handleDeselect = async (value: ClassifyLabel) => {
  try {
    await removeDocFromClassify({
      classifyId: value.key || '',
      docId: props.docId,
    })
    storeClassify.refreshClassifyList()
    updateOfflineClassifyInfos(
      currentTags.value.filter((tag) => tag.key !== value.key)
    )
  } catch (error) {}
}

const selectRef = ref()

onClickOutside(selectRef, (event: PointerEvent) => {
  if (
    (event.target as HTMLElement)?.closest('.ant-select-dropdown--multiple ')
  ) {
    return
  }
  inputValue.value = ''
  handleToggleInput(false)
})

const handleSearch = (val: string) => {
  console.log('handleSearch', val)
  inputValue.value = val
}

const handleKeydown = async (event: KeyboardEvent) => {
  console.log('handleKeydown', event)

  if (event.code !== 'Enter' && event.key !== 'Enter') {
    return
  }

  await delay(300)

  if (Date.now() - latestSelectTiming < 600) {
    return
  }

  handleEnter()
}
</script>
<style lang="less" scoped>
.classify-cell {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  height: 100%;
  min-height: 32px;
  overflow: hidden;
  .classify-tag {
    margin-top: 4px;
    margin-bottom: 4px;
  }
  .classify-add {
    cursor: pointer;
    color: var(--site-theme-text-tertiary);
    display: none;
    background: var(--site-theme-background-tertiary);
    padding: 0 6px;
    margin-top: 4px;
    margin-bottom: 4px;
    line-height: 24px;
    height: 24px;
    .icon-add {
      font-size: 12px;
    }
  }
  &:hover {
    cursor: pointer;
    .classify-add {
      display: flex;
      align-items: center;
    }
  }
  .select {
    border: 1px solid var(--site-theme-border-color);
    margin: 4px 0;
    cursor: pointer;
  }
}
</style>
<style lang="less">
.classify-cell {
  .ant-select {
    padding-left: 0 !important;
    .ant-select-selector {
      border: 0 !important;
      outline: 0 !important;
      background-color: var(--site-theme-background-primary) !important;
    }
    .ant-select-selection-item {
      color: var(--site-theme-text-primary) !important;
    }
  }
}
</style>
