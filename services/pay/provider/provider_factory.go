package provider

import (
	"fmt"

	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pay/model"
)

// PaymentProviderFactory 支付提供商工厂，用于创建和管理不同的支付提供商
type PaymentProviderFactory struct {
	providers map[string]PaymentProvider
	logger    logging.Logger
}

// NewPaymentProviderFactory 创建一个新的支付提供商工厂
func NewPaymentProviderFactory(logger logging.Logger) *PaymentProviderFactory {
	return &PaymentProviderFactory{
		providers: make(map[string]PaymentProvider),
		logger:    logger,
	}
}

// RegisterStripeProvider 注册Stripe支付提供商
func (f *PaymentProviderFactory) RegisterStripeProvider(apiKey, webhookSecret string) {
	f.providers[model.PaymentChannelStripe] = NewStripeProvider(apiKey, webhookSecret, f.logger)
	f.logger.Info("msg", "Stripe支付提供商已注册")
}

// RegisterProvider 注册自定义支付提供商
func (f *PaymentProviderFactory) RegisterProvider(provider PaymentProvider) {
	providerName := provider.GetName()
	f.providers[providerName] = provider
	f.logger.Info("msg", "支付提供商已注册", "provider", providerName)
}

// GetProvider 获取指定名称的支付提供商
func (f *PaymentProviderFactory) GetProvider(name string) (PaymentProvider, error) {
	provider, exists := f.providers[name]
	if !exists {
		return nil, fmt.Errorf("支付提供商 '%s' 未注册", name)
	}
	return provider, nil
}

// GetAllProviders 获取所有已注册的支付提供商
func (f *PaymentProviderFactory) GetAllProviders() map[string]PaymentProvider {
	return f.providers
}
