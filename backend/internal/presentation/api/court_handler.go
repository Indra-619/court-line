package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/book-lapangan/backend/internal/domain/entity"
	"github.com/your-username/book-lapangan/backend/internal/usecase"
)

type CourtHandler struct {
	courtUseCase usecase.CourtUseCase
}

func NewCourtHandler(u usecase.CourtUseCase) *CourtHandler {
	return &CourtHandler{
		courtUseCase: u,
	}
}

func (h *CourtHandler) GetAllCourts(c *gin.Context) {
	result := h.courtUseCase.GetAllCourts(c.Request.Context())
	if result.IsFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Value})
}

func (h *CourtHandler) GetCourtByID(c *gin.Context) {
	id := c.Param("id")
	result := h.courtUseCase.GetCourtByID(c.Request.Context(), id)
	if result.IsFailure() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Value})
}

func (h *CourtHandler) CreateCourt(c *gin.Context) {
	var input entity.Court // Reusing entity for simplicity, or use DTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.courtUseCase.CreateCourt(c.Request.Context(), &input)
	if result.IsFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": result.Value})
}

func (h *CourtHandler) UpdateCourt(c *gin.Context) {
	id := c.Param("id")
	var input entity.Court
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = id

	result := h.courtUseCase.UpdateCourt(c.Request.Context(), &input)
	if result.IsFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": result.Value})
}

func (h *CourtHandler) DeleteCourt(c *gin.Context) {
	id := c.Param("id")
	result := h.courtUseCase.DeleteCourt(c.Request.Context(), id)
	if result.IsFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": result.Value})
}
