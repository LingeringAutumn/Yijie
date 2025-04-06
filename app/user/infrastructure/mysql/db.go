package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/LingeringAutumn/Yijie/app/user/domain/model"
	"github.com/LingeringAutumn/Yijie/app/user/domain/repository"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
)

// userDB impl domain.UserDB defined domain
type userDB struct {
	client *gorm.DB
}

func NewUserDB(client *gorm.DB) repository.UserDB {
	return &userDB{client: client}
}

func (db *userDB) CreateUser(ctx context.Context, u *model.User) (int64, error) {
	// 将 entity 转换成 mysql 这边的 model
	user := User{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Phone:    u.Phone,
	}
	// TODO 我不确定我们是否要主动生成雪花ID
	if err := db.client.Create(&user).Error; err != nil {
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create user: %v", err)
	}
	return user.Uid, nil
}

// IsUserExist 方法用于检查指定用户名的用户是否存在于数据库中。
// 它接收一个上下文对象 ctx 用于控制数据库操作的生命周期，以及一个字符串类型的用户名 username 作为查询条件。
// 返回值为一个布尔值，表示用户是否存在，以及一个错误对象，如果发生错误则返回具体的错误信息，否则返回 nil。
func (db *userDB) IsUserExist(ctx context.Context, username string) (bool, error) {
	// 声明一个 User 类型的变量 user，用于存储从数据库中查询到的用户记录。
	// User 结构体通常定义了与数据库中用户表对应的字段。
	var user User
	// 明确指定要操作的数据库表名，constants.UserTableName 是一个常量，代表用户表的名称。
	// 使用 WithContext 方法为数据库操作设置上下文，确保操作可以被正确取消或超时控制。
	// 使用 Where 方法构建 SQL 查询的 WHERE 子句，过滤出 username 字段等于传入的 username 变量值的记录。
	// 使用 First 方法从查询结果中获取第一条记录，并将其存储到 user 变量的内存地址中。
	// 如果查询过程中出现错误，错误信息将被存储在 err 变量中。
	err := db.client.WithContext(ctx).Table(constants.UserTableName).Where("username = ?", username).First(&user).Error
	// 检查查询过程中是否出现错误。
	if err != nil {
		// 判断错误类型是否为 gorm.ErrRecordNotFound，该错误表示在数据库中没有找到符合条件的记录。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果是记录未找到的错误，说明数据库中不存在该用户名对应的用户，返回 false 和 nil 错误。
			return false, nil
		}
		// 如果不是记录未找到的错误，说明数据库操作过程中出现了其他问题，例如数据库连接失败、SQL 语法错误等。
		// 使用 errno.Errorf 方法创建一个自定义的错误对象，包含错误码 errno.InternalDatabaseErrorCode 和错误信息。
		// 错误信息中包含了原始的错误信息，方便后续排查问题。
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query user: %v", err)
	}
	// 如果查询过程中没有出现错误，说明成功从数据库中找到了符合条件的记录，即用户存在。
	// 返回 true 和 nil 错误。
	return true, nil
}

func (db *userDB) GetUserById(ctx context.Context, uid int64) (*model.User, error) {
	var user User
	err := db.client.WithContext(ctx).Table(constants.UserTableName).Where("id = ?", uid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceUserNotExist, "mysql: user %d not exist", uid)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query user: %v", err)
	}
	resp := &model.User{
		Uid:      uid,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Phone:    user.Phone,
	}
	return resp, nil
}

func (db *userDB) GetUserProfileInfoById(ctx context.Context, uid int64) (*model.UserProfileResponse, error) {
	var userProfileResp UserProfileResponse
	err := db.client.WithContext(ctx).Table(constants.UserTableName).Where("id = ?", uid).First(&userProfileResp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceUserNotExist, "mysql: user %d not exist", uid)
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query user: %v", err)
	}
	resp := &model.UserProfileResponse{
		Uid:             userProfileResp.Uid,
		Username:        userProfileResp.Username,
		Email:           userProfileResp.Email,
		Phone:           userProfileResp.Phone,
		Avatar:          userProfileResp.Avatar,
		Bio:             userProfileResp.Bio,
		MembershipLevel: userProfileResp.MembershipLevel,
		Point:           userProfileResp.Point,
		Team:            userProfileResp.Team,
	}
	return resp, nil
}

// StoreUserAvatar 方法用于存储用户头像图片信息到数据库中。
// 接收一个上下文对象 ctx，用于控制数据库操作的生命周期，比如在需要时可以取消操作。
// 接收一个指向 model.Image 结构体的指针 image，包含了要存储或更新的图片相关信息。
// 该方法返回一个错误对象，如果操作成功则返回 nil，若失败则返回具体的错误信息。
func (db *userDB) StoreUserAvatar(ctx context.Context, image *model.Image) error {
	// `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '图片ID，主键，自增',
	// `uid` BIGINT NOT NULL COMMENT '用户ID，关联 users 表，标识图片的上传用户',
	// 声明一个 model.Image 类型的变量 previousImage，用于存储从数据库中查询到的之前的图片记录。
	// 如果查询到了记录，它将包含之前已存在的与当前用户（由 image.Uid 标识）相关的图片信息。
	var previousImage model.Image
	// 使用数据库客户端 db.client，结合传入的上下文 ctx，指定要操作的表为 constants.ImageTableName 所代表的表。
	// 通过 WHERE 条件筛选出 `uid` 字段值等于 image.Uid 的记录，并将查询到的第一条记录填充到 previousImage 中。
	// 将查询过程中可能出现的错误存储在 err 变量中。
	err := db.client.WithContext(ctx).Table(constants.ImageTableName).Where("uid =?", image.Uid).First(&previousImage).Error
	if err != nil {
		// 检查错误是否为 gorm.ErrRecordNotFound，这表示在数据库中没有找到符合条件的记录，即图片不存在。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果图片不存在，则使用数据库客户端 db.client，结合上下文 ctx，指定表为 constants.ImageTableName。
			// 执行创建操作，将传入的 image 结构体作为新的记录插入到数据库中。
			// 将创建过程中可能出现的错误存储在 err 变量中。
			err = db.client.WithContext(ctx).Table(constants.ImageTableName).Create(&image).Error
			if err != nil {
				// 如果创建失败，使用 errno.Errorf 方法生成一个自定义错误。
				// 错误码为 errno.InternalDatabaseErrorCode，错误信息包含了具体的错误描述和原始错误。
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to store avatar: %v", err)
			}
		}
		// 如果错误不是 gorm.ErrRecordNotFound，说明在查询过程中出现了其他错误，
		// 这里没有进一步处理其他错误类型，只是让错误继续向上层函数传递。
	} else {
		// 如果没有错误，说明查询到了之前的图片记录，将之前记录的 ImageID 赋值给当前的 image.ImageID。
		// 这样做是为了在更新操作时能够正确地定位到要更新的记录。
		image.ImageID = previousImage.ImageID
		// 使用数据库客户端 db.client，结合上下文 ctx，指定表为 constants.ImageTableName。
		// 使用 Model 方法指定操作的模型为 model.Image{}，通过 WHERE 条件筛选出 `image_id` 字段值等于 image.ImageID 的记录。
		// 执行更新操作，将 image 结构体中的字段值更新到数据库中对应的记录上。
		// 将更新过程中可能出现的错误存储在 err 变量中。
		err = db.client.WithContext(ctx).Table(constants.ImageTableName).
			Model(&model.Image{}).Where("image_id =?", image.ImageID).Updates(image).Error
		if err != nil {
			// 如果更新失败，使用 errno.Errorf 方法生成一个自定义错误。
			// 错误码为 errno.InternalDatabaseErrorCode，错误信息包含了具体的错误描述和原始错误。
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update avatar: %v", err)
		}
	}
	// 如果整个操作过程没有出现错误，返回 nil 表示操作成功。
	return nil
}

// StoreUserProfile TODO id是user的id,uid是image的userid
func (db *userDB) StoreUserProfile(ctx context.Context, userProfileRequest *model.UserProfileRequest, uid int64, image *model.Image) (*model.UserProfileResponse, error) {
	var userProfileResponse UserProfileResponse
	err := db.client.WithContext(ctx).Table(constants.UserTableName).Where("id = ?", uid).First(&userProfileResponse).Error
	if err != nil {
		// 如果个人信息不存在，则创建
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r, err := db.ConvertDefaultUPReqToUPResp(userProfileRequest, image.Url)
			err = db.client.WithContext(ctx).Table(constants.UserTableName).Create(r).Error
			if err != nil {
				return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create user profile: %v", err)
			}
			userProfileResponse.Uid = r.Uid
			userProfileResponse.Username = r.Username
			userProfileResponse.Email = r.Email
			userProfileResponse.Phone = r.Phone
			userProfileResponse.Avatar = r.Avatar
			userProfileResponse.Bio = r.Bio
			userProfileResponse.MembershipLevel = r.MembershipLevel
			userProfileResponse.Point = r.Point
			userProfileResponse.Team = r.Team
		}
	} else {
		// 如果存在，就更新
		err = db.client.WithContext(ctx).Table(constants.UserTableName).Where("id = ?", uid).Updates(userProfileResponse).Error
		if err != nil {
			return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update user profile: %v", err)
		}
	}
	resp := &model.UserProfileResponse{
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

func (db *userDB) ConvertDefaultUPReqToUPResp(req *model.UserProfileRequest, avatarUrl string) (*UserProfileResponse, error) {
	if req == nil {
		return nil, errno.Errorf(errno.ParamVerifyErrorCode, "ConvertDefaultUPReqToUPResp: failed to convert user profile request to user profile response")
	}
	return &UserProfileResponse{
		Uid:             req.Uid,
		Username:        req.Username,
		Email:           req.Email,
		Phone:           req.Phone,
		Avatar:          avatarUrl,
		Bio:             req.Bio,
		MembershipLevel: constants.MembershipLevelFreeCode,
		Point:           0,
		Team:            "",
	}, nil
}
