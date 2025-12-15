package oss

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	configPkg "github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
)

// S3Storage AWS S3存储实现
type S3Storage struct {
	client                *s3.Client
	presignClient         *s3.PresignClient
	presignDownloadClient *s3.PresignClient
	buckets               map[string]BucketConfig
	defaultBucket         string
	config                *configPkg.Config
}

// NewS3Storage
func NewS3Storage(cfg *configPkg.Config) (*S3Storage, error) {
	// 为了代码清晰，创建一个配置的快捷方式
	s3Cfg := cfg.OSS.S3

	// 创建一个AWS配置加载器选项的切片，用于动态构建配置
	var awsLoadOptions []func(*config.LoadOptions) error

	// 1. 设置区域 (Region)
	awsLoadOptions = append(awsLoadOptions, config.WithRegion(s3Cfg.Region))

	// 2. 动态设置凭证提供者 (Credentials)
	if s3Cfg.AccessKeyID != "" && s3Cfg.SecretAccessKey != "" {
		staticCreds := credentials.NewStaticCredentialsProvider(s3Cfg.AccessKeyID, s3Cfg.SecretAccessKey, "")
		awsLoadOptions = append(awsLoadOptions, config.WithCredentialsProvider(staticCreds))
	}

	// 3. 动态设置端点解析器 (Endpoint) - 用于兼容MinIO
	if s3Cfg.Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               s3Cfg.Endpoint,
					SigningRegion:     s3Cfg.Region,
					HostnameImmutable: true,
				}, nil
			},
		)
		awsLoadOptions = append(awsLoadOptions, config.WithEndpointResolverWithOptions(customResolver))
	}

	// 加载AWS配置
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), awsLoadOptions...)
	if err != nil {
		return nil, errors.BizWrap("failed to load AWS config!", err)
	}

	// 创建S3客户端，并根据配置动态决定是否强制使用路径样式
	s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		if s3Cfg.ForcePathStyle {
			o.UsePathStyle = true
		}
	})

	// 创建预签名客户端
	presignClient := s3.NewPresignClient(s3Client)

	// 4. 创建专门用于下载的、绝对“干净”的客户端
	s3ClientForDownload := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		// 【最终解决方案 - Part 1】禁用高优先级的 RetryMode
		o.RetryMode = ""
		// 【最终解决方案 - Part 2】同时禁用底层的 Retryer
		o.Retryer = aws.NopRetryer{}

		if s3Cfg.ForcePathStyle {
			o.UsePathStyle = true
		}
	})
	presignDownloadClient := s3.NewPresignClient(s3ClientForDownload)

	// 【核心修改】动态转换桶配置，不再硬编码
	buckets := make(map[string]BucketConfig)
	for logicalName, bucketCfg := range s3Cfg.Buckets {
		// logicalName 将是 "public", "pdf", "temp" 等
		// bucketCfg 是从 YAML 加载的 config.BucketConfig
		buckets[logicalName] = BucketConfig{
			Name:       bucketCfg.Name,
			Public:     bucketCfg.Public,
			Versioning: bucketCfg.Versioning,
			// 根据生命周期配置动态判断是否为临时桶
			IsTemp: bucketCfg.LifecycleDays > 0,
			// 将天数转换为秒
			Expiration: int64(bucketCfg.LifecycleDays) * 24 * 3600,
		}
	}

	// 初始化存储实例
	storage := &S3Storage{
		client:                s3Client,
		presignClient:         presignClient,
		presignDownloadClient: presignDownloadClient,
		buckets:               buckets,
		defaultBucket:         "public", // 依然可以保留一个默认桶的逻辑名称
		config:                cfg,
	}

	// 确保所有桶存在并配置
	// 警告：在生产环境中，建议禁用此功能或通过开关控制
	if err := storage.initializeBuckets(context.Background()); err != nil {
		return nil, errors.BizWrap("failed to initialize buckets!", err)
	}
	return storage, nil
}

// initializeBuckets 初始化存储桶
func (s *S3Storage) initializeBuckets(ctx context.Context) error {
	for _, bucketConfig := range s.buckets {
		// 检查桶是否存在
		_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketConfig.Name),
		})

		if err != nil {
			// 检查是否是 404 错误（bucket 不存在）
			var notFound *types.NotFound
			if errors.As(err, &notFound) {
				// 桶确实不存在，创建它
				createBucketInput := &s3.CreateBucketInput{
					Bucket: aws.String(bucketConfig.Name),
				}

				// SDK 会自动根据 region 配置处理 LocationConstraint
				_, err = s.client.CreateBucket(ctx, createBucketInput)
				if err != nil {
					// 检查是否是 BucketAlreadyOwnedByYou 错误（并发创建）
					var alreadyExists *types.BucketAlreadyOwnedByYou
					if errors.As(err, &alreadyExists) {
						// bucket 已存在，忽略错误
						return nil
					}
					return errors.BizWrap("failed to create bucket ", err)
				}
			} else {
				// 其他错误（权限、网络等），记录日志但继续
				// 假设 bucket 可能已存在，跳过创建
				fmt.Printf("Warning: HeadBucket failed for %s: %v, assuming bucket exists\n",
					bucketConfig.Name, err)
			}
		}

		// 设置桶策略（如果是公开桶）
		if bucketConfig.Public {
			policy := fmt.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {"AWS": ["*"]},
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::%s/*"]
					}
				]
			}`, bucketConfig.Name)

			_, err = s.client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
				Bucket: aws.String(bucketConfig.Name),
				Policy: aws.String(policy),
			})
			if err != nil {
				return errors.BizWrap("failed to set bucket policy for ", err)
			}
		}

		// 设置版本控制
		if bucketConfig.Versioning {
			_, err = s.client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
				Bucket: aws.String(bucketConfig.Name),
				VersioningConfiguration: &types.VersioningConfiguration{
					Status: types.BucketVersioningStatusEnabled,
				},
			})
			if err != nil {
				return errors.BizWrap("failed to enable versioning for ", err)
			}
		}

		// 设置生命周期规则（如果是临时桶且配置了过期时间）
		if bucketConfig.IsTemp && bucketConfig.Expiration > 0 {
			days := int32(bucketConfig.Expiration / (24 * 3600)) // 转换为天数
			if days < 1 {
				days = 1
			}

			_, err = s.client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
				Bucket: aws.String(bucketConfig.Name),
				LifecycleConfiguration: &types.BucketLifecycleConfiguration{
					Rules: []types.LifecycleRule{
						{
							ID:     aws.String("expire-objects"),
							Status: types.ExpirationStatusEnabled,
							Filter: &types.LifecycleRuleFilter{
								// 将之前作为类型的 Prefix，现在设置为字段
								Prefix: aws.String(""), // 指定一个空前缀，表示应用于桶内的所有对象
							},
							Expiration: &types.LifecycleExpiration{
								Days: aws.Int32(days),
							},
						},
					},
				},
			})
			if err != nil {
				return errors.BizWrap("failed to set lifecycle for ", err)
			}
		}
	}

	return nil
}

// Upload 上传对象
func (s *S3Storage) Upload(ctx context.Context, bucket, objectName string, reader io.Reader, objectSize int64, opts map[string]string, userMetadata map[string]string) error {
	bucketName, err := s.getBucketName(bucket)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:   aws.String(bucketName),
		Key:      aws.String(objectName),
		Body:     reader,
		Metadata: userMetadata,
	}

	// 设置内容类型
	if contentType, ok := opts["content-type"]; ok {
		input.ContentType = aws.String(contentType)
	}

	_, err = s.client.PutObject(ctx, input)
	if err != nil {
		return errors.BizWrap("put object error: ", err)
	}

	return nil
}

// Download 下载对象
func (s *S3Storage) Download(ctx context.Context, bucket, objectName string) (io.Reader, error) {
	bucketName, err := s.getBucketName(bucket)
	if err != nil {
		return nil, err
	}

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, fmt.Errorf("get object error: %w", err)
	}

	return result.Body, nil
}

// Delete 删除对象
func (s *S3Storage) Delete(ctx context.Context, bucket, objectName string) error {
	bucketName, err := s.getBucketName(bucket)
	if err != nil {
		return err
	}

	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return errors.BizWrap("delete object error: ", err)
	}

	return nil
}

// GetObjectURL 获取对象的临时URL
func (s *S3Storage) GetObjectURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error) {
	bucketName, err := s.getBucketName(bucket)
	if err != nil {
		return "", err
	}
	bucketConfig, ok := s.buckets[bucket]
	if !ok {
		return "", errors.Biz("bucket not found")
	}

	if bucketConfig.Public {
		return "", errors.Biz("bucket is not private, cannot generate presigned URL")
	}

	request, err := s.presignDownloadClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", errors.BizWrap("get presigned url error！", err)
	}
	// 替换URL中的endpoint
	urlStr := request.URL
	if s.config.OSS.S3.Upload.PublicDownloadEndpoint != "" {
		urlStr = strings.Replace(urlStr, s.config.OSS.S3.Endpoint, s.config.OSS.S3.Upload.PublicDownloadEndpoint, 1)
	}

	return urlStr, nil
}

// GetPermanentURL 获取对象的永久URL（仅适用于公开桶）
func (s *S3Storage) GetPermanentURL(ctx context.Context, bucket, objectName string) (string, error) {
	bucketName, err := s.getBucketName(bucket)
	if err != nil {
		return "", err
	}
	bucketConfig, ok := s.buckets[bucket]
	if !ok {
		return "", errors.Biz("bucket not found")
	}

	if !bucketConfig.Public {
		return "", errors.Biz("bucket is not public, cannot generate permanent URL")
	}

	// 构建永久URL
	endpoint := s.config.OSS.S3.Upload.PublicDownloadEndpoint
	if endpoint == "" {
		endpoint = s.config.OSS.S3.Endpoint
	}
	//如果endpoint为空，那就是AWS的配置
	if endpoint == "" {
		//使用Region来构建URL
		region := s.config.OSS.S3.Region
		if region == "" {
			// 如果连 region 都没有，就无法构建 URL
			return "", errors.Biz("cannot generate permanent URL without endpoint or region configuration")
		}
		// 按照 AWS 的 Virtual-Hosted–Style URL 格式构建 (更推荐)   这里如果带bucketName会导致后端无法下载
		// return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectName), nil
		return fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", region, bucketName, objectName), nil
	}

	// protocol := "http"
	// if s.config.OSS.S3.UseSSL {
	// 	protocol = "https"
	// }

	permanentURL := fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectName)
	return permanentURL, nil
}

// GeneratePreSignedUpload 生成预签名上传URL（使用 PUT 方法）
// 这是 AWS SDK v2 的标准和推荐做法。
func (s *S3Storage) GeneratePreSignedUpload(ctx context.Context, bucketType, objectKey, contentType string, fileSize int64, metadata map[string]string) (*PreSignedUploadResponse, error) {
	bucketName, err := s.getBucketName(bucketType)
	if err != nil {
		return nil, err
	}

	// 在这里，您可以先对客户端传来的 fileSize 进行校验
	if fileSize > s.config.OSS.S3.Upload.MaxFileSize {
		return nil, errors.Biz("file size exceeds the limit")
	}
	if fileSize <= 0 {
		return nil, errors.Biz("invalid file size")
	}

	// 1. 准备 PutObject 的输入参数
	// 我们告诉SDK，将要发生的这个PUT操作的参数是什么
	putObjectInput := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(contentType),
		Metadata:      metadata,
		ContentLength: aws.Int64(fileSize),
	}

	// 2. 使用 PresignClient.PresignPutObject 生成预签名URL
	// 第二个参数 func(opts *s3.PresignOptions) 用于设置过期时间
	request, err := s.presignClient.PresignPutObject(ctx, putObjectInput, func(opts *s3.PresignOptions) {
		// 设置过期时间，如果配置不存在则使用默认值
		if s.config.OSS.S3.Upload.PresignedURLExpires > 0 {
			opts.Expires = time.Duration(s.config.OSS.S3.Upload.PresignedURLExpires) * time.Second
		} else {
			opts.Expires = 30 * time.Minute // 默认30分钟
		}
	})
	if err != nil {
		return nil, errors.BizWrap("failed to generate presigned PUT URL: ", err)
	}
	finalURL := request.URL
	publicEndpoint := s.config.OSS.S3.Upload.PublicUploadEndpoint
	if publicEndpoint != "" {
		finalURL = strings.Replace(finalURL, s.config.OSS.S3.Endpoint, publicEndpoint, 1)
	}
	// useSSL := s.config.OSS.S3.UseSSL
	// if useSSL {
	// 	parsedSDKURL.Scheme = "https"
	// } else {
	// 	parsedSDKURL.Scheme = "http"
	// }

	headers := make(map[string]string)
	for key, values := range request.SignedHeader {
		// http.Header 的值是 []string，但对于上传请求，通常每个 Header 只有一个值。
		if len(values) > 0 {
			// 将 Header 的 key 转换为标准的 HTTP Header 格式（例如 "Content-Type"）
			// 并赋值给 map。
			headers[key] = values[0]
		}
	}
	return &PreSignedUploadResponse{
		NeedUpload: true,
		UploadInfo: &PreSignedUploadInfo{
			URL:     finalURL, // 这是客户端需要上传到的URL
			Method:  "PUT",    // 方法是PUT
			Headers: headers,  // 客户端上传时需要包含的HTTP头
		},
	}, nil
}

func (s *S3Storage) GetBucketConfig(bucket string) (BucketConfig, error) {
	if bucketConfig, ok := s.buckets[bucket]; ok {
		return bucketConfig, nil // 直接返回 bucketConfig 值
	}
	// 如果找不到，返回一个空的 BucketConfig 和一个错误
	return BucketConfig{}, errors.Biz("bucket not found!")
}

// getBucketName 获取实际的桶名称
func (s *S3Storage) getBucketName(bucket string) (string, error) {
	if bucketConfig, ok := s.buckets[bucket]; ok {
		return bucketConfig.Name, nil
	}
	return "", errors.Biz("bucket not found!")
}

// Close 关闭存储客户端
func (s *S3Storage) Close() error {
	return nil
}
