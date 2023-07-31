package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func EncryptByPass(pass string, data []byte) string {
	bPass := []byte(pass)

	key := sha256.Sum256(bPass)
	block, _ := aes.NewCipher(key[:])
	gcm, _ := cipher.NewGCM(block)

	nonce := md5.Sum(bPass)

	ret := gcm.Seal(nil, nonce[:12], data, nil)
	return Base64Encode(ret)
}

func DecryptByPass(pass string, b64cipherData string) ([]byte, error) {
	cipherText, err := Base64Decode(b64cipherData)
	if err != nil {
		return nil, err
	}

	bPass := []byte(pass)

	key := sha256.Sum256(bPass)
	block, _ := aes.NewCipher(key[:])
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := md5.Sum(bPass)
	return gcm.Open(nil, nonce[:12], cipherText, nil)
}

func GenerateECDHKeyPair() (priKey, pubKey []byte, err error) {
	privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	publicKey := privateKey.PublicKey()

	return privateKey.Bytes(), publicKey.Bytes(), nil
}

func Base64Encode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}

func ECDH(priKey, pubKey []byte) ([]byte, error) {
	privateKey, err := ecdh.X25519().NewPrivateKey(priKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := ecdh.X25519().NewPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	return privateKey.ECDH(publicKey)
}

func Encrypt(sharedKey, data []byte) (string, error) {
	if len(sharedKey) != 32 {
		return "", fmt.Errorf("invalid key size: %d", len(sharedKey))
	}

	block, _ := aes.NewCipher(sharedKey)
	gcm, _ := cipher.NewGCM(block)

	nonce := md5.Sum(sharedKey)

	ret := gcm.Seal(nil, nonce[:12], data, nil)
	return Base64Encode(ret), nil
}

func Decrypt(sharedKey []byte, b64cipherData string) ([]byte, error) {
	cipherText, err := Base64Decode(b64cipherData)
	if err != nil {
		return nil, err
	}

	block, _ := aes.NewCipher(sharedKey)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := md5.Sum(sharedKey)
	return gcm.Open(nil, nonce[:12], cipherText, nil)
}
