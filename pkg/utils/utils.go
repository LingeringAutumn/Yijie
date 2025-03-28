package utils

import (
	"bytes"
	"errors"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"io"
	"mime/multipart"
)

// CheckImageFileType 检查文件格式是否合规，接收一个 multipart.FileHeader 指针作为参数，返回文件类型和是否合规的布尔值
func CheckImageFileType(header *multipart.FileHeader) (string, bool) {
	// 打开文件
	file, err := header.Open()
	if err != nil {
		// 若打开文件失败，返回空字符串和 false
		return "", false
	}
	// 使用 defer 确保文件在函数结束时关闭，并捕获并处理关闭文件时可能发生的错误
	defer func() {
		if err := file.Close(); err != nil {
			logger.Errorf("utils.CheckImageFileType: failed to close file: %v", err.Error())
		}
	}()

	// 创建一个字节切片，用于存储读取的文件头部信息
	buffer := make([]byte, constants.CheckFileTypeBufferSize)
	// 读取文件头部信息到 buffer 中
	_, err = file.Read(buffer)
	if err != nil {
		// 若读取文件头部信息失败，返回空字符串和 false
		return "", false
	}

	// 使用 filetype 库的 Match 函数来判断文件类型
	kind, _ := filetype.Match(buffer)

	// 检查文件类型是否为 jpg 或 png
	switch kind {
	case types.Get("jpg"):
		// 若为 jpg 类型，返回 "jpg" 和 true
		return "jpg", true
	case types.Get("png"):
		// 若为 png 类型，返回 "png" 和 true
		return "png", true
	default:
		// 若不是 jpg 或 png 类型，返回空字符串和 false
		return "", false
	}
}

// GetImageFileType 获得图片格式，接收一个字节切片指针作为参数，返回文件类型和可能的错误
func GetImageFileType(fileBytes *[]byte) (string, error) {
	// 截取字节切片的前 constants.CheckFileTypeBufferSize 个字节作为文件头部信息
	buffer := (*fileBytes)[:constants.CheckFileTypeBufferSize]

	// 使用 filetype 库的 Match 函数来判断文件类型
	kind, _ := filetype.Match(buffer)

	// 检查文件类型是否为 jpg 或 png
	switch kind {
	case types.Get("jpg"):
		// 若为 jpg 类型，返回 "jpg" 和 nil
		return "jpg", nil
	case types.Get("png"):
		// 若为 png 类型，返回 "png" 和 nil
		return "png", nil
	default:
		// 若不是 jpg 或 png 类型，返回空字符串和内部服务错误
		return "", errno.InternalServiceError
	}
}

// FileToBytes 将文件转换为字节切片，接收一个 multipart.FileHeader 指针作为参数，返回二维字节切片和可能的错误
func FileToBytes(file *multipart.FileHeader) (ret []byte, err error) {
	// 检查文件是否为空
	if file == nil {
		// 若文件为空，返回 nil 和参数缺失错误
		return nil, errno.ParamMissingError.WithError(err)
	}

	// 打开文件
	fileOpen, err := file.Open()
	if err != nil {
		// 若打开文件失败，返回 nil 和操作系统操作错误，并附带错误信息
		return nil, errno.OSOperationError.WithMessage(err.Error())
	}
	// 确保文件在函数结束时关闭
	defer fileOpen.Close()

	// 使用 bytes.Buffer 存储所有读取的数据
	var buffer bytes.Buffer
	buf := make([]byte, constants.FileStreamBufferSize)

	for {
		n, err := fileOpen.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errno.InternalServiceError.WithMessage(err.Error())
		}
		// 只写入实际读取的 n 个字节
		buffer.Write(buf[:n])
	}

	return buffer.Bytes(), nil
}
