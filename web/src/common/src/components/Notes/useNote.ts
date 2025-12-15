import { ComputedRef, Ref, computed, ref, watch } from 'vue';
import { SortTypeEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/protocol/note/request/_GetMyNoteMarkListReq';
// import { GetAnnotateTagsResp } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Tag';
import { GetAnnotateTagsResponse} from 'go-sea-proto/gen/ts/note/PaperNoteAnnotateTag';
// import {
  // GetWordListByFolderIdResponse,
  // WordInfo,
  // NoteManageFolderInfo,
  // NoteManageDocInfo,
  // GetSummaryResponse,
  // GetWordsResponse,
  // GetSummaryListByFolderIdResponse,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import {GetWordListByFolderIdResponse, NoteManageFolderInfo, NoteManageDocInfo, GetSummaryListByFolderIdResponse } from 'go-sea-proto/gen/ts/note/NoteManage';
import {
  getAnnotateTags,
  getMarkTagListByFolderId,
  getMyNoteMarkList,
  getSummaryNote,
  getWordListByFolderId,
  getSummaryListByFolderId,
  getWordNotes,
  NoteAnnotation,
  MyNoteMarkListResponse,
} from '~/src/api/note';
import { GetNoteSummaryByNoteIdResponse } from 'go-sea-proto/gen/ts/note/NoteSummary'
import { GetNoteWordsByNoteIdResponse, WordInfo } from 'go-sea-proto/gen/ts/note/NoteWord'
import { goSummaryPage } from '@common/utils/url';
import {
  ColorKey,
  styleMap,
  NoteSubTypes,
  NoteStyleDefine,
  NoteBreadcrumb,
  NoteFolder,
  NoteFolderExtra,
  NoteTag,
} from './types';
import { useLocalStorage } from '@vueuse/core';

export const FOLDER_INDENT = 215;
export const NOTE_TAG_ALL_KEY = 'NOTE_TAG_ALL_KEY';
export const NOTE_FOLDER_ALL_KEY = 'NOTE_FOLDER_ALL_KEY';
export const NOTE_VOCABULARY_SIZE = 'NOTE_VOCABULARY_SIZE';

export const noteStyleList: NoteStyleDefine[] = Object.keys(styleMap).map(
  (key) => ({
    ...styleMap[key as unknown as ColorKey],
    type: key as unknown as ColorKey,
  })
);

export const noteStyleMap = Object.freeze(
  noteStyleList.reduce(
    (o, sd) => {
      o[sd.type] = sd;
      return o;
    },
    {} as Record<ColorKey, NoteStyleDefine>
  )
);

export const removePrefix = (key: string) => {
  return key.split('-').pop()!;
};

export const autoDeprecate = <T extends (...args: any[]) => Promise<any>>(
  func: T
) => {
  let latest: symbol | null = null;
  return async function (this: any, ...args: any[]) {
    const current = Symbol('autoDeprecateCurrent');
    latest = current;
    const result = await func.apply(this, args);
    if (latest === current) {
      return result;
    }

    return new Promise(() => null);
  };
};

export const convertNoteFolder = (
  folder: NoteManageFolderInfo,
  prefix = ''
): NoteFolder => {
  const key = `${prefix}-${folder.folderId}`;
  if (!folder.childrenFolders) {
    folder.childrenFolders = [];
  }

  const children = folder.childrenFolders
    .map((child) => convertNoteFolder(child, key))
    .concat(
      folder?.docInfos?.map((doc) => ({
        key: doc.noteId,
        title: doc.docName,
        noteId: doc.noteId,
        docId: doc.docId,
        isDoc: true,
        children: [],
        noteWordCount: doc.noteWordCount,
        noteAnnotateCount: doc.noteAnnotateCount,
      })) || []
    );

  return {
    key,
    title: folder.name,
    count: folder.count,
    children,
    noteWordCount: folder.noteWordCount,
    noteAnnotateCount: folder.noteAnnotateCount,
  };
};

export const useNote = (
  getFolderTree: ComputedRef<() => Promise<any>>,
  title: string,
  noteType: Ref<NoteSubTypes>,
  isTabbed = false
) => {
  const noteAllFolder = convertNoteFolder({
    folderId: '0',
    name: title,
    childrenFolders: [],
    count: 0,
    docInfos: [],
    noteWordCount: 0,
    noteAnnotateCount: 0,
  });
  const noteFolderLoading = ref(false);
  const noteFolderList = ref<NoteFolder[]>([noteAllFolder]);
  const noteFolderMap = computed(() => {
    const extraMap: Record<NoteFolder['key'], NoteFolderExtra> = {};
    const path: number[] = [];
    // 回溯算法 DFS 后序遍历
    const backtrack = (folder: NoteFolder, index: number): void => {
      path.push(index);

      folder.children?.forEach(backtrack);

      const width = FOLDER_INDENT - (path.length - 1) * 19;
      extraMap[folder.key] = {
        path: [...path],
        width,
        title: folder.title,
        docInfos: folder.docInfos,
        docId: folder.docId,
        isDoc: !!folder.isDoc,
      };

      path.pop();
    };

    noteFolderList.value?.forEach(backtrack);

    return extraMap;
  });
  const noteFolderSelected = ref<NoteFolder['key']>(noteAllFolder.key);
  const noteFolderSelectedCache = ref<Record<NoteSubTypes, string>>({
    [NoteSubTypes.Summary]: '0',
    [NoteSubTypes.Vocabulary]: '0',
    [NoteSubTypes.Annotation]: '0',
  });
  watch(
    noteFolderSelected,
    () =>
      (noteFolderSelectedCache.value[noteType.value] = noteFolderSelected.value)
  );
  const noteFolderExpanded = ref<NoteFolder['title'][]>([]);
  const emptyFolderExpanded = ref<NoteFolder['key'][]>([]);
  const noteBreadcrumbList = computed<NoteBreadcrumb[]>(() => {
    const list = [
      {
        key: noteAllFolder.key,
        title: noteAllFolder.title,
      },
    ];

    if (
      !noteFolderSelected.value ||
      noteFolderSelected.value === noteAllFolder.key ||
      !(noteFolderSelected.value in noteFolderMap.value)
    ) {
      return list;
    }

    let children = noteFolderList.value;
    const { path } = noteFolderMap.value[noteFolderSelected.value];
    for (const p of path) {
      const folder = children[p];
      list.push({
        key: folder.key,
        title: folder.title,
      });
      ({ children } = folder);
    }

    return list;
  });

  const userTagList = ref<GetAnnotateTagsResponse['tags']>();
  const noteTagList = ref<NoteTag[]>([]);
  const noteTagSelected = ref<NoteTag['tagId'][]>([]);
  const noteTagEditId = ref<NoteTag['tagId']>('');
  const noteTagEditName = ref<NoteTag['tagName']>('');

  const noteSort = ref<SortTypeEnum>(SortTypeEnum.NEWEST);
  const noteStyle = ref<Record<ColorKey, boolean>>({
    [ColorKey.blue]: true,
    [ColorKey.green]: true,
    [ColorKey.yellow]: true,
    [ColorKey.orange]: true,
    [ColorKey.red]: true,
  });
  const noteQuote = ref(true);

  const summaryPageSize = ref(10);
  const summaryPageNumber = ref(1);
  const summaryTotal = ref(0);
  const summaryError = ref('');
  const summaryList = ref<NoteManageDocInfo[]>([]);
  const summaryNote = ref<GetNoteSummaryByNoteIdResponse>();
  const wordList = ref<WordInfo[]>([]);
  const wordListError = ref('');
  const wordTotal = ref(0);
  const wordPageNumber = ref(1);
  const wordPageSize = useLocalStorage(NOTE_VOCABULARY_SIZE, 20);
  const wordLoading = ref(true);
  const noteLoading = ref(true);
  const noteList = ref<NoteAnnotation[]>([]);
  const noteListError = ref('');
  const noteTotal = ref(0);
  const notePageNumber = ref(1);
  const notePageSize = ref(10);
  const noteTotalPage = computed(() =>
    Math.ceil(noteTotal.value / notePageSize.value)
  );
  const noteSearchKeyword = ref('');

  const noteEditId = ref('');
  const noteEditContent = ref('');
  const autoDeprecateGetMyNoteMarkList = autoDeprecate(getMyNoteMarkList);

  const noteState = {
    getFolderTree,
    noteAllFolder,
    noteFolderLoading,
    noteFolderList,
    noteFolderMap,
    noteFolderSelected,
    noteFolderSelectedCache,
    noteFolderExpanded,
    emptyFolderExpanded,
    noteBreadcrumbList,
    clickNoteFolder,
    noteTagList,
    userTagList,
    noteTagSelected,
    noteTagEditId,
    noteTagEditName,
    fetchUserTagList,
    fetchNoteTagList,
    fetchWordList,
    summaryError,
    summaryList,
    summaryNote,
    summaryPageNumber,
    summaryPageSize,
    summaryTotal,
    fetchSummary,
    wordList,
    wordListError,
    wordTotal,
    wordPageNumber,
    wordPageSize,
    wordLoading,
    noteSort,
    noteStyle,
    noteQuote,
    noteList,
    noteListError,
    noteTotal,
    noteTotalPage,
    notePageNumber,
    notePageSize,
    noteLoading,
    noteSearchKeyword,
    noteEditId,
    noteEditContent,
    fetchNoteList,
    refreshNoteExplorer,
  };
  // if (process.client) {
  //   console.warn(noteState);
  // }
  return noteState;

  async function refreshNoteExplorer() {
    if (!getFolderTree.value) {
      return;
    }

    noteFolderLoading.value = true;
    try {
      const response = await getFolderTree.value();
      const children = (response?.folderInfos || response?.folderList)
        .map((folder: NoteManageFolderInfo) => convertNoteFolder(folder))
        .concat(
          (response?.unclassifiedDocInfos || [])?.map(
            (doc: NoteManageDocInfo) => ({
              key: doc.noteId,
              title: doc.docName,
              docId: doc.docId,
              isDoc: true,
              children: [],
            })
          ) || []
        );
      noteAllFolder.count = response?.totalDocCount || response?.total;
      noteAllFolder.noteWordCount = response?.totalDocCount || response?.total;
      noteAllFolder.noteAnnotateCount = response?.totalDocCount || response?.total;

      noteFolderList.value = [noteAllFolder, ...children];
    } catch (e) {}
    noteFolderLoading.value = false;

    const equalFolder = (key1: NoteFolder['key'], key2: NoteFolder['key']) => {
      const id1 = removePrefix(key1);
      const id2 = removePrefix(key2);
      return id1 === id2;
    };

    const newKeys = Object.keys(noteFolderMap.value);
    const filterKeyList = (keyList: NoteFolder['key'][]) => {
      return keyList
        .map((oldKey) => {
          const newKey = newKeys.find((key) => equalFolder(key, oldKey));
          return newKey!;
        })
        .filter(Boolean);
    };
    noteFolderExpanded.value = filterKeyList(noteFolderExpanded.value);
    emptyFolderExpanded.value = filterKeyList(emptyFolderExpanded.value);

    const cachedSelectedKey = noteFolderSelectedCache.value[noteType.value];
    const newSelectedKey =
      newKeys.find((key) => equalFolder(key, cachedSelectedKey)) ??
      noteAllFolder.key;

    clickNoteFolder([newSelectedKey]);
  }

  async function clickNoteFolder([key]: NoteFolder['key'][]) {
    switch (noteType.value) {
      case NoteSubTypes.Summary: {
        if (isTabbed || !noteFolderMap.value[key]?.isDoc) {
          noteFolderSelected.value = key;
          // 文献才需要获取总结 文件夹则不需要
          fetchSummary(key ? 1 : undefined);
        } else {
          goSummaryPage({
            noteId: removePrefix(key),
          });
        }
        break;
      }
      case NoteSubTypes.Vocabulary: {
        noteFolderSelected.value = key;
        fetchWordList(key ? 1 : undefined);
        break;
      }
      case NoteSubTypes.Annotation: {
        noteFolderSelected.value = key;
        fetchNoteTagList();
        if (key) {
          noteSearchKeyword.value = '';
          noteTagEditId.value = '';
          noteEditId.value = '';
          noteTagEditName.value = '';
          noteEditContent.value = '';
          fetchNoteList(1);
        } else {
          // 禁止取消选择目录，改为刷新目录
          fetchNoteList();
        }
        break;
      }
      default: {
        break;
      }
    }
  }

  async function fetchUserTagList() {
    userTagList.value = [];
    const tags = await getAnnotateTags({ onlyUsed: false });
    userTagList.value = tags;
  }

  async function fetchNoteTagList() {
    const response = await getMarkTagListByFolderId({
      folderId: removePrefix(noteFolderSelected.value),
    });
    noteTagList.value = response;
    noteTagSelected.value = noteTagSelected.value.filter((id) =>
      noteTagList.value.some((mark) => mark.tagId === id)
    );
    if (!noteTagList.value.some((tag) => tag.tagId === noteTagEditId.value)) {
      noteTagEditId.value = '';
    }
  }

  async function fetchSummary(pageNumber = summaryPageNumber.value) {
    noteLoading.value = true;

    let response: null | GetNoteSummaryByNoteIdResponse | GetSummaryListByFolderIdResponse =
      null;
    const isDoc = noteFolderMap.value[noteFolderSelected.value]?.isDoc;
    try {
      summaryTotal.value = 0;
      if (isDoc) {
        response = await getSummaryNote({
          noteId: removePrefix(noteFolderSelected.value),
        });
        summaryNote.value = response;
      } else {
        response = await getSummaryListByFolderId({
          currentPage: pageNumber,
          pageSize: summaryPageSize.value,
          folderId: removePrefix(noteFolderSelected.value),
        });
        summaryList.value = response.docInfos || [];
        summaryTotal.value = response.total;
      }
    } catch (error) {
      console.warn('拉取总结失败', error);
      noteLoading.value = false;
      typeof error === 'string'
        ? error
        : error instanceof Object
          ? (error as any).message
          : '';
      return;
    }
    summaryPageNumber.value = pageNumber;
    summaryError.value = '';
    noteLoading.value = false;
  }

  async function fetchWordList(pageNumber = wordPageNumber.value) {
    noteLoading.value = true;
    // noteEditId.value = '';

    let response: null | GetNoteWordsByNoteIdResponse | GetWordListByFolderIdResponse =
      null;
    const isDoc = noteFolderMap.value[noteFolderSelected.value].isDoc;
    try {
      if (isDoc) {
        response = await getWordNotes({
          currentPage: pageNumber,
          pageSize: wordPageSize.value,
          noteId: removePrefix(noteFolderSelected.value),
        });
      } else {
        response = await getWordListByFolderId({
          currentPage: pageNumber,
          pageSize: wordPageSize.value,
          folderId: removePrefix(noteFolderSelected.value),
        });
      }
    } catch (error) {
      console.warn('拉取单词列表失败', error);
      noteLoading.value = false;
      wordListError.value =
        typeof error === 'string'
          ? error
          : error instanceof Object
            ? (error as any).message
            : '';
      return;
    }

    noteListError.value = '';

    wordList.value = response.words;
    wordTotal.value = response.total;
    wordPageNumber.value = pageNumber;

    noteLoading.value = false;
  }

  async function fetchNoteList(pageNumber = notePageNumber.value) {
    noteLoading.value = true;
    noteEditId.value = '';
    const isDoc = noteFolderMap.value[noteFolderSelected.value].isDoc;

    const searchContent = noteSearchKeyword.value;
    const styleIdList = Object.keys(noteStyle.value)
      .filter((key) => noteStyle.value[key as unknown as ColorKey])
      .map(Number);

    let response: null | MyNoteMarkListResponse = null;
    try {
      if (isDoc) {
        response = await autoDeprecateGetMyNoteMarkList({
          currentPage: pageNumber,
          pageSize: notePageSize.value,
          docId: noteFolderMap.value[noteFolderSelected.value].docId || '',
          tagIdList: noteTagSelected.value,
          sortType: noteSort.value,
          searchContent,
          styleIdList:
            styleIdList.length === noteStyleList.length ? [] : styleIdList,
        });
      } else {
        response = await autoDeprecateGetMyNoteMarkList({
          currentPage: pageNumber,
          pageSize: notePageSize.value,
          folderId: removePrefix(noteFolderSelected.value),
          tagIdList: noteTagSelected.value,
          sortType: noteSort.value,
          searchContent,
          styleIdList:
            styleIdList.length === noteStyleList.length ? [] : styleIdList,
        });
      }
      if (!response) {
        throw Error('No response');
      }
    } catch (error) {
      noteLoading.value = false;
      noteListError.value =
        typeof error === 'string'
          ? error
          : error instanceof Object
            ? (error as any).message
            : '';
      return;
    }

    noteListError.value = '';
    noteList.value =
      pageNumber === 1
        ? response.annotationModelList
        : noteList.value.concat(response.annotationModelList);
    noteTotal.value = response.total;
    notePageNumber.value = pageNumber;
    noteLoading.value = false;
  }
};

export const getCharWidth: (char: string, fontSize?: number) => number =
  (() => {
    const charList = Array(129)
      .fill(null)
      .map((_, index) => index)
      .map((code) => String.fromCharCode(code));

    const FULL_CHARACTER_SAMPLE = '全';
    charList.push(FULL_CHARACTER_SAMPLE);

    const charWidthMap: Record<number, Record<string, number>> = {};
    const REPEAT = 10;
    const createSpan = (char: string) => {
      const span = document.createElement('span');
      if (char === ' ') {
        span.innerHTML = Array(REPEAT).fill('&nbsp;').join('');
      } else {
        span.innerText = Array(REPEAT)
          // -会识别错误，用+代替
          .fill(char === '-' ? '+' : char)
          .join('');
      }
      return span;
    };
    const createWidthMap = (fontSize: number): void => {
      charWidthMap[fontSize] = {};
      const div = document.createElement('div');
      div.style.fontSize = `${fontSize}px`;
      div.style.opacity = String(0);
      const fragment = document.createDocumentFragment();
      const spanList = charList.map(createSpan);
      spanList.forEach((span) => {
        fragment.appendChild(span);
      });
      div.appendChild(fragment);
      document.body.appendChild(div);
      charList.forEach((char, index) => {
        charWidthMap[fontSize][char] = spanList[index].offsetWidth / REPEAT;
      });
      document.body.removeChild(div);
    };

    const DEFAULT_FONTSIZE = 14;
    const getServerSideCharWidth = (char: string, fontSize: number) => {
      const DEFAULT_HALF = 7;
      const DEFAULT_FULL = DEFAULT_HALF * 2;
      const code = char.charCodeAt(0);
      const width = code >= 0 && code <= 128 ? DEFAULT_HALF : DEFAULT_FULL;
      const scale = fontSize / DEFAULT_FONTSIZE;
      return width * scale;
    };

    return (char: string, fontSize = DEFAULT_FONTSIZE) => {
      if (typeof document === 'undefined') {
        return getServerSideCharWidth(char, fontSize);
      }

      if (!(fontSize in charWidthMap)) {
        createWidthMap(fontSize);
      }

      return char in charWidthMap[fontSize]
        ? charWidthMap[fontSize][char]
        : charWidthMap[fontSize][FULL_CHARACTER_SAMPLE];
    };
  })();

export const NOTES_CLASSNAME = 'readpaper-notes-container';
