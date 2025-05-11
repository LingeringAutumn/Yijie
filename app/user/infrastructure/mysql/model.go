package mysql

import (
	"github.com/LingeringAutumn/Yijie/pkg/constants"
)

// User 是 mysql 【独有】的，和 db 中的表数据一一对应，和 entities 层的 User 的作用域不一样

// User 表示 users 表中的用户基本信息。
type User struct {
	Uid      int64  `gorm:"column:id;primaryKey" json:"uid"`      // 用户ID，对应数据库中的 id 字段
	Username string `gorm:"column:username" json:"username"`      // 用户名
	Password string `gorm:"column:password_hash" json:"password"` // 密码哈希值
	Email    string `gorm:"column:email" json:"email"`            // 邮箱
	Phone    string `gorm:"column:phone" json:"phone"`            // 手机号
}

// UserInfo 是用户的基础信息响应
type UserInfo struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
}

// UserProfileRequest 是用户更新资料的请求体
type UserProfileRequest struct {
	Uid      int64  `json:"uid"`      // 用户 ID
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 邮箱
	Phone    string `json:"phone"`    // 手机号
	Bio      string `json:"bio"`      // 个人简介
}

// UserProfileResponse 是用户资料的完整响应体
type UserProfileResponse struct {
	Uid             int64  `json:"uid" gorm:"-"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Avatar          string `json:"avatar" gorm:"column:avatar_url"` // 头像 URL
	Bio             string `json:"bio"`                             // 个人简介
	MembershipLevel int64  `json:"member"`                          // 会员等级
	Point           int64  `json:"point"`                           // 当前积分
	Team            string `json:"team"`                            // 团队信息
}

// UpdateUserProfileResponse 是更新用户资料后的响应
type UpdateUserProfileResponse = UserProfileResponse

// GetUserProfileResponse 是获取用户资料的响应
type GetUserProfileResponse = UserProfileResponse

// Image 表示用户上传的图片信息
type Image struct {
	ImageID int64  `json:"image_id" gorm:"column:image_id;primaryKey;autoIncrement"`
	Uid     int64  `json:"uid" gorm:"column:uid;not null;index"`
	Url     string `json:"url" gorm:"column:url;type:varchar(255);not null"`
}

func (User) TableName() string {
	return constants.UserTableName
}

func (Image) TableName() string {
	return constants.ImageTableName
}
