import { ViewerController } from '@idea/pdf-annotate-viewer';

// 获取目录的其中一条
const getCatalogItems = (item: any, res: any[]) => {
  const catalogItem = {
    title: item.title,
    dest: item.dest,
    child: [],
  };

  item.items.forEach((i: any) => {
    getCatalogItems(i, catalogItem.child);
  });

  res.push(catalogItem);
};

// 从 PDF 文件读取目录
export const getPDFCatalog = async (viewer: ViewerController) => {
  const catalog = { title: '', pageNum: 0, child: [] };
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const { pdfDocument } = viewer;
  const outLine = await pdfDocument.getOutline();

  // 文件无自带目录，直接返回
  if (!outLine) {
    console.error('文件无自带目录', outLine);
    return catalog;
  }

  outLine.forEach(async (item: any) => {
    getCatalogItems(item, catalog.child);
  });

  return catalog;
};
