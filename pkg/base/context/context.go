package context

import (
	"context"
	// 引入字节跳动的云元信息包，用于处理上下文的持久化值等相关操作
	// 这里的 metainfo.WithPersistentValue 和 metainfo.GetPersistentValue 用于设置和获取上下文中的持久化值
	// 该包可能提供了一些云环境下上下文元信息管理的功能
	"github.com/bytedance/gopkg/cloud/metainfo"
	// 引入 Kitex 框架中关于网络协议（nphttp2）元数据相关的包
	// 用于处理 Kitex 框架中在网络传输（如 HTTP/2）时的元数据操作
	// 这里的 metadata.FromIncomingContext 和 metadata.AppendToOutgoingContext 用于从传入上下文获取元数据和向传出上下文追加元数据
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/metadata"
)

// newContext 函数用于创建一个新的上下文，在原有的上下文基础上添加一个持久化的值。
// ctx 是传入的原始上下文对象，后续操作基于此上下文。
// key 是用于标识该值的键，在上下文中通过这个键来访问对应的值。
// value 是要添加到上下文中的具体值，这里是字符串类型。
// 返回值是添加了指定键值对后的新上下文对象。
func newContext(ctx context.Context, key string, value string) context.Context {
	return metainfo.WithPersistentValue(ctx, key, value)
}

// fromContext 函数用于从给定的上下文中获取指定键对应的持久化值。
// ctx 是传入的上下文对象，期望从中获取值。
// key 是用于标识要获取的值的键。
// 返回值第一个是获取到的字符串类型的值，如果没有获取到对应的值则返回空字符串。
// 第二个返回值是一个布尔类型，表示是否成功获取到了值，true 表示成功获取，false 表示未获取到。
func fromContext(ctx context.Context, key string) (string, bool) {
	return metainfo.GetPersistentValue(ctx, key)
}

// streamFromContext 函数用于在流式传输的上下文中获取指定键对应的值。
// 注释中提到流式传输 ctx 不能用传统方式传递，具体参考链接：https://www.cloudwego.io/zh/docs/kitex/tutorials/advanced-feature/metainfo/#kitex-grpc-metadata
// 由于目前只传输了 uid，所以返回的第一个值只取了获取到的数组中的第一位。
// 如果后续传输的数据结构发生变化，需要对这部分代码进行修正。
// ctx 是传入的流式传输上下文对象。
// key 是用于标识要获取的值的键。
// 返回值第一个是获取到的字符串类型的值（从数组第一位取出），如果没有获取到对应的值则返回空字符串。
// 第二个返回值是一个布尔类型，表示是否成功获取到了值，true 表示成功获取，false 表示未获取到。
func streamFromContext(ctx context.Context, key string) (string, bool) {
	md, success := metadata.FromIncomingContext(ctx)
	return md.Get(key)[0], success
}

// streamAppendContext 函数用于在流式传输的上下文中追加一个指定键值对。
// ctx 是传入的流式传输上下文对象。
// key 是用于标识要追加的值的键。
// value 是要追加到上下文中的具体值，这里是字符串类型。
// 返回值是追加了指定键值对后的新上下文对象。
func streamAppendContext(ctx context.Context, key string, value string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, key, value)
}
