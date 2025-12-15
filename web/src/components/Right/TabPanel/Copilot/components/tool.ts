import markdownIt from 'markdown-it';
import textMath from 'markdown-it-texmath';
import katex from 'katex';

export const formatQuestion = (answer: string, anchors?: string[]) => {
  if (!anchors?.length) {
    return answer;
  }
  anchors.forEach((anchor) => {
    const idx = answer.indexOf(anchor);
    if (idx === -1) {
      return;
    }
    const prevChar = answer[idx - 1];
    if (prevChar === '\n') {
      answer = answer.replace(
        anchor,
        `<a class="js-copilot-question copilot-question-anchor">${anchor}</a>`
      );
    } else {
      answer = answer.replace(
        anchor,
        `<div><a class="js-copilot-question copilot-question-anchor">${anchor}</a></div>`
      );
    }
  });
  return answer;
};

textMath.katex = katex;

export const md = markdownIt({
  html: true,
  breaks: true,
});

md.use(textMath);
