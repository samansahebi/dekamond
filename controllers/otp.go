package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samansahebi/dekamond/models"
)

var jwtKey = []byte("secret_key")

func GenerateJWT(phone string) (string, error) {

	claims := &jwt.MapClaims{
		"phone": phone,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization") // read "Authorization" header
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("phone", claims["phone"])
		} else {
			c.JSON(401, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func SendOTP(c *gin.Context) {
	var user models.User
	var count int64

	rand.Seed(time.Now().UnixNano())
	otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&models.OTP{}).Where("created_at >= ? AND phone_number = ?", time.Now().Add(-10*time.Minute), user.PhoneNumber).Count(&count)
	if count >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many requests"})
		return
	}
	if err := models.DB.Create(&models.OTP{PhoneNumber: user.PhoneNumber, OTPCode: otpCode}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(otpCode)
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP Code has been sent",
	})
}

func VerifyOTP(c *gin.Context) {
	var otp models.OTP
	var user models.User
	var req struct {
		PhoneNumber string
		OTPCode     string
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Where("phone_number = ?", req.PhoneNumber).Order("created_at desc").First(&otp)

	if time.Since(otp.CreatedAt) > 2*time.Minute {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP Code has expired"})
		return
	}

	result := models.DB.Where("phone_number = ?", otp.PhoneNumber).First(&user)

	if result.Error != nil {
		user = models.User{PhoneNumber: otp.PhoneNumber}
		models.DB.Create(&user)
	}

	token, err := GenerateJWT(req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
