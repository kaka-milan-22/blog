package mid

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

func Authorize(e *casbin.Enforcer) gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.FullPath() == "/v1/permission" ||
			c.FullPath() == "/v2/hello" ||
			strings.Contains(c.FullPath(), "swagger") ||
			strings.Contains(c.FullPath(), "favicon.ico") {
			c.Next()
			return
		}

		//获取请求的URI
		obj := c.Request.URL.RequestURI()
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色
		sub := "admin"

		//判断策略中是否存在
		if ok, _ := e.EnforceSafe(sub, obj, act); ok {
			fmt.Println("恭喜您,权限验证通过")
			c.Next()
		} else {
			fmt.Println("很遗憾,权限验证没有通过")

			c.JSON(http.StatusForbidden, gin.H{"message": "No permission"})
			c.Abort()
		}
	}
}
