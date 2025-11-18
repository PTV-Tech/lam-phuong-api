package location

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"lam-phuong-api/internal/response"
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

// ListLocations godoc
// @Summary      List all locations
// @Description  Get a list of all locations (requires authentication)
// @Tags         locations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  location.LocationsResponseWrapper  "Locations retrieved successfully"
// @Failure      401  {object}  response.ErrorResponse  "Unauthorized"
// @Router       /locations [get]
func (h *Handler) ListLocations(c *gin.Context) {
	locations := h.repo.List()
	response.Success(c, http.StatusOK, locations, "Locations retrieved successfully")
}

// CreateLocation godoc
// @Summary      Create a new location
// @Description  Create a new location with name and optional slug. If slug is not provided, it will be generated from the name. (requires authentication)
// @Tags         locations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        location  body      locationPayload  true  "Location payload"
// @Success      201       {object}  location.LocationResponseWrapper  "Location created successfully"
// @Failure      400       {object}  response.ErrorResponse  "Validation error"
// @Failure      401       {object}  response.ErrorResponse  "Unauthorized"
// @Failure      500       {object}  response.ErrorResponse  "Internal server error"
// @Router       /locations [post]
func (h *Handler) CreateLocation(c *gin.Context) {
	var payload locationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]interface{}{
			"validation_error": err.Error(),
		})
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
		response.InternalError(c, "Failed to create location: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, created, "Location created successfully")
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

// DeleteLocationBySlug godoc
// @Summary      Delete a location by slug
// @Description  Delete a location using its slug (requires authentication)
// @Tags         locations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug  path      string  true  "Location slug"
// @Success      200   {object}  response.Response  "Location deleted successfully"
// @Failure      400   {object}  response.ErrorResponse  "Validation error"
// @Failure      401   {object}  response.ErrorResponse  "Unauthorized"
// @Failure      404   {object}  response.ErrorResponse  "Location not found"
// @Router       /locations/{slug} [delete]
func (h *Handler) DeleteLocationBySlug(c *gin.Context) {
	slugParam := c.Param("slug")
	if slugParam == "" {
		response.BadRequest(c, "Slug is required", nil)
		return
	}

	normalizedSlug := slug.Make(slugParam)
	if normalizedSlug == "" {
		response.ValidationError(c, "Invalid slug format", nil)
		return
	}

	if ok := h.repo.DeleteBySlug(normalizedSlug); !ok {
		response.NotFound(c, "Location")
		return
	}

	response.SuccessNoContent(c, "Location deleted successfully")
}
