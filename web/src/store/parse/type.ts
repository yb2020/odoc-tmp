import { FigureAndTableReferenceMarker, PdfFigureAndTableInfo } from "@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse";
import { TippyVueItem } from "~/src/dom/tippy";
import { Nullable } from "~/src/typings/global";

export interface FigureTippyTriggerItem {
  id: string,
  triggerEle: Element,
  pdfId: string,

}

export interface PdfFigureAndTableInfosType {
  list: Nullable<PdfFigureAndTableInfo[]>,
  list0: Nullable<FigureAndTableReferenceMarker[]>,
  error: Nullable<Error>,
  pending: boolean,
}

export interface ParseState {
  figureTippyTriggers: (FigureTippyTriggerItem & { item: TippyVueItem })[],
  pdfFigureAndTableInfos: Record<string, PdfFigureAndTableInfosType>,
}