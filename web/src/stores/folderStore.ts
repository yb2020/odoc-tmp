import { defineStore } from 'pinia';
import { nextTick } from 'vue';
import { addFolder, getDocsIndex } from '../api/material';
import { delay } from '@idea/aiknowledge-special-util/delay';
// import { UserDocFolderInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/UserDocManage';
import { UserDocFolderInfo } from 'go-sea-proto/gen/ts/doc/userDoc';

export interface PaperFolder {
  key: string;
  title: string;
  children: PaperFolder[];
  path: number[];
}

export interface FolderState {
  folderList: PaperFolder[];
  newFolderParent: PaperFolder | null;
  newFolderName: string;
}

let inited: null | Promise<void> = null;

export const useFolderStore = defineStore('folder', {
  state: (): FolderState => ({
    folderList: [],
    newFolderParent: null,
    newFolderName: '',
  }),
  actions: {
    async fetchFolderTree() {
      const res = await getDocsIndex({});
      const list = res!.folderInfos.map((child, index) =>
        convertFolder(child, [index])
      );
      this.folderList = list;

      function convertFolder(
        folder: UserDocFolderInfo,
        path: number[]
      ): PaperFolder {
        const key = folder.folderId;
        const children =
          folder.childrenFolders?.map((child, index) =>
            convertFolder(child, [...path, index])
          ) ?? [];
        return {
          key,
          title: folder.name,
          children,
          path,
        };
      }
    },

    init() {
      if (!inited) {
        inited = (async () => {
          try {
            await this.fetchFolderTree();
          } catch (error) {
            inited = null;
          }
        })();
      }
    },

    cancelFolder() {
      this.newFolderParent = null;
      this.newFolderName = '';
    },
    async newFolder(folder: PaperFolder) {
      this.cancelFolder();
      this.newFolderParent = folder;
      await nextTick();
      await delay(300);
      this.getFolderInput()?.focus();
    },
    async submitFolder() {
      const parent = this.newFolderParent as PaperFolder;

      await addFolder({
        parentId: parent.key,
        name: this.newFolderName,
        level: parent.path.length,
        sort: parent.children.length,
        oldFolderItems: parent.children.map((child, sort) => ({
          id: child.key,
          sort,
        })),
      });

      this.cancelFolder();
      return this.fetchFolderTree();
    },
    getFolderInput() {
      return document.querySelector<HTMLInputElement>(
        `input[data-folder-key="${this.newFolderParent?.key}"]`
      );
    },
  },
});

export const getCharLengthList = (string: string) => {
  const lengthList = Array(string.length).fill(1);
  for (let i = 0; i < string.length; i += 1) {
    const code = string.charCodeAt(i);
    if (code < 0 || code > 128) lengthList[i] += 1;
  }

  return lengthList;
};

export const limit30 = (string: string) => {
  const lengthList = getCharLengthList(string);
  const charList = [];
  let sum = 0;
  for (let i = 0; i < string.length; i += 1) {
    if (sum >= 30) {
      charList.push('…');
      break;
    }

    charList.push(string[i]);
    sum += lengthList[i];
  }

  return charList.join('');
};
