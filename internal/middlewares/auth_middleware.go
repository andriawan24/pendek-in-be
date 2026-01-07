package middlewares

import (
	"strings"

	"github.com/andriawan24/link-short/internal/utils"
	"github.com/gin-gonic/gin"
)

func RequiredAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondUnauthorized(ctx, "Unauthorized")
			ctx.Abort()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.RespondUnauthorized(ctx, "Unauthorized")
			ctx.Abort()
			return
		}

		token := headerParts[1]

		claim, err := utils.ParseToken(token)
		if err != nil {
			utils.RespondUnauthorized(ctx, "Unauthorized: "+err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claim.UserId)
		ctx.Next()
	}
}
