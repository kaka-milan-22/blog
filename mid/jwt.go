package mid

import (
	"blog/config"
	"net/http"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
)

var jc = config.Cfg.Authentication

// AuthClaims 是 claims struct
type AuthClaims struct {
	UserId uint64 `json:"userId"`
	jwt.StandardClaims
}

// JWTAuth 鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 登录不需 要密码
		if c.FullPath() == "/v1/login" ||
			c.FullPath() == "/v1/permission" ||
			c.FullPath() == "/v2/hello" ||
			strings.Contains(c.FullPath(), "swagger") ||
			strings.Contains(c.FullPath(), "favicon.ico") {
			c.Next()
			return
		}
		// 获取请求头中 token，实际是一个完整被签名过的 token；a complete, signed token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusForbidden, "No token. You don't have permission!")
			c.Abort()
			return
		}

		// 解析拿到完整有效的 token，里头包含解析后的 3 segment
		token, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusForbidden, "Invalid token! You don't have permission!")
			c.Abort()
			return
		}
		// 获取 token 中的 claims
		claims, ok := token.Claims.(*AuthClaims)
		if !ok {
			c.JSON(http.StatusForbidden, "Invalid token!")
			c.Abort()
			return
		}

		// 将 claims 中的用户信息存储在 context 中
		c.Set("userId", claims.UserId)

		// 这里执行路由 HandlerFunc
		c.Next()
	}
}

// ParseToken 解析请求头中的 token string，转换成被解析后的 jwt.Token
func ParseToken(tokenStr string) (*jwt.Token, error) {
	// 解析 token string 拿到 token jwt.Token
	return jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(tk *jwt.Token) (interface{}, error) {
		return []byte(jc.JwtKey), nil
	})
}

func GenerateToken(userId uint64, expireTime time.Time) (string, error) {
	// 创建一个 claim
	claim := AuthClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 签名时间
			IssuedAt: time.Now().Unix(),
			// 签名颁发者
			Issuer: "abcnull",
			// 签名主题
			Subject: "gindemo",
		},
	}

	// 使用指定的签名加密方式创建 token，有 1，2 段内容，第 3 段内容没有加上
	noSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// 使用 secretKey 密钥进行加密处理后拿到最终 token string
	return noSignedToken.SignedString([]byte(jc.JwtKey))
}
