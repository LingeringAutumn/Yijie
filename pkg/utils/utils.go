package utils

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"

	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
)

// GetMysqlDSN 会拼接 mysql 的 DSN（Data Source Name），DSN 是一种用于连接 MySQL 数据库的字符串，包含了数据库连接所需的各种信息。
func GetMysqlDSN() (string, error) {
	// 检查配置中的 MySQL 部分是否为空。
	// 如果为空，说明没有找到 MySQL 数据库的配置信息，无法拼接 DSN。
	// 此时返回一个空字符串和一个错误信息，提示配置未找到。
	if config.Mysql == nil {
		return "", errors.New("config not found")
	}

	// 开始拼接 DSN 字符串。
	// 使用 strings.Join 函数将多个字符串片段连接成一个完整的 DSN 字符串。
	// 下面是 DSN 字符串的各个组成部分：
	// 1. config.Mysql.Username：MySQL 数据库的用户名，用于身份验证。
	// 2. ":"：用户名和密码之间的分隔符。
	// 3. config.Mysql.Password：MySQL 数据库的密码，用于身份验证。
	// 4. "@tcp("：表示使用 TCP 协议进行连接。
	// 5. config.Mysql.Addr：MySQL 数据库服务器的地址，通常是 IP 地址和端口号，例如 "127.0.0.1:3306"。
	// 6. ")/"：结束 TCP 地址的指定，并指定要连接的数据库。
	// 7. config.Mysql.Database：要连接的 MySQL 数据库的名称。
	// 8. "?charset=" + config.Mysql.Charset + "&parseTime=true"：
	//    - "?charset="：指定字符集的参数起始符号。
	//    - config.Mysql.Charset：指定连接使用的字符集，例如 "utf8mb4"。
	//    - "&parseTime=true"：表示允许将 MySQL 的日期和时间类型自动解析为 Go 语言的时间类型。
	dsn := strings.Join([]string{
		config.Mysql.Username, ":", config.Mysql.Password,
		"@tcp(", config.Mysql.Addr, ")/",
		config.Mysql.Database, "?charset=" + config.Mysql.Charset + "&parseTime=true",
	}, "")

	// 返回拼接好的 DSN 字符串和一个 nil 错误，表示没有发生错误。
	return dsn, nil
}

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

// AddrCheck 会检查当前的监听地址是否已被占用
// 它通过尝试在指定地址上创建一个 TCP 监听器来判断地址是否可用
// 如果能成功创建监听器，说明地址未被占用；反之则表示地址已被占用
func AddrCheck(addr string) bool {
	// 尝试在指定地址上创建一个 TCP 监听器
	// net.Listen 函数接收两个参数，第一个参数 "tcp" 表示使用 TCP 协议
	// 第二个参数 addr 是要监听的地址
	// 函数返回一个监听器对象 l 和一个错误对象 err
	l, err := net.Listen("tcp", addr)
	// 如果创建监听器时出现错误，说明该地址可能已被其他程序占用
	if err != nil {
		// 返回 false 表示地址不可用
		return false
	}
	// 使用 defer 关键字确保在函数返回之前关闭监听器
	// 这是为了避免资源泄漏，即使后续代码出现异常，也能保证监听器被关闭
	defer func() {
		// 尝试关闭监听器
		if err := l.Close(); err != nil {
			// 如果关闭监听器时出现错误，使用日志记录错误信息
			// logger.Errorf 是一个自定义的日志记录函数，用于记录错误级别的日志
			logger.Errorf("utils.AddrCheck: failed to close listener: %v", err.Error())
		}
	}()
	// 如果能成功创建监听器，说明地址未被占用，返回 true
	return true
}

// GetAvailablePort 会尝试获取可用的监听地址
// 该函数会遍历配置文件中提供的地址列表，检查每个地址是否可用
// 如果找到可用地址，则返回该地址；如果没有找到可用地址，则返回错误
func GetAvailablePort() (string, error) {
	// 检查配置文件中定义的地址列表是否为空
	// config.Service.AddrList 是配置文件中存储的地址列表，如果为空则说明没有提供有效的地址
	if config.Service.AddrList == nil {
		// 如果地址列表为空，返回一个错误，提示配置文件中的地址列表为空
		return "", errors.New("utils.GetAvailablePort: config.Service.AddrList is nil")
	}
	// 遍历配置文件中提供的地址列表
	for _, addr := range config.Service.AddrList {
		// 调用 AddrCheck 函数检查当前地址是否可用
		// AddrCheck 函数应该是一个自定义函数，用于检查地址是否可以被监听
		if ok := AddrCheck(addr); ok {
			// 如果当前地址可用，返回该地址和 nil 错误
			return addr, nil
		}
	}
	// 如果遍历完所有地址都没有找到可用地址，返回一个错误，提示没有从配置文件中找到可用端口
	return "", errors.New("utils.GetAvailablePort: not available port from config")
}

// CheckVideoFileType 检查视频格式
func CheckVideoFileType(header *multipart.FileHeader) (string, bool) {
	// 打开文件
	file, err := header.Open()
	if err != nil {
		// 若打开文件失败，返回空字符串和 false
		return "", false
	}
	// 使用 defer 确保文件在函数结束时关闭，并捕获并处理关闭文件时可能发生的错误
	defer func() {
		if err := file.Close(); err != nil {
			logger.Errorf("utils.CheckVideoFileType: failed to close file: %v", err.Error())
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
	case types.Get("mp4"):
		// 若为 mp4 类型，返回 "mp4" 和 true
		return "mp4", true
	default:
		// 若不是 mp4 或 avi 类型，返回空字符串和 false
		return "", false
	}
}
