package helper

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	docpb "github.com/yb2020/odoc-proto/gen/go/doc"
	"github.com/yb2020/odoc/services/doc/model"
	paperModel "github.com/yb2020/odoc/services/paper/model"
)

// FilterUserDocsByClassify 根据分类标签过滤用户文档
func FilterUserDocsByClassify(req *docpb.GetDocListReq, userDocs []model.UserDoc, docClassifyRelations []model.DocClassifyRelation) []model.UserDoc {
	if req == nil || len(req.ClassifyIds) == 0 || len(userDocs) == 0 {
		return userDocs
	}

	// 构建分类 ID 集合
	classifyIdSet := make(map[string]bool)
	for _, classifyId := range req.ClassifyIds {
		classifyIdSet[classifyId] = true
	}

	// 构建文档 ID 到分类 ID 的映射
	docClassifyMap := make(map[string][]string)
	for _, relation := range docClassifyRelations {
		if _, exists := docClassifyMap[relation.DocId]; !exists {
			docClassifyMap[relation.DocId] = make([]string, 0)
		}
		docClassifyMap[relation.DocId] = append(docClassifyMap[relation.DocId], relation.ClassifyId)
	}

	// 过滤文档
	filteredDocs := make([]model.UserDoc, 0)
	for _, doc := range userDocs {
		// 获取文档的分类 ID 列表
		classifyIds, exists := docClassifyMap[doc.Id]
		if !exists {
			continue
		}

		// 检查是否包含请求中的任一分类 ID
		for _, classifyId := range classifyIds {
			if classifyIdSet[classifyId] {
				filteredDocs = append(filteredDocs, doc)
				break
			}
		}
	}

	return filteredDocs
}

// FilterUserDocsByAuthor 根据作者过滤用户文档
func FilterUserDocsByAuthor(req *docpb.GetDocListReq, userDocs []model.UserDoc) []model.UserDoc {
	if req == nil || len(req.AuthorInfos) == 0 || len(userDocs) == 0 {
		return userDocs
	}

	// 转换作者名称为小写，用于不区分大小写的比较
	authorInfosLower := make([]string, 0, len(req.AuthorInfos))
	for _, author := range req.AuthorInfos {
		authorInfosLower = append(authorInfosLower, strings.ToLower(author))
	}

	filteredDocs := make([]model.UserDoc, 0)
	for _, doc := range userDocs {
		// 获取文档的作者信息
		displayAuthor := GetUserDocDisplayAuthor(&doc)
		if displayAuthor == "" {
			continue
		}

		// 转换为小写进行比较
		displayAuthorLower := strings.ToLower(displayAuthor)

		// 检查是否包含请求中的任一作者
		matched := false
		for _, authorInfo := range authorInfosLower {
			if strings.Contains(displayAuthorLower, authorInfo) {
				matched = true
				break
			}
		}

		if matched {
			filteredDocs = append(filteredDocs, doc)
		}
	}

	return filteredDocs
}

// FilterUserDocsByVenue 根据收录情况过滤用户文档
func FilterUserDocsByVenue(req *docpb.GetDocListReq, userDocs []model.UserDoc) []model.UserDoc {
	if req == nil || len(req.VenueInfos) == 0 || len(userDocs) == 0 {
		return userDocs
	}

	// 转换收录情况为小写，用于不区分大小写的比较
	venueInfosLower := make([]string, 0, len(req.VenueInfos))
	for _, venue := range req.VenueInfos {
		venueInfosLower = append(venueInfosLower, strings.ToLower(venue))
	}

	filteredDocs := make([]model.UserDoc, 0)
	for _, doc := range userDocs {
		// 获取文档的收录情况
		displayVenue := GetUserDocDisplayVenue(&doc)
		if displayVenue == "" {
			continue
		}

		// 转换为小写进行比较
		displayVenueLower := strings.ToLower(displayVenue)

		// 检查是否包含请求中的任一收录情况
		matched := false
		for _, venueInfo := range venueInfosLower {
			if strings.Contains(displayVenueLower, venueInfo) {
				matched = true
				break
			}
		}

		if matched {
			filteredDocs = append(filteredDocs, doc)
		}
	}

	return filteredDocs
}

// FilterUserDocsByFolder 根据文件夹过滤用户文档
func FilterUserDocsByFolder(req *docpb.GetDocListReq, userDocs []model.UserDoc, userDocFolderRelations []model.UserDocFolderRelation) []model.UserDoc {
	if req == nil || req.FolderId == nil || *req.FolderId == "0" || len(userDocs) == 0 {
		return userDocs
	}

	folderId := *req.FolderId

	// 构建文档ID到文件夹ID的映射
	docFolderMap := make(map[string][]string)
	for _, relation := range userDocFolderRelations {
		if _, exists := docFolderMap[relation.DocId]; !exists {
			docFolderMap[relation.DocId] = make([]string, 0)
		}
		docFolderMap[relation.DocId] = append(docFolderMap[relation.DocId], relation.FolderId)
	}

	// 过滤文档
	filteredDocs := make([]model.UserDoc, 0)
	for _, doc := range userDocs {
		// 获取文档的文件夹ID列表
		folderIds, exists := docFolderMap[doc.Id]
		if !exists {
			continue
		}

		// 检查是否包含请求中的文件夹ID
		for _, id := range folderIds {
			if id == folderId {
				filteredDocs = append(filteredDocs, doc)
				break
			}
		}
	}

	return filteredDocs
}

// highlightText 高亮匹配的文本部分，用<em>标签包裹
func highlightText(text, searchContent string) string {
	if text == "" || searchContent == "" {
		return text
	}

	lowerText := strings.ToLower(text)
	lowerSearchContent := strings.ToLower(searchContent)

	// 如果不包含搜索内容，直接返回原文本
	if !strings.Contains(lowerText, lowerSearchContent) {
		return text
	}

	result := ""
	index := 0
	for {
		i := strings.Index(strings.ToLower(text[index:]), lowerSearchContent)
		if i == -1 {
			result += text[index:]
			break
		}

		// 添加匹配前的文本
		result += text[index : index+i]

		// 添加带高亮的匹配文本
		matchedText := text[index+i : index+i+len(searchContent)]
		result += "<em>" + matchedText + "</em>"

		// 更新索引
		index += i + len(searchContent)
	}

	return result
}

// FilterUserDocsBySearchContent 根据搜索内容过滤用户文档
func FilterUserDocsBySearchContent(req *docpb.GetDocListReq, userDocs []model.UserDoc, docSearchResultMap map[string]*docpb.DocSearchResult) []model.UserDoc {
	if req == nil || req.SearchContent == nil || *req.SearchContent == "" || len(userDocs) == 0 {
		return userDocs
	}
	searchContent := *req.SearchContent
	lowerSearchContent := strings.ToLower(searchContent)
	filteredDocs := make([]model.UserDoc, 0)
	for _, doc := range userDocs {
		matched := false
		docSearchResult := &docpb.DocSearchResult{}
		// 在文档名称中搜索
		if strings.Contains(strings.ToLower(doc.DocName), lowerSearchContent) {
			matched = true
			highlightedName := highlightText(doc.DocName, searchContent)
			docSearchResult.HitDocName = &highlightedName
		}
		// 在作者中搜索
		// displayAuthor := doc.AuthorDesc
		// if displayAuthor != "" && strings.Contains(strings.ToLower(displayAuthor), lowerSearchContent) {
		// 	matched = true
		// 	highlightedAuthor := highlightText(displayAuthor, searchContent)
		// 	docSearchResult.HitAuthor = &highlightedAuthor
		// }
		// 在收录情况中搜索
		displayVenue := doc.Venue
		if displayVenue != "" && strings.Contains(strings.ToLower(displayVenue), lowerSearchContent) {
			matched = true
			highlightedVenue := highlightText(displayVenue, searchContent)
			docSearchResult.HitVenue = &highlightedVenue
		}
		// 在备注中搜索
		if doc.Remark != "" && strings.Contains(strings.ToLower(doc.Remark), lowerSearchContent) {
			matched = true
			highlightedRemark := highlightText(doc.Remark, searchContent)
			docSearchResult.HitRemark = &highlightedRemark
		}
		//在jcr分区搜索
		if doc.UserEditedJcrPartion != "" && strings.Contains(strings.ToLower(doc.UserEditedJcrPartion), lowerSearchContent) {
			matched = true
			highlightedJcr := highlightText(doc.UserEditedJcrPartion, searchContent)
			docSearchResult.HitJcrVenuePartion = &highlightedJcr
		}
		//在发布日期搜索
		publishDate := doc.PublishDate
		if publishDate != "" && strings.Contains(strings.ToLower(publishDate), lowerSearchContent) {
			matched = true
			highlightedDate := highlightText(publishDate, searchContent)
			docSearchResult.HitPublishDate = &highlightedDate
		}
		if matched {
			docSearchResultMap[doc.Id] = docSearchResult
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs
}

// GetUserDocDisplayAuthor 获取文档显示作者信息
func GetUserDocDisplayAuthor(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 优先使用显示作者
	if doc.DisplayAuthors != "" {
		return doc.DisplayAuthors
	}

	// 其次使用元数据作者
	if doc.MetaAuthors != "" {
		return doc.MetaAuthors
	}

	// 最后使用图表作者
	return doc.GraphAuthors
}

//TODO
// func GetDisplayAuthorByUserDoc(doc *model.UserDoc) docpb.DisplayAuthor {
// 	if doc == nil {
// 		return docpb.DisplayAuthor{}
// 	}
// 	displayAuthor := docpb.DisplayAuthor{
// 		RollbackEnable: false,
// 	}
// 	//查询作者信息
// 	if doc.AuthorDesc != "" {
// 		authorInfos := strings.Split(doc.AuthorDesc, ",")
// 		displayAuthor.OriginAuthors = authorInfos
// 	}
// 	//查询显示作者信息
// 	if doc.DisplayAuthors != "" {
// 		authorInfos := strings.Split(doc.DisplayAuthors, ",")
// 		displayAuthor.Authors = authorInfos
// 	}
// 	// 最后使用图表作者
// 	return displayAuthor
// }

// GetUserDocDisplayAuthors 获取文档的作者信息列表
func GetUserDocDisplayAuthors(doc *model.UserDoc) []string {
	if doc == nil {
		return []string{}
	}
	var authorInfoResults []string
	//查询作者信息
	authInfoStr := doc.AuthorDesc
	//查询作者信息
	if authInfoStr != "" {
		var authorInfos []*docpb.UserDocAuthorInfo
		err := json.Unmarshal([]byte(authInfoStr), &authorInfos)
		if err != nil {
			return []string{}
		}
		for _, authorInfo := range authorInfos {
			authorInfoResults = append(authorInfoResults, authorInfo.Literal)
		}
	}
	return authorInfoResults
}

// GetUserDocDisplayVenue 获取文档显示收录情况
func GetUserDocDisplayVenue(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 优先使用显示发表场所
	if doc.DisplayVenues != "" {
		return doc.DisplayVenues
	}

	// 其次使用元数据发表场所
	if doc.MetaVenues != "" {
		return doc.MetaVenues
	}

	// 最后使用图表发表场所
	return doc.GraphVenues
}

// GetUserDocDisplayVenueInfos 获取文档的发表场所信息列表
func GetUserDocDisplayVenueInfos(doc *model.UserDoc) []string {
	if doc == nil {
		return []string{}
	}

	// 获取发表场所字符串
	venueStr := GetUserDocDisplayVenue(doc)
	if venueStr == "" {
		return []string{}
	}

	// 分割发表场所字符串，通常发表场所之间用逗号分隔
	venues := strings.Split(venueStr, ",")

	// 去除空格并过滤空字符串
	result := make([]string, 0, len(venues))
	for _, venue := range venues {
		venue = strings.TrimSpace(venue)
		if venue != "" {
			result = append(result, venue)
		}
	}

	return result
}

// getUserDocDisplayPage 获取文档显示页码信息
func getUserDocDisplayPage(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 优先使用用户编辑的页码
	if doc.PageEdited && doc.UserEditedPage != "" {
		return doc.UserEditedPage
	}

	// 其次使用元数据页码
	return doc.MetaPage
}

// getUserDocDisplayDocType 获取文档显示类型
func getUserDocDisplayDocType(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 优先使用用户编辑的文档类型
	if doc.DocTypeEdited && doc.UserEditedDocType != "" {
		return doc.UserEditedDocType
	}

	// 其次使用元数据文档类型
	return doc.MetaDocType
}

// getUserDocDisplayDoi 获取文档显示DOI信息
func getUserDocDisplayDoi(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 优先使用用户编辑的DOI
	if doc.DoiEdited && doc.UserEditedDoi != "" {
		return doc.UserEditedDoi
	}

	// 其次使用元数据DOI
	return doc.MetaDoi
}

// getUserDocDisplayVolume 获取文档显示卷号信息
func getUserDocDisplayVolume(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 优先使用用户编辑的卷号
	if doc.VolumeEdited && doc.UserEditedVolume != "" {
		return doc.UserEditedVolume
	}

	// 其次使用元数据卷号
	return doc.MetaVolume
}

// getUserDocDisplayIssue 获取文档显示期号信息
func getUserDocDisplayIssue(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 优先使用用户编辑的期号
	if doc.IssueEdited && doc.UserEditedIssue != "" {
		return doc.UserEditedIssue
	}

	// 其次使用元数据期号
	return doc.MetaIssue
}

// getUserDocDisplayUrl 获取文档显示URL信息
func getUserDocDisplayUrl(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 使用元数据URL
	return doc.MetaUrl
}

// getUserDocDisplayLanguage 获取文档显示语言信息
func getUserDocDisplayLanguage(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 使用元数据语言
	return doc.MetaLanguage
}

// getUserDocDisplayEventInfo 获取文档显示事件信息
func getUserDocDisplayEventInfo(doc *model.UserDoc) *docpb.UserDocDisplayEventInfo {
	if doc == nil || (doc.MetaEventTitle == "" && doc.MetaEventPlace == "" && doc.MetaEventDate == "") {
		return nil
	}

	// 创建事件信息对象
	eventTitle := doc.MetaEventTitle
	eventPlace := doc.MetaEventPlace

	eventInfo := &docpb.EventInfo{
		EventTitle: &eventTitle,
		EventPlace: &eventPlace,
	}

	// 设置事件日期
	if doc.MetaEventDate != "" {
		eventInfo.EventDate = []string{doc.MetaEventDate}
	}

	return &docpb.UserDocDisplayEventInfo{
		EventInfo: eventInfo,
	}
}

// GetDocClassifyInfoMap 构建文献-标签信息Map
func GetDocClassifyInfoMap(docClassifyRelations []model.DocClassifyRelation, userDocClassifies []model.UserDocClassify) map[string][]*docpb.UserDocClassifyInfo {
	docClassifiesMap := make(map[string][]*docpb.UserDocClassifyInfo)

	// 按文档ID分组关系
	docRelationMap := make(map[string][]model.DocClassifyRelation)
	for _, relation := range docClassifyRelations {
		docRelationMap[relation.DocId] = append(docRelationMap[relation.DocId], relation)
	}

	// 构建分类ID到分类对象的映射
	classifyMap := make(map[string]model.UserDocClassify)
	for _, classify := range userDocClassifies {
		classifyMap[classify.Id] = classify
	}

	// 为每个文档构建分类信息
	for docId, relations := range docRelationMap {
		var classifyInfos []*docpb.UserDocClassifyInfo

		// 获取该文档的所有分类ID
		classifyIds := make(map[string]bool)
		for _, relation := range relations {
			classifyIds[relation.ClassifyId] = true
		}

		// 构建分类信息
		for classifyId := range classifyIds {
			if classify, ok := classifyMap[classifyId]; ok {
				classifyInfo := &docpb.UserDocClassifyInfo{
					ClassifyId:   classify.Id,
					ClassifyName: classify.Name,
				}
				classifyInfos = append(classifyInfos, classifyInfo)
			}
		}

		docClassifiesMap[docId] = classifyInfos
	}

	return docClassifiesMap
}

// getUserDocDisplayPublishDate 获取文档显示发布日期
func getUserDocDisplayPublishDate(doc *model.UserDoc) string {
	if doc == nil {
		return ""
	}

	// 优先使用显示发布日期
	if doc.DisplayPublishDate != "" {
		return doc.DisplayPublishDate
	}

	// 其次使用元数据发布日期
	if doc.MetaPublishDate != "" {
		return doc.MetaPublishDate
	}

	// 最后使用图表发布日期
	return doc.GraphPublishDate
}

// SortDocList 根据请求中的排序类型对文档列表进行排序
func SortDocList(userDocs []model.UserDoc, req *docpb.GetDocListReq) {
	if len(userDocs) == 0 {
		return
	}

	// 获取排序类型，默认为自定义排序
	sortType := docpb.UserDocListSortType(5) // CUSTOM_SORT
	if req.SortType != nil {
		sortType = *req.SortType
	}

	// 根据排序类型进行排序
	switch sortType {
	case docpb.UserDocListSortType(1): // LAST_READ
		// 按最后阅读时间排序，最近阅读的文档排在前面
		sort.Slice(userDocs, func(i, j int) bool {
			// 如果两个文档都没有阅读记录，则按创建时间排序
			if userDocs[i].LastReadTime.IsZero() && userDocs[j].LastReadTime.IsZero() {
				return userDocs[i].CreatedAt.After(userDocs[j].CreatedAt)
			}
			// 如果只有一个文档没有阅读记录，则有阅读记录的排在前面
			if userDocs[i].LastReadTime.IsZero() {
				return false
			}
			if userDocs[j].LastReadTime.IsZero() {
				return true
			}
			// 如果两个文档都有阅读记录，则按最后阅读时间排序
			return userDocs[i].LastReadTime.After(userDocs[j].LastReadTime)
		})

	case docpb.UserDocListSortType(2): // LAST_ADD
		// 按添加时间排序，最近添加的文档排在前面
		sort.Slice(userDocs, func(i, j int) bool {
			return userDocs[i].CreatedAt.After(userDocs[j].CreatedAt)
		})

	case docpb.UserDocListSortType(3): // DOC_NAME
		// 按文档名称排序，按字母顺序
		sort.Slice(userDocs, func(i, j int) bool {
			return strings.ToLower(userDocs[i].DocName) < strings.ToLower(userDocs[j].DocName)
		})

	case docpb.UserDocListSortType(4): // PUBLISH_DATE
		// 按发布日期排序，最近发布的文档排在前面
		sort.Slice(userDocs, func(i, j int) bool {
			// 获取显示发布日期
			dateI := userDocs[i].DisplayPublishDate
			if dateI == "" {
				dateI = userDocs[i].MetaPublishDate
			}
			if dateI == "" {
				dateI = userDocs[i].GraphPublishDate
			}

			dateJ := userDocs[j].DisplayPublishDate
			if dateJ == "" {
				dateJ = userDocs[j].MetaPublishDate
			}
			if dateJ == "" {
				dateJ = userDocs[j].GraphPublishDate
			}

			// 如果两个文档都没有发布日期，则按创建时间排序
			if dateI == "" && dateJ == "" {
				return userDocs[i].CreatedAt.After(userDocs[j].CreatedAt)
			}
			// 如果只有一个文档没有发布日期，则有发布日期的排在前面
			if dateI == "" {
				return false
			}
			if dateJ == "" {
				return true
			}
			// 如果两个文档都有发布日期，则按发布日期排序
			// 注意：这里的比较是字符串比较，假设日期格式是 YYYY-MM-DD
			return dateI > dateJ
		})

	case docpb.UserDocListSortType(6): // IMPORTANCE_SCORE
		// 按重要性评分排序，评分高的文档排在前面
		sort.Slice(userDocs, func(i, j int) bool {
			if userDocs[i].ImportanceScore == userDocs[j].ImportanceScore {
				// 如果评分相同，则按创建时间排序
				return userDocs[i].CreatedAt.After(userDocs[j].CreatedAt)
			}
			return userDocs[i].ImportanceScore > userDocs[j].ImportanceScore
		})

	default: // docpb.UserDocListSortType(5) // CUSTOM_SORT
		// 按自定义排序，即按 sort 字段排序
		sort.Slice(userDocs, func(i, j int) bool {
			if userDocs[i].Sort == userDocs[j].Sort {
				// 如果排序值相同，则按创建时间排序
				return userDocs[i].CreatedAt.After(userDocs[j].CreatedAt)
			}
			return userDocs[i].Sort < userDocs[j].Sort
		})
	}

	return
}

// GetDescendantFolderIds 获取指定文件夹的所有子文件夹ID
func GetDescendantFolderIds(folderId string, folders []model.UserDocFolder, result *[]string) {
	for _, folder := range folders {
		if folder.ParentId == folderId {
			*result = append(*result, folder.Id)
			GetDescendantFolderIds(folder.Id, folders, result)
		}
	}
}

// GetPbUserDocInfo 将用户文档模型对象转换为Protobuf的UserDocInfo对象
func GetPbUserDocInfo(ctx context.Context, userDoc *model.UserDoc, docClassifiesMap map[string][]*docpb.UserDocClassifyInfo, paperJcrService interface{}) (*docpb.UserDocInfo, error) {
	// 收录情况
	displayVenue := GetUserDocDisplayVenue(userDoc)
	// 发布日期
	displayPublishDate := getUserDocDisplayPublishDate(userDoc)
	// 作者信息
	displayAuthorObj := GetDisplayAuthorByUserDocToUserEdited(userDoc)
	var displayPublishDateObj *docpb.UserDocDisplayPublishDate
	if displayPublishDate != "" {
		displayPublishDateObj = &docpb.UserDocDisplayPublishDate{
			PublishDate: displayPublishDate,
		}
	}
	var displayVenueObj *docpb.UserDocDisplayVenue
	if displayVenue != "" {
		displayVenueObj = &docpb.UserDocDisplayVenue{
			VenueInfos: []string{displayVenue},
		}
	}
	// 获取分区信息
	var displayPartition string
	if userDoc.UserEditedJcrPartion != "" {
		displayPartition = userDoc.UserEditedJcrPartion
	} else {
		// 尝试从JCR实体获取分区
		if paperJcrService != nil {
			// 类型断言，获取GetPaperJcrEntityByVenue方法
			if jcrService, ok := paperJcrService.(interface {
				GetPaperJcrEntityByVenue(ctx context.Context, venue string) (*paperModel.PaperJcrEntity, error)
			}); ok {
				entity, err := jcrService.GetPaperJcrEntityByVenue(ctx, userDoc.Venue)
				if err == nil && entity != nil && entity.JcrPartion != "" {
					displayPartition = entity.JcrPartion
				}
			}
		}
	}
	// 设置最后阅读时间
	var lastReadTime uint64
	if !userDoc.LastReadTime.IsZero() {
		lastReadTime = uint64(userDoc.LastReadTime.UnixNano() / int64(time.Millisecond))
	}
	// 创建分区显示对象
	var displayPartitionObj *docpb.UserDocDisplayPartition
	if displayPartition != "" {
		partitionValue := displayPartition // 创建一个新的变量避免命名冲突
		displayPartitionObj = &docpb.UserDocDisplayPartition{
			Partition: &partitionValue,
		}
	}
	// 组织数据
	return &docpb.UserDocInfo{
		DocId:              userDoc.Id,
		DocName:            userDoc.DocName,
		Sort:               uint32(userDoc.Sort),
		PaperId:            userDoc.PaperId,
		PdfId:              &[]string{userDoc.PdfId}[0],
		Remark:             &userDoc.Remark,
		CreateDate:         uint64(userDoc.CreatedAt.UnixNano() / int64(time.Millisecond)),
		LastReadTime:       &lastReadTime,
		NoteId:             &[]string{userDoc.NoteId}[0],
		NewPaper:           false,
		ClassifyInfos:      docClassifiesMap[userDoc.Id],
		DisplayAuthor:      displayAuthorObj,
		DisplayPublishDate: displayPublishDateObj,
		DisplayVenue:       displayVenueObj,
		DisplayPartition:   displayPartitionObj,
	}, nil
}

// GetAllInfoPbUserDocInfo 将用户文档模型对象转换为包含全部信息的Protobuf UserDocInfo对象
func GetAllInfoPbUserDocInfo(ctx context.Context, doc *model.UserDoc, docClassifiesMap map[string][]*docpb.UserDocClassifyInfo,
	docSearchResultMap map[string]*docpb.DocSearchResult, lastReadDoc *model.UserDoc, paperJcrService interface{}) *docpb.UserDocInfo {
	// 页码信息
	displayPage := GetUserDocPageToUserEdited(doc)
	// 创建文档显示信息
	displayAuthor := GetDisplayAuthorByUserDocToUserEdited(doc)
	// 收录情况
	displayVenue := GetUserDocDisplayVenueToUserEdited(doc)
	// 语言信息
	displayLanguage := GetUserDocLanguageToUserEdited(doc)
	// 发布日期
	displayPublishDate := GetPublishDateToUserEdited(doc)
	// 文档类型
	displayDocType := GetUserDocTypeToUserEdited(doc)
	// DOI信息
	displayDoi := GetUserDocDoiToUserEdited(doc)
	// 卷号信息
	displayVolume := GetUserDocVolumeToUserEdited(doc)
	// 期号信息
	displayIssue := GetUserDocIssueToUserEdited(doc)
	// 事件信息
	displayEventInfo := GetUserDocEventInfoToUserEdited(doc)
	// 影响因子
	displayImpactFactor := GetUserDocImpactFactorToUserEdited(doc)
	// 分区信息
	displayPartition := GetUserDocPartitionToUserEdited(doc)

	// 论文存储库状态
	var paperRepositoryStatus docpb.PaperRepositoryStatus
	if doc.PaperRepositoryStatus == "" {
		if doc.PaperId != "0" {
			paperRepositoryStatus = docpb.PaperRepositoryStatus_IN_REPOSITORY
		} else {
			paperRepositoryStatus = docpb.PaperRepositoryStatus_NOT_IN_REPOSITORY
		}
	} else {
		// 尝试解析状态字符串
		if doc.PaperRepositoryStatus == "IN_REPOSITORY" {
			paperRepositoryStatus = docpb.PaperRepositoryStatus_IN_REPOSITORY
		} else {
			paperRepositoryStatus = docpb.PaperRepositoryStatus_NOT_IN_REPOSITORY
		}
	}
	// 阅读进度
	progress := int32(doc.Progress)
	if doc.ReadingStatus == "READING" && doc.Progress == 0 {
		progress = 1
	}
	// 是否是最近阅读的文献
	isLatestRead := lastReadDoc != nil && lastReadDoc.Id == doc.Id
	// 创建文档信息对象
	var lastReadTime uint64
	if !doc.LastReadTime.IsZero() {
		lastReadTime = uint64(doc.LastReadTime.UnixNano() / int64(time.Millisecond))
	}
	// 解析阅读状态
	var docReadingStatus docpb.DocReadingStatus
	switch doc.ReadingStatus {
	case "READING":
		docReadingStatus = docpb.DocReadingStatus_READING
	case "READ":
		docReadingStatus = docpb.DocReadingStatus_READ
	default:
		docReadingStatus = docpb.DocReadingStatus_UNREAD
	}
	// 组织数据
	return &docpb.UserDocInfo{
		DocId:                 doc.Id,
		DocName:               doc.DocName,
		Sort:                  uint32(doc.Sort),
		PaperId:               doc.PaperId,
		PdfId:                 &[]string{doc.PdfId}[0],
		Remark:                &doc.Remark,
		CreateDate:            uint64(doc.CreatedAt.UnixNano() / int64(time.Millisecond)),
		LastReadTime:          &lastReadTime,
		IsLatestRead:          isLatestRead,
		NoteId:                &[]string{doc.NoteId}[0],
		NewPaper:              doc.NewPaper,
		ClassifyInfos:         docClassifiesMap[doc.Id],
		DisplayAuthor:         displayAuthor,
		DisplayPublishDate:    displayPublishDate,
		DisplayVenue:          displayVenue,
		PaperRepositoryStatus: paperRepositoryStatus,
		SearchResult:          docSearchResultMap[doc.Id],
		HasAttachment:         false, // 根据要求忽略附件相关逻辑
		DisplayPage:           displayPage,
		DisplayDocType:        displayDocType,
		DisplayDoi:            displayDoi,
		DisplayVolume:         displayVolume,
		DisplayIssue:          displayIssue,
		DisplayLanguage:       displayLanguage,
		DisplayEventInfo:      displayEventInfo,
		JcrVenuePartion:       displayPartition,
		ImpactOfFactor:        displayImpactFactor,
		ImportantanceScore:    int32(doc.ImportanceScore),
		FillExtMeta:           doc.FillExtMeta,
		DocReadingStatus:      docReadingStatus,
		Progress:              &progress,
		ParsedStatus:          docpb.UserDocParsedStatusEnum(doc.ParseStatus),
	}
}

// GetUserDocDetailInfo 将UserDocInfo对象转换为DocDetailInfo对象
func GetUserDocDetailInfo(ctx context.Context, userDocInfo *docpb.UserDocInfo, userDoc *model.UserDoc) (*docpb.DocDetailInfo, error) {
	// 创建DocDetailInfo对象
	docDetailInfo := &docpb.DocDetailInfo{}
	// 设置基本字段
	docDetailInfo.DocId = userDocInfo.DocId
	docDetailInfo.DocName = userDocInfo.DocName
	docDetailInfo.PaperId = userDocInfo.PaperId
	// 设置分区信息
	displayVenue := GetVenueStrByUserDoc(userDoc)
	docDetailInfo.DisplayVenue = displayVenue
	//作者信息
	displayAuthor := GetDisplayAuthorStrByUserDoc(userDoc)
	docDetailInfo.DisplayAuthor = displayAuthor
	//发布日期
	displayPublishDate := GetPublishDateStrByUserDoc(userDoc)
	docDetailInfo.DisplayPublishDate = displayPublishDate
	docDetailInfo.PdfId = *userDocInfo.PdfId
	//设置状态
	docDetailInfo.ParsedStatus = docpb.UserDocParsedStatusEnum(userDoc.ParseStatus)

	return docDetailInfo, nil
}

// ================   这里tm的这么写是因为要符合之前的数据结构   ============================
// 根据userDoc获取displayAuthor
func GetDisplayAuthorByUserDocToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayAuthor {
	if userDoc == nil {
		return &docpb.UserDocDisplayAuthor{}
	}
	//userEdited 如果AuthorDesc和MetaAuthors不一致则为true
	userEdited := false
	if userDoc.AuthorDesc != userDoc.MetaAuthors {
		userEdited = true
	}
	displayAuthor := docpb.UserDocDisplayAuthor{
		UserEdited: userEdited,
	}
	//查询作者信息
	if userDoc.AuthorDesc != "" {
		var authorInfos []*docpb.UserDocAuthorInfo
		err := json.Unmarshal([]byte(userDoc.AuthorDesc), &authorInfos)
		if err != nil {
			return &displayAuthor
		}
		displayAuthor.AuthorInfos = authorInfos
	}
	//查询显示作者信息
	if userDoc.DisplayAuthors != "" {
		//userDoc.MetaAuthors 转对象数组
		var authorInfos []*docpb.UserDocAuthorInfo
		err := json.Unmarshal([]byte(userDoc.MetaAuthors), &authorInfos)
		if err != nil {
			return &displayAuthor
		}
		displayAuthor.OriginAuthorInfos = authorInfos
	}
	return &displayAuthor
}

func GetDisplayAuthorStrByUserDoc(userDoc *model.UserDoc) *docpb.DisplayAuthor {
	if userDoc == nil {
		return &docpb.DisplayAuthor{}
	}
	displayAuthor := docpb.DisplayAuthor{
		RollbackEnable: false,
	}
	//查询作者信息
	if userDoc.AuthorDesc != "" {
		var authorInfos []*docpb.UserDocAuthorInfo
		err := json.Unmarshal([]byte(userDoc.AuthorDesc), &authorInfos)
		if err != nil {
			return &displayAuthor
		}
		authors := []string{}
		for _, authorInfo := range authorInfos {
			authors = append(authors, authorInfo.Literal)
		}
		displayAuthor.Authors = authors
	}
	//查询显示作者信息
	if userDoc.DisplayAuthors != "" {
		//userDoc.MetaAuthors 转对象数组
		var authorInfos []*docpb.UserDocAuthorInfo
		err := json.Unmarshal([]byte(userDoc.MetaAuthors), &authorInfos)
		if err != nil {
			return &displayAuthor
		}
		originAuthors := []string{}
		for _, authorInfo := range authorInfos {
			originAuthors = append(originAuthors, authorInfo.Literal)
		}
		displayAuthor.OriginAuthors = originAuthors
	}
	return &displayAuthor
}

func GetDisplaySimpleAuthorByUserDoc(userDoc *model.UserDoc) *docpb.DisplaySimpleAuthors {
	if userDoc == nil {
		return &docpb.DisplaySimpleAuthors{}
	}
	//设置这个RollbackEnable 如果AuthorDesc和MetaAuthors不一致则为true
	rollbackEnable := false
	if userDoc.AuthorDesc != userDoc.MetaAuthors {
		rollbackEnable = true
	}
	displayAuthor := docpb.DisplaySimpleAuthors{
		RollbackEnable: rollbackEnable,
	}
	//查询作者信息
	if userDoc.AuthorDesc != "" {
		var authorInfos []*docpb.UserDocAuthorInfo
		err := json.Unmarshal([]byte(userDoc.AuthorDesc), &authorInfos)
		if err != nil {
			return &displayAuthor
		}
		authorInfoResults := []*docpb.BaseAuthorInfo{}
		for _, authorInfo := range authorInfos {
			authorInfoResults = append(authorInfoResults, &docpb.BaseAuthorInfo{
				Name: authorInfo.Literal,
			})
		}
		displayAuthor.Authors = authorInfoResults
	}
	//查询显示作者信息
	if userDoc.DisplayAuthors != "" {
		//userDoc.MetaAuthors 转对象数组
		var authorInfos []*docpb.UserDocAuthorInfo
		err := json.Unmarshal([]byte(userDoc.MetaAuthors), &authorInfos)
		if err != nil {
			return &displayAuthor
		}
		originAuthorInfos := []*docpb.BaseAuthorInfo{}
		for _, authorInfo := range authorInfos {
			originAuthorInfos = append(originAuthorInfos, &docpb.BaseAuthorInfo{
				Name: authorInfo.Literal,
			})
		}
		displayAuthor.OriginAuthors = originAuthorInfos
	}
	return &displayAuthor
}

func BaseAuthorInfoToUserDocAuthorInfo(baseAuthorInfos []*docpb.BaseAuthorInfo) []*docpb.UserDocAuthorInfo {
	var authorInfos []*docpb.UserDocAuthorInfo
	if len(baseAuthorInfos) == 0 {
		return authorInfos
	}
	for _, baseAuthorInfo := range baseAuthorInfos {
		// 手动找出姓和名  这里还得区分中英文  暂时不做
		authorName := baseAuthorInfo.Name
		if authorName == "" {
			continue
		}
		authorInfos = append(authorInfos, &docpb.UserDocAuthorInfo{
			Literal: authorName,
		})
	}
	return authorInfos
}
func GetPublishDateToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayPublishDate {
	if userDoc == nil {
		return nil
	}
	//userEdited 如果PublishDate和MetaPublishDate不一致则为true
	userEdited := false
	if userDoc.PublishDate != userDoc.MetaPublishDate {
		userEdited = true
	}
	publishDate := docpb.UserDocDisplayPublishDate{
		PublishDate:       userDoc.PublishDate,
		OriginPublishDate: userDoc.MetaPublishDate,
		UserEdited:        userEdited,
	}
	return &publishDate
}
func GetPublishDateStrByUserDoc(userDoc *model.UserDoc) *docpb.DisplayPublishDate {
	if userDoc == nil {
		return nil
	}
	//设置这个RollbackEnable 如果PublishDate和MetaPublishDate不一致则为true
	rollbackEnable := false
	if userDoc.PublishDate != userDoc.MetaPublishDate {
		rollbackEnable = true
	}
	publishDate := docpb.DisplayPublishDate{
		PublishDate:       userDoc.PublishDate,
		OriginPublishDate: userDoc.MetaPublishDate,
		RollbackEnable:    rollbackEnable,
	}
	return &publishDate
}

func GetVenueStrByUserDoc(userDoc *model.UserDoc) *docpb.DisplayVenue {
	if userDoc == nil {
		return nil
	}
	//设置这个RollbackEnable 如果Venue和MetaVenues不一致则为true
	rollbackEnable := false
	if userDoc.Venue != userDoc.MetaVenues {
		rollbackEnable = true
	}
	venue := docpb.DisplayVenue{
		Venue:          userDoc.Venue,
		OriginVenue:    userDoc.MetaVenues,
		RollbackEnable: rollbackEnable,
	}
	return &venue
}

func GetUserDocDisplayVenueToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayVenue {
	if userDoc == nil {
		return nil
	}
	//设置这个RollbackEnable 如果Venue和MetaVenues不一致则为true
	UserEdited := false
	if userDoc.Venue != userDoc.MetaVenues {
		UserEdited = true
	}
	venue := docpb.UserDocDisplayVenue{
		VenueInfos:       []string{userDoc.Venue},
		OriginVenueInfos: []string{userDoc.MetaVenues},
		UserEdited:       UserEdited,
	}
	return &venue
}

func GetUserDocLanguageToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayLanguage {
	if userDoc == nil {
		return nil
	}
	language := userDoc.MetaLanguage
	UserEdited := false
	userDocDisplayLanguage := docpb.UserDocDisplayLanguage{
		Language:       &language,
		OriginLanguage: &language,
		UserEdited:     &UserEdited,
	}
	return &userDocDisplayLanguage
}

func GetUserDocPageToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayPage {
	if userDoc == nil {
		return nil
	}
	UserEdited := false
	if userDoc.UserEditedPage != userDoc.MetaPage {
		UserEdited = true
	}
	userDocDisplayPage := docpb.UserDocDisplayPage{
		Page:       &userDoc.UserEditedPage,
		OriginPage: &userDoc.MetaPage,
		UserEdited: &UserEdited,
	}
	return &userDocDisplayPage
}

func GetUserDocTypeToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayDocType {
	if userDoc == nil {
		return nil
	}
	UserEdited := false
	if userDoc.UserEditedDocType != userDoc.MetaDocType {
		UserEdited = true
	}
	userDocDisplayDocType := docpb.UserDocDisplayDocType{
		DocType:       &userDoc.UserEditedDocType,
		OriginDocType: &userDoc.MetaDocType,
		UserEdited:    &UserEdited,
	}
	return &userDocDisplayDocType
}

func GetUserDocDoiToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayDoi {
	if userDoc == nil {
		return nil
	}
	UserEdited := false
	if userDoc.UserEditedDoi != userDoc.MetaDoi {
		UserEdited = true
	}
	userDocDisplayDoi := docpb.UserDocDisplayDoi{
		Doi:        &userDoc.UserEditedDoi,
		OriginDoi:  &userDoc.MetaDoi,
		UserEdited: &UserEdited,
	}
	return &userDocDisplayDoi
}

func GetUserDocVolumeToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayVolume {
	if userDoc == nil {
		return nil
	}
	UserEdited := false
	if userDoc.UserEditedVolume != userDoc.MetaVolume {
		UserEdited = true
	}
	userDocDisplayVolume := docpb.UserDocDisplayVolume{
		Volume:       &userDoc.UserEditedVolume,
		OriginVolume: &userDoc.MetaVolume,
		UserEdited:   &UserEdited,
	}
	return &userDocDisplayVolume
}

func GetUserDocIssueToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayIssue {
	if userDoc == nil {
		return nil
	}
	UserEdited := false
	if userDoc.UserEditedIssue != userDoc.MetaIssue {
		UserEdited = true
	}
	userDocDisplayIssue := docpb.UserDocDisplayIssue{
		Issue:       &userDoc.UserEditedIssue,
		OriginIssue: &userDoc.MetaIssue,
		UserEdited:  &UserEdited,
	}
	return &userDocDisplayIssue
}

func GetUserDocEventInfoToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayEventInfo {
	if userDoc == nil {
		return nil
	}
	// 没有事件，这里暂时为空
	UserEdited := false
	userDocDisplayEventInfo := docpb.UserDocDisplayEventInfo{
		UserEdited: &UserEdited,
	}
	return &userDocDisplayEventInfo
}

// TODO： 因为在论文解析的时候是没有影响因子的  这里使用空，后续如果有逻辑进行补充的话需要修改此方法
func GetUserDocImpactFactorToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayJcrImpactFactor {
	if userDoc == nil {
		return nil
	}
	UserEdited := false
	if userDoc.UserEditedImpactOfFactor != 0 {
		UserEdited = true
	}
	userDocDisplayImpactFactor := docpb.UserDocDisplayJcrImpactFactor{
		ImpactOfFactor: &userDoc.UserEditedImpactOfFactor,
		UserEdited:     &UserEdited,
	}
	return &userDocDisplayImpactFactor
}

// TODO： 因为在论文解析的时候是没有分区的  这里使用空，后续如果有逻辑进行补充的话需要修改此方法
func GetUserDocPartitionToUserEdited(userDoc *model.UserDoc) *docpb.UserDocDisplayJcrVenuePartion {
	if userDoc == nil {
		return nil
	}
	UserEdited := false
	if userDoc.UserEditedJcrPartion != userDoc.MetaPartition {
		UserEdited = true
	}
	userDocDisplayPartition := docpb.UserDocDisplayJcrVenuePartion{
		JcrVenuePartion:       &userDoc.UserEditedJcrPartion,
		OriginJcrVenuePartion: &userDoc.MetaPartition,
		UserEdited:            &UserEdited,
	}
	return &userDocDisplayPartition
}

// 论文解析状态转换进度比
func GetUserDocParsedStatusEnum(userDoc *docpb.UserDocInfo) string {
	progress := "0%"
	if userDoc == nil {
		return progress
	}

	switch userDoc.ParsedStatus {
	case docpb.UserDocParsedStatusEnum_READY:
		progress = "0%"
	case docpb.UserDocParsedStatusEnum_REPARSE:
		progress = "0%"
	case docpb.UserDocParsedStatusEnum_DOWNLOADING:
		progress = "5%"
	case docpb.UserDocParsedStatusEnum_DOWNLOADED:
		progress = "10%"
	case docpb.UserDocParsedStatusEnum_DOWNLOAD_FAILED:
		progress = "15%"
	case docpb.UserDocParsedStatusEnum_UPLOADING:
		progress = "20%"
	case docpb.UserDocParsedStatusEnum_UPLOADED:
		progress = "25%"
	case docpb.UserDocParsedStatusEnum_UPLOAD_FAILED:
		progress = "30%"
	case docpb.UserDocParsedStatusEnum_GENERATING_BASE_DATA:
		progress = "40%"
	case docpb.UserDocParsedStatusEnum_BASE_DATA_GENERATED:
		progress = "45%"
	case docpb.UserDocParsedStatusEnum_BASE_DATA_GENERATE_FAILED:
		progress = "40%"
	case docpb.UserDocParsedStatusEnum_PARSING_HEADER_DATA:
		progress = "50%"
	case docpb.UserDocParsedStatusEnum_HEADER_DATA_PARSED:
		progress = "60%"
	case docpb.UserDocParsedStatusEnum_HEADER_DATA_PARSE_FAILED:
		progress = "50%"
	case docpb.UserDocParsedStatusEnum_PARSING_CONTENT_DATA:
		progress = "70%"
	case docpb.UserDocParsedStatusEnum_CONTENT_DATA_PARSED:
		progress = "80%"
		if userDoc.EmbeddingStatus == docpb.UserDocParsedStatusEnum_EMBEDDED {
			progress = "100%"
		}
	case docpb.UserDocParsedStatusEnum_CONTENT_DATA_PARSE_FAILED:
		progress = "70%"
	case docpb.UserDocParsedStatusEnum_PARSE_FAILED:
		progress = "0%"
	default:
		// unknown 或其他未知状态
		progress = "0%"
	}

	return progress
}
