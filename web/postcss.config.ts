// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import postcssPxToViewport from 'postcss-px-to-viewport';
import postcssThemeColor from './build/plugins/postcss-plugin-theme-color';
import { Config } from 'postcss-load-config';

const config: Config = {
  plugins: [
    require('postcss-import'),
    require('tailwindcss/nesting'),
    require('tailwindcss'),
    require('autoprefixer'),
    postcssPxToViewport({
      unitToConvert: 'rpx',
      viewportWidth: 375,
      unitPrecision: 5,
      propList: ['*'],
      viewportUnit: 'vw',
      fontViewportUnit: 'vw',
      selectorBlackList: [],
      minPixelValue: 1,
      exclude: [],
    }),
    postcssThemeColor(),
  ],
};

export default config;
