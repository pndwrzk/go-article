package utils

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pndwrzk/go-article/internal/common/constants"
	"github.com/pndwrzk/go-article/internal/common/response"
)

func BindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBind(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = fmt.Sprintf("failed on the '%s' tag", fe.Tag())
			}
			response.Failed(c, constants.CodeBadRequest, "Validation failed", out)
			return false
		}
		response.Failed(c, constants.CodeBadRequest, err.Error(), nil)
		return false
	}
	return true
}
