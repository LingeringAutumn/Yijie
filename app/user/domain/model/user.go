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
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	// TODO 有点怀疑这里
	Avatar          int64  `json:"avatar"`
	Bio             string `json:"bio"`
	MembershipLevel int64  `json:"member"`
	Point           int64  `json:"point"`
	Team            string `json:"team"`
}
