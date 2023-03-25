package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var sessions = map[string]Session{}

type Session struct {
	username string
	expiry   time.Time
}

func (s Session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func Signin(c *gin.Context) {
	var cred Credentials
	err := c.ShouldBindJSON(&cred)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	expectedPassword, ok := users[cred.Username]
	if !ok || expectedPassword != cred.Password {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	sessionToken := uuid.NewString()
	expiredAt := time.Now().Add(120 * time.Second)

	sessions[sessionToken] = Session{
		username: cred.Username,
		expiry:   expiredAt,
	}

	c.SetCookie("session_token", sessionToken, 120, "/", "localhost", false, false)
}

func Welcome(c *gin.Context) {

}

func Refresh(c *gin.Context) {

}

func Logout(c *gin.Context) {

}

func main() {
	router := gin.Default()
	router.POST("/signin", Signin)
	router.GET("/welcome", Welcome)
	router.POST("/refresh", Refresh)
	router.POST("/logout", Logout)
	router.Run("localhost:8001")
}
