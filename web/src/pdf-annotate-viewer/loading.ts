import { css } from "linaria"
import cbem from "./css/style"

export const LoadingSVG = (size?: number) => {
  size = size || 40
  return `
    <svg version="1.1" id="loading" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" width="${size}px" height="${size}px" viewBox="0 0 50 50" xml:space="preserve">
      <path fill="var(--site-theme-color-primary, #1f71e0)" d="M43.935,25.145c0-10.318-8.364-18.683-18.683-18.683c-10.318,0-18.683,8.365-18.683,18.683h4.068c0-8.071,6.543-14.615,14.615-14.615c8.072,0,14.615,6.543,14.615,14.615H43.935z">
        <animateTransform attributeType="xml"
          attributeName="transform"
          type="rotate"
          from="0 25 25"
          to="360 25 25"
          dur="0.6s"
          repeatCount="indefinite"/>
      </path>
    </svg>
  `
}

export const LoadingBem = cbem('loading')

export const LoadingProgressBem = LoadingBem('progress')



export const loadingClassname = css`
  &${LoadingBem.toSelector()} {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    background: rgba(56, 58, 61, 0.2);

    ${LoadingProgressBem.toSelector()} {
      margin-top: 8px;
      font-size: 13px;
      color: #62738c;
    }
  }

  svg#loading {
    path, rect {
      fill: #1f71e0;
    }
  }
`
