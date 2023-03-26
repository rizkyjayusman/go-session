package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	sessionName       = "session_token"
	expiryInSecond    = 120
	isSessionSecure   = false
	isSessionHttpOnly = false

	baseUrl  = "localhost"
	basePath = "/"
	port     = 8000
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
		return
	}

	expectedPassword, ok := users[cred.Username]
	if !ok || expectedPassword != cred.Password {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.NewString()
	expiredAt := time.Now().Add(expiryInSecond * time.Second)

	sessions[sessionToken] = Session{
		username: cred.Username,
		expiry:   expiredAt,
	}

	c.SetCookie(sessionName, sessionToken, expiryInSecond, basePath, baseUrl, isSessionSecure, isSessionHttpOnly)
}

func Welcome(c *gin.Context) {
	sessionToken, err := c.Cookie(sessionName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userSession, exists := sessions[sessionToken]
	fmt.Println(sessions)
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Writer.Write([]byte(fmt.Sprintf("Welcome, %s!", userSession.username)))
}

func Refresh(c *gin.Context) {
	sessionToken, err := c.Cookie(sessionName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userSession, exists := sessions[sessionToken]
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if userSession.isExpired() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	newSessionToken := uuid.NewString()
	expiredAt := time.Now().Add(expiryInSecond * time.Second)

	sessions[newSessionToken] = Session{
		username: userSession.username,
		expiry:   expiredAt,
	}

	delete(sessions, sessionToken)
	c.SetCookie(sessionName, newSessionToken, expiryInSecond, basePath, baseUrl, isSessionSecure, isSessionHttpOnly)
}

func Logout(c *gin.Context) {
	sessionToken, err := c.Cookie(sessionName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	delete(sessions, sessionToken)
	c.SetCookie(sessionName, "", 0, basePath, baseUrl, isSessionSecure, isSessionHttpOnly)
}

func main() {
	router := gin.Default()
	router.POST("/signin", Signin)
	router.GET("/welcome", Welcome)
	router.POST("/refresh", Refresh)
	router.POST("/logout", Logout)
	router.Run(fmt.Sprintf("%s:%v", baseUrl, port))
}
