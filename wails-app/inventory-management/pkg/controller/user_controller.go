package controller

import (
	"github.com/labstack/echo/v4"
	"inventory-management/pkg/controller/request"
	"inventory-management/pkg/controller/response"
	"inventory-management/pkg/domain"
	"inventory-management/pkg/service"
	"net/http"
	"time"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{userService}
}

func (controller *UserController) RegisterUserRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.POST("/login", controller.Login)
	e.POST("/signup", controller.SignUp)
	e.POST("/logout", controller.Logout)

	// Protected routes
	e.GET("/profile", controller.GetProfile, authMiddleware)
	e.PUT("/profile/password", controller.UpdatePassword, authMiddleware)
}

func (controller *UserController) Login(c echo.Context) error {
	var loginRequest request.LoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data to the user structure"))
	}

	token, err := controller.userService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	// Set cookie for browser-based access
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = false // Allow JavaScript access for Wails app
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	// Also return token in response body for Wails app
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (controller *UserController) SignUp(c echo.Context) error {
	var signUpRequest request.SignUpRequest
	err := c.Bind(&signUpRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data to the user structure"))
	}

	err = controller.userService.SignUp(signUpRequest.ToDtoModel())
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusCreated)
}

func (controller *UserController) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-24 * time.Hour)
	cookie.MaxAge = -1
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}

func (controller *UserController) GetProfile(c echo.Context) error {
	// Get username from JWT claims
	claims := c.Get("user").(*domain.Claims)
	username := claims.Username

	// Fetch user profile
	userProfile, err := controller.userService.GetUserByUsername(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NewErrorResponse("User not found"))
	}

	return c.JSON(http.StatusOK, response.ToUserProfileResponse(userProfile))
}

func (controller *UserController) UpdatePassword(c echo.Context) error {
	// Get username from JWT claims
	claims := c.Get("user").(*domain.Claims)
	username := claims.Username

	// Bind request
	var updatePasswordRequest request.UpdatePasswordRequest
	err := c.Bind(&updatePasswordRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request"))
	}

	// Update password
	err = controller.userService.UpdatePassword(username, updatePasswordRequest.OldPassword, updatePasswordRequest.NewPassword)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
