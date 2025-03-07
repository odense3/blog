package handler

import (
	"blog/internal/adapter/handler/request"
	"blog/internal/adapter/handler/response"
	"blog/internal/core/domain/entity"
	"blog/internal/core/service"
	"blog/lib/conv"
	validatorLib "blog/lib/validator"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentByID(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	EditContentByID(c *fiber.Ctx) error
	DeleteContent(c *fiber.Ctx) error
	UploadImageR2(c *fiber.Ctx) error

	// FE
	GetContentWithQuery(c *fiber.Ctx) error
	GetContentDetail(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

// GetContentDetail implements ContentHandler.
func (ch *contentHandler) GetContentDetail(c *fiber.Ctx) error {
	idParam := c.Params("contentId")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetContentDetail - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.contentService.GetContentByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetContentDetail - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	contentResponse := response.SuccessContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		Status:       result.Status,
		CategoryID:   result.CategoryID,
		UserID:       result.UserID,
		CreatedAt:    result.CreatedAt.Format(time.RFC3339),
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
	}

	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Success"
	successResponseDefault.Data = contentResponse
	return c.JSON(successResponseDefault)
}

// GetContentWithQuery implements ContentHandler.
func (ch *contentHandler) GetContentWithQuery(c *fiber.Ctx) error {
	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[HANDLER] GetContentWithQuery - 1"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()
			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	limit := 6
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code = "[HANDLER] GetContentWithQuery - 2"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()
			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "desc"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	categoryID := 0
	if c.Query("categoryID") != "" {
		categoryID, err = conv.StringToInt(c.Query("categoryID"))
		if err != nil {
			code = "[HANDLER] GetContentWithQuery - 3"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()
			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	reqEntity := entity.QueryString{
		Limit:      limit,
		Page:       page,
		OrderBy:    orderBy,
		OrderType:  orderType,
		Search:     search,
		Status:     "PUBLISH",
		CategoryID: int64(categoryID),
	}

	results, totalData, totalPages, err := ch.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] GetContentWithQuery - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	contentResponses := []response.SuccessContentResponse{}
	for _, result := range results {
		contentResponse := response.SuccessContentResponse{
			ID:           result.ID,
			Title:        result.Title,
			Excerpt:      result.Excerpt,
			Description:  result.Description,
			Image:        result.Image,
			Tags:         result.Tags,
			Status:       result.Status,
			CategoryID:   result.CategoryID,
			UserID:       result.UserID,
			CreatedAt:    result.CreatedAt.Format(time.RFC3339),
			CategoryName: result.Category.Title,
			Author:       result.User.Name,
		}

		contentResponses = append(contentResponses, contentResponse)
	}

	successResponseDefault.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Contents fetched successfully"
	successResponseDefault.Data = contentResponses

	return c.JSON(successResponseDefault)
}

// CreateContent implements ContentHandler.
func (ch *contentHandler) CreateContent(c *fiber.Ctx) error {
	var req request.ContentRequest
	claims := c.Locals("user").(*entity.JwtData)
	UserID := claims.UserID

	if UserID == 0 {
		code = "[HANDLER] CreateContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] CreateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		UserID:      int64(UserID),
	}

	err = ch.contentService.CreateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	successResponseDefault.Data = nil
	successResponseDefault.Pagination = nil
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Content created successfully"

	return c.Status(fiber.StatusCreated).JSON(successResponseDefault)
}

// DeleteContent implements ContentHandler.
func (ch *contentHandler) DeleteContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	UserID := claims.UserID

	if UserID == 0 {
		code = "[HANDLER] DeleteContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentId")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] DeleteContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = ch.contentService.DeleteContent(c.Context(), id)
	if err != nil {
		code = "[HANDLER] DeleteContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	successResponseDefault.Data = nil
	successResponseDefault.Pagination = nil
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Content deleted successfully"

	return c.JSON(successResponseDefault)
}

// EditContentByID implements ContentHandler.
func (ch *contentHandler) EditContentByID(c *fiber.Ctx) error {
	var req request.ContentRequest
	claims := c.Locals("user").(*entity.JwtData)
	UserID := claims.UserID

	if UserID == 0 {
		code = "[HANDLER] EditContentByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] EditContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] EditContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	idParam := c.Params("contentId")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] EditContentByID - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		ID:          id,
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		UserID:      int64(UserID),
	}

	err = ch.contentService.EditContentByID(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] EditContentByID - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	successResponseDefault.Data = nil
	successResponseDefault.Pagination = nil
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Content updated successfully"

	return c.JSON(successResponseDefault)
}

// GetContentByID implements ContentHandler.
func (ch *contentHandler) GetContentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	UserID := claims.UserID

	if UserID == 0 {
		code = "[HANDLER] GetContentByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentId")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.contentService.GetContentByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	contentResponse := response.SuccessContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		Status:       result.Status,
		CategoryID:   result.CategoryID,
		UserID:       result.UserID,
		CreatedAt:    result.CreatedAt.Format(time.RFC3339),
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
	}

	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Success"
	successResponseDefault.Data = contentResponse
	return c.JSON(successResponseDefault)
}

// GetContents implements ContentHandler.
func (ch *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	UserID := claims.UserID

	if UserID == 0 {
		code = "[HANDLER] GetContents - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code := "[HANDLER] GetContents - 2"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Invalid page number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	limit := 10
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code := "[HANDLER] GetContents - 3"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Invalid limit number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "desc"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	categoryID := 0
	if c.Query("categoryID") != "" {
		categoryID, err = conv.StringToInt(c.Query("categoryID"))
		if err != nil {
			code := "[HANDLER] GetContents - 4"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Invalid category ID"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	reqEntity := entity.QueryString{
		Limit:      limit,
		Page:       page,
		OrderBy:    orderBy,
		OrderType:  orderType,
		Search:     search,
		CategoryID: int64(categoryID),
	}

	results, totalData, totalPages, err := ch.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] GetContent - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	contentResponses := []response.SuccessContentResponse{}
	for _, result := range results {
		contentResponse := response.SuccessContentResponse{
			ID:           result.ID,
			Title:        result.Title,
			Excerpt:      result.Excerpt,
			Description:  result.Description,
			Image:        result.Image,
			Tags:         result.Tags,
			Status:       result.Status,
			CategoryID:   result.CategoryID,
			UserID:       result.UserID,
			CreatedAt:    result.CreatedAt.Format(time.RFC3339),
			CategoryName: result.Category.Title,
			Author:       result.User.Name,
		}

		contentResponses = append(contentResponses, contentResponse)
	}

	successResponseDefault.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}
	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Contents fetched successfully"
	successResponseDefault.Data = contentResponses

	return c.JSON(successResponseDefault)
}

// UploadImageR2 implements ContentHandler.
func (ch *contentHandler) UploadImageR2(c *fiber.Ctx) error {
	var req request.FileUploadRequest
	claims := c.Locals("user").(*entity.JwtData)
	UserID := claims.UserID

	if UserID == 0 {
		code = "[HANDLER] UploadImageR2 - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	file, err := c.FormFile("image")
	if err != nil {
		code = "[HANDLER] UploadImageR2 - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = c.SaveFile(file, fmt.Sprintf("./temp/content/%s", file.Filename)); err != nil {
		code = "[HANDLER] UploadImageR2 - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	req.Image = fmt.Sprintf("./temp/content/%s", file.Filename)
	reqEntity := entity.FileUploadEntity{
		Name: fmt.Sprintf("%d-%d", int64(UserID), time.Now().UnixNano()),
		Path: req.Image,
	}

	imageUrl, err := ch.contentService.UploadImageR2(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] UploadImageR2 - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if req.Image != "" {
		err = os.Remove(req.Image)
		if err != nil {
			code = "[HANDLER] UploadImageR2 - 5"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()
			return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
		}
	}

	urlImageResp := map[string]interface{}{
		"urlImage": imageUrl,
	}

	successResponseDefault.Meta.Status = true
	successResponseDefault.Meta.Message = "Image uploaded successfully"
	successResponseDefault.Data = urlImageResp

	return c.JSON(successResponseDefault)
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{
		contentService: contentService,
	}
}
