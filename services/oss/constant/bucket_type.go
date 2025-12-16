package constant

import (
	"strings"

	"github.com/yb2020/odoc/config"
	pb "github.com/yb2020/odoc/proto/gen/go/oss"
)

// BucketTypeToEnum 根据桶名称获取对应的枚举值
func BucketTypeToEnum(cfg *config.Config, bucketName string) pb.OSSBucketEnum {
	// 首先尝试从配置中找到对应的逻辑名称
	for logicalName, bucketCfg := range cfg.OSS.S3.Buckets {
		if bucketCfg.Name == bucketName {
			// 尝试将逻辑名称转换为枚举
			// 假设逻辑名称与枚举名称的关系是：逻辑名称是小写，枚举名称是大写
			enumName := strings.ToUpper(logicalName)

			// 遍历所有枚举值，查找匹配项
			for _, v := range pb.OSSBucketEnum_value {
				enum := pb.OSSBucketEnum(v)
				if enum.String() == enumName {
					return enum
				}
			}
		}
	}
	// 如果没有找到匹配的桶，返回默认值
	return pb.OSSBucketEnum_PUBLIC
}

// EnumToBucketType 将枚举转换为逻辑桶名称
func EnumToBucketType(cfg *config.Config, enum pb.OSSBucketEnum) string {
	// 将枚举名称转换为小写作为逻辑名称
	logicalName := strings.ToLower(enum.String())

	// 验证逻辑名称是否存在于配置中
	if _, ok := cfg.OSS.S3.Buckets[logicalName]; ok {
		return logicalName // 返回逻辑名称，不是实际桶名称
	}
	// 如果找不到，返回默认的逻辑名称
	if _, ok := cfg.OSS.S3.Buckets["public"]; ok {
		return "public"
	}
	return ""
}
