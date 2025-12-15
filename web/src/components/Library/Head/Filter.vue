<template>
  <div flex="200px" class="title">
    <Dropdown
      v-model:visible="isFilterDropdownVisible"
      overlay-class-name="classify-header"
      :disabled="optionList.length === 0"
      :destroy-popup-on-hide="true"
    >
      <a class="ant-dropdown-link" @click.prevent>
        {{ text }}
        <span v-if="indeterminate"
          >{{ $t('home.library.selected')
          }}{{ checked.length + (searchInput ? 0 : +emptyChecked) }}</span
        >
        <DownOutlined
          v-if="optionList.length > 0"
          style="transition: 0.24s"
          :style="{
            transform: isFilterDropdownVisible ? 'rotate(-180deg)' : '',
          }"
        />
      </a>
      <template #overlay>
        <Menu class="classify-list">
          <Item>
            <Input
              v-if="useKeyword"
              v-model:value="keyword"
              :placeholder="$t('home.library.enterAuthName')"
            />
            <Checkbox
              v-else
              class="classify-checkbox"
              :checked="allChecked"
              :indeterminate="indeterminate"
              @change="$emit('allToggle')"
              >{{ $t('home.library.selectAll') }}</Checkbox
            >
          </Item>
          <Divider />
          <Item v-for="item in viewList" :key="item" class="classify-item">
            <Checkbox
              class="classify-checkbox"
              :checked="checked.includes(item)"
              @change="$emit('toggle', item)"
            >
              <span class="text">{{ item }}</span>
            </Checkbox>
          </Item>
          <Divider v-if="!searchInput" />
          <Item>
            <Checkbox
              v-if="!searchInput"
              class="classify-checkbox"
              :checked="emptyChecked"
              @change="$emit('emptyToggle')"
              >{{ $t('home.library.unmarked') }}{{ text }}</Checkbox
            >
          </Item>
        </Menu>
      </template>
    </Dropdown>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { DownOutlined } from '@ant-design/icons-vue'
import { Dropdown, Menu, Checkbox, Input } from 'ant-design-vue'

const { Item, Divider } = Menu

const props = defineProps({
  text: {
    type: String,
    required: true,
  },
  useKeyword: {
    type: Boolean,
    default: false,
  },
  searchInput: String,
  optionList: {
    type: Array as () => string[],
    required: true,
  },
  checked: {
    type: Array as () => string[],
    required: true,
  },
  emptyChecked: {
    type: Boolean,
    required: true,
  },
  allChecked: {
    type: Boolean,
    required: true,
  },
  indeterminate: {
    type: Boolean,
    required: true,
  },
})

defineEmits<{
  (event: 'toggle', value: string): void
  (event: 'allToggle'): void
  (event: 'emptyToggle'): void
}>()

const isFilterDropdownVisible = ref(false)
const keyword = ref('')
const viewList = computed(() => {
  if (!props.useKeyword) {
    return props.optionList
  }

  let count = 30
  const list = props.optionList.filter((option: any) => {
    if (props.checked.includes(option)) {
      return true
    }

    if (count === 0) {
      return false
    }

    if (!keyword.value || option.includes(keyword.value)) {
      count -= 1
      return true
    }

    return false
  })
  return list
})


</script>

<style lang="less">
.title {
  .ant-dropdown-link {
    color: var(--site-theme-text-primary) !important; /* 使用主题文本颜色替代蓝色 */
  }
}

.classify-header {
  box-shadow: var(--site-theme-shadow-2) !important;
  border: 1px solid var(--site-theme-border-color) !important;
  
  .ant-dropdown-menu {
    background-color: var(--site-theme-background) !important;
  }
  
  .ant-dropdown-menu-item {
    color: var(--site-theme-text-primary) !important;
    
    &:hover {
      background-color: var(--site-theme-background-hover) !important;
    }
  }
  
  .ant-checkbox-wrapper {
    color: var(--site-theme-text-primary) !important;
  }
  
  .ant-input {
    background-color: var(--site-theme-background) !important;
    border-color: var(--site-theme-border-color) !important;
    color: var(--site-theme-text-primary) !important;
    
    &::placeholder {
      color: var(--site-theme-placeholder-color) !important;
    }
  }
  
  .ant-divider {
    border-color: var(--site-theme-divider-light) !important;
    opacity: 0.6;
  }
}

.classify-list {
  overflow: auto;
  max-height: 400px;
  background-color: var(--site-theme-background) !important;
  
  .classify-item {
    position: relative;
    
    .text {
      color: var(--site-theme-text-primary);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      display: inline-block;
      max-width: 155px;
      vertical-align: bottom;
    }
  }
}

.classify-checkbox-delete {
  position: absolute;
  right: 5px;
  top: 4px;
  width: 20px;
  height: 20px;
  color: var(--site-theme-text-tertiary);
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  > * {
    font-size: 12px;
  }
}
</style>
