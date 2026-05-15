package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"musicapp/backend/internal/services"
	"musicapp/backend/pkg/apierror"
)

type UploadHandler struct {
	uploadSvc *services.UploadService
}

func NewUploadHandler(uploadSvc *services.UploadService) *UploadHandler {
	return &UploadHandler{uploadSvc: uploadSvc}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		apierror.New(400, "FILE_REQUIRED", "file field is required").Respond(c)
		return
	}
	defer file.Close()

	job, err := h.uploadSvc.UploadFile(userID, file, header)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, job)
}

func (h *UploadHandler) ImportYouTube(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var body struct {
		URL string `json:"url" binding:"required,url"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	job, err := h.uploadSvc.ImportYouTube(userID, body.URL)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, job)
}

func (h *UploadHandler) GetUploadJob(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apierror.New(400, "INVALID_ID", "id must be a valid UUID").Respond(c)
		return
	}

	job, err := h.uploadSvc.GetJob(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, job)
}
