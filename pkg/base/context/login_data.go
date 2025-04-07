package context

import (
	"context"
	"strconv"

	// 这里的 constants.LoginDataKey 是用于标识登录数据在 context 中存储的键名
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
)

// WithLoginData 将LoginData加入到context中，通过metainfo传递到RPC server
// ctx 是传入的上下文对象，后续操作都基于此上下文
// uid 是用户的ID，类型为int64
// 返回值是添加了登录数据后的新上下文对象
func WithLoginData(ctx context.Context, uid int64) context.Context {
	// 调用 newContext 函数，将用户ID转换为字符串后，以 constants.LoginDataKey 为键，添加到上下文中
	return newContext(ctx, constants.LoginDataKey, strconv.FormatInt(uid, 10))
}

// GetLoginData 从context中取出LoginData
// ctx 是传入的上下文对象，期望从中获取登录数据
// 返回值第一个是取出并转换后的用户ID（int64类型），第二个是错误信息，如果获取或转换过程出错会返回非nil的错误
func GetLoginData(ctx context.Context) (int64, error) {
	// 调用 fromContext 函数，尝试从上下文中以 constants.LoginDataKey 为键取出对应的值
	// 这里 fromContext 函数是自定义的用于从上下文中取值的函数
	// user 是取出的值（字符串类型），ok 表示是否成功取出
	user, ok := fromContext(ctx, constants.LoginDataKey)
	if !ok {
		// 如果没有成功取出值，使用 errno 包创建一个新的错误对象返回
		// 错误码为 errno.ParamMissingErrorCode，错误信息为 "Failed to get header in context"
		return -1, errno.NewErrNo(errno.ParamMissingErrorCode, "Failed to get header in context")
	}

	// 将取出的字符串类型的用户ID转换为int64类型
	value, err := strconv.ParseInt(user, 10, 64)
	if err != nil {
		// 如果转换过程出错，使用 errno 包创建一个新的错误对象返回
		// 错误码为 errno.InternalServiceErrorCode，错误信息为 "Failed to get header in context when parse loginData"
		return -1, errno.NewErrNo(errno.InternalServiceErrorCode, "Failed to get header in context when parse loginData")
	}
	return value, nil
}

// GetStreamLoginData 流式传输传递ctx, 获取loginData
// ctx 是传入的上下文对象，用于流式传输场景下获取登录数据
// 返回值第一个是取出并转换后的用户ID（int64类型），第二个是错误信息，如果获取或转换过程出错会返回非nil的错误
func GetStreamLoginData(ctx context.Context) (int64, error) {
	// 调用 streamFromContext 函数，尝试从上下文中以 constants.LoginDataKey 为键取出对应的值
	// 这里假设 streamFromContext 函数是自定义的用于在流式传输场景下从上下文中取值的函数
	// uid 是取出的值（字符串类型），success 表示是否成功取出
	uid, success := streamFromContext(ctx, constants.LoginDataKey)
	if !success {
		// 如果没有成功取出值，使用 errno 包创建一个新的错误对象返回
		// 错误码为 errno.ParamMissingErrorCode，错误信息为 "Failed to get info in context"
		return -1, errno.NewErrNo(errno.ParamMissingErrorCode, "Failed to get info in context")
	}

	// 将取出的字符串类型的用户ID转换为int64类型
	value, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		// 如果转换过程出错，使用 errno 包创建一个新的错误对象返回
		// 错误码为 errno.InternalServiceErrorCode，错误信息为 "Failed to get info in context when parse loginData"
		return -1, errno.NewErrNo(errno.InternalServiceErrorCode, "Failed to get info in context when parse loginData")
	}
	return value, nil
}

// SetStreamLoginData 流式传输传递ctx, 设置ctx值
// ctx 是传入的上下文对象，用于流式传输场景下设置登录数据
// uid 是用户的ID，类型为int64
// 返回值是添加了登录数据后的新上下文对象
func SetStreamLoginData(ctx context.Context, uid int64) context.Context {
	// 将用户ID转换为字符串
	value := strconv.FormatInt(uid, 10)
	// 调用 streamAppendContext 函数，以 constants.LoginDataKey 为键，将转换后的字符串值添加到上下文中
	// 这里假设 streamAppendContext 函数是自定义的用于在流式传输场景下更新上下文的函数
	return streamAppendContext(ctx, constants.LoginDataKey, value)
}
