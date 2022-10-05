package article

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	v2 := e.Group("/v2")
	{
		v2.GET("/article", helloHandler)
		v2.POST("/article", AddArticle)
		v2.GET("/article/:tag", GetUserByTag)
		v2.GET("/article/", GetAllArticles)
		v2.GET("/hello", helloHandler)
		v2.GET("/hello1", helloHandler)

	}

}
