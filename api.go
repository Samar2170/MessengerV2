package main

import (
	"encoding/json"
	"net/http"
	"time"

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

func createService(c echo.Context) error {
	serviceName := c.FormValue("service")
	var claims jwtCustomClaims
	user := c.Get("user").(*jwt.Token)
	tmp, _ := json.Marshal(user.Claims)
	_ = json.Unmarshal(tmp, &claims)
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

func StartEcho() {
	e := echo.New()
	e.POST("/signup", signup)
	e.POST("/login", login)
	e.POST("/create-service", createService)
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
