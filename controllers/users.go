package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samansahebi/dekamond/models"
)

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	models.DB.Where("id = ?", id).First(&user)
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func Users(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	var total int64
	var users []models.User
	models.DB.Model(&models.User{}).Count(&total)

	offset := (page - 1) * limit
	models.DB.Limit(limit).Offset(offset).Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
		"data":       users,
	})
}

func Search(c *gin.Context) {
	phone := c.Query("PhoneNumber") // e.g. /users/search?phone=1234567890

	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone query parameter is required"})
		return
	}

	var user models.User
	if err := models.DB.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
