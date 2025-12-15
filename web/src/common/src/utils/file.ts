import _ from 'lodash';

export const BYTES_KB = 1024 / 8;
export const BYTES_MB = BYTES_KB * 1024;
export const BYTES_GB = BYTES_MB * 1024;

export const bytes2GB = (size = 0, fix = 0) => {
  return `${_.round(size / BYTES_GB, fix)}`;
};

export const bytes2MB = (size = 0, fix = 0) => {
  return `${_.round(size / BYTES_MB, fix)}`;
};

export const bytes2KB = (size = 0, fix = 0) => {
  return `${_.round(size / BYTES_KB, fix)}`;
};

export const formatSize = (bytes = 0) => {
  if (bytes < BYTES_KB) {
    return `${bytes}B`;
  }
  if (bytes < BYTES_MB) {
    return `${bytes2KB(bytes)}KB`;
  }
  if (bytes < BYTES_GB) {
    return `${bytes2MB(bytes)}MB`;
  }

  return `${bytes2GB(bytes)}GB`;
};

export const formatBitSize = (bits = 0) => {
  return formatSize(bits / 8);
};
