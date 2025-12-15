import { createApp } from 'vue';
import PDFRightMarker from '~/src/components/PDFRightMarker/index.vue';
import Antd from 'ant-design-vue';
import { ViewerEvent, ViewerController } from '@idea/pdf-annotate-viewer';
import { CopyrightProps } from '../components/Copyright/type';
import { Nullable } from '../typings/global';
import { debounce } from 'lodash-es'
import i18n from '../locals/i18n';

export const createPDFSourceCopyrightVue = ({
  pageViewerDiv,
  pdfWebview,
  copyrightProps,
  isGroupPdf,
  paperId,
  isPrivatePaper,
}: {
  pageViewerDiv: HTMLDivElement;
  pdfWebview: ViewerController;
  copyrightProps: Nullable<CopyrightProps>;
  paperId: string;
  isPrivatePaper: boolean;
  isGroupPdf?: boolean;
}) => {
  const changeTranslate = (pageDiv: HTMLDivElement, scale: number) => {
    const pageWidth = pageDiv.offsetWidth
    const parentWidth = (pageDiv.offsetParent as HTMLDivElement).offsetWidth
    const sourceDiv = pageViewerDiv.querySelector('.js-copyright') as HTMLDivElement;
    if (sourceDiv) {
      sourceDiv.style.transform = `translateX(${-Math.min(pageWidth, parentWidth) / 2}px)`;
      setTimeout(() => {
        sourceDiv.style.display = 'flex';
      }, 0);
    }
    const sourceRightDivs = pageViewerDiv.querySelectorAll('.js-source-right') as NodeListOf<HTMLDivElement>;
    if (sourceRightDivs?.length) {
      sourceRightDivs.forEach(rightDiv => {
        rightDiv.style.transform = `translateX(${(pageWidth > parentWidth ? -parentWidth / 2 + pageWidth : pageWidth / 2)}px)`;
        setTimeout(() => {
          rightDiv.style.display = 'block';
        }, 0);
      })
      
    }
  }
  pdfWebview.addEventListener(ViewerEvent.PAGE_RENDERED, (info) => {
    
    if (info.pageNumber === 1 && info.source.div.offsetParent) {
      changeTranslate(info.source.div, info.source.scale)
    }
  });

  const app = createApp(PDFRightMarker, {
    copyright: copyrightProps,
    isGroupPdf,
    paperId,
    isPrivatePaper,
  });
  app.use(Antd);
  app.use(i18n)
  app.mount(pageViewerDiv);

  // const sourceDiv = pageViewerDiv.querySelector('.js-copyright') as HTMLDivElement;
  // if (sourceDiv) {
  //   sourceDiv.style.display = 'flex';
  // }

  new ResizeObserver(debounce(() => {
    const pageView = pdfWebview.getDocumentViewer().getPdfViewer().getPageView(0)
    if (pageView.div && pageView.div.offsetParent) {
      changeTranslate(pageView.div, pageView.scale)
    }
  }, 100, {
    leading: false,
    trailing: true,
  })).observe(pageViewerDiv)

  return app;
};
