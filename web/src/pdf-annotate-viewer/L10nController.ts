import * as pdfjsViewer from '@idea/pdfjs-dist/web/pdf_viewer';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';

export default class L10nController {
  private l10n;
  constructor(lang?: string) {
    const link = document.createElement('link');
    link.setAttribute('rel', 'resource')
    link.setAttribute('type', 'application/l10n')
    link.setAttribute('href', 'https://nuxt.cdn.readpaper.com/pdfjs-dist%402.13.216/local/locale.properties')
    document.head.appendChild(link)
    this.l10n = new pdfjsViewer.GenericL10n(lang || 'EN_US')
  }

  getL10n() {
    return this.l10n
  }
}

const loadMap = new Map<string, L10nController>()

export const createL10nController = (lang = 'EN_US') => {
  if (loadMap.has(lang)) {
    return loadMap.get(lang)!
  }
  const l10nController = new L10nController(lang)
  loadMap.set(lang, l10nController)
  return l10nController
}

