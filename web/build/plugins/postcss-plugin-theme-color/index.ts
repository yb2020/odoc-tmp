import postcss from 'postcss';
import { getModifiableCSSDeclaration } from './css-rules';

const postcssPluginThemeColor = () => {
  return {
    postcssPlugin: 'postcss-theme-color',
    Once(root: postcss.Root /*, { Rule }*/) {
      const themeRules: postcss.Rule[] = [];
      root.walkDecls((decl) => {
        if (
          (decl.parent as postcss.Rule).selectors
            ?.join(' ')
            ?.indexOf('.no-rp-theme') >= 0
        ) {
          return;
        }
        const modifiedValue = getModifiableCSSDeclaration(
          decl.prop,
          decl.value
        );
        if (!modifiedValue) {
          return;
        }

        decl.replaceWith(decl.clone({ value: modifiedValue }));
        // if (!isColorStyle(decl.prop)) {
        //   return
        // }
        // const value = decl.value;

        // const rgba = parseColorWithCache(value)

        // if (!rgba) {
        //   return;
        // }

        // [ThemeType.beige, ThemeType.green, ThemeType.dark].forEach((theme) => {
        //   const colorValue = modifyColor(rgba, theme)
        //   if (!colorValue) {
        //     return
        //   }
        //   const selectors = processSelectors((decl.parent as postcss.Rule).selectors, `[data-theme="${theme}"]`)
        //   const themeDecl = decl.clone({ value: colorValue, });
        //   const themeRule: postcss.Rule = new Rule({selector: selectors,})
        //   themeRule.append(themeDecl)
        //   themeRules.push(themeRule)
        // })
      });
      themeRules.forEach((rule) => {
        root.append(rule);
      });
    },
  };
};

export default postcssPluginThemeColor;
