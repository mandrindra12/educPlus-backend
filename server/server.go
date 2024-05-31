package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mada.h/educplus/database"
	"mada.h/educplus/mailsender"
	"mada.h/educplus/models"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func ListenAndServe() {
	gin.SetMode(gin.ReleaseMode)
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

	// give the client a new access token if the previous has expired
	router.GET("/api/newToken/:refresh", func(ctx *gin.Context) {
		refreshToken := ctx.Param("refresh")
		if newToken, err := GenerateNewAccessToken(refreshToken); err != nil {
			ctx.AbortWithStatus(400)
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{"new_access_token": newToken})
			return
		}
	})

	// handle login
	router.POST("/api/login", func(c *gin.Context) {
		var admin models.Credentials
		if err := c.ShouldBindJSON(&admin); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		} else {
			if database.Authenticate(admin.Username, admin.Password) {
				accessToken, refreshToken, err := GenerateToken(admin.Username)
				if err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
				return
			} else {
				c.JSON(http.StatusNotFound, gin.H{"message": "username or password incorrect"})
				return
			}
		}
	})
	// handle email registration
	router.POST("/api/subscription", func(c *gin.Context) {
		mail := models.Mail{}
		if err := c.ShouldBindJSON(&mail); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else {
			database.RegisterMail(mail)
			c.JSON(http.StatusCreated, gin.H{"hello": "subscription successful"})
			return
		}
	})
	// handle event updates and send mails to customers in case of brand new event
	router.POST("/api/newEvent", func(c *gin.Context) {
		event := models.Event{}
		// if err := c.ShouldBindJSON(&event); err != nil {
		// c.JSON(http.StatusInternalServerError, nil)
		// return
		// } else {
		emails := database.GetAllMail()
		event.Title = c.PostForm("title")
		event.Description = c.PostForm("description")
		// file, _ := c.FormFile("file")
		// database.RegisterEvent(event)
		// emails := database.GetAllMail()
		if err := mailsender.SendMail(event, emails); err != nil {
			c.JSON(http.StatusRequestTimeout, gin.H{"message": "mailsender failed"})
		}
		c.JSON(http.StatusCreated, nil)
		return
		// }
	})
	// route for testing cookie
	router.GET("/api/login", func(c *gin.Context) {
		if token, err := ExtractTokenFromRequest(c); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else {
			if claims, err := ValidateJWTToken(token); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			} else {
				c.JSON(http.StatusAccepted, gin.H{"username": claims["Username"]})
			}
		}

		return
	})
	// route to retrieve all events
	router.GET("/api/events", func(c *gin.Context) {
		events := database.GetAllEvents()
		println(len(events))
		return
	})
	// route for testing the route protection
	router.GET("/protected", func(c *gin.Context) {
		if isAdmin, err := IsAdmin(c); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"": "nothing to do here"})
			return
		} else if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied by the server"})
			return
		} else {
			c.JSON(http.StatusAccepted, gin.H{"sudo": "here you go!"})
			return
		}
	})
	// handle 404
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	})
	// RUN RUN RUN ...
	router.Run(":8080")
}

func GenerateToken(username string) (string, string, error) {
	expirationTime := time.Now().Add(time.Minute * 5)
	refreshTime := time.Now().Add(time.Hour * 2)
	loginClaims := &models.Claims{
		Username: username,
		Type:     "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshClaims := &models.Claims{
		Username: username,
		Type:     "refresh",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTime.Unix(),
		},
	}
	loginToken := jwt.NewWithClaims(jwt.SigningMethodHS256, loginClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	loginTokenString, err := loginToken.SignedString(jwtKey)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	} else {
		return loginTokenString, refreshTokenString, nil
	}
}

func ExtractTokenFromRequest(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
		return "", errors.New("invalid authorization header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func IsAdmin(c *gin.Context) (bool, error) {
	token, err := ExtractTokenFromRequest(c)
	if err != nil {
		return false, err
	}
	claims, err := ValidateJWTToken(token)
	if claims["Type"] == "access" {
		return true, nil
	}
	return false, nil

}

func ValidateJWTToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch ve.Errors {
			case jwt.ValidationErrorMalformed:
				return nil, errors.New("malformed token")
			case jwt.ValidationErrorUnverifiable:
				return nil, errors.New("token could not be verified")
			case jwt.ValidationErrorSignatureInvalid:
				return nil, errors.New("invalid token signature")
			case jwt.ValidationErrorExpired:
				return nil, errors.New("expired token")
			default:
				return nil, errors.New("invalid token")
			}
		} else {
			return nil, err
		}
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GenerateNewAccessToken(refreshToken string) (string, error) {
	if claims, err := ValidateJWTToken(refreshToken); err != nil {
		return "", err
	} else {
		username := fmt.Sprintf("%s", claims["Username"])
		newAccessToken, _, err := GenerateToken(username)
		if err != nil {
			return "", err
		} else {
			return newAccessToken, nil
		}
	}
}
