package mysql

import (
	"database/sql"
	"fmt"
	"gameapp/entity"
)

func (d MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	row := d.db.QueryRow(`select id from users where phone_number=?`, phoneNumber)

	var id uint
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("can't scan query result: %w", err)
	}

	return false, nil
}

func (d MySQLDB) Register(u entity.User) (entity.User, error) {
	result, err := d.db.Exec(`INSERT INTO users(name, phone_number) VALUES (?, ?)`, u.Name, u.PhoneNumber)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't insert user to db: %w", err)
	}

	id, _ := result.LastInsertId()
	u.ID = uint(id)

	return u, nil
}
