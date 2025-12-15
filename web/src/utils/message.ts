export const postMessage = (params: { event: string; params: any }) => {
  top?.postMessage(params, '*')
}

export const ELECTRON_CLIENT_EVENT = {
  WS_MESSAGE: 'electron-client-event-ws-message',
  LOGIN_CHANGE: 'electron-client-event-login-change',
}

export enum ELECTRON_LOGIN_STATUS {
  Login,
  Logout,
}

export enum ELECTRON_WS_MESSAGE_SCENE {
  docChange = 'doc_change',
  groupChange = 'group_change',
  folderChange = 'folder_change',
  noteTabActive = 'note_tab_active',
  officialMsg = 'OFFICIAL_MESSAGE',
}

export const useWsMessage = (
  scene: ELECTRON_WS_MESSAGE_SCENE,
  callback: (json: any) => void
) => {
  const emptyFunction = () => {}

  if (typeof window === 'undefined') {
    return emptyFunction
  }

  const { promiseIpc } = window as any
  if (!promiseIpc) {
    return emptyFunction
  }

  const listener = (data: { scene: string; content: string }) => {
    let json: null | { type: string } = null
    try {
      json = JSON.parse(data.content)
    } catch (error) {
      console.error('electron消息解析出错', error)
      return
    }

    if (json && json.type === scene) {
      callback(json)
    }
  }

  promiseIpc.on(ELECTRON_CLIENT_EVENT.WS_MESSAGE, listener)
  return () => {
    promiseIpc.off(ELECTRON_CLIENT_EVENT.WS_MESSAGE, listener)
  }
}
