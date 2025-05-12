package service

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/LingeringAutumn/Yijie/app/user/domain/model"
	"github.com/LingeringAutumn/Yijie/config"
	userData "github.com/LingeringAutumn/Yijie/pkg/base/context"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

// EncryptPassword 是 UserService 结构体的一个方法，用于对输入的明文密码进行加密。
// 它使用 bcrypt 算法生成加密后的密码，并在加密过程中处理可能出现的错误。
// 参数 pwd 是需要加密的明文密码。
// 返回值是加密后的密码字符串和一个错误对象，如果加密成功，错误对象为 nil。
func (svc *UserService) EncryptPassword(pwd string) (string, error) {
	// 调用 bcrypt 包的 GenerateFromPassword 函数对明文密码进行加密。
	// 第一个参数是将输入的字符串密码转换为字节切片，因为该函数需要字节切片类型的输入。
	// 第二个参数是一个常量，表示 bcrypt 算法的计算成本，成本越高加密越安全但耗时也越长。
	// passwordDigest 是生成的加密后的密码，类型为字节切片。
	// err 是在加密过程中可能出现的错误。
	passwordDigest, err := bcrypt.GenerateFromPassword([]byte(pwd), constants.UserDefaultEncryptPasswordCost)
	// 检查加密过程中是否发生错误。
	if err != nil {
		// 如果发生错误，使用 errno 包的 NewErrNo 函数创建一个自定义的错误对象。
		// errno.InternalServiceErrorCode 是一个错误码，表示内部服务错误。
		// fmt.Sprintf 用于生成包含错误信息的字符串，其中包含原始密码和具体的错误信息。
		// 返回一个空字符串和自定义的错误对象。
		return "", errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("encrypt password failed, pwd: %s, err: %v", pwd, err))
	}

	// 如果加密过程没有发生错误，将加密后的字节切片密码转换为字符串。
	// 返回加密后的密码字符串和 nil 表示没有错误。
	return string(passwordDigest), nil
}

func (svc *UserService) CreateUser(ctx context.Context, u *model.User) (int64, error) {
	uid, err := svc.db.CreateUser(ctx, u)
	if err != nil {
		return 0, fmt.Errorf("create user failed: %w", err)
	}
	return uid, nil
}

// GetUserById 这个是核对密码的时候获取密码的
func (svc *UserService) GetUserById(ctx context.Context, uid int64) (*model.User, error) {
	u, err := svc.db.GetUserById(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}
	return u, nil
}

// GetUserById 这个是核对密码的时候获取密码的
func (svc *UserService) GetUserByName(ctx context.Context, username string) (*model.User, error) {
	u, err := svc.db.GetUserByName(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}
	return u, nil
}

// GetUserProfileInfoById 这个是查询和更新个人信息的时候来获取个人信息的
func (svc *UserService) GetUserProfileInfoById(ctx context.Context, uid int64) (*model.UserProfileResponse, error) {
	userInfo, err := svc.db.GetUserProfileInfoById(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("service get user profile info failed: %w", err)
	}
	return userInfo, nil
}

func (svc *UserService) IsUserExist(ctx context.Context, username string) (bool, error) {
	exist, err := svc.db.IsUserExist(ctx, username)
	if err != nil {
		return false, fmt.Errorf("check user exist failed: %w", err)
	}
	return exist, nil
}

// CheckPassword 是 UserService 结构体的一个方法，用于验证输入的明文密码是否与加密后的密码相匹配。
// 参数 passwordDigest 是存储在数据库中的加密后的密码。
// 参数 password 是用户输入的明文密码。
// 返回值是一个错误对象，如果密码匹配，返回 nil；如果密码不匹配，返回一个自定义的错误对象。
func (svc *UserService) CheckPassword(passwordDigest, password string) error {
	// 调用 bcrypt 包的 CompareHashAndPassword 函数来比较加密后的密码和明文密码。
	// 该函数接收两个字节切片作为参数，所以需要将输入的字符串转换为字节切片。
	// 如果比较结果不为 nil，说明明文密码与加密后的密码不匹配。
	if bcrypt.CompareHashAndPassword([]byte(passwordDigest), []byte(password)) != nil {
		// 当密码不匹配时，使用 errno 包的 NewErrNo 函数创建一个自定义的错误对象。
		// errno.ServiceWrongPassword 是一个错误码，表示密码错误。
		// 错误信息 "wrong password" 明确告知调用者密码输入错误。
		return errno.NewErrNo(errno.ServiceWrongPassword, "wrong password")
	}

	// 如果密码匹配，返回 nil 表示没有错误。
	return nil
}

func (svc *UserService) UploadProfile(ctx context.Context, user *model.UserProfileRequest, avatar []byte) (*model.UpdateUserProfileResponse, error) {
	// 1. 把头像的二进制字节流传到MinIO服务器上
	// log.Printf("upload profile input: %+v", user)
	// log.Printf(">>>> UploadProfile input: %+v", user)

	var imageId string
	// TODO 这个地方没问题吗？
	uid, err := svc.GetUserId(ctx)
	log.Printf(">>>> Got UID: %v, err: %v", uid, err)
	imageId = fmt.Sprintf("%d", uid)
	err = utils.MinioClientGlobal.UploadFile(constants.ImageBucket, imageId, constants.Location, constants.ImageType, avatar)
	if err != nil {
		return nil, fmt.Errorf("upload image failed: %w", err)
	}

	// 2. 获取头像对应的url
	url := fmt.Sprintf("%s/%s/%d", config.Minio.Endpoint, constants.ImageBucket, uid)

	// 3. 储存头像到image内部
	image := &model.Image{
		Url: url,
		Uid: uid,
	}
	err = svc.db.StoreUserAvatar(ctx, image)
	if err != nil {
		return nil, fmt.Errorf("store user avatar failed: %w", err)
	}
	// 4. 储存其它内容
	log.Printf(">>>> UploadProfile received: UID=%d, UserProfile=%+v", uid, user)
	userProfileResponse, err := svc.db.StoreUserProfile(ctx, user, uid, image)
	if err != nil {
		log.Printf(">>>> StoreUserProfile error: %v", err)
		return nil, fmt.Errorf("store user profile failed: %w", err)
	}
	log.Printf(">>>> StoreUserProfile returned: %+v", userProfileResponse)
	// 5. 更改返回类型
	resp := &model.UpdateUserProfileResponse{
		Uid:             userProfileResponse.Uid,
		Username:        userProfileResponse.Username,
		Email:           userProfileResponse.Email,
		Phone:           userProfileResponse.Phone,
		Avatar:          userProfileResponse.Avatar,
		Bio:             userProfileResponse.Bio,
		MembershipLevel: userProfileResponse.MembershipLevel,
		Point:           userProfileResponse.Point,
		Team:            userProfileResponse.Team,
	}
	return resp, nil
}

func (svc *UserService) GetUserId(ctx context.Context) (uid int64, err error) {
	if uid, err = userData.GetLoginData(ctx); err != nil {
		return 0, fmt.Errorf("get user id failed: %w", err)
	}
	return uid, err
}
