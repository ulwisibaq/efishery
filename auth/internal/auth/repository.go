package auth

import "github.com/ulwisibaq/efishery/auth/internal/models"

type MysqlRepository interface {
	CreateUser(userReq models.UserRequest) (id int64, err error)
	GetRoleIdByName(roleName string) (roleId int, err error)
	GetUserByPhone(phone string) (resp models.User, err error)
}
