import { levenshtein } from 'wuzzy';

export const splitIntoWords = (s: string) => {
  // 仅支持中文分字/用空格及符号分词，其他语言作为整句处理
  return s
    .replaceAll(/[\u4E00-\u9FFF]+/g, (match) => {
      return `#${[...match].join('|')}#`;
    })
    .split(/[\s`~!@#$%^&*()=+[\]{}\\|;:'",.<>/?]/)
    .filter(Boolean);
};

export const findIdxBySimilarity = (
  query: string,
  content: string,
  similarity = 80,
  suffixShift: number | boolean = false
) => {
  const l = query.length;
  const words = splitIntoWords(query);
  // 先找出开头单词（字）匹配位置列表
  let start = 0;
  let bgIdx = -1;
  let endIdx = -1;
  while (start + l <= content.length) {
    const idx = content.indexOf(words[0], start);
    if (idx !== -1) {
      const str = content.substring(idx, idx + l);
      // const strWords = this.splitIntoWords(str);
      const rate = levenshtein(splitIntoWords(str), splitIntoWords(query));
      if (rate * 100 >= similarity) {
        bgIdx = idx;
        endIdx = idx + l;
        if (suffixShift) {
          const lastWord = words[words.length - 1];
          const offset =
            typeof suffixShift === 'number' && suffixShift > 0
              ? suffixShift
              : 10;
          for (let i = 0; i < offset && i + endIdx < content.length; i++) {
            const end = endIdx + i;
            const suffix = content.substring(end - lastWord.length, end);
            if (suffix === lastWord) {
              endIdx = end;
              break;
            }
          }
        }
        break;
      }
      start = idx + 1;
    } else {
      break;
    }
  }
  return {
    bgIdx,
    endIdx,
  };
};
