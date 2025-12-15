export interface NoteSlice {
  page: number;
  pageEl?: HTMLElement;
  uuid: string;
  // 因PDF改变主题、大小等事件会触发DOM重渲染
  // 会干掉NotesSvg通过teleport渲染到PDF的元素
  // 无法用v-show缓存DOM节点
  // hidden?: boolean
  locked?: boolean;
}
