import cbem from './style';

export const ViewerBem = cbem('viewer');

export const classname = ViewerBem.toString();

// 定义一个样式标签，将原来的 CSS 内容添加到页面中
const styleTag = document.createElement('style');
styleTag.textContent = `
  background-color: #52565a;
  height: 100%;
  &:hover.large-scrollbar-y {
    ::-webkit-scrollbar:vertical {
      width: 12px;
      height: 12px;
    }
  }

  &:hover.large-scrollbar-x {
    ::-webkit-scrollbar:horizontal {
      width: 12px;
      height: 12px;
    }
  }

  &:hover::-webkit-scrollbar {
    width: 8px;
    height: 16px;
  }

  ::-webkit-scrollbar {
    width: 8px;
    height: 8px;
  }
  ::-webkit-scrollbar-button {
    display: none;
  }
  ::-webkit-scrollbar-corner {
    display: none;
  }
  /*定义滚动条轨道
    内阴影+圆角*/
  ::-webkit-scrollbar-track {
    /* -webkit-box-shadow:inset 0 0 6px rgba(0,0,0,0.3); */
    /* border-radius:10px; */
    box-shadow: none;
    background-color: transparent;
  }
  /* 定义滑块 */
  ::-webkit-scrollbar-thumb {
    border-radius: 6px;
    /* -webkit-box-shadow:inset 0 0 6px rgba(0,0,0,.3); */
    background-color: #aaa;
    transition: background-color 0.2s linear, width 0.2s ease-in-out;
    &:hover {
      background-color: #999;
    }
  }

  .pdfViewer .page {
    margin: 0 auto;
    border: none !important;
    border-bottom: 1px solid #000 !important;
    border-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAYAAAByUDbMAAAA1ElEQVQ4jbWUWw6EIAxFy2NFs/8NzR4UJhpqLsdi5mOmSSMUOfYWqv3S0gMr4XlYH/64gZa/gN3ANYA7KAXALt4ktoQ5MI9YxqaG8bWmsIysMuT6piSQCa4whZThCu8CM4zP9YJaKci9jicPq3NcBWYoPMGUlhG7ivtkB+gVyFY75wXghOvh8t5mto1Mdim6e+MBqH6XsY+YAwjpq3vGF7weTWQptLEDVCZvPTMl5JZZsdh47FHW6qFMyvLYqjcnmdFfY9Xk/KDOlzCusX2mi/ofM7MPkzBcSp4Q1/wAAAAASUVORK5CYII=) 6 6 repeat !important;
        
  }

  .ps__rail-x {
    z-index: 999;
  }
  .ps__rail-y {
    z-index: 999;
  }
`;

// 当文档准备好时添加样式
if (typeof document !== 'undefined') {
  document.head.appendChild(styleTag);
}
