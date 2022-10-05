package article

import (
	"blog/app/model"

	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello article!",
	})
}

type CreateArticleInput struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
	UserID  int      `json:"uid"`
}

func AddArticle(c *gin.Context) {
	var input CreateArticleInput
	var article model.Article
	if err := c.ShouldBindJSON(&input); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article = model.Article{
		Title:   input.Title,
		Content: input.Content,
		Tags:    fmt.Sprint(strings.Join(input.Tags, ",")),
		UserID:  uint(input.UserID),
	}

	model.DB.AutoMigrate(&model.Article{})
	r1 := model.DB.Create(&article)
	if r1.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprint(article.Title, " insert failed!"),
			"desc":    r1.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprint(article.Title, " insert sucessfully!"),
	})
}

func GetUserByTag(c *gin.Context) {
	var a1 model.Article
	tag := fmt.Sprint("%", c.Param("tag"), "%")
	if err := model.DB.Where("tags LIKE?", tag).First(&a1).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": a1})

}

/*
	gorm默认不会添加外键约束，必须做一次Migrate
	https://blog.csdn.net/Xiang_lhh/article/details/113746104
	user := model.User{}
	model.DB.Model(&model.Article{}).AddForeignKey("user_id", "Users(id)", "CASCADE", "CASCADE")
	model.DB.Where("user_name = ?", "Kaka").First(&user)
	fmt.Print(user)
	model.DB.Model(&user).Related(&user.Articles)
*/

func GetAllArticles(c *gin.Context) {
	var a1 []model.Article
	model.DB.Limit(10).Find(&a1)
	c.JSON(http.StatusOK, gin.H{"data": a1})
}
