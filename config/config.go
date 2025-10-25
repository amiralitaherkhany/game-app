package config

import (
	"gameapp/repository/mysql"
	"gameapp/service/authservice"
)

type Config struct {
	HTTPServer HTTPServer
	Auth       authservice.Config
	DB         mysql.Config
}

type HTTPServer struct {
	Port uint
}
