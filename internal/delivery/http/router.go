package http

import (
	"github.com/cukhoaimon/SimpleBank/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

func (handler *Handler) SetupRouter() {
	router := gin.Default()
	router.POST("/api/v1/user", handler.createUser)
	router.POST("/api/v1/user/login", handler.loginUser)
	router.POST("/api/v1/user/token/renew_access", handler.renewAccessTokenUser)

	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(handler.TokenMaker))

	authRoutes.GET("/api/v1/account", handler.listAccount)
	authRoutes.GET("/api/v1/account/:id", handler.getAccount)
	authRoutes.POST("/api/v1/account", handler.createAccount)

	authRoutes.POST("/api/v1/transfer", handler.createTransfer)

	handler.Router = router
}
