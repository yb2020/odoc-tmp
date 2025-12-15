import _ from 'lodash'
import fs from 'fs'
import path from 'path'
import langJson from '../lang.properties.json';


type Obj = { [key: string]: any }; // 定义一个 Object 类型

const dfs = (tree: Obj, targetObj: Obj, valIdx: number) => {
  const stack = [{ tree, keys: Object.keys(tree), idx: 0, target: targetObj, }];
  while (stack.length) {
    const item = stack[stack.length - 1];
    const { tree, keys, idx, target } = item
    if (idx === keys.length) {
      stack.pop();
      continue;
    }
    const key = keys[idx];
    const val = tree[key];
    if ( _.isPlainObject(val)) {
      target[key] = {}
      stack.push({ tree: val, keys: Object.keys(val), idx: 0, target:  target[key]});
    } else {
      target[key] = val[valIdx]
    }
    item.idx++;
  }
}

const writeJson = () => {
  const langs = langJson.lang;
  fs.writeFileSync(path.join(__dirname, '../src/locals/lang.json'), JSON.stringify({ lang: langs }, null, 2));
  langs.forEach((lang, idx) => {
    const json = {};
    dfs(langJson, json, idx);
    fs.writeFileSync(path.join(__dirname, '../src/locals/' + lang + '.json'), JSON.stringify(json, null, 2));
  })

}

writeJson()





