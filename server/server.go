package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mada.h/educplus/database"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

// struct that would map the request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string
	jwt.StandardClaims
}

func ListenAndServe() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// handle login
	router.POST("/api/login", func(c *gin.Context) {
		var admin Credentials
		if err := c.ShouldBindJSON(&admin); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		} else {
			if database.Authenticate(admin.Username, admin.Password) {
				ManageCookies(c, admin)
				c.JSON(http.StatusOK, gin.H{"message": "loggin successful"})
				return
			} else {
				c.JSON(http.StatusNotFound, gin.H{"message": "username or password incorrect"})
				return
			}
		}
		// c.JSON(200, gin.H{"user": "use"})
		// return
	})
	router.GET("/api/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": "get request"})
		return
	})
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	})
	router.Run(":8080")
}

func ManageCookies(c *gin.Context, admin Credentials) {
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
		Username: admin.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.SetCookie("token", tokenString, int(300), "/", "", false, true)
	return
}
