// Code generated by hertz generator.

package user_behavior

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	user_behavior "github.com/LingeringAutumn/Yijie/app/gateway/model/api/user_behavior"
)

// LikeVideo .
// @router api/v1/like/video [POST]
func LikeVideo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user_behavior.VideoLikeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(user_behavior.VideoLikeResponse)

	c.JSON(consts.StatusOK, resp)
}
