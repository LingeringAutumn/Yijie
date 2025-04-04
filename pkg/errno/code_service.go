package errno

// 业务强相关, 范围是 1000-9999
// User
const (
	ServiceWrongPassword = 1000 + iota
	ServiceUserExist
	ServiceUserNotExist

	ErrRecordNotFound
)
