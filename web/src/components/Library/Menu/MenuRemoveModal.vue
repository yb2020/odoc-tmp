<template>
  <Modal
    v-model:visible="visible"
    :ok-text="
      !isFolder && parent && !removeFromAllFolder
        ? $t('home.global.remove')
        : $t('home.global.delete')
    "
    :cancel-text="$t('home.global.cancel')"
    width="480px"
    icon="info-circle"
    class="ant-modal-confirm ant-modal-confirm-confirm move-delete"
    :footer="null"
    :closable="false"
    @ok="onOk"
    @cancel="onCancel"
  >
    <div class="ant-modal-confirm-body-wrapper">
      <div class="ant-modal-confirm-body">
        <InfoCircleOutlined />
        <span v-if="isFolder" class="ant-modal-confirm-title">
          {{ $t('home.library.delFolderComfirm')
          }}<span style="color: #1f71e0">{{ itemList[0].title }}</span
          >{{ $t('home.library.andFolderLiterature') }}
        </span>
        <template v-else>
          <template v-if="itemList.length === 1">
            <span v-if="parent" class="ant-modal-confirm-title">
              {{ $t('home.library.sure')
              }}<span style="color: #1f71e0">{{ itemList[0].title }}</span
              >{{
                $t('home.library.removeFromFolder', {
                  title: parent.title,
                })
              }}
            </span>
            <span v-else class="ant-modal-confirm-title">
              {{ $t('home.library.delCompleteConfirm1')
              }}<span style="color: #1f71e0">{{ itemList[0].title }}</span
              >{{ $t('home.library.delCompleteConfirm2') }}
            </span>
          </template>
          <template v-else>
            <span v-if="parent" class="ant-modal-confirm-title">
              {{ $t('home.library.sure')
              }}<span style="color: #1f71e0">{{
                $t('home.library.literatureNum', {
                  num: itemList.length,
                })
              }}</span
              >{{
                $t('home.library.removeFromFolder', {
                  title: parent.title,
                })
              }}
            </span>
            <span v-else class="ant-modal-confirm-title">
              {{ $t('home.library.delCompleteConfirm1')
              }}<span style="color: #1f71e0">{{
                $t('home.library.literatureNum', {
                  num: itemList.length,
                })
              }}</span
              >{{ $t('home.library.delCompleteConfirm2') }}
            </span>
          </template>
        </template>

        <div class="ant-modal-confirm-content">
          <p
            v-if="hasAttachment && (!parent || removeFromAllFolder)"
            class="text-rp-neutral-8"
          >
            {{ $t('home.library.delTogether') }}
          </p>
          <div
            v-if="!isFolder && parent"
            :style="{
              color: '#73716f',
              width: '175px',
              position: 'absolute',
              marginTop: '30px',
              marginLeft: '-40px',
            }"
          >
            <input
              id="checkbox"
              v-model="removeFromAllFolder"
              type="checkbox"
              style="width: 16px; height: 16px"
            /><span
              style="position: absolute; margin-top: -3px; margin-left: 5px"
              >{{ $t('home.library.removeComplete') }}</span
            >
          </div>
        </div>
      </div>

      <div class="ant-modal-confirm-btns">
        <button type="button" class="ant-btn" @click="onCancel">
          <span>{{ $t('home.global.cancel') }}</span></button
        ><button type="button" class="ant-btn ant-btn-primary" @click="onOk">
          <span>{{
            !isFolder && parent && !removeFromAllFolder
              ? $t('home.global.remove')
              : $t('home.global.delete')
          }}</span>
        </button>
      </div>
    </div>
  </Modal>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { InfoCircleOutlined } from '@ant-design/icons-vue'
import { removePrefix } from '@/stores/library'
import {
  deleteDoc,
  deleteFolder,
  removeDocFromFolder,
} from '@/api/document'
import { useI18n } from 'vue-i18n'

export interface RemoveItem {
  title: string
  key: string
}

const props = defineProps({
  itemList: {
    type: Array as () => RemoveItem[],
    required: true,
  },
  remove: {
    type: Boolean,
    default: false,
  },
  parent: {
    type: Object as () => RemoveItem | null,
    default: null,
  },
  isFolder: {
    type: Boolean,
    default: false,
  },
  hasAttachment: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits<{
  (event: 'close'): void
  (event: 'refresh'): void
  (event: 'dirty'): void
}>()

const $t = useI18n().t

const visible = ref(true)
const removeFromAllFolder = ref(false)

const onOk = async () => {
  visible.value = false

  emit('dirty')

  const idList = props.itemList.map((item) => removePrefix(item.key))
  if (props.isFolder) {
    await deleteFolder({
      folderIds: idList,
    })
    message.success($t('home.global.deleteSucc') as string)
  } else if (!props.parent || removeFromAllFolder.value) {
    await deleteDoc({
      docIds: idList,
    })
    message.success($t('home.global.deleteSucc') as string)
  } else {
    await removeDocFromFolder({
      removedDocItems: idList.map((docId) => ({
        docId,
        folderId: removePrefix(props.parent!.key),
      })),
      isHierarchicallyRemove: props.remove,
    })
    message.success($t('home.global.removeSucc') as string)
  }

  emit('refresh')
  clear()
}

const onCancel = () => {
  visible.value = false
  setTimeout(clear, 200)
}

const clear = () => {
  emit('close')
  removeFromAllFolder.value = false
}
</script>
<style>
.move-delete .ant-modal-body .ant-modal-confirm-body-wrapper .ant-btn-primary {
  color: #fff !important;
  background-color: #e66045;
  border-color: #e66045 !important;
}
.move-delete .ant-modal-confirm-title {
  line-height: 24px;
}
.move-delete .anticon-info-circle {
  font-size: 16px;
  line-height: 27px;
}
.move-delete .ant-modal-body {
  padding: 32px 40px 24px 40px;
}
</style>
