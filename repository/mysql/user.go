package mysql

import (
	"database/sql"
	"gameapp/entity"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
)

func (d MySQLDB) GetUserByID(userID uint) (entity.User, error) {
	row := d.db.QueryRow(`select id, name, phone_number, password from users where id=?`, userID)

	var user entity.User
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{},
				richerror.
					New("mysql.GetUserByID").
					WithErr(err).
					WithMessage(errmsg.ErrorMsgNotFound).
					WithKind(richerror.KindNotFound)
		}
		return entity.User{},
			richerror.
				New("mysql.GetUserByID").
				WithErr(err).
				WithMessage(errmsg.ErrorMsgCantScanQueryResult).
				WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d MySQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	row := d.db.QueryRow(`select id, name, phone_number, password from users where phone_number=?`, phoneNumber)

	var user entity.User
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{},
				richerror.
					New("mysql.GetUserByPhoneNumber").
					WithErr(err).
					WithMessage(errmsg.ErrorMsgNotFound).
					WithKind(richerror.KindNotFound)
		}
		return entity.User{},
			richerror.
				New("mysql.GetUserByPhoneNumber").
				WithErr(err).
				WithMessage(errmsg.ErrorMsgCantScanQueryResult).
				WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	row := d.db.QueryRow(`select id from users where phone_number=?`, phoneNumber)

	var id uint
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, richerror.
			New("mysql.GetUserByPhoneNumber").
			WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d MySQLDB) Register(u entity.User) (entity.User, error) {
	result, err := d.db.Exec(`INSERT INTO users(name, phone_number, password) VALUES (?, ?, ?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{},
			richerror.
				New("mysql.Register").
				WithErr(err).
				WithMessage("can't insert user to db").
				WithKind(richerror.KindUnexpected)
	}

	id, _ := result.LastInsertId()
	u.ID = uint(id)

	return u, nil
}
