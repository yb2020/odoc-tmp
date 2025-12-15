package idgen

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base62Count = int64(len(base62Chars))

const bucketChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const bucketCount = int64(len(bucketChars))

// Encode 将 uint64 数字编码为 Base62 字符串
func EncodeBase62(id int64) (string, error) {
	if id == 0 {
		return "0", errors.New("id cannot be 0")
	}

	var sb strings.Builder
	for id > 0 {
		r := id % base62Count
		sb.WriteByte(base62Chars[r])
		id /= base62Count
	}

	// 结果需要反转
	return reverse(sb.String()), nil
}

// Decode 将 Base62 字符串解码为 uint64 数字
func DecodeBase62(s string) (int64, error) {
	var n int64
	for _, char := range s {
		i := strings.IndexRune(base62Chars, char)
		if i == -1 {
			return 0, fmt.Errorf("invalid character in Base62 string: %c", char)
		}
		n = n*base62Count + int64(i)
	}
	return n, nil
}

// EncodeUUIDToBase62 将 UUID 编码为 Base62 字符串
// UUID 是 128 位，编码后约为 22 个字符
func EncodeUUIDToBase62(id uuid.UUID) string {
	// 将 UUID 的 16 字节转换为大整数
	bytes := id[:]
	bigInt := new(big.Int).SetBytes(bytes)

	if bigInt.Sign() == 0 {
		return "0"
	}

	var sb strings.Builder
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for bigInt.Cmp(zero) > 0 {
		bigInt.DivMod(bigInt, base, mod)
		sb.WriteByte(base62Chars[mod.Int64()])
	}

	// 结果需要反转
	return reverse(sb.String())
}

// DecodeBase62ToUUID 将 Base62 字符串解码为 UUID
func DecodeBase62ToUUID(s string) (uuid.UUID, error) {
	var id uuid.UUID
	if s == "" {
		return id, errors.New("empty string")
	}

	bigInt := new(big.Int)
	base := big.NewInt(62)

	for _, char := range s {
		i := strings.IndexRune(base62Chars, char)
		if i == -1 {
			return id, fmt.Errorf("invalid character in Base62 string: %c", char)
		}
		bigInt.Mul(bigInt, base)
		bigInt.Add(bigInt, big.NewInt(int64(i)))
	}

	// 将大整数转换回 16 字节
	bytes := bigInt.Bytes()
	// UUID 是 16 字节，需要左侧补零
	if len(bytes) < 16 {
		padded := make([]byte, 16)
		copy(padded[16-len(bytes):], bytes)
		bytes = padded
	} else if len(bytes) > 16 {
		return id, errors.New("decoded value too large for UUID")
	}

	copy(id[:], bytes)
	return id, nil
}

// EncodeStringUUIDToBase62 将字符串格式的 UUID 编码为 Base62
func EncodeStringUUIDToBase62(uuidStr string) (string, error) {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID string: %w", err)
	}
	return EncodeUUIDToBase62(id), nil
}

// DecodeBase62ToStringUUID 将 Base62 字符串解码为字符串格式的 UUID
func DecodeBase62ToStringUUID(s string) (string, error) {
	id, err := DecodeBase62ToUUID(s)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// reverse 用于反转字符串
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func GetShardDirectory(uniqueID string) string {
	// 1. 对输入ID进行SHA1哈希计算
	hasher := sha1.New()
	hasher.Write([]byte(uniqueID))
	hashBytes := hasher.Sum(nil)

	// 2. 将哈希结果的字节转换为一个大整数
	hashInt := new(big.Int)
	hashInt.SetBytes(hashBytes)

	// 3. 对大整数进行取模运算
	// remainder = hashInt % shardCount
	remainder := new(big.Int)
	remainder.Mod(hashInt, big.NewInt(bucketCount))

	// 4. 使用模的结果作为索引，从字符表中获取字符
	shardIndex := remainder.Int64()
	return string(bucketChars[shardIndex])
}

func testMain() {
	// 假设这是你的雪花算法生成的ID
	snowflakeID := int64(372077392774041600)

	// 编码
	shortID, _ := EncodeBase62(snowflakeID)

	fmt.Printf("原始雪花ID: %d\n", snowflakeID)
	fmt.Printf("Base62短ID: %s\n", shortID)
	fmt.Printf("短ID长度: %d\n", len(shortID))

	// 解码（验证）
	originalID, _ := DecodeBase62(shortID)
	fmt.Printf("解码后的ID: %d\n", originalID)
	fmt.Printf("解码是否一致: %v\n", originalID == snowflakeID)

	// UUID v7 Base62 编码示例
	uuidV7, _ := uuid.NewV7()
	base62UUID := EncodeUUIDToBase62(uuidV7)
	fmt.Printf("原始UUID v7: %s\n", uuidV7.String())
	fmt.Printf("Base62编码: %s\n", base62UUID)
	fmt.Printf("编码长度: %d\n", len(base62UUID))

	// 解码验证
	decodedUUID, _ := DecodeBase62ToUUID(base62UUID)
	fmt.Printf("解码后UUID: %s\n", decodedUUID.String())
	fmt.Printf("解码是否一致: %v\n", decodedUUID == uuidV7)
}
