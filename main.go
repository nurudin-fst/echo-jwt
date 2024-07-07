package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var jwtSecret = []byte("oiye897erowlkehrh8ieyr") // Ganti dengan secret yang lebih aman

// JWTClaims mendefinisikan payload untuk token JWT
type JWTClaims struct {
	JwtIss   string `json:"iss"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	e := echo.New()

	// Middleware Logger dan Recover
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Endpoint Public
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Echo with JWT!")
	})

	// Endpoint untuk login
	e.POST("/login", login)

	// Endpoint yang memerlukan autentikasi
	r := e.Group("/restricted")
	var IsLoggedIn = echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecret),
	})

	r.GET("", restricted, IsLoggedIn)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler Login
func login(c echo.Context) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.ErrBadRequest
	}

	// Validasi username dan password sederhana (sebaiknya gunakan DB di produksi)
	if req.Username == "admin" && req.Password == "password" {
		// Buat token JWT
		claims := &JWTClaims{
			JwtIss:   "gcAzjNDX9swAIX6oltBpCFcIoKP0sF9N",
			Username: req.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // Token valid 1 jam
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(jwtSecret)
		if err != nil {
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": signedToken,
		})
	}

	return echo.ErrUnauthorized
}

// Handler untuk endpoint yang membutuhkan autentikasi
func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return c.JSON(http.StatusOK, claims)
}
