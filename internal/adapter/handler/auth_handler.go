package handler

import (
	"blog/internal/adapter/handler/request"
	"blog/internal/adapter/handler/response"
	"blog/internal/core/domain/entity"
	"blog/internal/core/service"
	validatorLib "blog/lib/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var err error
var code string
var errorResp response.ErrorResponseDefault
var validate = validator.New()

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

func (a *authHandler) Login(c *fiber.Ctx) error {
	req := request.LoginRequest{}
	resp := response.SuccessAuthResponse{}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] Login - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] Login - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqLogin := entity.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := a.authService.GetUserByEmail(c.Context(), reqLogin)
	if err != nil {
		code = "[HANDLER] Login - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		if err.Error() == "Invalid password" {
			return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
		}

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Login successful"
	resp.AccessToken = result.AccessToken
	resp.ExpiresAt = result.ExpiresAt

	return c.JSON(resp)
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}
