package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"musicapp/backend/internal/services"
	"musicapp/backend/pkg/apierror"
)

type SearchHandler struct {
	searchSvc *services.SearchService
}

func NewSearchHandler(searchSvc *services.SearchService) *SearchHandler {
	return &SearchHandler{searchSvc: searchSvc}
}

func (h *SearchHandler) Search(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		apierror.New(400, "QUERY_REQUIRED", "q query param is required").Respond(c)
		return
	}

	results, err := h.searchSvc.Search(userID, q)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, results)
}
