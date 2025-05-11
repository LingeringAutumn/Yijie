package pack

import (
	dmodel "github.com/LingeringAutumn/Yijie/app/user/domain/model"
	kmodel "github.com/LingeringAutumn/Yijie/kitex_gen/model"
)

// BuildUser 将 entities 定义的 User 实体转换成 idl 定义的 RPC 交流实体，类似 dto
func BuildUser(user *dmodel.User) *kmodel.UserInfo {
	return &kmodel.UserInfo{
		UserId: user.Uid,
		Name:   user.Username,
	}
}

func BuildUserProfileRequest(user *dmodel.UserProfileRequest) *kmodel.UserProfileReq {
	return &kmodel.UserProfileReq{
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Bio:      user.Bio,
	}
}

func BuildUpdateUserProfileResponse(user *dmodel.UpdateUserProfileResponse) *kmodel.UserProfileResp {
	return &kmodel.UserProfileResp{
		Username:        user.Username,
		Email:           user.Email,
		Phone:           user.Phone,
		Avatar:          user.Avatar,
		Bio:             user.Bio,
		MembershipLevel: user.MembershipLevel,
		Point:           user.Point,
		Team:            user.Team,
	}
}

func BuildGetUserProfileResponse(user *dmodel.GetUserProfileResponse) *kmodel.UserProfileResp {
	return &kmodel.UserProfileResp{
		Username:        user.Username,
		Email:           user.Email,
		Phone:           user.Phone,
		Avatar:          user.Avatar,
		Bio:             user.Bio,
		MembershipLevel: user.MembershipLevel,
		Point:           user.Point,
		Team:            user.Team,
	}
}

func BuildImage(image *dmodel.Image) *kmodel.Image {
	return &kmodel.Image{
		ImageId:  image.ImageID,
		ImageUrl: image.Url,
	}
}
