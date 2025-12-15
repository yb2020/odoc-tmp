/*
 * Created Date: March 17th 2022, 4:37:01 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: March 17th 2022, 4:37:01 pm
 */
const REGEXP = /[xy]/g;
const PATTERN = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx';

function replacement(c: string) {
  const r = (Math.random() * 16) | 0;
  const v = c == 'x' ? r : (r & 0x3) | 0x8;
  return v.toString(16);
}

export default function uuid() {
  return PATTERN.replace(REGEXP, replacement);
}
