<template>
  <div
    class="collect-panel js-collect-panel"
    @click="(e) => stopPropagation && e.stopPropagation()"
  >
    <div style="display: flex; justify-content: space-between" class="title">
      <span
        >{{ title }}&nbsp;
        {{
          addedTags && addedTags.length > 0 ? `(${addedTags.length})` : ''
        }}</span
      >
      <CloseOutlined v-if="!fromAnnotate" @click="handleClose" />
    </div>
    <Tag
      v-for="item in tags"
      :key="item.folderId"
      style="cursor: pointer"
      :class="{ checked: item.isContainCurrentPaper }"
      @click="handleCollectBtn(item)"
      >{{ item.name }}
    </Tag>
    <div>
      <span
        v-if="!addtag"
        class="add"
        style="cursor: pointer"
        @click="addothertag"
        >+ &nbsp;{{ newText }}</span
      >
      <Input
        v-else
        ref="input"
        v-model:value.trim="newTag"
        class="input"
        autofocus
        type="text"
        @pressEnter="handleCreateSuccess"
        @blur="cancel"
      />
    </div>
  </div>
</template>
<script lang="ts">
import { computed, defineComponent, nextTick, PropType, ref } from 'vue'
import { Input, Tag } from 'ant-design-vue'
import _ from 'lodash'
import { SimpleFolderInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc'
import {
  getTopLevelFolderList,
  addFolder,
  attachPaperToFolder,
  removeDocFromFolder,
} from '@/api/document'
import { CloseOutlined } from '@ant-design/icons-vue'
export interface ToggleClassifyEmitData {
  type: 'remove' | 'add'
  classifyItems: SimpleFolderInfo[]
}

export default defineComponent({
  components: {
    Input,
    CloseOutlined,
    Tag,
  },
  props: {
    title: {
      type: String,
      default: '',
    },
    newText: {
      type: String,
      default: '',
    },
    paperId: {
      type: String,
      default: '',
    },
    onClose: {
      props: Function as PropType<() => void>,
      default: null,
    },
    fromAnnotate: {
      type: Boolean,
      default: false,
    },
    stopPropagation: Boolean,
  },
  setup(props, { emit }) {
    const addtag = ref(false)
    const newTag = ref('')
    const input = ref<HTMLElement>()
    const tags = ref<SimpleFolderInfo[]>([])
    const fetchClassifyData = async () => {
      const res = await getTopLevelFolderList({
        paperId: props.paperId,
      })

      tags.value = (res.data || []).filter((item) => _.trim(item.name))
    }
    const addedTags = computed(() =>
      tags.value?.filter((item) => item.isContainCurrentPaper)
    )
    const addothertag = () => {
      addtag.value = true
      nextTick(() => {
        input.value?.focus()
      })
    }
    const handleCollectBtn = async (data: SimpleFolderInfo) => {
      try {
        if (data.isContainCurrentPaper) {
          await removeDocFromFolder({
            removedDocItems: [{ docId: data.docId, folderId: data.folderId }],
            isHierarchicallyRemove: false,
          })
          data.isContainCurrentPaper = false
        } else {
          await attachPaperToFolder({
            paperId: props.paperId,
            folderId: data.folderId,
          })
          data.isContainCurrentPaper = true
        }

        fetchClassifyData()
      } catch (error) {}
    }
    const cancel = () => {
      addtag.value = false
      newTag.value = ''
    }
    const handleCreateSuccess = async () => {
      try {
        addtag.value = false
        if (!newTag.value) return
        await addFolder({
          parentId: '0',
          name: newTag.value,
          level: 0,
          sort: tags.value.length,
          oldFolderItems: (tags.value || []).map((item, index) => {
            return {
              id: item.folderId,
              sort: index,
            }
          }),
        })

        newTag.value = ''
        fetchClassifyData()
      } catch (error) {}
    }

    const handleClose = () => {
      emit('close')
      if (props.onClose) {
        ;(props.onClose as any)()
      }
    }

    if (props.fromAnnotate) {
      fetchClassifyData()
    }

    return {
      handleCollectBtn,
      handleCreateSuccess,
      addtag,
      newTag,
      tags,
      addothertag,
      input,
      cancel,
      fetchClassifyData,
      handleClose,
      addedTags,
    }
  },
})
</script>
<style lang="less" scoped>
.collect-panel {
  width: 400px;
  max-height: 300px;
  overflow: auto;
  background: #fff;
  border-radius: 4px 4px 0 0;
  padding: 16px;
  .title {
    margin-bottom: 16px;
    font-size: 13px;
    font-weight: 400;
    color: #73716f;
    line-height: 22px;
  }
  .add {
    font-size: 12px;
    font-weight: 400;
    color: #1f71e0;
    line-height: 16px;
    margin-top: 8px;
  }
  .checked {
    position: relative;
    border-color: #1f71e0 !important;
    color: #1f71e0 !important;
    &::after {
      content: '';
      position: absolute;
      right: -1px;
      bottom: -1px;
      border: 5px solid #1f71e0;
      border-color: transparent #1f71e0 #1f71e0 transparent;
    }
  }
  .input {
    width: 200px;
  }
  /deep/span.ant-tag {
    padding: 8px 10px;
    border-radius: 2px;
    border: 1px solid #dce0e6;
    font-size: 12px;
    font-weight: 400;
    color: #73716f;
    line-height: 16px;
    margin: 0 10px 10px 0;
  }
}
</style>
