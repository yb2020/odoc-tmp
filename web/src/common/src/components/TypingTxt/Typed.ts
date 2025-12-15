import mitt from 'mitt';

interface TypedOptions {
  strings: Typed['strings'][number][];
  typeSpeed?: number;
  startDelay?: number;
  pauseTime?: number;
  typingTime?: number;
}

export type TypingPayload = {
  curString: string;
  fullString: string;
  curIndex: number;
};

type TypedEventMap = Record<'onStop', void> &
  Record<'onBegin', void> &
  Record<'onTyping', TypingPayload>;

export const wait = (ms: number) =>
  new Promise((resolve) => {
    window.setTimeout(() => {
      resolve(true);
    }, ms);
  });

const noop = () => {
  // do nothing
};

const humanizer = (speed: number) => {
  return Math.round((Math.random() * speed) / 2) + speed;
};

class TypingOnString {
  private onEndHandler = noop;
  private onTypingHandler = (str: string) => {
    // do nothing
    console.log(str);
  };

  private curIndex = 0;

  private timer = 0;

  private string = '';
  private speed = 0;
  private time = 0;

  private starTime = Date.now();

  constructor(str: string, speed: number, time?: number) {
    this.string = str;
    this.speed = speed ?? 0;
    this.time = time ?? 0;
  }

  typing() {
    let humanizedSpeed = humanizer(this.speed);
    const leftTime = this.time - Math.ceil(Date.now() - this.starTime);
    const leftLength = this.string.length - this.curIndex;
    let numChars =
      this.time === -1
        ? 1
        : leftTime <= 0
          ? this.string.length
          : Math.floor(leftLength / (leftTime / (this.speed || 1)));
    if (numChars <= 0) {
      numChars = 1;
    }
    if (this.time !== -1 && humanizedSpeed >= leftTime) {
      numChars = this.string.length;
      humanizedSpeed = leftTime;
    }

    this.curIndex += numChars;
    this.onTypingHandler(this.string.slice(0, this.curIndex));

    if (this.curIndex >= this.string.length) {
      this.onEndHandler();
      return;
    }

    this.timer = window.setTimeout(() => {
      this.typing();
    }, humanizedSpeed);
  }

  destroy() {
    if (this.timer) {
      window.clearTimeout(this.timer);
    }
  }

  onEnd(cb: TypingOnString['onEndHandler']) {
    this.onEndHandler = cb;
  }

  onTyping(cb: TypingOnString['onTypingHandler']) {
    this.onTypingHandler = cb;
  }
}

export class Typed {
  private typeSpeed = 0;
  private startDelay = 0;
  private typingTime = 0;

  private strings: {
    string: string; // 打字内容
    typeSpeed?: number; // 打字速度，设置后表示每个字的打字速度
    typingTime?: number; // 打字时间，设置后表示必须在这个time时间内打完
    pauseTime?: number; // typing完毕后的暂停时间
  }[] = [];

  private timer = 0;

  private emitter = mitt<TypedEventMap>();

  private curRenderedIndex = -1;

  private cursor: HTMLSpanElement | null = null;

  private typingOnString: TypingOnString | null = null;

  constructor(options: TypedOptions) {
    this.strings = options.strings;
    this.typeSpeed = options.typeSpeed || this.typeSpeed;
    this.startDelay = options.startDelay || this.startDelay;
    this.typingTime = options.typingTime || this.typingTime;
    this.begin();
  }

  private begin() {
    this.emitter.emit('onBegin');
    const startDelay = this.startDelay;
    this.timer = window.setTimeout(() => {
      this.typingLine();
    }, startDelay);
  }

  private typingLine() {
    if (this.curRenderedIndex === this.strings.length - 1) {
      this.emitter.emit('onStop');
      return;
    }
    // 如果当前正在打字，就不要创建新的打字了
    if (this.typingOnString) {
      return;
    }

    const startTime = Date.now();

    this.curRenderedIndex++;
    const str = this.strings[this.curRenderedIndex];
    const typeSpeed = str.typeSpeed || this.typeSpeed;
    const pauseTime = str.pauseTime || 0;
    const typingTime = str.typingTime || this.typingTime;
    const typing = new TypingOnString(str.string, typeSpeed, typingTime);
    this.typingOnString = typing;
    typing.onTyping((curString) => {
      const fullString =
        this.strings
          .slice(0, this.curRenderedIndex)
          .map((str) => str.string)
          .join('') + curString;
      this.emitter.emit('onTyping', {
        curString: curString,
        fullString,
        curIndex: this.curRenderedIndex,
      });
    });
    typing.onEnd(() => {
      console.log(str.string, Date.now() - startTime);
      this.typingOnString = null;
      this.timer = window.setTimeout(() => {
        this.typingLine();
      }, pauseTime);
    });
    typing.typing();
  }

  flushStrings(str: Typed['strings'][number]) {
    this.strings.push(str);
    // flushStrings后，如果当前没有在打字，就开始打字
    this.typingLine();
  }

  destroy() {
    if (this.timer) {
      window.clearTimeout(this.timer);
    }
    if (this.typingOnString) {
      this.typingOnString.destroy();
    }
  }

  onTyping(handler: (payload: TypingPayload) => void) {
    this.emitter.on('onTyping', handler);
    return this;
  }

  onStop(handler: () => void) {
    this.emitter.on('onStop', handler);
    return this;
  }

  insertCursor(container: HTMLElement) {
    this.cursor = document.createElement('span');
    this.cursor.classList.add('typed-cursor');
    this.cursor.innerHTML = '|';
    this.cursor.setAttribute('aria-hidden', 'true');
    container.appendChild(this.cursor);
    return this;
  }

  blinkingCursor(blinking: boolean) {
    if (!this.cursor) {
      return;
    }
    if (blinking) {
      this.cursor.style.display = 'inline';
    } else {
      this.cursor.style.display = 'none';
    }
  }

  private static _hasInit = false;

  static init() {
    if (Typed._hasInit) {
      return;
    }
    const style = document.createElement('style');
    style.setAttribute('type', 'text/css');
    style.innerHTML = `
      .typed-cursor {
        opacity: 1;
        animation: typedjsBlink 0.7s infinite;
        -webkit-animation: typedjsBlink 0.7s infinite;
      }

      @keyframes typedjsBlink {
        0% { opacity: 1; }
        50% { opacity: 0; }
        100% { opacity: 1; }
      }
    `;
    document.head.appendChild(style);
    Typed._hasInit = true;
  }
}

Typed.init();
