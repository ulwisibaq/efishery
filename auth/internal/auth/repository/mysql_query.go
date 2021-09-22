package repository

const (
	GetRoleIdByNameQuery = `
		SELECT
			id
		FROM 
			role
		WHERE 
			name = ?
	`

	CreateUserQuery = `
		INSERT INTO user (name, phone, password, role_id, created_at)
		VALUES(?, ?, ?, ?, ?);
	`

	GetUserByPhoneQuery = `
		SELECT 
			u.id as id,
			u.name as name,
			u.phone as phone, 
			r.name as role,
			u.password as password,
			u.created_at as created_at 
		FROM 
			user u join role r on r.id = u.role_id
		WHERE u.phone = ?
	`
)
