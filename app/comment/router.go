package comment

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	v3 := e.Group("/v3")
	{
		v3.GET("/hello", helloHandler)
		v3.POST("/comment", AddComment)
		v3.GET("/migration", MigrateDB)
		v3.GET("/comment/:aid", GetCommentsByAid)
		v3.DELETE("/commenst/:id", DeleteComment)

	}

}
