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

type PlaylistHandler struct {
	playlistSvc *services.PlaylistService
}

func NewPlaylistHandler(playlistSvc *services.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{playlistSvc: playlistSvc}
}

func (h *PlaylistHandler) ListPlaylists(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	p := pagination.FromQuery(c)

	playlists, total, err := h.playlistSvc.ListPlaylists(userID, p)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, pagination.Page[models.Playlist]{Items: playlists, Total: total, Page: p.Page, Limit: p.Limit})
}

func (h *PlaylistHandler) CreatePlaylist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input services.CreatePlaylistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	playlist, err := h.playlistSvc.CreatePlaylist(userID, input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, playlist)
}

func (h *PlaylistHandler) GetPlaylist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	playlist, err := h.playlistSvc.GetPlaylist(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, playlist)
}

func (h *PlaylistHandler) UpdatePlaylist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	var input services.UpdatePlaylistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	playlist, err := h.playlistSvc.UpdatePlaylist(id, userID, input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, playlist)
}

func (h *PlaylistHandler) DeletePlaylist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	if err := h.playlistSvc.DeletePlaylist(id, userID); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PlaylistHandler) ListPlaylistTracks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	tracks, err := h.playlistSvc.ListPlaylistTracks(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, tracks)
}

func (h *PlaylistHandler) AddTrackToPlaylist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	var input services.AddTrackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	if err := h.playlistSvc.AddTrack(id, userID, input); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *PlaylistHandler) RemoveTrackFromPlaylist(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}
	trackID, err := uuid.Parse(c.Param("trackID"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "trackID must be a valid UUID").Respond(c)
		return
	}

	if err := h.playlistSvc.RemoveTrack(id, trackID, userID); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PlaylistHandler) ReorderPlaylistTracks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	var input services.ReorderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	if err := h.playlistSvc.ReorderTracks(id, userID, input); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
