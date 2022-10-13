package main

import (
	"encoding/json"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func signup(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")
	if username == "" || password == "" || email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "All fields are required",
		})
	}

	user := User{
		Username: username,
		Password: password,
		Email:    email,
	}
	err2 := user.Create()
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "User already exists",
		})
	}
	return c.JSON(http.StatusOK, user)
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	user, err := FindUser(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if user.Password != password {
		return c.String(http.StatusInternalServerError, "Invalid Password")
	}
	claims := &jwtCustomClaims{
		user.ID,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})

}
func unwrapToken(token *jwt.Token) (*jwtCustomClaims, error) {
	var claims jwtCustomClaims
	tmp, _ := json.Marshal(token.Claims)
	_ = json.Unmarshal(tmp, &claims)
	return &claims, nil
}

func createService(c echo.Context) error {
	serviceName := c.FormValue("service")
	user := c.Get("user").(*jwt.Token)
	claims, err := unwrapToken(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// var claims jwtCustomClaims
	// tmp, _ := json.Marshal(user.Claims)
	// _ = json.Unmarshal(tmp, &claims)
	if serviceName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "All fields are required",
		})
	}
	service := Service{Name: serviceName, UserID: claims.Id}
	Db.Create(&service)
	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "Service Created",
	})
}

func botMessage(e echo.Context) error {
	serviceName := e.FormValue("service")
	message := e.FormValue("message")
	user := e.Get("user").(*jwt.Token)
	claims, err := unwrapToken(user)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	service, err := GetServiceByUser(serviceName, claims.Id)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	subscibers, err := GetSubscriberByServiceId(service.ID)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	for _, subscriber := range subscibers {
		msg := tgbotapi.NewMessage(subscriber.ChatID, message)
		Tgbot.Send(msg)
	}
	return e.JSON(http.StatusOK, map[string]string{
		"message": "Message sent",
	})
}

func StartEcho() {
	e := echo.New()
	e.POST("/signup", signup)
	e.POST("/login", login)
	e.POST("/create-service", createService)
	e.POST("/bot-message", botMessage)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		TokenLookup: "header:Authorization",
		ContextKey:  "user",
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/signup" || c.Path() == "/login" {
				return true
			}
			return false
		},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
