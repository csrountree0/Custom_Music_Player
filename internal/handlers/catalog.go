package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"musicapp/backend/internal/models"
	"musicapp/backend/internal/services"
	"musicapp/backend/pkg/apierror"
	"musicapp/backend/pkg/pagination"
)

type CatalogHandler struct {
	catalogSvc *services.CatalogService
}

func NewCatalogHandler(catalogSvc *services.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalogSvc: catalogSvc}
}

func (h *CatalogHandler) ListArtists(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	p := pagination.FromQuery(c)

	artists, total, err := h.catalogSvc.ListArtists(userID, p)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, pagination.Page[models.Artist]{Items: artists, Total: total, Page: p.Page, Limit: p.Limit})
}

func (h *CatalogHandler) CreateArtist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input services.CreateArtistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	artist, err := h.catalogSvc.CreateArtist(userID, input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, artist)
}

func (h *CatalogHandler) GetArtist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	artist, err := h.catalogSvc.GetArtist(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, artist)
}

func (h *CatalogHandler) UpdateArtist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	var input services.UpdateArtistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	artist, err := h.catalogSvc.UpdateArtist(id, userID, input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, artist)
}

func (h *CatalogHandler) DeleteArtist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	if err := h.catalogSvc.DeleteArtist(id, userID); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CatalogHandler) ListAlbums(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	p := pagination.FromQuery(c)

	albums, total, err := h.catalogSvc.ListAlbums(userID, p)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, pagination.Page[models.Album]{Items: albums, Total: total, Page: p.Page, Limit: p.Limit})
}

func (h *CatalogHandler) CreateAlbum(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input services.CreateAlbumInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	album, err := h.catalogSvc.CreateAlbum(userID, input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, album)
}

func (h *CatalogHandler) GetAlbum(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	album, err := h.catalogSvc.GetAlbum(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, album)
}

func (h *CatalogHandler) UpdateAlbum(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	var input services.UpdateAlbumInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	album, err := h.catalogSvc.UpdateAlbum(id, userID, input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, album)
}

func (h *CatalogHandler) DeleteAlbum(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	if err := h.catalogSvc.DeleteAlbum(id, userID); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CatalogHandler) ListTracks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	p := pagination.FromQuery(c)

	tracks, total, err := h.catalogSvc.ListTracks(userID, p)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, pagination.Page[models.Track]{Items: tracks, Total: total, Page: p.Page, Limit: p.Limit})
}

func (h *CatalogHandler) GetTrack(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	track, err := h.catalogSvc.GetTrack(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, track)
}

func (h *CatalogHandler) UpdateTrack(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	var input services.UpdateTrackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	track, err := h.catalogSvc.UpdateTrack(id, userID, input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, track)
}

func (h *CatalogHandler) DeleteTrack(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	if err := h.catalogSvc.DeleteTrack(id, userID); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
