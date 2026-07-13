package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"todo-api/repositories"
	"todo-api/services"
	"todo-api/utils"
)

var userService = services.NewUserService(repositories.NewUserRepository())

type AuthController struct {
	BaseController
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *AuthController) Register() {
	var req RegisterRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidRequestBody, nil)
		return
	}

	user, err := userService.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNameRequired),
			errors.Is(err, services.ErrEmailRequired),
			errors.Is(err, services.ErrPasswordRequired),
			errors.Is(err, services.ErrPasswordTooShort),
			errors.Is(err, services.ErrInvalidEmail):
			utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, err.Error(), nil)
		case errors.Is(err, services.ErrEmailExists):
			utils.SendJSONResponse(c.Ctx, http.StatusConflict, false, err.Error(), nil)
		default:
			utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, utils.MsgFailedToRegisterUser, nil)
		}
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusCreated, true, utils.MsgUserRegistered, user)
}

func (c *AuthController) Login() {
	var req LoginRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidRequestBody, nil)
		return
	}

	token, err := userService.LoginUser(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailRequired),
			errors.Is(err, services.ErrPasswordRequired):
			utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, err.Error(), nil)
		case errors.Is(err, services.ErrInvalidCredentials):
			utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, err.Error(), nil)
		default:
			utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, utils.MsgFailedToLogin, nil)
		}
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, utils.MsgLoginSuccessful, LoginResponse{Token: token})
}
