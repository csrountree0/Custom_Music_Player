package middleware

import (
	"github.com/gin-gonic/gin"

	"musicapp/backend/pkg/apierror"
	"musicapp/backend/pkg/token"
)

func AuthRequired(tokenSvc *token.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			apierror.ErrUnauthorized.Respond(c)
			return
		}

		tokenStr, err := token.ExtractFromHeader(header)
		if err != nil {
			apierror.New(401, "INVALID_TOKEN_FORMAT", err.Error()).Respond(c)
			return
		}

		claims, err := tokenSvc.Validate(tokenStr)
		if err != nil {
			apierror.New(401, "INVALID_TOKEN", "token is invalid or expired").Respond(c)
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
