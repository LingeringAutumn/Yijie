package mw

import (
	"context"
	"github.com/LingeringAutumn/Yijie/app/gateway/pack"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"

	metainfoContext "github.com/LingeringAutumn/Yijie/pkg/base/context"
)

// Auth 函数是一个中间件处理器，负责校验用户身份。
// 它会从请求头中提取 token 并进行相应处理，在调用 Next 方法时会携带 token 类型。
// 该函数返回一个 app.HandlerFunc 类型的处理器函数。
func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从请求上下文中的请求头里获取名为 constants.AuthHeader 的值，并将其转换为字符串类型，赋值给 token 变量。
		// 这里假设 constants.AuthHeader 是一个定义好的常量，表示存储 token 的请求头名称。
		token := string(c.GetHeader(constants.AuthHeader))
		// 调用 utils.CheckToken 函数来校验 token 的有效性，该函数返回三个值：token 类型、用户 ID（uid）和错误信息。
		// 这里忽略了 token 类型（用 _ 表示丢弃该返回值），只关注用户 ID 和错误信息。
		_, uid, err := utils.CheckToken(token)
		if err != nil {
			// 如果校验 token 时出现错误，调用 pack.RespError 函数将错误信息封装并返回给客户端。
			// 然后调用 c.Abort() 方法中止后续的处理流程，并返回，不再执行后续代码。
			pack.RespError(c, err)
			c.Abort()
			return
		}

		// 调用 utils.CreateAllToken 函数，传入用户 ID（uid），生成访问令牌（access）和刷新令牌（refresh）。
		// 该函数返回生成的令牌以及可能出现的错误信息。
		access, refresh, err := utils.CreateAllToken(uid)
		if err != nil {
			// 如果生成令牌时出现错误，同样调用 pack.RespError 函数封装并返回错误信息，
			// 调用 c.Abort() 中止后续处理流程并返回。
			pack.RespError(c, err)
			c.Abort()
			return
		}

		// 实现规范化服务透传，不需要中间进行编解码。
		// 调用 metainfoContext.WithLoginData 函数，将用户 ID（uid）添加到上下文中，用于后续服务间传递用户身份信息。
		ctx = metainfoContext.WithLoginData(ctx, uid)
		// 调用 metainfoContext.SetStreamLoginData 函数，在流式传输的上下文中也设置用户 ID（uid），
		// 以确保在流式传输场景下也能传递用户身份信息。
		ctx = metainfoContext.SetStreamLoginData(ctx, uid)
		// 将生成的访问令牌（access）设置到请求头中，键为 constants.AccessTokenHeader。
		// 这里假设 constants.AccessTokenHeader 是一个定义好的常量，表示存储访问令牌的请求头名称。
		c.Header(constants.AccessTokenHeader, access)
		// 将生成的刷新令牌（refresh）设置到请求头中，键为 constants.RefreshTokenHeader。
		// 这里假设 constants.RefreshTokenHeader 是一个定义好的常量，表示存储刷新令牌的请求头名称。
		c.Header(constants.RefreshTokenHeader, refresh)
		// 调用 c.Next(ctx) 方法，将处理流程传递给下一个中间件或处理器，同时将更新后的上下文（ctx）传递下去。
		c.Next(ctx)
	}
}
