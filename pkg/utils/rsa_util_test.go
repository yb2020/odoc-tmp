package utils

import (
	"testing"
)

// 使用配置中的RSA密钥对
var publicKey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLd1/A1PM9nXJwNpewStEsNcB126iLLPTnAnSo95QIg89J2C8mi04n+STN24QlDERVeIiBESjqImhlyhTNOdKmFaVEnPglINALaa4visQMccBwVYqCR2AebkhP83Dx0GJp0ywecJqd62kYLhSmnikDh+GYVCA/ohn80+1L4V6LhwIDAQAB"
var privateKey = "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAMt3X8DU8z2dcnA2l7BK0Sw1wHXbqIss9OcCdKj3lAiDz0nYLyaLTif5JM3bhCUMRFV4iIERKOoiaGXKFM050qYVpUSc+CUg0Atpri+KxAxxwHBVioJHYB5uSE/zcPHQYmnTLB5wmp3raRguFKaeKQOH4ZhUID+iGfzT7UvhXouHAgMBAAECgYEAjxdS5gBdWHXEJ5qdL0ROuvLKeZiTfd2OFnCprrL/DsX0IBDDiC3sNzyGX6gD1TI9VIbCKVLyHUc5eGyYGIST2SzLLNP/nYa9j2ZvCa/1c3sZ/vyyoo4GXBKW+lXkRoJZ4DTBz/6+4EzNT+glGEzOZSq/xQL9a4i+0nFg0OoiGrECQQDoXptHQ03UlMBw1Oj6tjL2asX/3Vzec7Ey//xpwHqAUR4fNhsXKbBJNLKYkW4M7Z1h+mhJQP4YlV8Ckt7UWVmdAkEA4ChQiVN508NI5+fQ9eGYeHizo9EqKSS2xBG4z6BVgWkHdW4VNYWYgmnqItT3qOPDO4x5sF37LHsOb8olAXNScwJAdLjRFwLf3aC66fKI9ScAgncv7k6rj7JdmFit2hEtd7dHgjYTdZcjTiKCc9DZjvTs0YKPT/ytpnuhthFAjTo0oQJAE6x2JRdmgeeJ5pC6DlqWfzxYx+/7u1C1mc/UYKS53HnTZcMbqW7oS8nv+s6mTfRvljJmG8yj1uuWAMnFJbNxcQJBAOIuin6vurrE4v/BFoStKSVZCb8xZBO2YHz3Y9zMgEU4x+IHue3k97j8M7uaVY5a68udBOgIyjlHbd4oSsKV5TU="

func TestRSAEncryptionDecryption(t *testing.T) {
	rsaUtil := &RSAUtil{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	// 测试加密和解密
	originalText := "Hello, RSA encryption in Go!"

	// 加密
	encryptedBase64, err := rsaUtil.EncryptBase64(originalText)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	// 解密
	decryptedText, err := rsaUtil.DecryptBase64(encryptedBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	// 验证结果
	if decryptedText != originalText {
		t.Errorf("解密结果与原文不匹配, 期望: %s, 实际: %s", originalText, decryptedText)
	}
}

func TestRSASignatureVerification(t *testing.T) {
	rsaUtil := &RSAUtil{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	// 测试签名和验证
	originalText := "Hello, RSA signature in Go!"

	// 签名
	signatureBase64, err := rsaUtil.SignBase64(originalText)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}

	// 验证签名
	err = rsaUtil.VerifyBase64(originalText, signatureBase64)
	if err != nil {
		t.Fatalf("验证签名失败: %v", err)
	}

	// 验证错误的签名
	err = rsaUtil.VerifyBase64("Modified text", signatureBase64)
	if err == nil {
		t.Error("验证错误的签名应该失败，但却成功了")
	}
}

func TestGenerateKeyPair(t *testing.T) {
	// 生成密钥对
	publicKey, privateKey, err := GenerateKeyPair(1024)
	if err != nil {
		t.Fatalf("生成密钥对失败: %v", err)
	}

	if publicKey == "" || privateKey == "" {
		t.Error("生成的密钥对不应为空")
	}
}
