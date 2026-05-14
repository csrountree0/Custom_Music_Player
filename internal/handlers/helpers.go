package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"musicapp/backend/pkg/apierror"
)

func respondError(c *gin.Context, err error) {
	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		apiErr.Respond(c)
		return
	}
	apierror.ErrInternal.Respond(c)
}
