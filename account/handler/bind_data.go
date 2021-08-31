package handler

import (
	"github.com/ELOATS/studyMem/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
)

// used to help extra validation errors
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func bindData(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)

		if errs, ok := err.(validator.ValidationErrors); ok {
			var invalidArgs []invalidArgument

			for _, val := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					Field: val.Field(),
					Value: val.Value().(string),
					Tag:   val.Tag(),
					Param: val.Param(),
				})
			}

			err := apperrors.NewBadRequest("Invalid request parameters.")
			c.JSON(err.Status(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return false
		}
		fallBack := apperrors.NewInternal()
		c.JSON(fallBack.Status(), gin.H{"error": fallBack})
		return false
	}
	return true
}
