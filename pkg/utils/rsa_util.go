package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// RSAUtil 提供RSA加密和解密功能
type RSAUtil struct {
	PublicKey  string
	PrivateKey string
}

// NewRSAUtil 创建一个新的RSAUtil实例
func NewRSAUtil(publicKey, privateKey string) *RSAUtil {
	return &RSAUtil{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

// GenerateKeyPair 生成RSA密钥对
func GenerateKeyPair(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// 将私钥转换为PEM格式
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	}
	privateKeyPEM := pem.EncodeToMemory(privateKeyBlock)

	// 将公钥转换为PEM格式
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyDER,
	}
	publicKeyPEM := pem.EncodeToMemory(publicKeyBlock)

	return string(publicKeyPEM), string(privateKeyPEM), nil
}

// LoadPublicKeyFromBase64 从Base64字符串加载公钥
func LoadPublicKeyFromBase64(publicKeyBase64 string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("解码公钥失败: %w", err)
	}

	// 尝试直接解析公钥
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		// 尝试解析PEM格式
		block, _ := pem.Decode(publicKeyBytes)
		if block != nil {
			publicKey, err = x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("解析公钥失败: %w", err)
			}
		} else {
			// 尝试将Base64解码的数据包装成PKIX格式
			pubKey, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
			if err != nil {
				return nil, fmt.Errorf("解析公钥失败，尝试了多种格式: %w", err)
			}
			return pubKey, nil
		}
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("公钥类型错误，不是RSA公钥")
	}

	return rsaPublicKey, nil
}

// LoadPrivateKeyFromBase64 从Base64字符串加载私钥
func LoadPrivateKeyFromBase64(privateKeyBase64 string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("解码私钥失败: %w", err)
	}

	// 首先尝试PKCS1格式（最常见）
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err == nil {
		return privateKey, nil
	}

	// 尝试PKCS8格式
	pk8, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err == nil {
		privateKey, ok := pk8.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("私钥类型错误，不是RSA私钥")
		}
		return privateKey, nil
	}

	// 尝试PEM格式
	block, _ := pem.Decode(privateKeyBytes)
	if block != nil {
		// 尝试PKCS1
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err == nil {
			return privateKey, nil
		}

		// 尝试PKCS8
		pk8, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err == nil {
			privateKey, ok := pk8.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("私钥类型错误，不是RSA私钥")
			}
			return privateKey, nil
		}
	}

	return nil, fmt.Errorf("解析私钥失败，尝试了PKCS1和PKCS8格式: %w", err)
}

// EncryptWithPublicKey 使用公钥加密数据
func EncryptWithPublicKey(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha1.New(), rand.Reader, publicKey, data, nil)
}

// DecryptWithPrivateKey 使用私钥解密数据
func DecryptWithPrivateKey(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptOAEP(sha1.New(), rand.Reader, privateKey, ciphertext, nil)
}

// EncryptBase64 加密字符串并返回Base64编码的结果
func (r *RSAUtil) EncryptBase64(plainText string) (string, error) {
	publicKey, err := LoadPublicKeyFromBase64(r.PublicKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := EncryptWithPublicKey([]byte(plainText), publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptBase64 解密Base64编码的密文
func (r *RSAUtil) DecryptBase64(cipherTextBase64 string) (string, error) {
	privateKey, err := LoadPrivateKeyFromBase64(r.PrivateKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	plaintext, err := DecryptWithPrivateKey(ciphertext, privateKey)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// SignWithPrivateKey 使用私钥对数据进行签名
func SignWithPrivateKey(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha1.Sum(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed[:])
}

// VerifyWithPublicKey 使用公钥验证签名
func VerifyWithPublicKey(data []byte, signature []byte, publicKey *rsa.PublicKey) error {
	hashed := sha1.Sum(data)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashed[:], signature)
}

// SignBase64 对数据进行签名并返回Base64编码的结果
func (r *RSAUtil) SignBase64(data string) (string, error) {
	privateKey, err := LoadPrivateKeyFromBase64(r.PrivateKey)
	if err != nil {
		return "", err
	}

	signature, err := SignWithPrivateKey([]byte(data), privateKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyBase64 验证Base64编码的签名
func (r *RSAUtil) VerifyBase64(data string, signatureBase64 string) error {
	publicKey, err := LoadPublicKeyFromBase64(r.PublicKey)
	if err != nil {
		return err
	}

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return err
	}

	return VerifyWithPublicKey([]byte(data), signature, publicKey)
}
