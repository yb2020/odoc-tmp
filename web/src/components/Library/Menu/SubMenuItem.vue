<template>
  <div>
    <SubMenu
      :key="menuInfo.key"
      v-bind="props"
      @titleClick="handleTitleClick"
      xxvon="$listeners"
    >
      <template #title>
        <span>
          <a v-if="detectDisable(menuInfo.key) || isProcessing" disabled
            ><span class="folder-name" :style="{ opacity: isProcessing ? 0.5 : 1 }">{{ menuInfo.title }}</span></a
          >
          <a v-else
            ><span class="folder-name">{{ menuInfo.title }}</span></a
          >
        </span>
      </template>
      <template v-for="item in menuInfo.children">
        <SubMenuItem
          v-if="!item.isLeaf"
          :key="item.key"
          :menu-info="item"
          :disabled="detectDisable(item.key)"
          :detect-disable="detectDisable"
          :move-or-copy="moveOrCopy"
          :new-folder-parent-key="newFolderParentKey"
          :create-folder="createFolder"
          :cancel-folder="cancelFolder"
          :submit-folder="submitFolder"
          :type="type"
          @lock="emit('lock', [menuInfo.key, item.key])"
        />
        <!-- @unlock="emit('unlock', [menuInfo.key, item.key])" -->
      </template>
      <Item v-if="newFolderParentKey !== menuInfo.key">
        <span
          class="folder-name blue-6-c"
          @click.prevent.stop="createFolder(menuInfo.key)"
        >
          <PlusOutlined />{{ $t('home.library.newFolder') }}
        </span>
      </Item>
      <Item v-else>
        <input
          v-model="newFolderName"
          :data-folder-key="menuInfo.key"
          class="blue-6-b"
          @focus="emit('lock', [menuInfo.key])"
          @click.stop
          @keypress.enter.stop="submit()"
          @keypress.esc.stop="cancelFolder()"
        />
        <!-- Windows原生输入法中文输入会触发blur -->
        <!-- @blur="emit('unlock', menuInfo.key)" -->
      </Item>
    </SubMenu>
  </div>
</template>
<script lang="ts" setup>
import { PlusOutlined } from '@ant-design/icons-vue'
import { Menu } from 'ant-design-vue'
import { ref, computed } from 'vue'
import SubMenuItem from './SubMenuItem.vue'
import { useI18n } from 'vue-i18n'
import { Key } from 'ant-design-vue/lib/_util/type'
import { useLibraryMenu } from '@/stores/library/menu'

const { Item, SubMenu } = Menu
const storeLibraryMenu = useLibraryMenu()

export interface MenuInfo {
  key: string
  title: string
  children: MenuInfo[]
  isLeaf: boolean
}

const $t = useI18n().t

const props = defineProps({
  type: {
    type: Number,
    required: true,
  },
  menuInfo: {
    type: Object,
    default: () => ({}),
  },
  moveOrCopy: {
    type: Function,
    default: () => {},
  },
  detectDisable: {
    type: Function,
    required: true,
  },
  newFolderParentKey: {
    type: String,
    default: '',
  },
  createFolder: {
    type: Function,
    default: () => {},
  },
  cancelFolder: {
    type: Function,
    default: () => {},
  },
  submitFolder: {
    type: Function,
    default: () => {},
  },
})

const emit = defineEmits<{
  (e: 'lock', k: Key[]): void
  (e: 'unlock', k: Key[]): void
}>()

const newFolderName = ref('')

// 根据操作类型判断是否正在处理中
const isProcessing = computed(() => {
  return props.type === 1 ? storeLibraryMenu.isMoving : storeLibraryMenu.isCopying
})

const handleTitleClick = () => {
  // 如果正在处理中，阻止点击
  if (isProcessing.value) {
    return
  }
  console.log('[DEBUG][SubMenuItem] 点击文件夹:', props.menuInfo.title, '键值:', props.menuInfo.key, 'type值:', props.type);
  props.moveOrCopy(props.type, { key: props.menuInfo.key });
}

const submit = async () => {
  await props.submitFolder(newFolderName.value)
  newFolderName.value = ''
}
</script>
<style lang="less" scoped>
.folder-name {
  display: inline-block;
  min-width: 120px;
}
input {
  outline: 0;
  border-style: solid;
  border-radius: 2px;
}
</style>

<style lang="less">
@import '@/common/assets/common-style.less';
</style>
