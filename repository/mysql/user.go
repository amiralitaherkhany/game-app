package mysql

import "gameapp/entity"

func (d DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {}

func (d DB) Register(u entity.User) (entity.User, error) {}
