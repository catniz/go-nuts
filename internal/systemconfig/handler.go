package systemconfig

import (
	"github.com/gin-gonic/gin"
)

type handlerImpl struct{}

func (u *handlerImpl) RegisterRoutes(router gin.IRouter) {
	r := router.Group(u.Path())
	r.GET("/all", u.GetAllSystemConfig)
}

var _ = &handlerImpl{}

func (u *handlerImpl) Path() string {
	return "/systemconfig"
}

func (u *handlerImpl) GetAllSystemConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "GetAllSystemConfig",
	})
}
