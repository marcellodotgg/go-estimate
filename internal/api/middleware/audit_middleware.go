package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Audit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := ctx.Cookie("user_id")

		if err != nil {
			userID = uuid.NewString()
			ctx.SetCookie("user_id", userID, 3600*24*365, "/", os.Getenv("COOKIE_URL"), true, true)
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
