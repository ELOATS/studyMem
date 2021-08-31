package handler

import (
	"github.com/ELOATS/studyMem/account/model"
	"github.com/ELOATS/studyMem/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"log"
)

type SignupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupReq

	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.UserService.Signup(c, u)
	if err != nil {
		log.Printf("Fail to sign up user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{"error": err})
		return
	}
}
