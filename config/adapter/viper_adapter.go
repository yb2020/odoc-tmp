package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/yb2020/odoc/config"
)

// ViperConfig 使用Viper的配置管理器
type ViperConfig struct {
	viper *viper.Viper
}

// NewViperConfig 创建一个新的Viper配置管理器
func NewViperConfig() *ViperConfig {
	return &ViperConfig{
		viper: viper.New(),
	}
}

// LoadConfigWithViper 使用Viper加载配置并支持环境变量覆盖
func LoadConfigWithViper(path string) (*config.Config, error) {
	// 1. 先使用 config.LoadConfig 读取 YAML 配置
	// 这样可以确保带下划线的键被正确处理
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// 2. 创建一个新的 Viper 实例用于环境变量处理
	v := viper.New()

	// 3. 启用环境变量支持
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 4. 处理各类环境变量并应用到配置对象
	processServerEnv(v, cfg)
	processDatabaseEnv(v, cfg)
	processRedisEnv(v, cfg)
	processOSSEnv(v, cfg)
	processRocketMQEnv(v, cfg)
	processPDFEnv(v, cfg)
	processDifyEnv(v, cfg)
	processFullTextTranslateEnv(v, cfg)
	processOcrEnv(v, cfg)
	processSquidEnv(v, cfg)

	// 5. 更新全局配置
	config.SetConfig(cfg)

	return cfg, nil
}

// processServerEnv 处理服务器环境变量
func processServerEnv(v *viper.Viper, cfg *config.Config) {
	if port := v.GetInt("SERVER_PORT"); port != 0 {
		cfg.Server.Port = port
	}
	if host := v.GetString("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if timeout := v.GetInt("SERVER_TIMEOUT"); timeout != 0 {
		cfg.Server.Timeout = timeout
	}
}

// processDatabaseEnv 处理数据库环境变量
func processDatabaseEnv(v *viper.Viper, cfg *config.Config) {
	// PostgreSQL 环境变量
	if host := v.GetString("DATABASE_HOST"); host != "" {
		cfg.Database.Postgres.Host = host
	}
	if port := v.GetInt("DATABASE_PORT"); port != 0 {
		cfg.Database.Postgres.Port = port
	}
	if user := v.GetString("DATABASE_USER"); user != "" {
		cfg.Database.Postgres.User = user
	}
	if password := v.GetString("DATABASE_PASSWORD"); password != "" {
		cfg.Database.Postgres.Password = password
	}
	if dbName := v.GetString("DATABASE_DBNAME"); dbName != "" {
		cfg.Database.Postgres.DBName = dbName
	}
	if maxIdleConns := v.GetInt("DATABASE_MAX_IDLE_CONNS"); maxIdleConns != 0 {
		cfg.Database.Postgres.MaxIdleConns = maxIdleConns
	}
	if maxOpenConns := v.GetInt("DATABASE_MAX_OPEN_CONNS"); maxOpenConns != 0 {
		cfg.Database.Postgres.MaxOpenConns = maxOpenConns
	}
	// SQLite 环境变量
	if dbPath := v.GetString("SQLITE_DB_PATH"); dbPath != "" {
		cfg.Database.SQLite.DBPath = dbPath
	}
}

// processRedisEnv 处理Redis环境变量
func processRedisEnv(v *viper.Viper, cfg *config.Config) {
	if host := v.GetString("REDIS_HOST"); host != "" {
		cfg.Redis.Host = host
	}
	if port := v.GetInt("REDIS_PORT"); port != 0 {
		cfg.Redis.Port = port
	}
	if password := v.GetString("REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}
	if db := v.GetInt("REDIS_DB"); db != 0 {
		cfg.Redis.DB = db
	}
}

// processOSSEnv 处理OSS环境变量
func processOSSEnv(v *viper.Viper, cfg *config.Config) {
	// 处理 S3 配置
	if endpoint := v.GetString("OSS_S3_ENDPOINT"); endpoint != "" {
		cfg.OSS.S3.Endpoint = endpoint
	}
	if useSSL := v.GetBool("OSS_S3_USE_SSL"); useSSL {
		cfg.OSS.S3.UseSSL = useSSL
	}
	if location := v.GetString("OSS_S3_LOCATION"); location != "" {
		cfg.OSS.S3.Location = location
	}
	if accessKey := v.GetString("OSS_S3_ACCESS_KEY_ID"); accessKey != "" {
		cfg.OSS.S3.AccessKeyID = accessKey
	}
	if secretKey := v.GetString("OSS_S3_SECRET_ACCESS_KEY"); secretKey != "" {
		cfg.OSS.S3.SecretAccessKey = secretKey
	}
	if region := v.GetString("OSS_S3_REGION"); region != "" {
		cfg.OSS.S3.Region = region
	}
	if forcePathStyle := v.GetBool("OSS_S3_FORCE_PATH_STYLE"); forcePathStyle {
		cfg.OSS.S3.ForcePathStyle = forcePathStyle
	}
}

// processRocketMQEnv 处理RocketMQ环境变量
func processRocketMQEnv(v *viper.Viper, cfg *config.Config) {
	if enabled := v.GetBool("ROCKETMQ_ENABLED"); enabled {
		cfg.RocketMQ.Enabled = enabled
	}
	if nameServer := v.GetString("ROCKETMQ_NAME_SERVER"); nameServer != "" {
		cfg.RocketMQ.NameServer = nameServer
	}
	if grpcAddress := v.GetString("ROCKETMQ_GRPC_ADDRESS"); grpcAddress != "" {
		cfg.RocketMQ.GrpcAddress = grpcAddress
	}
}

// processPDFEnv 处理PDF相关环境变量
func processPDFEnv(v *viper.Viper, cfg *config.Config) {
	if grobidURL := v.GetString("PDF_PARSE_GROBID_URL"); grobidURL != "" {
		cfg.PDF.Parse.Grobid.URL = grobidURL
	}
	if mineruURL := v.GetString("PDF_PARSE_MINERU_URL"); mineruURL != "" {
		cfg.PDF.Parse.Mineru.URL = mineruURL
	}
}

// processDifyEnv 处理Dify相关环境变量
func processDifyEnv(v *viper.Viper, cfg *config.Config) {
	if difyApiBaseUrl := v.GetString("DIFY_API_BASE_URL"); difyApiBaseUrl != "" {
		cfg.Dify.ApiBaseUrl = difyApiBaseUrl
	}
	if difyDatasetId := v.GetString("DIFY_DATASET_ID"); difyDatasetId != "" {
		cfg.Dify.Datasets.Doc2DifyIntegrationDataset.Id = difyDatasetId
	}
	if difyDatasetName := v.GetString("DIFY_DATASET_NAME"); difyDatasetName != "" {
		cfg.Dify.Datasets.Doc2DifyIntegrationDataset.Name = difyDatasetName
	}
	if difyDatasetApiKey := v.GetString("DIFY_DATASET_API_KEY"); difyDatasetApiKey != "" {
		cfg.Dify.Datasets.Doc2DifyIntegrationDataset.ApiKey = difyDatasetApiKey
	}
	if difyChatflowSummaryId := v.GetString("DIFY_CHATFLOW_SUMMARY_ID"); difyChatflowSummaryId != "" {
		cfg.Dify.Chatflows.Doc2DifyIntegrationChatWorkflow.SummarySinglePaper.Id = difyChatflowSummaryId
	}
	if difyChatflowSummaryName := v.GetString("DIFY_CHATFLOW_SUMMARY_NAME"); difyChatflowSummaryName != "" {
		cfg.Dify.Chatflows.Doc2DifyIntegrationChatWorkflow.SummarySinglePaper.Name = difyChatflowSummaryName
	}
	if difyChatflowSummaryApiKey := v.GetString("DIFY_CHATFLOW_SUMMARY_API_KEY"); difyChatflowSummaryApiKey != "" {
		cfg.Dify.Chatflows.Doc2DifyIntegrationChatWorkflow.SummarySinglePaper.ApiKey = difyChatflowSummaryApiKey
	}
	if difyChatflowCopilotId := v.GetString("DIFY_CHATFLOW_COPILOT_ID"); difyChatflowCopilotId != "" {
		cfg.Dify.Chatflows.Doc2DifyIntegrationChatWorkflow.SinglePaperCopilotChat.Id = difyChatflowCopilotId
	}
	if difyChatflowCopilotName := v.GetString("DIFY_CHATFLOW_COPILOT_NAME"); difyChatflowCopilotName != "" {
		cfg.Dify.Chatflows.Doc2DifyIntegrationChatWorkflow.SinglePaperCopilotChat.Name = difyChatflowCopilotName
	}
	if difyChatflowCopilotApiKey := v.GetString("DIFY_CHATFLOW_COPILOT_API_KEY"); difyChatflowCopilotApiKey != "" {
		cfg.Dify.Chatflows.Doc2DifyIntegrationChatWorkflow.SinglePaperCopilotChat.ApiKey = difyChatflowCopilotApiKey
	}
	if difyWorkflowName := v.GetString("DIFY_WORKFLOW_NAME"); difyWorkflowName != "" {
		cfg.Dify.Workflows.Doc2DifyIntegrationWorkflow.Name = difyWorkflowName
	}
	if difyWorkflowApiKey := v.GetString("DIFY_WORKFLOW_API_KEY"); difyWorkflowApiKey != "" {
		cfg.Dify.Workflows.Doc2DifyIntegrationWorkflow.ApiKey = difyWorkflowApiKey
	}
}

// processFullTextTranslateEnv 处理全文翻译环境变量
func processFullTextTranslateEnv(v *viper.Viper, cfg *config.Config) {
	if baseURL := v.GetString("FULLTEXT_TRANSLATE_BASE_URL"); baseURL != "" {
		cfg.Translate.FullTextTranslate.TranslateBaseURL = baseURL
	}
}

// processOcrEnv 处理 OCR 环境变量
func processOcrEnv(v *viper.Viper, cfg *config.Config) {
	if extractTextURL := v.GetString("TRANSLATE_OCR_EXTRACT_TEXT_URL"); extractTextURL != "" {
		cfg.Translate.OCR.ExtractTextURL = extractTextURL
	}
}

// GlobalViperConfig 全局Viper配置实例
var GlobalViperConfig = NewViperConfig()

// LoadConfigWithGlobalViper 使用全局Viper实例加载配置
func LoadConfigWithGlobalViper(path string) (*config.Config, error) {
	// 设置配置文件路径
	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)
	name := strings.TrimSuffix(file, ext)

	GlobalViperConfig.viper.SetConfigName(name)
	GlobalViperConfig.viper.SetConfigType(strings.TrimPrefix(ext, "."))
	GlobalViperConfig.viper.AddConfigPath(dir)

	// 读取配置文件
	if err := GlobalViperConfig.viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 启用环境变量支持
	GlobalViperConfig.viper.AutomaticEnv()
	GlobalViperConfig.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 将Viper配置转换为Config结构体
	var config config.Config
	if err := GlobalViperConfig.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// processSquidEnv 处理 Squid 环境变量
func processSquidEnv(v *viper.Viper, cfg *config.Config) {
	if proxyUrl := v.GetString("SQUID_PROXY_URL"); proxyUrl != "" {
		cfg.Squid.ProxyUrl = proxyUrl
	}
	if username := v.GetString("SQUID_USERNAME"); username != "" {
		cfg.Squid.Username = username
	}
	if password := v.GetString("SQUID_PASSWORD"); password != "" {
		cfg.Squid.Password = password
	}
	if timeout := v.GetInt("SQUID_TIMEOUT"); timeout != 0 {
		cfg.Squid.Timeout = timeout
	}
}
