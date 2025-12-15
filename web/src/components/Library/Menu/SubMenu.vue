<template>
  <div>
    <Dropdown
      v-model:visible="visible"
      :trigger="['click']"
      :showAction="['mouseenter']"
      @visible-change="emit('visibleChange', $event)"
    >
      <div v-if="hasSlot">
        <slot></slot>
      </div>
      <div v-else class="ant-dropdown-link">
        {{ type === 1 ? $t('home.library.moveTo') : $t('home.library.copyTo') }}
      </div>
      <template #overlay>
        <Menu
          mode="vertical"
          selectable
          class="thin-scroll"
          :class="{
            'first-select': !hasSlot,
          }"
          :open-keys="openKeys"
          @update:open-keys="handleOpenKeys"
        >
          <SubMenuItem
            v-for="item in folderList"
            :key="`${type}-${item.key}`"
            :menu-info="item"
            :disabled="detectDisable(item.key)"
            :detect-disable="detectDisable"
            :move-or-copy="moveOrCopy"
            :new-folder-parent-key="newFolderParentKey"
            :create-folder="createFolder"
            :cancel-folder="cancelFolder"
            :submit-folder="submitFolder"
            :type="type"
            @lock="toggleLockKeys($event)"
            @unlock="toggleLockKeys($event, false)"
          />
          <Item
            v-if="newFolderParentKey !== allKey"
            class="readpaper-submenu-addfolder-0"
          >
            <span
              class="folder-name blue-6-c"
              @click.prevent.stop="createFolder(allKey)"
            >
              <PlusOutlined />{{ $t('home.library.newFolder') }}
            </span>
          </Item>
          <Item v-else class="readpaper-submenu-addfolder-1">
            <input
              v-model="newFolderName"
              :data-folder-key="allKey"
              class="blue-6-b"
              @focus="openKeys = []"
              @click.prevent.stop
              @keypress.enter.stop="submit()"
              @keypress.esc.stop="cancelFolder()"
            />
          </Item>
        </Menu>
      </template>
    </Dropdown>
  </div>
</template>

<script lang="ts" setup>
import { difference } from 'lodash-es'
import { computed, nextTick, ref, watch } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { Dropdown, Menu } from 'ant-design-vue'
import { delay } from '@idea/aiknowledge-special-util'
import SubMenuItem from './SubMenuItem.vue'
import {
  allKey,
  seperatorKey,
  LiteratureNode,
  useLibraryIndex,
  removePrefix,
  newFolderKey,
  noteDirty,
} from '@/stores/library'
import { addFolder } from '@/api/document'
import { useLibraryMenu } from '@/stores/library/menu'
import { Key } from 'ant-design-vue/lib/_util/type'

const { Item } = Menu

const props = defineProps({
  type: { type: Number, default: 0 },
  active: { type: Boolean, default: false },
  hasSlot: {
    type: Boolean,
    default: false,
  },
  detectDisable: {
    type: Function as unknown as () => (...args: any[]) => boolean,
    required: true,
  },
  moveOrCopy: {
    type: Function as unknown as () => (...args: any[]) => void | Promise<void>,
    required: true,
  },
})

const emit = defineEmits<{
  (x: 'visibleChange', visible: boolean): void
}>()

const visible = ref(false)
watch(
  () => props.active,
  (v) => {
    visible.value = v
  }
)

const openKeys = ref<Key[]>([])
const lockKeys = ref<Key[]>([])
const toggleLockKeys = (key: Key[], v?: boolean) => {
  lockKeys.value = v !== false ? key : []
}
const handleOpenKeys = (keys: Key[]) => {
  if (
    !difference(lockKeys.value, openKeys.value).length &&
    (!keys.length || !difference(keys, lockKeys.value).length)
  ) {
    openKeys.value = [...keys, ...lockKeys.value]
  } else {
    const lastKey = keys[keys.length - 1]
    openKeys.value = keys.filter((x) => `${lastKey}`.includes(`${x}`))
    lockKeys.value = []
  }
}

const storeLibraryIndex = useLibraryIndex()
const storeLibraryMenu = useLibraryMenu()

const recursive = (children: LiteratureNode[]): LiteratureNode[] => {
  return children
    .filter((node) => !node.isLeaf)
    .map((node) => {
      return {
        ...node,
        children: recursive(node.children),
      }
    })
}

const folderList = computed(() => {
  console.log('[DEBUG][SubMenu] type值:', props.type);
  
  // 获取原始列表
  const nodes = storeLibraryIndex.libraryIndexList.filter(
    (node) => node.key !== allKey && node.key !== seperatorKey
  );
  
  // 使用Map确保键的唯一性
  const uniqueMap = new Map();
  nodes.forEach(node => {
    if (!uniqueMap.has(node.key)) {
      uniqueMap.set(node.key, node);
    }
  });
  
  // 转换为数组并递归处理子节点
  const uniqueNodes = Array.from(uniqueMap.values());
  const result = recursive(uniqueNodes);
  
  console.log('[DEBUG][SubMenu] 文件夹列表(去重后):', JSON.stringify(result.map(item => ({ key: item.key, title: item.title }))));
  return result;
})

const newFolderParentKey = ref('')

const createFolder = async (parentKey: string) => {
  newFolderParentKey.value = parentKey

  await nextTick()
  await delay(300)
  document
    .querySelector<HTMLInputElement>(
      `input[data-folder-key="${newFolderParentKey.value}"]`
    )
    ?.focus()
}

const cancelFolder = () => {
  newFolderParentKey.value = ''
}

const submitNewFolder = async (parentKey: string, name: string) => {
  let children: LiteratureNode[]

  if (parentKey === allKey) {
    children = storeLibraryIndex.libraryIndexList
      .filter(({ isLeaf }) => !isLeaf)
      .filter(({ key }) => key !== allKey && key !== seperatorKey)
  } else {
    const family = storeLibraryMenu.getLiteratureNodeFamily(parentKey)

    children = family.node.children
      .filter(({ isLeaf }) => !isLeaf)
      .filter(({ key }) => key !== newFolderKey)
  }

  const req: Parameters<typeof addFolder>[0] = {
    name,
    parentId: removePrefix(parentKey),
    level: storeLibraryIndex.libraryIndexExtra[parentKey].path.length,
    sort: children.length,
    oldFolderItems: children.map(({ key }, sort) => ({
      id: removePrefix(key),
      sort,
    })),
  }

  await addFolder(req)
  await storeLibraryIndex.fetchLibraryIndex()
  noteDirty.is = true
}

const submitFolder = async (name: string) => {
  await submitNewFolder(newFolderParentKey.value, name)
  cancelFolder()
}

const newFolderName = ref('')
const submit = async () => {
  await submitFolder(newFolderName.value)
  newFolderName.value = ''
}
</script>
<style lang="less" scoped>
.first-select {
  position: absolute;
  margin-top: -30px;
  margin-left: 95px;
  max-height: 300px;
  overflow-y: auto;
}
.ant-dropdown-link::after {
  content: '>';
  margin-left: 25px;
}
</style>
<style>
.readpaper-submenu-addfolder-0 {
  position: relative;
  padding: 0 !important;
  .folder-name {
    display: inline-block;
    width: 100%;
    height: 100%;
    padding: 5px 12px;
  }

}
.readpaper-submenu-addfolder-1 {
  position: relative;
  display: flex;
  align-items: center;
  input {
    flex: 1 1 100%;
    outline: 0;
    border-style: solid;
    border-radius: 2px;
    color: #000; /* 添加黑色文本颜色 */
    background-color: #fff; /* 添加白色背景 */
    padding: 2px 5px; /* 添加内边距使文本更易读 */
  }
}
</style>
<style lang="less">
@import '@/common/assets/common-style.less';
</style>
