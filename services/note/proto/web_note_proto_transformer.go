package proto

import (
	"strings"

	common "github.com/yb2020/odoc-proto/gen/go/common"
	pb "github.com/yb2020/odoc-proto/gen/go/note"
	"github.com/yb2020/odoc/pkg/errors"
)

/**
 * Web笔记标注协议 <--> DB通用标注协议
 */
type WebNoteProtoTransformer struct {
}

// Web标注协议 -> 统一标注协议 AnnotationRawModel 是 ToAnnotationRawModel 的别名，保持向后兼容性
func (s *WebNoteProtoTransformer) AnnotationRawModel(webAnnotation *pb.WebNoteAnnotationModel) (*pb.AnnotationRawModel, error) {
	return s.ToAnnotationRawModel(webAnnotation)
}

// Web标注协议 -> 统一标注协议
func (s *WebNoteProtoTransformer) ToAnnotationRawModel(webAnnotation *pb.WebNoteAnnotationModel) (*pb.AnnotationRawModel, error) {
	if webAnnotation == nil {
		return nil, errors.Biz("Web标注模型不能为空")
	}

	// 创建AnnotationRawModel
	annotationRaw := &pb.AnnotationRawModel{
		Type:        uint32(webAnnotation.Type),
		GroupId:     webAnnotation.GroupId,
		IsHighlight: webAnnotation.IsHighlight,
	}

	// 根据标注类型处理
	switch webAnnotation.Type {
	case common.IDEAAnnotateType_IDEAAnnotateTypeComment:
		//文字选区标注类型
		comment, err := s.fillUpCommentModel(webAnnotation)
		if err != nil {
			return nil, err
		}
		annotationRaw.Comment = comment
	case common.IDEAAnnotateType_IDEAAnnotateTypeRect:
		//图片选区标注类型
		rect, err := s.fillUpRectModel(webAnnotation)
		if err != nil {
			return nil, err
		}
		annotationRaw.Rect = rect
	case common.IDEAAnnotateType_IDEAAnnotateTypeTextBox:
		//文本框标注类型
		textBox, err := s.fillUpTextBoxModel(webAnnotation)
		if err != nil {
			return nil, err
		}
		annotationRaw.TextBoxModel = textBox
	default:
		return nil, errors.Biz("不支持的标注类型")
	}

	return annotationRaw, nil
}

/*
*  填充注释模型
*  WebNoteAnnotationModel.select(WebNoteSelectAnnotation) --> AnnotationRawModel.comment(CommentRawModel)
 */
func (s *WebNoteProtoTransformer) fillUpCommentModel(webAnnotation *pb.WebNoteAnnotationModel) (*pb.CommentRawModel, error) {
	selectAnnot := webAnnotation.Select
	if selectAnnot == nil {
		return nil, errors.Biz("缺少文字标注信息")
	}

	commentRaw := &pb.CommentRawModel{
		Id:         selectAnnot.Uuid,
		Idea:       selectAnnot.Idea,
		NoteId:     selectAnnot.DocumentId,
		Type:       "comment",
		RectStr:    selectAnnot.RectStr,
		AnnotateId: selectAnnot.Uuid,
		StyleId:    selectAnnot.StyleId,
		Fill:       new(NoteAnnotationStyleTransformer).GetColorByStyleId(selectAnnot.StyleId),
		Page:       webAnnotation.PageNumber,
		PdfId:      webAnnotation.PdfId,
	}

	if webAnnotation.Position != -1 {
		commentRaw.Sort = uint32(webAnnotation.Position)
	}

	// 补充选区
	rectOptions := selectAnnot.Rectangle
	if len(rectOptions) > 0 {
		rectangles := make([]*pb.Rectangle, 0)
		for _, rectOption := range rectOptions {
			if rectOption != nil {
				rectangles = append(rectangles, s.getRectFromOptions(rectOption))
			}
		}
		commentRaw.Rectangles = rectangles
	}
	return commentRaw, nil
}

/*
*  图片标注
*  WebNoteAnnotationModel.rect(WebNoteRectAnnotation) --> AnnotationRawModel.rect(RectRawModel)
 */
func (s *WebNoteProtoTransformer) fillUpRectModel(webAnnotation *pb.WebNoteAnnotationModel) (*pb.RectRawModel, error) {
	rectAnnot := webAnnotation.Rect
	if rectAnnot == nil {
		return nil, errors.Biz("缺少图片标注信息")
	}

	picUrl := rectAnnot.PicUrl
	if picUrl != "" {
		// Go中没有直接的URL验证器，这里可以使用正则表达式或其他方式验证
		// 简单验证URL格式
		if !strings.HasPrefix(picUrl, "http://") && !strings.HasPrefix(picUrl, "https://") {
			return nil, errors.Biz("图片链接不符合规范")
		}

		if len(picUrl) > 2048 {
			return nil, errors.Biz("图片链接过长")
		}
	}

	// 创建RectRawModel
	rectModel := &pb.RectRawModel{
		Id:       rectAnnot.Uuid,
		PicUrl:   picUrl,
		IsDelete: false,
		NoteId:   rectAnnot.DocumentId,
		Type:     "rect",
		StyleId:  rectAnnot.StyleId,
		Idea:     rectAnnot.Idea,
		Page:     webAnnotation.PageNumber,
		PdfId:    webAnnotation.PdfId,
	}

	// 设置样式相关字段
	styleTransformer := &NoteAnnotationStyleTransformer{}
	color := styleTransformer.GetColorByStyleId(rectAnnot.StyleId)
	rectModel.LineColor = color
	rectModel.StrokeDasharray = "5,5"
	rectModel.Fill = color
	rectModel.FillOpacity = "0.08"
	rectModel.Color = color

	// 设置排序值
	if webAnnotation.Position != -1 {
		rectModel.Sort = uint32(webAnnotation.Position)
	}

	// 处理矩形区域
	if rectAnnot.Rectangle != nil {
		rectangle := s.getRectFromOptions(rectAnnot.Rectangle)
		rectModel.Rectangles = []*pb.Rectangle{rectangle}
	}

	return rectModel, nil
}

/*
*  填充文本框模型
*  WebNoteAnnotationModel.textBox(AnnotateTextBox) --> AnnotationRawModel.textBoxModel(TextBoxModel)
 */
func (s *WebNoteProtoTransformer) fillUpTextBoxModel(webAnnotation *pb.WebNoteAnnotationModel) (*pb.TextBoxModel, error) {
	textBoxAnnot := webAnnotation.TextBox
	if textBoxAnnot == nil {
		return nil, errors.Biz("文本框不能为空")
	}

	// 创建TextBoxModel
	texBoxModel := &pb.TextBoxModel{
		AnnotateTextBox: textBoxAnnot,
		Page:            webAnnotation.PageNumber,
		NoteId:          webAnnotation.NoteId,
	}

	// 设置ID
	if textBoxAnnot.Id != "0" {
		texBoxModel.Id = textBoxAnnot.Id
	}

	return texBoxModel, nil
}

func (s *WebNoteProtoTransformer) getRectFromOptions(rect *common.RectOptions) *pb.Rectangle {
	rectangle := &pb.Rectangle{
		Height:     rect.Height,
		Width:      rect.Width,
		X:          rect.X,
		Y:          rect.Y,
		PageNumber: rect.PageNumber,
		Rotation:   rect.Rotation,
	}
	return rectangle
}

func (s *WebNoteProtoTransformer) WebAnnotationModel(annotationRawModel *pb.AnnotationRawModel) (*pb.WebNoteAnnotationModel, error) {
	return s.ToWebAnnotationModel(annotationRawModel)
}

// WebAnnotationModel 统一标注协议转Web标注协议
func (s *WebNoteProtoTransformer) ToWebAnnotationModel(annotationRawModel *pb.AnnotationRawModel) (*pb.WebNoteAnnotationModel, error) {
	if annotationRawModel == nil {
		return nil, errors.Biz("标注原始模型不能为空")
	}

	// 创建Web标注模型
	webNoteAnnotationModel := &pb.WebNoteAnnotationModel{
		Type: common.IDEAAnnotateType(annotationRawModel.Type),
	}

	// 设置高亮标志
	if annotationRawModel.IsHighlight != nil {
		webNoteAnnotationModel.IsHighlight = annotationRawModel.IsHighlight
	}

	// 根据标注类型处理
	switch common.IDEAAnnotateType(annotationRawModel.Type) {
	case common.IDEAAnnotateType_IDEAAnnotateTypeComment: // 评论类型
		model := annotationRawModel.Comment
		if model == nil {
			return nil, errors.Biz("评论类型标注缺少Comment字段")
		}

		webNoteAnnotationModel.PageNumber = model.Page

		// 填充选择模型
		selectAnnotation, err := s.fillUpWebNoteAnnotationModel(annotationRawModel)
		if err != nil {
			return nil, err
		}
		webNoteAnnotationModel.Select = selectAnnotation

		// 转换论文ID
		paperId := model.PaperId
		webNoteAnnotationModel.PaperId = paperId
		// noteID
		noteId := model.NoteId
		webNoteAnnotationModel.NoteId = noteId
		// pdfID
		webNoteAnnotationModel.PdfId = model.PdfId

		// 转换ID
		id := model.Id
		webNoteAnnotationModel.Id = id
		// 设置创建时间
		webNoteAnnotationModel.CreateDate = model.CreateDate

	case common.IDEAAnnotateType_IDEAAnnotateTypeRect: // 矩形类型
		model := annotationRawModel.Rect
		if model == nil {
			return nil, errors.Biz("矩形类型标注缺少Rect字段")
		}

		webNoteAnnotationModel.PageNumber = model.Page

		// 填充矩形模型
		rect, err := s.fillUpWebRectModel(annotationRawModel)
		if err != nil {
			return nil, err
		}
		webNoteAnnotationModel.Rect = rect

		// 转换论文ID
		paperId := model.PaperId
		webNoteAnnotationModel.PaperId = paperId

		// noteID
		noteId := model.NoteId
		webNoteAnnotationModel.NoteId = noteId
		// pdfID
		webNoteAnnotationModel.PdfId = model.PdfId

		// 转换ID
		id := model.Id
		webNoteAnnotationModel.Id = id
		// 设置创建时间
		webNoteAnnotationModel.CreateDate = model.CreateDate

	// case common.IDEAAnnotateType_IDEAAnnotateTypeIOSDraw: // 绘图类型
	// 	drawNoteRawModel := annotationRawModel.DrawNoteRawModel
	// 	if drawNoteRawModel == nil {
	// 		return nil, errors.Biz("绘图类型标注缺少DrawNoteRawModel字段")
	// 	}

	// 	drawNoteFileUrl := drawNoteRawModel.DrawNoteFileUrl

	// 	// 检查是否为iOS特殊绘图文件
	// 	if s.isIOSSpecialDrawFile(drawNoteFileUrl) {
	// 		return nil, errors.Biz("不支持的iOS特殊绘图文件")
	// 	}

	// 	webNoteAnnotationModel.PageNumber = drawNoteRawModel.Page

	// 	// 转换论文ID
	// 	paperId, err := strconv.ParseInt(drawNoteRawModel.PaperId, 10, 64)
	// 	if err == nil {
	// 		webNoteAnnotationModel.PaperId = uint64(paperId)
	// 	}

	// 	// noteID
	// 	noteId, err := strconv.ParseInt(drawNoteRawModel.NoteId, 10, 64)
	// 	if err == nil {
	// 		webNoteAnnotationModel.PaperId = uint64(noteId)
	// 	}

	// 	// 转换ID
	// 	id, err := strconv.ParseInt(drawNoteRawModel.Id, 10, 64)
	// 	if err == nil {
	// 		webNoteAnnotationModel.Id = uint64(id)
	// 	}

	// 	// 创建绘图模型
	// 	webNoteAnnotationModel.DrawAnnotation = &pb.WebNoteAnnotateDraw{
	// 		NsDataUrl: drawNoteFileUrl,
	// 	}

	case common.IDEAAnnotateType_IDEAAnnotateTypeTextBox: // 文本框类型
		textBoxModel := annotationRawModel.TextBoxModel
		if textBoxModel == nil {
			return nil, errors.Biz("文本框类型标注缺少TextBoxModel字段")
		}

		// 处理文本框模型
		if textBoxModel.AnnotateTextBox != nil {
			textBoxModel.AnnotateTextBox.Id = textBoxModel.Id
			webNoteAnnotationModel.TextBox = textBoxModel.AnnotateTextBox
			webNoteAnnotationModel.PageNumber = textBoxModel.Page
		}

		// noteID
		noteId := textBoxModel.NoteId
		webNoteAnnotationModel.NoteId = noteId

		// 转换ID
		webNoteAnnotationModel.Id = textBoxModel.Id

	default:
		return nil, errors.Biz("不支持的标注类型")
	}

	return webNoteAnnotationModel, nil
}

// fillUpWebNoteAnnotationModelFromRaw 从原始标注模型填充Web选择模型
func (t *WebNoteProtoTransformer) fillUpWebNoteAnnotationModel(annotationRawModel *pb.AnnotationRawModel) (*pb.WebNoteSelectAnnotation, error) {
	if annotationRawModel == nil {
		return nil, errors.Biz("标注原始模型不能为空")
	}

	comment := annotationRawModel.Comment
	if comment == nil {
		return nil, errors.Biz("评论类型标注缺少Comment字段")
	}

	// 创建Web选择模型
	selectAnnotation := &pb.WebNoteSelectAnnotation{
		Uuid:       comment.Id,
		DocumentId: comment.NoteId,
		Idea:       comment.Idea,
		StyleId:    comment.StyleId,
	}

	// 处理矩形区域
	if comment.RectStr != "" {
		selectAnnotation.RectStr = comment.RectStr
	}

	// 补充选区信息，以Web形式输出
	if comment.Rectangles != nil {
		webRectOptions := make([]*common.RectOptions, 0)
		for _, rectangle := range comment.Rectangles {
			// 创建RectOptions
			rectOption := &common.RectOptions{
				X:          rectangle.X,
				Y:          rectangle.Y,
				Width:      rectangle.Width,
				Height:     rectangle.Height,
				PageNumber: uint32(comment.Page),
				Rotation:   rectangle.Rotation,
			}
			webRectOptions = append(webRectOptions, rectOption)
		}
		selectAnnotation.Rectangle = webRectOptions
	}

	return selectAnnotation, nil
}

// fillUpWebRectModel 从原始标注模型填充Web矩形模型
func (t *WebNoteProtoTransformer) fillUpWebRectModel(annotationRawModel *pb.AnnotationRawModel) (*pb.WebNoteRectAnnotation, error) {
	if annotationRawModel == nil {
		return nil, errors.Biz("标注原始模型不能为空")
	}

	rect := annotationRawModel.Rect
	if rect == nil {
		return nil, errors.Biz("矩形类型标注缺少Rect字段")
	}

	// 创建Web矩形模型
	webRect := &pb.WebNoteRectAnnotation{
		Uuid:       rect.Id,
		DocumentId: rect.NoteId,
		PicUrl:     rect.PicUrl,
		Idea:       rect.Idea,
		StyleId:    rect.StyleId,
	}

	// 补充选区信息，以Web格式输出
	if rect.Rectangles == nil || len(rect.Rectangles) != 1 {
		return nil, errors.Biz("DB标注数据选区信息异常，id:" + rect.Id)
	}

	// 创建RectOptions
	rectangle := rect.Rectangles[0]
	rectOption := &common.RectOptions{
		X:          rectangle.X,
		Y:          rectangle.Y,
		Width:      rectangle.Width,
		Height:     rectangle.Height,
		PageNumber: uint32(rect.Page),
		Rotation:   rectangle.Rotation,
	}
	webRect.Rectangle = rectOption

	return webRect, nil
}

// isIOSSpecialDrawFile 检查是否为iOS特殊绘图文件
// func (t *WebNoteProtoTransformer) isIOSSpecialDrawFile(urlStr string) bool {
// 	// 如果URL为空，直接返回false
// 	if urlStr == "" {
// 		return false
// 	}

// 	// 解析URL
// 	parsedURL, err := url.Parse(urlStr)
// 	if err != nil {
// 		// URL解析失败，直接返回false
// 		return false
// 	}

// 	// 获取路径
// 	path := parsedURL.Path
// 	// 获取文件名
// 	slashIndex := strings.LastIndex(path, "/")
// 	var fileName string
// 	if slashIndex >= 0 && slashIndex < len(path)-1 {
// 		fileName = path[slashIndex+1:]
// 	} else {
// 		fileName = path
// 	}

// 	// 检查文件名是否以.drawdata结尾
// 	return strings.HasSuffix(fileName, ".drawdata")
// }
