package main

import (
	"fmt"
	"log"
	"net/http"
	"rizkyjayusman/go-session/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	expiredAt := time.Now().Add(time.Duration(config.ExpiryInSecond) * time.Second)

	sessions[sessionToken] = Session{
		username: cred.Username,
		expiry:   expiredAt,
	}

	c.SetCookie(config.SessionName, sessionToken, config.ExpiryInSecond, config.BasePath, config.BaseUrl, config.IsSessionSecure, config.IsSessionHttpOnly)
}

func Welcome(c *gin.Context) {
	sessionToken, err := c.Cookie(config.SessionName)
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
	sessionToken, err := c.Cookie(config.SessionName)
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
	expiredAt := time.Now().Add(time.Duration(config.ExpiryInSecond) * time.Second)

	sessions[newSessionToken] = Session{
		username: userSession.username,
		expiry:   expiredAt,
	}

	delete(sessions, sessionToken)
	c.SetCookie(config.SessionName, newSessionToken, config.ExpiryInSecond, config.BasePath, config.BaseUrl, config.IsSessionSecure, config.IsSessionHttpOnly)
}

func Logout(c *gin.Context) {
	sessionToken, err := c.Cookie(config.SessionName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	delete(sessions, sessionToken)
	c.SetCookie(config.SessionName, "", 0, config.BasePath, config.BaseUrl, config.IsSessionSecure, config.IsSessionHttpOnly)
}

var config util.Config

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("the db fail to run!")
	}
	fmt.Println(db)

	router := gin.Default()
	router.POST("/signin", Signin)
	router.GET("/welcome", Welcome)
	router.POST("/refresh", Refresh)
	router.POST("/logout", Logout)
	router.Run(fmt.Sprintf("%s:%v", config.BaseUrl, config.PORT))
}
