
package utils

import "strings"

var (
	specialTLDs = []string{"gov.cn", "edu.cn", "com.cn", "org.cn", "net.cn"}
	validTlds   = []string{
		"com", "net", "org", "gov", "cc", "biz", "info", "cn", "co",
		"ai", "me", "top", "xyz", "vip", "tv", "wang", "shop", "red",
		"ltd", "ink", "pro", "kim", "club", "online", "tech", "store",
		"site", "live", "fun", "design", "work", "pub", "group", "link",
		"news", "name", "wiki", "ren", "space", "mobi", "asia",
	}
	businessValidTlds = []string{
		"paper.uat.idea.edu.cn",
		"paper.dev.idea.edu.cn",
		"paper.dev.aiteam.cc",
		"paper.uat.aiteam.cc",
	}
)

// GetMainDomain 从 host 中提取主域名，逻辑从 Java 的 DomainUtil 迁移而来
func GetMainDomain(host string) string {
	// 检查是否是特殊业务域名
	for _, businessTld := range businessValidTlds {
		if strings.HasSuffix(host, businessTld) {
			parts := strings.Split(businessTld, ".")
			if len(parts) > 2 {
				// 例如 a.b.c.d -> c.d
				return strings.Join(parts[len(parts)-2:], ".")
			}
			return businessTld
		}
	}

	// 检查是否是特殊 TLDs，例如 .com.cn, .edu.cn
	for _, specialTld := range specialTLDs {
		if strings.HasSuffix(host, "."+specialTld) {
			// 去掉 .com.cn 后剩下的部分
			remaining := strings.TrimSuffix(host, "."+specialTld)
			parts := strings.Split(remaining, ".")
			if len(parts) > 0 {
				// 取最后一部分 + specialTld
				return parts[len(parts)-1] + "." + specialTld
			}
		}
	}

	// 通用 TLD 规则
	parts := strings.Split(host, ".")
	for _, tld := range validTlds {
		if len(parts) > 1 && parts[len(parts)-1] == tld {
			return parts[len(parts)-2] + "." + parts[len(parts)-1]
		}
	}

	// 如果没有匹配到任何规则，返回原始 host
	return host
}