<template>
  <div class="list-author-edit-container" @click.stop>
    <div class="list-author-edit-header">
      <div class="title">
        {{ $t('home.library.viewAuthorInfo') }}
      </div>
      <div class="header-icons">
        <Tooltip placement="top">
          <template #title>
            <span>{{ $t('home.library.uptoAuthor') }}</span>
          </template>
          <QuestionCircleOutlined class="icon" />
        </Tooltip>
        <CloseOutlined class="icon" @click="cancelAuthorList()" />
      </div>
    </div>
    <div class="list-author-edit-body list-thin-scroll">
      <div class="list-author-edit-caption">
        {{ $t('home.library.authorName') }}
      </div>
      <div class="list-author-edit-list">
        <div
          v-for="(item, index) in authorList"
          :key="getAuthorKey(item, index)"
          class="list-author-edit-item"
        >
          <div v-if="authorEditingKey !== getAuthorKey(item, index)" class="author-display">
            <div
              :title="item.name"
              class="list-author-edit-view"
              :class="{ 'list-author-identify': isIdentify(item.id) }"
              :style="{
                color: item.isAuthentication ? '#1F71E0' : '#4E5969',
              }"
              @click="goAuthorPage(item.id)"
            >
              {{ item.name }}
            </div>
            <div class="icon-wrap">
              <EditOutlined @click.stop="editAuthor(item, index)" />
              <DeleteOutlined @click.stop="deleteAuthor(item)" />
            </div>
          </div>
          <input
            v-else
            v-model.trim="authorEditingName"
            autofocus
            :data-metadata-author="getAuthorKey(item, index)"
            @blur="enterAuthor()"
            @keyup.enter="enterAuthor()"
          />
        </div>
        <div class="list-author-edit-add" @click="addAuthor()">
          <PlusOutlined style="font-size: 10px" />
          {{ $t('home.library.addAuthor') }}
        </div>
      </div>
    </div>
    <div class="list-author-edit-footer">
      <button
        style="background: #f0f2f5; color: #4e5969"
        @click="cancelAuthorList()"
      >
        {{ $t('home.global.cancel') }}
      </button>
      <button
        style="background: #1f71e0; color: white"
        :disabled="!authorListDirty"
        @click="submitAuthorList()"
      >
        {{ $t('home.global.ok') }}
      </button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue'
import { Tooltip } from 'ant-design-vue'
import {
  CloseOutlined,
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons-vue'
import { BaseAuthorInfo } from 'go-sea-proto/gen/ts/doc/UserDoc'
import { getAuthors, updateAuthors } from '@/api/material'
import { goPathPage } from '@/common/src/utils/url'

import { useLibraryList } from '@/stores/library/list'

const EMPTY_ID = '0'
const isIdentify = (id: string) => id && id !== EMPTY_ID
const goAuthorPage = (id: string) => {
  if (isIdentify(id)) {
    goPathPage(`/author/${id}`)
  }
}

const storeLibraryList = useLibraryList()

const authorListBackup = ref<BaseAuthorInfo[]>([])
const authorList = ref<BaseAuthorInfo[]>([])
const authorListDirty = computed(() => {
  return !(
    authorListBackup.value.length === authorList.value.length &&
    authorListBackup.value.every((bak, index) => {
      const item = authorList.value[index]
      return bak.id === item.id && bak.name === item.name
    })
  )
})
const authorEditingKey = ref('')
const authorEditingName = ref('')
const getAuthorKey = (item: BaseAuthorInfo, index: number) => {
  return `${item.id}-${index}`
}
const addAuthor = () => {
  if (authorEditingKey.value) {
    return
  }

  const newAuthor: BaseAuthorInfo = {
    id: '',
    name: '',
    isAuthentication: false,
  }
  authorList.value.push(newAuthor)
  editAuthor(newAuthor, authorList.value.length - 1)
}
const editAuthor = async (item: BaseAuthorInfo, index: number) => {
  authorEditingName.value = item.name
  authorEditingKey.value = getAuthorKey(item, index)
  await nextTick()
  document
    .querySelector<HTMLInputElement>(
      `input[data-metadata-author="${authorEditingKey.value}"]`
    )
    ?.focus()
}
const enterAuthor = () => {
  if (!authorEditingKey.value) {
    return
  }

  const author = authorList.value.find(
    (item, index) => getAuthorKey(item, index) === authorEditingKey.value
  )!

  const { value } = authorEditingName
  cancelEdit()

  if (!value) {
    if (!author.id && !author.name) {
      deleteAuthor(author)
    }
  } else if (author.name !== value) {
    author.name = value
    author.id = ''
    author.isAuthentication = false
    authorList.value = [...authorList.value]
  }
}
const cancelEdit = () => {
  authorEditingKey.value = ''
  authorEditingName.value = ''
}

const deleteAuthor = (item: BaseAuthorInfo) => {
  authorList.value = authorList.value.filter((author) => author !== item)
}
const submitAuthorList = async () => {
  const newAuthor = await updateAuthors({
    docId: storeLibraryList.paperListAuthorEdit,
    authors: authorList.value.map((item) => {
      if (item.id) {
        return item
      }

      return {
        name: item.name,
      } as BaseAuthorInfo
    }),
  })

  storeLibraryList.authorRefreshOptionsList()

  if (!newAuthor) {
    return
  }

  const doc = storeLibraryList.paperListAll.find(
    (item) => item.docId === storeLibraryList.paperListAuthorEdit
  )!
  doc.displayAuthor.authorInfos = newAuthor.authors.map((item) => ({
    literal: item.name,
    isAuthentication: item.isAuthentication,
  }))
  doc.displayAuthor!.userEdited = newAuthor.rollbackEnable
  storeLibraryList.paperListAll = [...storeLibraryList.paperListAll]
  cancelAuthorList()
}

const cancelAuthorList = () => {
  storeLibraryList.paperListAuthorEdit = ''
}

onMounted(async () => {
  const response = await getAuthors({
    docId: storeLibraryList.paperListAuthorEdit,
  })
  authorListBackup.value = response.displayAuthor!.authors
  authorList.value = authorListBackup.value.map((item) => ({
    ...item,
  }))
  if (typeof document !== 'undefined') {
    document.body.addEventListener('click', cancelAuthorList)
  }
})
onUnmounted(() => {
  if (typeof document !== 'undefined') {
    document.body.removeEventListener('click', cancelAuthorList)
  }
})
</script>

<style lang="less" scoped>
.list-author-edit-container {
  width: 300px;
  position: absolute;
  top: 0;
  right: 0;
  bottom: 60px;
  display: flex;
  flex-direction: column;
  background: white;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.08);
  z-index: 2;
  padding: 16px;

  .list-author-edit-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 12px;
    border-bottom: 1px solid #e5e6eb;
    margin-bottom: 8px;

    .title {
      font-size: 14px;
      font-weight: 600;
      color: #1d2129;
    }

    .header-icons {
      display: flex;
      align-items: center;
      gap: 16px;
      .icon {
        font-size: 16px;
        color: #86909c;
        cursor: pointer;
      }
    }
  }

  .list-author-edit-body {
    flex: 1 1 100%;
    overflow-y: auto;

    .list-author-edit-caption {
      height: 38px;
      margin-top: 4px;
      margin-left: 16px;
      color: #4e5969;
      display: flex;
      align-items: center;
    }

    .list-author-edit-list {
      padding-left: 8px;

      .list-author-edit-add {
        padding-left: 8px;
        margin-top: 14px;
        padding-bottom: 8px;
        font-size: 12px;
        color: #1f71e0;
        cursor: pointer;

        > i {
          margin-right: 8px;
        }
      }

      .list-author-edit-item {
        height: 32px;
        display: flex;
        align-items: center;
        border-radius: 4px;
        padding: 0 8px;

        &:hover {
          background-color: #f0f2f5;
        }

        .author-display {
          display: flex !important;
          justify-content: space-between !important;
          align-items: center !important;
          width: 100% !important;
          position: relative !important;
        }

        .list-author-edit-view {
          flex: 1 !important;
          white-space: nowrap !important;
          text-overflow: ellipsis !important;
          overflow: hidden !important;
          min-width: 0 !important;
          padding-right: 50px !important;
        }

        .icon-wrap {
          display: flex !important;
          align-items: center !important;
          gap: 8px !important;
          flex-shrink: 0 !important;
          position: absolute !important;
          right: 0 !important;
          top: 50% !important;
          transform: translateY(-50%) !important;

          .anticon {
            cursor: pointer;
            font-size: 14px;
            color: #4e5969;
            padding: 4px;
            border-radius: 2px;
            transition: all 0.2s;

            &:hover {
              background-color: rgba(0, 0, 0, 0.04);
            }

            &.anticon-delete {
              color: #e66045;
            }
          }
        }

        .list-author-identify {
          text-decoration: underline;
          cursor: pointer;
        }

        input {
          flex: 1 1 auto;
          min-width: 0;
          height: 28px;
          border: 1px solid #1f71e0;
          outline: 0;
          border-radius: 4px;
          padding: 0 8px;
        }
      }
    }
  }

  .list-author-edit-footer {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 12px;
    padding: 8px 0 0 0;
    border-top: 1px solid #e5e6eb;
    margin-top: auto;
  }
}
</style>
