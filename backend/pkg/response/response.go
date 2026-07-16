package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamtbay/tyr-fintech/pkg/apperrors"
)

func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		c.JSON(appErr.Code, gin.H{
			"success": false,
			"error":   appErr.Msg,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   "unexpected error occured",
	})
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

func Error(c *gin.Context, status int, errMessage string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   errMessage,
	})
}
