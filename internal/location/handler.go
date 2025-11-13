package location

import (
	"log"
	"net/http"

	"lam-phuong-api/internal/airtable"
	"lam-phuong-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

// Handler exposes HTTP handlers for the location resource.
type Handler struct {
	repo           Repository
	airtableClient *airtable.Client
	airtableTable  string
}

// NewHandler creates a handler with the provided repository and config.
func NewHandler(repo Repository, cfg *config.Config) (*Handler, error) {
	// Create Airtable client from config
	airtableClient, err := cfg.NewAirtableClient()
	if err != nil {
		return nil, err
	}

	return &Handler{
		repo:           repo,
		airtableClient: airtableClient,
		airtableTable:  cfg.Airtable.LocationsTableName,
	}, nil
}

// RegisterRoutes attaches location routes to the supplied router group.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/locations", h.ListLocations)
	router.GET("/locations/:id", h.GetLocation)
	router.POST("/locations", h.CreateLocation)
	router.PUT("/locations/:id", h.UpdateLocation)
	router.DELETE("/locations/:id", h.DeleteLocation)
}

func (h *Handler) ListLocations(c *gin.Context) {
	c.JSON(http.StatusOK, h.repo.List())
}

func (h *Handler) GetLocation(c *gin.Context) {
	id := c.Param("id")
	location, ok := h.repo.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}

	c.JSON(http.StatusOK, location)
}

func (h *Handler) CreateLocation(c *gin.Context) {
	var payload locationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate slug from name if not provided
	locationSlug := payload.Slug
	if locationSlug == "" {
		locationSlug = slug.Make(payload.Name)
	}

	// Validate status if provided, default to active
	status := StatusActive
	if payload.Status != "" {
		status = Status(payload.Status)
		if !status.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status. must be 'Active' or 'Disabled'"})
			return
		}
	}

	location := Location{
		Name:   payload.Name,
		Slug:   locationSlug,
		Status: status,
	}

	// Create in repository first
	created, err := h.repo.Create(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save to Airtable
	airtableFields := created.ToAirtableFields()
	airtableRecord, err := h.airtableClient.CreateRecord(c.Request.Context(), h.airtableTable, airtableFields)
	if err != nil {
		// Log error but don't fail the request - location is already created in repo
		log.Printf("Failed to save location to Airtable: %v", err)
		// Optionally, you could return an error here if you want to ensure Airtable sync
	} else {
		// Update the created location with Airtable ID if needed
		created.ID = airtableRecord.ID
		log.Printf("Location saved to Airtable with ID: %s", airtableRecord.ID)
	}

	c.JSON(http.StatusCreated, created)
}

func (h *Handler) UpdateLocation(c *gin.Context) {
	id := c.Param("id")

	var payload locationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status if provided
	var status Status
	if payload.Status != "" {
		status = Status(payload.Status)
		if !status.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status. must be 'active' or 'disabled'"})
			return
		}
	} else {
		// Get existing location to preserve status if not provided
		existing, ok := h.repo.Get(id)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
			return
		}
		status = existing.Status
	}

	location := Location{
		Name:   payload.Name,
		Slug:   payload.Slug,
		Status: status,
	}

	updated, ok := h.repo.Update(id, location)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *Handler) DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	if ok := h.repo.Delete(id); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

type locationPayload struct {
	Name   string `json:"name" binding:"required"` // Required
	Slug   string `json:"slug"`                    // Optional, will be generated from name if not provided
	Status string `json:"status"`                  // Optional, defaults to "Active"
}
