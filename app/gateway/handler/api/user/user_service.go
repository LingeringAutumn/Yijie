// Code generated by hertz generator.

package user

import (
	"context"

	user "github.com/LingeringAutumn/Yijie/app/gateway/model/api/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Register .
// @router api/v1/user/register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.RegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(user.RegisterResponse)

	c.JSON(consts.StatusOK, resp)
}

// Login .
// @router api/v1/user/login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(user.LoginResponse)

	c.JSON(consts.StatusOK, resp)
}

// UpdateUserProfile .
// @router api/v1/user/profile/update [PUT]
func UpdateUserProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UpdateUserProfileRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(user.UpdateUserProfileResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetUserProfile .
// @router api/v1/user/profile/get [GET]
func GetUserProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.GetUserProfileRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(user.GetUserProfileResponse)

	c.JSON(consts.StatusOK, resp)
}
