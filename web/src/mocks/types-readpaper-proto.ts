// Mock types for @idea/types-readpaper-proto
export namespace UserDocManage {
  export interface UserDocDisplayAuthor {
    authorInfos: Array<{
      literal: string;
      [key: string]: any;
    }>;
    [key: string]: any;
  }
}

// Mock types for RedDot
export namespace RedDot {
  export enum FunctionTypeEnum {
    UNKNOWN = 0,
    LIBRARY = 1,
    NOTES = 2,
    SEARCH = 3,
    PAPER = 4,
    PROFILE = 5,
    WORKBENCH = 6,
    PREMIUM = 7,
    TEAM = 8,
    COLOR_PATTERN = 9,
  }
}
