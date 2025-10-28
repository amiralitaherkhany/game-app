package httpserver

import (
	"gameapp/pkg/httpmsg"
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
		message, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, message)
	}

	return c.JSON(http.StatusCreated, resp)
}

func (s Server) userLoginHandler(c echo.Context) error {
	req := new(userservice.LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := s.userSvc.Login(*req)
	if err != nil {
		message, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, message)
	}

	return c.JSON(http.StatusOK, resp)
}

func (s Server) userProfileHandler(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")

	claims, err := s.authSvc.ParseAccessToken(authHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	resp, err := s.userSvc.GetProfile(
		userservice.GetProfileRequest{
			UserID: claims.UserID,
		},
	)
	if err != nil {
		message, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, message)
	}

	return c.JSON(http.StatusOK, resp)
}
