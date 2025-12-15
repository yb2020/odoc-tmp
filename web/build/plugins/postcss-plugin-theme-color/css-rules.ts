import { parseColorWithCache, rgbToHexString } from './color';

export enum ThemeType {
  beige = 'beige',
  dark = 'dark',
  green = 'green',
}

const themeVarPrefix = '--rp-theme';

export function modifyBorderColor(value: string) {
  const needChangedBorderColors = ['#e9ebf0', '#d9d9d9', '#e4e7ed', '#f5f7fa'];

  for (let i = 0; i < needChangedBorderColors.length; i += 1) {
    const color = needChangedBorderColors[i];
    const reg = new RegExp(color, 'ig');
    if (reg.test(value)) {
      return value.replace(
        color,
        `var(${themeVarPrefix}-bd-${color.substring(1)}, ${color})`
      );
    }
  }

  return null;
}

export function modifyForegroundColor(value: string) {
  const rgba = parseColorWithCache(value);
  if (!rgba) {
    return null;
  }
  const key = rgbToHexString(rgba);
  const needChangedFgColors = [
    '#ffffff',
    '#e4e6e9',
    '#f5f5f5d9',
    '#f5f7fa',
    '#e4e7ed',
    '#1f71e0',
    '#4e5969',
    '#1d2129',
    '#000000d9',
    '#262625',
    '#000000a3',
    '#00000059',
    '#73716f',
    '#00000040',
    '#1d2229',
    '#000000a6',
    '#e5e6eb',
  ];

  if (needChangedFgColors.includes(key)) {
    return `var(${themeVarPrefix}-fg-${key.slice(1)}, ${key})`;
  }
}

export function modifyBackgroundColor(value: string) {
  const rgba = parseColorWithCache(value);
  if (!rgba) {
    return null;
  }
  const needChangedBgColors = [
    '#ffffff',
    '#e4e6e9',
    '#f5f5f5d9',
    '#f5f7fa',
    '#e4e7ed',
    '#1f71e0',
    '#4e5969',
    '#1d2129',
    '#f5f5f5',
    '#6f77bb',
    '#f0f2f5',
    '#f7f8fa',
    '#e8f5ff',
    '#e5e6eb',
  ];

  const key = rgbToHexString(rgba);
  if (needChangedBgColors.includes(key)) {
    return `var(${themeVarPrefix}-bg-${key.slice(1)}, ${key})`;
  }
}

export function getModifiableCSSDeclaration(property: string, value: string) {
  if (property.startsWith('--')) {
    return null;
  }
  if (
    property === 'color-scheme' ||
    property === '-webkit-print-color-adjust'
  ) {
    return null;
  }

  if (/^background*/.test(property)) {
    return modifyBackgroundColor(value);
  }

  if (/^border*/.test(property)) {
    return modifyBorderColor(value);
  }

  return modifyForegroundColor(value);
}
