package main

import (
	"fmt"
	"gameapp/repository/mysql"
)

func main() {
	mysqlRepo := mysql.New()
	is, err := mysqlRepo.IsPhoneNumberUnique("0914")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(is)
}
