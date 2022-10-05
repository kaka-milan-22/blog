package user

import (
	"blog/app/model"
	"blog/config"
	"blog/mid"
	"blog/routers"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var jc = config.Cfg.Authentication

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Passwd   string `json:"passwd" binding:"required"`
	Github   string `json:"github"`
}

type LoginUserInput struct {
	Username string `json:"username" binding:"required"`
	Passwd   string `json:"passwd" binding:"required"`
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": model.DB,
	})
}

// @Description 添加用户
// @Tags Blog文档
// @Accept json
// @Param username formData string true "用户名"
// @Param passwd formData string true "密码"
// @Param email formData string true "邮箱"
// @Param github formData string false "github"
// @Success 200 {string} string "message": fmt.Sprint(userJson.Username, " insert sucessfully!")"
// @Failure 400 {string} string "message": fmt.Sprint(userJson.Username, " insert failed!"),"
// @Router /user [post]
func AddUser(c *gin.Context) {
	var userJson CreateUserInput
	if err := c.ShouldBindJSON(&userJson); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passHash, _ := HashPassword(userJson.Passwd)
	u1 := model.User{
		UserName: userJson.Username,
		Email:    userJson.Email,
		Passwd:   passHash,
		Github:   userJson.Github,
	}
	result := model.DB.Where("user_name = ?", u1.UserName).Or("email=?", u1.Email).First(&u1)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprint(userJson.Username, " is already Exists!"),
		})
		return
	}

	model.DB.AutoMigrate(&model.User{})
	r1 := model.DB.Create(&u1)
	if r1.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprint(userJson.Username, " insert failed!"),
			"desc":    r1.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprint(userJson.Username, " insert sucessfully!"),
	})
}

// @Description 通过ID查询用户信息
// @Tags Blog文档
// @Accept json
// @Param int path string true "id"
// @Success 200 {string} string {"error": "Record not found!"}
// @Failure 400 {string} string {"data": user}
// @Router /user/:id [get]
func GetUserById(c *gin.Context) {
	var u1 model.User

	if err := model.DB.Where("id = ?", c.Param("id")).First(&u1).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": u1})

}

func GetAllUsers(c *gin.Context) {
	var users []model.User
	model.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// DELETE /usrs/:id
// Delete a user
func DeleteUser(c *gin.Context) {
	// Get model if exist
	var u1 model.User
	if err := model.DB.Where("id = ?", c.Param("id")).First(&u1).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	model.DB.Delete(&u1)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func LoginUser(c *gin.Context) {
	// Get model if exist
	var u1 model.User
	var loginJson LoginUserInput
	if err := c.ShouldBindJSON(&loginJson); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := model.DB.Where("user_name = ?", loginJson.Username).First(&u1).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found!"})
		return
	}

	if CheckPasswordHash(loginJson.Passwd, u1.Passwd) {
		expTime := time.Now().Add(60 * time.Duration(jc.ValidHour) * time.Hour)
		tokenString, tokenErr := mid.GenerateToken(uint64(u1.ID), expTime)
		if tokenErr != nil {
			c.JSON(http.StatusOK, gin.H{"message": tokenErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"message": "wrong password!"})
}

// add permission of default
func AddPermission(c *gin.Context) {
	fmt.Println("增加Policy")
	if ok := routers.Enforcer.AddPolicy("admin", "/v1/login", "POST"); !ok {
		fmt.Println("Policy已经存在")
	} else {
		fmt.Println("增加成功")
	}
}
