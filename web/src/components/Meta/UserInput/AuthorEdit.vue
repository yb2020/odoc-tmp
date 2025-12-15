<template>
  <div class="list-author-edit-container" @click.stop>
    <div class="list-author-edit-header">
      <div class="list-author-edit-title">
        {{ $t('meta.authorEditor.title') }}
      </div>
      <CloseOutlined @click="cancelAuthorList()" />
    </div>
    <div class="list-author-edit-body metadata-scroll">
      <!-- <div class="list-author-edit-caption">{{ $t('viewer.authorName') }}</div> -->
      <div class="list-author-edit-list">
        <div
          v-for="(item, index) in authorList"
          :key="getAuthorKey(item, index)"
          class="list-author-edit-item"
        >
          <div
            v-if="authorEditingKey !== getAuthorKey(item, index)"
            :title="item.literal"
            class="list-author-edit-view"
            :class="[
              isIdentify(item.id || '') && 'list-author-identify',
              item.isAuthentication
                ? 'list-author-authed'
                : 'list-author-not-authed',
            ]"
            @click="goAuthorPage(item.id || '')"
          >
            {{ item.literal }}
          </div>
          <EditOutlined
            v-if="authorEditingKey !== getAuthorKey(item, index)"
            @click="editAuthor(item, index)"
          />
          <DeleteOutlined
            v-if="authorEditingKey !== getAuthorKey(item, index)"
            @click="deleteAuthor(item)"
          />
          <input
            v-if="authorEditingKey === getAuthorKey(item, index)"
            v-model.trim="authorEditingName"
            autofocus
            :data-metadata-author="getAuthorKey(item, index)"
            @blur="enterAuthor()"
            @keyup.enter="enterAuthor()"
          />
        </div>
        <div class="list-author-edit-add" @click="addAuthor()">
          <PlusOutlined style="font-size: 10px"></PlusOutlined>
          {{ $t('meta.add') }}
        </div>
      </div>
    </div>
    <div class="list-author-edit-footer">
      <button @click="cancelAuthorList()">{{ $t('meta.cancel') }}</button>
      <button :disabled="!authorListDirty" @click="submitAuthorList()">
        {{ $t('meta.confirm') }}
      </button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, nextTick, watch } from 'vue'
import {
  CloseOutlined,
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'
import { goPathPage } from '@common/utils/url'
import { getDomainOrigin } from '@common/utils/env'
import { AuthorInfo } from 'go-sea-proto/gen/ts/doc/CSL'

const props = defineProps<{
  withAuthorLink: boolean
  authorListBackup: AuthorInfo[]
}>()

const emit = defineEmits<{
  (event: 'cancel'): void
  (event: 'update', newAuthor: AuthorInfo[]): void
}>()

const EMPTY_ID = '0'
const isIdentify = (id: string) => id && id !== EMPTY_ID && props.withAuthorLink
const goAuthorPage = (id: string) => {
  if (isIdentify(id)) {
    goPathPage(`${getDomainOrigin()}/author/${id}`)
  }
}

const getAuthorList = () =>
  props.authorListBackup.map((item) => ({
    ...item,
  }))
const authorList = ref<AuthorInfo[]>(getAuthorList())
watch(
  () => props.authorListBackup,
  () => {
    authorList.value = getAuthorList()
  }
)

const authorListDirty = computed(() => {
  return !(
    props.authorListBackup.length === authorList.value.length &&
    props.authorListBackup.every((bak, index) => {
      const item = authorList.value[index]
      return bak.id === item.id && bak.literal === item.literal
    })
  )
})
const authorEditingKey = ref('')
const authorEditingName = ref('')
const getAuthorKey = (item: AuthorInfo, index: number) => {
  return `${item.id}-${index}`
}

const addAuthor = () => {
  if (authorEditingKey.value) {
    return
  }

  const newAuthor: AuthorInfo = {
    id: '',
    literal: '',
    isAuthentication: false,
  }
  authorList.value.push(newAuthor)
  editAuthor(newAuthor, authorList.value.length - 1)
  console.log(authorList.value.length)
}

const editAuthor = async (item: AuthorInfo, index: number) => {
  authorEditingName.value = item.literal
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
    if (!author.id && !author.literal) {
      deleteAuthor(author)
    }
  } else if (author.literal !== value) {
    author.literal = value
    author.id = ''
    author.isAuthentication = false
    authorList.value = [...authorList.value]
  }
}
const cancelEdit = () => {
  authorEditingKey.value = ''
  authorEditingName.value = ''
}

const deleteAuthor = (item: AuthorInfo) => {
  authorList.value = authorList.value.filter((author) => author !== item)
}
const submitAuthorList = async () => {
  emit('update', authorList.value)
  cancelAuthorList()
}

const cancelAuthorList = () => {
  emit('cancel')
}
</script>

<style lang="less">
.list-author-edit-container {
  width: 200px;
  color: #1d2229;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  box-shadow: -1px 0 3px rgba(0, 0, 0, 0.12);
  z-index: 2;
  .list-author-edit-header {
    flex: 0 0 48px;
    padding: 0 16px;
    display: flex;
    align-items: center;
    border-bottom: 1px solid #e5e6eb;
    .list-author-edit-title {
      color: #1d2229;
      font-weight: bold;
      flex: 1;
    }
  }
  .list-author-edit-footer {
    flex: 0 0 48px;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    padding-right: 16px;
    border-top: 1px solid #e5e6eb;
    button {
      margin-left: 12px;
      width: 64px;
      height: 24px;
      border: 0;
      outline: 0;
      border-radius: 2px;
      font-size: 12px;
      cursor: pointer;
      &:disabled {
        cursor: not-allowed;
        opacity: 0.7;
      }
      &:first-of-type {
        background: #f0f2f5;
        color: #4e5969;
      }
      &:last-of-type {
        background: #1f71e0;
        color: #fff;
      }
    }
  }
  .list-author-edit-body {
    flex: 1 1 100%;
    max-height: 500px;
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
      padding-right: 8px;
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
        justify-content: flex-end;
        align-items: center;
        border-radius: 2px;
        position: relative;
        overflow: hidden;
        > span {
          display: none;
          margin-right: 14px;
          cursor: pointer;
          z-index: 3;
          &:first-of-type {
            color: #4e5969;
          }
          &:last-of-type {
            color: #e66045;
          }
        }
        &:hover {
          background-color: #f0f2f5;
          > span {
            display: initial !important;
          }
          .list-author-edit-view {
            width: 70%;
          }
        }
        .list-author-edit-view {
          padding-left: 8px;
          position: absolute;
          left: 0;
          width: 100%;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
        .list-author-identify {
          text-decoration: underline;
          cursor: pointer;
        }

        .list-author-authed {
          color: #1f71e0;
        }
        .list-author-not-authed {
          color: #4e5969;
        }

        input {
          padding-left: 8px;
          position: absolute;
          left: 0;
          background-color: #ffffff;
          border: 1px solid #1f71e0;
          outline: 0;
          height: 100%;
          width: 100%;
          border-radius: 2px;
        }
      }
    }
  }
}

[data-theme='dark'] {
  .list-author-edit-container {
    background-color: #222326;
    .list-author-edit-header {
      border-bottom: 1px solid #ffffff26;
      .list-author-edit-title {
        color: #ffffffd9 !important;
      }
    }
    .list-author-edit-body {
      .list-author-edit-list {
        .list-author-edit-item {
          &:hover {
            background-color: #ffffff14 !important;
          }
        }
      }
    }
    .list-author-edit-footer {
      border-top: 1px solid #ffffff26;
      button {
        &:first-of-type {
          background-color: #667180 !important;
        }
      }
    }
  }
}
</style>
