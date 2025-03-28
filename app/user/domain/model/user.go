package model

// User 用于在 handler --- use case --- infrastructure 之间传递数据的实体类
// 目的是方便 use case 操作对应的业务
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

type UserProfile struct {
}
