package ap

import (
	"encoding/json"
	"fmt"
	"github.com/MatticNote/MatticNote/activitypub"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func apUserController(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}
	if misc.IsAPAcceptHeader(c) {
		return RenderUser(c, targetUuid)
	} else {
		targetUser, err := internal.GetLocalUser(targetUuid)
		if err != nil {
			switch err {
			case internal.ErrNoSuchUser:
				return fiber.ErrNotFound
			case internal.ErrUserGone:
				return fiber.ErrGone
			default:
				return err
			}
		}
		return c.Redirect(fmt.Sprintf("/@%s", targetUser.Username))
	}
}

func RenderUser(c *fiber.Ctx, targetUuid uuid.UUID) error {
	c.Set("Content-Type", "application/activity+json; charset=utf-8")

	actor, err := activitypub.RenderActor(targetUuid)
	if err != nil {
		return err
	}

	body, err := json.Marshal(actor)
	if err != nil {
		return err
	}

	return c.Send(body)
}
