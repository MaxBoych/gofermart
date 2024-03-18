package cookie

import (
	"github.com/MaxBoych/gofermart/internal/user/user_models"
	"github.com/MaxBoych/gofermart/pkg/errs"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func SetCookie(c *fiber.Ctx, params user_models.CookieData) {
	cookie := &fiber.Cookie{
		Name:     params.Name,
		Value:    params.Value,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		Domain:   params.Domain,
		Expires:  params.Expires,
	}

	if strings.HasPrefix(c.Hostname(), "localhost") {
		cookie.Domain = "localhost"
	}

	c.Cookie(cookie)
}

func GetCookie(c *fiber.Ctx, name string) (string, error) {
	cookie := c.Cookies(name)
	if cookie == "" {
		return "", errs.HttpErrCookieIsEmpty
	}
	return cookie, nil
}
