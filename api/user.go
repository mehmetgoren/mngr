package api

import (
	"github.com/gin-gonic/gin"
	"mngr/models"
	"mngr/reps"
	"net/http"
)

func RegisterUserEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.POST("/login", func(ctx *gin.Context) {
		var lu models.LoginUserViewModel
		if err := ctx.BindJSON(&lu); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, err := rb.UserRep.Login(&lu)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if u != nil {
			rb.Users[u.Token] = u
		}
		ctx.JSON(http.StatusOK, u)
	})

	router.POST("/registeruser", func(ctx *gin.Context) {
		var ru models.RegisterUserViewModel
		if err := ctx.BindJSON(&ru); err != nil {
			ctx.JSON(http.StatusBadRequest, false)
			return
		}
		if ru.Password != ru.RePassword {
			ctx.JSON(http.StatusBadRequest, false)
			return
		}
		u, err := rb.UserRep.Login(&models.LoginUserViewModel{Username: ru.Username, Password: ru.Password})
		if u != nil && err == nil {
			ctx.JSON(http.StatusNotFound, true)
			return
		}
		u, err = rb.UserRep.Register(&ru)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, false)
			return
		}
		if u != nil {
			rb.Users[u.Token] = u
		}
		ctx.JSON(http.StatusOK, true)
	})
}
