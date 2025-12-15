package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	pkgi18n "github.com/yb2020/odoc/pkg/i18n"
)

// WebhookConfig 保持不变
type WebhookConfig struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	AuthToken string `json:"auth_token" yaml:"auth_token"`
}

// S3Config S3配置
// S3Config 结构体现在使用 map 来处理 buckets
type S3Config struct {
	Region          string `json:"region" yaml:"region"`
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	AccessKeyID     string `json:"access_key_id" yaml:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key"`
	UseSSL          bool   `json:"use_ssl" yaml:"use_ssl"`
	ForcePathStyle  bool   `json:"force_path_style" yaml:"force_path_style"`
	Location        string `json:"location" yaml:"location"` // 确保这个字段存在，因为您的代码中使用了它

	// 【核心修改】将固定的结构体改为 map
	// map 的键将是 "public", "pdf", "temp" 等逻辑名称
	Buckets map[string]BucketConfig `json:"buckets" yaml:"buckets"`

	Upload  UploadConfig  `json:"upload" yaml:"upload"`
	Webhook WebhookConfig `json:"webhook" yaml:"webhook"`
}

// BucketConfig 定义了单个存储桶的配置属性
// 这个结构体直接映射 YAML 文件中的字段
type BucketConfig struct {
	Name          string `json:"name" yaml:"name"`
	Public        bool   `json:"public" yaml:"public"`
	Versioning    bool   `json:"versioning" yaml:"versioning"`
	LifecycleDays int    `json:"lifecycle_days" yaml:"lifecycle_days"` // 0 表示不过期
}

// UploadConfig 保持不变
type UploadConfig struct {
	PublicUploadEndpoint   string `json:"public_upload_endpoint" yaml:"public_upload_endpoint"`
	PublicDownloadEndpoint string `json:"public_download_endpoint" yaml:"public_download_endpoint"`
	PresignedURLExpires    int    `json:"presigned_url_expires" yaml:"presigned_url_expires"`
	MaxFileSize            int64  `json:"max_file_size" yaml:"max_file_size"`
}

// S3Provider S3提供者配置
type S3Provider struct {
	Config S3Config `json:"config" yaml:"config"` // S3配置
}

type OSSConfig struct {
	S3    S3Config    `json:"s3" yaml:"s3"`
	Local LocalConfig `json:"local" yaml:"local"`
}

type ServiceConfig struct {
	Type string `json:"type" yaml:"type"`
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	BasePath string `json:"basePath" yaml:"basePath"` // 本地存储基础路径
}

// WebsiteInitData defines the structure for initial website data.
type WebsiteInitData struct {
	Name    string `json:"name" yaml:"name"`
	URL     string `json:"url" yaml:"url"`
	IconURL string `json:"iconUrl" yaml:"iconUrl"`
}

// WebsiteConfig holds the configuration for the website.
type WebsiteConfig struct {
	InitData []WebsiteInitData `json:"initData" yaml:"initData"`
}

// NavConfig holds all navigation-related configurations.
type NavConfig struct {
	Website WebsiteConfig `json:"website" yaml:"website"`
}

// Config holds all configuration for our application
type Config struct {
	Server struct {
		Port    int    `json:"port" yaml:"port"`
		Host    string `json:"host" yaml:"host"`
		Timeout int    `json:"timeout" yaml:"timeout"` // in seconds
		GRPC    struct {
			Port int    `json:"port" yaml:"port"`
			Host string `json:"host" yaml:"host"`
		} `json:"grpc" yaml:"grpc"`
	} `json:"server" yaml:"server"`

	Logging struct {
		Level      string `json:"level" yaml:"level"`
		Format     string `json:"format" yaml:"format"`         // 支持 "json", "logfmt", "springboot"
		Path       string `json:"path" yaml:"path"`             // 日志文件路径
		MaxSize    int    `json:"maxSize" yaml:"maxSize"`       // 单个日志文件最大大小（MB）
		MaxAge     int    `json:"maxAge" yaml:"maxAge"`         // 日志文件保留天数
		MaxBackups int    `json:"maxBackups" yaml:"maxBackups"` // 保留的旧日志文件数量
	} `json:"logging" yaml:"logging"`

	Tracing struct {
		Enabled     bool    `json:"enabled" yaml:"enabled"`
		ServiceName string  `json:"serviceName" yaml:"serviceName"`
		JaegerURL   string  `json:"jaegerUrl" yaml:"jaegerUrl"`
		SampleRate  float64 `json:"sampleRate" yaml:"sampleRate"`
	} `json:"tracing" yaml:"tracing"`

	Metrics struct {
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Path    string `json:"path" yaml:"path"`
	} `json:"metrics" yaml:"metrics"`

	ErrorNotification struct {
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Type    string `json:"type" yaml:"type"` // email, slack, etc.
		URL     string `json:"url" yaml:"url"`
	} `json:"errorNotification" yaml:"errorNotification"`

	Database struct {
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Type    string `json:"type" yaml:"type"` // postgres, sqlite

		// PostgreSQL 配置
		Postgres struct {
			Host            string `json:"host" yaml:"host"`
			Port            int    `json:"port" yaml:"port"`
			User            string `json:"user" yaml:"user"`
			Password        string `json:"password" yaml:"password"`
			DBName          string `json:"dbname" yaml:"dbname"`
			SSLMode         string `json:"sslmode" yaml:"sslmode"`
			TimeZone        string `json:"timezone" yaml:"timezone"`
			MaxIdleConns    int    `json:"maxIdleConns" yaml:"maxIdleConns"`
			MaxOpenConns    int    `json:"maxOpenConns" yaml:"maxOpenConns"`
			ConnMaxLifetime int    `json:"connMaxLifetime" yaml:"connMaxLifetime"` // 秒
			LogLevel        string `json:"logLevel" yaml:"logLevel"`
		} `json:"postgres" yaml:"postgres"`

		// SQLite 配置
		SQLite struct {
			DBPath          string `json:"dbPath" yaml:"dbPath"` // 数据库文件路径
			MaxIdleConns    int    `json:"maxIdleConns" yaml:"maxIdleConns"`
			MaxOpenConns    int    `json:"maxOpenConns" yaml:"maxOpenConns"`
			ConnMaxLifetime int    `json:"connMaxLifetime" yaml:"connMaxLifetime"` // 秒
			LogLevel        string `json:"logLevel" yaml:"logLevel"`
		} `json:"sqlite" yaml:"sqlite"`
	} `json:"database" yaml:"database"`

	OSS OSSConfig `json:"oss" yaml:"oss"` // 对象存储配置

	Service ServiceConfig `json:"service" yaml:"service"` // 服务配置

	// RocketMQ配置
	RocketMQ RocketMQConfig `json:"rocketmq" yaml:"rocketmq"`

	Redis struct {
		Enabled         bool   `json:"enabled" yaml:"enabled"`
		Host            string `json:"host" yaml:"host"`
		Port            int    `json:"port" yaml:"port"`
		Password        string `json:"password" yaml:"password"`
		DB              int    `json:"db" yaml:"db"`
		PoolSize        int    `json:"poolSize" yaml:"poolSize"`
		MinIdleConns    int    `json:"minIdleConns" yaml:"minIdleConns"`
		DialTimeout     int    `json:"dialTimeout" yaml:"dialTimeout"`   // 秒
		ReadTimeout     int    `json:"readTimeout" yaml:"readTimeout"`   // 秒
		WriteTimeout    int    `json:"writeTimeout" yaml:"writeTimeout"` // 秒
		MaxConnAge      int    `json:"maxConnAge" yaml:"maxConnAge"`     // 秒
		MaxRetries      int    `json:"maxRetries" yaml:"maxRetries"`
		MinRetryBackoff int    `json:"minRetryBackoff" yaml:"minRetryBackoff"` // 毫秒
		MaxRetryBackoff int    `json:"maxRetryBackoff" yaml:"maxRetryBackoff"` // 毫秒
	} `json:"redis" yaml:"redis"`

	// Cache 缓存配置
	Cache struct {
		Type       string `json:"type" yaml:"type"`             // redis 或 memory
		Expiration int    `json:"expiration" yaml:"expiration"` // 默认过期时间（秒）
	} `json:"cache" yaml:"cache"`

	Translate struct {
		Text struct {
			MaxLength int `json:"maxLength" yaml:"maxLength"`
			Special   struct {
				Words struct {
					NeedReplace struct {
						List []string `json:"list" yaml:"list"`
					} `json:"needreplace" yaml:"needreplace"`
				} `json:"words" yaml:"words"`
			} `json:"special" yaml:"special"`
			Channel struct {
				Youdao struct {
					AppId        string `json:"app_id" yaml:"appId"`
					AppSecretKey string `json:"app_secret_key" yaml:"appSecretKey"`
					ApiUrl       string `json:"api_url" yaml:"apiUrl"`
				} `json:"youdao" yaml:"youdao"`
				Google struct {
					BaseURL       string `json:"baseurl" yaml:"baseurl"`
					FreeUri       string `json:"freeUri" yaml:"freeUri"`
					SimpleFreeUri string `json:"simpleFreeUri" yaml:"simpleFreeUri"`
					Token         string `json:"token" yaml:"token"`
				} `json:"google" yaml:"google"`
			} `json:"channel" yaml:"channel"`
			DefaultSourceLanguage string `json:"defaultSourceLanguage" yaml:"defaultSourceLanguage"`
			DefaultTargetLanguage string `json:"defaultTargetLanguage" yaml:"defaultTargetLanguage"`
		} `json:"text" yaml:"text"`
		Pronunciation struct {
			Crawl struct {
				BD struct {
					URL string `json:"url" yaml:"url"`
				} `json:"bd" yaml:"bd"`
				YD struct {
					Result struct {
						URL string `json:"url" yaml:"url"`
					} `json:"result" yaml:"result"`
					Pronounce struct {
						URL string `json:"url" yaml:"url"`
					} `json:"pronounce" yaml:"pronounce"`
				} `json:"yd" yaml:"yd"`
			} `json:"crawl" yaml:"crawl"`
		} `json:"pronunciation" yaml:"pronunciation"`
		OCR struct {
			ExtractTextURL string  `json:"extractTextURL" yaml:"extractTextURL"`
			UploadImageURL string  `json:"uploadImageURL" yaml:"uploadImageURL"`
			ProgressURL    string  `json:"progressURL" yaml:"progressURL"`
			JoinURL        string  `json:"joinURL" yaml:"joinURL"`
			DataURL        string  `json:"dataURL" yaml:"dataURL"`
			FileURL        string  `json:"fileURL" yaml:"fileURL"`
			Timeout        int     `json:"timeout" yaml:"timeout"`
			ReturnWordBox  string  `json:"returnWordBox" yaml:"returnWordBox"`
			TextScore      float64 `json:"textScore" yaml:"textScore"`
			BoxThresh      float64 `json:"boxThresh" yaml:"boxThresh"`
			UnclipRatio    float64 `json:"unclipRatio" yaml:"unclipRatio"`
			MaxSideLen     int     `json:"maxSideLen" yaml:"maxSideLen"`
			UseDet         struct {
				EngineType string `json:"engineType" yaml:"engineType"`
				LangDet    string `json:"langDet" yaml:"langDet"`
				ModelType  string `json:"modelType" yaml:"modelType"`
				Version    string `json:"version" yaml:"version"`
			} `json:"useDet" yaml:"useDet"`
			UseCls struct {
				EngineType string `json:"engineType" yaml:"engineType"`
				LangCls    string `json:"langCls" yaml:"langCls"`
				ModelType  string `json:"modelType" yaml:"modelType"`
				Version    string `json:"version" yaml:"version"`
			} `json:"useCls" yaml:"useCls"`
			UseRec struct {
				EngineType string `json:"engineType" yaml:"engineType"`
				LangRec    string `json:"langRec" yaml:"langRec"`
				ModelType  string `json:"modelType" yaml:"modelType"`
				Version    string `json:"version" yaml:"version"`
			} `json:"useRec" yaml:"useRec"`
		} `json:"ocr" yaml:"ocr"`
		Glossary struct {
			MaxSize   int `json:"maxSize" yaml:"maxSize"`
			MaxLength int `json:"maxLength" yaml:"maxLength"`
		} `json:"glossary" yaml:"glossary"`
		FullTextTranslate struct {
			TranslateBaseURL       string   `json:"translateBaseURL" yaml:"translateBaseURL"`
			TranslateServiceURI    string   `json:"translateServiceURI" yaml:"translateServiceURI"`
			TranslateServiceFixURI string   `json:"translateServiceFixURI" yaml:"translateServiceFixURI"`
			UseWatermark           bool     `json:"useWatermark" yaml:"useWatermark"`
			HistoryVisibleDays     int      `json:"historyVisibleDays" yaml:"historyVisibleDays"`
			TranslateTimeOut       int      `json:"translateTimeOut" yaml:"translateTimeOut"`
			DuplicateDelaySeconds  int      `json:"duplicateDelaySeconds" yaml:"duplicateDelaySeconds"`
			FileSizeLimit          int      `json:"fileSizeLimit" yaml:"fileSizeLimit"`
			FilePageLimit          int      `json:"filePageLimit" yaml:"filePageLimit"`
			AllowedLang            string   `json:"allowedLang" yaml:"allowedLang"`
			MockFailSwitch         bool     `json:"mockFailSwitch" yaml:"mockFailSwitch"`
			InternalOrderChannel   []string `json:"internal-order-channel" yaml:"internal-order-channel"`
		} `json:"fullTextTranslate" yaml:"fullTextTranslate"`
	} `json:"translate" yaml:"translate"`

	OAuth2 struct {
		AppID string `json:"app_id" yaml:"app_id"`
		JWT   struct {
			Secret        string `json:"secret" yaml:"secret"`
			Issuer        string `json:"issuer" yaml:"issuer"`
			Expiry        int    `json:"expiry" yaml:"expiry"`               // 秒
			RefreshExpiry int    `json:"refreshExpiry" yaml:"refreshExpiry"` // 秒
		} `json:"jwt" yaml:"jwt"`
		TokenStorage struct {
			Type                string `json:"type" yaml:"type"`                               // 存储类型，目前支持redis
			MaxTokensPerUser    int    `json:"maxTokensPerUser" yaml:"maxTokensPerUser"`       // 每个用户最多可以有多少个有效令牌
			TokenHeaderName     string `json:"tokenHeaderName" yaml:"tokenHeaderName"`         // 令牌头名称
			CookieTokenLabel    string `json:"cookieTokenLabel" yaml:"cookieTokenLabel"`       // token cookie 标签格式
			CookieTokenUIDLabel string `json:"cookieTokenUIDLabel" yaml:"cookieTokenUIDLabel"` // token uid cookie 标签格式
			RedisKeyPrefix      struct {
				Token        string `json:"token" yaml:"token"`               // 令牌信息的键前缀
				AccessToken  string `json:"accessToken" yaml:"accessToken"`   // 访问令牌索引的键前缀
				RefreshToken string `json:"refreshToken" yaml:"refreshToken"` // 刷新令牌索引的键前缀
				UserTokens   string `json:"userTokens" yaml:"userTokens"`     // 用户令牌集合的键前缀
			} `json:"redisKeyPrefix" yaml:"redisKeyPrefix"`
			AuthCode struct {
				Lifetime int `json:"lifetime" yaml:"lifetime"` // 认证码有效期（秒）
			} `json:"authCode" yaml:"authCode"`
		} `json:"tokenStorage" yaml:"tokenStorage"`
		RSA struct {
			PublicKey  string `json:"publicKey" yaml:"publicKey"`   // RSA公钥
			PrivateKey string `json:"privateKey" yaml:"privateKey"` // RSA私钥
		} `json:"RSA" yaml:"RSA"`
		GoogleOAuth2 struct {
			ClientID        string   `json:"client_id" yaml:"client_id"`         // Google OAuth2 Client ID
			ClientSecret    string   `json:"client_secret" yaml:"client_secret"` // Google OAuth2 Client Secret
			RedirectPath    string   `json:"redirect_path" yaml:"redirect_path"` // Google OAuth2 Redirect Path (e.g., /api/oauth2/google/callback)
			Scopes          []string `json:"scopes" yaml:"scopes"`
			LoginSuccessURL string   `json:"login_success_url" yaml:"login_success_url"` // Google OAuth2 Login Success URL
		} `json:"googleOAuth2" yaml:"googleOAuth2"`
		ResourceProtection struct {
			PublicPaths []string `json:"publicPaths" yaml:"publicPaths"` // 公开资源路径前缀，不需要验证
			AdminPaths  []string `json:"adminPaths" yaml:"adminPaths"`   // 管理员资源路径前缀，需要管理员角色
			AdminRoles  []string `json:"adminRoles" yaml:"adminRoles"`   // 管理员角色
		} `json:"resourceProtection" yaml:"resourceProtection"`
		// 服务账户配置，用于服务间通信
		ServiceAccount struct {
			TokenExpiry     int      `json:"tokenExpiry" yaml:"tokenExpiry"`         // 服务令牌有效期（秒）
			TokenHeaderName string   `json:"tokenHeaderName" yaml:"tokenHeaderName"` // 服务令牌头名称
			Roles           []string `json:"roles" yaml:"roles"`                     // 服务角色
			Services        []struct {
				Name   string `json:"name" yaml:"name"`     // 服务名称
				Id     string `json:"id" yaml:"id"`         // 服务ID
				Secret string `json:"secret" yaml:"secret"` // 服务密钥
			} `json:"services" yaml:"services"`
			TokenStorage struct {
				RedisKeyPrefix string `json:"redisKeyPrefix" yaml:"redisKeyPrefix"` // 服务令牌的键前缀
			} `json:"tokenStorage" yaml:"tokenStorage"`
			ResourceProtection struct {
				PublicPaths []string `json:"publicPaths" yaml:"publicPaths"` // 公开资源路径前缀，不需要验证
			} `json:"resourceProtection" yaml:"resourceProtection"`
		} `json:"serviceAccount" yaml:"serviceAccount"`
	} `json:"oauth2" yaml:"oauth2"`

	LLM struct {
		UseChannel string `json:"useChannel" yaml:"useChannel"`
		Channel    struct {
			Deepseek struct {
				URL    string `json:"url" yaml:"url"`
				APIKey string `json:"apiKey" yaml:"apiKey"`
			} `json:"deepseek" yaml:"deepseek"`
			Gpt4oMini struct {
				URL    string `json:"url" yaml:"url"`
				APIKey string `json:"apiKey" yaml:"apiKey"`
			} `json:"gpt4o_mini" yaml:"gpt4o_mini"`
			Qwen struct {
				URL    string `json:"url" yaml:"url"`
				APIKey string `json:"apiKey" yaml:"apiKey"`
			} `json:"qwen" yaml:"qwen"`
		} `json:"channel" yaml:"channel"`
	} `json:"llm" yaml:"llm"`

	Dify struct {
		ApiBaseUrl string `json:"apiBaseUrl" yaml:"apiBaseUrl"`
		Httpclient struct {
			Timeout               int `json:"timeout" yaml:"timeout"`
			ResponseHeaderTimeout int `json:"responseHeaderTimeout" yaml:"responseHeaderTimeout"`
		} `json:"httpclient" yaml:"httpclient"`
		Datasets struct {
			Doc2DifyIntegrationDataset struct {
				Id     string `json:"id" yaml:"id"`
				Name   string `json:"name" yaml:"name"`
				ApiKey string `json:"apiKey" yaml:"apiKey"`
				Meta   []struct {
					Name string `json:"name" yaml:"name"`
					Type string `json:"type" yaml:"type"`
				} `json:"meta" yaml:"meta"`
			} `json:"doc-2-dify-integration-dataset" yaml:"doc-2-dify-integration-dataset"`
		} `json:"datasets" yaml:"datasets"`
		Chatflows struct {
			Doc2DifyIntegrationChatWorkflow struct {
				SummarySinglePaper struct {
					Id     string `json:"id" yaml:"id"`
					Name   string `json:"name" yaml:"name"`
					ApiKey string `json:"apiKey" yaml:"apiKey"`
				} `json:"summary-single-paper" yaml:"summary-single-paper"`
				SinglePaperCopilotChat struct {
					Id     string `json:"id" yaml:"id"`
					Name   string `json:"name" yaml:"name"`
					ApiKey string `json:"apiKey" yaml:"apiKey"`
				} `json:"single-paper-copilot-chat" yaml:"single-paper-copilot-chat"`
			} `json:"doc-2-dify-integration-chat-workflow" yaml:"doc-2-dify-integration-chat-workflow"`
		} `json:"chatflows" yaml:"chatflows"`
		Workflows struct {
			Doc2DifyIntegrationWorkflow struct {
				Name   string `json:"name" yaml:"name"`
				ApiKey string `json:"apiKey" yaml:"apiKey"`
			} `json:"doc-2-dify-integration-workflow" yaml:"doc-2-dify-integration-workflow"`
		} `json:"workflows" yaml:"workflows"`
	} `json:"dify" yaml:"dify"`

	Copilot struct {
		SummaryPaper struct {
			HandlerFrom                string   `json:"handler-from" yaml:"handler-from"`
			HandlerVersion             string   `json:"handler-version" yaml:"handler-version"`
			SummaryHandlerSplitStrings []string `json:"summary-handler-split-strings" yaml:"summary-handler-split-strings"`
		} `json:"summary-paper" yaml:"summary-paper"`
	} `json:"copilot" yaml:"copilot"`
	Membership MembershipConfig `json:"membership" yaml:"membership"`

	// 支付配置
	Pay struct {
		Stripe struct {
			IsEnable           bool   `json:"isEnable" yaml:"isEnable"`
			PublishableKey     string `json:"publishableKey" yaml:"publishableKey"`
			SecretKey          string `json:"secretKey" yaml:"secretKey"`
			WebhookSecret      string `json:"webhookSecret" yaml:"webhookSecret"`
			CheckoutSuccessURL string `json:"checkoutSuccessURL" yaml:"checkoutSuccessURL"`
			CheckoutCancelURL  string `json:"checkoutCancelURL" yaml:"checkoutCancelURL"`
		} `json:"stripe" yaml:"stripe"`
	} `json:"pay" yaml:"pay"`

	// 调度器配置
	Scheduler struct {
		Jobs struct {
			MembershipExpiredJob struct {
				Spec   string `json:"spec" yaml:"spec"`
				Key    string `json:"key" yaml:"key"`
				Expiry int    `json:"expiry" yaml:"expiry"`
			} `json:"membership-expired-job" yaml:"membership-expired-job"`
			CreditPayConfirmExpiredJob struct {
				Spec   string `json:"spec" yaml:"spec"`
				Key    string `json:"key" yaml:"key"`
				Expiry int    `json:"expiry" yaml:"expiry"`
			} `json:"credit-pay-confirm-expired-job" yaml:"credit-pay-confirm-expired-job"`
		} `json:"jobs" yaml:"jobs"`
	} `json:"scheduler" yaml:"scheduler"`

	// PDF相关配置
	PDF PDFConfig `json:"pdf" yaml:"pdf"`

	// 个人配置
	Personal PersonalConfig `json:"personal" yaml:"personal"`

	// 网站配置
	Nav NavConfig `json:"nav" yaml:"nav"`

	// 调试相关配置
	Debug struct {
		// 是否启用请求日志记录
		EnableRequestLogging bool `json:"enableRequestLogging" yaml:"enableRequestLogging"`
		// 是否记录请求体
		LogRequestBody bool `json:"logRequestBody" yaml:"logRequestBody"`
		// 是否记录响应体
		LogResponseBody bool `json:"logResponseBody" yaml:"logResponseBody"`
		// 最大记录的请求体大小（字节）
		MaxRequestBodySize int `json:"maxRequestBodySize" yaml:"maxRequestBodySize"`
	} `json:"debug" yaml:"debug"`

	// 文章问题配置
	PaperQuestion *PaperQuestionConfig `json:"paperQuestion" yaml:"paperQuestion"`

	Squid struct {
		ProxyUrl string `json:"proxy_url" yaml:"proxy_url"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		Timeout  int    `json:"timeout" yaml:"timeout"`
	} `json:"squid" yaml:"squid"`

	LoginPage struct {
		IsGoogleLoginEnabled           bool `json:"is_google_login_enabled" yaml:"is_google_login_enabled"`
		IsUsernamePasswordLoginEnabled bool `json:"is_username_password_login_enabled" yaml:"is_username_password_login_enabled"`
		IsRegisterEnabled              bool `json:"is_register_enabled" yaml:"is_register_enabled"`
		IsForgetPasswordEnabled        bool `json:"is_forget_password_enabled" yaml:"is_forget_password_enabled"`
	} `json:"login_page" yaml:"login_page"`

	DocMetaInfoSearch struct {
		TransformDoiDocTypeMapRel  map[string]string `json:"transformDoiDocTypeMapRel" yaml:"transformDoiDocTypeMapRel"`
		TransformDoiLanguageMapRel map[string]string `json:"transformDoiLanguageMapRel" yaml:"transformDoiLanguageMapRel"`
		DocTypeInfoListJsonDesc    string            `json:"docTypeInfoListJsonDesc" yaml:"docTypeInfoListJsonDesc"`
		Doi                        struct {
			QueryDoiInfoUrl      string `json:"query_doi_info_url" yaml:"query_doi_info_url"`
			QueryPaperDoiInfoUrl string `json:"query_paper_doi_info_url" yaml:"query_paper_doi_info_url"`
		} `json:"doi" yaml:"doi"`
	} `json:"doc_meta_info_search" yaml:"doc_meta_info_search"`
}

// RocketMQConfig RocketMQ配置
type RocketMQConfig struct {
	Enabled     bool   `mapstructure:"enabled" yaml:"enabled"`
	NameServer  string `mapstructure:"name-server" yaml:"name-server"`
	GrpcAddress string `mapstructure:"grpc-address" yaml:"grpc-address"`
	AccessKey   string `mapstructure:"access-key" yaml:"access-key"`
	SecretKey   string `mapstructure:"secret-key" yaml:"secret-key"`
	Client      struct {
		LogLevel         string `mapstructure:"log-level" yaml:"log-level"`
		RequestTimeout   int    `mapstructure:"request-timeout" yaml:"request-timeout"`
		RetryTimes       int    `mapstructure:"retry-times" yaml:"retry-times"`
		RetryInterval    int    `mapstructure:"retry-interval" yaml:"retry-interval"`
		RetryMaxInterval int    `mapstructure:"retry-max-interval" yaml:"retry-max-interval"`
		MaxMessageSize   int    `mapstructure:"max-message-size" yaml:"max-message-size"`
	} `mapstructure:"client" yaml:"client"`
	Topic struct {
		UploadCallback struct {
			Name              string `mapstructure:"name" yaml:"name"`
			Group             string `mapstructure:"group" yaml:"group"`
			Tag               string `mapstructure:"tag" yaml:"tag"`
			BatchSize         int    `mapstructure:"batch-size" yaml:"batch-size"`
			MaxReconsumeTimes int    `mapstructure:"max-reconsume-times" yaml:"max-reconsume-times"`
			ConsumeTimeout    int    `mapstructure:"consume-timeout" yaml:"consume-timeout"`
			ManualAck         bool   `mapstructure:"manual-ack" yaml:"manual-ack"`
			MaxAwaitTime      int    `mapstructure:"max-await-time" yaml:"max-await-time"`
		} `mapstructure:"upload-callback" yaml:"upload-callback"`
		ParsePdfHeader struct {
			Name              string `mapstructure:"name" yaml:"name"`
			Group             string `mapstructure:"group" yaml:"group"`
			Tag               string `mapstructure:"tag" yaml:"tag"`
			BatchSize         int    `mapstructure:"batch-size" yaml:"batch-size"`
			MaxReconsumeTimes int    `mapstructure:"max-reconsume-times" yaml:"max-reconsume-times"`
			ConsumeTimeout    int    `mapstructure:"consume-timeout" yaml:"consume-timeout"`
			ManualAck         bool   `mapstructure:"manual-ack" yaml:"manual-ack"`
			MaxAwaitTime      int    `mapstructure:"max-await-time" yaml:"max-await-time"`
		} `mapstructure:"parse-pdf-header" yaml:"parse-pdf-header"`
		ParsePdfText struct {
			Name              string `mapstructure:"name" yaml:"name"`
			Group             string `mapstructure:"group" yaml:"group"`
			Tag               string `mapstructure:"tag" yaml:"tag"`
			BatchSize         int    `mapstructure:"batch-size" yaml:"batch-size"`
			MaxReconsumeTimes int    `mapstructure:"max-reconsume-times" yaml:"max-reconsume-times"`
			ConsumeTimeout    int    `mapstructure:"consume-timeout" yaml:"consume-timeout"`
			ManualAck         bool   `mapstructure:"manual-ack" yaml:"manual-ack"`
			MaxAwaitTime      int    `mapstructure:"max-await-time" yaml:"max-await-time"`
		} `mapstructure:"parse-pdf-text" yaml:"parse-pdf-text"`
		FullTextTranslateUploadHandler struct {
			Name              string `mapstructure:"name" yaml:"name"`
			Group             string `mapstructure:"group" yaml:"group"`
			Tag               string `mapstructure:"tag" yaml:"tag"`
			BatchSize         int    `mapstructure:"batch-size" yaml:"batch-size"`
			MaxReconsumeTimes int    `mapstructure:"max-reconsume-times" yaml:"max-reconsume-times"`
			ConsumeTimeout    int    `mapstructure:"consume-timeout" yaml:"consume-timeout"`
			ManualAck         bool   `mapstructure:"manual-ack" yaml:"manual-ack"`
			MaxAwaitTime      int    `mapstructure:"max-await-time" yaml:"max-await-time"`
		} `mapstructure:"full-text-translate-upload-handler" yaml:"full-text-translate-upload-handler"`
		FullTextTranslateProgress struct {
			Name              string `mapstructure:"name" yaml:"name"`
			Group             string `mapstructure:"group" yaml:"group"`
			Tag               string `mapstructure:"tag" yaml:"tag"`
			BatchSize         int    `mapstructure:"batch-size" yaml:"batch-size"`
			MaxReconsumeTimes int    `mapstructure:"max-reconsume-times" yaml:"max-reconsume-times"`
			ConsumeTimeout    int    `mapstructure:"consume-timeout" yaml:"consume-timeout"`
			ManualAck         bool   `mapstructure:"manual-ack" yaml:"manual-ack"`
			MaxAwaitTime      int    `mapstructure:"max-await-time" yaml:"max-await-time"`
		} `mapstructure:"full-text-translate-progress" yaml:"full-text-translate-progress"`
		Event struct {
			Doc2DifyIntegrationEvent struct {
				Name              string `mapstructure:"name" yaml:"name"`
				Group             string `mapstructure:"group" yaml:"group"`
				Tag               string `mapstructure:"tag" yaml:"tag"`
				BatchSize         int    `mapstructure:"batch-size" yaml:"batch-size"`
				MaxReconsumeTimes int    `mapstructure:"max-reconsume-times" yaml:"max-reconsume-times"`
				ConsumeTimeout    int    `mapstructure:"consume-timeout" yaml:"consume-timeout"`
				ManualAck         bool   `mapstructure:"manual-ack" yaml:"manual-ack"`
				MaxAwaitTime      int    `mapstructure:"max-await-time" yaml:"max-await-time"`
			} `mapstructure:"doc-2-dify-integration-event" yaml:"doc-2-dify-integration-event"`
		} `mapstructure:"event" yaml:"event"`
	} `mapstructure:"topic" yaml:"topic"`
}

// PDFConfig PDF相关配置
type PDFConfig struct {
	RePattern string `json:"rePattern" yaml:"rePattern"` // PDF解析图表和参考文献的正则表达式
	Download  struct {
		TempDirectory         string `json:"tempDirectory" yaml:"tempDirectory"`                 // PDF下载临时文件目录
		TempDownloadDirectory string `json:"tempDownloadDirectory" yaml:"tempDownloadDirectory"` // PDF Url下载的临时文件目录
		MineruImageDirectory  string `json:"mineruImageDirectory" yaml:"mineruImageDirectory"`   // MinerU的图片目录
	} `json:"download" yaml:"download"`

	Parse struct {
		Grobid struct {
			URL              string `json:"url" yaml:"url"`                           // Grobid服务地址
			HeaderDocument   string `json:"headerDocument" yaml:"headerDocument"`     // Grobid头部文档解析地址
			FulltextDocument string `json:"fulltextDocument" yaml:"fulltextDocument"` // Grobid全文文档解析地址
			Timeout          int    `json:"timeout" yaml:"timeout"`                   // Grobid超时 单位：分钟
		} `json:"grobid" yaml:"grobid"`
		Mineru struct {
			URL         string `json:"url" yaml:"url"`                 // MinerU服务地址
			UploadURL   string `json:"uploadURL" yaml:"uploadURL"`     // MinerU上传地址
			ProgressURL string `json:"progressURL" yaml:"progressURL"` // MinerU上传进度地址
			JoinURL     string `json:"joinURL" yaml:"joinURL"`         // MinerU加入解析队列
			DataURL     string `json:"dataURL" yaml:"dataURL"`         // MinerU获取解析结果
			FileURL     string `json:"fileURL" yaml:"fileURL"`         // MinerU文件地址
			Timeout     int    `json:"timeout" yaml:"timeout"`         // MinerU超时 单位：分钟
			MaxPage     int    `json:"maxPage" yaml:"maxPage"`         // MinerU解析的最大页数
		} `json:"mineru" yaml:"mineru"`
	} `json:"parse" yaml:"parse"`
}

type PersonalConfig struct {
	LatestReadSize int `json:"latestReadSize" yaml:"latestReadSize"` // 最近阅读文献数量
}

// 定义一个清晰的结构来表示中英文问题
type Question struct {
	Id int    `json:"id" yaml:"id"`
	Zh string `json:"zh" yaml:"zh"`
	En string `json:"en" yaml:"en"`
}

// 定义列表中的每一项
type QuestionItem struct {
	// 将 YAML 中的 "question" 字段直接映射到这个 Question 结构体上
	Question Question `json:"question" yaml:"question"`
}

// 文章问题主配置结构体
type PaperQuestionConfig struct {
	DefaultQuestionCount int `json:"defaultQuestionCount" yaml:"defaultQuestionCount"`
	// 切片现在包含的是具名的、结构清晰的 QuestionItem
	DefaultQuestionList []QuestionItem `json:"defaultQuestionList" yaml:"defaultQuestionList"`
	// 问题ID到问题的映射，方便快速查询
	QuestionMap map[int]Question `json:"-" yaml:"-"`
}

// InitQuestionMap 初始化问题映射
func (c *PaperQuestionConfig) InitQuestionMap() {
	c.QuestionMap = make(map[int]Question)
	for _, item := range c.DefaultQuestionList {
		c.QuestionMap[item.Question.Id] = item.Question
	}
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	ext := strings.ToLower(filepath.Ext(path))

	if ext != ".yaml" && ext != ".yml" {
		return nil, fmt.Errorf("unsupported config file format: %s, only YAML (.yaml, .yml) is supported", ext)
	}

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode YAML config file: %w", err)
	}
	if config.PaperQuestion != nil {
		config.PaperQuestion.InitQuestionMap()
	}

	// 从环境变量读取 PDF 下载目录配置
	if envTempDir := os.Getenv("PDF_DOWNLOAD_TEMP_PDF_DIR"); envTempDir != "" {
		config.PDF.Download.TempDirectory = envTempDir
	}
	if envDownloadDir := os.Getenv("PDF_DOWNLOAD_TEMP_DOWNLOAD_DIR"); envDownloadDir != "" {
		config.PDF.Download.TempDownloadDirectory = envDownloadDir
	}
	if envMineruDir := os.Getenv("PDF_DOWNLOAD_MINERU_IMAGE_DIR"); envMineruDir != "" {
		config.PDF.Download.MineruImageDirectory = envMineruDir
	}

	return &config, nil
}

// GlobalConfig returns a default configuration
func GlobalConfig() *Config {
	config := &Config{}

	// Server defaults
	config.Server.Port = 8080
	config.Server.Host = "0.0.0.0"
	config.Server.Timeout = 30
	config.Server.GRPC.Port = 50052
	config.Server.GRPC.Host = "0.0.0.0"

	// Logging defaults
	config.Logging.Level = "info"
	config.Logging.Format = "springboot" // 默认使用 Spring Boot 风格的日志格式

	// Tracing defaults
	config.Tracing.Enabled = true
	config.Tracing.ServiceName = "go-sea-service"
	config.Tracing.JaegerURL = "http://localhost:14268/api/traces"
	config.Tracing.SampleRate = 0.1

	// Metrics defaults
	config.Metrics.Enabled = true
	config.Metrics.Path = "/metrics"

	// Error notification defaults
	config.ErrorNotification.Enabled = false
	config.ErrorNotification.Type = "slack"

	// Database defaults
	config.Database.Enabled = true
	config.Database.Type = "postgres"
	// PostgreSQL defaults
	config.Database.Postgres.Host = "localhost"
	config.Database.Postgres.Port = 5432
	config.Database.Postgres.User = "postgres"
	config.Database.Postgres.Password = "postgres"
	config.Database.Postgres.DBName = "go_sea"
	config.Database.Postgres.SSLMode = "disable"
	config.Database.Postgres.TimeZone = "Asia/Shanghai"
	config.Database.Postgres.MaxIdleConns = 10
	config.Database.Postgres.MaxOpenConns = 100
	config.Database.Postgres.ConnMaxLifetime = 3600
	config.Database.Postgres.LogLevel = "info"
	// SQLite defaults
	config.Database.SQLite.DBPath = "./data/app.db"
	config.Database.SQLite.MaxIdleConns = 1
	config.Database.SQLite.MaxOpenConns = 1
	config.Database.SQLite.ConnMaxLifetime = 3600
	config.Database.SQLite.LogLevel = "info"

	// MinIO 配置必须在配置文件中完整定义，不设置任何默认值
	// 如果配置缺失，系统将无法启动

	// Redis defaults
	config.Redis.Enabled = false
	config.Redis.Host = "localhost"
	config.Redis.Port = 6379
	config.Redis.Password = ""
	config.Redis.DB = 0
	config.Redis.PoolSize = 10
	config.Redis.MinIdleConns = 5
	config.Redis.DialTimeout = 5
	config.Redis.ReadTimeout = 3
	config.Redis.WriteTimeout = 3
	config.Redis.MaxConnAge = 3600
	config.Redis.MaxRetries = 3
	config.Redis.MinRetryBackoff = 8
	config.Redis.MaxRetryBackoff = 512

	// Cache defaults
	config.Cache.Type = "redis"    // 默认使用 redis，本地版可改为 memory
	config.Cache.Expiration = 1800 // 默认 30 分钟

	// OAuth2 defaults
	config.OAuth2.AppID = "go-sea"
	config.OAuth2.JWT.Secret = "go-sea-sjdpoi2(&(*)&%290371jd-oauth2-secret"
	config.OAuth2.JWT.Issuer = "go-sea"
	config.OAuth2.JWT.Expiry = 86400
	config.OAuth2.JWT.RefreshExpiry = 86400
	config.OAuth2.TokenStorage.Type = "redis"
	config.OAuth2.TokenStorage.MaxTokensPerUser = 5
	config.OAuth2.TokenStorage.TokenHeaderName = "Access-Control-Request-Token"
	config.OAuth2.TokenStorage.CookieTokenLabel = "%s-token"
	config.OAuth2.TokenStorage.CookieTokenUIDLabel = "%s-UID"
	config.OAuth2.TokenStorage.RedisKeyPrefix.Token = "token:%s"
	config.OAuth2.TokenStorage.RedisKeyPrefix.AccessToken = "access_token:%s"
	config.OAuth2.TokenStorage.RedisKeyPrefix.RefreshToken = "refresh_token:%s"
	config.OAuth2.TokenStorage.RedisKeyPrefix.UserTokens = "user_tokens:%d"
	config.OAuth2.ResourceProtection.PublicPaths = []string{"/api/public/", "/api/oauth2/token", "/api/oauth2/refresh", "/api/oauth2/auth_code", "/api/oauth2/sign_in", "/api/oauth2/validate"}
	config.OAuth2.ResourceProtection.AdminPaths = []string{"/admin"}
	config.OAuth2.ResourceProtection.AdminRoles = []string{"admin"}

	// 服务账户默认配置
	config.OAuth2.ServiceAccount.TokenExpiry = 2592000 // 30天
	config.OAuth2.ServiceAccount.TokenHeaderName = "Service-Authorization"
	config.OAuth2.ServiceAccount.Roles = []string{"ROLE_SERVICE"}
	config.OAuth2.ServiceAccount.TokenStorage.RedisKeyPrefix = "service_token:"

	// 翻译配置默认值
	config.Translate.Text.Special.Words.NeedReplace.List = []string{
		"\\p{Cntrl}",
		"\\u0000",
	}
	config.Translate.Text.MaxLength = 2000
	config.Translate.Text.Channel.Youdao.AppId = "3b32b4416f33f420"
	config.Translate.Text.Channel.Youdao.AppSecretKey = "wns0F0fn2jkbFt7PNbOj3D4SIEyB7tfY"
	config.Translate.Text.Channel.Youdao.ApiUrl = "https://openapi.youdao.com/api"
	config.Translate.Text.Channel.Google.BaseURL = "http://47.236.50.140:8090"
	config.Translate.Text.Channel.Google.FreeUri = "/gg/trn"
	config.Translate.Text.Channel.Google.SimpleFreeUri = "/gg/ss"

	config.Translate.Text.DefaultSourceLanguage = pkgi18n.LanguageEnUS
	config.Translate.Text.DefaultTargetLanguage = pkgi18n.LanguageZhCN

	// 术语库配置默认值
	config.Translate.Glossary.MaxSize = 1000
	config.Translate.Glossary.MaxLength = 100

	// 全文翻译配置默认值
	config.Translate.FullTextTranslate.TranslateBaseURL = "http://localhost:8891"
	config.Translate.FullTextTranslate.TranslateServiceURI = "/translate"
	config.Translate.FullTextTranslate.TranslateServiceFixURI = "/retrans"
	config.Translate.FullTextTranslate.HistoryVisibleDays = 30
	config.Translate.FullTextTranslate.TranslateTimeOut = 60
	config.Translate.FullTextTranslate.DuplicateDelaySeconds = 300
	config.Translate.FullTextTranslate.FileSizeLimit = 104857600
	config.Translate.FullTextTranslate.FilePageLimit = 1000
	config.Translate.FullTextTranslate.AllowedLang = strings.Join(pkgi18n.GetSupportedLanguages(), ",")
	config.Translate.FullTextTranslate.MockFailSwitch = false
	config.Translate.FullTextTranslate.InternalOrderChannel = []string{"googleFree", "googleSimpleFree", "youdao"}

	// 调试相关配置默认值
	config.Debug.EnableRequestLogging = false
	config.Debug.LogRequestBody = false
	config.Debug.LogResponseBody = false
	config.Debug.MaxRequestBodySize = 1024

	//设置文件目录默认值
	config.PDF.Download.TempDirectory = "temp/pdf"
	config.PDF.Download.TempDownloadDirectory = "temp/download"
	config.PDF.Download.MineruImageDirectory = "images"

	//设置RocketMQ配置默认值
	config.RocketMQ.Client.LogLevel = "ERROR"
	config.RocketMQ.Client.RequestTimeout = 30000
	config.RocketMQ.Client.RetryTimes = 3
	config.RocketMQ.Client.RetryInterval = 3000
	config.RocketMQ.Client.RetryMaxInterval = 15000
	config.RocketMQ.Client.MaxMessageSize = 4194304

	return config
}
