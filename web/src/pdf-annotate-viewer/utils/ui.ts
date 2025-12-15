import { css } from "linaria"

const containerClassname = css`
  -webkit-overflow-scrolling: touch;
  left: 0;
  right: 0;
  bottom: 0;
  top: 0;
  position: absolute;
`

export const createContainer = (wrapper: HTMLDivElement) => {
  const container = document.createElement('div')
  //container.classList.add(containerClassname);
  // 使用内联样式:TODO
  Object.assign(container.style, {
    position: 'absolute',
    left: '0',
    right: '0',
    top: '0',
    bottom: '0',
    webkitOverflowScrolling: 'touch',
  });
  wrapper.appendChild(container)
  wrapper.style.position = 'relative'
  return container
}