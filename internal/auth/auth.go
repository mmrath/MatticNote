package auth

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal/oauth"
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ory/fosite"
	"strings"
)

const (
	LoginUserLocal       = "loginUser"
	AuthorizeMethodLocal = "authorizeMethod"
	OAuthClientIDLocal   = "oauthClientID"
	OAuthTokenLocal      = "oauthToken"
)

type AuthorizeMethod string

const (
	JWT   AuthorizeMethod = "jwt"
	OAuth AuthorizeMethod = "oauth"
)

func AuthenticationUser(c *fiber.Ctx) error {
	// JWT token authentication
	headerSplit := strings.Split(c.Get(signature.AuthHeaderName, ""), " ")
	if len(headerSplit) > 0 && strings.TrimSpace(headerSplit[0]) == signature.AuthSchemeName {
		token, ok := c.Locals(signature.JWTContextKey).(*jwt.Token)
		if ok {
			claim := token.Claims.(jwt.MapClaims)
			usr, err := user.GetLocalUser(uuid.MustParse(claim["sub"].(string)))
			if err == nil {
				c.Locals(LoginUserLocal, usr)
				c.Locals(AuthorizeMethodLocal, JWT)
			}
		} else {
			return fiber.ErrUnauthorized
		}
	} else if len(headerSplit) == 2 && strings.TrimSpace(headerSplit[0]) == fosite.BearerAccessToken {
		token := strings.TrimSpace(headerSplit[1])
		authorizedUser, clientId, err := oauth.APIIntrospect(token)
		if err != nil {
			switch {
			case errors.Is(err, fosite.ErrTokenExpired),
				errors.Is(err, fosite.ErrLoginRequired),
				errors.Is(err, fosite.ErrInactiveToken):
				c.Status(fiber.StatusUnauthorized)
				return nil
			default:
				return err
			}
		}
		c.Locals(LoginUserLocal, authorizedUser)
		c.Locals(AuthorizeMethodLocal, OAuth)
		c.Locals(OAuthClientIDLocal, clientId)
		c.Locals(OAuthTokenLocal, token)
	}

	return c.Next()
}
