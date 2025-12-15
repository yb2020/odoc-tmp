import { PageViewport } from "@idea/pdfjs-dist/types/src/display/display_utils";
import { TextRectCoordinate } from "..";
import { TextContentBound } from "../type";
import { getCoordinateAfterRotatateAxis } from "./bound";
import { calcCharPoint } from "./text";

const getTargetIndexChar = (index: number, bound: TextContentBound) => {
  const { charsArray } = bound.geom;
  for (let i = 0; i < charsArray.length; i += 1) {
    const chars = charsArray[i]
    if (chars.str.length >= index + 1) {
      return i
    }
  }
  return charsArray.length - 1
}

export const getHightlightRectCoordinate = ({ textContentItems, matches, matchesLength, viewport }: {
  textContentItems: TextContentBound[],
  matches: number[],
  matchesLength: number[],
  viewport: PageViewport
}) => {
  if (!matches) {
    return [];
  }


  let i = 0,
    iIndex = 0;
  const end = textContentItems.length - 1;
  const result = [];

  for (let m = 0, mm = matches.length; m < mm; m++) {
    // Calculate the start position.
    let matchIdx = matches[m];

    // Loop over the divIdxs.
    while (i !== end && matchIdx >= iIndex + textContentItems[i].str.length) {
      iIndex += textContentItems[i].str.length;
      i++;
    }

    if (i === textContentItems.length) {
      console.error("Could not find a matching mapping");
    }

    const match = {
      begin: {
        divIdx: i,
        offset: matchIdx - iIndex,
      },
      end: {
        divIdx: -1,
        offset: 0,
      },
    };

    // Calculate the end position.
    matchIdx += matchesLength[m];

    // Somewhat the same array as above, but use > instead of >= to get
    // the end position right.
    while (i !== end && matchIdx > iIndex + textContentItems[i].str.length) {
      iIndex += textContentItems[i].str.length;
      i++;
    }

    match.end = {
      divIdx: i,
      offset: matchIdx - iIndex,
    };
    result.push(match);
  }
  const coordinates: TextRectCoordinate[] = []

  result.forEach(({ begin, end }) => {
    if (begin.divIdx === end.divIdx) {
      const bound = textContentItems[begin.divIdx]
      const rotateLeftTopPoint =
        bound.offset.m ?
          getCoordinateAfterRotatateAxis({ x: bound.offset.m[4], y: bound.offset.m[5] }, bound.offset.m) :
          { x: bound.offset.left, y: bound.offset.top }

      const startChar = bound.geom.charsArray[getTargetIndexChar(begin.offset, bound) - 1];
      const startNextChar = bound.geom.charsArray[getTargetIndexChar(begin.offset, bound)]
      const endChar = bound.geom.charsArray[getTargetIndexChar(end.offset - 1, bound)]
      const startPoint = startChar ? getCoordinateAfterRotatateAxis(calcCharPoint({ viewport, transform: startChar.transfrom }), bound.offset.m) : rotateLeftTopPoint
      const endPoint = getCoordinateAfterRotatateAxis(calcCharPoint({ viewport, transform: endChar.transfrom }), bound.offset.m)
      let dxx = 0;
      if (startChar && startNextChar?.spaceWidth) {
        dxx = startNextChar.spaceWidth * viewport.scale
      }
      coordinates.push({
        x: bound.offset.m ? bound.offset.m[4] : startPoint.x + dxx,
        y: bound.offset.m ? bound.offset.m[5] : bound.offset.top,
        width: endPoint.x - startPoint.x - dxx,
        height: bound.offset.size[1],
        text: endChar.str.substring(startChar?.str.length - 1 || 0),
        offset: bound.offset.m ? {
          angle: bound.offset.angle,
          dx: startPoint.x - rotateLeftTopPoint.x,
          dy: 0,
        }: undefined,
        shouldScaleText: bound.shouldScaleText,
      })
    } else {
      // TODO angle !== 0
      const startBound = textContentItems[begin.divIdx];
      const startChar = startBound.geom.charsArray[getTargetIndexChar(begin.offset, startBound) - 1];
      const startNextChar = startBound.geom.charsArray[getTargetIndexChar(begin.offset, startBound)]
      const startPoint = startChar ? calcCharPoint({ viewport, transform: startChar.transfrom }) : { x: startBound.offset.left, y: startBound.offset.top }
      let dxx = 0;
      if (startChar && startNextChar?.spaceWidth) {
        dxx = startNextChar.spaceWidth * viewport.scale
      }
      coordinates.push({
        x: startPoint.x + dxx,
        y: startBound.offset.top,
        width: startBound.offset.left + startBound.offset.size[0] - startPoint.x - dxx,
        height: startBound.offset.size[1],
        text: startBound.str.substring(startChar?.str.length - 1 || 0),
        shouldScaleText: startBound.shouldScaleText,
      })
      for (let i = begin.divIdx + 1; i <= end.divIdx - 1; i += 1) {
        const bound = textContentItems[i]
        coordinates.push({
          x: bound.offset.left,
          y: bound.offset.top,
          width: bound.offset.size[0],
          height: bound.offset.size[1],
          text: bound.str,
          shouldScaleText: bound.shouldScaleText,
        })
      }
      const endBound = textContentItems[end.divIdx];
      const endChar = endBound.geom.charsArray[getTargetIndexChar(end.offset - 1, endBound)];
      const endPoint = calcCharPoint({ viewport, transform: endChar.transfrom })
      coordinates.push({
        x: endBound.offset.left,
        y: endBound.offset.top,
        width: endPoint.x - endBound.offset.left,
        height: endBound.offset.size[1],
        text: endChar.str,
        shouldScaleText: endBound.shouldScaleText,
      })
    }
  })
  return coordinates;
}