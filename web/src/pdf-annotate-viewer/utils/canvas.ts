import { PageViewport } from "@idea/pdfjs-dist";
import { TextRectCoordinate } from "../type";

export const drawTextRectCoordinate = ({
  ctx,
  coordinate,
  color,
}: {
  ctx: CanvasRenderingContext2D;
  coordinate: TextRectCoordinate;
  color: string;
}) => {
  ctx.fillStyle = color;

  const { x, y, width, height, offset } = coordinate;
  if (offset && offset.angle) {
    ctx.save();
    ctx.translate(x, y);
    ctx.rotate(offset.angle);
    ctx.translate(offset.dx, offset.dy);
    if ((ctx as any).origFillRect) {
      (ctx as any).origFillRect(0, 0, width, height);
    } else {
      ctx.fillRect(0, 0, width, height);
    }
    ctx.restore();
  } else {
    if ((ctx as any).origFillRect) {
      (ctx as any).origFillRect(x || 0, y || 0, width || 0, height || 0);
    } else {
      ctx.fillRect(x || 0, y || 0, width || 0, height || 0);
    }
   
  }
};

export const createCanvas = (viewport: PageViewport, canvasEle?: HTMLCanvasElement) => {
  const width = Math.floor(viewport.width)
  const height = Math.floor(viewport.height)
  const canvas = canvasEle || document.createElement('canvas')
  canvas.style.width = width + 'px'
  canvas.style.height = height + 'px'
  canvas.style.position = 'absolute'
  canvas.setAttribute('width', width.toString())
  canvas.setAttribute('height', height.toString())
  return canvas
}