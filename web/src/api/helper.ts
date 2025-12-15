export function isEmptyPaperId(paperId: string) {
  return !paperId || paperId === '0';
}

export function checkOpenPaper(paperId: string, isPrivatePaper: boolean) {
  return !isEmptyPaperId(paperId) && !isPrivatePaper;
}
