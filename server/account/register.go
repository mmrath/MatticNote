package account

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
)

type registerUserStruct struct {
	Username          string `validate:"required,username"`
	Email             string `validate:"required,email"`
	Password          string `validate:"required,min=8"`
	ReCaptchaResponse string `form:"g-recaptcha-response"`
	// TODO: CAPTCHAなどの対策用のフォーム内容も含める
}

func registerUserGet(c *fiber.Ctx) error {
	if c.Cookies(signature.JWTAuthCookieName, "") != "" {
		return c.Redirect(c.Query("next", "/web/"))
	}

	return registerUserView(c)
}

func registerUserView(c *fiber.Ctx, errors ...string) error {
	return c.Status(fiber.StatusOK).Render(
		"account/register",
		fiber.Map{
			"CSRFFormName": misc.CSRFFormName,
			"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
			"errors":       errors,
			"reCaptchaKey": config.Config.Server.RecaptchaSiteKey,
		},
		"_layout/account",
	)
}

func registerPost(c *fiber.Ctx) error {
	formData := new(registerUserStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return registerUserView(c, errs...)
	}

	_, err := user.RegisterLocalUser(formData.Email, formData.Username, formData.Password, false)
	if err != nil {
		if err == user.ErrUserExists {
			return registerUserView(c, "Username or email is already taken")
		} else {
			return err
		}
	}

	return c.Redirect("/account/login?created=true")
}
