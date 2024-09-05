package internal

import "github.com/gin-gonic/gin"

type Handler interface {
	Path() string
	RegisterRoutes(router gin.IRouter)
}
