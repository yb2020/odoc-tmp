/* eslint-disable no-use-before-define */
import { computed, ref } from 'vue'
import { DocDetailInfo } from 'go-sea-proto/gen/ts/doc/ClientDoc'
import { useLocalStorage } from '@vueuse/core'
import {
  UserDocInfo,
  UserDocListSortType,
} from 'go-sea-proto/gen/ts/doc/UserDocManage'

import { newFolderKey, useLibraryIndex } from '.'
import {
  getAuthorOptionList,
  getClassifyOptionList,
  getDocListNew,
  getJcrOptionList,
  getVenueOptionList,
} from '@/api/document'

import { defineStore } from 'pinia'
import { useClassify } from '../classify'
import { liteThrottle, withDeprecate } from '@/common/src/utils/aiknowledge-special-util'
import { UI, USER_CENTER } from '@/common/src/constants/storage-keys'

const autoDeprecateGetDocListByFolderId = withDeprecate(getDocListNew)

export const useLibraryList = defineStore('libraryList', () => {
  const storeLibraryIndex = useLibraryIndex()
  const storeClassify = useClassify()

  const currentSortType = useLocalStorage<UserDocListSortType>(
    USER_CENTER.LITERATURE_SORT,
    UserDocListSortType.LAST_ADD
  )
  const dropdownSortType = ref(
    currentSortType.value === UserDocListSortType.LAST_ADD ||
      currentSortType.value === UserDocListSortType.LAST_READ
      ? currentSortType.value
      : UserDocListSortType.LAST_READ
  )

  const searchInput = ref('')
  const paperListAll = ref<Required<UserDocInfo>[]>([])
  const paperListPageNumber = ref(1)
  const paperListPageSize = useLocalStorage(UI.PAPER_LIST_PAGE_SIZE, 10)
  const paperListTotal = useLocalStorage(UI.PAPER_LIST_TOTAL, 10)
  const paperListLoading = ref(true)
  const paperListEmpty = computed(() => {
    return (
      !paperListTotal.value &&
      !searchInput.value &&
      !paperListClassifyParam.value.length &&
      !paperListAuthorFilter.param.value.length &&
      !paperListVenueFilter.param.value.length &&
      !paperListJcrFilter.param.value.length &&
      !paperListImpactFactorNoLimit.value &&
      paperListImpactFactorMin.value === undefined &&
      paperListImpactFactorMax.value === undefined
    )
  })

  const refreshClassiyAuthorVenuePaperList = async (goToFirstPage = false) => {
    // 1. getDocIndex - 通过 fetchLibraryIndex 调用
    const docIndex = storeLibraryIndex.fetchLibraryIndex()
    
    // 2. getDocRelatedClassifyList
    const classify = paperListClassifyRefresh()
    
    // 3. getDocRelatedAuthorList
    const author = paperListAuthorFilter.refreshOptionsList()
    
    // 4. getDocRelatedVenueList
    const venue = paperListVenueFilter.refreshOptionsList()
    
    // 5. partions?folderId=0 - 这个接口可能在其他地方调用
    // 确保这个接口被调用
    const partions = getJcrOptionList({ folderId: '0' })
    
    // 6. getUserAllClassifyList
    const jcr = paperListJcrFilter.refreshOptionsList()
    
    // 等待所有接口完成
    await Promise.all([docIndex, classify, author, venue, partions, jcr])
    
    // 7. getDocList - 通过 getFilesByFolderId 调用
    if (goToFirstPage) {
      await getFilesByFolderId(1)
    } else {
      await getFilesByFolderId()
    }
  }

  const getFilesByFolderId = async (
    number = paperListPageNumber.value,
    size = paperListPageSize.value
  ) => {
    await storeClassify.initClassifyList()

    paperListLoading.value = true

    if (storeLibraryIndex.rawFolderId === newFolderKey) {
      paperListAll.value = []
      paperListLoading.value = false
      paperListPageNumber.value = 1
      paperListTotal.value = 0
      paperListPageSize.value = size
      paperListClassifyChecked.value = []
      paperListChecked.value = []
      return
    }

    let impactOfFactorRange: number[] = [0, 99999]
    let onlyShowDocsWithImpactOfFactor = true

    const isNumber = (value: unknown) => {
      return typeof value === 'number' && !isNaN(value)
    }

    if (
      isNumber(paperListImpactFactorMin.value) ||
      isNumber(paperListImpactFactorMax.value)
    ) {
      if (isNumber(paperListImpactFactorMin.value)) {
        impactOfFactorRange[0] = paperListImpactFactorMin.value as number
      }
      if (isNumber(paperListImpactFactorMax.value)) {
        impactOfFactorRange[1] = paperListImpactFactorMax.value as number
      }
    } else if (!paperListImpactFactorNoLimit.value) {
      impactOfFactorRange = []
      onlyShowDocsWithImpactOfFactor = false
    }

    const params = {
      currentPage: number,
      pageSize: size,
      folderId: storeLibraryIndex.rawFolderId,
      /** 排序方式,0=最近添加,1=最近阅读 */
      sortType: currentSortType.value,
      ascSort: paperListLocalSortDirection.value === 1,
      classifyIds: paperListClassifyParam.value,
      authorInfos: paperListAuthorFilter.param.value,
      venueInfos: paperListVenueFilter.param.value,
      ...(searchInput.value && {
        searchContent: searchInput.value,
      }),
      jcrPartions: paperListJcrFilter.param.value,
      impactOfFactorRange,
      onlyShowDocsWithImpactOfFactor,
    };
    try {
      const result = await autoDeprecateGetDocListByFolderId(params)      
      // 直接使用 result，因为 autoDeprecateGetDocListByFolderId 返回的是 { docList, total }
      paperListPageNumber.value = number
      paperListPageSize.value = size
      paperListTotal.value = result.total
      paperListAll.value = result.docList as Required<UserDocInfo>[]
      paperListLoading.value = false
      paperListClassifyChecked.value = paperListClassifyChecked.value.filter(
        (id) => paperListClassifyList.value.some((classify) => classify.id === id)
      )
      paperListChecked.value = []
    } catch (error) {
      console.error('[ERROR] getFilesByFolderId - 获取文档列表失败:', error);
      paperListLoading.value = false
    }
  }

  const paperListClassifyChecked = ref<string[]>([])
  const paperListClassifyEmptyChecked = ref(false)
  const paperListClassifyList = ref<
    {
      id: string
      name: string
    }[]
  >([])
  const paperListClassifyRefresh = async () => {
    const list = await getClassifyOptionList({
      folderId: storeLibraryIndex.rawFolderId,
    })
    paperListClassifyList.value = list.map((item) => ({
      id: item.classifyId,
      name: item.classifyName,
    }))
    paperListClassifyChecked.value = paperListClassifyChecked.value.filter(
      (item) => {
        return paperListClassifyList.value.some((option) => option.id === item)
      }
    )
  }

  const paperListClassifyAllChecked = computed(() => {
    const normalAll = paperListClassifyList.value.every((classify) =>
      paperListClassifyChecked.value.includes(classify.id)
    )
    return searchInput.value
      ? normalAll
      : normalAll && paperListClassifyEmptyChecked.value
  })
  const paperListClassifyIndeterminate = computed(() => {
    return (
      !paperListClassifyAllChecked.value &&
      (searchInput.value
        ? paperListClassifyChecked.value.length > 0
        : paperListClassifyChecked.value.length > 0 ||
          paperListClassifyEmptyChecked.value)
    )
  })
  const paperListClassifyParam = computed(() => {
    return !paperListClassifyEmptyChecked.value
      ? paperListClassifyChecked.value
      : paperListClassifyAllChecked.value
        ? []
        : [...paperListClassifyChecked.value, '0']
  })

  const paperListLocalSort = ref<Extract<
    TableHeadKey,
    'docName' | 'publishDate'
  > | null>(null)
  const paperListLocalSortDirection = useLocalStorage<1 | -1>(
    UI.PAPER_LIST_SORT_DIRECTION,
    1
  )

  const paperListClassifyToggle = (id: string) => {
    if (paperListClassifyChecked.value.includes(id)) {
      paperListClassifyChecked.value = paperListClassifyChecked.value.filter(
        (cid) => cid !== id
      )
    } else {
      paperListClassifyChecked.value.push(id)
    }

    getFilesByFolderId(1)
  }
  const paperListClassifyAllToggle = () => {
    if (paperListClassifyAllChecked.value) {
      paperListClassifyChecked.value = []
      paperListClassifyEmptyChecked.value = false
    } else {
      paperListClassifyChecked.value = paperListClassifyList.value.map(
        (classify) => classify.id
      )
      paperListClassifyEmptyChecked.value = true
    }

    getFilesByFolderId(1)
  }
  const paperListClassifyEmptyToggle = () => {
    paperListClassifyEmptyChecked.value = !paperListClassifyEmptyChecked.value
    getFilesByFolderId(1)
  }

  const paperListAuthorFilter = useFilter(() => {
    return getAuthorOptionList({
      folderId: storeLibraryIndex.rawFolderId,
    })
  })
  const {
    optionList: authorOptionList,
    checked: authorChecked,
    emptyChecked: authorEmptyChecked,
    allChecked: authorAllChecked,
    indeterminate: authorIndeterminate,
    toggle: authorToggle,
    allToggle: authorAllToggle,
    emptyToggle: authorEmptyToggle,
    param: authorParam,
    refreshOptionsList: authorRefreshOptionsList,
  } = paperListAuthorFilter
  const paperListVenueFilter = useFilter(() => {
    return getVenueOptionList({
      folderId: storeLibraryIndex.rawFolderId,
    })
  })
  const {
    optionList: venueOptionList,
    checked: venueChecked,
    emptyChecked: venueEmptyChecked,
    allChecked: venueAllChecked,
    indeterminate: venueIndeterminate,
    toggle: venueToggle,
    allToggle: venueAllToggle,
    emptyToggle: venueEmptyToggle,
    param: venueParam,
    refreshOptionsList: venueRefreshOptionsList,
  } = paperListVenueFilter
  const paperListJcrFilter = useFilter(() => {
    return getJcrOptionList({
      folderId: storeLibraryIndex.rawFolderId,
    })
  })
  const {
    optionList: jcrOptionList,
    checked: jcrChecked,
    emptyChecked: jcrEmptyChecked,
    allChecked: jcrAllChecked,
    indeterminate: jcrIndeterminate,
    toggle: jcrToggle,
    allToggle: jcrAllToggle,
    emptyToggle: jcrEmptyToggle,
    param: jcrParam,
    refreshOptionsList: jcrRefreshOptionsList,
  } = paperListJcrFilter

  const paperHeadList = ref(getTableHeadList())
  const paperHeadSync = liteThrottle(
    () => {
      setTableHeadList(paperHeadList.value)
    },
    2000,
    true
  )

  const paperHeadVisibleList = computed(() => {
    return paperHeadList.value.filter((item) => item.visible && item.key !== 'operation')
  })

  const paperHeadExtra = computed(() => {
    const extra: Partial<PaperHeadExtra> = {}
    paperHeadList.value.forEach((head) => {
      extra[head.key] = {
        style: {
          flexBasis: `${head.width}px`,
          flexShrink: '0',
          flexGrow: '0',
          // ...(head.key === 'docName' && {
          //   flexGrow: '1',
          //   overflow: 'hidden',
          // }),
        },
      }
    })

    return extra as PaperHeadExtra
  })

  const paperListChecked = ref<UserDocInfo['docId'][]>([])
  const paperListCheckedMap = computed(() => {
    const map: Record<UserDocInfo['docId'], boolean> = {}
    paperListAll.value.forEach((item) => {
      map[item.docId] = false
    })
    paperListChecked.value.forEach((id) => {
      if (id in map) {
        map[id] = true
      }
    })
    return map
  })
  const paperListAuthorEdit = ref('')
  const paperListPublishEdit = ref('')
  const paperListPublishYear = ref<number | undefined>(new Date().getFullYear())
  const paperListPublishMonth = ref<number | undefined>(undefined)
  const paperListPublishDate = ref<number | undefined>(undefined)
  const paperListVenueEdit = ref('')
  const paperListVenueEditContent = ref('')
  const paperListJcrEdit = ref('')
  const paperListJcrEditContent = ref('')
  const paperListImpactFactorEdit = ref('')
  const paperListImpactFactorEditContent = ref('')

  const paperListImpactFactorNoLimit = ref(false)
  const paperListImpactFactorMin = ref<number | undefined>(undefined)
  const paperListImpactFactorMax = ref<number | undefined>(undefined)

  return {
    searchInput,
    paperListAll,
    paperListPageNumber,
    paperListPageSize,
    paperListTotal,
    paperListEmpty,
    paperListLoading,
    refreshClassiyAuthorVenuePaperList,
    getFilesByFolderId,
    paperListClassifyChecked,
    paperListClassifyList,
    paperListClassifyEmptyChecked,
    paperListClassifyAllChecked,
    paperListClassifyIndeterminate,
    paperListClassifyToggle,
    paperListClassifyAllToggle,
    paperListClassifyEmptyToggle,
    paperListClassifyParam,
    paperListClassifyRefresh,

    authorOptionList,
    authorChecked,
    authorEmptyChecked,
    authorAllChecked,
    authorIndeterminate,
    authorToggle,
    authorAllToggle,
    authorEmptyToggle,
    authorParam,
    authorRefreshOptionsList,
    venueOptionList,
    venueChecked,
    venueEmptyChecked,
    venueAllChecked,
    venueIndeterminate,
    venueToggle,
    venueAllToggle,
    venueEmptyToggle,
    venueParam,
    venueRefreshOptionsList,
    jcrOptionList,
    jcrChecked,
    jcrEmptyChecked,
    jcrAllChecked,
    jcrIndeterminate,
    jcrToggle,
    jcrAllToggle,
    jcrEmptyToggle,
    jcrParam,
    jcrRefreshOptionsList,

    paperHeadList,
    paperHeadVisibleList,
    paperHeadExtra,
    paperHeadSync,
    paperListLocalSort,
    paperListLocalSortDirection,
    paperListChecked,
    paperListCheckedMap,
    paperListAuthorEdit,
    paperListPublishEdit,
    paperListPublishYear,
    paperListPublishMonth,
    paperListPublishDate,
    paperListVenueEdit,
    paperListVenueEditContent,
    paperListJcrEdit,
    paperListJcrEditContent,
    paperListImpactFactorEdit,
    paperListImpactFactorEditContent,
    currentSortType,
    dropdownSortType,
    paperListImpactFactorNoLimit,
    paperListImpactFactorMin,
    paperListImpactFactorMax,
  }

  function useFilter(getOptionList: () => Promise<string[]>) {
    const optionList = ref<string[]>([])
    const checked = ref<string[]>([])
    const emptyChecked = ref(false)
    const allChecked = computed(() => {
      const normalAll = optionList.value.every((author) =>
        checked.value.includes(author)
      )
      return normalAll && emptyChecked.value
      // return searchInput.value ? normalAll : normalAll && emptyChecked.value
    })
    const indeterminate = computed(() => {
      return (
        !allChecked.value && (checked.value.length > 0 || emptyChecked.value)
      )
    })
    const toggle = (value: string) => {
      if (checked.value.includes(value)) {
        checked.value = checked.value.filter((item) => item !== value)
      } else {
        checked.value.push(value)
      }

      getFilesByFolderId(1)
    }
    const allToggle = () => {
      if (allChecked.value) {
        checked.value = []
        emptyChecked.value = false
      } else {
        checked.value = optionList.value.slice()
        emptyChecked.value = true
      }

      getFilesByFolderId(1)
    }
    const emptyToggle = () => {
      emptyChecked.value = !emptyChecked.value
      getFilesByFolderId(1)
    }
    const param = computed(() => {
      return !emptyChecked.value
        ? checked.value
        : allChecked.value
          ? []
          : [...checked.value, '']
    })
    const refreshOptionsList = async () => {
      optionList.value = await getOptionList()
      checked.value = checked.value.filter((item) =>
        optionList.value.includes(item)
      )
    }
    return {
      optionList,
      checked,
      emptyChecked,
      allChecked,
      indeterminate,
      toggle,
      allToggle,
      emptyToggle,
      param,
      refreshOptionsList,
    }
  }
})

export type TableHeadKey =
  | keyof Pick<DocDetailInfo, 'authors' | 'publishDate'>
  | keyof Pick<
      UserDocInfo,
      | 'docName'
      | 'classifyInfos'
      | 'remark'
      | 'parseProgress'
      | 'displayVenue'
      | 'jcrVenuePartion'
      | 'impactOfFactor'
      | 'importantanceScore'
    >
  // | 'operation'

export const sortToKey: Record<UserDocListSortType, TableHeadKey | ''> = {
  [UserDocListSortType.UNRECOGNIZED]: '',
  [UserDocListSortType.DEFAULT]: '',
  [UserDocListSortType.LAST_ADD]: '',
  [UserDocListSortType.LAST_READ]: '',
  [UserDocListSortType.CUSTOM_SORT]: '',
  [UserDocListSortType.DOC_NAME]: 'docName',
  [UserDocListSortType.PUBLISH_DATE]: 'publishDate',
  [UserDocListSortType.IMPACT_OF_FACTOR]: 'impactOfFactor',
  [UserDocListSortType.IMPORTANCE_SCORE]: 'importantanceScore',
}

export interface TableHeadDetail {
  key: TableHeadKey
  visible: boolean
  width: number
}

export type PaperHeadExtra = Record<
  TableHeadKey,
  {
    style: Partial<CSSStyleDeclaration>
  }
>

export const TableHeadName: Record<TableHeadKey, string> = {
  docName: 'home.library.tableHead.docName',
  classifyInfos: 'home.library.tableHead.classifyInfos',
  authors: 'home.library.tableHead.authors',
  remark: 'home.library.tableHead.remark',
  parseProgress: 'home.library.tableHead.parseProgress',
  publishDate: 'home.library.tableHead.publishDate',
  displayVenue: 'home.library.tableHead.displayVenue',
  jcrVenuePartion: 'home.library.tableHead.jcrVenuePartion',
  impactOfFactor: 'home.library.tableHead.impactOfFactor',
  importantanceScore: 'home.library.tableHead.importantanceScore',
  // operation: 'home.library.tableHead.operation',
}

export const defaultTableHeadWidth: Record<TableHeadKey, number> = {
  docName: 200,
  classifyInfos: 116,
  authors: 116,
  remark: 100,
  parseProgress: 100,
  publishDate: 150,
  displayVenue: 146,
  jcrVenuePartion: 100,
  impactOfFactor: 130,
  importantanceScore: 200,
  // operation: 100,
}

const orderList: TableHeadKey[] = [
  'docName',
  'classifyInfos',
  'authors',
  'remark',
  'parseProgress',
  'publishDate',
  'displayVenue',
  'jcrVenuePartion',
  'impactOfFactor',
  'importantanceScore',
  // 'operation',
]

const PAPER_HEAD_LIST = 'paperHeadList'

const getDefaultTableHeadList = () =>
  orderList.map<TableHeadDetail>((key) => ({
    key,
    visible: true,
    width: defaultTableHeadWidth[key],
  }))

const getTableHeadList = (): TableHeadDetail[] => {
  const defaultTableHeadList = getDefaultTableHeadList()
  const listString = localStorage.getItem(PAPER_HEAD_LIST)

  if (!listString) {
    return defaultTableHeadList
  }

  let paperHeadList: TableHeadDetail[]
  try {
    paperHeadList = JSON.parse(listString)
  } catch (error) {
    return defaultTableHeadList
  }

  if (!Array.isArray(paperHeadList)) {
    return defaultTableHeadList
  }

  const newHeadList = defaultTableHeadList.filter((item) =>
    paperHeadList.every((head) => head.key !== item.key)
  )

  if (newHeadList.length) {
    paperHeadList.push(...newHeadList)
    setTableHeadList(paperHeadList)
  }

  return paperHeadList
}

const setTableHeadList = (list: TableHeadDetail[]) => {
  localStorage.setItem(PAPER_HEAD_LIST, JSON.stringify(list))
}
