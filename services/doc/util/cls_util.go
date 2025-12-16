package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	pb "github.com/yb2020/odoc/proto/gen/go/doc"
)

// CSLType 表示CSL文档类型
type CSLType string

// CSL文档类型常量
const (
	CSLTypeArticle               CSLType = "article"
	CSLTypeArticleJournal        CSLType = "article-journal"
	CSLTypeArticleMagazine       CSLType = "article-magazine"
	CSLTypeArticleNewspaper      CSLType = "article-newspaper"
	CSLTypeBill                  CSLType = "bill"
	CSLTypeBook                  CSLType = "book"
	CSLTypeBroadcast             CSLType = "broadcast"
	CSLTypeChapter               CSLType = "chapter"
	CSLTypeDataset               CSLType = "dataset"
	CSLTypeEntry                 CSLType = "entry"
	CSLTypeEntryDictionary       CSLType = "entry-dictionary"
	CSLTypeEntryEncyclopedia     CSLType = "entry-encyclopedia"
	CSLTypeFigure                CSLType = "figure"
	CSLTypeGraphic               CSLType = "graphic"
	CSLTypeInterview             CSLType = "interview"
	CSLTypeLegalCase             CSLType = "legal_case"
	CSLTypeLegislation           CSLType = "legislation"
	CSLTypeManuscript            CSLType = "manuscript"
	CSLTypeMap                   CSLType = "map"
	CSLTypeMotionPicture         CSLType = "motion_picture"
	CSLTypeMusicalScore          CSLType = "musical_score"
	CSLTypePamphlet              CSLType = "pamphlet"
	CSLTypePaperConference       CSLType = "paper-conference"
	CSLTypePatent                CSLType = "patent"
	CSLTypePersonalCommunication CSLType = "personal_communication"
	CSLTypePost                  CSLType = "post"
	CSLTypePostWeblog            CSLType = "post-weblog"
	CSLTypeReport                CSLType = "report"
	CSLTypeReview                CSLType = "review"
	CSLTypeReviewBook            CSLType = "review-book"
	CSLTypeSong                  CSLType = "song"
	CSLTypeSpeech                CSLType = "speech"
	CSLTypeThesis                CSLType = "thesis"
	CSLTypeTreaty                CSLType = "treaty"
	CSLTypeWebpage               CSLType = "webpage"
)

// String 返回CSLType的字符串表示
func (t CSLType) String() string {
	return string(t)
}

// CSLTypeFromString 从字符串转换为CSLType
func CSLTypeFromString(str string) (CSLType, error) {
	switch str {
	case "article":
		return CSLTypeArticle, nil
	case "article-journal", "article journal":
		return CSLTypeArticleJournal, nil
	case "article-magazine", "article magazine":
		return CSLTypeArticleMagazine, nil
	case "article-newspaper", "article newspaper":
		return CSLTypeArticleNewspaper, nil
	case "bill":
		return CSLTypeBill, nil
	case "book":
		return CSLTypeBook, nil
	case "broadcast":
		return CSLTypeBroadcast, nil
	case "chapter":
		return CSLTypeChapter, nil
	case "dataset":
		return CSLTypeDataset, nil
	case "entry":
		return CSLTypeEntry, nil
	case "entry-dictionary", "entry dictionary":
		return CSLTypeEntryDictionary, nil
	case "entry-encyclopedia", "entry encyclopedia":
		return CSLTypeEntryEncyclopedia, nil
	case "figure":
		return CSLTypeFigure, nil
	case "graphic":
		return CSLTypeGraphic, nil
	case "interview":
		return CSLTypeInterview, nil
	case "legal_case":
		return CSLTypeLegalCase, nil
	case "legislation":
		return CSLTypeLegislation, nil
	case "manuscript":
		return CSLTypeManuscript, nil
	case "map":
		return CSLTypeMap, nil
	case "motion_picture":
		return CSLTypeMotionPicture, nil
	case "musical_score":
		return CSLTypeMusicalScore, nil
	case "pamphlet":
		return CSLTypePamphlet, nil
	case "paper-conference", "paper conference":
		return CSLTypePaperConference, nil
	case "patent":
		return CSLTypePatent, nil
	case "personal_communication":
		return CSLTypePersonalCommunication, nil
	case "post":
		return CSLTypePost, nil
	case "post-weblog", "post weblog":
		return CSLTypePostWeblog, nil
	case "report":
		return CSLTypeReport, nil
	case "review":
		return CSLTypeReview, nil
	case "review-book", "review book":
		return CSLTypeReviewBook, nil
	case "song":
		return CSLTypeSong, nil
	case "speech":
		return CSLTypeSpeech, nil
	case "thesis":
		return CSLTypeThesis, nil
	case "treaty":
		return CSLTypeTreaty, nil
	case "webpage":
		return CSLTypeWebpage, nil
	default:
		return "", fmt.Errorf("未知的CSLType: %s", str)
	}
}

// CSLName 表示CSL中的名称
type CSLName struct {
	Family              string `json:"family,omitempty"`
	Given               string `json:"given,omitempty"`
	DroppingParticle    string `json:"dropping-particle,omitempty"`
	NonDroppingParticle string `json:"non-dropping-particle,omitempty"`
	Suffix              string `json:"suffix,omitempty"`
	CommaPrefix         *bool  `json:"comma-prefix,omitempty"`
	CommaSuffix         *bool  `json:"comma-suffix,omitempty"`
	StaticOrdering      *bool  `json:"static-ordering,omitempty"`
	StaticParticles     *bool  `json:"static-particles,omitempty"`
	Literal             string `json:"literal,omitempty"`
	ParseNames          *bool  `json:"parse-names,omitempty"`
	IsInstitution       *bool  `json:"isInstitution,omitempty"`
}

// NewCSLName 创建一个新的CSLName实例
func NewCSLName() *CSLName {
	return &CSLName{}
}

// NewCSLNameWithValues 使用指定值创建一个新的CSLName实例
func NewCSLNameWithValues(family, given, droppingParticle, nonDroppingParticle, suffix string,
	commaPrefix, commaSuffix, staticOrdering, staticParticles *bool,
	literal string, parseNames, isInstitution *bool) *CSLName {

	return &CSLName{
		Family:              family,
		Given:               given,
		DroppingParticle:    droppingParticle,
		NonDroppingParticle: nonDroppingParticle,
		Suffix:              suffix,
		CommaPrefix:         commaPrefix,
		CommaSuffix:         commaSuffix,
		StaticOrdering:      staticOrdering,
		StaticParticles:     staticParticles,
		Literal:             literal,
		ParseNames:          parseNames,
		IsInstitution:       isInstitution,
	}
}

// CSLNameFromJSON 从JSON对象创建CSLName
func CSLNameFromJSON(obj map[string]interface{}) *CSLName {
	name := NewCSLName()

	if v, ok := obj["family"]; ok && v != nil {
		name.Family = fmt.Sprintf("%v", v)
	}

	if v, ok := obj["given"]; ok && v != nil {
		name.Given = fmt.Sprintf("%v", v)
	}

	if v, ok := obj["dropping-particle"]; ok && v != nil {
		name.DroppingParticle = fmt.Sprintf("%v", v)
	}

	if v, ok := obj["non-dropping-particle"]; ok && v != nil {
		name.NonDroppingParticle = fmt.Sprintf("%v", v)
	}

	if v, ok := obj["suffix"]; ok && v != nil {
		name.Suffix = fmt.Sprintf("%v", v)
	}

	if v, ok := obj["comma-prefix"]; ok && v != nil {
		b := toBool(v)
		name.CommaPrefix = &b
	}

	if v, ok := obj["comma-suffix"]; ok && v != nil {
		b := toBool(v)
		name.CommaSuffix = &b
	}

	if v, ok := obj["static-ordering"]; ok && v != nil {
		b := toBool(v)
		name.StaticOrdering = &b
	}

	if v, ok := obj["static-particles"]; ok && v != nil {
		b := toBool(v)
		name.StaticParticles = &b
	}

	if v, ok := obj["literal"]; ok && v != nil {
		name.Literal = fmt.Sprintf("%v", v)
	}

	if v, ok := obj["parse-names"]; ok && v != nil {
		b := toBool(v)
		name.ParseNames = &b
	}

	if v, ok := obj["isInstitution"]; ok && v != nil {
		b := toBool(v)
		name.IsInstitution = &b
	}

	// 处理多语言支持
	if v, ok := obj["multi"]; ok && v != nil {
		if multi, ok := v.(map[string]interface{}); ok {
			if main, ok := multi["main"]; ok && main != nil {
				if mainStr, ok := main.(string); ok {
					mainStr = strings.ToLower(mainStr)
					if len(mainStr) >= 2 {
						mainStr = mainStr[:2]
						if mainStr == "hu" || mainStr == "vi" {
							b := true
							name.StaticOrdering = &b
						}
					}
				}
			}
		}
	}

	return name
}

// ToJSON 将CSLName转换为JSON
func (n *CSLName) ToJSON() (map[string]interface{}, error) {
	data, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// toBool 将各种类型转换为布尔值
func toBool(v interface{}) bool {
	switch val := v.(type) {
	case string:
		b, _ := strconv.ParseBool(val)
		return b
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(val).Int() != 0
	case float32, float64:
		return reflect.ValueOf(val).Float() != 0
	case bool:
		return val
	default:
		return false
	}
}

// 以下是BibTeX导出相关的辅助方法

// GetContainerTitle 获取容器标题（如期刊名称）
func GetContainerTitle(metaInfoSimpleVo *pb.DocMetaInfoSimpleVo) string {
	if metaInfoSimpleVo == nil || len(metaInfoSimpleVo.GetContainerTitle()) == 0 {
		return ""
	}
	return strings.Join(metaInfoSimpleVo.GetContainerTitle(), ",")
}

// GetMonth 获取发布月份
func GetMonth(metaInfoSimpleVo *pb.DocMetaInfoSimpleVo) int {
	if metaInfoSimpleVo == nil || metaInfoSimpleVo.GetPublishTimestamp() == 0 {
		return 0
	}
	// 将毫秒时间戳转换为秒，并显式转换为int64类型
	t := time.Unix(int64(metaInfoSimpleVo.GetPublishTimestamp()/1000), 0)
	// Go的time.Month是从1开始的，不需要+1
	return int(t.Month())
}

// GetYear 获取发布年份
func GetYear(metaInfoSimpleVo *pb.DocMetaInfoSimpleVo) int {
	if metaInfoSimpleVo == nil || metaInfoSimpleVo.GetPublishTimestamp() == 0 {
		return 0
	}
	// 显式转换为int64类型
	t := time.Unix(int64(metaInfoSimpleVo.GetPublishTimestamp()/1000), 0)
	return t.Year()
}

// GetAuthor 获取作者列表并转换为CSL格式
func GetAuthor(metaInfoSimpleVo *pb.DocMetaInfoSimpleVo) []*CSLName {
	if metaInfoSimpleVo == nil || len(metaInfoSimpleVo.GetAuthorList()) == 0 {
		return []*CSLName{}
	}

	authors := make([]*CSLName, 0, len(metaInfoSimpleVo.GetAuthorList()))
	for _, author := range metaInfoSimpleVo.GetAuthorList() {
		// 创建作者的JSON表示
		authorMap := make(map[string]interface{})

		// 将Proto对象转换为map
		authorBytes, err := json.Marshal(author)
		if err != nil {
			continue
		}

		// 解析为map
		if err := json.Unmarshal(authorBytes, &authorMap); err != nil {
			continue
		}

		// 确保literal字段为空
		delete(authorMap, "literal")

		// 从map创建CSLName
		authors = append(authors, CSLNameFromJSON(authorMap))
	}

	return authors
}

// GetType 获取文档类型
func GetType(metaInfoSimpleVo *pb.DocMetaInfoSimpleVo) CSLType {
	if metaInfoSimpleVo == nil {
		return CSLTypeArticleJournal
	}

	docType := metaInfoSimpleVo.GetDocType()
	if docType == "" {
		return CSLTypeArticleJournal
	}

	cslType, err := CSLTypeFromString(docType)
	if err != nil {
		return CSLTypeArticleJournal
	}

	return cslType
}

// yyyy-mm-dd时间字符串转时间戳
func GetTimestampByDate(date string) int64 {
	if date == "" {
		return 0
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0
	}
	return t.Unix()
}
