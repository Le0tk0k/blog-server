package web

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"

	"github.com/Le0tk0k/blog-server/log"
	"github.com/Le0tk0k/blog-server/service"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	userService service.UserService
}

// NewAuthMiddleware はAuthMiddlewareを返す
func NewAuthMiddleware(userService service.UserService) AuthMiddleware {
	return AuthMiddleware{
		userService: userService,
	}
}

func (m *AuthMiddleware) Login(c echo.Context) error {
	logger := log.New()
	req := new(userJSON)
	if err := c.Bind(req); err != nil {
		logger.Errorj(map[string]interface{}{"message": "failed to bind", "error": err.Error()})
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	email := req.Email
	password := req.Password

	user, err := m.userService.GetUser(email)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if user.Email != email || err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	t, err := token.SignedString([]byte(os.Getenv("AUTH_KEY")))
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	c.SetCookie(&http.Cookie{
		Name:   "jwt",
		Value:  t,
		Path:   "/",
		MaxAge: 60 * 60,
		// Todo デプロイしてからどうにかする
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	return c.JSON(http.StatusOK, map[string]string{
		"message": "successfully logged in",
	})
}

type userJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
