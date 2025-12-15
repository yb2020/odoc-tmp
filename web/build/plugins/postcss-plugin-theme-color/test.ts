import { parseColorWithCache } from './color';

const test1 = (value: string) => {
  const rgba = parseColorWithCache(value);
  console.log(rgba);
};

test1('rgb(31 113 224 / var(--tw-border-opacity))');
