<template>
  <Layout
    :style="{
      maxHeight: 'calc(100vh - 40px)',
    }"
    :title="info.refIdx ? info.refIdx : `Page ${info.pageNum}`"
    :isDing="isDing"
    :tippyHandler="tippyHandler"
    group="figure"
    class="js-ignore-pdf-wheel-scale"
  >
    <PerfectScrollbar
      v-if="shown"
      ref="scrollRef"
      :options="{
        suppressScrollX: true,
        wheelPropagation: true,
        // handlers: ['keyboard', 'touch', 'drag-thumb', 'click-rail'],
      }"
      class="wrapper"
    >
      <div class="content">
        <div class="img-wrap">
          <img
            class="img"
            :src="info.url"
            :alt="info.desc"
            :title="$t('message.pressCtrlAndScrollToZoom')"
          >
          <div class="action">
            <!-- <div class="similar bg" @click="handleSimilar">
              <i
                class="aiknowledge-icon icon-picture-search"
              ></i>
              <div class="similar-text">{{ $t('viewer.similar') }}</div>
            </div> -->

            <a
              :href="info.url + '?attname='"
              download
            >
              <div class="download bg">
                <download-outlined style="color: #fff" />
              </div>
            </a>
          </div>
        </div>
        <p class="js-interact-drag-ignore">
          {{ info.desc }}
        </p>
      </div>
    </PerfectScrollbar>
  </Layout>
</template>
<script lang="ts" setup>
import { nextTick, onMounted, onUnmounted, ref } from 'vue';
import { PdfFigureAndTableInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import Layout from './Layout/index.vue';
import ZoomController from '~/src/dom/zoom';
// import { getImgMd5FromUrl } from '~/src/api/imgSearch';
// import { goPathPage } from '~/src/common/src/utils/url';
import {
  PageType,
  reportPopupNoteCiteFigureImpression,
} from '~/src/api/report';
import { currentNoteInfo } from '~/src/store';
import { DownloadOutlined } from '@ant-design/icons-vue';

const DEFAULT_SCALE_DELTA = 0.1;

const MIN_SCALE = 0.5;
const MAX_SCALE = 4.0;

const zoomIn = (ticks: number, evt: WheelEvent) => {
  const viewer = (evt.target as HTMLElement).closest(
    '.js-figure-tippy-viewer'
  ) as HTMLDivElement;
  let newScale = 1;
  do {
    // newScale = (newScale * DEFAULT_SCALE_DELTA).toFixed(2);
    newScale += DEFAULT_SCALE_DELTA;
    newScale = Math.round(newScale * 10) / 10;
    newScale = Math.min(MAX_SCALE, newScale);
  } while (--ticks > 0 && newScale < MAX_SCALE);
  const rect = viewer.getBoundingClientRect();
  if (rect) {
    viewer.style.width =
      Math.min(viewer.offsetWidth * newScale, window.innerWidth - rect.left) +
      'px';
    viewer.style.height =
      Math.min(viewer.offsetHeight * newScale, window.innerHeight - rect.top) +
      'px';
  } else {
    viewer.style.width =
      Math.min(viewer.offsetWidth * newScale, window.innerWidth) + 'px';
    viewer.style.height =
      Math.min(viewer.offsetHeight * newScale, window.innerHeight) + 'px';
  }
};

const zoomOut = (ticks: number, evt: WheelEvent) => {
  const viewer = (evt.target as HTMLElement).closest(
    '.js-figure-tippy-viewer'
  ) as HTMLDivElement;
  let newScale = 1;
  do {
    // newScale = (newScale / DEFAULT_SCALE_DELTA).toFixed(2);
    newScale -= DEFAULT_SCALE_DELTA;
    newScale = Math.round(newScale * 10) / 10;
    newScale = Math.max(MIN_SCALE, newScale);
  } while (--ticks > 0 && newScale > MIN_SCALE);
  viewer.style.width = Math.max(viewer.offsetWidth * newScale, 200) + 'px';
  viewer.style.height = Math.max(viewer.offsetHeight * newScale, 200) + 'px';
};

const props = defineProps<{
  info: PdfFigureAndTableInfo & { refContent: string };
  isDing: import('vue').Ref<boolean>; // 新增：响应式的钉住状态
  tippyHandler: (event: 'ding' | 'close' | 'unding' | 'lock') => void;
}>();

const shown = ref(false);

onMounted(() => {
  reportPopupNoteCiteFigureImpression({
    page_type: PageType.note,
    type_parameter: currentNoteInfo.value?.pdfId || '',
    popup_type: props.info.refIdx,
    popup_id: 'none',
    order_num:
      props.info?.refContent || props.info?.refIdx?.split('_')[1] || '',
  });

  document.addEventListener('wheel', zoomImage, {
    passive: false,
  });
  /**
   * 这里要设置shown的原因是，如果一开始就加载PerfectScrollbar，会导致滚动bar位于left的位置
   * 猜测是tippy在shown之前会导致这里PerfectScrollbar判断样式有问题
   */
  shown.value = true;
  nextTick(() => {
    updateScrollBar();
  });
});

onUnmounted(() => {
  document.removeEventListener('wheel', zoomImage);
});

const zoomController = new ZoomController({
  isValid: (evt) => {
    const dom = (evt.target as HTMLElement).closest(
      '.js-figure-tippy-viewer'
    ) as HTMLDivElement;
    return !!dom;
  },
  zoomIn,
  zoomOut,
});

const zoomImage = (evt: WheelEvent) => {
  scrollRef.value?.update();
  zoomController.onWheel(evt);
};

const scrollRef = ref();

const updateScrollBar = () => {
  const img = new Image();
  img.onload = function () {
    nextTick(() => {
      scrollRef.value?.update();
    });
  };
  img.onerror = function () {
    nextTick(() => {
      scrollRef.value?.update();
    });
  };
  img.src = props.info.url;
};

// const handleSimilar = async () => {
//   const res = await getImgMd5FromUrl({ picUrl: props.info.url });

//   reportElementClick({
//     page_type: PageType.note,
//     type_parameter: currentNoteInfo.value?.pdfId || '',
//     element_name: ElementClick.similar_picture,
//   });

//   goPathPage(`/img-search?token=${res.picMd5}`);
// };
</script>
<style lang="less" scoped>
.wrapper {
  max-height: calc(100vh - 72px);
  overflow-y: auto;
  height: 100%;
  .content {
    padding: 10px;

    p {
      margin-top: 10px;
      margin-bottom: 0;
    }

    .action {
      right: 0;
      bottom: 0;
      position: absolute;
      display: none;
      align-items: center;
      justify-content: center;
      .bg {
        background: rgba(0, 0, 0, 0.35);
        border-radius: 20px;
        cursor: pointer;
      }
    }

    .similar {
      padding: 3px 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-style: normal;
      font-weight: 400;
      font-size: 12px;
      line-height: 18px;
      color: rgba(255, 255, 255, 0.85);
      .similar-text {
        margin-left: 5px;
      }
      .icon-picture-search {
        font-size: 16px;
        height: 16px;
        line-height: 16px;
      }
    }

    .download {
      margin-left: 16px;
      line-height: 16px;
      color: #fff;
      padding: 4px;
    }

    .img-wrap {
      width: 100%;
      position: relative;
      &:hover {
        .action {
          display: flex;
        }
      }
    }
  }

  .img {
    width: 100%;
  }
}
</style>
