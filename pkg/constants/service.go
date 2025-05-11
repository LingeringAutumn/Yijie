package constants

const (
	GatewayServiceName      = "gateway"
	UserServiceName         = "user"
	VideoServiceName        = "video"
	CommentServiceName      = "comment"
	UserBehaviorServiceName = "user-behavior"
	ChatServiceName         = "chat"
)

// UserService
const (
	UserMaximumPasswordLength      = 72 // DO NOT EDIT (ref: bcrypt.GenerateFromPassword)
	UserMinimumPasswordLength      = 5
	UserDefaultEncryptPasswordCost = 10
	UserTestId                     = 1
)
