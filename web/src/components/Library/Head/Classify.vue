<template>
  <div flex="200px" class="title">
    <Dropdown
      v-model:visible="isClassifyDropdownVisible"
      overlay-class-name="classify-header"
      :disabled="storeLibraryList.paperListClassifyList.length === 0"
    >
      <a class="ant-dropdown-link" @click.prevent>
        {{ $t('home.library.tag') }}
        <span v-if="storeLibraryList.paperListClassifyIndeterminate"
          >{{ $t('home.library.selected')
          }}{{
            storeLibraryList.paperListClassifyChecked.length +
            +(
              storeLibraryList.searchInput &&
              storeLibraryList.paperListClassifyEmptyChecked
            )
          }}</span
        >
        <DownOutlined
          v-if="storeLibraryList.paperListClassifyList.length > 0"
          style="transition: 0.24s"
          :style="{
            transform: isClassifyDropdownVisible ? 'rotate(-180deg)' : '',
          }"
        />
      </a>
      <template #overlay>
        <Menu class="classify-list">
          <Item>
            <Checkbox
              class="classify-checkbox"
              :checked="storeLibraryList.paperListClassifyAllChecked"
              :indeterminate="storeLibraryList.paperListClassifyIndeterminate"
              @change="paperListClassifyAllToggle"
              >{{ $t('home.library.selectAll') }}</Checkbox
            >
          </Item>
          <Divider />
          <Item
            v-for="classify in storeLibraryList.paperListClassifyList"
            :key="classify.id"
            class="classify-item"
          >
            <Checkbox
              class="classify-checkbox"
              :checked="
                storeLibraryList.paperListClassifyChecked.includes(classify.id)
              "
              @change="paperListClassifyToggle(classify.id)"
            >
              <Tag :closable="false">
                {{ classify.name }}
              </Tag>
              <div
                class="classify-checkbox-delete"
                @click.stop.prevent="removeClassify(classify.id, classify.name)"
              >
                <CloseOutlined />
              </div>
            </Checkbox>
          </Item>
          <Divider v-if="!storeLibraryList.searchInput" />
          <Item>
            <Checkbox
              v-if="!storeLibraryList.searchInput"
              class="classify-checkbox"
              :checked="storeLibraryList.paperListClassifyEmptyChecked"
              @change="paperListClassifyEmptyToggle()"
              >{{ $t('home.library.noTag') }}</Checkbox
            >
          </Item>
        </Menu>
      </template>
    </Dropdown>
  </div>
</template>

<script lang="ts" setup>
import { message, Modal, Dropdown, Menu, Checkbox } from 'ant-design-vue'
import {
  DownOutlined,
  CloseOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons-vue'
import { ref, h } from 'vue'
import { useLibraryList } from '@/stores/library/list'
import { useClassify } from '@/stores/classify'
import Tag from './Tag.vue'
import { deleteClassify } from '@/api/document'

const { Item, Divider } = Menu

const storeLibraryList = useLibraryList()
const storeClassify = useClassify()

const {
  paperListClassifyToggle,
  paperListClassifyAllToggle,
  paperListClassifyEmptyToggle,
  paperListClassifyRefresh,
  getFilesByFolderId,
} = storeLibraryList

const isClassifyDropdownVisible = ref(false)

const removeClassify = (classifyId: string, classifyName: string) => {
  Modal.confirm({
    title: `确定彻底删除标签 [${classifyName}] 并从所有文献中移除该标签？`,
    content: '',
    icon: h(InfoCircleOutlined),
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      try {
        await deleteClassify({
          classifyId,
        })
        storeClassify.refreshClassifyList()
        message.success('删除成功！')

        await paperListClassifyRefresh()
        getFilesByFolderId()
      } catch (err) {
        message.error('删除失败！请稍后再试')
      }
    },
    onCancel() {},
  })
}
</script>

<style lang="less">
.classify-header {
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.1), 0 3px 6px 0 rgba(0, 0, 0, 0.12);
  border: 1px solid #dfe6f0;
}
.classify-list {
  overflow: auto;
  max-height: 400px;
  .classify-item {
    position: relative;
  }
}
.classify-checkbox-delete {
  position: absolute;
  right: 5px;
  top: 4px;
  width: 20px;
  height: 20px;
  color: #979797;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  > * {
    font-size: 12px;
  }
}
</style>
