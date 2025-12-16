package proto

import (
	"encoding/json"
	"strconv"

	"github.com/yb2020/odoc/pkg/errors"
	common "github.com/yb2020/odoc/proto/gen/go/common"
	pb "github.com/yb2020/odoc/proto/gen/go/note"

	pdfModel "github.com/yb2020/odoc/services/pdf/model"
)

// AnnotationModelTool 注释模型工具类
type AnnotationModelTool struct {
}

// GetNoteIdFromAnnotationRawModel 从注释原始模型中获取笔记ID
// 返回笔记ID和是否存在的标志
func (t *AnnotationModelTool) GetNoteIdFromAnnotationRawModel(annotationRawModel *pb.AnnotationRawModel) (string, bool) {
	if annotationRawModel == nil {
		return "0", false
	}

	switch {
	case annotationRawModel.Comment != nil:
		// 从评论类型中获取笔记ID
		noteId := annotationRawModel.Comment.NoteId
		if noteId == "" {
			return "0", false
		}
		return noteId, true
	case annotationRawModel.Rect != nil:
		// 从矩形类型中获取笔记ID
		noteId := annotationRawModel.Rect.NoteId
		if noteId == "" {
			return "0", false
		}
		return noteId, true
	case annotationRawModel.TextBoxModel != nil:
		// 从文本框类型中获取笔记ID
		noteId := annotationRawModel.TextBoxModel.NoteId
		if noteId == "" {
			return "0", false
		}
		return noteId, true
	default:
		return "0", false
	}
}

// GetPageNumberFromAnnotationRawModel 从注释原始模型中获取页码
// 返回页码和是否存在的标志
func (t *AnnotationModelTool) GetPageNumberFromAnnotationRawModel(annotationRawModel *pb.AnnotationRawModel) (uint32, bool) {
	if annotationRawModel == nil {
		return 0, false
	}

	switch {
	case annotationRawModel.Comment != nil:
		// 从评论类型中获取页码
		page := annotationRawModel.Comment.Page
		if page <= 0 {
			return 0, false
		}
		return page, true
	case annotationRawModel.Rect != nil:
		// 从矩形类型中获取页码
		page := annotationRawModel.Rect.Page
		if page <= 0 {
			return 0, false
		}
		return page, true
	case annotationRawModel.DrawNoteRawModel != nil:
		// 从绘图笔记类型中获取页码
		page := annotationRawModel.DrawNoteRawModel.Page
		if page <= 0 {
			return 0, false
		}
		return page, true
	case annotationRawModel.TextBoxModel != nil:
		// 从文本框类型中获取页码
		page := annotationRawModel.TextBoxModel.Page
		if page <= 0 {
			return 0, false
		}
		return page, true
	default:
		return 0, false
	}
}

// GetSortOfAnnotationRawModel 从注释原始模型中获取排序值
// 如果无法获取排序值，则返回0
func (t *AnnotationModelTool) GetSortOfAnnotationRawModel(annotationRawModel *pb.AnnotationRawModel) uint32 {
	if annotationRawModel == nil {
		return 0
	}

	switch {
	case annotationRawModel.Comment != nil:
		// 从评论类型中获取排序值
		return annotationRawModel.Comment.Sort
	case annotationRawModel.Rect != nil:
		// 从矩形类型中获取排序值
		return annotationRawModel.Rect.Sort
	case annotationRawModel.DrawNoteRawModel != nil, annotationRawModel.TextBoxModel != nil:
		// 绘图笔记和文本框类型默认返回0
		return 0
	default:
		return 0
	}
}

// SetSortOfAnnotationRawModel 设置注释原始模型的排序值
func (t *AnnotationModelTool) SetSortOfAnnotationRawModel(annotationRawModel *pb.AnnotationRawModel, sort uint32) {
	if annotationRawModel == nil {
		return
	}

	switch {
	case annotationRawModel.Comment != nil:
		// 设置评论类型的排序值
		annotationRawModel.Comment.Sort = sort
	case annotationRawModel.Rect != nil:
		// 设置矩形类型的排序值
		annotationRawModel.Rect.Sort = sort
		// 注意：绘图笔记和文本框类型不设置排序值，与Java版本保持一致
	}
}

// PdfMark --> AnnotationRawModel
func (t *AnnotationModelTool) ToAnnotationRawModel(pdfMark *pdfModel.PdfMark) (*pb.AnnotationRawModel, error) {
	if pdfMark == nil {
		return nil, errors.Biz("PDF标记不能为空")
	}

	// 根据标注类型处理
	switch pdfMark.Type {
	case int(common.IDEAAnnotateType_IDEAAnnotateTypeComment): // 评论类型

		annotationRawModel := &pb.AnnotationRawModel{
			Type:        uint32(pdfMark.Type),
			IsHighlight: &pdfMark.IsHighlight,
			UpdateTime:  uint64(pdfMark.UpdatedAt.Unix()),
			UserId:      pdfMark.CreatorId,
		}
		// 尝试从CommentContent字段解析出更多信息
		if pdfMark.CommentContent != "" {
			var commentRaw *pb.CommentRawModel
			if err := json.Unmarshal([]byte(pdfMark.CommentContent), &commentRaw); err == nil {

				commentRaw.Id = pdfMark.Id
				commentRaw.NoteId = pdfMark.NoteId
				commentRaw.PaperId = pdfMark.PaperId
				commentRaw.PdfId = pdfMark.PdfId

				// 设置Idea字段
				commentRaw.Idea = pdfMark.Idea

				// 设置Fill字段，使用NoteAnnotationStyleTransformer获取颜色
				styleTransformer := &NoteAnnotationStyleTransformer{}
				commentRaw.Fill = styleTransformer.GetColorByStyleId(uint32(pdfMark.StyleId))
				annotationRawModel.Comment = commentRaw
			}
		}

		return annotationRawModel, nil

	case int(common.IDEAAnnotateType_IDEAAnnotateTypeRect): // 矩形类型

		annotationRawModel := &pb.AnnotationRawModel{
			Type:        uint32(pdfMark.Type),
			IsHighlight: &pdfMark.IsHighlight,
			UpdateTime:  uint64(pdfMark.UpdatedAt.Unix()),
			UserId:      pdfMark.CreatorId,
		}
		// 尝试从Content字段解析出更多信息
		if pdfMark.RectContent != "" {
			var rectRaw *pb.RectRawModel
			if err := json.Unmarshal([]byte(pdfMark.RectContent), &rectRaw); err == nil {

				rectRaw.Id = pdfMark.Id
				rectRaw.NoteId = pdfMark.NoteId
				rectRaw.PaperId = pdfMark.PaperId
				rectRaw.PdfId = pdfMark.PdfId

				// 设置Idea字段
				rectRaw.Idea = pdfMark.Idea

				// 设置Fill字段，使用NoteAnnotationStyleTransformer获取颜色
				styleTransformer := &NoteAnnotationStyleTransformer{}
				rectRaw.Fill = styleTransformer.GetColorByStyleId(uint32(pdfMark.StyleId))

				rectRaw.PicUrl = pdfMark.PicUrl

				annotationRawModel.Rect = rectRaw
			}
		}

		return annotationRawModel, nil

	case int(common.IDEAAnnotateType_IDEAAnnotateTypeTextBox): // 文本框类型
		annotationRawModel := &pb.AnnotationRawModel{
			Type:        uint32(pdfMark.Type),
			IsHighlight: &pdfMark.IsHighlight,
			UpdateTime:  uint64(pdfMark.UpdatedAt.Unix()),
			UserId:      pdfMark.CreatorId,
		}
		// 尝试从CommentContent字段解析出更多信息
		if pdfMark.TextBoxContent != "" {
			var textBoxRaw *pb.TextBoxModel
			if err := json.Unmarshal([]byte(pdfMark.TextBoxContent), &textBoxRaw); err == nil {

				textBoxRaw.Id = pdfMark.Id
				textBoxRaw.NoteId = pdfMark.NoteId

				annotationRawModel.TextBoxModel = textBoxRaw
			}
		}

		return annotationRawModel, nil
	default:
		return nil, errors.Biz("不支持的标注类型: " + strconv.Itoa(pdfMark.Type))
	}

	// return nil, nil
}

// AnnotationRawModel --> PdfMark
func (s *AnnotationModelTool) ToPdfMark(annotationRawModel *pb.AnnotationRawModel) (*pdfModel.PdfMark, error) {
	if annotationRawModel == nil {
		return nil, errors.Biz("标注原始模型不能为空")
	}

	// 根据标注类型处理
	switch annotationRawModel.Type {
	case uint32(common.IDEAAnnotateType_IDEAAnnotateTypeComment): // 注释
		return s.toPdfMarkOfComment(annotationRawModel)

	case uint32(common.IDEAAnnotateType_IDEAAnnotateTypeRect): // 矩形
		return s.toPdfMarkOfRect(annotationRawModel)

	// case uint32(common.IDEAAnnotateType_IDEAAnnotateTypeIOSDraw): //暂无此类型
	// 	return s.toPdfMarkOfDrawNote(annotationRawModel)

	case uint32(common.IDEAAnnotateType_IDEAAnnotateTypeTextBox): //文本框
		return s.toPdfMarkOfTextBox(annotationRawModel)

	default:
		return nil, errors.Biz("不支持的标注类型")
	}
}

// toPdfMarkOfComment 抽取评论类型标注
func (t *AnnotationModelTool) toPdfMarkOfComment(annotationRawModel *pb.AnnotationRawModel) (*pdfModel.PdfMark, error) {
	comment := annotationRawModel.Comment
	if comment == nil {
		return nil, errors.Biz("评论类型标注缺少Comment字段")
	}

	// 序列化标注模型
	commentContent, _ := json.Marshal(comment)

	// 创建保存请求
	pdfMark := &pdfModel.PdfMark{
		Idea:           comment.Idea,
		PaperId:        comment.PaperId,
		NoteId:         comment.NoteId,
		PdfId:          comment.PdfId,
		CommentContent: string(commentContent),
		KeyContent:     comment.RectStr,
		StyleId:        int(comment.StyleId),
		Sort:           int(comment.Sort),
		IsHighlight:    false,
		Page:           int(comment.Page),
	}

	if comment.Id != "" {
		pdfMark.Id = comment.Id
	}
	// 设置高亮标志
	if annotationRawModel.IsHighlight != nil {
		pdfMark.IsHighlight = *annotationRawModel.IsHighlight
	}

	// 注意：在Go中，我们通常通过上下文传递用户ID，而不是通过静态方法获取
	// 这里假设用户ID会在后续的处理中设置

	return pdfMark, nil
}

// toPdfMarkOfRect 抽取矩形类型标注
func (t *AnnotationModelTool) toPdfMarkOfRect(annotationRawModel *pb.AnnotationRawModel) (*pdfModel.PdfMark, error) {
	rect := annotationRawModel.Rect
	if rect == nil {
		return nil, errors.Biz("矩形类型标注缺少Rect字段")
	}

	// 序列化标注模型
	recContent, _ := json.Marshal(rect)

	// 创建保存请求
	pdfMark := &pdfModel.PdfMark{
		PaperId:     rect.PaperId,
		NoteId:      rect.NoteId,
		PdfId:       rect.PdfId,
		PicUrl:      rect.PicUrl,
		Idea:        rect.Idea,
		RectContent: string(recContent),
		StyleId:     int(rect.StyleId),
		Sort:        int(rect.Sort),
		IsHighlight: false,
		Type:        int(common.IDEAAnnotateType_IDEAAnnotateTypeRect),
		Page:        int(rect.Page),
	}
	if rect.Id != "" {
		pdfMark.Id = rect.Id
	}

	// 设置高亮标志
	if annotationRawModel.IsHighlight != nil {
		pdfMark.IsHighlight = *annotationRawModel.IsHighlight
	}

	// 注意：在Go中，我们通常通过上下文传递用户ID，而不是通过静态方法获取
	// 这里假设用户ID会在后续的处理中设置

	return pdfMark, nil
}

// toPdfMarkOfDrawNote 抽取绘图类型标注
// func (t *AnnotationModelTool) toPdfMarkOfDrawNote(annotationRawModel *pb.AnnotationRawModel) (*pdfModel.PdfMark, error) {
// 	if annotationRawModel == nil {
// 		return nil, errors.Biz("标注原始模型不能为空")
// 	}

// 	drawNoteRawModel := annotationRawModel.DrawNoteRawModel
// 	if drawNoteRawModel == nil {
// 		return nil, errors.Biz("绘图类型标注缺少DrawNoteRawModel字段")
// 	}

// 	// 转换ID字段
// 	paperId, _ := strconv.ParseInt(drawNoteRawModel.PaperId, 10, 64)
// 	noteId, _ := strconv.ParseInt(drawNoteRawModel.NoteId, 10, 64)
// 	id, _ := strconv.ParseInt(drawNoteRawModel.Id, 10, 64)

// 	// 序列化标注模型
// 	content, _ := json.Marshal(annotationRawModel)

// 	// 创建保存请求
// 	pdfMark := &pdfModel.PdfMark{
// 		PaperId:     paperId,
// 		NoteId:      noteId,
// 		Content:     string(content),
// 		Page:        int(drawNoteRawModel.Page),
// 		IsHighlight: false,
// 		Type:        int(common.IDEAAnnotateType_IDEAAnnotateTypeIOSDraw),
// 	}
// 	if id != 0 {
// 		pdfMark.Id = id
// 	}

// 	// 设置高亮标志
// 	if annotationRawModel.IsHighlight != nil {
// 		pdfMark.IsHighlight = *annotationRawModel.IsHighlight
// 	}

// 	// 注意：在Go中，我们通常通过上下文传递用户ID，而不是通过静态方法获取
// 	// 这里假设用户ID会在后续的处理中设置

// 	return pdfMark, nil
// }

// toPdfMarkOfTextBox 抽取文本框类型标注
func (t *AnnotationModelTool) toPdfMarkOfTextBox(annotationRawModel *pb.AnnotationRawModel) (*pdfModel.PdfMark, error) {
	if annotationRawModel == nil {
		return nil, errors.Biz("标注原始模型不能为空")
	}

	textBoxModel := annotationRawModel.TextBoxModel
	if textBoxModel == nil {
		return nil, errors.Biz("文本框类型标注缺少TextBoxModel字段")
	}

	// 序列化标注模型
	textBoxcontent, _ := json.Marshal(textBoxModel)

	// 创建保存请求
	pdfMark := &pdfModel.PdfMark{
		NoteId:         textBoxModel.NoteId,
		TextBoxContent: string(textBoxcontent),
		Page:           int(textBoxModel.Page),
		IsHighlight:    false,
		Type:           int(common.IDEAAnnotateType_IDEAAnnotateTypeTextBox),
	}

	if textBoxModel.Id != "" {
		pdfMark.Id = textBoxModel.Id
	}

	// 设置高亮标志
	if annotationRawModel.IsHighlight != nil {
		pdfMark.IsHighlight = *annotationRawModel.IsHighlight
	}

	// 注意：在Go中，我们通常通过上下文传递用户ID，而不是通过静态方法获取
	// 这里假设用户ID会在后续的处理中设置

	return pdfMark, nil
}
