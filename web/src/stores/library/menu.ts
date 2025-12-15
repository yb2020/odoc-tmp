// import {
//   MoveFolderOrDocReq,
//   MoveDocOrFolderToAnotherFolderReq,
//   CopyDocOrFolderToAnotherFolderReq,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/client/ClientDoc'
import {
  MoveFolderOrDocReq,
  MoveDocOrFolderToAnotherFolderReq,
  CopyDocOrFolderToAnotherFolderReq,
} from 'go-sea-proto/gen/ts/doc/UserDocFolder'

import {
  addFolder,
  copyDocOrFolderToAnotherFolder,
  moveDocOrFolderToAnotherFolder,
  moveFolderOrDoc,
  updateDoc,
  updateFolder,
} from '@/api/document'
import { goPathPage, goPdfPage } from '@/common/src/utils/url'
import { defineStore } from 'pinia'
import { useLibraryList } from './list'
import {
  LiteratureNode,
  LiteratureNodeFamily,
  allKey,
  newFolderKey,
  noteDirty,
  removePrefix,
  seperatorKey,
  useLibraryIndex,
} from '.'
import { computed, nextTick, ref } from 'vue'

export enum DropToPosition {
  before = -1,
  middle = 0,
  after = 1,
}

export const FOLDER_INDENT = 215

export const useLibraryMenu = defineStore('libraryMenu', () => {
  const storeLibraryIndex = useLibraryIndex()
  const storeLibraryList = useLibraryList()
  const {
    refreshClassiyAuthorVenuePaperList,
    getFilesByFolderId,
    paperListClassifyRefresh,
  } = storeLibraryList

  const loadTree = ref(true)
  const isMoving = ref(false)
  const isCopying = ref(false)
  const selectedNodeFamily = computed(() => {
    if (
      !storeLibraryIndex.selectedKey ||
      storeLibraryIndex.selectedKey === allKey ||
      storeLibraryIndex.selectedKey === seperatorKey
    ) {
      return null
    }

    return getLiteratureNodeFamily(storeLibraryIndex.selectedKey)
  })
  const selectedPaper = computed(() => {
    return storeLibraryList.paperListAll.find(
      (x) => x.docId === selectedNodeFamily.value?.node.docId
    )
  })
  const expandedKeys = ref<LiteratureNode['key'][]>([])
  const dragKey = ref<LiteratureNode['key']>('')

  return {
    loadTree,
    selectedNodeFamily,
    selectedPaper,
    expandedKeys,
    dragKey,
    isMoving,
    isCopying,
    clickItem,
    openMenu,
    getLiteratureNodeFamily,
    moveTo,
    copyTo,
    mouseDown,
    dropTo,
    renameStart,
    renameCancel,
    renameSubmit,
    addSiblingFolderStart,
    addChildFolderStart,
    addFolderCancel,
    addFolderSubmit,
  }

  /** 根据key返回节点自身、节点兄弟数组、节点索引、节点父节点 */
  function getLiteratureNodeFamily(
    key: LiteratureNode['key']
  ): LiteratureNodeFamily {
    if (!key || key === allKey || key === seperatorKey) {
      throw new Error('非常规节点key')
    }

    const [first, ...rest] = storeLibraryIndex.libraryIndexExtra[key].path
    const last = rest.pop()

    if (last === undefined) {
      return {
        parent: null,
        children: storeLibraryIndex.libraryIndexList,
        node: storeLibraryIndex.libraryIndexList[first],
        index: first,
      }
    }

    let parent = storeLibraryIndex.libraryIndexList[first]
    for (const p of rest) {
      parent = parent.children[p]
    }

    return {
      parent,
      children: parent.children,
      node: parent.children[last],
      index: last,
    }
  }

  function clickItem(event: any) {
    const key =
      event.currentTarget.children[0].children[0].getAttribute('data-key')
    if (key === seperatorKey) {
      return
    }

    storeLibraryIndex.selectedKey = key
    storeLibraryList.searchInput = ''

    if (key === allKey) {
      return
    }

    const clickFamily = getLiteratureNodeFamily(key)
    console.log({ clickFamily })

    if (!clickFamily.node.isLeaf) {
      return
    }

    if (clickFamily.node.pdfId) {
      goPdfPage({ pdfId: clickFamily.node.pdfId })
    } else {
      goPathPage(`/paper/${clickFamily.node.paperId}`)
    }
  }

  function openMenu(event: any) {
    storeLibraryIndex.selectedKey = event.node.eventKey
  }

  /** 右键菜单-移动到 */
  async function moveTo(
    targetFolderId: string,
    nodeFamily = selectedNodeFamily.value as LiteratureNodeFamily
  ) {
    // 防止重复点击
    if (isMoving.value) {
      return
    }

    isMoving.value = true
    try {
      const req: MoveDocOrFolderToAnotherFolderReq = {
        movedDocItems: [],
        movedFolderIds: [],
        targetFolderId: removePrefix(targetFolderId),
        /** true=将选中的文献从当前文件夹及其子文件夹下迁移出去,false=只将选中的文献从当前文件夹下迁移出去 */
        isHierarchicallyMoveDoc: false,
      }

      if (nodeFamily.node.isLeaf) {
        req.movedDocItems.push({
          docId: removePrefix(nodeFamily.node.key),
          sourceFolderId: removePrefix(nodeFamily.parent?.key ?? allKey),
        })
      } else {
        req.movedFolderIds.push(removePrefix(nodeFamily.node.key))
      }

      noteDirty.is = true
      await moveDocOrFolderToAnotherFolder(req)
      storeLibraryIndex.fetchLibraryIndex()
      refreshClassiyAuthorVenuePaperList()
    } finally {
      isMoving.value = false
    }
  }

  /** 右键菜单-复制到 */
  async function copyTo(
    targetFolderId: string,
    nodeFamily = selectedNodeFamily.value as LiteratureNodeFamily
  ) {
    // 防止重复点击
    if (isCopying.value) {
      return
    }

    isCopying.value = true
    try {
      const req: CopyDocOrFolderToAnotherFolderReq = {
        docIds: [],
        folderIds: [],
        targetFolderId: removePrefix(targetFolderId),
      }

      if (nodeFamily.node.isLeaf) {
        req.docIds.push(removePrefix(nodeFamily.node.key))
      } else {
        req.folderIds.push(removePrefix(nodeFamily.node.key))
      }

      noteDirty.is = true
      await copyDocOrFolderToAnotherFolder(req)
      storeLibraryIndex.fetchLibraryIndex()
      getFilesByFolderId()
      paperListClassifyRefresh()
      storeLibraryList.authorRefreshOptionsList()
      storeLibraryList.venueRefreshOptionsList()
      storeLibraryList.jcrRefreshOptionsList()
    } finally {
      isCopying.value = false
    }
  }

  function mouseDown(key: LiteratureNode['key']) {
    const detectTimeoutId = setTimeout(() => {
      const clearDrag = () => {
        document.removeEventListener('mouseup', clearDrag)
        dragKey.value = ''
      }
      document.addEventListener('mouseup', clearDrag)
      dragKey.value = key

      clearDetect()
    }, 320)

    const clearDetect = () => {
      clearTimeout(detectTimeoutId)
      document.removeEventListener('mouseup', clearDetect)
    }
    document.addEventListener('mouseup', clearDetect)
  }

  function dropTo(
    dropKey: LiteratureNode['key'],
    position: DropToPosition,
    dragKeyValue = dragKey.value
  ) {
    const dragFamily = getLiteratureNodeFamily(dragKeyValue)
    console.log({ dragFamily })

    noteDirty.is = true

    // 拖到全部文献或根列表
    if (!dropKey || dropKey === allKey) {
      return dropBetweenFoldersAndDocuments(
        allKey,
        storeLibraryIndex.libraryIndexList
      )
    }

    const dropFamily = getLiteratureNodeFamily(dropKey)
    console.log({ dropFamily, position })

    const dragPath =
      storeLibraryIndex.libraryIndexExtra[dragFamily.node.key].path
    const dropPath = storeLibraryIndex.libraryIndexExtra[dropKey].path
    if (dragPath.every((p, i) => p === dropPath[i])) {
      console.log('父目录不能放入子目录')
      return null
    }

    if (dropFamily.node.isLeaf) {
      // 文献拖到文献
      if (dragFamily.node.isLeaf) {
        // 放在前方
        if (
          position === DropToPosition.before ||
          position === DropToPosition.middle
        ) {
          return dropAsPreviousSibling()
          // 放在后方
        } else {
          return dropAsNextSibling()
        }
        // 目录拖到文献
        // 目录不能插在文献与文献之间，所以一律放到目录同级的目录后、文献前
      } else {
        return dropBetweenFoldersAndDocuments()
      }
    }

    if (dragFamily.node.isLeaf) {
      // 文献拖到目录
      // 位置正中，放入目录里面
      if (position === DropToPosition.middle) {
        dragFamily.children.splice(dragFamily.index, 1)
        const dropIndex = dropFamily.node.children.filter(
          (item) => !item.isLeaf
        ).length
        dropFamily.node.children.splice(dropIndex, 0, dragFamily.node)
        storeLibraryIndex.libraryIndexList = [
          ...storeLibraryIndex.libraryIndexList,
        ]
        return formatMoveFolderOrDocReq(
          dragFamily,
          dropFamily.node.key,
          dropFamily.node.children
        )
        // 位置偏上或偏下
        // 文献不能插在目录与目录之间，所以一律放到目录同级的目录后、文献前
      } else {
        return dropBetweenFoldersAndDocuments()
      }
    }

    // 目录拖到目录
    // 位置正中，拖动的目录放入目标目录内
    if (position === DropToPosition.middle) {
      dragFamily.children.splice(dragFamily.index, 1)
      dropFamily.node.children.unshift(dragFamily.node)
      storeLibraryIndex.libraryIndexList = [
        ...storeLibraryIndex.libraryIndexList,
      ]
      return formatMoveFolderOrDocReq(
        dragFamily,
        dropFamily.node.key,
        dropFamily.node.children
      )
      // 位置偏上，放到目标目录前面
    } else if (position === DropToPosition.before) {
      return dropAsPreviousSibling()
      // 位置偏下，放到目标目录后面
    } else {
      return dropAsNextSibling()
    }

    function dropAsPreviousSibling() {
      dragFamily.children.splice(dragFamily.index, 1)
      // 拖动节点和放入节点可能原本就在同一个目录，删除可能导致后者索引变化
      // 所以要重新计算放入节点的索引
      const safeDropIndex = dropFamily.children.findIndex(
        (item) => item === dropFamily.node
      )
      dropFamily.children.splice(safeDropIndex, 0, dragFamily.node)
      storeLibraryIndex.libraryIndexList = [
        ...storeLibraryIndex.libraryIndexList,
      ]
      return formatMoveFolderOrDocReq(
        dragFamily,
        dropFamily.parent?.key ?? allKey,
        dropFamily.children
      )
    }

    function dropAsNextSibling() {
      dragFamily.children.splice(dragFamily.index, 1)
      // 拖动节点和放入节点可能原本就在同一个目录，删除可能导致后者索引变化
      // 所以要重新计算放入节点的索引
      const safeDropIndex = dropFamily.children.findIndex(
        (item) => item === dropFamily.node
      )
      dropFamily.children.splice(safeDropIndex + 1, 0, dragFamily.node)
      storeLibraryIndex.libraryIndexList = [
        ...storeLibraryIndex.libraryIndexList,
      ]
      return formatMoveFolderOrDocReq(
        dragFamily,
        dropFamily.parent?.key ?? allKey,
        dropFamily.children
      )
    }

    function dropBetweenFoldersAndDocuments(
      dropFamilyParentKey = dropFamily.parent?.key ?? allKey,
      dropFamilyChildren = dropFamily.children
    ) {
      dragFamily.children.splice(dragFamily.index, 1)
      const dropIndex = dropFamilyChildren.filter((item) => !item.isLeaf).length
      dropFamilyChildren.splice(dropIndex, 0, dragFamily.node)
      storeLibraryIndex.libraryIndexList = [
        ...storeLibraryIndex.libraryIndexList,
      ]
      return formatMoveFolderOrDocReq(
        dragFamily,
        dropFamilyParentKey,
        dropFamilyChildren
      )
    }

    async function formatMoveFolderOrDocReq(
      dragFamily: LiteratureNodeFamily,
      targetFolderId: LiteratureNode['key'],
      children: LiteratureNode[]
    ) {
      const req: MoveFolderOrDocReq = {
        type: dragFamily.node.isLeaf ? 0 : 1,
        sourceFolderId: removePrefix(dragFamily.parent?.key ?? allKey),
        targetFolderId: removePrefix(targetFolderId),
        movedIds: [removePrefix(dragFamily.node.key)],
        targetFolderItems: children
          .filter((node) => node.isLeaf === dragFamily.node.isLeaf)
          .map((node, sort) => ({
            id: removePrefix(node.key),
            sort,
          })),
      }

      await moveFolderOrDoc(req)
      storeLibraryIndex.fetchLibraryIndex()
      refreshClassiyAuthorVenuePaperList()
    }
  }

  async function renameStart() {
    (selectedNodeFamily.value as LiteratureNodeFamily).node.isEditing = true
    storeLibraryIndex.libraryIndexList = [...storeLibraryIndex.libraryIndexList]
    await nextTick()
    document
      .getElementById(`list-tree-input:${storeLibraryIndex.selectedKey}`)
      ?.focus()
  }

  function renameCancel() {
    if (!selectedNodeFamily.value?.node.isEditing) {
      return
    }

    selectedNodeFamily.value.node.isEditing = false
  }

  function renameSubmit(event: any) {
    if (!selectedNodeFamily.value?.node.isEditing) {
      return
    }

    console.log(event.target.value)
    const value = (event.target.value as string).trim()
    const family = selectedNodeFamily.value as LiteratureNodeFamily
    const id = removePrefix(family.node.key)
    ;(async () => {
      if (family.node.isLeaf) {
        const result = await updateDoc({
          docId: id,
          docName: value,
        })

        if (!result) {
          return
        }

        storeLibraryIndex.fetchLibraryIndex()
        getFilesByFolderId()
      } else {
        await updateFolder({
          folderId: id,
          name: value,
        })

        storeLibraryIndex.fetchLibraryIndex()
      }

      family.node.title = value
      family.node.isEditing = false
      storeLibraryIndex.libraryIndexList = [
        ...storeLibraryIndex.libraryIndexList,
      ]
    })()
    noteDirty.is = true
  }

  function getNewFolder(): LiteratureNode {
    return {
      isLeaf: false,
      children: [],
      key: newFolderKey,
      title: '',
      isEditing: true,
      hoverType: 0,
      docCount: 0,
    }
  }

  async function addSiblingFolderStart() {
    const newFolder = getNewFolder()
    ;(selectedNodeFamily.value as LiteratureNodeFamily).children.splice(
      (selectedNodeFamily.value as LiteratureNodeFamily).index + 1,
      0,
      newFolder
    )
    storeLibraryIndex.libraryIndexList = [...storeLibraryIndex.libraryIndexList]
    await nextTick()
    document.getElementById(`list-tree-input:${newFolder.key}`)?.focus()
  }

  async function addChildFolderStart() {
    const newFolder = getNewFolder()
    if (selectedNodeFamily.value) {
      selectedNodeFamily.value.node.children.unshift(newFolder)
      expandedKeys.value.push(storeLibraryIndex.selectedKey)
      storeLibraryIndex.libraryIndexList = [
        ...storeLibraryIndex.libraryIndexList,
      ]
    } else if (storeLibraryIndex.selectedKey === allKey) {
      storeLibraryIndex.libraryIndexList.splice(1, 0, newFolder)
      await nextTick()
      expandedKeys.value.push(storeLibraryIndex.selectedKey)
    }

    await nextTick()
    document.getElementById(`list-tree-input:${newFolderKey}`)?.focus()
  }

  function addFolderCancel() {
    const family = getLiteratureNodeFamily(newFolderKey)
    if (!family.node.isEditing) {
      return
    }

    storeLibraryIndex.selectedKey = ''
    family.children.splice(family.index, 1)
    storeLibraryIndex.libraryIndexList = [...storeLibraryIndex.libraryIndexList]
  }

  async function addFolderSubmit(event: any) {
    const family = getLiteratureNodeFamily(newFolderKey)
    if (!family.node.isEditing) {
      return
    }

    const value = (event.target.value as string).trim()
    if (!value) {
      return addFolderCancel()
    }

    family.node.isEditing = false
    family.node.title = value
    const req: Parameters<typeof addFolder>[0] = {
      name: value,
      parentId: removePrefix(family.parent?.key ?? allKey),
      level: storeLibraryIndex.libraryIndexExtra[newFolderKey].path.length - 1,
      sort: family.index,
      oldFolderItems: family.children
        .filter((child) => !child.isLeaf)
        .filter((child) => child.key !== newFolderKey)
        .map((child, index) => ({
          id: removePrefix(child.key),
          sort: index,
        })),
    }

    storeLibraryIndex.libraryIndexList = [...storeLibraryIndex.libraryIndexList]
    // console.log({ req })
    // return

    const result: any = await addFolder(req)
    getFilesByFolderId()
    await storeLibraryIndex.fetchLibraryIndex()
    const key = Object.keys(storeLibraryIndex.libraryIndexExtra).find((k) =>
      k.includes(result.data.data.id)
    )
    storeLibraryIndex.selectedKey = key || ''
    noteDirty.is = true
  }
})
