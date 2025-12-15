import * as pdfjsLib from '@idea/pdfjs-dist';
import { TextItem } from "@idea/pdfjs-dist/types/src/display/api";
import { PageViewport } from "@idea/pdfjs-dist/types/src/display/display_utils";

const Util = pdfjsLib.Util

const ascentCache = new Map()

const DEFAULT_FONT_SIZE = 30;
const DEFAULT_FONT_ASCENT = 0.8;

export const getAscent = (fontFamily: string, ctx: CanvasRenderingContext2D) => {
  const cachedAscent = ascentCache.get(fontFamily);
  if (cachedAscent) {
    return cachedAscent;
  }

  ctx.save();
  ctx.font = `${DEFAULT_FONT_SIZE}px ${fontFamily}`;
  const metrics = ctx.measureText("");

  // Both properties aren't available by default in Firefox.
  let ascent = metrics.fontBoundingBoxAscent;
  let descent = Math.abs(metrics.fontBoundingBoxDescent);
  if (ascent) {
    ctx.restore();
    const ratio = ascent / (ascent + descent);
    ascentCache.set(fontFamily, ratio);
    return ratio;
  }

  // Try basic heuristic to guess ascent/descent.
  // Draw a g with baseline at 0,0 and then get the line
  // number where a pixel has non-null red component (starting
  // from bottom).
  ctx.strokeStyle = "red";
  ctx.clearRect(0, 0, DEFAULT_FONT_SIZE, DEFAULT_FONT_SIZE);
  ctx.strokeText("g", 0, 0);
  let pixels = ctx.getImageData(
    0,
    0,
    DEFAULT_FONT_SIZE,
    DEFAULT_FONT_SIZE
  ).data;
  descent = 0;
  for (let i = pixels.length - 1 - 3; i >= 0; i -= 4) {
    if (pixels[i] > 0) {
      descent = Math.ceil(i / 4 / DEFAULT_FONT_SIZE);
      break;
    }
  }

  // Draw an A with baseline at 0,DEFAULT_FONT_SIZE and then get the line
  // number where a pixel has non-null red component (starting
  // from top).
  ctx.clearRect(0, 0, DEFAULT_FONT_SIZE, DEFAULT_FONT_SIZE);
  ctx.strokeText("A", 0, DEFAULT_FONT_SIZE);
  pixels = ctx.getImageData(0, 0, DEFAULT_FONT_SIZE, DEFAULT_FONT_SIZE).data;
  ascent = 0;
  for (let i = 0, ii = pixels.length; i < ii; i += 4) {
    if (pixels[i] > 0) {
      ascent = DEFAULT_FONT_SIZE - Math.floor(i / 4 / DEFAULT_FONT_SIZE);
      break;
    }
  }

  ctx.restore();

  if (ascent) {
    const ratio = ascent / (ascent + descent);
    ascentCache.set(fontFamily, ratio);
    return ratio;
  }

  ascentCache.set(fontFamily, DEFAULT_FONT_ASCENT);
  return DEFAULT_FONT_ASCENT;
}

const AllWhitespaceRegexp = /^\s+$/g;
export const getTextProperties = ({ viewport, geom, ctx, styles }: {
  viewport: PageViewport,
  geom: TextItem,
  ctx: CanvasRenderingContext2D,
  styles: Record<string, { ascent: number, descent: number, fontFamily: string, vertical: boolean }>
}) => {

  let shouldScaleText = false;

  if (
    geom.str.length > 1 ||
    (!AllWhitespaceRegexp.test(geom.str))
  ) {
    shouldScaleText = true;
  } else if (geom.str !== " " && geom.transform[0] !== geom.transform[3]) {
    const absScaleX = Math.abs(geom.transform[0]),
      absScaleY = Math.abs(geom.transform[3]);
    // When the horizontal/vertical scaling differs significantly, also scale
    // even single-char text to improve highlighting (fixes issue11713.pdf).
    if (
      absScaleX !== absScaleY &&
      Math.max(absScaleX, absScaleY) / Math.min(absScaleX, absScaleY) > 1.5
    ) {
      shouldScaleText = true;
    }
  }

  // Initialize all used properties to keep the caches monomorphic.

  const tx = Util.transform(viewport.transform, geom.transform);
  let angle = Math.atan2(tx[1], tx[0]);
  const style = styles[geom.fontName];
  if (style.vertical) {
    angle += Math.PI / 2;
  }
  const fontHeight = Math.hypot(tx[2], tx[3]);
  const fontAscent = fontHeight * getAscent(style.fontFamily, ctx);

  let left, top;
  if (angle === 0) {
    left = tx[4];
    top = tx[5] - fontAscent;
  } else {
    left = tx[4] + fontAscent * Math.sin(angle);
    top = tx[5] - fontAscent * Math.cos(angle);
  }


  if (geom.str !== "") {
    let angleCos = 1,
      angleSin = 0;
    if (angle !== 0) {
      angleCos = Math.cos(angle);
      angleSin = Math.sin(angle);
    }
    const divWidth =
      (style.vertical ? geom.height : geom.width) * viewport.scale;
    const divHeight = fontHeight;

    let m, b;
    if (angle !== 0) {
      m = [angleCos, angleSin, -angleSin, angleCos, left, top];
      b = Util.getAxialAlignedBoundingBox([0, 0, divWidth, divHeight], m);
    } else {
      b = [left, top, left + divWidth, top + divHeight];
    }

    return {
      offset: {
        left: b[0],
        top: b[1],
        right: b[2],
        bottom: b[3],
        size: [divWidth, divHeight] as [number, number],
        m,
        angle,
      },
      style: {
        ...style,
        fontSize: fontHeight,
      },
      shouldScaleText,
    }
  }
  return null
}

export const calcCharPoint = ({ viewport, transform }: {
  viewport: PageViewport,
  transform: number[],
  }) => {
  const tx = Util.transform(viewport.transform, transform);
  return {
    x: tx[4],
    y: tx[5],
  }
}


const CHARACTERS_TO_NORMALIZE = {
  "\u2010": "-", // Hyphen
  "\u2018": "'", // Left single quotation mark
  "\u2019": "'", // Right single quotation mark
  "\u201A": "'", // Single low-9 quotation mark
  "\u201B": "'", // Single high-reversed-9 quotation mark
  "\u201C": '"', // Left double quotation mark
  "\u201D": '"', // Right double quotation mark
  "\u201E": '"', // Double low-9 quotation mark
  "\u201F": '"', // Double high-reversed-9 quotation mark
  "\u00BC": "1/4", // Vulgar fraction one quarter
  "\u00BD": "1/2", // Vulgar fraction one half
  "\u00BE": "3/4", // Vulgar fraction three quarters
};

// These diacritics aren't considered as combining diacritics
// when searching in a document:
//   https://searchfox.org/mozilla-central/source/intl/unicharutil/util/is_combining_diacritic.py.
// The combining class definitions can be found:
//   https://www.unicode.org/reports/tr44/#Canonical_Combining_Class_Values
// Category 0 corresponds to [^\p{Mn}].
// const DIACRITICS_EXCEPTION = new Set([
//   // UNICODE_COMBINING_CLASS_KANA_VOICING
//   // https://www.compart.com/fr/unicode/combining/8
//   0x3099, 0x309a,
//   // UNICODE_COMBINING_CLASS_VIRAMA (under 0xFFFF)
//   // https://www.compart.com/fr/unicode/combining/9
//   0x094d, 0x09cd, 0x0a4d, 0x0acd, 0x0b4d, 0x0bcd, 0x0c4d, 0x0ccd, 0x0d3b,
//   0x0d3c, 0x0d4d, 0x0dca, 0x0e3a, 0x0eba, 0x0f84, 0x1039, 0x103a, 0x1714,
//   0x1734, 0x17d2, 0x1a60, 0x1b44, 0x1baa, 0x1bab, 0x1bf2, 0x1bf3, 0x2d7f,
//   0xa806, 0xa82c, 0xa8c4, 0xa953, 0xa9c0, 0xaaf6, 0xabed,
//   // 91
//   // https://www.compart.com/fr/unicode/combining/91
//   0x0c56,
//   // 129
//   // https://www.compart.com/fr/unicode/combining/129
//   0x0f71,
//   // 130
//   // https://www.compart.com/fr/unicode/combining/130
//   0x0f72, 0x0f7a, 0x0f7b, 0x0f7c, 0x0f7d, 0x0f80,
//   // 132
//   // https://www.compart.com/fr/unicode/combining/132
//   0x0f74,
// ]);


const DIACRITICS_REG_EXP = /\p{M}+/gu;


let normalizationRegex: RegExp | null = null;
export function normalize(text: string) {
  // The diacritics in the text or in the query can be composed or not.
  // So we use a decomposed text using NFD (and the same for the query)
  // in order to be sure that diacritics are in the same order.

  if (!normalizationRegex) {
    // Compile the regular expression for text normalization once.
    const replace = Object.keys(CHARACTERS_TO_NORMALIZE).join("");
    normalizationRegex = new RegExp(
      `([${replace}])|(\\p{M}+(?:-\\n)?)|(\\S-\\n)|(\\n)`,
      "gum"
    );
  }

  // The goal of this function is to normalize the string and
  // be able to get from an index in the new string the
  // corresponding index in the old string.
  // For example if we have: abCd12ef456gh where C is replaced by ccc
  // and numbers replaced by nothing (it's the case for diacritics), then
  // we'll obtain the normalized string: abcccdefgh.
  // So here the reverse map is: [0,1,2,2,2,3,6,7,11,12].

  // The goal is to obtain the array: [[0, 0], [3, -1], [4, -2],
  // [6, 0], [8, 3]].
  // which can be used like this:
  //  - let say that i is the index in new string and j the index
  //    the old string.
  //  - if i is in [0; 3[ then j = i + 0
  //  - if i is in [3; 4[ then j = i - 1
  //  - if i is in [4; 6[ then j = i - 2
  //  ...
  // Thanks to a binary search it's easy to know where is i and what's the
  // shift.
  // Let say that the last entry in the array is [x, s] and we have a
  // substitution at index y (old string) which will replace o chars by n chars.
  // Firstly, if o === n, then no need to add a new entry: the shift is
  // the same.
  // Secondly, if o < n, then we push the n - o elements:
  // [y - (s - 1), s - 1], [y - (s - 2), s - 2], ...
  // Thirdly, if o > n, then we push the element: [y - (s - n), o + s - n]

  // Collect diacritics length and positions.
  const rawDiacriticsPositions: any[] = [];
  let m;
  while ((m = DIACRITICS_REG_EXP.exec(text)) !== null) {
    rawDiacriticsPositions.push([m[0].length, m.index]);
  }

  let normalized = text.normalize("NFD");
  const positions = [[0, 0]];
  let k = 0;
  let shift = 0;
  let shiftOrigin = 0;
  let eol = 0;
  let hasDiacritics = false;

  normalized = normalized.replace(
    normalizationRegex,
    (match, p1, p2, p3, p4, i) => {
      i -= shiftOrigin;
      if (p1) {
        // Maybe fractions or quotations mark...
        const replacement = CHARACTERS_TO_NORMALIZE[match as keyof typeof CHARACTERS_TO_NORMALIZE];
        const jj = replacement.length;
        for (let j = 1; j < jj; j++) {
          positions.push([i - shift + j, shift - j]);
        }
        shift -= jj - 1;
        return replacement;
      }

      if (p2) {
        const hasTrailingDashEOL = p2.endsWith("\n");
        const len = hasTrailingDashEOL ? p2.length - 2 : p2.length;

        // Diacritics.
        hasDiacritics = true;
        let jj = len;
        if (i + eol === rawDiacriticsPositions[k]?.[1]) {
          jj -= rawDiacriticsPositions[k][0];
          ++k;
        }

        for (let j = 1; j < jj + 1; j++) {
          // i is the position of the first diacritic
          // so (i - 1) is the position for the letter before.
          positions.push([i - 1 - shift + j, shift - j]);
        }
        shift -= jj;
        shiftOrigin += jj;

        if (hasTrailingDashEOL) {
          // Diacritics are followed by a -\n.
          // See comments in `if (p3)` block.
          i += len - 1;
          positions.push([i - shift + 1, 1 + shift]);
          shift += 1;
          shiftOrigin += 1;
          eol += 1;
          return p2.slice(0, len);
        }

        return p2;
      }

      if (p3) {
        // "[A-z]-\n" is removed because an hyphen at the end of a line
        // with not a space before is likely here to mark a break
        // in a word.
        // The \n isn't in the original text so here y = i, n = 1 and o = 2.
        let char = p3.replace('-\n', '')
        if (!/[A-z]/.test(char)) {
          char += '-'
        }
        const step = p3.length - char.length

        positions.push([i - shift + char.length, step + shift])
        shift += step
        shiftOrigin += step
        eol += 1;
        return char
      }

      // p4
      // eol is replaced by space: "foo\nbar" is likely equivalent to
      // "foo bar".
      positions.push([i - shift + 1, shift - 1]);
      shift -= 1;
      shiftOrigin += 1;
      eol += 1;
      return " ";
    }
  );

  positions.push([normalized.length, shift]);

  return [normalized, positions, hasDiacritics];
}