package main

import (
	"fmt"
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
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	userSession, exists := sessions[sessionToken]
	fmt.Println(sessions)
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Writer.Write([]byte(fmt.Sprintf("Welcome, %s!", userSession.username)))
}

func Refresh(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	userSession, exists := sessions[sessionToken]
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if userSession.isExpired() {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	newSessionToken := uuid.NewString()
	expiredAt := time.Now().Add(120 * time.Second)

	sessions[newSessionToken] = Session{
		username: userSession.username,
		expiry:   expiredAt,
	}

	delete(sessions, sessionToken)
	c.SetCookie("session_token", newSessionToken, 120, "/", "localhost", false, false)
}

func Logout(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	delete(sessions, sessionToken)
	c.SetCookie("session_token", "", 0, "/", "localhost", false, false)
}

func main() {
	router := gin.Default()
	router.POST("/signin", Signin)
	router.GET("/welcome", Welcome)
	router.POST("/refresh", Refresh)
	router.POST("/logout", Logout)
	router.Run("localhost:8000")
}
