package utils

import (
	"strings"
)

// DesensitizationUtil provides methods for data masking/desensitization
// Converted from Java to Go

// Left only shows the first n characters, others are masked with asterisks
// Example: Left("李明", 1) returns "李*"
func Left(fullName string, index int) string {
	if len(fullName) == 0 {
		return ""
	}

	if index > len(fullName) {
		index = len(fullName)
	}

	name := fullName[:index]
	return name + strings.Repeat("*", len(fullName)-index)
}

// Around masks the middle part of a string, keeping the first 'index' characters
// and the last 'end' characters visible
// Example: Around("11012345678", 3, 2) returns "110*****78"
func Around(name string, index int, end int) string {
	if len(name) == 0 {
		return ""
	}

	if index > len(name) {
		index = len(name)
	}

	if end > len(name) {
		end = len(name)
	}

	if index+end > len(name) {
		end = len(name) - index
	}

	prefix := name[:index]
	suffix := ""
	if end > 0 {
		suffix = name[len(name)-end:]
	}

	masked := strings.Repeat("*", len(name)-index-end)
	return prefix + masked + suffix
}

// Right only shows the last n characters, others are masked with asterisks
// Example: Right("1234567890", 4) returns "******7890"
func Right(num string, end int) string {
	if len(num) == 0 {
		return ""
	}

	if end > len(num) {
		end = len(num)
	}

	suffix := num[len(num)-end:]
	return strings.Repeat("*", len(num)-end) + suffix
}

// MobileEncrypt masks the middle 4 digits of a mobile phone number
// Example: MobileEncrypt("13812345678") returns "138****5678"
func MobileEncrypt(mobile string) string {
	if len(mobile) != 11 {
		return mobile
	}

	return mobile[:3] + "****" + mobile[7:]
}

// IDEncrypt masks the middle part of an ID card number, showing only first 3 and last 4 digits
// Example: IDEncrypt("123456789012345678") returns "123***********5678"
func IDEncrypt(id string) string {
	if len(id) < 8 {
		return id
	}

	return id[:3] + strings.Repeat("*", len(id)-7) + id[len(id)-4:]
}

// IDPassport masks a passport number, showing only first 2 and last 3 digits
// Example: IDPassport("AB1234567") returns "AB*****567"
func IDPassport(id string) string {
	if len(id) < 8 {
		return id
	}

	return id[:2] + strings.Repeat("*", len(id)-5) + id[len(id)-3:]
}

// IDPassportWithSize masks the last n digits of an ID/passport
// Example: IDPassportWithSize("AB1234567", 3) returns "AB123****"
func IDPassportWithSize(id string, sensitiveSize int) string {
	if len(id) == 0 {
		return ""
	}

	if sensitiveSize > len(id) {
		sensitiveSize = len(id)
	}

	return id[:len(id)-sensitiveSize] + strings.Repeat("*", sensitiveSize)
}
