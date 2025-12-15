import { ViewerController, ViewerControllerOptions } from './ViewerController';
import _PDFViewerColorToneAddon from './addon/theme'
import * as pdfjsLib from '@idea/pdfjs-dist';
import '@idea/pdfjs-dist/web/pdf_viewer.css';
import 'perfect-scrollbar/css/perfect-scrollbar.css';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';

export const initialPdfjsLib = ({
  pdfjsWorkerSrc,
}: {
  pdfjsWorkerSrc: string;
}) => {
  pdfjsLib.GlobalWorkerOptions.workerSrc = pdfjsWorkerSrc;
};

export const createPDFWebview = (options: ViewerControllerOptions, pdfjsWorkerSrc?: string) => {
  if (pdfjsWorkerSrc) {
    initialPdfjsLib({pdfjsWorkerSrc})
  }
  options.pdfDocumentParams.language = options.pdfDocumentParams.language || Language.EN_US
  const viewer = new ViewerController(options);
  return viewer;
};

export * from './type';

export * from './ScrollController';

export const PDFViewerColorToneAddon = _PDFViewerColorToneAddon

export {
  useMouseDownMoveUpDoubleClick,
  detectCtrlPressing,
  JS_IGNORE_MOUSE_OUTSIDE,
} from './utils/mouse';
