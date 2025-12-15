export enum PAGE_ROUTE_NAME {
  WORKBENCH = 'workbench',//工作台，相当于登录后的首页
  LOGIN = 'login',
  NOTE = 'note', //笔记页
  EXCEPTION = 'error',
  EXCEPTION_400 = '400',
  EXCEPTION_404 = '404',
  EXCEPTION_SERVER_ERROR = '500',
  EXCEPTION_PAGENOTFOUND = 'PageNotFound',
  FORBIDDEN = 'forbidden',
  CHAT = 'chat',
  WRITE = 'write',
  LIBRARY = 'library', // 文献库页面
  LIBRARY_IN_WORKBENCH = 'library_in_workbench', // 工作台中的文献库页面
  NOTES = 'notes', // 笔记列表页面
  NOTES_IN_WORKBENCH = 'notes_in_workbench', // 工作台中的笔记列表页面
  RECENT_READING_IN_WORKBENCH = 'recent_reading_in_workbench', // 工作台中的最近阅读页面
  STRIPEPAY = 'stripepay',
}
