package mineru

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	pb "github.com/yb2020/odoc/proto/gen/go/parsed"
	"github.com/yb2020/odoc/services/parse/constant"
	"github.com/yb2020/odoc/services/parse/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// 处理minerU的images文件 extractedDir: 解压后的目录 imagesDir: images文件夹
func HandleImageFiles(extractedDir, imagesDir string) (map[string]*pb.ImageRecord, error) {

	imagesFiles, err := os.ReadDir(imagesDir)
	if err != nil {
		return nil, errors.BizWrap("read images directory failed", err)
	}
	// 使用map存储图片记录，key为文件名，value为图片上传记录
	var imageRecords = make(map[string]*pb.ImageRecord)
	for _, file := range imagesFiles {
		// 获取图片文件名
		filename := file.Name()
		//读取图片文件
		imageFile, err := os.Open(filepath.Join(imagesDir, filename))
		if err != nil {
			return nil, errors.BizWrap("read image file failed", err)
		}
		defer imageFile.Close()
		// 保存图片记录到map中，使用文件名作为key
		imageRecord := &pb.ImageRecord{
			Id: idgen.GenerateUUID(),
		}
		imageRecords[filename] = imageRecord
	}
	return imageRecords, nil
}

// 解析mineru的json文件   返回middle 和 content_list 的json文件字节数组
func GetMineruContentListAndMiddleJsonBytes(extractedDir string) ([]byte, []byte, error) {

	// 列出解压目录中的文件
	files, err := os.ReadDir(extractedDir)
	if err != nil {
		return nil, nil, errors.BizWrap("read extracted directory failed", err)
	}

	var contentListJsonBytes []byte
	var middleJsonBytes []byte
	for _, file := range files {
		// 文件名称
		filename := file.Name()
		// 根据文件的扩展名识别文件类型
		if !strings.HasSuffix(filename, constant.MineruJsonSuffix) {
			continue
		}
		// 如果是json文件则继续判断json文件的类型
		if strings.Contains(filename, constant.MineruJsonContentList) {
			// 读取文件内容
			contentListJsonBytes, err = os.ReadFile(filepath.Join(extractedDir, filename))
			if err != nil {
				return nil, nil, errors.BizWrap("read content_list.json file failed", err)
			}
			continue
		}
		if strings.Contains(filename, constant.MineruJsonMiddle) {
			// 读取文件内容
			middleJsonBytes, err = os.ReadFile(filepath.Join(extractedDir, filename))
			if err != nil {
				return nil, nil, errors.BizWrap("read middle.json file failed", err)
			}
			continue
		}
		break
	}
	return contentListJsonBytes, middleJsonBytes, nil
}

// 处理middle.json文件
func HandleMiddleJsonByte(middleJsonBytes []byte, imageRecords *map[string]*pb.ImageRecord, formulaRecords *map[string]*pb.FormulaRecord, contentTitles []*pb.ContentTitle) ([]*pb.Paragraph, *pb.DocumentMetadata, []*pb.PageBlockData, error) {

	var mineruPdfInfo pb.MineruPdfInfo
	// 解析middle.json
	if err := json.Unmarshal(middleJsonBytes, &mineruPdfInfo); err != nil {
		return nil, nil, nil, errors.BizWrap("parse middle.json failed", err)
	}
	pages := mineruPdfInfo.PdfInfo

	var paragraphs []*pb.Paragraph
	//
	var catalogueItems []*pb.CatalogueItem
	// 最新的title
	var latestTitle string
	// 段落的index
	var paragraphIndex int = 1

	//解析page属性
	var docPages []*pb.PageInfo
	//全文翻译的页面块
	var allPageBlocks []*pb.PageBlockData
	//解析全文和图表 抽取段落信息，图表信息，目录信息
	for _, page := range pages {
		if page.PageSize == nil {
			continue
		}

		pagePageSize := page.PageSize
		// 转换成 float64
		middlePageSizeWidth := pagePageSize[0]
		middlePageSizeHeight := pagePageSize[1]
		blocks := page.ParaBlocks
		if len(blocks) == 0 {
			continue
		}
		//转换成PageInfo
		docPageInfo := &pb.PageInfo{
			PageNumber: page.PageIdx,
			Width:      middlePageSizeWidth,
			Height:     middlePageSizeHeight,
		}
		docPages = append(docPages, docPageInfo)
		//处理块，专门用于全文翻译
		pageBlocks := HandleParaBlocksForFullText(middlePageSizeWidth, middlePageSizeHeight, blocks, imageRecords, formulaRecords, int(page.PageIdx), &latestTitle, &paragraphIndex, contentTitles)
		if len(pageBlocks) > 0 {
			//这里需要将数据进行存储并上传到oss
			// 这里的pageidx需要+1，因为pageidx从0开始，但是pageblockdata的pageindex从1开始
			pageBlockData := &pb.PageBlockData{
				PageIndex:  int32(page.PageIdx) + 1,
				PageBlocks: pageBlocks,
			}
			allPageBlocks = append(allPageBlocks, pageBlockData)
		}
		//处理段落
		paragraphBlocks := HandleParaBlocks(middlePageSizeWidth, middlePageSizeHeight, blocks, imageRecords, formulaRecords, int(page.PageIdx), &latestTitle, &paragraphIndex)
		if len(paragraphBlocks) == 0 {
			continue
		}
		paragraphs = append(paragraphs, paragraphBlocks...)
		//处理目录
		catalogueItem := ParsePageCatalogueItems(page)
		catalogueItems = append(catalogueItems, catalogueItem...)

	}

	//过滤掉paragraphs中的paragraph为nil的
	order := 1
	var filteredParagraphs []*pb.Paragraph
	for _, paragraph := range paragraphs {
		if paragraph != nil {
			paragraph.Order = int32(order)
			//去除标题前后空格
			paragraph.SectionTitle = strings.TrimSpace(paragraph.SectionTitle)
			//去除section_id前后空格
			paragraph.SectionId = strings.TrimSpace(paragraph.SectionId)
			filteredParagraphs = append(filteredParagraphs, paragraph)
			order++
		}
	}
	//从段落中抽出标题和作者
	authors := HandleParagraphsToAuthors(filteredParagraphs)
	//从目录中抽出标题  因为这里的更贴合实际
	title := HandleCatalogueItemsToTitle(catalogueItems)

	//从段落中抽取摘要和致谢
	abstract, acknowledgments := HandleAbstractAndAcknowledgments(filteredParagraphs)

	//目录整理 使其是一个层级的目录结构 这个方法实际测试发现精度不够，姑且保留，但是不使用
	// catalogue := util.HandleCatalogueItems(catalogueItems)
	//目录整理 使其是一个层级的目录结构
	catalogue := handleCatalogueItems(catalogueItems)
	//按照正常思维无法解析目录的时候
	if len(catalogue) == 0 {
		catalogue = handleCatalogueItemsNoRule(catalogueItems)
	}
	//从段落中抽取参考文献
	// references := HandleReferences(paragraphs)
	metadata := &pb.DocumentMetadata{
		Abstract:       abstract,
		Acknowledgment: acknowledgments,
		Catalogue:      catalogue,
		Title:          title,
		Authors:        authors,
		Pages:          docPages,
	}

	return filteredParagraphs, metadata, allPageBlocks, nil
}

// 处理content_list.json文件
func HandleContentListJsonByte(contentListJsonBytes []byte, imageRecords *map[string]*pb.ImageRecord, formulaRecords *map[string]*pb.FormulaRecord) (*pb.DocumentMetadata, []*pb.ContentTitle, error) {

	var figureTables []*pb.FigureTable
	var formulas []*pb.Formula
	//标题
	var contentTitles []*pb.ContentTitle
	//content_list.json文件的json数据结构数组
	var contentList []pb.MineruContentList
	// 解析content_list.json
	if err := json.Unmarshal(contentListJsonBytes, &contentList); err != nil {
		return nil, nil, errors.BizWrap("parse content_list.json failed", err)
	}
	caser := cases.Title(language.English)
	// 遍历contentList
	for i := range contentList {
		content := &contentList[i]
		switch content.Type {
		case constant.ContentListJsonImage:
			if image := handleContentListImage(content, imageRecords, caser); image != nil {
				figureTables = append(figureTables, image)
			}
		case constant.ContentListJsonTable:
			if table := handleContentListTable(content, imageRecords, caser); table != nil {
				figureTables = append(figureTables, table)
			}
		case constant.ContentListJsonEquation:
			if formula := handleContentListEquation(content, formulaRecords); formula != nil {
				formulas = append(formulas, formula)
			}
		default:
			if content.TextLevel > 0 {
				// 这里的标题级别大于0  则代表在当前文件中。它是被识别成了一个标题的
				if title := handleContentListTitle(content, caser); title != nil {
					contentTitles = append(contentTitles, title)
				}
			}
		}

	}

	return &pb.DocumentMetadata{
		FiguresAndTables: figureTables,
		Formulas:         formulas,
	}, contentTitles, nil
}

// 处理content_list中的图片类型
func handleContentListImage(content *pb.MineruContentList, imageRecords *map[string]*pb.ImageRecord, caser cases.Caser) *pb.FigureTable {
	//过滤掉imgpath为空的
	if content.ImgPath == "" {
		return nil
	}
	image := &pb.FigureTable{
		Type: constant.PdfOssTypeFigure,
	}
	if len(content.ImgCaption) > 0 {
		image.RefContent = caser.String(content.ImgCaption[0])
	}
	//图片需要进行对应，从ImgPath中提取文件名
	imageFileName := filepath.Base(content.ImgPath)
	// 查找对应的图片记录
	if imageRecord, exists := (*imageRecords)[imageFileName]; exists {
		image.Id = imageRecord.Id
		// image.Bbox = imageRecord.Bbox //这个bbox在这里的时候还没有被赋值
		image.SectionTitle = imageRecord.SectionTitle
		image.SectionId = imageRecord.SectionId
		image.RefBbox = imageRecord.RefBbox
	}
	return image
}

// 处理content_list中的表格类型
func handleContentListTable(content *pb.MineruContentList, imageRecords *map[string]*pb.ImageRecord, caser cases.Caser) *pb.FigureTable {
	//过滤掉imgpath为空的
	if content.ImgPath == "" {
		return nil
	}
	table := &pb.FigureTable{
		Id:   idgen.GenerateUUID(),
		Type: constant.PdfOssTypeTable,
	}
	if len(content.TableCaption) > 0 {
		table.RefContent = caser.String(content.TableCaption[0])
	}
	//table需要进行对应，如果有ImgPath字段
	if content.ImgPath != "" {
		imageFileName := filepath.Base(content.ImgPath)
		// 查找对应的图片记录
		if imageRecord, exists := (*imageRecords)[imageFileName]; exists {
			table.Id = imageRecord.Id
			// table.Bbox = imageRecord.Bbox //这个bbox在这里的时候还没有被赋值
			table.SectionTitle = imageRecord.SectionTitle
			table.SectionId = imageRecord.SectionId
			table.RefBbox = imageRecord.RefBbox
		}
	}
	return table
}

// 处理content_list中的公式类型
func handleContentListEquation(content *pb.MineruContentList, formulaRecords *map[string]*pb.FormulaRecord) *pb.Formula {
	formula := &pb.Formula{
		Id: idgen.GenerateUUID(),
	}
	// 公式内容为$$\n*\n$$
	formulaText := content.Text
	if len(formulaText) > 0 {
		//去除公式的前缀和后缀
		formulaText = strings.TrimPrefix(formulaText, "$$\n")
		formulaText = strings.TrimSuffix(formulaText, "\n$$")
		formula.RefContent = formulaText
		//查找对应的公式记录
		if formulaRecord, exists := (*formulaRecords)[formulaText]; exists {
			// formula.Bbox = formulaRecord.Bbox //这个bbox在这里的时候还没有被赋值
			formula.SectionTitle = formulaRecord.SectionTitle
			formula.SectionId = formulaRecord.SectionId
		}
		return formula
	}
	return nil
}

func handleContentListTitle(content *pb.MineruContentList, caser cases.Caser) *pb.ContentTitle {
	// 过滤掉空文本
	if content.Text == "" {
		return nil
	}
	title := &pb.ContentTitle{
		Level:      content.TextLevel,
		Text:       strings.TrimSpace(caser.String(content.Text)),
		PageNumber: content.PageIdx + 1,
	}
	return title
}

// 解析页面中的目录信息
func ParsePageCatalogueItems(page *pb.MineruPage) []*pb.CatalogueItem {
	//转换成 float64
	middlePageSizeWidth := page.PageSize[0]
	middlePageSizeHeight := page.PageSize[1]
	pageNum := page.PageIdx + 1
	var catalogueItems []*pb.CatalogueItem
	for _, block := range page.ParaBlocks {
		if block.Type == constant.MiddleBlockTypeTitle {
			//
			var title string
			lines := block.Lines
			if len(lines) == 0 {
				continue
			}
			for _, line := range lines {
				//如果存在多条  则全部取出 并且用空格连接
				for _, span := range line.Spans {
					//先去除span.content的前后的空格
					spanContent := strings.TrimSpace(span.Content)
					title += spanContent + " "
				}
			}
			//去除title前后的空格
			title = strings.TrimSpace(title)
			// 更严格的正则：仅当数字+空格+数字时才删除前面的部分  这个主要是为了处理章节前面带页码或者行号的情况，如"29 1 Introduction"
			title = strings.TrimSpace(regexp.MustCompile(`^\d+\s+`).ReplaceAllString(title, ""))
			//bbox
			bbox := pb.BBox{
				X0:           float64(block.Bbox[0]),
				Y0:           float64(block.Bbox[1]),
				X1:           float64(block.Bbox[2]),
				Y1:           float64(block.Bbox[3]),
				OriginHeight: float64(middlePageSizeHeight),
				OriginWidth:  float64(middlePageSizeWidth),
				PageNumber:   int32(pageNum),
			}
			catalogueItem := &pb.CatalogueItem{
				Bbox:  &bbox,
				Title: title,
			}
			catalogueItems = append(catalogueItems, catalogueItem)
		}
	}
	return catalogueItems
}

// 抽取目录， 不带规则
func handleCatalogueItemsNoRule(catalogueItems []*pb.CatalogueItem) []*pb.CatalogueItem {
	if len(catalogueItems) == 0 {
		return make([]*pb.CatalogueItem, 0)
	}
	// 1. 设置目录项的顺序
	var order int = 1
	//过滤第一个目录，因为第一个目录大概率是标题
	catalogueItems = catalogueItems[1:]
	for _, item := range catalogueItems {
		item.Order = int32(order)
		item.Level = "1"
		item.TitleOrder = item.Title
		item.Child = []*pb.CatalogueItem{}
		order++
	}
	return catalogueItems
}

// 抽取目录 带规则
func handleCatalogueItems(catalogueItems []*pb.CatalogueItem) []*pb.CatalogueItem {
	if len(catalogueItems) == 0 {
		return make([]*pb.CatalogueItem, 0)
	}
	// 1. 过滤掉不带数字标题的目录对象
	var filteredItems []*pb.CatalogueItem
	// 2. 设置目录项的顺序
	var order int = 1
	// 匹配整个标题的正则
	titleRe := regexp.MustCompile(`^\s*\d+(?:\.\d+)*\.?\s*[^\s\d].*`)

	for _, item := range catalogueItems {
		itemTitle := item.Title
		item.Title = itemTitle
		if titleRe.MatchString(itemTitle) {
			// 拆分格式化标题
			titleOrder, _ := util.SplitNormalizedTitle(itemTitle)
			level := len(strings.Split(titleOrder, "."))
			item.Level = strconv.Itoa(level)
			item.TitleOrder = titleOrder
			item.Child = []*pb.CatalogueItem{}
			item.Order = int32(order)
			order++
			filteredItems = append(filteredItems, item)
		}
	}

	// 2. 构建目录树
	// 创建根节点映射，用于快速查找
	rootItems := make([]*pb.CatalogueItem, 0)
	itemMap := make(map[string]*pb.CatalogueItem)

	// 按照标题序号排序
	sort.Slice(filteredItems, func(i, j int) bool {
		return filteredItems[i].TitleOrder < filteredItems[j].TitleOrder
	})

	// 构建目录树
	for _, item := range filteredItems {
		titleOrder := item.TitleOrder

		// 查找父节点
		lastDotIndex := strings.LastIndex(titleOrder, ".")
		if lastDotIndex == -1 {
			// 这是一级目录
			rootItems = append(rootItems, item)
			itemMap[titleOrder] = item
		} else {
			// 这是子目录
			parentOrder := titleOrder[:lastDotIndex]
			if parent, exists := itemMap[parentOrder]; exists {
				// 将当前目录添加为父目录的子目录
				parent.Child = append(parent.Child, item)
			} else {
				// 找不到父节点，作为根节点处理
				rootItems = append(rootItems, item)
			}
			itemMap[titleOrder] = item
		}
	}
	return rootItems
}

// 全文翻译专用。处理页面块的数据
func HandleParaBlocksForFullText(
	middlePageSizeWidth float64,
	middlePageSizeHeight float64,
	blocks []*pb.Block,
	imageRecords *map[string]*pb.ImageRecord,
	formulaRecords *map[string]*pb.FormulaRecord,
	pageIdx int,
	title *string,
	paragraphIndex *int,
	contentTitles []*pb.ContentTitle,
) []*pb.PageBlock {

	//定义页面块对象
	var PageBlocks []*pb.PageBlock
	//遍历blocks
	var index int = 1
	for _, block := range blocks {
		//定义页面块
		pageBlock := &pb.PageBlock{
			LinesDeleted: false,
		}
		switch block.Type {
		case constant.MiddleBlockTypeTitle:
			pageBlock.Type = pb.PageBlockType_PAGE_BLOCK_TITLE
		case constant.ContentListJsonText:
			pageBlock.Type = pb.PageBlockType_PAGE_BLOCK_PLAIN_TEXT
		case constant.ContentListJsonImage:
			pageBlock.Type = pb.PageBlockType_PAGE_BLOCK_FIGURE
		case constant.ContentListJsonTable:
			pageBlock.Type = pb.PageBlockType_PAGE_BLOCK_TABLE
		case constant.MiddleBlockInterlineEquation:
			pageBlock.Type = pb.PageBlockType_PAGE_BLOCK_ISOLATE_FORMULA
		default:
			pageBlock.Type = pb.PageBlockType_PAGE_BLOCK_PLAIN_TEXT
		}
		if pageBlock.Type == pb.PageBlockType_PAGE_BLOCK_FIGURE ||
			pageBlock.Type == pb.PageBlockType_PAGE_BLOCK_TABLE {
			//图表类型的单独处理
			imageAndTablesBlocks := HandleImageAndTablesForFullText(middlePageSizeWidth, middlePageSizeHeight, block.Blocks, pageIdx, pageBlock.Type, index)
			PageBlocks = append(PageBlocks, imageAndTablesBlocks...)
			index += len(imageAndTablesBlocks)
			continue
		}
		//块的bbox坐标
		bbox := pb.BBox{
			X0:           float64(block.Bbox[0]),
			Y0:           float64(block.Bbox[1]),
			X1:           float64(block.Bbox[2]),
			Y1:           float64(block.Bbox[3]),
			OriginHeight: float64(middlePageSizeHeight),
			OriginWidth:  float64(middlePageSizeWidth),
			PageNumber:   int32(pageIdx + 1),
		}
		pageBlock.Bbox = &bbox
		//设置块的index
		pageBlock.Index = int32(index)
		index++
		//是否删除行
		pageBlock.LinesDeleted = block.LinesDeleted
		//内容 （只需要纯文本数据就好）
		pageBlock.Texts = HandleTextBlockForFullText(middlePageSizeWidth, middlePageSizeHeight, pageIdx, block, pageBlock, contentTitles)
		PageBlocks = append(PageBlocks, pageBlock)
	}
	return PageBlocks
}

func HandleParaBlocks(
	middlePageSizeWidth float64,
	middlePageSizeHeight float64,
	blocks []*pb.Block,
	imageRecords *map[string]*pb.ImageRecord,
	formulaRecords *map[string]*pb.FormulaRecord,
	pageIdx int,
	title *string,
	paragraphIndex *int,
) []*pb.Paragraph {
	//定义段落对象
	var paragraphs []*pb.Paragraph

	//遍历blocks
	for _, block := range blocks {
		switch block.Type {
		case constant.MiddleBlockTypeTitle:
			*title = HandleTitleBlock(middlePageSizeWidth, middlePageSizeHeight, pageIdx, block)
		case constant.ContentListJsonText:
			paragraph := HandleTextBlock(middlePageSizeWidth, middlePageSizeHeight, pageIdx, block, *title, paragraphIndex)
			paragraphs = append(paragraphs, paragraph)
			*paragraphIndex++
		case constant.ContentListJsonImage:
			paragraph := HandleImageAndTables(middlePageSizeWidth, middlePageSizeHeight, block.Blocks, pageIdx, imageRecords, constant.ContentListJsonImage, *title)
			paragraphs = append(paragraphs, paragraph)
			*paragraphIndex++
		case constant.ContentListJsonTable:
			paragraph := HandleImageAndTables(middlePageSizeWidth, middlePageSizeHeight, block.Blocks, pageIdx, imageRecords, constant.ContentListJsonTable, *title)
			paragraphs = append(paragraphs, paragraph)
			*paragraphIndex++
		case constant.MiddleBlockInterlineEquation:
			paragraph := HandleEquation(middlePageSizeWidth, middlePageSizeHeight, block, pageIdx, formulaRecords, *title)
			paragraphs = append(paragraphs, paragraph)
			*paragraphIndex++
		}
	}
	return paragraphs
}

func HandleTitleBlock(
	middlePageSizeWidth float64,
	middlePageSizeHeight float64,
	pageIdx int,
	block *pb.Block,
) string {
	var title string

	if len(block.Lines) == 0 {
		return ""
	}
	for _, line := range block.Lines {
		//如果存在多条  则全部取出 并且用空格连接
		for _, span := range line.Spans {
			//先去除span.content的前后的空格
			spanContent := strings.TrimSpace(span.Content)
			title += spanContent + " "
		}
	}
	//更严格的正则：仅当数字+空格+数字时才删除前面的部分  这个主要是为了处理章节前面带页码或者行号的情况，如"29 1 Introduction"
	title = strings.TrimSpace(regexp.MustCompile(`^\d+\s+`).ReplaceAllString(title, ""))
	title = util.NormalizeTitle(title)
	return title
}

// 处理文本块，专门给全文翻译用，因为这里和embbebbing的处理方式不同
func HandleTextBlockForFullText(
	middlePageSizeWidth float64,
	middlePageSizeHeight float64,
	pageIdx int,
	block *pb.Block,
	pageBlock *pb.PageBlock,
	contentTitles []*pb.ContentTitle,
) []*pb.BlockText {
	//判断lines的长度
	if len(block.Lines) == 0 {
		return nil
	}
	//将每行的span内容拼接成paragraph
	var texts []*pb.BlockText
	for _, line := range block.Lines {
		text := &pb.BlockText{
			CrossPage: false,
		}
		bbox := pb.BBox{
			X0:           float64(line.Bbox[0]),
			Y0:           float64(line.Bbox[1]),
			X1:           float64(line.Bbox[2]),
			Y1:           float64(line.Bbox[3]),
			OriginHeight: float64(middlePageSizeHeight),
			OriginWidth:  float64(middlePageSizeWidth),
			PageNumber:   int32(pageIdx + 1),
		}
		text.Bbox = &bbox
		//目前还没有发现多个span的情况，不知道多个span的情况是什么样的，所以这里默认只有一个。
		//当块的类型是标题时，这里是存在多个span的，
		if len(line.Spans) == 1 {
			text.Text = strings.TrimSpace(line.Spans[0].Content)
			//当前行是否跨页
			text.CrossPage = line.Spans[0].CrossPage
			// if strings.Contains(text.Text, "References") {
			// 	log.Printf("Found References block: %s", text.Text)
			// }
		}
		if len(line.Spans) > 1 {
			for _, span := range line.Spans {
				text.Text += " " + span.Content
				//当前行是否跨页
				text.CrossPage = span.CrossPage
			}
			text.Text = strings.TrimSpace(text.Text)

		}
		//这里如果text.Text的内容和contentTitles数组中的任意一项的value的开头一样，则需要设置当前级别
		if text.Text != "" && pageBlock.Type == pb.PageBlockType_PAGE_BLOCK_TITLE {
			for _, contentTitle := range contentTitles {
				if strings.HasPrefix(normalizeForCompare(text.Text), normalizeForCompare(contentTitle.Text)) {
					pageBlock.Level = contentTitle.Level
					break
				}
			}
		}

		texts = append(texts, text)
	}
	return texts
}

func HandleImageAndTablesForFullText(
	middlePageSizeWidth float64,
	middlePageSizeHeight float64,
	imageAndTables []*pb.Block,
	pageIdx int,
	blockType pb.PageBlockType,
	index int,
) []*pb.PageBlock {
	if len(imageAndTables) == 0 {
		return nil
	}
	//根据type的定义查找是否存在描述
	var pageBlocks []*pb.PageBlock
	for _, block := range imageAndTables {

		//bbox
		bbox := pb.BBox{
			X0:           float64(block.Bbox[0]),
			Y0:           float64(block.Bbox[1]),
			X1:           float64(block.Bbox[2]),
			Y1:           float64(block.Bbox[3]),
			OriginHeight: float64(middlePageSizeHeight),
			OriginWidth:  float64(middlePageSizeWidth),
			PageNumber:   int32(pageIdx + 1),
		}
		//pageBlock
		pageBlock := &pb.PageBlock{
			Bbox:         &bbox,
			Index:        int32(index),
			LinesDeleted: false,
		}
		if strings.Contains(block.Type, constant.MiddleBlockImageAndTableBody) {
			//
			pageBlock.Type = blockType
			index++
			pageBlocks = append(pageBlocks, pageBlock)
		}
		// 这里判断type是否包含了caption  如果包含了，设置当前为文本块，翻译的时候需要处理
		if strings.Contains(block.Type, constant.MiddleBlockCaption) {
			pageBlock.Type = pb.PageBlockType_PAGE_BLOCK_PLAIN_TEXT
			pageBlock.Texts = HandleTextBlockForFullText(middlePageSizeWidth, middlePageSizeHeight, pageIdx, block, pageBlock, nil)
			index++
			pageBlocks = append(pageBlocks, pageBlock)
		}
	}
	return pageBlocks
}

func HandleTextBlock(
	middlePageSizeWidth float64,
	middlePageSizeHeight float64,
	pageIdx int,
	block *pb.Block,
	title string,
	index *int,
) *pb.Paragraph {

	//判断lines的长度
	if len(block.Lines) == 0 {
		return nil
	}
	//解析段落的页码大小
	lenPageSize := len(block.PageSize)
	//如果当前块中存在页码大小信息则使用当前的，如果不存在则使用外部的
	if lenPageSize > 0 {
		middlePageSizeWidth = block.PageSize[0]
		middlePageSizeHeight = block.PageSize[1]
	}
	//解析bbox
	bbox := pb.BBox{
		X0:           float64(block.Bbox[0]),
		Y0:           float64(block.Bbox[1]),
		X1:           float64(block.Bbox[2]),
		Y1:           float64(block.Bbox[3]),
		OriginHeight: float64(middlePageSizeHeight),
		OriginWidth:  float64(middlePageSizeWidth),
		PageNumber:   int32(pageIdx + 1),
	}
	var paragraph pb.Paragraph
	var contentBuilder strings.Builder
	//将每行的span内容拼接成paragraph
	for _, line := range block.Lines {
		for _, span := range line.Spans {
			contentBuilder.WriteString(span.Content)
		}
	}
	//段落内容
	content := contentBuilder.String()
	//判断段落的语言
	language := util.DetectLanguageFallback(content)
	/** 句子切割 这里根据中英文的标点符号进行切割
	TODO: 这里拆分句子的时候用的是段落的坐标，不够准确，特别是段落存在分页的情况
	正确的做法应该把每行的坐标传入，然后根据每行的坐标来拆分句子，保证句子的坐标定位到行
	实现：可以通过段落的每行坐标计算出段落的初始坐标，然后和每行所增加的坐标和页码，后续处理的时候的就可以根据这些信息计算出坐标
	*/
	sentences := HandleParagraphToSentences(content, language, &bbox)
	// 解析setion_title 和section_id
	paragraph.SectionTitle = title
	paragraph.SectionId = title

	//text
	paragraph.Text = &pb.Text{
		Text:      content,
		Sentences: sentences,
		Bbox:      &bbox,
	}

	paragraph.Type = pb.ParagraphType_TEXT
	//references 信息暂时不处理
	return &paragraph
}

// 处理image和table
func HandleImageAndTables(
	middlePageSizeWidth float64,
	middlePageSizeHeight float64,
	imageAndTables []*pb.Block,
	pageIdx int,
	imageRecords *map[string]*pb.ImageRecord,
	handleType string,
	title string,
) *pb.Paragraph {
	var block *pb.Block
	var line *pb.Line
	blockType := handleType + constant.MiddleBlockImageAndTableBody

	//图片坐标
	var bbox *pb.BBox
	//标题坐标
	var captionBbox *pb.BBox

	var caption *pb.Block
	captionType := handleType + constant.MiddleBlockCaption

	for _, it := range imageAndTables {
		switch it.Type {
		case blockType:
			block = it
		case captionType:
			caption = it
		}
	}
	if block == nil {
		return nil
	}
	if block.Bbox != nil {
		bbox = &pb.BBox{
			X0:           float64(block.Bbox[0]),
			Y0:           float64(block.Bbox[1]),
			X1:           float64(block.Bbox[2]),
			Y1:           float64(block.Bbox[3]),
			OriginHeight: float64(middlePageSizeHeight),
			OriginWidth:  float64(middlePageSizeWidth),
			PageNumber:   int32(pageIdx + 1),
		}
	}
	paragraph := &pb.Paragraph{}
	if caption != nil {
		captionBbox = &pb.BBox{
			X0:           float64(caption.Bbox[0]),
			Y0:           float64(caption.Bbox[1]),
			X1:           float64(caption.Bbox[2]),
			Y1:           float64(caption.Bbox[3]),
			OriginHeight: float64(middlePageSizeHeight),
			OriginWidth:  float64(middlePageSizeWidth),
			PageNumber:   int32(pageIdx + 1),
		}
		var spanText string
		for _, line := range caption.Lines {
			for _, span := range line.Spans {
				spanText += span.Content
			}
		}
		figureTable := &pb.FigureTable{
			Id:           idgen.GenerateUUID(),
			SectionTitle: title,
			SectionId:    title,
			RefContent:   spanText,
			RefBbox:      captionBbox,
			Bbox:         bbox,
		}
		paragraph = &pb.Paragraph{
			FigureTable:  figureTable,
			SectionTitle: title,
			SectionId:    title,
		}
		switch handleType {
		case constant.ContentListJsonImage:
			paragraph.Type = pb.ParagraphType_IMAGE
		case constant.ContentListJsonTable:
			paragraph.Type = pb.ParagraphType_TABLE
		}
	}
	if len(block.Lines) == 0 {
		return nil
	}
	line = block.Lines[0]
	for _, span := range line.Spans {
		if span.Type != handleType {
			continue
		}
		imageBbox := pb.BBox{
			X0:           float64(span.Bbox[0]),
			Y0:           float64(span.Bbox[1]),
			X1:           float64(span.Bbox[2]),
			Y1:           float64(span.Bbox[3]),
			OriginHeight: float64(middlePageSizeHeight),
			OriginWidth:  float64(middlePageSizeWidth),
			PageNumber:   int32(pageIdx + 1),
		}
		imageName := span.ImagePath
		// 查找对应的图片记录
		if imageRecord, exists := (*imageRecords)[imageName]; exists {
			//插入坐标
			imageRecord.Bbox = &imageBbox
			//插入section_title
			imageRecord.SectionTitle = title
			//插入section_id
			imageRecord.SectionId = title
			//插入ref_bbox（标题坐标）
			imageRecord.RefBbox = captionBbox
		}
	}
	return paragraph
}

// 处理equation
func HandleEquation(
	middlePageSizeWidth float64,
	middlePageSizeHeight float64,
	block *pb.Block,
	pageIdx int,
	formulaRecords *map[string]*pb.FormulaRecord,
	title string,
) *pb.Paragraph {
	//判断lines的长度
	if len(block.Lines) == 0 {
		return nil
	}
	var ParagraphFormula *pb.Formula
	//定义bbox
	bbox := pb.BBox{
		X0:           float64(block.Bbox[0]),
		Y0:           float64(block.Bbox[1]),
		X1:           float64(block.Bbox[2]),
		Y1:           float64(block.Bbox[3]),
		OriginHeight: float64(middlePageSizeHeight),
		OriginWidth:  float64(middlePageSizeWidth),
		PageNumber:   int32(pageIdx + 1),
	}
	// 遍历lines
	for _, line := range block.Lines {
		for _, span := range line.Spans {
			if span.Type != constant.MiddleBlockInterlineEquation {
				continue
			}
			// 新建公式对象
			formula := &pb.FormulaRecord{
				PageIdx:      int32(pageIdx + 1),
				Content:      span.Content,
				Bbox:         &bbox,
				SectionTitle: title,
				SectionId:    title,
			}
			ParagraphFormula = &pb.Formula{
				Id:           idgen.GenerateUUID(),
				Bbox:         &bbox,
				SectionTitle: title,
				SectionId:    title,
				RefContent:   span.Content,
			}
			//插入公式记录
			(*formulaRecords)[span.Content] = formula
		}
	}
	paragraph := &pb.Paragraph{
		Type:         pb.ParagraphType_FORMULA,
		Formula:      ParagraphFormula,
		SectionTitle: title,
		SectionId:    title,
	}
	return paragraph
}

func HandleParagraphToSentences(content string, language string, bbox *pb.BBox) []*pb.Sentence {
	var sentences []*pb.Sentence
	var index int = 1
	//根据中英文的句子符号进行切割  TODO : 这里解析出来的句子少了bbox信息
	splitString := ""
	switch language {
	case constant.LanguageEnUS:
		splitString = "."
	case constant.LanguageZhCN:
		splitString = "。"
	default:
		splitString = "."
	}
	sentenceStrs := strings.Split(content, splitString)
	for _, sentenceStr := range sentenceStrs {
		if sentenceStr == "" {
			continue
		}
		sentences = append(sentences, &pb.Sentence{
			Text:  sentenceStr + splitString,
			Index: int32(index),
			Bbox:  bbox,
		})
		index++
	}
	return sentences
}

func HandleAbstractAndAcknowledgments(paragraphs []*pb.Paragraph) (*pb.Abstract, *pb.Acknowledgment) {
	var abstractParagraphs []*pb.Paragraph
	var acknowledgmentsParagraphs []*pb.Paragraph
	for _, paragraph := range paragraphs {
		if paragraph.SectionTitle == "" || paragraph.Text == nil || paragraph.Text.Text == "" {
			continue
		}
		// 忽略大小写
		sectionTitle := strings.ToLower(paragraph.SectionTitle)
		text := strings.ToLower(paragraph.Text.Text)
		if strings.Contains(sectionTitle, "abstract") ||
			strings.Contains(sectionTitle, "摘要") ||
			strings.HasPrefix(text, "abstract") ||
			strings.HasPrefix(text, "摘要") {
			abstractParagraphs = append(abstractParagraphs, paragraph)
		}
		if strings.Contains(sectionTitle, "acknowledgements") ||
			strings.Contains(sectionTitle, "致谢") ||
			strings.HasPrefix(text, "acknowledgements") ||
			strings.HasPrefix(text, "致谢") {
			acknowledgmentsParagraphs = append(acknowledgmentsParagraphs, paragraph)
		}
	}
	//abstractParagraphs, acknowledgmentsParagraphs
	abstract := HandleParagraphsToAbstract(abstractParagraphs)
	acknowledgments := HandleParagraphsToAcknowledgments(acknowledgmentsParagraphs)
	return abstract, acknowledgments
}

func HandleParagraphsToAbstract(paragraphs []*pb.Paragraph) *pb.Abstract {
	if len(paragraphs) == 0 {
		return nil
	}
	abstract := &pb.Abstract{}
	//bbox
	bbox := pb.BBox{
		X0:           float64(paragraphs[0].Text.Bbox.X0),
		Y0:           float64(paragraphs[0].Text.Bbox.Y0),
		X1:           float64(paragraphs[0].Text.Bbox.X1),
		Y1:           float64(paragraphs[0].Text.Bbox.Y1),
		OriginHeight: float64(paragraphs[0].Text.Bbox.OriginHeight),
		OriginWidth:  float64(paragraphs[0].Text.Bbox.OriginWidth),
		PageNumber:   paragraphs[0].Text.Bbox.PageNumber,
	}
	abstract.Bbox = &bbox
	//text内容
	var contentBuilder strings.Builder
	for _, paragraph := range paragraphs {
		contentBuilder.WriteString(paragraph.Text.Text)
	}
	abstract.Text = contentBuilder.String()
	return abstract
}

func HandleParagraphsToAcknowledgments(paragraphs []*pb.Paragraph) *pb.Acknowledgment {
	if len(paragraphs) == 0 {
		return nil
	}
	acknowledgments := &pb.Acknowledgment{}
	//bbox
	bbox := pb.BBox{
		X0:           float64(paragraphs[0].Text.Bbox.X0),
		Y0:           float64(paragraphs[0].Text.Bbox.Y0),
		X1:           float64(paragraphs[0].Text.Bbox.X1),
		Y1:           float64(paragraphs[0].Text.Bbox.Y1),
		OriginHeight: float64(paragraphs[0].Text.Bbox.OriginHeight),
		OriginWidth:  float64(paragraphs[0].Text.Bbox.OriginWidth),
		PageNumber:   paragraphs[0].Text.Bbox.PageNumber,
	}
	acknowledgments.Bbox = &bbox
	//text内容
	var contentBuilder strings.Builder
	for _, paragraph := range paragraphs {
		contentBuilder.WriteString(paragraph.Text.Text)
	}
	acknowledgments.Text = contentBuilder.String()
	return acknowledgments
}

func HandleParagraphsToAuthors(paragraphs []*pb.Paragraph) []*pb.Author {
	if len(paragraphs) == 0 {
		return nil
	}
	bbox := paragraphs[0].Text.Bbox
	//解析作者信息
	authors := extractNames(paragraphs[0].Text.Text, bbox)
	return authors

}
func HandleCatalogueItemsToTitle(catalogueItems []*pb.CatalogueItem) *pb.Title {
	if len(catalogueItems) == 0 {
		return nil
	}
	bbox := catalogueItems[0].Bbox
	title := &pb.Title{
		Bbox: bbox,
		Text: catalogueItems[0].Title,
	}
	return title
}

func extractNames(input string, bbox *pb.BBox) []*pb.Author {
	// 定义正则表达式
	englishName := `([A-Z][a-zA-Z.-]+)\s+([A-Z][a-zA-Z.-]+)`
	chineseName := `([\p{Han}]{1})([\p{Han}]{1,2})`
	pattern := fmt.Sprintf("(%s|%s)", englishName, chineseName)

	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	// 查找所有匹配项
	matches := re.FindAllStringSubmatch(input, -1)

	var result []*pb.Author
	for _, match := range matches {
		var nameInfo pb.Author
		// var surname, givenName string

		if match[1] != "" { // 英文名
			nameInfo.FullName = strings.TrimSpace(match[1])
			nameInfo.GivenName = strings.TrimSpace(match[2])
			nameInfo.Surname = strings.TrimSpace(match[3])
		} else if match[4] != "" { // 中文名
			nameInfo.FullName = match[4]
			// 中文名处理：第一个字符是姓，其余是名
			nameInfo.Surname = string([]rune(match[4])[0])
			nameInfo.GivenName = string([]rune(match[4])[1:])
		}

		// 清理姓名中的特殊字符
		// clean := func(s string) string {
		// 	return regexp.MustCompile(`[^a-zA-Z\p{Han}]`).ReplaceAllString(s, "")
		// }

		// nameInfo.GivenName = clean(givenName)
		// nameInfo.Surname = clean(surname)

		// 确保不为空
		if nameInfo.FullName != "" {
			result = append(result, &nameInfo)
		}
		nameInfo.Bbox = bbox
	}

	return result
}

func normalizeForCompare(s string) string {
	// normalizeForCompare 标准化字符串用于比较：转小写 + 压缩多余空格
	var spaceRegex = regexp.MustCompile(`\s+`)
	return strings.ToLower(spaceRegex.ReplaceAllString(strings.TrimSpace(s), " "))
}
