package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func main() {
	// 生成RSA密钥对，密钥长度为2048位，通常2048位在当前具有较好的安全性
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("生成密钥对失败:", err)
		return
	}

	// 将私钥转换为PKCS#8格式的字节切片
	pkcs8PrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		fmt.Println("转换私钥为PKCS#8格式失败:", err)
		return
	}

	// 创建PEM格式的私钥块
	privatePemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8PrivateKey,
	}

	// 使用PEM编码私钥
	privatePemEncoded := pem.EncodeToMemory(privatePemBlock)

	// 提取公钥并将其转换为标准格式
	publicKey := privateKey.PublicKey
	pkixPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		fmt.Println("转换公钥格式失败:", err)
		return
	}

	// 创建PEM格式的公钥块
	publicPemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkixPublicKey,
	}

	// 使用PEM编码公钥
	publicPemEncoded := pem.EncodeToMemory(publicPemBlock)

	// 输出PKCS#8格式的私钥和公钥
	fmt.Println("PKCS#8格式私钥:")
	fmt.Println(string(privatePemEncoded))
	fmt.Println("公钥:")
	fmt.Println(string(publicPemEncoded))
}
