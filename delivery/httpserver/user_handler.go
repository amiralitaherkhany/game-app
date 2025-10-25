package httpserver

import (
	"gameapp/service/userservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s Server) userRegisterHandler(c echo.Context) error {
	req := new(userservice.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := s.userSvc.Register(*req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}
