package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pndwrzk/go-article/internal/common/constants"
	"github.com/pndwrzk/go-article/internal/common/dto"
)

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, dto.BaseResponse{
		Code:    constants.CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code string, message string) {
	c.JSON(http.StatusInternalServerError, dto.BaseResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func Failed(c *gin.Context, code string, message string, data interface{}) {
	c.JSON(http.StatusBadRequest, dto.BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, message string, data interface{}, meta *dto.Meta) {
	c.JSON(http.StatusOK, dto.BaseResponse{
		Code:    constants.CodeSuccess,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
