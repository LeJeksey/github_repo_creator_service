package marcopolo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const polo = "polo"

func Marco(ctx *gin.Context) {
	ctx.String(http.StatusOK, polo)
}
