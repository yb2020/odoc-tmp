import { Module, Store as VuexStore, MutationTree, CommitOptions, ActionContext, ActionTree } from 'vuex';
import { FigureTippyTriggerItem, ParseState, PdfFigureAndTableInfosType } from './type';
import { RootState } from '../types';
import { createFigureTippyVue, destroyTippyVues } from '~/src/dom/tippy';
import {  getFiguresAndTablesFinal } from "@/api/parse";

export enum ParseMutationTypes {
  SET_FIGURE_VIEWER_ITEM = 'setFigureViewerItem',
  // DEL_FIGURE_VIEWER_ITEM = 'delFigureViewerItem',
  // CLEAR_FIGURE_VIEWER_ITEM = 'clearFigureViewerItem',
  SET_PDF_FIGURE_AND_TABLE_INFOS = 'setPdfFigureAndTableInfos',
}

export type Mutations<S = ParseState> = {
  [ParseMutationTypes.SET_FIGURE_VIEWER_ITEM](state: S, payload: FigureTippyTriggerItem): void;
  // [ParseMutationTypes.DEL_FIGURE_VIEWER_ITEM](state: S, payload: string): void;
  [ParseMutationTypes.SET_PDF_FIGURE_AND_TABLE_INFOS](state: S, payload: {
    pdfId: string;
    infos: Partial<PdfFigureAndTableInfosType>
  }): void;
  // [ParseMutationTypes.CLEAR_FIGURE_VIEWER_ITEM](state: S): void;
};

const mutations: MutationTree<ParseState> & Mutations = {
  [ParseMutationTypes.SET_FIGURE_VIEWER_ITEM](state, payload) {
    const id = payload.id  // 现在 id 是 url
    const infos = state.pdfFigureAndTableInfos[payload.pdfId]
    
    // 优先使用 url 查找，如果 id 看起来像 url（包含 http），则用 url 查找
    // 否则用 refIdx 查找（兼容 markers.ts 的调用）
    const item = id.startsWith('http') 
      ? infos?.list?.find(item => item.url === id)
      : infos?.list?.find(item => item.refIdx === id)
      
    if (!item) {
      return
    }
    const item0 = infos?.list0?.find(item => item.refIdx === item.refIdx)
    const exit = state.figureTippyTriggers.find((current) => current.id === id && current.pdfId === payload.pdfId);
    if (exit) {
      const isDing = exit.item.tippy.popper?.querySelector('em.icondingzhu')
      if (isDing || exit.triggerEle === payload.triggerEle) {
        exit.item.tippy.show()
        return;
      }
      destroyTippyVues([exit.item])
      const idx = state.figureTippyTriggers.findIndex(item => item === exit)
      state.figureTippyTriggers.splice(idx, 1);
    }
    const tippyItem = createFigureTippyVue({ info: {
      ...item,
      refContent: item0?.refContent ?? '',
    }, triggerEle: payload.triggerEle })
    state.figureTippyTriggers.push({
      pdfId: payload.pdfId,      
      id: id,
      item: tippyItem,
      triggerEle: payload.triggerEle,
    });

  },
  
  [ParseMutationTypes.SET_PDF_FIGURE_AND_TABLE_INFOS](state, payload) {
    const infos = state.pdfFigureAndTableInfos[payload.pdfId] || {}
    state.pdfFigureAndTableInfos[payload.pdfId] = {
      ...infos,
      ...payload.infos,
    }
  },
};

export type ParseStore<S = ParseState> = Omit<VuexStore<S>, 'commit'> & {
  commit<K extends keyof Mutations, P extends Parameters<Mutations[K]>[1]>(
    key: K,
    payload: P,
    options?: CommitOptions
  ): ReturnType<Mutations[K]>;
};

type AugmentedActionContext = {
  commit<K extends keyof Mutations>(
    key: K,
    payload: Parameters<Mutations[K]>[1]
  ): ReturnType<Mutations[K]>;
} & Omit<ActionContext<ParseState, RootState>, 'commit'>;

export enum ParseActionTypes {
  GET_PDF_FIGURE_AND_TABLE_INFOS = 'getPdfFigureAndTableInfos',
}

interface Actions {
  [ParseActionTypes.GET_PDF_FIGURE_AND_TABLE_INFOS]: (context: AugmentedActionContext, pdfId: string) => Promise<void>;
}

const actions: ActionTree<ParseState, RootState> & Actions = {
  async [ParseActionTypes.GET_PDF_FIGURE_AND_TABLE_INFOS]({ commit }, pdfId) {
    commit(ParseMutationTypes.SET_PDF_FIGURE_AND_TABLE_INFOS, { pdfId, infos: { pending: true, error: null } })
    try {
      const res = await getFiguresAndTablesFinal({
        pdfId,
        pageReq: {
          pageSize: 100,
          pageNum: 1,
        },
      });
      commit(ParseMutationTypes.SET_PDF_FIGURE_AND_TABLE_INFOS, {
        pdfId,
        infos: {
          list: res.list,
          pending: false,
        },
      });
    } catch (error) {
      console.error(error)
      commit(ParseMutationTypes.SET_PDF_FIGURE_AND_TABLE_INFOS, {
        pdfId,
        infos: {
          error: error as Error,
          pending: false,
        },
      });
    }
  },
};

export const ParseModule: Module<ParseState, RootState> = {
  namespaced: true,
  state: () => ({
    figureTippyTriggers: [],
    pdfFigureAndTableInfos: {},
  }),
  mutations,
  actions,
};
