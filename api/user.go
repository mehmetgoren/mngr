package api

import (
	"github.com/gin-gonic/gin"
	"mngr/models"
	"mngr/ws"
	"net/http"
	"time"
)

func RegisterUserEndpoints(router *gin.Engine, holders *ws.Holders) {
	rb := holders.Rb

	logoutUser := func(user *models.User, triggerLogout bool) {
		rb.RemoveUser(user.Token)
		holders.UserLogout(user.Token, triggerLogout)
	}

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
		logoutUser(u, true)
		time.Sleep(1 * time.Second)
		if u != nil {
			rb.AddUser(u)
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
		ctx.JSON(http.StatusOK, true)
	})

	router.GET("/users", func(ctx *gin.Context) {
		services, err := rb.UserRep.GetUsers()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, services)
	})

	router.DELETE("/users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		user, err := rb.UserRep.GetUser(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user == nil {
			ctx.JSON(http.StatusOK, 0)
			return
		}
		result, err := rb.UserRep.RemoveById(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			logoutUser(user, true)
			ctx.JSON(http.StatusOK, result)
		}
	})

	router.POST("/logoutuser", func(ctx *gin.Context) {
		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(user.Token) == 0 {
			ctx.JSON(http.StatusNotFound, false)
			return
		}
		logoutUser(&user, false)
		ctx.JSON(http.StatusOK, true)
	})
}
