import { defineStore } from 'pinia'

import {
  GetDocIndexResponse,
  SimpleUserDocInfo,
  UserDocFolderInfo,
} from 'go-sea-proto/gen/ts/doc/userDoc'
import { ref } from 'vue'
import { getDocIndex } from '@/api/document'

export interface BreadCrumb {
  docName: string
  folderId: string
}

export interface LiteratureNode {
  children: LiteratureNode[]
  isLeaf: boolean
  key: string
  title: string
  paperId?: string
  pdfId?: string
  isEditing: boolean
  hoverType: number
  docCount: number
  disabled?: boolean
  docId?: string
}

export interface LiteratureExtra {
  path: number[]
  count: number
  width: number
}

export type LibraryIndexExtraMap = Record<
  LiteratureNode['key'],
  LiteratureExtra
>

export interface LiteratureNodeFamily {
  parent: LiteratureNode | null
  children: LiteratureNode[]
  node: LiteratureNode
  index: number
}

export const noteDirty = {
  is: false,
}

export const convertFolder = (
  folder: UserDocFolderInfo,
  prefix = ''
): LiteratureNode => {
  const key = `${prefix}-${folder.folderId}`
  const childrenFolders =
    folder.childrenFolders?.map((child) => convertFolder(child, key)) ?? []
  const childrenDocs =
    folder.docInfos?.map((child) => convertDocument(child, key)) ?? []
  return {
    isLeaf: false,
    docCount: folder.docCount,
    children: [...childrenFolders, ...childrenDocs],
    key,
    title: folder.name,
    isEditing: false,
    hoverType: 0,
  }
}

export const convertDocument = (
  doc: SimpleUserDocInfo,
  prefix = ''
): LiteratureNode => {
  return {
    isLeaf: true,
    docCount: 0,
    children: [],
    key: `${prefix}-${doc.docId}`,
    title: doc.docName,
    paperId: doc.paperId,
    pdfId: doc.pdfId,
    isEditing: false,
    hoverType: 0,
    docId: doc.docId,
  }
}

export const seperatorKey = '07'
export const allKey = '0'
export const newFolderKey = 'literature_tree_new_folder'
export const getAllNode = (totalDocCount = 0): LiteratureNode => ({
  isLeaf: false,
  children: [],
  key: allKey,
  title: '全部文献',
  isEditing: false,
  hoverType: 0,
  docCount: totalDocCount,
})

export enum DropToPosition {
  before = -1,
  middle = 0,
  after = 1,
}

export const FOLDER_INDENT = 215

export interface LibraryIndex {
  selectedKey: string
  libraryIndexList: LiteratureNode[]
  uploaderVisible: boolean
}

export const useLibraryIndex = defineStore('libraryIndex', {
  state(): LibraryIndex {
    return {
      selectedKey: '',
      libraryIndexList: [getAllNode()],
      uploaderVisible: false,
    }
  },
  getters: {
    libraryIndexExtra() {
      const extraMap: LibraryIndexExtraMap = {}
      const path: number[] = []
      // 回溯算法 DFS 后序遍历
      const backtrack = (menu: LiteratureNode, index: number): number => {
        let count = 0
        path.push(index)

        const width = FOLDER_INDENT - (path.length - 1) * 19
        if (menu.isLeaf) {
          count += 1
          extraMap[menu.key] = {
            path: [...path],
            count: 1,
            width,
          }
        } else {
          menu.children.map(backtrack).forEach((cnt) => (count += cnt))
          extraMap[menu.key] = {
            path: [...path],
            count,
            width,
          }
        }

        path.pop()
        return count
      }

      let allCount = 0
      this.libraryIndexList.map(backtrack).forEach((cnt) => (allCount += cnt))
      extraMap[allKey] = {
        path: [],
        count: allCount,
        width: FOLDER_INDENT,
      }

      return extraMap
    },
    breadCrumbList() {
      const list: BreadCrumb[] = [
        {
          docName: getAllNode().title,
          folderId: allKey,
        },
      ]

      if (
        !this.selectedKey ||
        this.selectedKey === allKey ||
        this.selectedKey === seperatorKey
      ) {
        return list
      }

      const extra = this.libraryIndexExtra[this.selectedKey]
      if (!extra) {
        return list
      }

      const path = extra.path.slice()
      const last = path.pop() as number
      let children = this.libraryIndexList
      for (const p of path) {
        const folder = children[p]
        list.push({
          docName: folder.title,
          folderId: folder.key,
        })
        children = folder.children
      }

      if (!children[last].isLeaf) {
        list.push({
          docName: children[last].title,
          folderId: children[last].key,
        })
      }

      return list
    },
    currentFolder(): BreadCrumb {
      return this.breadCrumbList.slice(-1).pop() as BreadCrumb
    },
    rawFolderId(): string {
      return removePrefix(this.currentFolder.folderId)
    },
  },
  actions: {
    async fetchLibraryIndex() {
      let res: GetDocIndexResponse
      try {
        res = await getDocIndex({})
      } catch (error) {
        // TODO
        return
      }

      const folderList = res.folderInfos.map((child) => convertFolder(child))
      const docList = res.unclassifiedDocInfos.map((child) =>
        convertDocument(child)
      )
      const all = getAllNode(res.totalDocCount)
      const seperator: LiteratureNode = {
        isLeaf: false,
        children: [],
        title: '以下文献尚未分类',
        key: seperatorKey,
        isEditing: false,
        hoverType: 0,
        docCount: 0,
        disabled: true,
      }

      this.libraryIndexList = [all, ...folderList, seperator, ...docList]
    },
  },
})

export const removePrefix = (key: string) => {
  return key.split('-').pop() as string
}

export const useRemoveModalVisible = () => {
  const removeModalVisible = ref(false)

  return {
    removeModalVisible,
  }
}
