package mysql

import "github.com/LingeringAutumn/Yijie/pkg/constants"

// User 是 mysql 【独有】的，和 db 中的表数据一一对应，和 entities 层的 User 的作用域不一样

type User struct {
	Uid      int64
	Username string
	Password string
	Email    string
	Phone    string
}

type UserInfo struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
}

type UserProfileRequest struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	// TODO传进来的avatar头像文件本身
	Avatar []byte `json:"avatar"`
	Bio    string `json:"bio"`
}

type UserProfileResponse struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	// TODO返回的avatar是它的url
	Avatar          string `json:"avatar"`
	Bio             string `json:"bio"`
	MembershipLevel int64  `json:"member"`
	Point           int64  `json:"point"`
	Team            string `json:"team"`
}

type UpdateUserProfileResponse struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	// TODO返回的avatar是它的url
	Avatar          string `json:"avatar"`
	Bio             string `json:"bio"`
	MembershipLevel int64  `json:"member"`
	Point           int64  `json:"point"`
	Team            string `json:"team"`
}

type GetUserProfileResponse struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	// TODO返回的avatar是它的url
	Avatar          string `json:"avatar"`
	Bio             string `json:"bio"`
	MembershipLevel int64  `json:"member"`
	Point           int64  `json:"point"`
	Team            string `json:"team"`
}

type Image struct {
	Uid     int64  `json:"uid"`
	ImageID string `json:"image_id"`
	Url     string `json:"url"`
}

func (User) TableName() string {
	return constants.UserTableName
}

func (Image) TableName() string {
	return constants.ImageTableName
}
