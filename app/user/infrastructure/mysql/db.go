package mysql

import (
	"context"
	"github.com/LingeringAutumn/Yijie/app/user/domain/model"
	"github.com/LingeringAutumn/Yijie/app/user/domain/repository"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"gorm.io/gorm"
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
	if err := db.client.Create(&user).Error; err != nil {
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode,"mysql: failed to create user: %v", err)
	}
	return u.Uid,nil
}

func
