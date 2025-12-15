import { PDFDocumentProxy } from "@idea/pdfjs-dist"
import { PDFViewer, PageViewport } from "@idea/pdfjs-dist/types/web/pdf_viewer"
import { convertToRegExpString, getOriginalIndex, normalize } from "./utils/matchChars"
import { TextRectCoordinate } from "./type"
import TextHighlighter from "./TextHighlighter"

export default class WordMatchController {
  pdfDocument: PDFDocumentProxy
  pdfViewer: PDFViewer

  private pageContents = [] as { content: string, diffs: number[][], hasDiacritics: boolean }[]

  constructor(pdfDocument: PDFDocumentProxy, pdfViewer: PDFViewer) {
    this.pdfDocument = pdfDocument
    this.pdfViewer = pdfViewer
  }

  async updatePageContents() {
    const p: Promise<typeof this.pageContents[0]>[] = []
      const pdfDocument = this.pdfDocument as PDFDocumentProxy
      const pageNums = pdfDocument.numPages
      for (let i = 1; i <= pageNums; i++) {
        p.push(pdfDocument.getPage(i).then(pdfPage => pdfPage.getTextContent()).then(textContent => {
          const strBuf = [];
          for (const textItem of textContent.items) {
            strBuf.push((textItem as any).str);
            if ((textItem as any).hasEOL) {
              strBuf.push("\n");
            }
          }
          // Store the normalized page content (text items) as one string.
          const [
            pageContent,
            pageDiff,
            hasDiacritic,
          ] = normalize(strBuf.join(""));
          return {
            content: pageContent as string,
            diffs: pageDiff as number[][],
            hasDiacritics: hasDiacritic as boolean,
          };
        }).catch(reason => {
          console.error(
            `Unable to get text content for page ${i + 1}`,
            reason
          );
          return {
            content: '',
            diffs: [] as number[][],
            hasDiacritics: false as boolean,
          }
        }))
      }
      this.pageContents = await Promise.all(p);
  }


  async getWordMatchRects(words: string[], pageIndex?: number) {

    if (!this.pageContents.length) {
      await this.updatePageContents()
    }

    const pageContents = this.pageContents
  
    const caseSensitive = false;
  
    const pdfViewer = this.pdfViewer as PDFViewer
  
    const results: { 
      matches: {
        [pageIndex: number]: { matches: number[]; matchesLength: number[], rects: TextRectCoordinate[], viewport: PageViewport }
      }, 
      word: string,
    }[] = [];

    let startPageIndex = pageIndex === void 0 ? 0 : pageIndex;
    const endPageIndex = pageIndex === void 0 ? pdfViewer.pagesCount - 1 : pageIndex;

    const p = [] as Promise<void>[]

    for(; startPageIndex <= endPageIndex; startPageIndex++) {
      const pageContent = pageContents[startPageIndex];
      const diffs = pageContent.diffs;
  
      const textHighlighter: TextHighlighter | null =
        pdfViewer.getPageView(startPageIndex)?.textLayer?.highlighter;
  
  
      if (!textHighlighter) {
        // 说明当前页没有被渲染出来
        if (pageIndex === void 0) {
          // 如果没有指定页码，说明是获取全部页码的匹配结果，那么就不用等待当前页渲染出来，直接跳过
          continue;
        } else {
          throw Error(`page ${startPageIndex} is not rendered`)
        }
      }

      
      
      words.forEach((word, wordIndex) => {
        const [isUnicode, q] = convertToRegExpString(word, pageContent.hasDiacritics, false)
        const flags = `g${isUnicode ? "u" : ""}${caseSensitive ? "" : "i"}`;
        const query = q ? new RegExp(q as string, flags) : null;
        if (!query) {
          return;
        }
        const wordResults = {
          matches: [] as number[],
          matchesLength: [] as number[],
        };
        let match;
        while ((match = query.exec(pageContent.content)) !== null) {
          const [matchPos, matchLen] = getOriginalIndex(
            diffs,
            match.index,
            match[0].length
          );
          if (matchLen) {
            wordResults.matches.push(matchPos);
            wordResults.matchesLength.push(matchLen);
          }
        }

        if (!wordResults.matches.length) {
          return
        }


        results[wordIndex] = results[wordIndex] || {
          matches: {},
          word,
        }
        p.push((async (startPageIndex: number, wordIndex: number, textHighlighter: TextHighlighter, wordResults: {
          matches: number[];
          matchesLength: number[];
      }) => {
          await textHighlighter.readyPromise.promise;
          results[wordIndex].matches[startPageIndex] = {
            matches: wordResults.matches,
            matchesLength: wordResults.matchesLength,
            rects: textHighlighter._convertMatches(wordResults.matches, wordResults.matchesLength),
            viewport: textHighlighter.viewport!,
          };
          return void 0;
        })(startPageIndex, wordIndex, textHighlighter, wordResults))
      });

    }

    await Promise.all(p)

    return results
  
  }

}