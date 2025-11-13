package location

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

// Handler exposes HTTP handlers for the location resource.
type Handler struct {
	repo Repository
}

// NewHandler creates a handler with the provided repository.
func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// RegisterRoutes attaches location routes to the supplied router group.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/locations", h.ListLocations)
	router.POST("/locations", h.CreateLocation)
	router.DELETE("/locations/:slug", h.DeleteLocationBySlug)
}

func (h *Handler) ListLocations(c *gin.Context) {
	c.JSON(http.StatusOK, h.repo.List())
}

func (h *Handler) CreateLocation(c *gin.Context) {
	var payload locationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate slug from name if not provided
	locationSlug := payload.Slug
	if locationSlug != "" {
		locationSlug = slug.Make(locationSlug)
	} else {
		locationSlug = slug.Make(payload.Name)
	}

	locationSlug = ensureUniqueSlug(h.repo, locationSlug)

	location := Location{
		Name: payload.Name,
		Slug: locationSlug,
	}

	// Create in repository (repository handles Airtable sync if configured)
	created, err := h.repo.Create(c.Request.Context(), location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

type locationPayload struct {
	Name string `json:"name" binding:"required"` // Required
	Slug string `json:"slug"`                    // Optional, will be generated from name if not provided
}

func ensureUniqueSlug(repo Repository, baseSlug string) string {
	if baseSlug == "" {
		baseSlug = "location"
	}

	existingSlugs := make(map[string]struct{})
	for _, loc := range repo.List() {
		existingSlugs[loc.Slug] = struct{}{}
	}

	if _, exists := existingSlugs[baseSlug]; !exists {
		return baseSlug
	}

	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s-%d", baseSlug, i)
		if _, exists := existingSlugs[candidate]; !exists {
			return candidate
		}
	}
}

func (h *Handler) DeleteLocationBySlug(c *gin.Context) {
	slugParam := c.Param("slug")
	if slugParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	normalizedSlug := slug.Make(slugParam)
	if normalizedSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slug"})
		return
	}

	if ok := h.repo.DeleteBySlug(normalizedSlug); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
