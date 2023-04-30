package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sherpalden/go-saas-template/external/auth"
	"github.com/sherpalden/go-saas-template/internal/api/httpi"
	"github.com/sherpalden/go-saas-template/internal/app_error"
)

type AuthMiddleware struct {
	authClient auth.AuthClient
}

func NewAuthMiddleware(
	authClient auth.AuthClient,
) AuthMiddleware {
	return AuthMiddleware{
		authClient: authClient,
	}
}

// Handle auth requests
func (am AuthMiddleware) Handle(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	accessToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
	authUser, err := am.authClient.VerifyAccessToken(accessToken)
	if err != nil {
		err = app_error.New(err, http.StatusUnauthorized)
		httpi.HandleError(ctx, err)
		ctx.Abort()
		return
	}
	ctx.Set("auth_user", authUser)
	ctx.Next()
}

func (am AuthMiddleware) AdminOnly(ctx *gin.Context) {
	authUser, exists := ctx.Get("auth_user")
	if !exists || !authUser.(map[string]interface{})["is_admin"].(bool) {
		err := app_error.New(errors.New("Unauthorised"), http.StatusUnauthorized)
		httpi.HandleError(ctx, err)
		ctx.Abort()
		return
	}

	ctx.Next()
}
