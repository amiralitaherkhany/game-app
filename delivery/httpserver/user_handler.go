package httpserver

import (
	"gameapp/dto"
	"gameapp/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s Server) userRegisterHandler(c echo.Context) error {
	req := new(dto.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err, fieldErrors := s.userValidator.ValidateRegisterRequest(*req)
	if err != nil {
		message, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, echo.Map{
			"message": message,
			"errors":  fieldErrors,
		})
	}

	resp, err := s.userSvc.Register(*req)
	if err != nil {
		message, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, message)
	}

	return c.JSON(http.StatusCreated, resp)
}

func (s Server) userLoginHandler(c echo.Context) error {
	req := new(dto.LoginRequest)
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
		dto.GetProfileRequest{
			UserID: claims.UserID,
		},
	)
	if err != nil {
		message, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, message)
	}

	return c.JSON(http.StatusOK, resp)
}
