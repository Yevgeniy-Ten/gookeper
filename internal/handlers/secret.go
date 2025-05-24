package handlers

import (
	"net/http"
	"strconv"

	"gokeeper/internal/models"
	"gokeeper/internal/service"

	"github.com/gin-gonic/gin"
)

type SecretHandler struct {
	secretService *service.SecretService
}

func NewSecretHandler(secretService *service.SecretService) *SecretHandler {
	return &SecretHandler{secretService: secretService}
}

func getUserIDFromContext(c *gin.Context) int64 {
	id, _ := c.Get("user_id")
	return id.(int64)
}

func (h *SecretHandler) Create(c *gin.Context) {
	var input models.CreateSecretDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	userID := getUserIDFromContext(c)
	secret, err := h.secretService.Create(c.Request.Context(), userID, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, secret)
}

func (h *SecretHandler) List(c *gin.Context) {
	userID := getUserIDFromContext(c)
	secrets, err := h.secretService.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, secrets)
}

func (h *SecretHandler) Get(c *gin.Context) {
	userID := getUserIDFromContext(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	secret, err := h.secretService.Get(c.Request.Context(), userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if secret == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
		return
	}
	c.JSON(http.StatusOK, secret)
}

func (h *SecretHandler) Delete(c *gin.Context) {
	userID := getUserIDFromContext(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	if err := h.secretService.Delete(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
