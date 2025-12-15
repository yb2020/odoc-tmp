import { PageViewport } from "@idea/pdfjs-dist/types/src/display/display_utils";
import { BaseMousePoint, OffsetCoordinate, TextContentBound, TextContentChar, TextRectCoordinate } from "../type";
import { calcCharPoint } from "./text";

export const getCoordinateAfterRotatateAxis = (point: BaseMousePoint, m?: number[]) => {
  if (!m) {
    return point
  }
  // m的值[angleCos, angleSin, -angleSin, angleCos, left, top]
  const { x, y } = point
  return {
    x: x * m[0] + y * m[1],
    y: x * m[2] + y * m[3],
  }
}

export const getPointsOfRotatateRect = (offset: OffsetCoordinate) => {
  const { m } = offset
  if (!m) {
    // [topLeft, topRight, bottomRight, bottomLeft]
    return [
      { x: offset.left, y: offset.top },
      { x: offset.right, y: offset.top },
      { x: offset.right, y: offset.bottom },
      { x: offset.left, y: offset.bottom },
    ] as const
  }

  const [width, height] = offset.size
  const p1 = {
    x: m[4],
    y: m[5],
  }
  const p2 = {
    x: p1.x + width * m[0],
    y: p1.y + width * m[1],
  }

  const p3 = {
    x: p1.x - height * m[1],
    y: p1.y + height * m[0],
  }

  const p4 = {
    x: p2.x + p3.x - p1.x,
    y: p2.y + p3.y - p1.y,
  }

  // [topLeft, topRight, bottomRight, bottomLeft]
  return [p1, p2, p4, p3] as const

}

export const getBoundOffsetAfterRotatateAxis = (offset: OffsetCoordinate): OffsetCoordinate => {
  const { m } = offset
  if (!m) {
    return offset
  }

  const [p1, , p3] = getPointsOfRotatateRect(offset)

  const pTopLeft = getCoordinateAfterRotatateAxis(p1, m)
  const pRightBottom = getCoordinateAfterRotatateAxis(p3, m)

  return {
    ...offset,
    left: pTopLeft.x,
    right: pRightBottom.x,
    top: pTopLeft.y,
    bottom: pRightBottom.y,
  }
}

export const isInRotationRect = (offsetCoordinate: OffsetCoordinate, point: BaseMousePoint, buff: number) => {
  /**
   * 矩形有旋转
   * 坐标轴：左上角[0, 0]
   * 矩形不是平行于坐标轴的，那么把坐标轴旋转一下，让矩形平行，
   * 需要计算矩形顶点和目标点在旋转后坐标轴的坐标
   */
  const { x, y } = getCoordinateAfterRotatateAxis(point, offsetCoordinate.m)
  const offset = getBoundOffsetAfterRotatateAxis(offsetCoordinate)

  const { left, right, top, bottom } = offset

  if (
    x >= left - buff &&
    x <= right + buff &&
    y >= top - buff &&
    y <= bottom + buff
  ) {
    return true
  }
  return false

}

export const isUpDownRotationRect = ({ pos, offsetCoordinate, point, buff }: {
  pos: 'up' | 'down', offsetCoordinate: OffsetCoordinate, point: BaseMousePoint, buff: number
}) => {
  const { y } = getCoordinateAfterRotatateAxis(point, offsetCoordinate.m)
  const offset = getBoundOffsetAfterRotatateAxis(offsetCoordinate)

  const { top, bottom } = offset
  if (pos === 'up') {
    return y < bottom - buff
  }
  return y > top + buff
}

export const isInRotationBounds = (bounds: TextContentBound[], point: BaseMousePoint, buff: number) => {
  for (let i = 0; i < bounds.length; i += 1) {
    const bound = bounds[i]
    if (isInRotationRect(bound.offset, point, buff)) {
      return true
    }
  }
  return false
}

export const getNearestDistanceBetweenPointAndRect = (offsetCoordinate: OffsetCoordinate, point: BaseMousePoint) => {
  let distance
  if (
    isInRotationRect(offsetCoordinate, point, 0)
  ) {
    distance = 0
  } else {
    const points = getPointsOfRotatateRect(offsetCoordinate)
    const distances = points.map(p => {
      return Math.pow(point.x - p.x, 2) + Math.pow(point.y - p.y, 2)
    })

    distance = Math.min(...distances)

  }
  return distance
}

export const findTargetTextPosition = ({ bound, point, include, viewport, ctx }: {
  bound: TextContentBound, 
  point: BaseMousePoint, 
  include?: boolean,
  viewport: PageViewport,
  ctx: CanvasRenderingContext2D,
}) => {

  /**
   * 文本有旋转角度，那么就会有一个问题
   */


  const rotatePoint = getCoordinateAfterRotatateAxis(point, bound.offset.m)

  const rotateLeftTopPoint = 
    bound.offset.m ?
      getCoordinateAfterRotatateAxis({ x: bound.offset.m[4], y: bound.offset.m[5] }, bound.offset.m) :
      { x: bound.offset.left, y: bound.offset.top }

  const targetRange = rotatePoint.x - rotateLeftTopPoint.x

  let { charsArray } = bound.geom


  if (targetRange <= 0) {
    return {
      text: '',
      width: 0,
    }
  }

  if (!charsArray) {
    charsArray = []
  }
  

  const chars: { idx: number, char: TextContentChar }[] = []

  let curWidth = 0
  
  for(let i = 0; i < charsArray.length; i += 1) {
    if (!charsArray[i].transfrom) {
      chars.push({
        idx: i,
        char: charsArray[i],
      })
      continue
    }
    const charPoint = calcCharPoint({viewport, transform: charsArray[i].transfrom})
    const left = getCoordinateAfterRotatateAxis(charPoint, bound.offset.m).x
    
    const width = left - rotateLeftTopPoint.x


    if (include) {

      if (width >= targetRange) {
        break
      }
      chars.push({
        idx: i,
        char: charsArray[i],
      })
      curWidth = width;
    } else {
      chars.push({
        idx: i,
        char: charsArray[i],
      })
      curWidth = width;
      if (width >= targetRange) {
        break
      }     
    }
  }

  const last = chars[chars.length - 1]

  if (last) {
    if (include) {
      const [, nextEndChar] = [charsArray[last.idx], charsArray[last.idx + 1]]
      if (nextEndChar) {
        curWidth += nextEndChar.spaceWidth * viewport.scale
      }
    }
    return {
      text: last.char.str,
      width: curWidth,
      
    }
  }

  ctx.font = `${bound.style.fontSize} ${bound.style.fontFamily}`
  const scale = bound.offset.size[0] / ctx.measureText(bound.str).width
  let text = ''
  const str = bound.str
  let width = 0
  
  for (let i = 0; i < str.length; i += 1) {
    text += str[i]
    width = ctx.measureText(text).width * scale
    if (width >= targetRange) {
      break
    }
  }
  if (include && width > targetRange) {
    text = text.substring(0, text.length - 1)
  }

  return {
    text,
    width: ctx.measureText(text).width * scale,
    
  }
  
  
}



const getStartTextRectCoordinate = (
  { bound, ctx, point, viewport }: {
    point: BaseMousePoint,
    bound: TextContentBound,
    ctx: CanvasRenderingContext2D,
    viewport: PageViewport
  }
): TextRectCoordinate => {
  const { offset, str } = bound

  const { m } = offset

  const pos = findTargetTextPosition({
    bound,
    include: true,
    point,
    viewport,
    ctx,
  })

  let text = str.substring(pos.text.length)

  if (bound.geom.hasEOL) {
    text += '\n';
  }

  /**
   * 这里的[x,y]坐标是旋转后矩形的左上角顶点坐标，
   * [bound.left,bound.top]是无旋转角度的矩形的左上角顶点坐标
   * [m[4],m[5]]是旋转角度的矩形的左上角顶点坐标
   * canvas画矩形的时候，对于有旋转角度的，按如下处理：
   * 1.先translate画布到旋转后的左上角位置
   * 2.然后rotate画布angle角度
   * 3.接下来将继续translate画布dx的距离
   * 4.最后画矩形[0,0,width,height]
   * 对于没有旋转角度的，不移动画布位置，因为[x,y]已经包含了偏移量
   */
  if (m) {
    return {
      x: m[4],
      y: m[5],
      width: bound.offset.size[0] - pos.width,
      height: bound.offset.size[1],
      offset: {
        angle: bound.offset.angle,
        dx: pos.width,
        dy: 0,
      },
      text,
      shouldScaleText: bound.shouldScaleText,
    }
  }


  return {
    x: offset.left + pos.width,
    y: offset.top,
    width: bound.offset.size[0] - pos.width,
    height: bound.offset.size[1],
    text,
    shouldScaleText: bound.shouldScaleText,
  }
}

const getEndTextRectCoordinate = (
  { bound, ctx, point, viewport }: {
    point: BaseMousePoint,
    bound: TextContentBound,
    ctx: CanvasRenderingContext2D,
    viewport: PageViewport
  }
): TextRectCoordinate => {
  const { offset } = bound

  const { m } = offset

  const pos = findTargetTextPosition({
    bound,
    include: false,
    point,
    viewport,
    ctx,
  })
  const text = pos.text
   if (m) {
    return {
      x: m[4],
      y: m[5],
      width: pos.width,
      height: bound.offset.size[1],
      offset: {
        angle: bound.offset.angle,
        dx: 0,
        dy: 0,
      },
      text,
      shouldScaleText: bound.shouldScaleText,
    }
  }

  return {
    x: offset.left,
    y: offset.top,
    width: pos.width,
    height: bound.offset.size[1],
    text,
    shouldScaleText: bound.shouldScaleText,
  }
}

export const getTextRectCoordinate = (
  { bound, ctx, points, viewport }: {
    points: (BaseMousePoint | null)[],
    bound: TextContentBound,
    ctx: CanvasRenderingContext2D,
    viewport: PageViewport
  }
): TextRectCoordinate => {

  const { offset, str } = bound

  const { m } = offset

  const [p1, p2] = points

  if (p1 && !p2) {
    return getStartTextRectCoordinate({
      bound,
      ctx,
      point: p1,
      viewport,
    })
  }

  if (!p1 && p2) {
    return getEndTextRectCoordinate({
      bound,
      ctx,
      point: p2,
      viewport,
    })
  }

  if (p1 && p2) {
    const rP1 = getCoordinateAfterRotatateAxis(p1, offset.m)
    const rP2 = getCoordinateAfterRotatateAxis(p2, offset.m)
    const [sP, eP] = rP1.x >= rP2.x ? [p2, p1] as const : [p1, p2] as const

    const sPos = findTargetTextPosition({
      bound,
      include: true,
      point: sP,
      viewport,
      ctx,
    })
    const ePos = findTargetTextPosition({
      bound,
      point: eP,
      viewport,
      ctx,
    })
    const text = str.substring(sPos.text.length, ePos.text.length)
    if (m) {
      return {
        x: m[4],
        y: m[5],
        width: ePos.width - sPos.width,
        height: offset.size[1],
        offset: {
          angle: offset.angle,
          dx: sPos.width,
          dy: 0,
        },
        text,
        shouldScaleText: bound.shouldScaleText,
      }
    }
  
    return {
      x: offset.left + sPos.width,
      y: offset.top,
      width: ePos.width - sPos.width,
      height: offset.size[1],
      text,
      shouldScaleText: bound.shouldScaleText,
    }
  }

  let text = str

  if (bound.geom.hasEOL) {
    text += '\n';
  }

  if (m) {
    return {
      x: m[4],
      y: m[5],
      width: offset.size[0],
      height: offset.size[1],
      offset: {
        angle: offset.angle,
        dx: 0,
        dy: 0,
      },
      text,
      shouldScaleText: bound.shouldScaleText,
    }
  }

  return {
    x: offset.left,
    y: offset.top,
    width: offset.size[0] ,
    height: offset.size[1],
    text,
    shouldScaleText: bound.shouldScaleText,
  }

}

export const getWordRectCoordinate = (
  { bound, point, viewport }: {
    point: BaseMousePoint,
    bound: TextContentBound,
    viewport: PageViewport
  }
): TextRectCoordinate | null => {
  const rotatePoint = getCoordinateAfterRotatateAxis(point, bound.offset.m)

  const rotateLeftTopPoint = 
    bound.offset.m ?
      getCoordinateAfterRotatateAxis({ x: bound.offset.m[4], y: bound.offset.m[5] }, bound.offset.m) :
      { x: bound.offset.left, y: bound.offset.top }

  const targetRange = rotatePoint.x - rotateLeftTopPoint.x

  const { charsArray } = bound.geom

  const chars: ({ idx: number, char: TextContentChar } | null)[] = []

  let i = 0;

  for(; i < charsArray.length; i += 1) {
    if (!charsArray[i].transfrom) {
      continue
    }

    const charPoint = calcCharPoint({viewport, transform: charsArray[i].transfrom})
    const left = getCoordinateAfterRotatateAxis(charPoint, bound.offset.m).x
    
    const width = left - rotateLeftTopPoint.x


    if (width >= targetRange) {
      break
    }
  }

  let idx = i

  while( idx >= 0) {
    const cur = charsArray[idx]
    const pre = charsArray[idx - 1]
    if (!pre) {
      chars.push(/[a-zA-Z0-9]/.test(cur.str) ? null: {
        idx: idx,
        char: cur,
      })
      break
    } else {
      const diffStr = cur.str.substring(pre.str.length)
      if (!/^[a-zA-Z0-9]+$/g.test(diffStr)) {
        if (/[a-zA-Z0-9]/.test(diffStr)) { 
          chars.push({
            idx: idx - 1,
            char: pre,
          })
        } else {
          chars.push({
            idx: idx,
            char: cur,
          })
        }
        break;
      } else {
        idx -= 1;
      }
    }
  }

  idx = i

  while( idx <= charsArray.length - 1) {
    const cur = charsArray[idx]
    const next = charsArray[idx + 1]
    if (!next) {
      chars.push({
        idx,
        char: cur,
      })
      break
    } else {
      const diffStr = next.str.substring(cur.str.length)
      if (!/^[a-zA-Z0-9]+$/g.test(diffStr)) {
        chars.push({
          idx,
          char: cur,
        })
        break;
      } else {
        idx += 1;
      }
    }
  }


  if (!chars[1]) {
    return null
  }


  let rStartPoint: null | BaseMousePoint = null

  if (!chars[0]) {
    rStartPoint = rotateLeftTopPoint
  } else {
    rStartPoint = getCoordinateAfterRotatateAxis(calcCharPoint({viewport, transform: chars[0].char.transfrom}), bound.offset.m)
  }

   const rEndPoint = getCoordinateAfterRotatateAxis(calcCharPoint({viewport, transform: chars[1].char.transfrom}), bound.offset.m)



  const { m } = bound.offset


  const text =  chars[1].char.str.substring(chars[0]?.char.str.length || 0).trim()


  let dxx = 0;

  if (chars[0]) {
    const [, curStartChar] = [charsArray[chars[0].idx], charsArray[chars[0].idx + 1]]
    if (curStartChar) {
      dxx = curStartChar.spaceWidth * viewport.scale
    }
  }

  const curWidth = rEndPoint.x - rStartPoint.x - dxx


  if (m) {
    return {
      x: m[4],
      y: m[5],
      width: curWidth,
      height: bound.offset.size[1],
      offset: {
        angle: bound.offset.angle,
        dx: rStartPoint.x + dxx - rotateLeftTopPoint.x,
        dy: 0,
      },
      text,
      shouldScaleText: bound.shouldScaleText,
    }
  }


  return {
    x: rStartPoint.x + dxx,
    y: bound.offset.top,
    width: curWidth,
    height: bound.offset.size[1],
    text,
    shouldScaleText: bound.shouldScaleText,
  }
}

export const optimizeRects = (coordinateList: TextRectCoordinate[]) => {


  // coordinates.reduce((prev, cur) => {
  //   const last = prev[prev.length - 1]
  //   if (!last) {
  //     prev.push(cur)
  //   } else if (last.y === cur.y && last.x + last.width <= cur.x) {
  //     last.width = cur.x - last.x + cur.width
  //     last.text += cur.text
  //   } else {
  //     prev.push(cur)
  //   }
  //   return prev
  // }, mergedCoordinates)

  const optimizeNumber = (number: number) => Math.round(number * 100) / 100
  // fix: 优化选区，将相近的一行bound进行合并，防止出现数据过长的问题
  coordinateList.forEach((coordinate) => {
    coordinate.x = optimizeNumber(coordinate.x);
    coordinate.y = optimizeNumber(coordinate.y);
    coordinate.width = optimizeNumber(coordinate.width);
    coordinate.height = optimizeNumber(coordinate.height);
  })
}
