package util

import (
	pb "github.com/yb2020/odoc-proto/gen/go/parsed"
	v8 "github.com/yb2020/odoc/services/parse/util/grobid/v8"
)

// TODO：这里的封装方式不合理，由于时间原因，无法使用ai生成合适的代码，暂时先这样
// 支持的解析器版本
const (
	GrobidParserV8 = "v8"
	// 默认版本
	DefaultGrobidParser = GrobidParserV8
)

// GrobidParser 定义了GROBID解析器的通用接口
type GrobidParser interface {
	// 解析文档头部信息
	ParseDocumentHeaderXml(xmlContent []byte) (*pb.DocumentHeader, error)
	// 解析完整文档
	ParseDocument(xmlContent []byte) (*pb.DocumentMetadata, *pb.FullDocument, error)
	// 获取解析器版本
	Version() string
}

// NewGrobidParser 创建指定版本的GROBID解析器
func NewGrobidParser(version string) GrobidParser {
	// 如果未指定版本，使用默认版本
	if version == "" {
		version = DefaultGrobidParser
	}
	// 创建相应版本的解析器
	switch version {
	case GrobidParserV8:
		return newGrobidParserV8()
	}
	return nil
}

// V8版本解析器实现
type grobidParserV8 struct{}

func newGrobidParserV8() GrobidParser {
	return &grobidParserV8{}
}

func (p *grobidParserV8) Version() string {
	return GrobidParserV8
}

func (p *grobidParserV8) ParseDocumentHeaderXml(xmlContent []byte) (*pb.DocumentHeader, error) {
	// 调用V8版本的解析函数
	return v8.ParseDocumentHeaderXml(xmlContent)
}

func (p *grobidParserV8) ParseDocument(xmlContent []byte) (*pb.DocumentMetadata, *pb.FullDocument, error) {
	// 调用V8版本的解析函数
	return v8.ParseDocument(xmlContent)
}
