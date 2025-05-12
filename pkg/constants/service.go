package constants

const (
	GatewayServiceName       = "gateway"
	UserServiceName          = "user"
	VideoServiceName         = "video"
	CommentServiceName       = "comment"
	UserBehaviourServiceName = "user_behaviour"
	ChatServiceName          = "chat"
)

// UserService
const (
	UserMaximumPasswordLength      = 72 // DO NOT EDIT (ref: bcrypt.GenerateFromPassword)
	UserMinimumPasswordLength      = 5
	UserDefaultEncryptPasswordCost = 10
	UserTestId                     = 1
)
