package main

import (
	"time"

	"github.com/gin-gonic/gin"
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
	router.Run("localhost:8080")
}
