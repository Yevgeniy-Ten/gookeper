package handlers

import (
	"net/http"

	"gokeeper/internal/config"
	"gokeeper/internal/models"
	"gokeeper/internal/pkg/jwt"
	"gokeeper/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
	cfg         *config.Config
}

func NewUserHandler(userService *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		cfg:         cfg,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input models.CreateUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &input)
	if err != nil {
		switch err {
		case service.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var input models.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.Login(c.Request.Context(), &input)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		}
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Login, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}
