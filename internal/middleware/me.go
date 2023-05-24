package middleware

import (
	"fmt"
	"mangosteen/config/queries"
	"mangosteen/internal/database"
	"mangosteen/internal/jwt_helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getMe(c)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Set("me", user)
		c.Next()
	}
}

func getMe(c *gin.Context) (queries.User, error) {
	var user queries.User
	auth := c.GetHeader("Authorization")
	if len(auth) < 8 {
		return user, fmt.Errorf("JWT为空")
	}
	jwtString := auth[7:]
	if len(jwtString) == 0 {
		return user, fmt.Errorf("JWT为空")
	}
	t, err := jwt_helper.Parse(jwtString)
	if err != nil {
		return user, fmt.Errorf("无效的JWT")
	}
	m, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return user, fmt.Errorf("无效的JWT")
	}
	userID, ok := m["user_id"].(float64)
	if !ok {
		return user, fmt.Errorf("无效的JWT")
	}
	userIDInt := int32(userID)
	if err != nil {
		return user, fmt.Errorf("无效的JWT")
	}
	q := database.NewQuery()
	user, err = q.FindUser(c, userIDInt)
	if err != nil {
		return user, fmt.Errorf("无效的JWT")
	}
	return user, nil
}
