package util

import (
	"encoding/json"
	"fmt"

	"github.com/yb2020/odoc/pkg/errors"
	pb "github.com/yb2020/odoc/proto/gen/go/parsed"
)

// DeserializePdfMetadata 将字节数组反序列化为PDF元数据对象
func DeserializePdfMetadata(data []byte) (*pb.DocumentMetadata, error) {
	var metadata pb.DocumentMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, errors.BizWrap("反序列化PDF元数据失败", err)
	}
	return &metadata, nil
}

// DeserializePdfParagraphs 将字节数组反序列化为PDF段落列表
func DeserializePdfParagraphs(data []byte) (*pb.FullDocument, error) {
	fmt.Println("data", string(data))
	var paragraphs struct {
		Paragraphs []*pb.Paragraph `json:"paragraphs"`
	}
	if err := json.Unmarshal(data, &paragraphs); err != nil {
		return nil, errors.BizWrap("反序列化PDF段落失败", err)
	}
	return &pb.FullDocument{
		Paragraphs: paragraphs.Paragraphs,
	}, nil
}
