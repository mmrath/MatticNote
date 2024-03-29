package account

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/MatticNote/MatticNote/server/account/oauth"
	"github.com/MatticNote/MatticNote/server/account/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/session"
	"strings"
	"time"
)

var login2faSession *session.Store

func csrfErrorView(c *fiber.Ctx, _ error) error {
	return c.Status(fiber.StatusForbidden).Render(
		"csrf_failed",
		fiber.Map{},
	)
}

func ConfigureRoute(r fiber.Router) {
	r.Use(csrf.New(csrf.Config{
		Next: func(c *fiber.Ctx) bool {
			if strings.HasPrefix(c.Path(), "/account/oauth/") {
				return c.Path() != "/account/oauth/authorize"
			} else {
				return false
			}
		},
		KeyLookup:      fmt.Sprintf("form:%s", misc.CSRFFormName),
		CookieName:     "_csrf",
		CookiePath:     "/account",
		CookieSecure:   config.Config.Server.CookieSecure,
		CookieHTTPOnly: true,
		ContextKey:     misc.CSRFContextKey,
		ErrorHandler:   csrfErrorView,
		KeyGenerator: func() string {
			return misc.GenToken(32)
		},
		Storage: config.GetFiberRedisMemory(),
	}))
	login2faSession = session.New(session.Config{
		Expiration:     5 * time.Minute,
		Storage:        config.GetFiberRedisMemory(),
		CookiePath:     "/account",
		CookieSecure:   config.Config.Server.CookieSecure,
		CookieHTTPOnly: true,
		KeyGenerator: func() string {
			return misc.GenToken(32)
		},
	})

	r.Get("/register", registerUserGet)
	r.Post("/register",
		limiter.New(limiter.Config{
			Next: func(_ *fiber.Ctx) bool {
				return config.Config.Server.DisableAccountRateLimit
			},
			Max: int(config.Config.Server.AccountRegistrationLimitCount),
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_ACCTREG-%s", c.IP())
			},
			Expiration: 24 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		registerPost,
	)

	r.Get("/login", loginUserGet)
	r.Post("/login",
		limiter.New(limiter.Config{
			Next: func(_ *fiber.Ctx) bool {
				return config.Config.Server.DisableAccountRateLimit
			},
			Max: 30,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_LOGIN-%s", c.IP())
			},
			Expiration: 30 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		loginPost,
	)
	r.Get("/login/2fa", login2faGet)
	r.Post("/login/2fa", login2faPost)

	r.Get("/logout", destroySession)

	r.Get("/forgot", forgotPasswordGet)
	r.Post("/forgot",
		limiter.New(limiter.Config{
			Next: func(_ *fiber.Ctx) bool {
				return config.Config.Server.DisableAccountRateLimit
			},
			Max: 10,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_FORGOT-%s", c.IP())
			},
			Expiration: 1 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		forgotPasswordPost,
	)
	r.Get("/forgot/:token",
		limiter.New(limiter.Config{
			Next: func(_ *fiber.Ctx) bool {
				return config.Config.Server.DisableAccountRateLimit
			},
			Max: 30,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_FPTF-%s", c.IP())
			},
			Expiration: 3 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				c.Status(fiber.StatusTooManyRequests)
				return nil
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		forgotPasswordResetGet,
	)
	r.Post("/forgot/:token",
		limiter.New(limiter.Config{
			Next: func(_ *fiber.Ctx) bool {
				return config.Config.Server.DisableAccountRateLimit
			},
			Max: 30,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_FPTP-%s", c.IP())
			},
			Expiration: 3 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				c.Status(fiber.StatusTooManyRequests)
				return nil
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		forgotPasswordResetPost,
	)

	r.Get("/issue_confirm_mail", issueConfirmGet)
	r.Post("/issue_confirm_mail",
		limiter.New(limiter.Config{
			Next: func(_ *fiber.Ctx) bool {
				return config.Config.Server.DisableAccountRateLimit
			},
			Max: 10,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_ISSCM-%s", c.IP())
			},
			Expiration: 1 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		issueConfirmPost,
	)

	r.Get("/verify/:token",
		limiter.New(limiter.Config{
			Next: func(_ *fiber.Ctx) bool {
				return config.Config.Server.DisableAccountRateLimit
			},
			Max: 30,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_VFTK-%s", c.IP())
			},
			Expiration: 1 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				c.Status(fiber.StatusTooManyRequests)
				return nil
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		verifyMail,
	)

	oauth.ConfigureRoute(r.Group("/oauth"))
	settings.ConfigureRoute(r.Group("/settings"))
}
