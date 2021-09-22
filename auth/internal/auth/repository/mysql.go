package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/ulwisibaq/efishery/auth/internal/auth"
	"github.com/ulwisibaq/efishery/auth/internal/models"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(
	db *sqlx.DB,
) auth.MysqlRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar AuthRepository) CreateUser(userReq models.UserRequest) (id int64, err error) {

	res, err := ar.db.Exec(
		CreateUserQuery,
		userReq.Name,
		userReq.Phone,
		userReq.Password,
		userReq.RoleId,
		userReq.Created,
	)
	if err != nil {
		return
	}

	id, _ = res.LastInsertId()

	return

}

func (ar AuthRepository) GetRoleIdByName(roleName string) (roleId int, err error) {

	err = ar.db.QueryRowx(
		GetRoleIdByNameQuery,
		roleName,
	).Scan(&roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			return roleId, nil
		}
		return
	}

	return
}

func (ar AuthRepository) GetUserByPhone(phone string) (resp models.User, err error) {

	err = ar.db.QueryRowx(GetUserByPhoneQuery, phone).StructScan(&resp)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}
		return
	}
	return
}
