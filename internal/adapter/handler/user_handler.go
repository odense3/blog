package handler

import (
	"blog/internal/adapter/handler/request"
	"blog/internal/adapter/handler/response"
	"blog/internal/core/domain/entity"
	"blog/internal/core/service"
	validatorLib "blog/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserHandler interface {
	UpdatePassword(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

// GetUserById implements UserHandler.
func (u *userHandler) GetUserById(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := int64(claims.UserID)

	if userID == 0 {
		code = "[HANDLER] GetUserById - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	user, err := u.userService.GetUserById(c.Context(), userID)
	if err != nil {
		code = "[HANDLER] GetUserById - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	resp := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Success"
	successResponseDefault.Data = resp

	return c.JSON(successResponseDefault)
}

// UpdatePassword implements UserHandler.
func (u *userHandler) UpdatePassword(c *fiber.Ctx) error {
	var req request.UpdatePasswordRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := int64(claims.UserID)

	if userID == 0 {
		code = "[HANDLER] GetUserById - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[HANDLER] UpdatePassword - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request payload"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = u.userService.UpdatePassword(c.Context(), req.NewPassword, userID)
	if err != nil {
		code = "[HANDLER] UpdatePassword - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Password updated successfully"
	successResponseDefault.Data = nil

	return c.JSON(successResponseDefault)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}
