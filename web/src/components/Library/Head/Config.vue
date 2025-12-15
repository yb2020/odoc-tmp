<template>
  <div class="literature-paper-list-config dot">
    <Dropdown
      v-model:visible="dropdownVisible"
      overlay-class-name="dropdownMenu"
      placement="bottomRight"
    >
      <a class="ant-dropdown-link dropdownAction" @click.prevent>
        <MoreOutlined class="iconmore" />
      </a>
      <template #overlay>
        <Menu>
          <Item
            v-for="(item, index) in storeLibraryList.paperHeadList"
            :key="item.key"
            draggable="true"
            @dragstart="dragStart(index)"
            @dragend="dragEnd()"
            @dragover.prevent
            @dragenter="dragEnter(index)"
            @drop="drop()"
          >
            <div
              class="config-item"
              :style="
                dragging === index && dragStyled
                  ? {
                      visibility: 'hidden',
                    }
                  : {}
              "
            >
              <img class="config-item-icon" :src="HOLDER_OUTLINED" alt="" />
              <div class="config-item-name">
                {{ $t(TableHeadName[item.key]) }}
              </div>
              <CheckOutlined
                v-if="item.key !== 'docName'"
                style="cursor: pointer"
                :style="{
                  color: item.visible ? '#4E5969' : '#c9cdd4',
                }"
                @click="toggleHead(index)"
              />
              <div
                v-if="index === 0"
                v-show="dragging > index && dragHovering === index"
                class="config-item-insert"
                style="top: -5px"
              ></div>
              <div
                v-show="
                  (dragging < index && dragHovering === index) ||
                  (dragging > index &&
                    dragHovering === index + 1 &&
                    dragging !== dragHovering)
                "
                class="config-item-insert"
                style="bottom: -5px"
              ></div>
            </div>
          </Item>
        </Menu>
      </template>
    </Dropdown>
  </div>
</template>

<script lang="ts" setup>
import { ref, nextTick } from 'vue'
import { Menu, Dropdown } from 'ant-design-vue'
import { CheckOutlined, MoreOutlined } from '@ant-design/icons-vue'

import {
  useLibraryList,
  TableHeadName,
  TableHeadDetail,
} from '../../../stores/library/list'
import HOLDER_OUTLINED from './holder-outlined.svg'

const { Item } = Menu

const storeLibraryList = useLibraryList()
const dropdownVisible = ref(false)

const dragging = ref(NaN)
const dragHovering = ref(NaN)
const dragStyled = ref(false)

const dragStart = async (index: number) => {
  dragging.value = index
  await nextTick()
  await new Promise((resolve) => setTimeout(resolve, 100))
  dragStyled.value = true
}
const dragEnd = () => {
  dragging.value = NaN
  dragStyled.value = false
  dragHovering.value = NaN
}
const dragEnter = (index: number) => {
  dragHovering.value = index
}
const drop = () => {
  if (isNaN(dragging.value) || dragging.value === dragHovering.value) {
    return
  }

  const toIndex =
    dragging.value < dragHovering.value
      ? dragHovering.value + 1
      : dragHovering.value

  const item = storeLibraryList.paperHeadList[dragging.value]
  const list: (TableHeadDetail | null)[] = [...storeLibraryList.paperHeadList]
  list[dragging.value] = null
  list.splice(toIndex, 0, item)
  storeLibraryList.paperHeadList = list.filter(Boolean) as TableHeadDetail[]
  dragEnd()
}

const toggleHead = (index: number) => {
  storeLibraryList.paperHeadList.splice(index, 1, {
    ...storeLibraryList.paperHeadList[index],
    visible: !storeLibraryList.paperHeadList[index].visible,
  })
}
</script>

<style lang="less">
.config-item {
  position: relative;
  width: 104px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  cursor: grab;
  .config-item-icon {
    width: 12px;
    flex: 0 0 12px;
  }
  .config-item-name {
    margin-left: 8px;
    flex: 1 1 100%;
  }
}

.config-item-insert {
  position: absolute;
  width: 100%;
  background-color: skyblue;
  height: 2px;
}
</style>
