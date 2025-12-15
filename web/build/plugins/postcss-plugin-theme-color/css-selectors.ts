const roots = ['html', ':root']

let uniqueRoots = roots
  if (uniqueRoots.includes('html')) {
    uniqueRoots = uniqueRoots.filter(i => i !== ':root')
  }



function escapeRegExp(string) {
  return string.replace(/[$()*+.?[\\\]^{|}-]/g, '\\$&')
}

function replaceAll(string, find, replace) {
  return string.replace(new RegExp(escapeRegExp(find), 'g'), replace)
}

export function processSelectors(selectors, add) {
  return selectors.map(selector => {
    let changed = false
    for (const root of roots) {
      if (selector.includes(root)) {
        changed = true
        selector = replaceAll(selector, root, root + add)
      }
    }
    if (!changed) {
      selector = uniqueRoots
        .map(root => `${root}${add} ${selector}`)
        .join(',')
    }
    return selector
  })
}