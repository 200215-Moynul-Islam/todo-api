package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"todo-api/repositories"
	"todo-api/services"
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

func (c *AuthController) Register() {
	var req RegisterRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.SendError(http.StatusBadRequest, "Invalid request body")
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
			c.SendError(http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrEmailExists):
			c.SendError(http.StatusConflict, err.Error())
		default:
			c.SendError(http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	c.SendSuccess(http.StatusCreated, "User registered successfully", user)
}
