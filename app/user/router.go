package user

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	v1 := e.Group("/v1")
	{
		v1.GET("/tt", helloHandler)
		v1.POST("/user", AddUser)
		v1.POST("/login", LoginUser)
		v1.GET("/users/:id", GetUserById)
		v1.GET("/users", GetAllUsers)
		v1.DELETE("/users/:id", DeleteUser) // new
		v1.POST("/permission", AddPermission)
	}

}
