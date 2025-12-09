package validator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name    string   `json:"user_name" binding:"required"`
	Email   string   `form:"email_addr" binding:"required,email"`
	Token   string   `header:"X-Token" binding:"required"`
	Bio     string   `label:"个人简介" binding:"required"`
	Phone   string   `binding:"required"`
	Profile *Profile `json:"profile" binding:"required"`
}

type Profile struct {
	NickName string `json:"nick_name" binding:"required"`
	Avatar   string `json:"avatar" binding:"required,url"`
}

func CreateUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		// key point: pass the struct to get the friendly field name
		errorMsg := FmtValidationErrors(err, user)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
