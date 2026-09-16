package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/domain/repository"
	"github.com/Indra-619/court-line/backend/internal/models"
)

// CourtHandler serves the /api/courts endpoints on top of an injected
// repository; it never touches the database driver directly.
type CourtHandler struct {
	repo repository.CourtRepository
}

// NewCourtHandler builds a CourtHandler backed by the given repository.
func NewCourtHandler(repo repository.CourtRepository) *CourtHandler {
	return &CourtHandler{repo: repo}
}

// GetCourts returns all courts
func (h *CourtHandler) GetCourts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	courts, err := h.repo.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courts"})
		return
	}

	if courts == nil {
		courts = []*entity.Court{}
	}

	c.JSON(http.StatusOK, gin.H{"data": courts})
}

// GetCourtByID returns a single court by ID
func (h *CourtHandler) GetCourtByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	if _, err := parseHexID(idParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	court, err := h.repo.FindByID(ctx, idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": court})
}

// CreateCourt creates a new court
func (h *CourtHandler) CreateCourt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var input models.CreateCourtInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	court := &entity.Court{
		Name:         input.Name,
		Type:         input.Type,
		Location:     input.Location,
		Description:  input.Description,
		PricePerHour: input.PricePerHour,
		ImageURL:     input.ImageURL,
		Facilities:   input.Facilities,
		IsAvailable:  true,
	}

	if err := h.repo.Create(ctx, court); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create court"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": court})
}

// UpdateCourt updates an existing court
func (h *CourtHandler) UpdateCourt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	if _, err := parseHexID(idParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	var input models.CreateCourtInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	court := &entity.Court{
		ID:           idParam,
		Name:         input.Name,
		Type:         input.Type,
		Location:     input.Location,
		Description:  input.Description,
		PricePerHour: input.PricePerHour,
		ImageURL:     input.ImageURL,
		Facilities:   input.Facilities,
	}

	if err := h.repo.Update(ctx, court); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update court"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Court updated successfully"})
}

// DeleteCourt deletes a court
func (h *CourtHandler) DeleteCourt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	if _, err := parseHexID(idParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	if err := h.repo.Delete(ctx, idParam); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete court"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Court deleted successfully"})
}
