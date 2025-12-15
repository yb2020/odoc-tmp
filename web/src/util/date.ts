export const convertDate = (modifyDate?: number) => {
  if (!modifyDate) return '';

  if (modifyDate.toString().length == 10) {
    modifyDate = modifyDate * 1000;
  }

  modifyDate = parseInt(modifyDate.toString());

  const jsdate = modifyDate ? new Date(modifyDate) : new Date();
  const nowdate = new Date(); //开始时间
  const tempdate = nowdate.getTime() - jsdate.getTime(); //时间差的毫秒数

  //计算出相差天数
  const days = Math.floor(tempdate / (24 * 3600 * 1000));
  if (days >= 1 && days < 10) {
    return days + '天前';
  }

  //计算出小时数
  const hours = Math.floor(tempdate / (3600 * 1000)); //计算天数后剩余的毫秒数
  if (hours >= 1 && hours < 24) {
    return hours + '小时前';
  }

  //计算相差分钟数
  const minutes = Math.floor(tempdate / (60 * 1000)); //计算小时数后剩余的毫秒数
  if (minutes >= 1 && minutes < 60) {
    return minutes + '分钟前';
  }

  //计算相差秒数
  const seconds = Math.floor(tempdate / 1000);
  if (seconds < 60) {
    return seconds + '秒钟前';
  }
  return jsdate.toLocaleString('chinese', { hour12: false });
};