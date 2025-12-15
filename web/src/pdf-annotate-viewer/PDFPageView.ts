import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer';
import TextLayerBuilder from './TextLayerBuilder';
import { TextLayerBuilderOptions } from '@idea/pdfjs-dist/types/web/text_layer_builder'
import { TextHighlighterOptions } from '@idea/pdfjs-dist/types/web/text_highlighter';
import TextHighlighter from './TextHighlighter';

export default class PDFPageView extends pdfjsViewer.PDFPageView {
  createTextLayerBuilder(options: TextLayerBuilderOptions) {
    return new TextLayerBuilder(options as any) as any;
  }
  createTextHighlighter(options: TextHighlighterOptions) {
    return new TextHighlighter(options) as any;
  }
}