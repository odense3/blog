package handler

import (
	"blog/internal/adapter/handler/request"
	"blog/internal/adapter/handler/response"
	"blog/internal/core/domain/entity"
	"blog/internal/core/service"
	"blog/lib/conv"
	validatorLib "blog/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var successResponseDefault response.SuccessResponseDefault

type CategoryHandler interface {
	GetCategories(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	EditCategoryByID(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error

	GetCategoryFE(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

// GetCategoryFE implements CategoryHandler.
func (ch *categoryHandler) GetCategoryFE(c *fiber.Ctx) error {
	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID:            result.ID,
			Title:         result.Title,
			Slug:          result.Slug,
			CreatedByName: result.User.Name,
		}

		categoryResponses = append(categoryResponses, categoryResponse)
	}

	successResponseDefault.Pagination = nil
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Gategories fetched successfully"
	successResponseDefault.Data = categoryResponses

	return c.JSON(successResponseDefault)
}

// CreateCategory implements CategoryHandler.
func (ch *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] CreateCategory - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateCategory - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] CreateCategory - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = ch.categoryService.CreateCategory(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] Create Category - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	successResponseDefault.Data = nil
	successResponseDefault.Pagination = nil
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Category created successfully"

	return c.JSON(successResponseDefault)
}

// DeleteCategoryByID implements CategoryHandler.
func (ch *categoryHandler) DeleteCategory(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] GetCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("categoryId")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = ch.categoryService.DeleteCategory(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	successResponseDefault.Data = nil
	successResponseDefault.Pagination = nil
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Category deleted successfully"

	return c.JSON(successResponseDefault)
}

// EditCategoryByID implements CategoryHandler.
func (ch *categoryHandler) EditCategoryByID(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] EditCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[HANDLER] EditCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request payload"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] EditCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	idParam := c.Params("categoryID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] EditCategoryByID - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		ID:    id,
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = ch.categoryService.EditCategoryByID(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] EditCategoryByID - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	successResponseDefault.Data = nil
	successResponseDefault.Pagination = nil
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Category updated successfully"

	return c.JSON(successResponseDefault)
}

// GetCategories implements CategoryHandler.
func (ch *categoryHandler) GetCategories(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetCategories - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID:            result.ID,
			Title:         result.Title,
			Slug:          result.Slug,
			CreatedByName: result.User.Name,
		}

		categoryResponses = append(categoryResponses, categoryResponse)
	}

	successResponseDefault.Pagination = nil
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Gategories fetched successfully"
	successResponseDefault.Data = categoryResponses

	return c.JSON(successResponseDefault)
}

// GetCategoryByID implements CategoryHandler.
func (ch *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] GetCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("categoryId")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.categoryService.GetCategoryByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponse := response.SuccessCategoryResponse{
		ID:            result.ID,
		Title:         result.Title,
		Slug:          result.Slug,
		CreatedByName: result.User.Name,
	}

	successResponseDefault.Pagination = nil
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Category fetched details successfully"
	successResponseDefault.Data = categoryResponse

	return c.JSON(successResponseDefault)
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
	}
}
