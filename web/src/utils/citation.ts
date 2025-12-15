import _ from 'lodash'
import moment from 'moment'
import Cite from 'citation-js'
import '@citation-js/plugin-enw'
import {
  paper2MetaData,
  citeprocSys,
  renderBibl,
  PaperData,
  OutputFormat,
  RENDER_ERR,
  CitationRenderStatus,
} from './citation-types'

// 重新导出 CitationRenderStatus 以保持向后兼容
export { CitationRenderStatus } from './citation-types'

export enum CitationStyle {
  BIBTEX = 'BibTex',
  ENDNOTE = 'EndNote',
  MLA = 'MLA',
  APA = 'APA',
  CHICAGO = 'Chicago',
  GBT7714 = 'GB/T 7714',
}

// export type CitationStyleType =

export type StructuredAuthorName = {
  lastName?: string
  firstName?: string
  middleNames: Array<string>
}

export interface PaperRecord {
  title: string
  publishDate: string
  primaryVenue: string
  authorList: { name: string }[]
}

const TERMINATING_PUNCTUATION = { '!': true, '.': true, '?': true }

// BibTeX escaped characters per http://www.bibtex.org/SpecialSymbols/
// https://link.springer.com/content/pdf/bbm%3A978-3-319-06425-3%2F1.pdf
// http://tug.ctan.org/info/symbols/comprehensive/symbols-letter.pdf
const BibTexSpecialChars = [
  ['{', '\\{'],
  ['}', '\\}'],
  ['$', '\\$'],
  ['&', '\\&'],
  ['%', '\\%'],
  ['#', '\\#'],
  ['_', '\\_'],
  ['Ä', '{\\"A}'],
  ['ä', '{\\"a}'],
  ['Ë', '{\\"E}'],
  ['ë', '{\\"e}'],
  ['Ï', '{\\"I}'],
  ['ï', '{\\"i}'],
  ['Ö', '{\\"O}'],
  ['ö', '{\\"o}'],
  ['Ü', '{\\"U}'],
  ['ü', '{\\"u}'],
  ['Ÿ', '{\\"Y}'],
  ['ÿ', '{\\"y}'],
  ['ß', '{\\ss}'],
  ['Ø', '{\\O}'],
  ['ø', '{\\o}'],
  ['Á', "{\\'A}"],
  ['á', "{\\'a}"],
  ['Ć', "{\\'C}"],
  ['ć', "{\\'c}"],
  ['É', "{\\'E}"],
  ['é', "{\\'e}"],
  ['Í', "{\\'I}"],
  ['í', "{\\'i}"],
  ['Ó', "{\\'O}"],
  ['ó', "{\\'o}"],
  ['Ú', "{\\'U}"],
  ['ú', "{\\'u}"],
  ['Ý', "{\\'Y}"],
  ['ý', "{\\'y}"],
  ['À', '{\\`A}'],
  ['à', '{\\`a}'],
  ['È', '{\\`E}'],
  ['è', '{\\`e}'],
  ['Ì', '{\\`I}'],
  ['ì', '{\\`i}'],
  ['Ò', '{\\`O}'],
  ['ò', '{\\`o}'],
  ['Ù', '{\\`U}'],
  ['ù', '{\\`u}'],
  ['Â', '{\\^A}'],
  ['â', '{\\^a}'],
  ['Ê', '{\\^E}'],
  ['ê', '{\\^e}'],
  ['Î', '{\\^I}'],
  ['î', '{\\^i}'],
  ['Ô', '{\\^O}'],
  ['ô', '{\\^o}'],
  ['Û', '{\\^u}'],
  ['û', '{\\^u}'],
  ['Ã', '{\\~A}'],
  ['ã', '{\\~a}'],
  ['Ñ', '{\\~N}'],
  ['ñ', '{\\~n}'],
  ['Õ', '{\\~O}'],
  ['õ', '{\\~o}'],
  ['Æ', '{\\aE}'],
  ['æ', '{\\ae}'],
  ['Œ', '{\\OE}'],
  ['œ', '{\\oe}'],
  ['Č', '{\\vC}'],
  ['č', '{\\vc}'],
  ['Š', '{\\vS}'],
  ['š', '{\\vs}'],
  ['ž', '{\\vz}'],
  ['Å', '{\\AA}'],
  ['å', '{\\aa}'],
  ['™', '{\\texttrademark}'],
  ['®', '{\\textregistered}'],
  ['©', '{\\textcopyright}'],
  ['α', '$\\alpha$'],
  ['β', '$\\beta$'],
  ['γ', '$\\gamma$'],
  ['Γ', '$\\Gamma$'],
  ['δ', '$\\delta$'],
  ['Δ', '$\\Delta$'],
  ['ε', '$\\epsilon$'],
  ['ζ', '$\\zeta$'],
  ['η', '$\\eta$'],
  ['θ', '$\\theta$'],
  ['Θ', '$\\Theta$'],
  ['ι', '$\\iota$'],
  ['κ', '$\\kappa$'],
  ['λ', '$\\lambda$'],
  ['Λ', '$\\Lambda$'],
  ['μ', '$\\mu$'],
  ['ν', '$\\nu$'],
  ['ξ', '$\\xi$'],
  ['ο', '$\\omicron$'],
  ['π', '$\\pi$'],
  ['Π', '$\\Pi$'],
  ['ρ', '$\\rho$'],
  ['σ', '$\\sigma$'],
  ['Σ', '$\\Sigma$'],
  ['τ', '$\\tau$'],
  ['υ', '$\\upsilon$'],
  ['Υ', '$\\Upsilon$'],
  ['φ', '$\\phi$'],
  ['Φ', '$\\Phi$'],
  ['χ', '$\\chi$'],
  ['ψ', '$\\psi$'],
  ['Ψ', '$\\Psi$'],
  ['ω', '$\\omega$'],
  ['Ω', '$\\Omega$'],
]

export function encodeBibTexComponent(str: string): string {
  return BibTexSpecialChars.reduce((curStr, [character, replacement]) => {
    return curStr.replace(new RegExp(`[${character}]`, 'g'), replacement)
  }, str + '')
}

export function authorNameLastThenFirst(
  structuredAuthorName: StructuredAuthorName,
  initialize?: boolean
): string {
  let formatted: string
  const { firstName, lastName, middleNames } = structuredAuthorName
  if (lastName && firstName) {
    const middleName = middleNames.join(' ')
    if (initialize) {
      const middleInitial =
        middleName.length > 0 ? middleName.substr(0, 1).toUpperCase() + '.' : ''
      const firstNameInitial = firstName.substr(0, 1).toUpperCase() + '.'
      formatted = `${lastName}, ${firstNameInitial}${middleInitial}`
    } else {
      formatted = `${lastName}, ${firstName} ${middleName}`
    }
  } else {
    // Fallback to something
    formatted = lastName || firstName || middleNames[0] || ''
  }
  return formatted.trim()
}

export function authorNameFirstThenLast(
  structuredAuthorName: StructuredAuthorName
): string {
  let formatted: string
  const { firstName, lastName, middleNames } = structuredAuthorName
  if (lastName && firstName) {
    formatted = [firstName, ...middleNames].concat([lastName]).join(' ')
  } else {
    formatted = lastName || firstName || middleNames[0] || ''
  }
  return formatted.trim()
}

/**
 * Returns the paper title with the following formatting changes:
 *  - removes wrapping whitespace
 *  - appends a period, unless the title ends in a '.', '!' or '?'
 * @param {string} title the title
 * @returns {string}
 */
function formattedPaperTitle(title: string): string {
  // use a map for fast lookup -- if the title ends in one of these we don't append a trailing '.'
  const withoutWhitespace = _.trim(title) as string
  if (
    TERMINATING_PUNCTUATION[withoutWhitespace.substr(-1) as '!' | '.' | '?']
  ) {
    return `${withoutWhitespace}.`
  } else {
    return withoutWhitespace
  }
}

function getYearByPublishDate(date: string): string {
  let year = ''

  if (date !== '0') {
    const publishDate = moment(date)
    year = publishDate?.year() ? `${publishDate?.year()}` : ''
  }

  return year
}

export function getMLACitation(
  paper: PaperRecord,
  abbreviateAuthors: boolean
): string {
  const fields = []

  // Authors
  if (!_.isEmpty(paper.authorList)) {
    const authorsLen = paper.authorList.length
    if (abbreviateAuthors && authorsLen >= 3) {
      const firstAuthor = paper.authorList[0]
      fields.push(`${firstAuthor.name} et al.`)
    } else {
      fields.push(
        paper.authorList
          .map((author, index) => {
            if (index === 0) {
              const formattedName = author.name
              if (authorsLen === 1) {
                return `${formattedName}.`
              } else {
                return formattedName
              }
            } else {
              const formatted = author.name
              if (index === authorsLen - 1) {
                return ` and ${formatted}.`
              } else {
                return `, ${formatted}`
              }
            }
          })
          .join('')
      )
    }
  }

  // Title
  fields.push(`“${formattedPaperTitle(paper.title)}”`)

  // Journal / Venue Information and Year
  const year = `(${getYearByPublishDate(paper.publishDate)})`

  if (paper.primaryVenue) {
    const pages = ': n. pag'
    fields.push(`<em>${paper.primaryVenue}</em>${year}${pages}.`)
  } else {
    year && fields.push(`${year}.`)
  }

  return fields.join(' ')
}

export function getAPACitation(paper: PaperRecord): string {
  const fields = []

  // Authors
  if (!_.isEmpty(paper.authorList)) {
    fields.push(
      paper.authorList
        .map((author, index) => {
          const formatted = author.name
          if (index > 0 && index === paper.authorList.length - 1) {
            return `& ${formatted}`
          } else {
            return formatted
          }
        })
        .join(', ')
    )
  }

  // Year
  const year = getYearByPublishDate(paper.publishDate)
  fields.push(year ? `(${year}).` : '')

  // Title
  fields.push(formattedPaperTitle(paper.title))

  if (paper.primaryVenue) {
    fields.push(`<em>${paper.primaryVenue}</em>.`)
  }

  return fields.join(' ')
}

const PatternNonWordCharacters = /\W/g
const PatternOneOrMoreSpaces = /\s/

export function getBibTexCitation(paper: PaperRecord): string {
  // 转换为 PaperData 格式
  const paperData: PaperData = {
    docName: paper.title,
    authorList: paper.authorList?.map(author => ({
      literal: author.name,
      given: '', // 如果有分离的名字可以在这里处理
      family: '', // 如果有分离的姓氏可以在这里处理
    })) || [],
    primaryVenue: paper.primaryVenue,
    publishDate: paper.publishDate,
    // 添加更多字段支持
    doi: (paper as any).doi,
    url: (paper as any).url,
    language: (paper as any).language || 'en-US',
  };

  // 使用 citation-js 生成 BibTeX
  try {
    const metaData = paper2MetaData(paperData);
    const cite = new Cite(metaData);
    return cite.format('bibtex');
  } catch (error) {
    console.error('BibTeX generation failed, falling back to custom format:', error);
    
    // 降级到原有逻辑作为备选方案
  const hasJournal = !!paper.primaryVenue
  const type = hasJournal ? 'article' : 'inproceedings'
  const firstAuthor = paper.authorList?.[0]

  // For the third portion of the id we'd like to use the first word of the title and the first
  // letter of the second and third words (this is per user feedback)
  const titleWords: Array<string> = paper.title
    ? paper.title.split(PatternOneOrMoreSpaces)
    : []
  const abbreviatedTitle = [
    titleWords.shift(),
    titleWords.shift(),
    titleWords.shift(),
  ]
    .filter((str) => !!str)
    .map((str, i) => (i > 0 ? (str as string).substr(0, 1).toUpperCase() : str))
    .join('')

  const publishYear = getYearByPublishDate(paper.publishDate)

  const idParts: Array<string> = [publishYear, abbreviatedTitle]
  if (firstAuthor?.name) {
    idParts.unshift(firstAuthor.name)
  }

  const id = idParts
    .map((str) => str.replace(PatternNonWordCharacters, ''))
    .join('')
    .substr(0, 45)
  const authors: string = (paper.authorList || [])
    .map((author) => author.name)
    .join(' and ')

  const journalOrVenue = hasJournal
    ? `journal={${encodeBibTexComponent(paper.primaryVenue)}}`
      : `booktitle={${encodeBibTexComponent(paper.primaryVenue)}}`
  const year = `year={${encodeBibTexComponent(publishYear)}}`

  const fields = [
    `title={${encodeBibTexComponent(paper.title)}}`,
    `author={${encodeBibTexComponent(authors)}}`,
    journalOrVenue,
    year,
  ].filter((v) => !!v)

  const fieldsStr =
    fields.reduce((str, field) => `${str},\n  ${field || ''}`, '') + '\n'
  return `@${type}{${encodeBibTexComponent(id)}${fieldsStr}}`
  }
}

export function getEndNoteCitation(paperData: PaperData): string {
  const metaData = paper2MetaData(paperData)
  return Cite(metaData).format('enw', { format: 'text', lineEnding: '\n' })
}

export function getCitation(paper: PaperRecord, style: CitationStyle): string {
  transformPaperData(paper)
  switch (style) {
    case CitationStyle.CHICAGO:
    case CitationStyle.MLA: {
      return getMLACitation(paper, CitationStyle.MLA === style)
    }
    case CitationStyle.APA: {
      return getAPACitation(paper)
    }
    case CitationStyle.BIBTEX: {
      return getBibTexCitation(paper)
    }
    default: {
      throw new Error(`Unknown citation style: ${style}`)
    }
  }
}

export function getFilename(style: CitationStyle): string {
  switch (style) {
    case CitationStyle.BIBTEX:
      return 'citation.bib'
    case CitationStyle.ENDNOTE:
      return 'citation.enw'
    default:
      throw new Error(`Unknown citation style: ${style}`)
  }
}

export function transformPaperData(originData: any) {
  if (!originData.authorList && originData.authors) {
    originData.authorList = originData.authors.map((name: string) => ({ name }))
  }

  if (!originData.title && originData.docName) {
    originData.title = originData.docName
  }
}

interface warningMessageParams {
  [RENDER_ERR.NO_BIBL_STYLE]: string
}

export const getRenderCitation = async (
  style: string,
  paperData: PaperData,
  $t: any,
  outputFormat: OutputFormat,
  lang: string
) => {
  if (paperData) {
    const res = paper2MetaData(paperData)

    citeprocSys.addItem(res.id, res, true)

    const warningMessage: warningMessageParams = {
      [RENDER_ERR.NO_BIBL_STYLE]: $t(
        'home.paper.cite.warningMessage.renderEmpty'
      ) as string,
    }

    try {
      const result = await renderBibl(style, [], [], [res], outputFormat, lang)

      console.log('引文渲染数据', result)

      let txt = result.txt

      if (style.includes('bibtex')) {
        txt = txt
          ?.replace(/}, /g, '}, \n ')
          .replace(/, title={/, ',  \n title={')
          .replace(/} }/, '} \n }')
      }

      if (!txt) {
        return {
          status: CitationRenderStatus.empty,
          data: warningMessage[result.msgCode as keyof warningMessageParams],
        }
      }

      return {
        status: CitationRenderStatus.successful,
        data: txt?.replace(/<\/?em>/g, '').trim(),
      }
    } catch (error) {
      console.log('引文渲染失败', error)

      return {
        status: CitationRenderStatus.failed,
        data: (error as Error)?.message || '',
      }
    }
  }
}
