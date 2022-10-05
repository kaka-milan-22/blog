package comment

import (
	"blog/app/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello article!",
	})
}

type CreateCommentInput struct {
	Content   string `json:"content" binding:"required"`
	UserID    int    `json:"uid"`
	ArticleID int    `json:"aid"`
}

func AddComment(c *gin.Context) {
	var input CreateCommentInput
	var comment model.Comment

	if err := c.ShouldBindJSON(&input); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment = model.Comment{
		Content:   input.Content,
		UserID:    uint(input.UserID),
		ArticleID: uint(input.ArticleID),
	}

	model.DB.AutoMigrate(&model.Comment{})
	r1 := model.DB.Create(&comment)
	if r1.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Comment insert failed!",
			"desc":    r1.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment insert sucessfully!",
	})
}

func MigrateDB(c *gin.Context) {
	model.DB.Model(&model.Comment{}).AddForeignKey("user_id", "Users(id)", "CASCADE", "CASCADE")
	model.DB.Model(&model.Comment{}).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")
	c.JSON(http.StatusOK, gin.H{
		"message": "the Database Migration  sucessfully!",
	})
}

func GetCommentsByAid(c *gin.Context) {
	var a1 model.Article

	if err := model.DB.Where("id = ?", c.Param("aid")).First(&a1).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article Record not found!"})
		return
	}
	model.DB.Model(&a1).Related(&a1.Comments)
	c.JSON(http.StatusOK, gin.H{"data": a1})

}

// DELETE /comments/:id
// Delete a  comment
func DeleteComment(c *gin.Context) {
	// Get model if exist
	var c1 model.Comment
	if err := model.DB.Where("id = ?", c.Param("id")).First(&c1).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	model.DB.Delete(&c1)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
