package factory

import (
	"errors"

	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/interfaces"
)

// DocMetaInfoHandlerServiceFactory 文档元信息处理器工厂，用于创建和管理不同的文档元信息处理器
type DocMetaInfoHandlerServiceFactory struct {
	providers map[string]interfaces.IDocMetaInfoHandlerService
	logger    logging.Logger
}

// NewDocMetaInfoHandlerServiceFactory 创建一个新的文档元信息处理器工厂
func NewDocMetaInfoHandlerServiceFactory(logger logging.Logger) *DocMetaInfoHandlerServiceFactory {
	return &DocMetaInfoHandlerServiceFactory{
		providers: make(map[string]interfaces.IDocMetaInfoHandlerService),
		logger:    logger,
	}
}

// Register 注册一个文档元信息处理器
func (f *DocMetaInfoHandlerServiceFactory) Register(provider interfaces.IDocMetaInfoHandlerService) {
	f.providers[provider.GetName()] = provider
}

// GetProvider 获取一个合适的文档元信息处理器
func (f *DocMetaInfoHandlerServiceFactory) GetProvider(typeName string) (interfaces.IDocMetaInfoHandlerService, error) {
	for _, provider := range f.providers {
		if provider.GetName() == typeName {
			return provider, nil
		}
	}
	f.logger.Error("msg", "no provider found for type: ", "type", typeName)
	return nil, errors.New("no provider found for type: " + typeName)
}
