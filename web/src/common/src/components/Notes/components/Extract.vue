<template>
  <div class="word-wrap h-full flex flex-col">
    <NoteSearch :noteState="noteState" />
    <NoteTag :noteState="noteState" />
    <a-spin
      :spinning="noteLoading"
      class="h-full"
    >
      <div
        class="flex justify-center mx-4 overflow-x-hidden pt-4 extract-content-wrap"
      >
        <!-- 留出10px误差空间 -->
        <a-row
          :gutter="[20, 20]"
          class="flex-1"
          :style="{ maxWidth: column * 388 + 10 + 'px' }"
        >
          <a-col
            v-for="(item, index) in cardList"
            :key="index"
            flex="0 0 388px"
          >
            <ExtractCard
              v-for="data in item"
              :key="data.uuid"
              class="mb-5"
              :noteState="noteState"
              :data="data"
            />
          </a-col>
        </a-row>
      </div>
    </a-spin>
  </div>
</template>

<script lang="ts" setup>
import lodash from 'lodash';
import { onMounted, onUnmounted, ref, watch } from 'vue';
import NoteSearch from './NoteSearch.vue';
import NoteTag from './NoteTag.vue';
import ExtractCard from './ExtractCard.vue';
import { useNote } from '../useNote';
import { NoteAnnotation } from '@common/api/note';

const props = defineProps<{
  noteState: ReturnType<typeof useNote>;
}>();

const loadingHeight = ref<number>(Number.POSITIVE_INFINITY);
const column = ref<number>(2);
const cardList = ref<NoteAnnotation[][]>([[]]); // 各列的数据
const { noteList, notePageNumber, noteTotalPage, fetchNoteList, noteLoading } =
  props.noteState;

function equallyCard() {
  // 平分数据
  const listData = noteList.value;
  const cardListTemp: NoteAnnotation[][] = [
    ...Array.from({ length: column.value }, () => []),
  ];

  listData?.forEach((item: NoteAnnotation, index: number) => {
    const i = index % column.value;
    cardListTemp[i].push(item);
  });
  cardList.value = [...cardListTemp];
}

watch(
  () => [cardList.value],
  () => {
    setTimeout(() => {
      // 计算最小高度
      const ele = document.querySelectorAll('.extract-content-wrap')?.[0];
      const colEle = ele?.querySelectorAll('.ant-col');
      let height = Number.POSITIVE_INFINITY;
      colEle?.forEach((item: Element) => {
        const extractList = item?.querySelectorAll('.extract-card-wrap');
        if (extractList?.length > 0) {
          const e = extractList[extractList.length - 1];
          if (
            e &&
            height >
              (e as HTMLElement).offsetTop + (e as HTMLElement).offsetHeight
          ) {
            height =
              (e as HTMLElement).offsetTop + (e as HTMLElement).offsetHeight;
          }
        }
      });
      loadingHeight.value = height;
      if (
        loadingHeight.value < ele.clientHeight &&
        !noteLoading.value &&
        notePageNumber.value < noteTotalPage.value
      ) {
        // 如果当前底部留有空隙且还有下一页则加一页数据
        fetchNoteList(notePageNumber.value + 1);
      }
    }, 500);
  },
  { immediate: false, deep: true }
);
watch(
  () => [noteList.value, column.value],
  () => {
    equallyCard();
  },
  { immediate: false, deep: true }
);

const handleScroll = lodash.debounce(() => {
  const ele = document.querySelectorAll('.extract-content-wrap')?.[0];
  // 最小列总高度
  // 滚动距离
  // 可视高度
  if (
    loadingHeight.value - ele.clientHeight - ele.scrollTop <= 1 &&
    !noteLoading.value &&
    notePageNumber.value < noteTotalPage.value
  ) {
    // 向下翻页
    fetchNoteList(notePageNumber.value + 1);
  }
}, 300);

const countColumnNum = lodash.debounce((ele: NodeListOf<Element>) => {
  column.value = Math.floor(ele?.[0]?.clientWidth / 388);
}, 300);

onMounted(() => {
  const ele = document.querySelectorAll('.extract-content-wrap');
  countColumnNum(ele);
  window.addEventListener('resize', () => countColumnNum(ele), false);
  ele?.[0]?.addEventListener('scroll', handleScroll, false);
});
onUnmounted(() => {
  const ele = document.querySelectorAll('.extract-content-wrap');
  window.removeEventListener('resize', () => countColumnNum(ele), false);
  ele?.[0]?.removeEventListener('scroll', handleScroll, false);
});
</script>

<style lang="less" scoped>
.word-wrap {
  background: #f0f2f5;
  .extract-content-wrap {
    padding-bottom: 30px;
    box-sizing: border-box;
  }
  :deep(.ant-spin-nested-loading) {
    flex: 1;
    overflow: auto;
    .ant-spin-container {
      padding-bottom: 30px;
      box-sizing: border-box;
    }
  }
  .card {
    display: flex;
    flex-direction: row;
    justify-content: space-around;
    .card-item {
      // visibility: hidden;
      margin-bottom: 20px;
      // text-align: center;
      width: 388px;
      // border-radius: 16px;
    }
    .visibles {
      visibility: visible;
    }
  }
  ::-webkit-scrollbar {
    width: 6px; //修改垂直滚动条宽度
    height: 6px; //修改水平滚动条宽度
  }

  ::-webkit-scrollbar-thumb {
    border-radius: 10px;
    background: transparent;
  }

  &:hover {
    ::-webkit-scrollbar-thumb {
      background: rgba(201, 205, 212, 1);
      &:hover {
        background-color: rgba(170, 170, 170, 0.8);
      }
    }
  }
}
</style>
