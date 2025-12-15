/* eslint-disable @typescript-eslint/ban-ts-comment */
// @ts-nocheck

import { binarySearchFirstItem } from '@idea/pdfjs-dist/lib/web/ui_utils'



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
const DIACRITICS_EXCEPTION = new Set([
  // UNICODE_COMBINING_CLASS_KANA_VOICING
  // https://www.compart.com/fr/unicode/combining/8
  0x3099, 0x309a,
  // UNICODE_COMBINING_CLASS_VIRAMA (under 0xFFFF)
  // https://www.compart.com/fr/unicode/combining/9
  0x094d, 0x09cd, 0x0a4d, 0x0acd, 0x0b4d, 0x0bcd, 0x0c4d, 0x0ccd, 0x0d3b,
  0x0d3c, 0x0d4d, 0x0dca, 0x0e3a, 0x0eba, 0x0f84, 0x1039, 0x103a, 0x1714,
  0x1734, 0x17d2, 0x1a60, 0x1b44, 0x1baa, 0x1bab, 0x1bf2, 0x1bf3, 0x2d7f,
  0xa806, 0xa82c, 0xa8c4, 0xa953, 0xa9c0, 0xaaf6, 0xabed,
  // 91
  // https://www.compart.com/fr/unicode/combining/91
  0x0c56,
  // 129
  // https://www.compart.com/fr/unicode/combining/129
  0x0f71,
  // 130
  // https://www.compart.com/fr/unicode/combining/130
  0x0f72, 0x0f7a, 0x0f7b, 0x0f7c, 0x0f7d, 0x0f80,
  // 132
  // https://www.compart.com/fr/unicode/combining/132
  0x0f74,
]);
let DIACRITICS_EXCEPTION_STR; // Lazily initialized, see below.

const DIACRITICS_REG_EXP = /\p{M}+/gu;
const SPECIAL_CHARS_REG_EXP =
  /([.*+?^${}()|[\]\\])|(\p{P})|(\s+)|(\p{M})|(\p{L})/gu;
// const NOT_DIACRITIC_FROM_END_REG_EXP = /([^\p{M}])\p{M}*$/u;
// const NOT_DIACRITIC_FROM_START_REG_EXP = /^\p{M}*([^\p{M}])/u;

// The range [AC00-D7AF] corresponds to the Hangul syllables.
// The few other chars are some CJK Compatibility Ideographs.
const SYLLABLES_REG_EXP = /[\uAC00-\uD7AF\uFA6C\uFACF-\uFAD1\uFAD5-\uFAD7]+/g;
const SYLLABLES_LENGTHS = new Map();
// When decomposed (in using NFD) the above syllables will start
// with one of the chars in this regexp.
const FIRST_CHAR_SYLLABLES_REG_EXP =
  "[\\u1100-\\u1112\\ud7a4-\\ud7af\\ud84a\\ud84c\\ud850\\ud854\\ud857\\ud85f]";

const NFKC_CHARS_TO_NORMALIZE = new Map();

let noSyllablesRegExp = null;
let withSyllablesRegExp = null;


export function normalize(text: string) {
  // The diacritics in the text or in the query can be composed or not.
  // So we use a decomposed text using NFD (and the same for the query)
  // in order to be sure that diacritics are in the same order.

  // Collect syllables length and positions.
  const syllablePositions = [];
  let m;
  while ((m = SYLLABLES_REG_EXP.exec(text)) !== null) {
    let { index } = m;
    for (const char of m[0]) {
      let len = SYLLABLES_LENGTHS.get(char);
      if (!len) {
        len = char.normalize("NFD").length;
        SYLLABLES_LENGTHS.set(char, len);
      }
      syllablePositions.push([len, index++]);
    }
  }

  let normalizationRegex;
  if (syllablePositions.length === 0 && noSyllablesRegExp) {
    normalizationRegex = noSyllablesRegExp;
  } else if (syllablePositions.length > 0 && withSyllablesRegExp) {
    normalizationRegex = withSyllablesRegExp;
  } else {
    // Compile the regular expression for text normalization once.
    const replace = Object.keys(CHARACTERS_TO_NORMALIZE).join("");
    const toNormalizeWithNFKC =
      "\u2460-\u2473" + // Circled numbers.
      "\u24b6-\u24ff" + // Circled letters/numbers.
      "\u3244-\u32bf" + // Circled ideograms/numbers.
      "\u32d0-\u32fe" + // Circled ideograms.
      "\uff00-\uffef"; // Halfwidth, fullwidth forms.

    // 3040-309F: Hiragana
    // 30A0-30FF: Katakana
    const CJK = "(?:\\p{Ideographic}|[\u3040-\u30FF])";
    const regexp = `([${replace}])|([${toNormalizeWithNFKC}])|(\\p{M}+(?:-\\n)?)|(\\S-\\n)|(${CJK}\\n)|(\\n)`;

    if (syllablePositions.length === 0) {
      // Most of the syllables belong to Hangul so there are no need
      // to search for them in a non-Hangul document.
      // We use the \0 in order to have the same number of groups.
      normalizationRegex = noSyllablesRegExp = new RegExp(
        regexp + "|(\\u0000)",
        "gum"
      );
    } else {
      normalizationRegex = withSyllablesRegExp = new RegExp(
        regexp + `|(${FIRST_CHAR_SYLLABLES_REG_EXP})`,
        "gum"
      );
    }
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
  const rawDiacriticsPositions = [];
  while ((m = DIACRITICS_REG_EXP.exec(text)) !== null) {
    rawDiacriticsPositions.push([m[0].length, m.index]);
  }

  let normalized = text.normalize("NFD");
  const positions = [[0, 0]];
  let rawDiacriticsIndex = 0;
  let syllableIndex = 0;
  let shift = 0;
  let shiftOrigin = 0;
  let eol = 0;
  let hasDiacritics = false;

  normalized = normalized.replace(
    normalizationRegex,
    (match, p1, p2, p3, p4, p5, p6, p7, i) => {
      i -= shiftOrigin;
      if (p1) {
        // Maybe fractions or quotations mark...
        const replacement = CHARACTERS_TO_NORMALIZE[p1];
        const jj = replacement.length;
        for (let j = 1; j < jj; j++) {
          positions.push([i - shift + j, shift - j]);
        }
        shift -= jj - 1;
        return replacement;
      }

      if (p2) {
        // Use the NFKC representation to normalize the char.
        let replacement = NFKC_CHARS_TO_NORMALIZE.get(p2);
        if (!replacement) {
          replacement = p2.normalize("NFKC");
          NFKC_CHARS_TO_NORMALIZE.set(p2, replacement);
        }
        const jj = replacement.length;
        for (let j = 1; j < jj; j++) {
          positions.push([i - shift + j, shift - j]);
        }
        shift -= jj - 1;
        return replacement;
      }

      if (p3) {
        const hasTrailingDashEOL = p3.endsWith("\n");
        const len = hasTrailingDashEOL ? p3.length - 2 : p3.length;

        // Diacritics.
        hasDiacritics = true;
        let jj = len;
        if (i + eol === rawDiacriticsPositions[rawDiacriticsIndex]?.[1]) {
          jj -= rawDiacriticsPositions[rawDiacriticsIndex][0];
          ++rawDiacriticsIndex;
        }

        for (let j = 1; j <= jj; j++) {
          // i is the position of the first diacritic
          // so (i - 1) is the position for the letter before.
          positions.push([i - 1 - shift + j, shift - j]);
        }
        shift -= jj;
        shiftOrigin += jj;

        if (hasTrailingDashEOL) {
          // Diacritics are followed by a -\n.
          // See comments in `if (p4)` block.
          i += len - 1;
          positions.push([i - shift + 1, 1 + shift]);
          shift += 1;
          shiftOrigin += 1;
          eol += 1;
          return p3.slice(0, len);
        }

        return p3;
      }

      if (p4) {
        // "X-\n" is removed because an hyphen at the end of a line
        // with not a space before is likely here to mark a break
        // in a word.
        // The \n isn't in the original text so here y = i, n = 1 and o = 2.
        positions.push([i - shift + 1, 1 + shift]);
        shift += 1;
        shiftOrigin += 1;
        eol += 1;
        return p4.charAt(0);
      }

      if (p5) {
        // An ideographic at the end of a line doesn't imply adding an extra
        // white space.
        positions.push([i - shift + 1, shift]);
        shiftOrigin += 1;
        eol += 1;
        return p5.charAt(0);
      }

      if (p6) {
        // eol is replaced by space: "foo\nbar" is likely equivalent to
        // "foo bar".
        positions.push([i - shift + 1, shift - 1]);
        shift -= 1;
        shiftOrigin += 1;
        eol += 1;
        return " ";
      }

      // p7
      if (i + eol === syllablePositions[syllableIndex]?.[1]) {
        // A syllable (1 char) is replaced with several chars (n) so
        // newCharsLen = n - 1.
        const newCharLen = syllablePositions[syllableIndex][0] - 1;
        ++syllableIndex;
        for (let j = 1; j <= newCharLen; j++) {
          positions.push([i - (shift - j), shift - j]);
        }
        shift -= newCharLen;
        shiftOrigin += newCharLen;
      }
      return p7;
    }
  );

  positions.push([normalized.length, shift]);

  return [normalized, positions, hasDiacritics];
}

// Determine the original, non-normalized, match index such that highlighting of
// search results is correct in the `textLayer` for strings containing e.g. "½"
// characters; essentially "inverting" the result of the `normalize` function.
export function getOriginalIndex(diffs, pos, len) {
  if (!diffs) {
    return [pos, len];
  }

  const start = pos;
  const end = pos + len;
  let i = binarySearchFirstItem(diffs, x => x[0] >= start);
  if (diffs[i][0] > start) {
    --i;
  }

  let j = binarySearchFirstItem(diffs, x => x[0] >= end, i);
  if (diffs[j][0] > end) {
    --j;
  }

  return [start + diffs[i][1], len + diffs[j][1] - diffs[i][1]];
}

export function convertToRegExpString(query: string, hasDiacritics: boolean, matchDiacritics: boolean) {
  let isUnicode = false;
  query = query.replace(
    SPECIAL_CHARS_REG_EXP,
    (
      match,
      p1 /* to escape */,
      p2 /* punctuation */,
      p3 /* whitespaces */,
      p4 /* diacritics */,
      p5 /* letters */
    ) => {
      // We don't need to use a \s for whitespaces since all the different
      // kind of whitespaces are replaced by a single " ".

      if (p1) {
        // Escape characters like *+?... to not interfer with regexp syntax.
        return `[ ]*\\${p1}[ ]*`;
      }
      if (p2) {
        // Allow whitespaces around punctuation signs.
        return `[ ]*${p2}[ ]*`;
      }
      if (p3) {
        // Replace spaces by \s+ to be sure to match any spaces.
        return "[ ]+";
      }
      if (matchDiacritics) {
        return p4 || p5;
      }

      if (p4) {
        // Diacritics are removed with few exceptions.
        return DIACRITICS_EXCEPTION.has(p4.charCodeAt(0)) ? p4 : "";
      }

      // A letter has been matched and it can be followed by any diacritics
      // in normalized text.
      if (hasDiacritics) {
        isUnicode = true;
        return `${p5}\\p{M}*`;
      }
      return p5;
    }
  );

  const trailingSpaces = "[ ]*";
  if (query.endsWith(trailingSpaces)) {
    // The [ ]* has been added in order to help to match "foo . bar" but
    // it doesn't make sense to match some whitespaces after the dot
    // when it's the last character.
    query = query.slice(0, query.length - trailingSpaces.length);
  }

  if (matchDiacritics) {
    // aX must not match aXY.
    if (hasDiacritics) {
      DIACRITICS_EXCEPTION_STR ||= String.fromCharCode(
        ...DIACRITICS_EXCEPTION
      );

      isUnicode = true;
      query = `${query}(?=[${DIACRITICS_EXCEPTION_STR}]|[^\\p{M}]|$)`;
    }
  }

  return [isUnicode, query];
}

// /**
//  * 
//  * @param words 需要查询坐标的单词
//  */
// const matchChars = async (words: string[], viewerController: ViewerController) => {
//   const startTimestamp = Date.now()
//   const p: Promise<any>[] = []
//   const pdfDocument = viewerController.getPdfDocument() as PDFDocumentProxy
//   const pageNums = pdfDocument.numPages
//   for (let i = 1; i <= pageNums; i++) {
//     p.push(pdfDocument.getPage(i).then(pdfPage => pdfPage.getTextContent()).then(textContent => {
//       const strBuf = [];
//       for (const textItem of textContent.items) {
//         strBuf.push(textItem.str);
//         if (textItem.hasEOL) {
//           strBuf.push("\n");
//         }
//       }
//        // Store the normalized page content (text items) as one string.
//       const [
//         pageContent,
//         pageDiff,
//         hasDiacritic,
//       ] = normalize(strBuf.join(""));
//       return {
//         content: pageContent,
//         diffs: pageDiff,
//         hasDiacritics: hasDiacritic,
//       };
//     }, reason => {
//       console.error(
//         `Unable to get text content for page ${i + 1}`,
//         reason
//       );
//       return {
//         content: '',
//         diffs: [],
//         hasDiacritics: false,
//       }
//     }))
//   }

//   const pageContents = await Promise.all(p);

//   const caseSensitive = false;

//   const pdfViewer = viewerController.getDocumentViewer()!.getPdfViewer()

//   const results: { 
//     matches: {
//       [pageIndex: number]: { matches: number[]; matchesLength: number[], rects: TextRectCoordinate[] }
//     }, 
//     word: string,
//     viewport: PageViewport 
//   }[] = [];
//   pageContents.forEach((pageContent, pageIndex) => {
//     const diffs = pageContent.diffs;

//     const textHighlighter: TextHighlighter | null =
//       pdfViewer.getPageView(pageIndex)?.textLayer?.highlighter;


//     if (!textHighlighter) {
//       // 说明当前页没有被渲染出来
//       return;
//     }
    
//     words.forEach((word, wordIndex) => {
//       const [isUnicode, q] = convertToRegExpString(word, pageContent.hasDiacritics, false)
//       const flags = `g${isUnicode ? "u" : ""}${caseSensitive ? "" : "i"}`;
//       const query = q ? new RegExp(q as string, flags) : null;
//       if (!query) {
//         return null;
//       }
//       const wordResults = {
//         matches: [] as number[],
//         matchesLength: [] as number[],
//       };
//       let match;
//       while ((match = query.exec(pageContent.content)) !== null) {
//         const [matchPos, matchLen] = getOriginalIndex(
//           diffs,
//           match.index,
//           match[0].length
//         );
//         if (matchLen) {
//           wordResults.matches.push(matchPos);
//           wordResults.matchesLength.push(matchLen);
//         }
//       }
//       results[wordIndex] = results[wordIndex] || {
//         matches: {
//           [pageIndex]: {
//             matches: [] as number[],
//             matchesLength: [] as number[],
//             rects: [] as TextRectCoordinate[],
//             viewport: textHighlighter.viewport!,
//           },
//         },
//         word,
//       }
//       results[wordIndex].matches[pageIndex] = {
//         matches: wordResults.matches,
//         matchesLength: wordResults.matchesLength,
//         rects: textHighlighter._convertMatches(wordResults.matches, wordResults.matchesLength),
//       };
//     });
//   })

//   const endTimestamp = Date.now()
//   console.log(results, endTimestamp - startTimestamp)

// }

// export default matchChars