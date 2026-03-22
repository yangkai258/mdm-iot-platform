package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// AESKeyLength AES-256 requires 32 bytes
const AESKeyLength = 32

// getAESKey 获取 AES 加密密钥，从环境变量或使用默认密钥（仅用于开发）
func getAESKey() []byte {
	keyStr := os.Getenv("AES_ENCRYPTION_KEY")
	if keyStr == "" {
		// 开发环境默认密钥，仅作为 fallback
		// 生产环境务必设置 AES_ENCRYPTION_KEY 环境变量
		return []byte("mdm-secret-key-32-bytes-long!!")
	}
	// 确保持密钥长度正确（截断或填充）
	key := []byte(keyStr)
	if len(key) < AESKeyLength {
		// 填充到 32 字节
		padded := make([]byte, AESKeyLength)
		copy(padded, key)
		return padded
	}
	return key[:AESKeyLength]
}

// EncryptAES AES 加密（返回 base64 编码）
func EncryptAES(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key := getAESKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建 GCM 失败: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成 nonce 失败: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES AES 解密（输入 base64 编码）
func DecryptAES(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("Base64 解码失败: %w", err)
	}

	key := getAESKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建 GCM 失败: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("密文长度不足")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %w", err)
	}

	return string(plaintext), nil
}
