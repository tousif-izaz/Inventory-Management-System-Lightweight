package controller

import (
	"github.com/labstack/echo/v4"
	"inventory-management/pkg/controller/request"
	"inventory-management/pkg/controller/response"
	"inventory-management/pkg/domain"
	"inventory-management/pkg/middleware"
	"inventory-management/pkg/service"
	"net/http"
	"strconv"
)

type CategoryController struct {
	categoryService service.ICategoryService
}

func NewCategoryController(categoryService service.ICategoryService) *CategoryController {
	return &CategoryController{categoryService}
}

func (controller *CategoryController) RegisterCategoryRoutes(e *echo.Echo) {
	categoriesGroup := e.Group("/categories")
	categoriesGroup.Use(middleware.AuthMiddleware)

	categoriesGroup.GET("", controller.GetAllCategories)
	categoriesGroup.POST("", controller.AddNewCategory)
	categoriesGroup.GET("/:id", controller.GetCategoryByID)
	categoriesGroup.GET("/:id/children", controller.GetChildCategories)
	categoriesGroup.PUT("/:id", controller.UpdateCategoryById)
	categoriesGroup.DELETE("/:id", controller.DeleteCategoryById)
}

func (controller *CategoryController) GetAllCategories(c echo.Context) error {
	categories, err := controller.categoryService.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.ToCategoryResponseList(categories))
}

func (controller *CategoryController) GetCategoryByID(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no category id specified"))
	}

	categoryId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: category id must be an integer"))
	}

	category, err := controller.categoryService.GetCategoryByID(categoryId)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.ToCategoryResponseList([]*domain.Category{category})[0])
}

func (controller *CategoryController) GetChildCategories(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no category id specified"))
	}

	categoryId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: category id must be an integer"))
	}

	categories, err := controller.categoryService.GetChildCategories(categoryId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.ToCategoryResponseList(categories))
}

func (controller *CategoryController) AddNewCategory(c echo.Context) error {
	addCategoryRequest := new(request.CreateCategoryRequest)

	err := c.Bind(addCategoryRequest)
	if err != nil || addCategoryRequest == nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data to the category structure"))
	}

	err = controller.categoryService.Add(addCategoryRequest.ToDTO())
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusCreated)
}

func (controller *CategoryController) UpdateCategoryById(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no category id specified"))
	}

	categoryId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: category id must be an integer"))
	}

	updateCategoryRequest := new(request.UpdateCategoryRequest)
	err = c.Bind(updateCategoryRequest)
	if err != nil || updateCategoryRequest == nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data to the category structure"))
	}

	err = controller.categoryService.UpdateCategoryById(updateCategoryRequest.ToDTO(), categoryId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (controller *CategoryController) DeleteCategoryById(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no category id specified"))
	}

	categoryId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: category id must be an integer"))
	}

	err = controller.categoryService.DeleteById(categoryId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}
