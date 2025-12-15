// import {
  // ReferenceMarker,
  // FigureAndTableReferenceMarker,
  // PdfFigureAndTableInfo,
  // PdfBBox,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import { ReferenceMarker, FigureAndTableReferenceMarker, PdfFigureAndTableInfo, PdfBBox } from 'go-sea-proto/gen/ts/pdf/PdfParse';
import { message } from 'ant-design-vue';
import { ViewerEvent, ViewerController } from '@idea/pdf-annotate-viewer';
import {
  firstFigureAndTableFetch,
  firstReferenceMarkersFetch,
  getFiguresAndTablesFinal,
  getReferenceMarkersFinal,
} from '../api/parse';
// import MarkersVue from '../components/Markers/index.vue'
import { useStore } from '../store';
import {
  createReferenceTippyVue,
  destroyTippyVues,
  TippyVueItem,
} from './tippy';
import { ParseMutationTypes } from '../store/parse';
import {
  FigureTippyTriggerItem,
  PdfFigureAndTableInfosType,
} from '../store/parse/type';
// import { useI18n } from 'vue-i18n';
import i18n from '../locals/i18n';

export const scaleMarker = (
  item: ReferenceMarker | FigureAndTableReferenceMarker,
  viewportWidth: number
): Omit<PdfBBox, 'originWidth' | 'originHeight'> | null => {
  if (!item.bbox) {
    return null;
  }
  const scale = viewportWidth / item.bbox.originWidth;

  return {
    x0: scale * item.bbox.x0,
    y0: scale * item.bbox.y0,
    x1: scale * item.bbox.x1,
    y1: scale * item.bbox.y1,
  };
};

// function isReferenceMarker(marker: ReferenceMarker | FigureAndTableReferenceMarker): marker is ReferenceMarker {
//   return !!(marker as ReferenceMarker).paperId
// }

const appendDivs = (
  {
    referenceMarkers,
    figureAndTableMarkers,
    pageDiv,
    viewportWidth,
    pdfId,
  }: {
    referenceMarkers: ReferenceMarker[];
    figureAndTableMarkers: FigureAndTableReferenceMarker[];
    pageDiv: HTMLDivElement;
    viewportWidth: number;
    pdfId: string;
  },
  cacheReferenceTippyVueItems: TippyVueItem[]
) => {
  if (!referenceMarkers.length && !figureAndTableMarkers.length) {
    console.warn(`no reference and figure markers for div`, pageDiv);
    return;
  }
  let wrapper = pageDiv.querySelector('.idea-marker-layer') as HTMLDivElement;
  if (wrapper) {
    pageDiv.removeChild(wrapper);
  }

  wrapper = document.createElement('div');
  wrapper.classList.add('idea-marker-layer');
  pageDiv.appendChild(wrapper);

  const fragment = document.createDocumentFragment();
  [
    { uuid: 'reference', markers: referenceMarkers },
    { uuid: 'figure', markers: figureAndTableMarkers },
  ].forEach((items) => {
    if (!items.markers.length) {
      console.warn(`no ${items.uuid} markers for div`, pageDiv);
    } else {
      items.markers.forEach((item, i) => {
        const marker = scaleMarker(item, viewportWidth);
        if (marker) {
          const span = document.createElement('span');
          span.setAttribute('data-marker-type', items.uuid);
          span.classList.add('marker');
          span.style.top = marker.y0 + 'px';
          span.style.left = marker.x0 + 'px';
          span.style.width = marker.x1 - marker.x0 + 'px';
          span.style.height = marker.y1 - marker.y0 + 'px';
          span.setAttribute('data-index', `${i}`);
          fragment.appendChild(span);
        }
      });
    }
  });

  wrapper.appendChild(fragment);

  wrapper.addEventListener('click', (evt) => {
    const span = (evt.target as HTMLElement).closest('.marker');
    if (span) {
      const type = span.getAttribute('data-marker-type') as
        | 'reference'
        | 'figure';
      const index = parseInt(span.getAttribute('data-index') || '');

      if (type === 'reference') {
        const marker = referenceMarkers[index];
        if (!marker) {
          return;
        }
        const prevTippyInstance = cacheReferenceTippyVueItems?.[index];
        if (prevTippyInstance) {
          prevTippyInstance.tippy.show();
          return;
        }
        cacheReferenceTippyVueItems[index] = createReferenceTippyVue({
          marker,
          triggerEle: span,
        });
      } else {
        const marker = figureAndTableMarkers[index];
        if (!marker) {
          return;
        }
        const store = useStore();
        const payload: FigureTippyTriggerItem = {
          triggerEle: span,
          id: marker.refIdx,
          pdfId,
        };
        store.commit(
          `parse/${ParseMutationTypes.SET_FIGURE_VIEWER_ITEM}`,
          payload
        );
      }
    }
  });
};

export default function buildMarkers(
  pdfViewer: ViewerController,
  pdfId: string
) {
  let referenceMarkers: ReferenceMarker[] = [];
  let figureAndTableMarkers: FigureAndTableReferenceMarker[] = [];

  // const { t } = useI18n();
  const t = i18n.global.t;


  const fetchParseData = async () => {
    const formatData = (
      res0: {
        referenceMarkers: ReferenceMarker[];
        figureAndTableMarkers: FigureAndTableReferenceMarker[];
      },
      res1: {
        list: PdfFigureAndTableInfo[];
        total: number;
      }
    ) => {
      const validFigureAndTableMarkers = res0.figureAndTableMarkers.filter(
        (item) => {
          return res1.list.find((info) => info.refIdx === item.refIdx);
        }
      );
      referenceMarkers = res0.referenceMarkers;
      figureAndTableMarkers = validFigureAndTableMarkers;
      const store = useStore();
      store.commit(
        `parse/${ParseMutationTypes.SET_PDF_FIGURE_AND_TABLE_INFOS}`,
        {
          pdfId,
          infos: {
            list: res1.list,
            list0: res0.figureAndTableMarkers,
          },
        } as {
          pdfId: string;
          infos: Partial<PdfFigureAndTableInfosType>;
        }
      );
    };
    try {
      const [res0, res1] = await Promise.all([
        firstReferenceMarkersFetch({ pdfId }),
        firstFigureAndTableFetch({
          pdfId,
          pageReq: {
            pageSize: 100,
            pageNum: 1,
          },
        }),
      ]);

      if (res0 && res1) {
        formatData(res0, res1);
        return;
      }
    } catch (error) {}
    const hideLoading = message.loading(t('message.parsingDocumentTip'), 0);
    try {
      const [finalRes0, finalRes1] = await Promise.all([
        getReferenceMarkersFinal({ pdfId }),
        getFiguresAndTablesFinal({
          pdfId,
          pageReq: {
            pageSize: 100,
            pageNum: 1,
          },
        }),
      ]);
      formatData(finalRes0, finalRes1);
      hideLoading();
      message.success(t('message.parsingDocumentFinishedTip'));
    } catch (error) {
      hideLoading();
      message.error(t('message.acceptToparseDocumentTip'));
      const store = useStore();
      store.commit(
        `parse/${ParseMutationTypes.SET_PDF_FIGURE_AND_TABLE_INFOS}`,
        {
          pdfId,
          infos: { error },
        }
      );
    }
  };

  let isFetched: Promise<void>;

  const cacheMap: Record<string, TippyVueItem[]> = {};

  pdfViewer.addEventListener(ViewerEvent.PAGES_INIT, (info) => {
    info.source.container.addEventListener('scroll', () => {
      Object.keys(cacheMap).forEach((key) => {
        const list = cacheMap[key];
        list?.forEach((item) => {
          item.tippy.hide();
        });
      });
    });
  });

  pdfViewer.addEventListener(ViewerEvent.PAGE_RENDERED, async (info) => {
    if (!isFetched) {
      isFetched = fetchParseData();
    }
    await isFetched;

    const pageDiv = info.source.div;
    let viewportWidth = pdfViewer
      ?.getDocumentViewer()
      .getPdfViewer()
      .getPageView(0)?.viewport.width;

    if (!viewportWidth) {
      viewportWidth =
        (
          pdfViewer
            ?.getDocumentViewer()
            .container?.querySelector(
              '.page[data-page-number="1"]'
            ) as HTMLDivElement
        )?.offsetWidth || pageDiv.offsetWidth;
    }

    const validReferenceMarkers = referenceMarkers.filter(
      (m) => m.pageNum === info.pageNumber
    );
    const validFigureAndTableMarkers = figureAndTableMarkers.filter(
      (m) => m.pageNum === info.pageNumber
    );

    const key = info.pageNumber + '-' + 'reference';
    const prevTippyVues = cacheMap[key];
    destroyTippyVues(prevTippyVues);
    cacheMap[key] = [];

    appendDivs(
      {
        pageDiv,
        referenceMarkers: validReferenceMarkers,
        figureAndTableMarkers: validFigureAndTableMarkers,
        viewportWidth,
        pdfId,
      },
      cacheMap[key]
    );
  });
}
