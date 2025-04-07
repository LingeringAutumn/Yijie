package service

import (
	"github.com/LingeringAutumn/Yijie/app/user/domain/repository"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

type UserService struct {
	db    repository.UserDB
	redis repository.UserRedis
	rpc   repository.UserRPC
	sf    *utils.Snowflake
}

// TODO User这边先不写RPC
/*func NewUserService(db repository.UserDB, redis repository.UserRedis, rpc repository.UserRPC, sf *utils.Snowflake) *UserService {
	if db == nil {
		panic("userService`s db should not be nil")
	}
	if redis == nil {
		panic("userService`s redis should not be nil")
	}
	if rpc == nil {
		panic("userService`s rpc should not be nil")
	}
	if sf == nil {
		panic("userService`s sf should not be nil")
	}
	svc := &UserService{
		db:    db,
		redis: redis,
		rpc:   rpc,
		sf:    sf,
	}
	return svc
}*/

func NewUserService(db repository.UserDB, redis repository.UserRedis, sf *utils.Snowflake) *UserService {
	if db == nil {
		panic("userService`s db should not be nil")
	}
	if redis == nil {
		panic("userService`s redis should not be nil")
	}
	if sf == nil {
		panic("userService`s sf should not be nil")
	}
	svc := &UserService{
		db:    db,
		redis: redis,
		sf:    sf,
	}
	return svc
}
