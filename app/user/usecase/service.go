package usecase

import (
	"context"
	"fmt"

	"github.com/LingeringAutumn/Yijie/app/user/domain/model"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
)

func (uc *userUseCase) RegisterUser(ctx context.Context, u *model.User) (uid int64, err error) {
	// 注意: 这里使用 uc 调用了 DB, 但显然这个方法其他地方也可能会用的上, 所以可以考虑包装在 service 里面
	exist, err := uc.svc.IsUserExist(ctx, u.Username)
	if err != nil {
		// 这里返回了 fmt.Errorf 而不是 errno 的原因是 db.IsUserExist 返回的已经是 errno 了
		// 这里是用 %w 占位符做了一层 wrap, 其实这个 error 的底部(origin error) 还是 errno 类型的
		return 0, fmt.Errorf("check user exist failed: %w", err)
	}
	if exist {
		return 0, errno.NewErrNo(errno.ServiceUserExist, "user already exist")
	}
	u.Password, err = uc.svc.EncryptPassword(u.Password)
	if err != nil {
		return 0, err
	}
	// 这里没有直接调用 db.CreateUser 是因为 svc.CreateUser 包含了一点业务逻辑, 这些细节不需要被 useCase 知道
	uid, err = uc.svc.CreateUser(ctx, u)
	if err != nil {
		return
	}
	return uid, nil
}

func (uc *userUseCase) LoginUser(ctx context.Context, user *model.User) (*model.User, error) {
	userData, err := uc.svc.GetUserByName(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("get user info failed: %w", err)
	}
	if err = uc.svc.CheckPassword(userData.Password, user.Password); err != nil {
		return nil, err
	}
	return userData, nil
}

// Todo
func (uc *userUseCase) UpdateUserProfile(ctx context.Context, user *model.UserProfileRequest, avatar []byte) (*model.UpdateUserProfileResponse, error) {
	userProfile, err := uc.svc.UploadProfile(ctx, user, avatar)
	if err != nil {
		return nil, fmt.Errorf("update user profile failed: %w", err)
	}
	return userProfile, nil
}

func (uc *userUseCase) GetUserProfile(ctx context.Context, uid int64) (*model.GetUserProfileResponse, error) {
	u, err := uc.svc.GetUserProfileInfoById(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("usecase get user profile info failed: %w", err)
	}
	userProfile := &model.GetUserProfileResponse{
		Uid:             u.Uid,
		Username:        u.Username,
		Email:           u.Email,
		Phone:           u.Phone,
		Avatar:          u.Avatar,
		Bio:             u.Bio,
		MembershipLevel: u.MembershipLevel,
		Point:           u.Point,
		Team:            u.Team,
	}
	return userProfile, nil
}
