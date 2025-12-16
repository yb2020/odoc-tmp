package v8

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/beevik/etree"
	pb "github.com/yb2020/odoc/proto/gen/go/parsed"
	"github.com/yb2020/odoc/services/parse/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// 定义错误
var ErrInvalidDocument = errors.New("invalid document")

const (
	XMLNS = "http://www.w3.org/XML/1998/namespace"
	NS    = "http://www.tei-c.org/ns/1.0"
)

func ParseDocumentHeaderXml(xmlContent []byte) (*pb.DocumentHeader, error) {
	tree := etree.NewDocument()
	if err := tree.ReadFromBytes(xmlContent); err != nil {
		return nil, err
	}

	tei := tree.Root()
	if tei == nil {
		return nil, ErrInvalidDocument
	}

	// 解析基本的GROBID信息
	header := tei.FindElement(".//teiHeader")
	if header == nil {
		return nil, ErrInvalidDocument
	}

	applicationTags := header.FindElements(".//appInfo/application")
	if len(applicationTags) == 0 {
		return nil, ErrInvalidDocument
	}

	var (
		applicationTag = applicationTags[0]
		version        = strings.TrimSpace(applicationTag.SelectAttr("version").Value)
		ts             = strings.TrimSpace(applicationTag.SelectAttr("when").Value)
	)

	// 创建扩展文档对象
	doc := &pb.DocumentHeader{
		GrobidVersion: version,
		GrobidTs:      ts,
	}
	// 解析标题
	titleTag := header.FindElement(`.//titleStmt/title`)
	if titleTag != nil {
		doc.Title = &pb.Title{
			//标题名称
			Text: strings.Join(iterTextTrimSpace(titleTag), " "),
			//bbox
			Bbox: parseBBox(titleTag, nil),
		}
	}
	// 解析作者信息
	title, authors, dois, date, fileSHA256 := parseHeader(header)
	if title != nil {
		doc.Title = title
	}
	if authors != nil {
		doc.Authors = authors
	}
	if dois != nil {
		doc.Dois = dois
	}
	if date != "" {
		doc.Date = date
	}
	if fileSHA256 != "" {
		doc.FileSHA256 = fileSHA256
	}
	// 解析语言
	textTag := tei.FindElement(`.//text`)
	if textTag != nil {
		if lang := textTag.SelectAttrValue("lang", ""); lang != "" {
			doc.Lang = lang
		}
	}
	// 解析摘要
	abstractTag := tei.FindElement(`.//profileDesc/abstract`)
	if abstractTag != nil {
		doc.Abstract = &pb.Abstract{
			Text: strings.Join(iterTextTrimSpace(abstractTag), " "),
			Bbox: parseBBox(abstractTag, nil),
		}
	}

	return doc, nil
}

// ParseDocument 从XML读取器解析文档并转换为ExtendedDocument结构
func ParseDocument(xmlContent []byte) (*pb.DocumentMetadata, *pb.FullDocument, error) {
	tree := etree.NewDocument()
	if err := tree.ReadFromBytes(xmlContent); err != nil {
		return nil, nil, err
	}

	tei := tree.Root()
	if tei == nil {
		return nil, nil, ErrInvalidDocument
	}

	// 解析基本的GROBID信息
	header := tei.FindElement(".//teiHeader")
	if header == nil {
		return nil, nil, ErrInvalidDocument
	}

	applicationTags := header.FindElements(".//appInfo/application")
	if len(applicationTags) == 0 {
		return nil, nil, ErrInvalidDocument
	}

	var (
		applicationTag = applicationTags[0]
		version        = strings.TrimSpace(applicationTag.SelectAttr("version").Value)
		ts             = strings.TrimSpace(applicationTag.SelectAttr("when").Value)
	)

	// 创建扩展文档对象
	doc := &pb.DocumentMetadata{
		GrobidVersion: version,
		GrobidTs:      ts,
	}

	// 解析页面信息
	doc.Pages = parsePageInfo(tei)

	// 解析标题
	titleTag := header.FindElement(`.//titleStmt/title`)
	if titleTag != nil {
		doc.Title = &pb.Title{
			Text: strings.Join(iterTextTrimSpace(titleTag), " "),
			Bbox: parseBBox(titleTag, doc.Pages),
		}
	}

	// 解析作者信息
	title, authors, dois, date, fileSHA256 := parseHeader(header)
	if title != nil {
		doc.Title = title
	}
	if authors != nil {
		doc.Authors = authors
	}
	if dois != nil {
		doc.Dois = dois
	}
	if date != "" {
		doc.Date = date
	}
	if fileSHA256 != "" {
		doc.FileSHA256 = fileSHA256
	}

	// 解析语言
	textTag := tei.FindElement(`.//text`)
	if textTag != nil {
		if lang := textTag.SelectAttrValue("lang", ""); lang != "" {
			doc.Lang = lang
		}
	}

	// 解析摘要
	abstractTag := tei.FindElement(`.//profileDesc/abstract`)
	if abstractTag != nil {
		doc.Abstract = &pb.Abstract{
			Text: strings.Join(iterTextTrimSpace(abstractTag), " "),
			Bbox: parseBBox(abstractTag, doc.Pages),
		}
	}

	// 解析引用
	doc.References = parseReferences(textTag, doc.Pages)
	// 解析图表引用标记
	doc.FigureAndTableMarkers = parseFigureAndTableMarkers(string(xmlContent), doc.Pages)
	// 解析图表标记
	// doc.FigureAndTableMarkers = parseFigureTableMarkers(tei, doc.Pages)

	// 解析目录结构
	doc.Catalogue = parseCatalogue(tei, doc.Pages)

	// 解析致谢
	doc.Acknowledgment = parseAcknowledgement(textTag, doc.Pages)
	// 解析段落
	paragraphs, refMarkers := parseParagraphs(tei, doc.Pages)
	doc.ReferenceMarkers = refMarkers
	// 返回解析结果
	return doc, &pb.FullDocument{
		Paragraphs: paragraphs,
	}, nil
}

// parsePageInfo 解析文档的页面信息
func parsePageInfo(tei *etree.Element) []*pb.PageInfo {
	if tei == nil {
		return nil
	}

	var pages []*pb.PageInfo

	// 查找facsimile元素
	facsimile := tei.FindElement(`.//facsimile`)
	if facsimile == nil {
		return nil
	}

	// 解析每个surface元素（每个页面）
	surfaces := facsimile.FindElements(`.//surface`)
	for _, surface := range surfaces {
		// 获取页码
		pageNumAttr := surface.SelectAttr("n")
		if pageNumAttr == nil {
			continue
		}

		pageNum, err := strconv.Atoi(pageNumAttr.Value)
		if err != nil {
			continue
		}

		// 获取页面宽度和高度
		lrxAttr := surface.SelectAttr("lrx")
		lryAttr := surface.SelectAttr("lry")

		if lrxAttr == nil || lryAttr == nil {
			continue
		}

		width, err := strconv.ParseFloat(lrxAttr.Value, 64)
		if err != nil {
			continue
		}

		height, err := strconv.ParseFloat(lryAttr.Value, 64)
		if err != nil {
			continue
		}

		// 创建页面信息对象
		page := &pb.PageInfo{
			PageNumber: int32(pageNum),
			Width:      width,
			Height:     height,
		}

		pages = append(pages, page)
	}

	return pages
}

// getPageDimensions 根据页码获取页面尺寸
func getPageDimensions(pageNum int, pages []*pb.PageInfo) (width, height float64) {
	if pages == nil {
		return 0, 0
	}

	// 查找对应页码的页面信息
	for _, page := range pages {
		if page.PageNumber == int32(pageNum) {
			return page.Width, page.Height
		}
	}

	// 如果找不到对应页码，但有其他页面信息，返回第一个页面的尺寸
	if len(pages) > 0 {
		return pages[0].Width, pages[0].Height
	}

	return 0, 0
}

// parseHeader 解析文档头部信息
func parseHeader(header *etree.Element) (*pb.Title, []*pb.Author, []*pb.Doi, string, string) {
	if header == nil {
		return nil, nil, nil, "", ""
	}

	var title *pb.Title
	var date string
	var dois []*pb.Doi
	var authors []*pb.Author
	var fileSHA256 string
	// 解析标题
	titleTag := header.FindElement(`.//titleStmt/title`)
	if titleTag != nil {
		title = &pb.Title{
			Text: strings.Join(iterTextTrimSpace(titleTag), " "),
			Bbox: parseBBox(titleTag, nil),
		}
	}

	// 解析日期
	dateTag := header.FindElement(`.//publicationStmt/date`)
	if dateTag != nil {
		date = dateTag.Text()
	}

	// DOI
	idnoTags := header.FindElements(`.//sourceDesc/biblStruct/idno`)
	for _, idnoTag := range idnoTags {
		//md5也是在这个里面，type = MD5
		if idnoTag.SelectAttrValue("type", "") == "arXiv" {
			dois = append(dois, &pb.Doi{
				Doi:  idnoTag.Text(),
				Type: "arXiv",
			})
		}
		if idnoTag.SelectAttrValue("type", "") == "MD5" {
			dois = append(dois, &pb.Doi{
				Doi:  idnoTag.Text(),
				Type: "MD5",
			})
			fileSHA256 = idnoTag.Text()
		}
	}
	// 解析作者
	authorTags := header.FindElements(`.//sourceDesc/biblStruct/analytic/author`)
	for _, authorTag := range authorTags {
		author := parseAuthor(authorTag)
		if author != nil {
			authors = append(authors, author)
		}
	}

	return title, authors, dois, date, fileSHA256
}

// parseAuthor 解析作者信息
func parseAuthor(elem *etree.Element) *pb.Author {
	if elem == nil {
		return nil
	}

	// 解析姓名
	persNameTag := elem.FindElement(`./persName`) // 只查找直接子元素，不递归
	if persNameTag == nil {
		return nil // 如果没有persName标签，返回nil
	}

	author := &pb.Author{}

	// 解析全名
	fullNameText := strings.Join(iterTextTrimSpace(persNameTag), " ")
	if fullNameText == "" {
		// 如果没有有效的全名，尝试从forename和surname构建
		forenameTag := persNameTag.FindElement(`./forename`)
		surnameTag := persNameTag.FindElement(`./surname`)

		forename := ""
		surname := ""

		if forenameTag != nil {
			forename = forenameTag.Text()
			author.GivenName = forename
		}

		if surnameTag != nil {
			surname = surnameTag.Text()
			author.Surname = surname
		}

		// 如果名和姓都为空，返回nil
		if forename == "" && surname == "" {
			return nil
		}

		// 构建全名
		author.FullName = strings.TrimSpace(forename + " " + surname)
	} else {
		author.FullName = fullNameText

		// 解析名
		forenameTag := persNameTag.FindElement(`./forename`)
		if forenameTag != nil {
			author.GivenName = forenameTag.Text()
		}

		// 解析姓
		surnameTag := persNameTag.FindElement(`./surname`)
		if surnameTag != nil {
			author.Surname = surnameTag.Text()
		}
	}

	// 解析邮箱
	emailTag := elem.FindElement(`./email`) // 只查找直接子元素
	if emailTag != nil {
		author.Email = emailTag.Text()
	}

	// 解析位置信息（如果有）
	author.Bbox = parseBBox(persNameTag, nil)

	// 最终检查：如果所有主要字段都为空，返回nil
	if author.FullName == "" && author.GivenName == "" &&
		author.Surname == "" && author.Email == "" {
		return nil
	}

	return author
}

// parseReferences 解析引用列表
func parseReferences(textTag *etree.Element, pages []*pb.PageInfo) []*pb.Reference {
	if textTag == nil {
		return nil
	}
	var references []*pb.Reference
	divs := textTag.FindElements(`.//back/div`)
	if len(divs) == 0 {
		return nil
	}
	var referencesTag *etree.Element
	for _, div := range divs {
		if div.SelectAttrValue("type", "") == "references" {
			referencesTag = div
		}
	}
	if referencesTag == nil {
		return nil
	}
	refTags := referencesTag.FindElements(`.//listBibl/biblStruct`)

	for i, refTag := range refTags {
		ref := &pb.Reference{
			RefIdx: fmt.Sprintf("b%d", i),
		}

		// 解析标题
		titleTag := refTag.FindElement(`.//title`)
		if titleTag != nil {
			ref.Title = strings.Join(iterTextTrimSpace(titleTag), " ")
		}
		// 解析位置信息
		ref.Bbox = parseBBox(refTag, pages)

		// 解析作者信息
		var authors []*pb.Author
		for _, authorTag := range refTag.FindElements(`.//monogr/author`) {
			author := parseReferenceAuthor(authorTag)
			if author != nil {
				authors = append(authors, author)
			}
		}
		ref.Authors = authors

		// 解析发布日期
		dateTag := refTag.FindElement(`.//monogr/imprint/date`)
		if dateTag != nil {
			ref.PublishDate = strings.Join(iterTextTrimSpace(dateTag), " ")
		}

		// 解析arXiv ID
		arxivIDTag := refTag.FindElement(`.//idno[@type="arXiv"]`)
		if arxivIDTag != nil {
			ref.ArxivId = strings.TrimSpace(arxivIDTag.Text())
			// 移除前缀"arXiv:"
			if strings.HasPrefix(ref.ArxivId, "arXiv:") {
				ref.ArxivId = ref.ArxivId[6:]
			}
		}

		// 解析原始引用文本
		rawRefTags := refTag.FindElements(`.//note`)
		for _, rawRefTag := range rawRefTags {
			if rawRefTag.SelectAttrValue("type", "") == "raw_reference" {
				ref.ContentText = strings.Join(iterTextTrimSpace(rawRefTag), " ")
			}
		}

		references = append(references, ref)
	}

	return references
}

// parseReferenceAuthor 解析引用中的作者信息
func parseReferenceAuthor(elem *etree.Element) *pb.Author {
	if elem == nil {
		return nil
	}

	// 解析作者姓名
	persNameTag := elem.FindElement(`.//persName`)
	if persNameTag == nil {
		return nil
	}

	author := &pb.Author{}

	// 解析姓
	surnameTag := persNameTag.FindElement(`.//surname`)
	if surnameTag != nil {
		author.Surname = strings.TrimSpace(surnameTag.Text())
	}

	// 解析名
	forenameTag := persNameTag.FindElement(`.//forename`)
	if forenameTag != nil {
		author.GivenName = strings.TrimSpace(forenameTag.Text())
	}

	// 组合全名
	if author.GivenName != "" && author.Surname != "" {
		author.FullName = author.GivenName + " " + author.Surname
	} else if author.Surname != "" {
		author.FullName = author.Surname
	} else if author.GivenName != "" {
		author.FullName = author.GivenName
	}

	// 解析邮箱
	emailTag := elem.FindElement(`.//email`)
	if emailTag != nil {
		author.Email = strings.TrimSpace(emailTag.Text())
	}

	// 如果没有任何有效信息，返回nil
	if author.FullName == "" && author.Email == "" {
		return nil
	}

	return author
}

// parseCatalogue 解析目录结构
// createAndAddParentItem 创建并添加缺失的父级目录项
func createAndAddParentItem(itemFirstChar string, divIndex int, divTags []*etree.Element, headTag *etree.Element,
	pages []*pb.PageInfo, lastItem *pb.CatalogueItem, currentItem *pb.CatalogueItem, catalogueItems *[]*pb.CatalogueItem) {

	// 默认使用当前节点的位置信息和基本属性创建父节点
	supplementItem := &pb.CatalogueItem{
		Order:      currentItem.Order,
		Level:      "1",
		TitleOrder: itemFirstChar,
	}

	// 尝试从上一个节点的最后一个段落的最后一个句子中获取标题和位置信息
	if divIndex > 0 {
		prevDivTag := divTags[divIndex-1]
		pTags := prevDivTag.FindElements(`./p`)

		if len(pTags) > 0 {
			lastPTag := pTags[len(pTags)-1]
			sTags := lastPTag.FindElements(`./s`)

			// 如果找到了单个句子，使用它的内容作为父级目录标题
			if len(sTags) == 1 {
				sTag := sTags[0]
				title := strings.TrimSpace(sTag.Text())
				lastItemOrder := lastItem.Order + 1
				if title == "" {
					supplementItem.Title = itemFirstChar
				} else {
					//判断内容是否为数字开头 并且数字和当前currentItem.Order序号一样
					if !strings.HasPrefix(title, itemFirstChar) {
						supplementItem.Title = fmt.Sprintf("%s %s", itemFirstChar, title)
					} else {
						supplementItem.Title = title
					}
				}
				supplementItem.Order = lastItemOrder

				// 使用句子的坐标信息作为父级目录的位置信息
				if coordsAttr := sTag.SelectAttr("coords"); coordsAttr != nil {
					sentenceBbox := parseBBox(sTag, pages)
					if sentenceBbox != nil {
						supplementItem.Bbox = sentenceBbox
					}
				}
				// 标记这个句子已被用作目录项，后续解析时可以跳过
				sTag.CreateAttr("used-as-catalogue", "true")
				//格式化标题
				supplementItem.Title = util.NormalizeTitle(supplementItem.Title)
				//根据空格拆分
				matches := strings.Split(supplementItem.Title, " ")
				if len(matches) == 2 {
					supplementItem.TitleOrder = matches[0]
				}
				*catalogueItems = append(*catalogueItems, supplementItem)
				currentItem.Order = currentItem.Order + 1
				return
			}
		}
	}

	// 如果无法从上一个节点获取信息，使用当前节点的信息
	//格式化标题
	supplementItem.Title = util.NormalizeTitle(supplementItem.Title)
	//根据空格拆分
	matches := strings.Split(supplementItem.Title, " ")
	if len(matches) == 2 {
		supplementItem.TitleOrder = matches[0]
	}
	supplementItem.Bbox = parseBBox(headTag, pages)
	*catalogueItems = append(*catalogueItems, supplementItem)
	currentItem.Order = currentItem.Order + 1
}

func parseCatalogue(tei *etree.Element, pages []*pb.PageInfo) []*pb.CatalogueItem {
	if tei == nil {
		return nil
	}
	var catalogueItems []*pb.CatalogueItem

	// 如果没有找到content标签或其中没有目录项，尝试从正文中提取
	divTags := tei.FindElements(`.//text/body/div`)
	// 用于跟踪上一个目录项
	var lastItem *pb.CatalogueItem

	for i, divTag := range divTags {
		// 解析章节标题
		headTag := divTag.FindElement(`./head`)
		if headTag == nil {
			continue
		}
		title := strings.Join(iterTextTrimSpace(headTag), " ")
		item := &pb.CatalogueItem{
			Order: int32(i + 1), // 设置原始顺序
		}

		// 获取章节序号
		if nAttr := headTag.SelectAttr("n"); nAttr != nil {
			item.TitleOrder = nAttr.Value
			// 根据n属性中的数字和点的数量确定目录级别
			// 例如：1或1.为一级目录，1.1或1.1.为二级目录，以此类推
			trimmedValue := strings.TrimRight(nAttr.Value, ".")
			parts := strings.Split(trimmedValue, ".")
			item.Level = strconv.Itoa(len(parts))
			//设置标题 如果标题前面没有序号  则增加序号
			if !strings.HasPrefix(title, item.TitleOrder) {
				item.Title = fmt.Sprintf("%s %s", item.TitleOrder, title)
			}
		} else {
			// 如果没有n属性，默认不属于目录
			// item.Level = "1"
			continue
		}
		// 处理目录跳级情况
		if item.Level != "1" && lastItem != nil && item.TitleOrder != "" && lastItem.TitleOrder != "" {
			// 获取当前项和上一项的第一个字符
			itemFirstChar := item.TitleOrder[0:1]
			lastItemFirstChar := lastItem.TitleOrder[0:1]
			// 判断第一个字符是否不同，表示可能存在目录跳级情况
			if itemFirstChar != lastItemFirstChar {
				// 创建并添加缺失的父级目录项
				createAndAddParentItem(itemFirstChar, i, divTags, headTag, pages, lastItem, item, &catalogueItems)
				fmt.Printf("\u68c0测到目录跳级: 从 %s 到 %s\n", lastItem.TitleOrder, item.TitleOrder)
			}
		}
		// 格式化标题
		// item.Title = util.NormalizeTitle(item.Title)
		item.FormattedTitle = util.NormalizeTitle(item.Title)
		// 拆分格式化标题
		titleOrder, _ := util.SplitNormalizedTitle(item.FormattedTitle)
		item.TitleOrder = titleOrder
		// 解析位置信息
		item.Bbox = parseBBox(headTag, pages)
		catalogueItems = append(catalogueItems, item)
		//记录最新的目录项
		lastItem = item
	}
	//清空lastItem
	lastItem = nil
	// 如果找到了目录项，构建层级结构
	if len(catalogueItems) > 0 {
		return buildLayeredCatalogue(catalogueItems)
	}
	return nil
}

func buildLayeredCatalogue(items []*pb.CatalogueItem) []*pb.CatalogueItem {
	// 首先按照Order字段排序，确保同级目录项保持原始顺序
	sort.Slice(items, func(i, j int) bool {
		return items[i].Order < items[j].Order
	})

	var result []*pb.CatalogueItem

	// 创建一个映射，用于存储每个层级的项目
	levelItems := make(map[int][]*pb.CatalogueItem)

	// 按层级分组
	for _, item := range items {
		if item.Level == "" {
			// 如果没有层级信息，作为顶级项
			result = append(result, item)
			continue
		}

		// 尝试将level解析为整数
		level, err := strconv.Atoi(item.Level)
		if err != nil {
			// 如果无法解析为整数，作为顶级项
			result = append(result, item)
			continue
		}

		if level <= 0 {
			// 非法层级值，作为顶级项
			result = append(result, item)
			continue
		}

		// 将项目添加到对应层级的列表中
		levelItems[level] = append(levelItems[level], item)
	}

	// 处理一级目录项
	if level1Items, ok := levelItems[1]; ok {
		result = append(result, level1Items...)
	}

	// 处理层级关系
	// 遍历每个层级（从2开始）
	for level := 2; level <= len(levelItems)+1; level++ {
		items, ok := levelItems[level]
		if !ok {
			continue
		}

		// 遍历当前层级的所有项
		for _, item := range items {
			// 找到合适的父项
			var foundParent bool

			// 从当前项的位置向前查找最近的上一级项
			for i := len(levelItems[level-1]) - 1; i >= 0; i-- {
				potentialParent := levelItems[level-1][i]

				// 检查顺序号是否合适（父项的顺序号应小于子项）
				if potentialParent.Order < item.Order {
					// 找到合适的父项，添加为子项
					potentialParent.Child = append(potentialParent.Child, item)
					foundParent = true
					break
				}
			}

			// 如果找不到合适的父项，作为顶级项
			if !foundParent {
				result = append(result, item)
			}
		}
	}

	return result
}

func parseAcknowledgement(textTag *etree.Element, pages []*pb.PageInfo) *pb.Acknowledgment {
	// 只查找back标签下的div标签
	divs := textTag.FindElements(`.//back/div`)
	if len(divs) == 0 {
		return nil
	}
	var acknowledgementTag *etree.Element
	for _, div := range divs {
		if div.SelectAttrValue("type", "") == "acknowledgement" {
			acknowledgementTag = div
		}
	}
	if acknowledgementTag == nil {
		return nil
	}
	// 获取致谢内容的全文
	text := strings.Join(iterTextTrimSpace(acknowledgementTag), " ")

	// 查找第一个 p 标签，用于获取坐标
	pTag := acknowledgementTag.FindElement(`.//p`)
	var bbox *pb.BBox

	if pTag != nil {
		// 使用第一个 p 标签的坐标
		bbox = parseBBox(pTag, pages)
	} else {
		// 如果没有 p 标签，使用整个致谢标签的坐标
		bbox = parseBBox(acknowledgementTag, pages)
	}

	return &pb.Acknowledgment{
		Text: text,
		Bbox: bbox,
	}
}

// parseBBox 从元素的coords属性中解析边界框信息
func parseBBox(elem *etree.Element, pages []*pb.PageInfo) *pb.BBox {
	if elem == nil {
		return nil
	}

	// 在GROBID中，边界框信息存储在coords属性中
	coordsAttr := elem.SelectAttr("coords")
	if coordsAttr == nil {
		return nil
	}

	// coords属性格式为：page,x,y,width,height;page,x,y,width,height;...
	// 我们只取第一个边界框（如果有多个的话）
	coordsValue := coordsAttr.Value

	// 如果有多个边界框（用分号分隔），只取第一个
	parts := strings.Split(coordsValue, ";")
	if len(parts) == 0 {
		return nil
	}

	// 解析第一个边界框的5个属性
	attributes := strings.Split(parts[0], ",")
	if len(attributes) < 5 {
		return nil // 格式不正确
	}

	// 解析各个属性
	pageNum, err := strconv.Atoi(attributes[0])
	if err != nil {
		return nil
	}

	x, err := strconv.ParseFloat(attributes[1], 64)
	if err != nil {
		return nil
	}

	y, err := strconv.ParseFloat(attributes[2], 64)
	if err != nil {
		return nil
	}

	width, err := strconv.ParseFloat(attributes[3], 64)
	if err != nil {
		return nil
	}

	height, err := strconv.ParseFloat(attributes[4], 64)
	if err != nil {
		return nil
	}

	// 获取页面尺寸
	pageWidth, pageHeight := getPageDimensions(pageNum, pages)

	// 创建并返回边界框对象
	// 注意：GROBID的坐标系是左上角为原点，x向右，y向下
	return &pb.BBox{
		X0:           x,
		Y0:           y,
		X1:           x + width,
		Y1:           y + height,
		OriginHeight: pageHeight,
		OriginWidth:  pageWidth,
		PageNumber:   int32(pageNum),
	}
}

// iterText 递归返回所有文本元素，按文档顺序
func iterText(elem *etree.Element) (result []string) {
	if elem == nil {
		return
	}
	result = append(result, elem.Text())
	for _, ch := range elem.ChildElements() {
		result = append(result, iterText(ch)...)
	}
	result = append(result, elem.Tail())
	return result
}

// iterTextTrimSpace 递归返回所有子文本元素，按文档顺序，去除所有空白
func iterTextTrimSpace(elem *etree.Element) (result []string) {
	if elem == nil {
		return
	}
	for _, v := range iterText(elem) {
		c := strings.TrimSpace(v)
		if len(c) == 0 {
			continue
		}
		result = append(result, c)
	}
	return result
}

// SectionInfo 保存章节信息及其顺序
type SectionInfo struct {
	Items     map[string]*pb.CatalogueItem // 章节ID到章节项的映射
	OrderList []string                     // 保持章节ID的顺序
}

// extractSectionMap 提取文档中的所有章节信息到映射中
func extractSectionMap(elem *etree.Element, sectionInfo *SectionInfo, parentPath string) {
	if elem == nil {
		return
	}

	// 查找所有div元素
	divElems := elem.FindElements(".//div")
	for _, divElem := range divElems {
		// 查找章节标题
		head := divElem.FindElement("./head")
		if head != nil {
			title := strings.Join(iterTextTrimSpace(head), " ")

			// 获取章节序号
			titleOrder := ""
			if nAttr := head.SelectAttr("n"); nAttr != nil {
				titleOrder = nAttr.Value
				// 如果标题中不包含序号，添加序号到标题
				if !strings.Contains(title, titleOrder) {
					title = titleOrder + " " + title
				}
			}

			// 创建目录项
			item := &pb.CatalogueItem{
				Title:      title,
				TitleOrder: titleOrder,
				Bbox:       parseBBox(head, nil),
			}

			// 生成唯一ID
			id := divElem.SelectAttrValue("xml:id", "")
			if id == "" {
				// 如果没有ID，使用标题作为ID
				id = title
			}

			// 将章节信息添加到映射中
			sectionInfo.Items[id] = item
			// 记录章节ID的顺序
			sectionInfo.OrderList = append(sectionInfo.OrderList, id)

			// 递归处理子章节
			for _, childDiv := range divElem.FindElements("./div") {
				extractSectionMap(childDiv, sectionInfo, parentPath)
			}
		}
	}
}

// parseParagraphs 解析文档中的段落
func parseParagraphs(tei *etree.Element, pages []*pb.PageInfo) ([]*pb.Paragraph, []*pb.RefMarker) {
	if tei == nil {
		return nil, nil
	}

	var paragraphs []*pb.Paragraph
	var index int = 0

	// 查找文档正文
	textElem := tei.FindElement(".//text")
	if textElem == nil {
		return nil, nil
	}

	// 查找正文中的所有段落
	bodyElem := textElem.FindElement(".//body")
	if bodyElem == nil {
		return nil, nil
	}
	var refMarkers []*pb.RefMarker
	// 首先提取所有章节信息，创建章节映射和有序列表
	sectionInfo := &SectionInfo{
		Items:     make(map[string]*pb.CatalogueItem),
		OrderList: make([]string, 0),
	}
	extractSectionMap(bodyElem, sectionInfo, "")

	// 查找所有段落元素
	pElems := bodyElem.FindElements(".//p")
	// 按顺序处理段落，先按章节顺序组织
	paragraphsBySection := make(map[string][]*etree.Element)
	for _, pElem := range pElems {
		// 跳过引用部分的段落
		if isInReference(pElem) {
			continue
		}
		// 获取段落所属的章节ID
		sectionID, _ := findSectionTitle(pElem)
		paragraphsBySection[sectionID] = append(paragraphsBySection[sectionID], pElem)
	}

	// 按章节顺序处理段落
	for _, sectionID := range sectionInfo.OrderList {
		for _, pElem := range paragraphsBySection[sectionID] {
			// 跳过引用中的段落
			if isInReference(pElem) {
				continue
			}
			// 创建段落对象
			paragraph := &pb.Paragraph{
				Order: int32(index),
			}
			index++
			// 解析段落文本，处理引用标记
			text, refInfos, _ := parseTextWithReferences(pElem, pages)
			// 将引用信息转换为RefMarker
			for _, refInfo := range refInfos {
				refMarkers = append(refMarkers, &pb.RefMarker{
					RefIdx:     refInfo.Target,
					Bbox:       refInfo.Bbox,
					RefContent: refInfo.Text,
				})
			}
			paragraph.Text = text
			paragraph.References = refInfos
			// 设置段落所属的章节信息
			paragraph.SectionId = sectionID
			// 从章节映射中获取标题
			if sectionItem, ok := sectionInfo.Items[sectionID]; ok {
				paragraph.SectionTitle = sectionItem.Title
			}
			// 解析段落中的句子
			sentences := parseSentences(pElem, pages)
			paragraph.Text.Sentences = sentences
			//段落类型
			paragraph.Type = pb.ParagraphType_TEXT
			// 按照文档顺序添加段落
			paragraphs = append(paragraphs, paragraph)
		}
	}
	// 确保段落按照位置排序
	sortParagraphsByPosition(paragraphs)

	return paragraphs, refMarkers
}

// isInReference 检查段落元素是否在引用部分内
func isInReference(elem *etree.Element) bool {
	if elem == nil {
		return false
	}

	// 向上查找父元素，检查是否在引用部分内
	parent := elem.Parent()
	for parent != nil {
		if parent.Tag == "listBibl" {
			return true
		}
		parent = parent.Parent()
	}

	return false
}

// findSectionTitle 查找段落所属的章节标题和ID
func findSectionTitle(elem *etree.Element) (string, string) {
	if elem == nil {
		return "", ""
	}

	// 向上查找最近的div元素
	parent := elem.Parent()
	for parent != nil {
		if parent.Tag == "div" {
			// 查找该div下的head元素
			head := parent.FindElement("./head")
			if head != nil {
				title := strings.Join(iterTextTrimSpace(head), " ")

				// 获取章节序号
				if nAttr := head.SelectAttr("n"); nAttr != nil {
					titleOrder := strings.TrimSpace(nAttr.Value)
					// 如果标题中不包含序号，添加序号到标题
					if !strings.Contains(title, titleOrder) {
						title = titleOrder + " " + title
					}
				}

				// 获取章节ID
				id := parent.SelectAttrValue("xml:id", "")
				if id == "" {
					// 如果没有ID，使用标题作为ID 去除前后空格的id
					id = strings.TrimSpace(title)
				}

				return title, id
			}
			break
		}
		parent = parent.Parent()
	}

	return "", ""
}

// parseSentences 解析段落中的句子
func parseSentences(elem *etree.Element, pages []*pb.PageInfo) []*pb.Sentence {
	if elem == nil {
		return nil
	}

	var sentences []*pb.Sentence

	// 查找所有句子元素
	sElems := elem.FindElements(".//s")
	index := 0
	for _, sElem := range sElems {

		//由于grobid存在一个解析bug，目录跳级的情况使用了句子作为最后一个目录，所以这里需要进行过滤
		sentence := strings.Join(iterTextTrimSpace(sElem), " ")
		if sElem.SelectAttrValue("used-as-catalogue", "") != "" {
			continue
		}

		if sentence != "" {
			bbox := parseBBox(sElem, pages)
			sentences = append(sentences, &pb.Sentence{
				Text:  sentence,
				Bbox:  bbox,
				Index: int32(index),
			})
		}
		index++
	}

	// 如果没有明确的句子标记，则将整个段落作为一个句子
	if len(sentences) == 0 {
		text := strings.Join(iterTextTrimSpace(elem), " ")
		if text != "" {
			sentences = append(sentences, &pb.Sentence{
				Text:  text,
				Bbox:  parseBBox(elem, pages),
				Index: int32(index),
			})
		}
	}

	return sentences
}

// determineParagraphType 确定段落类型
func determineParagraphType(elem *etree.Element) string {
	if elem == nil {
		return "normal"
	}

	// 检查段落的类型属性
	pType := elem.SelectAttrValue("type", "")
	if pType != "" {
		return pType
	}

	// 检查是否在摘要部分
	parent := elem.Parent()
	for parent != nil {
		if parent.Tag == "abstract" {
			return "abstract"
		}
		parent = parent.Parent()
	}

	// 默认为普通段落
	return "normal"
}

// parseTextWithReferences 解析段落文本，处理引用标记
func parseTextWithReferences(pElem *etree.Element, pages []*pb.PageInfo) (*pb.Text, []*pb.RefInfo, []string) {
	if pElem == nil {
		return nil, nil, nil
	}
	bbox := parseBBox(pElem, pages)
	var refInfos []*pb.RefInfo
	var refIds []string

	// 获取段落的纯文本内容（不包括XML标签）
	// 这里我们只获取直接文本内容，不包括子元素中的文本
	var textContent string
	for _, child := range pElem.Child {
		if token, ok := child.(*etree.CharData); ok {
			textContent += strings.TrimSpace(token.Data)
		}
	}

	// 如果文本内容为空或只包含空白字符，可能所有内容都在子元素中
	if strings.TrimSpace(textContent) == "" {
		// 获取所有非引用元素的文本
		for _, child := range pElem.ChildElements() {
			if child.Tag != "ref" {
				textContent += " " + strings.Join(iterTextTrimSpace(child), " ")
			}
		}
		textContent = strings.TrimSpace(textContent)
	}

	// 如果仍然为空，获取整个元素的文本
	if textContent == "" {
		textContent = strings.Join(iterTextTrimSpace(pElem), " ")
	}

	// 收集所有引用信息
	refElems := pElem.FindElements(".//ref[@type='bibr']")
	for _, refElem := range refElems {
		// 获取引用ID
		target := refElem.SelectAttrValue("target", "")
		if target != "" && strings.HasPrefix(target, "#") {
			target = target[1:] // 去掉前缀 '#'
		}

		// 获取引用文本
		refText := strings.Join(iterTextTrimSpace(refElem), " ")

		// 解析边界框
		bbox := parseBBox(refElem, pages)

		// 添加到引用信息列表
		refInfos = append(refInfos, &pb.RefInfo{
			Text:   refText,
			Target: target,
			Bbox:   bbox,
		})

		// 添加到引用ID列表
		if target != "" {
			refIds = append(refIds, target)
		}
	}

	return &pb.Text{
		Bbox: bbox,
		Text: textContent,
	}, refInfos, refIds
}

// sortParagraphsByPosition 确保段落按照页码和位置排序
func sortParagraphsByPosition(paragraphs []*pb.Paragraph) {
	sort.Slice(paragraphs, func(i, j int) bool {
		// 如果没有位置信息，保持原有顺序
		return paragraphs[i].Order < paragraphs[j].Order
	})
}

func parseFigureAndTableMarkers(xmlContent string, pages []*pb.PageInfo) []*pb.RefMarker {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(xmlContent))
	if err != nil {
		return nil
	}
	caser := cases.Title(language.English)
	var refMarkers []*pb.RefMarker
	doc.Find("ref").Each(func(i int, s *goquery.Selection) {
		if refType, exists := s.Attr("type"); exists && (refType == "table" || refType == "figure") {
			value := s.Text()
			bbox := parseBBoxFromSelection(s, pages)
			refIdx := caser.String(fmt.Sprintf("%s %s", refType, value))
			marker := &pb.RefMarker{
				RefIdx:     refIdx,
				RefContent: value,
				Bbox:       bbox,
			}
			refMarkers = append(refMarkers, marker)
		}
	})

	return refMarkers
}

// parseBBoxFromSelection 从 goquery.Selection 中解析边界框信息
func parseBBoxFromSelection(elem *goquery.Selection, pages []*pb.PageInfo) *pb.BBox {
	if elem == nil {
		return nil
	}

	// 在GROBID中，边界框信息存储在coords属性中
	coordsValue, exists := elem.Attr("coords")
	if !exists {
		return nil
	}

	// coords属性格式为：page,x,y,width,height;page,x,y,width,height;...
	// 我们只取第一个边界框（如果有多个的话）

	// 如果有多个边界框（用分号分隔），只取第一个
	parts := strings.Split(coordsValue, ";")
	if len(parts) == 0 {
		return nil
	}

	// 解析第一个边界框的5个属性
	attributes := strings.Split(parts[0], ",")
	if len(attributes) < 5 {
		return nil // 格式不正确
	}

	// 解析各个属性
	pageNum, err := strconv.Atoi(attributes[0])
	if err != nil {
		return nil
	}

	x, err := strconv.ParseFloat(attributes[1], 64)
	if err != nil {
		return nil
	}

	y, err := strconv.ParseFloat(attributes[2], 64)
	if err != nil {
		return nil
	}

	width, err := strconv.ParseFloat(attributes[3], 64)
	if err != nil {
		return nil
	}

	height, err := strconv.ParseFloat(attributes[4], 64)
	if err != nil {
		return nil
	}

	// 获取页面尺寸
	pageWidth, pageHeight := getPageDimensions(pageNum, pages)

	// 创建并返回边界框对象
	// 注意：GROBID的坐标系是左上角为原点，x向右，y向下
	return &pb.BBox{
		X0:           x,
		Y0:           y,
		X1:           x + width,
		Y1:           y + height,
		OriginHeight: pageHeight,
		OriginWidth:  pageWidth,
		PageNumber:   int32(pageNum),
	}
}
