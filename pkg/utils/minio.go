// 定义包名为 utils
package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	// 导入 MinIO Go 客户端库
	"github.com/minio/minio-go"
)

// MinioClientGlobal 是一个全局的 MinioClient 实例，用于在整个程序中共享 MinIO 客户端
var MinioClientGlobal *MinioClient

// MinioClient 结构体封装了 MinIO 客户端，方便对 MinIO 进行操作
type MinioClient struct {
	// 指向 MinIO 客户端的指针
	Client *minio.Client
}

// InitMinioClient 用于初始化 MinIO 客户端，并设置全局的 MinioClient 实例
// 参数 endpoint 是 MinIO 服务的地址
// 参数 accessKeyID 是访问 MinIO 的密钥 ID
// 参数 secretAccessKey 是访问 MinIO 的密钥
// 返回值为错误信息，如果初始化成功则返回 nil
func InitMinioClient(endpoint, accessKeyID, secretAccessKey string) error {
	// 调用 minio.New 函数创建一个新的 MinIO 客户端实例
	// 最后一个参数 false 表示不使用 SSL 连接
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		// 如果创建客户端失败，使用 log.Fatalf 输出错误信息并终止程序
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}
	// 设置全局的 MinioClient 实例，将新创建的客户端封装到 MinioClient 结构体中
	MinioClientGlobal = &MinioClient{Client: client}
	return nil
}

// UploadFile 用于将文件上传到 MinIO 存储桶
// 参数 bucketName 是存储桶的名称
// 参数 objectName 是文件在存储桶中的对象名称
// 参数 Location 是存储桶的地理位置
// 参数 ContentType 是文件的内容类型
// 参数 file 是要上传的文件的字节切片
// 返回值为错误信息，如果上传成功则返回 nil
func (m *MinioClient) UploadFile(bucketName, objectName, Location, ContentType string, file []byte) error {
	// 创建一个字节读取器，用于读取要上传的文件内容
	reader := bytes.NewReader(file)
	// 检查指定的存储桶是否存在
	exist, err := m.Client.BucketExists(bucketName)
	if err != nil {
		// 处理检查存储桶是否存在时的错误
		return fmt.Errorf("failed to check if bucket %s exists: %w", bucketName, err)
	}
	if !exist {
		// 如果存储桶不存在，则创建该存储桶
		err := m.Client.MakeBucket(bucketName, Location)
		if err != nil {
			// 如果创建存储桶失败，返回错误信息
			return fmt.Errorf("failed to create bucket %s: %v", bucketName, err)
		}
	}
	// 创建 PutObjectOptions 结构体，设置文件的内容类型
	options := minio.PutObjectOptions{ContentType: ContentType}
	// 调用 MinIO 客户端的 PutObject 方法上传文件
	// -1 表示自动计算文件大小
	n, err := m.Client.PutObject(bucketName, objectName, reader, -1, options)
	if err != nil {
		// 如果上传失败，使用 log.Printf 输出错误信息，并返回错误信息
		log.Printf("Failed to upload %s: %v", objectName, err)
		return fmt.Errorf("failed to upload %s: %v", objectName, err)
	}
	// 如果上传成功，使用 log.Printf 输出上传成功的信息
	log.Printf("Successfully uploaded %s of size %d", objectName, n)
	return nil
}

// DownloadFile 下载文件
// 参数 bucketName 是存储桶的名称
// 参数 objectName 是文件在存储桶中的对象名称
// 参数 filePath 是要下载到本地的文件路径
// 返回值为错误信息，如果下载成功则返回 nil
func (m *MinioClient) DownloadFile(bucketName, objectName, filePath string) error {
	// 创建本地文件
	// file, err := os.Create(filePath)
	// 使用 os.OpenFile 并设置标志，避免不必要的文件覆盖
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return fmt.Errorf("failed to create local file %s: %w", filePath, err)
	}
	// 确保在函数结束时关闭文件
	defer file.Close()

	// 下载存储桶中的文件到本地
	err = m.Client.FGetObject(bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to download %s from bucket %s: %w", objectName, bucketName, err)
	}

	// 输出下载成功的信息
	fmt.Println("Successfully downloaded", objectName)
	return nil
}

// DeleteFile 删除文件
// 参数 bucketName 是存储桶的名称
// 参数 objectName 是文件在存储桶中的对象名称
// 返回值为布尔值表示是否删除成功，以及错误信息，如果删除成功则错误信息为 nil
func (m *MinioClient) DeleteFile(bucketName, objectName string) (bool, error) {
	// 删除存储桶中的文件
	err := m.Client.RemoveObject(bucketName, objectName)
	if err != nil {
		return false, fmt.Errorf("failed to delete %s from bucket %s: %w", objectName, bucketName, err)
	}
	return true, nil
}

// ListObjects 列出文件
// 参数 bucketName 是存储桶的名称
// 参数 prefix 是对象名称的前缀，用于过滤对象
// 返回值为对象名称的切片和错误信息，如果列出成功则错误信息为 nil
func (m *MinioClient) ListObjects(bucketName, prefix string) ([]string, error) {
	// 用于存储对象名称的切片
	var objectNames []string

	// 遍历 MinIO 客户端返回的对象迭代器
	for object := range m.Client.ListObjects(bucketName, prefix, true, nil) {
		if object.Err != nil {
			// 如果迭代过程中出现错误，返回 nil 和错误信息
			return nil, fmt.Errorf("failed to list objects in bucket %s with prefix %s: %w", bucketName, prefix, object.Err)
		}

		// 将对象的键（即对象名称）添加到切片中
		objectNames = append(objectNames, object.Key)
	}

	return objectNames, nil
}

// GetPresignedGetObject 返回对象的 url 地址，有效期时间为 expires
// 参数 bucketName 是存储桶的名称
// 参数 objectName 是文件在存储桶中的对象名称
// 参数 expires 是预签名 URL 的有效期
// 返回值为预签名 URL 的字符串和错误信息，如果生成成功则错误信息为 nil
func (m *MinioClient) GetPresignedGetObject(bucketName string, objectName string, expires time.Duration) (string, error) {
	// 调用 MinIO 客户端的 PresignedGetObject 方法生成预签名 URL
	object, err := m.Client.PresignedGetObject(bucketName, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to generate presigned URL for %s in bucket %s: %w", objectName, bucketName, err)
	}
	// 返回预签名 URL 的字符串形式
	return object.String(), nil
}
