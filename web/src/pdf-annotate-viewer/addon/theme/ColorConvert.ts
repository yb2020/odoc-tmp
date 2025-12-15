import BaseColor from 'colorjs.io'
import { ColorTypes } from 'colorjs.io/types/src/color';

export default class ColorConvert extends BaseColor {
  public static white = new ColorConvert('#FFFFFF');

  constructor(...args: [ColorTypes] | ConstructorParameters<typeof BaseColor> ) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    super(...args)
  }
  get lightness() {
    return this.lab[0];
  }

  get chroma() {
    const [, a, b] = this.lab;
    return Math.sqrt(a ** 2 + b ** 2);
  }

  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  deltaE(color: ColorConvert) {
    return Math.sqrt(
      this.lab.reduce((a: number, c: number, i: number) => {
        if (isNaN(c) || isNaN(color.lab[i])) {
          return a;
        }
        return a + (color.lab[i] - c) ** 2;
      }, 0)
    );
  }

  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  range(color: ColorConvert) {
    function interpolate(start: number, end: number, p: number) {
      if (isNaN(start)) {
        return end;
      }
      if (isNaN(end)) {
        return start;
      }
      return start + (end - start) * p;
    }
    return (p: number) => {
      const coords = this.lab.map((start: number, i: number) => {
        const end = color.lab[i];
        return interpolate(start, end, p);
      }) as [number, number, number];
      return new ColorConvert("lab", coords);
    };
  }

  toHex(alpha?: number) {
    const coords = this.srgb.map((c: number) => Math.round(c * 255));
    if (alpha !== undefined && alpha < 1) {
      return `rgba(${coords[0]}, ${coords[1]}, ${coords[2]}, ${alpha})`
    }
    return rgbToHex(coords)
  }
}

export const rgbToHex = (rgb: number[]) => {
  const hex = rgb
      .map((c) => {
        c = Math.min(Math.max(c, 0), 255);
        return c.toString(16).padStart(2, "0");
      })
      .join("");
    return "#" + hex;
}
