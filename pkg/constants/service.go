package constants

const (
	GatewayServiceName = "gateway"
	UserServiceName    = "user"
)

// UserService
const (
	UserMaximumPasswordLength      = 72 // DO NOT EDIT (ref: bcrypt.GenerateFromPassword)
	UserMinimumPasswordLength      = 5
	UserDefaultEncryptPasswordCost = 10
	UserTestId                     = 1
)
