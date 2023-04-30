package httpi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sherpalden/go-saas-template/internal/app_error"
)

func JSON(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, gin.H{"data": data})
}

func JSONWithMetaData(ctx *gin.Context, statusCode int, data interface{}, paginationReponse PaginationResponse) {
	ctx.JSON(statusCode, gin.H{"data": data, "meta_data": paginationReponse})
}

func HandleError(ctx *gin.Context, err error) {
	if errResp, ok := err.(*app_error.AppError); ok {
		ctx.JSON(int(errResp.Code), gin.H{"error": errResp})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": app_error.AppError{Code: http.StatusInternalServerError, Message: err.Error()}})
}
